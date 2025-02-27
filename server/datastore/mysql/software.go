package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/jmoiron/sqlx"
)

const (
	maxSoftwareNameLen             = 255
	maxSoftwareVersionLen          = 255
	maxSoftwareSourceLen           = 64
	maxSoftwareBundleIdentifierLen = 255

	maxSoftwareReleaseLen = 64
	maxSoftwareVendorLen  = 32
	maxSoftwareArchLen    = 16
)

func truncateString(str string, length int) string {
	if len(str) > length {
		return str[:length]
	}
	return str
}

func softwareToUniqueString(s fleet.Software) string {
	ss := []string{s.Name, s.Version, s.Source, s.BundleIdentifier}
	// Release, Vendor and Arch fields were added on a migration,
	// thus we only include them in the string if at least one of them is defined.
	if s.Release != "" || s.Vendor != "" || s.Arch != "" {
		ss = append(ss, s.Release, s.Vendor, s.Arch)
	}
	return strings.Join(ss, "\u0000")
}

func uniqueStringToSoftware(s string) fleet.Software {
	parts := strings.Split(s, "\u0000")

	// Release, Vendor and Arch fields were added on a migration,
	// If one of them is defined, then they are included in the string.
	var release, vendor, arch string
	if len(parts) > 4 {
		release = truncateString(parts[4], maxSoftwareReleaseLen)
		vendor = truncateString(parts[5], maxSoftwareVendorLen)
		arch = truncateString(parts[6], maxSoftwareArchLen)
	}

	return fleet.Software{
		Name:             truncateString(parts[0], maxSoftwareNameLen),
		Version:          truncateString(parts[1], maxSoftwareVersionLen),
		Source:           truncateString(parts[2], maxSoftwareSourceLen),
		BundleIdentifier: truncateString(parts[3], maxSoftwareBundleIdentifierLen),

		Release: release,
		Vendor:  vendor,
		Arch:    arch,
	}
}

func softwareSliceToMap(softwares []fleet.Software) map[string]fleet.Software {
	result := make(map[string]fleet.Software)
	for _, s := range softwares {
		result[softwareToUniqueString(s)] = s
	}
	return result
}

// UpdateHostSoftware updates the software list of a host.
// The update consists of deleting existing entries that are not in the given `software`
// slice, updating existing entries and inserting new entries.
func (ds *Datastore) UpdateHostSoftware(ctx context.Context, hostID uint, software []fleet.Software) error {
	return ds.withRetryTxx(ctx, func(tx sqlx.ExtContext) error {
		return applyChangesForNewSoftwareDB(ctx, tx, hostID, software, ds.minLastOpenedAtDiff)
	})
}

func nothingChanged(current, incoming []fleet.Software, minLastOpenedAtDiff time.Duration) bool {
	if len(current) != len(incoming) {
		return false
	}

	currentMap := make(map[string]fleet.Software)
	for _, s := range current {
		currentMap[softwareToUniqueString(s)] = s
	}
	for _, s := range incoming {
		cur, ok := currentMap[softwareToUniqueString(s)]
		if !ok {
			return false
		}

		// if the incoming software has a last opened at timestamp and it differs
		// significantly from the current timestamp (or there is no current
		// timestamp), then consider that something changed.
		if s.LastOpenedAt != nil {
			if cur.LastOpenedAt == nil {
				return false
			}

			oldLast := *cur.LastOpenedAt
			newLast := *s.LastOpenedAt
			if newLast.Sub(oldLast) >= minLastOpenedAtDiff {
				return false
			}
		}
	}

	return true
}

func (ds *Datastore) ListSoftwareByHostIDShort(ctx context.Context, hostID uint) ([]fleet.Software, error) {
	return listSoftwareByHostIDShort(ctx, ds.reader, hostID)
}

func listSoftwareByHostIDShort(
	ctx context.Context,
	db sqlx.QueryerContext,
	hostID uint,
) ([]fleet.Software, error) {
	q := `
SELECT
    s.id,
    s.name,
    s.version,
    s.source,
    s.bundle_identifier,
    s.release,
    s.vendor,
    s.arch,
    hs.last_opened_at
FROM
    software s
    JOIN host_software hs ON hs.software_id = s.id
WHERE
    hs.host_id = ?
`
	var softwares []fleet.Software
	err := sqlx.SelectContext(ctx, db, &softwares, q, hostID)
	if err != nil {
		return nil, err
	}

	return softwares, nil
}

func applyChangesForNewSoftwareDB(
	ctx context.Context,
	tx sqlx.ExtContext,
	hostID uint,
	software []fleet.Software,
	minLastOpenedAtDiff time.Duration,
) error {
	currentSoftware, err := listSoftwareByHostIDShort(ctx, tx, hostID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "loading current software for host")
	}

	if nothingChanged(currentSoftware, software, minLastOpenedAtDiff) {
		return nil
	}

	current := softwareSliceToMap(currentSoftware)
	incoming := softwareSliceToMap(software)

	if err = deleteUninstalledHostSoftwareDB(ctx, tx, hostID, current, incoming); err != nil {
		return err
	}

	if err = insertNewInstalledHostSoftwareDB(ctx, tx, hostID, current, incoming); err != nil {
		return err
	}

	if err = updateModifiedHostSoftwareDB(ctx, tx, hostID, current, incoming, minLastOpenedAtDiff); err != nil {
		return err
	}

	return nil
}

// delete host_software that is in current map, but not in incoming map.
func deleteUninstalledHostSoftwareDB(
	ctx context.Context,
	tx sqlx.ExecerContext,
	hostID uint,
	currentMap map[string]fleet.Software,
	incomingMap map[string]fleet.Software,
) error {
	var deletesHostSoftware []interface{}
	deletesHostSoftware = append(deletesHostSoftware, hostID)

	for currentKey, curSw := range currentMap {
		if _, ok := incomingMap[currentKey]; !ok {
			deletesHostSoftware = append(deletesHostSoftware, curSw.ID)
		}
	}
	if len(deletesHostSoftware) <= 1 {
		return nil
	}
	sql := fmt.Sprintf(
		`DELETE FROM host_software WHERE host_id = ? AND software_id IN (%s)`,
		strings.TrimSuffix(strings.Repeat("?,", len(deletesHostSoftware)-1), ","),
	)
	if _, err := tx.ExecContext(ctx, sql, deletesHostSoftware...); err != nil {
		return ctxerr.Wrap(ctx, err, "delete host software")
	}

	return nil
}

func getOrGenerateSoftwareIdDB(ctx context.Context, tx sqlx.ExtContext, s fleet.Software) (uint, error) {
	getExistingID := func() (int64, error) {
		var existingID int64
		if err := sqlx.GetContext(ctx, tx, &existingID,
			"SELECT id FROM software "+
				"WHERE name = ? AND version = ? AND source = ? AND `release` = ? AND "+
				"vendor = ? AND arch = ? AND bundle_identifier = ? LIMIT 1",
			s.Name, s.Version, s.Source, s.Release, s.Vendor, s.Arch, s.BundleIdentifier,
		); err != nil {
			return 0, err
		}
		return existingID, nil
	}

	switch id, err := getExistingID(); {
	case err == nil:
		return uint(id), nil
	case errors.Is(err, sql.ErrNoRows):
		// OK
	default:
		return 0, ctxerr.Wrap(ctx, err, "get software")
	}

	_, err := tx.ExecContext(ctx,
		"INSERT INTO software "+
			"(name, version, source, `release`, vendor, arch, bundle_identifier) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?) "+
			"ON DUPLICATE KEY UPDATE bundle_identifier=VALUES(bundle_identifier)",
		s.Name, s.Version, s.Source, s.Release, s.Vendor, s.Arch, s.BundleIdentifier,
	)
	if err != nil {
		return 0, ctxerr.Wrap(ctx, err, "insert software")
	}

	// LastInsertId sometimes returns 0 as it's dependent on connections and how mysql is
	// configured.
	switch id, err := getExistingID(); {
	case err == nil:
		return uint(id), nil
	case errors.Is(err, sql.ErrNoRows):
		return 0, doRetryErr
	default:
		return 0, ctxerr.Wrap(ctx, err, "get software")
	}
}

// insert host_software that is in incoming map, but not in current map.
func insertNewInstalledHostSoftwareDB(
	ctx context.Context,
	tx sqlx.ExtContext,
	hostID uint,
	currentMap map[string]fleet.Software,
	incomingMap map[string]fleet.Software,
) error {
	var insertsHostSoftware []interface{}

	incomingOrdered := make([]string, 0, len(incomingMap))
	for s := range incomingMap {
		incomingOrdered = append(incomingOrdered, s)
	}
	sort.Strings(incomingOrdered)

	for _, s := range incomingOrdered {
		if _, ok := currentMap[s]; !ok {
			id, err := getOrGenerateSoftwareIdDB(ctx, tx, uniqueStringToSoftware(s))
			if err != nil {
				return err
			}
			sw := incomingMap[s]
			insertsHostSoftware = append(insertsHostSoftware, hostID, id, sw.LastOpenedAt)
		}
	}

	if len(insertsHostSoftware) > 0 {
		values := strings.TrimSuffix(strings.Repeat("(?,?,?),", len(insertsHostSoftware)/3), ",")
		sql := fmt.Sprintf(`INSERT IGNORE INTO host_software (host_id, software_id, last_opened_at) VALUES %s`, values)
		if _, err := tx.ExecContext(ctx, sql, insertsHostSoftware...); err != nil {
			return ctxerr.Wrap(ctx, err, "insert host software")
		}
	}

	return nil
}

// update host_software when incoming software has a significantly more recent
// last opened timestamp (or didn't have on in currentMap). Note that it only
// processes software that is in both current and incoming maps, as the case
// where it is only in incoming is already handled by
// insertNewInstalledHostSoftwareDB.
func updateModifiedHostSoftwareDB(
	ctx context.Context,
	tx sqlx.ExtContext,
	hostID uint,
	currentMap map[string]fleet.Software,
	incomingMap map[string]fleet.Software,
	minLastOpenedAtDiff time.Duration,
) error {
	const stmt = `UPDATE host_software SET last_opened_at = ? WHERE host_id = ? AND software_id = ?`

	var keysToUpdate []string
	for key, newSw := range incomingMap {
		curSw, ok := currentMap[key]
		if !ok || newSw.LastOpenedAt == nil {
			// software must also exist in current map, and new software must have a
			// last opened at timestamp (otherwise we don't overwrite the old one)
			continue
		}

		if curSw.LastOpenedAt == nil || (*newSw.LastOpenedAt).Sub(*curSw.LastOpenedAt) >= minLastOpenedAtDiff {
			keysToUpdate = append(keysToUpdate, key)
		}
	}
	sort.Strings(keysToUpdate)

	for _, key := range keysToUpdate {
		curSw, newSw := currentMap[key], incomingMap[key]
		if _, err := tx.ExecContext(ctx, stmt, newSw.LastOpenedAt, hostID, curSw.ID); err != nil {
			return ctxerr.Wrap(ctx, err, "update host software")
		}
	}

	return nil
}

var dialect = goqu.Dialect("mysql")

// listSoftwareDB returns software installed on hosts. Use opts for pagination, filtering, and controlling
// fields populated in the returned software.
func listSoftwareDB(
	ctx context.Context,
	q sqlx.QueryerContext,
	opts fleet.SoftwareListOptions,
) ([]fleet.Software, error) {
	sql, args, err := selectSoftwareSQL(opts)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "sql build")
	}

	var results []softwareCVE
	if err := sqlx.SelectContext(ctx, q, &results, sql, args...); err != nil {
		return nil, ctxerr.Wrap(ctx, err, "select host software")
	}

	var softwares []fleet.Software
	ids := make(map[uint]int) // map of ids to index into softwares
	for _, result := range results {
		result := result // create a copy because we need to take the address to fields below

		idx, ok := ids[result.ID]
		if !ok {
			idx = len(softwares)
			softwares = append(softwares, result.Software)
			ids[result.ID] = idx
		}

		// handle null cve from left join
		if result.CVE != nil {
			cveID := *result.CVE
			cve := fleet.CVE{
				CVE:         cveID,
				DetailsLink: fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cveID),
			}
			if opts.IncludeCVEScores {
				cve.CVSSScore = &result.CVSSScore
				cve.EPSSProbability = &result.EPSSProbability
				cve.CISAKnownExploit = &result.CISAKnownExploit
			}
			softwares[idx].Vulnerabilities = append(softwares[idx].Vulnerabilities, cve)
		}
	}

	return softwares, nil
}

// softwareCVE is used for left joins with cve
type softwareCVE struct {
	fleet.Software
	CVE              *string  `db:"cve"`
	CVSSScore        *float64 `db:"cvss_score"`
	EPSSProbability  *float64 `db:"epss_probability"`
	CISAKnownExploit *bool    `db:"cisa_known_exploit"`
}

func selectSoftwareSQL(opts fleet.SoftwareListOptions) (string, []interface{}, error) {
	ds := dialect.
		From(goqu.I("software").As("s")).
		Select(
			"s.id",
			"s.name",
			"s.version",
			"s.source",
			"s.bundle_identifier",
			"s.release",
			"s.vendor",
			"s.arch",
			"scv.cpe_id", // for join on sub query
			goqu.COALESCE(goqu.I("scp.cpe"), "").As("generated_cpe"),
		)

	if opts.HostID != nil || opts.TeamID != nil {
		ds = ds.
			Join(
				goqu.I("host_software").As("hs"),
				goqu.On(
					goqu.I("hs.software_id").Eq(goqu.I("s.id")),
				),
			)
	}

	if opts.HostID != nil {
		ds = ds.
			SelectAppend("hs.last_opened_at").
			Where(goqu.I("hs.host_id").Eq(opts.HostID))
	}

	if opts.TeamID != nil {
		ds = ds.
			Join(
				goqu.I("hosts").As("h"),
				goqu.On(
					goqu.I("hs.host_id").Eq(goqu.I("h.id")),
				),
			).
			Where(goqu.I("h.team_id").Eq(opts.TeamID))
	}

	if opts.VulnerableOnly {
		ds = ds.
			Join(
				goqu.I("software_cpe").As("scp"),
				goqu.On(
					goqu.I("s.id").Eq(goqu.I("scp.software_id")),
				),
			).
			Join(
				goqu.I("software_cve").As("scv"),
				goqu.On(goqu.I("scp.id").Eq(goqu.I("scv.cpe_id"))),
			)
	} else {
		ds = ds.
			LeftJoin(
				goqu.I("software_cpe").As("scp"),
				goqu.On(
					goqu.I("s.id").Eq(goqu.I("scp.software_id")),
				),
			).
			LeftJoin(
				goqu.I("software_cve").As("scv"),
				goqu.On(goqu.I("scp.id").Eq(goqu.I("scv.cpe_id"))),
			)
	}

	if opts.IncludeCVEScores {
		ds = ds.
			LeftJoin(
				goqu.I("cve_scores").As("c"),
				goqu.On(goqu.I("c.cve").Eq(goqu.I("scv.cve"))),
			).
			SelectAppend(
				goqu.MAX("c.cvss_score").As("cvss_score"),                 // for ordering
				goqu.MAX("c.epss_probability").As("epss_probability"),     // for ordering
				goqu.MAX("c.cisa_known_exploit").As("cisa_known_exploit"), // for ordering
			)
	}

	if match := opts.MatchQuery; match != "" {
		match = likePattern(match)
		ds = ds.Where(
			goqu.Or(
				goqu.I("s.name").ILike(match),
				goqu.I("s.version").ILike(match),
				goqu.I("scv.cve").ILike(match),
			),
		)
	}

	if opts.WithHostCounts {
		ds = ds.
			Join(
				goqu.I("software_host_counts").As("shc"),
				goqu.On(goqu.I("s.id").Eq(goqu.I("shc.software_id"))),
			).
			Where(goqu.I("shc.hosts_count").Gt(0)).
			SelectAppend(
				goqu.I("shc.hosts_count"),
				goqu.I("shc.updated_at").As("counts_updated_at"),
			)

		if opts.TeamID != nil {
			ds = ds.Where(goqu.I("shc.team_id").Eq(opts.TeamID))
		} else {
			ds = ds.Where(goqu.I("shc.team_id").Eq(0))
		}
	}

	ds = ds.GroupBy(
		"s.id",
		"scv.cpe_id",
		"generated_cpe",
	)

	// Pagination is a bit more complex here due to left join with software_cve table and aggregated columns from cve_scores table.
	// Apply order by again after joining on sub query
	ds = appendListOptionsToSelect(ds, opts.ListOptions)

	// join on software_cve and cve_scores after apply pagination using the sub-query above
	ds = dialect.From(ds.As("s")).
		Select(
			"s.id",
			"s.name",
			"s.version",
			"s.source",
			"s.bundle_identifier",
			"s.release",
			"s.vendor",
			"s.arch",
			"s.generated_cpe",
			// omit s.cpe_id
			"scv.cve",
		).
		LeftJoin(
			goqu.I("software_cve").As("scv"),
			goqu.On(goqu.I("scv.cpe_id").Eq(goqu.I("s.cpe_id"))),
		).
		LeftJoin(
			goqu.I("cve_scores").As("c"),
			goqu.On(goqu.I("c.cve").Eq(goqu.I("scv.cve"))),
		)

	// select optional columns
	if opts.IncludeCVEScores {
		ds = ds.SelectAppend(
			"c.cvss_score",
			"c.epss_probability",
			"c.cisa_known_exploit",
		)
	}

	if opts.HostID != nil {
		ds = ds.SelectAppend(
			goqu.I("s.last_opened_at"),
		)
	}

	if opts.WithHostCounts {
		ds = ds.SelectAppend(
			goqu.I("s.hosts_count"),
			goqu.I("s.counts_updated_at"),
		)
	}

	ds = appendOrderByToSelect(ds, opts.ListOptions)

	return ds.ToSQL()
}

func countSoftwareDB(
	ctx context.Context,
	q sqlx.QueryerContext,
	opts fleet.SoftwareListOptions,
) (int, error) {
	opts.ListOptions = fleet.ListOptions{
		MatchQuery: opts.MatchQuery,
	}

	sql, args, err := selectSoftwareSQL(opts)
	if err != nil {
		return 0, ctxerr.Wrap(ctx, err, "sql build")
	}

	sql = `SELECT COUNT(DISTINCT s.id) FROM (` + sql + `) AS s`

	var count int
	if err := sqlx.GetContext(ctx, q, &count, sql, args...); err != nil {
		return 0, ctxerr.Wrap(ctx, err, "count host software")
	}

	return count, nil
}

func (ds *Datastore) LoadHostSoftware(ctx context.Context, host *fleet.Host) error {
	software, err := listSoftwareDB(ctx, ds.reader, fleet.SoftwareListOptions{HostID: &host.ID})
	if err != nil {
		return err
	}
	host.Software = software
	return nil
}

type softwareIterator struct {
	rows *sqlx.Rows
}

func (si *softwareIterator) Value() (*fleet.Software, error) {
	dest := fleet.Software{}
	err := si.rows.StructScan(&dest)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func (si *softwareIterator) Err() error {
	return si.rows.Err()
}

func (si *softwareIterator) Close() error {
	return si.rows.Close()
}

func (si *softwareIterator) Next() bool {
	return si.rows.Next()
}

func (ds *Datastore) AllSoftwareWithoutCPEIterator(ctx context.Context) (fleet.SoftwareIterator, error) {
	sql := `SELECT s.* FROM software s LEFT JOIN software_cpe sc on (s.id=sc.software_id) WHERE sc.id is null`
	// The rows.Close call is done by the caller once iteration using the
	// returned fleet.SoftwareIterator is done.
	rows, err := ds.reader.QueryxContext(ctx, sql) //nolint:sqlclosecheck
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "load host software")
	}
	return &softwareIterator{rows: rows}, nil
}

func (ds *Datastore) AddCPEForSoftware(ctx context.Context, software fleet.Software, cpe string) error {
	_, err := addCPEForSoftwareDB(ctx, ds.writer, software, cpe)
	return err
}

func addCPEForSoftwareDB(ctx context.Context, exec sqlx.ExecerContext, software fleet.Software, cpe string) (uint, error) {
	sql := `INSERT INTO software_cpe (software_id, cpe) VALUES (?, ?)`
	res, err := exec.ExecContext(ctx, sql, software.ID, cpe)
	if err != nil {
		return 0, ctxerr.Wrap(ctx, err, "insert software cpe")
	}
	id, _ := res.LastInsertId() // cannot fail with the mysql driver
	return uint(id), nil
}

func (ds *Datastore) AllCPEs(ctx context.Context) ([]string, error) {
	sql := `SELECT cpe FROM software_cpe`
	var cpes []string
	err := sqlx.SelectContext(ctx, ds.reader, &cpes, sql)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "loads cpes")
	}
	return cpes, nil
}

// InsertCVEForCPE inserts the cve into software_cve, linking it to all the
// provided cpes. It returns the number of new rows inserted or an error. If
// the CVE already existed for all CPEs, it would return 0, nil.
func (ds *Datastore) InsertCVEForCPE(ctx context.Context, cve string, cpes []string) (int64, error) {
	var totalCount int64
	for _, cpe := range cpes {
		var ids []uint
		err := sqlx.Select(ds.writer, &ids, `SELECT id FROM software_cpe WHERE cpe = ?`, cpe)
		if err != nil {
			return 0, err
		}

		values := strings.TrimSuffix(strings.Repeat("(?,?),", len(ids)), ",")
		sql := fmt.Sprintf(`INSERT IGNORE INTO software_cve (cpe_id, cve) VALUES %s`, values)

		var args []interface{}
		for _, id := range ids {
			args = append(args, id, cve)
		}
		res, err := ds.writer.ExecContext(ctx, sql, args...)
		if err != nil {
			return 0, ctxerr.Wrap(ctx, err, "insert software cve")
		}
		count, _ := res.RowsAffected()
		totalCount += count
	}

	return totalCount, nil
}

func (ds *Datastore) ListSoftware(ctx context.Context, opt fleet.SoftwareListOptions) ([]fleet.Software, error) {
	return listSoftwareDB(ctx, ds.reader, opt)
}

func (ds *Datastore) CountSoftware(ctx context.Context, opt fleet.SoftwareListOptions) (int, error) {
	return countSoftwareDB(ctx, ds.reader, opt)
}

// ListVulnerableSoftwareBySource lists all the vulnerable software that matches the given source.
func (ds *Datastore) ListVulnerableSoftwareBySource(ctx context.Context, source string) ([]fleet.SoftwareWithCPE, error) {
	var softwareCVEs []struct {
		fleet.Software
		CPE  uint   `db:"cpe_id"`
		CVEs string `db:"cves"`
	}
	if err := sqlx.SelectContext(ctx, ds.reader, &softwareCVEs, `
SELECT
    s.*,
    scv.cpe_id,
    GROUP_CONCAT(scv.cve SEPARATOR ',') as cves
FROM
    software s
    JOIN software_cpe scp ON scp.software_id = s.id
    JOIN software_cve scv ON scv.cpe_id = scp.id
WHERE
    s.source = ?
GROUP BY
    scv.cpe_id
	`, source); err != nil {
		return nil, ctxerr.Wrapf(ctx, err, "listing vulnerable software by source")
	}
	software := make([]fleet.SoftwareWithCPE, 0, len(softwareCVEs))
	for _, sc := range softwareCVEs {
		for _, cve := range strings.Split(sc.CVEs, ",") {
			sc.Vulnerabilities = append(sc.Vulnerabilities, fleet.CVE{
				CVE:         cve,
				DetailsLink: fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cve),
			})
		}
		software = append(software, fleet.SoftwareWithCPE{
			Software: sc.Software,
			CPEID:    sc.CPE,
		})
	}
	return software, nil
}

// DeleteVulnerabilitiesByCPECVE deletes the given list of vulnerabilities identified by CPE+CVE.
func (ds *Datastore) DeleteVulnerabilitiesByCPECVE(ctx context.Context, vulnerabilities []fleet.SoftwareVulnerability) error {
	if len(vulnerabilities) == 0 {
		return nil
	}

	sql := fmt.Sprintf(
		`DELETE FROM software_cve WHERE (cpe_id, cve) IN (%s)`,
		strings.TrimSuffix(strings.Repeat("(?,?),", len(vulnerabilities)), ","),
	)
	var args []interface{}
	for _, vulnerability := range vulnerabilities {
		args = append(args, vulnerability.CPEID, vulnerability.CVE)
	}
	if _, err := ds.writer.ExecContext(ctx, sql, args...); err != nil {
		return ctxerr.Wrapf(ctx, err, "deleting vulnerable software")
	}
	return nil
}

func (ds *Datastore) SoftwareByID(ctx context.Context, id uint, includeCVEScores bool) (*fleet.Software, error) {
	q := dialect.From(goqu.I("software").As("s")).
		Select(
			"s.id",
			"s.name",
			"s.version",
			"s.source",
			"s.bundle_identifier",
			"s.release",
			"s.vendor",
			"s.arch",
			"scv.cve",
		).
		LeftJoin(
			goqu.I("software_cpe").As("scp"),
			goqu.On(
				goqu.I("s.id").Eq(goqu.I("scp.software_id")),
			),
		).
		LeftJoin(
			goqu.I("software_cve").As("scv"),
			goqu.On(goqu.I("scp.id").Eq(goqu.I("scv.cpe_id"))),
		)

	if includeCVEScores {
		q = q.
			LeftJoin(
				goqu.I("cve_scores").As("c"),
				goqu.On(goqu.I("c.cve").Eq(goqu.I("scv.cve"))),
			).
			SelectAppend(
				"c.cvss_score",
				"c.epss_probability",
				"c.cisa_known_exploit",
			)
	}

	q = q.Where(goqu.I("s.id").Eq(id))

	sql, args, err := q.ToSQL()
	if err != nil {
		return nil, err
	}

	var results []softwareCVE
	err = sqlx.SelectContext(ctx, ds.reader, &results, sql, args...)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "get software")
	}

	if len(results) == 0 {
		return nil, ctxerr.Wrap(ctx, notFound("Software").WithID(id))
	}

	var software fleet.Software
	for i, result := range results {
		result := result // create a copy because we need to take the address to fields below

		if i == 0 {
			software = result.Software
		}

		if result.CVE != nil {
			cveID := *result.CVE
			cve := fleet.CVE{
				CVE:         cveID,
				DetailsLink: fmt.Sprintf("https://nvd.nist.gov/vuln/detail/%s", cveID),
			}
			if includeCVEScores {
				cve.CVSSScore = &result.CVSSScore
				cve.EPSSProbability = &result.EPSSProbability
				cve.CISAKnownExploit = &result.CISAKnownExploit
			}
			software.Vulnerabilities = append(software.Vulnerabilities, cve)
		}
	}

	return &software, nil
}

// CalculateHostsPerSoftware calculates the number of hosts having each
// software installed and stores that information in the software_host_counts
// table.
//
// After aggregation, it cleans up unused software (e.g. software installed
// on removed hosts, software uninstalled on hosts, etc.)
func (ds *Datastore) CalculateHostsPerSoftware(ctx context.Context, updatedAt time.Time) error {
	const (
		resetStmt = `
      UPDATE software_host_counts
      SET hosts_count = 0, updated_at = ?`

		// team_id is added to the select list to have the same structure as
		// the teamCountsStmt, making it easier to use a common implementation
		globalCountsStmt = `
      SELECT count(*), 0 as team_id, software_id
      FROM host_software
      WHERE software_id > 0
      GROUP BY software_id`

		teamCountsStmt = `
      SELECT count(*), h.team_id, hs.software_id
      FROM host_software hs
      INNER JOIN hosts h
      ON hs.host_id = h.id
      WHERE h.team_id IS NOT NULL AND hs.software_id > 0
      GROUP BY hs.software_id, h.team_id`

		insertStmt = `
      INSERT INTO software_host_counts
        (software_id, hosts_count, team_id, updated_at)
      VALUES
        %s
      ON DUPLICATE KEY UPDATE
        hosts_count = VALUES(hosts_count),
        updated_at = VALUES(updated_at)`

		valuesPart = `(?, ?, ?, ?),`

		cleanupSoftwareStmt = `
      DELETE s
      FROM software s
      LEFT JOIN software_host_counts shc
      ON s.id = shc.software_id
      WHERE
        shc.software_id IS NULL OR
        (shc.team_id = 0 AND shc.hosts_count = 0)`

		cleanupTeamStmt = `
      DELETE shc
      FROM software_host_counts shc
      LEFT JOIN teams t
      ON t.id = shc.team_id
      WHERE
        shc.team_id > 0 AND
        t.id IS NULL`
	)

	// first, reset all counts to 0
	if _, err := ds.writer.ExecContext(ctx, resetStmt, updatedAt); err != nil {
		return ctxerr.Wrap(ctx, err, "reset all software_host_counts to 0")
	}

	// next get a cursor for the global and team counts for each software
	stmtLabel := []string{"global", "team"}
	for i, countStmt := range []string{globalCountsStmt, teamCountsStmt} {
		rows, err := ds.reader.QueryContext(ctx, countStmt)
		if err != nil {
			return ctxerr.Wrapf(ctx, err, "read %s counts from host_software", stmtLabel[i])
		}
		defer rows.Close()

		// use a loop to iterate to prevent loading all in one go in memory, as it
		// could get pretty big at >100K hosts with 1000+ software each. Use a write
		// batch to prevent making too many single-row inserts.
		const batchSize = 100
		var batchCount int
		args := make([]interface{}, 0, batchSize*4)
		for rows.Next() {
			var (
				count  int
				teamID uint
				sid    uint
			)

			if err := rows.Scan(&count, &teamID, &sid); err != nil {
				return ctxerr.Wrapf(ctx, err, "scan %s row into variables", stmtLabel[i])
			}

			args = append(args, sid, count, teamID, updatedAt)
			batchCount++

			if batchCount == batchSize {
				values := strings.TrimSuffix(strings.Repeat(valuesPart, batchCount), ",")
				if _, err := ds.writer.ExecContext(ctx, fmt.Sprintf(insertStmt, values), args...); err != nil {
					return ctxerr.Wrapf(ctx, err, "insert %s batch into software_host_counts", stmtLabel[i])
				}

				args = args[:0]
				batchCount = 0
			}
		}
		if batchCount > 0 {
			values := strings.TrimSuffix(strings.Repeat(valuesPart, batchCount), ",")
			if _, err := ds.writer.ExecContext(ctx, fmt.Sprintf(insertStmt, values), args...); err != nil {
				return ctxerr.Wrapf(ctx, err, "insert last %s batch into software_host_counts", stmtLabel[i])
			}
		}
		if err := rows.Err(); err != nil {
			return ctxerr.Wrapf(ctx, err, "iterate over %s host_software counts", stmtLabel[i])
		}
		rows.Close()
	}

	// remove any unused software (global counts = 0)
	if _, err := ds.writer.ExecContext(ctx, cleanupSoftwareStmt); err != nil {
		return ctxerr.Wrap(ctx, err, "delete unused software")
	}

	// remove any software count row for teams that don't exist anymore
	if _, err := ds.writer.ExecContext(ctx, cleanupTeamStmt); err != nil {
		return ctxerr.Wrap(ctx, err, "delete software_host_counts for non-existing teams")
	}
	return nil
}

// HostsByCPEs returns a list of all hosts that have the software corresponding
// to at least one of the CPEs installed. It returns a minimal represention of
// matching hosts.
func (ds *Datastore) HostsByCPEs(ctx context.Context, cpes []string) ([]*fleet.HostShort, error) {
	queryStmt := `
    SELECT DISTINCT
      h.id,
      h.hostname
    FROM
      hosts h
    INNER JOIN
      host_software hs
    ON
      h.id = hs.host_id
    INNER JOIN
      software_cpe scp
    ON
      hs.software_id = scp.software_id
    WHERE
      scp.cpe IN (?)
    ORDER BY
      h.id`

	stmt, args, err := sqlx.In(queryStmt, cpes)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "building query args")
	}
	var hosts []*fleet.HostShort
	if err := sqlx.SelectContext(ctx, ds.reader, &hosts, stmt, args...); err != nil {

		return nil, ctxerr.Wrap(ctx, err, "select hosts by cpes")
	}
	return hosts, nil
}

func (ds *Datastore) HostsByCVE(ctx context.Context, cve string) ([]*fleet.HostShort, error) {
	query := `
SELECT
    DISTINCT(h.id), h.hostname
FROM
    hosts h
    JOIN host_software hs ON h.id = hs.host_id
    JOIN software_cpe scp ON scp.software_id = hs.software_id
    JOIN software_cve scv ON scv.cpe_id = scp.id
WHERE
    scv.cve = ?
ORDER BY
    h.id
`

	var hosts []*fleet.HostShort
	if err := sqlx.SelectContext(ctx, ds.reader, &hosts, query, cve); err != nil {
		return nil, ctxerr.Wrap(ctx, err, "select hosts by cves")
	}
	return hosts, nil
}

func (ds *Datastore) InsertCVEScores(ctx context.Context, cveScores []fleet.CVEScore) error {
	query := `
INSERT INTO cve_scores (cve, cvss_score, epss_probability, cisa_known_exploit)
VALUES %s
ON DUPLICATE KEY UPDATE
    cvss_score = VALUES(cvss_score),
    epss_probability = VALUES(epss_probability),
    cisa_known_exploit = VALUES(cisa_known_exploit)
`

	batchSize := 500
	for i := 0; i < len(cveScores); i += batchSize {
		end := i + batchSize
		if end > len(cveScores) {
			end = len(cveScores)
		}

		batch := cveScores[i:end]

		valuesFrag := strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?), ", len(batch)), ", ")
		var args []interface{}
		for _, score := range batch {
			args = append(args, score.CVE, score.CVSSScore, score.EPSSProbability, score.CISAKnownExploit)
		}

		query := fmt.Sprintf(query, valuesFrag)

		_, err := ds.writer.ExecContext(ctx, query, args...)
		if err != nil {
			return ctxerr.Wrap(ctx, err, "insert cve scores")
		}
	}

	return nil
}

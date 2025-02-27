// Package service holds the implementation of the fleet interface and HTTP
// endpoints for the API
package service

import (
	"context"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/fleetdm/fleet/v4/server/authz"
	"github.com/fleetdm/fleet/v4/server/config"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/logging"
	"github.com/fleetdm/fleet/v4/server/service/async"
	"github.com/fleetdm/fleet/v4/server/sso"
	kitlog "github.com/go-kit/kit/log"
)

var _ fleet.Service = (*Service)(nil)

// Service is the struct implementing fleet.Service. Create a new one with NewService.
type Service struct {
	ds             fleet.Datastore
	task           *async.Task
	carveStore     fleet.CarveStore
	resultStore    fleet.QueryResultStore
	liveQueryStore fleet.LiveQueryStore
	logger         kitlog.Logger
	config         config.FleetConfig
	clock          clock.Clock
	license        fleet.LicenseInfo

	osqueryLogWriter *logging.OsqueryLogger

	mailService     fleet.MailService
	ssoSessionStore sso.SessionStore

	failingPolicySet fleet.FailingPolicySet

	authz *authz.Authorizer

	jitterMu *sync.Mutex
	jitterH  map[time.Duration]*jitterHashTable

	geoIP fleet.GeoIP
}

func (s *Service) LookupGeoIP(ctx context.Context, ip string) *fleet.GeoLocation {
	return s.geoIP.Lookup(ctx, ip)
}

// NewService creates a new service from the config struct
func NewService(
	ctx context.Context,
	ds fleet.Datastore,
	task *async.Task,
	resultStore fleet.QueryResultStore,
	logger kitlog.Logger,
	osqueryLogger *logging.OsqueryLogger,
	config config.FleetConfig,
	mailService fleet.MailService,
	c clock.Clock,
	sso sso.SessionStore,
	lq fleet.LiveQueryStore,
	carveStore fleet.CarveStore,
	license fleet.LicenseInfo,
	failingPolicySet fleet.FailingPolicySet,
	geoIP fleet.GeoIP,
) (fleet.Service, error) {
	authorizer, err := authz.NewAuthorizer()
	if err != nil {
		return nil, fmt.Errorf("new authorizer: %w", err)
	}

	svc := &Service{
		ds:               ds,
		task:             task,
		carveStore:       carveStore,
		resultStore:      resultStore,
		liveQueryStore:   lq,
		logger:           logger,
		config:           config,
		clock:            c,
		osqueryLogWriter: osqueryLogger,
		mailService:      mailService,
		ssoSessionStore:  sso,
		license:          license,
		failingPolicySet: failingPolicySet,
		authz:            authorizer,
		jitterH:          make(map[time.Duration]*jitterHashTable),
		jitterMu:         new(sync.Mutex),
		geoIP:            geoIP,
	}
	return validationMiddleware{svc, ds, sso}, nil
}

func (s *Service) SendEmail(mail fleet.Email) error {
	return s.mailService.SendEmail(mail)
}

type validationMiddleware struct {
	fleet.Service
	ds              fleet.Datastore
	ssoSessionStore sso.SessionStore
}

// getAssetURL simply returns the base url used for retrieving image assets from fleetdm.com.
func getAssetURL() template.URL {
	return template.URL("https://fleetdm.com/images/permanent")
}

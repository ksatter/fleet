## Fleet 4.14.0 (May 9, 2022)

* Added beta support for Jira integration. This allows users to configure Fleet to
  automatically create a Jira issue when a new vulnerability (CVE) is detected on
  your hosts.

* Added a "Show query" button on the live query results page. This allows users to double-check the
  syntax used and compare this to their results without leaving the current view.

* Added a [Postman
  Collection](https://www.postman.com/fleetdm/workspace/fleet/collection/18010889-c5604fe6-7f6c-44bf-a60c-46650d358dde?ctx=documentation)
  for the Fleet API. This allows users to easily interact with Fleet's API routes so that they can
  build and test integrations.

* Added beta support for Fleet Desktop on Linux. Fleet Desktop allows the device user to see
information about their device. To add Fleet Desktop to a Linux device, first add the
`--fleet-desktop` flag to the `fleectl package` command to generate a Fleet-osquery installer that
includes Fleet Desktop. Then, open this installer on the device.

* Added `last_opened_at` property, for macOS software, to the **Host details** API route (`GET /hosts/{id}`).

* Improved the **Settings** pages in the the Fleet UI.

* Improved error message retuned when running `fleetctl query` command with missing or misspelled hosts.

* Improved the empty states and forms on the **Policies** page, **Queries** page, and **Host details** page in the Fleet UI.

- All duration settings returned by `fleetctl get config --include-server-config` were changed from
nanoseconds to an easy to read format.

* Fixed a bug in which the "Bundle identifier" tooltips displayed on **Host details > Software** did not
  render correctly.

* Fixed a bug in which the Fleet UI would render an empty Google Chrome profiles on the **Host details** page.

* Fixed a bug in which the Fleet UI would error when entering the "@" characters in the **Search targets** field.

* Fixed a bug in which a scheduled query would display the incorrect name when editing the query on
  the **Schedule** page.

* Fixed a bug in which a deprecation warning would be displayed when generating a `deb` or `rpm`
  Fleet-osquery package when running the `fleetctl package` command.

* Fixed a bug that caused panic errors when running the `fleet serve --debug` command.

## Fleet 4.13.2 (Apr 25, 2022)

* Fixed a bug with os versions not being updated. Affected deployments using MySQL < 5.7.22 or equivalent AWS RDS Aurora < 2.10.1.

## Fleet 4.13.1 (Apr 20, 2022)

* Fixed an SSO login issue introduced in 4.13.0.

* Fixed authorization errors encountered on the frontend login and live query pages.

## Fleet 4.13.0 (Apr 18, 2022)

### This is a security release.

* **Security**: Fixed several post-authentication authorization issues. Only Fleet Premium users that
  have team users are affected. Fleet Free users do not have access to the teams feature and are
  unaffected. See the following security advisory for details: https://github.com/fleetdm/fleet/security/advisories/GHSA-pr2g-j78h-84cr

* Improved performance of software inventory on Windows hosts.

* Added `basic​_auth.username` and `basic_auth.password` [Prometheus configuration options](https://fleetdm.com/docs/deploying/configuration#prometheus). The `GET
/metrics` API route is now disabled if these configuration options are left unspecified. 

* Fleet Premium: Add ability to specify a team specific "Destination URL" for policy automations.
This allows the user to configure Fleet to send a webhook request to a unique location for
policies that belong to a specific team. Documentation on what data is included the webhook
request and when the webhook request is sent can be found here on [fleedm.com/docs](https://fleetdm.com/docs/using-fleet/automations#vulnerability-automations)

* Added the ability to see the total number of hosts with a specific macOS version (ex. 12.3.1) on the
**Home > macOS** page. This information is also available via the [`GET /os_versions` API route](https://fleetdm.com/docs/using-fleet/rest-api#get-host-os-versions).

* Added the ability to sort live query results in the Fleet UI.

* Added a "Vulnerabilities" column to **Host details > Software** page. This allows the user see and search for specific vulnerabilities (CVEs) detected on a specific host.

* Updated vulnerability automations to fire anytime a vulnerability (CVE), that is detected on a
  host, was published to the
  National Vulnerability Database (NVD) in the last 30 days, is detected on a host. In previous
  versions of Fleet, vulnerability automations would fire anytime a CVE was published to NVD in the
  last 2 days.

* Updated the **Policies** page to ask the user to wait to see accurate passing and failing counts for new and recently edited policies.

* Improved API-only (integration) users by removing the requirement to reset these users' passwords
  before use. Documentation on how to use API-only users can be found here on [fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/fleetctl-cli#using-fleetctl-with-an-api-only-user).

* Improved the responsiveness of the Fleet UI by adding tablet screen width support for the **Software**,
  **Queries**, **Schedule**, **Policies**, **Host details**, **Settings > Teams**, and **Settings > Users** pages.

* Added Beta support for integrating with Jira to automatically create a Jira issue when a
  new vulnerability (CVE) is detected on a host in Fleet. 

* Added Beta support for Fleet Desktop on Windows. Fleet Desktop allows the device user to see
information about their device. To add Fleet Desktop to a Windows device, first add the
`--fleet-desktop` flag to the `fleectl package` command to generate a Fleet-osquery installer that
includes Fleet Desktop. Then, open this installer on the device.

* Fixed a bug in which downloading [Fleet's vulnerability database](https://github.com/fleetdm/nvd) failed if the destination directory specified
was not in the `tmp/` directory.

* Fixed a bug in which the "Updated at" time was not being updated for the "Mobile device management
(MDM) enrollment" and "Munki versions" information on the **Home > macOS** page.

* Fixed a bug in which Fleet would consider Docker network interfaces to be a host's primary IP address.

* Fixed a bug in which tables in the Fleet UI would present misaligned buttons.

* Fixed a bug in which Fleet failed to connect to Redis in standalone mode.
## Fleet 4.12.1 (Apr 4, 2022)

* Fixed a bug in which a user could not log in with basic authentication. This only affects Fleet deployments that use a [MySQL read replica](https://fleetdm.com/docs/deploying/configuration#my-sql).

## Fleet 4.12.0 (Mar 24, 2022)

* Added ability to update which platform (macOS, Windows, Linux) a policy is checked on.

* Added ability to detect compatibility for custom policies.

* Increased the default session duration to 5 days. Session duration can be updated using the
  [`session_duration` configuration option](https://fleetdm.com/docs/deploying/configuration#session-duration).

* Added ability to see the percentage of hosts that responded to a live query.

* Added ability for user's with [admin permissions](https://fleetdm.com/docs/using-fleet/permissions#user-permissions) to update any user's password.

* Added [`content_type_value` Kafka REST Proxy configuration
  option](https://fleetdm.com/docs/deploying/configuration#kafkarest-content-type-value) to allow
  the use of different versions of the Kafka REST Proxy.

* Added [`database_path` GeoIP configuration option](https://fleetdm.com/docs/deploying/configuration#database-path) to specify a GeoIP database. When configured,
  geolocation information is presented on the **Host details** page and in the `GET /hosts/{id}` API route.

* Added ability to retrieve a host's public IP address. This information is available on the **Host
  details** page and `GET /hosts/{id}` API route.

* Added instructions and materials needed to add hosts to Fleet using [plain osquery](https://fleetdm.com/docs/using-fleet/adding-hosts#plain-osquery). These instructions
can be found in **Hosts > Add hosts > Advanced** in the Fleet UI.

* Added Beta support for Fleet Desktop on macOS. Fleet Desktop allows the device user to see
  information about their device. To add Fleet Desktop to a macOS device, first add the
  `--fleet-desktop` flag to the `fleectl package` command to generate a Fleet-osquery installer that
  includes Fleet Desktop. Then, open this installer on the device.

* Reduced the noise of osquery status logs by only running a host vital query, which populate the
**Host details** page, when the query includes tables that are compatible with a specific host.

* Fixed a bug on the **Edit pack** page in which the "Select targets" element would display the hover effect for the wrong target.

* Fixed a bug on the **Software** page in which software items from deleted hosts were not removed.

* Fixed a bug in which the platform for Amazon Linux 2 hosts would be displayed incorrectly.

## Fleet 4.11.0 (Mar 7, 2022)

* Improved vulnerability processing to reduce the number of false positives for RPM packages on Linux hosts.

* Fleet Premium: Added a `teams` key to the `packs` yaml document to allow adding teams as targets when using CI/CD to manage query packs.

* Fleet premium: Added the ability to retrieve configuration for a specific team with the `fleetctl get team --name
<team-name-here>` command.

* Removed the expiration for API tokens for API-only users. API-only users can be created using the
  `fleetctl user create --api-only` command.

* Improved performance of the osquery query used to collect software inventory for Linux hosts.

* Updated the activity feed on the **Home page** to include add, edit, and delete policy activities.
  Activity information is also available in the `GET /activities` API route.

* Updated Kinesis logging plugin to append newline character to raw message bytes to properly format NDJSON for downstream consumers.

* Clarified why the "Performance impact" for some queries is displayed as "Undetermined" in the Fleet
  UI.

* Added instructions for using plain osquery to add hosts to Fleet in the Fleet View these instructions by heading to **Hosts > Add hosts > Advanced**.

* Fixed a bug in which uninstalling Munki from one or more hosts would result in inaccurate Munki
  versions displayed on the **Home > macOS** page.

* Fixed a bug in which a user, with access limited to one or more teams, was able to run a live query
against hosts in any team. This bug is not exposed in the Fleet UI and is limited to users of the
`POST run` API route. 

* Fixed a bug in the Fleet UI in which the "Select targets" search bar would not return the expected hosts.

* Fixed a bug in which global agent options were not updated correctly when editing these options in
the Fleet UI.

* Fixed a bug in which the Fleet UI would incorrectly tag some URLs as invalid.

* Fixed a bug in which the Fleet UI would attempt to connect to an SMTP server when SMTP was disabled.

* Fixed a bug on the Software page in which the "Hosts" column was not filtered by team.

* Fixed a bug in which global maintainers were unable to add and edit policies that belonged to a
  specific team.

* Fixed a bug in which the operating system version for some Linux distributions would not be
displayed properly.

* Fixed a bug in which configuring an identity provider name to a value shorter than 4 characters was
not allowed.

* Fixed a bug in which the avatar would not appear in the top navigation.


## Fleet 4.10.0 (Feb 13, 2022)

* Upgraded Go to 1.17.7 with security fixes for crypto/elliptic (CVE-2022-23806), math/big (CVE-2022-23772), and cmd/go (CVE-2022-23773). These are not likely to be high impact in Fleet deployments, but we are upgrading in an abundance of caution.

* Added aggregate software and vulnerability information on the new **Software** page.

* Added ability to see how many hosts have a specific vulnerable software installed on the
  **Software** page. This information is also available in the `GET /api/v1/fleet/software` API route.

* Added ability to send a webhook request if a new vulnerability (CVE) is
found on at least one host. Documentation on what data is included the webhook
request and when the webhook request is sent can be found here on [fleedm.com/docs](https://fleetdm.com/docs/using-fleet/automations#vulnerability-automations).

* Added aggregate Mobile Device Management and Munki data on the **Home** page.

* Added email and URL validation across the entire Fleet UI.

* Added ability to filter software by "Vulnerable" on the **Host details** page.

* Updated standard policy templates to use new naming convention. For example, "Is FileVault enabled on macOS
devices?" is now "Full disk encryption enabled (macOS)."

* Added db-innodb-status and db-process-list to `fleetctl debug` command.

* Fleet Premium: Added the ability to generate a Fleet installer and manage enroll secrets on the **Team details**
  page. 

* Added the ability for users with the observer role to view which platforms (macOS, Windows, Linux) a query
  is compatible with.

* Improved the experience for editing queries and policies in the Fleet UI.

* Improved vulnerability processing for NPM packages.

* Added supports triggering a webhook for newly detected vulnerabilities with a list of affected hosts.

* Added filter software by CVE.

* Added the ability to disable scheduled query performance statistics.

* Added the ability to filter the host summary information by platform (macOS, Windows, Linux) on the **Home** page.

* Fixed a bug in Fleet installers for Linux in which a computer restart would stop the host from
  reporting to Fleet.

* Made sure ApplyTeamSpec only works with premium deployments.

* Disabled MDM, Munki, and Chrome profile queries on unsupported platforms to reduce log noise.

* Properly handled paths in CVE URL prefix.

## Fleet 4.9.1 (Feb 2, 2022)

### This is a security release.

* **Security**: Fixed a vulnerability in Fleet's SSO implementation that could allow a malicious or compromised SAML Service Provider (SP) to log into Fleet as an existing Fleet user. See https://github.com/fleetdm/fleet/security/advisories/GHSA-ch68-7cf4-35vr for details.

* Allowed MSI packages generated by `fleetctl package` to reinstall on Windows without uninstall.

* Fixed a bug in which a team's scheduled queries didn't render correctly on the **Schedule** page.

* Fixed a bug in which a new policy would always get added to "All teams" rather than the selected team.

## Fleet 4.9.0 (Jan 21, 2022)

* Added ability to apply a `policy` yaml document so that GitOps workflows can be used to create and
  modify policies.

* Added ability to run a live query that returns 1,000+ results in the Fleet UI by adding
  client-side pagination to the results table.

* Improved the accuracy of query platform compatibility detection by adding recognition for queries
  with the `WITH` expression.

* Added ability to open a page in the Fleet UI in a new tab by "right-clicking" an item in the navigation.

* Improved the [live query API route (`GET /api/v1/queries/run`)](https://fleetdm.com/docs/using-fleet/rest-api#run-live-query) so that it successfully return results for Fleet
  instances using a load balancer by reducing the wait period to 25 seconds.

* Improved performance of the Fleet UI by updating loading states and reducing the number of requests
  made to the Fleet API.

* Improved performance of the MySQL database by updating the queries used to populate host vitals and
  caching the results.

* Added [`read_timeout` Redis configuration
  option](https://fleetdm.com/docs/deploying/configuration#redis-read-timeout) to customize the
  maximum amount of time Fleet should wait to receive a response from a Redis server.

* Added [`write_timeout` Redis configuration
  option](https://fleetdm.com/docs/deploying/configuration#redis-write-timeout) to customize the
  maximum amount of time Fleet should wait to send a command to a Redis server.

* Fixed a bug in which browser extensions (Google Chrome, Firefox, and Safari) were not included in
  software inventory.

* Improved the security of the **Organization settings** page by preventing the browser from requesting
  to save SMTP credentials.

* Fixed a bug in which an existing pack's targets were not cleaned up after deleting hosts, labels, and teams.

* Fixed a bug in which non-existent queries and policies would not return a 404 not found response.

### Performance

* Our testing demonstrated an increase in max devices served in our load test infrastructure to 70,000 from 60,000 in v4.8.0.

#### Load Test Infrastructure

* Fleet server
  * AWS Fargate
  * 2 tasks with 1024 CPU units and 2048 MiB of RAM.

* MySQL
  * Amazon RDS
  * db.r5.2xlarge

* Redis
  * Amazon ElastiCache 
  * cache.m5.large with 2 replicas (no cluster mode)

#### What was changed to accomplish these improvements?

* Optimized the updating and fetching of host data to only send and receive the bare minimum data
  needed. 

* Reduced the number of times host information is updated by caching more data.

* Updated cleanup jobs and deletion logic.

#### Future improvements

* At maximum DB utilization, we found that some hosts fail to respond to live queries. Future releases of Fleet will improve upon this.

## Fleet 4.8.0 (Dec 31, 2021)

* Added ability to configure Fleet to send a webhook request with all hosts that failed a
  policy. Documentation on what data is included the webhook
  request and when the webhook request is sent can be found here on [fleedm.com/docs](https://fleetdm.com/docs/using-fleet/automations).

* Added ability to find a user's device in Fleet by filtering hosts by email associated with a Google Chrome
  profile. Requires the [macadmins osquery
  extension](https://github.com/macadmins/osquery-extension) which comes bundled in [Fleet's osquery
  installers](https://fleetdm.com/docs/using-fleet/adding-hosts#osquery-installer). 
  
* Added ability to see a host's Google Chrome profile information using the [`GET
  api/v1/fleet/hosts/{id}/device_mapping` API
  route](https://fleetdm.com/docs/using-fleet/rest-api#get-host-device-mapping).

* Added ability to see a host's mobile device management (MDM) enrollment status, MDM server URL,
  and Munki version on a host's **Host details** page. Requires the [macadmins osquery
  extension](https://github.com/macadmins/osquery-extension) which comes bundled in [Fleet's osquery
  installers](https://fleetdm.com/docs/using-fleet/adding-hosts#osquery-installer). 

* Added ability to see a host's MDM and Munki information with the [`GET
  api/v1/fleet/hosts/{id}/macadmins` API
  route](https://fleetdm.com/docs/using-fleet/rest-api#list-mdm-and-munki-information-if-available).

* Improved the handling of certificates in the `fleetctl package` command by adding a check for a
  valid PEM file.

* Updated [Prometheus Go client library](https://github.com/prometheus/client_golang) which
  results in the following breaking changes to the [`GET /metrics` API
  route](https://fleetdm.com/docs/using-fleet/monitoring-fleet#metrics):
  `http_request_duration_microseconds` is now `http_request_duration_seconds_bucket`,
  `http_request_duration_microseconds_sum` is now `http_request_duration_seconds_sum`,
  `http_request_duration_microseconds_count` is now `http_request_duration_seconds_count`,
  `http_request_size_bytes` is now `http_request_size_bytes_bucket`, and `http_response_size_bytes`
  is now `http_response_size_bytes_bucket`

* Improved performance when searching and sorting hosts in the Fleet UI.

* Improved performance when running a live query feature by reducing the load on Redis.

* Improved performance when viewing software installed across all hosts in the Fleet
  UI.

* Fixed a bug in which the Fleet UI presented the option to download an undefined certificate in the "Generate installer" instructions.

* Fixed a bug in which database migrations failed when using MariaDB due to a migration introduced in Fleet 4.7.0.

* Fixed a bug that prevented hosts from checking in to Fleet when Redis was down.

## Fleet 4.7.0 (Dec 14, 2021)

* Added ability to create, modify, or delete policies in Fleet without modifying saved queries. Fleet
  4.7.0 introduces breaking changes to the `/policies` API routes to separate policies from saved
  queries in Fleet. These changes will not affect any policies previously created or modified in the
  Fleet UI.

* Turned on vulnerability processing for all Fleet instances with software inventory enabled.
  [Vulnerability processing in Fleet](https://fleetdm.com/docs/using-fleet/vulnerability-processing)
  provides the ability to see all hosts with specific vulnerable software installed. 

* Improved the performance of the "Software" table on the **Home** page.

* Improved performance of the MySQL database by changing the way a host's users information   is saved.

* Added ability to select from a library of standard policy templates on the **Policies** page. These
  pre-made policies ask specific "yes" or "no" questions about your hosts. For example, one of
  these policy templates asks "Is Gatekeeper enabled on macOS devices?"

* Added ability to ask whether or not your hosts have a specific operating system installed by
  selecting an operating system policy on the **Host details** page. For example, a host that is
  running macOS 12.0.1 will present a policy that asks "Is macOS 12.0.1 installed on macOS devices?"

* Added ability to specify which platform(s) (macOS, Windows, and/or Linux) a policy is checked on.

* Added ability to generate a report that includes which hosts are answering "Yes" or "No" to a 
  specific policy by running a policy's query as a live query.

* Added ability to see the total number of installed software software items across all your hosts.

* Added ability to see an example scheduled query result that is sent to your configured log
  destination. Select "Schedule a query" > "Preview data" on the **Schedule** page to see the 
  example scheduled query result.

* Improved the host's users information by removing users without login shells and adding users 
  that are not associated with a system group.

* Added ability to see a Fleet instance's missing migrations with the `fleetctl debug migrations`
  command. The `fleet serve` and `fleet prepare db` commands will now fail if any unknown migrations
  are detected.

* Added ability to see syntax errors as your write a query in the Fleet UI.

* Added ability to record a policy's resolution steps that can be referenced when a host answers "No" 
  to this policy.

* Added server request errors to the Fleet server logs to allow for troubleshooting issues with the 
Fleet server in non-debug mode.

* Increased default login session length to 24 hours.

* Fixed a bug in which software inventory and disk space information was not retrieved for Debian hosts.

* Fixed a bug in which searching for targets on the **Edit pack** page negatively impacted performance of 
  the MySQL database.

* Fixed a bug in which some Fleet migrations were incompatible with MySQL 8.

* Fixed a bug that prevented the creation of osquery installers for Windows (.msi) when a non-default 
  update channel is specified.

* Fixed a bug in which the "Software" table on the home page did not correctly filtering when a
  specific team was selected on the **Home** page.

* Fixed a bug in which users with "No access" in Fleet were presented with a perpetual 
  loading state in the Fleet UI.

## Fleet 4.6.2 (Nov 30, 2021)

* Improved performance of the **Home** page by removing total hosts count from the "Software" table.

* Improved performance of the **Queries** page by adding pagination to the list of queries.

* Fixed a bug in which the "Shell" column of the "Users" table on the **Host details** page would sometimes fail to update.

* Fixed a bug in which a host's status could quickly alternate between "Online" and "Offline" by increasing the grace period for host status.

* Fixed a bug in which some hosts would have a missing `host_seen_times` entry.

* Added an `after` parameter to the [`GET /hosts` API route](https://fleetdm.com/docs/using-fleet/rest-api#list-hosts) to allow for cursor pagination.

* Added a `disable_failing_policies` parameter to the [`GET /hosts` API route](https://fleetdm.com/docs/using-fleet/rest-api#list-hosts) to allow the API request to respond faster if failing policies count information is not needed.

## Fleet 4.6.1 (Nov 21, 2021)

* Fixed a bug (introduced in 4.6.0) in which Fleet used progressively more CPU on Redis, resulting in API and UI slowdowns and inconsistency.

* Made `fleetctl apply` fail when the configuration contains invalid fields.

## Fleet 4.6.0 (Nov 18, 2021)

* Fleet Premium: Added ability to filter aggregate host data such as platforms (macOS, Windows, and Linux) and status (online, offline, and new) the **Home** page. The aggregate host data is also available in the [`GET /host_summary API route`](https://fleetdm.com/docs/using-fleet/rest-api#get-hosts-summary).

* Fleet Premium: Added ability to move pending invited users between teams.

* Fleet Premium: Added `fleetctl updates rotate` command for rotation of keys in the updates system. The `fleetctl updates` command provides the ability to [self-manage an agent update server](https://fleetdm.com/docs/deploying/fleetctl-agent-updates).

* Enabled the software inventory by default for new Fleet instances. The software inventory feature can be turned on or off using the [`enable_software_inventory` configuration option](https://fleetdm.com/docs/using-fleet/vulnerability-processing#setup).

* Updated the JSON payload for the host status webhook by renaming the `"message"` property to `"text"` so that the payload can be received and displayed in Slack.

* Removed the deprecated `app_configs` table from Fleet's MySQL database. The `app_config_json` table has replaced it.

* Improved performance of the policies feature for Fleet instances with over 100,000 hosts.

* Added instructions in the Fleet UI for generating an osquery installer for macOS, Linux, or Windows. Documentation for generating an osquery installer and distributing the installer to your hosts to add them to Fleet can be found here on [fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/adding-hosts)

* Added ability to see all the software, and filter by vulnerable software, installed across all your hosts on the **Home** page. Each software's `name`, `version`, `hosts_count`, `vulnerabilities`, and more is also available in the [`GET /software` API route](https://fleetdm.com/docs/using-fleet/rest-api#software) and `fleetctl get software` command.

* Added ability to add, edit, and delete enroll secrets on the **Hosts** page.

* Added ability to see aggregate host data such as platforms (macOS, Windows, and Linux) and status (online, offline, and new) the **Home** page.

* Added ability to see all of the queries scheduled to run on a specific host on the **Host details** page immediately after a query is added to a schedule or pack.

* Added a "Shell" column to the "Users" table on the **Host details** page so users can now be filtered to see only those who have logged in.

* Packaged osquery's `certs.pem` in `fleetctl package` to improve TLS compatibility.

* Added support for packaging an osquery flagfile with `fleetctl package --osquery-flagfile`.

* Used "Fleet osquery" rather than "Orbit osquery" in packages generated by `fleetctl package`.

* Clarified that a policy in Fleet is a yes or no question you can ask about your hosts by replacing "Passing" and "Failing" text with "Yes" and "No" respectively on the **Policies** page and **Host details** page.

* Added ability to see the original author of a query on the **Query** page.

* Improved the UI for the "Software" table and "Policies" table on the **Host details** page so that it's easier to pivot to see all hosts with a specific software installed or answering "No" to a specific policy.

* Fixed a bug in which modifying a specific target for a live query, in target selector UI, would deselect a different target.

* Fixed a bug in which the user was navigated to a non existent page, in the Fleet UI, after saving a pack.

* Fixed a bug in which long software names in the "Software" table caused the bundle identifier tooltip to be inaccessible.

## Fleet 4.5.1 (Nov 10, 2021)

* Fixed performance issues with search filtering on manage queries page.

* Improved correctness and UX for query platform compatibility.

* Fleet Premium: Shows correct hosts when a team is selected.

* Fixed a bug preventing login for new SSO users.

* Added always return the `disabled` value in the `GET /api/v1/fleet/packs/{id}` API (previously it was
  sometimes left out).

## Fleet 4.5.0 (Nov 1, 2021)

* Fleet Premium: Added a Team admin user role. This allows users to delegate the responsibility of managing team members in Fleet. Documentation for the permissions associated with the Team admin and other user roles can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/permissions).

* Added Apache Kafka logging plugin. Documentation for configuring Kafka as a logging plugin can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#kafka-rest-proxy-logging). Thank you to Joseph Macaulay for adding this capability.

* Added support for [MinIO](https://min.io/) as a file carving backend. Documentation for configuring MinIO as a file carving backend can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/fleetctl-cli#minio). Thank you to Chandra Majumdar and Ben Edwards for adding this capability.

* Added support for generating a `.pkg` osquery installer on Linux without dependencies (beyond Docker) with the `fleetctl package` command.

* Improved the performance of vulnerability processing by making the process consume less RAM. Documentation for the vulnerability processing feature can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/vulnerability-processing).

* Added the ability to run a live query and receive results using only the Fleet REST API with a `GET /api/v1/fleet/queries/run` API route. Documentation for this new API route can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/rest-api#run-live-query).

* Added ability to see whether a specific host is "Passing" or "Failing" a policy on the **Host details** page. This information is also exposed in the `GET api/v1/fleet/hosts/{id}` API route. In Fleet, a policy is a "yes" or "no" question you can ask of all your hosts.

* Added the ability to quickly see the total number of "Failing" policies for a particular host on the **Hosts** page with a new "Issues" column. Total "Issues" are also revealed on a specific host's **Host details** page.

* Added the ability to see which platforms (macOS, Windows, Linux) a specific query is compatible with. The compatibility detected by Fleet is estimated based on the osquery tables used in the query.

* Added the ability to see whether your queries have a "Minimal," "Considerable," or "Excessive" performance impact on your hosts. Query performance information is only collected when a query runs as a scheduled query.

  * Running a "Minimal" query, very frequently, has little to no impact on your host's performance.

  * Running a "Considerable" query, frequently, can have a noticeable impact on your host's performance.

  * Running an "Excessive" query, even infrequently, can have a significant impact on your host’s performance.

* Added the ability to see a list of hosts that have a specific software version installed by selecting a software version on a specific host's **Host details** page. Software inventory is currently under a feature flag. To enable this feature flag, check out the [feature flag documentation](https://fleetdm.com/docs/deploying/configuration#feature-flags).

* Added the ability to see all vulnerable software detected across all your hosts with the `GET /api/v1/fleet/software` API route. Documentation for this new API route can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/rest-api#software).

* Added the ability to see the exact number of hosts that selected filters on the **Hosts** page. This ability is also available when using the `GET api/v1/fleet/hosts/count` API route.

* Added ability to automatically "Refetch" host vitals for a particular host without manually reloading the page.

* Added ability to connect to Redis with TLS. Documentation for configuring Fleet to use a TLS connection to the Redis server can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#redis-use-tls).

* Added `cluster_read_from_replica` Redis to specify whether or not to prefer readying from a replica when possible. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#redis-cluster-read-from-replica).

* Improved experience of the Fleet UI by preventing autocomplete in forms.

* Fixed a bug in which generating an `.msi` osquery installer on Windows would fail with a permission error.

* Fixed a bug in which turning on the host expiry setting did not remove expired hosts from Fleet.

* Fixed a bug in which the Software inventory for some host's was missing `bundle_identifier` information.

## Fleet 4.4.3 (Oct 21, 2021)

* Cached AppConfig in redis to speed up requests and reduce MySQL load.

* Fixed migration compatibility with MySQL GTID replication.

* Improved performance of software listing query.

* Improved MSI generation compatibility (for macOS M1 and some Virtualization configurations) in `fleetctl package`.

## Fleet 4.4.2 (Oct 14, 2021)

* Fixed migration errors under some MySQL configurations due to use of temporary tables.

* Fixed pagination of hosts on host dashboard.

* Optimized HTTP requests on host search.

## Fleet 4.4.1 (Oct 8, 2021)

* Fixed database migrations error when updating from 4.3.2 to 4.4.0. This did not effect upgrades
  between other versions and 4.4.0.

* Improved logging of errors in fleet serve.

## Fleet 4.4.0 (Oct 6, 2021)

* Fleet Premium: Teams Schedules show inherited queries from All teams (global) Schedule.

* Fleet Premium: Team Maintainers can modify and delete queries, and modify the Team Schedule.

* Fleet Premium: Team Maintainers can delete hosts from their teams.

* `fleetctl get hosts` now shows host additional queries if there are any.

* Update default homepage to new dashboard.

* Added ability to bulk delete hosts based on manual selection and applied filters.

* Added display macOS bundle identifiers on software table if available.

* Fixed scroll position when navigating to different pages.

* Fleet Premium: When transferring a host from team to team, clear the Policy results for that host.

* Improved stability of host vitals (fix cases of dropping users table, disk space).

* Improved performance and reliability of Policy database migrations.

* Provided a more clear error when a user tries to delete a query that is set in a Policy.

* Fixed query editor Delete key and horizontal scroll issues.

* Added cleaner buttons and icons on Manage Hosts Page.

## Fleet 4.3.2 (Sept 29, 2021)

* Improved database performance by reducing the amount of MySQL database queries when a host checks in.

* Fixed a bug in which users with the global maintainer role could not edit or save queries. In, Fleet 4.0.0, the Admin, Maintainer, and Observer user roles were introduced. Documentation for the permissions associated with each role can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/permissions). 

* Fixed a bug in which policies were checked about every second and add a `policy_update_interval` osquery configuration option. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#osquery-policy-update-interval).

* Fixed a bug in which edits to a query’s name, description, SQL did not appear until the user refreshed the Edit query page.

* Fixed a bug in which the hosts count for a label returned 0 after modifying a label’s name or description.

* Improved error message when attempting to create or edit a user with an email that already exists.

## Fleet 4.3.1 (Sept 21, 2021)

* Added `fleetctl get software` command to list all software and the detected vulnerabilities. The Vulnerable software feature is currently in Beta. For information on how to configure the Vulnerable software feature and how exactly Fleet processes vulnerabilities, check out the [Vulnerability processing documentation](https://fleetdm.com/docs/using-fleet/vulnerability-processing).

* Added `fleetctl vulnerability-data-stream` command to sync the vulnerabilities processing data streams by hand.

* Added `disable_data_sync` vulnerabilities configuration option to avoid downloading the data streams. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#disable-data-sync).

* Only shows observers the queries they have permissions to run on the **Queries** page. In, Fleet 4.0.0, the Admin, Maintainer, and Observer user roles were introduced. Documentation for the permissions associated with each role can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/using-fleet/permissions). 

* Added `connect_retry_attempts` Redis configuration option to retry failed connections. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#redis-connect-retry-attempts).

* Added `cluster_follow_redirections` Redis configuration option to follow cluster redirections. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#redis-cluster-follow-redirections).

* Added `max_jitter_percent` osquery configuration option to prevent all hosts from returning data at roughly the same time. Note that this improves the Fleet server performance, but it will now take longer for new labels to populate. Documentation for this configuration option can be found [here on fleetdm.com/docs](https://fleetdm.com/docs/deploying/configuration#osquery-max-jitter-percent).

* Improved the performance of database migrations.

* Reduced database load for label membership recording.

* Added fail early if the process does not have permissions to write to the logging file.

* Added completely skip trying to save a host's users and software inventory if it's disabled to reduce database load. 

* Fixed a bug in which team maintainers were unable to run live queries against the hosts assigned to their team(s).

* Fixed a bug in which a blank screen would intermittently appear on the **Hosts** page.

* Fixed a bug detecting disk space for hosts.

## Fleet 4.3.0 (Sept 13, 2021)

* Added Policies feature for detecting device compliance with organizational policies.

* Run/edit query experience has been completely redesigned.

* Added support for MySQL read replicas. This allows the Fleet server to scale to more hosts.

* Added configurable webhook to notify when a specified percentage of hosts have been offline for over the specified amount of days.

* Added `fleetctl package` command for building Orbit packages.

* Added enroll secret dialog on host dashboard.

* Exposed free disk space in gigs and percentage for hosts.

* Added 15-minute interval option on Schedule page.

* Cleaned up advanced options UI.

* 404 and 500 page now include buttons for Osquery community Slack and to file an issue

* Updated all empty and error states for cleaner UI.

* Added warning banners in Fleet UI and `fleetctl` for license expiration.

* Rendered query performance information on host vitals page pack section.

* Improved performance for app loading.

* Made team schedule names more user friendly and hide the stats for global and team schedules when showing host pack stats.

* Displayed `query_name` in when referencing scheduled queries for more consistent UI/UX.

* Query action added for observers on host vitals page.

* Added `server_settings.debug_host_ids` to gather more detailed information about what the specified hosts are sending to fleet.

* Allowed deeper linking into the Fleet application by saving filters in URL parameters.

* Renamed Basic Tier to Premium Tier, and Core Tier to Free Tier.

* Improved vulnerability detection compatibility with database configurations.

* MariaDB compatibility fixes: add explicit foreign key constraint and on cascade delete for host_software to allow for hosts with software to be deleted.

* Fixed migration that was incompatible with MySQL primary key requirements (default on DigitalOcean MySQL 5.8).

* Added 30 second SMTP timeout for mail configuration.

* Fixed display of platform Labels on manage hosts page

* Fixed a bug recording scheduled query statistics.

* When a label is removed, ignore query executions for that label.

* Added fleet serve config to change the redis connection timeout and keep alive interval.

* Removed hardcoded limits in label searches when targeting queries.

* Allow host users to be readded.

* Moved email template images from github to fleetdm.com.

* Fixed bug rendering CPU in host vitals.

* Updated the schema for host_users to allow for bulk inserts without locking, and allow for users without unique uid.

* When using dynamic vulnerability processing node, try to create the vulnerability.databases-path.

* Fixed `fleetctl get host <hostname>` to properly output JSON when the command line flag is supplied i.e `fleetctl get host --json foobar`

## Fleet 4.2.4 (Sept 2, 2021)

* Fixed a bug in which live queries would fail for deployments that use Redis Cluster.

* Fixed a bug in which some new Fleet deployments don't include the default global agent options. Documentation for global and team agent options can be found [here](https://fleetdm.com/docs/using-fleet/configuration-files#agent-options).

* Improved how a host's `users` are stored in MySQL to prevent deadlocks. This information is available in the "Users" table on each host's **Host details** page and in the `GET /api/v1/fleet/hosts/{id}` API route.

## Fleet 4.2.3 (Aug 23, 2021)

* Added ability to troubleshoot connection issues with the `fleetctl debug connection` command.

* Improved compatibility with MySQL variants (MariaDB, Aurora, etc.) by removing usage of JSON_ARRAYAGG.

* Fixed bug in which live queries would stop returning results if more than 5 seconds goes by without a result. This bug was introduced in 4.2.1.

* Eliminated double-logging of IP addresses in osquery endpoints.

* Update host details after transferring a host on the details page.

* Logged errors in osquery endpoints to improve debugging.

## Fleet 4.2.2 (Aug 18, 2021)

* Added a new built in label "All Linux" to target all hosts that run any linux flavor.

* Allowed finer grained configuration of the vulnerability processing capabilities.

* Fixed performance issues when updating pack contents.

* Fixed a build issue that caused external network access to panic in certain Linux distros (Ubuntu).

* Fixed rendering of checkboxes in UI when modals appear.

* Orbit: synced critical file writes to disk.

* Added "-o" flag to fleetctl convert command to ensure consistent output rather than relying on shell redirection (this was causing issues with file encodings).

* Fixed table column wrapping for manage queries page.

* Fixed wrapping in Label pills.

* Side panels in UI have a fresher look, Teams/Roles UI greyed out conditionally.

* Improved sorting in UI tables.

* Improved detection of CentOS in label membership.

## Fleet 4.2.1 (Aug 14, 2021)

* Fixed a database issue with MariaDB 10.5.4.

* Displayed updated team name after edit.

* When a connection from a live query websocket is closed, Fleet now timeouts the receive and handles the different cases correctly to not hold the connection to Redis.

* Added read live query results from Redis in a thread safe manner.

* Allows observers and maintainers to refetch a host in a team they belong to.

## Fleet 4.2.0 (Aug 11, 2021)

* Added the ability to simultaneously filter hosts by status (`online`, `offline`, `new`, `mia`) and by label on the **Hosts** page.

* Added the ability to filter hosts by team in the Fleet UI, fleetctl CLI tool, and Fleet API. *Available for Fleet Basic customers*.

* Added the ability to create a Team schedule in Fleet. The Schedule feature was released in Fleet 4.1.0. For more information on the new Schedule feature, check out the [Fleet 4.1.0 release blog post](https://blog.fleetdm.com/fleet-4-1-0-57dfa25e89c1). *Available for Fleet Basic customers*.

* Added Beta Vulnerable software feature which surfaces vulnerable software on the **Host details** page and the `GET /api/v1/fleet/hosts/{id}` API route. For information on how to configure the Vulnerable software feature and how exactly Fleet processes vulnerabilities, check out the [Vulnerability processing documentation](https://fleetdm.com/docs/using-fleet/vulnerability-processing#vulnerability-processing).

* Added the ability to see which logging destination is configured for Fleet in the Fleet UI. To see this information, head to the **Schedule** page and then select "Schedule a query." Configured logging destination information is also available in the `GET api/v1/fleet/config` API route.

* Improved the `fleetctl preview` experience by downloading Fleet's standard query library and loading the queries into the Fleet UI.

* Improved the user interface for the **Packs** page and **Queries** page in the Fleet UI.

* Added the ability to modify scheduled queries in your Schedule in Fleet. The Schedule feature was released in Fleet 4.1.0. For more information on the new Schedule feature, check out the [Fleet 4.1.0 release blog post](https://blog.fleetdm.com/fleet-4-1-0-57dfa25e89c1).

* Added the ability to disable the Users feature in Fleet by setting the new `enable_host_users` key to `true` in the `config` yaml, configuration file. For documentation on using configuration files in yaml syntax, check out the [Using yaml files in Fleet](https://fleetdm.com/docs/using-fleet/configuration-files#using-yaml-files-in-fleet) documentation.

* Improved performance of the Software inventory feature. Software inventory is currently under a feature flag. To enable this feature flag, check out the [feature flag documentation](https://fleetdm.com/docs/deploying/configuration#feature-flags).

* Improved performance of inserting `pack_stats` in the database. The `pack_stats` information is used to display "Frequency" and "Last run" information for a specific host's scheduled queries. You can find this information on the **Host details** page.

* Improved Fleet server logging so that it is more uniform.

* Fixed a bug in which a user with the Observer role was unable to run a live query.

* Fixed a bug that prevented the new **Home** page from being displayed in some Fleet instances.

* Fixed a bug that prevented accurate sorting issues across multiple pages on the **Hosts** page.

## Fleet 4.1.0 (Jul 26, 2021)

The primary additions in Fleet 4.1.0 are the new Schedule and Activity feed features.

Scheduled lets you add queries which are executed on your devices at regular intervals without having to understand or configure osquery query packs. For experienced Fleet and osquery users, the ability to create new, and modify existing, query packs is still available in the Fleet UI and fleetctl command-line tool. To reach the **Packs** page in the Fleet UI, head to **Schedule > Advanced**.

Activity feed adds the ability to observe when, and by whom, queries are changes, packs are created, live queries are run, and more. The Activity feed feature is located on the new Home page in the Fleet UI. Select the logo in the top right corner of the Fleet UI to navigate to the new **Home** page.

### New features breakdown

* Added ability to create teams and update their respective agent options and enroll secrets using the new `teams` yaml document and fleetctl. Available in Fleet Basic.

* Added a new **Home** page to the Fleet UI. The **Home** page presents a breakdown of the enrolled hosts by operating system.

* Added a "Users" table on the **Host details** page. The `username` information displayed in the "Users" table, as well as the `uid`, `type`, and `groupname` are available in the Fleet REST API via the `/api/v1/fleet/hosts/{id}` API route.

* Added ability to create a user without an invitation. You can now create a new user by heading to **Settings > Users**, selecting "Create user," and then choosing the "Create user" option.

* Added ability to search and sort installed software items in the "Software" table on the **Host details** page. 

* Added ability to delete a user from Fleet using a new `fleetctl user delete` command.

* Added ability to retrieve hosts' `status`, `display_text`, and `labels` using the `fleetctl get hosts` command.

* Added a new `user_roles` yaml document that allows users to manage user roles via fleetctl. Available in Fleet Basic.

* Changed default ordering of the "Hosts" table in the Fleet UI to ascending order (A-Z).

* Improved performance of the Software inventory feature by reducing the amount of inserts and deletes are done in the database when updating each host's
software inventory.

* Removed YUM and APT sources from Software inventory.

* Fixed an issue in which disabling SSO at the organization level would not disable SSO for all users.

* Fixed an issue with data migrations in which enroll secrets are duplicated after the `name` column was removed from the `enroll_secrets` table.

* Fixed an issue in which it was not possible to clear host settings by applying the `config` yaml document. This allows users to successfully remove the `additional_queries` property after adding it.

* Fixed printing of failed record count in AWS Kinesis/Firehose logging plugins.

* Fixed compatibility with GCP Memorystore Redis due to missing CLUSTER command.


## Fleet 4.0.1 (Jul 01, 2021)

* Fixed an issue in which migrations failed on MariaDB MySQL.

* Allowed `http` to be used when configuring `fleetctl` for `localhost`.

* Fixed a bug in which Team information was missing for hosts looked up by Label. 

## Fleet 4.0.0 (Jun 29, 2021)

The primary additions in Fleet 4.0.0 are the new Role-based access control (RBAC) and Teams features. 

RBAC adds the ability to define a user's access to features in Fleet. This way, more individuals in an organization can utilize Fleet with appropriate levels of access.

* Check out the [permissions documentation](https://github.com/fleetdm/fleet/blob/2f42c281f98e39a72ab4a5125ecd26d303a16a6b/docs/1-Using-Fleet/9-Permissions.md) for a breakdown of the new user roles.

Teams adds the ability to separate hosts into exclusive groups. This way, users can easily act on consistent groups of hosts. 

* Read more about the Teams feature in [the documentation here](https://github.com/fleetdm/fleet/blob/2f42c281f98e39a72ab4a5125ecd26d303a16a6b/docs/1-Using-Fleet/10-Teams.md).

### New features breakdown

* Added the ability to define a user's access to features in Fleet by introducing the Admin, Maintainer, and Observer roles. Available in Fleet Core.

* Added the ability to separate hosts into exclusive groups with the Teams feature. The Teams feature is available for Fleet Basic customers. Check out the list below for the new functionality included with Teams:

* Teams: Added the ability to enroll hosts to one team using team specific enroll secrets.

* Teams: Added the ability to manually transfer hosts to a different team in the Fleet UI.

* Teams: Added the ability to apply unique agent options to each team. Note that "osquery options" have been renamed to "agent options."

* Teams: Added the ability to grant users access to one or more teams. This allows you to define a user's access to specific groups of hosts in Fleet.

* Added the ability to create an API-only user. API-only users cannot access the Fleet UI. These users can access all Fleet API endpoints and `fleetctl` features. Available in Fleet Core.

* Added Redis cluster support. Available in Fleet Core.

* Fixed a bug that prevented the columns chosen for the "Hosts" table from persisting after logging out of Fleet.

### Upgrade plan

Fleet 4.0.0 is a major release and introduces several breaking changes and database migrations. The following sections call out changes to consider when upgrading to Fleet 4.0.0:

* The structure of Fleet's `.tar.gz` and `.zip` release archives have changed slightly. Deployments that use the binary artifacts may need to update scripts or tooling. The `fleetdm/fleet` Docker container maintains the same API.

* Use strictly `fleet` in Fleet's configuration, API routes, and environment variables. Users must update all usage of `kolide` in these items (deprecated since Fleet 3.8.0).

* Changeed your SAML SSO URI to use fleet instead of kolide . This is due to the changes to Fleet's API routes outlined in the section above.

* Changeed configuration option `server_tlsprofile` to `server_tls_compatibility`. This options previously had an inconsistent key name.

* Replaced the use of the `api/v1/fleet/spec/osquery/options` with `api/v1/fleet/config`. In Fleet 4.0.0, "osquery options" are now called "agent options." The new agent options are moved to the Fleet application config spec file and the `api/v1/fleet/config` API endpoint.

* Enrolled secrets no longer have "names" and are now either global or for a specific team. Hosts no longer store the “name” of the enroll secret that was used. Users that want to be able to segment hosts (for configuration, queries, etc.) based on the enrollment secret should use the Teams feature in Fleet Basic.

* JWT encoding is no longer used for session keys. Sessions now default to expiring in 4 hours of inactivity. `auth_jwt_key` and `auth_jwt_key_file` are no longer accepted as configuration.

* The `username` artifact has been removed in favor of the more recognizable `name` (Full name). As a result the `email` artifact is now used for uniqueness in Fleet. Upon upgrading to Fleet 4.0.0, existing users will have the `name` field populated with `username`. SAML users may need to update their username mapping to match user emails.

* As of Fleet 4.0.0, Fleet Device Management Inc. periodically collects anonymous information about your instance. Sending usage statistics is turned off by default for users upgrading from a previous version of Fleet. Read more about the exact information collected [here](https://github.com/fleetdm/fleet/blob/2f42c281f98e39a72ab4a5125ecd26d303a16a6b/docs/1-Using-Fleet/11-Usage-statistics.md).

## Fleet 4.0.0 RC3 (Jun 25, 2021)

Primarily teste the new release workflows. Relevant changelog will be updated for Fleet 4.0. 

## Fleet 4.0.0 RC2 (Jun 18, 2021)

The primary additions in Fleet 4.0.0 are the new Role-based access control (RBAC) and Teams features. 

RBAC adds the ability to define a user's access to features in Fleet. This way, more individuals in an organization can utilize Fleet with appropriate levels of access.

* Check out the [permissions documentation](https://github.com/fleetdm/fleet/blob/5e40afa8ba28fc5cdee813dfca53b84ee0ee65cd/docs/1-Using-Fleet/8-Permissions.md) for a breakdown of the new user roles.

Teams adds the ability to separate hosts into exclusive groups. This way, users can easily act on consistent groups of hosts. 

* Read more about the Teams feature in [the documentation here](https://github.com/fleetdm/fleet/blob/5e40afa8ba28fc5cdee813dfca53b84ee0ee65cd/docs/1-Using-Fleet/9-Teams.md).

### New features breakdown

* Added the ability to define a user's access to features in Fleet by introducing the Admin, Maintainer, and Observer roles. Available in Fleet Core.

* Added the ability to separate hosts into exclusive groups with the Teams feature. The Teams feature is available for Fleet Basic customers. Check out the list below for the new functionality included with Teams:

* Teams: Added the ability to enroll hosts to one team using team specific enroll secrets.

* Teams: Added the ability to manually transfer hosts to a different team in the Fleet UI.

* Teams: Added the ability to apply unique agent options to each team. Note that "osquery options" have been renamed to "agent options."

* Teams: Added the ability to grant users access to one or more teams. This allows you to define a user's access to specific groups of hosts in Fleet.

* Added the ability to create an API-only user. API-only users cannot access the Fleet UI. These users can access all Fleet API endpoints and `fleetctl` features. Available in Fleet Core.

* Added Redis cluster support. Available in Fleet Core.

* Fixed a bug that prevented the columns chosen for the "Hosts" table from persisting after logging out of Fleet.

### Upgrade plan

Fleet 4.0.0 is a major release and introduces several breaking changes and database migrations. 

* Use strictly `fleet` in Fleet's configuration, API routes, and environment variables. Users must update all usage of `kolide` in these items (deprecated since Fleet 3.8.0).

* Changed configuration option `server_tlsprofile` to `server_tls_compatability`. This option previously had an inconsistent key name.

* Replaced the use of the `api/v1/fleet/spec/osquery/options` with `api/v1/fleet/config`. In Fleet 4.0.0, "osquery options" are now called "agent options." The new agent options are moved to the Fleet application config spec file and the `api/v1/fleet/config` API endpoint.

* Enrolled secrets no longer have "names" and are now either global or for a specific team. Hosts no longer store the “name” of the enroll secret that was used. Users that want to be able to segment hosts (for configuration, queries, etc.) based on the enrollment secret should use the Teams feature in Fleet Basic.

* `auth_jwt_key` and `auth_jwt_key_file` are no longer accepted as configuration. 

* JWT encoding is no longer used for session keys. Sessions now default to expiring in 4 hours of inactivity.

### Known issues


There are currently no known issues in this release. However, we recommend only upgrading to Fleet 4.0.0-rc2 for testing purposes. Please file a GitHub issue for any issues discovered when testing Fleet 4.0.0!

## Fleet 4.0.0 RC1 (Jun 10, 2021)

The primary additions in Fleet 4.0.0 are the new Role-based access control (RBAC) and Teams features. 

RBAC adds the ability to define a user's access to information and features in Fleet. This way, more individuals in an organization can utilize Fleet with appropriate levels of access. Check out the [permissions documentation](https://fleetdm.com/docs/using-fleet/permissions) for a breakdown of the new user roles and their respective capabilities.

Teams adds the ability to separate hosts into exclusive groups. This way, users can easily observe and apply operations to consistent groups of hosts. Read more about the Teams feature in [the documentation here](https://fleetdm.com/docs/using-fleet/teams).

There are several known issues that will be fixed for the stable release of Fleet 4.0.0. Therefore, we recommend only upgrading to Fleet 4.0.0 RC1 for testing purposes. Please file a GitHub issue for any issues discovered when testing Fleet 4.0.0!

### New features breakdown

* Added the ability to define a user's access to information and features in Fleet by introducing the Admin, Maintainer, and Observer roles.

* Added the ability to separate hosts into exclusive groups with the Teams feature. The Teams feature is available for Fleet Basic customers. Check out the list below for the new functionality included with Teams:

* Added the ability to enroll hosts to one team using team specific enroll secrets.

* Added the ability to manually transfer hosts to a different team in the Fleet UI.

* Added the ability to apply unique agent options to each team. Note that "osquery options" have been renamed to "agent options."

* Added the ability to grant users access to one or more teams. This allows you to define a user's access to specific groups of hosts in Fleet.

### Upgrade plan

Fleet 4.0.0 is a major release and introduces several breaking changes and database migrations. 

* Used strictly `fleet` in Fleet's configuration, API routes, and environment variables. This means that you must update all usage of `kolide` in these items. The backwards compatibility introduced in Fleet 3.8.0 is no longer valid in Fleet 4.0.0.

* Changed configuration option `server_tlsprofile` to `server_tls_compatability`. This options previously had an inconsistent key name.

* Replaced the use of the `api/v1/fleet/spec/osquery/options` with `api/v1/fleet/config`. In Fleet 4.0.0, "osquery options" are now called "agent options." The new agent options are moved to the Fleet application config spec file and the `api/v1/fleet/config` API endpoint.

* Enrolled secrets no longer have "names" and are now either global or for a specific team. Hosts no longer store the “name” of the enroll secret that was used. Users that want to be able to segment hosts (for configuration, queries, etc.) based on the enrollment secret should use the Teams feature in Fleet Basic.

* `auth_jwt_key` and `auth_jwt_key_file` are no longer accepted as configuration. 

* JWT encoding is no longer used for session keys. Sessions now default to expiring in 4 hours of inactivity.

### Known issues

* Query packs cannot be targeted to teams.

## Fleet 3.13.0 (Jun 3, 2021)

* Improved performance of the `additional_queries` feature by moving `additional` query results into a separate table in the MySQL database. Please note that the `/api/v1/fleet/hosts` API endpoint now return only the requested `additional` columns. See documentation on the changes to the hosts API endpoint [here](https://github.com/fleetdm/fleet/blob/06b2e564e657492bfbc647e07eb49fd4efca5a03/docs/1-Using-Fleet/3-REST-API.md#list-hosts).

* Fixed a bug in which running a live query in the Fleet UI would return no results and the query would seem "hung" on a small number of devices.

* Improved viewing live query errors in the Fleet UI by including the “Errors” table in the full screen view.

* Improved `fleetctl preview` experience by adding the `fleetctl preview reset` and `fleetctl preview stop` commands to reset and stop simulated hosts running in Docker.

* Added several improvements to the Fleet UI including additional contrast on checkboxes and dropdown pills.

## Fleet 3.12.0 (May 19, 2021)

* Added scheduled queries to the _Host details_ page. Surface the "Name", "Description", "Frequency", and "Last run" information for each query in a pack that apply to a specific host.

* Improved the freshness of host vitals by adding the ability to "refetch" the data on the _Host details_ page.

* Added ability to copy log fields into Google Cloud Pub/Sub attributes. This allows users to use these values for subscription filters.

* Added ability to duplicate live query results in Redis. When the `redis_duplicate_results` configuration option is set to `true`, all live query results will be copied to an additional Redis Pub/Sub channel named LQDuplicate.

* Added ability to controls the server-side HTTP keepalive property. Turning off keepalives has helped reduce outstanding TCP connections in some deployments.

* Fixed an issue on the _Packs_ page in which Fleet would incorrectly handle the configured `server_url_prefix`.

## Fleet 3.11.0 (Apr 28, 2021)

* Improved Fleet performance by batch updating host seen time instead of updating synchronously. This improvement reduces MySQL CPU usage by ~33% with 4,000 simulated hosts and MySQL running in Docker.

* Added support for software inventory, introducing a list of installed software items on each host's respective _Host details_ page. This feature is flagged off by default (for now). Check out [the feature flag documentation for instructions on how to turn this feature on](https://fleetdm.com/docs/deploying/configuration#software-inventory).

* Added Windows support for `fleetctl` agent autoupdates. The `fleetctl updates` command provides the ability to self-manage an agent update server. Available for Fleet Basic customers.

* Made running common queries more convenient by adding the ability to select a saved query directly from a host's respective _Host details_ page.

* Fixed an issue on the _Query_ page in which Fleet would override the CMD + L browser hotkey.

* Fixed an issue in which a host would display an unreasonable time in the "Last fetched" column.

## Fleet 3.10.1 (Apr 6, 2021)

* Fixed a frontend bug that prevented the "Pack" page and "Edit pack" page from rendering in the Fleet UI. This issue occurred when the `platform` key, in the requested pack's configuration, was set to any value other than `darwin`, `linux`, `windows`, or `all`.

## Fleet 3.10.0 (Mar 31, 2021)

* Added `fleetctl` agent auto-updates beta which introduces the ability to self-manage an agent update server. Available for Fleet Basic customers.

* Added option for Identity Provider-Initiated (IdP-initiated) Single Sign-On (SSO).

* Improved logging. All errors are logged regardless of log level, some non-errors are logged regardless of log level (agent enrollments, runs of live queries etc.), and all other non-errors are logged on debug level.

* Improved login resilience by adding rate-limiting to login and password reset attempts and preventing user enumeration.

* Added Fleet version and Go version in the My Account page of the Fleet UI.

* Improved `fleetctl preview` to ensure the latest version of Fleet is fired up on every run. In addition, the Fleet UI is now accessible without having to click through browser security warning messages.

* Added prefer storing IPv4 addresses for host details.

## Fleet 3.9.0 (Mar 9, 2021)

* Added configurable host identifier to help with duplicate host enrollment scenarios. By default, Fleet's behavior does not change (it uses the identifier configured in osquery's `--host_identifier` flag), but for users with overlapping host UUIDs changing `--osquery_host_identifier` to `instance` may be helpful. 

* Made cool-down period for host enrollment configurable to control load on the database in scenarios in which hosts are using the same identifier. By default, the cooldown is off, reverting to the behavior of Fleet <=3.4.0. The cooldown can be enabled with `--osquery_enroll_cooldown`.

* Refreshed the Fleet UI with a new layout and horizontal navigation bar.

* Trimmed down the size of Fleet binaries.

* Improved handling of config_refresh values from osquery clients.

* Fixed an issue with IP addresses and host additional info dropping.

## Fleet 3.8.0 (Feb 25, 2021)

* Added search, sort, and column selection in the hosts dashboard.

* Added AWS Lambda logging plugin.

* Improved messaging about number of hosts responding to live query.

* Updated host listing API endpoints to support search.

* Added fixes to the `fleetctl preview` experience.

* Fixed `denylist` parameter in scheduled queries.

* Fixed an issue with errors table rendering on live query page.

* Deprecated `KOLIDE_` environment variable prefixes in favor of `FLEET_` prefixes. Deprecated prefixes continue to work and the Fleet server will log warnings if the deprecated variable names are used. 

* Deprecated `/api/v1/kolide` routes in favor of `/api/v1/fleet`. Deprecated routes continue to work and the Fleet server will log warnings if the deprecated routes are used. 

* Added Javascript source maps for development.

## Fleet 3.7.1 (Feb 3, 2021)

* Changed the default `--server_tls_compatibility` to `intermediate`. The new settings caused TLS connectivity issues for users in some environments. This new default is a more appropriate balance of security and compatibility, as recommended by Mozilla.

## Fleet 3.7.0 (Feb 3, 2021)

### This is a security release.

* **Security**: Fixed a vulnerability in which a malicious actor with a valid node key can send a badly formatted request that causes the Fleet server to exit, resulting in denial of service. See https://github.com/fleetdm/fleet/security/advisories/GHSA-xwh8-9p3f-3x45 and the linked content within that advisory.

* Added new Host details page which includes a rich view of a specific host’s attributes.

* Revealed live query errors in the Fleet UI and `fleetctl` to help target and diagnose hosts that fail.

* Added Helm chart to make it easier for users to deploy to Kubernetes.

* Added support for `denylist` parameter in scheduled queries.

* Added debug flag to `fleetctl` that enables logging of HTTP requests and responses to stderr.

* Improved the `fleetctl preview` experience to include adding containerized osquery agents, displaying login information, creating a default directory, and checking for Docker daemon status.

* Added improved error handling in host enrollment to make debugging issues with the enrollment process easier.

* Upgraded TLS compatibility settings to match Mozilla.

* Added comments in generated flagfile to add clarity to different features being configured.

* Fixed a bug in Fleet UI that allowed user to edit a scheduled query after it had been deleted from a pack.


## Fleet 3.6.0 (Jan 7, 2021)

* Added the option to set up an S3 bucket as the storage backend for file carving.

* Built Docker container with Fleet running as non-root user.

* Added support to read in the MySQL password and JWT key from a file.

* Improved the `fleetctl preview` experience by automatically completing the setup process and configuring fleetctl for users.

* Restructured the documentation into three top-level sections titled "Using Fleet," "Deployment," and "Contribution."

* Fixed a bug that allowed hosts to enroll with an empty enroll secret in new installations before setup was completed.

* Fixed a bug that made the query editor render strangely in Safari.

## Fleet 3.5.1 (Dec 14, 2020)

### This is a security release.

* **Security**: Introduced XML validation library to mitigate Go stdlib XML parsing vulnerability effecting SSO login. See https://github.com/fleetdm/fleet/security/advisories/GHSA-w3wf-cfx3-6gcx and the linked content within that advisory.

Follow up: Rotated `--auth_jwt_key` to invalidate existing sessions. Audit for suspicious activity in the Fleet server.

* **Security**: Prevents new queries from using the SQLite `ATTACH` command. This is a mitigation for the osquery vulnerability https://github.com/osquery/osquery/security/advisories/GHSA-4g56-2482-x7q8.

Follow up: Audit existing saved queries and logs of live query executions for possible malicious use of `ATTACH`. Upgrade osquery to 4.6.0 to prevent `ATTACH` queries from executing.

* Update icons and fix hosts dashboard for wide screen sizes.

## Fleet 3.5.0 (Dec 10, 2020)

* Refresh the Fleet UI with new colors, fonts, and Fleet logos.

* All releases going forward will have the fleectl.exe.zip on the release page.

* Added documentation for the authentication Fleet REST API endpoints.

* Added FAQ answers about the stress test results for Fleet, configuring labels, and resetting auth tokens.

* Fixed a performance issue users encountered when multiple hosts shared the same UUID by adding a one minute cooldown.

* Improved the `fleetctl preview` startup experience.

* Fixed a bug preventing the same query from being added to a scheduled pack more than once in the Fleet UI.


## Fleet 3.4.0 (Nov 18, 2020)

* Added NPM installer for `fleetctl`. Install via `npm install -g osquery-fleetctl`.

* Added `fleetctl preview` command to start a local test instance of the Fleet server with Docker.

* Added `fleetctl debug` commands and API endpoints for debugging server performance.

* Added additional_info_filters parameter to get hosts API endpoint for filtering returned additional_info.

* Updated package import paths from github.com/kolide/fleet to github.com/fleetdm/fleet.

* Added first of the Fleet REST API documentation.

* Added documentation on monitoring with Prometheus.

* Added documentation to FAQ for debugging database connection errors.

* Fixed fleetctl Windows compatibility issues.

* Fixed a bug preventing usernames from containing the @ symbol.

* Fixed a bug in 3.3.0 in which there was an unexpected database migration warning.

## Fleet 3.3.0 (Nov 05, 2020)

With this release, Fleet has moved to the new github.com/fleetdm/fleet
repository. Please follow changes and releases there.

* Added file carving functionality.

* Added `fleetctl user create` command.

* Added osquery options editor to admin pages in UI.

* Added `fleetctl query --pretty` option for pretty-printing query results. 

* Added ability to disable packs with `fleetctl apply`.

* Improved "Add New Host" dialog to walk the user step-by-step through host enrollment.

* Improved 500 error page by allowing display of the error.

* Added partial transition of branding away from "Kolide Fleet".

* Fixed an issue with case insensitive enroll secret and node key authentication.

* Fixed an issue with `fleetctl query --quiet` flag not actually suppressing output.


## Fleet 3.2.0 (Aug 08, 2020)

* Added `stdout` logging plugin.

* Added AWS `kinesis` logging plugin.

* Added compression option for `filesystem` logging plugin.

* Added support for Redis TLS connections.

* Added osquery host identifier to EnrollAgent logs.

* Added osquery version information to output of `fleetctl get hosts`.

* Added hostname to UI delete host confirmation modal.

* Updated osquery schema to 4.5.0.

* Updated osquery versions available in schedule query UI.

* Updated MySQL driver.

* Removed support for (previously deprecated) `old` TLS profile.

* Fixed cleanup of queries in bad state. This should resolve issues in which users experienced old live queries repeatedly returned to hosts. 

* Fixed output kind of `fleetctl get options`.

## Fleet 3.1.0 (Aug 06, 2020)

* Added configuration option to set Redis database (`--redis_database`).

* Added configuration option to set MySQL connection max lifetime (`--mysql_conn_max_lifetime`).

* Added support for printing a single enroll secret by name.

* Fixed bug with label_type in older fleetctl yaml syntax.

* Fixed bug with URL prefix and Edit Pack button. 

## Kolide Fleet 3.0.0 (Jul 23, 2020)

* Backend performance overhaul. The Fleet server can now handle hundreds of thousands of connected hosts.

* Pagination implemented in the web UI. This makes the UI usable for any host count supported by the backend.

* Added capability to collect "additional" information from hosts. Additional queries can be set to be updated along with the host detail queries. This additional information is returned by the API.

* Removed extraneous network interface information to optimize server performance. Users that require this information can use the additional queries functionality to retrieve it.

* Added "manual" labels implementation. Static labels can be set by providing a list of hostnames with `fleetctl`.

* Added JSON output for `fleetctl get` commands.

* Added `fleetctl get host` to retrieve details for a single host.

* Updated table schema for osquery 4.4.0.

* Added support for multiple enroll secrets.

* Logging verbosity reduced by default. Logs are now much less noisy.

* Fixed import of github.com/kolide/fleet Go packages for consumers outside of this repository.

## Kolide Fleet 2.6.0 (Mar 24, 2020)

* Added server logging for X-Forwarded-For header.

* Added `--osquery_detail_update_interval` to set interval of host detail updates.
  Set this (along with `--osquery_label_update_interval`) to a longer interval
  to reduce server load in large deployments.

* Fixed MySQL deadlock errors by adding retries and backoff to transactions.

## Kolide Fleet 2.5.0 (Jan 26, 2020)

* Added `fleetctl goquery` command to bring up the github.com/AbGuthrie/goquery CLI.

* Added ability to disable live queries in web UI and `fleetctl`.

* Added `--query-name` option to `fleetctl query`. This allows using the SQL from a saved query.

* Added `--mysql-protocol` flag to allow connection to MySQL by domain socket.

* Improved server logging. Add logging for creation of live queries. Add username information to logging for other endpoints.

* Allows CREATE queries in the web UI.

* Fixed a bug in which `fleetctl query` would exit before any results were returned when latency to the Fleet server was high.

* Fixed an error initializing the Fleet database when MySQL does not have event permissions.

* Deprecated "old" TLS profile.

## Kolide Fleet 2.4.0 (Nov 12, 2019)

* Added `--server_url_prefix` flag to configure a URL prefix to prepend on all Fleet URLs. This can be useful to run fleet behind a reverse-proxy on a hostname shared with other services.

* Added option to automatically expire hosts that have not checked in within a certain number of days. Configure this in the "Advanced Options" of "App Settings" in the browser UI.

* Added ability to search for hosts by UUID when targeting queries.

* Allows SAML IdP name to be as short as 4 characters.

## Kolide Fleet 2.3.0 (Aug 14, 2019)

### This is a security release.

* Security: Upgraded Go to 1.12.8 to fix CVE-2019-9512, CVE-2019-9514, and CVE-2019-14809.

* Added capability to export packs, labels, and queries as yaml in `fleetctl get` with the `--yaml` flag. Include queries with a pack using `--with-queries`.

* Modified email templates to load image assets from Github CDN rather than Fleet server (fixes broken images in emails when Fleet server is not accessible from email clients).

* Added warning in query UI when Redis is not functioning.

* Fixed minor bugs in frontend handling of scheduled queries.

* Minor styling changes to frontend.


## Kolide Fleet 2.2.0 (Jul 16, 2019)

* Added GCP PubSub logging plugin. Thanks to Michael Samuel for adding this capability.

* Improved escaping for target search in live query interface. It is now easier to target hosts with + and - characters in the name.

* Server and browser performance improved to reduced loading of hosts in frontend. Host status will only update on page load when over 100 hosts are present.

* Utilized details sent by osquery in enrollment request to more quickly display details of new hosts. Also fixes a bug in which hosts could not complete enrollment if certain platform-dependent options were used.

* Fixed a bug in which the default query runs after targets are edited.

## Kolide Fleet 2.1.2 (May 30, 2019)

* Prevented sending of SMTP credentials over insecure connection

* Added prefix generated SAML IDs with 'id' (improves compatibility with some IdPs)

## Kolide Fleet 2.1.1 (Apr 25, 2019)

* Automatically pulls AWS STS credentials for Firehose logging if they are not specified in config.

* Fixed bug in which log output did not include newlines separating characters.

* Fixed bug in which the default live query was run when navigating to a query by URL.

* Updated logic for setting primary NIC to ignore link-local or loopback interfaces.

* Disabled editing of logged in user email in admin panel (instead, use the "Account Settings" menu in top left).

* Fixed a panic resulting from an invalid config file path.

## Kolide Fleet 2.1.0 (Apr 9, 2019)

* Added capability to log osquery status and results to AWS Firehose. Note that this deprecated some existing logging configuration (`--osquery_status_log_file` and `--osquery_result_log_file`). Existing configurations will continue to work, but will be removed at some point.

* Automatically cleans up "incoming hosts" that do not complete enrollment.

* Fixed bug with SSO requests that caused issues with some IdPs.

* Hid built-in platform labels that have no hosts.

* Fixed references to Fleet documentation in emails.

* Minor improvements to UI in places where editing objects is disabled.

## Kolide Fleet 2.0.2 (Jan 17, 2019)

* Improved performance of `fleetctl query` with high host counts.

* Added `fleetctl get hosts` command to retrieve a list of enrolled hosts.

* Added support for Login SMTP authentication method (Used by Office365).

* Added `--timeout` flag to `fleetctl query`.

* Added query editor support for control-return shortcut to run query.

* Allowed preselection of hosts by UUID in query page URL parameters.

* Allowed username to be specified in `fleetctl setup`. Default behavior remains to use email as username.

* Fixed conversion of integers in `fleetctl convert`.

* Upgraded major Javascript dependencies.

* Fixed a bug in which query name had to be specified in pack yaml.

## Kolide Fleet 2.0.1 (Nov 26, 2018)

* Fixed a bug in which deleted queries appeared in pack specs returned by fleetctl.

* Fixed a bug getting entities with spaces in the name.

## Kolide Fleet 2.0.0 (Oct 16, 2018)

* Stable release of Fleet 2.0.

* Supports custom certificate authorities in fleetctl client.

* Added support for MySQL 8 authentication methods.

* Allows INSERT queries in editor.

* Updated UI styles.

* Fixed a bug causing migration errors in certain environments.

See changelogs for release candidates below to get full differences from 1.0.9
to 2.0.0.

## Kolide Fleet 2.0.0 RC5 (Sep 18, 2018)

* Fixed a security vulnerability that would allow a non-admin user to elevate privileges to admin level.

* Fixed a security vulnerability that would allow a non-admin user to modify other user's details.

* Reduced the information that could be gained by an admin user trying to port scan the network through the SMTP configuration.

* Refactored and add testing to authorization code.

## Kolide Fleet 2.0.0 RC4 (August 14, 2018)

* Exposed the API token (to be used with fleetctl) in the UI.

* Updated autocompletion values in the query editor.

* Fixed a longstanding bug that caused pack targets to sometimes update incorrectly in the UI.

* Fixed a bug that prevented deletion of labels in the UI.

* Fixed error some users encountered when migrating packs (due to deleted scheduled queries).

* Updated favicon and UI styles.

* Handled newlines in pack JSON with `fleetctl convert`.

* Improved UX of fleetctl tool.

* Fixed a bug in which the UI displayed the incorrect logging type for scheduled queries.

* Added support for SAML providers with whitespace in the X509 certificate.

* Fixed targeting of packs to individual hosts in the UI.

## Kolide Fleet 2.0.0 RC3 (June 21, 2018)

* Fixed a bug where duplicate queries were being created in the same pack but only one was ever delivered to osquery. A migration was added to delete duplicate queries in packs created by the UI.
  * It is possible to schedule the same query with different options in one pack, but only via the CLI.
  * If you thought you were relying on this functionality via the UI, note that duplicate queries will be deleted when you run migrations as apart of a cleanup fix. Please check your configurations and make sure to create any double-scheduled queries via the CLI moving forward.

* Fixed a bug in which packs created in UI could not be loaded by fleetctl.

* Fixed a bug where deleting a query would not delete it from the packs that the query was scheduled in.

## Kolide Fleet 2.0.0 RC2 (June 18, 2018)

* Fixed errors when creating and modifying packs, queries and labels in UI.

* Fixed an issue with the schema of returned config JSON.

* Handled newlines when converting query packs with fleetctl convert.

* Added last seen time hover tooltip in Fleet UI.

* Fixed a null pointer error when live querying via fleetctl.

* Explicitly set timezone in MySQL connection (improves timestamp consistency).

* Allowed native password auth for MySQL (improves compatibility with Amazon RDS).

## Kolide Fleet 2.0.0 (currently preparing for release)

The primary new addition in Fleet 2 is the new `fleetctl` CLI and file-format, which dramatically increases the flexibility and control that administrators have over their osquery deployment. The CLI and the file format are documented [in the Fleet documentation](https://fleetdm.com/docs/using-fleet/fleetctl-cli).

### New Features

* New `fleetctl` CLI for managing your entire osquery workflow via CLI, API, and source controlled files!
  * You can use `fleetctl` to manage osquery packs, queries, labels, and configuration.

* In addition to the CLI, Fleet 2.0.0 introduces a new file format for articulating labels, queries, packs, options, etc. This format is designed for composability, enabling more effective sharing and re-use of intelligence.

```yaml
apiVersion: v1
kind: query
spec:
  name: pending_updates
  query: >
    select value
    from plist
    where
      path = "/Library/Preferences/ManagedInstalls.plist" and
      key = "PendingUpdateCount" and
      value > "0";
```

* Run live osquery queries against arbitrary subsets of your infrastructure via the `fleetctl query` command.

* Use `fleetctl setup`, `fleetctl login`, and `fleetctl logout` to manage the authentication life-cycle via the CLI.

* Use `fleetctl get`, `fleetctl apply`, and `fleetctl delete` to manage the state of your Fleet data.

* Manage any osquery option you want and set platform-specific overrides with the `fleetctl` CLI and file format.

### Upgrade Plan

* Managing osquery options via the UI has been removed in favor of the more flexible solution provided by the CLI. If you have customized your osquery options with Fleet, there is [a database migration](./server/datastore/mysql/migrations/data/20171212182458_MigrateOsqueryOptions.go) which will port your existing data into the new format when you run `fleet prepare db`. To download your osquery options after migrating your database, run `fleetctl get options > options.yaml`. Further modifications to your options should occur in this file and it should be applied with `fleetctl apply -f ./options.yaml`.

## Kolide Fleet 1.0.8 (May 3, 2018)

* Osquery 3.0+ compatibility!

* Included RFC822 From header in emails (for email authentication)

## Kolide Fleet 1.0.7 (Mar 30, 2018)

* Now supports FileAccesses in FIM configuration.

* Now populates network interfaces on windows hosts in host view.

* Added flags for configuring MySQL connection pooling limits.

* Fixed bug in which shard and removed keys are dropped in query packs returned to osquery clients.

* Fixed handling of status logs with unexpected fields.

## Kolide Fleet 1.0.6 (Dec 4, 2017)

* Added remote IP in the logs for all osqueryd/launcher requests. (#1653)

* Fixed bugs that caused logs to sometimes be omitted from the logwriter. (#1636, #1617)

* Fixed a bug where request bodies were not being explicitly closed. (#1613)

* Fixed a bug where SAML client would create too many HTTP connections. (#1587)

* Fixed bug in which default query was run instead of entered query. (#1611)

* Added pagination to the Host browser pages for increased performance. (#1594)

* Fixed bug rendering hosts when clock speed cannot be parsed. (#1604)

## Kolide Fleet 1.0.5 (Oct 17, 2017)

* Renamed the binary from kolide to fleet.

* Added support for Kolide Launcher managed osquery nodes.

* Removed license requirements.

* Updated documentation link in the sidebar to point to public GitHub documentation.

* Added FIM support.

* Title on query page correctly reflects new or edit mode.

* Fixed issue on new query page where last query would be submitted instead of current.

* Fixed issue where user menu did not work on Firefox browser.

* Fixed issue cause SSO to fail for ADFS.

* Fixed issue validating signatures in nested SAML assertions..

## Kolide 1.0.4 (Jun 1, 2017)

* Added feature that allows users to import existing Osquery configuration files using the [configimporter](https://github.com/kolide/configimporter) utility.

* Added support for Osquery decorators.

* Added SAML single sign on support.

* Improved online status detection.

  The Kolide server now tracks the `distributed_interval` and `config_tls_refresh` values for each individual host (these can be different if they are set via flagfile and not through Kolide), to ensure that online status is represented as accurately as possible.

* Kolide server now requires `--auth_jwt_key` to be specified at startup.

  If no JWT key is provided by the user, the server will print a new suggested random JWT key for use.

* Fixed bug in which deleted packs were still displayed on the query sidebar.

* Fixed rounding error when showing % of online hosts.

* Removed --app_token_key flag.

* Fixed issue where heavily loaded database caused host authentication failures.

* Fixed issue where osquery sends empty strings for integer values in log results.

## Kolide 1.0.3 (April 3, 2017)

* Log rotation is no longer the default setting for Osquery status and results logs. To enable log rotation use the `--osquery_enable_log_rotation` flag.

* Added a debug endpoint for collecting performance statistics and profiles.

  When `kolide serve --debug` is used, additional handlers will be started to provide access to profiling tools. These endpoints are authenticated with a randomly generated token that is printed to the Kolide logs at startup. These profiling tools are not intended for general use, but they may be useful when providing performance-related bug reports to the Kolide developers.

* Added a workaround for CentOS6 detection.

  Osquery 2.3.2 incorrectly reports an empty value for `platform` on CentOS6 hosts. We added a workaround to properly detect platform in Kolide, and also [submitted a fix](https://github.com/facebook/osquery/pull/3071) to upstream osquery.

* Ensured hosts enroll in labels immediately even when `distributed_interval` is set to a long interval.

* Optimizations reduce the CPU and DB usage of the manage hosts page.

* Managed packs page now loads much quicker when a large number of hosts are enrolled.

* Fixed bug with the "Reset Options" button.

* Fixed 500 error resulting from saving unchanged options.

* Improved validation for SMTP settings.

* Added command line support for `modern`, `intermediate`, and `old` TLS configuration
profiles. The profile is set using the following command line argument.
```
--server_tls_compatibility=modern
```
See https://wiki.mozilla.org/Security/Server_Side_TLS for more information on the different profile options.

* The Options Configuration item in the sidebar is now only available to admin users.

  Previously this item was visible to non-admin users and if selected, a blank options page would be displayed since server side authorization constraints prevent regular users from viewing or changing options.

* Improved validation for the Kolide server URL supplied in setup and configuration.

* Fixed an issue importing osquery configurations with numeric values represented as strings in JSON.

## Kolide 1.0.2 (March 14, 2017)

* Fixed an issue adding additional targets when querying a host

* Shows loading spinner while newly added Host Details are saved

* Shows a generic computer icon when when referring to hosts with an unknown platform instead of the text "All"

* Kolide will now warn on startup if there are database migrations not yet completed.

* Kolide will prompt for confirmation before running database migrations.

  To disable this, use `kolide prepare db --no-prompt`.

* Kolide now supports emoji, so you can 🔥 to your heart's content.

* When setting the platform for a scheduled query, selecting "All" now clears individually selected platforms.

* Updated Host details cards UI

* Lowered HTTP timeout settings.

  In an effort to provide a more resilient web server, timeouts are more strictly enforced by the Kolide HTTP server (regardless of whether or not you're using the built-in TLS termination).

* Hardened TLS server settings.

  For customers using Kolide's built-in TLS server (if the `server.tls` configuration is `true`), the server was hardened to only accept modern cipher suites as recommended by [Mozilla](https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility).

* Improve the mechanism used to calculate whether or not hosts are online.

  Previously, hosts were categorized as "online" if they had been seen within the past 30 minutes. To make the "online" status more representative of reality, hosts are marked "online" if the Kolide server has heard from them within two times the lowest polling interval as described by the Kolide-managed osquery configuration. For example, if you've configured osqueryd to check-in with Kolide every 10 seconds, only hosts that Kolide has heard from within the last 20 seconds will be marked "online".

* Updated Host details cards UI

* Added support for rotating the osquery status and result log files by sending a SIGHUP signal to the kolide process.

* Fixed Distributed Query compatibility with load balancers and Safari.

  Customers running Kolide behind a web balancer lacking support for websockets were unable to use the distributed query feature. Also, in certain circumstances, Safari users with a self-signed cert for Kolide would receive an error. This release add a fallback mechanism from websockets using SockJS for improved compatibility.

* Fixed issue with Distributed Query Pack results full screen feature that broke the browser scrolling abilities.

* Fixed bug in which host counts in the sidebar did not match up with displayed hosts.

## Kolide 1.0.1 (February 27, 2017)

* Fixed an issue that prevented users from replacing deleted labels with a new label of the same name.

* Improved the reliability of IP and MAC address data in the host cards and table.

* Added full screen support for distributed query results.

* Enabled users to double click on queries and packs in a table to see their details.

* Reprompted for a password when a user attempts to change their email address.

* Automatically decorates the status and result logs with the host's UUID and hostname.

* Fixed an issue where Kolide users on Safari were unable to delete queries or packs.

* Improved platform detection accuracy.

  Previously Kolide was determining platform based on the OS of the system osquery was built on instead of the OS it was running on. Please note: Offline hosts may continue to report an erroneous platform until they check-in with Kolide.

* Fixed bugs where query links in the pack sidebar pointed to the wrong queries.

* Improved MySQL compatibility with stricter configurations.

* Allows users to edit the name and description of host labels.

* Added basic table autocompletion when typing in the query composer.

* Now support MySQL client certificate authentication. More details can be found in the [Configuring the Fleet binary docs](./docs/infrastructure/configuring-the-fleet-binary.md).

* Improved security for user-initiated email address changes.

  This improvement ensures that only users who own an email address and are logged in as the user who initiated the change can confirm the new email.

  Previously it was possible for Administrators to also confirm these changes by clicking the confirmation link.

* Fixed an issue where the setup form rejects passwords with certain characters.

  This change resolves an issue where certain special characters like "." where rejected by the client-side JS that controls the setup form.

* Now automatically logs in the user once initial setup is completed.

### v0.11
##### September 19, 2021

#### :tada: Features

- Added metrics cleanup TTLS for cleaning up metrics from memory
- Support existing secret in helm chart
- Added support for extra environment variables

#### :beetle: Bug Fixes

- Fixed metrics overriding issue
- Helm testing is optional now

### v0.10
##### April 12, 2021

#### :beetle: Bug Fixes

- Fixed duplication of metrics in worker capacity used

#### :tada: Features

- Histogram metrics has boolean switch to enable or disable
- Optimized helm charts with latest standards

### v0.9
##### January 13, 2021

#### :beetle: Bug Fixes

- Fixed interface to a nil conversion error in Druid HTTP Endpoint
- Invalid argument(0) to Intn goroutine

#### :tada: Features

- Added datasource row count as metrics
- Docker compose setup to fix
- Added new task count metrics

### v0.8
##### July 11, 2020

#### :beetle: Bug Fixes

- Changed Port env variable to remove k8s default env variables conflicts
- Fixed servicemonitor file in helm template
- Refactored metrics ingestion from podName to hostName
- Fixed druid exporter name in helm template

### v0.7
##### June 22, 2020

#### :beetle: Bug Fixes

- Fixed JSON parser failure issue
- Fixed duplicate task metric's server 500 issue

### v0.6
##### June 11, 2020

#### :tada: Features

- Added TLS verify skip support flag

#### :beetle: Bug Fixes

- Bad username and password error statement for authentication request failure

### v0.5
##### June 7, 2020

#### :tada: Features

- Added HTTP basic auth support
- Added HTTP TLS support
- Revamped README with useful information
- Added log format and log level support

#### :beetle: Bug Fixes

- Fixed exporter panic issue in case druid url is not available

### v0.4
##### May 21, 2020

#### :tada: Features

- Add datasource label in druid emitted metrics
- Updated new logo

### v0.3
##### May 19, 2020

#### :tada: Features

- Added version flag to determine version

#### :beetle: Bug Fixes

- Fixed duplicate values error
- Replaced int with float as some of the druid emitted values are float

### v0.2
##### May 3, 2020

#### :tada: Features

- Druid http emitter based metrics support
- Mux handler for webserver
- Structured logging with logger module
- Flag and environment variable support
- Main page addition

### v0.1
##### April 24, 2020

#### :tada: Features

- Initial release
- Druid API based metrics
- Dockerization
- Initial lebel of documentation
- Initial Grafana Dashboard

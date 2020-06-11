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

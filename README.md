## Druid Exporter

[![CircleCI](https://circleci.com/gh/opstree/druid-exporter.svg?style=shield)](https://circleci.com/gh/opstree/druid-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/opstree/druid-exporter)](https://goreportcard.com/report/github.com/opstree/druid-exporter)
[![Maintainability](https://api.codeclimate.com/v1/badges/f3d9db298411361ca84a/maintainability)](https://codeclimate.com/github/opstree/druid-exporter/maintainability)
[![Docker Repository on Quay](https://img.shields.io/badge/container-ready-green "Docker Repository on Quay")](https://quay.io/repository/opstree/redis-operator)
[![Apache License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A Golang based exporter captures druid API related metrics and receives druid-emitting HTTP JSON data and converts it into the Prometheus time series format.

[Grafana Dashboard](https://grafana.com/grafana/dashboards/12155)

## Purpose

The aim of creating this druid exporter was to capture all of the metrics that druid exposes or emits. [JMX Exporter](https://github.com/prometheus/jmx_exporter) can be used to obtain JVM based metrics.

You can find examples of JMX exporter metrics [here](https://gist.github.com/iamabhishek-dubey/5ef19d3db9deb25475a80c9ff5c79262)

## Features

- Configuration values with flags and environment variables
- JSON centered logging system
- Druid API based metrics
  - Health Status
  - Datasource
  - Segments
  - Supervisors
  - Tasks
- Druid HTTP Emitted metrics
  - Broker
  - Historical
  - Ingestion(Kafka)
  - Coordination
  - Sys

## Available Options or Flags

See the help page with `--help`

```shell
$ ./druid-exporter --help
usage: druid-exporter [<flags>]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
      --druid.user=""      HTTP basic auth username. (Only if it is set)
      --druid.password=""  HTTP basic auth password. (Only if it is set)
      --cert=""            A pem encoded certificate file. (Only if tls is configured)
      --key=""             A pem encoded key file. (Only if tls is configured)
      --ca=""              A pem encoded CA's certificate file. (Only if tls is configured)
  -d, --druid.uri="http://druid.opstreelabs.in"  
                           URL of druid router or coordinator
  -p, --port="8080"        Port to listen druid exporter. (Default - 8080)
  -l, --log.level="info"   Log level for druid exporter. (Default: info)
  -f, --log.format="text"  Log format for druid exporter, text or json. (Default: text)
      --version            Show application version.
```

| **Option** | **Default Value** | **Environment Variable** | **Description** |
|------------|-------------------|--------------------------|-----------------|
| --help | - | - | Show context-sensitive help (also try --help-long and --help-man) |
| --druid.uri | http://druid.opstreelabs.in | DRUID_URL | URL of druid's coordinator service or router service |
| --debug | false | - | Enable to log into debug mode |
| --port | 8080 | DRUID_EXPORTER_PORT | Listening port of the druid exporter |

## Druid Configuration Changes

There are some changes needed in the druid cluster to exploit full capabilities of druid exporter. Druid emits the metrics to different emitters. So, in druid database, we must allow the http emitter.

If you are using the druid properties file you must add this entry to the file `common.properties`:-

```properties
druid.emitter.http.recipientBaseUrl=http://<druid_exporter_url>:<druid_exporter_port>/druid
druid.emitter=http
```

In case configuration of druid are managed by environment variables:-

```properties
druid_emitter_http_recipientBaseUrl=http://<druid_exporter_url>:<druid_exporter_port>/druid
druid_emitter=http
```

## Installing

Druid exporter can be download from [release](https://github.com/opstree/druid-exporter/releases)

To run the druid exporter:-

```shell
# Export the Druid Coordinator or Router URL
export DRUID_URL="http://druid.opstreelabs.in"
export DRUID_EXPORTER_PORT="8080"

./druid-exporter [<flags>]
```

## Building From Source

Requires 1.13 => go version to compile code from source.

```shell
make build-code
```

## Building Docker Image

This druid exporter has support for docker as well. The docker image can simply built by

```shell
make build-image
```

For running the druid exporter from docker image:-

```shell
# Execute docker run command
docker run -itd --name druid-exporter -e DRUID_URL="http://druid.opstreelabs.in" \
quay.io/opstree/druid-exporter:latest
```

## Kubernetes Deployment

The Kubernetes deployment and service manifests are present under the **[manifests](./manifets)** directory and you can deploy it on Kubernetes from there.

To deploy it on Kubernetes we need some basic sets of command:-

```shell
# Kubernetes deployment creation
kubectl apply -f manifests/deployment.yaml -n my_awesome_druid_namespace

# Kubernetes service creation
kubectl apply -f manifests/service.yaml -n my_awesome_druid_namespace
```

## Dashboard Screenshot

<p align="center">
  <img src="./static/grafana-screenshot.svg">
</p>

## Roadmap

- [ ] Add docker compose setup for druid and druid exporter
- [ ] Unit test cases should be in place
- [ ] Integration test cases should be in place

## Development

Please see our [development documentation](./DEVELOPMENT.md)

## Release

Please see our [release documentation](./CHANGELOG.md) for details

## Contact

If you have any suggestion or query. Contact us at

opensource@opstree.com

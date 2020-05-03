## Druid Exporter

[![CircleCI](https://circleci.com/gh/opstree/druid-exporter.svg?style=shield)](https://circleci.com/gh/opstree/druid-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/opstree/druid-exporter)](https://goreportcard.com/report/github.com/opstree/druid-exporter)
[![Maintainability](https://api.codeclimate.com/v1/badges/f3d9db298411361ca84a/maintainability)](https://codeclimate.com/github/opstree/druid-exporter/maintainability)
[![Docker Repository on Quay](https://img.shields.io/badge/container-ready-green "Docker Repository on Quay")](https://quay.io/repository/opstree/redis-operator)
[![Apache License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A Golang based exporter to scrap druid metrics in Prometheus format.

[Grafana Dashboard](https://grafana.com/grafana/dashboards/12155)

## Installing

Druid exporter can be download from [release](https://github.com/opstree/druid-exporter/releases)

To run the druid exporter:-

```shell
# Export the Druid Coordinator or Router URL
export DRUID_URL="http://druid.opstreelabs.in"

./druid-exporter
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

## Available Data Groups

These are the available datagroups present in druid exporter.

|**Name**|**Description**|
|--------|---------------|
| druid_health_status | To check if druid cluster is healthy or not |
| druid_datasource | All datasources present in druid cluster |
| druid_tasks | All druid's supervisors tasks status |
| druid_supervisors | Complete information of druid's supervisors |
| druid_segment_count | How many segments are available in each datasource |
| druid_segment_size | Size of druid segments in each datasource |
| druid_segment_replicated_size | Replicated size of druid segments in each datasource |

## Development

Please see our [development documentation](./DEVELOPMENT.md)

## Release

Please see our [release documentation](./CHANGELOG.md) for details

## Contact

If you have any suggestion or query. Contact us at

opensource@opstree.com

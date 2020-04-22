<p align="left">
  <img src="./static/druid-exporter-logo.svg" height="60" width="128">
</p>

# Druid Exporter

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

|**Name**|**Description**|
|--------|---------------|
| druid_health_status | To check if druid cluster is healthy or not |
| druid_datasource | All datasources present in druid cluster |
| druid_tasks | All druid's supervisors tasks status |
| druid_supervisors | Complete information of druid's supervisors |
| druid_segement_count | How many segments are available in each datasource |
| druid_segement_size | Size of druid segments in each datasource |
| druid_segement_replicated_size | Replicated size of druid segments in each datasource |


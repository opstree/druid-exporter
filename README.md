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
docker run -itd --name druid-exporter -e DRUID_URL="http://druid.opstreelabs.in" quay.io/opstree/druid-exporter:latest
```
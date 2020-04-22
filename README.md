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
export DRUID_URL=http://druid.opstreelabs.in

./druid-exporter
```

package main

import (
	"net/http"
	"druid-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	data := collector.Collector()

	prometheus.MustRegister(data)
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	http.ListenAndServe(":8080", nil)
}

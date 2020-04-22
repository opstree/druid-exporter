package main

import (
	"druid-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	data := collector.Collector()
	prometheus.MustRegister(data)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

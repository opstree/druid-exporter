package main

import (
	"druid-exporter/listener"
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
	http.Handle("/druid/endpoint", listener.ListenerEndpoint)
	log.Printf("Opstree's druid exporter is listing on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

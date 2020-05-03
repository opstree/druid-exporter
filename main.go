package main

import (
	"druid-exporter/listener"
	// "druid-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	druidEmittedData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "druid_emitted_metrics",
			Help: "Druid emitted metrics from druid emitter",
		},[]string{"metric_name", "service", "host"},
	)
)

func init() {
	// data := collector.Collector()
	// prometheus.MustRegister(data)
	prometheus.MustRegister(druidEmittedData)
}

func main() {
	router := mux.NewRouter()
	router.Handle("/druid", listener.ListenerEndpoint(druidEmittedData))
	router.Handle("/metrics", promhttp.Handler())
	// prometheus.Register(druidEmittedData)
	log.Printf("Opstree's druid exporter is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

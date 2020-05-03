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
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "druid_emitted_metrics",
			Help: "Druid emitted metrics from druid emitter",
		},
		[]string{"metric"},
	)
)

func init() {
	prometheus.MustRegister(hdFailures)
}

func main() {
	router := mux.NewRouter()
	router.Handle("/druid", listener.ListenerEndpoint(hdFailures))
	router.Handle("/metrics", promhttp.Handler())
	prometheus.Register(hdFailures)
	log.Printf("Opstree's druid exporter is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
	// data := collector.Collector()
	// recieveData := listener.Collector()
	// prometheus.MustRegister(data)
	// prometheus.MustRegister(recieveData)
	// http.Handle("/druid", listener.Collect())
	// http.Handle("/metrics", promhttp.Handler())
	// log.Printf("Opstree's druid exporter is listening on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

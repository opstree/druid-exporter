package main

import (
	"druid-exporter/listener"
	"druid-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug   = kingpin.Flag("debug", "Enable debug mode.").Bool()
	port = kingpin.Flag("port", "Port for druid exporter").Default("8080").OverrideDefaultFromEnvar("DRUID_EXPORTER_PORT").Short('p').String()
	druidEmittedData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "druid_emitted_metrics",
			Help: "Druid emitted metrics from druid emitter",
		},[]string{"metric_name", "service", "host"},
	)
)

func init() {
	getDruidAPIdata := collector.Collector()
	prometheus.MustRegister(getDruidAPIdata)
	prometheus.MustRegister(druidEmittedData)
}

func main() {
	kingpin.Parse()
	router := mux.NewRouter()
	router.Handle("/druid", listener.ListenerEndpoint(druidEmittedData))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Druid Exporter</title></head>
			<body>
			<h1>Druid Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Printf("Opstree's druid exporter is listening on :" + *port)
	log.Fatal(http.ListenAndServe(":" + *port, router))
}

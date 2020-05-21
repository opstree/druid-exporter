package main

import (
	"druid-exporter/collector"
	"druid-exporter/listener"
	"druid-exporter/logger"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

var (
	debug            = kingpin.Flag("debug", "Enable debug mode.").Bool()
	port             = kingpin.Flag("port", "Port for druid exporter").Default("8080").OverrideDefaultFromEnvar("DRUID_EXPORTER_PORT").Short('p').String()
	druidEmittedData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "druid_emitted_metrics",
			Help: "Druid emitted metrics from druid emitter",
		}, []string{"metric_name", "service", "host", "datasource"},
	)
)

func init() {
	getDruidAPIdata := collector.Collector()
	prometheus.MustRegister(getDruidAPIdata)
	prometheus.MustRegister(druidEmittedData)
}

func main() {
	druidLogger := logger.GetLoggerInterface()
	kingpin.Version("0.3")
	kingpin.Parse()
	router := mux.NewRouter()
	router.Handle("/druid", listener.DruidHTTPEndpoint(druidEmittedData))
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
	level.Info(druidLogger).Log("msg", "Druid exporter started listening on :"+*port)
	level.Error(druidLogger).Log("msg", http.ListenAndServe("0.0.0.0:"+*port, router))
}

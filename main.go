package main

import (
	"druid-exporter/collector"
	"druid-exporter/listener"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
)

var (
	port             = kingpin.Flag("port", "Port to listen druid exporter. (Default - 8080)").Default("8080").OverrideDefaultFromEnvar("DRUID_EXPORTER_PORT").Short('p').String()
	logLevel         = kingpin.Flag("log.level", "Log level for druid exporter. (Default: info)").Default("info").OverrideDefaultFromEnvar("LOG_LEVEL").Short('l').String()
	logFormat        = kingpin.Flag("log.format", "Log format for druid exporter, text or json. (Default: text)").Default("text").OverrideDefaultFromEnvar("LOG_FORMAT").Short('f').String()
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
	kingpin.Version("0.3")
	kingpin.Parse()
	parsedLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Errorf("log-level flag has invalid value %s", *logLevel)
	} else {
		logrus.SetLevel(parsedLevel)
	}
	if *logFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
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
	logrus.Infof("Druid exporter started listening on: %v", *port)
	logrus.Infof("Metrics endpoint - http://0.0.0.0:%v/metrics", *port)
	logrus.Infof("Druid emitter endpoint - http://0.0.0.0:%v/druid", *port)
	http.ListenAndServe("0.0.0.0:"+*port, router)
}

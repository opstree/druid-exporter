package main

import (
	"druid-exporter/collector"
	"druid-exporter/listener"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Flag(
		"port",
		"Port to listen druid exporter, EnvVar - DRUID_EXPORTER_PORT. (Default - 8080)",
	).Default("8080").OverrideDefaultFromEnvar("DRUID_EXPORTER_PORT").Short('p').String()
	logLevel = kingpin.Flag(
		"log.level",
		"Log level for druid exporter, EnvVar - LOG_LEVEL. (Default: info)",
	).Default("info").OverrideDefaultFromEnvar("LOG_LEVEL").Short('l').String()
	logFormat = kingpin.Flag(
		"log.format",
		"Log format for druid exporter, text or json, EnvVar - LOG_FORMAT. (Default: text)",
	).Default("text").OverrideDefaultFromEnvar("LOG_FORMAT").Short('f').String()
	druidEmittedData = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "druid_emitted_metrics",
			Help: "Druid emitted metrics from druid emitter",
		}, []string{"host", "metric_name", "service", "datasource"},
	)
)

func init() {
	prometheus.MustRegister(druidEmittedData)
}

func main() {
	kingpin.Version("0.5")
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

	dnsCache := cache.New(5*time.Minute, 10*time.Minute)
	router := mux.NewRouter()
	getDruidAPIdata := collector.Collector()
	handlerFunc := newHandler(*getDruidAPIdata)
	router.Handle("/druid", listener.DruidHTTPEndpoint(druidEmittedData, dnsCache))
	router.Handle("/metrics", promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
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

func newHandler(metrics collector.MetricCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		getDruidAPIdata := collector.Collector()
		registry.MustRegister(getDruidAPIdata)
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

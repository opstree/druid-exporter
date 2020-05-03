package listener

import (
	"time"
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/prometheus/client_golang/prometheus"
)

// DruidEmittedData is the data structure which druid sends over HTTP
type DruidEmittedData struct {
	Feed           string    `json:"feed"`
	Timestamp      time.Time `json:"timestamp"`
	Service        string    `json:"service"`
	Host           string    `json:"host"`
	Version        string    `json:"version"`
	Metric         string    `json:"metric"`
	Value          int       `json:"value"`
	GcGen          []string  `json:"gcGen"`
	GcGenSpaceName string    `json:"gcGenSpaceName"`
	GcName         []string  `json:"gcName"`
}

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus         *prometheus.Desc
	DataSourceCount           *prometheus.Desc
	DruidTasks                *prometheus.Desc
	DruidSupervisors          *prometheus.Desc
	DruidSegmentCount         *prometheus.Desc
	DruidSegmentSize          *prometheus.Desc
	DruidSegmentReplicateSize *prometheus.Desc
	DruidEmitted              *prometheus.Desc
	Handler                   func(w http.ResponseWriter, r *http.Request)
}

// ListenerEndpoint is the endpoint to listen all druid metrics
func (collector *MetricCollector)ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var ch chan<- prometheus.Metric
	var druidData []DruidEmittedData
	if req.Method == "POST" {
		jsonDecoder := json.NewDecoder(req.Body)
		err := jsonDecoder.Decode(&druidData)
		if err != nil {
			log.Error().Msg("Error while decoding JSON sent by druid")
		}
	}
	for _, data := range druidData {
		ch <- prometheus.MustNewConstMetric(collector.DruidEmitted, prometheus.GaugeValue, float64(data.Value), data.Metric, data.Service)
	}
}

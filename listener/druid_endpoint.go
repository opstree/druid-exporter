package listener

import (
	"druid-exporter/logger"
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strings"
	"time"
)

// DruidEmittedData is the data structure which druid sends over HTTP
type DruidEmittedData struct {
	Feed      string    `json:"feed"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Host      string    `json:"host"`
	Version   string    `json:"version"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
}

// DruidHTTPEndpoint is the endpoint to listen all druid metrics
func DruidHTTPEndpoint(gauge *prometheus.GaugeVec) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		druidLogger := logger.GetLoggerInterface()
		var druidData []DruidEmittedData
		if req.Method == "POST" {
			jsonDecoder := json.NewDecoder(req.Body)
			err := jsonDecoder.Decode(&druidData)
			if err != nil {
				level.Error(druidLogger).Log("msg", "Error in decoding JSON sent by druid", "err", err)
			}
			for _, data := range druidData {
				gauge.With(prometheus.Labels{
					"metric_name": strings.Replace(data.Metric, "/", "-", 3),
					"service":     strings.Replace(data.Service, "/", "-", 3),
					"host":        data.Host}).Set(data.Value)
			}
			level.Info(druidLogger).Log("msg", "Successfully recieved data from druid emitter")
		}
	})
}

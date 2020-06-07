package listener

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// DruidEmittedData is the data structure which druid sends over HTTP
type DruidEmittedData struct {
	Feed       string    `json:"feed"`
	Timestamp  time.Time `json:"timestamp"`
	Service    string    `json:"service"`
	Host       string    `json:"host"`
	Version    string    `json:"version"`
	Metric     string    `json:"metric"`
	DataSource string    `json:"dataSource"`
	Value      float64   `json:"value"`
}

// DruidHTTPEndpoint is the endpoint to listen all druid metrics
func DruidHTTPEndpoint(gauge *prometheus.GaugeVec) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var druidData []DruidEmittedData
		if req.Method == "POST" {
			jsonDecoder := json.NewDecoder(req.Body)
			err := jsonDecoder.Decode(&druidData)
			if err != nil {
				logrus.Debugf("Error decoding JSON sent by druid: %v", err)
			}
			for _, data := range druidData {
				gauge.With(prometheus.Labels{
					"metric_name": strings.Replace(data.Metric, "/", "-", 3),
					"service":     strings.Replace(data.Service, "/", "-", 3),
					"host":        data.Host,
					"datasource":  data.DataSource,
				}).Set(data.Value)
			}
			logrus.Debugf("Successfully recieved data from druid emitter")
		}
	})
}

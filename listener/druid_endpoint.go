package listener

import (
	"druid-exporter/collector"
	"encoding/json"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// DruidHTTPEndpoint is the endpoint to listen all druid metrics
func DruidHTTPEndpoint(gauge *prometheus.GaugeVec) http.HandlerFunc {
	podNames := &sync.Map{}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var druidData []map[string]interface{}
		reqHeader, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if req.Method == "POST" && reqHeader == "application/json" {
			output, err := ioutil.ReadAll(req.Body)
			defer req.Body.Close()
			if err != nil {
				logrus.Errorf("Unable to read JSON reponse: %v", err)
			}
			err = json.Unmarshal(output, &druidData)
			if err != nil {
				logrus.Debugf("Error decoding JSON sent by druid: %v", err)
				logrus.Debugf("%v", druidData)
			}
			for _, data := range druidData {
				if data["dataSource"] != nil {
					metricName := fmt.Sprintf("%v", data["metric"])
					serviceName := fmt.Sprintf("%v", data["service"])
					host := fmt.Sprintf("%v", data["host"])
					value, _ := strconv.ParseFloat(fmt.Sprintf("%v", data["value"]), 64)

					var datasources []string
					if rawDatasource, ok := data["dataSource"].(string); ok {
						datasources = []string{rawDatasource}
					} else if rawDatasources, ok := data["dataSource"].([]interface{}); ok {
						datasources = make([]string, len(rawDatasources))
						for i, rawDatasource := range rawDatasources {
							datasources[i] = rawDatasource.(string)
						}
					}
					for _, datasource := range datasources {
						host := strings.Split(host, ":")[0]
						podName := collector.ToPodName(host)
						if podName == "" {
							if rawPodName, has := podNames.Load(host); has {
								podName = rawPodName.(string)
							}
						} else {
							podNames.Store(host, podName)
						}
						gauge.With(prometheus.Labels{
							"metric_name": strings.Replace(metricName, "/", "-", 3),
							"service":     strings.Replace(serviceName, "/", "-", 3),
							"datasource":  datasource,
							"pod":         podName,
						}).Set(value)
					}
				} else {
					metricName := fmt.Sprintf("%v", data["metric"])
					serviceName := fmt.Sprintf("%v", data["service"])
					host := fmt.Sprintf("%v", data["host"])
					value, _ := strconv.ParseFloat(fmt.Sprintf("%v", data["value"]), 64)
					gauge.With(prometheus.Labels{
						"metric_name": strings.Replace(metricName, "/", "-", 3),
						"service":     strings.Replace(serviceName, "/", "-", 3),
						"datasource":  host,
						"pod":         collector.ToPodName(strings.Split(host, ":")[0]),
					}).Set(value)
				}
			}
			logrus.Debugf("Successfully collected data from druid emitter")
		}
	})
}

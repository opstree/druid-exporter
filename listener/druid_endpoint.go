package listener

import (
	"druid-exporter/collector"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// DruidHTTPEndpoint is the endpoint to listen all druid metrics
func DruidHTTPEndpoint(gauge *prometheus.GaugeVec, dnsCache *cache.Cache) http.HandlerFunc {
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
				metric := data["metric"].(string)
				service := data["service"].(string)
				hostname := data["host"].(string)
				value, _ := data["value"].(float64)

				// Reverse DNS Lookup
				// Mutates dnsCache
				hostValue := strings.Split(hostname, ":")[0]
				dnsLookupValue := collector.ReverseDNSLookup(hostValue, dnsCache)

				host := strings.Replace(hostname, hostValue, dnsLookupValue, 1) // Adding back port

				if datasource, ok := data["dataSource"]; ok {
					if arrDatasource, ok := datasource.([]interface{}); ok { // Array datasource
						for _, entryDatasource := range arrDatasource {
							gauge.With(prometheus.Labels{
								"metric_name": strings.Replace(metric, "/", "-", 3),
								"service":     strings.Replace(service, "/", "-", 3),
								"datasource":  entryDatasource.(string),
								"host":        host,
							}).Set(value)
						}
					} else { // Single datasource
						gauge.With(prometheus.Labels{
							"metric_name": strings.Replace(metric, "/", "-", 3),
							"service":     strings.Replace(service, "/", "-", 3),
							"datasource":  datasource.(string),
							"host":        host,
						}).Set(value)
					}
				} else { // Missing datasource case
					gauge.With(prometheus.Labels{
						"metric_name": strings.Replace(metric, "/", "-", 3),
						"service":     strings.Replace(service, "/", "-", 3),
						"datasource":  "",
						"host":        host,
					}).Set(value)
				}
			}
			logrus.Infof("Successfully collected data from druid emitter, %s", druidData[0]["service"].(string))
		}
	})
}

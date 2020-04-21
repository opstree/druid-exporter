package collector

import (
	"encoding/json"
	"druid-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus        *prometheus.Desc
}

type DataSources struct {
	DataSource []string
}

// GetDruidMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	return utils.GetDruidHealth("http://52.172.156.84:8081/status/health")
}

// GetDruidDatasource returns the datasources of druid
func GetDruidDatasource() DataSources{
	respData, _ := GetDruidResponse("http://52.172.156.84:8081/druid/coordinator/v1/metadata/datasources")
	generic, _ := json.Unmarshal(respData, &DataSources)
	return generic
}

// Describe will associate the value for druid exporter
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DruidHealthStatus
}

// Collector return the defined metrics
func Collector() *MetricCollector{
	return &MetricCollector{
		DruidHealthStatus: prometheus.NewDesc("druid_health_status",
			"Health of Druid, 1 is healthy 0 is not",
			nil, prometheus.Labels{
				"druid": "health",
			},
		),
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.DruidHealthStatus, prometheus.CounterValue, GetDruidHealthMetrics())
}

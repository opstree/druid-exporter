package collector

import (
	"druid-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus        *prometheus.Desc
}

// GetDruidMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() {
	return utils.GetDruidHealth("http://52.172.156.84:8081/status/health")
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

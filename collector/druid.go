package collector

import (
	"druid-exporter/utils"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"os"
)

// GetDruidHealthMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	druidHealthURL := os.Getenv("DRUID_URL") + healthURL
	log.Info().Str("Query Type", "Health").Msg("Successfully made a request to get healthcheck")
	return utils.GetHealth(druidHealthURL)
}

// GetDruidSegmentData returns the datasources of druid
func GetDruidSegmentData() SegementInterface {
	druidSegmentURL := os.Getenv("DRUID_URL") + segmentDataURL
	responseData, err := utils.GetResponse(druidSegmentURL, "Segment")
	if err != nil {
		log.Error().Str("Query Type", "Segment").Msg("Error while making request on provided URL")
	}
	log.Info().Str("Query Type", "Segment").Msg("Successfully executed the request to get segment data")
	var metric SegementInterface
	json.Unmarshal(responseData, &metric)
	return metric
}

// GetDruidData return all the tasks and its state
func GetDruidData(pathURL string) []map[string]interface{} {
	druidURL := os.Getenv("DRUID_URL") + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		log.Error().Str("Query Type", pathURL).Msg("Error while making request on provided URL")
	}
	log.Info().Str("Query Type", pathURL).Msg("Successfully executed the request")
	var metric []map[string]interface{}
	json.Unmarshal(responseData, &metric)
	return metric
}

// Describe will associate the value for druid exporter
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DruidHealthStatus
	ch <- collector.DataSourceCount
	ch <- collector.DruidTasks
	ch <- collector.DruidSupervisors
	ch <- collector.DruidSegmentCount
	ch <- collector.DruidSegmentSize
	ch <- collector.DruidSegmentReplicateSize
	ch <- collector.DruidEmittedData
}

// Collector return the defined metrics
func Collector() *MetricCollector {
	return &MetricCollector{
		DruidHealthStatus: prometheus.NewDesc("druid_health_status",
			"Health of Druid, 1 is healthy 0 is not",
			nil, prometheus.Labels{
				"druid": "health",
			},
		),
		DataSourceCount: prometheus.NewDesc("druid_datasource",
			"Datasources present",
			[]string{"datasource"}, nil,
		),
		DruidTasks: prometheus.NewDesc("druid_tasks",
			"Druid tasks status",
			[]string{"datasource", "index_group_id", "task_status", "created_time"}, nil,
		),
		DruidSupervisors: prometheus.NewDesc("druid_supervisors",
			"Druid supervisors status",
			[]string{"supervisor_name", "healthy", "state"}, nil,
		),
		DruidSegmentCount: prometheus.NewDesc("druid_segment_count",
			"Druid segment count",
			[]string{"datasource_name"}, nil,
		),
		DruidSegmentSize: prometheus.NewDesc("druid_segment_size",
			"Druid segment size",
			[]string{"datasource_name"}, nil,
		),
		DruidSegmentReplicateSize: prometheus.NewDesc("druid_segment_replicated_size",
			"Druid segment replicated size",
			[]string{"datasource_name"}, nil,
		),
		DruidEmittedData: prometheus.NewDesc("druid_http_emitter_metrics",
			"Druid emitted data",
			[]string{"metric_name", "service_name"}, nil,
		),
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.DruidHealthStatus, prometheus.CounterValue, GetDruidHealthMetrics())
	for _, data := range GetDruidSegmentData() {
		ch <- prometheus.MustNewConstMetric(collector.DataSourceCount, prometheus.GaugeValue, float64(1), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentCount, prometheus.GaugeValue, float64(data.Properties.Segments.Count), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentSize, prometheus.GaugeValue, float64(data.Properties.Segments.Size), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentReplicateSize, prometheus.GaugeValue, float64(data.Properties.Segments.ReplicatedSize), data.Name)
	}
	for _, data := range GetDruidData(tasksURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidTasks, prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["dataSource"]), fmt.Sprintf("%v", data["groupId"]), fmt.Sprintf("%v", data["status"]), fmt.Sprintf("%v", data["createdTime"]))
	}
	for _, data := range GetDruidData(supervisorURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors, prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["id"]), fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}
}

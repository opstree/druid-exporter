package collector

import (
	"fmt"
	"encoding/json"
	"druid-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus        *prometheus.Desc
	DataSourceCount          *prometheus.Desc
	DruidTasks               *prometheus.Desc
	DruidSupervisors         *prometheus.Desc
	DruidSegmentCount        *prometheus.Desc
}

// SegementInterface is the interface for parsing segments data
type SegementInterface []struct {
	Name       string `json:"name"`
	Properties struct {
		Tiers struct {
			DefaultTier struct {
				Size           int64 `json:"size"`
				ReplicatedSize int64 `json:"replicatedSize"`
				SegmentCount   int   `json:"segmentCount"`
			} `json:"_default_tier"`
		} `json:"tiers"`
		Segments struct {
			MaxTime        time.Time `json:"maxTime"`
			Size           int64     `json:"size"`
			MinTime        time.Time `json:"minTime"`
			Count          int       `json:"count"`
			ReplicatedSize int64     `json:"replicatedSize"`
		} `json:"segments"`
	} `json:"properties"`
}

// GetDruidMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	return utils.GetDruidHealth("http://52.172.156.84:8081/status/health")
}

// GetDruidDatasource returns the datasources of druid
func GetDruidSegmentData() []string{
	respData, _ := utils.GetDruidResponse("http://52.172.156.84:8081/druid/coordinator/v1/datasources?simple")

	var metric SegementInterface
	json.Unmarshal(respData, &metric)
	return metric
}

// GetDruidTasks() return all the tasks and its state
func GetDruidTasks() []map[string]interface{} {
	respData, _ := utils.GetDruidResponse("http://52.172.156.84:8081/druid/indexer/v1/tasks")
	var metric []map[string]interface{}
	json.Unmarshal(respData, &metric)
	return metric
}

// GetDruidTasks() return all the tasks and its state
func GetDruidSupervisors() []map[string]interface{} {
	respData, _ := utils.GetDruidResponse("http://52.172.156.84:8081/druid/indexer/v1/supervisor?full")
	var metric []map[string]interface{}
	json.Unmarshal(respData, &metric)
	return metric
}


// Describe will associate the value for druid exporter
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DruidHealthStatus
	ch <- collector.DataSourceCount
	ch <- collector.DruidTasks
	ch <- collector.DruidSupervisors
	ch <- collector.DruidSegmentCount
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
		DruidSegmentCount: prometheus.NewDesc("druid_segement_count",
			"Druid segment count",
			[]string{"name"}, nil,
		),
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.DruidHealthStatus, prometheus.CounterValue, GetDruidHealthMetrics())
	for _, data := range SegementInterface() {
		ch <- prometheus.MustNewConstMetric(collector.DataSourceCount, prometheus.GaugeValue, float64(1), data.Name)
	}
	for _, data := range GetDruidTasks() {
		ch <- prometheus.MustNewConstMetric(collector.DruidTasks, prometheus.GaugeValue, float64(1), fmt.Sprintf("%v",data["dataSource"]), fmt.Sprintf("%v", data["groupId"]), fmt.Sprintf("%v", data["status"]), fmt.Sprintf("%v", data["createdTime"]))
	}
	for _, data := range GetDruidSupervisors() {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors, prometheus.GaugeValue, float64(1), fmt.Sprintf("%v",data["id"]), fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}
	for _, data := range SegementInterface() {
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentCount, prometheus.GaugeValue, float64(data.Properties.Segments.Count), data.Name)
	}
}

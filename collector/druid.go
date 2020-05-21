package collector

import (
	"druid-exporter/logger"
	"druid-exporter/utils"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
	"strconv"
)

var (
	druid       = kingpin.Flag("druid.uri", "URL of druid router or coordinator").Default("http://druid.opstreelabs.in").OverrideDefaultFromEnvar("DRUID_URL").Short('d').String()
	druidLogger = logger.GetLoggerInterface()
)

// GetDruidHealthMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	kingpin.Parse()
	druidHealthURL := *druid + healthURL
	level.Info(druidLogger).Log("msg", "Successfully retrieved the data for druid healthcheck")
	return utils.GetHealth(druidHealthURL)
}

// GetDruidSegmentData returns the datasources of druid
func GetDruidSegmentData() SegementInterface {
	kingpin.Parse()
	druidSegmentURL := *druid + segmentDataURL
	responseData, err := utils.GetResponse(druidSegmentURL, "Segment")
	if err != nil {
		level.Error(druidLogger).Log("msg", "Cannot retrieve data for druid segments", "err", err)
	}
	level.Info(druidLogger).Log("msg", "Successfully retrieved the data for druid segment")
	var metric SegementInterface
	json.Unmarshal(responseData, &metric)
	return metric
}

// GetDruidData return all the tasks and its state
func GetDruidData(pathURL string) []map[string]interface{} {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		level.Error(druidLogger).Log("msg", "Cannot retrieve data for druid's supervisors tasks", "err", err)
	}
	level.Info(druidLogger).Log("msg", "Successfully retrieved the data for druid's supervisors tasks")
	var metric []map[string]interface{}
	json.Unmarshal(responseData, &metric)
	return metric
}

// Describe will associate the value for druid exporter
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DruidHealthStatus
	ch <- collector.DataSourceCount
	ch <- collector.DruidSupervisors
	ch <- collector.DruidSegmentCount
	ch <- collector.DruidSegmentSize
	ch <- collector.DruidSegmentReplicateSize
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
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.DruidHealthStatus,
		prometheus.CounterValue, GetDruidHealthMetrics())
	for _, data := range GetDruidSegmentData() {
		ch <- prometheus.MustNewConstMetric(collector.DataSourceCount,
			prometheus.GaugeValue, float64(1), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentCount,
			prometheus.GaugeValue, float64(data.Properties.Segments.Count), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentSize,
			prometheus.GaugeValue, float64(data.Properties.Segments.Size), data.Name)
		ch <- prometheus.MustNewConstMetric(collector.DruidSegmentReplicateSize,
			prometheus.GaugeValue, float64(data.Properties.Segments.ReplicatedSize), data.Name)
	}
	for _, data := range GetDruidData(supervisorURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors,
			prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["id"]),
			fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}
}

// CollectTaskMetrics will capture the druid tasks metrics
func CollectTaskMetrics(gauge *prometheus.GaugeVec) {
	for _, data := range GetDruidData(tasksURL) {
		value, err := strconv.ParseFloat(fmt.Sprintf("%v", data["duration"]), 64)

		if err != nil {
			level.Debug(druidLogger).Log("msg", "Unable to parse the duration value", "err", err)
		}
		gauge.With(prometheus.Labels{
			"datasource_name": fmt.Sprintf("%v", data["dataSource"]),
			"group_id":        fmt.Sprintf("%v", data["groupId"]),
			"task_status":     fmt.Sprintf("%v", data["status"]),
			"created_time":    fmt.Sprintf("%v", data["createdTime"]),
		}).Set(value)
	}
}

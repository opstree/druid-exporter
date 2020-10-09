package collector

import (
	"druid-exporter/utils"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	druid = kingpin.Flag(
		"druid.uri",
		"URL of druid router or coordinator, EnvVar - DRUID_URL",
	).Default("http://druid.opstreelabs.in").OverrideDefaultFromEnvar("DRUID_URL").Short('d').String()
)

// GetDruidHealthMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	kingpin.Parse()
	druidHealthURL := *druid + healthURL
	logrus.Debugf("Successfully collected the data for druid healthcheck")
	return utils.GetHealth(druidHealthURL)
}

// GetDruidSegmentData returns the datasources of druid
func GetDruidSegmentData() SegementInterface {
	kingpin.Parse()
	druidSegmentURL := *druid + segmentDataURL
	responseData, err := utils.GetResponse(druidSegmentURL, "Segment")
	if err != nil {
		logrus.Errorf("Cannot collect data for druid segments: %v", err)
		return nil
	}
	logrus.Debugf("Successfully collected the data for druid segment")
	var metric SegementInterface
	err = json.Unmarshal(responseData, &metric)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Druid segment's metric data, %v", metric)
	return metric
}

// GetDruidData return all the tasks and its state
func GetDruidData(pathURL string) []map[string]interface{} {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		logrus.Errorf("Cannot collect data for druid's supervisors: %v", err)
		return nil
	}
	logrus.Debugf("Successfully collected the data for druid's supervisors")
	var metric []map[string]interface{}
	err = json.Unmarshal(responseData, &metric)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Druid supervisor's metric data, %v", metric)
	return metric
}

// GetDruidTasksData return all the tasks and its state
func GetDruidTasksData(pathURL string) TasksInterface {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		logrus.Errorf("Cannot retrieve data for druid's tasks: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved the data for druid's tasks")
	var metric TasksInterface
	err = json.Unmarshal(responseData, &metric)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Druid tasks's metric data, %v", metric)
	return metric
}

// getDruidWorkersData return all the workers and its state
func getDruidWorkersData(pathURL string) []worker {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		logrus.Errorf("Cannot retrieve data for druid's workers: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved the data for druid's workers")
	var workers []worker
	err = json.Unmarshal(responseData, &workers)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Druid workers's metric data, %v", workers)

	return workers
}

// Describe will associate the value for druid exporter
func (collector *MetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.DruidHealthStatus
	ch <- collector.DataSourceCount
	ch <- collector.DruidSupervisors
	ch <- collector.DruidSegmentCount
	ch <- collector.DruidSegmentSize
	ch <- collector.DruidWorkers
	ch <- collector.DruidTasks
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
		DruidWorkers: prometheus.NewDesc("druid_workers_capacity_used",
			"Druid workers capacity used",
			[]string{"pod", "version"}, nil,
		),
		DruidTasks: prometheus.NewDesc("druid_tasks_duration",
			"Druid tasks duration and state",
			[]string{"pod", "datasource", "task_id", "groupd_id", "task_status", "created_time"}, nil,
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

	workers := getDruidWorkersData(workersURL)

	for _, worker := range workers {
		ch <- prometheus.MustNewConstMetric(collector.DruidWorkers,
			prometheus.GaugeValue, float64(worker.CurrCapacityUsed), worker.hostname(), worker.Worker.Version)
	}

	for _, data := range GetDruidTasksData(tasksURL) {
		hostname := ""
		for _, worker := range workers {
			for _, task := range worker.RunningTasks {
				if task == data.ID {
					hostname = worker.hostname()
					break
				}
			}
			if hostname != "" {
				break
			}
		}
		if hostname == "" {
			if len(workers) != 0 {
				hostname = workers[rand.Intn(len(workers))].hostname()
			}
		}
		ch <- prometheus.MustNewConstMetric(collector.DruidTasks,
			prometheus.GaugeValue, data.Duration, hostname, data.DataSource, data.ID, data.GroupID, data.Status, data.CreatedTime)
	}

	for _, data := range GetDruidData(supervisorURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors,
			prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["id"]),
			fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}
}

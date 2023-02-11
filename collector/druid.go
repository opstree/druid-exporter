package collector

import (
	"druid-exporter/utils"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"strings"
)

var (
	druid = kingpin.Flag(
		"druid.uri",
		"URL of druid router or coordinator, EnvVar - DRUID_URL",
	).Default("http://druid.opstreelabs.in").OverrideDefaultFromEnvar("DRUID_URL").Short('d').String()
)

const (
	LABEL_NAME_LENGTH_LIMIT = 90
	EXCEPTION               = "Exception"
	FAILED                  = "FAILED"
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

// GetDruidDataSourcesTotalRows returns the amount of rows in each datasource
func GetDruidDataSourcesTotalRows(pathURL string) DataSourcesTotalRows {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetSQLResponse(druidURL, totalRowsSQL)
	if err != nil {
		logrus.Errorf("Cannot retrieve data for druid's datasources rows: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved the data for druid's datasources rows")
	var datasources DataSourcesTotalRows
	err = json.Unmarshal(responseData, &datasources)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Druid datasources total rows, %v", datasources)
	return datasources
}

// GetDruidTasksStatusCount returns count of different tasks by status
func GetDruidTasksStatusCount(pathURL string) TaskStatusMetric {
	kingpin.Parse()
	druidURL := *druid + pathURL
	responseData, err := utils.GetResponse(druidURL, pathURL)
	if err != nil {
		logrus.Errorf("Cannot retrieve data for druid's workers: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved the data for druid task: %v", pathURL)
	var taskCount TaskStatusMetric
	err = json.Unmarshal(responseData, &taskCount)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}
	logrus.Debugf("Successfully collected tasks status count: %v", pathURL)
	return taskCount
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
	ch <- collector.DruidRunningTasks
	ch <- collector.DruidWaitingTasks
	ch <- collector.DruidCompletedTasks
	ch <- collector.DruidPendingTasks
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
			[]string{"pod", "version", "ip"}, nil,
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
		DruidDataSourcesTotalRows: prometheus.NewDesc("druid_datasource_total_rows",
			"Number of rows in a datasource",
			[]string{"datasource_name", "source"}, nil),
		DruidRunningTasks: prometheus.NewDesc("druid_running_tasks",
			"Druid running tasks count",
			nil, nil,
		),
		DruidWaitingTasks: prometheus.NewDesc("druid_waiting_tasks",
			"Druid waiting tasks count",
			nil, nil,
		),
		DruidCompletedTasks: prometheus.NewDesc("druid_completed_tasks",
			"Druid completed tasks count",
			nil, nil,
		),
		DruidPendingTasks: prometheus.NewDesc("druid_pending_tasks",
			"Druid pending tasks count",
			nil, nil,
		),
		DruidFailedTasks: prometheus.NewDesc("druid_failed_tasks",
			"Druid failed tasks count",
			nil, nil,
		),
		DruidTaskCapacity: prometheus.NewDesc("druid_task_capacity",
			"Druid task capacity",
			nil, nil,
		),
		DruidTaskErrors: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "druid_task_errors",
				Help: "Druid task errors",
			},
			[]string{"error_msg"},
		),
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {
	var errorMap map[string]float64 = make(map[string]float64)
	var failedTaskCount float64 = 0

	ch <- prometheus.MustNewConstMetric(collector.DruidHealthStatus,
		prometheus.CounterValue, GetDruidHealthMetrics())
	for _, data := range GetDruidSegmentData() {
		ch <- prometheus.MustNewConstMetric(collector.DataSourceCount,
			prometheus.GaugeValue, float64(1), data.Name)
		if data.Properties.Segments.Count != 0 {
			ch <- prometheus.MustNewConstMetric(collector.DruidSegmentCount,
				prometheus.GaugeValue, float64(data.Properties.Segments.Count), data.Name)
		}
		if data.Properties.Segments.Size != 0 {
			ch <- prometheus.MustNewConstMetric(collector.DruidSegmentSize,
				prometheus.GaugeValue, float64(data.Properties.Segments.Size), data.Name)
		}
		if data.Properties.Segments.ReplicatedSize != 0 {
			ch <- prometheus.MustNewConstMetric(collector.DruidSegmentReplicateSize,
				prometheus.GaugeValue, float64(data.Properties.Segments.ReplicatedSize), data.Name)
		}
	}

	ch <- prometheus.MustNewConstMetric(collector.DruidRunningTasks,
		prometheus.GaugeValue, float64(len(GetDruidTasksStatusCount(runningTask))))
	ch <- prometheus.MustNewConstMetric(collector.DruidWaitingTasks,
		prometheus.GaugeValue, float64(len(GetDruidTasksStatusCount(waitingTask))))
	ch <- prometheus.MustNewConstMetric(collector.DruidCompletedTasks,
		prometheus.GaugeValue, float64(len(GetDruidTasksStatusCount(completedTask))))
	ch <- prometheus.MustNewConstMetric(collector.DruidPendingTasks,
		prometheus.GaugeValue, float64(len(GetDruidTasksStatusCount(pendingTask))))

	workers := getDruidWorkersData(workersURL)

	taskCapacity := 0
	for _, worker := range workers {
		taskCapacity += worker.Worker.Capacity
		ch <- prometheus.MustNewConstMetric(collector.DruidWorkers,
			prometheus.GaugeValue, float64(worker.CurrCapacityUsed), worker.hostname(), worker.Worker.Version, worker.Worker.IP)
	}

	ch <- prometheus.MustNewConstMetric(collector.DruidTaskCapacity, prometheus.GaugeValue, float64(taskCapacity))

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
		if len(data.ErrorMsg) > 0 {
			error_label := errorMsgLabel(data.ErrorMsg)
			errorMap[error_label]++
			if data.Status == FAILED {
				failedTaskCount++
			}
		}
		ch <- prometheus.MustNewConstMetric(collector.DruidTasks,
			prometheus.GaugeValue, data.Duration, hostname, data.DataSource, data.ID, data.GroupID, data.Status, data.CreatedTime)
	}

	for key, val := range errorMap {
		m := collector.DruidTaskErrors.With(prometheus.Labels{"error_msg": key})
		m.Set(val)
		ch <- m
	}

	ch <- prometheus.MustNewConstMetric(collector.DruidFailedTasks, prometheus.GaugeValue, failedTaskCount)

	for _, data := range GetDruidData(supervisorURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors,
			prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["id"]),
			fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}

	for _, data := range GetDruidDataSourcesTotalRows(sqlURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidDataSourcesTotalRows, prometheus.GaugeValue, float64(data.TotalRows), data.Datasource, data.Source)
	}
}

func errorMsgLabel(errorMsg string) string {
	var label string
	i := strings.Index(errorMsg, EXCEPTION)
	if i > 0 {
		label = errorMsg[0 : i+9]
	} else {
		label = strings.ReplaceAll(errorMsg, " ", "")
	}

	if len(label) > LABEL_NAME_LENGTH_LIMIT { // truncate label
		label = label[:LABEL_NAME_LENGTH_LIMIT]
	}

	fmt.Println(label)
	return label
}

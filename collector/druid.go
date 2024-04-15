package collector

import (
	"druid-exporter/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	druid = kingpin.Flag(
		"druid.uri",
		"URL of druid router or coordinator, EnvVar - DRUID_URL",
	).Default("http://localhost:8888").OverrideDefaultFromEnvar("DRUID_URL").Short('d').String()

	maxCompletedTasks = kingpin.Flag(
		"maxCompletedTasks",
		"Max Results of completed Tasks (Default: 50)",
	).Default("50").OverrideDefaultFromEnvar("MAX_COMPLETED_TASKS").String()
)

// GetDruidHealthMetrics returns the set of metrics for druid
func GetDruidHealthMetrics() float64 {
	kingpin.Parse()
	druidHealthURL := *druid + healthURL
	logrus.Debugf("Successfully collected the data for druid healthcheck")
	return utils.GetHealth(druidHealthURL)
}

// GetDruidSegmentData returns the datasources of druid
func GetDruidSegmentData() SegmentInterface {
	kingpin.Parse()
	druidSegmentURL := *druid + segmentDataURL
	responseData, err := utils.GetResponse(druidSegmentURL, "Segment")
	if err != nil {
		logrus.Errorf("Cannot collect data for druid segments: %v", err)
		return nil
	}
	logrus.Debugf("Successfully collected the data for druid segment")
	var metric SegmentInterface
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
	max := *maxCompletedTasks
	druidURL := *druid + pathURL
	pathURL = pathURL + fmt.Sprintf("&max=%s", max)
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

// GetDruidHistoricalFreespace returns the sum of freespace of historical nodes
func GetDruidHistoricalFreespace(pathURL string, dnsCache *cache.Cache) DruidHistoricalFreeSpace {
	kingpin.Parse()
	druidURL := *druid + pathURL
	const query = `
	SELECT
		"server" AS "host",
		"server_type" AS "server_type",
		"host" as "ip",
		"max_size" - "curr_size" as "free_size"
	FROM sys.servers
	WHERE 
	"server_type" = 'historical'
	`
	responseData, err := utils.GetSQLResponse(druidURL, query)
	if err != nil {
		logrus.Errorf("Cannot retrieve data from druid about freespace: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved from druid about total freespace")
	var freeSpace DruidHistoricalFreeSpace
	err = json.Unmarshal(responseData, &freeSpace)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}

	for i, _ := range freeSpace {
		freeSpace[i].POD = utils.ReverseDNSLookup(freeSpace[i].IP, dnsCache)
	}

	logrus.Debugf("Druid Historical total free space, %v", freeSpace)
	return freeSpace
}

func GetDruidHistoricalUsagePercent(pathURL string, dnsCache *cache.Cache) DruidHistoricalUsagePercent {
	kingpin.Parse()
	druidURL := *druid + pathURL
	const query = `
	SELECT
  "server" AS "host",
  "server_type" AS "server_type",
  "host" AS "ip",
  CAST ("curr_size" AS FLOAT) / CAST ("max_size" AS FLOAT) * 100.0 AS "usage_percent"
  FROM sys.servers
  WHERE 
  "server_type" = 'historical'
	`
	responseData, err := utils.GetSQLResponse(druidURL, query)
	if err != nil {
		logrus.Errorf("Cannot retrieve data from druid about usage: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved from druid about usage")
	var usage DruidHistoricalUsagePercent
	err = json.Unmarshal(responseData, &usage)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}

	for i, _ := range usage {
		usage[i].POD = utils.ReverseDNSLookup(usage[i].IP, dnsCache)
	}

	logrus.Debugf("Druid Historical total free space, %v", usage)
	return usage
}

func GetDruidHistoricalUsageAbsolute(pathURL string, dnsCache *cache.Cache) DruidHistoricalUsageAbsolute {
	kingpin.Parse()
	druidURL := *druid + pathURL
	const query = `
	SELECT
  "server" AS "host",
  "server_type" AS "server_type",
  "host" AS "ip",
	"curr_size" AS "usage_absolute"
  FROM sys.servers
  WHERE 
  "server_type" = 'historical'
	`
	responseData, err := utils.GetSQLResponse(druidURL, query)
	if err != nil {
		logrus.Errorf("Cannot retrieve data from druid about usage: %v", err)
		return nil
	}
	logrus.Debugf("Successfully retrieved from druid about usage")
	var usage DruidHistoricalUsageAbsolute
	err = json.Unmarshal(responseData, &usage)
	if err != nil {
		logrus.Errorf("Cannot parse JSON data: %v", err)
		return nil
	}

	for i, _ := range usage {
		usage[i].POD = utils.ReverseDNSLookup(usage[i].IP, dnsCache)
	}

	logrus.Debugf("Druid Historical usage, %v", usage)
	return usage
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

		DruidHistoricalFreeSpace: prometheus.NewDesc("druid_historical_free_space",
			"Freespace of all historicals per node",
			[]string{"host", "server_type", "ip", "pod"}, nil),

		DruidHistoricalUsagePercent: prometheus.NewDesc("druid_historical_usage_percent",
			"Usage of all historicals per node in Percent",
			[]string{"host", "server_type", "ip", "pod"}, nil),

		DruidHistoricalUsageAbsolute: prometheus.NewDesc("druid_historical_usage_absolute",
			"Absolute Usage of all historicals per node",
			[]string{"host", "server_type", "ip", "pod"}, nil),

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
		DruidTaskCapacity: prometheus.NewDesc("druid_task_capacity",
			"Druid task capacity",
			nil, nil,
		),
	}
}

// Collect will collect all the metrics
func (collector *MetricCollector) Collect(ch chan<- prometheus.Metric) {

	dnsCache := cache.New(5*time.Minute, 10*time.Minute)

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
		ch <- prometheus.MustNewConstMetric(collector.DruidTasks,
			prometheus.GaugeValue, data.Duration, hostname, data.DataSource, data.ID, data.GroupID, data.Status, data.CreatedTime)
	}

	for _, data := range GetDruidData(supervisorURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidSupervisors,
			prometheus.GaugeValue, float64(1), fmt.Sprintf("%v", data["id"]),
			fmt.Sprintf("%v", data["healthy"]), fmt.Sprintf("%v", data["detailedState"]))
	}

	for _, data := range GetDruidDataSourcesTotalRows(sqlURL) {
		ch <- prometheus.MustNewConstMetric(collector.DruidDataSourcesTotalRows, prometheus.GaugeValue, float64(data.TotalRows), data.Datasource, data.Source)
	}

	for _, data := range GetDruidHistoricalFreespace(sqlURL, dnsCache) {
		ch <- prometheus.MustNewConstMetric(collector.DruidHistoricalFreeSpace, prometheus.GaugeValue, float64(data.FreeSize), data.Host, data.ServerType, data.IP, data.POD)
	}

	for _, data := range GetDruidHistoricalUsagePercent(sqlURL, dnsCache) {
		ch <- prometheus.MustNewConstMetric(collector.DruidHistoricalUsagePercent, prometheus.GaugeValue, float64(data.UsagePercent), data.Host, data.ServerType, data.IP, data.POD)
	}
	for _, data := range GetDruidHistoricalUsageAbsolute(sqlURL, dnsCache) {
		ch <- prometheus.MustNewConstMetric(collector.DruidHistoricalUsageAbsolute, prometheus.GaugeValue, float64(data.UsageAbsolute), data.Host, data.ServerType, data.IP, data.POD)
	}
}

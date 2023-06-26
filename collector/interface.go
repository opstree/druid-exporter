package collector

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	healthURL           = "/status/health"
	segmentDataURL      = "/druid/coordinator/v1/datasources?simple"
	tasksURL            = "/druid/indexer/v1/tasks"
	workersURL          = "/druid/indexer/v1/workers"
	supervisorURL       = "/druid/indexer/v1/supervisor?full"
	sqlURL              = "/druid/v2/sql"
	pendingTask         = "/druid/indexer/v1/pendingTasks"
	runningTask         = "/druid/indexer/v1/runningTasks"
	waitingTask         = "/druid/indexer/v1/waitingTasks"
	completedTask       = "/druid/indexer/v1/completeTasks"
	compactionStatusURL = "/druid/coordinator/v1/compaction/status"
)

const totalRowsSQL = `select SEG.datasource, SUP.source,
SUM(SEG."num_rows") FILTER (WHERE (SEG.is_published = 1 AND SEG.is_overshadowed = 0) OR SEG.is_realtime = 1) AS total_rows
from sys.segments SEG
inner join sys.supervisors SUP ON SEG.datasource=SUP.supervisor_id
group by SEG.datasource, SUP.source`

const segmentData = `SELECT datasource,
    COUNT(*) FILTER (WHERE (is_published = 1 AND is_overshadowed = 0) OR is_realtime = 1) AS num_segments,
	COUNT(*) FILTER (WHERE is_published = 1 AND is_overshadowed = 0 AND is_available = 0) AS num_segments_to_load,
	COUNT(*) FILTER (WHERE is_available = 1 AND NOT ((is_published = 1 AND is_overshadowed = 0) OR is_realtime = 1)) AS num_segments_to_drop,
	CASE WHEN
	 SUM("num_rows") FILTER (WHERE is_published = 1 AND is_overshadowed = 0) <> 0
	THEN
	 (SUM("size") FILTER (WHERE is_published = 1 AND is_overshadowed = 0) / SUM("num_rows") FILTER (WHERE is_published = 1 AND is_overshadowed = 0))
	ELSE 0
	END AS avg_row_size
FROM sys.segments
GROUP BY 1 ORDER BY 1`

// MetricCollector includes the list of metrics
type MetricCollector struct {
	DruidHealthStatus                 *prometheus.Desc
	DataSourceCount                   *prometheus.Desc
	DruidWorkers                      *prometheus.Desc
	DruidTasks                        *prometheus.Desc
	DruidSupervisors                  *prometheus.Desc
	DruidSegmentCount                 *prometheus.Desc
	DruidSegmentSize                  *prometheus.Desc
	DruidSegmentReplicateSize         *prometheus.Desc
	DruidDataSourcesTotalRows         *prometheus.Desc
	DruidDataSourcesAverageRowSize    *prometheus.Desc
	DruidDataSourcesNumSegments       *prometheus.Desc
	DruidDataSourcesNumSegmentsToLoad *prometheus.Desc
	DruidDataSourcesNumSegmentsToDrop *prometheus.Desc
	DruidRunningTasks                 *prometheus.Desc
	DruidWaitingTasks                 *prometheus.Desc
	DruidCompletedTasks               *prometheus.Desc
	DruidPendingTasks                 *prometheus.Desc
	DruidPendingIngestTasks           *prometheus.Desc
	DruidFailedTasks                  *prometheus.Desc
	DruidTaskCapacity                 *prometheus.Desc
	DruidTaskErrors                   *prometheus.GaugeVec
	DruidBytesCompaction              *prometheus.Desc
	DruidSegmentCountCompaction       *prometheus.Desc
	DruidIntervalCountCompaction      *prometheus.Desc
}

// DataSourcesSegmentData shows average row size from each datasource
type DataSourcesSegmentData []struct {
	Datasource        string `json:"datasource"`
	AvgRowSize        int64  `json:"avg_row_size"`
	NumSegments       int64  `json:"num_segments"`
	NumSegmentsToLoad int64  `json:"num_segments_to_load"`
	NumSegmentsToDrop int64  `json:"num_segments_to_drop"`
}

// DataSourcesTotalRows shows total rows from each datasource
type DataSourcesTotalRows []struct {
	Datasource string `json:"datasource"`
	Source     string `json:"source"`
	TotalRows  int64  `json:"total_rows"`
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

// TasksInterface is the interface for parsing druid tasks data
type TasksInterface []struct {
	ID               string  `json:"id"`
	GroupID          string  `json:"groupId"`
	Type             string  `json:"type"`
	CreatedTime      string  `json:"createdTime"`
	StatusCode       string  `json:"statusCode"`
	Status           string  `json:"status"`
	RunnerStatusCode string  `json:"runnerStatusCode"`
	Duration         float64 `json:"duration"`
	DataSource       string  `json:"dataSource"`
	ErrorMsg         string  `json:"errorMsg,omitempty"`
}

// TaskStatusMetric is the interface for tasks status
type TaskStatusMetric []struct {
	NameDataSource string `json:"dataSource"`
	StatusCode     string `json:"statusCode"`
	Type           string `json:"type"`
}

type worker struct {
	Worker struct {
		Host     string `json:"host"`
		Version  string `json:"version"`
		IP       string `json:"ip"`
		Capacity int    `json:"capacity"`
	}
	CurrCapacityUsed int      `json:"currCapacityUsed"`
	RunningTasks     []string `json:"runningTasks"`
}

// CompactionInterface is the interface for parsing compaction status data
type CompactionStatusInterface []struct {
	DataSource                      string  `json:"dataSource"`
	ScheduleStatus                  string  `json:"scheduleStatus"`
	BytesAwaitingCompaction         float64 `json:"bytesAwaitingCompaction"`
	BytesCompacted                  float64 `json:"bytesCompacted"`
	BytesSkipped                    float64 `json:"bytesSkipped"`
	SegmentCountAwaitingCompaction  float64 `json:"segmentCountAwaitingCompaction"`
	SegmentCountCompacted           float64 `json:"segmentCountCompacted"`
	SegmentCountSkipped             float64 `json:"segmentCountSkipped"`
	IntervalCountAwaitingCompaction float64 `json:"intervalCountAwaitingCompaction"`
	IntervalCountCompacted          float64 `json:"intervalCountCompacted"`
	IntervalCountSkipped            float64 `json:"intervalCountSkipped"`
}

func (w worker) hostname() string {
	return strings.Split(w.Worker.IP, ".")[0]
}

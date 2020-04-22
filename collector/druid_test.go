package collector

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestGetDruidHealthMetrics(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDruidHealthMetrics(); got != tt.want {
				t.Errorf("GetDruidHealthMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDruidSegmentData(t *testing.T) {
	tests := []struct {
		name string
		want SegementInterface
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDruidSegmentData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDruidSegmentData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDruidData(t *testing.T) {
	type args struct {
		pathURL string
	}
	tests := []struct {
		name string
		args args
		want []map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDruidData(tt.args.pathURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDruidData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricCollector_Describe(t *testing.T) {
	type args struct {
		ch chan<- *prometheus.Desc
	}
	tests := []struct {
		name      string
		collector *MetricCollector
		args      args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.collector.Describe(tt.args.ch)
		})
	}
}

func TestCollector(t *testing.T) {
	tests := []struct {
		name string
		want *MetricCollector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Collector(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricCollector_Collect(t *testing.T) {
	type args struct {
		ch chan<- prometheus.Metric
	}
	tests := []struct {
		name      string
		collector *MetricCollector
		args      args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.collector.Collect(tt.args.ch)
		})
	}
}

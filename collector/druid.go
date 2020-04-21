package collector

import (
	"fmt"
	"druid-exporter/utils"
)

// GetDruidMetrics returns the set of metrics for druid
func GetDruidMetrics() {
	fmt.Println(utils.GetDruidHealth("http://52.172.156.84:8081/status/health"))
}

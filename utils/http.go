package utils

import (
	"druid-exporter/logger"
	"github.com/go-kit/kit/log/level"
	"io/ioutil"
	"net/http"
)

// GetHealth returns that druid is healthy or not
func GetHealth(url string) float64 {
	druidLogger := logger.GetLoggerInterface()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		level.Error(druidLogger).Log("msg", "Cannot create GET request for druid healthcheck", "err", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		level.Error(druidLogger).Log("msg", "Error while making GET request for druid healthcheck", "err", err)
	}
	level.Info(druidLogger).Log("msg", "GET request is successful on druid healthcheck", "url", url)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return 1
	}

	return 0
}

// GetResponse will return API response for druid
func GetResponse(url string, queryType string) ([]byte, error) {
	druidLogger := logger.GetLoggerInterface()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		level.Error(druidLogger).Log("msg", "Cannot create http request", "err", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		level.Error(druidLogger).Log("msg", "Error while making http request", "err", err)
	}

	defer resp.Body.Close()
	level.Info(druidLogger).Log("msg", "GET request is successful for druid api", "url", url)

	return ioutil.ReadAll(resp.Body)
}

package utils

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// GetHealth returns that druid is healthy or not
func GetHealth(url string) float64 {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Cannot create GET request for druid healthcheck: %v", err)
		return 0
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Error on GET request for druid healthcheck: %v", err)
		return 0
	}
	logrus.Debugf("Successful healthcheck request for druid - %v", url)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return 1
	}

	return 0
}

// GetResponse will return API response for druid
func GetResponse(url string, queryType string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Cannot create http request: %v", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Errorf("Error on making http request for druid: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	logrus.Errorf("Successful GET request on Druid API - %v", url)

	return ioutil.ReadAll(resp.Body)
}

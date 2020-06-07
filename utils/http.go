package utils

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
)

var (
	user     = kingpin.Flag("druid.user", "HTTP basic auth username. (Only if it is set)").Default("").OverrideDefaultFromEnvar("DRUID_USER").String()
	password = kingpin.Flag("druid.password", "HTTP basic auth password. (Only if it is set)").Default("").OverrideDefaultFromEnvar("DRUID_PASSWORD").String()
	certFile = kingpin.Flag("cert", "A pem encoded certificate file. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CERT_FILE").String()
	keyFile  = kingpin.Flag("key", "A pem encoded key file. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CERT_KEY").String()
	caFile   = kingpin.Flag("ca", "A pem encoded CA's certificate file. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CA_CERT_FILE").String()
)

// GetHealth returns that druid is healthy or not
func GetHealth(url string) float64 {
	kingpin.Parse()
	client, err := generateTLSConfig()
	if err != nil {
		logrus.Errorf("Cannot generate http client: %v", err)
		return 0
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Cannot create GET request for druid healthcheck: %v", err)
		return 0
	}
	if *user != "" && *password != "" {
		req.SetBasicAuth(*user, *password)
	}
	resp, err := client.Do(req)
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
	kingpin.Parse()
	client, err := generateTLSConfig()
	if err != nil {
		logrus.Errorf("Cannot generate http client: %v", err)
		return nil, err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("Cannot create http request: %v", err)
		return nil, err
	}

	if *user != "" && *password != "" {
		req.SetBasicAuth(*user, *password)
	}

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error on making http request for druid: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	logrus.Debugf("Successful GET request on Druid API - %v", url)

	return ioutil.ReadAll(resp.Body)
}

func generateTLSConfig() (*http.Client, error) {
	kingpin.Parse()

	if *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			logrus.Errorf("Unable to load certificate file: %v", err)
			return nil, err
		}
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			logrus.Errorf("Unable to load CA's certificate file: %v", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client := &http.Client{Transport: transport}
		return client, nil
	}
	client := &http.Client{}
	return client, nil
}

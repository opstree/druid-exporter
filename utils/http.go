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
	user        = kingpin.Flag("druid.user", "HTTP basic auth username, EnvVar - DRUID_USER. (Only if it is set)").Default("").OverrideDefaultFromEnvar("DRUID_USER").String()
	password    = kingpin.Flag("druid.password", "HTTP basic auth password, EnvVar - DRUID_PASSWORD. (Only if it is set)").Default("").OverrideDefaultFromEnvar("DRUID_PASSWORD").String()
	insecureTLS = kingpin.Flag("insecure.tls.verify", "Boolean flag to skip TLS verification, EnvVar - INSECURE_TLS_VERIFY.").OverrideDefaultFromEnvar("INSECURE_TLS_VERIFY").Bool()
	certFile    = kingpin.Flag("tls.cert", "A pem encoded certificate file, EnvVar - CERT_FILE. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CERT_FILE").String()
	keyFile     = kingpin.Flag("tls.key", "A pem encoded key file, EnvVar - CERT_KEY. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CERT_KEY").String()
	caFile      = kingpin.Flag("tls.ca", "A pem encoded CA's certificate file, EnvVar - CA_CERT_FILE. (Only if tls is configured)").Default("").OverrideDefaultFromEnvar("CA_CERT_FILE").String()
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
		logrus.Debugf("Successful GET request on Druid API - %v", url)
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
		logrus.Errorf("Possible issue can be with Druid's URL, Username or Password")
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		logrus.Debugf("Successful GET request on Druid API - %v", url)
	} else {
		logrus.Errorf("Druid's API response is not 200, Status Code - %v", resp.StatusCode)
		logrus.Errorf("Possible issue can be with Druid's URL, Username or Password")
	}

	return ioutil.ReadAll(resp.Body)
}

func generateTLSConfig() (*http.Client, error) {
	kingpin.Parse()

	if *certFile != "" && *keyFile != "" && *caFile != "" {
	  // mutual TLS, verify server and client
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
	if *caFile != "" {
	  // TLS, verify server
		caCert, err := ioutil.ReadFile(*caFile)
		if err != nil {
			logrus.Errorf("Unable to load CA's certificate file: %v", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client := &http.Client{Transport: transport}
		return client, nil
	} 
	if *insecureTLS == true {
	  // TLS, no server verification
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		client := &http.Client{Transport: transport}
		return client, nil
	} 
	// http, no TLS
	client := &http.Client{}
	return client, nil
}

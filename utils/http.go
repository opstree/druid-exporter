package utils

import (
	"net/http"
	"io/ioutil"
	"github.com/rs/zerolog/log"
)

// GetHealth returns that druid is healthy or not
func GetHealth(url string) float64 {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal().Str("URL", url).Msg("Error while generating request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Str("URL", url).Msg("Error on GET request")
	}
	log.Info().Str("Method", resp.Request.Method).Str("Response", resp.Status).Str("Query Type", "Health").Msg("GET request is successful on specified URL")
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return 1
	}

	return 0
}

// GetResponse will return API response for druid
func GetResponse(url string, queryType string) ([]byte, error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal().Str("URL", url).Msg("Error while generating request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal().Str("URL", url).Msg("Error while making request")
	}

	defer resp.Body.Close()
	log.Info().Str("Method", resp.Request.Method).Str("Response", resp.Status).Str("Query Type", queryType).Msg("GET request is successful on specified URL")

	return ioutil.ReadAll(resp.Body)
}

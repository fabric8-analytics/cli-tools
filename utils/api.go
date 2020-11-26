package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

// GenericPayload is Generic Interface of Request Payload
type GenericPayload interface{}

// HTTPRequestType is request type HTTPRequest Method accepts
type HTTPRequestType struct {
	Payload         GenericPayload `json:"payload,omitempty"`
	Method          string         `json:"method,omitempty"`
	Endpoint        string         `json:"endpoint,omitempty"`
	ThreeScaleToken string         `json:"threeScale,omitempty"`
	Host            string         `json:"host,omitempty"`
}

// buildAPIURL builds API Endpoint URL
func buildAPIURL(host string, endpoint string, threeScale string) url.URL {
	log.Debug().Msgf("Building API Url.")
	APIHost, err := url.Parse(host)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to Parse Host URL")
	}
	url := url.URL{Host: APIHost.Hostname(), Path: endpoint}
	url.Scheme = "https"
	q := url.Query()
	q.Set("user_key", threeScale)
	url.RawQuery = q.Encode()
	log.Debug().Msgf("Success: Building API Url.")
	return url
}

// HTTPRequest is generic method for HTTP Requests to server
func HTTPRequest(data HTTPRequestType) *http.Response {
	log.Debug().Msgf("Executing HTTPRequest.")
	client := &http.Client{}
	url := buildAPIURL(data.Host, data.Endpoint, data.ThreeScaleToken)
	payload, _ := json.Marshal(&data.Payload)
	req, err := http.NewRequest(data.Method, url.String(), bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to build request")
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to reach the server. hint: Check your Internet connection.")
	}
	log.Debug().Msgf("Success HTTPRequest.")
	return res
}

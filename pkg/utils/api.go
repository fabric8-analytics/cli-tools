package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
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
	UserID          string         `json:"user_id,omitempty"`
	Client          string         `json:"client,omitempty"`
}

// BuildReportLink builds stack report UI Link
func BuildReportLink(stackID string) string {
	log.Debug().Msgf("Building Report Url.")
	APIHost, err := url.Parse(ActualHost)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to Parse Host URL")
	}
	endpoint := fmt.Sprintf("api/v2/stack-report/%s", stackID)
	reportURL := url.URL{Host: APIHost.Hostname(), Path: endpoint}
	reportURL.Scheme = "https"
	log.Debug().Msgf("Success Building Report Url.")
	return reportURL.String()
}

// buildAPIURL builds API Endpoint URL
func buildAPIURL(host string, endpoint string, threeScale string) *url.URL {
	log.Debug().Msgf("Building API Url.")
	APIHost, err := url.Parse(host)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to Parse Host URL")
	}
	APIHost.Path = endpoint
	q := APIHost.Query()
	q.Set("user_key", threeScale)
	APIHost.RawQuery = q.Encode()
	log.Debug().Msgf("Success: Building API Url.")
	return APIHost
}

// HTTPRequest is generic method for HTTP Requests to server
func HTTPRequest(data HTTPRequestType) *http.Response {
	log.Debug().Msgf("Executing HTTPRequest.")
	client := &http.Client{}
	apiURL := buildAPIURL(data.Host, data.Endpoint, data.ThreeScaleToken)
	payload, _ := json.Marshal(&data.Payload)
	req, err := http.NewRequest(data.Method, apiURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to build request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("uuid", data.UserID)
	req.Header.Add("client", data.Client)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to reach the server. hint: Check your Internet connection.")
	}
	log.Debug().Msgf("Success HTTPRequest.")
	return res
}

// HTTPRequestMultipart is generic method for HTTP Multipart Requests to server
func HTTPRequestMultipart(data HTTPRequestType, w *multipart.Writer, buf *bytes.Buffer) *http.Response {
	log.Debug().Msgf("Executing HTTPRequestMultipart.")
	client := &http.Client{}
	apiURL := buildAPIURL(data.Host, data.Endpoint, data.ThreeScaleToken)
	req, err := http.NewRequest(data.Method, apiURL.String(), buf)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to build request")
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("uuid", data.UserID)
	req.Header.Set("client", data.Client)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to reach the server. hint: Check your Internet connection.")
	}
	log.Debug().Msgf("Success HTTPRequestMultipart.")
	return res
}

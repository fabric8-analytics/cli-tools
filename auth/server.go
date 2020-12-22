package auth

import (
	"encoding/json"
	"net/http"

	utils "github.com/fabric8-analytics/cli-tools/utils"
	"github.com/rs/zerolog/log"
)

// RequestServerType is a argtype of RequestServer func
type RequestServerType struct {
	SynkToken       string
	UserID          string
	Host            string
	ThreeScaleToken string
}

// Payload is format for auth requests to Server
type Payload struct {
	utils.GenericPayload
	SnykToken string `json:"snyk_api_token,omitempty"`
	UUID      string `json:"user_id,omitempty"`
}

// UserResponse is format for Response from Server
type UserResponse struct {
	UUID    string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int32  `json:"status,omitempty"`
}

// RequestServer is auth request to CRDA server
func RequestServer(requestParams RequestServerType) string {
	log.Debug().Msgf("Executing Request Server.")
	var payload Payload
	var resData UserResponse
	requestData := utils.HTTPRequestType{
		Payload:         payload,
		Method:          http.MethodPost,
		Endpoint:        utils.APIUsers,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
	}

	if requestParams.UserID == "" {
		// If userID is not present, get one from server.
		apiResponse := utils.HTTPRequest(requestData)
		resData = validateResponse(apiResponse)

		requestParams.UserID = resData.UUID

	}
	if requestParams.SynkToken != "" {
		// If Snyk Token is present in input.
		payload = Payload{
			UUID:      requestParams.UserID,
			SnykToken: requestParams.SynkToken,
		}
		requestData.Payload = payload
		requestData.Method = http.MethodPut
		resData = validateResponse(utils.HTTPRequest(requestData))
	}
	log.Debug().Msgf("Successfully executed RequestServer.")
	return requestParams.UserID
}

// validateResponse validates API Response and Generate Auth Specific API Errors, if required.
func validateResponse(apiResponse *http.Response) UserResponse {
	log.Debug().Msgf("Executing validateResponse.")
	var body UserResponse
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf(body.Message)
	}
	log.Debug().Msgf("Successfully executed validateResponse.")
	return body
}

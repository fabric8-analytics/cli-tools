package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"

	"github.com/fabric8-analytics/cli-tools/pkg/utils"
	"github.com/rs/zerolog/log"
)

// RequestServerType is a argtype of RequestServer func
type RequestServerType struct {
	SynkToken       string
	UserID          string
	Host            string
	ThreeScaleToken string
	Client          string
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

// API Constants
const (
	APIUsers = "/user"
)

// RequestServer is auth request to CRDA server
func RequestServer(ctx context.Context, requestParams RequestServerType) (string, error) {
	log.Debug().Msgf("Executing Request Server.")
	var payload Payload

	requestData := utils.HTTPRequestType{
		Payload:         payload,
		Method:          http.MethodPost,
		Endpoint:        APIUsers,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
		Client:          requestParams.Client,
	}

	if requestParams.UserID == "" {
		// If userID is not present, get one from server.
		apiResponse := utils.HTTPRequest(requestData)
		resData, err := validateResponse(apiResponse)
		if err != nil {
			return "", err
		}
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
		_, err := validateResponse(utils.HTTPRequest(requestData))
		if err != nil {
			return "", err
		}
		telemetry.SetSnykTokenAssociation(ctx, true)
	}
	log.Debug().Msgf("Successfully executed RequestServer.")
	return requestParams.UserID, nil
}

// validateResponse validates API Response and Generate Auth Specific API Errors, if required.
func validateResponse(apiResponse *http.Response) (*UserResponse, error) {
	log.Debug().Msgf("Executing validateResponse.")
	var body UserResponse
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		return nil, errors.New(body.Message)
	}
	if err != nil {
		log.Error().Msg("Unable to decode Snyk Token. Please try again.")
		return nil, err
	}
	log.Debug().Msgf("Successfully executed validateResponse.")
	return &body, nil
}

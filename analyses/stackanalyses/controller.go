package stackanalyses

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/jpillora/backoff"
	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/analyses/golang"
	"github.com/fabric8-analytics/cli-tools/analyses/maven"
	"github.com/fabric8-analytics/cli-tools/analyses/npm"
	"github.com/fabric8-analytics/cli-tools/analyses/pypi"
	"github.com/fabric8-analytics/cli-tools/analyses/summary"
	"github.com/fabric8-analytics/cli-tools/analyses/verbose"
	"github.com/fabric8-analytics/cli-tools/utils"
)

// Controller is a control structure.
type Controller struct {
	// an implemented Matcher
	m         driver.StackAnalysisInterface
	fileStats *driver.ReadManifestResponse
}

// API Constants
const (
	APIStackAnalyses = "/api/v2/stack-analyses"
	RegisteredStatus = "REGISTERED"
)

//StackAnalyses is main controller function for analyse command. This function is responsible for all communications between cmd and custom packages.
func StackAnalyses(ctx context.Context, requestParams driver.RequestType, jsonOut bool, verboseOut bool) (bool, error) {
	log.Debug().Msgf("Executing StackAnalyses.")
	var hasVul bool
	matcher, err := GetMatcher(requestParams.RawManifestFile)
	if err != nil {
		return hasVul, err
	}
	mc := NewController(matcher)
	mc.fileStats = mc.buildFileStats(requestParams.RawManifestFile)
	postResponse, err := mc.postRequest(requestParams, mc.fileStats.DepsTreePath)
	if err != nil {
		return hasVul, err
	}
	getResponse, err := mc.getRequest(requestParams, postResponse)
	if err != nil {
		return hasVul, err
	}
	verboseEligible := getResponse.RegistrationStatus == RegisteredStatus
	showVerboseMsg := verboseOut && !verboseEligible

	if verboseOut && verboseEligible {
		hasVul = verbose.ProcessVerbose(ctx, getResponse, jsonOut)
	} else {
		hasVul = summary.ProcessSummary(ctx, getResponse, jsonOut, showVerboseMsg)
	}
	telemetry.SetEcosystem(ctx, mc.fileStats.Ecosystem)
	log.Debug().Msgf("Success StackAnalyses.")
	return hasVul, nil
}

// postRequest performs Stack Analyses POST Request to CRDA server.
func (mc *Controller) postRequest(requestParams driver.RequestType, filePath string) (*driver.PostResponseType, error) {
	log.Debug().Msgf("Executing: postRequest.")
	manifest := &bytes.Buffer{}
	requestData := utils.HTTPRequestType{
		Method:          http.MethodPost,
		Endpoint:        APIStackAnalyses,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
		UserID:          requestParams.UserID,
	}
	writer := multipart.NewWriter(manifest)
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fw, err := writer.CreateFormFile("manifest", mc.m.DepsTreeFileName())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, fd)
	if err != nil {
		return nil, err
	}
	_ = writer.WriteField("ecosystem", mc.m.Ecosystem())
	_ = writer.WriteField("file_path", "/tmp/bin")
	err = writer.Close()
	if err != nil {
		return nil, errors.New("error closing Buffer Writer in Stack Analyses Request")
	}
	log.Debug().Msgf("Hitting: Stack Analyses Post API.")
	apiResponse := utils.HTTPRequestMultipart(requestData, writer, manifest)
	body, err := mc.validatePostResponse(apiResponse)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Got Stack Analyses Post Response Stack Id: %s", body.ID)
	log.Debug().Msgf("Success: postRequest.")
	return body, nil
}

// getRequest performs Stack Analyses GET Request to CRDA Server.
func (mc *Controller) getRequest(requestParams driver.RequestType, postResponse *driver.PostResponseType) (*driver.GetResponseType, error) {
	log.Debug().Msgf("Executing: getRequest.")
	polling := &backoff.Backoff{
		Min:    5 * time.Second,
		Max:    120 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	var apiResponse *http.Response
	requestData := utils.HTTPRequestType{
		Method:          http.MethodGet,
		Endpoint:        APIStackAnalyses + "/" + postResponse.ID,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
		UserID:          requestParams.UserID,
	}
	for {
		d := polling.Duration()
		log.Debug().Msgf("Sleeping for %s", d)
		time.Sleep(d)
		apiResponse = utils.HTTPRequest(requestData)
		if apiResponse.StatusCode != http.StatusAccepted {
			// Break when server returns anything other than 202.
			break
		}
		log.Debug().Msgf("Retrying...")
	}
	body, err := mc.validateGetResponse(apiResponse)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// validatePostResponse validates Stack Analyses POST API Response.
func (mc *Controller) validatePostResponse(apiResponse *http.Response) (*driver.PostResponseType, error) {
	log.Debug().Msgf("Executing validatePostResponse.")
	var body driver.PostResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)

	// In Case of Authentication Failure, json is not return from API, Need to catch before decoding.
	if apiResponse.StatusCode == http.StatusForbidden {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Error().Msgf("Stack Analyses Post Request Failed.  Please check auth token and try again.")
		return nil, fmt.Errorf("invalid authentication token")
	}

	if err != nil {
		return nil, err
	}
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Error().Msgf("Stack Analyses Post Request Failed with status code %d.  Please retry after sometime. If issue persists, Please raise at https://github.com/fabric8-analytics/cli-tools/issues.\"", apiResponse.StatusCode)
		return nil, fmt.Errorf("message from server: %s", body.Error)
	}
	log.Debug().Msgf("Success validatePostResponse.")
	return &body, nil
}

// validateGetResponse validates Stack Analyses GET API Response.
func (mc *Controller) validateGetResponse(apiResponse *http.Response) (*driver.GetResponseType, error) {
	log.Debug().Msgf("Executing validateGetResponse.")
	var body driver.GetResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Error().Msgf("Stack Analyses Get Request Failed with status code %d.  Please retry after sometime. If issue persists, Please raise at https://github.com/fabric8-analytics/cli-tools/issues.\"", apiResponse.StatusCode)
		return nil, fmt.Errorf("message from server: %s", body.Error)
	}
	log.Debug().Msgf("Success validateGetResponse.")
	return &body, err
}

// NewController is a constructor for a Controller
func NewController(m driver.StackAnalysisInterface) *Controller {
	return &Controller{
		m: m,
	}
}

// defaultMatchers is a variable containing all the matchers.
var defaultMatchers = []driver.StackAnalysisInterface{
	&pypi.Matcher{},
	&maven.Matcher{},
	&golang.Matcher{},
	&npm.Matcher{},
}

// GetMatcher returns ecosystem specific matcher
func GetMatcher(manifestFile string) (driver.StackAnalysisInterface, error) {
	for _, matcher := range defaultMatchers {
		if matcher.IsSupportedManifestFormat(manifestFile) {
			return matcher, nil
		}
	}
	return nil, errors.New("ecosystem not supported yet")
}

func (mc *Controller) buildFileStats(manifestFile string) *driver.ReadManifestResponse {
	stats := &driver.ReadManifestResponse{
		Ecosystem:        mc.m.Ecosystem(),
		RawFileName:      GetManifestName(manifestFile),
		RawFilePath:      manifestFile,
		DepsTreePath:     mc.m.GeneratorDependencyTree(manifestFile),
		DepsTreeFileName: mc.m.DepsTreeFileName(),
	}
	return stats
}

// GetManifestName extracts manifest name from user input path
func GetManifestName(manifestFile string) string {
	stats, err := os.Stat(manifestFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error fetching manifest name.")
	}
	return stats.Name()
}

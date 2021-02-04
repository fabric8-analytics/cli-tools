package stackanalyses

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jpillora/backoff"
	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/analyses/golang"
	"github.com/fabric8-analytics/cli-tools/analyses/maven"
	"github.com/fabric8-analytics/cli-tools/analyses/npm"
	"github.com/fabric8-analytics/cli-tools/analyses/pypi"
	"github.com/fabric8-analytics/cli-tools/analyses/summary"
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
)

//StackAnalyses Performs Full Stack Analyses
func StackAnalyses(requestParams driver.RequestType, jsonOut bool) bool {
	log.Debug().Msgf("Executing StackAnalyses.")
	matcher, err := GetMatcher(requestParams.RawManifestFile)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	mc := NewController(matcher)
	mc.fileStats = mc.buildFileStats(requestParams.RawManifestFile)
	postResponse := mc.postRequest(requestParams, mc.fileStats.DepsTreePath)
	getResponse := mc.getRequest(requestParams, postResponse)
	hasVul := summary.ProcessSummary(getResponse, jsonOut)
	log.Debug().Msgf("Success StackAnalyses.")
	return hasVul
}

// GetManifestFilePath sets file path
func (mc *Controller) GetManifestFilePath(input string) string {
	path, err := filepath.Abs(input)
	if err != nil {
		log.Fatal().Msgf("Invalid Path of Manifest file. Only Absolute path is allowed.")
	}
	return path
}

// postRequest performs Stack Analyses POST Request to CRDA server.
func (mc *Controller) postRequest(requestParams driver.RequestType, filePath string) driver.PostResponseType {
	log.Debug().Msgf("Executing: postRequest.")
	manifest := &bytes.Buffer{}
	requestData := utils.HTTPRequestType{
		Method:          http.MethodPost,
		Endpoint:        APIStackAnalyses,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
	}
	writer := multipart.NewWriter(manifest)
	fd, err := os.Open(filePath)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	defer fd.Close()

	fw, err := writer.CreateFormFile("manifest", mc.m.DepsTreeFileName())
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	_ = writer.WriteField("ecosystem", mc.m.Ecosystem())
	_ = writer.WriteField("file_path", "/tmp/bin")
	err = writer.Close()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error closing Buffer Writer in Stack Analyses Request.")
	}
	log.Debug().Msgf("Hitting: Stack Analyses Post API.")
	apiResponse := utils.HTTPRequestMultipart(requestData, writer, manifest)
	body := mc.validatePostResponse(apiResponse)
	log.Debug().Msgf("Got Stack Analyses Post Response Stack Id: %s", body.ID)
	log.Debug().Msgf("Success: postRequest.")
	return body
}

// getRequest performs Stack Analyses GET Request to CRDA Server.
func (mc *Controller) getRequest(requestParams driver.RequestType, postResponse driver.PostResponseType) driver.GetResponseType {
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
	body := mc.validateGetResponse(apiResponse)
	return body
}

// validatePostResponse validates Stack Analyses POST API Response.
func (mc *Controller) validatePostResponse(apiResponse *http.Response) driver.PostResponseType {
	log.Debug().Msgf("Executing validatePostResponse.")
	var body driver.PostResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf("Stack Analyses Post Request Failed. Please retry after sometime. If issue persists, Please raise at https://github.com/fabric8-analytics/cli-tools/issues.")
	}
	log.Debug().Msgf("Success validatePostResponse.")
	return body
}

// validateGetResponse validates Stack Analyses GET API Response.
func (mc *Controller) validateGetResponse(apiResponse *http.Response) driver.GetResponseType {
	log.Debug().Msgf("Executing validateGetResponse.")
	var body driver.GetResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf("Stack Analyses Request Failed. Please retry after sometime. If issue persists, Please raise at https://github.com/fabric8-analytics/cli-tools/issues.")
	}
	log.Debug().Msgf("Success validateGetResponse.")
	return body
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
		RawFileName:      mc.getManifestName(manifestFile),
		RawFilePath:      mc.GetManifestFilePath(manifestFile),
		DepsTreePath:     mc.m.GeneratorDependencyTree(manifestFile),
		DepsTreeFileName: mc.m.DepsTreeFileName(),
	}
	return stats
}

func (mc *Controller) getManifestName(manifestFile string) string {
	stats, err := os.Stat(manifestFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error fetching manifest name.")
	}
	return stats.Name()
}

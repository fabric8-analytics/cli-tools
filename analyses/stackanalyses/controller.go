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

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/analyses/pypi"
	"github.com/fabric8-analytics/cli-tools/utils"
	"github.com/jpillora/backoff"
)

// Controller is a control structure.
type Controller struct {
	// an implemented Matcher
	m         driver.StackAnalysisInterface
	fileStats *driver.ReadManifestResponse
}

//StackAnalyses Performs Full Stack Analyses
func StackAnalyses(requestParams driver.RequestType) driver.GetResponseType {
	log.Info().Msgf("Performing full Stack Analyses. Please wait...")
	log.Debug().Msgf("Executing StackAnalyses.")
	b := &backoff.Backoff{
		Min:    5 * time.Second,
		Max:    120 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	matcher, err := GetMatcher(requestParams.Ecosystem)
	if err != nil {
		log.Fatal().Msgf(err.Error())
	}
	mc := NewController(matcher)
	mc.fileStats = mc.buildFileStats(requestParams.RawManifestFile)
	if !mc.m.IsSupportedManifestFormat(mc.fileStats.RawFileName) {
		log.Fatal().Msgf("Manifest File not supported.")
	}
	postResponse := mc.postRequest(requestParams, mc.fileStats.DepsTreePath)
	getResponse := mc.getRequest(requestParams, postResponse, b)
	log.Debug().Msgf("Success StackAnalyses.")
	return getResponse
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
		Endpoint:        utils.APISA,
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
		log.Fatal().Err(err).Msgf("Error closing Buffer Writer in SA Request.")
	}
	log.Debug().Msgf("Hitting: SA Post API.")
	apiResponse := utils.HTTPRequestMultipart(requestData, writer, manifest)
	body := mc.validatePostResponse(apiResponse)
	log.Debug().Msgf("Got SA Post Response Stack Id: %s", body.ID)
	log.Debug().Msgf("Success: postRequest.")
	return body
}

// getRequest performs Stack Analyses GET Request to CRDA Server.
func (mc *Controller) getRequest(requestParams driver.RequestType, saPost driver.PostResponseType, back *backoff.Backoff) driver.GetResponseType {
	log.Debug().Msgf("Executing: getRequest.")
	requestData := utils.HTTPRequestType{
		Method:          http.MethodGet,
		Endpoint:        utils.APISA + "/" + saPost.ID,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
	}
	d := back.Duration()
	log.Debug().Msgf("Sleeping for %s", d)
	time.Sleep(d)
	apiResponse := utils.HTTPRequest(requestData)
	if apiResponse.StatusCode == http.StatusAccepted {
		// Retry till server respond 200 or Timeout Error or Exponential Backoff limit hit.
		log.Debug().Msgf("Retrying...")
		mc.getRequest(requestParams, saPost, back)
	}
	body := mc.validateGetResponse(apiResponse)
	return body
}

// validateResponse validates API Response.
func (mc *Controller) validatePostResponse(apiResponse *http.Response) driver.PostResponseType {
	log.Debug().Msgf("Executing validatePostResponse.")
	var body driver.PostResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success validatePostResponse.")
	return body
}

// validateGetResponse validates API Response.
func (mc *Controller) validateGetResponse(apiResponse *http.Response) driver.GetResponseType {
	log.Debug().Msgf("Executing validateGetResponse.")
	var body driver.GetResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf("SA Request Failed. Please retry after sometime. If issue persists, Please raise at https://github.com/fabric8-analytics/cli-tools/issues.")
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
}

// GetMatcher returns ecosystem specific matcher
func GetMatcher(ecosystem string) (driver.StackAnalysisInterface, error) {
	for _, matcher := range defaultMatchers {
		if matcher.Filter(ecosystem) {
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
		log.Fatal().Err(err).Msgf("Error")
	}
	return stats.Name()
}

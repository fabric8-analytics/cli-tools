package stackanalyses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/pypi"
	"github.com/fabric8-analytics/cli-tools/utils"
	"github.com/jpillora/backoff"
)

// RequestType is a argtype of RequestServer func
type RequestType struct {
	UserID          string
	Host            string
	ThreeScaleToken string
	ShellPath       string
	RawManifestFile string
	DepsTreePath    string
}

// PostResponseType is a argtype of RequestServer func
type PostResponseType struct {
	SubmittedAt string `json:"submitted_at,omitempty"`
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
}

// GetResponseType is a argtype of RequestServer func
type GetResponseType struct {
	AnalysedDeps    []interface{}          `json:"analyzed_dependencies"`
	Ecosystem       string                 `json:"ecosystem"`
	Recommendation  map[string]interface{} `json:"recommendation"`
	LicenseAnalyses map[string]interface{} `json:"license_analysis"`
}

// ReadManifestResponse is arg type of readManifest func
type ReadManifestResponse struct {
	DepsTreePath     string `json:"manifest,omitempty"`
	RawFileName      string `json:"file,omitempty"`
	Eco              string `json:"ecosystem,omitempty"`
	DepsTreeFileName string `json:"deps_tree,omitempty"`
}

//StackAnalyses Performs Full Stack Analyses
func StackAnalyses(requestParams RequestType) GetResponseType {
	log.Info().Msgf("Performing full Stack Analyses. Please wait...")
	log.Debug().Msgf("Executing StackAnalyses.")
	b := &backoff.Backoff{
		Min:    5 * time.Second,
		Max:    120 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	fileStats := readManifest(requestParams.ShellPath, requestParams.RawManifestFile)
	postResponse := postRequest(requestParams, fileStats)
	getResponse := getRequest(requestParams, postResponse, b)
	log.Debug().Msgf("Success StackAnalyses.")
	return getResponse
}

// postRequest performs Stack Analyses POST Request to CRDA server.
func postRequest(requestParams RequestType, fileStats ReadManifestResponse) PostResponseType {
	log.Debug().Msgf("Executing: postRequest.")
	manifest := &bytes.Buffer{}
	requestData := utils.HTTPRequestType{
		Method:          http.MethodPost,
		Endpoint:        utils.APISA,
		ThreeScaleToken: requestParams.ThreeScaleToken,
		Host:            requestParams.Host,
	}
	writer := multipart.NewWriter(manifest)
	fd, err := os.Open(fileStats.DepsTreePath)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	defer fd.Close()

	fw, err := writer.CreateFormFile("manifest", fileStats.DepsTreeFileName)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	_, err = io.Copy(fw, fd)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	_ = writer.WriteField("ecosystem", fileStats.Eco)
	_ = writer.WriteField("file_path", "/tmp/bin")
	err = writer.Close()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error closing Buffer Writer in SA Request.")
	}
	apiResponse := utils.HTTPRequestMultipart(requestData, writer, manifest)
	body := validatePostResponse(apiResponse)
	log.Debug().Msgf("Success: postRequest.")
	return body
}

// getRequest performs Stack Analyses GET Request to CRDA Server.
func getRequest(requestParams RequestType, saPost PostResponseType, back *backoff.Backoff) GetResponseType {
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
		log.Debug().Msgf("Retying...")
		getRequest(requestParams, saPost, back)
	}
	body := validateGetResponse(apiResponse)
	return body
}

// validateResponse validates API Response.
func validatePostResponse(apiResponse *http.Response) PostResponseType {
	log.Debug().Msgf("Executing validatePostResponse.")
	var body PostResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success validatePostResponse.")
	return body
}

// validateGetResponse validates API Response.
func validateGetResponse(apiResponse *http.Response) GetResponseType {
	log.Debug().Msgf("Executing validateGetResponse.")
	var body GetResponseType
	err := json.NewDecoder(apiResponse.Body).Decode(&body)
	if apiResponse.StatusCode != http.StatusOK {
		log.Debug().Msgf("Status from Server: %d", apiResponse.StatusCode)
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success validateGetResponse.")
	return body
}

// readManifest Manifest File validator and reader.
func readManifest(shellPath string, manifestFile string) ReadManifestResponse {
	log.Debug().Msgf("Executing readManifest.")
	stats, err := os.Stat(manifestFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error")
	}
	fileStats := ReadManifestResponse{
		RawFileName: stats.Name(),
	}
	msg := fmt.Sprintf("Support for %s is coming soon. Thanks for your Patience. :)", fileStats.RawFileName)
	switch fileStats.RawFileName {
	case "requirements.txt":
		fileStats.DepsTreePath = pypi.GeneratePylist(shellPath, manifestFile)
		fileStats.Eco = "pypi"
		fileStats.DepsTreeFileName = "pylist.json"
		return fileStats
	case "go.mod":
		log.Info().Err(err).Msgf(msg)
		os.Exit(1)
	case "pom.xml":
		log.Info().Err(err).Msgf(msg)
		os.Exit(1)
	case "package.json":
		log.Info().Err(err).Msgf(msg)
		os.Exit(1)
	default:
		log.Fatal().Err(err).Msgf("Manifest file not supported. Please try again with one of following: requirements.txt, go.mod, pom.xml or package.json.")
	}
	log.Debug().Msgf("Success readManifest.")
	return fileStats
}

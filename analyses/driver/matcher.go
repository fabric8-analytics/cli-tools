package driver

import (
	"net/url"
)

// RequestType is a argtype of RequestServer func
type RequestType struct {
	UserID          string
	Host            string
	ThreeScaleToken string
	RawManifestFile string
	DepsTreePath    string
}

// PostResponseType is a argtype of RequestServer func
type PostResponseType struct {
	SubmittedAt string `json:"submitted_at,omitempty"`
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
}

// VulnerabilitiesType is a Vulnerability Response structure
type VulnerabilitiesType struct {
	CveID    []interface{} `json:"cve_id"`
	Cvss     int           `json:"cvss"`
	CvssV3   string        `json:"cvss_v3"`
	Cwes     []interface{} `json:"cwes"`
	ID       string        `json:"id"`
	Severity string        `json:"severity"`
	Title    string        `json:"title"`
	URL      url.URL       `json:"url"`
}

// Transitives type for Transitives
type Transitives struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// VulnerableDependencies is type for Transitives in direct Dependencies
type VulnerableDependencies struct {
	PrivateVulnerabilities []VulnerabilitiesType `json:"private_vulnerabilities"`
	PublicVulnerabilities  []VulnerabilitiesType `json:"public_vulnerabilities"`
}

// AnalysedDepsType is type for Analysed Deps API Response
type AnalysedDepsType struct {
	Transitives            []Transitives            `json:"dependencies"`
	Ecosystem              string                   `json:"ecosystem"`
	LatestVersion          string                   `json:"latest_version"`
	Licenses               []interface{}            `json:"licenses"`
	Name                   string                   `json:"name"`
	PrivateVulnerabilities []VulnerabilitiesType    `json:"private_vulnerabilities"`
	PublicVulnerabilities  []VulnerabilitiesType    `json:"public_vulnerabilities"`
	RecommendedVersion     string                   `json:"recommended_version"`
	Version                string                   `json:"version"`
	VulnerableDependencies []VulnerableDependencies `json:"vulnerable_dependencies"`
}

// GetResponseType is a argtype of RequestServer func
type GetResponseType struct {
	AnalysedDeps []AnalysedDepsType `json:"analyzed_dependencies"`
}

// ReadManifestResponse is arg type of readManifest func
type ReadManifestResponse struct {
	DepsTreePath     string `json:"manifest,omitempty"`
	RawFileName      string `json:"file,omitempty"`
	RawFilePath      string `json:"filepath,omitempty"`
	Ecosystem        string `json:"ecosystem,omitempty"`
	DepsTreeFileName string `json:"deps_tree,omitempty"`
}

// StackAnalysisInterface is implemented by each ecosystem
type StackAnalysisInterface interface {
	DepsTreeFileName() string
	Ecosystem() string
	IsSupportedManifestFormat(string) bool
	GeneratorDependencyTree(string) string
}

package driver

import "net/url"

// RequestType is a argtype of RequestServer func
type RequestType struct {
	UserID          string
	Host            string
	ThreeScaleToken string
	RawManifestFile string
	DepsTreePath    string
}

// CIFormat is SA Result output format required by CI
type CIFormat struct {
	TotalScannedDependencies                   int          `json:"total_scanned_dependencies"`
	DirectDependenciesWithKnownVulnerabilities int          `json:"direct_dependencies_with_known_vulnerabilites"`
	DirectDependenciesWithSynkAdvisories       int          `json:"direct_dependencies_with_synk_advisories"`
	VulnerabilitiesPerSeverity                 SeverityType `json:"vulnerabilities_per_severity"`
}

// PostResponseType is a argtype of RequestServer func
type PostResponseType struct {
	SubmittedAt string `json:"submitted_at,omitempty"`
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
}

// SeverityType is Possible Types of Severities from Server
type SeverityType struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
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

// AnalysedDepsType is type for Analysed Deps API Response
type AnalysedDepsType struct {
	Dependencies           []interface{}          `json:"dependencies"`
	Ecosystem              string                 `json:"ecosystem"`
	Github                 map[string]interface{} `json:"github"`
	LatestVersion          string                 `json:"latest_version"`
	Licenses               []interface{}          `json:"licenses"`
	Name                   string                 `json:"name"`
	PrivateVulnerabilities []VulnerabilitiesType  `json:"private_vulnerabilities"`
	PublicVulnerabilities  []VulnerabilitiesType  `json:"public_vulnerabilities"`
	RecommendedVersion     string                 `json:"recommended_version"`
	URL                    url.URL                `json:"url"`
	Version                string                 `json:"version"`
	VulnerableDependencies []interface{}          `json:"vulnerable_dependencies"`
}

// GetResponseType is a argtype of RequestServer func
type GetResponseType struct {
	AnalysedDeps        []AnalysedDepsType     `json:"analyzed_dependencies"`
	Ecosystem           string                 `json:"ecosystem"`
	Recommendation      map[string]interface{} `json:"recommendation"`
	LicenseAnalyses     map[string]interface{} `json:"license_analysis"`
	UnknownDependencies map[string]interface{} `json:"unknown_dependencies"`
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

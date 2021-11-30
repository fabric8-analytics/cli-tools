package driver

// RequestType is a argtype of RequestServer func
type RequestType struct {
	UserID          string
	Host            string
	ThreeScaleToken string
	RawManifestFile string
	DepsTreePath    string
	Client          string
	Ignore          map[string]map[string][]string
}

// PostResponseType is a argtype of RequestServer func
type PostResponseType struct {
	SubmittedAt string `json:"submitted_at,omitempty"`
	Status      string `json:"status,omitempty"`
	ID          string `json:"id,omitempty"`
	Error       string `json:"error,omitempty"`
}

// VulnerabilitiesType is a Vulnerability Response structure
type VulnerabilitiesType struct {
	CveID    []string `json:"cve_ids"`
	Cvss     float32  `json:"cvss"`
	ID       string   `json:"id"`
	Severity string   `json:"severity"`
	Title    string   `json:"title"`
	URL      string   `json:"url"`
	Kind     string   `json:"kind"`
}

// Transitives type for Transitives
type Transitives struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// AnalysedDepsType is type for Analysed Deps API Response
type AnalysedDepsType struct {
	Transitives            []Transitives         `json:"dependencies"`
	Ecosystem              string                `json:"ecosystem"`
	LatestVersion          string                `json:"latest_version"`
	Licenses               []interface{}         `json:"licenses"`
	Name                   string                `json:"name"`
	PrivateVulnerabilities []VulnerabilitiesType `json:"private_vulnerabilities"`
	IgnoredTransitiveVulns int                   `json:"ignored_trans_vulnerability_count,omitempty"`
	IgnoredVulns           int                   `json:"ignored_vulnerability_count,omitempty"`
	PublicVulnerabilities  []VulnerabilitiesType `json:"public_vulnerabilities"`
	RecommendedVersion     string                `json:"recommended_version"`
	Version                string                `json:"version"`
	VulnerableDependencies []AnalysedDepsType    `json:"vulnerable_dependencies"`
}

// GetResponseType is a argtype of RequestServer func
type GetResponseType struct {
	AnalysedDeps       []AnalysedDepsType `json:"analyzed_dependencies"`
	RegistrationStatus string             `json:"registration_status"`
	StackID            string             `json:"external_request_id"`
	Error              string             `json:"error,omitempty"`
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
	IgnoreVulnerabilities(string) (map[string][]string, error)
}

package driver

// RequestType is a argtype of RequestServer func
type RequestType struct {
	UserID          string
	Ecosystem       string
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
	RawFilePath      string `json:"filepath,omitempty"`
	Ecosystem        string `json:"ecosystem,omitempty"`
	DepsTreeFileName string `json:"deps_tree,omitempty"`
}

// StackAnalysisInterface is implemented by each ecosystem
type StackAnalysisInterface interface {
	DepsTreeFileName() string
	Filter(string) bool
	Ecosystem() string
	IsSupportedManifestFormat(string) bool
	GeneratorDependencyTree(string) string
}

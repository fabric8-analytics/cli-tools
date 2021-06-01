package verbose

import (
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
)

// StackVerbose is SA Result Verbose output
type StackVerbose struct {
	Dependencies                   []DependenciesType `json:"analysed_dependencies"`
	TotalDirectVulnerabilities     int                `json:"total_direct_vulnerabilities"`
	TotalTransitives               int                `json:"total_transitives_scanned"`
	TotalTransitiveVulnerabilities int                `json:"transitive_vulnerabilities"`
	Severity                       SeverityType       `json:"severity"`
	ReportLink                     string             `json:"report_link"`
}

// VulnerabilityType type for Vulnerability in verbose output
type VulnerabilityType struct {
	Severity string `json:"severity"`
	ID       string `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

// DependenciesType verbose output
type DependenciesType struct {
	Name                             string              `json:"name"`
	Version                          string              `json:"version"`
	Transitives                      []DependenciesType  `json:"transitives"`
	LatestVersion                    string              `json:"latest_version"`
	RecommendedVersion               string              `json:"recommended_version"`
	PubliclyAvailableVulnerabilities []VulnerabilityType `json:"publicly_available_vulnerabilities"`
	VulnerabilitiesUniqueToSynk      []VulnerabilityType `json:"vulnerabilities_unique_with_snyk"`
	VulnerableTransitives            []DependenciesType  `json:"vulnerable_transitives"`
}

// SeverityType is Possible Types of Severities from Server
type SeverityType struct {
	Low      []driver.VulnerabilitiesType `json:"low"`
	Medium   []driver.VulnerabilitiesType `json:"medium"`
	High     []driver.VulnerabilitiesType `json:"high"`
	Critical []driver.VulnerabilitiesType `json:"critical"`
}

// ProcessVulnerabilities is arg type of processVulnerabilities
type ProcessVulnerabilities struct {
	AnalysedDependencies           []DependenciesType
	TotalTransitives               int
	TotalDirectVulnerabilities     int
	TotalTransitiveVulnerabilities int
	Severities                     SeverityType
}

// CustomColors maintain state of custom colors
type CustomColors struct {
	Green func(...interface{}) string
	White func(...interface{}) string
	Cyan  func(...interface{}) string
	Red   func(...interface{}) string
}

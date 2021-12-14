package summary

// SeverityType is Possible Types of Severities from Server
type SeverityType struct {
	Low      int `json:"low"`
	Medium   int `json:"medium"`
	High     int `json:"high"`
	Critical int `json:"critical"`
}

// StackSummary is SA Result Summary output
type StackSummary struct {
	TotalScannedDependencies              int    `json:"total_scanned_dependencies"`
	TotalScannedTransitiveDependencies    int    `json:"total_scanned_transitives"`
	TotalVulnerabilities                  int    `json:"total_vulnerabilities"`
	PubliclyAvailableVulnerabilities      int    `json:"publicly_available_vulnerabilities"`
	VulnerabilitiesUniqueToSynk           int    `json:"vulnerabilities_unique_to_synk"`
	DirectVulnerableDependencies          int    `json:"direct_vulnerable_dependencies"`
	LowVulnerabilities                    int    `json:"low_vulnerabilities"`
	MediumVulnerabilities                 int    `json:"medium_vulnerabilities"`
	HighVulnerabilities                   int    `json:"high_vulnerabilities"`
	CriticalVulnerabilities               int    `json:"critical_vulnerabilities"`
	ReportLink                            string `json:"report_link"`
	TotalDirectVulnerabilitiesIgnored     int    `json:"total_direct_vulns_ignored"`
	TotalTransitiveVulnerabilitiesIgnored int    `json:"total_transitive_vulns_ignored"`
	SnykTokenStatus						  bool   `json:"snyk_token_status"`
}

// ProcessVulnerabilities is arg type of processVulnerabilities
type ProcessVulnerabilities struct {
	PublicVul                    int
	PrivateVul                   int
	DirectVulnerableDependencies int
	TotalTransitives             int
	TotalDirectIgnored           int
	TotalTransitiveIgnored       int
	Severities                   SeverityType
}

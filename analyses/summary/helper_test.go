package summary

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

func data() driver.GetResponseType {
	var GetResponse = &driver.GetResponseType{
		AnalysedDeps: []driver.AnalysedDepsType{
			{
				Ecosystem:     "pypi",
				LatestVersion: "7.1.2",
				Name:          "click",
				PublicVulnerabilities: []driver.VulnerabilitiesType{
					{
						Severity: "medium",
						ID:       "ABC-PYTHON-CODECOV-12345",
						Title:    "Command Injection",
					},
				},
			},
		},
	}
	return *GetResponse
}

func TestProcessSummary(t *testing.T) {
	got := ProcessSummary(data(), false)
	if got != true {
		t.Errorf("Error in ProcessSummary.")
	}
}

func TestGetResultSummary(t *testing.T) {
	got := getResultSummary(data())
	want := &StackSummary{
		TotalScannedDependencies:           1,
		TotalScannedTransitiveDependencies: 0,
		TotalVulnerabilities:               1,
		CommonlyKnownVulnerabilities:       1,
		VulnerabilitiesUniqueToSynk:        0,
		DirectVulnerableDependencies:       1,
		CriticalVulnerabilities:            0,
		HighVulnerabilities:                0,
		MediumVulnerabilities:              1,
		LowVulnerabilities:                 0,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Vuln mismatch (-want, +got):\n%s", diff)
	}
}

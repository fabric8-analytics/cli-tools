package summary

import (
	"context"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

func data() *driver.GetResponseType {
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
		StackID: "123456789",
	}
	return GetResponse
}

func TestProcessSummary(t *testing.T) {
	var ctx = telemetry.NewContext(context.Background())
	got := ProcessSummary(ctx, data(), false, false)
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
		PubliclyAvailableVulnerabilities:   1,
		VulnerabilitiesUniqueToSynk:        0,
		DirectVulnerableDependencies:       1,
		CriticalVulnerabilities:            0,
		HighVulnerabilities:                0,
		MediumVulnerabilities:              1,
		LowVulnerabilities:                 0,
		ReportLink:                         "https://recommender.api.openshift.io/api/v2/stack-report/123456789",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Vuln mismatch (-want, +got):\n%s", diff)
	}
}

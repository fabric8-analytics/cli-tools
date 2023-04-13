package verbose

import (
	"context"
	"encoding/json"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
)

func data() *driver.GetResponseType {
	var body driver.GetResponseType
	// json.NewDecoder(apiResponse.Body).Decode(&body)
	plan, _ := os.ReadFile("testdata/getresponse.json")
	_ = json.Unmarshal(plan, &body)
	return &body
}
func verboseData() *StackVerbose {
	var body StackVerbose
	plan, _ := os.ReadFile("testdata/verbosedata.json")
	_ = json.Unmarshal(plan, &body)
	return &body
}

func TestProcessSummary(t *testing.T) {
	var ctx = telemetry.NewContext(context.Background())
	got := ProcessVerbose(ctx, data(), false)
	if got != true {
		t.Errorf("Error in ProcessSummary.")
	}
}

func TestGetResultSummary(t *testing.T) {
	got := getVerboseResult(data())
	want := verboseData()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Vuln mismatch (-want, +got):\n%s", diff)
	}
}

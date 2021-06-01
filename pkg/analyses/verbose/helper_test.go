package verbose

import (
	"context"
	"encoding/json"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

func data() *driver.GetResponseType {
	var body driver.GetResponseType
	// json.NewDecoder(apiResponse.Body).Decode(&body)
	plan, _ := ioutil.ReadFile("testdata/getresponse.json")
	json.Unmarshal(plan, &body)
	return &body
}
func verboseData() *StackVerbose {
	var body StackVerbose
	plan, _ := ioutil.ReadFile("testdata/verbosedata.json")
	json.Unmarshal(plan, &body)
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

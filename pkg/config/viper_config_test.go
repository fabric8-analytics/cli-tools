package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestViperUnMarshal(t *testing.T) {
	os.Setenv("HOST", "host")
	os.Setenv("AUTH_TOKEN", "token")
	os.Setenv("CONSENT_TELEMETRY", "false")
	os.Setenv("CRDA_KEY", "abc")
	want := &viperConfigs{
		Host:             "host",
		AuthToken:        "token",
		ConsentTelemetry: "false",
		CrdaKey:          "abc",
	}
	got := ViperUnMarshal()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Error in ViperUnMarshal func. (-want, +got):\n%s", diff)
	}
}

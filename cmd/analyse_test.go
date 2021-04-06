package cmd

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestGetAbsPath tests Absolute Path func.
func TestGetAbsPath(t *testing.T) {
	got := getAbsPath("testdata/requirements.txt")
	want, _ := filepath.Abs("testdata/requirements.txt")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Error in getAbsPath func. (-want, +got):\n%s", diff)
	}
}

package pypi

import (
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/rs/zerolog/log"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for Pypi
type Matcher struct {
	FilePath string
}

// Name implements driver.Matcher.
func (*Matcher) Name() string { return "python" }

// Filter implements driver.Filter.
func (*Matcher) Filter(ecosystem string) bool { return ecosystem == "python" }

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "pypi" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "pylist.json" }

// GetManifestFilePath sets file path
func (*Matcher) GetManifestFilePath(input string) string {
	path, err := filepath.Abs(input)
	if err != nil {
		log.Fatal().Msgf("Invalid Path of Manifest file. Only Absolute path is allowed.")
	}
	return path
}

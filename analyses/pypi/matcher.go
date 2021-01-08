package pypi

// Matcher implements driver.Matcher Interface for Pypi

import (
	"path/filepath"
	"strings"

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
func (*Matcher) Filter(ecosystem string) bool { return ecosystem == "pypi" }

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

// GeneratorDependencyTree creates pylist.json from requirements.txt
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate Pylist")
	m.getPylistGenerator()
	pathToPylist := m.buildDepsTree("generate_pylist.py", manifestFilePath)
	log.Debug().Msgf("Success: Generate Pylist")
	return pathToPylist
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	s := strings.Split(filename, ".")
	name, ext := s[0], s[1]
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		log.Debug().Msgf("Success: IsSupportedManifestFormat")
		return true
	}
	log.Debug().Msgf("Success: IsSupportedManifestFormat")
	return false
}

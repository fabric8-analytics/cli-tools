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

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "pypi" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "pylist.json" }

// GeneratorDependencyTree creates pylist.json from requirements.txt
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate Pylist")
	pylistGenerator := m.getPylistGenerator(filepath.Join("/tmp", "generate_pylist.py"))
	pathToPylist := m.buildDepsTree(pylistGenerator, manifestFilePath)
	log.Debug().Msgf("Success: Generate Pylist")
	return pathToPylist
}

// IsSupportedManifestFormat checks for Supported File Formats
func (*Matcher) IsSupportedManifestFormat(manifestFile string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	basename := filepath.Base(manifestFile)
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext)
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		log.Debug().Msgf("Success: IsSupportedManifestFormat")
		return true
	}
	log.Debug().Msgf("Success: IsSupportedManifestFormat")
	return false
}

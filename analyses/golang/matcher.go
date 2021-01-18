package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for Golang
type Matcher struct{}

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "golang" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "golist.json" }

// GeneratorDependencyTree creates golist.json from go.mod
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate golist.json")
	golang, err := exec.LookPath("go")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Golang is installed. Hint: Check same by executing: go version \n")
	}
	manifestDir := filepath.Dir(manifestFilePath)
	treePath, _ := filepath.Abs(filepath.Join(os.TempDir(), m.DepsTreeFileName()))
	genRepo := "github.com/fabric8-analytics/cli-tools/gomanifest"
	getGenerator := exec.Command(golang, "get", "-u", genRepo)
	dependencyTree := exec.Command(golang, "run", genRepo, manifestDir, treePath, golang)
	log.Debug().Msgf("Generator Command: %s", getGenerator)
	log.Debug().Msgf("dependencyTree Command: %s", dependencyTree)
	if err := getGenerator.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	if err := dependencyTree.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success: buildDepsTree")
	return treePath
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	basename := filepath.Base(filename)
	match, _ := regexp.MatchString("go.mod$", basename)
	log.Debug().Bool("regex", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

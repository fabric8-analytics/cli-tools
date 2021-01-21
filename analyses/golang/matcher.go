package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	gomanifest "github.com/fabric8-analytics/cli-tools/gomanifest/manifest"
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
	golang, err := exec.LookPath("go")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Golang is installed. Hint: Check same by executing: go version \n")
	}
	manifestDir := filepath.Dir(manifestFilePath)
	treePath, _ := filepath.Abs(filepath.Join(os.TempDir(), m.DepsTreeFileName()))
	gomanifest.Generate(golang, manifestDir, treePath)
	return treePath
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	match := strings.HasSuffix(filename, "go.mod")
	log.Debug().Bool("match", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

package golang

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
	gomanifest "github.com/fabric8-analytics/cli-tools/gomanifest/generator"
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

// generate is Generator for Golist.
func generate(goExePath string, goModPath string, goManifestPath string) error {
	// Start generating manifest data.
	cmd, err := gomanifest.RunGoList(goExePath, goModPath)
	if err != nil {
		return err
	}

	depPackages, err := gomanifest.GetDeps(cmd)
	if err != nil {
		return err
	}

	manifest := gomanifest.BuildManifest(&depPackages)
	// Create out file.
	f, err := os.Create(goManifestPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = manifest.Write(f)
	if err != nil {
		return err
	}
	return nil
}

// GeneratorDependencyTree creates golist.json from go.mod
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate golist.json")
	golang, err := exec.LookPath("go")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Golang is installed. Hint: Check same by executing: go version \n")
	}
	manifestDir := filepath.Dir(manifestFilePath)
	treePath, _ := filepath.Abs(filepath.Join(os.TempDir(), m.DepsTreeFileName()))
	generate(golang, manifestDir, treePath)
	log.Debug().Msgf("Success: Generate golist.json")
	return treePath
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	match := strings.HasSuffix(filename, "go.mod")
	log.Debug().Bool("match", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

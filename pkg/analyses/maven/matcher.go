package maven

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for Maven
type Matcher struct{}

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "maven" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "dependencies.txt" }

// GeneratorDependencyTree creates dependencies.txt from pom.xml
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate dependencies.txt")
	maven, err := exec.LookPath("mvn")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Maven is installed. Hint: Check same by executing: mvn --version \n")
	}
	treePath, _ := filepath.Abs(filepath.Join(os.TempDir(), m.DepsTreeFileName()))
	outcmd := fmt.Sprintf("-DoutputFile=%s", treePath)
	cleanRepo := exec.Command(maven, "--quiet", "clean", "-f", manifestFilePath)
	dependencyTree := exec.Command(maven, "--quiet", "org.apache.maven.plugins:maven-dependency-plugin:3.0.2:tree", "-f", manifestFilePath, outcmd, "-DoutputType=dot", "-DappendOutput=true")
	log.Debug().Msgf("Clean Repo Command: %s", cleanRepo)
	log.Debug().Msgf("dependencyTree Command: %s", dependencyTree)
	if err := cleanRepo.Run(); err != nil {
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
	match, _ := regexp.MatchString("pom.xml$", basename)
	log.Debug().Bool("regex", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

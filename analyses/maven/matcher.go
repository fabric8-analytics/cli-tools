package maven

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/rs/zerolog/log"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for Maven
type Matcher struct {
	FilePath string
}

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
	cmd := exec.Command(maven, "--quiet", "clean", "-f", manifestFilePath)
	cmd2 := exec.Command(maven, "--quiet", "org.apache.maven.plugins:maven-dependency-plugin:3.0.2:tree", "-f", manifestFilePath, outcmd, "-DoutputType=dot", "-DappendOutput=true")
	log.Debug().Msgf("Command: %s", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	if err := cmd2.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success: buildDepsTree")
	return treePath
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	basename := filepath.Base(filename)
	ext := filepath.Ext(basename)
	name := strings.TrimSuffix(basename, ext)
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		log.Debug().Msgf("Success: Manifest file is supported.")
		return true
	}
	log.Debug().Msgf("Success: Manifest file is not supported.")
	return false
}

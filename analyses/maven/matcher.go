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

// Filter implements driver.Filter.
func (*Matcher) Filter(ecosystem string) bool { return ecosystem == "maven" }

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
	// Command to execute dependencies.txt :
	// mvn --quiet clean -f "absolute/path/to/pom.xml" && mvn --quiet org.apache.maven.plugins:maven-dependency-plugin:3.0.2:tree
	// -f "absolute/path/to/pom.xml" -DoutputFile="output/path/to/dependencies.txt" -DoutputType=dot -DappendOutput=true

	outcmd := fmt.Sprintf("-DoutputFile=%s", m.DepsTreeFileName())
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
	depsList, _ := filepath.Abs(m.DepsTreeFileName())
	log.Debug().Msgf("Success: buildDepsTree")
	return depsList
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	s := strings.Split(filename, ".")
	name, ext := s[0], s[1]
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		log.Debug().Msgf("Success: Manifest file is supported.")
		return true
	}
	log.Debug().Msgf("Success: Manifest file is not supported.")
	return false
}

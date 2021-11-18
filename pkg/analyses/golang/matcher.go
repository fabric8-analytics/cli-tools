package golang

import (
	"errors"
	"fmt"
	"github.com/fabric8-analytics/cli-tools/pkg/utils"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	gomanifest "github.com/fabric8-analytics/cli-tools/gomanifest/generator"
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
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

func (m *Matcher) IgnoreVulnerabilities(manifestPath string) (map[string][]string, error) {

	//Ignore Vulnerabilities to be implemented for golang manifests
	log.Debug().Msgf("Extracting Packages and Vulnerabilities to Ignore.")
	ignoreVulnerabilities := make(map[string][]string)
	manifestFile, err := ioutil.ReadFile(manifestPath)

	if err != nil {
		return ignoreVulnerabilities, err
	}

	manifestFileContents := string(manifestFile) // convert contents of manifest file from []byte to a string

	//split the manifest file on the newline character to process each package, version, and vulnerabilities to ignore separately
	packagesWithVersions := strings.Split(manifestFileContents, "\n")
	for _, packageWithVersion := range packagesWithVersions {
		if strings.Contains(packageWithVersion, utils.CRDAIGNORE) { //check if crdaignore comment is present in the line
			packageWithVersion = strings.TrimSpace(packageWithVersion)
			pkgVersionIgnore := strings.Split(packageWithVersion, " ") //split package, version and crdaignore for further processing
			packageName := pkgVersionIgnore[0]                         // retrieve package from the slice after splitting
			var listOfVulnerabilities []string

			//if list of vulnerabilities are not provided, then return an empty list of vulnerabilities for that package
			if strings.Contains(pkgVersionIgnore[len(pkgVersionIgnore)-1], utils.CRDAIGNORE) {
				ignoreVulnerabilities[packageName] = listOfVulnerabilities
				continue
			}

			//check if crdaignore with vulnerabilities is in the right format
			if !(strings.Contains(packageWithVersion, "[") && (strings.Contains(packageWithVersion, "]"))) {
				err := errors.New("invalid 'crdaignore' format. Please enter vulnerabilities to ignore in the appropriate format")
				return ignoreVulnerabilities, err
			}

			//extract list of vulnerabilities between '[' and ']'
			vulnerabilitiesToIgnore := packageWithVersion[strings.Index(packageWithVersion, "[")+1 : strings.Index(packageWithVersion, "]")]
			vulnerabilitiesToIgnore = strings.ReplaceAll(vulnerabilitiesToIgnore, " ", "")
			fmt.Println(vulnerabilitiesToIgnore)
			listOfVulnerabilities = strings.Split(vulnerabilitiesToIgnore, ",")
			fmt.Println(listOfVulnerabilities)
			ignoreVulnerabilities[packageName] = listOfVulnerabilities
		}
	}
	return ignoreVulnerabilities, nil
}

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
	err = generate(golang, manifestDir, treePath)
	if err == nil {
		log.Debug().Msgf("Success: Generate golist.json")
	} else {
		log.Fatal().Err(err).Msg("Failed to Generate golist.json")
	}

	return treePath
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	match := strings.HasSuffix(filename, "go.mod")
	log.Debug().Bool("match", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

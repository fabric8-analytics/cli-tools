package pypi

// Matcher implements driver.Matcher Interface for Pypi

import (
	"errors"
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/pkg/utils"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for Pypi
type Matcher struct{}

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "pypi" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "pylist.json" }

func (m *Matcher) IgnoreVulnerabilities(manifestPath string) (map[string][]string, error) {

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
			pkgVersion := pkgVersionIgnore[0]                          // retrieve package and version in the format 'package==version' from the slice after splitting
			pkgVersionSplit := strings.Split(pkgVersion, "==")
			packageName := pkgVersionSplit[0]
			var listOfVulnerabilities []string

			//if list of vulnerabilities are not provided, then return an empty list of vulnerabilities for that package
			if strings.Contains(pkgVersionIgnore[len(pkgVersionIgnore)-1], utils.CRDAIGNORE) {
				ignoreVulnerabilities[packageName] = make([]string, 0)
				continue
			}

			if !(strings.Contains(packageWithVersion, "[") && (strings.Contains(packageWithVersion, "]"))) {
				err := errors.New("invalid 'crdaignore' format. Please enter vulnerabilities to ignore in the appropriate format")
				return ignoreVulnerabilities, err
			}

			vulnerabilitiesToIgnore := packageWithVersion[strings.Index(packageWithVersion, "[")+1 : strings.Index(packageWithVersion, "]")]
			vulnerabilitiesToIgnore = strings.ReplaceAll(vulnerabilitiesToIgnore, " ", "")
			listOfVulnerabilities = strings.Split(vulnerabilitiesToIgnore, ",")
			ignoreVulnerabilities[packageName] = listOfVulnerabilities
		}
	}
	return ignoreVulnerabilities, nil
}

// GeneratorDependencyTree creates pylist.json from requirements.txt
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate Pylist")
	pylistGenerator := m.getPylistGenerator(filepath.Join(os.TempDir(), "generate_pylist.py"))
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

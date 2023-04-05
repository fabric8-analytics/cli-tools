package npm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
)

var (
	_ driver.StackAnalysisInterface = (*Matcher)(nil)
)

// Matcher is State Object for NPM
type Matcher struct{}

// Ecosystem implements driver.Matcher.
func (*Matcher) Ecosystem() string { return "npm" }

// DepsTreeFileName implements driver.Matcher.
func (*Matcher) DepsTreeFileName() string { return "npmlist.json" }

// GeneratorDependencyTree creates npmlist.json from package.json
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate npmlist.json")
	npm, err := exec.LookPath("npm")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure `npm` is installed. Hint: Check same by executing: npm --version \n")
	}
	treePath := filepath.Join(os.TempDir(), m.DepsTreeFileName())
	prefix := fmt.Sprintf("--prefix=%s", filepath.Dir(manifestFilePath))
	npmList := exec.Command(npm, "list", prefix, "--prod", "--json")

	outfile, err := os.Create(treePath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating npmlist.json file.")
	}
	defer outfile.Close()
	npmList.Stdout = outfile

	log.Debug().Msgf("Dependency Tree Command: %s", npmList)

	var stderr bytes.Buffer
	npmList.Stderr = &stderr

	if err := npmList.Run(); err != nil {
		log.Debug().Msg("ERROR - Failed to Execute " + npmList.String() + "\n" + stderr.String())
		log.Fatal().Err(err).Msgf("Missing Dependencies. Hint: Please install the required dependencies with \"npm install\" from the directory of the manifest file")
	}
	_ = npmList.Wait()

	log.Debug().Msgf("Success: buildDepsTree at %s", treePath)
	return treePath
}

func (m *Matcher) IgnoreVulnerabilities(manifestPath string) (map[string][]string, error) {
	//Ignore Vulnerabilities for npm manifests
	log.Debug().Msgf("Extracting Packages and Vulnerabilities to Ignore.")
	manifestFile, err := ioutil.ReadFile(manifestPath)

	if err != nil {
		return nil, err
	}

	type CRDAIgnore struct {
		Ignore map[string]map[string][]string `json:"crdaignore,omitempty"`
	}

	var crdaIgnore CRDAIgnore
	if err := json.Unmarshal(manifestFile, &crdaIgnore); err != nil {
		err = errors.New("invalid 'crdaignore' format. Please enter packages and vulnerabilities to ignore in the right format")
		return nil, err
	}
	var ignoreAllVulnerabilityPackages []string
	for pkg, ignoreList := range crdaIgnore.Ignore["packages"] {
		if len(ignoreList) == 1 && ignoreList[0] == "*" {
			ignoreAllVulnerabilityPackages = append(ignoreAllVulnerabilityPackages, pkg)
		}
	}

	for _, pkg := range ignoreAllVulnerabilityPackages {
		crdaIgnore.Ignore["packages"][pkg] = make([]string, 0)
	}

	return crdaIgnore.Ignore["packages"], nil
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	match := strings.HasSuffix(filename, "package.json")
	log.Debug().Bool("match", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

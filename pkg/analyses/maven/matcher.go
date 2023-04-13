package maven

import (
	"bytes"
	"fmt"
	"github.com/fabric8-analytics/cli-tools/pkg/analyses/driver"
	"github.com/fabric8-analytics/cli-tools/pkg/utils"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
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

	var stdout bytes.Buffer
	dependencyTree.Stdout = &stdout

	log.Debug().Msgf("Clean Repo Command: %s", cleanRepo)
	log.Debug().Msgf("dependencyTree Command: %s", dependencyTree)
	if err := cleanRepo.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}

	if err := dependencyTree.Run(); err != nil {
		log.Debug().Err(err).Msg("ERROR - Failed to execute: " + dependencyTree.String() + "\n" + stdout.String())
		log.Fatal().Err(err).Msgf("Missing, or Unable to Resolve Certain Dependencies or Artifacts")
	}
	log.Debug().Msgf("Success: buildDepsTree")
	return treePath
}

func (m *Matcher) IgnoreVulnerabilities(manifestPath string) (map[string][]string, error) {
	log.Debug().Msgf("Extracting Packages and Vulnerabilities to Ignore.")
	ignoreVulnerabilities := make(map[string][]string)
	manifestFile, err := os.ReadFile(manifestPath)

	if err != nil {
		return ignoreVulnerabilities, err
	}

	manifestFileContents := string(manifestFile)               // convert contents of manifest file from []byte to a string
	r := regexp.MustCompile(`(?s)<dependency>.*</dependency>`) //Search for all valid dependencies
	matches := r.FindAllString(manifestFileContents, -1)
	dependencies := strings.Split(matches[0], "<dependency>") //split such that each dependencies details are separately acquired
	for _, dependency := range dependencies {
		if strings.Contains(dependency, utils.CRDAIGNORE) {
			groupIdRegEx := regexp.MustCompile(`<groupId>(.*)</groupId>`)
			groupIdRegExMatch := groupIdRegEx.FindAllStringSubmatch(dependency, -1)
			groupId := groupIdRegExMatch[0][1] //extract the group ID of the dependency

			artifactRegEx := regexp.MustCompile(`<artifactId>(.*)</artifactId>`)
			artifactRegExMatch := artifactRegEx.FindAllStringSubmatch(dependency, -1)
			artifactId := artifactRegExMatch[0][1] //extract the artifact ID of the dependency

			packageName := groupId + ":" + artifactId            //create package name to send with group ID and artifact ID
			dependencyDetails := strings.Split(dependency, "\n") //split to check which line 'crdaignore' is present in
			var listOfVulnerabilities []string
			for _, dependencyDetail := range dependencyDetails {
				if strings.Contains(dependencyDetail, utils.CRDAIGNORE) {
					if !(strings.Contains(dependencyDetail, "[") && strings.Contains(dependencyDetail, "]")) { //if list of vulnerabilities not entered, or
						//entered incorrectly, return an empty slice of vulnerabilities for that package
						ignoreVulnerabilities[packageName] = make([]string, 0)
						break
					}
					vulnerabilitiesToIgnore := dependencyDetail[strings.Index(dependencyDetail, "[")+1 : strings.Index(dependencyDetail, "]")] //extract list of vulnerabilities
					vulnerabilitiesToIgnore = strings.ReplaceAll(vulnerabilitiesToIgnore, " ", "")                                             //process to remove unnecessary empty spaces
					listOfVulnerabilities = strings.Split(vulnerabilitiesToIgnore, ",")                                                        //get vulnerabilities separately in a slice
					ignoreVulnerabilities[packageName] = listOfVulnerabilities
				}
			}
		}
	}
	return ignoreVulnerabilities, nil
}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	basename := filepath.Base(filename)
	match, _ := regexp.MatchString("pom.xml$", basename)
	log.Debug().Bool("regex", match).Str("path", filename).Msg("IsSupportedManifest")
	return match
}

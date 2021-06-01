package verbose

import (
	"errors"
	"fmt"
	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	"github.com/owenrumney/go-sarif/models"
	"github.com/owenrumney/go-sarif/sarif"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type RegexDependencyLocator struct {
	FileContent         string
	Ecosystem           string
	EndIndices          []int
	DependencyNodeIndex []int
}

func ProcessSarif(analysedResult *driver.GetResponseType, manifestFilePath string) (bool, error) {
	var hasVuln bool
	report, err := sarif.New(sarif.Version210)

	if err != nil {
		log.Fatal().Msg("Error forming SARIF respose")
		return false, errors.New("unable to create SARIF file")
	}

	run := report.AddRun("CRDA", "https://github.com/fabric8-analytics")

	regexDependencyLocator := RegexDependencyLocator{}
	if len(analysedResult.AnalysedDeps) == 0 {
		log.Fatal().Msg("Dependencies have not been analysed")
		return false, errors.New("dependencies have not been analysed")
	}
	err = regexDependencyLocator.SetUp(manifestFilePath, analysedResult.AnalysedDeps[0].Ecosystem)

	if err != nil {
		return false, errors.New("unable to setup dependency locator")
	}

	manifestParts := strings.Split(manifestFilePath, string(os.PathSeparator))
	manifest := manifestParts[len(manifestParts) - 1]
	for _, dep := range analysedResult.AnalysedDeps {
		line, column := regexDependencyLocator.LocateDependency(dep.Name)
		for _, publicVuln := range dep.PublicVulnerabilities {
			addVulnToReport(run, publicVuln, manifest, line, column)
			hasVuln = true
		}
		for _, privateVuln := range dep.PrivateVulnerabilities {
			addVulnToReport(run, privateVuln, manifest, line, column)
			hasVuln = true
		}
	}

	report.Write(os.Stdout)
	return hasVuln, nil
}

func addVulnToReport(run *models.Run, vuln driver.VulnerabilitiesType, manifestFilePath string, line int, column int) {
	rule := run.AddRule(vuln.ID).
		WithHelpURI(vuln.URL).WithDescription(vuln.Title)

	run.AddResult(rule.ID).
		WithMessage(vuln.Title).
		WithLevel(vuln.Severity).
		WithLocationDetails(manifestFilePath, line, column)
}


func (r *RegexDependencyLocator) SetUp(manifestFilePath string, ecosystem string) error{
	content, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		log.Fatal().Msg("Unable to load manifest File " + manifestFilePath)
		return fmt.Errorf("unable to load manifest file %s" ,manifestFilePath)
	}

	r.FileContent = string(content)
	r.Ecosystem = ecosystem
	newLineRegex, _ := regexp.Compile("\n")

	lines := newLineRegex.FindAllStringIndex(r.FileContent, -1)
	r.EndIndices = make([]int, len(lines))
	// Finding the end index for each end of line
	for i, line := range lines {
		r.EndIndices[i] = line[1]
	}

	dependenciesRegex, _ := regexp.Compile(getDependencyNodeRegex(r.Ecosystem))
	// Find the index for the start of the dependency node ( <dependencies> in case of pom.xml,
	// dependencies: in case of package.json)
	r.DependencyNodeIndex = dependenciesRegex.FindStringIndex(r.FileContent)

	return nil

}

func (r *RegexDependencyLocator) LocateDependency(dependency string) (int, int){
	// In case of maven the dependency consists groupId and artifactID
	// Picking up the artifact ID as the dependency
	if r.Ecosystem == "maven" {
		dependencyParts := strings.Split(dependency, ":")
		dependency = dependencyParts[len(dependencyParts) - 1]
	}
	// Adding the actual dependency in to the dependency regex
	dependencyRegexStr := strings.Replace(getDependencyRegex(r.Ecosystem), "?", dependency, 1)
	dependencyRegex, _ := regexp.Compile(dependencyRegexStr)
	dependencyIndex := dependencyRegex.FindStringIndex(r.FileContent)

	var lineNum int
	var column int

	// Check if the dependency index is within the dependency node
	if r.DependencyNodeIndex[0] < dependencyIndex[0] && dependencyIndex[0] < r.DependencyNodeIndex[1] {
		for i, val := range r.EndIndices {
			// Getting the line num and column number of the dependency
			if val <= dependencyIndex[0] && dependencyIndex[0] < r.EndIndices[i+1] {
				lineNum = i + 2
				column = dependencyIndex[0] - val + 2
				break
			}
		}
	}
	return lineNum, column
}

func getDependencyNodeRegex(ecosystem string) string{
	switch ecosystem {
	case "npm":
		return "\"dependencies\"[\\s]*:[\\s\n]*{[\\s\na-z\"\\d.,\\-:^]*}"
	case "maven":
		return "<dependencies[\\s\\n]*>[\\s]*[\\s\\n]*[<>!\\s\\w\\/${}\"\\d.,\\-:^]*<\\/dependencies[\\s\\n]*>"
	default:
		return "[\\s\\S]*"

	}
}

func getDependencyRegex(ecosystem string) string{
	switch ecosystem{
	case "npm":
		return "\"?\"[\\s]*:[\\s\n]*\"[\\d\\.^\\-\\w]*\""
	case "maven":
		return "?[\\s\\n]*<\\/artifactId[\\s\\n]*>"
	default:
		return "?[\\s]*=="

	}
}


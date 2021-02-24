package verbose

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

var cusColor = &CustomColors{
	Green: color.New(color.FgGreen, color.Bold).SprintFunc(),
	White: color.New(color.FgWhite, color.Bold).SprintFunc(),
	Cyan:  color.New(color.FgCyan).SprintFunc(),
	Red:   color.New(color.FgRed, color.Bold).SprintFunc(),
}

// ProcessVerbose processes verbose results and decides STDOUT format
func ProcessVerbose(analysedResult driver.GetResponseType, jsonOut bool) bool {
	out := getVerboseResult(analysedResult)
	if jsonOut {
		outputVerboseJSON(out)
	} else {
		outputVerbosePlain(out)
	}
	return out.TotalDirectVulnerabilities+out.TotalTransitiveVulnerabilities > 0
}

// getVerboseResult prepares verbose struct
func getVerboseResult(analysedResult driver.GetResponseType) *StackVerbose {
	data := processVulnerabilities(analysedResult.AnalysedDeps)
	out := &StackVerbose{
		Dependencies:                   data.AnalysedDependencies,
		TotalTransitives:               data.TotalTransitives,
		TotalDirectVulnerabilities:     data.TotalDirectVulnerabilities,
		TotalTransitiveVulnerabilities: data.TotalTransitiveVulnerabilities,
		Exploits:                       data.Severities,
	}
	return out
}

// processVulnerabilities prepares data for verbose STDOUT. It iterate over getResponse and picks up selective data for verbose STDOUT
func processVulnerabilities(analysedDeps []driver.AnalysedDepsType) ProcessVulnerabilities {
	processedData := &ProcessVulnerabilities{}
	for _, dep := range analysedDeps {
		dependency := getDependencyData(dep)
		dependency.CommonlyKnownVulnerabilities = getVulnerabilities(dep.PublicVulnerabilities)
		dependency.VulnerabilitiesUniqueToSynk = getVulnerabilities(dep.PrivateVulnerabilities)
		for _, trans := range dep.VulnerableDependencies {
			transitive := getDependencyData(trans)
			transitive.CommonlyKnownVulnerabilities = getVulnerabilities(trans.PublicVulnerabilities)
			transitive.VulnerabilitiesUniqueToSynk = getVulnerabilities(trans.PrivateVulnerabilities)
			dependency.VulnerableTransitives = append(dependency.VulnerableTransitives, transitive)
			processedData.Severities = getSeverity(trans.PublicVulnerabilities, processedData.Severities)
			processedData.Severities = getSeverity(trans.PrivateVulnerabilities, processedData.Severities)
			processedData.TotalTransitiveVulnerabilities += len(trans.PublicVulnerabilities) + len(trans.PrivateVulnerabilities)
		}

		processedData.AnalysedDependencies = append(processedData.AnalysedDependencies, dependency)
		processedData.TotalTransitives += len(dep.Transitives)

		processedData.TotalDirectVulnerabilities += len(dep.PublicVulnerabilities) + len(dep.PrivateVulnerabilities)
		processedData.Severities = getSeverity(dep.PublicVulnerabilities, processedData.Severities)
		processedData.Severities = getSeverity(dep.PrivateVulnerabilities, processedData.Severities)
	}
	return *processedData
}

// getDependencyData selects selective fields from Dependency Details for verbose STDOUT.
func getDependencyData(dep driver.AnalysedDepsType) DependenciesType {
	dependency := &DependenciesType{
		Name:               dep.Name,
		Version:            dep.Version,
		LatestVersion:      dep.LatestVersion,
		RecommendedVersion: dep.RecommendedVersion,
	}
	for _, trans := range dep.Transitives {
		transitive := &DependenciesType{
			Name:    trans.Name,
			Version: trans.Version,
		}
		dependency.Transitives = append(dependency.Transitives, *transitive)
	}
	return *dependency
}

//getVulnerabilities selects selective fields from Vulnerability details for verbose STDOUT.
func getVulnerabilities(vul []driver.VulnerabilitiesType) []VulnerabilityType {
	var allVuls []VulnerabilityType
	for _, v := range vul {
		myVul := &VulnerabilityType{
			ID:       v.ID,
			Severity: v.Severity,
			Title:    v.Title,
			URL:      v.URL,
		}
		allVuls = append(allVuls, *myVul)
	}
	return allVuls
}

// getSeverity prepares Severity struct
func getSeverity(vulnerability []driver.VulnerabilitiesType, severity SeverityType) SeverityType {
	for _, vul := range vulnerability {
		switch vul.Severity {
		case "critical":
			severity.Critical = append(severity.Critical, vul)
			break
		case "high":
			severity.High = append(severity.High, vul)
			break
		case "medium":
			severity.Medium = append(severity.Medium, vul)
			break
		case "low":
			severity.Low = append(severity.Low, vul)
			break
		}
	}
	return severity
}

// outputVerboseJSON STDOUT verbose output as JSON
func outputVerboseJSON(result *StackVerbose) {
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatal().Msg("Error forming CLI JSON Response.")
	}
	fmt.Fprintln(os.Stdout, string(b))
}

// outputVerbosePlain STDOUT Headers and Footers in Verbose Plain Output
func outputVerbosePlain(result *StackVerbose) {
	fmt.Fprint(os.Stdout, "Verbose Report for given Stack:\n\n")

	fmt.Fprint(os.Stdout,
		color.WhiteString("Scanned %d Dependencies and %d Transitives,", len(result.Dependencies), result.TotalTransitives),
		fmt.Sprintf(cusColor.Red(" Found %d Issues\n\n"), result.TotalDirectVulnerabilities+result.TotalTransitiveVulnerabilities),
	)
	fmt.Fprintln(os.Stdout, cusColor.Green("Fixable Issues:"))
	outputVulDeps(result.Dependencies)
	fmt.Fprint(os.Stdout, "(Powered by Snyk)\n\n")
}

// outputVulDeps STDOUT Vulnerable dependencies in verbose format
func outputVulDeps(deps []DependenciesType) {
	for _, dep := range deps {
		pkgName := fmt.Sprintf("%s@%s", cusColor.White(dep.Name), cusColor.White(dep.Version))

		if len(dep.CommonlyKnownVulnerabilities)+len(dep.VulnerabilitiesUniqueToSynk) > 0 {
			fmt.Fprint(os.Stdout,
				fmt.Sprintf("\tUpgrade %s ", pkgName),
				fmt.Sprintf("to %s@%s\n", cusColor.White(dep.Name), cusColor.White(dep.RecommendedVersion)),
			)
			dep.CommonlyKnownVulnerabilities = append(dep.CommonlyKnownVulnerabilities, dep.VulnerabilitiesUniqueToSynk...)
			outputVulType(dep.CommonlyKnownVulnerabilities, pkgName, pkgName)
		}
		if len(dep.VulnerableTransitives) > 0 {
			fmt.Fprint(os.Stdout,
				fmt.Sprintf(cusColor.Cyan("\n\t Issues in Transitives:\n")),
			)

			for _, trans := range dep.VulnerableTransitives {
				transName := fmt.Sprintf("%s@%s", cusColor.White(trans.Name), cusColor.White(trans.Version))
				fmt.Fprint(os.Stdout,
					fmt.Sprintf("\t \u2712 %s->%s\n", pkgName, transName),
				)
				trans.CommonlyKnownVulnerabilities = append(trans.CommonlyKnownVulnerabilities, trans.VulnerabilitiesUniqueToSynk...)
				outputVulType(trans.CommonlyKnownVulnerabilities, transName, pkgName)
				fmt.Fprintln(os.Stdout, "")
			}
		}
		fmt.Fprintln(os.Stdout, "")
	}
}

//outputVulType sorts vulnerabilty in order of severity and STDOUT
func outputVulType(allVuls []VulnerabilityType, pkg string, directDep string) {
	sort.Slice(allVuls, func(i, j int) bool {
		return allVuls[i].Severity < allVuls[j].Severity
	})
	for _, vul := range allVuls {
		attr1, severityText := getSeverityIntesity(vul.Severity)
		base := color.New(attr1).SprintFunc()
		bold := color.New(attr1, color.Bold).SprintFunc()
		fmt.Fprint(os.Stdout,
			fmt.Sprintf(bold("\t \u2718 %s"), vul.Title),
			fmt.Sprintf(base("[%s]"), severityText),
			fmt.Sprintf(" [%s] in %s\n", vul.URL, pkg),
			fmt.Sprintf("\t introduced by %s\n", directDep),
		)
	}
}

// getSeverityIntesity maps intensity of a Severity to Human Readable text and corresponding color.
func getSeverityIntesity(severity string) (color.Attribute, string) {
	myColor := color.FgHiCyan
	vulText := "Severity Unknown"
	switch severity {
	case "critical":
		myColor = color.FgHiRed
		vulText = "Critical Severity"
		break
	case "high":
		myColor = color.FgHiMagenta
		vulText = "High Severity"
		break
	case "medium":
		myColor = color.FgHiYellow
		vulText = "Medium Severity"
		break
	case "low":
		myColor = color.FgHiBlue
		vulText = "Low Severity"
		break
	}
	return myColor, vulText

}

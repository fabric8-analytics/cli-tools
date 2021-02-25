package summary

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
)

// ProcessSummary processes summary results, return true if Vul found
func ProcessSummary(analysedResult driver.GetResponseType, jsonOut bool, showVerboseMsg bool) bool {
	out := getResultSummary(analysedResult)
	if jsonOut {
		outputSummaryJSON(out)
	} else {
		outputSummaryPlain(out, showVerboseMsg)
	}
	return out.TotalVulnerabilities > 0
}

// GetResultSummary processes result Summary
func getResultSummary(analysedResult driver.GetResponseType) *StackSummary {
	totalDepsScanned := len(analysedResult.AnalysedDeps)
	data := processVulnerabilities(analysedResult.AnalysedDeps)
	out := &StackSummary{
		TotalScannedDependencies:           totalDepsScanned,
		TotalScannedTransitiveDependencies: data.TotalTransitives,
		TotalVulnerabilities:               data.PublicVul + data.PrivateVul,
		DirectVulnerableDependencies:       data.DirectVulnerableDependencies,
		CommonlyKnownVulnerabilities:       data.PublicVul,
		VulnerabilitiesUniqueToSynk:        data.PrivateVul,
		CriticalVulnerabilities:            data.Severities.Critical,
		HighVulnerabilities:                data.Severities.High,
		MediumVulnerabilities:              data.Severities.Medium,
		LowVulnerabilities:                 data.Severities.Low,
	}
	return out
}

// processVulnerabilities calculates Total Direct Public Vulnerabilities in Response
func processVulnerabilities(analysedDeps []driver.AnalysedDepsType) ProcessVulnerabilities {
	processedData := &ProcessVulnerabilities{}
	for _, dep := range analysedDeps {
		publicVul := len(dep.PublicVulnerabilities)
		privateVul := len(dep.PrivateVulnerabilities)
		if publicVul+privateVul > 0 {
			processedData.DirectVulnerableDependencies++
		}
		processedData.TotalTransitives += len(dep.Transitives)
		processedData.Severities = getSeverity(dep.PublicVulnerabilities, processedData.Severities)
		processedData.Severities = getSeverity(dep.PrivateVulnerabilities, processedData.Severities)
		for _, transVul := range dep.VulnerableDependencies {
			publicVul += len(transVul.PublicVulnerabilities)
			privateVul += len(transVul.PrivateVulnerabilities)
			processedData.Severities = getSeverity(transVul.PublicVulnerabilities, processedData.Severities)
			processedData.Severities = getSeverity(transVul.PrivateVulnerabilities, processedData.Severities)
		}
		processedData.PublicVul += publicVul
		processedData.PrivateVul += privateVul
	}
	return *processedData
}

// getSeverity calculates total severities in Vulnerabilities
func getSeverity(vulnerability []driver.VulnerabilitiesType, severity SeverityType) SeverityType {
	for _, vul := range vulnerability {
		switch vul.Severity {
		case "critical":
			severity.Critical++
			break
		case "high":
			severity.High++
			break
		case "medium":
			severity.Medium++
			break
		case "low":
			severity.Low++
			break
		}
	}
	return severity
}

// outputSummaryJSON stdout analyses summary output as JSON
func outputSummaryJSON(result *StackSummary) {
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatal().Msg("Error forming CLI JSON Response.")
	}
	fmt.Fprintln(os.Stdout, string(b))
}

// outputSummaryPlain stdout analyses summary output as JSON
func outputSummaryPlain(result *StackSummary, verboseMsg bool) {
	yellow := color.New(color.FgHiYellow, color.Bold).SprintFunc()
	white := color.New(color.FgHiWhite, color.Bold).SprintFunc()
	red := color.New(color.FgHiRed, color.Bold).SprintFunc()
	blue := color.New(color.FgHiBlue, color.Bold).SprintFunc()
	magenta := color.New(color.FgHiMagenta, color.Bold).SprintFunc()
	fmt.Fprint(os.Stdout, "Summary Report for Analyses:\n\n")
	fmt.Fprint(os.Stdout,
		white("Total Scanned Dependencies: "), white(result.TotalScannedDependencies), "\n",
		white("Total Scanned Transitive Dependencies: "), white(result.TotalScannedTransitiveDependencies), "\n",
		white("Direct Vulnerable Dependencies: "), white(result.DirectVulnerableDependencies), "\n",
		white("Total Vulnerabilities: "), white(result.TotalVulnerabilities), "\n",
		white("Commonly Known Vulnerabilities: "), white(result.CommonlyKnownVulnerabilities), "\n",
		white("Vulnerabilities Unique to Synk: "), white(result.VulnerabilitiesUniqueToSynk), "\n",
		red("Critical Vulnerabilities: "), red(result.CriticalVulnerabilities), "\n",
		magenta("High Vulnerabilities: "), magenta(result.HighVulnerabilities), "\n",
		yellow("Medium Vulnerabilities: "), yellow(result.MediumVulnerabilities), "\n",
		blue("Low Vulnerabilities: "), blue(result.LowVulnerabilities), "\n\n",
	)
	fmt.Fprint(os.Stdout, "(Powered by Snyk)\n\n")
	if verboseMsg {
		fmt.Fprint(os.Stdout, "tip: Register with Snyk and add token by running `crda auth`.\n")
		return
	}
	fmt.Fprint(os.Stdout, "tip: Use --verbose for detailed report.\n")
}

package tests

import (
	. "github.com/onsi/ginkgo"
)

// PrCheckSuite runs checks on every PR
func PrCheckSuite() {

	When("Run Crda version", TestCRDAVersion)
	When("Test CRDA auth command", TestCRDAauth)
	When("Test for invalid path throws error", TestInvalidPath)
	When("Test for invalid command", TestInvalidCommand)
	When("Test for invalid flag throws error", TestInvalidFlag)
	When("Run Crda help", TestCRDAHelp)
	When("Validate CRDA completion command by os", TestCRDACompletion)
	When("Validate there is a help page for all commands", TestCRDAallCommandsHelp)
	When("Run Crda analyse without any file", TestCRDAanalyseWithoutFile) 
	When("Test CRDA analyse Go", GolangTestSuitePR)
	When("Test CRDA analyse Maven", MavenTestSuitePR)
	When("Test CRDA analyse Npm", NpmTestSuitePR)
	When("Test CRDA analyse Pypi", PypiTestSuitePR) 

}

// NightlySuite runs checks on every PR
func NightlySuite() {

	When("Test CRDA auth command", TestCRDAauth)
	When("Test for invalid path throws error", TestInvalidPath)
	When("Test for invalid command", TestInvalidCommand)
	When("Test for invalid flag throws error", TestInvalidFlag)
	When("Run Crda version", TestCRDAVersion)
	When("Run Crda help", TestCRDAHelp)
	When("Validate CRDA completion command by os", TestCRDACompletion)
	When("Validate there is a help page for all commands", TestCRDAallCommandsHelp)
	When("Run Crda analyse without any file", TestCRDAanalyseWithoutFile)
	When("Test CRDA analyse Maven", MavenTestSuite)
	When("Test CRDA analyse Npm", NpmTestSuite)
	When("Test CRDA analyse Go", GolangTestSuite)
	When("Test CRDA analyse Pypi", PypiTestSuite)

}

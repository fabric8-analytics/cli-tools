package tests


import (
	. "github.com/onsi/ginkgo"
)


func PR_CHECK_SUITE(){

	When("Run Crda version", TestCRDA_version)
	//When("Test CRDA auth command", TestCRDA_auth)
	When("Test for invalid path throws error", TestInvalidPath)
	When("Test for invalid command", TestInvalidCommand)
	When("Test for invalid flag throws error", TestInvalidFlag)
	When("Run Crda help", TestCRDA_help)
	When("Validate CRDA completion command by os", TestCRDA_completion)
	When("Validate there is a help page for all commands", TestCRDA_all_commands_help)
	When("Run Crda analyse without any file", TestCRDA_analyse_without_file)
	When("Test CRDA analyse Go", GolangTestSuitePR)
	When("Test CRDA analyse Maven", MavenTestSuitePR)
	When("Test CRDA analyse Npm", NpmTestSuitePR)
	When("Test CRDA analyse Pypi", PypiTestSuitePR)

}




func NIGHTLY_SUITE(){

	When("Test CRDA auth command", TestCRDA_auth)
	When("Test CRDA analyse Go", GolangTestSuite)
	When("Test for invalid path throws error", TestInvalidPath)
	When("Test for invalid command", TestInvalidCommand)
	When("Test for invalid flag throws error", TestInvalidFlag)
	When("Run Crda version", TestCRDA_version)
	When("Run Crda help", TestCRDA_help)
	When("Validate CRDA completion command by os", TestCRDA_completion)
	When("Validate there is a help page for all commands", TestCRDA_all_commands_help)
	When("Run Crda analyse without any file", TestCRDA_analyse_without_file) 
    When("Test CRDA analyse Maven", MavenTestSuite)
	When("Test CRDA analyse Npm", NpmTestSuite)
	When("Test CRDA analyse Pypi", PypiTestSuite)

}
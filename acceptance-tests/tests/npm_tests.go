package tests

import (
	. "github.com/onsi/ginkgo"
)

func NpmTestSuitePR(){
	BeforeEach(func(){
		file = "/package.json"
		target = "/package.json"
	})
	When("I test analyse for npm with vulns", TestCRDA_analyse_npm)
	When("Test analyse with relative path", TestCRDA_analyse_with_relative_path_npm)

}

func NpmTestSuite(){
	BeforeEach(func(){
		file = "/package.json"
		target = "/package.json"
	})
	When("I test analyse for npm with vulns", TestCRDA_analyse_npm)
	When("I test analyse for npm with vulns with json", TestCRDA_analyse_npm_json)
	When("I test analyse for npm with vulns with verbose", TestCRDA_analyse_npm_verbose)
	When("I test analyse for npm with vulns with debug", TestCRDA_analyse_npm_debug)
	When("I test analyse for npm with vulns with all flags set true", TestCRDA_analyse_npm_all_flags)
	When("I test analyse for npm without vulns with json", TestCRDA_analyse_npm_json_no_vulns)
	When("Test analyse with relative path", TestCRDA_analyse_with_relative_path_npm)

}

func TestCRDA_analyse_npm() {
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem", Validate_analse)
	It("I should Cleanup", Cleanup_npm)
}

func TestCRDA_analyse_npm_json() {
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem with vulns", Validate_analse_json_vulns)
	It("I should Cleanup", Cleanup_npm)
}

func TestCRDA_analyse_npm_json_no_vulns() {
	BeforeEach(func(){
		file = "/vulns.json"
	})
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem without vulns", Validate_analse_json_no_vulns)
	It("I should Cleanup", Cleanup_npm)
}

func TestCRDA_analyse_npm_verbose() {
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem with vulns json and verbose", Validate_analse_vuln_verbose)
	It("I should Cleanup", Cleanup_npm)
}

 func TestCRDA_analyse_npm_debug() {
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem with vulns with debug", Validate_analse_vuln_debug)
	It("I should Cleanup", Cleanup_npm)
} 

func TestCRDA_analyse_npm_all_flags() {
	It("I should Get Absolute Path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
	It("I should Validate analyse for npm ecosystem with vulns and all flags set true", Validate_analse_all_flags)
	It("I should Cleanup", Cleanup_npm)
}


func TestCRDA_analyse_with_relative_path_npm() {
	BeforeEach(func() {
		file = "/package.json"
		target = "/package.json"
	})
	When("I Test for analyse command with relative path npm", func() {
		It("I should Copy Manifest", Copy_manifests)
		It("I should Install Dependencies to run Stack analyses", Install_npm_deps)
		It("Should be able to run analyse without error",  RunAnalyseRelative)
		It("I should Cleanup", Cleanup_npm)
	})
}
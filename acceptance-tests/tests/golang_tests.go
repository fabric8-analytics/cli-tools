package tests


import (
	. "github.com/onsi/ginkgo"
)

func BasicTestGO() {

	It("Should be able to get the absolute path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should be able to copy the main file", Copy_maninfile)
	It("I should able to run go mod tidy", Run_go_mod_tidy)
	It("I Should be able to run analyse without error", Validate_analse)
	It("I should perform cleanup", Cleanup_go)

}
func TestCRDA_analyse_with_relative_path_go() {
	When("I Test for analyse command with relative path npm", func() {
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("Should be able to run analyse without error",  RunAnalyseRelative)
		It("I should Cleanup", Cleanup_go)
	})
}

func GolangTestSuitePR() {
	BeforeEach(func() {
		file = "/go.mod.template"
		target = "/go.mod"
	})
	When("I test analyse command for Go with vulns", BasicTestGO)
	When("I test analyse command for Go relative path", TestCRDA_analyse_with_relative_path_go)

}

func GolangTestSuite() {
	BeforeEach(func() {
		file = "/go.mod.template"
		target = "/go.mod"
	})
	When("I test analyse command for Go with vulns", BasicTestGO)
	When("I test analyse command for Go relative path", TestCRDA_analyse_with_relative_path_go)
	When("I test analyse command for Go with vulns json", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse_json_vulns)
		It("I should perform cleanup", Cleanup_go)
	})
	When("I test analyse command for Go with vulns and verbose", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse_vuln_verbose)
		It("I should perform cleanup", Cleanup_go)

	})
	When("I test analyse command for Go with vulns and debug", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse_vuln_debug)
		It("I should perform cleanup", Cleanup_go)

	})
	When("I test analyse command for Go with vulns and all flags true", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse_all_flags)
		It("I should perform cleanup", Cleanup_go)
	})
	When("I test analyse command for Go without vulns", func() {
		BeforeEach(func() {
			file = "/go2.mod.template"
			go_main_file = "/main2.go.template"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse)
		It("I should perform cleanup", Cleanup_go)
	})
	When("I test analyse command for Go without vulns json", func() {
		BeforeEach(func() {
			file = "/go2.mod.template"
			go_main_file = "/main2.go.template"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should be able to copy the main file", Copy_maninfile)
		It("I should able to run go mod tidy", Run_go_mod_tidy)
		It("I Should be able to run analyse without error", Validate_analse_json_no_vulns)
		It("I should perform cleanup", Cleanup_go)
	})
}


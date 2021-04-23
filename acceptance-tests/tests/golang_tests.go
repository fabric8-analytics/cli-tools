package tests

import "github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"

// BasicTestGO tests Basic Go functionality
func BasicTestGO() {

	It("I should Copy Manifest", CopyManifests)
	It("I should be able to copy the main file", CopyManinfile)
	It("I should able to run go mod tidy", RunGoModTidy)
	It("I Should be able to run analyse without error", ValidateAnalse)
	It("I should perform cleanup", Cleanup)

}

// TestCRDAanalyseWithAbsolutePathGo runs Basic Absolute path test
func TestCRDAanalyseWithAbsolutePathGo() {
	When("I Test for analyse command with absolute path go", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("Should be able to run analyse without error", RunAnalyseAbsolute)
		It("I should Cleanup", Cleanup)
	})
}

// GolangTestSuitePR tests golang on each PR
func GolangTestSuitePR() {
	BeforeEach(func() {
		
		// file = "/go.mod.template"
		// target = "/go.mod"
		file, target = helper.CommonBeforeEach("/go.mod.template", "go")
	})
	When("I test analyse command for Go with vulns", BasicTestGO)
	When("I test analyse command for Go absolute path", TestCRDAanalyseWithAbsolutePathGo)

}

// GolangTestSuite works on the nightly test suite
func GolangTestSuite() {
	BeforeEach(func() {
		// file = "/go.mod.template"
		// target = "/go.mod"
		file, target = helper.CommonBeforeEach("/go.mod.template", "go")
	})
	When("I test analyse command for Go with vulns", BasicTestGO)
	When("I test analyse command for Go absolute path", TestCRDAanalyseWithAbsolutePathGo)
	When("I test analyse command for Go with vulns json", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseJSONVulns)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for Go with vulns and verbose", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseVulnVerbose)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Go with vulns and debug", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseVulnDebug)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Go with vulns and all flags true", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseAllFlags)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for Go without vulns", func() {
		BeforeEach(func() {
			// file = "/go2.mod.template"
			GoMainFile = "/main2.go.template"
			file, target = helper.CommonBeforeEach("/go2.mod.template", "go")
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseNoVulns)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for Go without vulns json", func() {
		BeforeEach(func() {
			// file = "/go2.mod.template"
			GoMainFile = "/main2.go.template"
			file, target = helper.CommonBeforeEach("/go2.mod.template", "go")
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should be able to copy the main file", CopyManinfile)
		It("I should able to run go mod tidy", RunGoModTidy)
		It("I Should be able to run analyse without error", ValidateAnalseJSONNoVulns)
		It("I should perform cleanup", Cleanup)
	})
}

package tests

import "github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"

// BasicTestMaven tests Basic maven functionality
func BasicTestMaven() {
	It("Copy Manifest", CopyManifests)
	It("Should be able to run analyse without error", ValidateAnalseNoVulns)
	It("I should perform cleanup", Cleanup)
}

func BasicTestMavenVulns() {
	BeforeEach(func() {
		file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
	})
	It("Copy Manifest", CopyManifests)
	It("Should be able to run analyse without error", ValidateAnalse)
	It("I should perform cleanup", Cleanup)
}

// TestCRDAanalyseWithAbsolutePathMvn tests with absolute path
func TestCRDAanalyseWithAbsolutePathMvn() {
	When("I Test for analyse command with absolute path maven", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", RunAnalyseAbsolute)
		It("I should Cleanup", Cleanup)
	})
}

// MavenTestSuitePR runs on each PR
func MavenTestSuitePR() {
	BeforeEach(func() {
		file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
	})
	When("I test analyse command for Maven with no vulns", BasicTestMaven)
	When("I test analyse command for Maven with vulns", BasicTestMavenVulns)
	When("I test analyse command for Maven absolute path", TestCRDAanalyseWithAbsolutePathMvn)

}

// MavenTestSuite runs on a nightly basis
func MavenTestSuite() {
	BeforeEach(func() {
		file, target = helper.CommonBeforeEach("/pom.xml", "maven")
	})
	When("I test analyse command for Maven with no vulns", BasicTestMaven)
	When("I test analyse command for Maven with vulns", BasicTestMavenVulns)
	When("I test analyse command for Maven absolute path", TestCRDAanalyseWithAbsolutePathMvn)
	When("I test analyse command for Maven with no vulns json", func() {
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalseJSONNoVulns)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for Maven with vulns", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalse)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Maven with vulns and json", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalseJSONVulns)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Maven with vulns and verbose", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalseVulnVerbose)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Maven with vulns and debug", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalseVulnDebug)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for Maven with vulns and all flags true", func() {
		BeforeEach(func() {
			file, target = helper.CommonBeforeEach("/pom2.xml", "maven")
		})
		It("Copy Manifest", CopyManifests)
		It("Should be able to run analyse without error", ValidateAnalseAllFlags)
		It("I should perform cleanup", Cleanup)

	})
}

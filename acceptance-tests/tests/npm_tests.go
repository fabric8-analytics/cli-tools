package tests


// NpmTestSuitePR runs on each PR check
func NpmTestSuitePR() {
	BeforeEach(func() {
		file = "/package.json"
		target = "/package.json"
	})
	When("I test analyse for npm with vulns", TestCRDAanalyseNpm)
	When("Test analyse with absolute path", TestCRDAanalyseWithAbsolutePathNpm)

}

// NpmTestSuite runs on a Nightly Check
func NpmTestSuite() {
	BeforeEach(func() {
		file = "/package.json"
		target = "/package.json"
	})
	When("I test analyse for npm with vulns", TestCRDAanalyseNpm)
	When("I test analyse for npm with vulns with json", TestCRDAanalyseNpmJSON)
	When("I test analyse for npm with vulns with verbose", TestCRDAanalyseNpmVerbose)
	When("I test analyse for npm with vulns with debug", TestCRDAanalyseNpmDebug)
	When("I test analyse for npm with vulns with all flags set true", TestCRDAanalyseNpmAllFlags)
	When("I test analyse for npm without vulns with json", TestCRDAanalyseNpmJSONNoVulns)
	When("Test analyse with absolute path", TestCRDAanalyseWithAbsolutePathNpm)

}

// TestCRDAanalyseNpm tests Basic npm functionality
func TestCRDAanalyseNpm() {
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem", ValidateAnalse)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseNpmJSON tests functionality with json
func TestCRDAanalyseNpmJSON() {
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem with vulns", ValidateAnalseJSONVulns)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseNpmJSONNoVulns tests functionality with json and no vulns
func TestCRDAanalyseNpmJSONNoVulns() {
	BeforeEach(func() {
		file = "/vulns.json"
	})
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem without vulns", ValidateAnalseJSONNoVulns)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseNpmVerbose tests functionality with verbose
func TestCRDAanalyseNpmVerbose() {
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem with vulns json and verbose", ValidateAnalseVulnVerbose)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseNpmDebug tests functionality with debug
func TestCRDAanalyseNpmDebug() {
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem with vulns with debug", ValidateAnalseVulnDebug)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseNpmAllFlags tests functionality with all flags set true
func TestCRDAanalyseNpmAllFlags() {
	It("I should Copy Manifest", CopyManifests)
	It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
	It("I should Validate analyse for npm ecosystem with vulns and all flags set true", ValidateAnalseAllFlags)
	It("I should Cleanup", Cleanup)
}

// TestCRDAanalyseWithAbsolutePathNpm tests functionality with abs path
func TestCRDAanalyseWithAbsolutePathNpm() {
	BeforeEach(func() {
		file = "/package.json"
		target = "/package.json"
	})
	When("I Test for analyse command with absolute path npm", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", CopyManifests)
		It("I should Install Dependencies to run Stack analyses", InstallNpmDeps)
		It("Should be able to run analyse without error", RunAnalyseAbsolute)
		It("I should Cleanup", Cleanup)
	})
}

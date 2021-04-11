package tests

// BasicTestPypi tests Basic functionality
func BasicTestPypi() {
	It("I should Copy Manifest", CopyManifests)
	It("I should able to run pip install", RunPipInstall)
	It("I Should be able to run analyse without error", ValidateAnalse)
	It("I should perform cleanup", Cleanup)

}

// TestCRDAanalyseWithAbsolutePathPypi tests functionality with abs path
func TestCRDAanalyseWithAbsolutePathPypi() {
	When("I Test for analyse command with Absolute path pypi", func() {
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall) // While Running local replace the occurence of RunPipInstall with InitVirtualEnv func
		It("Should be able to run analyse without error", RunAnalyseAbsolute)
		It("I should Cleanup", Cleanup)
	})
}

// PypiTestSuitePR runs on each PR check
func PypiTestSuitePR() {
	BeforeEach(func() {
		file = "/requirements.txt"
		target = "/requirements.txt"
	})
	When("I test analyse command for pypi with vulns", BasicTestPypi)
	When("I test analyse command for pypi absolute path", TestCRDAanalyseWithAbsolutePathPypi)

}

// PypiTestSuite runs on a nightly basis
func PypiTestSuite() {
	BeforeEach(func() {
		file = "/requirements.txt"
		target = "/requirements.txt"
	})
	When("I test analyse command for pypi with vulns", BasicTestPypi)
	When("I test analyse command for pypi absolute path", TestCRDAanalyseWithAbsolutePathPypi)
	When("I test analyse command for pypi with vulns json", func() {
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalseJSONVulns)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for pypi with vulns and verbose", func() {
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalseVulnVerbose)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for pypi with vulns and debug", func() {
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalseVulnDebug)
		It("I should perform cleanup", Cleanup)

	})
	When("I test analyse command for pypi with vulns and all flags true", func() {
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalseAllFlags)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for pypi without vulns", func() {
		BeforeEach(func() {
			file = "/requirements2.txt"
		})
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalse)
		It("I should perform cleanup", Cleanup)
	})
	When("I test analyse command for pypi without vulns json", func() {
		BeforeEach(func() {
			file = "/requirements2.txt"
		})
		It("I should Copy Manifest", CopyManifests)
		It("I should able to run pip install", RunPipInstall)
		It("I Should be able to run analyse without error", ValidateAnalseJSONNoVulns)
		It("I should perform cleanup", Cleanup)
	})

}

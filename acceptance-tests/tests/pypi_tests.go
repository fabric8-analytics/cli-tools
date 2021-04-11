package tests


import (
	. "github.com/onsi/ginkgo"
)

func BasicTestPypi() {
	It("Should be able to get the absolute path", GetAbsPath)
	It("I should Copy Manifest", Copy_manifests)
	It("I should able to run pip install", Init_virtual_env)
	It("I Should be able to run analyse without error", Validate_analse_pypi)
	It("I should perform cleanup", Cleanup_pypi)

}

func TestCRDA_analyse_with_relative_path_pypi() {
	When("I Test for analyse command with relative path npm", func() {
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("Should be able to run analyse without error",  RunAnalyseRelative)
		It("I should Cleanup", Cleanup_pypi)
	})
}

func PypiTestSuitePR() {
	BeforeEach(func() {
		file = "/requirements.txt"
		target = "/requirements.txt"
	})
	When("I test analyse command for pypi with vulns", BasicTestPypi)
	When("I test analyse command for pypi relative path",TestCRDA_analyse_with_relative_path_pypi)

}

func PypiTestSuite() {
	BeforeEach(func() {
		file = "/requirements.txt"
		target = "/requirements.txt"
	})
	When("I test analyse command for pypi with vulns", BasicTestPypi)
	When("I test analyse command for pypi relative path",TestCRDA_analyse_with_relative_path_pypi)
	When("I test analyse command for pypi with vulns json", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi_json_vulns)
		It("I should perform cleanup", Cleanup_pypi)
	})
	When("I test analyse command for pypi with vulns and verbose", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi_vuln_verbose)
		It("I should perform cleanup", Cleanup_pypi)

	})
	When("I test analyse command for pypi with vulns and debug", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi_vuln_debug)
		It("I should perform cleanup", Cleanup_pypi)

	})
	When("I test analyse command for pypi with vulns and all flags true", func() {
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi_all_flags)
		It("I should perform cleanup", Cleanup_pypi)
	})
	When("I test analyse command for pypi without vulns", func() {
		BeforeEach(func() {
			file = "/requirements2.txt"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi)
		It("I should perform cleanup", Cleanup_pypi)
	})
	When("I test analyse command for pypi without vulns json", func() {
		BeforeEach(func() {
			file = "/requirements2.txt"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("I should Copy Manifest", Copy_manifests)
		It("I should able to run pip install", Init_virtual_env)
		It("I Should be able to run analyse without error", Validate_analse_pypi_json_no_vulns)
		It("I should perform cleanup", Cleanup_pypi)
	})

}


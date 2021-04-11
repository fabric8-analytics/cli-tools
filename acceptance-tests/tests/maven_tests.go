package tests

import (
	. "github.com/onsi/ginkgo"
)

func BasicTestMaven(){
	It("Should be able to get the absolute path", GetAbsPath)
	It("Copy Manifest", Copy_manifests)
	It("Should be able to run analyse without error", Validate_analse)
	It("I should perform cleanup", Cleanup_mvn)
}

func TestCRDA_analyse_with_relative_path_mvn() {
	When("I Test for analyse command with relative path npm", func() {
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error",  RunAnalyseRelative)
		It("I should Cleanup", Cleanup_mvn)
	})
}

func MavenTestSuitePR() {
	BeforeEach(func(){
		file = "/pom.xml"
		target = "/pom.xml"
	})
	When("I test analyse command for Maven with no vulns", BasicTestMaven)
	When("I test analyse command for Maven relative path", TestCRDA_analyse_with_relative_path_mvn)

}

func MavenTestSuite(){
	BeforeEach(func(){
		file = "/pom.xml"
		target = "/pom.xml"
	})
	When("I test analyse command for Maven with no vulns", BasicTestMaven)
	When("I test analyse command for Maven relative path", TestCRDA_analyse_with_relative_path_mvn)
	When("I test analyse command for Maven with no vulns json", func(){
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse_json_no_vulns)
		It("I should perform cleanup", Cleanup_mvn)
	})
	When("I test analyse command for Maven with vulns", func(){
		BeforeEach(func(){
			file = "/pom2.xml"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse)
		It("I should perform cleanup", Cleanup_mvn)

	})
	When("I test analyse command for Maven with vulns and json", func(){
		BeforeEach(func(){
			file = "/pom2.xml"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse_json_vulns)
		It("I should perform cleanup", Cleanup_mvn)

	})
	When("I test analyse command for Maven with vulns and verbose", func(){
		BeforeEach(func(){
			file = "/pom2.xml"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse_vuln_verbose)
		It("I should perform cleanup", Cleanup_mvn)

	})
	When("I test analyse command for Maven with vulns and debug", func(){
		BeforeEach(func(){
			file = "/pom2.xml"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse_vuln_debug)
		It("I should perform cleanup", Cleanup_mvn)

	})
	When("I test analyse command for Maven with vulns and all flags true", func(){
		BeforeEach(func(){ 
			file = "/pom2.xml"
		})
		It("Should be able to get the absolute path", GetAbsPath)
		It("Copy Manifest", Copy_manifests)
		It("Should be able to run analyse without error", Validate_analse_all_flags)
		It("I should perform cleanup", Cleanup_mvn)

	})
}


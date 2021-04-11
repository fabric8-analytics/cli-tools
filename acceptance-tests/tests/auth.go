package tests

import (
	"os"
	"os/exec"
	acclog "github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
)

// TestCRDAauth implements Test Cases to test auth crda command
func TestCRDAauth() {
	When("I run crda auth without snyk token", func() {
		It("should throw error", func() {
			cmd := exec.Command("./crda", "auth", "snyk-token")
			stdout, err := cmd.Output()
			acclog.InfoLogger.Println(string(stdout))
			Expect(err).To(HaveOccurred())
		})
	})
	When("I run crda auth with invalid snyk token", func() {
		It("it should throw error", func() {
			cmd := exec.Command("./crda", "auth", "snyk-token", "invalid-token")
			stdout, err := cmd.Output()
			acclog.InfoLogger.Println(string(stdout))
			Expect(err).To(HaveOccurred())
		})

	})
	When("I run crda auth with valid snyk token", func() {
		validToken := os.Getenv("snyk_token")
		It("it should not throw error", func() {
			cmd := exec.Command("./crda", "auth", "--snyk-token", string(validToken))
			stdout, err := cmd.Output()
			acclog.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())
		})

	})
}

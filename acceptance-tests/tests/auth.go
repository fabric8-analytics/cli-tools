package tests


import (
	"os"
	"os/exec"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCRDA_auth() {
	When("I run crda auth without snyk token", func() {
		It("should throw error", func() {
			cmd := exec.Command("./crda", "auth", "snyk-token")
			stdout, err := cmd.Output()
			acc_log.InfoLogger.Println(string(stdout))
			Expect(err).To(HaveOccurred())
		})
	})
	When("I run crda auth with invalid snyk token", func() {
		It("it should throw error", func() {
			cmd := exec.Command("./crda", "auth", "snyk-token", "invalid-token")
			stdout, err := cmd.Output()
			acc_log.InfoLogger.Println(string(stdout))
			Expect(err).To(HaveOccurred())
		})

	})
	When("I run crda auth with valid snyk token", func() {
		valid_token := os.Getenv("snyk_token")
		It("it should not throw error", func() {
			cmd := exec.Command("./crda", "auth", "--snyk-token", string(valid_token))
			stdout, err := cmd.Output()
			acc_log.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())
		})

	})
}

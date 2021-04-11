package acceptance_tests_test

import (
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	acclog "github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AcceptanceTests Suite")
}

var _ = BeforeSuite(func() {
	acclog.Initlog()
	helper.CreateDataDir()
	helper.CheckforSynkToken()
})

var _ = AfterSuite(func() {
	helper.CleanupSuite()
})
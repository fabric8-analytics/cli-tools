package acceptance_tests_test

import (
	"testing"
	"time"

	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	acclog "github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AcceptanceTests Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(1 * time.Minute)
	acclog.Initlog()
	helper.CreateDataDir()
	helper.CheckforSynkToken()
})

var _ = AfterSuite(func() {
	helper.CleanupSuite()
})
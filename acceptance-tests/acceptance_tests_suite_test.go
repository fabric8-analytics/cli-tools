package acceptance_tests_test

import (
	"testing"
	"time"

	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AcceptanceTests Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(2 * time.Minute)
	helper.CreateDataDir()
	helper.CheckforSynkToken()
})

var _ = AfterSuite(func() {
	helper.CleanupSuite()
})
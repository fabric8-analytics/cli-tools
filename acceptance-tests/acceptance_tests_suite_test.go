package acceptance_tests_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	
)

func TestAcceptanceTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AcceptanceTests Suite")
}


var _ = BeforeSuite(func() {
	acc_log.Init_log()
	helper_cli.CreateData_dir()
})

var _ = AfterSuite(func() {
	helper_cli.Cleanup_suite()
})
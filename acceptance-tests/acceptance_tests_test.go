package acceptance_tests_test

import (
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/tests"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("AcceptanceTests", func() {

	Describe("PR ACCEPTANCE TESTS", tests.PrCheckSuite)

	Describe("NIGHTLY SUITE", tests.NightlySuite)

})

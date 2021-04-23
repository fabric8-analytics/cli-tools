package tests

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	acclog "github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	//"github.com/onsi/gomega/gbytes"
	//"github.com/onsi/gomega/gexec"
)

// Done Declarations for Ginkgo DSL
type Done ginkgo.Done

// Benchmarker for ginkgo
type Benchmarker ginkgo.Benchmarker

// GinkgoWriter this are all ginkgo vars to be used in tests
var (
	GinkgoWriter                          = ginkgo.GinkgoWriter
	GinkgoRandomSeed                      = ginkgo.GinkgoRandomSeed
	GinkgoParallelNode                    = ginkgo.GinkgoParallelNode
	GinkgoT                               = ginkgo.GinkgoT
	CurrentGinkgoTestDescription          = ginkgo.CurrentGinkgoTestDescription
	RunSpecs                              = ginkgo.RunSpecs
	RunSpecsWithDefaultAndCustomReporters = ginkgo.RunSpecsWithDefaultAndCustomReporters
	RunSpecsWithCustomReporters           = ginkgo.RunSpecsWithCustomReporters
	Skip                                  = ginkgo.Skip
	Fail                                  = ginkgo.Fail
	GinkgoRecover                         = ginkgo.GinkgoRecover
	Describe                              = ginkgo.Describe
	FDescribe                             = ginkgo.FDescribe
	PDescribe                             = ginkgo.PDescribe
	XDescribe                             = ginkgo.XDescribe
	Context                               = ginkgo.Context
	FContext                              = ginkgo.FContext
	PContext                              = ginkgo.PContext
	XContext                              = ginkgo.XContext
	When                                  = ginkgo.When
	FWhen                                 = ginkgo.FWhen
	PWhen                                 = ginkgo.PWhen
	XWhen                                 = ginkgo.XWhen
	It                                    = ginkgo.It
	FIt                                   = ginkgo.FIt
	PIt                                   = ginkgo.PIt
	XIt                                   = ginkgo.XIt
	Specify                               = ginkgo.Specify
	FSpecify                              = ginkgo.FSpecify
	PSpecify                              = ginkgo.PSpecify
	XSpecify                              = ginkgo.XSpecify
	By                                    = ginkgo.By
	Measure                               = ginkgo.Measure
	FMeasure                              = ginkgo.FMeasure
	PMeasure                              = ginkgo.PMeasure
	XMeasure                              = ginkgo.XMeasure
	BeforeSuite                           = ginkgo.BeforeSuite
	AfterSuite                            = ginkgo.AfterSuite
	SynchronizedBeforeSuite               = ginkgo.SynchronizedBeforeSuite
	SynchronizedAfterSuite                = ginkgo.SynchronizedAfterSuite
	BeforeEach                            = ginkgo.BeforeEach
	JustBeforeEach                        = ginkgo.JustBeforeEach
	JustAfterEach                         = ginkgo.JustAfterEach
	AfterEach                             = ginkgo.AfterEach
	RegisterFailHandler                   = gomega.RegisterFailHandler
	RegisterFailHandlerWithT              = gomega.RegisterFailHandlerWithT
	RegisterTestingT                      = gomega.RegisterTestingT
	InterceptGomegaFailures               = gomega.InterceptGomegaFailures
	Ω                                     = gomega.Ω
	Expect                                = gomega.Expect
	ExpectWithOffset                      = gomega.ExpectWithOffset
	Eventually                            = gomega.Eventually
	EventuallyWithOffset                  = gomega.EventuallyWithOffset
	Consistently                          = gomega.Consistently
	ConsistentlyWithOffset                = gomega.ConsistentlyWithOffset
	SetDefaultEventuallyTimeout           = gomega.SetDefaultEventuallyTimeout
	SetDefaultEventuallyPollingInterval   = gomega.SetDefaultEventuallyPollingInterval
	SetDefaultConsistentlyDuration        = gomega.SetDefaultConsistentlyDuration
	SetDefaultConsistentlyPollingInterval = gomega.SetDefaultConsistentlyPollingInterval
	NewWithT                              = gomega.NewWithT
	NewGomegaWithT                        = gomega.NewGomegaWithT
	Default                               = gomega.Default
	Equal                = gomega.Equal
	BeEquivalentTo       = gomega.BeEquivalentTo
	BeIdenticalTo        = gomega.BeIdenticalTo
	BeNil                = gomega.BeNil
	BeTrue               = gomega.BeTrue
	BeFalse              = gomega.BeFalse
	HaveOccurred         = gomega.HaveOccurred
	Succeed              = gomega.Succeed
	MatchError           = gomega.MatchError
	BeClosed             = gomega.BeClosed
	Receive              = gomega.Receive
	BeSent               = gomega.BeSent
	MatchRegexp          = gomega.MatchRegexp
	ContainSubstring     = gomega.ContainSubstring
	HavePrefix           = gomega.HavePrefix
	HaveSuffix           = gomega.HaveSuffix
	MatchJSON            = gomega.MatchJSON
	MatchXML             = gomega.MatchXML
	MatchYAML            = gomega.MatchYAML
	BeEmpty              = gomega.BeEmpty
	HaveLen              = gomega.HaveLen
	HaveCap              = gomega.HaveCap
	BeZero               = gomega.BeZero
	ContainElement       = gomega.ContainElement
	BeElementOf          = gomega.BeElementOf
	ConsistOf            = gomega.ConsistOf
	ContainElements      = gomega.ContainElements
	HaveKey              = gomega.HaveKey
	HaveKeyWithValue     = gomega.HaveKeyWithValue
	BeNumerically        = gomega.BeNumerically
	BeTemporally         = gomega.BeTemporally
	BeAssignableToTypeOf = gomega.BeAssignableToTypeOf
	Panic                = gomega.Panic
	PanicWith            = gomega.PanicWith
	BeAnExistingFile     = gomega.BeAnExistingFile
	BeARegularFile       = gomega.BeARegularFile
	BeADirectory         = gomega.BeADirectory
	HaveHTTPStatus       = gomega.HaveHTTPStatus
	And                  = gomega.And
	SatisfyAll           = gomega.SatisfyAll
	Or                   = gomega.Or
	SatisfyAny           = gomega.SatisfyAny
	Not                  = gomega.Not
	WithTransform        = gomega.WithTransform
)

// Path and all the other global vars
var (
	Path        string = "/data"
	pwd         string
	err         error
	manifests   string = "/manifests"
	file        string = "/package.json"
	target      string = "/package.json"
	GoMainFile  string = "/main.go.template"
	GoMainFileT string = "/main.go"
)

// getCRDAcmd returns the command according to os
func getCRDAcmd() string {
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return "./crda"
	} else if runtime.GOOS == "windows" {
		return "./crda.exe"
	}
	return "crda"
}

// GetAbsPath returns absolute path
func GetAbsPath() {
	pwd, err = helper.Getabspath(Path)
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// InitVirtualEnv makes a new virtual env for python
func InitVirtualEnv() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; python3 -m venv env; source env/bin/activate; pip3 install -r requirements.txt;")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// InstallNpmDeps runs npm install
func InstallNpmDeps() {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// RunGoModTidy runs the tidy command
func RunGoModTidy() {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// RunPipInstall runs pip install command
func RunPipInstall() {
	cmd := exec.Command("pip", "install", "-r", "requirements.txt")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// ValidateAnalse runs analyse command
func ValidateAnalse() {
	session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse","data"+target)
	fmt.Println(GinkgoWriter, string(session))
}

// ValidateAnalseJSONVulns runs analyse command with json
func ValidateAnalseJSONVulns() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "data"+target, "-j")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acclog.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

// ValidateAnalseJSONNoVulns runs analyse command with json no vulns
func ValidateAnalseJSONNoVulns() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "data"+target, "-j")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acclog.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

// ValidateAnalseVulnVerbose runs analyse command with verbose
func ValidateAnalseVulnVerbose() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "data"+target, "-v")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())

}

// ValidateInvalidFilePath runs analyse command with invalid file path
func ValidateInvalidFilePath() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "/package.json", "-v")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	Expect(err).To(HaveOccurred())
}

// ValidateInvalidCommand runs invalid command
func ValidateInvalidCommand() {
	cmd := exec.Command(getCRDAcmd(), "analysess")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())
}

// ValidateInvalidFlag runs analyse command with invalid flag
func ValidateInvalidFlag() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "-ghghd")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	Expect(err).To(HaveOccurred())
	e := err.(*exec.ExitError)
	Expect(e.ExitCode()).To(Equal(1))
}

// ValidateAnalseVulnDebug runs analyse command with debug
func ValidateAnalseVulnDebug() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "data"+target, "-d")
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			acclog.InfoLogger.Printf(scanner.Text())
		}
		done <- true
	}()
	cmd.Start()
	<-done
	err = cmd.Wait()
	Expect(err).NotTo(HaveOccurred())
}

// ValidateAnalseAllFlags runs analyse command with all flags set true
func ValidateAnalseAllFlags() {
	cmd := exec.Command(getCRDAcmd(), "analyse", "data"+target, "-v", "-j", "-d")
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acclog.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

// Cleanup cleans the data dir 
func Cleanup(){
	err := helper.Cleanup("data/*")
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
}

// CopyManifests copies manifests to data dir
func CopyManifests() {
	dir1, err := helper.Getabspath(manifests)
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
	dir2, err := helper.Getabspath(Path)
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
	acclog.InfoLogger.Println(dir1 + file)
	acclog.InfoLogger.Println(dir2 + target)
	err = helper.CopyContentstoTarget("manifests"+file, "data"+target)
	Expect(err).NotTo(HaveOccurred())
}

// CopyManinfile copies go main file to target
func CopyManinfile() {
	dir1, err := helper.Getabspath(manifests)
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
	dir2, err := helper.Getabspath(Path)
	if err != nil {
		acclog.ErrorLogger.Println(err)
	}
	acclog.InfoLogger.Println(dir1 + GoMainFile)
	acclog.InfoLogger.Println(dir2 + GoMainFileT)
	err = helper.CopyContentstoTarget("manifests"+GoMainFile, "data"+GoMainFileT)
	Expect(err).NotTo(HaveOccurred())
}

// RunAnalyseAbsolute runs analyse with abs path
func RunAnalyseAbsolute() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(getCRDAcmd(), "analyse", "data"+target)
	} else {
		cmd = exec.Command(getCRDAcmd(), "analyse", pwd+target)
	}
	stdout, err := cmd.Output()
	acclog.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acclog.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

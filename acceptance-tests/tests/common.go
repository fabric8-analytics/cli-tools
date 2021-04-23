package tests

import (
	"os/exec"
	"runtime"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	. "github.com/onsi/gomega"
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
		helper.PrintWithGinkgo(err.Error())
	}
}

// InitVirtualEnv makes a new virtual env for python
func InitVirtualEnv() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; python3 -m venv env; source env/bin/activate; pip3 install -r requirements.txt;")
	stdout, err := cmd.Output()
	helper.PrintWithGinkgo(string(stdout))
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
}

// InstallNpmDeps runs npm install
func InstallNpmDeps() {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	helper.PrintWithGinkgo(string(stdout))
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
}

// RunGoModTidy runs the tidy command
func RunGoModTidy() {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	helper.PrintWithGinkgo(string(stdout))
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
}

// RunPipInstall runs pip install command
func RunPipInstall() {
	cmd := exec.Command("pip", "install", "-r", "requirements.txt")
	cmd.Dir = "data"
	stdout, err := cmd.Output()
	helper.PrintWithGinkgo(string(stdout))
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
}

// ValidateAnalse runs analyse command
func ValidateAnalse() {
	session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse","data"+target)
	helper.PrintWithGinkgo(session)
}

// ValidateAnalseNoVulns runs analyse command
func ValidateAnalseNoVulns() {
	session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse","data"+target)
	helper.PrintWithGinkgo(session)
}

// ValidateAnalseJSONVulns runs analyse command with json
func ValidateAnalseJSONVulns() {
	session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse","data"+target, "-j")
	helper.PrintWithGinkgo(session)

}

// ValidateAnalseJSONNoVulns runs analyse command with json no vulns
func ValidateAnalseJSONNoVulns() {
	session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse","data"+target, "-j")
	helper.PrintWithGinkgo(session)


}

// ValidateAnalseVulnVerbose runs analyse command with verbose
func ValidateAnalseVulnVerbose() {
	session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse","data"+target, "-v")
	helper.PrintWithGinkgo(session)

}

// ValidateInvalidFilePath runs analyse command with invalid file path
func ValidateInvalidFilePath() {
	session := helper.CmdShouldFailWithExit1(getCRDAcmd(), "analyse","/package.json")
	helper.PrintWithGinkgo(session)
	Expect(string(session)).To(ContainSubstring("invalid file path: /package.json"))
}

// ValidateInvalidCommand runs invalid command
func ValidateInvalidCommand() {
	session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyseess")
	helper.PrintWithGinkgo(session)
}

// ValidateInvalidFlag runs analyse command with invalid flag
func ValidateInvalidFlag() {
	session := helper.CmdShouldFailWithExit1(getCRDAcmd(), "analyse","-y")
	helper.PrintWithGinkgo(session)
	Expect(string(session)).To(ContainSubstring("unknown shorthand flag: 'y' in -y"))
}

// ValidateAnalseVulnDebug runs analyse command with debug
func ValidateAnalseVulnDebug() {
	session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse","data"+target, "-d")
	helper.PrintWithGinkgo(session)
}

// ValidateAnalseAllFlags runs analyse command with all flags set true
func ValidateAnalseAllFlags() {
	session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse", "data"+target, "-v", "-j", "-d")
	helper.PrintWithGinkgo(session)
}

// Cleanup cleans the data dir 
func Cleanup(){
	err := helper.Cleanup("data/*")
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
}

// CopyManifests copies manifests to data dir
func CopyManifests() {
	dir1, err := helper.Getabspath(manifests)
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
	dir2, err := helper.Getabspath(Path)
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
	helper.PrintWithGinkgo(dir1+ file)
	helper.PrintWithGinkgo(dir2 + target)
	err = helper.CopyContentstoTarget("manifests"+file, "data"+target)
	Expect(err).NotTo(HaveOccurred())
}

// CopyManinfile copies go main file to target
func CopyManinfile() {
	dir1, err := helper.Getabspath(manifests)
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
	dir2, err := helper.Getabspath(Path)
	if err != nil {
		helper.PrintWithGinkgo(err.Error())
	}
	helper.PrintWithGinkgo(dir1 + GoMainFile)
	helper.PrintWithGinkgo(dir2 + GoMainFileT)
	err = helper.CopyContentstoTarget("manifests"+GoMainFile, "data"+GoMainFileT)
	Expect(err).NotTo(HaveOccurred())
}

// RunAnalyseAbsolute runs analyse with abs path
func RunAnalyseAbsolute() {
	if runtime.GOOS == "windows" {
		session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse", "data"+target)
		helper.PrintWithGinkgo(session)
	} else {
		session := helper.CmdShouldPassWithExit2(getCRDAcmd(), "analyse", pwd+target)
		helper.PrintWithGinkgo(session)
	}
	

}

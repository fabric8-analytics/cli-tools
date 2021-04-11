package tests

import (
	"os/exec"
	"runtime"

	acclog "github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	
)

// TestCRDAVersion checks for version command
func TestCRDAVersion() {

	It("Runs and Validate CLI version", func() {
		cmd := exec.Command(getCRDAcmd(), "version")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})

}

// TestInvalidPath checks for invalid path error
func TestInvalidPath() {
	It("Should throw error if i send invalid file path", ValidateInvalidFilePath)
}

// TestInvalidCommand checks for invalid sub command
func TestInvalidCommand() {
	It("Should throw error when run an invalid command", ValidateInvalidCommand)
}

// TestInvalidFlag checks for an invalid flag
func TestInvalidFlag() {
	It("Should throw an error when set an invalid flag", ValidateInvalidFlag)
}

// TestCRDAHelp verifies the help command
func TestCRDAHelp() {
	It("Runs and Validate Help command", func() {
		cmd := exec.Command(getCRDAcmd(), "help")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})

}

// TestCRDACompletion verifies the completion command
func TestCRDACompletion() {
	It("Runs and Validate completion command", func() {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			cmd := exec.Command(getCRDAcmd(), "completion", "bash")
			stdout, err := cmd.Output()
			acclog.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())
		} else if runtime.GOOS == "windows" {
			cmd := exec.Command(getCRDAcmd(), "completion", "powershell")
			stdout, err := cmd.Output()
			acclog.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())

		} else {
			Skip("No supporting operating system")
		}
	})
}

// TestCRDAallCommandsHelp verifies if there is a help page for all sub commands
func TestCRDAallCommandsHelp() {
	It("analyse command has help page", func() {
		cmd := exec.Command(getCRDAcmd(), "analyse", "--help")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())
	})
	It("auth command has help page", func() {
		cmd := exec.Command(getCRDAcmd(), "auth", "--help")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("completion command has help page", func() {
		cmd := exec.Command(getCRDAcmd(), "completion", "--help")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("version command has help page", func() {
		cmd := exec.Command(getCRDAcmd(), "version", "--help")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("help command has help page", func() {
		cmd := exec.Command(getCRDAcmd(), "help", "bash")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
}

// TestCRDAanalyseWithoutFile veifies error when no file is provided
func TestCRDAanalyseWithoutFile() {
	It("Validate analyse without flile throws error", func() {
		cmd := exec.Command(getCRDAcmd(), "analyse")
		stdout, err := cmd.Output()
		acclog.InfoLogger.Println(string(stdout))
		Expect(err).To(HaveOccurred())

	})
}

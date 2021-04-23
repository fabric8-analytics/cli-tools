package tests

import (
	"fmt"
	"runtime"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
)

// TestCRDAVersion checks for version command
func TestCRDAVersion() {

	It("Runs and Validate CLI version", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse", "version")
		fmt.Println(GinkgoWriter, string(session))

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
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse", "help")
		fmt.Println(GinkgoWriter, string(session))

	})

}

// TestCRDACompletion verifies the completion command
func TestCRDACompletion() {
	It("Runs and Validate completion command", func() {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			_ = helper.CmdShouldPassWithoutError(getCRDAcmd(), "completion", "bash")
		} else if runtime.GOOS == "windows" {
			_ = helper.CmdShouldPassWithoutError(getCRDAcmd(), "completion", "powershell")
		} else {
			Skip("No supporting operating system")
		}
	})
}

// TestCRDAallCommandsHelp verifies if there is a help page for all sub commands
func TestCRDAallCommandsHelp() {
	It("analyse command has help page", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "analyse", "--help")
		fmt.Println(GinkgoWriter, string(session))
	})
	It("auth command has help page", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "auth", "--help")
		fmt.Println(GinkgoWriter, string(session))

	})
	It("completion command has help page", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "completion", "--help")
		fmt.Println(GinkgoWriter, string(session))

	})
	It("version command has help page", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "version", "--help")
		fmt.Println(GinkgoWriter, string(session))

	})
	It("help command has help page", func() {
		session := helper.CmdShouldPassWithoutError(getCRDAcmd(), "help", "--help")
		fmt.Println(GinkgoWriter, string(session))

	})
}

// TestCRDAanalyseWithoutFile veifies error when no file is provided
func TestCRDAanalyseWithoutFile() {
	It("Validate analyse without flile throws error", func() {
		session := helper.CmdShouldFailWithExit1(getCRDAcmd(), "analyse")
		fmt.Println(GinkgoWriter, session)
		Expect(string(session)).To(ContainSubstring("requires valid manifest file path"))

	})
}

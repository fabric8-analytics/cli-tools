package tests

import (
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
    . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os/exec"
	"runtime"
	
)

func TestCRDA_version() {

	It("Runs and Validate CLI version", func() {
		cmd := exec.Command("./crda", "version")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})

}

func TestInvalidPath() {
	It("Should throw error if i send invalid file path", Validate_invalid_file_path)
}

func TestInvalidCommand() {
	It("Should throw error when run an invalid command", Validate_invalid_command)
}

func TestInvalidFlag() {
	It("Should throw an error when set an invalid flag", Validate_invalid_flag)
}

func TestCRDA_help() {
	It("Runs and Validate Help command", func() {
		cmd := exec.Command("./crda", "help")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})

}

func TestCRDA_completion() {
	It("Runs and Validate completion command", func() {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
			cmd := exec.Command("./crda", "completion", "bash")
			stdout, err := cmd.Output()
			acc_log.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())
		} else if runtime.GOOS == "windows" {
			cmd := exec.Command("./crda", "completion", "powershell")
			stdout, err := cmd.Output()
			acc_log.InfoLogger.Println(string(stdout))
			Expect(err).NotTo(HaveOccurred())

		} else {
			Skip("No supporting operating system")
		}
	})
}

func TestCRDA_all_commands_help() {
	It("analyse command has help page", func() {
		cmd := exec.Command("./crda", "analyse", "--help")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())
	})
	It("auth command has help page", func() {
		cmd := exec.Command("./crda", "auth", "--help")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("completion command has help page", func() {
		cmd := exec.Command("./crda", "completion", "--help")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("version command has help page", func() {
		cmd := exec.Command("./crda", "version", "--help")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
	It("help command has help page", func() {
		cmd := exec.Command("./crda", "help", "bash")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).NotTo(HaveOccurred())

	})
}

func TestCRDA_analyse_without_file() {
	It("Validate analyse without flile throws error", func() {
		cmd := exec.Command("./crda", "analyse")
		stdout, err := cmd.Output()
		acc_log.InfoLogger.Println(string(stdout))
		Expect(err).To(HaveOccurred())

	})
}



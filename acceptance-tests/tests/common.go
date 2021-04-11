package tests

import (
	"bufio"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/helper"
	"github.com/fabric8-analytics/cli-tools/acceptance-tests/log"
	. "github.com/onsi/gomega"
	"os/exec"
)

var (
	Path           string = "/data"
	pwd            string
	err            error
	manifests      string = "/manifests"
	file           string = "/package.json"
	target         string = "/package.json"
	go_main_file   string = "/main.go.template"
	go_main_file_t string = "/main.go"
)

func GetAbsPath() {
	pwd, err = helper_cli.Get_abs_path(Path)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Init_virtual_env() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; python3 -m venv env; source env/bin/activate; pip3 install -r requirements.txt;")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func GetPath() {
	if pwd == "" {
		GetAbsPath()
	}
}

func Install_npm_deps() {
	GetPath()
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; npm i;")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Run_go_mod_tidy() {
	GetPath()
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; go mod tidy;")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Run_pip_install() {
	GetPath()
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; pip install -r requirements.txt;")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Validate_analse() {
	cmd := exec.Command("crda", "analyse", pwd+target)
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_pypi() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target)
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_pypi_json_vulns() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target+" -j")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}
func Validate_analse_pypi_json_no_vulns() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target+" -j")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())

}
func Validate_analse_pypi_vuln_verbose() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target+" -v")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_pypi_vuln_debug() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target+" -d")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_pypi_all_flags() {
	cmd := exec.Command("/bin/sh", "-c", "cd "+pwd+"; source env/bin/activate; ./crda"+" analyse "+pwd+target+" -j -d -v")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_json_vulns() {
	cmd := exec.Command("./crda", "analyse", pwd+target, "-j")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_json_no_vulns() {
	cmd := exec.Command("./crda", "analyse", pwd+target, "-j")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Validate_analse_vuln_verbose() {
	cmd := exec.Command("./crda", "analyse", pwd+target, "-v")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())

}

func Validate_invalid_file_path() {
	cmd := exec.Command("./crda", "analyse", "/package.json", "-v")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).To(HaveOccurred())
}

func Validate_invalid_command() {
	cmd := exec.Command("./crda", "analysess")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())
}

func Validate_invalid_flag() {
	cmd := exec.Command("./crda", "analyse", "-ghghd")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).To(HaveOccurred())
	e := err.(*exec.ExitError)
	Expect(e.ExitCode()).To(Equal(1))
}

func Validate_analse_vuln_debug() {
	cmd := exec.Command("./crda", "analyse", pwd+target, "-d")
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			acc_log.InfoLogger.Printf(scanner.Text())
		}
		done <- true
	}()
	cmd.Start()
	<-done
	err = cmd.Wait()
	Expect(err).NotTo(HaveOccurred())
}

func Validate_analse_all_flags() {
	cmd := exec.Command("./crda", "analyse", pwd+target, "-v", "-j", "-d")
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	e := err.(*exec.ExitError)
	acc_log.InfoLogger.Println(e.ExitCode())
	Expect(e.ExitCode()).To(Equal(2))

}

func Cleanup_npm() {
	GetPath()
	err := helper_cli.Cleanup_npm(pwd)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Cleanup_mvn() {
	GetPath()
	err := helper_cli.Cleanup_mvn(pwd)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Cleanup_go() {
	GetPath()
	err := helper_cli.Cleanup_go(pwd)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}
func Cleanup_pypi() {
	GetPath()
	err := helper_cli.Cleanup_pypi(pwd)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
}

func Copy_manifests() {
	dir1, err := helper_cli.Get_abs_path(manifests)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
	dir2, err := helper_cli.Get_abs_path(Path)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
	acc_log.InfoLogger.Println(dir1 + file)
	acc_log.InfoLogger.Println(dir2 + target)
	err = helper_cli.Copy_contents_to_target(dir1+file, dir2+target)
	Expect(err).NotTo(HaveOccurred())
}

func Copy_maninfile() {
	dir1, err := helper_cli.Get_abs_path(manifests)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
	dir2, err := helper_cli.Get_abs_path(Path)
	if err != nil {
		acc_log.ErrorLogger.Println(err)
	}
	acc_log.InfoLogger.Println(dir1 + go_main_file)
	acc_log.InfoLogger.Println(dir2 + go_main_file_t)
	err = helper_cli.Copy_contents_to_target(dir1+go_main_file, dir2+go_main_file_t)
	Expect(err).NotTo(HaveOccurred())
}

func RunAnalyseRelative() {
	
	cmd := exec.Command("./crda", "analyse", "data"+target)
	stdout, err := cmd.Output()
	acc_log.InfoLogger.Println(string(stdout))
	Expect(err).NotTo(HaveOccurred())

}

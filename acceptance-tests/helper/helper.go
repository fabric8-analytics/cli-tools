package helper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func runningCmd(cmd *exec.Cmd) string {
	prog := filepath.Base(cmd.Path)
	return fmt.Sprintf("Running %s with args %v", prog, cmd.Args)
}

func CmdRunner(program string, args ...string) *gexec.Session {
	//prefix ginkgo verbose output with program name
	prefix := fmt.Sprintf("[%s] ", filepath.Base(program))
	prefixWriter := gexec.NewPrefixedWriter(prefix, GinkgoWriter)
	command := exec.Command(program, args...)
	fmt.Fprintln(GinkgoWriter, runningCmd(command))
	session, err := gexec.Start(command, prefixWriter, prefixWriter)
	Expect(err).NotTo(HaveOccurred())
	return session
}

func CmdShouldPassWithExit2(program string, args ...string) string {
	session := CmdRunner(program, args...)
	Eventually(session).Should(gexec.Exit(2), runningCmd(session.Command))
	return string(session.Out.Contents())
}

func CmdShouldPassWithExit1(program string, args ...string) string {
	session := CmdRunner(program, args...)
	Eventually(session).Should(gexec.Exit(1), runningCmd(session.Command))
	return string(session.Out.Contents())
}

func CmdShouldFailWithExit1(program string, args ...string) string {
	session := CmdRunner(program, args...)
	Eventually(session).Should(gexec.Exit(1), runningCmd(session.Command))
	return string(session.Err.Contents())
}

func CmdShouldPassWithoutError(program string, args ...string) string {
	session := CmdRunner(program, args...)
	Eventually(session).ShouldNot(gexec.Exit(1),runningCmd(session.Command))
	Eventually(session).ShouldNot(gexec.Exit(2),runningCmd(session.Command))
	Eventually(session).Should(gexec.Exit(), runningCmd(session.Command))
	return string(session.Out.Contents())
}

// Getabspath Gets the Absolute Path
func Getabspath(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty file path supplied")
	}
	cmd := exec.Command("pwd")
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.New("cannot run PWD")
	}
	pwd := string(stdout)
	pwd = strings.TrimSpace(pwd)
	result := pwd + path
	return result, nil

}

// CommonBeforeEach is a common before each function
func CommonBeforeEach(file string, target string) (string, string) {
	
	if target == "npm" {
		return file, "/package.json"
	}else if target == "pypi"{
		return file, "/requirements.txt"
	}else if target == "go"{
		return file, "/go.mod"
	}else if target == "maven"{
		return file, "/pom.xml"
	}else {
		return "", ""
	}
}

// CheckforSynkToken Checks for Snyk Token var in env
func CheckforSynkToken() {
	_, present := os.LookupEnv("snyk_token")
	if !present {
		fmt.Printf("snyk token env variable not present")
		os.Exit(1)
	} else {
		fmt.Println("Snyk Token var present.. Starting Test Suite")
	}
}

// CreateDataDir creates data directory
func CreateDataDir() error {
	err := os.Mkdir("data", 0755)
	if err != nil {
		return err
	}
	return nil
}

// CleanupSuite deletes files used by tests
func CleanupSuite() error {
	err := os.RemoveAll("/data")
	if err != nil {
		return err
	}
	return nil
}

// Cleanup global cleanup function
func Cleanup(path string) error {
	contents, err := filepath.Glob(path)
    if err != nil {
        return err
    }
    for _, item := range contents {
        err = os.RemoveAll(item)
        if err != nil {
            return err
        }
    }
    return nil
}

// CopyContentstoTarget Copies the files to target
func CopyContentstoTarget(filename string, target string) error {
	from, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

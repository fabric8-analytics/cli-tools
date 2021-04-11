package helper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
)

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

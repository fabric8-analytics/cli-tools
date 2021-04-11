package helper_cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Get_abs_path(filepath string) (string, error) {
	if filepath == "" {
		return "", errors.New("empty file path supplied")
	}
	cmd := exec.Command("pwd")
	stdout, err := cmd.Output()
	if err != nil {
		return "", errors.New("cannot run PWD")
	}
	pwd := string(stdout)
	pwd = strings.TrimSpace(pwd)
	result := pwd + filepath
	return result, nil

}

func Cleanup_npm(filepath string) error {
	fmt.Println(filepath)
	err := os.RemoveAll(filepath + "/node_modules/")
	err1 := os.Remove(filepath + "/package-lock.json")
	err2 := os.Remove(filepath + "/package.json")
	if err != nil || err1 != nil || err2 != nil{
		return err
	}
	return nil
}

func CreateData_dir() error {
	err := os.Mkdir("data", 0755)
	if err != nil {
		return err
	}
	return nil
}

func Cleanup_suite() error {
	err := os.RemoveAll("/data")
	if err != nil {
		return err
	}
	return nil
}

func Cleanup_mvn(filepath string) error {
	fmt.Println(filepath)
	err := os.Remove(filepath + "/pom.xml")
	if err != nil{
		return err
	}
	return nil
}

func Cleanup_go(filepath string) error {
	fmt.Println(filepath)
	err := os.Remove(filepath + "/go.sum")
	err1 := os.Remove(filepath + "/go.mod")
	err2 := os.Remove(filepath + "/main.go")
	if err != nil || err1 != nil || err2 != nil{
		return err
	}
	return nil
}

func Cleanup_pypi(filepath string) error {
	fmt.Println(filepath)
	err := os.Remove(filepath + "/requirements.txt")
	err1 := os.RemoveAll(filepath + "/env/")
	if err != nil||err1 != nil{
		return err
	}
	return nil
}


func Copy_contents_to_target(filename string, target string) error {
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



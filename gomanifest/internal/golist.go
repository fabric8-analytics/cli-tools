package internal

import (
	"io"
	"os/exec"
)

// GoListCmd ... Go list command structure.
type GoListCmd struct {
	cmd    *exec.Cmd
	output io.ReadCloser
}

// RunGoList ... Actual function that executes go list command and returns output as string.
func RunGoList(cwd string, goExec string) (*GoListCmd, error) {
	goList := exec.Command(goExec, "list", "-json", "-deps", "-mod=readonly", "./...")
	goList.Dir = cwd
	output, err := goList.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = goList.Start(); err != nil {
		return nil, err
	}
	return &GoListCmd{cmd: goList, output: output}, nil
}

// ReadCloser implements internal.GoList
func (list *GoListCmd) ReadCloser() io.ReadCloser {
	return list.output
}

// Wait implements internal.GoList
func (list *GoListCmd) Wait() error {
	return list.cmd.Wait()
}

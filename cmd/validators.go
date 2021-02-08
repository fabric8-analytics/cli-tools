package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func checkFileExists(file string) bool {
	isAbsPath := filepath.IsAbs(file)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return isAbsPath
}

func validateFileArg(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires absolute file path	")
	}
	if checkFileExists(args[0]) {
		return nil
	}
	return fmt.Errorf("Invalid file path: %s", args[0])
}

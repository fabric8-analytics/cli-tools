package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func validateFileArg(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires valid manifest file path")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return fmt.Errorf("invalid file path: %s", args[0])
	}
	return nil
}

func validateFlagValues(flag string) error {
	validValues := []string{"jenkins", "terminal", "tekton", "gh-actions", "intellij"}
	for _, item := range validValues {
		if item == flag {
			return nil
		}
	}
	return fmt.Errorf("invalid client flag value. Please select one out of %v", validValues)
}

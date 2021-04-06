package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func validateFileArg(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Requires valid manifest file path.")
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Invalid file path: %s", args[0]))
	}
	return nil
}

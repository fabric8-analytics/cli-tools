package cmd

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config SUBCOMMAND [flags]",
	Short: "Modify crda configuration",
	Long:  `Modify crda configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

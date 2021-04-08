package cmd

import (
	"fmt"

	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get binary version.",
	Long:  `Command to output version of a binary.`,
	Run:   getVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion(_ *cobra.Command, args []string) {
	version := sa.BuildVersion()
	fmt.Println(version)
}

package cmd

import (
	"fmt"

	"github.com/fabric8-analytics/cli-tools/pkg/version"
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
	version := version.BuildVersion()
	fmt.Println(version)
}

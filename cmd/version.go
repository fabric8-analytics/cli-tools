package cmd

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get binary version.",
	Long:  `Command to output version of a binary.`,
	Run:   getVersion,
}

// Version string populated by releaser Job
var (
	Version string = "0.0.0"
	// commitHash contains the current Git revision.
	CommitHash string = "abcd"
	// Timestamp contains the UnixTimestamp of the binary build.
	Timestamp string = "0"
	// VendorInfo contains vendor notes about the current build.
	VendorInfo string = "Local Build"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion(cmd *cobra.Command, args []string) {

	version := "v" + Version
	if CommitHash != "" {
		version += "-" + CommitHash
	}
	osArch := runtime.GOOS + "/" + runtime.GOARCH

	versionString := fmt.Sprintf("%s %s ",
		version, osArch)

	if Timestamp != "" {
		i, err := strconv.ParseInt(Timestamp, 10, 64)
		if err == nil {
			tm := time.Unix(i, 0)
			versionString += fmt.Sprintf(" BuildDate: %s", tm.Format(time.RFC1123))
		}
	}
	if VendorInfo != "" {
		versionString += fmt.Sprintf("  Vendor: %s", VendorInfo)
	}

	fmt.Println(versionString)
}

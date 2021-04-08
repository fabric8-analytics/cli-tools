package version

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// Version string populated by releaser Job
var (
	// The current version of crda cli
	version string = "0.0.0"

	// commitHash contains the current Git revision.
	commitHash string = "abcd"

	// Timestamp contains the UnixTimestamp of the binary build.
	timestamp string = "0"

	// VendorInfo contains vendor notes about the current build.
	vendorInfo string = "Local Build"
)

// GetCRDAVersion returns CRDA CLI Version
func GetCRDAVersion() string {
	return version
}

// BuildVersion builds version string for command output
func BuildVersion() string {
	version := "v" + version
	version += "-" + commitHash
	osArch := runtime.GOOS + "/" + runtime.GOARCH

	versionString := fmt.Sprintf("%s %s ",
		version, osArch)

	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err == nil {
		tm := time.Unix(i, 0)
		versionString += fmt.Sprintf(" BuildDate: %s", tm.Format(time.RFC1123))
	}
	versionString += fmt.Sprintf("  Vendor: %s", vendorInfo)
	return versionString
}

package stackanalyses

import (
	"github.com/fabric8-analytics/cli-tools/pkgs/version"
)

// BuildVersion builds version string and passed to caller
func BuildVersion() string {
	return version.BuildVersion()
}

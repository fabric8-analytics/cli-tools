package maven

// This File contains Utility functions of Maven Ecosystem

import "regexp"

// checkName checks for valid file name.
func checkName(name string) bool {
	match1, _ := regexp.MatchString("pom?", name)
	match2, _ := regexp.MatchString("maven?", name)
	if match1 || match2 {
		return true
	}
	return false
}

// checkExt checks for valid file extension.
func checkExt(ext string) bool {
	switch ext {
	case
		"xml",
		"txt":
		return true
	}
	return false
}

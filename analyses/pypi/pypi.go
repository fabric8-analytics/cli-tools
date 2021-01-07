package pypi

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var pypiSupportedFormats = []string{
	"*-requirements.txt",
	"requirements-*.txt",
	"requirements.txt",
	"constraints-*.txt",
	"*-constraints.txt",
	"*-requirements.in",
	"requirements-*.in"}

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	s := strings.Split(filename, ".")
	name, ext := s[0], s[1]
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		return true
	}
	return false
}

// checkName checks for valid file name.
func checkName(name string) bool {
	match1, _ := regexp.MatchString("requirements?", name)
	match2, _ := regexp.MatchString("constraints?", name)
	fmt.Println(match1)
	if match1 || match2 {
		return true
	}
	return false
}

// checkExt checks for valid file extension.
func checkExt(ext string) bool {
	switch ext {
	case
		"in",
		"txt":
		return true
	}
	return false
}

// GeneratorDependencyTree creates pylist.json from requirements.txt
func (m *Matcher) GeneratorDependencyTree(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate Pylist")
	pathToPylist := m.buildDepsTree("analyses/pypi/generate_pylist.py", manifestFilePath)
	log.Debug().Msgf("Success: Generate Pylist")
	return pathToPylist
}

// buildDepsTree generates final Deps Tree and saves it as pylist.json
func (m *Matcher) buildDepsTree(generatePylist string, manifestFilePath string) string {
	log.Debug().Msgf("Execute: buildDepsTree")
	command := fmt.Sprintf("cat %s | python - %s %s", generatePylist, manifestFilePath, m.DepsTreeFileName())
	shell, isSet := os.LookupEnv("SHELL")
	if !isSet {
		log.Fatal().Msgf("Please set $SHELL value to current shell.")
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	pylist, _ := filepath.Abs(m.DepsTreeFileName())
	log.Debug().Msgf("Success: buildDepsTree")
	return pylist
}

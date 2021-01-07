package pypi

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// IsSupportedManifestFormat checks for Supported Formats
func (*Matcher) IsSupportedManifestFormat(filename string) bool {
	log.Debug().Msgf("Executing: IsSupportedManifestFormat")
	s := strings.Split(filename, ".")
	name, ext := s[0], s[1]
	isExtSupported := checkExt(ext)
	isNameSupported := checkName(name)
	if isExtSupported && isNameSupported {
		log.Debug().Msgf("Success: IsSupportedManifestFormat")
		return true
	}
	log.Debug().Msgf("Success: IsSupportedManifestFormat")
	return false
}

// checkName checks for valid file name.
func checkName(name string) bool {
	match1, _ := regexp.MatchString("requirements?", name)
	match2, _ := regexp.MatchString("constraints?", name)
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
	m.getPylistGenerator()
	pathToPylist := m.buildDepsTree("generate_pylist.py", manifestFilePath)
	log.Debug().Msgf("Success: Generate Pylist")
	return pathToPylist
}

// getPylistGenerator generates `generate_pylist.py from `generatepylist.go`
func (m *Matcher) getPylistGenerator() error {
	log.Debug().Msgf("Executing: getPylistGenerator")
	// Generating generate_pylist.py
	err := ioutil.WriteFile("generate_pylist.py", CodeForPylist, 0644)
	if err != nil {
		log.Fatal().Msg("Error Generating generate_pylist.py")
	}
	log.Debug().Msgf("Done: getPylistGenerator")
	return nil
}

// buildDepsTree generates final Deps Tree and saves it as pylist.json
func (m *Matcher) buildDepsTree(generatePylist string, manifestFilePath string) string {
	log.Debug().Msgf("Execute: buildDepsTree")
	python, err := exec.LookPath("python")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Python v3.6+ is installed. Hint: Check same by executing python --version \n")
	}
	cmd := exec.Command(python, generatePylist, manifestFilePath, m.DepsTreeFileName())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	pylist, _ := filepath.Abs(m.DepsTreeFileName())
	log.Debug().Msgf("Success: buildDepsTree")
	return pylist
}

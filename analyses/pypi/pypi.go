package pypi

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

// GeneratePylist creates pylist.json from requirements.txt
func GeneratePylist(manifestFilePath string) string {
	log.Debug().Msgf("Executing: Generate Pylist")
	crdaTempPath := "/tmp/crda/"
	err := copyToTemp("analyses/pypi/generate_pylist.py", crdaTempPath)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	err = setUpEnv(crdaTempPath, manifestFilePath)
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	pathToPylist := buildDepsTree(crdaTempPath, manifestFilePath)
	log.Debug().Msgf("Success: Generate Pylist")
	return pathToPylist
}

//copyToTemp moves necessary files for generating pylist.json to tmp folder
func copyToTemp(src string, crdaTempPath string) error {
	log.Debug().Msgf("Executing: copyToTemp")
	absSource, _ := filepath.Abs(src)
	sourceFileStat, err := os.Stat(absSource)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	dst := crdaTempPath + sourceFileStat.Name()
	source, err := os.Open(absSource)
	if err != nil {
		return err
	}
	defer source.Close()
	os.MkdirAll(crdaTempPath, os.ModePerm)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Success: copyToTemp")
	return err
}

// setUpEnv sets up virtual env and install dependencies.
func setUpEnv(crdaTempPath string, manifestFilePath string) error {
	log.Debug().Msgf("Executing: setUpEnv")
	pyExecPath, err := exec.LookPath("python")
	if err != nil {
		log.Fatal().Err(err).Msgf("Please make sure Python v3.6+ is installed. Hint: Check same by executing `python --version`\n.")
	}
	os.MkdirAll(crdaTempPath, os.ModePerm)
	cmd := exec.Command(pyExecPath, "-m", "venv", "venv")
	log.Debug().Msgf("Executing command for Env. Setup %s", cmd.String())

	cmd.Stderr = os.Stdout
	cmd.Dir = crdaTempPath
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success: setUpEnv")
	return nil

}

// buildDepsTree generates final Deps Tree and saves it to pylist.json
func buildDepsTree(crdaTempDir string, manifestFilePath string) string {
	log.Debug().Msgf("Execute: buildDepsTree")
	pathToPylist := manifestFilePath + " pylist.json"
	command := fmt.Sprintf("cat %s/generate_pylist.py | python - %s", crdaTempDir, pathToPylist)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Dir = crdaTempDir
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	log.Debug().Msgf("Success: buildDepsTree")
	return crdaTempDir + "pylist.json"
}

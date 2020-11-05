package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDataFolder = "./../testdata/"
const testOutputManifest = "test_golist.json"

func TestMainWithInvalidNumOfArgs(t *testing.T) {
	cmd := exec.Command("go", "run", "build.go")
	cmd.Dir, _ = os.Getwd()
	err := cmd.Run()

	assert.Equal(t, "exit status 1", err.Error(), "Expecting error upon missing arguments")
}

func TestMainWithInvalidFolder(t *testing.T) {
	cmd := exec.Command("go", "run", "build.go", testDataFolder+"dummy", "/temp/go.list")
	cmd.Dir, _ = os.Getwd()
	err := cmd.Run()

	assert.Equal(t, "exit status 1", err.Error(), "Expecting error upon wrong source folder")
}

func TestMainHappyCase(t *testing.T) {
	manifestFilePath := testDataFolder + testOutputManifest
	defer os.Remove(manifestFilePath)

	cmd := exec.Command("go", "run", "build.go", testDataFolder, manifestFilePath)
	cmd.Dir, _ = os.Getwd()
	err := cmd.Run()

	// Still this input will generate a empty manifest file.
	assert.Equal(t, "exit status 1", err.Error(), "Expecting no errors")
}

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const testDataFolder = "./../testdata/"
const testOutputManifest = "test_manifest.json"

func readFileContentForTesting(fileName string) string {
	content, err := ioutil.ReadFile(testDataFolder + fileName)
	if err != nil {
		log.Fatal().Msgf("Exception: %v", err)
	}

	return string(content)
}

func TestMainWithInvalidNumOfArgs(t *testing.T) {
	manifestFilePath := testDataFolder + testOutputManifest

	defer os.Remove(manifestFilePath)

	os.Args = []string{"go_build_manifest/"}
	main()

	_, err := os.Stat(manifestFilePath)
	assert.NotEqual(t, nil, err, "Output manifest file size missmatch")
}

func TestMainWithInvalidFolder(t *testing.T) {
	manifestFilePath := testDataFolder + testOutputManifest

	defer os.Remove(manifestFilePath)

	os.Args = []string{"go_build_manifest/", testDataFolder + "dummy", manifestFilePath}
	main()

	_, err := os.Stat(manifestFilePath)
	assert.NotEqual(t, nil, err, "Output manifest file size missmatch")
}

func TestMainHappyCase(t *testing.T) {
	manifestFilePath := testDataFolder + testOutputManifest

	defer os.Remove(manifestFilePath)

	os.Args = []string{"go_build_manifest/", testDataFolder, manifestFilePath}
	main()

	// Read output json and check for its size
	output := readFileContentForTesting(manifestFilePath)
	assert.Equal(t, 40, len(output), "Output manifest file size missmatch")
}

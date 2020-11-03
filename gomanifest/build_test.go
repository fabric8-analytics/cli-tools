package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const dataFolder = "./../testdata/"
const outputManifest = "test_manifest.json"

func readFileContent(fileName string) string {
	content, err := ioutil.ReadFile(dataFolder + fileName)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func TestTransformationVerionSemVer(t *testing.T) {
	assert.Equal(t, "2.5.8", transformVersion("2.5.8"), "Semver positive transformation failed")
	assert.Equal(t, "3.2.0", transformVersion("v3.2.0"), "Semver 'v' transformation failed")
	assert.Equal(t, "3.2.0", transformVersion("v3.2.0+incompatible"), "Semver with incompatible transformation failed")
	assert.Equal(t, "3.2.0-alpha", transformVersion("v3.2.0-alpha+incompatible"), "Semver with alpha + incompatible transformation failed")
	assert.Equal(t, "3.2.0-beta1.5", transformVersion("v3.2.0-beta1.5"), "Semver with beta version transformation failed")
	assert.Equal(t, "3.2.0-beta1.2", transformVersion("v3.2.0-beta1.2+incompatible"), "Semver with beta version + incompatible transformation failed")
	assert.Equal(t, "3.2.0-20201023112233-abcd1234abcd", transformVersion("v3.2.0-20201023112233-abcd1234abcd"), "Pseudo version transformation failed")
	assert.Equal(t, "3.2.0-20201023112233-abcd1234abcd", transformVersion("v3.2.0-20201023112233-abcd1234abcd+alpha"), "Pseudo version with alpha transformation failed")
}

type FakeGoListCmd struct {
	mock.Mock
}

func (mock *FakeGoListCmd) Run() (string, error) {
	args := mock.Called()
	return args.String(0), args.Error(1)
}

var goDepsTestData = readFileContent("godeps.txt")

func TestProcessDepsDataFailCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("Run").Return("", errors.New("TEST :: Go list failure"))

	goList := &GoList{Command: fakeGoListCmd}
	assert.NotEqual(t, 0, goList.Get(), "Expect to handle go list command failure")
}

func TestProcessDepsDataHappyCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("Run").Return(goDepsTestData, nil)

	goList := &GoList{Command: fakeGoListCmd}
	assert.Equal(t, 0, goList.Get(), "Expect to handle go list command failure")
	assert.Equal(t, 12, len(depsPackages), "Package count check failed")
}

func TestBuildManifest(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("Run").Return(goDepsTestData, nil)

	goList := &GoList{Command: fakeGoListCmd}
	assert.Equal(t, 0, goList.Get(), "Expect to handle go list command failure")
	assert.Equal(t, 12, len(depsPackages), "Package count check failed")

	manifestFilePath := dataFolder + outputManifest
	buildManifest(manifestFilePath)

	// Read output json and check for its size
	output := readFileContent(outputManifest)
	assert.Equal(t, 264, len(output), "Output manifest file size missmatch")

	defer os.Remove(manifestFilePath)
}

func TestMainWithInvalidNumOfArgs(t *testing.T) {
	manifestFilePath := dataFolder + outputManifest
	os.Args = []string{"go_build_manifest/"}
	main()

	_, err := os.Stat(manifestFilePath)
	assert.NotEqual(t, nil, err, "Output manifest file size missmatch")

	defer os.Remove(manifestFilePath)
}

func TestMainWithInvalidFolder(t *testing.T) {
	manifestFilePath := dataFolder + outputManifest
	os.Args = []string{"go_build_manifest/", dataFolder + "dummy", manifestFilePath}
	main()

	_, err := os.Stat(manifestFilePath)
	assert.NotEqual(t, nil, err, "Output manifest file size missmatch")

	defer os.Remove(manifestFilePath)
}

func TestMainHappyCase(t *testing.T) {
	manifestFilePath := dataFolder + outputManifest

	os.Args = []string{"go_build_manifest/", dataFolder, manifestFilePath}
	main()

	// Read output json and check for its size
	output := readFileContent(manifestFilePath)
	assert.Equal(t, 264, len(output), "Output manifest file size missmatch")

	defer os.Remove(manifestFilePath)
}

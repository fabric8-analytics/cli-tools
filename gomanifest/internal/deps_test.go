package internal

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testDataFolder = "./../../testdata/"
const testOutputManifest = "test_manifest.json"

func readFileContentForTesting(fileName string) string {
	content, err := ioutil.ReadFile(testDataFolder + fileName)
	if err != nil {
		log.Fatal().Msgf("Exception: %v", err)
	}

	return string(content)
}

var goDepsTestData = readFileContentForTesting("godeps.txt")

type FakeGoListCmd struct {
	mock.Mock
}

func (mock *FakeGoListCmd) Run() (string, error) {
	args := mock.Called()
	return args.String(0), args.Error(1)
}

func TestProcessDepsDataFailCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("Run").Return("", errors.New("TEST :: Go list failure"))

	goList := &GoList{Command: fakeGoListCmd}
	_, err := goList.Get()
	assert.NotEqual(t, nil, err, "Expect to handle go list command failure")
}

func TestProcessDepsDataHappyCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("Run").Return(goDepsTestData, nil)

	goList := &GoList{Command: fakeGoListCmd}
	depsPackages, err := goList.Get()
	assert.Equal(t, nil, err, "Expect to handle go list command failure")
	assert.Equal(t, 12, len(depsPackages), "Package count check failed")
}

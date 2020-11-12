package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testDataFolder = "./../../testdata/"

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

// ReadCloser implements internal.GoList
func (mock *FakeGoListCmd) ReadCloser() io.ReadCloser {
	args := mock.Called()
	return ioutil.NopCloser(strings.NewReader(args.String(0)))
}

// Wait implements internal.GoList
func (mock *FakeGoListCmd) Wait() error {
	args := mock.Called()
	return args.Error(0)
}

func TestProcessDepsDataFailCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("ReadCloser").Return("")
	fakeGoListCmd.On("Wait").Return(fmt.Errorf("TEST :: Go list failure"))

	_, err := GetDeps(fakeGoListCmd)
	assert.NotEqual(t, nil, err, "Expect to handle go list command failure")
}

func TestProcessDepsDataHappyCase(t *testing.T) {
	fakeGoListCmd := &FakeGoListCmd{}
	fakeGoListCmd.On("ReadCloser").Return(goDepsTestData)
	fakeGoListCmd.On("Wait").Return(nil)

	depPackages, err := GetDeps(fakeGoListCmd)
	assert.Equal(t, nil, err, "Expect to handle go list command failure")
	assert.Equal(t, 12, len(depPackages), "Package count check failed")
}

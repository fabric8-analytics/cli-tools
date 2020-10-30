package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

const dataFolder = "./../testdata/"
const outputManifest = "test_manifest.json"

func TestListContainsPositive(t *testing.T) {
	existingList := []string{"abc", "def", "hij"}

	if !listContains(existingList, "abc") {
		t.Errorf("Contains check failed")
	}
}

func TestListContainsNegative(t *testing.T) {
	existingList := []string{"abc", "def", "hij"}

	if listContains(existingList, "rdt") {
		t.Errorf("Does not contain check failed")
	}
}

func TestTransformationVerionSemVer(t *testing.T) {
	if transformVersion("2.5.8") != "2.5.8" {
		t.Errorf("Semver positive transformation failed")
	}

	if transformVersion("v3.2.0") != "3.2.0" {
		t.Errorf("Semver 'v' transformation failed")
	}

	if transformVersion("v3.2.0+incompatible") != "3.2.0" {
		t.Errorf("Semver with incompatible transformation failed")
	}

	if transformVersion("v3.2.0-alpha+incompatible") != "3.2.0-alpha" {
		t.Errorf("Semver with alpha + incompatible transformation failed")
	}

	if transformVersion("v3.2.0-beta1.5") != "3.2.0-beta1.5" {
		t.Errorf("Semver with beta version transformation failed")
	}

	if transformVersion("v3.2.0-beta1.2+incompatible") != "3.2.0-beta1.2" {
		t.Errorf("Semver with beta version + incompatible transformation failed")
	}

	if transformVersion("v3.2.0-20201023112233-abcd1234abcd") != "3.2.0-20201023112233-abcd1234abcd" {
		t.Errorf("Pseudo version transformation failed")
	}

	if transformVersion("v3.2.0-20201023112233-abcd1234abcd+alpha") != "3.2.0-20201023112233-abcd1234abcd" {
		t.Errorf("Pseudo version with alpha transformation failed")
	}
}

var mockedExitStatus = 0
var mockedStdoutMod string = ""
var mockedStdoutList string = ""

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestExecCommandHelper", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	es := strconv.Itoa(mockedExitStatus)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1",
		"CMD=" + args[0],
		"STDOUT_MOD=" + mockedStdoutMod,
		"STDOUT_LIST=" + mockedStdoutList,
		"EXIT_STATUS=" + es}
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("CMD") == "mod" {
		fmt.Fprintf(os.Stdout, os.Getenv("STDOUT_MOD"))
	} else {
		fmt.Fprintf(os.Stdout, os.Getenv("STDOUT_LIST"))
	}
	i, _ := strconv.Atoi(os.Getenv("EXIT_STATUS"))
	os.Exit(i)
}

func readFileContent(fileName string) string {
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(dataFolder + fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	return text
}

func TestProcessGraphFailGoModGraph(t *testing.T) {
	mockedExitStatus = -1
	mockedStdoutMod = ""
	mockedStdoutList = ""
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if -1 != processGraphData() {
		t.Errorf("Expect go graph to fail")
	}
}

func TestProcessGraphHappyCase(t *testing.T) {
	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = ""
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if 0 != processGraphData() {
		t.Errorf("Expect go graph to pass")
	}

	if len(directDependencies) != 1 {
		t.Errorf("Expect direct dependencies to be %d, but found %d", 1, len(directDependencies))
	}
}

func TestProcessDepsDataFailGoList(t *testing.T) {
	mockedExitStatus = -1
	mockedStdoutMod = ""
	mockedStdoutList = ""
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if -1 != processDepsData() {
		t.Errorf("Expect go list -deps -json to fail")
	}
}

func TestProcessDepsDataHappyCase(t *testing.T) {
	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if 0 != processDepsData() {
		t.Errorf("Expect go list -deps -json to pass")
	}

	if len(goPackages.Packages) != 45 {
		t.Errorf("Expect packages to be %d, but found %d", 45, len(goPackages.Packages))
	}
}

func TestDirectDeps(t *testing.T) {
	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if 0 != processGraphData() {
		t.Errorf("Expect go graph to pass")
	}

	if 0 != processDepsData() {
		t.Errorf("Expect go list -deps -json to pass")
	}

	// Build direct deps
	buildDirectDependencies()

	if totalDirectModuleDependencies != 2 {
		t.Errorf("Expect direct module deps to be %d, but found %d", 2, totalDirectModuleDependencies)
	}

	if totalImports != 4 {
		t.Errorf("Expect import count to be %d, but found %d", 3, totalImports)
	}

	if totalDirectDependencies != 3 {
		t.Errorf("Expect direct package deps to be %d, but found %d", 2, totalDirectDependencies)
	}
}

func TestTransDeps(t *testing.T) {
	totalTransitivesDependencies = 0

	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if 0 != processGraphData() {
		t.Errorf("Expect go graph to pass")
	}

	if 0 != processDepsData() {
		t.Errorf("Expect go list -deps -json to pass")
	}

	// Build direct deps
	buildDirectDependencies()

	// Build trans deps
	buildTransitiveDeps()

	if totalTransitivesDependencies != 11 {
		t.Errorf("Expect transitive deps to be %d, but found %d", 11, totalTransitivesDependencies)
	}
}

func TestBuildManifest(t *testing.T) {
	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	if 0 != processDepsData() {
		t.Errorf("Expect go list -deps -json to pass")
	}

	// Build direct deps
	buildDirectDependencies()

	// Build trans deps
	buildTransitiveDeps()

	manifestFilePath = dataFolder + outputManifest
	// Build manifest file
	buildManifest()

	// Read output json and check for its size
	output := readFileContent(outputManifest)
	if len(output) != 2697 {
		t.Errorf("Expect manifest file size of %d bytes, but found %d bytes", 2697, len(output))
	}

	defer os.Remove(manifestFilePath)
}

func TestGenerateManifest(t *testing.T) {
	totalDirectDependencies = 0
	totalTransitivesDependencies = 0

	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	sourceRootFolder = dataFolder
	manifestFilePath = dataFolder + outputManifest
	generateManifest()

	defer os.Remove(manifestFilePath)

	if totalDirectDependencies != 3 {
		t.Errorf("Expect direct deps to be %d, but found %d", 3, totalDirectDependencies)
	}

	if totalTransitivesDependencies != 11 {
		t.Errorf("Expect transitive deps to be %d, but found %d", 11, totalTransitivesDependencies)
	}
}

func TestMainWithInvalidNumOfArgs(t *testing.T) {
	sourceRootFolder = ""
	manifestFilePath = ""
	outManifest := dataFolder + outputManifest
	os.Args = []string{"go_build_manifest/"}
	main()

	defer os.Remove(outManifest)

	if sourceRootFolder != "" {
		t.Errorf("Expect source folder to be empty, but found %s", sourceRootFolder)
	}

	if manifestFilePath != "" {
		t.Errorf("Expect manifest file path to be empty, but found %s", manifestFilePath)
	}
}

func TestMainWithInvalidFolder(t *testing.T) {
	sourceRootFolder = ""
	manifestFilePath = ""

	outManifest := dataFolder + outputManifest
	os.Args = []string{"go_build_manifest/", dataFolder + "dummy", outManifest}
	main()

	defer os.Remove(outManifest)

	if sourceRootFolder != "" {
		t.Errorf("Expect source folder to be empty, but found %s", sourceRootFolder)
	}

	if manifestFilePath != "" {
		t.Errorf("Expect manifest file path to be empty, but found %s", manifestFilePath)
	}
}

func TestMainHappyCase(t *testing.T) {
	sourceRootFolder = ""
	manifestFilePath = ""

	outManifest := dataFolder + outputManifest

	mockedExitStatus = 0
	mockedStdoutMod = readFileContent("gograph.txt")
	mockedStdoutList = readFileContent("godeps.txt")
	executeCommand = fakeExecCommand
	defer func() { executeCommand = exec.Command }()

	os.Args = []string{"go_build_manifest/", dataFolder, outManifest}
	main()

	defer os.Remove(outManifest)

	if sourceRootFolder != dataFolder {
		t.Errorf("Expect source folder is expected to be %s, but found %s", dataFolder, sourceRootFolder)
	}

	if manifestFilePath != outManifest {
		t.Errorf("Expect manifest file path is expected to be %s, but found %s", outManifest, manifestFilePath)
	}
}

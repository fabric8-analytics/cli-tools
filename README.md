# cli-tools
This repo would server as an interface between the different CRDA clients and the platform. It contains tools that will be used by clients inorder to generate required input for platform APIs. One such tool is `gomanifest`.

## gomanifest CLI tool
This tool shall be used by client that needs to trigger stack analyses request for golang software stack. The tool generate a manifest file that should be passed to stack analyses request for CRDA platform API. 

### Usage 
Tool can be used as shown below:

```
# Run command to generate manifest file.
go run github.com/fabric8-analytics/cli-tools/gomanifest /absolute/path/to/source/folder/containing/go.mod/ /absolute/path/to/save/generated/golist.json

```

It take two paths (1) Absolute source code path containing go.mod file. and (2) Absolute path to save manifest file in json format.

Ensure command is ran with proper previlage to read the source folder and to save the manifest file.
The tool intenall uses `go list` command to fetch source code dependency data. Ensure that go project is free from any dependency errors / issues.

### Contribution
To make changes in this tool you need to install `go` and development environment for executing go commands. Get the source from the repository.

#### Execute
To execute tool locally with developer changes, you can execute below command

```
go run <go tool source code>/gomanifest <Sample go project under test> <Path to manifest file>
```

#### Test
This tool has unit test which are packged along with source code. Required test data can be stored under `testdata` folder. Unit tests can be executes be using below command:

```

go test -v -cover ./...

```

Execute above command at root of the source tree, it runs all test cases and provides oneliner output in code coverage.
Sample output is as below:

```
$ go test -cover -v ./...
=== RUN   TestMainWithInvalidNumOfArgs
--- PASS: TestMainWithInvalidNumOfArgs (0.13s)
=== RUN   TestMainWithInvalidFolder
--- PASS: TestMainWithInvalidFolder (0.13s)
=== RUN   TestMainHappyCase
--- PASS: TestMainHappyCase (0.14s)
PASS
coverage: 0.0% of statements
ok  	github.com/fabric8-analytics/cli-tools/gomanifest	(cached)	coverage: 0.0% of statements
=== RUN   TestProcessDepsDataFailCase
{"level":"error","time":"2020-11-05T14:07:35+05:30","message":"`go list` command failed, clean dependencies using `go mod tidy` command"}
--- PASS: TestProcessDepsDataFailCase (0.00s)
=== RUN   TestProcessDepsDataHappyCase
{"level":"info","time":"2020-11-05T14:07:35+05:30","message":"Total packages: \t\t12"}
--- PASS: TestProcessDepsDataHappyCase (0.00s)
=== RUN   TestTransformationVerionSemVer
--- PASS: TestTransformationVerionSemVer (0.00s)
=== RUN   TestBuildManifest
{"level":"info","time":"2020-11-05T14:07:35+05:30","message":"Source code imports: \t13"}
{"level":"info","time":"2020-11-05T14:07:35+05:30","message":"Direct dependencies: \t1"}
--- PASS: TestBuildManifest (0.00s)
=== RUN   TestSaveManifest
{"level":"info","time":"2020-11-05T14:07:35+05:30","message":"Source code imports: \t13"}
{"level":"info","time":"2020-11-05T14:07:35+05:30","message":"Direct dependencies: \t1"}
--- PASS: TestSaveManifest (0.00s)
PASS
coverage: 79.4% of statements
ok  	github.com/fabric8-analytics/cli-tools/gomanifest/internal	(cached)	coverage: 79.4% of statements
```
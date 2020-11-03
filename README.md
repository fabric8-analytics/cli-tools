# cli-tools
This repo would server as an interface between the different CRDA clients and the platform. It contains tools that will be used by clients inorder to generate required input for platform APIs. One such tool is `gomanifest`.

## gomanifest CLI tool
This tool shall be used by client that needs to trigger stack analyses request for golang software stack. The tool generate a manifest file that should be passed to stack analyses request for CRDA platform API. 

### Usage 
Tool can be used as shown below:

```
# Run command to generate manifest file.
go run github.com/fabric8-analytics/cli-tools/gomanifest /absolute/path/to/source/folder/containing/go.mod/ /absolute/path/to/save/generated/manifestfile.json

```

It take two paths (1) Absolute source code path containing go.mod file. and (2) Absolute path to save manifest file in json format.

Ensure command is ran with proper previlage to read the source folder and to save the manifest file.
The tool intenall uses `go list` command to fetch source code dependency data. Ensure that go project is free from any dependency errors / issues.

### Test
This tool has unit test which are packged along with source code. Required test data can be stored under `testdata` folder. Unit tests can be executes be using below command:

```

go test -v -cover ./...

```

Execute above command at root of the source tree, it runs all test cases and provides oneliner output in code coverage.
Sample output is as below:

```
[dhpatel@dhpatel cli-tools]$ go test -v -cover ./...
=== RUN   TestTransformationVerionSemVer
--- PASS: TestTransformationVerionSemVer (0.00s)
=== RUN   TestProcessDepsDataFailCase
2020/11/03 13:12:04 ERROR :: Command `go list -json -deps ./...` failed, resolve project build errors. TEST :: Go list failure
--- PASS: TestProcessDepsDataFailCase (0.00s)
=== RUN   TestProcessDepsDataHappyCase
2020/11/03 13:12:04 Packages in deps: 48
2020/11/03 13:12:04 Filter package count: 12
--- PASS: TestProcessDepsDataHappyCase (0.00s)
=== RUN   TestBuildManifest
2020/11/03 13:12:04 Packages in deps: 48
2020/11/03 13:12:04 Filter package count: 12
2020/11/03 13:12:04 Source code imports: 22
--- PASS: TestBuildManifest (0.01s)
=== RUN   TestMainWithInvalidNumOfArgs
2020/11/03 13:12:04 Error :: Invalid arguments for the command.
2020/11/03 13:12:04 Usage :: go run github.com/dgpatelgit/gobuildmanifest <Absolute source root folder path containing go.mod> <Output file path>.json
2020/11/03 13:12:04 Example :: go run github.com/dgpatelgit/gobuildmanifest /home/user/goproject/root/folder /home/user/gomanifest.json
--- PASS: TestMainWithInvalidNumOfArgs (0.00s)
=== RUN   TestMainWithInvalidFolder
2020/11/03 13:12:04 ERROR :: Invalid source folder path :: ./../testdata/dummy
--- PASS: TestMainWithInvalidFolder (0.00s)
=== RUN   TestMainHappyCase
2020/11/03 13:12:04 Building manifest file for :: ./../testdata/
Outp 
2020/11/03 13:12:04 Packages in deps: 0
2020/11/03 13:12:04 Filter package count: 12
2020/11/03 13:12:04 Source code imports: 22
2020/11/03 13:12:04 Success :: Manifest file generated and stored at ./../testdata/test_manifest.json
--- PASS: TestMainHappyCase (0.00s)
PASS
coverage: 93.0% of statements
ok  	github.com/fabric8-analytics/cli-tools/gomanifest	0.013s	coverage: 93.0% of statements
[dhpatel@dhpatel cli-tools]$
```
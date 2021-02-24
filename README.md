# cli-tools
<p align="center">
    <a alt="GoReport" href="https://goreportcard.com/report/github.com/fabric8-analytics/cli-tools">
        <img alt="GoReport" src="https://goreportcard.com/badge/github.com/fabric8-analytics/cli-tools">
    </a>
    <a href="https://github.com/fabric8-analytics/cli-tools/actions?query=workflow%3ACI">
        <img alt="CI" src="https://github.com/fabric8-analytics/cli-tools/workflows/CI/badge.svg">
    </a>
      <a href="https://codecov.io/gh/fabric8-analytics/cli-tools">
        <img src="https://codecov.io/gh/fabric8-analytics/cli-tools/branch/main/graph/badge.svg?token=AN4JV0BPW1"/>
      </a>
</p>

This repo would server as an interface between the different CRDA clients and the platform. It contains tools that will be used by clients inorder to generate required input for platform APIs. One such tool is `gomanifest`.
### Tools and Packages:

  * `CRDA Cli`:  CLI Tools to interact with CRDA Platform. [Learn more](docs/cli_README.md)
  * `GoManifest`: Go Package used in Stack Analyses for Golang ecosystem. [Learn more](docs/gomanifest.md)

### Contribution
To make changes in this tool you need to install `go` and development environment for executing go commands. Get the source from the repository.

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

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


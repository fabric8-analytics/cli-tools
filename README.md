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

This repo would serve as an interface between the different CRDA clients and the platform. It contains tools that will be used by clients inorder to generate required input for platform APIs. One such tool is `gomanifest`.
### Tools and Packages:

  * `CRDA Cli`:  CLI Tools to interact with CRDA Platform. [Learn more](docs/cli_README.md)
  * `GoManifest`: Go Package used in Stack Analyses for Golang ecosystem. [Learn more](docs/gomanifest.md)

## Ignoring Vulnerabilities during analysis

If users wish to ignore vulnerabilities for a dependency, it can be done by adding "crdaignore" as a comment in manifest file for Python, Maven, Golang and as a JSON for Node ecosystem because Node manifest files dont support comments hence it has to be given inside a JSON.
If "crdaignore" is followed by a list of comma separated Snyk vulnerability IDs then only listed vulnerabilities will be ignored during analysis, in case "crdaignore" is not followed by any list all vulnerabilities present in package will be ignored.

# Examples

# Python
Ignore all vulnerabilities in fastapi and few for flask

```
fastapi==0.36.0 #crdaignore
sceptre==2.2.1
flask==1.0 #crdaignore [<Snyk vulnerability ID 1 >, <Snyk vulnerability ID 2 >]
```
# Golang
Ignore all the security vulnerabilities present in the "ginkgo" and "pax-go" dependencies in a golang manifest file.
```
	code.cloudfoundry.org/archiver v0.0.0-20170223024658-7291196139d7
	github.com/googleapis/gax-go v1.0.3 //crdaignore [<Snyk vulnerability ID 1 >]
	github.com/googleapis/gax-go/v2 v2.0.5
	github.com/onsi/ginkgo v1.14.2 // indirect crdaignore 
	github.com/onsi/gomega v1.10.3 // indirect 

```
# Maven
Ignore all vulnerabilities of the dependency "junit:junit". 

```
 <dependency>
      <groupId>junit</groupId>  <!--crdaignore-->
      <artifactId>junit</artifactId>
      <version>3.8.1</version>
 </dependency>
```
Note: To ignore vulnerabilities for a dependency in a maven manifest file, insert "crdaignore" in comments against the group id, artifact id, or version of that particular dependency.

# Node
Ignore all the security vulnerabilities for "bootstrap" and a set of vulnerabilities for the "lodash" package.
```
"crdaignore": {
			"packages": {
				"bootstrap": [
					"*"
				],
				"lodash": [<Snyk vulnerability ID 1 >]
			}
	},
```
A sample npm manifest file with the security vulnerabilities to ignore during analysis:
```
{
		"name": "node-js-sample",
		"version": "0.2.0",
		"description": "A sample Node.js app using Express 4",
		"main": "index.js",
		"scripts": {
				"start": "node index.js"
		},
		"dependencies": {
				"ansi-styles": "3.2.1",
				"escape-string-regexp": "1.0.5",
				"supports-color": "5.5.0",
				"cordova-plugin-camera": "4.1.0",
				"bootstrap": "4.1.1",
				"libnmap": "0.4.15",
				"lodash": "4.17.11",
				"html-purify": "1.1.0"
		},
		"engines": {
				"node": "4.0.0"
		},
		"crdaignore": {
			"packages": {
				"bootstrap": [
					"*"
				],
				"lodash": ["vulnerability 1"]
			}
		},
		"repository": {
				"type": "git",
				"url": "https://github.com/heroku/node-js-sample"
		},
		"keywords": [
				"node",
				"heroku",
				"express"
		],
		"author": "Mark Pundsack",
		"contributors": [
				"Zeke Sikelianos <zeke@sikelianos.com> (http://zeke.sikelianos.com)"
		],
		"license": "MIT"
}
```

### Contribution
To make changes in this tool you need to install `go` and development environment for executing go commands. Get the source from the repository.

#### Test
This tool has unit test which are packged along with source code. Required test data can be stored under `testdata` folder. 
Execute above command at root of the source tree, it runs all test cases and provides oneliner output in code coverage.

`make test`

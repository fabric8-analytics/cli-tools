## CRDA CLI

#### CRDA CLI is a tool used to perform dependency stack scanning, right from CLI. 
As of now, CLI supports following stacks:
- Node (NPM): `package.json` 
- Golang (Go): `go.mod`
- Java (Maven): `pom.xml`
- Python (pip): `requirements.txt`

![screenshot of summary](https://github.com/fabric8-analytics/cli-tools/blob/b407d2a7c595a47e3126ad62a816dc107bd148d2/summary.png)
![screenshot of analyse](https://github.com/fabric8-analytics/cli-tools/blob/71198735d0dee3173ed3082a5ab1dee41dfa9ce8/analyse.png)

### Installation:
- Select, Download and Install latest binary from [Releases](https://github.com/fabric8-analytics/cli-tools/releases)

#### curl

- ##### For Linux
```
$ curl -s -L https://github.com/fabric8-analytics/cli-tools/releases/download/v0.0.1/crda_0.0.1_Linux_64bit.tar.gz | tar xvz -C .
```
- ##### For Linux - Fedora/CentOS/RHEL
```
$ curl -s -L https://github.com/fabric8-analytics/cli-tools/releases/download/v0.0.1/crda_0.0.1_Linux-64bit.rpm 
```
- ##### For MacOS
```
$ curl -s -L https://github.com/fabric8-analytics/cli-tools/releases/download/v0.0.1/crda_0.0.1_macOS_64bit.tar.gz | tar xvz -C .
```
- ##### For MacOS - Apple Silicon
```
$ curl -s -L https://github.com/fabric8-analytics/cli-tools/releases/download/v0.0.1/crda_0.0.1_macOS_ARM64.tar.gz | tar xvz -C .
```
- ##### For Windows
Click [here](https://github.com/fabric8-analytics/cli-tools/releases/download/v0.0.1/crda_0.0.1_Windows_64bit.tar.gz) to start download.

### Usage:
Executable supports following commands:

* Please install manifest dependencies first to have correct CLI behaviour.

- `crda auth`: This command is used to enable user to Authenticate with CRDA Server. It outputs a unique UUID. This command generates and saves `crda_key` in `$HOME/.crda/config.yaml`

    Supported Flags:
    * `--synk-token` (string) (OPTIONAL): Can be obtained from [here](https://app.snyk.io/login?utm_campaign=Code-Ready-Analytics-2020&utm_source=code_ready&code_ready=FF1B53D9-57BE-4613-96D7-1D06066C38C9). If not set, Freemium a/c with limited functionality will be created.
    Please note, New Token generated is confidential and is mapped to your Synk Account. Keep it safe!
    * `--help` (Optional): Command level Help.


- `crda analyse`: Command to perform Full Stack Analyses. 
    Supported Arguments:
    * (string) (Required): Manifest file Absolute Path. Ex: for Node, usually its `/path/to/package.json`, similarly `/path/to/pom.xml`for Java.

    * `--help` (Optional): Command level Help.

- `crda version`: This outputs version details of Binary.


#### Global Flags:
- `--debug`: (bool) (Optional): Debug Flag. Enables Debug Logs
- `--no-color`: (bool) (Optional): Toggles colors in output.
- `--help` : help about binary functionalities.

### EXIT CODES

Possible exit codes and their meaning:

- 0: success, no vulns found
- 1: failure, try to re-run command
- 2: action_needed, vulns found


#### Build:

```go
go build -o crda
```


### Issue:
Got Issues..? We got your back. Tell Us here: [Raise Issue](https://github.com/fabric8-analytics/cli-tools/issues) 

### Feedback: 
We Love stars, just like you do.  
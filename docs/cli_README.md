## CRDA CLI

#### CRDA CLI is a CLI Tool to interact with CRDA Platform. 

### Installation:
- Select, Download and Install latest binary from [Releases](https://github.com/fabric8-analytics/cli-tools/releases)

### Usage:
Executable supports following commands:

- `crda auth`: This command is used to enable user to Authenticate with CRDA Server. This command generates and saves Auth Token in `$HOME/.crda/config.yaml`
    Supported Flags:
    * `--synk_token` (string) (OPTIONAL): Can be obtained from [here](https://app.snyk.io/login?utm_campaign=Code-Ready-Analytics-2020&utm_source=code_ready&code_ready=FF1B53D9-57BE-4613-96D7-1D06066C38C9). If not set, Freemium a/c with limited functionality will be created.
    Please note, New Token generated is confidential and is mapped to your Synk Account. Keep it safe!

- `crda analyse`: Command to perform Full Stack Analyses. 
    Supported Flags:
    * `--file`: (string) (Required): Manifest file Absolute Path. Ex: for Node, usually its `/path/to/package.json`, similarly `/path/to/pom.xml`for Java.


#### Global Flags:
- `--debug`: (bool) (Optional): Debug Flag. Enables Debug Logs

#### Execute:
To execute tool locally with developer changes, you can execute below command

```go
go run main.go <any-cmd> <flags>
```

#### Build:

```go
go build -o crda
```


### Issue:
Got Issues..? We got your back. Tell Us here: [Raise Issue](https://github.com/fabric8-analytics/cli-tools/issues) 

### Feedback: 
We Love stars, just like you do.  
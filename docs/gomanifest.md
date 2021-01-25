
## gomanifest CLI tool
This tool shall be used by client that needs to trigger stack analyses request for golang software stack. The tool generate a manifest file that should be passed to stack analyses request for CRDA platform API.

### Usage
Tool can be used as shown below:

```
# Run command to generate manifest file.
go run github.com/fabric8-analytics/cli-tools/gomanifest /absolute/path/to/source/folder/containing/go.mod/ /absolute/path/to/save/generated/golist.json [/absolute/path/to/goexecutable/go]
```

It take three arguments as below:
- Absolute source code path containing go.mod file.
- Absolute path to save manifest file in json format.
- [Optional] This is a optional parameter to provide absolute path to `go` executable. This might be required in case gomanifest command is used in a script that executes in a process which does not have `PATH` variable set to locate `go` executable.

Sample executable script will be as below:
```
# Get the latest gomanifest package
go get -u github.com/fabric8-analytics/cli-tools/gomanifest

# Execute gomanifest to generate manifest (golist.json)
go run github.com/fabric8-analytics/cli-tools/gomanifest /home/user1/example/source /home/user1/golist.json /usr/local/bin/go
```

Ensure command is ran with proper previlage to read the source folder and to save the manifest file.
The tool intenall uses `go list` command to fetch source code dependency data. Ensure that go project is free from any dependency errors / issues.

#### Execute
To execute tool locally with developer changes, you can execute below command

```
go run <go tool source code>/gomanifest <Sample go project under test> <Path to manifest file> [<Go executable path>]
```

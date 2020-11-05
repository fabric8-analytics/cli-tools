package internal

import (
	"encoding/json"
	"io"
	"os/exec"

	"github.com/rs/zerolog/log"
)

/* Structure required for parsing `go list -deps -json ./...` json parsing */

// DepModule ... Module structure from deps json
type DepModule struct {
	Path    string     `json:"Path"`
	Main    bool       `json:"Main"`
	Version string     `json:"Version"`
	Replace *DepModule `json:"Replace"`
}

// DepPackage ... Package structure from deps json
type DepPackage struct {
	Root       string    `json:"Root"`
	ImportPath string    `json:"ImportPath"`
	Module     DepModule `json:"Module"`
	Standard   bool      `json:"Standard"`
	Imports    []string  `json:"Imports"`
	Deps       []string  `json:"Deps"`
}

// GoListCmd ... Go list command structure.
type GoListCmd struct {
	CWD string
}

// Run ... Actual function that executes go list command and returns output as string.
func (goListCmd GoListCmd) Run() (io.ReadCloser, error) {
	GoListGoListDeps := exec.Command("go", "list", "-json", "-deps", "./...")
	GoListGoListDeps.Dir = goListCmd.CWD
	defer GoListGoListDeps.Start()
	GoListGoListDeps.Wait()
	return GoListGoListDeps.StdoutPipe()
}

// GoListCmdInterface ... Interface to be implemented to execute go list command.
type GoListCmdInterface interface {
	Run() (io.ReadCloser, error)
}

// GoList ... Structure that handle go list data and extract required packages.
type GoList struct {
	Command GoListCmdInterface
}

// parseDepsJSON ... Parse json output of go list -dep command.
func parseDepsJSON(out chan<- DepPackage, in io.Reader) {
	defer close(out)

	dec := json.NewDecoder(in)
	for {
		var m DepPackage
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Error().Msg("Parsing json data failed.")
			panic(err)
		}
		out <- m
	}
}

// Get ... Get deps data through go list deps command and converts json into objects.
func (goList *GoList) Get() (map[string]DepPackage, error) {
	var depPackagesMap = make(map[string]DepPackage)

	output, err := goList.Command.Run()
	if err != nil {
		log.Error().Msgf("`go list` command failed, clean dependencies using `go mod tidy` command")
		return nil, err
	}

	ch := make(chan DepPackage, 0)
	go parseDepsJSON(ch, output)

	// Preprocess and remove all standard packages.
	for pckg := range ch {
		// Exclude standard packages
		if !pckg.Standard {
			depPackagesMap[pckg.ImportPath] = pckg
		}
	}
	log.Info().Msgf("Total packages: \t\t%d", len(depPackagesMap))

	return depPackagesMap, nil
}

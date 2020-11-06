package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

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

// GoList ... Interface to be implemented to execute go list command.
type GoList interface {
	ReadCloser() io.ReadCloser
	Wait() error
}

// parseDepsJSON ... Parse json output of go list -dep command.
func parseDepsJSON(dataC chan<- DepPackage, errorC chan<- error, in io.ReadCloser) {
	defer close(dataC)
	defer close(errorC)
	defer in.Close()

	dec := json.NewDecoder(in)
	for {
		var m DepPackage
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			errorC <- err
			break
		}
		dataC <- m
	}
}

// GetDeps ... GetDeps converts GoList output into DepPackage format.
func GetDeps(cmd GoList) (map[string]DepPackage, error) {
	var depPackagesMap = make(map[string]DepPackage)

	dataC := make(chan DepPackage, 0)
	errorC := make(chan error, 1)
	go parseDepsJSON(dataC, errorC, cmd.ReadCloser())

	for pckg := range dataC {
		// Exclude standard packages
		if !pckg.Standard {
			depPackagesMap[pckg.ImportPath] = pckg
		}
	}

	// Check for json Decode error.
	if err := <-errorC; err != nil {
		return nil, err
	}
	// Wait for the `go list` command to complete.
	if err := cmd.Wait(); err != nil {
		return nil, errors.New(fmt.Sprintf("%v: `go list` failed, use `go mod tidy` to known more", err))
	}
	log.Info().Msgf("Total packages: \t\t%d", len(depPackagesMap))
	return depPackagesMap, nil
}

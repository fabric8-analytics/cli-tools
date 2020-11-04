package internal

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

/* Structure required for parsing `go list -deps -json ./...` json parsing */

// DepPackages ... Package list structure from deps json
type DepPackages struct {
	Packages []DepPackage `json:"Packages"`
}

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
func (goListCmd GoListCmd) Run() (string, error) {
	GoListGoListDeps := exec.Command("go", "list", "-json", "-deps", "./...")
	GoListGoListDeps.Dir = goListCmd.CWD
	output, err := GoListGoListDeps.Output()

	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GoListCmdInterface ... Interface to be implemented to execute go list command.
type GoListCmdInterface interface {
	Run() (string, error)
}

// GoList ... Structure that handle go list data and extract required packages.
type GoList struct {
	Command GoListCmdInterface
}

// Get ... Get deps data through go list deps command and converts json into objects.
func (goList *GoList) Get() (map[string]DepPackage, error) {
	output, err := goList.Command.Run()

	if err != nil {
		log.Error().Msgf("`go list` command failed, clean dependencies using `go mod tidy` command")
		return nil, err
	}

	goListDepsData := string(output)
	goListDepsData = `{"Packages": [` + strings.ReplaceAll(goListDepsData, "}\n{", "},\n{") + "]}"

	var depPackages DepPackages
	json.Unmarshal([]byte(goListDepsData), &depPackages)
	log.Info().Msgf("Packages in deps: %d", len(depPackages.Packages))

	var depPackagesMap = make(map[string]DepPackage)

	// Preprocess and remove all standard packages.
	for _, pckg := range depPackages.Packages {
		// Exclude standard packages
		if pckg.Standard == false {
			depPackagesMap[pckg.ImportPath] = pckg
		}
	}
	log.Info().Msgf("Filter package count: %d", len(depPackagesMap))

	return depPackagesMap, nil
}

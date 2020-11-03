package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

/* MANIFEST FILE VERSION */
const manifestVersion = "v1"

/* Structure required for parsing `go list -deps -json ./...` json parsing */

// GoPackages ... Package list structure from deps json
type GoPackages struct {
	Packages []GoPackage `json:"Packages"`
}

// GoModule ... Module structure from deps json
type GoModule struct {
	Path    string    `json:"Path"`
	Main    bool      `json:"Main"`
	Version string    `json:"Version"`
	Replace *GoModule `json:"Replace"`
}

// GoPackage ... Package structure from deps json
type GoPackage struct {
	Root       string   `json:"Root"`
	ImportPath string   `json:"ImportPath"`
	Module     GoModule `json:"Module"`
	Standard   bool     `json:"Standard"`
	Imports    []string `json:"Imports"`
	Deps       []string `json:"Deps"`
}

/* Structure for output the manifest */

// ManifestTransitiveDeps ... Transitive details in final manifest
type ManifestTransitiveDeps struct {
	Name    string `json:"package"`
	Version string `json:"version"`
}

// MainfestDirectDeps ... Direct dependency details.
type MainfestDirectDeps struct {
	Name         string                   `json:"package"`
	Version      string                   `json:"version"`
	Dependencies []ManifestTransitiveDeps `json:"deps"`
}

// Manifest ... Final manifest file structure
type Manifest struct {
	Version  string               `json:"version"`
	Main     string               `json:"main"`
	Packages []MainfestDirectDeps `json:"packages"`
}

// Deps package from `go list -deps -json ./...` command
var depsPackages = make(map[string]GoPackage)

// transformVersion ... Converts the golang version string into semver without 'v' and appended text after '+'
func transformVersion(inVersion string) string {
	var outVersion string = strings.Replace(inVersion, "v", "", 1)

	return strings.Split(outVersion, "+")[0]
}

// getPackageName ... Utility function to convert package + module data into package name used by manifest.
func getPackageName(depPackage GoPackage) string {
	// Get module / package@module
	if depPackage.ImportPath != depPackage.Module.Path {
		return depPackage.ImportPath + "@" + depPackage.Module.Path
	}

	return depPackage.ImportPath
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

	fmt.Println("Outp", string(output))
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
func (goList *GoList) Get() int {
	output, err := goList.Command.Run()

	if err != nil {
		log.Println("ERROR :: Command `go list -json -deps ./...` failed, resolve project build errors.", err)
		return -1
	}

	goListDepsData := string(output)
	goListDepsData = `{"Packages": [` + strings.ReplaceAll(goListDepsData, "}\n{", "},\n{") + "]}"

	var goPackages GoPackages
	json.Unmarshal([]byte(goListDepsData), &goPackages)
	log.Println("Packages in deps:", len(goPackages.Packages))

	// Preprocess and remove all standard packages.
	for i := 0; i < len(goPackages.Packages); i++ {
		// Exclude standard packages
		if goPackages.Packages[i].Standard == false {
			depsPackages[goPackages.Packages[i].ImportPath] = goPackages.Packages[i]
		}
	}
	log.Println("Filter package count:", len(depsPackages))

	return 0
}

// getTransitives ... Returns a clean list of deps
func getTransitives(deps []string) []ManifestTransitiveDeps {
	var manifestDependencies = make([]ManifestTransitiveDeps, 0)
	for i := 0; i < len(deps); i++ {
		if depPackage, ok := depsPackages[deps[i]]; ok {
			if depPackage.Module.Main == false {
				manifestDependencies = append(manifestDependencies, ManifestTransitiveDeps{
					getPackageName(depPackage),
					transformVersion(depPackage.Module.Version),
				})
			}
		}
	}
	return manifestDependencies
}

// buildManifest ... Build direct & transitive dependencies.
func buildManifest(manifestFilePath string) {
	var manifest Manifest = Manifest{manifestVersion, "", make([]MainfestDirectDeps, 0)}

	// Get direct imports from current source.
	var sourceImports = make(map[string]bool, 0)
	for _, pckg := range depsPackages {
		// Include only packages with project ROOT
		if pckg.Module.Main == true {
			// Set main module if not set.
			if manifest.Main == "" {
				manifest.Main = pckg.Module.Path
			}

			// Added imports as direct dependencies
			for _, imp := range pckg.Imports {
				// Add imports if not added yet.
				if _, ok := sourceImports[imp]; !ok {
					// Added imports that are non-standard (or present in deps packages) and
					// which are having main module as 'false'
					if depPackage, ok := depsPackages[imp]; ok {
						if depPackage.Module.Main == false {
							manifest.Packages = append(manifest.Packages,
								MainfestDirectDeps{
									getPackageName(depPackage),
									transformVersion(depPackage.Module.Version),
									getTransitives(depPackage.Deps),
								})
						}
					}

					sourceImports[imp] = true
				}
			}
		}
	}
	log.Println("Source code imports:", len(sourceImports))

	d, err := json.Marshal(manifest)
	if err != nil {
		panic(err)
	}
	var directDependenciesJSON string = string(d)

	f, err := os.Create(manifestFilePath)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(string(directDependenciesJSON))
	if err != nil {
		panic(err)
	}
	f.Sync()

	defer f.Close()
}

func main() {
	if len(os.Args) != 3 {
		log.Println("Error :: Invalid arguments for the command.")
		log.Println("Usage :: go run github.com/fabric8-analytics/cli-tools/gomanifest <Absolute source root folder path containing go.mod> <Output file path>.json")
		log.Println("Example :: go run github.com/fabric8-analytics/cli-tools/gomanifest /home/user/goproject/root/folder /home/user/gomanifest.json")
	} else {
		_, err := os.Stat(os.Args[1])
		if err != nil {
			log.Println("ERROR :: Invalid source folder path ::", os.Args[1])
		} else {
			log.Println("Building manifest file for ::", os.Args[1])

			goListCmd := &GoListCmd{CWD: os.Args[1]}
			goList := &GoList{Command: goListCmd}
			if goList.Get() == 0 {
				buildManifest(os.Args[2])
				log.Println("Success :: Manifest file generated and stored at", os.Args[2])
			} else {
				log.Fatalln("ERROR :: Could not run go list command, clean dependencies using `go mod tidy` command")
			}
		}
	}
}

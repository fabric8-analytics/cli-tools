package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
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

// Transitive ... Transitive details
type Transitive struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Included bool     `json:"include"`
	Packages []string `json:"packages"`
}

// DirectPackage ... Direct package details
type DirectPackage struct {
	Name        string       `json:"name"`
	Transitives []Transitive `json:"transitives"`
}

// DirectDependency ... Direct dependency details
type DirectDependency struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Included    bool            `json:"include"`
	Packages    []DirectPackage `json:"packages"`
	Transitives []Transitive    `json:"transitives"`
}

// ManifestTransitive ... Transitive details in final manifest
type ManifestTransitive struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
}

// Manifest ... Final manifest file structure
type Manifest struct {
	Version 	 string 			`json:"version"`
	Main         string             `json:"main"`
	Dependencies []DirectDependency `json:"dependencies"`
}

// Source root folder, set via command line ARGS
var sourceRootFolder = ""

// Destination manifest file path, set via command line ARGS
var manifestFilePath = ""

var goPackages GoPackages

var mainModule string = ""
var directDependencies = make(map[string]DirectDependency)
var totalDirectModuleDependencies = 0
var totalDependencyPackages = 0
var totalImports = 0
var totalDirectDependencies = 0
var totalTransitivesDependencies = 0

func contains(s []string, searchterm string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processGraphData() int {
	// Get graph data
	cmdGoModGraph := exec.Command("go", "mod", "graph")
	cmdGoModGraph.Dir = sourceRootFolder
	if output, err := cmdGoModGraph.Output(); err != nil {
		fmt.Printf("ERROR :: Command `go mod graph` failed, resolve project build errors. %s\n", err)
		return -1
	} else {
		for _, value := range strings.Split(string(output), "\n") {
			if len(value) > 0 {
				// Extract direct dependency from go.mod
				pc := strings.Split(value, " ")
				if !strings.Contains(pc[0], "@") {
					mainModule = pc[0]
					mv := strings.Split(pc[1], "@")
					directDependencies[mv[0]] = DirectDependency{mv[0], mv[1], false, make([]DirectPackage, 0), make([]Transitive, 0)}
					totalDirectModuleDependencies++
				}
			}
		}
		fmt.Println("Direct module dependencies:", totalDirectModuleDependencies)
	}

	return 0
}

func processDepsData() int {
	// Switch to code directory and get graph
	cmdGoListDeps := exec.Command("go", "list", "-json", "-deps", "./...")
	cmdGoListDeps.Dir = sourceRootFolder
	if output, err := cmdGoListDeps.Output(); err != nil {
		fmt.Println("ERROR :: Command `go list -json -deps ./...` failed, resolve project build errors.", err)
		return -1
	} else {
		goListDepsData := string(output)
		goListDepsData = "{\"Packages\": [" + strings.ReplaceAll(goListDepsData, "}\n{", "},\n{") + "]}"

		json.Unmarshal([]byte(goListDepsData), &goPackages)
		totalDependencyPackages = len(goPackages.Packages)
		fmt.Println("Packages in deps:", totalDependencyPackages)
	}
	return 0
}

func buildDirectDependencies() {
	// Get direct imports from current source.
	var sourceImports []string
	for i := 0; i < len(goPackages.Packages); i++ {
		// Exclude standard packages and include only packages with project ROOT
		if goPackages.Packages[i].Standard == false && !strings.Contains(goPackages.Packages[i].Root, "@") {
			for _, imp := range goPackages.Packages[i].Imports {
				if !contains(sourceImports, imp) {
					sourceImports = append(sourceImports, imp)
				}
			}
		}
	}
	totalImports = len(sourceImports)
	fmt.Println("Source code imports:", totalImports)

	for mk, mod := range directDependencies {
		for _, imp := range sourceImports {
			if imp == mod.Name || strings.HasPrefix(imp, mod.Name+"/") {
				om := directDependencies[mk]
				if imp == mod.Name {
					om.Included = true
				} else {
					om.Packages = append(directDependencies[mk].Packages, DirectPackage{imp, make([]Transitive, 0)})
				}
				directDependencies[mk] = om
				totalDirectDependencies++
			}
		}
	}
	fmt.Println("Direct dependencies from imports:", totalDirectDependencies)
}

func findAndAddTransitive(importPath string, transitivies []Transitive) []Transitive {
	for i := 0; i < len(goPackages.Packages); i++ {
		if goPackages.Packages[i].Standard == false && goPackages.Packages[i].ImportPath == importPath {
			var foundModule = false
			for t := 0; t < len(transitivies); t++ {
				if transitivies[t].Name == goPackages.Packages[i].Module.Path {
					if importPath == goPackages.Packages[i].Module.Path {
						transitivies[t].Included = true
					} else {
						transitivies[t].Packages = append(transitivies[t].Packages, importPath)
					}
					foundModule = true
					totalTransitivesDependencies++
				}
			}

			if !foundModule {
				var newTrans = Transitive{
					goPackages.Packages[i].Module.Path,
					goPackages.Packages[i].Module.Version,
					importPath == goPackages.Packages[i].Module.Path,
					make([]string, 0)}
				if importPath != goPackages.Packages[i].Module.Path {
					newTrans.Packages = append(newTrans.Packages, importPath)
				}
				transitivies = append(transitivies, newTrans)
				totalTransitivesDependencies++
			}
		}
	}
	return transitivies
}

func getTransitiveDetails(modPath string, importPath string) []Transitive {
	var transitivies = make([]Transitive, 0)
	for i := 0; i < len(goPackages.Packages); i++ {
		if goPackages.Packages[i].ImportPath == importPath {
			if goPackages.Packages[i].Standard == true {
				fmt.Println("Skipping strandard import ::", importPath)
				break
			}

			for _, dv := range goPackages.Packages[i].Deps {
				transitivies = findAndAddTransitive(dv, transitivies)
			}
		}
	}
	return transitivies
}

func buildTransitiveDeps() {
	for k, ddeps := range directDependencies {
		var dm = directDependencies[k]
		if ddeps.Included {
			dm.Transitives = getTransitiveDetails(ddeps.Name, ddeps.Name)
		}
		dm.Packages = make([]DirectPackage, 0)

		for _, pckg := range ddeps.Packages {
			pckg.Transitives = getTransitiveDetails(ddeps.Name, pckg.Name)
			dm.Packages = append(dm.Packages, pckg)
		}
		directDependencies[k] = dm
	}
	fmt.Println("Total transitive dependencies:", totalTransitivesDependencies)
}

func buildManifest() {
	var manifest Manifest = Manifest{manifestVersion, mainModule, make([]DirectDependency, 0)}
	for mk := range directDependencies {
		manifest.Dependencies = append(manifest.Dependencies, directDependencies[mk])
	}

	d, err := json.Marshal(manifest)
	check(err)
	var directDependenciesJSON string = string(d)

	f, err := os.Create(manifestFilePath)
	check(err)
	_, err = f.WriteString(string(directDependenciesJSON))
	check(err)
	f.Sync()

	defer f.Close()
	fmt.Println("Success :: Manifest generated & stored at", manifestFilePath)
}

func flatenTransitives(transtivies []Transitive) string {
	var returnString = ""
	var firstEntry = true
	for _, t := range transtivies {
		if t.Included {
			if !firstEntry {
				returnString = returnString + ","
			} else {
				firstEntry = false
			}
			//returnString = returnString + "{ \"" + t.Name + "\": {\"name\": \"" + t.Name + "\", \"version\": \"" + t.Version + "\"} }"
			returnString = returnString + "{\"name\": \"" + t.Name + "\", \"version\": \"" + t.Version + "\"}"
		}

		for _, p := range t.Packages {
			if !firstEntry {
				returnString = returnString + ","
			} else {
				firstEntry = false
			}
			//returnString = returnString + "{ \"" + p + "@" + t.Name + "\": {\"name\": \"" + p + "@" + t.Name + "\", \"version\": \"" + t.Version + "\"} }"
			returnString = returnString + "{\"name\": \"" + p + "@" + t.Name + "\", \"version\": \"" + t.Version + "\"}"
		}
	}

	return "[" + returnString + "]"
}

func buildManifest2() {
	var directDependenciesJSON string = "{\"version\": \"" + manifestVersion + "\", \"main\": \"" + mainModule + "\", \"packages\": ["

	var firstEntry = true
	for _, mod := range directDependencies {
		if mod.Included {
			if !firstEntry {
				directDependenciesJSON = directDependenciesJSON + ","
			} else {
				firstEntry = false
			}
			directDependenciesJSON = directDependenciesJSON + "{\"name\": \"" + mod.Name + "\", \"version\": \"" + mod.Version + "\","
			directDependenciesJSON = directDependenciesJSON + "\"dependencies\": " + flatenTransitives(mod.Transitives)
			directDependenciesJSON = directDependenciesJSON + "}"
		}

		for _, pckg := range mod.Packages {
			if !firstEntry {
				directDependenciesJSON = directDependenciesJSON + ","
			} else {
				firstEntry = false
			}
			directDependenciesJSON = directDependenciesJSON + "{\"name\": \"" + pckg.Name + "@" + mod.Name + "\", \"version\": \"" + mod.Version + "\","
			directDependenciesJSON = directDependenciesJSON + "\"dependencies\": " + flatenTransitives(pckg.Transitives)
			directDependenciesJSON = directDependenciesJSON + "}"
		}
	}
	directDependenciesJSON = directDependenciesJSON + "]}"

	f, err := os.Create(manifestFilePath)
	check(err)
	_, err = f.WriteString(string(directDependenciesJSON))
	check(err)
	f.Sync()

	defer f.Close()
	fmt.Println("Success :: Manifest generated and stored at", manifestFilePath)
}

func main() {
	if (len(os.Args) != 3) {
		fmt.Println("Error :: Invalid arguments for the command.")
		fmt.Println("Usage :: go run github.com/dgpatelgit/gobuildmanifest <Absolute source root folder path containing go.mod> <Output file path>.json")
		fmt.Println("")
		fmt.Println("Example :: go run github.com/dgpatelgit/gobuildmanifest /home/user/goproject/root/folder /home/user/gomanifest.json")
	} else {
		_, err := os.Stat(os.Args[1])
		if (err != nil) {
			fmt.Println("Invalid source folder path ::", os.Args[1])
		} else {
			fmt.Println("Building manifest file for ::", os.Args[1])
			sourceRootFolder = os.Args[1]
			manifestFilePath = os.Args[2]
			if processGraphData() == 0 {
				if processDepsData() == 0 {
					buildDirectDependencies()
					buildTransitiveDeps()
					buildManifest2()
				}
			}
		}
	}
}

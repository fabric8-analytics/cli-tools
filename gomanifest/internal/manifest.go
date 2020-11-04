package internal

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

/* MANIFEST FILE VERSION */
const manifestVersion = "v1"

// Dependency ... Direct and transitive dependency details.
type Dependency struct {
	Name         string       `json:"package"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"deps,omitempty"`
}

// Manifest ... Final manifest file structure
type Manifest struct {
	Version  string       `json:"version"`
	Main     string       `json:"main"`
	Packages []Dependency `json:"packages"`
}

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

// getTransitives ... Returns a clean list of deps
func getTransitives(deps []string, depsPackages map[string]GoPackage) []Dependency {
	var dependencies = make([]Dependency, 0)
	for _, dep := range deps {
		if depPackage, ok := depsPackages[dep]; ok {
			if depPackage.Module.Main == false {
				dependencies = append(dependencies, Dependency{
					getPackageName(depPackage),
					transformVersion(depPackage.Module.Version),
					nil,
				})
			}
		}
	}
	return dependencies
}

// BuildManifest ... Build direct & transitive dependencies.
func BuildManifest(depsPackages map[string]GoPackage) Manifest {
	var manifest Manifest = Manifest{manifestVersion, "", make([]Dependency, 0)}

	// Get direct imports from current source.
	var sourceImports = make(map[string]bool, 0)
	for _, pckg := range depsPackages {
		// Include only packages with project ROOT
		if pckg.Module.Main {
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
								Dependency{
									getPackageName(depPackage),
									transformVersion(depPackage.Module.Version),
									getTransitives(depPackage.Deps, depsPackages),
								})
						}
					}

					sourceImports[imp] = true
				}
			}
		}
	}
	log.Info().Msgf("Source code imports: %d", len(sourceImports))

	return manifest
}

// SaveManifestFile ... Save the given manifest data into a file.
func SaveManifestFile(manifest Manifest, manifestFilePath string) error {
	d, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	var directDependenciesJSON string = string(d)

	f, err := os.Create(manifestFilePath)
	if err != nil {
		return err
	}
	_, err = f.WriteString(string(directDependenciesJSON))
	if err != nil {
		return err
	}
	f.Sync()

	defer f.Close()

	return nil
}

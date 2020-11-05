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

// Save ... Save manifest object into a given save path.
func (manifest Manifest) Save(savePath string) error {
	d, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	f, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(d)
	if err != nil {
		return err
	}
	f.Sync()

	return nil
}

// transformVersion ... Converts the golang version string into semver without 'v' and appended text after '+'
func transformVersion(inVersion string) string {
	var outVersion string = strings.Replace(inVersion, "v", "", 1)

	return strings.Split(outVersion, "+")[0]
}

// getPackageName ... Utility function to convert package + module data into package name used by manifest.
func getPackageName(depPackage DepPackage) string {
	// Get module / package@module
	if depPackage.ImportPath != depPackage.Module.Path {
		return depPackage.ImportPath + "@" + depPackage.Module.Path
	}

	return depPackage.ImportPath
}

// newDirectDependency ... Return a new direct dependency for a given go package.
func newDirectDependency(depPackage DepPackage, depPackages *map[string]DepPackage) Dependency {
	return Dependency{
		getPackageName(depPackage),
		transformVersion(depPackage.Module.Version),
		getTransitives(depPackage.Deps, depPackages),
	}
}

// newTransitiveDependency ... Build and returns a new transitive dependency for a go package.
func newTransitiveDependency(depPackage DepPackage) Dependency {
	return Dependency{
		getPackageName(depPackage),
		transformVersion(depPackage.Module.Version),
		nil,
	}
}

// getTransitives ... Returns a clean list of deps
func getTransitives(deps []string, depPackages *map[string]DepPackage) []Dependency {
	var dependencies = make([]Dependency, 0)
	for _, dep := range deps {
		if depPackage, ok := (*depPackages)[dep]; ok {
			if !depPackage.Module.Main {
				dependencies = append(dependencies,
					newTransitiveDependency(depPackage))
			}
		}
	}
	return dependencies
}

// BuildManifest ... Build direct & transitive dependencies.
func BuildManifest(depPackages *map[string]DepPackage) Manifest {
	var manifest Manifest = Manifest{manifestVersion, "", make([]Dependency, 0)}

	// Get direct imports from current source.
	var sourceImports = make(map[string]bool, 0)
	for _, pckg := range *depPackages {
		// Skip dependent packages while scanning for "imports"
		if !pckg.Module.Main {
			continue
		}

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
				if depPackage, ok := (*depPackages)[imp]; ok {
					if !depPackage.Module.Main {
						manifest.Packages = append(manifest.Packages,
							newDirectDependency(depPackage, depPackages))
					}
				}

				sourceImports[imp] = true
			}
		}
	}
	log.Info().Msgf("Source code imports: \t%d", len(sourceImports))
	log.Info().Msgf("Direct dependencies: \t%d", len(manifest.Packages))

	return manifest
}

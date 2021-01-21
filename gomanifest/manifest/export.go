package manifest

import (
	"os"

	"github.com/fabric8-analytics/cli-tools/gomanifest/internal"
	"github.com/rs/zerolog/log"
)

func Generate(goExePath string, goModPath string, goManifestPath string) error {
	// Start generating manifest data.
	cmd, err := internal.RunGoList(goExePath, goModPath)
	if err != nil {
		return err
	}

	depPackages, err := internal.GetDeps(cmd)
	if err != nil {
		return err
	}

	manifest := internal.BuildManifest(&depPackages)
	// Check for empty manifest file.
	if manifest.Main == "" {
		log.Error().Msg("Empty manifest generated, correct project dependencies using `go mod tidy` command")
	}

	// Create out file.
	f, err := os.Create(goManifestPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = manifest.Write(f)
	if err != nil {
		return err
	}
	return nil
}

package main

import (
	"os"

	"github.com/fabric8-analytics/cli-tools/gomanifest/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set debug level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Validate required number of parameters.
	if len(os.Args) != 3 {
		log.Error().Msg("invalid arguments for the command")
		log.Info().Msgf("Usage :: %s <Path to source folder> <Output file path>/golist.json", os.Args[0])
		log.Info().Msgf("Example :: %s /home/user/goproject/root/folder /home/user/gomanifest.json", os.Args[0])
		os.Exit(1)
	}

	// Ensure source path exists.
	_, err := os.Stat(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("invalid")
		os.Exit(2)
	}

	// Start generating manifest data.
	log.Info().Msgf("Started analysing go project at %s", os.Args[1])
	cmd, err := internal.RunGoList(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("`go list` failed")
		os.Exit(3)
	}

	depPackages, err := internal.GetDeps(cmd)
	if err != nil {
		log.Error().Err(err).Msg("get deps")
		os.Exit(3)
	}

	manifest := internal.BuildManifest(&depPackages)
	// Check for empty manifest file.
	if manifest.Main == "" {
		log.Error().Msg("Empty manifest generated, correct project dependencies using `go mod tidy` command")
		os.Exit(4)
	}

	// Create out file.
	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Error().Err(err).Msg("unable to write")
		os.Exit(4)
	}
	defer f.Close()

	err = manifest.Write(f)
	if err != nil {
		log.Error().Err(err).Msg("unable to write")
		os.Exit(5)
	}

	log.Info().Msgf("Manifest file generated and stored at %s", os.Args[2])
}

package main

import (
	"context"
	"os"

	"github.com/fabric8-analytics/cli-tools/gomanifest/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	// Set debug level to info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	ctx = log.Logger.WithContext(ctx)

	// Validate required number of parameters.
	if len(os.Args) != 3 {
		log.Error().Msg("invalid arguments for the command")
		log.Info().Msg("Usage :: go run github.com/fabric8-analytics/cli-tools/gomanifest <Path to source folder> <Output file path>/golist.json")
		log.Info().Msg("Example :: go run github.com/fabric8-analytics/cli-tools/gomanifest /home/user/goproject/root/folder /home/user/gomanifest.json")
		os.Exit(-1)
	}

	// Ensure source path exists.
	_, err := os.Stat(os.Args[1])
	if err != nil {
		log.Error().Msgf("Invalid source folder path: %s", os.Args[1])
		os.Exit(-2)
	}

	// Start generating manifest data.
	log.Info().Msgf("Started analysing go project at %s", os.Args[1])
	goList := &internal.GoList{Command: &internal.GoListCmd{CWD: os.Args[1]}}
	depPackages, err := goList.Get()
	if err != nil {
		log.Error().Msgf("Exception raised: %v", err)
		os.Exit(-3)
	}

	manifest := internal.BuildManifest(&depPackages)
	err = manifest.Save(os.Args[2])
	if err != nil {
		log.Error().Msgf("Exception raised: %v", err)
		os.Exit(-4)
	}

	log.Info().Msgf("Manifest file generated and stored at %s", os.Args[2])
	os.Exit(0)
}

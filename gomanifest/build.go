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

	if len(os.Args) != 3 {
		log.Error().Msg("invalid arguments for the command")
		log.Info().Msg("Usage :: go run github.com/fabric8-analytics/cli-tools/gomanifest <Path to source folder> <Output file path>/golist.json")
		log.Info().Msg("Example :: go run github.com/fabric8-analytics/cli-tools/gomanifest /home/user/goproject/root/folder /home/user/gomanifest.json")
	} else {
		_, err := os.Stat(os.Args[1])
		if err != nil {
			log.Error().Msgf("Invalid source folder path: %s", os.Args[1])
		} else {

			goListCmd := &internal.GoListCmd{CWD: os.Args[1]}
			goList := &internal.GoList{Command: goListCmd}
			if depsPackages, err := goList.Get(); err != nil {
				log.Error().Msgf("Exception raised: %v", err)
			} else {
				if err := internal.SaveManifestFile(internal.BuildManifest(depsPackages), os.Args[2]); err != nil {
					log.Error().Msgf("Exception raised: %v", err)
				} else {
					log.Info().Msgf("Manifest file generated and stored at %s", os.Args[2])
				}
			}
		}
	}
}

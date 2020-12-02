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
	argsLen := len(os.Args)
	if argsLen != 3 && argsLen != 4 {
		log.Error().Msg("invalid arguments for the command")
		log.Info().Msgf("Usage :: go run github.com/fabric8-analytics/cli-tools/gomanifest <Path to source folder> <Output file path>/golist.json [<Go executable path>]")
		log.Info().Msgf("Example :: go run github.com/fabric8-analytics/cli-tools/gomanifest /home/user/goproject/root/folder /home/user/golist.json /usr/local/bin/go")
		os.Exit(1)
	}

	// Ensure source path exists.
	_, err := os.Stat(os.Args[1])
	if err != nil {
		log.Error().Err(err).Msg("invalid")
		os.Exit(2)
	}

	// Extract go executable path, if given
	goExec := "go"
	if argsLen == 4 {
		goExec = os.Args[3]
	}

	// Start generating manifest data.
	log.Info().Msgf("Started analysing go project at [%s] using go exec [%s]", os.Args[1], goExec)
	cmd, err := internal.RunGoList(os.Args[1], goExec)
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

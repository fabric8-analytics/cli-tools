package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonOut bool
var verboseOut bool

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "Get detailed report of vulnerabilities.",
	Long: `Get detailed report of vulnerabilities. Supported ecosystems are Pypi (Python), Maven (Java), Npm (Node) and Golang (Go).
If stack has Vulnerabilities, command will exit with status code 2.`,
	Args:     validateFileArg,
	PostRunE: destructor,
	RunE:     runAnalyse,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Set output format to JSON.")
	analyseCmd.Flags().BoolVarP(&verboseOut, "verbose", "v", false, "Detailed Analyses Report.")
}

// destructor deletes intermediary files used to have stack analyses
func destructor(_ *cobra.Command, _ []string) error {
	log.Debug().Msgf("Running Destructor.\n")
	if debug {
		// Keep intermediary files, when on debug
		log.Debug().Msgf("Skipping file clearance on Debug Mode.\n")
		return nil
	}
	intermediaryFiles := []string{"generate_pylist.py", "pylist.json", "dependencies.txt", "golist.json", "npmlist.json"}
	for _, file := range intermediaryFiles {
		file = filepath.Join(os.TempDir(), file)
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				// If file doesn't exists, continue
				continue
			}
			return err
		}
		e := os.Remove(file)
		if e != nil {
			return fmt.Errorf("error clearing files %s", file)
		}
	}
	return nil
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) error {
	telemetry.SetFlag(cmd.Context(), "json", jsonOut)
	telemetry.SetFlag(cmd.Context(), "verbose", verboseOut)
	if !viper.IsSet("crda_key") {
		telemetry.SetExitCode(cmd.Context(), 1)
		return errors.New(
			"please run `crda auth` command first")
	}
	telemetry.SetCrdaKey(cmd.Context(), viper.GetString("crda_key"))
	manifestPath := getAbsPath(args[0])
	requestParams := driver.RequestType{
		UserID:          viper.GetString("crda_key"),
		ThreeScaleToken: viper.GetString("auth_token"),
		Host:            viper.GetString("host"),
		RawManifestFile: manifestPath,
	}
	if !jsonOut {
		fmt.Println("Analysing your Dependency Stack! Please wait...")
	}
	name := sa.GetManifestName(manifestPath)
	hasVul, err := sa.StackAnalyses(cmd.Context(), requestParams, jsonOut, verboseOut)
	telemetry.SetManifest(cmd.Context(), name)
	if err != nil {
		telemetry.SetExitCode(cmd.Context(), 1)
		return err
	}
	if hasVul {
		// Stack has vulnerability, exit with 2 code
		exitCode = 2
		telemetry.SetExitCode(cmd.Context(), exitCode)
	}
	return nil
}

// getAbsPath converts relative path to Abs
func getAbsPath(givenPath string) string {
	manifestPath, err := filepath.Abs(givenPath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Unable to convert to Absolute file path.")
	}
	return manifestPath
}

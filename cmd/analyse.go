package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/fatih/color"
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
	Run:     runAnalyse,
	Args:    validateFileArg,
	PostRun: destructor,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Set output format to JSON.")
	analyseCmd.Flags().BoolVarP(&verboseOut, "verbose", "v", false, "Detailed Analyses Report.")
}

// destructor deletes intermediary files used to have stack analyses
func destructor(cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Running Destructor.\n")
	if debug {
		// Keep intermediary files, when on debug
		log.Debug().Msgf("Skipping file clearance on Debug Mode.\n")
		return
	}
	intermediaryFiles := []string{"generate_pylist.py", "pylist.json", "dependencies.txt", "golist.json", "npmlist.json"}
	for _, file := range intermediaryFiles {
		file = filepath.Join(os.TempDir(), file)
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				// If file doesn't exists, continue
				continue
			}
		}
		e := os.Remove(file)
		if e != nil {
			log.Fatal().Msgf("Error clearing files %s", file)
		}
	}
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) {
	if !viper.IsSet("crda_key") {
		fmt.Fprintln(os.Stdout,
			color.RedString("\u2718 "),
			"Please run `crda auth` command first.",
		)
		os.Exit(1)
	}
	requestParams := driver.RequestType{
		UserID:          viper.GetString("crda_key"),
		ThreeScaleToken: viper.GetString("auth_token"),
		Host:            viper.GetString("host"),
		RawManifestFile: args[0],
	}
	if !jsonOut {
		fmt.Fprintln(os.Stdout, "Analysing your Dependency Stack! Please wait...")
	}
	hasVul := sa.StackAnalyses(requestParams, jsonOut, verboseOut)
	if hasVul && jsonOut {
		// Stack has vulnerability, exit with 2 code
		os.Exit(2)
	}
}

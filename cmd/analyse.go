package cmd

import (
	"os"
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifestFile string
var jsonOut bool
var flagNoColor bool

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "Provides detailed report of vulnerabilities.",
	Long: `Provides detailed report of vulnerabilities. Supported ecosystems are Pypi (Python), Maven (Java), Npm (Node) and Golang (Go).
If stack has Vulnerabilities, command will exit with status code 2.`,
	Run:     runAnalyse,
	PostRun: destructor,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.PersistentFlags().StringVarP(&manifestFile, "file", "f", "", "Manifest file absolute path.")
	analyseCmd.MarkPersistentFlagRequired("file")
	analyseCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "Set output format to JSON.")
	analyseCmd.Flags().BoolVarP(&flagNoColor, "no-color", "c", false, "Toggle colors in output.")
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
	if flagNoColor {
		color.NoColor = true
	}
	requestParams := driver.RequestType{
		UserID:          viper.GetString("crda-key"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
		RawManifestFile: manifestFile,
	}
	if !jsonOut {
		log.Info().Msgf("Executing Stack Analyses! Please wait... ")
	}
	hasVul := sa.StackAnalyses(requestParams, jsonOut)
	if hasVul && jsonOut {
		// Stack has vulnerability, exit with 2 code
		os.Exit(2)
	}
}

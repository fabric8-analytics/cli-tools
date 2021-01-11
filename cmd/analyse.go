package cmd

import (
	"os"
	"path/filepath"

	"github.com/fabric8-analytics/cli-tools/analyses/driver"
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifestFile string

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:     "analyse",
	Short:   "analyse performs Stack Analyses",
	Long:    `analyse performs full Stack Analyses. Supported Ecosystems are Pypi, Maven, npm and golang.`,
	Run:     runAnalyse,
	PostRun: destructor,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.PersistentFlags().StringVarP(&manifestFile, "file", "f", "", "Manifest file absolute path.")
	analyseCmd.MarkPersistentFlagRequired("file")
}

// destructor deletes intermediatery files used to have stack analyses
func destructor(cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Running Destructor.\n")
	if debug {
		// Keep intermediatery files, when on debug
		return
	}
	intermediataryFiles := []string{"generate_pylist.py", "pylist.json", "dependencies.txt", "golist.json", "npmlist.json"}
	for _, file := range intermediataryFiles {
		file = filepath.Join("/tmp", file)
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
	log.Debug().Msgf("Completed Destructor.\n")
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) {
	requestParams := driver.RequestType{
		UserID:          viper.GetString("crda-key"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
		RawManifestFile: manifestFile,
	}
	stackAnalysesResponse := sa.StackAnalyses(requestParams)
	log.Info().Msgf("Stack Analyses Response:\n %s \n\n", stackAnalysesResponse.AnalysedDeps)
	log.Info().Msgf("Successfully completed Stack Analyses.\n")
}

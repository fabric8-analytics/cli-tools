package cmd

import (
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var manifestFile string

// analyseCmd represents the analyse command
var analyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "analyse performs Stack Analyses",
	Long:  `analyse performs full Stack Analyses. Supported Ecosystems are Pypi, Maven, npm and golang.`,
	Run:   runAnalyse,
}

func init() {
	rootCmd.AddCommand(analyseCmd)
	analyseCmd.PersistentFlags().StringVarP(&manifestFile, "file", "f", "", "Manifest file absolute path.")
	analyseCmd.MarkPersistentFlagRequired("file")
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) {
	requestParams := sa.RequestType{
		UserID:          viper.GetString("crda-key"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
		RawManifestFile: manifestFile,
	}
	saResponse := sa.StackAnalyses(requestParams)
	log.Info().Msgf("Stack Analyses Response:\n %s \n\n", saResponse.AnalysedDeps)
	log.Info().Msgf("Successfully completed Stack Analyses.\n")
}

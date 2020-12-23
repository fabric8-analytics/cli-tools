package cmd

import (
	sa "github.com/fabric8-analytics/cli-tools/analyses/stackanalyses"
	constants "github.com/fabric8-analytics/cli-tools/utils"
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
	analyseCmd.PersistentFlags().String("shell-path", constants.Shell, "Shell Path.")
	analyseCmd.MarkPersistentFlagRequired("file")
	viper.BindPFlag("shell-path", analyseCmd.PersistentFlags().Lookup("shell-path"))
}

//runAnalyse is controller func for analyses cmd.
func runAnalyse(cmd *cobra.Command, args []string) {
	requestParams := sa.SARequestType{
		UserID:          viper.GetString("crda-key"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
		ShellPath:       viper.GetString("shell-path"),
		RawManifestFile: manifestFile,
	}
	saResponse := sa.StackAnalyses(requestParams)
	log.Info().Msgf("Stack Analyses Response:\n %s \n\n", saResponse.AnalysedDeps)
	log.Info().Msgf("Successfully completed Stack Analyses.\n")
}

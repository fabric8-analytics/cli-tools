package cmd

import (
	auth "github.com/fabric8-analytics/cli-tools/auth"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showUUID bool

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:              "auth",
	Short:            "Links uuid with Snyk token.",
	Long:             `command maps Snyk Token with UUID on crda server and outputs UUID to use for Authentication.`,
	Run:              main,
	PersistentPreRun: preRun,
}

func preRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("snyk-token", rootCmd.PersistentFlags().Lookup("snyk-token"))
}

func init() {
	authCmd.Flags().BoolVarP(&showUUID, "show-token", "s", false, "Show token to STDOUT")
	rootCmd.AddCommand(authCmd)
}

func main(cmd *cobra.Command, args []string) {
	log.Debug().Msgf("Executing Auth command.")
	requestParams := auth.RequestServerType{
		UserID:          viper.GetString("crda-key"),
		SynkToken:       viper.GetString("snyk-token"),
		ThreeScaleToken: viper.GetString("auth-token"),
		Host:            viper.GetString("host"),
	}
	userID := auth.RequestServer(requestParams)

	log.Info().Msgf("Successfully Registered.\n")
	if showUUID {
		log.Info().Msgf("Please update CI env with:")
		log.Info().Msgf("crda-key: %s\n", userID)
		log.Info().Msgf("This token is confidential and exculsive to your Snyk Id.")
	}
	viper.Set("crda-key", userID)
	viper.WriteConfig()
	log.Debug().Msgf("Successfully Executed Auth command.")
}

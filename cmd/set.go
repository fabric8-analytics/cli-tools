package cmd

import (
	"fmt"
	"github.com/fabric8-analytics/cli-tools/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set CONFIG-KEY VALUE",
	Short: "Set a crda configuration property",
	Long:  `Sets a crda configuration property`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Error().Msgf("please provide a configuration property and its value as in 'crda config set CONFIG-KEY VALUE'")
			return
		}
		viper.Set(args[0], args[1])
		err := viper.WriteConfig()
		if err != nil {
			log.Error().Msgf("unable to write config health")
			return
		}
		value := viper.Get(args[0])
		config.ViperUnMarshal()
		if value != args[1] {
			log.Error().Msgf("unable to set configuration value")
			return
		}
		fmt.Println("successfully set configuration value")
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}

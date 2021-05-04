package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get CONFIG-KEY",
	Short: "Get active crda configurations",
	Long:  `Get active crda configurations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			c := viper.AllSettings()
			bs, err := yaml.Marshal(c)
			if err != nil {
				log.Error().Msgf("unable to marshal config to YAML")
				return err
			}
			fmt.Println(string(bs))
			return nil
		}
		key := args[0]
		v := viper.IsSet(args[0])
		if !v {
			return fmt.Errorf("configuration property '%s' does not exist", key)
		}
		value := viper.Get(args[0])
		fmt.Println(key, ":", value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(getCmd)
}

package cmd

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	constants "github.com/fabric8-analytics/cli-tools/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags
var (
	debug       bool
	cfgFile     string
	flagNoColor bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crda",
	Short: "Cli tool to interact with CRDA Platform.",
	Long: `Cli tool to interact with CRDA Platform. This tool performs token Authentication and verbose Analyses of Dependency Stack. 
	
	Authenticated token can be used as Auth Token to interact with CRDA Platform.`,
	Args: cobra.ExactValidArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msgf("Error Executing crda command. Please raise at https://github.com/fabric8-analytics/cli-tools/issues, if issue persists.")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", constants.Debug, "Sets Log level to Debug.")
	rootCmd.PersistentFlags().BoolVarP(&flagNoColor, "no-color", "c", false, "Toggle colors in output.")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	color.NoColor = flagNoColor

	// Log Level Settings
	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// config.yaml Settings
	log.Debug().Msgf("Setting Configuration files")
	configName := "config"
	configHome, err := homedir.Dir()
	crdaFolder := ".crda"
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
	configPath := filepath.Join(configHome, crdaFolder)
	configFilePath := filepath.Join(configPath, configName)
	// Search config file in home path.
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv() // read in environment variables that match

	//Handle Reading Config Files Error
	if err := viper.ReadInConfig(); err != nil {
		log.Debug().Msgf("Error reading config file.")
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, Creating one
			if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
				log.Debug().Msgf("Creating config file.")
				os.MkdirAll(configPath, os.ModePerm)
				_, err := os.Create(configFilePath + ".yaml")
				if err != nil {
					log.Fatal().Err(err).Msgf(err.Error())
				}
			}
			viper.SetConfigFile(configFilePath + ".yaml")
		} else {
			// Config file was found but another error was produced
			log.Fatal().Err(err).Msgf(err.Error())
		}
	}
	if !viper.IsSet("host") {
		viper.Set("host", constants.Host)
	}
	if !viper.IsSet("auth-token") {
		viper.Set("auth-token", constants.AuthToken)
	}
	viper.WriteConfig()
	log.Debug().Msgf("Using config file %s.\n", viper.ConfigFileUsed())
	log.Debug().Msgf("Successfully configured config files %s.", viper.ConfigFileUsed())
}

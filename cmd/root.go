package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	crdaConfig "github.com/fabric8-analytics/cli-tools/pkg/config"
	"github.com/fabric8-analytics/cli-tools/pkg/segment"
	"github.com/fabric8-analytics/cli-tools/pkg/telemetry"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fabric8-analytics/cli-tools/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flags
var (
	debug       bool
	flagNoColor bool
	client      string
)

// Variables
var (
	segmentClient *segment.Client
	exitCode      = 0
	ctx           context.Context
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crda",
	Short: "Cli tool to interact with CRDA Platform.",
	Long: `Cli tool to interact with CRDA Platform. This tool performs token Authentication and verbose Analyses of Dependency Stack. 
	
	Authenticated token can be used as Auth Token to interact with CRDA Platform.`,
	Args:          cobra.ExactValidArgs(1),
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	attachMiddleware([]string{}, rootCmd)
	ctx = telemetry.NewContext(context.Background())
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		// CLI Errors
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		exitCode = 1
	}
	err := segmentClient.Close()
	if err != nil {
		return
	}
	os.Exit(exitCode)
}

func init() {
	var err error
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", utils.Debug, "Sets Log level to Debug.")
	rootCmd.PersistentFlags().BoolVarP(&flagNoColor, "no-color", "c", false, "Toggle colors in output.")
	rootCmd.PersistentFlags().StringVarP(&client, "client", "m", "terminal", "Client [tekton/jenkins/gh-actions]")

	// Initiate segment client
	if segmentClient, err = segment.NewClient(); err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
	}
}

func executeWithLogging(fullCmd string, input func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		log.Debug().Msgf("Running '%s'", fullCmd)
		startTime := time.Now()
		err := input(cmd, args)
		pushToSegment(fullCmd, startTime, err)
		return err
	}
}

func attachMiddleware(names []string, cmd *cobra.Command) {
	if cmd.HasSubCommands() {
		for _, command := range cmd.Commands() {
			attachMiddleware(append(names, cmd.Name()), command)
		}
	} else if cmd.RunE != nil {
		fullCmd := strings.Join(append(names, cmd.Name()), " ")
		src := cmd.RunE
		cmd.RunE = executeWithLogging(fullCmd, src)
	}
}

func pushToSegment(event string, startTime time.Time, err error) {
	if serr := segmentClient.Upload(rootCmd.Context(), event, time.Since(startTime), err); serr != nil {
		log.Info().Msgf("Cannot send data to telemetry: %v", serr)
	}
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
				err := os.MkdirAll(configPath, os.ModePerm)
				if err != nil {
					return
				}
				_, err = os.Create(configFilePath + ".yaml")
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
	if !viper.IsSet("crda_host") {
		viper.Set("crda_host", utils.CRDAHost)
	}
	if !viper.IsSet("crda_auth_token") {
		viper.Set("crda_auth_token", utils.CRDAAuthToken)
	}

	err = viper.WriteConfig()
	if err != nil {
		return
	}
	err = validateFlagValues(client)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	telemetry.SetClient(ctx, client)
	crdaConfig.ViperUnMarshal()
	log.Debug().Msgf("Using config file %s.\n", viper.ConfigFileUsed())
	log.Debug().Msgf("Successfully configured config files %s.", viper.ConfigFileUsed())
}

// askTelemetryConsent fires Telemetry Consent Prompt
func askTelemetryConsent() {
	if !viper.IsSet("consent_telemetry") {
		response := telemetry.GetTelemetryConsent()
		viper.Set("consent_telemetry", strconv.FormatBool(response))
		err := viper.WriteConfig()
		if err != nil {
			log.Error().Msgf("unable to write config")
			return
		}
	}
}

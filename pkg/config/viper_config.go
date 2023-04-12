package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type viperConfigs struct {
	Host             string `mapstructure:"CRDA_HOST" yaml:"crda_host"`
	AuthToken        string `mapstructure:"CRDA_AUTH_TOKEN" yaml:"crda_auth_token"`
	CrdaKey          string `mapstructure:"CRDA_KEY" yaml:"crda_key"`
	ConsentTelemetry string `mapstructure:"CONSENT_TELEMETRY" yaml:"consent_telemetry"`
}

// ActiveConfigValues Maintain state of viper configurations
var ActiveConfigValues = &viperConfigs{}

// ViperUnMarshal loads viper configs into ActiveConfigValues
func ViperUnMarshal() *viperConfigs {
	viper.AutomaticEnv()
	// Have to bind individual variables: https://github.com/spf13/viper/issues/188
	_ = viper.BindEnv("CONSENT_TELEMETRY")
	_ = viper.BindEnv("CRDA_AUTH_TOKEN")
	_ = viper.BindEnv("CRDA_HOST")
	_ = viper.BindEnv("CRDA_KEY")
	err := viper.Unmarshal(&ActiveConfigValues)
	if err != nil {
		log.Fatal().Msgf("unable to decode into struct, %v", err)
	}
	return ActiveConfigValues
}

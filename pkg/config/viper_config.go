package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type viperConfigs struct {
	Host             string `mapstructure:"host"`
	AuthToken        string `mapstructure:"auth_token"`
	ConsentTelemetry string `mapstructure:"consent_telemetry"`
}

var ActiveConfigValues = &viperConfigs{}

func ViperUnMarshal() {

	err := viper.Unmarshal(&ActiveConfigValues)
	if err != nil {
		log.Fatal().Msgf("unable to decode into struct, %v", err)
	}
}

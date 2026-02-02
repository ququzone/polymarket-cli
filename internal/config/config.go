package config

import (
	"github.com/spf13/viper"
)

type BuilderConfig struct {
	APIKey     string `mapstructure:"api_key"`
	Passphrase string `mapstructure:"passphrase"`
	APISecret  string `mapstructure:"api_secret"`
}

type Config struct {
	Builder        BuilderConfig `mapstructure:"builder"`
	DataAPIBaseURL string        `mapstructure:"data_api_base_url"`
	PrivateKey     string        `mapstructure:"private_key"`
}

var AppCfg *Config

func Init() {
	AppCfg = &Config{
		Builder: BuilderConfig{
			APIKey:     viper.GetString("builder.api_key"),
			Passphrase: viper.GetString("builder.passphrase"),
			APISecret:  viper.GetString("builder.api_secret"),
		},
		DataAPIBaseURL: viper.GetString("data_api_base_url"),
		PrivateKey:     viper.GetString("private_key"),
	}

	if AppCfg.DataAPIBaseURL == "" {
		AppCfg.DataAPIBaseURL = "https://data-api.polymarket.com"
	}
}

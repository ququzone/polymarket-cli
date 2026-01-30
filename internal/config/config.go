package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	APIKey       string
	APISecret    string
	Debug        bool
	OutputFormat string // json, table, etc.
}

var AppCfg *Config

func Init() {
	AppCfg = &Config{
		APIKey:       viper.GetString("api_key"),
		APISecret:    viper.GetString("api_secret"),
		Debug:        viper.GetBool("debug"),
		OutputFormat: viper.GetString("output_format"),
	}

	if AppCfg.OutputFormat == "" {
		AppCfg.OutputFormat = "table"
	}
}

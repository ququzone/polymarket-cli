package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	APIKey         string
	APISecret      string
	DataAPIBaseURL string
}

var AppCfg *Config

func Init() {
	AppCfg = &Config{
		APIKey:         viper.GetString("api_key"),
		APISecret:      viper.GetString("api_secret"),
		DataAPIBaseURL: viper.GetString("data_api_base_url"),
	}

	if AppCfg.DataAPIBaseURL == "" {
		AppCfg.DataAPIBaseURL = "https://data-api.polymarket.com"
	}
}

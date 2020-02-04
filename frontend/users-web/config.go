package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ApiAddress string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("ApiAddress", "API_ADDRESS")
	viper.BindEnv("ServerPort", "SERVER_PORT")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

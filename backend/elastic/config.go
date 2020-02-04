package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	ElasticAddress string
	KubeMQHost     string
	KubeMQPort     int
	Channel        string
	Group          string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("ElasticAddress", "ELASTIC_ADDRESS")
	viper.BindEnv("KubeMQHost", "KUBEMQ_HOST")
	viper.BindEnv("KubeMQPort", "KUBEMQ_POST")
	viper.BindEnv("Channel", "CHANNEL")
	viper.BindEnv("Group", "GROUP")

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

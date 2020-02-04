package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress   string
	KubeMQHost     string
	KubeMQPort     int
	Channel        string
	Group          string
	HistoryChannel string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("RedisAddress", "REDIS_ADDRESS")
	viper.BindEnv("KubeMQHost", "KUBEMQ_HOST")
	viper.BindEnv("KubeMQPort", "KUBEMQ_POST")
	viper.BindEnv("Channel", "CHANNEL")
	viper.BindEnv("Group", "GROUP")
	viper.BindEnv("HistoryChannel", "HISTORY_CHANNEL")

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

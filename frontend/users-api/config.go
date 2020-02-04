package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	KubeMQHost     string
	KubeMQPort     int
	UsersChannel   string
	CacheChannel   string
	AuditChannel   string
	HistoryChannel string
	LogsChannel    string
	ConfigChannel  string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("Port", "PORT")
	viper.BindEnv("KubeMQHost", "KUBEMQ_HOST")
	viper.BindEnv("KubeMQPort", "KUBEMQ_PORT")
	viper.BindEnv("UsersChannel", "USERS_CHANNEL")
	viper.BindEnv("CacheChannel", "CACHE_CHANNEL")
	viper.BindEnv("AuditChannel", "AUDIT_CHANNEL")
	viper.BindEnv("HistoryChannel", "HISTORY_CHANNEL")
	viper.BindEnv("LogsChannel", "LOGS_CHANNEL")
	viper.BindEnv("ConfigChannel", "CONFIG_CHANNEL")

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

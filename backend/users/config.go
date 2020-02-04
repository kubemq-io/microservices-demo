package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	KubeMQHost          string
	KubeMQPort          int
	PostgresHost        string
	PostgresPort        int
	PostgresUser        string
	PostgresPassword    string
	PostgresDB          string
	UsersChannel        string
	CacheChannel        string
	AuditChannel        string
	HistoryChannel      string
	LogsChannel         string
	ConfigChannel       string
	NotificationChannel string
	Group               string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	viper.BindEnv("PostgresHost", "POSTGRES_HOST")
	viper.BindEnv("PostgresPort", "POSTGRES_PORT")
	viper.BindEnv("PostgresUser", "POSTGRES_USER")
	viper.BindEnv("PostgresPassword", "POSTGRES_PASSWORD")
	viper.BindEnv("PostgresDB", "POSTGRES_DB")
	viper.BindEnv("KubeMQHost", "KUBEMQ_HOST")
	viper.BindEnv("KubeMQPort", "KUBEMQ_POST")
	viper.BindEnv("UsersChannel", "USERS_CHANNEL")
	viper.BindEnv("CacheChannel", "CACHE_CHANNEL")
	viper.BindEnv("AuditChannel", "AUDIT_CHANNEL")
	viper.BindEnv("HistoryChannel", "HISTORY_CHANNEL")
	viper.BindEnv("LogsChannel", "LOGS_CHANNEL")
	viper.BindEnv("ConfigChannel", "CONFIG_CHANNEL")
	viper.BindEnv("NotificationChannel", "NOTIFICATION_CHANNEL")
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

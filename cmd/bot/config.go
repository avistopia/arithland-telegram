package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Logger   LoggerConfig
	Database DatabaseConfig
	Telegram TelegramConfig
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Type string
	DSN  string
}

type TelegramConfig struct {
	Token string
}

func loadConfig() (*Config, error) {
	viper.SetDefault("Logger.Level", "info")

	viper.SetDefault("Database.Type", "sqlite")
	viper.SetDefault("Database.DSN", "data/db.sqlite")
	viper.SetDefault("Telegram.Token", "environment variable should be configured")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

package config

import (
	"os"
)

type Config struct {
	ApiPort				string
	DebugMode			string

	JWT					string
}

func Load() (*Config, error) {
	var conf *Config = &Config {
		ApiPort: os.Getenv("API_PORT"),
		DebugMode: os.Getenv("GIN_MODE"),
		JWT: os.Getenv("JWT_SECRET"),
	}

	return conf, nil
}
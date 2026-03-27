package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL		string
	DatabasePort	string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if (err != nil) {
		log.Println("[WARNING] .env file not found")
	}

	var config *Config = &Config {
		DatabaseURL: os.Getenv("DATABASE_URL"),
		DatabasePort: os.Getenv("DATABASE_PORT"),
	}

	return config, nil
}
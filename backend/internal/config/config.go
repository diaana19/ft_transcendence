package config

import (
	"os"

	_ "github.com/lib/pq" // registers the Postgres driver for database/sql
)

type Postgres struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Redis struct {
	Host string
	Port string
}

type Config struct {
	JWT                string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string
	FrontendURL        string
	Postgres           Postgres
	Redis              Redis
}

func Load() (*Config, error) {
	return &Config{
		JWT:                os.Getenv("JWT_SECRET"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		FrontendURL:        os.Getenv("FT_TRANSCENDENCE_URL"),
		Postgres: Postgres{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Redis: Redis{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}, nil
}

package config

import (
	"os"

	_ "github.com/lib/pq" // registers the Postgres driver for database/sql
)

type Config struct {
	JWT                string
	GithubClientID     string
	GithubClientSecret string
	GithubRedirectURL  string
	FrontendURL        string
}

func Load() (*Config, error) {
	conf := &Config{
		JWT:                os.Getenv("JWT_SECRET"),
		GithubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		FrontendURL:        os.Getenv("FT_TRANSCENDENCE_URL"),
	}

	return conf, nil
}

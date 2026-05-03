package config

import (
	"os"
	_ "github.com/lib/pq"
)

type Config struct {
	ApiPort				string
	DebugMode			string

	JWT					string
	GithubClientID		string
	GithubClientSecret	string
	GithubRedirectURL	string
	FrontendURL			string
}

func Load() (*Config, error) {
	var conf *Config = &Config {
		ApiPort: os.Getenv("API_PORT"),
		DebugMode: os.Getenv("GIN_MODE"),
		JWT: os.Getenv("JWT_SECRET"),
		GithubClientID: os.Getenv("GITHUB_CLIENT_ID"),
		GithubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL: os.Getenv("GITHUB_REDIRECT_URL"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}

	return conf, nil
}



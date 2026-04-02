package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DatabaseName		string
	DatabaseHost		string
	DatabasePort		string
	DatabaseUser		string
	DatabasePassword	string

	ApiPort				string
	DebugMode			string

	JWT					string
}

func Load() (*sql.DB, *Config, error) {
	var config *Config = &Config {
		DatabaseName: os.Getenv("DB_NAME"),
		DatabaseHost: os.Getenv("DB_HOST"),
		DatabasePort: os.Getenv("DB_PORT"),
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
		ApiPort: os.Getenv("API_PORT"),
		DebugMode: os.Getenv("GIN_MODE"),
		JWT: os.Getenv("JWT_SECRET"),
	}

	connString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        config.DatabaseHost, config.DatabasePort, config.DatabaseUser, config.DatabasePassword, config.DatabaseName,
    )

	DB, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	return DB, config, nil
}
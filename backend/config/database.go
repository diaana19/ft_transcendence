package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type DBConfig struct {
	DatabaseName		string
	DatabaseHost		string
	DatabasePort		string
	DatabaseUser		string
	DatabasePassword	string
}

func ConnectDB() (*sql.DB, error) {
	var conf *DBConfig = &DBConfig {
		DatabaseName: os.Getenv("DB_NAME"),
		DatabaseHost: os.Getenv("DB_HOST"),
		DatabasePort: os.Getenv("DB_PORT"),
		DatabaseUser: os.Getenv("DB_USER"),
		DatabasePassword: os.Getenv("DB_PASSWORD"),
	}

	connString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        conf.DatabaseHost, conf.DatabasePort, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName,
    )

	DB, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	defer DB.Close()

	return DB, nil
}
package config

import (
	"fmt"
	"os"

	"github.com/Transcendence/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DatabaseName		string
	DatabaseHost		string
	DatabasePort		string
	DatabaseUser		string
	DatabasePassword	string
}

func ConnectDB() (*gorm.DB, error) {
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
	
	DB, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	errA := DB.AutoMigrate(&models.User{})
	if errA != nil {
		return nil, errA
	}
	return DB, nil
}
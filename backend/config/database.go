package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Transcendence/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DatabaseName     string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
}

func ConnectDB() (*gorm.DB, error) {
	godotenv.Load(".env")

	conf := &DBConfig{
		DatabaseName:     os.Getenv("DB_NAME"),
		DatabaseHost:     os.Getenv("DB_HOST"),
		DatabasePort:     os.Getenv("DB_PORT"),
		DatabaseUser:     os.Getenv("DB_USER"),
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

	log.Println("Running AutoMigrate…")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.Post{},
		&models.Like{},
		&models.Reply{},
		&models.Repost{},
		// &models.Group{},
		&models.Message{},
		&models.Notification{},
	)
	if err != nil {
		log.Printf("AutoMigrate error: %v\n", err)
		return nil, err
	}
	log.Println("AutoMigrate completed successfully")

	return DB, nil
}

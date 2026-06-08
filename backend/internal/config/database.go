package config

import (
	"fmt"
	"log"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

func (pg *Postgres) Connect() (*gorm.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(pg.User),
		url.QueryEscape(pg.Password),
		pg.Host, pg.Port, pg.Name,
	)

	DB, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	log.Println("Running AutoMigrate…")
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.Post{},
		&models.PostReaction{},
		&models.Reply{},
		&models.ReplyReaction{},
		&models.Repost{},
		&models.Message{},
		&models.Notification{},
		&models.File{},
		&models.FileAccess{},
	); err != nil {
		return nil, fmt.Errorf("auto-migrate: %w", err)
	}
	log.Println("AutoMigrate completed successfully")

	return DB, nil
}

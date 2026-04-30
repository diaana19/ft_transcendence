package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/Transcendence/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type NotificationService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewNotificationService(db *gorm.DB, rdb *redis.Client) *NotificationService {
	return &NotificationService{db: db, rdb: rdb}
}

func (s *NotificationService) SendNotification(userID, userUsername, actorID, actorUsername, notifType, content string) error {
	var receiver models.User
	if err := s.db.Select("username").First(&receiver, "id = ?", userID).Error; err != nil {
		log.Printf("Warning: could not fetch receiver username: %v", err)
	}
	notif := models.Notification{
		ID:            uuid.New().String(),
		CreatedAt:     time.Now(),
		UserID:        userID,
		UserUsername:  receiver.Username,
		ActorID:       actorID,
		ActorUsername: actorUsername,
		Type:          notifType,
		Content:       content,
		Read:          false,
	}

	if err := s.db.Create(&notif).Error; err != nil {
		log.Printf("Error saving notification: %v", err)
		return err
	}

	payload, err := json.Marshal(map[string]interface{}{
		"type":         "notification",
		"notification": notif,
	})
	if err != nil {
		return err
	}

	s.rdb.Publish(context.Background(), "notifications:"+userID, string(payload))
	return nil
}

func (s *NotificationService) GetUnread(userID string) ([]models.Notification, error) {
	var notifs []models.Notification
	err := s.db.Where("user_id = ? AND read = false", userID).
		Order("created_at desc").
		Find(&notifs).Error
	return notifs, err
}

func (s *NotificationService) MarkAllRead(userID string) error {
	return s.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).
		Update("read", true).Error
}

package services

import (
	"context"
	"log"
	"time"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/google/uuid"
)

type NotificationService struct {
	// db  *gorm.DB
	// rdb *redis.Client
	repo   *repositories.NotificationRepositories
	pubsub *repositories.NotificationPubSub
}

func NewNotificationService(repo *repositories.NotificationRepositories, pubsub *repositories.NotificationPubSub) *NotificationService {
	return &NotificationService{repo: repo, pubsub: pubsub}
}

func (s *NotificationService) SendNotification(userID, userUsername, actorID, actorUsername, notifType, content string) error {
	username, err := s.repo.GetUsernameByID(userID)
	if err != nil {
		log.Printf("Warning could not fetch receiver username: %v", err)
	}
	notif := models.Notification{
		ID:            uuid.New().String(),
		CreatedAt:     time.Now(),
		UserID:        userID,
		UserUsername:  username,
		ActorID:       actorID,
		ActorUsername: actorUsername,
		Type:          notifType,
		Content:       content,
		Read:          false,
	}

	if err := s.repo.Create(&notif); err != nil {
		log.Printf("Error saving notification: %v", err)
		return err
	}

	log.Printf("[NotifService] Sending type=%q from actor=%q(%s) to user=%q(%s)",
		notifType, actorUsername, actorID, username, userID)
	if err := s.pubsub.PublishToUser(context.Background(), userID, &notif); err != nil {
		log.Printf("[NotifService] Error publishing notification to %s: %v", userID, err)
	}
	return nil
}

func (s *NotificationService) GetUnread(userID string) ([]models.Notification, error) {
	return s.repo.FindUnreadByUserID(userID)
}

func (s *NotificationService) MarkAllRead(userID string) error {
	return s.repo.MarkAllReadByUserID(userID)
}

package services

import (
	"context"
	"log"
	"time"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

// NotificationService creates and delivers notifications to users.
type NotificationService struct {
	repo   *repositories.NotificationRepositories
	pubsub *repositories.NotificationPubSub
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(
	repo *repositories.NotificationRepositories,
	pubsub *repositories.NotificationPubSub,
) *NotificationService {
	return &NotificationService{repo: repo, pubsub: pubsub}
}

// SendNotification saves a notification and publishes it to the user in real time.
func (s *NotificationService) SendNotification(
	userID, userUsername, actorID, actorUsername, notifType, content, postID string,
) error {
	if userUsername == "" {
		u, err := s.repo.GetUsernameByID(userID)
		if err != nil {
			log.Printf("Warning could not fetch receiver username: %v", err)
		}
		userUsername = u
	}

	notif := models.Notification{
		ID:            utils.NewID(),
		CreatedAt:     time.Now(),
		UserID:        userID,
		UserUsername:  userUsername,
		ActorID:       actorID,
		ActorUsername: actorUsername,
		Type:          notifType,
		Content:       content,
		PostID:        postID,
		Read:          false,
	}

	if err := s.repo.Create(&notif); err != nil {
		log.Printf("Error saving notification: %v", err)
		return err
	}

	log.Printf("[NotifService] Sending type=%q from actor=%q(%s) to user=%q(%s)",
		notifType, actorUsername, actorID, userUsername, userID)
	if err := s.pubsub.PublishToUser(context.Background(), userID, &notif); err != nil {
		log.Printf("[NotifService] Error publishing notification to %s: %v", userID, err)
	}
	return nil
}

func (s *NotificationService) GetAll(userID string) ([]models.Notification, error) {
	return s.repo.FindByUserID(userID, 50)
}

// GetUnread returns the unread notifications of the user.
func (s *NotificationService) GetUnread(userID string) ([]models.Notification, error) {
	return s.repo.FindUnreadByUserID(userID)
}

// MarkAllRead marks all notifications of the user as read.
func (s *NotificationService) MarkAllRead(userID string) error {
	return s.repo.MarkAllReadByUserID(userID)
}

// MarkRead marks one notification of the user as read.
func (s *NotificationService) MarkRead(userID, notifID string) error {
	return s.repo.MarkReadByID(userID, notifID)
}

func (s *NotificationService) Delete(userID, notifID string) error {
	if err := s.repo.DeleteByID(userID, notifID); err != nil {
		return err
	}
	if err := s.pubsub.PublishRemovalToUser(context.Background(), userID, notifID); err != nil {
		log.Printf("[NotifService] Error publishing notification removal to %s: %v", userID, err)
	}
	return nil
}

func (s *NotificationService) DeleteAll(userID string) error {
	if err := s.repo.DeleteAllByUserID(userID); err != nil {
		return err
	}
	if err := s.pubsub.PublishClearToUser(context.Background(), userID); err != nil {
		log.Printf("[NotifService] Error publishing notifications clear to %s: %v", userID, err)
	}
	return nil
}

func (s *NotificationService) RemoveFriendRequestNotification(userID, actorID string) {
	ids, err := s.repo.DeleteByActorAndType(userID, actorID, "friend_request")
	if err != nil {
		log.Printf("[NotifService] Error removing friend request notification: %v", err)
		return
	}
	for _, id := range ids {
		if err := s.pubsub.PublishRemovalToUser(context.Background(), userID, id); err != nil {
			log.Printf("[NotifService] Error publishing notification removal to %s: %v", userID, err)
		}
	}
}

// SendCommentNotification builds and sends a notification when someone comments on a post.
func (s *NotificationService) SendCommentNotification(
	receiverID, actorID, actorUsername, postID, commentPreview string,
) error {
	receiverUsername, err := s.repo.GetUsernameByID(receiverID)
	if err != nil {
		log.Printf("Warning could not get the username of receiverID: %v", err)
	}
	if len(commentPreview) > 50 {
		commentPreview = commentPreview[:50] + "..."
	}
	content := actorUsername + " comment on your post \"" + commentPreview + "\""
	return s.SendNotification(receiverID, receiverUsername, actorID, actorUsername, "comment", content, postID)
}

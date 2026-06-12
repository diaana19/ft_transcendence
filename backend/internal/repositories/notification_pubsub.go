package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/models"
)

// NotificationPubSub sends notifications to users using Redis pub/sub.
type NotificationPubSub struct {
	rdb *redis.Client
}

// NewNotificationPubSub creates a new NotificationPubSub using the given Redis client.
func NewNotificationPubSub(rdb *redis.Client) *NotificationPubSub {
	return &NotificationPubSub{rdb: rdb}
}

// PublishToUser publishes the notification to the user channel.
func (p *NotificationPubSub) PublishToUser(ctx context.Context, userID string, notif *models.Notification) error {
	payload, err := json.Marshal(map[string]any{
		"type":         "notification",
		"notification": notif,
	})
	if err != nil {
		return fmt.Errorf("marshal notification: %w", err)
	}
	channel := "notifications:" + userID
	log.Printf("[PubSub] Publishing to channel=%q type=%q actor=%q -> recipient=%q content=%q",
		channel, notif.Type, notif.ActorUsername, notif.UserUsername, notif.Content)
	if err := p.rdb.Publish(ctx, channel, string(payload)).Err(); err != nil {
		return fmt.Errorf("publish notification: %w", err)
	}
	return nil
}

func (p *NotificationPubSub) PublishRemovalToUser(ctx context.Context, userID, notifID string) error {
	payload, err := json.Marshal(map[string]any{
		"type":            "notification_removed",
		"notification_id": notifID,
	})
	if err != nil {
		return fmt.Errorf("marshal notification removal: %w", err)
	}
	channel := "notifications:" + userID
	if err := p.rdb.Publish(ctx, channel, string(payload)).Err(); err != nil {
		return fmt.Errorf("publish notification removal: %w", err)
	}
	return nil
}

func (p *NotificationPubSub) PublishClearToUser(ctx context.Context, userID string) error {
	payload, err := json.Marshal(map[string]any{
		"type": "notifications_cleared",
	})
	if err != nil {
		return fmt.Errorf("marshal notifications clear: %w", err)
	}
	channel := "notifications:" + userID
	if err := p.rdb.Publish(ctx, channel, string(payload)).Err(); err != nil {
		return fmt.Errorf("publish notifications clear: %w", err)
	}
	return nil
}

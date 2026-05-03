package repositories

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Transcendence/models"
	"github.com/redis/go-redis/v9"
)

type NotificationPubSub struct {
	rdb *redis.Client
}

func NewNotiticationPubSub(rdb *redis.Client) *NotificationPubSub {
	return &NotificationPubSub{rdb: rdb}
}

func (p *NotificationPubSub) PublishToUser(ctx context.Context, userID string, notif *models.Notification) error {
	payload, err := json.Marshal(map[string]interface{}{
		"type":         "notification",
		"notification": notif,
	})
	if err != nil {
		return err
	}
	channel := "notifications:" + userID
	log.Printf("[PubSub] Publishing to channel=%q type=%q actor=%q -> recipient=%q content=%q",
		channel, notif.Type, notif.ActorUsername, notif.UserUsername, notif.Content)
	return p.rdb.Publish(ctx, channel, string(payload)).Err()
}

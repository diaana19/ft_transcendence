package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Publish(rdb * redis.Client, channel, message string) error {
	ctx := context.Background()
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		log.Printf("Error: failed to publish to [%s]: %v\n", channel, err)
		return err
	}
	log.Printf("Publish to [%s]: %s\n", channel, message)
	return nil
}

func Subscribe(rdb *redis.Client, channel string, handler func(message string)) {
	ctx := context.Background()
	sub := rdb.Subscribe(ctx, channel)

	go func ()  {
		defer sub.Close()
		log.Printf("Subscribe to channel: [%s]\n", channel)
		for msg := range sub.Channel() {
			log.Printf("Received on [%s]: %s\n", msg.Channel, msg.Payload)
			if handler != nil {
				handler(msg.Payload)
			}
		}
	}()
}
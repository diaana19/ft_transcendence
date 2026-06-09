package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Publish sends a message to the Redis channel. It returns an error when the publish fails.
func Publish(rdb *redis.Client, channel, message string) error {
	ctx := context.Background()
	err := rdb.Publish(ctx, channel, message).Err()
	if err != nil {
		log.Printf("Error: failed to publish to [%s]: %v\n", channel, err)
		return fmt.Errorf("publish to %s: %w", channel, err)
	}
	log.Printf("Publish to [%s]: %s\n", channel, message)
	return nil
}

// RoundTrip publishes a message and waits to receive it back on the same channel. It is
// used as a health check, and returns an error when the payload does not match.
func RoundTrip(ctx context.Context, rdb *redis.Client, channel, message string) error {
	sub := rdb.Subscribe(ctx, channel)
	defer func() { _ = sub.Close() }()

	if _, err := sub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe handshake: %w", err)
	}

	if err := rdb.Publish(ctx, channel, message).Err(); err != nil {
		return fmt.Errorf("publish: %w", err)
	}

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return fmt.Errorf("receive message: %w", err)
	}
	if msg.Payload != message {
		return fmt.Errorf("payload mismatch: got %q want %q", msg.Payload, message)
	}
	return nil
}

// Subscribe listens to a Redis channel in a goroutine and calls handler for each message.
// The subscription stops when the context is canceled.
func Subscribe(ctx context.Context, rdb *redis.Client, channel string, handler func(message string)) {
	sub := rdb.Subscribe(ctx, channel)

	go func() {
		defer func() { _ = sub.Close() }()
		log.Printf("Subscribe to channel: [%s]\n", channel)
		for msg := range sub.Channel() {
			log.Printf("Received on [%s]: %s\n", msg.Channel, msg.Payload)
			if handler != nil {
				handler(msg.Payload)
			}
		}
	}()
}

package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Redis at %s: %w", addr, err)
	}

	log.Println("Connected to Redis successfully")
	return rdb, nil
}

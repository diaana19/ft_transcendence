package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func (r *Redis) Connect() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", r.Host, r.Port)

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

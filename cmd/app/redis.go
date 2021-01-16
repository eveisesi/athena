package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func makeRedis(cfg config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		MaxRetries:         5,
		IdleTimeout:        time.Second * 10,
		IdleCheckFrequency: time.Second * 5,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.WithError(err).Fatal("failed to ping redis server")
	}

	return redisClient
}

package redis

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

func NewClient() *redis.Client {
	opt, err := redis.ParseURL(config.Config.Redis.DSN)
	if err != nil {
		logger.Fatal("Failed to parse Redis DSN: " + err.Error())
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis: " + err.Error())
	}

	return client
}

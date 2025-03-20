package server

import "github.com/DaniilKalts/market-rest-api/pkg/redis"

func initRedis() *redis.TokenStore {
	redisClient := redis.NewClient()

	tokenStore := redis.NewTokenStore(redisClient)

	return tokenStore
}

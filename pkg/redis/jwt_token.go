package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type TokenStore struct {
	redisClient *redis.Client
}

func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{redisClient: client}
}

func (ts *TokenStore) SaveJWToken(userID int, token string) error {
	claims, err := jwt.ParseJWT(token)
	if err != nil {
		return err
	}

	expVal, ok := claims["exp"]
	if !ok {
		return errors.New("expiration time not found in token")
	}

	expFloat, ok := expVal.(float64)
	if !ok {
		return errors.New("expiration time is not a valid number")
	}

	expTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expTime)
	if duration <= 0 {
		return errors.New("token has already expired")
	}

	key := fmt.Sprintf("user:%d:jwt", userID)
	return ts.redisClient.Set(context.Background(), key, token, duration).Err()
}

func (ts *TokenStore) ValidateJWToken(userID int, token string) (bool, error) {
	key := fmt.Sprintf("user:%d:jwt", userID)
	storedToken, err := ts.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if storedToken == token {
		return true, nil
	}

	return false, nil
}

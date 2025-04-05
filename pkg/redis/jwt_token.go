package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type TokenStore interface {
	SaveJWToken(userID int, token string) error
	SaveJWTokens(userID int, accessToken, refreshToken string) error
	DeleteJWToken(userID int, token string) error
	DeleteJWTokens(userID int, accessToken, refreshToken string) error
	ValidateJWToken(userID int, token string) (bool, error)
}

type tokenStore struct {
	redisClient *redis.Client
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &tokenStore{redisClient: client}
}

func (ts *tokenStore) SaveJWToken(userID int, token string) error {
	claims, err := jwt.ParseJWT(token)
	if err != nil {
		return err
	}
	duration := time.Until(claims.ExpiresAt.Time)
	if duration <= 0 {
		return errors.New("token has already expired")
	}
	key := fmt.Sprintf("user:%d:jwt:%s", userID, claims.ID)
	return ts.redisClient.Set(context.Background(), key, token, duration).Err()
}

func (ts *tokenStore) SaveJWTokens(
	userID int, accessToken, refreshToken string,
) error {
	if err := ts.SaveJWToken(userID, accessToken); err != nil {
		return err
	}
	if err := ts.SaveJWToken(userID, refreshToken); err != nil {
		return err
	}
	return nil
}

func (ts *tokenStore) DeleteJWToken(userID int, token string) error {
	claims, err := jwt.ParseJWT(token)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("user:%d:jwt:%s", userID, claims.ID)
	return ts.redisClient.Del(context.Background(), key).Err()
}

func (ts *tokenStore) DeleteJWTokens(
	userID int, accessToken, refreshToken string,
) error {
	if err := ts.DeleteJWToken(userID, accessToken); err != nil {
		return err
	}
	if err := ts.DeleteJWToken(userID, refreshToken); err != nil {
		return err
	}
	return nil
}

func (ts *tokenStore) ValidateJWToken(userID int, token string) (bool, error) {
	claims, err := jwt.ParseJWT(token)
	if err != nil {
		return false, err
	}
	key := fmt.Sprintf("user:%d:jwt:%s", userID, claims.ID)
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

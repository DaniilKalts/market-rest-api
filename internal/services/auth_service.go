package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(email string, password string) (*models.User, error)
	Logout(userID int) error
	SaveUserToken(userID int, token string) error
	ValidateUserToken(userID int, token string) (bool, error)
}

type authService struct {
	repo        repositories.UserRepository
	redisClient *redis.Client
}

func NewAuthService(repo repositories.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{repo: repo, redisClient: redisClient}
}

func (r *authService) RegisterUser(user *models.User) error {
	existingUser, err := r.repo.GetByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if err := r.repo.Create(user); err != nil {
		return err
	}

	return nil
}

func (r *authService) AuthenticateUser(email string, password string) (*models.User, error) {
	user, err := r.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if _, err := jwt.CheckPassword(password, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authService) Logout(userID int) error {
	return nil
}

/* LOGIC TO HANDLE JWT TOKENS IN REDIS */

func (r *authService) SaveUserToken(userID int, token string) error {
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
	return r.redisClient.Set(context.Background(), key, token, duration).Err()
}

func (r *authService) ValidateUserToken(userID int, token string) (bool, error) {
	key := fmt.Sprintf("user:%d:jwt", userID)
	storedToken, err := r.redisClient.Get(context.Background(), key).Result()
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

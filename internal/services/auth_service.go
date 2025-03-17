package services

import (
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/auth"
)

type AuthService interface {
	RegisterUser(user *models.User) error
	AuthenticateUser(email string, password string) (*models.User, error)
	Logout(userID int) error
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
		return errors.New("User already exists")
	}

	hashedPassword, err := auth.HashPassword(user.Password)
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
	if _, err := auth.CheckPassword(password, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authService) Logout(userID int) error {
	return nil
}

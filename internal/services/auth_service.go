package services

import (
	"errors"

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
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (r *authService) RegisterUser(user *models.User) error {
	existingUser, err := r.repo.GetByEmail(user.Email)
	if err != nil {
		return errors.New("Error checking if user exists")
	}
	if existingUser != nil {
		return errors.New("User already exists")
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("Error hashing the password")
	}
	user.Password = hashedPassword

	if err := r.repo.Create(user); err != nil {
		return errors.New("Error creating a user")
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

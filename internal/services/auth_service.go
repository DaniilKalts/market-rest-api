package services

import (
	"errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
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

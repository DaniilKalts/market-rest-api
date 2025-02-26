package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/DaniilKalts/market-rest-api/models"
	"github.com/DaniilKalts/market-rest-api/repositories"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("Failed to hash password")
	}
	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(user)
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) UpdateUser(user *models.User) error {
	existingUser, err := s.repo.GetByID(user.ID)

	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("User not found")
	}

	if user.PasswordHash != "" && user.PasswordHash != existingUser.PasswordHash {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("Failed to hash password")
		}
		user.PasswordHash = string(hashedPassword)
	} else {
		user.PasswordHash = existingUser.PasswordHash
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id int) error {
	return s.repo.Delete(id)
}

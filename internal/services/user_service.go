package services

import (
	"errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type UserService interface {
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUserByID(id int, updateUserDTO *models.UpdateUser) error
	DeleteUserByID(id int) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) UpdateUserByID(userID int, updateUserDTO *models.UpdateUser) error {
	existingUser, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	if updateUserDTO.FirstName != nil {
		existingUser.FirstName = *updateUserDTO.FirstName
	}
	if updateUserDTO.LastName != nil {
		existingUser.LastName = *updateUserDTO.LastName
	}
	if updateUserDTO.Email != nil {
		existingUser.Email = *updateUserDTO.Email
	}
	if updateUserDTO.PhoneNumber != nil {
		existingUser.PhoneNumber = *updateUserDTO.PhoneNumber
	}
	if updateUserDTO.Password != nil || updateUserDTO.ConfirmPassword != nil {
		if updateUserDTO.Password == nil || updateUserDTO.ConfirmPassword == nil || *updateUserDTO.Password != *updateUserDTO.ConfirmPassword {
			return errors.New("passwords do not match")
		}
		if *updateUserDTO.Password != "" {
			hashedPassword, err := jwt.HashPassword(*updateUserDTO.Password)
			if err != nil {
				return err
			}
			existingUser.Password = hashedPassword
		}
	}

	return s.repo.Update(existingUser)
}

func (s *userService) DeleteUserByID(id int) error {
	return s.repo.Delete(id)
}

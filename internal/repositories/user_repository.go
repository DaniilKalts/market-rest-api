package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User

	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id int) error {
	res := r.db.Delete(&models.User{}, id)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("User not found")
	}

	return nil
}

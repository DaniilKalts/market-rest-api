package models

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

var phoneRegex = regexp.MustCompile(`^\+7[0-9]{10}$`)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func ValidatePhoneNumber(phoneNumber string) error {
	if !phoneRegex.MatchString(phoneNumber) {
		return errors.New("invalid phone number format for Kazakhstan")
	}

	return nil
}

type User struct {
	ID          int       `json:"id" gorm:"primaryKey" example:"1"`
	FirstName   string    `json:"first_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"John"`
	LastName    string    `json:"last_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"Doe"`
	Email       string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,email" example:"john.doe@example.com"`
	Password    string    `json:"password" gorm:"type:varchar(255);not null" binding:"required,min=8" example:"$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(12);not null" binding:"required" example:"+77051234567"`
	Role        Role      `json:"role" gorm:"type:varchar(10);not null;default:'user'" binding:"required,oneof=admin user" example:"user"`
	Cart        *Cart     `json:"cart,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := jwt.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	return
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12341234"`
}

type RegisterUser struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=30" example:"John"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30" example:"Doe"`
	Email           string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password        string `json:"password" binding:"required,min=8" example:"12341234"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8" example:"12341234"`
	PhoneNumber     string `json:"phone_number" binding:"required" example:"+77051234567"`
}

func (r *RegisterUser) Validate() error {
	if r.Password != r.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return ValidatePhoneNumber(r.PhoneNumber)
}

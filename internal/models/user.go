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
	FirstName   string    `json:"first_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"Martin"`
	LastName    string    `json:"last_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"Kalts"`
	Email       string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,email" example:"martin@gmail.com"`
	Password    string    `json:"password" gorm:"type:varchar(255);not null" binding:"required,min=8" example:"$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(12);not null" binding:"required" example:"+77007473472"`
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

type RegisterUser struct {
	FirstName       string `json:"first_name" binding:"required,min=2,max=30" example:"Martin"`
	LastName        string `json:"last_name" binding:"required,min=2,max=30" example:"Kalts"`
	Email           string `json:"email" binding:"required,email" example:"martin@gmail.com"`
	Password        string `json:"password" binding:"required,min=8" example:"12341234"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8" example:"12341234"`
	PhoneNumber     string `json:"phone_number" binding:"required" example:"+77007473472"`
}

func (r *RegisterUser) Validate() error {
	if r.Password != r.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return ValidatePhoneNumber(r.PhoneNumber)
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email" example:"martin@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12341234"`
}

type UpdateUser struct {
	FirstName       *string `json:"first_name" binding:"omitempty,min=2,max=30" example:"Martin"`
	LastName        *string `json:"last_name" binding:"omitempty,min=2,max=30" example:"Kalts"`
	Email           *string `json:"email" binding:"omitempty,email" example:"martin.programmer@gmail.com"`
	Password        *string `json:"password" binding:"omitempty,min=8" example:"12341234"`
	ConfirmPassword *string `json:"confirm_password" binding:"omitempty,min=8" example:"12341234"`
	PhoneNumber     *string `json:"phone_number" binding:"omitempty" example:"+77007473472"`
}

func (u *UpdateUser) Validate() error {
	if u.Password != nil || u.ConfirmPassword != nil {
		if u.Password == nil || u.ConfirmPassword == nil || *u.Password != *u.ConfirmPassword {
			return errors.New("passwords do not match")
		}
	}
	if u.PhoneNumber != nil {
		if err := ValidatePhoneNumber(*u.PhoneNumber); err != nil {
			return err
		}
	}
	return nil
}

package models

import "time"

type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	FirstName    string    `json:"first_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30"`
	LastName     string    `json:"last_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30"`
	Email        string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,email"`
	PasswordHash string    `json:"password_hash" gorm:"type:varchar(255);not null" binding:"required,min=8"`
	PhoneNumber  string    `json:"phone_number" gorm:"type:varchar(20)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

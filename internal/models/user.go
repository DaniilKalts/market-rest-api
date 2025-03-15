package models

import "time"

type User struct {
	ID          int       `json:"id" gorm:"primaryKey" example:"1"`
	FirstName   string    `json:"first_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"John"`
	LastName    string    `json:"last_name" gorm:"type:varchar(30);not null" binding:"required,min=2,max=30" example:"Doe"`
	Email       string    `json:"email" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,email" example:"john.doe@example.com"`
	Password    string    `json:"password" gorm:"type:varchar(255);not null" binding:"required,min=8" example:"$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20)" example:"+12345678901"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

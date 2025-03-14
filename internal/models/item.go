package models

import "time"

type Item struct {
	ID          int       `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,min=5,max=40" example:"T-shirt"`
	Description string    `json:"description" gorm:"type:varchar(255)" example:"A fashionable T-shirt with IITU logo."`
	Price       int       `json:"price" gorm:"not null" binding:"required,gte=100,lte=1000" example:"100"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

type RequestCreateItem struct {
	Name        string `json:"name" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,min=5,max=40" example:"T-shirt"`
	Description string `json:"description" gorm:"type:varchar(255)" example:"A fashionable T-shirt with IITU logo."`
	Price       int    `json:"price" gorm:"not null" binding:"required,gte=100,lte=1000" example:"100"`
}

type RequestUpdateItem struct {
	ID          int    `json:"id" gorm:"primaryKey" example:"1"`
	Name        string `json:"name" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,min=5,max=40" example:"T-shirt"`
	Description string `json:"description" gorm:"type:varchar(255)" example:"A fashionable T-shirt with IITU logo."`
	Price       int    `json:"price" gorm:"not null" binding:"required,gte=100,lte=1000" example:"100"`
}

type BadRequestError struct {
	Error string `json:"error" example:"Invalid request payload"`
}

type NotFoundError struct {
	Error string `json:"error" example:"Item not found"`
}

type InternalServerError struct {
	Error string `json:"error" example:"Failed to create item"`
}

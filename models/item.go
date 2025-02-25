package models

import "time"

type Item struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,min=5,max=40"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	Price       int       `json:"price" gorm:"not null" binding:"required,gte=100,lte=1000"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

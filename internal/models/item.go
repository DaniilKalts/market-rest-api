package models

import "time"

type Item struct {
	ID          int       `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required,min=5,max=40" example:"T-shirt"`
	Description string    `json:"description" gorm:"type:varchar(255)" example:"A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort."`
	Price       uint      `json:"price" gorm:"not null" binding:"required,gte=10,lte=100" example:"30"`
	Stock       uint      `json:"stock" gorm:"not null" binding:"required" example:"20"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

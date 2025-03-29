package models

import "time"

type Cart struct {
	ID        int        `json:"id" gorm:"primaryKey" example:"1"`
	UserID    int        `json:"user_id" gorm:"not null" example:"1"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

type CartItem struct {
	CartID    int       `json:"cart_id" gorm:"primaryKey;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" example:"1"`
	ItemID    int       `json:"item_id" gorm:"primaryKey;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" example:"1"`
	Quantity  uint      `json:"quantity" gorm:"not null" example:"2"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2025-02-25T12:37:32Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2025-02-25T12:37:32Z"`
}

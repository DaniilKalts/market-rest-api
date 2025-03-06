package server

import (
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/models"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Item{}, &models.User{})
}

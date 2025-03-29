package server

import (
	"errors"

	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Item{}, &models.User{}, &models.Cart{}, &models.CartItem{})

	var admin models.User

	err := db.Where("role = ?", models.RoleAdmin).First(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			admin := models.User{
				FirstName:   config.Config.Admin.FirstName,
				LastName:    config.Config.Admin.LastName,
				Email:       config.Config.Admin.Email,
				Password:    config.Config.Admin.Password,
				PhoneNumber: config.Config.Admin.PhoneNumber,
				Role:        models.RoleAdmin,
			}
			if err := db.Create(&admin).Error; err != nil {
				logger.Error("Failed to create admin user: " + err.Error())
			} else {
				logger.Info("Admin user created successfully")
			}
		}
	} else {
		logger.Info("Admin user already exists")
	}
}

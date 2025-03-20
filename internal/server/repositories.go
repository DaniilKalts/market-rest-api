package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"gorm.io/gorm"
)

func initRepositories(db *gorm.DB) (repositories.ItemRepository, repositories.UserRepository) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	return itemRepo, userRepo
}

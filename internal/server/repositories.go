package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"gorm.io/gorm"
)

func initRepositories(db *gorm.DB) (
	repositories.ItemRepository,
	repositories.UserRepository,
	repositories.CartRepository,
) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)
	cartRepo := repositories.NewCartRepository(db)

	return itemRepo, userRepo, cartRepo
}

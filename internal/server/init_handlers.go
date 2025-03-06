package server

import (
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

func initHandlers(db *gorm.DB) (*handlers.ItemHandler, *handlers.UserHandler) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)

	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)

	return itemHandler, userHandler
}

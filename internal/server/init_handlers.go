package server

import (
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func initHandlers(db *gorm.DB, tokenStore *redis.TokenStore) (*handlers.ItemHandler, *handlers.UserHandler, *handlers.AuthHandler) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, tokenStore)

	return itemHandler, userHandler, authHandler
}

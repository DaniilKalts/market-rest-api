package server

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

func initHandlers(db *gorm.DB, redisClient *redis.Client) (*handlers.ItemHandler, *handlers.UserHandler, *handlers.AuthHandler, *services.AuthService) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, redisClient)

	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	return itemHandler, userHandler, authHandler, &authService
}

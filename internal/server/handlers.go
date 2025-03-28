package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func initHandlers(
	itemService services.ItemService,
	userService services.UserService,
	authService services.AuthService,
	tokenStore *redis.TokenStore,
) (*handlers.ItemHandler, *handlers.UserHandler, *handlers.AuthHandler, *handlers.ProfileHandler) {
	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, tokenStore)
	profileHandler := handlers.NewProfileHandler(userService, authService)

	return itemHandler, userHandler, authHandler, profileHandler
}

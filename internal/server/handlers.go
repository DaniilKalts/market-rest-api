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
	cartService services.CartService,
	tokenStore *redis.TokenStore,
) (
	*handlers.ItemHandler,
	*handlers.UserHandler,
	*handlers.AuthHandler,
	*handlers.ProfileHandler,
	*handlers.CartHandler,
) {
	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, tokenStore)
	profileHandler := handlers.NewProfileHandler(userService, authService)
	cartHandler := handlers.NewCartHandler(itemService, cartService)

	return itemHandler, userHandler, authHandler, profileHandler, cartHandler
}

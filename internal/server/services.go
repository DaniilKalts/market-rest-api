package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func initServices(
	itemRepo repositories.ItemRepository,
	userRepo repositories.UserRepository,
	cartRepo repositories.CartRepository,
	tokenStore redis.TokenStore,
) (
	services.ItemService,
	services.UserService,
	services.AuthService,
	services.CartService,
) {
	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenStore)
	cartService := services.NewCartService(cartRepo, itemService)

	return itemService, userService, authService, cartService
}

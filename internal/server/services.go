package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func initServices(
	itemRepo repositories.ItemRepository,
	userRepo repositories.UserRepository,
	tokenStore *redis.TokenStore,
) (services.ItemService, services.UserService, services.AuthService) {
	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenStore)

	return itemService, userService, authService
}

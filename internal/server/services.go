package server

import (
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

func initServices(
	itemRepo repositories.ItemRepository,
	userRepo repositories.UserRepository,
) (services.ItemService, services.UserService, services.AuthService) {
	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	return itemService, userService, authService
}

package server

import (
	"net/http"

	"github.com/DaniilKalts/market-rest-api/internal/config"
)

func SetupServer() *http.Server {
	db := initDB()
	migrate(db)

	tokenStore := initRedis()

	itemRepository, userRepository, cartRepository := initRepositories(db)
	itemService, userService, authService, cartService := initServices(
		itemRepository,
		userRepository,
		cartRepository,
		tokenStore,
	)
	itemHandler, userHandler, authHandler, profileHandler, cartHandler := initHandlers(
		itemService,
		userService,
		authService,
		cartService,
		tokenStore,
	)

	port := config.Config.Server.Port
	if port == "" {
		port = "8080"
	}

	router := setupRouter(
		itemHandler,
		userHandler,
		authHandler,
		profileHandler,
		cartHandler,
	)

	srv := &http.Server{
		Addr:    ":" + config.Config.Server.Port,
		Handler: router,
	}

	return srv
}

package server

import (
	"net/http"

	"github.com/DaniilKalts/market-rest-api/internal/config"
)

func SetupServer() *http.Server {
	db := initDB()
	migrate(db)

	tokenStore := initRedis()

	itemRepository, userRepository := initRepositories(db)
	itemService, userService, authService := initServices(itemRepository, userRepository, tokenStore)
	itemHandler, userHandler, authHandler, profileHandler := initHandlers(
		itemService,
		userService,
		authService,
		tokenStore,
	)

	port := config.Config.Server.Port
	if port == "" {
		port = "8080"
	}

	router := setupRouter(itemHandler, userHandler, authHandler, profileHandler)

	srv := &http.Server{
		Addr:    ":" + config.Config.Server.Port,
		Handler: router,
	}

	return srv
}

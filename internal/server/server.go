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
	itemService, userService, authService := initServices(itemRepository, userRepository)
	itemHandler, userHandler, authHandler := initHandlers(
		itemService,
		userService,
		authService,
		tokenStore,
	)

	port := config.Config.Server.Port
	if port == "" {
		port = "8080"
	}

	router := setupRouter(itemHandler, userHandler, authHandler)

	srv := &http.Server{
		Addr:    ":" + config.Config.Server.Port,
		Handler: router,
	}

	return srv
}

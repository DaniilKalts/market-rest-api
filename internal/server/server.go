package server

import (
	"net/http"
	"time"

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
	)

	router := setupRouter(
		itemHandler,
		userHandler,
		authHandler,
		profileHandler,
		cartHandler,
	)

	srv := &http.Server{
		Addr:              ":" + config.Config.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv
}

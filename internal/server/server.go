package server

import (
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
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

	router := setupRouter(itemHandler, userHandler, authHandler)

	return router
}

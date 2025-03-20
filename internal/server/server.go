package server

import (
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	db := initDB()
	migrate(db)

	tokenStore := initRedis()

	itemHandler, userHandler, authHandler := initHandlers(db, tokenStore)

	router := setupRouter(itemHandler, userHandler, authHandler, tokenStore)

	return router
}

package server

import (
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	db := initDB()
	migrate(db)

	redisClient := initRedis()

	itemHandler, userHandler, authHandler := initHandlers(db, redisClient)
	return setupRouter(itemHandler, userHandler, authHandler)
}

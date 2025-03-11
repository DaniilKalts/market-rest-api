package server

import (
	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	db := initDB()
	migrate(db)

	itemHandler, userHandler, authHandler := initHandlers(db)
	return setupRouter(itemHandler, userHandler, authHandler)
}

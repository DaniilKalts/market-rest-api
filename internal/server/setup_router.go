package server

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
)

func setupRouter(itemHandler *handlers.ItemHandler, userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/item/create", itemHandler.CreateItem)
	router.GET("/item", itemHandler.GetItemByID)
	router.GET("/items", itemHandler.GetAllItems)
	router.PUT("/item/update", itemHandler.UpdateItem)
	router.DELETE("/item/delete", itemHandler.DeleteItem)

	router.POST("/user/create", userHandler.CreateUser)
	router.GET("/user", userHandler.GetUserByID)
	router.GET("/users", userHandler.GetAllUsers)
	router.PUT("/user/update", userHandler.UpdateUser)
	router.DELETE("/user/delete", userHandler.DeleteUser)

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL: "/docs/openapi.yaml",
	}, swaggerFiles.Handler))

	return router
}

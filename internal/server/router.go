package server

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/middlewares"
	"github.com/DaniilKalts/market-rest-api/internal/models"
)

func setupRouter(itemHandler *handlers.ItemHandler, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	tokenStore := initRedis()

	router.Use(middlewares.LoggerMiddleware())

	itemPublicRoutes := router.Group("/items")
	{
		itemPublicRoutes.GET("/:id", itemHandler.GetItemByID)
		itemPublicRoutes.GET("", itemHandler.GetAllItems)
	}

	itemPrivateRoutes := router.Group("/items")
	itemPrivateRoutes.Use(middlewares.JWTMiddleware(tokenStore))
	{
		itemPrivateRoutes.POST("", middlewares.BindBodyMiddleware(&models.Item{}), itemHandler.CreateItem)
		itemPrivateRoutes.PUT("/:id", middlewares.BindBodyMiddleware(&models.Item{}), itemHandler.UpdateItem)
		itemPrivateRoutes.DELETE("/:id", itemHandler.DeleteItem)
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middlewares.JWTMiddleware(tokenStore))
	{
		userRoutes.POST("", middlewares.BindBodyMiddleware(&models.User{}), userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.GET("", userHandler.GetAllUsers)
		userRoutes.PUT("/:id", middlewares.BindBodyMiddleware(&models.User{}), userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", middlewares.BindBodyMiddleware(&models.RegisterUser{}), authHandler.Register)
		authRoutes.POST("/login", middlewares.BindBodyMiddleware(&models.LoginUser{}), authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL: "/docs/openapi.yaml",
	}, swaggerFiles.Handler))

	return router
}

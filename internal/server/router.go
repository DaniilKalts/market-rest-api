package server

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/DaniilKalts/market-rest-api/internal/handlers"
	"github.com/DaniilKalts/market-rest-api/internal/middlewares"
	"github.com/DaniilKalts/market-rest-api/internal/models"
)

func setupRouter(itemHandler *handlers.ItemHandler, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, profileHandler *handlers.ProfileHandler) *gin.Engine {
	router := gin.Default()

	tokenStore := initRedis()

	router.Use(middlewares.LoggerMiddleware())

	itemPublicRoutes := router.Group("/items")
	{
		itemPublicRoutes.GET("/:id", itemHandler.HandleGetItemByID)
		itemPublicRoutes.GET("", itemHandler.HandleGetAllItems)
	}

	itemPrivateRoutes := router.Group("/items")
	itemPrivateRoutes.Use(middlewares.JWTMiddleware(tokenStore), middlewares.TokenStoreMiddleware(tokenStore))
	{
		itemPrivateRoutes.POST("", middlewares.BindBodyMiddleware(&models.Item{}), itemHandler.HandleCreateItem)
		itemPrivateRoutes.PUT("/:id", middlewares.BindBodyMiddleware(&models.Item{}), itemHandler.HandleUpdateItem)
		itemPrivateRoutes.DELETE("/:id", itemHandler.HandleDeleteItem)
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middlewares.JWTMiddleware(tokenStore), middlewares.TokenStoreMiddleware(tokenStore))
	{
		userRoutes.GET("/:id", userHandler.HandleGetUserByID)
		userRoutes.GET("", userHandler.HandleGetAllUsers)
		userRoutes.PUT("/:id", middlewares.BindBodyMiddleware(&models.UpdateUser{}), userHandler.HandleUpdateUserByID)
		userRoutes.DELETE("/:id", userHandler.HandleDeleteUser)
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", middlewares.BindBodyMiddleware(&models.RegisterUser{}), authHandler.HandleRegister)
		authRoutes.POST("/login", middlewares.BindBodyMiddleware(&models.LoginUser{}), authHandler.HandleLogin)
		authRoutes.POST("/logout", authHandler.HandleLogout)
		authRoutes.POST("/refresh", authHandler.HandleRefreshToken)
	}

	profileRoutes := router.Group("/profile")
	profileRoutes.Use(middlewares.JWTMiddleware(tokenStore), middlewares.TokenStoreMiddleware(tokenStore))
	{
		profileRoutes.GET("", profileHandler.HandleGetProfile)
		profileRoutes.PUT("", middlewares.BindBodyMiddleware(&models.UpdateUser{}), profileHandler.HandleUpdateProfile)
		profileRoutes.DELETE("", profileHandler.HandleDeleteProfile)
	}

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.CustomWrapHandler(&ginSwagger.Config{
		URL: "/docs/openapi.yaml",
	}, swaggerFiles.Handler))

	return router
}

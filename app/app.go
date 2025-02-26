package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/DaniilKalts/market-rest-api/handlers"
	"github.com/DaniilKalts/market-rest-api/logger"
	"github.com/DaniilKalts/market-rest-api/models"
	"github.com/DaniilKalts/market-rest-api/repositories"
	"github.com/DaniilKalts/market-rest-api/services"
)

func initDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Item{}, &models.User{})
}

func initHandlers(db *gorm.DB) (*handlers.ItemHandler, *handlers.UserHandler) {
	itemRepo := repositories.NewItemRepository(db)
	userRepo := repositories.NewUserRepository(db)

	itemService := services.NewItemService(itemRepo)
	userService := services.NewUserService(userRepo)

	itemHandler := handlers.NewItemHandler(itemService)
	userHandler := handlers.NewUserHandler(userService)

	return itemHandler, userHandler
}

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func SetupApplication() *gin.Engine {
	db := initDB()
	migrate(db)

	itemHandler, userHandler := initHandlers(db)
	return setupRouter(itemHandler, userHandler)
}

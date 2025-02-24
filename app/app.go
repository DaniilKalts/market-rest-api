package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

	router.POST("/item/create", func(c *gin.Context) {
		itemHandler.CreateItem(c.Writer, c.Request)
	})
	router.GET("/item", func(c *gin.Context) {
		itemHandler.GetItemByID(c.Writer, c.Request)
	})
	router.GET("/items", func(c *gin.Context) {
		itemHandler.GetAllItems(c.Writer, c.Request)
	})
	router.PUT("/item/update", func(c *gin.Context) {
		itemHandler.UpdateItem(c.Writer, c.Request)
	})
	router.DELETE("/item/delete", func(c *gin.Context) {
		itemHandler.DeleteItem(c.Writer, c.Request)
	})

	router.POST("/user/create", func(c *gin.Context) {
		userHandler.CreateUser(c.Writer, c.Request)
	})
	router.GET("/user", func(c *gin.Context) {
		userHandler.GetUserByID(c.Writer, c.Request)
	})
	router.GET("/users", func(c *gin.Context) {
		userHandler.GetAllUsers(c.Writer, c.Request)
	})
	router.PUT("/user/update", func(c *gin.Context) {
		userHandler.UpdateUser(c.Writer, c.Request)
	})
	router.DELETE("/user/delete", func(c *gin.Context) {
		userHandler.DeleteUser(c.Writer, c.Request)
	})

	return router
}

func SetupApplication() *gin.Engine {
	db := initDB()
	migrate(db)

	itemHandler, userHandler := initHandlers(db)
	return setupRouter(itemHandler, userHandler)
}

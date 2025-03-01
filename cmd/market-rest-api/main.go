package main

import (
	"os"

	_ "github.com/DaniilKalts/market-rest-api/docs"

	"github.com/DaniilKalts/market-rest-api/docs"
	"github.com/DaniilKalts/market-rest-api/internal/app"
	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/internal/logger"
)

func main() {
	config.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Info("Look")
	logger.Info(os.Getenv("PORT"))
	logger.Info(os.Getenv("DATABASE_DSN"))

	docs.SwaggerInfo.Title = "Market REST API"
	docs.SwaggerInfo.Description = "A REST API for managing market items and user accounts."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/"

	router := app.SetupApplication()

	logger.Info("Server is running on: http://localhost:" + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to run the server: " + err.Error())
	}
}

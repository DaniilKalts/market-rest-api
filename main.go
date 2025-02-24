package main

import (
	"os"

	"github.com/DaniilKalts/market-rest-api/app"
	"github.com/DaniilKalts/market-rest-api/logger"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := app.SetupApplication()

	logger.Info("Server is running on: http://localhost:" + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to run the server: " + err.Error())
	}
}

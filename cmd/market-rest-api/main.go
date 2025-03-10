package main

import (
	"os"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/internal/logger"
	"github.com/DaniilKalts/market-rest-api/internal/server"
)

func main() {
	config.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := server.SetupServer()

	logger.Info("Server is running on: http://localhost:" + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Failed to run the server: " + err.Error())
	}
}

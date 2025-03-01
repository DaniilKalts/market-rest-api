package config

import (
	"github.com/joho/godotenv"

	"github.com/DaniilKalts/market-rest-api/internal/logger"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		logger.Error("init: No .env file found " + err.Error())
	}
}

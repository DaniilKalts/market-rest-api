package main

import (
	"github.com/joho/godotenv"

	"github.com/DaniilKalts/market-rest-api/logger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Error("init: No .env file found " + err.Error())
	}
}

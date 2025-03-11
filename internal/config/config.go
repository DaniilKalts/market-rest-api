package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/DaniilKalts/market-rest-api/internal/logger"
)

type ServerConfig struct {
	Port   string
	Secret string
}

type DatabaseConfig struct {
	DSN string
}

type AppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
}

var Config AppConfig

func Load() {
	if err := godotenv.Load(); err != nil {
		logger.Error("init: No .env file found " + err.Error())
	}

	Config = AppConfig{
		Server: ServerConfig{
			Port:   os.Getenv("PORT"),
			Secret: os.Getenv("SECRET"),
		},
		Database: DatabaseConfig{
			DSN: os.Getenv("DATABASE_DSN"),
		},
	}
}

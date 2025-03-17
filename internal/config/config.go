package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

type ServerConfig struct {
	Port    string
	Secret  string
	BaseURL string
	Domain  string
}

type DatabaseConfig struct {
	DSN           string
	RedisDSN      string
	RedisPassword string
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
			Port:    os.Getenv("PORT"),
			Secret:  os.Getenv("SECRET"),
			BaseURL: os.Getenv("BASE_URL"),
			Domain:  os.Getenv("DOMAIN"),
		},
		Database: DatabaseConfig{
			DSN:           os.Getenv("DATABASE_DSN"),
			RedisDSN:      os.Getenv("REDIS_DSN"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		},
	}
}

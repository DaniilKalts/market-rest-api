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

type PostgresConfig struct {
	DSN string
}

type RedisConfig struct {
	DSN           string
	RedisPassword string
}

type AppConfig struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
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
		Postgres: PostgresConfig{
			DSN: os.Getenv("POSTGRES_DSN"),
		},
		Redis: RedisConfig{
			DSN:           os.Getenv("REDIS_DSN"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
		},
	}
}

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

type AdminConfig struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
}

type AppConfig struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Admin    AdminConfig
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
		Admin: AdminConfig{
			FirstName:   os.Getenv("ADMIN_FIRST_NAME"),
			LastName:    os.Getenv("ADMIN_LAST_NAME"),
			Email:       os.Getenv("ADMIN_EMAIL"),
			Password:    os.Getenv("ADMIN_PASSWORD"),
			PhoneNumber: os.Getenv("ADMIN_PHONE_NUMBER"),
		},
	}
}

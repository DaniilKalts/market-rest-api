package config

import (
	"os"
	"strings"

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

	endFields := map[string]string{
		"PORT":               Config.Server.Port,
		"SECRET":             Config.Server.Secret,
		"BASE_URL":           Config.Server.BaseURL,
		"DOMAIN":             Config.Server.Domain,
		"POSTGRES_DSN":       Config.Postgres.DSN,
		"REDIS_DSN":          Config.Redis.DSN,
		"REDIS_PASSWORD":     Config.Redis.RedisPassword,
		"ADMIN_FIRST_NAME":   Config.Admin.FirstName,
		"ADMIN_LAST_NAME":    Config.Admin.LastName,
		"ADMIN_EMAIL":        Config.Admin.Email,
		"ADMIN_PASSWORD":     Config.Admin.Password,
		"ADMIN_PHONE_NUMBER": Config.Admin.PhoneNumber,
	}

	missing := []string{}
	for key, value := range endFields {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		logger.Error("Missing required environment variables: " + strings.Join(missing, ", "))
		os.Exit(1)
	}
}

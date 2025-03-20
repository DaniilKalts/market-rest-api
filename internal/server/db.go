package server

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

func initDB() *gorm.DB {
	dsn := config.Config.Postgres.DSN

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database: " + err.Error())
	}

	return db
}

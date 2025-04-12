package integration

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"
)

func TestPostgresConnection(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("failed to load .env file:", err)
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		t.Skip("POSTGRES_DSN not set, skipping integration test")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	sqlDB.SetConnMaxLifetime(5 * time.Second)
	err = sqlDB.Ping()
	require.NoError(t, err)
}

func TestPostgresInvalidDSN(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("failed to load .env file:", err)
	}

	originalDSN := os.Getenv("POSTGRES_DSN")
	os.Setenv("POSTGRES_DSN", "invalid-postgres-dsn")
	defer os.Setenv("POSTGRES_DSN", originalDSN)

	db, err := gorm.Open(
		postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{},
	)
	if err != nil {
		wrappedErr := fmt.Errorf("%w: %v", errs.ErrInvalidPostgresDSN, err)
		require.Error(t, wrappedErr)
		require.True(t, errors.Is(wrappedErr, errs.ErrInvalidPostgresDSN))
		return
	}

	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	require.Error(t, err)
	wrappedErr := fmt.Errorf("%w: %v", errs.ErrPostgresConnectionFailed, err)
	require.True(t, errors.Is(wrappedErr, errs.ErrPostgresConnectionFailed))
}

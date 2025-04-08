package integration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"
)

func TestRedisConnection(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal("failed to load .env file:", err)
	}

	dsn := os.Getenv("REDIS_DSN")
	if dsn == "" {
		t.Skip("REDIS_DSN not set, skipping integration test")
	}

	opt, err := redis.ParseURL(dsn)
	require.NoError(t, err)

	client := redis.NewClient(opt)
	defer func() {
		if err := client.Close(); err != nil {
			t.Errorf("failed to close client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx).Err()
	require.NoError(t, err)
}

func TestRedisInvalidDSN(t *testing.T) {
	invalidDSN := "invalid-redis-dsn"

	_, err := redis.ParseURL(invalidDSN)
	require.Error(t, err)
	wrappedErr := fmt.Errorf("%w: %v", errs.ErrInvalidRedisDSN, err)
	require.True(t, errors.Is(wrappedErr, errs.ErrInvalidRedisDSN))
}

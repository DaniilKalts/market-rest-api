package middlewares

import (
	"errors"
	"net/http"

	"github.com/DaniilKalts/market-rest-api/pkg/redis"
	"github.com/gin-gonic/gin"
)

var (
	ErrTokenNotFound            = errors.New("token not found")
	ErrTokenTypeAssertionFailed = errors.New("token type assertion failed")
	ErrUnauthorizedToken        = errors.New("unauthorized or invalid token")
)

func TokenStoreMiddleware(tokenStore *redis.TokenStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDVal, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrUserIDNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		userID, ok := userIDVal.(int)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrUserIDTypeAssertionFailed.Error()},
			)
			ctx.Abort()
			return
		}

		tokenStringVal, exists := ctx.Get("tokenString")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrTokenNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		tokenString, ok := tokenStringVal.(string)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrTokenTypeAssertionFailed.Error()},
			)
			ctx.Abort()
			return
		}

		valid, err := tokenStore.ValidateJWToken(userID, tokenString)
		if err != nil || !valid {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrUnauthorizedToken.Error()},
			)
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

package middlewares

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
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
		claimsVal, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrClaimsNotFound.Error()},
			)
			ctx.Abort()
			return
		}
		claims, ok := claimsVal.(*jwt.Claims)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrInvalidClaims.Error()},
			)
			ctx.Abort()
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrUnauthorizedToken.Error()},
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

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

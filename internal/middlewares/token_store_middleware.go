package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func TokenStoreMiddleware(tokenStore *redis.TokenStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsVal, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrClaimsNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		claims, ok := claimsVal.(*jwt.Claims)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrInvalidClaims.Error()},
			)
			ctx.Abort()
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrUnauthorizedToken.Error()},
			)
			ctx.Abort()
			return
		}

		tokenStringVal, exists := ctx.Get("tokenString")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrTokenNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		tokenString, ok := tokenStringVal.(string)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrTokenTypeFailed.Error()},
			)
			ctx.Abort()
			return
		}

		valid, err := tokenStore.ValidateJWToken(userID, tokenString)
		if err != nil || !valid {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrUnauthorizedToken.Error()},
			)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

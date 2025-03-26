package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

func JWTMiddleware(tokenStore *redis.TokenStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or invalid token"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("tokenString", tokenString)
		ctx.Next()
	}
}

func TokenStoreMiddleware(tokenStore *redis.TokenStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
			ctx.Abort()
			return
		}

		tokenString, exists := ctx.Get("tokenString")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
			ctx.Abort()
			return
		}

		userID, convErr := strconv.Atoi(claims.(*jwt.Claims).Subject)
		if convErr != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			ctx.Abort()
			return
		}

		valid, err := tokenStore.ValidateJWToken(userID, tokenString.(string))
		if err != nil || !valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized or invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

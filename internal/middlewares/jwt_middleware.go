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
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("tokenString", tokenString)
		c.Next()
	}
}

func TokenStoreMiddleware(tokenStore *redis.TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "claims not found"})
			c.Abort()
			return
		}

		tokenString, exists := c.Get("tokenString")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
			c.Abort()
			return
		}

		userID, convErr := strconv.Atoi(claims.(*jwt.Claims).Subject)
		if convErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			c.Abort()
			return
		}

		valid, err := tokenStore.ValidateJWToken(userID, tokenString.(string))
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized or invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

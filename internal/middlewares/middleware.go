package middlewares

import (
	"net/http"
	"strings"

	"github.com/DaniilKalts/market-rest-api/internal/helpers"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := helpers.ExtractClaimsFromToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func FromGinContext(c *gin.Context) (*helpers.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	customClaims, ok := claims.(*helpers.Claims)

	return customClaims, ok
}

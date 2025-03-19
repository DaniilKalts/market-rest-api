package middlewares

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

func JWTMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Invalid signing method")
			}
			return []byte(config.Config.Server.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		userIDVal, ok := claims["sub"]
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in token"})
			c.Abort()
			return
		}

		userIDStr, ok := userIDVal.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			c.Abort()
			return
		}

		userID, convErr := strconv.Atoi(userIDStr)
		if convErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid used id in token"})
			c.Abort()
			return
		}

		valid, err := authService.ValidateUserToken(userID, tokenString)
		if err != nil || !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized or invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

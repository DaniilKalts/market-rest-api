package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

var (
	ErrAuthHeaderMissing         = errors.New("authorization header missing or invalid token")
	ErrUserIDNotFound            = errors.New("user id not found")
	ErrUserIDTypeAssertionFailed = errors.New("user id type assertion failed")
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrAuthHeaderMissing.Error()},
			)
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("tokenString", tokenString)
		ctx.Next()
	}
}

package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

var (
	ErrClaimsNotFound  = errors.New("claims not found")
	ErrInvalidClaims   = errors.New("invalid claims")
	ErrAdminAccessOnly = errors.New("admin access only")
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsValue, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrClaimsNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		claims, ok := claimsValue.(*jwt.Claims)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": ErrInvalidClaims.Error()},
			)
			ctx.Abort()
			return
		}

		if claims.Role != "admin" {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{"error": ErrAdminAccessOnly.Error()},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

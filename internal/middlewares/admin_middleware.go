package middlewares

import (
	"net/http"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsValue, exists := ctx.Get("claims")
		if !exists {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrClaimsNotFound.Error()},
			)
			ctx.Abort()
			return
		}

		claims, ok := claimsValue.(*jwt.Claims)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{"error": errs.ErrInvalidClaims.Error()},
			)
			ctx.Abort()
			return
		}

		if claims.Role != "admin" {
			ctx.JSON(
				http.StatusForbidden, gin.H{"error": errs.ErrAdminOnly.Error()},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

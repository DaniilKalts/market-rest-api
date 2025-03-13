package middlewares

import (
	"fmt"
	"time"

	"github.com/DaniilKalts/market-rest-api/internal/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		message := fmt.Sprintf("method %s path %s status %d duration %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration.String(),
		)

		if len(c.Errors) > 0 {
			logger.Error(message)
		} else {
			logger.Info(message)
		}
	}
}

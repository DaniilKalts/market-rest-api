package middlewares

import (
	"fmt"
	"time"

	"github.com/DaniilKalts/market-rest-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		message := fmt.Sprintf(
			"Request Log:\n  %-8s: %s\n  %-8s: %s\n  %-8s: %d\n  %-8s: %s",
			"Method", c.Request.Method,
			"Path", c.Request.URL.Path,
			"Status", c.Writer.Status(),
			"Duration", duration.String(),
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				message += fmt.Sprintf("\n  %-8s: %s", "Error", e.Error())
			}
			logger.Error(message)
		} else {
			logger.Info(message)
		}
	}
}

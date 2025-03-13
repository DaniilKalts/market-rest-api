package middlewares

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func BindBodyMiddleware(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		if err := c.ShouldBind(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
			c.Abort()
			return
		}
		c.Set("model", input)
		c.Next()
	}
}

package middlewares

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func BindBodyMiddleware(model interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		input := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		decoder := json.NewDecoder(ctx.Request.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("model", input)
		ctx.Next()
	}
}

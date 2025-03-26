package ginhelpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetContextValue[T any](ctx *gin.Context, key string) (T, error) {
	value, exists := ctx.Get(key)

	var typedValue T
	if !exists {
		return typedValue, fmt.Errorf("key %s not found", key)
	}

	typedValue, ok := value.(T)
	if !ok {
		return typedValue, fmt.Errorf("invalid type for key %s", key)
	}

	return typedValue, nil
}

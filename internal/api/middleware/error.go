package middleware

import (
	"fmt"
	"memoir-api/internal/api/dto"
	"memoir-api/internal/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware is a middleware that handles recovery from panics
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// 处理不同类型的 panic
				var errMsg string
				switch v := rec.(type) {
				case error:
					errMsg = v.Error()
				case string:
					errMsg = v
				default:
					errMsg = fmt.Sprintf("%v", v)
				}

				// 记录错误和堆栈
				logger.Error(nil, "Panic recovered",
					"panic", errMsg,
					"stack", string(debug.Stack()),
				)

				// 返回通用 500 错误
				response := dto.NewErrorResponse(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.JSON(http.StatusInternalServerError, response)
				c.Abort()
			}
		}()
		c.Next()
	}
}

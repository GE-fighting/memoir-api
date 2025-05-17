package middleware

import (
	"net/http"

	"memoir-api/internal/api/errors"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware is a middleware that handles recovery from panics
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				// logger.Error("Panic recovered", map[string]interface{}{"error": err})

				// Return a 500 error to the client
				c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Code:    errors.ErrCodeInternalServer,
					Message: "An unexpected error occurred",
				})

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}

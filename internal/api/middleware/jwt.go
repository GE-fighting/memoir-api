package middleware

import (
	"net/http"
	"strings"

	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware authenticates requests using a JWT token
func JWTAuthMiddleware(services service.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer token format check
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected 'Bearer TOKEN'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 使用JWT服务验证令牌并提取用户ID
		jwtService := services.JWT()
		userID, err := jwtService.ExtractUserID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID to context
		c.Set("userID", userID)
		c.Next()
	}
}

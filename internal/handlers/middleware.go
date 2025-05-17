package handlers

import (
	"errors"
	"fmt"
	"memoir-api/internal/logger"
	"memoir-api/internal/service"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Status  int    `json:"-"`                 // HTTP status code
	Code    string `json:"code"`              // Error code for clients
	Message string `json:"message"`           // User-friendly error message
	Details any    `json:"details,omitempty"` // Optional additional details
}

// Common error codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeValidation          = "VALIDATION_ERROR"
	ErrCodeResourceUnavailable = "RESOURCE_UNAVAILABLE"
)

// Common application errors
var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource conflict")
	ErrInternalServer      = errors.New("internal server error")
	ErrValidation          = errors.New("validation error")
	ErrResourceUnavailable = errors.New("resource unavailable")
)

// NewErrorResponse creates a new error response with the given parameters
func NewErrorResponse(status int, code, message string, details any) ErrorResponse {
	return ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// HandleError sends an appropriate error response based on the error type
func HandleError(c *gin.Context, err error, details ...any) {
	var detail any
	if len(details) > 0 {
		detail = details[0]
	}

	var response ErrorResponse

	switch {
	case errors.Is(err, ErrInvalidInput):
		response = NewErrorResponse(http.StatusBadRequest, ErrCodeBadRequest, err.Error(), detail)
	case errors.Is(err, ErrUnauthorized):
		response = NewErrorResponse(http.StatusUnauthorized, ErrCodeUnauthorized, err.Error(), detail)
	case errors.Is(err, ErrForbidden):
		response = NewErrorResponse(http.StatusForbidden, ErrCodeForbidden, err.Error(), detail)
	case errors.Is(err, ErrNotFound):
		response = NewErrorResponse(http.StatusNotFound, ErrCodeNotFound, err.Error(), detail)
	case errors.Is(err, ErrConflict):
		response = NewErrorResponse(http.StatusConflict, ErrCodeConflict, err.Error(), detail)
	case errors.Is(err, ErrValidation):
		response = NewErrorResponse(http.StatusBadRequest, ErrCodeValidation, err.Error(), detail)
	case errors.Is(err, ErrResourceUnavailable):
		response = NewErrorResponse(http.StatusServiceUnavailable, ErrCodeResourceUnavailable, err.Error(), detail)
	default:
		// Log unexpected errors but don't expose details to clients
		response = NewErrorResponse(http.StatusInternalServerError, ErrCodeInternalServer, "An unexpected error occurred", nil)
	}

	c.JSON(response.Status, response)
}

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware authenticates requests using a JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
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
		claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set claims to context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// validateToken validates a JWT token and returns its claims
func validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get secret key from environment
		jwtSecret := getJWTSecret()
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// JWTAuthWithServiceMiddleware JWT认证中间件
func JWTAuthWithServiceMiddleware(secretKey string, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Authorization头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			HandleError(c, ErrUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Bearer token格式验证
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			HandleError(c, ErrUnauthorized, "Invalid Authorization header format")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 解析和验证JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			HandleError(c, ErrUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// 提取Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			HandleError(c, ErrUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		// 检查令牌类型
		tokenType, _ := claims["typ"].(string)
		if tokenType != "access" {
			HandleError(c, ErrUnauthorized, "Invalid token type")
			c.Abort()
			return
		}

		// 检查过期时间
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			HandleError(c, ErrUnauthorized, "Token expired")
			c.Abort()
			return
		}

		// 获取用户ID
		userID, ok := claims["sub"].(float64)
		if !ok {
			HandleError(c, ErrUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		// 获取用户信息
		user, err := userService.GetUserByID(c, int64(userID))
		if err != nil {
			logger.Error(err, "Failed to get user", map[string]interface{}{"user_id": userID})
			HandleError(c, ErrUnauthorized, "User not found")
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user", user)
		c.Set("user_id", int64(userID))

		c.Next()
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware 恢复中间件，防止服务器因panic而崩溃
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// ErrorMiddleware is a middleware that handles recovery from panics
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				// logger.Error("Panic recovered", map[string]interface{}{"error": err})

				// Return a 500 error to the client
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Status:  http.StatusInternalServerError,
					Code:    ErrCodeInternalServer,
					Message: "An unexpected error occurred",
				})

				// Abort the request
				c.Abort()
			}
		}()

		c.Next()
	}
}

// getJWTSecret returns the JWT secret key from environment or a default one for development
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default for development only
		secret = "memoir-api-development-secret-key"
	}
	return secret
}

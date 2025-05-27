package middleware

import (
	"time"

	"memoir-api/internal/config"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ApplyMiddleware applies all middleware to the given router
func ApplyMiddleware(router *gin.Engine, cfg *config.Config) {
	// Apply recovery and error middleware
	router.Use(ErrorMiddleware())

	// Apply logger middleware with request details
	router.Use(LoggerMiddleware())

	// Apply CORS middleware
	router.Use(CorsMiddleware(cfg))

	// Apply request ID middleware
	router.Use(RequestIDMiddleware())

	// Apply body size limiting middleware
	router.Use(BodySizeLimitMiddleware(cfg.Server.MaxBodySize))
}

// LoggerMiddleware logs HTTP requests with details
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Get request ID
		requestID, _ := c.Get("requestID")

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get method
		method := c.Request.Method

		// Construct query if present
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		log.Info().
			Str("component", "http").
			Str("request_id", requestID.(string)).
			Int("status", statusCode).
			Str("method", method).
			Str("path", path).
			Str("ip", clientIP).
			Dur("latency", latency).
			Int("size", c.Writer.Size()).
			Msg("HTTP Request")
	}
}

// CorsMiddleware configures CORS settings
func CorsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID from header or generate a new one
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Set request ID in header and context
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)

		c.Next()
	}
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	// Use snowflake algorithm from our models package for ID generation
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(result)
}

// BodySizeLimitMiddleware adds a body size limit to each request
func BodySizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

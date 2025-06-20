package middleware

import (
	"time"

	"memoir-api/internal/config"
	"memoir-api/internal/logger"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		// Get or generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Request.Header.Set("X-Request-ID", requestID)
		}

		// Set requestID in the response header
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("requestID", requestID)

		// Create a request-scoped logger with the request ID
		reqLogger := logger.GetLogger("http").With("request_id", requestID)

		// Store the logger in the context
		c.Request = c.Request.WithContext(reqLogger.WithContext(c.Request.Context()))

		// Log the incoming request
		reqLogger.Info("Request started", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"query":  c.Request.URL.RawQuery,
			"ip":     c.ClientIP(),
		})

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Log the completed request
		logFields := map[string]interface{}{
			"status":  statusCode,
			"latency": latency,
			"size":    c.Writer.Size(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
		}

		// Determine log level based on status code
		if statusCode >= 500 {
			reqLogger.Error(c.Errors.Last(), "Request failed", logFields)
		} else if statusCode >= 400 {
			reqLogger.Warn("Request completed with client error", logFields)
		} else {
			reqLogger.Info("Request completed successfully", logFields)
		}
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

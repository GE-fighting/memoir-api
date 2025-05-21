package handlers

import (
	"net/http"
	"time"

	"memoir-api/internal/api/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthCheck represents the system's health check response
type HealthCheck struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// HealthCheckHandler returns a handler for health check endpoint
func HealthCheckHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		services := map[string]string{
			"api": "up",
		}

		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			services["database"] = "down"
		} else {
			if err := sqlDB.Ping(); err != nil {
				services["database"] = "down"
			} else {
				services["database"] = "up"
			}
		}

		// Determine overall status
		status := "healthy"
		for _, s := range services {
			if s != "up" {
				status = "degraded"
				break
			}
		}

		response := HealthCheck{
			Status:    status,
			Timestamp: time.Now(),
			Services:  services,
			Version:   "1.0.0", // You might want to fetch this from a version package or env variable
		}

		c.JSON(http.StatusOK, dto.NewSuccessResponse(response, "系统状态"))
	}
}

// PingHandler handles the ping endpoint request
func PingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.EmptySuccessResponse("pong"))
	}
}

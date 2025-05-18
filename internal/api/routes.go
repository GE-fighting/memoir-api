package api

import (
	"memoir-api/internal/api/middleware"
	"memoir-api/internal/handlers"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all API routes to the given router
func RegisterRoutes(router *gin.Engine, services service.Factory, db *gorm.DB) {
	// Apply middleware
	middleware.ApplyMiddleware(router)

	// Health check
	router.GET("/health", handlers.HealthCheckHandler(db))

	// Simple ping endpoint
	router.GET("/ping", handlers.PingHandler())

	// API v1 group
	v1 := router.Group("/api/v1")

	// Create Auth Handler
	authHandler := handlers.NewAuthHandler(services)

	// Auth routes
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	// Apply JWT auth middleware
	protected := v1.Group("")
	protected.Use(middleware.JWTAuthMiddleware(services))

	// User routes
	userRoutes := protected.Group("/users")
	{
		userRoutes.GET("/me", handlers.GetCurrentUserHandler(services))
		userRoutes.PUT("/me", handlers.UpdateUserHandler(services))
		userRoutes.PUT("/preferences", handlers.UpdateUserPreferencesHandler(services))
	}

	// Couple routes
	coupleRoutes := protected.Group("/couples")
	{
		coupleRoutes.GET("/", handlers.GetCoupleHandler(services))
		coupleRoutes.PUT("/", handlers.UpdateCoupleHandler(services))
		coupleRoutes.PUT("/settings", handlers.UpdateCoupleSettingsHandler(services))
	}

	// Timeline event routes
	eventRoutes := protected.Group("/events")
	{
		eventRoutes.GET("/", handlers.ListTimelineEventsHandler(services))
		eventRoutes.POST("/", handlers.CreateTimelineEventHandler(services))
		eventRoutes.GET("/:id", handlers.GetTimelineEventHandler(services))
		eventRoutes.PUT("/:id", handlers.UpdateTimelineEventHandler(services))
		eventRoutes.DELETE("/:id", handlers.DeleteTimelineEventHandler(services))
	}

	// Location routes
	locationRoutes := protected.Group("/locations")
	{
		locationRoutes.GET("/", handlers.ListLocationsHandler(services))
		locationRoutes.POST("/", handlers.CreateLocationHandler(services))
		locationRoutes.GET("/:id", handlers.GetLocationHandler(services))
		locationRoutes.PUT("/:id", handlers.UpdateLocationHandler(services))
		locationRoutes.DELETE("/:id", handlers.DeleteLocationHandler(services))
	}

	// Photos and videos routes
	mediaRoutes := protected.Group("/media")
	{
		mediaRoutes.GET("/", handlers.ListMediaHandler(services))
		mediaRoutes.POST("/", handlers.UploadMediaHandler(services))
		mediaRoutes.GET("/:id", handlers.GetMediaHandler(services))
		mediaRoutes.PUT("/:id", handlers.UpdateMediaHandler(services))
		mediaRoutes.DELETE("/:id", handlers.DeleteMediaHandler(services))
	}

	// Wishlist routes
	wishlistRoutes := protected.Group("/wishlist")
	{
		wishlistRoutes.GET("/", handlers.ListWishlistItemsHandler(services))
		wishlistRoutes.POST("/", handlers.CreateWishlistItemHandler(services))
		wishlistRoutes.GET("/:id", handlers.GetWishlistItemHandler(services))
		wishlistRoutes.PUT("/:id", handlers.UpdateWishlistItemHandler(services))
		wishlistRoutes.PUT("/:id/status", handlers.UpdateWishlistItemStatusHandler(services))
		wishlistRoutes.DELETE("/:id", handlers.DeleteWishlistItemHandler(services))
	}

	// OSS (Aliyun Object Storage Service) routes
	ossRoutes := protected.Group("/oss")
	{
		ossRoutes.GET("/token", handlers.GenerateSTSToken)
	}
}

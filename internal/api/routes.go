package api

import (
	"memoir-api/internal/api/middleware"
	"memoir-api/internal/config"
	"memoir-api/internal/handlers"
	"memoir-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all API routes to the given router
func RegisterRoutes(router *gin.Engine, services service.Factory, db *gorm.DB, cfg *config.Config) {
	// Apply middleware
	middleware.ApplyMiddleware(router, cfg)

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
		userRoutes.GET("/exist-couple", handlers.ExistCoupleHandler(services))
	}

	// Couple routes
	coupleRoutes := protected.Group("/couple")
	{
		coupleRoutes.POST("/create", handlers.CreateCoupleHandler(services))
		coupleRoutes.GET("/sts", handlers.GenerateCoupleSTSToken(services))
	}

	// Timeline event routes
	eventRoutes := protected.Group("/events")
	{
		eventRoutes.POST("/create", handlers.CreateTimelineEventHandler(services))
	}

	// Location routes
	locationRoutes := protected.Group("/locations")
	{
		locationRoutes.GET("/list", handlers.ListLocationsHandler(services))
		locationRoutes.POST("/create", handlers.CreateLocationHandler(services))
	}

	// Photos and videos routes
	mediaRoutes := protected.Group("/media")
	{
		mediaRoutes.POST("/create", handlers.CreatePhotoVideoHandler(services))
		mediaRoutes.GET("/page", handlers.PageQueryPersonalMediaHandler(services))
	}

	// 个人媒体路由
	// 注册个人媒体处
	personalMediaRoutes := protected.Group("/personal-media")
	{
		personalMediaRoutes.POST("/create", handlers.CreatePersonalMediaWithURLHandler(services))
		personalMediaRoutes.GET("/page", handlers.PageQueryPersonalMediaHandler(services))
	}

	// Wishlist routes
	wishlistRoutes := protected.Group("/wishlist")
	{
		wishlistRoutes.GET("/list", handlers.ListWishlistItemsHandler(services))
		wishlistRoutes.POST("/create", handlers.CreateWishlistItemHandler(services))
		wishlistRoutes.PUT("/update", handlers.UpdateWishlistItemHandler(services))
		wishlistRoutes.PUT("/:id/status", handlers.UpdateWishlistItemStatusHandler(services))
		wishlistRoutes.DELETE("/:id", handlers.DeleteWishlistItemHandler(services))
	}

	// 情侣相册路由
	albumRoutes := protected.Group("/albums")
	{
		albumRoutes.GET("/list", handlers.ListCoupleAlbumsHandler(services))
		albumRoutes.POST("/create", handlers.CreateCoupleAlbumHandler(services))
		albumRoutes.GET("/photos", handlers.GetCoupleAlbumWithPhotosHandler(services))
	}

	// OSS (Aliyun Object Storage Service) routes
	ossRoutes := protected.Group("/oss")
	{
		ossRoutes.GET("/token", handlers.GenerateSTSToken)
	}
}

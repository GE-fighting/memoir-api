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

	// Create Email Handler
	emailHandler := handlers.NewEmailHandler(services)

	// Auth routes
	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	// Email verification routes (public)
	emailRoutes := v1.Group("/email")
	{
		emailRoutes.POST("/verify", emailHandler.VerifyEmail)
		emailRoutes.POST("/resend-code", emailHandler.ResendVerificationCode)
		emailRoutes.POST("/forgot-password", emailHandler.ForgotPassword)
		emailRoutes.POST("/reset-password", emailHandler.ResetPassword)
	}

	// Protected routes
	// Apply JWT auth middleware
	protected := v1.Group("")
	protected.Use(middleware.JWTAuthMiddleware(services))

	//Dashboard routes
	dashboardRoutes := protected.Group("/dashboard")
	{
		dashboardRoutes.GET("", handlers.GetDashboardDataHandler(services))
	}

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
		coupleRoutes.GET("/info", handlers.GetCoupleInfoHandler(services))
	}

	// Timeline event routes
	eventRoutes := protected.Group("/events")
	{
		eventRoutes.POST("/create", handlers.CreateTimelineEventHandler(services))
		eventRoutes.GET("/page", handlers.PageTimelineEventsHandler(services))
		eventRoutes.GET("/:id", handlers.GetTimelineEventHandler(services))
		eventRoutes.DELETE("/:id", handlers.DeleteTimelineEventHandler(services))
		eventRoutes.PUT("/:id", handlers.UpdateTimelineEventHandler(services))
	}

	// Location routes
	locationRoutes := protected.Group("/locations")
	{
		locationRoutes.GET("/list", handlers.ListLocationsHandler(services))
		locationRoutes.GET("/:id", handlers.GetLocationHandler(services))
		locationRoutes.POST("/create", handlers.CreateLocationHandler(services))
		locationRoutes.DELETE("/:id", handlers.DeleteLocationHandler(services))
	}

	// Photos and videos routes
	mediaRoutes := protected.Group("/media")
	{
		mediaRoutes.POST("/create", handlers.CreatePhotoVideoHandler(services))
		mediaRoutes.GET("/page", handlers.ListPhotoVideoHandler(services))
	}

	// 个人媒体路由
	// 注册个人媒体处
	personalMediaRoutes := protected.Group("/personal-media")
	{
		personalMediaRoutes.POST("/create", handlers.CreatePersonalMediaWithURLHandler(services))
		personalMediaRoutes.GET("/page", handlers.PageQueryPersonalMediaHandler(services))
		personalMediaRoutes.DELETE("/:id", handlers.DeletePersonalMediaHandler(services))
	}

	// Wishlist routes
	wishlistRoutes := protected.Group("/wishlist")
	{
		wishlistRoutes.GET("/list", handlers.ListWishlistItemsHandler(services))
		wishlistRoutes.POST("/create", handlers.CreateWishlistItemHandler(services))
		wishlistRoutes.PUT("/update", handlers.UpdateWishlistItemHandler(services))
		wishlistRoutes.PUT("/:id/status", handlers.UpdateWishlistItemStatusHandler(services))
		wishlistRoutes.DELETE("/:id", handlers.DeleteWishlistItemHandler(services))
		wishlistRoutes.POST("/associateAttachments", handlers.AssociateAttachments(services))
	}

	// 情侣相册路由
	albumRoutes := protected.Group("/albums")
	{
		albumRoutes.GET("/list", handlers.ListCoupleAlbumsHandler(services))
		albumRoutes.POST("/create", handlers.CreateCoupleAlbumHandler(services))
		albumRoutes.GET("/photos", handlers.GetCoupleAlbumWithPhotosHandler(services))
		albumRoutes.DELETE("/:id", handlers.DeleteCoupleAlbumHandler(services))
		albumRoutes.POST("/deletePhotos", handlers.DeleteCoupleAlbumPhotosHandler(services))
		albumRoutes.GET("/all-media/page", handlers.PageCoupleMedia(services))
	}

	// 附件路由
	attachmentRoutes := protected.Group("/attachments")
	{
		attachmentRoutes.POST("/create", handlers.CreateAttachmentHandler(services))
		attachmentRoutes.GET("/:id", handlers.GetAttachmentHandler(services))
		attachmentRoutes.GET("/list", handlers.ListAttachmentsHandler(services))
		attachmentRoutes.DELETE("/:id", handlers.DeleteAttachmentHandler(services))
	}

	// OSS (Aliyun Object Storage Service) routes
	ossRoutes := protected.Group("/oss")
	{
		ossRoutes.GET("/token", handlers.GenerateSTSToken)
	}

	// 管理员路由 - 用于手动触发提醒功能
	adminRoutes := protected.Group("/reminders")
	{
		// 触发纪念日提醒
		adminRoutes.POST("/anniversary", handlers.TriggerAnniversaryRemindersHandler(services))
		// 触发节日提醒
		adminRoutes.POST("/festival", handlers.TriggerFestivalRemindersHandler(services))
	}
}

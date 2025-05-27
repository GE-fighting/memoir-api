package main

import (
	"context"
	"fmt"
	"memoir-api/internal/api"
	"memoir-api/internal/cache"
	"memoir-api/internal/config"
	"memoir-api/internal/db"
	"memoir-api/internal/logger"
	"memoir-api/internal/models"
	"memoir-api/internal/repository"
	"memoir-api/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {

	// Load configuration
	cfg := config.New()

	// Initialize logger
	logger.Initialize(cfg.Server.LogLevel)

	// Set Gin mode based on environment
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	dbConn, err := db.NewDB(cfg.DB.ConnectionString())
	if err != nil {
		logger.Fatal(err, "Failed to connect to database")
	}

	// Configure connection pool
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Fatal(err, "Failed to get database connection")
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate database schemas
	if err := autoMigrateDB(dbConn); err != nil {
		logger.Fatal(err, "Failed to migrate database")
	}

	// Initialize Redis
	if err := cache.Init(cfg); err != nil {
		logger.Fatal(err, "Failed to initialize Redis")
	}
	// Ensure Redis connection is closed on exit
	defer func() {
		if err := cache.Close(); err != nil {
			logger.Error(err, "Error closing Redis connection")
		} else {
			logger.Info("Redis connection closed")
		}
	}()

	// Initialize repositories
	repoFactory := repository.NewFactory(dbConn)

	// Initialize services
	serviceFactory := service.NewFactory(repoFactory)

	// Setup Gin router
	router := gin.Default()

	// Register API routes
	api.RegisterRoutes(router, serviceFactory, dbConn, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Starting Memoir API server on port %d...", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err, "Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create timeout context for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(err, "Server forced to shutdown")
	}

	// Close database connection
	if err := sqlDB.Close(); err != nil {
		logger.Fatal(err, "Error closing database connection")
	}

	logger.Info("Server exiting")
}

// autoMigrateDB handles database migrations
func autoMigrateDB(db interface{}) error {
	gormDB := db.(*gorm.DB)

	log.Info().Msg("Running database migrations...")

	// 使用AutoMigrate进行增量迁移（添加表/列，但不删除）
	if err := gormDB.AutoMigrate(
		&models.User{},
		&models.Couple{},
		&models.Location{},
		&models.TimelineEvent{},
		&models.PhotoVideo{},
		&models.Wishlist{},
		&models.PersonalMedia{},
		&models.CoupleAlbum{},
	); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}

package main

import (
	"context"
	"fmt"
	"memoir-api/internal/api"
	"memoir-api/internal/cache"
	"memoir-api/internal/config"
	"memoir-api/internal/db"
	"memoir-api/internal/logger"
	"memoir-api/internal/repository"
	"memoir-api/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
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
	dbConn, err := db.NewDB(&cfg.DB)
	if err != nil {
		logger.Fatal(err, "Failed to connect to database")
	}

	// Get SQL DB instance for graceful shutdown
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Fatal(err, "Failed to get database connection")
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

	// Create application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start email queue processor if enabled
	if cfg.Email.Enabled {
		logger.Info("启动邮件队列处理")
		go serviceFactory.Email().ProcessEmailQueue(ctx)
	}

	// 设置定时任务
	if cfg.Email.Enabled {
		logger.Info("设置纪念日和节日提醒定时任务")
		c := cron.New()

		// 每天早上9点检查纪念日
		_, err := c.AddFunc("0 9 * * *", func() {
			logger.Info("执行纪念日检查任务")
			err := serviceFactory.CoupleReminder().CheckAndSendAnniversaryReminders(context.Background())
			if err != nil {
				logger.Error(err, "纪念日检查任务失败")
			}
		})
		if err != nil {
			logger.Error(err, "添加纪念日检查任务失败")
		}

		// 每天早上10点检查节日
		_, err = c.AddFunc("0 10 * * *", func() {
			logger.Info("执行节日检查任务")
			err := serviceFactory.CoupleReminder().CheckAndSendFestivalReminders(context.Background())
			if err != nil {
				logger.Error(err, "节日检查任务失败")
			}
		})
		if err != nil {
			logger.Error(err, "添加节日检查任务失败")
		}

		// 启动定时任务
		c.Start()

		// 确保在程序退出时停止定时任务
		defer c.Stop()
	}

	// Setup Gin router
	router := gin.New()

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
		logger.Info("Starting Memoir API server", "port", cfg.Server.Port)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
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

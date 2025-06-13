package main

import (
	"flag"
	"fmt"
	"memoir-api/internal/config"
	"memoir-api/internal/db"
	"memoir-api/internal/logger"
	"memoir-api/internal/models"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func main() {
	// 定义命令行参数
	var (
		action = flag.String("action", "", "Migration action: up, down, status")
		help   = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 如果没有指定 action，显示帮助信息
	if *action == "" {
		fmt.Println("Error: No action specified")
		fmt.Println()
		showHelp()
		os.Exit(1)
	}

	// 加载配置
	cfg := config.New()

	// 初始化日志
	logger.Initialize(cfg.Server.LogLevel)

	// 连接数据库
	dbConn, err := db.NewDB(&cfg.DB)
	if err != nil {
		logger.Fatal(err, "Failed to connect to database")
	}

	// 获取底层 SQL 连接以便关闭
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Fatal(err, "Failed to get database connection")
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error(err, "Error closing database connection")
		}
	}()

	// 执行迁移操作
	switch *action {
	case "up":
		if err := migrateUp(dbConn); err != nil {
			logger.Fatal(err, "Migration up failed")
		}
	case "down":
		if err := migrateDown(dbConn); err != nil {
			logger.Fatal(err, "Migration down failed")
		}
	case "status":
		if err := migrateStatus(dbConn); err != nil {
			logger.Fatal(err, "Migration status check failed")
		}
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
}

// migrateUp 执行数据库迁移（创建/更新表）
func migrateUp(db *gorm.DB) error {
	log.Info().Msg("Running database migrations...")

	// 使用AutoMigrate进行增量迁移（添加表/列，但不删除）
	if err := db.AutoMigrate(
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

// migrateDown 回滚数据库迁移（删除表）
func migrateDown(db *gorm.DB) error {
	log.Info().Msg("Rolling back database migrations...")

	// 按依赖关系逆序删除表
	tables := []interface{}{
		&models.CoupleAlbum{},
		&models.PersonalMedia{},
		&models.Wishlist{},
		&models.PhotoVideo{},
		&models.TimelineEvent{},
		&models.Location{},
		&models.User{},
		&models.Couple{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.Warn().Err(err).Msgf("Failed to drop table for %T", table)
		} else {
			log.Info().Msgf("Dropped table for %T", table)
		}
	}

	log.Info().Msg("Database rollback completed")
	return nil
}

// migrateStatus 检查数据库迁移状态
func migrateStatus(db *gorm.DB) error {
	log.Info().Msg("Checking database migration status...")

	// 定义模型和对应的表名
	modelInfo := []struct {
		model     interface{}
		tableName string
	}{
		{&models.User{}, "users"},
		{&models.Couple{}, "couples"},
		{&models.Location{}, "locations"},
		{&models.TimelineEvent{}, "timeline_events"},
		{&models.PhotoVideo{}, "photo_videos"},
		{&models.Wishlist{}, "wishlists"},
		{&models.PersonalMedia{}, "personal_media"},
		{&models.CoupleAlbum{}, "couple_albums"},
	}

	for _, info := range modelInfo {
		if db.Migrator().HasTable(info.model) {
			log.Info().Msgf("✓ Table exists: %s", info.tableName)
		} else {
			log.Info().Msgf("✗ Table missing: %s", info.tableName)
		}
	}

	return nil
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -action string")
	fmt.Println("        Migration action: up, down, status (required)")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go -action=up     # Run migrations")
	fmt.Println("  go run cmd/migrate/main.go -action=down   # Rollback migrations")
	fmt.Println("  go run cmd/migrate/main.go -action=status # Check migration status")
}

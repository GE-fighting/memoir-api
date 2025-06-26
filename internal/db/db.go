package db

import (
	"context"
	"fmt"
	"memoir-api/internal/config"
	"memoir-api/internal/logger"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// CustomGormLogger implements gorm.io/gorm/logger.Interface
type CustomGormLogger struct {
	Logger        logger.Logger
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
	Colorful      bool
}

// LogMode sets the log level for the logger
func (l *CustomGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info prints info messages
func (l *CustomGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Logger.With("request_id", l.getRequestId(ctx)).Info(msg, data...)
	}
}

// Warn prints warn messages
func (l *CustomGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Logger.With("request_id", l.getRequestId(ctx)).Warn(fmt.Sprintf(msg, data...))
	}
}

// Error prints error messages
func (l *CustomGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Logger.With("request_id", l.getRequestId(ctx)).Error(nil, fmt.Sprintf(msg, data...))
	}
}

// Trace prints trace messages for SQL operations
func (l *CustomGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// Format SQL for better readability by removing escape characters
	formattedSQL := formatSQL(sql)

	customLogger := l.Logger.With(
		"request_id", l.getRequestId(ctx),
		"elapsed", elapsed.String(),
		"rows", rows,
	)

	// Log based on error and execution time
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error:
		customLogger.Error(err, "SQL error", "sql", formattedSQL)
	case elapsed > l.SlowThreshold && l.SlowThreshold > 0 && l.LogLevel >= gormlogger.Warn:
		customLogger.Warn("SLOW SQL", "sql", formattedSQL, "threshold", l.SlowThreshold)
	case l.LogLevel >= gormlogger.Info:
		customLogger.Info("SQL", "sql", formattedSQL)
	}
}

// formatSQL 格式化 SQL 语句，移除转义字符以提高可读性
func formatSQL(sql string) string {
	// 处理转义的双引号
	sql = strings.ReplaceAll(sql, `\"`, `"`)

	// 处理转义的单引号
	sql = strings.ReplaceAll(sql, `\'`, `'`)

	// 处理双反斜杠
	sql = strings.ReplaceAll(sql, `\\`, `\`)

	// 处理 PostgreSQL 特定的语法
	sql = strings.ReplaceAll(sql, `"`, `'`)

	// 移除多余的空格
	sql = strings.TrimSpace(sql)

	return sql
}

// getRequestId returns the request ID from context if available
func (l *CustomGormLogger) getRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 尝试从上下文中获取请求 ID
	requestID, ok := ctx.Value("requestID").(string)
	if !ok {
		return ""
	}

	return requestID
}

// NewGormLogger creates a new CustomGormLogger
func NewGormLogger(logLevel string) *CustomGormLogger {
	var level gormlogger.LogLevel
	switch logLevel {
	case "debug":
		level = gormlogger.Info
	case "info":
		level = gormlogger.Info
	case "warn":
		level = gormlogger.Warn
	case "error":
		level = gormlogger.Error
	default:
		level = gormlogger.Info
	}

	return &CustomGormLogger{
		Logger:        logger.GetLogger("gorm"),
		LogLevel:      level,
		SlowThreshold: 200 * time.Millisecond,
		Colorful:      true,
	}
}

// NewDB creates and returns a GORM database connection
func NewDB(dbConfig *config.DBConfig) (*gorm.DB, error) {
	connectionString := dbConfig.ConnectionString()

	// Create custom GORM logger
	gormLogger := NewGormLogger(dbConfig.LogLevel)

	// Open database connection
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 使用配置文件中的连接池参数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Minute)

	return db, nil
}

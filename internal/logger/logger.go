package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// loggerKey is the key used to store and retrieve the logger from a context.
const loggerKey = contextKey("logger")

// Logger is an interface for a structured logger, designed to be slog-native.
type Logger interface {
	// Core logging methods
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(err error, msg string, args ...any)
	Fatal(err error, msg string, args ...any)

	// Context methods
	WithContext(ctx context.Context) context.Context
	FromContext(ctx context.Context) Logger

	// Field methods
	With(args ...any) Logger
	WithComponent(component string) Logger
	WithError(err error) Logger

	// GetLogger returns the underlying slog.Logger.
	GetLogger() *slog.Logger
}

// logger is the concrete implementation of the Logger interface using slog.
type logger struct {
	slog *slog.Logger
}

// defaultLogger is the application's default logger instance.
var defaultLogger Logger = &logger{slog: slog.Default()}

// Initialize sets up the global logger with slog.
func Initialize(level string) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo // Default to Info level on parse error.
	}

	// 使用tint创建一个带颜色的日志处理器
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      lvl,
		TimeFormat: time.Kitchen, // 可以自定义时间格式，如 3:04:05PM
		AddSource:  false,        // 同样支持显示源码位置
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 为不同的字段设置不同的颜色，使用更优雅的配色方案
			if len(groups) == 0 {
				switch a.Key {
				case "component":
					// 深青色 (36) 用于组件名称 - 醒目但不刺眼
					return tint.Attr(36, a)
				case "request_id":
					// 浅灰色 (37) 用于请求ID - 重要但不需要太醒目
					return tint.Attr(37, a)
				case "sql":
					// 浅蓝色 (94) 用于SQL语句 - 易于阅读的代码颜色
					return tint.Attr(94, a)
				case "elapsed":
					// 浅黄色 (93) 用于执行时间 - 性能指标应该比较醒目
					return tint.Attr(93, a)
				case "rows":
					// 浅绿色 (92) 用于行数 - 表示结果数量
					return tint.Attr(92, a)
				case "threshold":
					// 浅紫色 (95) 用于阈值 - 配置值
					return tint.Attr(95, a)
				case slog.LevelKey:
					// 为不同的日志级别设置不同的颜色
					level := a.Value.Any().(slog.Level)
					switch {
					case level == slog.LevelDebug:
						return tint.Attr(90, a) // 灰色 - 不太重要的调试信息
					case level == slog.LevelInfo:
						return tint.Attr(32, a) // 绿色 - 正常信息
					case level == slog.LevelWarn:
						return tint.Attr(33, a) // 黄色 - 警告
					case level == slog.LevelError:
						return tint.Attr(31, a) // 红色 - 错误
					}
				}

				// 检查是否为错误类型，错误应该保持红色以便于识别
				if a.Value.Kind() == slog.KindAny {
					if _, ok := a.Value.Any().(error); ok {
						return tint.Attr(31, a) // 红色用于错误
					}
				}
			}
			return a
		},
	})

	slogLogger := slog.New(handler)

	// Set the global default logger for the application.
	slog.SetDefault(slogLogger)

	// Update our custom defaultLogger wrapper.
	defaultLogger = &logger{slog: slogLogger}
}

// GetLogger creates a logger with a "component" field.
func GetLogger(component string) Logger {
	return defaultLogger.WithComponent(component)
}

// GetLoggerFromContext retrieves a logger from the context.
// If no logger is found, it returns the default logger.
func GetLoggerFromContext(ctx context.Context) Logger {
	if ctx == nil {
		return defaultLogger
	}
	if l, ok := ctx.Value(loggerKey).(Logger); ok {
		return l
	}
	return defaultLogger
}

// Global logging functions that use the default logger.

// Debug logs a debug message.
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message.
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message.
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message.
func Error(err error, msg string, args ...any) {
	defaultLogger.Error(err, msg, args...)
}

// Fatal logs a fatal message and exits the application.
func Fatal(err error, msg string, args ...any) {
	defaultLogger.Fatal(err, msg, args...)
}

// WithContext adds the logger to the provided context.
func WithContext(ctx context.Context) context.Context {
	return defaultLogger.WithContext(ctx)
}

// FromContext retrieves a logger from the context or returns the default logger.
func FromContext(ctx context.Context) Logger {
	return GetLoggerFromContext(ctx)
}

// Implementation of Logger interface methods for the logger struct.

func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

func (l *logger) Error(err error, msg string, args ...any) {
	newArgs := make([]any, 0, len(args)+1)
	if err != nil {
		newArgs = append(newArgs, slog.Any("error", err))
	}
	newArgs = append(newArgs, args...)
	l.slog.Error(msg, newArgs...)
}

func (l *logger) Fatal(err error, msg string, args ...any) {
	l.Error(err, msg, args...)
	os.Exit(1)
}

func (l *logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func (l *logger) FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	if l, ok := ctx.Value(loggerKey).(Logger); ok {
		return l
	}
	return l
}

func (l *logger) With(args ...any) Logger {
	return &logger{slog: l.slog.With(args...)}
}

func (l *logger) WithComponent(component string) Logger {
	return &logger{slog: l.slog.With("component", component)}
}

func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l
	}
	return &logger{slog: l.slog.With("error", err)}
}

func (l *logger) GetLogger() *slog.Logger {
	return l.slog
}

package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// loggerKey is the key used to store and retrieve the logger from a context
const loggerKey = contextKey("logger")

// Logger represents a logger with a fluent API
type Logger interface {
	// Core logging methods
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(err error, msg string, fields ...map[string]interface{})
	Fatal(err error, msg string, fields ...map[string]interface{})

	// Context methods
	WithContext(ctx context.Context) context.Context
	FromContext(ctx context.Context) Logger

	// Field methods
	With(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithComponent(component string) Logger
	WithError(err error) Logger

	// Get the underlying zerolog.Logger for advanced use cases
	GetLogger() zerolog.Logger
}

// logger is the concrete implementation of the Logger interface
type logger struct {
	zerologger zerolog.Logger
}

// defaultLogger is the application's default logger instance
var defaultLogger Logger

// Initialize sets up the global logger with proper configuration
func Initialize(level string) {
	// Set logger level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure logger output
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatCaller: func(i interface{}) string {
			if i == nil {
				return ""
			}
			return fmt.Sprintf(" %s", i)
		},
	}

	// Set global logger
	zerologInstance := zerolog.New(output).
		With().
		Timestamp().
		CallerWithSkipFrameCount(3).
		Logger()

	// Set our default logger
	defaultLogger = &logger{zerologger: zerologInstance}
}

// GetLogger creates a logger with component context
func GetLogger(component string) Logger {
	return defaultLogger.WithComponent(component)
}

// GetLoggerFromContext retrieves a logger from the context
// If no logger is found in the context, returns the default logger
func GetLoggerFromContext(ctx context.Context) Logger {
	if ctx == nil {
		return defaultLogger
	}

	if loggerFromCtx, ok := ctx.Value(loggerKey).(Logger); ok {
		return loggerFromCtx
	}

	return defaultLogger
}

// SetOutput changes the output writer for testing purposes
func SetOutput(w io.Writer) {
	output := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
	}
	zerologInstance := zerolog.New(output).With().Timestamp().Caller().Logger()
	defaultLogger = &logger{zerologger: zerologInstance}
}

// Debug logs a debug message
func Debug(msg string, fields ...map[string]interface{}) {
	defaultLogger.Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...map[string]interface{}) {
	defaultLogger.Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...map[string]interface{}) {
	defaultLogger.Warn(msg, fields...)
}

// Error logs an error message
func Error(err error, msg string, fields ...map[string]interface{}) {
	defaultLogger.Error(err, msg, fields...)
}

// Fatal logs a fatal message and exits the application
func Fatal(err error, msg string, fields ...map[string]interface{}) {
	defaultLogger.Fatal(err, msg, fields...)
}

// WithContext adds the logger to the provided context
func WithContext(ctx context.Context) context.Context {
	return defaultLogger.WithContext(ctx)
}

// FromContext retrieves a logger from the context or returns the default logger
func FromContext(ctx context.Context) Logger {
	return GetLoggerFromContext(ctx)
}

// Implementation of Logger interface methods

func (l *logger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.zerologger.Debug()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func (l *logger) Info(msg string, fields ...map[string]interface{}) {
	event := l.zerologger.Info()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func (l *logger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.zerologger.Warn()
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func (l *logger) Error(err error, msg string, fields ...map[string]interface{}) {
	event := l.zerologger.Error().Err(err)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func (l *logger) Fatal(err error, msg string, fields ...map[string]interface{}) {
	event := l.zerologger.Fatal().Err(err)
	if len(fields) > 0 {
		for k, v := range fields[0] {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func (l *logger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func (l *logger) FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	if loggerFromCtx, ok := ctx.Value(loggerKey).(Logger); ok {
		return loggerFromCtx
	}

	return l
}

func (l *logger) With(key string, value interface{}) Logger {
	newZerologger := l.zerologger.With().Interface(key, value).Logger()
	return &logger{zerologger: newZerologger}
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	ctx := l.zerologger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return &logger{zerologger: ctx.Logger()}
}

func (l *logger) WithComponent(component string) Logger {
	return l.With("component", component)
}

func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l
	}
	newZerologger := l.zerologger.With().Err(err).Logger()
	return &logger{zerologger: newZerologger}
}

func (l *logger) GetLogger() zerolog.Logger {
	return l.zerologger
}

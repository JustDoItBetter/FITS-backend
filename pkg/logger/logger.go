package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// Init initializes the global logger based on environment and configuration
// Production uses JSON format for log aggregation, development uses console for readability
func Init(level, format string) error {
	var config zap.Config

	if format == "json" {
		// JSON format for production - enables log aggregation tools like ELK, Datadog
		config = zap.NewProductionConfig()
	} else {
		// Console format for development - human-readable output
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Parse log level from config
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		parsedLevel = zapcore.InfoLevel // Fallback to info if invalid level
	}
	config.Level = zap.NewAtomicLevelAt(parsedLevel)

	// Build logger
	logger, err := config.Build(zap.AddCallerSkip(1)) // Skip caller to show actual call site, not wrapper
	if err != nil {
		return err
	}

	globalLogger = logger
	return nil
}

// Get returns the global logger instance
// Panics if logger not initialized - ensures Init is called during app startup
func Get() *zap.Logger {
	if globalLogger == nil {
		// Fallback to stdout logger if not initialized (shouldn't happen in production)
		logger, _ := zap.NewProduction()
		return logger
	}
	return globalLogger
}

// Sync flushes any buffered log entries
// Should be called before application shutdown to prevent log loss
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// Helper functions for common log patterns

// Info logs an informational message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Debug logs a debug message (only visible when log level is debug)
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Fatal logs a fatal message and exits the application
// Use sparingly - only for unrecoverable errors during startup
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
	os.Exit(1)
}

// With creates a child logger with pre-configured fields
// Useful for adding request context (user_id, request_id) to all subsequent logs
func With(fields ...zap.Field) *zap.Logger {
	return Get().With(fields...)
}

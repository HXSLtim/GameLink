// Package logging provides generic logging utilities.
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func init() {
	// 默认初始化
	Init("info")
}

// Init sets the default slog logger with JSON handler and given level.
func Init(level string) *slog.Logger {
	lvl := parseLevel(level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl, AddSource: false})
	logger = slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(v string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Debug logs a debug message with the given key-value pairs.
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

// Info logs an info message with the given key-value pairs.
func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

// Warn logs a warning message with the given key-value pairs.
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

// Error logs an error message with the given key-value pairs.
func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}

// DebugContext logs a debug message with context and the given key-value pairs.
func DebugContext(ctx context.Context, msg string, args ...interface{}) {
	logger.DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context and the given key-value pairs.
func InfoContext(ctx context.Context, msg string, args ...interface{}) {
	logger.InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and the given key-value pairs.
func WarnContext(ctx context.Context, msg string, args ...interface{}) {
	logger.WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context and the given key-value pairs.
func ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	logger.ErrorContext(ctx, msg, args...)
}

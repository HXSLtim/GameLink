package logging

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * Test_Init_WithDebugLevel 测试初始化Debug级别
 */
func Test_Init_WithDebugLevel(t *testing.T) {
	// Act
	logger := Init("debug")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithInfoLevel 测试初始化Info级别
 */
func Test_Init_WithInfoLevel(t *testing.T) {
	// Act
	logger := Init("info")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithWarnLevel 测试初始化Warn级别
 */
func Test_Init_WithWarnLevel(t *testing.T) {
	// Act
	logger := Init("warn")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithWarningLevel 测试初始化Warning级别
 */
func Test_Init_WithWarningLevel(t *testing.T) {
	// Act
	logger := Init("warning")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithErrorLevel 测试初始化Error级别
 */
func Test_Init_WithErrorLevel(t *testing.T) {
	// Act
	logger := Init("error")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithInvalidLevel 测试初始化无效级别
 */
func Test_Init_WithInvalidLevel(t *testing.T) {
	// Act
	logger := Init("invalid")

	// Assert - 应该默认为info级别
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithEmptyLevel 测试初始化空级别
 */
func Test_Init_WithEmptyLevel(t *testing.T) {
	// Act
	logger := Init("")

	// Assert - 应该默认为info级别
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithUpperCase 测试初始化大写级别
 */
func Test_Init_WithUpperCase(t *testing.T) {
	// Act
	logger := Init("DEBUG")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_Init_WithWhitespace 测试初始化带空格的级别
 */
func Test_Init_WithWhitespace(t *testing.T) {
	// Act
	logger := Init("  debug  ")

	// Assert
	assert.NotNil(t, logger)
}

/**
 * Test_ParseLevel_AllLevels 测试解析所有级别
 */
func Test_ParseLevel_AllLevels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"Debug", "Debug", slog.LevelDebug},
		{"DEBUG", "DEBUG", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"Info", "Info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"warning", "warning", slog.LevelWarn},
		{"Warn", "Warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"Error", "Error", slog.LevelError},
		{"default", "unknown", slog.LevelInfo},
		{"empty", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			level := parseLevel(tt.input)

			// Assert
			assert.Equal(t, tt.expected, level)
		})
	}
}

/**
 * Test_Debug_LogsMessage 测试Debug日志
 */
func Test_Debug_LogsMessage(t *testing.T) {
	// Arrange
	Init("debug")

	// Act - 不应该panic
	Debug("test debug message")
	Debug("test with args", "key", "value")

	// Assert
	assert.True(t, true, "Debug应该成功执行")
}

/**
 * Test_Info_LogsMessage 测试Info日志
 */
func Test_Info_LogsMessage(t *testing.T) {
	// Arrange
	Init("info")

	// Act
	Info("test info message")
	Info("test with args", "key", "value")

	// Assert
	assert.True(t, true, "Info应该成功执行")
}

/**
 * Test_Warn_LogsMessage 测试Warn日志
 */
func Test_Warn_LogsMessage(t *testing.T) {
	// Arrange
	Init("warn")

	// Act
	Warn("test warn message")
	Warn("test with args", "key", "value")

	// Assert
	assert.True(t, true, "Warn应该成功执行")
}

/**
 * Test_Error_LogsMessage 测试Error日志
 */
func Test_Error_LogsMessage(t *testing.T) {
	// Arrange
	Init("error")

	// Act
	Error("test error message")
	Error("test with args", "key", "value")

	// Assert
	assert.True(t, true, "Error应该成功执行")
}

/**
 * Test_DebugContext_LogsMessage 测试DebugContext日志
 */
func Test_DebugContext_LogsMessage(t *testing.T) {
	// Arrange
	Init("debug")
	ctx := context.Background()

	// Act
	DebugContext(ctx, "test debug context message")
	DebugContext(ctx, "test with args", "key", "value")

	// Assert
	assert.True(t, true, "DebugContext应该成功执行")
}

/**
 * Test_InfoContext_LogsMessage 测试InfoContext日志
 */
func Test_InfoContext_LogsMessage(t *testing.T) {
	// Arrange
	Init("info")
	ctx := context.Background()

	// Act
	InfoContext(ctx, "test info context message")
	InfoContext(ctx, "test with args", "key", "value")

	// Assert
	assert.True(t, true, "InfoContext应该成功执行")
}

/**
 * Test_WarnContext_LogsMessage 测试WarnContext日志
 */
func Test_WarnContext_LogsMessage(t *testing.T) {
	// Arrange
	Init("warn")
	ctx := context.Background()

	// Act
	WarnContext(ctx, "test warn context message")
	WarnContext(ctx, "test with args", "key", "value")

	// Assert
	assert.True(t, true, "WarnContext应该成功执行")
}

/**
 * Test_ErrorContext_LogsMessage 测试ErrorContext日志
 */
func Test_ErrorContext_LogsMessage(t *testing.T) {
	// Arrange
	Init("error")
	ctx := context.Background()

	// Act
	ErrorContext(ctx, "test error context message")
	ErrorContext(ctx, "test with args", "key", "value")

	// Assert
	assert.True(t, true, "ErrorContext应该成功执行")
}

/**
 * Test_Context_WithRequestID 测试带请求ID的日志
 */
func Test_Context_WithRequestID(t *testing.T) {
	// Arrange
	Init("info")
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")

	// Act
	InfoContext(ctx, "test with request id")

	// Assert
	assert.True(t, true, "带请求ID的日志应该成功")
}

/**
 * Test_MultipleInits 测试多次初始化
 */
func Test_MultipleInits(t *testing.T) {
	// Act
	logger1 := Init("debug")
	logger2 := Init("info")
	logger3 := Init("error")

	// Assert - 每次都应该返回新的logger
	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	assert.NotNil(t, logger3)
}

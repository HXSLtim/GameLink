package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlogLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("记录成功的HTTP请求", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证日志输出
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.Equal(t, "INFO", logEntry["level"])
		assert.Equal(t, "http_request", logEntry["msg"])
		assert.Equal(t, float64(200), logEntry["status"])
		assert.Equal(t, "GET", logEntry["method"])
		assert.Equal(t, "/test", logEntry["path"])
		assert.NotEmpty(t, logEntry["duration"])
	})

	t.Run("记录带有RequestID的HTTP请求", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(RequestID()) // 先添加RequestID中间件
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证日志输出包含request_id
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.NotEmpty(t, logEntry["request_id"])
		assert.IsType(t, "", logEntry["request_id"])
	})

	t.Run("记录带有用户ID的HTTP请求", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", uint64(12345))
			c.Next()
		})
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证日志输出包含user_id
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.NotNil(t, logEntry["user_id"])
		assert.Equal(t, float64(12345), logEntry["user_id"]) // JSON数字解析为float64
	})

	t.Run("记录4xx错误为WARN级别", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		// 验证日志级别为WARN
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.Equal(t, "WARN", logEntry["level"])
		assert.Equal(t, float64(400), logEntry["status"])
	})

	t.Run("记录5xx错误为ERROR级别", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// 验证日志级别为ERROR
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.Equal(t, "ERROR", logEntry["level"])
		assert.Equal(t, float64(500), logEntry["status"])
	})

	t.Run("记录请求持续时间", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			// 模拟一些处理时间
			time.Sleep(10 * time.Millisecond)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证日志输出包含duration
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.NotEmpty(t, logEntry["duration"])
		durationStr := logEntry["duration"].(string)
		assert.Contains(t, durationStr, "ms")
	})
}

func TestSlogLogger_Context(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("使用请求上下文", func(t *testing.T) {
		var buf bytes.Buffer
		handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		slog.SetDefault(slog.New(handler))

		router := gin.New()
		router.Use(SlogLogger())
		router.GET("/test", func(c *gin.Context) {
			// 在请求上下文中添加值
			ctx := context.WithValue(c.Request.Context(), "trace_id", "test-trace-123")
			c.Request = c.Request.WithContext(ctx)
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证日志输出
		var logEntry map[string]interface{}
		require.NoError(t, json.Unmarshal(buf.Bytes(), &logEntry))

		assert.Equal(t, "INFO", logEntry["level"])
		assert.Equal(t, "http_request", logEntry["msg"])
	})
}

func BenchmarkSlogLogger(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// 使用discard writer避免I/O影响
	handler := slog.NewJSONHandler(&bytes.Buffer{}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	router := gin.New()
	router.Use(SlogLogger())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
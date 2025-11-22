package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("记录HTTP请求指标", func(t *testing.T) {
		// 重置注册表以避免冲突
		registry := prometheus.NewRegistry()
		prometheus.DefaultRegisterer = registry
		prometheus.DefaultGatherer = registry

		router := gin.New()
		router.Use(MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("使用空路径记录指标", func(t *testing.T) {
		// 重置注册表
		registry := prometheus.NewRegistry()
		prometheus.DefaultRegisterer = registry
		prometheus.DefaultGatherer = registry

		router := gin.New()
		router.Use(MetricsMiddleware())

		// 创建一个没有设置FullPath的请求
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 验证请求完成且没有panic
		if w.Code == 0 {
			t.Error("Expected a valid HTTP response")
		}
	})
}

// Note: MetricsHandler is not used anymore, metrics endpoint is registered in main.go

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("捕获panic并返回500错误", func(t *testing.T) {
		router := gin.New()
		router.Use(Recovery())
		router.GET("/api/panic", func(c *gin.Context) {
			panic("测试panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/api/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// 检查响应是否包含错误信息
		if w.Body.String() == "" {
			t.Error("Expected non-empty response body")
		}
	})

	t.Run("正常请求不受影响", func(t *testing.T) {
		router := gin.New()
		router.Use(Recovery())
		router.GET("/api/normal", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/normal", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}


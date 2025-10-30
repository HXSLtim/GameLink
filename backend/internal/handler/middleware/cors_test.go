package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 保存原始环境变量
	originalEnv := os.Getenv("APP_ENV")
	originalOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	defer func() {
		os.Setenv("APP_ENV", originalEnv)
		os.Setenv("CORS_ALLOWED_ORIGINS", originalOrigins)
	}()

	t.Run("允许的源-设置CORS头", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		os.Setenv("CORS_ALLOWED_ORIGINS", "")

		router := gin.New()
		router.Use(CORS())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 检查CORS头
		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Error("Expected Access-Control-Allow-Origin header to be set")
		}

		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("Expected Access-Control-Allow-Methods header to be set")
		}

		if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
			t.Error("Expected Access-Control-Allow-Credentials to be true")
		}
	})

	t.Run("处理OPTIONS预检请求", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		os.Setenv("CORS_ALLOWED_ORIGINS", "")

		router := gin.New()
		router.Use(CORS())
		router.OPTIONS("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
		})

		req := httptest.NewRequest(http.MethodOptions, "/api/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// OPTIONS请求应该返回204
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("特定源-匹配时设置CORS头", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		os.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com,https://app.example.com")

		router := gin.New()
		router.Use(CORS())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "https://example.com" {
			t.Errorf("Expected origin https://example.com, got %s", origin)
		}
	})

	t.Run("特定源-不匹配时不设置CORS头", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		os.Setenv("CORS_ALLOWED_ORIGINS", "https://example.com")

		router := gin.New()
		router.Use(CORS())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Origin", "https://attacker.com")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 不应该设置CORS头
		if w.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Error("Should not set CORS headers for unauthorized origin")
		}
	})

	t.Run("没有Origin头-正常处理", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		os.Setenv("CORS_ALLOWED_ORIGINS", "")

		router := gin.New()
		router.Use(CORS())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		// 不设置Origin头
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}


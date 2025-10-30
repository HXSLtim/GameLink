package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 保存原始环境变量
	originalEnv := os.Getenv("APP_ENV")
	originalToken := os.Getenv("ADMIN_TOKEN")
	defer func() {
		os.Setenv("APP_ENV", originalEnv)
		os.Setenv("ADMIN_TOKEN", originalToken)
	}()

	t.Run("生产环境未配置Token-拒绝访问", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		os.Setenv("ADMIN_TOKEN", "")

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}
	})

	t.Run("开发环境无Token-允许访问", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		os.Setenv("ADMIN_TOKEN", "")

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("开发环境带X-Admin-User-ID头", func(t *testing.T) {
		os.Setenv("APP_ENV", "development")
		os.Setenv("ADMIN_TOKEN", "")

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("X-Admin-User-ID", "123")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("配置Token-正确的Token-允许访问", func(t *testing.T) {
		testToken := "test-admin-token-12345"
		os.Setenv("APP_ENV", "production")
		os.Setenv("ADMIN_TOKEN", testToken)

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+testToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("配置Token-错误的Token-拒绝访问", func(t *testing.T) {
		testToken := "test-admin-token-12345"
		os.Setenv("APP_ENV", "production")
		os.Setenv("ADMIN_TOKEN", testToken)

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer wrong-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("配置Token-缺少Authorization头-拒绝访问", func(t *testing.T) {
		testToken := "test-admin-token-12345"
		os.Setenv("APP_ENV", "production")
		os.Setenv("ADMIN_TOKEN", testToken)

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("配置Token-错误的Bearer格式-拒绝访问", func(t *testing.T) {
		testToken := "test-admin-token-12345"
		os.Setenv("APP_ENV", "production")
		os.Setenv("ADMIN_TOKEN", testToken)

		router := gin.New()
		router.Use(AdminAuth())
		router.GET("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("Authorization", testToken) // 缺少 "Bearer " 前缀
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

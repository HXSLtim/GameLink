package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("默认配置应该设置所有安全头", func(t *testing.T) {
		router := gin.New()
		router.Use(SecurityHeaders())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// 验证所有安全头
		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
		assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "max-age=31536000; includeSubDomains", w.Header().Get("Strict-Transport-Security"))
		assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
		assert.NotEmpty(t, w.Header().Get("Permissions-Policy"))
	})

	t.Run("自定义配置应该覆盖默认值", func(t *testing.T) {
		config := SecurityHeadersConfig{
			XFrameOptions:           "SAMEORIGIN",
			XContentTypeOptions:     "nosniff",
			XXSSProtection:          "0",
			ContentSecurityPolicy:   "default-src 'none'",
			StrictTransportSecurity: "max-age=63072000",
			ReferrerPolicy:          "no-referrer",
			PermissionsPolicy:       "geolocation=(self)",
		}

		router := gin.New()
		router.Use(SecurityHeaders(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, "SAMEORIGIN", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "0", w.Header().Get("X-XSS-Protection"))
		assert.Equal(t, "default-src 'none'", w.Header().Get("Content-Security-Policy"))
		assert.Equal(t, "max-age=63072000", w.Header().Get("Strict-Transport-Security"))
		assert.Equal(t, "no-referrer", w.Header().Get("Referrer-Policy"))
		assert.Equal(t, "geolocation=(self)", w.Header().Get("Permissions-Policy"))
	})

	t.Run("空配置值不应该设置对应的header", func(t *testing.T) {
		config := SecurityHeadersConfig{
			XFrameOptions:           "",
			XContentTypeOptions:     "nosniff",
			XXSSProtection:          "",
			ContentSecurityPolicy:   "",
			StrictTransportSecurity: "",
			ReferrerPolicy:          "",
			PermissionsPolicy:       "",
		}

		router := gin.New()
		router.Use(SecurityHeaders(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Empty(t, w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Empty(t, w.Header().Get("X-XSS-Protection"))
		assert.Empty(t, w.Header().Get("Content-Security-Policy"))
		assert.Empty(t, w.Header().Get("Strict-Transport-Security"))
		assert.Empty(t, w.Header().Get("Referrer-Policy"))
		assert.Empty(t, w.Header().Get("Permissions-Policy"))
	})

	t.Run("SecureHeaders快捷方法应该使用默认配置", func(t *testing.T) {
		router := gin.New()
		router.Use(SecureHeaders())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	})

	t.Run("安全头应该在所有请求中设置", func(t *testing.T) {
		router := gin.New()
		router.Use(SecurityHeaders())
		
		router.GET("/get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "get"})
		})
		router.POST("/post", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "post"})
		})
		router.PUT("/put", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "put"})
		})
		router.DELETE("/delete", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "delete"})
		})

		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
		paths := []string{"/get", "/post", "/put", "/delete"}

		for i, method := range methods {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(method, paths[i], nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
			assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		}
	})

	t.Run("CSP配置应该正确设置", func(t *testing.T) {
		config := SecurityHeadersConfig{
			ContentSecurityPolicy: "default-src 'self'; script-src 'self' https://cdn.example.com; style-src 'self' 'unsafe-inline'",
		}

		router := gin.New()
		router.Use(SecurityHeaders(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		csp := w.Header().Get("Content-Security-Policy")
		assert.Contains(t, csp, "default-src 'self'")
		assert.Contains(t, csp, "script-src 'self' https://cdn.example.com")
		assert.Contains(t, csp, "style-src 'self' 'unsafe-inline'")
	})

	t.Run("HSTS配置应该正确设置", func(t *testing.T) {
		config := SecurityHeadersConfig{
			StrictTransportSecurity: "max-age=63072000; includeSubDomains; preload",
		}

		router := gin.New()
		router.Use(SecurityHeaders(config))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		hsts := w.Header().Get("Strict-Transport-Security")
		assert.Equal(t, "max-age=63072000; includeSubDomains; preload", hsts)
	})
}

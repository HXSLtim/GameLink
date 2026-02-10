package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuth_ProductionWithoutToken(t *testing.T) {
	// Set production environment
	os.Setenv("APP_ENV", "production")
	os.Setenv("ADMIN_TOKEN", "") // No token configured
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("ADMIN_TOKEN")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AdminAuth())
	router.GET("/admin/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test request should be rejected with 401 (not 503)
	req := httptest.NewRequest("GET", "/admin/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 Unauthorized when admin auth not configured in production")
	assert.Contains(t, w.Body.String(), "admin authentication required")
}

func TestAdminAuth_ProductionWithValidToken(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("ADMIN_TOKEN", "secure-admin-token-123")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("ADMIN_TOKEN")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AdminAuth())
	router.GET("/admin/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/admin/test", nil)
	req.Header.Set("Authorization", "Bearer secure-admin-token-123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 with valid token")
}

func TestAdminAuth_ProductionWithInvalidToken(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("ADMIN_TOKEN", "secure-admin-token-123")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("ADMIN_TOKEN")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AdminAuth())
	router.GET("/admin/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/admin/test", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 with invalid token")
}

func TestCORS_ProductionDefaultsToDeny(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	defer func() {
		os.Unsetenv("APP_ENV")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://malicious-site.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should NOT have CORS headers
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"), "Production should not allow unknown origins")
}

func TestCORS_ProductionWithExplicitOrigins(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("CORS_ALLOWED_ORIGINS", "https://gamelink.com,https://admin.gamelink.com")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("CORS_ALLOWED_ORIGINS")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		origin        string
		shouldAllow   bool
		expectedAllow string
	}{
		{"https://gamelink.com", true, "https://gamelink.com"},
		{"https://admin.gamelink.com", true, "https://admin.gamelink.com"},
		{"https://malicious-site.com", false, ""},
		{"https://subdomain.gamelink.com", false, ""}, // Not in list
	}

	for _, tt := range tests {
		t.Run(tt.origin, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			allowHeader := w.Header().Get("Access-Control-Allow-Origin")
			if tt.shouldAllow {
				assert.Equal(t, tt.expectedAllow, allowHeader, "Should allow configured origin")
			} else {
				assert.Empty(t, allowHeader, "Should not allow unconfigured origin")
			}
		})
	}
}

func TestCORS_DevelopmentDefaults(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	defer func() {
		os.Unsetenv("APP_ENV")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	tests := []struct {
		origin      string
		shouldAllow bool
	}{
		{"http://localhost:5173", true},
		{"http://localhost:3000", true},
		{"http://127.0.0.1:5173", true},
		{"http://127.0.0.1:3000", true},
		// Development mode uses wildcard ("*") to support LAN access,
		// so all origins are allowed. Production/staging blocks unknown origins.
		{"https://malicious-site.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.origin, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Origin", tt.origin)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			allowHeader := w.Header().Get("Access-Control-Allow-Origin")
			if tt.shouldAllow {
			assert.Equal(t, tt.origin, allowHeader, "Should allow all origins in development (wildcard)")
			} else {
				assert.Empty(t, allowHeader, "Should not allow external sites in development")
			}
		})
	}
}

func TestCORS_StagingDefaultsToDeny(t *testing.T) {
	os.Setenv("APP_ENV", "staging")
	os.Unsetenv("CORS_ALLOWED_ORIGINS")
	defer func() {
		os.Unsetenv("APP_ENV")
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORS())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://malicious-site.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Staging should also default to deny
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"), "Staging should not allow unknown origins")
}

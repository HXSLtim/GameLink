package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPrometheusMiddleware tests the Prometheus metrics middleware
func TestPrometheusMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create router with middleware
	router := gin.New()
	router.Use(PrometheusMiddleware(nil))

	router.GET("/fast", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "fast response"})
	})

	router.GET("/slow", func(c *gin.Context) {
		time.Sleep(1100 * time.Millisecond) // Simulate slow request
		c.JSON(http.StatusOK, gin.H{"message": "slow response"})
	})

	router.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	// Test fast request
	t.Run("fast request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/fast", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test error response
	t.Run("error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Note: We skip the actual slow request test in unit tests to avoid >1s delay
	// In integration tests, you would want to test the slow request logging
}

// TestMetricsAuth tests the metrics authentication middleware
func TestMetricsAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		config         MetricsAuthConfig
		clientIP       string
		expectedStatus int
	}{
		{
			name: "authentication disabled",
			config: MetricsAuthConfig{
				Enabled: false,
			},
			clientIP:       "192.168.1.100",
			expectedStatus: http.StatusOK,
		},
		{
			name: "localhost allowed",
			config: MetricsAuthConfig{
				Enabled:      true,
				AllowedCIDRs: []string{}, // Empty means localhost only
			},
			clientIP:       "127.0.0.1",
			expectedStatus: http.StatusOK,
		},
		{
			name: "remote IP denied",
			config: MetricsAuthConfig{
				Enabled:      true,
				AllowedCIDRs: []string{}, // Empty means localhost only
			},
			clientIP:       "192.168.1.100",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "allowed CIDR",
			config: MetricsAuthConfig{
				Enabled:      true,
				AllowedCIDRs: []string{"10.0.0.0/8", "192.168.0.0/16"},
			},
			clientIP:       "192.168.1.100",
			expectedStatus: http.StatusOK,
		},
		{
			name: "allowed CIDR - 10.x.x.x",
			config: MetricsAuthConfig{
				Enabled:      true,
				AllowedCIDRs: []string{"10.0.0.0/8"},
			},
			clientIP:       "10.0.0.5",
			expectedStatus: http.StatusOK,
		},
		{
			name: "denied CIDR",
			config: MetricsAuthConfig{
				Enabled:      true,
				AllowedCIDRs: []string{"10.0.0.0/8"},
			},
			clientIP:       "172.16.0.1",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				// Simulate client IP
				c.Request = &http.Request{
					RemoteAddr: tc.clientIP + ":12345",
				}
				c.Next()
			})
			router.Use(MetricsAuth(tc.config))

			router.GET("/metrics", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/metrics", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

// TestIsLocalhost tests the localhost detection function
func TestIsLocalhost(t *testing.T) {
	testCases := []struct {
		ip       string
		expected bool
	}{
		{"127.0.0.1", true},
		{"127.0.0.2", true},
		{"127.1.1.1", true},
		{"::1", true},
		{"localhost", true}, // Will fail ParseIP, but we handle it
		{"192.168.1.1", false},
		{"10.0.0.1", false},
		{"8.8.8.8", false},
	}

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			ip := net.ParseIP(tc.ip)
			if ip == nil {
				// Invalid IP (like "localhost")
				if tc.ip == "localhost" {
					// This is expected to fail ParseIP
					return
				}
				t.Fatalf("failed to parse IP: %s", tc.ip)
			}
			result := isLocalhost(ip)
			assert.Equal(t, tc.expected, result, "isLocalhost(%s) should be %v", tc.ip, tc.expected)
		})
	}
}

// TestGetClientIP tests the client IP extraction
func TestGetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name          string
		xForwardedFor string
		xRealIP       string
		remoteAddr    string
		expectedIP    string
	}{
		{
			name:       "direct connection",
			remoteAddr: "192.168.1.100:12345",
			expectedIP: "192.168.1.100",
		},
		{
			name:       "X-Real-IP header",
			xRealIP:    "10.0.0.5",
			remoteAddr: "192.168.1.100:12345",
			expectedIP: "10.0.0.5",
		},
		{
			name:          "X-Forwarded-For header",
			xForwardedFor: "10.0.0.6",
			remoteAddr:    "192.168.1.100:12345",
			expectedIP:    "10.0.0.6",
		},
		{
			name:          "X-Forwarded-For with multiple IPs",
			xForwardedFor: "10.0.0.7, 10.0.0.8, 10.0.0.9",
			remoteAddr:    "192.168.1.100:12345",
			expectedIP:    "10.0.0.7",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			req, _ := http.NewRequest("GET", "/test", nil)
			if tc.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tc.xForwardedFor)
			}
			if tc.xRealIP != "" {
				req.Header.Set("X-Real-IP", tc.xRealIP)
			}
			req.RemoteAddr = tc.remoteAddr

			c.Request = req

			result := getMetricsClientIP(c)
			assert.Equal(t, tc.expectedIP, result)
		})
	}
}

// TestDefaultMetricsAuthConfig tests the default configuration
func TestDefaultMetricsAuthConfig(t *testing.T) {
	config := DefaultMetricsAuthConfig()

	assert.True(t, config.Enabled, "authentication should be enabled by default")
	assert.Empty(t, config.AllowedCIDRs, "no CIDRs should be allowed by default (localhost only)")
	assert.NotEmpty(t, config.TrustedProxies, "trusted proxies should be configured")
}

// BenchmarkPrometheusMiddleware benchmarks the middleware performance
func BenchmarkPrometheusMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(PrometheusMiddleware(nil))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

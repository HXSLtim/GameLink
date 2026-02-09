package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis 设置测试用 Redis 服务器
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	s, err := miniredis.Run()
	require.NoError(t, err, "Failed to start miniredis")

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	// Test connection
	ctx := context.Background()
	err = client.Ping(ctx).Err()
	require.NoError(t, err, "Failed to connect to miniredis")

	return s, client
}

// testCache 实现 cache.Cache 接口用于测试
type testCache struct {
	client *redis.Client
}

func newTestCache(client *redis.Client) *testCache {
	return &testCache{client: client}
}

func (tc *testCache) Get(ctx context.Context, key string) (string, bool, error) {
	result, err := tc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return result, true, nil
}

func (tc *testCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return tc.client.Set(ctx, key, value, ttl).Err()
}

func (tc *testCache) Delete(ctx context.Context, key string) error {
	return tc.client.Del(ctx, key).Err()
}

func (tc *testCache) Close(ctx context.Context) error {
	return tc.client.Close()
}

func (tc *testCache) GetClient() *redis.Client {
	return tc.client
}

func (tc *testCache) GetRedisClient() interface{} {
	return tc.client
}

// TestRedisRateLimiter_NewRedisRateLimiter tests creating a new Redis rate limiter
func TestRedisRateLimiter_NewRedisRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := DefaultRedisRateLimiterConfig()
	config.Enabled = true

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err, "Failed to create rate limiter")
	require.NotNil(t, limiter, "Limiter should not be nil")

	assert.NotNil(t, limiter.loginLimiter, "Login limiter should be initialized")
	assert.NotNil(t, limiter.smsLimiter, "SMS limiter should be initialized")
	assert.NotNil(t, limiter.apiLimiter, "API limiter should be initialized")
	assert.NotNil(t, limiter.uploadLimiter, "Upload limiter should be initialized")
}

// TestRedisRateLimiter_LoginRateLimit tests login rate limiting (5 requests per 15 minutes)
func TestRedisRateLimiter_LoginRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RedisRateLimiterConfig{
		Enabled:       true,
		LoginRequests: 5,
		LoginWindow:   15 * time.Minute,
	}

	setupEngine := func() (*gin.Engine, func()) {
		s, client := setupTestRedis(t)
		cacheClient := newTestCache(client)
		limiter, err := NewRedisRateLimiter(cacheClient, config)
		require.NoError(t, err)

		engine := gin.New()
		engine.POST("/api/v1/auth/login", limiter.LoginRateLimit(), func(c *gin.Context) {
			c.JSON(200, gin.H{"success": true})
		})

		return engine, s.Close
	}

	t.Run("Allow requests within rate limit", func(t *testing.T) {
		engine, cleanup := setupEngine()
		defer cleanup()
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
			req.RemoteAddr = "192.168.1.100:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Request %d should be allowed", i+1)

			// Check rate limit headers
			expectedRemaining := 4 - i
			assert.Equal(t, strconv.Itoa(expectedRemaining), w.Header().Get("X-RateLimit-Remaining"))
		}
	})

	t.Run("Block requests over rate limit", func(t *testing.T) {
		engine, cleanup := setupEngine()
		defer cleanup()
		// Use same IP as first test to verify rate limiting works
		// Make 6 requests (5 allowed + 1 blocked)
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
			req.RemoteAddr = "192.168.1.100:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code, "Request %d should be allowed", i+1)
		}

		// 6th request should be blocked
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "6th request should be rate limited")

		// Check response body
		assert.Contains(t, w.Body.String(), "Rate limit exceeded")
		assert.Contains(t, w.Header().Get("Retry-After"), "900") // 15 minutes = 900 seconds
	})

	t.Run("Different IPs have separate limits", func(t *testing.T) {
		engine, cleanup := setupEngine()
		defer cleanup()
		ips := []string{"10.0.0.1", "10.0.0.2"}

		for _, ip := range ips {
			for i := 0; i < 5; i++ {
				req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
				req.RemoteAddr = ip + ":12345"
				w := httptest.NewRecorder()
				engine.ServeHTTP(w, req)

				assert.Equal(t, 200, w.Code, "IP %s request %d should be allowed", ip, i+1)
			}
		}
	})
}

// TestRedisRateLimiter_SMSRateLimit tests SMS rate limiting (3 requests per hour)
func TestRedisRateLimiter_SMSRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RedisRateLimiterConfig{
		Enabled:     true,
		SMSRequests: 3,
		SMSWindow:   1 * time.Hour,
	}

	setupEngine := func() (*gin.Engine, func()) {
		s, client := setupTestRedis(t)
		cacheClient := newTestCache(client)
		limiter, err := NewRedisRateLimiter(cacheClient, config)
		require.NoError(t, err)

		engine := gin.New()
		engine.POST("/api/v1/auth/sms", limiter.SMSRateLimit(), func(c *gin.Context) {
			c.JSON(200, gin.H{"success": true})
		})

		return engine, s.Close
	}

	t.Run("Allow 3 SMS requests", func(t *testing.T) {
		engine, cleanup := setupEngine()
		defer cleanup()
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/auth/sms", nil)
			req.RemoteAddr = "192.168.1.200:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "SMS request %d should be allowed", i+1)
		}
	})

	t.Run("Block 4th SMS request", func(t *testing.T) {
		engine, cleanup := setupEngine()
		defer cleanup()
		// Make 3 requests first to consume limit
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/auth/sms", nil)
			req.RemoteAddr = "192.168.1.200:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code, "SMS request %d should be allowed", i+1)
		}

		// 4th request should be blocked
		req, _ := http.NewRequest("POST", "/api/v1/auth/sms", nil)
		req.RemoteAddr = "192.168.1.200:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "4th SMS request should be rate limited")
		assert.Contains(t, w.Header().Get("Retry-After"), "3600") // 1 hour = 3600 seconds
	})
}

// TestRedisRateLimiter_APIRateLimit tests API rate limiting (60 requests per minute)
func TestRedisRateLimiter_APIRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := RedisRateLimiterConfig{
		Enabled:     true,
		APIRequests: 60,
		APIWindow:   1 * time.Minute,
	}

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	engine := gin.New()
	engine.GET("/api/v1/test", limiter.APIRateLimit(), func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Allow 60 API requests", func(t *testing.T) {
		for i := 0; i < 60; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/test", nil)
			req.RemoteAddr = "192.168.1.300:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "API request %d should be allowed", i+1)
		}
	})

	t.Run("Block 61st API request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		req.RemoteAddr = "192.168.1.301:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "61st API request should be rate limited")
	})
}

// TestRedisRateLimiter_UploadRateLimit tests upload rate limiting (10 requests per hour)
func TestRedisRateLimiter_UploadRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := RedisRateLimiterConfig{
		Enabled:        true,
		UploadRequests: 10,
		UploadWindow:   1 * time.Hour,
	}

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	engine := gin.New()
	engine.POST("/api/v1/upload", limiter.UploadRateLimit(), func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Allow 10 upload requests", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/upload", nil)
			req.RemoteAddr = "192.168.1.400:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Upload request %d should be allowed", i+1)
		}
	})

	t.Run("Block 11th upload request", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/upload", nil)
		req.RemoteAddr = "192.168.1.401:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "11th upload request should be rate limited")
	})
}

// TestRedisRateLimiter_Disabled tests that rate limiting can be disabled
func TestRedisRateLimiter_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := DefaultRedisRateLimiterConfig()
	config.Enabled = false // Disabled

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	engine := gin.New()
	engine.POST("/api/v1/auth/login", limiter.LoginRateLimit(), func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	// Should be able to make many requests without being limited
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
		req.RemoteAddr = "192.168.1.500:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code, "Request %d should succeed when rate limiting is disabled", i+1)
	}
}

// TestRedisRateLimiter_CustomRateLimit tests custom rate limiting
func TestRedisRateLimiter_CustomRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := DefaultRedisRateLimiterConfig()

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	// Create custom rate limiter: 2 requests per second
	customLimiter := limiter.CustomRateLimit(2, 1*time.Second, "custom")

	engine := gin.New()
	engine.GET("/api/v1/custom", customLimiter, func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Allow 2 custom requests", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/custom", nil)
			req.RemoteAddr = "192.168.1.600:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Custom request %d should be allowed", i+1)
		}
	})

	t.Run("Block 3rd custom request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/custom", nil)
		req.RemoteAddr = "192.168.1.601:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "3rd custom request should be rate limited")
	})

	t.Run("Allow requests after window expires", func(t *testing.T) {
		// Fast forward time by 1 second
		s.FastForward(1 * time.Second)

		req, _ := http.NewRequest("GET", "/api/v1/custom", nil)
		req.RemoteAddr = "192.168.1.602:12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code, "Request should be allowed after window expires")
	})
}

// TestRedisRateLimiter_UserBasedRateLimit tests user-based rate limiting
func TestRedisRateLimiter_UserBasedRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := RedisRateLimiterConfig{
		Enabled:       true,
		LoginRequests: 3,
		LoginWindow:   1 * time.Minute,
	}

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	engine := gin.New()
	// Middleware to set user_id
	setUserID := func(c *gin.Context) {
		c.Set("user_id", uint64(12345))
		c.Next()
	}

	engine.POST("/api/v1/auth/login", setUserID, limiter.LoginRateLimit(), func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("User-based rate limiting", func(t *testing.T) {
		// Make 3 requests with same user_id
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "User request %d should be allowed", i+1)
		}

		// 4th request should be blocked
		req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "4th user request should be rate limited")
	})
}

// TestRedisRateLimiter_GetClientIP tests client IP extraction
func TestRedisRateLimiter_GetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name          string
		xForwardedFor string
		xRealIP       string
		remoteAddr    string
		expectedIP    string
	}{
		{
			name:          "X-Forwarded-For header with multiple IPs",
			xForwardedFor: "203.0.113.1, 198.51.100.1",
			expectedIP:    "203.0.113.1",
		},
		{
			name:       "X-Real-IP header",
			xRealIP:    "198.51.100.50",
			expectedIP: "198.51.100.50",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
		{
			name:          "X-Forwarded-For priority over X-Real-IP",
			xForwardedFor: "10.0.0.1",
			xRealIP:       "10.0.0.2",
			expectedIP:    "10.0.0.1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			// Build headers
			headers := http.Header{}
			if tc.xForwardedFor != "" {
				headers.Set("X-Forwarded-For", tc.xForwardedFor)
			}
			if tc.xRealIP != "" {
				headers.Set("X-Real-IP", tc.xRealIP)
			}

			c.Request = &http.Request{
				RemoteAddr: tc.remoteAddr,
				Header:     headers,
			}

			ip := GetClientIP(c)
			assert.Equal(t, tc.expectedIP, ip)
		})
	}
}

// TestRedisRateLimiter_ResponseHeaders tests rate limit response headers
func TestRedisRateLimiter_ResponseHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s, client := setupTestRedis(t)
	defer s.Close()

	cacheClient := newTestCache(client)
	config := RedisRateLimiterConfig{
		Enabled:       true,
		LoginRequests: 10,
		LoginWindow:   1 * time.Minute,
	}

	limiter, err := NewRedisRateLimiter(cacheClient, config)
	require.NoError(t, err)

	engine := gin.New()
	engine.POST("/api/v1/auth/login", limiter.LoginRateLimit(), func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
	req.RemoteAddr = "192.168.1.700:12345"
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check response headers
	assert.Equal(t, "10", w.Header().Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
}

// TestRedisRateLimiter_DefaultConfig tests default configuration
func TestRedisRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRedisRateLimiterConfig()

	assert.True(t, config.Enabled, "Rate limiting should be enabled by default")
	assert.Equal(t, 5, config.LoginRequests, "Default login requests should be 5")
	assert.Equal(t, 15*time.Minute, config.LoginWindow, "Default login window should be 15 minutes")
	assert.Equal(t, 3, config.SMSRequests, "Default SMS requests should be 3")
	assert.Equal(t, 1*time.Hour, config.SMSWindow, "Default SMS window should be 1 hour")
	assert.Equal(t, 60, config.APIRequests, "Default API requests should be 60")
	assert.Equal(t, 1*time.Minute, config.APIWindow, "Default API window should be 1 minute")
	assert.Equal(t, 10, config.UploadRequests, "Default upload requests should be 10")
	assert.Equal(t, 1*time.Hour, config.UploadWindow, "Default upload window should be 1 hour")
}

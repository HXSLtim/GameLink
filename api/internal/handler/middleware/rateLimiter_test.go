package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRateLimiter_TokenBucket tests the TokenBucket algorithm
func TestRateLimiter_TokenBucket(t *testing.T) {
	t.Run("Allow request within rate", func(t *testing.T) {
		tb := NewTokenBucket(10, 10) // 10 tokens/sec, capacity 10
		allowed, remaining := tb.Allow()
		assert.True(t, allowed)
		assert.Equal(t, 9.0, remaining)
	})

	t.Run("Consume all tokens", func(t *testing.T) {
		tb := NewTokenBucket(5, 5) // 5 tokens capacity
		for i := 0; i < 5; i++ {
			allowed, _ := tb.Allow()
			assert.True(t, allowed, "Request %d should be allowed", i+1)
		}
		// 6th request should be denied
		allowed, remaining := tb.Allow()
		assert.False(t, allowed)
		assert.Equal(t, 0.0, remaining)
	})

	t.Run("Token refill over time", func(t *testing.T) {
		tb := NewTokenBucket(10, 10)
		// Consume all tokens
		for i := 0; i < 10; i++ {
			tb.Allow()
		}
		// Should be empty now
		allowed, _ := tb.Allow()
		assert.False(t, allowed)

		// Wait for refill (100ms should give 1 token at 10/sec rate)
		time.Sleep(150 * time.Millisecond)
		allowed, remaining := tb.Allow()
		assert.True(t, allowed)
		assert.Greater(t, remaining, 0.0)
	})

	t.Run("Burst capacity", func(t *testing.T) {
		tb := NewTokenBucket(1, 10) // 1 token/sec but 10 capacity
		// Should be able to make 10 requests immediately (burst)
		for i := 0; i < 10; i++ {
			allowed, _ := tb.Allow()
			assert.True(t, allowed, "Burst request %d should be allowed", i+1)
		}
		// 11th should fail
		allowed, _ := tb.Allow()
		assert.False(t, allowed)
	})
}

// TestRateLimiter_IPRateLimit tests IP-based rate limiting
func TestRateLimiter_IPRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:             true,
		IPRequestsPerSecond: 2, // Very low rate for testing
		RouteLimits:         map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		// Simulate rate limiter middleware for IP
		clientIP := c.GetHeader("X-Real-IP")
		tb := limiter.getIPLimiter(clientIP)
		if allowed, remaining := tb.Allow(); !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":   false,
				"code":      http.StatusTooManyRequests,
				"message":   "IP rate limit exceeded",
				"remaining": remaining,
			})
			c.Abort()
			return
		}
		c.Next()
	})
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Allow requests within rate limit", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Real-IP", "192.168.1.100")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
		}
	})

	t.Run("Block requests over rate limit", func(t *testing.T) {
		// Make 5 requests (burst capacity is 4, so 5th should be blocked)
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Real-IP", "192.168.1.101")
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			if i < 4 {
				assert.Equal(t, 200, w.Code, "Request %d should be allowed", i+1)
			} else {
				assert.Equal(t, 429, w.Code, "Request %d should be rate limited", i+1)
			}
		}
	})

	t.Run("Different IPs have separate limits", func(t *testing.T) {
		ips := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}

		for _, ip := range ips {
			// Each IP should be able to make 2 requests
			for i := 0; i < 2; i++ {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("X-Real-IP", ip)
				w := httptest.NewRecorder()
				engine.ServeHTTP(w, req)

				assert.Equal(t, 200, w.Code, "IP %s request %d should be allowed", ip, i+1)
			}
		}
	})
}

// TestRateLimiter_UserRateLimit tests user-based rate limiting
func TestRateLimiter_UserRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:               true,
		UserRequestsPerMinute: 10, // Low rate for testing
		RouteLimits:           map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		// Simulate user authentication
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr != "" {
			c.Set("user_id", uint64(1000)) // Simulate authenticated user
		}

		// Check user rate limit
		if userID, exists := getUserIDFromContext(c); exists {
			tb := limiter.getUserLimiter(userID)
			if allowed, remaining := tb.Allow(); !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success":   false,
					"code":      http.StatusTooManyRequests,
					"message":   "User rate limit exceeded",
					"remaining": remaining,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	})
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Authenticated user rate limiting", func(t *testing.T) {
		userID := uint64(12345)

		// Create 5 requests (should be allowed as rate is 10/min)
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("X-User-ID", "test")
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = req
			c.Set("user_id", userID)

			tb := limiter.getUserLimiter(userID)
			allowed, remaining := tb.Allow()

			assert.True(t, allowed, "Request %d should be allowed", i+1)
			assert.LessOrEqual(t, remaining, 10.0)
		}
	})

	t.Run("Different users have separate limits", func(t *testing.T) {
		userIDs := []uint64{1001, 1002, 1003}

		for _, uid := range userIDs {
			tb := limiter.getUserLimiter(uid)
			// Each user should get their own limiter
			for i := 0; i < 5; i++ {
				allowed, _ := tb.Allow()
				assert.True(t, allowed, "User %d request %d should be allowed", uid, i+1)
			}
		}
	})
}

// TestRateLimiter_RouteRateLimit tests route-specific rate limiting
func TestRateLimiter_RouteRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled: true,
		RouteLimits: map[string]RouteLimit{
			"/api/login": {
				Path:      "/api/login",
				Requests:  3,
				Window:    time.Minute,
				LimitByIP: true,
			},
		},
	}

	limiter := NewRateLimiter(config)

	engine := gin.New()
	engine.POST("/api/login", func(c *gin.Context) {
		path := c.FullPath()
		clientIP := c.ClientIP()

		// Check route limit
		if routeLimit, exists := config.RouteLimits[path]; exists {
			key := path + ":ip:" + clientIP
			tb := limiter.getRouteLimiter(key, routeLimit)
			if allowed, remaining := tb.Allow(); !allowed {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"success":   false,
					"message":   "Login rate limit exceeded",
					"remaining": remaining,
				})
				c.Abort()
				return
			}
		}
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Route-specific rate limiting", func(t *testing.T) {
		clientIP := "192.168.1.50"

		// Make 3 requests (should be allowed)
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("POST", "/api/login", nil)
			req.RemoteAddr = clientIP + ":12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Request %d should be allowed", i+1)
		}

		// 4th request should be blocked
		req, _ := http.NewRequest("POST", "/api/login", nil)
		req.RemoteAddr = clientIP + ":12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
	})
}

// TestRateLimiter_Whitelist tests IP and role whitelisting
func TestRateLimiter_Whitelist(t *testing.T) {
	config := RateLimiterConfig{
		Enabled:             true,
		IPRequestsPerSecond: 1, // Very restrictive
		WhitelistIPs:        []string{"192.168.1.100", "10.0.0.0/8"},
		WhitelistRoles:      []string{"admin", "superAdmin"},
	}

	limiter := NewRateLimiter(config)

	t.Run("IP whitelisting", func(t *testing.T) {
		testCases := []struct {
			ip       string
			expected bool
		}{
			{"192.168.1.100", true}, // Exact match
			{"192.168.1.1", false},  // Not whitelisted
			{"8.8.8.8", false},      // Not whitelisted
		}

		for _, tc := range testCases {
			result := limiter.isWhitelistedIP(tc.ip)
			assert.Equal(t, tc.expected, result, "IP %s whitelist check failed", tc.ip)
		}
	})

	t.Run("Role whitelisting", func(t *testing.T) {
		testCases := []struct {
			role     string
			expected bool
		}{
			{"admin", true},
			{"superAdmin", true},
			{"user", false},
			{"player", false},
		}

		for _, tc := range testCases {
			result := limiter.isWhitelistedRole(tc.role)
			assert.Equal(t, tc.expected, result, "Role %s whitelist check failed", tc.role)
		}
	})
}

// TestRateLimiter_Disabled tests that rate limiting can be disabled
func TestRateLimiter_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:             false, // Disabled
		IPRequestsPerSecond: 1,
	}

	engine := gin.New()
	engine.Use(RateLimit(config))
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	// Should be able to make many requests without being limited
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code, "Request %d should succeed when rate limiting is disabled", i+1)
	}
}

// TestRateLimiter_ConcurrentAccess tests thread safety
func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	tb := NewTokenBucket(100, 100) // High capacity for concurrent test

	var wg sync.WaitGroup
	numGoroutines := 10
	requestsPerGoroutine := 20

	// Track results
	results := make(chan bool, numGoroutines*requestsPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				allowed, _ := tb.Allow()
				results <- allowed
			}
		}()
	}

	wg.Wait()
	close(results)

	// Count allowed requests
	allowedCount := 0
	for allowed := range results {
		if allowed {
			allowedCount++
		}
	}

	// Should have exactly 100 allowed requests (bucket capacity)
	assert.Equal(t, 100, allowedCount, "Exactly 100 requests should be allowed")
}

// TestRateLimiter_ClientIP tests client IP extraction
func TestRateLimiter_ClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name          string
		xForwardedFor string
		xRealIP       string
		remoteAddr    string
		expectedIP    string
	}{
		{
			name:          "X-Forwarded-For header",
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

			// Build headers conditionally
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

			ip := getClientIP(c)
			assert.Equal(t, tc.expectedIP, ip)
		})
	}
}

// TestRateLimiter_DefaultConfig tests default configuration
func TestRateLimiter_DefaultConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.True(t, config.Enabled, "Rate limiting should be enabled by default")
	assert.Equal(t, 50.0, config.IPRequestsPerSecond)
	assert.Equal(t, 300, config.UserRequestsPerMinute)
	assert.NotEmpty(t, config.WhitelistIPs, "Should have whitelisted IPs")
	assert.NotEmpty(t, config.WhitelistRoles, "Should have whitelisted roles")
	assert.NotEmpty(t, config.RouteLimits, "Should have route-specific limits")

	// Check specific route limits
	loginLimit, exists := config.RouteLimits["/api/v1/auth/login"]
	require.True(t, exists, "Login route should have rate limit")
	assert.Equal(t, 10, loginLimit.Requests)
	assert.Equal(t, time.Minute, loginLimit.Window)
	assert.True(t, loginLimit.LimitByIP)

	registerLimit, exists := config.RouteLimits["/api/v1/auth/register"]
	require.True(t, exists, "Register route should have rate limit")
	assert.Equal(t, 20, registerLimit.Requests)
	assert.Equal(t, time.Hour, registerLimit.Window)
}

// TestRateLimiter_NewRateLimiter tests RateLimiter initialization
func TestRateLimiter_NewRateLimiter(t *testing.T) {
	t.Run("Default values for zero config", func(t *testing.T) {
		config := RateLimiterConfig{
			Enabled: true,
			// IPRequestsPerSecond and UserRequestsPerMinute are zero
		}

		limiter := NewRateLimiter(config)

		assert.NotNil(t, limiter)
		assert.NotNil(t, limiter.ipLimiters)
		assert.NotNil(t, limiter.userLimiters)
		assert.NotNil(t, limiter.routeLimiters)
		// Defaults should be applied
		assert.Equal(t, 10.0, limiter.config.IPRequestsPerSecond)
		assert.Equal(t, 60, limiter.config.UserRequestsPerMinute)
	})

	t.Run("Custom values", func(t *testing.T) {
		config := RateLimiterConfig{
			Enabled:               true,
			IPRequestsPerSecond:   100,
			UserRequestsPerMinute: 500,
			RouteLimits: map[string]RouteLimit{
				"/test": {Path: "/test", Requests: 100, Window: time.Minute},
			},
		}

		limiter := NewRateLimiter(config)

		assert.Equal(t, 100.0, limiter.config.IPRequestsPerSecond)
		assert.Equal(t, 500, limiter.config.UserRequestsPerMinute)
		assert.Len(t, limiter.config.RouteLimits, 1)
	})
}

// TestRateLimiter_Integration tests the full RateLimit middleware
func TestRateLimiter_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   5, // 5 requests per second
		UserRequestsPerMinute: 30,
		WhitelistIPs:          []string{"127.0.0.1"},
		RouteLimits: map[string]RouteLimit{
			"/api/v1/test": {
				Path:     "/api/v1/test",
				Requests: 3,
				Window:   time.Second,
			},
		},
	}

	engine := gin.New()
	engine.Use(RateLimit(config))
	engine.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})

	t.Run("Whitelisted IP bypasses rate limit", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/test", nil)
			req.RemoteAddr = "127.0.0.1:12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Whitelisted IP request %d should succeed", i+1)
		}
	})

	t.Run("Route-specific rate limiting", func(t *testing.T) {
		clientIP := "192.168.1.200"

		// First 3 requests should succeed
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", "/api/v1/test", nil)
			req.RemoteAddr = clientIP + ":12345"
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Request %d should succeed", i+1)
		}

		// 4th request should be rate limited
		req, _ := http.NewRequest("GET", "/api/v1/test", nil)
		req.RemoteAddr = clientIP + ":12345"
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code, "Request should be rate limited")

		// Check response headers
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}

// TestRateLimiter_Recovery tests recovery after rate limit window
func TestRateLimiter_Recovery(t *testing.T) {
	tb := NewTokenBucket(2, 2) // 2 tokens/sec, capacity 2

	// Consume all tokens
	for i := 0; i < 2; i++ {
		allowed, _ := tb.Allow()
		assert.True(t, allowed)
	}

	// Should be empty
	allowed, _ := tb.Allow()
	assert.False(t, allowed)

	// Wait for refill
	time.Sleep(600 * time.Millisecond)

	// Should have some tokens back
	allowed, remaining := tb.Allow()
	assert.True(t, allowed, "Should have refilled tokens")
	assert.Greater(t, remaining, 0.0)
}

// TestRateLimiter_getUserIDFromContext tests user ID extraction
func TestRateLimiter_getUserIDFromContext(t *testing.T) {
	t.Run("User ID exists", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint64(12345))
		userID, exists := getUserIDFromContext(c)
		assert.True(t, exists)
		assert.Equal(t, uint64(12345), userID)
	})

	t.Run("User ID does not exist", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		userID, exists := getUserIDFromContext(c)
		assert.False(t, exists)
		assert.Equal(t, uint64(0), userID)
	})

	t.Run("User ID wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", "not-a-uint64")
		userID, exists := getUserIDFromContext(c)
		assert.False(t, exists)
		assert.Equal(t, uint64(0), userID)
	})
}

// TestRateLimiter_getUserRoleFromContext tests user role extraction
func TestRateLimiter_getUserRoleFromContext(t *testing.T) {
	t.Run("Role exists", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "admin")
		role, exists := getUserRoleFromContext(c)
		assert.True(t, exists)
		assert.Equal(t, "admin", role)
	})

	t.Run("Role does not exist", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		role, exists := getUserRoleFromContext(c)
		assert.False(t, exists)
		assert.Equal(t, "", role)
	})

	t.Run("Role wrong type", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", 123)
		role, exists := getUserRoleFromContext(c)
		assert.False(t, exists)
		assert.Equal(t, "", role)
	})
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gamelink/internal/apierr"
)

func TestTokenBucket_Allow(t *testing.T) {
	// 创建令牌桶：每秒1个令牌，容量为2
	tb := NewTokenBucket(1.0, 2.0)

	// 第一次请求应该成功（剩余1个令牌）
	allowed, remaining := tb.Allow()
	assert.True(t, allowed)
	assert.Equal(t, 1.0, remaining)

	// 第二次请求应该成功（剩余0个令牌）
	allowed, remaining = tb.Allow()
	assert.True(t, allowed)
	assert.InDelta(t, 0.0, remaining, 0.01) // 允许很小的浮点误差

	// 第三次请求应该失败（没有令牌了）
	allowed, remaining = tb.Allow()
	assert.False(t, allowed)
	assert.InDelta(t, 0.0, remaining, 0.01) // 允许很小的浮点误差

	// 等待1.5秒，应该填充1.5个令牌
	time.Sleep(1500 * time.Millisecond)
	allowed, remaining = tb.Allow()
	assert.True(t, allowed)
	assert.InDelta(t, 0.5, remaining, 0.1) // 1.5 - 1 = 0.5
}

func TestRateLimiter_IPLimit(t *testing.T) {
	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   10, // 每秒10个请求
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{},
		WhitelistRoles:        []string{},
		RouteLimits:           map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	// 获取IP限流器
	ipLimiter := limiter.getIPLimiter("192.168.1.1")

	// 容量是20（rate * 2），消耗19个令牌（剩余1个）
	for i := 0; i < 19; i++ {
		allowed, _ := ipLimiter.Allow()
		assert.True(t, allowed)
	}

	// 第20个请求应该成功（还有1个令牌）
	allowed, _ := ipLimiter.Allow()
	assert.True(t, allowed)

	// 第21个请求应该失败（令牌耗尽）
	allowed, _ = ipLimiter.Allow()
	assert.False(t, allowed)

	// 等待0.3秒，应该填充3个令牌（10 * 0.3 = 3）
	time.Sleep(300 * time.Millisecond)

	// 现在应该允许3个请求
	for i := 0; i < 3; i++ {
		allowed, _ := ipLimiter.Allow()
		assert.True(t, allowed)
	}

	// 下一个请求应该失败
	allowed, _ = ipLimiter.Allow()
	assert.False(t, allowed)
}

func TestRateLimiter_UserLimit(t *testing.T) {
	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   10,
		UserRequestsPerMinute: 3, // 每分钟3个请求
		WhitelistIPs:          []string{},
		WhitelistRoles:        []string{},
		RouteLimits:           map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	// 获取用户限流器
	userLimiter := limiter.getUserLimiter(uint64(123))

	// 前3个请求应该成功
	for i := 0; i < 3; i++ {
		allowed, _ := userLimiter.Allow()
		assert.True(t, allowed)
	}

	// 第4个请求应该失败
	allowed, _ := userLimiter.Allow()
	assert.False(t, allowed)
}

func TestRateLimiter_WhitelistIP(t *testing.T) {
	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   1,
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{"192.168.1.100", "10.0.0."},
		WhitelistRoles:        []string{},
		RouteLimits:           map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	// 白名单IP应该返回true
	assert.True(t, limiter.isWhitelistedIP("192.168.1.100"))
	assert.True(t, limiter.isWhitelistedIP("10.0.0.5"))
	assert.False(t, limiter.isWhitelistedIP("192.168.1.101"))
}

func TestRateLimiter_WhitelistRole(t *testing.T) {
	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   10,
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{},
		WhitelistRoles:        []string{"super_admin", "admin"},
		RouteLimits:           map[string]RouteLimit{},
	}

	limiter := NewRateLimiter(config)

	assert.True(t, limiter.isWhitelistedRole("super_admin"))
	assert.True(t, limiter.isWhitelistedRole("admin"))
	assert.False(t, limiter.isWhitelistedRole("user"))
}

func TestRateLimit_Middleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   10,
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{},
		WhitelistRoles:        []string{},
		RouteLimits: map[string]RouteLimit{
			"/test/login": {
				Path:        "/test/login",
				Requests:    2,
				Window:      time.Minute,
				LimitByIP:   true,
				LimitByUser: false,
			},
		},
	}

	router := gin.New()
	router.Use(RateLimit(config))

	router.POST("/test/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 第1个请求应该成功
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/test/login", nil)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 第2个请求应该成功
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/test/login", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 第3个请求应该被限流（429）
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/test/login", nil)
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)

	// 检查响应头
	assert.Equal(t, "2", w3.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", w3.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, w3.Header().Get("X-RateLimit-Reset"))
}

func TestRateLimit_WhitelistIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   1,
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{"192.168.1.1"},
		WhitelistRoles:        []string{},
		RouteLimits: map[string]RouteLimit{
			"/test/whitelisted": {
				Path:        "/test/whitelisted",
				Requests:    1,
				Window:      time.Minute,
				LimitByIP:   true,
				LimitByUser: false,
			},
		},
	}

	router := gin.New()
	router.Use(RateLimit(config))

	router.GET("/test/whitelisted", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 白名单IP应该不受限流影响
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/whitelisted", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "白名单IP应该不受限流限制")
	}
}

func TestRateLimit_Disabled(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := RateLimiterConfig{
		Enabled:               false, // 禁用限流
		IPRequestsPerSecond:   1,
		UserRequestsPerMinute: 60,
		WhitelistIPs:          []string{},
		WhitelistRoles:        []string{},
		RouteLimits:           map[string]RouteLimit{},
	}

	router := gin.New()
	router.Use(RateLimit(config))

	router.GET("/test/disabled", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 发送多个请求，应该都成功（限流被禁用）
	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test/disabled", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "限流被禁用，所有请求应该成功")
	}
}

func TestGetClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		headers       map[string]string
		remoteAddr    string
		expectedIP    string
	}{
		{
			name:       "X-Forwarded-For",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.195, 70.41.3.18, 150.172.238.178"},
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "203.0.113.195",
		},
		{
			name:       "X-Real-IP",
			headers:    map[string]string{"X-Real-IP": "203.0.113.195"},
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "203.0.113.195",
		},
		{
			name:       "RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.1:12345",
			expectedIP: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", "/", nil)
			c.Request.RemoteAddr = tt.remoteAddr
			
			for k, v := range tt.headers {
				c.Request.Header.Set(k, v)
			}
			
			ip := getClientIP(c)
			assert.Equal(t, tt.expectedIP, ip)
		})
	}
}

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.True(t, config.Enabled)
	assert.Equal(t, 10.0, config.IPRequestsPerSecond)
	assert.Equal(t, 60, config.UserRequestsPerMinute)
	assert.Contains(t, config.WhitelistIPs, "127.0.0.1")
	assert.Contains(t, config.WhitelistRoles, "super_admin")
	assert.NotEmpty(t, config.RouteLimits)

	// 检查默认路由限流
	loginLimit, exists := config.RouteLimits["/api/v1/auth/login"]
	assert.True(t, exists)
	assert.Equal(t, 5, loginLimit.Requests)
	assert.Equal(t, time.Minute, loginLimit.Window)
	assert.True(t, loginLimit.LimitByIP)
	assert.False(t, loginLimit.LimitByUser)
}

func TestTooManyRequestsError(t *testing.T) {
	// 测试429错误
	err := apierr.TooManyRequests("Rate limit exceeded")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusTooManyRequests, err.Code)
	assert.Equal(t, "Rate limit exceeded", err.Message)
}
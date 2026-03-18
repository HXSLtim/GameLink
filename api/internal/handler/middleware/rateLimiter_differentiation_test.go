package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRateLimitConfig(t *testing.T) {
	config := DefaultRateLimitConfig()

	assert.Equal(t, 20.0, config.IPRequestsPerSecond, "Default IP rate should be 20 RPS")
	assert.Equal(t, 300, config.UserRequestsPerMinute, "Default user rate should be 300 RPM")
	assert.Contains(t, config.WhitelistIPs, "127.0.0.1")
	assert.Contains(t, config.WhitelistRoles, "superAdmin")

	// 验证路由限流规则
	loginLimit, exists := config.RouteLimits["/api/v1/auth/login"]
	assert.True(t, exists, "Login route limit should exist")
	assert.Equal(t, 10, loginLimit.Requests)
	assert.Equal(t, time.Minute, loginLimit.Window)
	assert.True(t, loginLimit.LimitByIP)

	registerLimit, exists := config.RouteLimits["/api/v1/auth/register"]
	assert.True(t, exists, "Register route limit should exist")
	assert.Equal(t, 20, registerLimit.Requests)
	assert.Equal(t, time.Hour, registerLimit.Window)

	refreshLimit, exists := config.RouteLimits["/api/v1/auth/refresh"]
	assert.True(t, exists, "Refresh route limit should exist")
	assert.Equal(t, 600, refreshLimit.Requests)
	assert.True(t, refreshLimit.LimitByUser)

	payLimit, exists := config.RouteLimits["/api/v1/user/orders/:id/pay"]
	assert.True(t, exists, "Payment route limit should exist")
	assert.Equal(t, 30, payLimit.Requests)
	assert.True(t, payLimit.LimitByUser)

	withdrawLimit, exists := config.RouteLimits["/api/v1/user/wallet/withdraw"]
	assert.True(t, exists, "Withdraw route limit should exist")
	assert.Equal(t, 10, withdrawLimit.Requests)
	assert.True(t, withdrawLimit.LimitByUser)
}

func TestAdminRateLimitConfig(t *testing.T) {
	config := AdminRateLimitConfig()

	assert.Equal(t, 5.0, config.IPRequestsPerSecond, "Admin rate should be 5 RPS")
	assert.Equal(t, 300, config.UserRequestsPerMinute, "Admin user rate should be 300 RPM")
	assert.Contains(t, config.WhitelistRoles, "superAdmin")
}

func TestPublicRateLimitConfig(t *testing.T) {
	config := PublicRateLimitConfig()

	assert.Equal(t, 20.0, config.IPRequestsPerSecond, "Public rate should be 20 RPS")
	assert.Equal(t, 1200, config.UserRequestsPerMinute, "Public user rate should be 1200 RPM")
}

func TestAuthRateLimitConfig(t *testing.T) {
	config := AuthRateLimitConfig()

	assert.Equal(t, 10.0, config.IPRequestsPerSecond, "Auth rate should be 10 RPS")
	assert.Equal(t, 600, config.UserRequestsPerMinute, "Auth user rate should be 600 RPM")

	// 验证认证路由限流
	loginLimit, exists := config.RouteLimits["/api/v1/auth/login"]
	assert.True(t, exists)
	assert.Equal(t, 10, loginLimit.Requests)
	assert.True(t, loginLimit.LimitByIP)

	registerLimit, exists := config.RouteLimits["/api/v1/auth/register"]
	assert.True(t, exists)
	assert.Equal(t, 20, registerLimit.Requests)
	assert.Equal(t, time.Hour, registerLimit.Window)
}

func TestTokenBucket(t *testing.T) {
	// 测试令牌桶基本功能
	tb := NewTokenBucket(10, 20) // 10 tokens/sec, capacity 20

	// 初始应该有满桶令牌
	allowed, remaining := tb.Allow()
	assert.True(t, allowed)
	assert.Equal(t, 19.0, remaining)

	// 连续消费令牌
	for i := 0; i < 19; i++ {
		allowed, _ = tb.Allow()
		assert.True(t, allowed)
	}

	// 桶空了，应该被限流
	allowed, remaining = tb.Allow()
	assert.False(t, allowed)
	assert.Equal(t, 0.0, remaining)

	// 等待令牌补充
	time.Sleep(200 * time.Millisecond) // 应该补充 2 个令牌
	allowed, _ = tb.Allow()
	assert.True(t, allowed)
}

func TestRateLimiterDifferentiation(t *testing.T) {
	// 验证差异化限流策略
	tests := []struct {
		name     string
		config   RateLimiterConfig
		expected float64
	}{
		{
			name:     "Admin - 5 RPS",
			config:   AdminRateLimitConfig(),
			expected: 5.0,
		},
		{
			name:     "Public - 20 RPS",
			config:   PublicRateLimitConfig(),
			expected: 20.0,
		},
		{
			name:     "Auth - 10 RPS",
			config:   AuthRateLimitConfig(),
			expected: 10.0,
		},
		{
			name:     "Default - 20 RPS",
			config:   DefaultRateLimitConfig(),
			expected: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.config.IPRequestsPerSecond)
		})
	}
}

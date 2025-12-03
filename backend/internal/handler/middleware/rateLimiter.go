package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"gamelink/pkg/apierr"
	"github.com/gin-gonic/gin"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// IP限流：每秒请求数
	IPRequestsPerSecond float64
	// 用户限流：每分钟请求数
	UserRequestsPerMinute int
	// 路由限流：特定路由的限流规则
	RouteLimits map[string]RouteLimit
	// 白名单IP
	WhitelistIPs []string
	// 白名单用户（角色）
	WhitelistRoles []string
	// 是否启用
	Enabled bool
}

// RouteLimit 路由限流配置
type RouteLimit struct {
	Path        string        // 路由路径
	Requests    int           // 请求数限制
	Window      time.Duration // 时间窗口
	LimitByIP   bool          // 是否按IP限流
	LimitByUser bool          // 是否按用户限流
}

// TokenBucket 令牌桶算法实现
type TokenBucket struct {
	rate       float64   // 每秒填充速率
	capacity   float64   // 桶容量
	tokens     float64   // 当前令牌数
	lastUpdate time.Time // 上次更新时间
	mu         sync.Mutex
}

// NewTokenBucket 创建新的令牌桶
func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastUpdate: time.Now(),
	}
}

// Allow 检查是否允许请求，返回是否允许和剩余令牌数
func (tb *TokenBucket) Allow() (bool, float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tb.lastUpdate = now

	// 填充令牌
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	// 消费令牌
	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true, tb.tokens
	}

	return false, tb.tokens
}

// RateLimiter 限流器
type RateLimiter struct {
	config RateLimiterConfig

	// IP限流器：map[ip]*TokenBucket
	ipLimiters map[string]*TokenBucket
	ipMutex    sync.RWMutex

	// 用户限流器：map[user_id]*TokenBucket
	userLimiters map[uint64]*TokenBucket
	userMutex    sync.RWMutex

	// 路由限流器：map[route_key]*TokenBucket
	routeLimiters map[string]*TokenBucket
	routeMutex    sync.RWMutex
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.IPRequestsPerSecond <= 0 {
		config.IPRequestsPerSecond = 10 // 默认每秒10个请求
	}
	if config.UserRequestsPerMinute <= 0 {
		config.UserRequestsPerMinute = 60 // 默认每分钟60个请求
	}
	if config.RouteLimits == nil {
		config.RouteLimits = make(map[string]RouteLimit)
	}

	return &RateLimiter{
		config:        config,
		ipLimiters:    make(map[string]*TokenBucket),
		userLimiters:  make(map[uint64]*TokenBucket),
		routeLimiters: make(map[string]*TokenBucket),
	}
}

// getIPLimiter 获取IP限流器
func (rl *RateLimiter) getIPLimiter(ip string) *TokenBucket {
	rl.ipMutex.Lock()
	defer rl.ipMutex.Unlock()

	if limiter, exists := rl.ipLimiters[ip]; exists {
		return limiter
	}

	// 创建新的限流器：默认容量为速率的2倍，允许突发
	capacity := rl.config.IPRequestsPerSecond * 2
	limiter := NewTokenBucket(rl.config.IPRequestsPerSecond, capacity)
	rl.ipLimiters[ip] = limiter
	return limiter
}

// getUserLimiter 获取用户限流器
func (rl *RateLimiter) getUserLimiter(userID uint64) *TokenBucket {
	rl.userMutex.Lock()
	defer rl.userMutex.Unlock()

	if limiter, exists := rl.userLimiters[userID]; exists {
		return limiter
	}

	// 创建新的限流器：每分钟限制转换为每秒速率
	rate := float64(rl.config.UserRequestsPerMinute) / 60.0
	capacity := float64(rl.config.UserRequestsPerMinute)
	limiter := NewTokenBucket(rate, capacity)
	rl.userLimiters[userID] = limiter
	return limiter
}

// getRouteLimiter 获取路由限流器
func (rl *RateLimiter) getRouteLimiter(key string, limit RouteLimit) *TokenBucket {
	rl.routeMutex.Lock()
	defer rl.routeMutex.Unlock()

	if limiter, exists := rl.routeLimiters[key]; exists {
		return limiter
	}

	// 创建新的限流器
	rate := float64(limit.Requests) / limit.Window.Seconds()
	capacity := float64(limit.Requests)
	limiter := NewTokenBucket(rate, capacity)
	rl.routeLimiters[key] = limiter
	return limiter
}

// isWhitelistedIP 检查IP是否在白名单
func (rl *RateLimiter) isWhitelistedIP(ip string) bool {
	for _, whitelistIP := range rl.config.WhitelistIPs {
		if ip == whitelistIP || strings.HasPrefix(ip, whitelistIP) {
			return true
		}
	}
	return false
}

// isWhitelistedRole 检查用户角色是否在白名单
func (rl *RateLimiter) isWhitelistedRole(role string) bool {
	for _, whitelistRole := range rl.config.WhitelistRoles {
		if role == whitelistRole {
			return true
		}
	}
	return false
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
	// 首先检查X-Forwarded-For头
	xForwardedFor := c.GetHeader("X-Forwarded-For")
	if xForwardedFor != "" {
		// 可能包含多个IP，取第一个
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 检查X-Real-IP头
	xRealIP := c.GetHeader("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 回退到RemoteAddr
	remoteAddr := c.Request.RemoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		return remoteAddr[:idx]
	}

	return remoteAddr
}

// RateLimit 返回限流中间件
func RateLimit(config RateLimiterConfig) gin.HandlerFunc {
	limiter := NewRateLimiter(config)

	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		path := c.FullPath()
		method := c.Request.Method
		clientIP := getClientIP(c)

		// 检查IP白名单
		if limiter.isWhitelistedIP(clientIP) {
			c.Next()
			return
		}

		// 1. 检查路由限流（最高优先级）
		if routeLimit, exists := config.RouteLimits[path]; exists {
			key := fmt.Sprintf("%s:%s:%s", path, method, clientIP)
			if routeLimit.LimitByUser {
				// 按用户限流
				if userID, exists := getUserIDFromContext(c); exists {
					key = fmt.Sprintf("%s:%s:user:%d", path, method, userID)
				}
			} else if routeLimit.LimitByIP {
				// 按IP限流
				key = fmt.Sprintf("%s:%s:ip:%s", path, method, clientIP)
			}

			limiter := limiter.getRouteLimiter(key, routeLimit)
			if allowed, remaining := limiter.Allow(); !allowed {
				c.Header("X-RateLimit-Limit", strconv.Itoa(routeLimit.Requests))
				c.Header("X-RateLimit-Remaining", "0")
				c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(routeLimit.Window).Unix(), 10))
				respondAPIError(c, apierr.TooManyRequests(fmt.Sprintf("Rate limit exceeded for %s", path)))
				c.Abort()
				return
			} else {
				c.Header("X-RateLimit-Limit", strconv.Itoa(routeLimit.Requests))
				c.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
			}

			// 路由限流通过后，继续检查其他限流
		}

		// 2. 检查用户限流
		if userID, exists := getUserIDFromContext(c); exists {
			// 检查角色白名单
			if role, roleExists := getUserRoleFromContext(c); roleExists {
				if limiter.isWhitelistedRole(role) {
					c.Next()
					return
				}
			}

			userLimiter := limiter.getUserLimiter(userID)
			if allowed, remaining := userLimiter.Allow(); !allowed {
				c.Header("X-RateLimit-Limit", strconv.Itoa(config.UserRequestsPerMinute))
				c.Header("X-RateLimit-Remaining", "0")
				c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))
				respondAPIError(c, apierr.TooManyRequests("User rate limit exceeded"))
				c.Abort()
				return
			} else {
				c.Header("X-RateLimit-Limit", strconv.Itoa(config.UserRequestsPerMinute))
				c.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
			}
		}

		// 3. 检查IP限流
		ipLimiter := limiter.getIPLimiter(clientIP)
		if allowed, remaining := ipLimiter.Allow(); !allowed {
			c.Header("X-RateLimit-Limit", strconv.Itoa(int(config.IPRequestsPerSecond*60)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
			respondAPIError(c, apierr.TooManyRequests("IP rate limit exceeded"))
			c.Abort()
			return
		} else {
			c.Header("X-RateLimit-Limit", strconv.Itoa(int(config.IPRequestsPerSecond*60)))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
		}

		c.Next()
	}
}

// DefaultRateLimitConfig 返回默认限流配置
func DefaultRateLimitConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Enabled:               true,
		IPRequestsPerSecond:   10, // 每秒10个请求
		UserRequestsPerMinute: 60, // 每分钟60个请求
		WhitelistIPs:          []string{"127.0.0.1", "::1"},
		WhitelistRoles:        []string{"super_admin"},
		RouteLimits: map[string]RouteLimit{
			// 登录接口限流：每分钟5次
			"/api/v1/auth/login": {
				Path:        "/api/v1/auth/login",
				Requests:    5,
				Window:      time.Minute,
				LimitByIP:   true,
				LimitByUser: false,
			},
			// 注册接口限流：每小时10次
			"/api/v1/auth/register": {
				Path:        "/api/v1/auth/register",
				Requests:    10,
				Window:      time.Hour,
				LimitByIP:   true,
				LimitByUser: false,
			},
			// 支付接口限流：每分钟20次
			"/api/v1/user/orders/:id/pay": {
				Path:        "/api/v1/user/orders/:id/pay",
				Requests:    20,
				Window:      time.Minute,
				LimitByIP:   false,
				LimitByUser: true,
			},
		},
	}
}

// getUserIDFromContext 从上下文获取用户ID
func getUserIDFromContext(c *gin.Context) (uint64, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	userID, ok := userIDVal.(uint64)
	return userID, ok
}

// getUserRoleFromContext 从上下文获取用户角色
func getUserRoleFromContext(c *gin.Context) (string, bool) {
	roleVal, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := roleVal.(string)
	return role, ok
}

// respondAPIError 响应API错误（中间件内部使用）
func respondAPIError(c *gin.Context, err error) {
	if apiErr, ok := err.(*apierr.APIError); ok {
		c.JSON(apiErr.Code, gin.H{
			"success":   false,
			"code":      apiErr.Code,
			"message":   apiErr.Message,
			"details":   apiErr.Details,
			"requestId": c.GetString("request_id"),
		})
		return
	}
	// fallback
	c.JSON(http.StatusInternalServerError, gin.H{
		"success":   false,
		"code":      http.StatusInternalServerError,
		"message":   err.Error(),
		"requestId": c.GetString("request_id"),
	})
}

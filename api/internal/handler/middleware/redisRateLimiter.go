package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gamelink/pkg/cache"
	"github.com/gin-gonic/gin"
	redislib "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redislimiter "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RedisRateLimiterConfig Redis 限流配置
type RedisRateLimiterConfig struct {
	// 登录接口限流配置
	LoginRequests int
	LoginWindow   time.Duration

	// 短信接口限流配置
	SMSRequests int
	SMSWindow   time.Duration

	// 通用 API 限流配置
	APIRequests int
	APIWindow   time.Duration

	// 上传接口限流配置
	UploadRequests int
	UploadWindow   time.Duration

	// 是否启用
	Enabled bool
}

// DefaultRedisRateLimiterConfig 返回默认 Redis 限流配置
func DefaultRedisRateLimiterConfig() RedisRateLimiterConfig {
	return RedisRateLimiterConfig{
		Enabled:        true,
		LoginRequests:  5, // 5次/15分钟
		LoginWindow:    15 * time.Minute,
		SMSRequests:    3, // 3次/小时
		SMSWindow:      1 * time.Hour,
		APIRequests:    60, // 60次/分钟
		APIWindow:      1 * time.Minute,
		UploadRequests: 10, // 10次/小时
		UploadWindow:   1 * time.Hour,
	}
}

// RedisRateLimiter Redis 限流器
type RedisRateLimiter struct {
	redisClient *redislib.Client
	config      RedisRateLimiterConfig

	// 预配置的限流实例
	loginLimiter  *limiter.Limiter
	smsLimiter    *limiter.Limiter
	apiLimiter    *limiter.Limiter
	uploadLimiter *limiter.Limiter
}

// NewRedisRateLimiter 创建新的 Redis 限流器
func NewRedisRateLimiter(cacheClient cache.Cache, config RedisRateLimiterConfig) (*RedisRateLimiter, error) {
	// 尝试从 cache.Cache 获取底层 Redis 客户端
	var redisClient *redislib.Client

	// 使用类型断言检查是否有 GetClient 方法
	type withGetClient interface {
		GetClient() *redislib.Client
	}

	if wc, ok := cacheClient.(withGetClient); ok {
		redisClient = wc.GetClient()
	} else {
		return nil, fmt.Errorf("cache client does not support GetClient method")
	}

	// 创建 Redis store
	store, err := redislimiter.NewStore(redisClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	// 创建各个限流器
	loginRate := limiter.Rate{
		Period: config.LoginWindow,
		Limit:  int64(config.LoginRequests),
	}
	loginLimiter := limiter.New(store, loginRate)

	smsRate := limiter.Rate{
		Period: config.SMSWindow,
		Limit:  int64(config.SMSRequests),
	}
	smsLimiter := limiter.New(store, smsRate)

	apiRate := limiter.Rate{
		Period: config.APIWindow,
		Limit:  int64(config.APIRequests),
	}
	apiLimiter := limiter.New(store, apiRate)

	uploadRate := limiter.Rate{
		Period: config.UploadWindow,
		Limit:  int64(config.UploadRequests),
	}
	uploadLimiter := limiter.New(store, uploadRate)

	return &RedisRateLimiter{
		redisClient:   redisClient,
		config:        config,
		loginLimiter:  loginLimiter,
		smsLimiter:    smsLimiter,
		apiLimiter:    apiLimiter,
		uploadLimiter: uploadLimiter,
	}, nil
}

// LoginRateLimit 登录接口限流中间件
// 限制：5次/15分钟
func (rl *RedisRateLimiter) LoginRateLimit() gin.HandlerFunc {
	return rl.createRateLimitMiddleware(rl.loginLimiter, "login", rl.config.LoginRequests, rl.config.LoginWindow)
}

// SMSRateLimit 短信接口限流中间件
// 限制：3次/小时
func (rl *RedisRateLimiter) SMSRateLimit() gin.HandlerFunc {
	return rl.createRateLimitMiddleware(rl.smsLimiter, "sms", rl.config.SMSRequests, rl.config.SMSWindow)
}

// APIRateLimit 通用 API 限流中间件
// 限制：60次/分钟
func (rl *RedisRateLimiter) APIRateLimit() gin.HandlerFunc {
	return rl.createRateLimitMiddleware(rl.apiLimiter, "api", rl.config.APIRequests, rl.config.APIWindow)
}

// UploadRateLimit 上传接口限流中间件
// 限制：10次/小时
func (rl *RedisRateLimiter) UploadRateLimit() gin.HandlerFunc {
	return rl.createRateLimitMiddleware(rl.uploadLimiter, "upload", rl.config.UploadRequests, rl.config.UploadWindow)
}

// createRateLimitMiddleware 创建限流中间件
func (rl *RedisRateLimiter) createRateLimitMiddleware(limiterInstance *limiter.Limiter, prefix string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		// 获取客户端标识（IP 或用户 ID）
		key := rl.getClientKey(c, prefix)

		ctx := context.Background()
		limiterContext, err := limiterInstance.Get(ctx, key)
		if err != nil {
			// Redis 错误时记录日志但不阻止请求（fail-open）
			c.Next()
			return
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(int(limiterContext.Remaining)))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(limiterContext.Reset, 10))

		// 检查是否超过限制
		if limiterContext.Reached {
			c.Header("Retry-After", strconv.Itoa(int(window.Seconds())))
			c.JSON(429, gin.H{
				"success":   false,
				"code":      429,
				"message":   fmt.Sprintf("Rate limit exceeded: %d requests per %v", limit, window),
				"requestId": c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// getClientKey 获取客户端限流 key
func (rl *RedisRateLimiter) getClientKey(c *gin.Context, prefix string) string {
	// 优先使用用户 ID（如果已认证）
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("ratelimit:%s:user:%v", prefix, userID)
	}

	// 否则使用 IP 地址
	clientIP := GetClientIP(c)
	return fmt.Sprintf("ratelimit:%s:ip:%s", prefix, clientIP)
}

// CustomRateLimit 自定义限流中间件
// 允许为特定路由创建自定义限流规则
func (rl *RedisRateLimiter) CustomRateLimit(requests int, window time.Duration, prefix string) gin.HandlerFunc {
	store, err := redislimiter.NewStore(rl.redisClient)
	if err != nil {
		// 如果创建失败，返回一个不限制的中间件
		return func(c *gin.Context) {
			c.Next()
		}
	}

	rate := limiter.Rate{
		Period: window,
		Limit:  int64(requests),
	}
	limiterInstance := limiter.New(store, rate)

	return rl.createRateLimitMiddleware(limiterInstance, prefix, requests, window)
}

// GetClientIP 获取客户端真实 IP
func GetClientIP(c *gin.Context) string {
	// 检查 X-Forwarded-For 头
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// 取第一个 IP
		for i, ch := range xff {
			if ch == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// 检查 X-Real-IP 头
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// 使用 RemoteAddr
	return c.ClientIP()
}

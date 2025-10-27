package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	rlOnce     sync.Once
	rlRPS      float64
	rlBurst    int
	rlMu       sync.Mutex
	rlLimiters map[string]*rate.Limiter
)

func initLimiter() {
	rlRPS = 20.0
	rlBurst = 40
	if v := os.Getenv("ADMIN_RATE_RPS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			rlRPS = f
		}
	}
	if v := os.Getenv("ADMIN_RATE_BURST"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			rlBurst = i
		}
	}
	rlLimiters = make(map[string]*rate.Limiter)
}

func getLimiter(key string) *rate.Limiter {
	rlMu.Lock()
	defer rlMu.Unlock()
	if l, ok := rlLimiters[key]; ok {
		return l
	}
	l := rate.NewLimiter(rate.Limit(rlRPS), rlBurst)
	rlLimiters[key] = l
	return l
}

// RateLimitAdmin applies a token-bucket rate limit per user or IP for admin endpoints.
func RateLimitAdmin() gin.HandlerFunc {
	rlOnce.Do(initLimiter)
	return func(c *gin.Context) {
		// Prefer authenticated user id when available, else fall back to client IP
		if uid, ok := c.Get("user_id"); ok {
			key := fmt.Sprintf("user:%v", uid)
			if !getLimiter(key).Allow() {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"success": false,
					"code":    http.StatusTooManyRequests,
					"message": "rate limit exceeded",
				})
				return
			}
			c.Next()
			return
		}
		key := "ip:" + c.ClientIP()
		if !getLimiter(key).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"code":    http.StatusTooManyRequests,
				"message": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

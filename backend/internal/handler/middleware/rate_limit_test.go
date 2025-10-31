package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimitAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 重置限流器状态
	rlOnce = *new(sync.Once)
	rlLimiters = nil

	// 设置测试环境变量
	os.Setenv("ADMIN_RATE_RPS", "10")
	os.Setenv("ADMIN_RATE_BURST", "20")
	defer func() {
		os.Unsetenv("ADMIN_RATE_RPS")
		os.Unsetenv("ADMIN_RATE_BURST")
	}()

	router := gin.New()
	router.Use(RateLimitAdmin())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	t.Run("未超限-允许请求", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
		}
	})

	t.Run("根据用户ID限流", func(t *testing.T) {
		routerWithUser := gin.New()
		routerWithUser.Use(func(c *gin.Context) {
			c.Set("user_id", uint64(123))
		})
		routerWithUser.Use(RateLimitAdmin())
		routerWithUser.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// 发送多个请求以触发限流
		successCount := 0
		limitedCount := 0

		for i := 0; i < 25; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			routerWithUser.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				limitedCount++
			}
		}

		// Burst 是 20，所以前 20 个应该成功
		if successCount < 20 {
			t.Errorf("期望至少 20 个成功请求, 得到 %d", successCount)
		}

		t.Logf("成功请求: %d, 被限流: %d", successCount, limitedCount)
	})

	t.Run("根据IP限流", func(t *testing.T) {
		// 重置限流器
		rlOnce = *new(sync.Once)
		rlLimiters = nil

		routerWithIP := gin.New()
		routerWithIP.Use(RateLimitAdmin())
		routerWithIP.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// 发送多个来自同一IP的请求
		successCount := 0
		limitedCount := 0

		for i := 0; i < 25; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.RemoteAddr = "192.168.1.100:12345"
			w := httptest.NewRecorder()

			routerWithIP.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				successCount++
			} else if w.Code == http.StatusTooManyRequests {
				limitedCount++
			}
		}

		if successCount < 20 {
			t.Errorf("期望至少 20 个成功请求, 得到 %d", successCount)
		}

		t.Logf("IP限流 - 成功请求: %d, 被限流: %d", successCount, limitedCount)
	})

	t.Run("超限后返回429", func(t *testing.T) {
		// 重置限流器为非常严格的限制
		rlOnce = *new(sync.Once)
		rlLimiters = nil
		os.Setenv("ADMIN_RATE_RPS", "1")
		os.Setenv("ADMIN_RATE_BURST", "2")
		defer func() {
			os.Setenv("ADMIN_RATE_RPS", "10")
			os.Setenv("ADMIN_RATE_BURST", "20")
		}()

		strictRouter := gin.New()
		strictRouter.Use(RateLimitAdmin())
		strictRouter.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		// 快速发送多个请求
		var found429 bool
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			strictRouter.ServeHTTP(w, req)

			if w.Code == http.StatusTooManyRequests {
				found429 = true
				// 验证响应格式
				body := w.Body.String()
				if len(body) == 0 || body[:1] != "{" {
					t.Error("期望返回 JSON 响应")
				}
				break
			}
		}

		if !found429 {
			t.Error("期望触发 429 Too Many Requests")
		}
	})
}

func TestInitLimiter(t *testing.T) {
	// 重置限流器
	rlOnce = *new(sync.Once)
	rlLimiters = nil

	t.Run("使用默认配置", func(t *testing.T) {
		os.Unsetenv("ADMIN_RATE_RPS")
		os.Unsetenv("ADMIN_RATE_BURST")

		rlOnce.Do(initLimiter)

		if rlRPS != 20.0 {
			t.Errorf("期望默认 RPS 20.0, 得到 %f", rlRPS)
		}
		if rlBurst != 40 {
			t.Errorf("期望默认 Burst 40, 得到 %d", rlBurst)
		}
		if rlLimiters == nil {
			t.Error("期望 rlLimiters 被初始化")
		}
	})

	t.Run("使用环境变量配置", func(t *testing.T) {
		rlOnce = *new(sync.Once)
		rlLimiters = nil

		os.Setenv("ADMIN_RATE_RPS", "50")
		os.Setenv("ADMIN_RATE_BURST", "100")
		defer func() {
			os.Unsetenv("ADMIN_RATE_RPS")
			os.Unsetenv("ADMIN_RATE_BURST")
		}()

		rlOnce.Do(initLimiter)

		if rlRPS != 50.0 {
			t.Errorf("期望 RPS 50.0, 得到 %f", rlRPS)
		}
		if rlBurst != 100 {
			t.Errorf("期望 Burst 100, 得到 %d", rlBurst)
		}
	})

	t.Run("无效的环境变量使用默认值", func(t *testing.T) {
		rlOnce = *new(sync.Once)
		rlLimiters = nil

		os.Setenv("ADMIN_RATE_RPS", "invalid")
		os.Setenv("ADMIN_RATE_BURST", "invalid")
		defer func() {
			os.Unsetenv("ADMIN_RATE_RPS")
			os.Unsetenv("ADMIN_RATE_BURST")
		}()

		rlOnce.Do(initLimiter)

		if rlRPS != 20.0 {
			t.Errorf("期望默认 RPS 20.0, 得到 %f", rlRPS)
		}
		if rlBurst != 40 {
			t.Errorf("期望默认 Burst 40, 得到 %d", rlBurst)
		}
	})

	t.Run("负数环境变量使用默认值", func(t *testing.T) {
		rlOnce = *new(sync.Once)
		rlLimiters = nil

		os.Setenv("ADMIN_RATE_RPS", "-10")
		os.Setenv("ADMIN_RATE_BURST", "-20")
		defer func() {
			os.Unsetenv("ADMIN_RATE_RPS")
			os.Unsetenv("ADMIN_RATE_BURST")
		}()

		rlOnce.Do(initLimiter)

		if rlRPS != 20.0 {
			t.Errorf("期望默认 RPS 20.0, 得到 %f", rlRPS)
		}
		if rlBurst != 40 {
			t.Errorf("期望默认 Burst 40, 得到 %d", rlBurst)
		}
	})
}

func TestGetLimiter(t *testing.T) {
	rlOnce = *new(sync.Once)
	rlOnce.Do(initLimiter)

	t.Run("创建新限流器", func(t *testing.T) {
		key := "test:user:123"
		limiter := getLimiter(key)

		if limiter == nil {
			t.Fatal("期望创建新的限流器")
		}

		// 验证限流器被缓存
		limiter2 := getLimiter(key)
		if limiter != limiter2 {
			t.Error("期望返回相同的限流器实例")
		}
	})

	t.Run("不同key创建不同限流器", func(t *testing.T) {
		limiter1 := getLimiter("user:1")
		limiter2 := getLimiter("user:2")

		if limiter1 == limiter2 {
			t.Error("期望不同key返回不同的限流器")
		}
	})

	t.Run("并发访问安全", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(id int) {
				key := "concurrent:user:" + string(rune(id))
				limiter := getLimiter(key)
				if limiter == nil {
					t.Error("并发访问时创建限流器失败")
				}
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

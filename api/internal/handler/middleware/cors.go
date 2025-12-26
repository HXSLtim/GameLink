package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 中间件处理跨域请求
//
// 技术原理：
// 1. 浏览器在发送跨域请求前会先发送OPTIONS预检请求
// 2. 后端需要返回正确的CORS头信息，告诉浏览器允许哪些跨域请求
// 3. 对于简单请求，浏览器直接发送请求并检查响应头
// 4. 对于复杂请求（如PUT、DELETE等），浏览器先发送OPTIONS预检请求
//
// 常见CORS头：
// - Access-Control-Allow-Origin: 允许的源（*表示所有源）
// - Access-Control-Allow-Methods: 允许的HTTP方法
// - Access-Control-Allow-Headers: 允许的请求头
// - Access-Control-Allow-Credentials: 是否允许携带认证信息
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 设置允许的源
		// 开发环境：允许所有源
		// 生产环境：应该配置具体的域名
		allowedOrigins := getAllowedOrigins()

		// 检查请求源是否在允许列表中
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			// 设置CORS响应头
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400") // 预检请求缓存24小时
		}

		// 处理预检请求（OPTIONS）
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getAllowedOrigins 根据环境获取允许的源
func getAllowedOrigins() []string {
	env := os.Getenv("APP_ENV")
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))

	// parse comma-separated list
	parse := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			if v := strings.TrimSpace(p); v != "" {
				out = append(out, v)
			}
		}
		return out
	}

	if env == "production" {
		// In production, default to EMPTY (deny) unless explicitly configured.
		// This prevents insecure wildcard. Set CORS_ALLOWED_ORIGINS to a comma-separated list.
		return parse(raw)
	}

	// Development: allow configured list if present, otherwise wildcard.
	if list := parse(raw); len(list) > 0 {
		return list
	}
	return []string{"*"}
}

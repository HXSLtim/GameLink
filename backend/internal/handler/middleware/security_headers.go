package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig 安全头配置
type SecurityHeadersConfig struct {
	// XFrameOptions X-Frame-Options header值
	XFrameOptions string
	// XContentTypeOptions X-Content-Type-Options header值
	XContentTypeOptions string
	// XXSSProtection X-XSS-Protection header值
	XXSSProtection string
	// ContentSecurityPolicy Content-Security-Policy header值
	ContentSecurityPolicy string
	// StrictTransportSecurity Strict-Transport-Security header值
	StrictTransportSecurity string
	// ReferrerPolicy Referrer-Policy header值
	ReferrerPolicy string
	// PermissionsPolicy Permissions-Policy header值
	PermissionsPolicy string
}

// DefaultSecurityHeadersConfig 默认安全头配置
var DefaultSecurityHeadersConfig = SecurityHeadersConfig{
	XFrameOptions:           "DENY",
	XContentTypeOptions:     "nosniff",
	XXSSProtection:          "1; mode=block",
	ContentSecurityPolicy:   "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'",
	StrictTransportSecurity: "max-age=31536000; includeSubDomains",
	ReferrerPolicy:          "strict-origin-when-cross-origin",
	PermissionsPolicy:       "geolocation=(), microphone=(), camera=()",
}

// SecurityHeaders 返回安全头中间件
func SecurityHeaders(config ...SecurityHeadersConfig) gin.HandlerFunc {
	cfg := DefaultSecurityHeadersConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// X-Frame-Options: 防止点击劫持
		if cfg.XFrameOptions != "" {
			c.Header("X-Frame-Options", cfg.XFrameOptions)
		}

		// X-Content-Type-Options: 防止MIME类型嗅探
		if cfg.XContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", cfg.XContentTypeOptions)
		}

		// X-XSS-Protection: 启用浏览器XSS过滤
		if cfg.XXSSProtection != "" {
			c.Header("X-XSS-Protection", cfg.XXSSProtection)
		}

		// Content-Security-Policy: 内容安全策略
		if cfg.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", cfg.ContentSecurityPolicy)
		}

		// Strict-Transport-Security: 强制HTTPS
		if cfg.StrictTransportSecurity != "" {
			c.Header("Strict-Transport-Security", cfg.StrictTransportSecurity)
		}

		// Referrer-Policy: 控制Referer信息
		if cfg.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", cfg.ReferrerPolicy)
		}

		// Permissions-Policy: 控制浏览器功能权限
		if cfg.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", cfg.PermissionsPolicy)
		}

		c.Next()
	}
}

// SecureHeaders 快捷方法，使用默认配置
func SecureHeaders() gin.HandlerFunc {
	return SecurityHeaders()
}

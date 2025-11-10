package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFConfig CSRF保护配置
type CSRFConfig struct {
	// TokenLength CSRF token长度（字节数）
	TokenLength int
	// CookieName CSRF cookie名称
	CookieName string
	// HeaderName CSRF header名称
	HeaderName string
	// FormFieldName CSRF表单字段名称
	FormFieldName string
	// Secret 用于加密token的密钥
	Secret string
	// CookiePath Cookie路径
	CookiePath string
	// CookieDomain Cookie域名
	CookieDomain string
	// CookieSecure 是否只在HTTPS下发送Cookie
	CookieSecure bool
	// CookieHTTPOnly 是否禁止JavaScript访问Cookie
	CookieHTTPOnly bool
	// CookieSameSite Cookie的SameSite属性
	CookieSameSite http.SameSite
	// TokenLookup Token查找位置，格式: "header:<name>" 或 "form:<name>"
	TokenLookup string
	// SkipCheck 跳过检查的函数
	SkipCheck func(*gin.Context) bool
}

// DefaultCSRFConfig 默认CSRF配置
var DefaultCSRFConfig = CSRFConfig{
	TokenLength:    32,
	CookieName:     "_csrf",
	HeaderName:     "X-CSRF-Token",
	FormFieldName:  "_csrf",
	CookiePath:     "/",
	CookieSecure:   true,
	CookieHTTPOnly: true,
	CookieSameSite: http.SameSiteStrictMode,
	TokenLookup:    "header:X-CSRF-Token",
}

var (
	// ErrCSRFTokenMissing CSRF token缺失错误
	ErrCSRFTokenMissing = errors.New("csrf token missing")
	// ErrCSRFTokenInvalid CSRF token无效错误
	ErrCSRFTokenInvalid = errors.New("csrf token invalid")
)

// CSRF 返回CSRF保护中间件
func CSRF(config ...CSRFConfig) gin.HandlerFunc {
	cfg := DefaultCSRFConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 设置默认值
	if cfg.TokenLength == 0 {
		cfg.TokenLength = DefaultCSRFConfig.TokenLength
	}
	if cfg.CookieName == "" {
		cfg.CookieName = DefaultCSRFConfig.CookieName
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = DefaultCSRFConfig.HeaderName
	}
	if cfg.FormFieldName == "" {
		cfg.FormFieldName = DefaultCSRFConfig.FormFieldName
	}
	if cfg.CookiePath == "" {
		cfg.CookiePath = DefaultCSRFConfig.CookiePath
	}
	if cfg.TokenLookup == "" {
		cfg.TokenLookup = DefaultCSRFConfig.TokenLookup
	}

	return func(c *gin.Context) {
		// 跳过检查
		if cfg.SkipCheck != nil && cfg.SkipCheck(c) {
			c.Next()
			return
		}

		// 生成或获取CSRF token
		token := getCSRFTokenFromCookie(c, cfg.CookieName)
		if token == "" {
			token = generateCSRFToken(cfg.TokenLength)
			setCSRFCookie(c, cfg, token)
		}

		// 将token存储到context中，供后续使用
		c.Set("csrf_token", token)

		// 对于非安全方法（POST, PUT, DELETE, PATCH），验证CSRF token
		if !isMethodSafe(c.Request.Method) {
			clientToken := extractCSRFToken(c, cfg)
			if clientToken == "" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"code":    http.StatusForbidden,
					"message": "CSRF token missing",
				})
				return
			}

			if !validateCSRFToken(token, clientToken) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"success": false,
					"code":    http.StatusForbidden,
					"message": "CSRF token invalid",
				})
				return
			}
		}

		c.Next()
	}
}

// generateCSRFToken 生成CSRF token
func generateCSRFToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}
	return base64.URLEncoding.EncodeToString(b)
}

// getCSRFTokenFromCookie 从Cookie中获取CSRF token
func getCSRFTokenFromCookie(c *gin.Context, cookieName string) string {
	token, err := c.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return token
}

// setCSRFCookie 设置CSRF Cookie
func setCSRFCookie(c *gin.Context, cfg CSRFConfig, token string) {
	c.SetSameSite(cfg.CookieSameSite)
	c.SetCookie(
		cfg.CookieName,
		token,
		3600*24, // 24小时
		cfg.CookiePath,
		cfg.CookieDomain,
		cfg.CookieSecure,
		cfg.CookieHTTPOnly,
	)
}

// extractCSRFToken 从请求中提取CSRF token
func extractCSRFToken(c *gin.Context, cfg CSRFConfig) string {
	// 优先从header中获取
	token := c.GetHeader(cfg.HeaderName)
	if token != "" {
		return token
	}

	// 从表单中获取
	token = c.PostForm(cfg.FormFieldName)
	if token != "" {
		return token
	}

	return ""
}

// validateCSRFToken 验证CSRF token
func validateCSRFToken(expected, actual string) bool {
	if expected == "" || actual == "" {
		return false
	}
	// 使用constant-time比较防止时序攻击
	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

// isMethodSafe 判断HTTP方法是否安全
func isMethodSafe(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// GetCSRFToken 从context中获取CSRF token（用于模板渲染）
func GetCSRFToken(c *gin.Context) string {
	if token, exists := c.Get("csrf_token"); exists {
		if tokenStr, ok := token.(string); ok {
			return tokenStr
		}
	}
	return ""
}

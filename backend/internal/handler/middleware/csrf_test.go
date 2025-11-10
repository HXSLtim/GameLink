package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCSRF(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GET请求不需要CSRF验证", func(t *testing.T) {
		router := gin.New()
		router.Use(CSRF())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// 应该设置CSRF cookie
		cookies := w.Result().Cookies()
		assert.NotEmpty(t, cookies)
		var csrfCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "_csrf" {
				csrfCookie = cookie
				break
			}
		}
		assert.NotNil(t, csrfCookie)
		assert.NotEmpty(t, csrfCookie.Value)
	})

	t.Run("POST请求缺少CSRF token应该被拒绝", func(t *testing.T) {
		router := gin.New()
		router.Use(CSRF())
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "CSRF token missing")
	})

	t.Run("POST请求带有效CSRF token应该通过", func(t *testing.T) {
		router := gin.New()
		// 使用测试配置，禁用Secure以便在测试环境工作
		router.Use(CSRF(CSRFConfig{
			CookieSecure: false,
		}))
		router.GET("/test-get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"token": GetCSRFToken(c)})
		})
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 先发GET请求获取CSRF token
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/test-get", nil)
		router.ServeHTTP(w1, req1)

		// 从响应body中获取token（这是未编码的原始token）
		var respBody map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &respBody)
		csrfToken, _ := respBody["token"].(string)
		assert.NotEmpty(t, csrfToken, "CSRF token should be set")

		// 发POST请求，带上CSRF token（header和cookie都要带）
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/test", nil)
		req2.Header.Set("X-CSRF-Token", csrfToken)
		// 必须带上cookie，因为服务器会从cookie读取token进行比较
		req2.AddCookie(&http.Cookie{
			Name:  "_csrf",
			Value: csrfToken,
		})
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code, "Response body: %s", w2.Body.String())
	})

	t.Run("POST请求带无效CSRF token应该被拒绝", func(t *testing.T) {
		router := gin.New()
		router.Use(CSRF())
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/test", nil)
		req.Header.Set("X-CSRF-Token", "invalid-token")
		req.AddCookie(&http.Cookie{
			Name:  "_csrf",
			Value: "different-token",
		})
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "CSRF token invalid")
	})

	t.Run("从表单字段中提取CSRF token", func(t *testing.T) {
		router := gin.New()
		router.Use(CSRF())
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 先获取token
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/test-get", nil)
		router.GET("/test-get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"token": GetCSRFToken(c)})
		})
		router.ServeHTTP(w1, req1)

		var csrfToken string
		cookies := w1.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "_csrf" {
				csrfToken = cookie.Value
				break
			}
		}

		// 通过表单提交token
		w2 := httptest.NewRecorder()
		formData := "_csrf=" + csrfToken
		req2 := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(formData))
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req2.AddCookie(&http.Cookie{
			Name:  "_csrf",
			Value: csrfToken,
		})
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
	})

	t.Run("自定义配置", func(t *testing.T) {
		config := CSRFConfig{
			CookieName:    "custom_csrf",
			HeaderName:    "X-Custom-CSRF",
			FormFieldName: "custom_csrf_token",
			CookieSecure:  false, // 测试环境禁用Secure
		}

		router := gin.New()
		router.Use(CSRF(config))
		router.GET("/test-get", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"token": GetCSRFToken(c)})
		})
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// 先发GET请求获取token
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/test-get", nil)
		router.ServeHTTP(w1, req1)

		// 从响应body中获取token（未编码的原始token）
		var respBody map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &respBody)
		csrfToken, _ := respBody["token"].(string)
		assert.NotEmpty(t, csrfToken, "Custom CSRF token should be set")

		// 使用自定义header提交
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/test", nil)
		req2.Header.Set("X-Custom-CSRF", csrfToken)
		req2.AddCookie(&http.Cookie{
			Name:  "custom_csrf",
			Value: csrfToken,
		})
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code, "Response body: %s", w2.Body.String())
	})

	t.Run("SkipCheck跳过验证", func(t *testing.T) {
		config := CSRFConfig{
			SkipCheck: func(c *gin.Context) bool {
				// 跳过/api/webhook路径的CSRF检查
				return strings.HasPrefix(c.Request.URL.Path, "/api/webhook")
			},
		}

		router := gin.New()
		router.Use(CSRF(config))
		router.POST("/api/webhook/payment", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/webhook/payment", nil)
		router.ServeHTTP(w, req)

		// 应该通过，不需要CSRF token
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGenerateCSRFToken(t *testing.T) {
	t.Run("生成的token应该不为空", func(t *testing.T) {
		token := generateCSRFToken(32)
		assert.NotEmpty(t, token)
	})

	t.Run("生成的token应该是唯一的", func(t *testing.T) {
		token1 := generateCSRFToken(32)
		token2 := generateCSRFToken(32)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("生成的token长度应该正确", func(t *testing.T) {
		token := generateCSRFToken(16)
		assert.NotEmpty(t, token)
		// Base64编码后的长度会比原始字节长度大
		assert.Greater(t, len(token), 16)
	})
}

func TestValidateCSRFToken(t *testing.T) {
	t.Run("相同的token应该验证通过", func(t *testing.T) {
		token := "test-token-123"
		assert.True(t, validateCSRFToken(token, token))
	})

	t.Run("不同的token应该验证失败", func(t *testing.T) {
		token1 := "test-token-123"
		token2 := "test-token-456"
		assert.False(t, validateCSRFToken(token1, token2))
	})

	t.Run("空token应该验证失败", func(t *testing.T) {
		assert.False(t, validateCSRFToken("", "test"))
		assert.False(t, validateCSRFToken("test", ""))
		assert.False(t, validateCSRFToken("", ""))
	})
}

func TestIsMethodSafe(t *testing.T) {
	safeMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	unsafeMethods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range safeMethods {
		t.Run(method+" should be safe", func(t *testing.T) {
			assert.True(t, isMethodSafe(method))
		})
	}

	for _, method := range unsafeMethods {
		t.Run(method+" should be unsafe", func(t *testing.T) {
			assert.False(t, isMethodSafe(method))
		})
	}
}

func TestGetCSRFToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("应该能从context中获取token", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		expectedToken := "test-csrf-token"
		c.Set("csrf_token", expectedToken)

		token := GetCSRFToken(c)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("context中没有token时应该返回空字符串", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		token := GetCSRFToken(c)
		assert.Empty(t, token)
	})

	t.Run("token类型不正确时应该返回空字符串", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("csrf_token", 12345) // 设置错误类型
		token := GetCSRFToken(c)
		assert.Empty(t, token)
	})
}

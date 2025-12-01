package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/auth"
)

func TestJWTAuth(t *testing.T) {
	// 设置测试用的JWT密钥
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")
	defer os.Unsetenv("JWT_SECRET_KEY")

	gin.SetMode(gin.TestMode)

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Equal(t, float64(http.StatusUnauthorized), response["code"].(float64))
		assert.Contains(t, response["message"].(string), "缺少Authorization")
	})

	t.Run("InvalidTokenFormat", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("ValidToken", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			userID, exists := GetUserID(c)
			assert.True(t, exists)
			assert.Equal(t, uint64(123), userID)

			role, exists := GetUserRole(c)
			assert.True(t, exists)
			assert.Equal(t, "user", role)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 生成有效Token
		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 生成过期Token
		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", -1*time.Hour)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		// 验证错误消息包含过期相关信息
		assert.Contains(t, response["message"].(string), "无效")
	})

	t.Run("MissingSecretKey", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET_KEY")
		defer os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("ShortSecretKey", func(t *testing.T) {
		os.Setenv("JWT_SECRET_KEY", "short")
		defer os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	})

	t.Run("TokenAutoRefreshNotTriggered", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 生成较长时间的Token(不会触发自动刷新)
		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 2*time.Hour)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证没有返回新Token(因为时间还很长)
		newToken := w.Header().Get("X-Refreshed-Token")
		assert.Empty(t, newToken)
	})

	t.Run("TokenRefreshRecommendation", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		// 生成1小时内过期的Token(50分钟)
		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 55*time.Minute)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// 验证是否提示刷新
		refreshRec := w.Header().Get("X-Token-Refresh-Recommendation")
		assert.Equal(t, "true", refreshRec)
	})
}

func TestRequireRole(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")
	defer os.Unsetenv("JWT_SECRET_KEY")

	gin.SetMode(gin.TestMode)

	t.Run("AdminAccessSuccess", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.Use(RequireRole("admin"))
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
		token, err := jwtManager.GenerateToken(123, "admin")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UserAccessAdminForbidden", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.Use(RequireRole("admin"))
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access"})
		})

		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("MultipleRoles", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.Use(RequireRole("admin", "moderator"))
		router.GET("/mod", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "mod access"})
		})

		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
		token, err := jwtManager.GenerateToken(123, "moderator")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/mod", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestOptionalAuth(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")
	defer os.Unsetenv("JWT_SECRET_KEY")

	gin.SetMode(gin.TestMode)

	t.Run("NoToken", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/public", func(c *gin.Context) {
			authenticated := IsAuthenticated(c)
			assert.False(t, authenticated)
			c.JSON(http.StatusOK, gin.H{"message": "public access"})
		})

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ValidToken", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/public", func(c *gin.Context) {
			authenticated := IsAuthenticated(c)
			assert.True(t, authenticated)

			userID, exists := GetUserID(c)
			assert.True(t, exists)
			assert.Equal(t, uint64(123), userID)

			c.JSON(http.StatusOK, gin.H{"message": "authenticated access"})
		})

		jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
		token, err := jwtManager.GenerateToken(123, "user")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/public", func(c *gin.Context) {
			authenticated := IsAuthenticated(c)
			assert.False(t, authenticated)
			c.JSON(http.StatusOK, gin.H{"message": "public access"})
		})

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("MissingSecretKey_ShouldFail", func(t *testing.T) {
		os.Unsetenv("JWT_SECRET_KEY")
		defer os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
		})

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 验证返回503错误,而不是允许访问
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["message"].(string), "认证服务配置错误")
	})

	t.Run("ShortSecretKey_ShouldFail", func(t *testing.T) {
		os.Setenv("JWT_SECRET_KEY", "short")
		defer os.Setenv("JWT_SECRET_KEY", "test-secret-key-that-is-32-characters-long")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
		})

		req := httptest.NewRequest(http.MethodGet, "/public", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// 验证返回503错误,而不是允许访问
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["message"].(string), "认证服务配置错误")
	})
}

func TestGetUserIDAndRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("GetUserIDSuccess", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint64(123))

		userID, exists := GetUserID(c)
		assert.True(t, exists)
		assert.Equal(t, uint64(123), userID)
	})

	t.Run("GetUserIDNotExists", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		userID, exists := GetUserID(c)
		assert.False(t, exists)
		assert.Equal(t, uint64(0), userID)
	})

	t.Run("GetUserIDWrongType", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", "123") // string instead of uint64

		userID, exists := GetUserID(c)
		assert.False(t, exists)
		assert.Equal(t, uint64(0), userID)
	})

	t.Run("GetUserRoleSuccess", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "admin")

		role, exists := GetUserRole(c)
		assert.True(t, exists)
		assert.Equal(t, "admin", role)
	})

	t.Run("IsAuthenticated", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("is_authenticated", true)

		assert.True(t, IsAuthenticated(c))
	})

	t.Run("IsNotAuthenticated", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		assert.False(t, IsAuthenticated(c))
	})
}

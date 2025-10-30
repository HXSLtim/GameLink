package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
)

func TestJWTAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 保存原始环境变量
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		os.Setenv("JWT_SECRET_KEY", originalSecret)
		os.Setenv("APP_ENV", originalEnv)
	}()

	testSecret := "test-jwt-secret-key-for-testing"
	os.Setenv("JWT_SECRET_KEY", testSecret)
	os.Setenv("APP_ENV", "development")

	jwtManager := auth.NewJWTManager(testSecret, 24*time.Hour)

	t.Run("成功验证有效Token", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(123, "user")

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			userID, _ := GetUserID(c)
			userRole, _ := GetUserRole(c)
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"role":    userRole,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("缺少Authorization头-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("无效的Token格式-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Authorization", "InvalidToken")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("无效的Token-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Token快过期-设置刷新提示头", func(t *testing.T) {
		// 创建一个即将过期的 token (1小时内)
		shortManager := auth.NewJWTManager(testSecret, 30*time.Minute)
		token, _ := shortManager.GenerateToken(123, "user")

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// 检查是否设置了刷新提示头
		if w.Header().Get("X-Token-Refresh") != "true" {
			t.Error("Expected X-Token-Refresh header to be set")
		}
	})

	t.Run("生产环境未配置JWT-拒绝访问", func(t *testing.T) {
		os.Setenv("JWT_SECRET_KEY", "")
		os.Setenv("APP_ENV", "production")

		router := gin.New()
		router.Use(JWTAuth())
		router.GET("/api/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		// 恢复设置
		os.Setenv("JWT_SECRET_KEY", testSecret)
		os.Setenv("APP_ENV", "development")
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("用户有正确角色-允许访问", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", "admin")
			c.Next()
		})
		router.Use(RequireRole("admin", "moderator"))
		router.GET("/api/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("用户角色不在允许列表-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", "user")
			c.Next()
		})
		router.Use(RequireRole("admin", "moderator"))
		router.GET("/api/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("用户未认证-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(RequireRole("admin"))
		router.GET("/api/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("用户角色类型错误-拒绝访问", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", 123) // 错误类型
			c.Next()
		})
		router.Use(RequireRole("admin"))
		router.GET("/api/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
		}
	})
}

func TestOptionalAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 保存原始环境变量
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		os.Setenv("JWT_SECRET_KEY", originalSecret)
		os.Setenv("APP_ENV", originalEnv)
	}()

	testSecret := "test-jwt-secret-key-for-testing"
	os.Setenv("JWT_SECRET_KEY", testSecret)
	os.Setenv("APP_ENV", "development")

	jwtManager := auth.NewJWTManager(testSecret, 24*time.Hour)

	t.Run("没有Token-允许继续（匿名访问）", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/api/public", func(c *gin.Context) {
			isAuth := IsAuthenticated(c)
			c.JSON(http.StatusOK, gin.H{
				"authenticated": isAuth,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("有效Token-设置用户信息", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(123, "user")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/api/public", func(c *gin.Context) {
			isAuth := IsAuthenticated(c)
			userID, _ := GetUserID(c)
			c.JSON(http.StatusOK, gin.H{
				"authenticated": isAuth,
				"user_id":       userID,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("无效Token-允许继续（匿名访问）", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/api/public", func(c *gin.Context) {
			isAuth := IsAuthenticated(c)
			c.JSON(http.StatusOK, gin.H{
				"authenticated": isAuth,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Token格式错误-允许继续（匿名访问）", func(t *testing.T) {
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/api/public", func(c *gin.Context) {
			isAuth := IsAuthenticated(c)
			c.JSON(http.StatusOK, gin.H{
				"authenticated": isAuth,
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("生产环境未配置JWT-允许继续（匿名访问）", func(t *testing.T) {
		os.Setenv("JWT_SECRET_KEY", "")
		os.Setenv("APP_ENV", "production")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/api/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest(http.MethodGet, "/api/public", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// 恢复设置
		os.Setenv("JWT_SECRET_KEY", testSecret)
		os.Setenv("APP_ENV", "development")
	})
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取用户ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint64(123))

		userID, ok := GetUserID(c)
		if !ok {
			t.Error("Expected to get user ID")
		}
		if userID != 123 {
			t.Errorf("Expected user ID 123, got %d", userID)
		}
	})

	t.Run("用户ID不存在", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, ok := GetUserID(c)
		if ok {
			t.Error("Expected not to get user ID")
		}
	})

	t.Run("用户ID类型错误", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", "123") // 错误类型

		_, ok := GetUserID(c)
		if ok {
			t.Error("Expected not to get user ID when type is wrong")
		}
	})
}

func TestGetUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("成功获取用户角色", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", "admin")

		role, ok := GetUserRole(c)
		if !ok {
			t.Error("Expected to get user role")
		}
		if role != "admin" {
			t.Errorf("Expected role 'admin', got '%s'", role)
		}
	})

	t.Run("用户角色不存在", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		_, ok := GetUserRole(c)
		if ok {
			t.Error("Expected not to get user role")
		}
	})

	t.Run("用户角色类型错误", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_role", 123) // 错误类型

		_, ok := GetUserRole(c)
		if ok {
			t.Error("Expected not to get user role when type is wrong")
		}
	})
}

func TestIsAuthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("用户已认证", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("is_authenticated", true)

		if !IsAuthenticated(c) {
			t.Error("Expected user to be authenticated")
		}
	})

	t.Run("用户未认证", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("is_authenticated", false)

		if IsAuthenticated(c) {
			t.Error("Expected user not to be authenticated")
		}
	})

	t.Run("认证标志不存在", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		if IsAuthenticated(c) {
			t.Error("Expected user not to be authenticated when flag doesn't exist")
		}
	})

	t.Run("认证标志类型错误", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("is_authenticated", "true") // 错误类型

		if IsAuthenticated(c) {
			t.Error("Expected user not to be authenticated when type is wrong")
		}
	})
}


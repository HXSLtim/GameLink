package middleware

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"

    "gamelink/internal/auth"
    "gamelink/internal/logging"
)

// JWTAuth JWT认证中间件
//
// 使用方法：
// router.Use(middleware.JWTAuth())
// 或者
// adminGroup.Use(middleware.JWTAuth())
func JWTAuth() gin.HandlerFunc {
	// 从环境变量获取JWT密钥
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		if os.Getenv("APP_ENV") == "production" {
			return func(c *gin.Context) {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"success": false,
					"code":    http.StatusServiceUnavailable,
					"message": "jwt not configured",
				})
			}
		}
		// 开发环境使用默认值
		secretKey = "gamelink-default-secret-key-change-in-production"
	}

	// Token有效期（24小时）
	tokenDuration := 24 * time.Hour

	// 创建JWT管理器
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "缺少Authorization头",
			})
			c.Abort()
			return
		}

		// 提取Token
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证Token
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "无效的Token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 检查Token是否过期
		if auth.IsTokenExpired(claims) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "Token已过期",
			})
			c.Abort()
			return
		}

        // 将用户信息存储到Context中，供后续处理使用
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Set("jwt_claims", claims)
        // 注入 actor 到 request context，便于服务层审计日志使用
        c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), claims.UserID))

		// 检查Token剩余时间，如果快要过期，在响应头中提示前端刷新Token
		remainingTime := auth.GetTokenRemainingTime(claims)
		if remainingTime < 1*time.Hour {
			c.Header("X-Token-Refresh", "true")
			c.Header("X-Token-Remaining", remainingTime.String())
		}

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
//
// 使用方法：
// router.Use(middleware.RequireRole("admin"))
// 或者
// adminGroup.Use(middleware.RequireRole("admin", "moderator"))
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "用户未认证",
			})
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "用户角色格式错误",
			})
			c.Abort()
			return
		}

		// 检查角色权限
		hasPermission := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
//
// 如果提供了Token则验证，如果没有提供Token则允许继续
// 适用于那些既可以登录访问也可以匿名访问的接口
func OptionalAuth() gin.HandlerFunc {
	// 从环境变量获取JWT密钥
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		if os.Getenv("APP_ENV") == "production" {
			// 生产环境未配置则视为未认证通过（但 optional 允许继续）
			return func(c *gin.Context) { c.Next() }
		}
		secretKey = "gamelink-default-secret-key-change-in-production"
	}

	tokenDuration := 24 * time.Hour
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有提供Token，继续执行
			c.Next()
			return
		}

		// 尝试验证Token
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Token格式错误，继续执行（匿名访问）
			c.Next()
			return
		}

		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			// Token无效，继续执行（匿名访问）
			c.Next()
			return
		}

		if auth.IsTokenExpired(claims) {
			// Token过期，继续执行（匿名访问）
			c.Next()
			return
		}

        // Token有效，将用户信息存储到Context中
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Set("jwt_claims", claims)
        c.Set("is_authenticated", true)
        // 注入 actor
        c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), claims.UserID))

		c.Next()
	}
}

// GetUserID 从Context中获取用户ID
func GetUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint64)
	return id, ok
}

// GetUserRole 从Context中获取用户角色
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}

// IsAuthenticated 检查用户是否已认证
func IsAuthenticated(c *gin.Context) bool {
	isAuthenticated, exists := c.Get("is_authenticated")
	if !exists {
		return false
	}

	authenticated, ok := isAuthenticated.(bool)
	return ok && authenticated
}

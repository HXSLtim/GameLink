package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/handler/resp"
	"gamelink/pkg/apierr"
	"gamelink/pkg/auth"
	"gamelink/pkg/cache"
	"gamelink/pkg/logging"
)

// JWTAuth JWT认证中间件
//
// 使用方法：
// router.Use(middleware.JWTAuth(secretKey))
// 或者
// adminGroup.Use(middleware.JWTAuth(secretKey))
func JWTAuth(secretKey string) gin.HandlerFunc {
	// 验证密钥长度
	if len(secretKey) < 32 {
		logging.Error("JWT secret too short, must be at least 32 characters")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"code":    http.StatusServiceUnavailable,
				"message": "认证服务配置错误，请联系管理员",
			})
		}
	}

	// Token有效期（24小时）
	tokenDuration := auth.DefaultTokenDuration

	// 创建JWT管理器
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.Error(c, apierr.Unauthorized("缺少Authorization头"))
			c.Abort()
			return
		}

		// 提取Token
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			resp.Error(c, apierr.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		// 验证Token
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			resp.Error(c, apierr.Unauthorized("无效的Token: "+err.Error()))
			c.Abort()
			return
		}

		// 检查Token是否过期
		if auth.IsTokenExpired(claims) {
			resp.Error(c, apierr.Unauthorized("Token已过期"))
			c.Abort()
			return
		}

		// 将用户信息存储到Context中，供后续处理使用
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)
		// 注入 actor 到 request context，便于服务层审计日志使用
		c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), claims.UserID))

		// 检查Token剩余时间，如果快要过期，自动刷新Token
		remainingTime := auth.GetTokenRemainingTime(claims)

		// 如果Token即将过期（15分钟内），自动刷新
		if remainingTime < auth.TokenAutoRefreshWindow {
			newToken, err := jwtManager.RefreshToken(claims)
			if err == nil {
				// 在响应头中返回新Token
				c.Header("X-Refreshed-Token", newToken)

				// 更新Context中的Token信息
				newClaims, verifyErr := jwtManager.VerifyToken(newToken)
				if verifyErr != nil {
					logging.Warn("Failed to verify refreshed token", "error", verifyErr, "user_id", claims.UserID)
				} else if newClaims != nil {
					c.Set("jwt_claims", newClaims)
					c.Set("user_id", newClaims.UserID)
					c.Set("user_role", newClaims.Role)
				}

				logging.Debug("Token auto-refreshed", "user_id", claims.UserID)
			} else {
				logging.Warn("Failed to refresh token", "error", err, "user_id", claims.UserID)
			}
		} else if remainingTime < auth.TokenRefreshRecommendationWindow {
			// 仍然保留提示，让前端可以主动刷新
			c.Header("X-Token-Refresh-Recommendation", "true")
			c.Header("X-Token-Remaining", remainingTime.String())
		}

		c.Next()
	}
}

// WSAuth WebSocket 认证中间件，支持 header / query / cookie
func WSAuth(secretKey string) gin.HandlerFunc {
	// 验证密钥长度
	if len(secretKey) < 32 {
		logging.Error("JWT secret too short, must be at least 32 characters")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"code":    http.StatusServiceUnavailable,
				"message": "认证服务配置错误，请联系管理员",
			})
		}
	}

	tokenDuration := auth.DefaultTokenDuration
	jwtManager := auth.NewJWTManager(secretKey, tokenDuration)

	return func(c *gin.Context) {
		tokenString := resolveWSToken(c)
		if tokenString == "" {
			resp.Error(c, apierr.Unauthorized("缺少Authorization头"))
			c.Abort()
			return
		}

		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			resp.Error(c, apierr.Unauthorized("无效的Token: "+err.Error()))
			c.Abort()
			return
		}
		if auth.IsTokenExpired(claims) {
			resp.Error(c, apierr.Unauthorized("Token已过期"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)
		c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), claims.UserID))
		c.Next()
	}
}

func resolveWSToken(c *gin.Context) string {
	authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authHeader != "" {
		if token, err := auth.ExtractTokenFromHeader(authHeader); err == nil {
			return token
		}
		if token := trimBearer(authHeader); token != "" {
			return token
		}
	}

	if token := strings.TrimSpace(c.Query("token")); token != "" {
		return trimBearer(token)
	}
	if token := strings.TrimSpace(c.Query("access_token")); token != "" {
		return trimBearer(token)
	}
	if cookie, err := c.Cookie("auth_token"); err == nil {
		return strings.TrimSpace(cookie)
	}
	return ""
}

func trimBearer(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		return strings.TrimSpace(token[7:])
	}
	return token
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
			resp.Error(c, apierr.Forbidden("权限不足"))
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
func OptionalAuth(secretKey string) gin.HandlerFunc {
	// 验证密钥长度
	if len(secretKey) < 32 {
		logging.Error("JWT secret too short, must be at least 32 characters")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"code":    http.StatusServiceUnavailable,
				"message": "认证服务配置错误，请联系管理员",
			})
		}
	}

	tokenDuration := auth.DefaultTokenDuration
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

// JWTAuthWithRevocation 带Token撤销功能的JWT认证中间件
//
// 与JWTAuth的区别：
// - 使用Redis检查Token是否已被撤销
// - 支持用户登出后立即失效Token
//
// 使用方法：
// router.Use(middleware.JWTAuthWithRevocation(secretKey, cache))
func JWTAuthWithRevocation(secretKey string, c cache.Cache) gin.HandlerFunc {
	// 验证密钥长度
	if len(secretKey) < 32 {
		logging.Error("JWT secret too short, must be at least 32 characters")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"code":    http.StatusServiceUnavailable,
				"message": "认证服务配置错误，请联系管理员",
			})
		}
	}

	// Token有效期（24小时）
	tokenDuration := auth.DefaultTokenDuration

	// 创建带缓存的JWT管理器
	jwtManager := auth.NewJWTManagerWithCache(secretKey, tokenDuration, c)

	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.Error(c, apierr.Unauthorized("缺少Authorization头"))
			c.Abort()
			return
		}

		// 提取Token
		tokenString, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			resp.Error(c, apierr.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		// 验证Token并检查撤销状态
		claims, err := jwtManager.VerifyTokenWithRevocation(c.Request.Context(), tokenString)
		if err != nil {
			resp.Error(c, apierr.Unauthorized("无效的Token: "+err.Error()))
			c.Abort()
			return
		}

		// 检查Token是否过期
		if auth.IsTokenExpired(claims) {
			resp.Error(c, apierr.Unauthorized("Token已过期"))
			c.Abort()
			return
		}

		// 将用户信息存储到Context中
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("jwt_claims", claims)
		c.Set("session_id", claims.SessionID)
		c.Set("jti", claims.JTI)
		// 注入 actor 到 request context
		c.Request = c.Request.WithContext(logging.WithActorUserID(c.Request.Context(), claims.UserID))

		// Token自动刷新逻辑（与JWTAuth相同）
		remainingTime := auth.GetTokenRemainingTime(claims)
		if remainingTime < auth.TokenAutoRefreshWindow {
			newToken, err := jwtManager.RefreshToken(claims)
			if err == nil {
				c.Header("X-Refreshed-Token", newToken)
				newClaims, verifyErr := jwtManager.VerifyToken(newToken)
				if verifyErr == nil && newClaims != nil {
					c.Set("jwt_claims", newClaims)
					c.Set("user_id", newClaims.UserID)
					c.Set("user_role", newClaims.Role)
				}
				logging.Debug("Token auto-refreshed", "user_id", claims.UserID)
			}
		} else if remainingTime < auth.TokenRefreshRecommendationWindow {
			c.Header("X-Token-Refresh-Recommendation", "true")
			c.Header("X-Token-Remaining", remainingTime.String())
		}

		c.Next()
	}
}

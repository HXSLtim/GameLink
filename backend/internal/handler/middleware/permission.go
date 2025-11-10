package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	permissionservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"
)

const (
	// UserIDKey 在 Gin Context 中存储用户 ID 的键
	UserIDKey = "user_id"
	// UserRoleKey 在 Gin Context 中存储用户角色的键
	UserRoleKey = "user_role"
	// UserPermissionsKey 在 Gin Context 中存储用户权限的键
	UserPermissionsKey = "user_permissions"
)

// PermissionMiddleware 权限中间件配置。
type PermissionMiddleware struct {
	jwtManager    *auth.JWTManager
	permissionSvc *permissionservice.PermissionService
	roleSvc       *roleservice.RoleService
}

// NewPermissionMiddleware 创建权限中间件实例。
func NewPermissionMiddleware(
	jwtManager *auth.JWTManager,
	permissionSvc *permissionservice.PermissionService,
	roleSvc *roleservice.RoleService,
) *PermissionMiddleware {
	return &PermissionMiddleware{
		jwtManager:    jwtManager,
		permissionSvc: permissionSvc,
		roleSvc:       roleSvc,
	}
}

// RequireAuth 要求用户已登录（验证 JWT）。
func (m *PermissionMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.authenticateRequest(c) {
			return
		}
		c.Next()
	}
}

// authenticateRequest 负责解析并验证 JWT，将用户信息写入 context，成功返回 true。
func (m *PermissionMiddleware) authenticateRequest(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	token, err := auth.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": "未授权：" + err.Error(),
		})
		return false
	}

	claims, err := m.jwtManager.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": "Token 无效：" + err.Error(),
		})
		return false
	}

	if auth.IsTokenExpired(claims) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    http.StatusUnauthorized,
			"message": "Token 已过期",
		})
		return false
	}

	c.Set(UserIDKey, claims.UserID)
	c.Set(UserRoleKey, claims.Role)
	return true
}

// RequireRole 要求用户拥有指定角色（向后兼容，使用旧的 role 字段）。
func (m *PermissionMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.authenticateRequest(c) {
			return
		}

		// 获取用户角色
		role, exists := c.Get(UserRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "无法获取用户角色",
			})
			return
		}

		// 检查角色
		if role.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "权限不足：需要 " + requiredRole + " 角色",
			})
			return
		}

		c.Next()
	}
}

// RequirePermission 要求用户拥有指定权限（使用 method+path 或 code）。
// 注意：此中间件假设在 group 级别已经执行了 RequireAuth()，不会重复执行认证。
func (m *PermissionMiddleware) RequirePermission(method model.HTTPMethod, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户 ID（应该已经由 RequireAuth 设置）
		userID, exists := c.Get(UserIDKey)
		if !exists {
			// 如果没有用户信息，说明认证中间件没有执行，返回 401
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "未授权：请先登录",
			})
			return
		}

		uid := userID.(uint64)

		// 检查是否为超级管理员（拥有所有权限）
		isSuperAdmin, err := m.roleSvc.CheckUserIsSuperAdmin(c.Request.Context(), uid)
		if err == nil && isSuperAdmin {
			// 超级管理员放行
			c.Next()
			return
		}

		// 使用 method+path 检查权限
		hasPermission, checkErr := m.permissionSvc.CheckUserHasPermission(
			c.Request.Context(),
			uid,
			method,
			path,
		)

		if checkErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "权限检查失败",
			})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "权限不足",
			})
			return
		}

		c.Next()
	}
}

// RequireAnyRole 要求用户拥有任一指定角色。
func (m *PermissionMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.authenticateRequest(c) {
			return
		}

		// 获取用户 ID
		userID, exists := c.Get(UserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "无法获取用户信息",
			})
			return
		}

		uid := userID.(uint64)

		// 检查用户是否拥有任一角色
		for _, roleSlug := range roles {
			hasRole, err := m.roleSvc.CheckUserHasRole(c.Request.Context(), uid, roleSlug)
			if err == nil && hasRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    http.StatusForbidden,
			"message": "权限不足：需要以下角色之一：" + strings.Join(roles, ", "),
		})
	}
}

// 注意：GetUserID 和 GetUserRole 已在 jwt_auth.go 中定义

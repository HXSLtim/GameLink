package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	permissionservice "gamelink/internal/service/admin"
	roleservice "gamelink/internal/service/admin"
	"gamelink/pkg/auth"
)

// PermissionCheckMode 权限检查模式
type PermissionCheckMode string

const (
	// PermissionCheckModeAny 任一权限满足即可
	PermissionCheckModeAny PermissionCheckMode = "any"
	// PermissionCheckModeAll 所有权限都必须满足
	PermissionCheckModeAll PermissionCheckMode = "all"
	// PermissionCheckModeExcept 排除某些权限（用户不能拥有这些权限）
	PermissionCheckModeExcept PermissionCheckMode = "except"
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
	whitelist     map[string]bool // 白名单路径（无需权限检查）
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
		whitelist:     make(map[string]bool),
	}
}

// AddToWhitelist 添加路径到白名单（无需权限检查）。
// 格式：METHOD:PATH，如 "GET:/api/v1/health"
func (m *PermissionMiddleware) AddToWhitelist(methodPath string) {
	m.whitelist[methodPath] = true
}

// AddPathsToWhitelist 批量添加路径到白名单。
func (m *PermissionMiddleware) AddPathsToWhitelist(methodPaths []string) {
	for _, mp := range methodPaths {
		m.whitelist[mp] = true
	}
}

// RemoveFromWhitelist 从白名单移除路径。
func (m *PermissionMiddleware) RemoveFromWhitelist(methodPath string) {
	delete(m.whitelist, methodPath)
}

// IsWhitelisted 检查路径是否在白名单中。
func (m *PermissionMiddleware) IsWhitelisted(method, path string) bool {
	key := method + ":" + path
	return m.whitelist[key]
}

// ClearWhitelist 清空白名单。
func (m *PermissionMiddleware) ClearWhitelist() {
	m.whitelist = make(map[string]bool)
}

// GetWhitelist 获取当前白名单列表。
func (m *PermissionMiddleware) GetWhitelist() []string {
	result := make([]string, 0, len(m.whitelist))
	for k := range m.whitelist {
		result = append(result, k)
	}
	return result
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
// 支持从 Authorization header 或 query parameter (token) 获取 JWT。
func (m *PermissionMiddleware) authenticateRequest(c *gin.Context) bool {
	var token string
	var err error

	// 优先从 Authorization header 获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token, err = auth.ExtractTokenFromHeader(authHeader)
	}

	// 如果 header 中没有或提取失败，尝试从 query parameter 获取（支持 WebSocket）
	if token == "" {
		token = c.Query("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "未授权：缺少Authorization头",
			})
			return false
		}
		err = nil // 从 query 获取成功，清除之前的错误
	}

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
		// 检查白名单
		if m.IsWhitelisted(string(method), path) {
			c.Next()
			return
		}

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
			// 记录未授权访问日志
			log.Printf("Unauthorized access attempt: userID=%d, method=%s, path=%s, clientIP=%s",
				uid, method, path, c.ClientIP())

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

// RequirePermissionAuto 自动检查当前请求的权限（使用请求的 method+path）。
// 支持白名单和超级管理员绕过。
func (m *PermissionMiddleware) RequirePermissionAuto() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath() // 使用路由模板路径，如 /api/v1/admin/users/:id

		// 检查白名单
		if m.IsWhitelisted(method, path) {
			c.Next()
			return
		}

		// 获取用户 ID（应该已经由 RequireAuth 设置）
		userID, exists := c.Get(UserIDKey)
		if !exists {
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
			model.HTTPMethod(method),
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
			log.Printf("Unauthorized access attempt: userID=%d, method=%s, path=%s, clientIP=%s",
				uid, method, path, c.ClientIP())

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

// RequirePermissionCode 要求用户拥有指定权限码。
// 注意：此中间件假设在 group 级别已经执行了 RequireAuth()，不会重复执行认证。
func (m *PermissionMiddleware) RequirePermissionCode(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户 ID（应该已经由 RequireAuth 设置）
		userID, exists := c.Get(UserIDKey)
		if !exists {
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

		// 使用权限码检查权限
		hasPermission, checkErr := m.permissionSvc.CheckUserHasPermissionCode(
			c.Request.Context(),
			uid,
			code,
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
			log.Printf("Unauthorized access attempt: userID=%d, code=%s, clientIP=%s",
				uid, code, c.ClientIP())

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "权限不足：需要 " + code + " 权限",
			})
			return
		}

		c.Next()
	}
}

// RequirePermissionCodes 要求用户拥有指定权限码（支持多种检查模式）。
// mode: "any" - 任一权限满足即可, "all" - 所有权限都必须满足, "except" - 排除某些权限
func (m *PermissionMiddleware) RequirePermissionCodes(codes []string, mode PermissionCheckMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户 ID（应该已经由 RequireAuth 设置）
		userID, exists := c.Get(UserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "未授权：请先登录",
			})
			return
		}

		uid := userID.(uint64)

		// 检查是否为超级管理员（拥有所有权限）
		// 注意：except 模式下超级管理员也需要检查
		if mode != PermissionCheckModeExcept {
			isSuperAdmin, err := m.roleSvc.CheckUserIsSuperAdmin(c.Request.Context(), uid)
			if err == nil && isSuperAdmin {
				// 超级管理员放行
				c.Next()
				return
			}
		}

		var hasPermission bool
		var checkErr error

		switch mode {
		case PermissionCheckModeAny:
			hasPermission, checkErr = m.permissionSvc.CheckUserHasAnyPermission(c.Request.Context(), uid, codes)
		case PermissionCheckModeAll:
			hasPermission, checkErr = m.permissionSvc.CheckUserHasAllPermissions(c.Request.Context(), uid, codes)
		case PermissionCheckModeExcept:
			hasPermission, checkErr = m.permissionSvc.CheckUserHasExceptPermissions(c.Request.Context(), uid, codes)
		default:
			// 默认使用 any 模式
			hasPermission, checkErr = m.permissionSvc.CheckUserHasAnyPermission(c.Request.Context(), uid, codes)
		}

		if checkErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    http.StatusInternalServerError,
				"message": "权限检查失败",
			})
			return
		}

		if !hasPermission {
			var message string
			switch mode {
			case PermissionCheckModeAny:
				message = "权限不足：需要以下权限之一：" + strings.Join(codes, ", ")
			case PermissionCheckModeAll:
				message = "权限不足：需要以下所有权限：" + strings.Join(codes, ", ")
			case PermissionCheckModeExcept:
				message = "权限冲突：不能拥有以下权限：" + strings.Join(codes, ", ")
			default:
				message = "权限不足"
			}

			log.Printf("Unauthorized access attempt: userID=%d, codes=%v, mode=%s, clientIP=%s",
				uid, codes, mode, c.ClientIP())

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": message,
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermissionCode 要求用户拥有任一指定权限码（any 模式的便捷方法）。
func (m *PermissionMiddleware) RequireAnyPermissionCode(codes ...string) gin.HandlerFunc {
	return m.RequirePermissionCodes(codes, PermissionCheckModeAny)
}

// RequireAllPermissionCodes 要求用户拥有所有指定权限码（all 模式的便捷方法）。
func (m *PermissionMiddleware) RequireAllPermissionCodes(codes ...string) gin.HandlerFunc {
	return m.RequirePermissionCodes(codes, PermissionCheckModeAll)
}

// RequireExceptPermissionCodes 要求用户不拥有任何指定权限码（except 模式的便捷方法）。
func (m *PermissionMiddleware) RequireExceptPermissionCodes(codes ...string) gin.HandlerFunc {
	return m.RequirePermissionCodes(codes, PermissionCheckModeExcept)
}

// 注意：GetUserID 和 GetUserRole 已在 jwt_auth.go 中定义

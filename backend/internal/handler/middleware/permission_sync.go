package middleware

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/service"
)

// APISyncConfig API 同步配置。
type APISyncConfig struct {
	// GroupFilter 只同步匹配该前缀的路由分组，如 "/api/v1/admin"
	GroupFilter string
	// SkipPaths 跳过的路径列表
	SkipPaths []string
	// DryRun 是否为演练模式（不实际写入数据库）
	DryRun bool
}

// SyncAPIPermissions 同步 API 路由到权限表。
// 在应用启动后调用，遍历所有路由并注册到 permissions 表。
func SyncAPIPermissions(router *gin.Engine, permissionSvc *service.PermissionService, cfg APISyncConfig) error {
	routes := router.Routes()

	var permissions []model.Permission
	for _, route := range routes {
		// 跳过不匹配的分组
		if cfg.GroupFilter != "" && !strings.HasPrefix(route.Path, cfg.GroupFilter) {
			continue
		}

		// 跳过指定路径
		if shouldSkip(route.Path, cfg.SkipPaths) {
			continue
		}

		// 提取分组（取路径的第三级，如 /api/v1/admin/games -> /admin/games）
		group := extractGroup(route.Path)

		// 生成语义化 code（如 admin.games.list）
		code := generatePermissionCode(route.Method, route.Path)

		perm := model.Permission{
			Method:      model.HTTPMethod(route.Method),
			Path:        route.Path,
			Code:        code,
			Group:       group,
			Description: fmt.Sprintf("%s %s", route.Method, route.Path),
		}

		permissions = append(permissions, perm)
	}

	if cfg.DryRun {
		log.Printf("[DryRun] Would sync %d permissions:", len(permissions))
		for _, p := range permissions {
			log.Printf("  - [%s] %s (code: %s, group: %s)", p.Method, p.Path, p.Code, p.Group)
		}
		return nil
	}

	// 批量 upsert 权限
	ctx := context.Background()
	for _, perm := range permissions {
		// 创建局部副本，避免 range 变量复用问题
		p := perm
		if err := permissionSvc.UpsertPermission(ctx, &p); err != nil {
			log.Printf("Failed to upsert permission %s %s: %v", p.Method, p.Path, err)
			// 继续处理其他权限
		}
	}

	log.Printf("Synced %d API permissions to database", len(permissions))
	return nil
}

// shouldSkip 检查路径是否应该跳过。
func shouldSkip(path string, skipPaths []string) bool {
	for _, skip := range skipPaths {
		if strings.HasPrefix(path, skip) {
			return true
		}
	}
	return false
}

// extractGroup 提取路由分组。
// 例如：/api/v1/admin/games/:id -> /admin/games
func extractGroup(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return ""
	}

	// 移除 /api/v1 前缀
	if parts[1] == "api" && strings.HasPrefix(parts[2], "v") {
		parts = parts[3:]
	}

	// 找到第一个非参数部分作为分组
	var groupParts []string
	for _, part := range parts {
		if part == "" {
			continue
		}
		// 跳过路径参数（:id, :slug 等）
		if strings.HasPrefix(part, ":") {
			break
		}
		groupParts = append(groupParts, part)
	}

	if len(groupParts) == 0 {
		return ""
	}

	return "/" + strings.Join(groupParts, "/")
}

// generatePermissionCode 生成权限的语义化标识。
// 例如：GET /api/v1/admin/games -> admin.games.list
//
//	POST /api/v1/admin/games -> admin.games.create
//	GET /api/v1/admin/games/:id -> admin.games.read
//	PUT /api/v1/admin/games/:id -> admin.games.update
//	DELETE /api/v1/admin/games/:id -> admin.games.delete
func generatePermissionCode(method, path string) string {
	parts := strings.Split(path, "/")

	// 移除空字符串和 api/v1 前缀
	var cleanParts []string
	for i, part := range parts {
		if part == "" || part == "api" {
			continue
		}
		if strings.HasPrefix(part, "v") && i == 2 {
			continue
		}
		// 跳过路径参数
		if strings.HasPrefix(part, ":") {
			continue
		}
		cleanParts = append(cleanParts, part)
	}

	if len(cleanParts) == 0 {
		return ""
	}

	// 判断是否有 ID 参数（资源详情操作）
	hasIDParam := strings.Contains(path, "/:id") || strings.Contains(path, "/:") && strings.HasSuffix(path, "/:id")

	// 根据 HTTP 方法和是否有 ID 参数确定操作类型
	var action string
	switch method {
	case "GET":
		if hasIDParam {
			action = "read"
		} else {
			action = "list"
		}
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = strings.ToLower(method)
	}

	// 组合：resource.action，如 admin.games.list
	return strings.Join(cleanParts, ".") + "." + action
}

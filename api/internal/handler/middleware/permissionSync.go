package middleware

import (
	"fmt"
	"hash/crc32"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gamelink/internal/model"
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
// 使用 GORM 批量 ON CONFLICT upsert，将 N×2 次 DB 调用压缩为几次批量操作。
func SyncAPIPermissions(router *gin.Engine, db *gorm.DB, cfg APISyncConfig) error {
	t := time.Now()
	routes := router.Routes()
	existingCodeOwners, err := loadPermissionCodeOwners(db)
	if err != nil {
		return fmt.Errorf("load existing permission code owners failed: %w", err)
	}

	var permissions []model.Permission
	seen := make(map[string]struct{}, len(routes))
	generatedCodeOwners := make(map[string]string, len(routes))
	for _, route := range routes {
		if cfg.GroupFilter != "" && !strings.HasPrefix(route.Path, cfg.GroupFilter) {
			continue
		}
		if shouldSkip(route.Path, cfg.SkipPaths) {
			continue
		}

		key := route.Method + ":" + route.Path
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		group := extractGroup(route.Path)
		routeKey := buildRouteKey(route.Method, route.Path)
		code := resolveUniquePermissionCode(route.Method, route.Path, routeKey, existingCodeOwners, generatedCodeOwners)
		generatedCodeOwners[code] = routeKey

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

	if len(permissions) == 0 {
		return nil
	}

	// 使用 GORM Clauses 批量 upsert。
	// 关键：以 method+path 作为冲突键，避免历史上“同 method+path 不同 code”导致的唯一索引冲突。
	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "method"}, {Name: "path"}},
		DoUpdates: clause.AssignmentColumns([]string{"code", "group", "description", "updated_at"}),
	}).CreateInBatches(permissions, 100).Error
	if err != nil {
		return fmt.Errorf("permission batch upsert failed: %w", err)
	}

	log.Printf("[startup] synced %d API permissions in %v", len(permissions), time.Since(t))
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

func buildRouteKey(method, path string) string {
	return strings.ToUpper(method) + ":" + path
}

func resolveUniquePermissionCode(
	method, path, routeKey string,
	existingCodeOwners map[string]string,
	generatedCodeOwners map[string]string,
) string {
	base := generatePermissionCode(method, path)
	if base == "" {
		base = "api.permission"
	}

	if isPermissionCodeAvailable(base, routeKey, existingCodeOwners, generatedCodeOwners) {
		return base
	}

	suffix := fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(routeKey)))
	candidate := base + "." + suffix
	if isPermissionCodeAvailable(candidate, routeKey, existingCodeOwners, generatedCodeOwners) {
		return candidate
	}

	// 理论上极少命中，作为最终保险。
	for index := 1; ; index++ {
		candidate = fmt.Sprintf("%s.%s.%d", base, suffix, index)
		if isPermissionCodeAvailable(candidate, routeKey, existingCodeOwners, generatedCodeOwners) {
			return candidate
		}
	}
}

func isPermissionCodeAvailable(
	code, routeKey string,
	existingCodeOwners map[string]string,
	generatedCodeOwners map[string]string,
) bool {
	if owner, exists := existingCodeOwners[code]; exists && owner != routeKey {
		return false
	}
	if owner, exists := generatedCodeOwners[code]; exists && owner != routeKey {
		return false
	}
	return true
}

func loadPermissionCodeOwners(db *gorm.DB) (map[string]string, error) {
	type permissionCodeOwner struct {
		Code   string
		Method string
		Path   string
	}

	var rows []permissionCodeOwner
	if err := db.Model(&model.Permission{}).
		Select("code, method, path").
		Where("code <> ''").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	owners := make(map[string]string, len(rows))
	for _, row := range rows {
		owners[row.Code] = buildRouteKey(row.Method, row.Path)
	}
	return owners, nil
}

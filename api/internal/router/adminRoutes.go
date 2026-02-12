package router

import (
	"context"
	"fmt"
	"log"

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
)

// AssignDefaultRolePermissions 为默认角色分配管理权限（默认仅 super_admin 拥有全部权限）。
func AssignDefaultRolePermissions(ctx context.Context, roleSvc *adminservice.RoleService, permService *adminservice.PermissionService) error {
	// 获取所有权限
	allPermissions, err := permService.ListPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to list permissions: %w", err)
	}

	if len(allPermissions) == 0 {
		log.Println("没有权限需要分配，跳过")
		return nil
	}

	// 提取所有权限 ID
	permissionIDs := make([]uint64, 0, len(allPermissions))
	for _, perm := range allPermissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// 为 super_admin 和 admin 分配所有权限
	roleSlugs := []string{
		string(model.RoleSlugSuperAdmin),
		string(model.RoleSlugAdmin),
	}

	for _, roleSlug := range roleSlugs {
		role, err := roleSvc.GetRoleBySlug(ctx, roleSlug)
		if err != nil {
			log.Printf("警告：未找到角色 %s，跳过: %v", roleSlug, err)
			continue
		}

		// 分配权限（替换现有权限）
		if err := roleSvc.AssignPermissionsToRole(ctx, role.ID, permissionIDs); err != nil {
			log.Printf("警告：为角色 %s 分配权限失败: %v", roleSlug, err)
			continue
		}

		log.Printf("已为角色 %s (id=%d) 分配 %d 个权限", roleSlug, role.ID, len(permissionIDs))
	}

	// 为客服主管补齐争议处理权限（基于 method+path，兼容 API 自动同步生成的权限码）
	disputePathsForCSLeader := []struct {
		method model.HTTPMethod
		path   string
	}{
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes"},
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes/pending"},
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes/:id"},
		{method: model.HTTPMethodPOST, path: "/api/v1/admin/disputes/:id/assign"},
		{method: model.HTTPMethodPOST, path: "/api/v1/admin/disputes/:id/resolve"},
	}

	if err := ensureRolePermissionsByMethodPath(ctx, roleSvc, permService, allPermissions, string(model.RoleSlugCSLeader), disputePathsForCSLeader); err != nil {
		log.Printf("警告：补齐 csLeader 争议权限失败: %v", err)
	}

	// 兼容旧版客服角色 customerService
	if err := ensureRolePermissionsByMethodPath(ctx, roleSvc, permService, allPermissions, string(model.RoleSlugCustomerService), disputePathsForCSLeader); err != nil {
		log.Printf("警告：补齐 customerService 争议权限失败: %v", err)
	}

	// 为客服专员补齐基础争议处理权限（可查看与处理，不含分配与高敏操作）
	disputePathsForCSAgent := []struct {
		method model.HTTPMethod
		path   string
	}{
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes"},
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes/pending"},
		{method: model.HTTPMethodGET, path: "/api/v1/admin/disputes/:id"},
		{method: model.HTTPMethodPOST, path: "/api/v1/admin/disputes/:id/resolve"},
	}
	if err := ensureRolePermissionsByMethodPath(ctx, roleSvc, permService, allPermissions, string(model.RoleSlugCSAgent), disputePathsForCSAgent); err != nil {
		log.Printf("警告：补齐 csAgent 基础争议权限失败: %v", err)
	}

	return nil
}

func ensureRolePermissionsByMethodPath(
	ctx context.Context,
	roleSvc *adminservice.RoleService,
	permService *adminservice.PermissionService,
	allPermissions []model.Permission,
	roleSlug string,
	required []struct {
		method model.HTTPMethod
		path   string
	},
) error {
	role, err := roleSvc.GetRoleBySlug(ctx, roleSlug)
	if err != nil {
		return fmt.Errorf("role %s not found: %w", roleSlug, err)
	}

	existing, err := permService.ListPermissionsByRoleID(ctx, role.ID)
	if err != nil {
		return fmt.Errorf("list role permissions failed: %w", err)
	}

	existingIDs := make(map[uint64]struct{}, len(existing))
	for _, perm := range existing {
		existingIDs[perm.ID] = struct{}{}
	}

	pathKeyToPermissionID := make(map[string]uint64, len(allPermissions))
	for _, perm := range allPermissions {
		key := fmt.Sprintf("%s:%s", perm.Method, perm.Path)
		pathKeyToPermissionID[key] = perm.ID
	}

	missingIDs := make([]uint64, 0, len(required))
	for _, rp := range required {
		key := fmt.Sprintf("%s:%s", rp.method, rp.path)
		permID, exists := pathKeyToPermissionID[key]
		if !exists {
			log.Printf("警告：权限不存在，跳过 %s", key)
			continue
		}
		if _, ok := existingIDs[permID]; ok {
			continue
		}
		missingIDs = append(missingIDs, permID)
	}

	if len(missingIDs) == 0 {
		return nil
	}

	if err := roleSvc.AddPermissionsToRole(ctx, role.ID, missingIDs); err != nil {
		return fmt.Errorf("add permissions to role %s failed: %w", roleSlug, err)
	}

	log.Printf("已为角色 %s (id=%d) 补齐 %d 个权限", roleSlug, role.ID, len(missingIDs))
	return nil
}

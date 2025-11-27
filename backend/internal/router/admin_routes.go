package router

import (
	"context"
	"fmt"
	"log"

	"gamelink/internal/model"
	permissionservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"
)

// AssignDefaultRolePermissions 为默认角色分配管理权限（默认仅 super_admin 拥有全部权限）。
func AssignDefaultRolePermissions(ctx context.Context, roleSvc *roleservice.RoleService, permService *permissionservice.PermissionService) error {
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

	// 为 super_admin 分配所有权限；admin 不再自动获得全部权限，避免与超管等价
	roleSlugs := []string{
		string(model.RoleSlugSuperAdmin),
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

	return nil
}

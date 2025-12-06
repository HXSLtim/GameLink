package role

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
	"gamelink/pkg/cache"
)

const (
	cacheKeyPermissionsByUser = "rbac:user_permissions:%d"
	cacheKeyPermissionsByRole = "rbac:role_permissions:%d"
)

var (
	ErrValidation = service.ErrValidation
	ErrNotFound   = service.ErrNotFound
)

// RoleService 提供角色管理的业务逻辑。
type RoleService struct {
	roles           repository.RoleRepository
	cache           cache.Cache
	permissionCache *cache.PermissionCache
}

// NewRoleService 创建角色服务实例。
func NewRoleService(roles repository.RoleRepository, c cache.Cache) *RoleService {
	return &RoleService{
		roles:           roles,
		cache:           c,
		permissionCache: cache.NewPermissionCache(c),
	}
}

// GetPermissionCache returns the permission cache instance.
// This is useful for external components that need to interact with the cache.
func (s *RoleService) GetPermissionCache() *cache.PermissionCache {
	return s.permissionCache
}

const (
	cacheKeyRoles       = "admin:roles"
	cacheKeyRolesByUser = "admin:roles:user:%d"
	cacheTTLRoles       = 30 * time.Minute
)

// ListRoles 获取所有角色列表。
func (s *RoleService) ListRoles(ctx context.Context) ([]model.RoleModel, error) {
	return s.roles.List(ctx)
}

// ListRolesPaged 分页获取角色列表。
func (s *RoleService) ListRolesPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	return s.roles.ListPaged(ctx, page, pageSize)
}

// ListRolesPagedWithFilter 分页获取角色列表（带过滤）。
func (s *RoleService) ListRolesPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	return s.roles.ListPagedWithFilter(ctx, page, pageSize, keyword, isSystem)
}

// ListRolesWithPermissions 获取角色列表，预加载权限。
func (s *RoleService) ListRolesWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	return s.roles.ListWithPermissions(ctx)
}

// GetRole 根据ID获取角色。
func (s *RoleService) GetRole(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return s.roles.Get(ctx, id)
}

// GetRoleWithPermissions 根据ID获取角色，预加载权限。
func (s *RoleService) GetRoleWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return s.roles.GetWithPermissions(ctx, id)
}

// GetRoleBySlug 根据Slug获取角色。
func (s *RoleService) GetRoleBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	return s.roles.GetBySlug(ctx, slug)
}

// CreateRole 创建角色。
func (s *RoleService) CreateRole(ctx context.Context, role *model.RoleModel) error {
	// 校验必填字段
	if role.Slug == "" || role.Name == "" {
		return fmt.Errorf("%w: slug and name are required", ErrValidation)
	}

	// 检查 slug 是否已存在
	existing, err := s.roles.GetBySlug(ctx, role.Slug)
	if err == nil && existing != nil {
		return fmt.Errorf("%w: role with slug %s already exists", ErrValidation, role.Slug)
	}

	// 创建角色
	if err := s.roles.Create(ctx, role); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateRoleCache()
	return nil
}

// UpdateRole 更新角色。
func (s *RoleService) UpdateRole(ctx context.Context, role *model.RoleModel) error {
	if role.ID == 0 {
		return fmt.Errorf("%w: role ID is required", ErrValidation)
	}

	// 检查是否为系统角色
	existing, err := s.roles.Get(ctx, role.ID)
	if err != nil {
		return err
	}

	if existing.IsSystem {
		// 系统角色只允许更新描述
		existing.Description = role.Description
		role = existing
	}

	if err := s.roles.Update(ctx, role); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateRoleCache()
	return nil
}

// DeleteRole 删除角色（系统角色不可删除）。
func (s *RoleService) DeleteRole(ctx context.Context, id uint64) error {
	if err := s.roles.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateRoleCache()
	return nil
}

// AssignPermissionsToRole 为角色分配权限（替换现有权限）。
func (s *RoleService) AssignPermissionsToRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if err := s.roles.AssignPermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 清除缓存并传播到所有拥有该角色的用户
	s.invalidateRoleCache()
	_ = s.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID)
	return nil
}

// AddPermissionsToRole 为角色添加权限（追加）。
func (s *RoleService) AddPermissionsToRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if err := s.roles.AddPermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 清除缓存并传播到所有拥有该角色的用户
	s.invalidateRoleCache()
	_ = s.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID)
	return nil
}

// RemovePermissionsFromRole 移除角色的权限。
func (s *RoleService) RemovePermissionsFromRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if err := s.roles.RemovePermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 清除缓存并传播到所有拥有该角色的用户
	s.invalidateRoleCache()
	_ = s.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID)
	return nil
}

// ListRolesByUserID 获取用户拥有的所有角色。
func (s *RoleService) ListRolesByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	cacheKey := fmt.Sprintf(cacheKeyRolesByUser, userID)

	// 尝试从缓存获取
	if value, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var roles []model.RoleModel
		if err := json.Unmarshal([]byte(value), &roles); err == nil {
			return roles, nil
		}
	}

	// 从数据库获取
	roles, err := s.roles.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(roles); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(data), cacheTTLRoles)
	}
	return roles, nil
}

// AssignRolesToUser 为用户分配角色（替换现有角色）。
func (s *RoleService) AssignRolesToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	if err := s.roles.AssignToUser(ctx, userID, roleIDs); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateUserRoleCache(userID)
	return nil
}

// RemoveRolesFromUser 移除用户的角色。
func (s *RoleService) RemoveRolesFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	if err := s.roles.RemoveFromUser(ctx, userID, roleIDs); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateUserRoleCache(userID)
	return nil
}

// CheckUserHasRole 检查用户是否拥有指定角色。
func (s *RoleService) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	return s.roles.CheckUserHasRole(ctx, userID, roleSlug)
}

// CheckUserIsSuperAdmin 检查用户是否为超级管理员。
func (s *RoleService) CheckUserIsSuperAdmin(ctx context.Context, userID uint64) (bool, error) {
	return s.roles.CheckUserHasRole(ctx, userID, string(model.RoleSlugSuperAdmin))
}

// invalidateRoleCache 清除角色相关缓存。
func (s *RoleService) invalidateRoleCache() {
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyRoles)
}

// invalidateUserRoleCache 清除用户角色缓存。
func (s *RoleService) invalidateUserRoleCache(userID uint64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cacheKeyRolesByUser, userID)
	_ = s.cache.Delete(ctx, cacheKey)

	// 同时清除用户权限缓存（使用新的权限缓存）
	_ = s.permissionCache.InvalidateUserCache(ctx, userID)
}

// invalidatePermissionCacheForRole 清除角色权限缓存。
func (s *RoleService) invalidatePermissionCacheForRole(roleID uint64) {
	ctx := context.Background()
	_ = s.permissionCache.InvalidateRoleCache(ctx, roleID)
}

// InvalidateRolePermissionsAndPropagateToUsers invalidates the role's permission cache
// and propagates the invalidation to all users who have that role.
// This should be called when role permissions are modified.
func (s *RoleService) InvalidateRolePermissionsAndPropagateToUsers(ctx context.Context, roleID uint64) error {
	return s.permissionCache.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID, s.roles)
}

// SetRoleParent sets the parent role for inheritance.
// It validates the inheritance depth and checks for circular inheritance.
func (s *RoleService) SetRoleParent(ctx context.Context, roleID uint64, parentID *uint64) error {
	// Validate that role exists
	role, err := s.roles.Get(ctx, roleID)
	if err != nil {
		return err
	}

	// If setting a parent, validate it
	if parentID != nil && *parentID > 0 {
		// Check that parent exists
		parent, err := s.roles.GetWithPermissions(ctx, *parentID)
		if err != nil {
			return fmt.Errorf("%w: parent role not found", ErrNotFound)
		}

		// Check for circular inheritance
		if err := s.ValidateNoCircularInheritance(ctx, roleID, *parentID); err != nil {
			return err
		}

		// Check max depth
		newLevel := parent.Level + 1
		if newLevel > model.MaxRoleInheritanceDepth {
			return model.ErrRoleMaxDepthExceeded
		}

		// Check if setting parent would cause children to exceed max depth
		if err := s.validateChildrenDepth(ctx, roleID, newLevel); err != nil {
			return err
		}
	}

	// Set the parent
	if err := s.roles.SetParent(ctx, roleID, parentID); err != nil {
		return err
	}

	// Invalidate caches
	s.invalidateRoleCache()
	s.invalidatePermissionCacheForRole(roleID)

	// Also invalidate cache for all users with this role
	s.invalidateUsersWithRoleCache(ctx, role.ID)

	return nil
}

// validateChildrenDepth checks if setting a new level would cause any children to exceed max depth.
func (s *RoleService) validateChildrenDepth(ctx context.Context, roleID uint64, newLevel int) error {
	children, err := s.roles.GetChildRoles(ctx, roleID)
	if err != nil {
		return err
	}

	for _, child := range children {
		childNewLevel := newLevel + 1
		if childNewLevel > model.MaxRoleInheritanceDepth {
			return model.ErrRoleMaxDepthExceeded
		}

		// Recursively check grandchildren
		if err := s.validateChildrenDepth(ctx, child.ID, childNewLevel); err != nil {
			return err
		}
	}

	return nil
}

// GetRoleInheritanceChain returns the inheritance chain from the given role up to the root.
// The chain is ordered from the given role to the root (child -> parent -> grandparent).
func (s *RoleService) GetRoleInheritanceChain(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	return s.roles.GetInheritanceChain(ctx, roleID)
}

// ValidateNoCircularInheritance checks if setting parentID as the parent of roleID would create a cycle.
func (s *RoleService) ValidateNoCircularInheritance(ctx context.Context, roleID uint64, parentID uint64) error {
	// A role cannot be its own parent
	if roleID == parentID {
		return model.ErrRoleCircularInheritance
	}

	// Get the inheritance chain of the proposed parent
	chain, err := s.roles.GetInheritanceChain(ctx, parentID)
	if err != nil {
		return err
	}

	// Check if roleID appears in the parent's inheritance chain
	for _, ancestor := range chain {
		if ancestor.ID == roleID {
			return model.ErrRoleCircularInheritance
		}
	}

	return nil
}

// GetEffectivePermissions returns all permissions for a role, including inherited permissions.
// Permissions are merged with child role permissions taking priority over parent permissions.
func (s *RoleService) GetEffectivePermissions(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	// Get the inheritance chain (from child to root)
	chain, err := s.roles.GetInheritanceChain(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Collect permissions from all roles in the chain
	// Process from root to child so child permissions override parent permissions
	permissionMap := make(map[uint64]model.Permission)

	// Reverse the chain to process from root to child
	for i := len(chain) - 1; i >= 0; i-- {
		role := chain[i]
		roleWithPerms, err := s.roles.GetWithPermissions(ctx, role.ID)
		if err != nil {
			continue
		}

		for _, perm := range roleWithPerms.Permissions {
			permissionMap[perm.ID] = perm
		}
	}

	// Convert map to slice
	permissions := make([]model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// GetChildRoles returns all direct child roles of the given role.
func (s *RoleService) GetChildRoles(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	return s.roles.GetChildRoles(ctx, roleID)
}

// invalidateUsersWithRoleCache invalidates the permission cache for all users with the given role.
func (s *RoleService) invalidateUsersWithRoleCache(ctx context.Context, roleID uint64) {
	// Get all users with this role
	role, err := s.roles.GetWithPermissions(ctx, roleID)
	if err != nil {
		return
	}

	// Invalidate cache for each user
	for _, user := range role.Users {
		s.invalidateUserRoleCache(user.ID)
	}
}

// GetUserEffectivePermissions returns all permissions for a user, merging permissions from all assigned roles.
// This includes inherited permissions from role hierarchies.
// Permissions are merged using union (all permissions from all roles).
// When roles have different priorities, higher priority role permissions are preferred.
func (s *RoleService) GetUserEffectivePermissions(ctx context.Context, userID uint64) ([]model.Permission, error) {
	// Get all roles for the user
	roles, err := s.ListRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Sort roles by priority (higher priority first)
	sortRolesByPriority(roles)

	// Collect all permissions from all roles (including inherited)
	permissionMap := make(map[uint64]model.Permission)

	for _, role := range roles {
		// Get effective permissions for this role (including inherited)
		rolePerms, err := s.GetEffectivePermissions(ctx, role.ID)
		if err != nil {
			continue
		}

		// Add to map (later roles with same permission ID will override)
		for _, perm := range rolePerms {
			// Only add if not already present (first role with higher priority wins)
			if _, exists := permissionMap[perm.ID]; !exists {
				permissionMap[perm.ID] = perm
			}
		}
	}

	// Convert map to slice
	permissions := make([]model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// sortRolesByPriority sorts roles by priority in descending order (higher priority first).
func sortRolesByPriority(roles []model.RoleModel) {
	for i := 0; i < len(roles)-1; i++ {
		for j := i + 1; j < len(roles); j++ {
			if roles[j].Priority > roles[i].Priority {
				roles[i], roles[j] = roles[j], roles[i]
			}
		}
	}
}

// MergePermissions merges multiple permission slices into one, removing duplicates.
// Permissions are identified by their ID.
func MergePermissions(permissionSets ...[]model.Permission) []model.Permission {
	permissionMap := make(map[uint64]model.Permission)

	for _, perms := range permissionSets {
		for _, perm := range perms {
			permissionMap[perm.ID] = perm
		}
	}

	result := make([]model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		result = append(result, perm)
	}

	return result
}

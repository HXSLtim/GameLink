package role

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
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
	roles repository.RoleRepository
	cache cache.Cache
}

// NewRoleService 创建角色服务实例。
func NewRoleService(roles repository.RoleRepository, cache cache.Cache) *RoleService {
	return &RoleService{
		roles: roles,
		cache: cache,
	}
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
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.roles.ListPaged(ctx, page, pageSize)
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

	// 清除缓存
	s.invalidateRoleCache()
	s.invalidatePermissionCacheForRole(roleID)
	return nil
}

// AddPermissionsToRole 为角色添加权限（追加）。
func (s *RoleService) AddPermissionsToRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if err := s.roles.AddPermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateRoleCache()
	s.invalidatePermissionCacheForRole(roleID)
	return nil
}

// RemovePermissionsFromRole 移除角色的权限。
func (s *RoleService) RemovePermissionsFromRole(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if err := s.roles.RemovePermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	// 清除缓存
	s.invalidateRoleCache()
	s.invalidatePermissionCacheForRole(roleID)
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

	// 同时清除用户权限缓存
	permCacheKey := fmt.Sprintf(cacheKeyPermissionsByUser, userID)
	_ = s.cache.Delete(ctx, permCacheKey)
}

// invalidatePermissionCacheForRole 清除角色权限缓存。
func (s *RoleService) invalidatePermissionCacheForRole(roleID uint64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(cacheKeyPermissionsByRole, roleID)
	_ = s.cache.Delete(ctx, cacheKey)
}

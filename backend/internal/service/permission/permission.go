package permission

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

var (
	ErrValidation = service.ErrValidation
	ErrNotFound   = service.ErrNotFound
)

// PermissionService 提供权限管理的业务逻辑。
type PermissionService struct {
	permissions repository.PermissionRepository
	cache       cache.Cache
}

// NewPermissionService 创建权限服务实例。
func NewPermissionService(permissions repository.PermissionRepository, cache cache.Cache) *PermissionService {
	return &PermissionService{
		permissions: permissions,
		cache:       cache,
	}
}

const (
	cacheKeyPermissions       = "admin:permissions"
	cacheKeyPermissionsByRole = "admin:permissions:role:%d"
	cacheKeyPermissionsByUser = "admin:permissions:user:%d"
	cacheTTLPermissions       = 30 * time.Minute
)

// ListPermissions 获取所有权限列表。
func (s *PermissionService) ListPermissions(ctx context.Context) ([]model.Permission, error) {
	return s.permissions.List(ctx)
}

// ListPermissionsPaged 分页获取权限列表。
func (s *PermissionService) ListPermissionsPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	page = repository.NormalizePage(page)
	pageSize = repository.NormalizePageSize(pageSize)
	return s.permissions.ListPaged(ctx, page, pageSize)
}

// ListPermissionsByGroup 按分组获取权限。
func (s *PermissionService) ListPermissionsByGroup(ctx context.Context, group string) ([]model.Permission, error) {
	grouped, err := s.permissions.ListByGroup(ctx)
	if err != nil {
		return nil, err
	}
	return grouped[group], nil
}

// GetPermission 根据ID获取权限。
func (s *PermissionService) GetPermission(ctx context.Context, id uint64) (*model.Permission, error) {
	return s.permissions.Get(ctx, id)
}

// CreatePermission 创建权限。
func (s *PermissionService) CreatePermission(ctx context.Context, permission *model.Permission) error {
	// 校验必填字段
	if permission.Method == "" || permission.Path == "" {
		return fmt.Errorf("%w: method and path are required", ErrValidation)
	}

	// 校验权限码格式（如果提供了权限码）
	if permission.Code != "" {
		if !permission.ValidateCode() {
			return ErrPermissionCodeInvalid
		}

		// 检查权限码是否已存在
		existing, err := s.permissions.GetByCode(ctx, permission.Code)
		if err == nil && existing != nil {
			return ErrPermissionCodeExists
		}
	}

	// 检查 method+path 是否已存在
	existing, err := s.permissions.GetByMethodAndPath(ctx, string(permission.Method), permission.Path)
	if err == nil && existing != nil {
		return fmt.Errorf("%w: permission with method %s and path %s already exists", ErrValidation, permission.Method, permission.Path)
	}

	// 创建权限
	if err := s.permissions.Create(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// UpdatePermission 更新权限。
// Note: Permission code cannot be modified after creation.
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	if permission.ID == 0 {
		return fmt.Errorf("%w: permission ID is required", ErrValidation)
	}

	// Get existing permission to check code immutability
	existing, err := s.permissions.Get(ctx, permission.ID)
	if err != nil {
		return err
	}

	// Prevent code modification if the existing permission has a code
	if existing.Code != "" && permission.Code != "" && existing.Code != permission.Code {
		return ErrPermissionCodeImmutable
	}

	// If setting a new code (existing was empty), validate format
	if existing.Code == "" && permission.Code != "" {
		if !permission.ValidateCode() {
			return ErrPermissionCodeInvalid
		}
		// Check if the new code already exists
		existingByCode, err := s.permissions.GetByCode(ctx, permission.Code)
		if err == nil && existingByCode != nil && existingByCode.ID != permission.ID {
			return ErrPermissionCodeExists
		}
	}

	if err := s.permissions.Update(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// DeletePermission 删除权限。
// Returns error if:
// - Permission is a system permission (IsSystem=true)
// - Permission is referenced by any roles
func (s *PermissionService) DeletePermission(ctx context.Context, id uint64) error {
	// Get the permission first to check constraints
	permission, err := s.permissions.Get(ctx, id)
	if err != nil {
		return err
	}

	// Check if it's a system permission
	if permission.IsSystem {
		return ErrPermissionIsSystem
	}

	// Check if the permission is referenced by any roles
	refCount, err := s.permissions.CountRoleReferences(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check role references: %w", err)
	}
	if refCount > 0 {
		return ErrPermissionInUse.WithDetails(fmt.Sprintf("该权限被 %d 个角色引用", refCount))
	}

	if err := s.permissions.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// DeletePermissionForce 强制删除权限（跳过引用检查，但仍不能删除系统权限）。
// This should only be used in special cases like cleanup operations.
func (s *PermissionService) DeletePermissionForce(ctx context.Context, id uint64) error {
	// Get the permission first to check if it's a system permission
	permission, err := s.permissions.Get(ctx, id)
	if err != nil {
		return err
	}

	// System permissions can never be deleted
	if permission.IsSystem {
		return ErrPermissionIsSystem
	}

	if err := s.permissions.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// CanDeletePermission checks if a permission can be deleted.
// Returns (canDelete, reason, error)
func (s *PermissionService) CanDeletePermission(ctx context.Context, id uint64) (bool, string, error) {
	permission, err := s.permissions.Get(ctx, id)
	if err != nil {
		return false, "", err
	}

	if permission.IsSystem {
		return false, "系统权限不可删除", nil
	}

	refCount, err := s.permissions.CountRoleReferences(ctx, id)
	if err != nil {
		return false, "", fmt.Errorf("failed to check role references: %w", err)
	}
	if refCount > 0 {
		return false, fmt.Sprintf("该权限被 %d 个角色引用", refCount), nil
	}

	return true, "", nil
}

// UpsertPermission 根据 method+path 存在则更新，不存在则创建。
func (s *PermissionService) UpsertPermission(ctx context.Context, permission *model.Permission) error {
	if err := s.permissions.UpsertByMethodPath(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// ListPermissionsByRoleID 获取指定角色拥有的所有权限。
func (s *PermissionService) ListPermissionsByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	cacheKey := fmt.Sprintf(cacheKeyPermissionsByRole, roleID)

	// 尝试从缓存获取
	if value, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var permissions []model.Permission
		if err := json.Unmarshal([]byte(value), &permissions); err == nil {
			return permissions, nil
		}
	}

	// 从数据库获取
	permissions, err := s.permissions.ListByRoleID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(permissions); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(data), cacheTTLPermissions)
	}
	return permissions, nil
}

// ListPermissionsByUserID 获取指定用户拥有的所有权限（通过角色）。
func (s *PermissionService) ListPermissionsByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	cacheKey := fmt.Sprintf(cacheKeyPermissionsByUser, userID)

	// 尝试从缓存获取
	if value, ok, err := s.cache.Get(ctx, cacheKey); err == nil && ok {
		var permissions []model.Permission
		if err := json.Unmarshal([]byte(value), &permissions); err == nil {
			return permissions, nil
		}
	}

	// 从数据库获取
	permissions, err := s.permissions.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if data, err := json.Marshal(permissions); err == nil {
		_ = s.cache.Set(ctx, cacheKey, string(data), cacheTTLPermissions)
	}
	return permissions, nil
}

// CheckUserHasPermission 检查用户是否拥有指定权限。
func (s *PermissionService) CheckUserHasPermission(ctx context.Context, userID uint64, method model.HTTPMethod, path string) (bool, error) {
	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Method == method && perm.Path == path {
			return true, nil
		}
	}

	return false, nil
}

// ListPermissionGroups 获取所有权限分组列表。
func (s *PermissionService) ListPermissionGroups(ctx context.Context) ([]string, error) {
	return s.permissions.ListGroups(ctx)
}

// ValidatePermissionCode validates a permission code format.
// Returns nil if valid, or an error describing the validation failure.
// Valid format: module.resource.action (three dot-separated lowercase segments)
func (s *PermissionService) ValidatePermissionCode(code string) error {
	if code == "" {
		return fmt.Errorf("%w: permission code is required", ErrValidation)
	}
	if !model.PermissionCodePattern.MatchString(code) {
		return ErrPermissionCodeInvalid
	}
	return nil
}

// CheckPermissionCodeExists checks if a permission code already exists.
func (s *PermissionService) CheckPermissionCodeExists(ctx context.Context, code string) (bool, error) {
	existing, err := s.permissions.GetByCode(ctx, code)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return existing != nil, nil
}

// invalidatePermissionCache 清除权限相关缓存。
func (s *PermissionService) invalidatePermissionCache() {
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyPermissions)
	_ = s.cache.Delete(ctx, cacheKeyPermissionTree)
	// 注意：用户和角色的权限缓存需要在分配权限时清除
}

const (
	cacheKeyPermissionTree = "admin:permissions:tree"
)

// GetPermissionTree returns all permissions organized in a tree structure.
// Uses caching to avoid repeated database queries.
// The tree is organized by Group, with parent-child relationships preserved.
func (s *PermissionService) GetPermissionTree(ctx context.Context) ([]*model.PermissionTreeNode, error) {
	// Try to get from cache
	if value, ok, err := s.cache.Get(ctx, cacheKeyPermissionTree); err == nil && ok {
		var tree []*model.PermissionTreeNode
		if err := json.Unmarshal([]byte(value), &tree); err == nil {
			return tree, nil
		}
	}

	// Fetch all permissions with efficient ordering
	permissions, err := s.permissions.ListWithChildren(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	// Build tree structure
	tree := model.BuildPermissionTree(permissions)

	// Cache the result
	if data, err := json.Marshal(tree); err == nil {
		_ = s.cache.Set(ctx, cacheKeyPermissionTree, string(data), cacheTTLPermissions)
	}

	return tree, nil
}

// GetPermissionTreeByGroup returns all permissions organized in a tree structure grouped by permission group.
func (s *PermissionService) GetPermissionTreeByGroup(ctx context.Context) ([]model.PermissionGroup, error) {
	// Fetch all permissions with efficient ordering
	permissions, err := s.permissions.ListWithChildren(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	// Build tree structure grouped by group
	return model.BuildPermissionTreeByGroup(permissions), nil
}

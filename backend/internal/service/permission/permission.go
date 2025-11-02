package permission

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
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
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
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	if permission.ID == 0 {
		return fmt.Errorf("%w: permission ID is required", ErrValidation)
	}

	if err := s.permissions.Update(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// DeletePermission 删除权限。
func (s *PermissionService) DeletePermission(ctx context.Context, id uint64) error {
	if err := s.permissions.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
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

// invalidatePermissionCache 清除权限相关缓存。
func (s *PermissionService) invalidatePermissionCache() {
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyPermissions)
	// 注意：用户和角色的权限缓存需要在分配权限时清除
}


package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/apierr"
	"gamelink/pkg/cache"
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

// PermissionListOptions 权限列表查询选项
type PermissionListOptions struct {
	Page     int
	PageSize int
	Keyword  string // 搜索关键词（匹配code, path, description）
	Method   string // HTTP方法过滤
	Group    string // 分组过滤
	IsSystem *bool  // 系统权限过滤
}

// ListPermissionsPagedWithFilter 分页获取权限列表（支持过滤）。
func (s *PermissionService) ListPermissionsPagedWithFilter(ctx context.Context, opts PermissionListOptions) ([]model.Permission, int64, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}
	return s.permissions.ListPagedWithFilter(ctx, opts.Page, opts.PageSize, opts.Keyword, opts.Method, opts.Group, opts.IsSystem)
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
		return apierr.BadRequest("method and path are required")
	}

	// 校验权限码格式（如果提供了权限码）
	if permission.Code != "" && !permission.ValidateCode() {
		return apierr.BadRequest("权限码格式无效，应为 module.resource.action")
	}

	// 检查 method+path 是否已存在
	existing, err := s.permissions.GetByMethodAndPath(ctx, string(permission.Method), permission.Path)
	if err == nil && existing != nil {
		return apierr.BadRequest("权限已存在: " + string(permission.Method) + " " + permission.Path)
	}

	// 检查权限码是否已存在
	if permission.Code != "" {
		existingByCode, err := s.permissions.GetByCode(ctx, permission.Code)
		if err == nil && existingByCode != nil {
			return apierr.BadRequest("权限码已存在: " + permission.Code)
		}
	}

	// 创建权限
	if err := s.permissions.Create(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// UpdatePermission 更新权限（全量更新，但权限码不可修改）。
func (s *PermissionService) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	if permission.ID == 0 {
		return apierr.BadRequest("permission ID is required")
	}

	// 获取现有权限
	existing, err := s.permissions.Get(ctx, permission.ID)
	if err != nil {
		return err
	}

	// 权限码创建后不可修改（Requirements 1.4）
	if permission.Code != "" && existing.Code != "" && permission.Code != existing.Code {
		return apierr.BadRequest("权限码创建后不可修改")
	}

	// 如果是新设置权限码，验证格式
	if permission.Code != "" && existing.Code == "" && !permission.ValidateCode() {
		return apierr.BadRequest("权限码格式无效，应为 module.resource.action")
	}

	// 保留原有的权限码（如果新值为空）
	if permission.Code == "" {
		permission.Code = existing.Code
	}

	if err := s.permissions.Update(ctx, permission); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// PartialUpdatePermission 部分更新权限（只更新非空字段）。
func (s *PermissionService) PartialUpdatePermission(ctx context.Context, id uint64, updates map[string]interface{}) (*model.Permission, error) {
	// 获取现有权限
	existing, err := s.permissions.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// 权限码创建后不可修改（Requirements 1.4）
	if code, ok := updates["code"]; ok {
		codeStr, _ := code.(string)
		if codeStr != "" && existing.Code != "" && codeStr != existing.Code {
			return nil, apierr.BadRequest("权限码创建后不可修改")
		}
		// 如果是新设置权限码，验证格式
		if codeStr != "" && existing.Code == "" {
			testPerm := &model.Permission{Code: codeStr}
			if !testPerm.ValidateCode() {
				return nil, apierr.BadRequest("权限码格式无效，应为 module.resource.action")
			}
		}
	}

	// 应用更新
	if group, ok := updates["group"]; ok {
		existing.Group = group.(string)
	}
	if description, ok := updates["description"]; ok {
		existing.Description = description.(string)
	}
	if code, ok := updates["code"]; ok {
		codeStr, _ := code.(string)
		if codeStr != "" {
			existing.Code = codeStr
		}
	}
	if sortOrder, ok := updates["sortOrder"]; ok {
		existing.SortOrder = int(sortOrder.(float64))
	}

	if err := s.permissions.Update(ctx, existing); err != nil {
		return nil, err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return existing, nil
}

// DeletePermission 删除权限（软删除）。
// 检查是否被角色引用，系统权限不可删除（Requirements 1.5）。
func (s *PermissionService) DeletePermission(ctx context.Context, id uint64) error {
	// 获取权限信息
	permission, err := s.permissions.Get(ctx, id)
	if err != nil {
		return err
	}

	// 系统权限不可删除
	if permission.IsSystem {
		return apierr.BadRequest("系统权限不可删除")
	}

	// 检查是否被角色引用
	refCount, err := s.permissions.CountRoleReferences(ctx, id)
	if err != nil {
		return err
	}
	if refCount > 0 {
		return apierr.BadRequest(fmt.Sprintf("权限被 %d 个角色引用，无法删除", refCount))
	}

	if err := s.permissions.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// DeletePermissionForce 强制删除权限（忽略引用检查，用于管理员确认后删除）。
func (s *PermissionService) DeletePermissionForce(ctx context.Context, id uint64) error {
	// 获取权限信息
	permission, err := s.permissions.Get(ctx, id)
	if err != nil {
		return err
	}

	// 系统权限不可删除
	if permission.IsSystem {
		return apierr.BadRequest("系统权限不可删除")
	}

	if err := s.permissions.Delete(ctx, id); err != nil {
		return err
	}

	// 清除缓存
	s.invalidatePermissionCache()
	return nil
}

// GetPermissionReferenceCount 获取权限被角色引用的数量。
func (s *PermissionService) GetPermissionReferenceCount(ctx context.Context, id uint64) (int64, error) {
	return s.permissions.CountRoleReferences(ctx, id)
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

// CheckUserHasPermissionCode 检查用户是否拥有指定权限码。
func (s *PermissionService) CheckUserHasPermissionCode(ctx context.Context, userID uint64, code string) (bool, error) {
	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Code == code {
			return true, nil
		}
	}

	return false, nil
}

// CheckUserHasAnyPermission 检查用户是否拥有任一指定权限码（any 模式）。
func (s *PermissionService) CheckUserHasAnyPermission(ctx context.Context, userID uint64, codes []string) (bool, error) {
	if len(codes) == 0 {
		return true, nil
	}

	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Build a set of user permission codes for efficient lookup
	userCodes := make(map[string]bool)
	for _, perm := range permissions {
		userCodes[perm.Code] = true
	}

	// Check if user has any of the required codes
	for _, code := range codes {
		if userCodes[code] {
			return true, nil
		}
	}

	return false, nil
}

// CheckUserHasAllPermissions 检查用户是否拥有所有指定权限码（all 模式）。
func (s *PermissionService) CheckUserHasAllPermissions(ctx context.Context, userID uint64, codes []string) (bool, error) {
	if len(codes) == 0 {
		return true, nil
	}

	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Build a set of user permission codes for efficient lookup
	userCodes := make(map[string]bool)
	for _, perm := range permissions {
		userCodes[perm.Code] = true
	}

	// Check if user has all of the required codes
	for _, code := range codes {
		if !userCodes[code] {
			return false, nil
		}
	}

	return true, nil
}

// CheckUserHasExceptPermissions 检查用户是否不拥有任何指定权限码（except 模式）。
// 返回 true 表示用户没有任何被排除的权限。
func (s *PermissionService) CheckUserHasExceptPermissions(ctx context.Context, userID uint64, excludedCodes []string) (bool, error) {
	if len(excludedCodes) == 0 {
		return true, nil
	}

	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Build a set of excluded codes for efficient lookup
	excludedSet := make(map[string]bool)
	for _, code := range excludedCodes {
		excludedSet[code] = true
	}

	// Check if user has any of the excluded codes
	for _, perm := range permissions {
		if excludedSet[perm.Code] {
			return false, nil
		}
	}

	return true, nil
}

// GetUserPermissionCodes 获取用户的所有权限码列表。
func (s *PermissionService) GetUserPermissionCodes(ctx context.Context, userID uint64) ([]string, error) {
	permissions, err := s.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		if perm.Code != "" {
			codes = append(codes, perm.Code)
		}
	}

	return codes, nil
}

// ListPermissionGroups 获取所有权限分组列表。
func (s *PermissionService) ListPermissionGroups(ctx context.Context) ([]string, error) {
	return s.permissions.ListGroups(ctx)
}

// GetPermissionTree 获取权限树形结构。
// 返回按父子关系组织的权限树。
func (s *PermissionService) GetPermissionTree(ctx context.Context) ([]*model.PermissionTreeNode, error) {
	permissions, err := s.permissions.ListWithChildren(ctx)
	if err != nil {
		return nil, err
	}
	return model.BuildPermissionTree(permissions), nil
}

// GetPermissionTreeByGroup 获取按分组组织的权限树形结构。
// 返回按分组分类的权限树。
func (s *PermissionService) GetPermissionTreeByGroup(ctx context.Context) ([]model.PermissionGroup, error) {
	permissions, err := s.permissions.ListWithChildren(ctx)
	if err != nil {
		return nil, err
	}
	return model.BuildPermissionTreeByGroup(permissions), nil
}

// invalidatePermissionCache 清除权限相关缓存。
func (s *PermissionService) invalidatePermissionCache() {
	ctx := context.Background()
	_ = s.cache.Delete(ctx, cacheKeyPermissions)
	// 注意：用户和角色的权限缓存需要在分配权限时清除
}

// PermissionBatchDeleteResult 权限批量删除结果
type PermissionBatchDeleteResult struct {
	SuccessCount int                       `json:"successCount"`
	FailedCount  int                       `json:"failedCount"`
	FailedPerms  []PermissionDeleteFailure `json:"failedPerms,omitempty"`
}

// PermissionDeleteFailure 权限删除失败详情
type PermissionDeleteFailure struct {
	PermissionID uint64 `json:"permissionId"`
	Reason       string `json:"reason"`
}

// BatchDeletePermissions 批量删除权限（系统权限不可删除，被引用的权限需要确认）
func (s *PermissionService) BatchDeletePermissions(ctx context.Context, ids []uint64, force bool) (*PermissionBatchDeleteResult, error) {
	if len(ids) == 0 {
		return nil, apierr.BadRequest("权限ID列表不能为空")
	}

	result := &PermissionBatchDeleteResult{
		FailedPerms: make([]PermissionDeleteFailure, 0),
	}

	for _, id := range ids {
		// 获取权限信息
		permission, err := s.permissions.Get(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedPerms = append(result.FailedPerms, PermissionDeleteFailure{
				PermissionID: id,
				Reason:       "权限不存在",
			})
			continue
		}

		// 系统权限不可删除
		if permission.IsSystem {
			result.FailedCount++
			result.FailedPerms = append(result.FailedPerms, PermissionDeleteFailure{
				PermissionID: id,
				Reason:       "系统权限不可删除",
			})
			continue
		}

		// 检查是否被角色引用
		if !force {
			refCount, err := s.permissions.CountRoleReferences(ctx, id)
			if err != nil {
				result.FailedCount++
				result.FailedPerms = append(result.FailedPerms, PermissionDeleteFailure{
					PermissionID: id,
					Reason:       "检查权限引用失败",
				})
				continue
			}
			if refCount > 0 {
				result.FailedCount++
				result.FailedPerms = append(result.FailedPerms, PermissionDeleteFailure{
					PermissionID: id,
					Reason:       fmt.Sprintf("权限被 %d 个角色引用，无法删除", refCount),
				})
				continue
			}
		}

		if err := s.permissions.Delete(ctx, id); err != nil {
			result.FailedCount++
			result.FailedPerms = append(result.FailedPerms, PermissionDeleteFailure{
				PermissionID: id,
				Reason:       err.Error(),
			})
		} else {
			result.SuccessCount++
		}
	}

	// 清除缓存
	if result.SuccessCount > 0 {
		s.invalidatePermissionCache()
	}

	return result, nil
}

// BatchDeletePermissionsWithResponse 批量删除权限（返回统一的BatchOperationResponse格式）
func (s *PermissionService) BatchDeletePermissionsWithResponse(ctx context.Context, ids []uint64, force bool) (*BatchDeletePermissionsResult, error) {
	if len(ids) == 0 {
		return nil, apierr.BadRequest("权限ID列表不能为空")
	}

	result := &BatchDeletePermissionsResult{
		SuccessItems: make([]uint64, 0),
		FailedItems:  make([]BatchDeletePermissionError, 0),
		TotalCount:   len(ids),
	}

	for _, id := range ids {
		// 获取权限信息
		permission, err := s.permissions.Get(ctx, id)
		if err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchDeletePermissionError{
				ID:      id,
				Message: "权限不存在",
			})
			continue
		}

		// 系统权限不可删除
		if permission.IsSystem {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchDeletePermissionError{
				ID:      id,
				Message: "系统权限不可删除",
			})
			continue
		}

		// 检查是否被角色引用
		if !force {
			refCount, err := s.permissions.CountRoleReferences(ctx, id)
			if err != nil {
				result.FailedCount++
				result.FailedItems = append(result.FailedItems, BatchDeletePermissionError{
					ID:      id,
					Message: "检查权限引用失败",
				})
				continue
			}
			if refCount > 0 {
				result.FailedCount++
				result.FailedItems = append(result.FailedItems, BatchDeletePermissionError{
					ID:      id,
					Message: fmt.Sprintf("权限被 %d 个角色引用，无法删除", refCount),
				})
				continue
			}
		}

		if err := s.permissions.Delete(ctx, id); err != nil {
			result.FailedCount++
			result.FailedItems = append(result.FailedItems, BatchDeletePermissionError{
				ID:      id,
				Message: err.Error(),
			})
		} else {
			result.SuccessCount++
			result.SuccessItems = append(result.SuccessItems, id)
		}
	}

	// 清除缓存
	if result.SuccessCount > 0 {
		s.invalidatePermissionCache()
	}

	return result, nil
}

// BatchDeletePermissionsResult 批量删除权限结果（统一格式）
type BatchDeletePermissionsResult struct {
	SuccessCount int                          `json:"success_count"`
	FailedCount  int                          `json:"failed_count"`
	TotalCount   int                          `json:"total_count"`
	FailedItems  []BatchDeletePermissionError `json:"failed_items,omitempty"`
	SuccessItems []uint64                     `json:"success_items,omitempty"`
}

// BatchDeletePermissionError 单个权限删除失败详情
type BatchDeletePermissionError struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

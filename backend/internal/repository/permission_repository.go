package repository

import (
	"context"

	"gamelink/internal/model"
)

// PermissionRepository 定义权限的数据访问操作。
type PermissionRepository interface {
	// List 获取所有权限列表
	List(ctx context.Context) ([]model.Permission, error)

	// ListPaged 分页获取权限列表
	ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error)

	// ListByGroup 按分组获取权限
	ListByGroup(ctx context.Context, group string) ([]model.Permission, error)

	// ListGroups 获取所有权限分组列表
	ListGroups(ctx context.Context) ([]string, error)

	// Get 根据ID获取权限
	Get(ctx context.Context, id uint64) (*model.Permission, error)

	// GetByMethodAndPath 根据 Method+Path 获取权限
	GetByMethodAndPath(ctx context.Context, method model.HTTPMethod, path string) (*model.Permission, error)

	// GetByCode 根据语义化Code获取权限
	GetByCode(ctx context.Context, code string) (*model.Permission, error)

	// Create 创建权限
	Create(ctx context.Context, permission *model.Permission) error

	// CreateBatch 批量创建权限（用于API自动注册）
	CreateBatch(ctx context.Context, permissions []model.Permission) error

	// Update 更新权限
	Update(ctx context.Context, permission *model.Permission) error

	// Delete 删除权限
	Delete(ctx context.Context, id uint64) error

	// UpsertByMethodPath 根据 method+path 存在则更新，不存在则创建
	UpsertByMethodPath(ctx context.Context, permission *model.Permission) error

	// ListByRoleID 获取指定角色拥有的所有权限
	ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error)

	// ListByUserID 获取指定用户拥有的所有权限（通过角色）
	ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error)
}

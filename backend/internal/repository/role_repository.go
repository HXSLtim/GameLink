package repository

import (
	"context"

	"gamelink/internal/model"
)

// RoleRepository 定义角色的数据访问操作。
type RoleRepository interface {
	// List 获取所有角色列表
	List(ctx context.Context) ([]model.RoleModel, error)

	// ListPaged 分页获取角色列表
	ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error)

	// ListWithPermissions 获取角色列表，预加载权限
	ListWithPermissions(ctx context.Context) ([]model.RoleModel, error)

	// Get 根据ID获取角色
	Get(ctx context.Context, id uint64) (*model.RoleModel, error)

	// GetWithPermissions 根据ID获取角色，预加载权限
	GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error)

	// GetBySlug 根据Slug获取角色
	GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error)

	// Create 创建角色
	Create(ctx context.Context, role *model.RoleModel) error

	// Update 更新角色
	Update(ctx context.Context, role *model.RoleModel) error

	// Delete 删除角色（系统角色不可删除）
	Delete(ctx context.Context, id uint64) error

	// AssignPermissions 为角色分配权限（替换现有权限）
	AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// AddPermissions 为角色添加权限（追加）
	AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// RemovePermissions 移除角色的权限
	RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error

	// ListByUserID 获取用户拥有的所有角色
	ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error)

	// AssignToUser 为用户分配角色
	AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error

	// RemoveFromUser 移除用户的角色
	RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error

	// CheckUserHasRole 检查用户是否拥有指定角色
	CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error)
}


package admin

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例。
func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) List(ctx context.Context) ([]model.RoleModel, error) {
	var roles []model.RoleModel
	err := r.db.WithContext(ctx).Order("is_system DESC, slug").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	var roles []model.RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RoleModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("is_system DESC, slug").
		Offset(offset).
		Limit(pageSize).
		Find(&roles).Error

	return roles, total, err
}

// ListPagedWithFilter 分页获取角色列表（支持过滤）
func (r *roleRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	var roles []model.RoleModel
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RoleModel{})

	// 关键词搜索（匹配name或slug）
	if keyword != "" {
		query = query.Where("name LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 系统角色过滤
	if isSystem != nil {
		query = query.Where("is_system = ?", *isSystem)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("is_system DESC, slug").
		Offset(offset).
		Limit(pageSize).
		Find(&roles).Error

	return roles, total, err
}

func (r *roleRepository) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	var roles []model.RoleModel
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		Order("is_system DESC, slug").
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	var role model.RoleModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &role, err
}

func (r *roleRepository) Create(ctx context.Context, role *model.RoleModel) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) Update(ctx context.Context, role *model.RoleModel) error {
	result := r.db.WithContext(ctx).Model(role).Updates(role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id uint64) error {
	// 检查是否为系统角色
	var role model.RoleModel
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return err
	}

	if role.IsSystem {
		return errors.New("cannot delete system role")
	}

	result := r.db.WithContext(ctx).Delete(&model.RoleModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *roleRepository) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有权限关联
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// 添加新的权限关联
		if len(permissionIDs) > 0 {
			rolePermissions := make([]model.RolePermission, len(permissionIDs))
			for i, permID := range permissionIDs {
				rolePermissions[i] = model.RolePermission{
					RoleID:       roleID,
					PermissionID: permID,
				}
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roleRepository) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	rolePermissions := make([]model.RolePermission, len(permissionIDs))
	for i, permID := range permissionIDs {
		rolePermissions[i] = model.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
	}

	// 使用事务批量插入，忽略已存在的记录
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, rp := range rolePermissions {
			// 检查是否存在
			var existing model.RolePermission
			err := tx.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
				First(&existing).Error
			if err == nil {
				// 已存在，跳过
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			// 不存在，创建
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepository) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&model.RolePermission{}).Error
}

func (r *roleRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	var roles []model.RoleModel
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.is_system DESC, roles.slug").
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除现有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 添加新的角色关联
		if len(roleIDs) > 0 {
			userRoles := make([]model.UserRole, len(roleIDs))
			for i, roleID := range roleIDs {
				userRoles[i] = model.UserRole{
					UserID: userID,
					RoleID: roleID,
				}
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *roleRepository) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&model.UserRole{}).Error
}

func (r *roleRepository) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.slug = ?", userID, roleSlug).
		Count(&count).Error
	return count > 0, err
}

// SetParent sets the parent role for a given role and updates the level.
func (r *roleRepository) SetParent(ctx context.Context, roleID uint64, parentID *uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get the current role
		var role model.RoleModel
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return repository.ErrNotFound
			}
			return err
		}

		// Calculate new level
		newLevel := 0
		if parentID != nil && *parentID > 0 {
			var parent model.RoleModel
			if err := tx.First(&parent, *parentID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return repository.ErrNotFound
				}
				return err
			}
			newLevel = parent.Level + 1
		}

		// Check max depth
		if newLevel > model.MaxRoleInheritanceDepth {
			return model.ErrRoleMaxDepthExceeded
		}

		// Update the role
		updates := map[string]interface{}{
			"parent_id": parentID,
			"level":     newLevel,
		}
		if err := tx.Model(&role).Updates(updates).Error; err != nil {
			return err
		}

		// Update all child roles' levels recursively
		return r.updateChildLevels(tx, roleID, newLevel)
	})
}

// updateChildLevels recursively updates the level of all child roles.
func (r *roleRepository) updateChildLevels(tx *gorm.DB, parentID uint64, parentLevel int) error {
	var children []model.RoleModel
	if err := tx.Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}

	for _, child := range children {
		newLevel := parentLevel + 1
		if newLevel > model.MaxRoleInheritanceDepth {
			return model.ErrRoleMaxDepthExceeded
		}

		if err := tx.Model(&child).Update("level", newLevel).Error; err != nil {
			return err
		}

		// Recursively update grandchildren
		if err := r.updateChildLevels(tx, child.ID, newLevel); err != nil {
			return err
		}
	}

	return nil
}

// GetInheritanceChain returns the inheritance chain from the given role up to the root.
// The chain is ordered from the given role to the root (child -> parent -> grandparent).
func (r *roleRepository) GetInheritanceChain(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	var chain []model.RoleModel

	// Get the starting role
	var role model.RoleModel
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	chain = append(chain, role)

	// Traverse up the inheritance chain
	currentRole := &role
	visited := make(map[uint64]bool)
	visited[role.ID] = true

	for currentRole.ParentID != nil && *currentRole.ParentID > 0 {
		parentID := *currentRole.ParentID

		// Check for circular reference (should not happen if validation is correct)
		if visited[parentID] {
			return nil, model.ErrRoleCircularInheritance
		}

		var parent model.RoleModel
		if err := r.db.WithContext(ctx).First(&parent, parentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break // Parent not found, stop traversal
			}
			return nil, err
		}

		chain = append(chain, parent)
		visited[parentID] = true
		currentRole = &parent
	}

	return chain, nil
}

// GetChildRoles returns all direct child roles of the given role.
func (r *roleRepository) GetChildRoles(ctx context.Context, roleID uint64) ([]model.RoleModel, error) {
	var children []model.RoleModel
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", roleID).
		Order("priority DESC, slug").
		Find(&children).Error
	return children, err
}

// UpdateLevel updates the level of a role.
func (r *roleRepository) UpdateLevel(ctx context.Context, roleID uint64, level int) error {
	result := r.db.WithContext(ctx).
		Model(&model.RoleModel{}).
		Where("id = ?", roleID).
		Update("level", level)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// GetUserIDsByRoleID returns all user IDs that have the specified role.
// This is used for cache invalidation propagation when role permissions change.
func (r *roleRepository) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	var userIDs []uint64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

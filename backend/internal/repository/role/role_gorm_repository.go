package role

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

// NewRoleRepository åå»ºè§è²ä»å¨å®ä¾ã?
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
	// æ£æ¥æ¯å¦ä¸ºç³»ç»è§è²
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
		// å é¤ç°ææéå³è
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}

		// æ·»å æ°çæéå³è
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

	// ä½¿ç¨äºå¡æ¹éæå¥ï¼å¿½ç¥å·²å­å¨çè®°å½?
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, rp := range rolePermissions {
			// æ£æ¥æ¯å¦å·²å­å¨
			var existing model.RolePermission
			err := tx.Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
				First(&existing).Error
			if err == nil {
				// å·²å­å¨ï¼è·³è¿
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			// ä¸å­å¨ï¼åå»º
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
		// å é¤ç°æè§è²å³è
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// æ·»å æ°çè§è²å³è
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

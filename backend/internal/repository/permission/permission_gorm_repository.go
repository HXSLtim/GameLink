package permission

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository åå»ºæéä»å¨å®ä¾ã?
func NewPermissionRepository(db *gorm.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) List(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).Order("group, method, path").Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Permission{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("group, method, path").
		Offset(offset).
		Limit(pageSize).
		Find(&permissions).Error

	return permissions, total, err
}

func (r *permissionRepository) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).
		Order("\"group\", method, path").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	// Group permissions by their group field
	grouped := make(map[string][]model.Permission)
	for _, p := range permissions {
		grouped[p.Group] = append(grouped[p.Group], p)
	}
	return grouped, nil
}

func (r *permissionRepository) ListGroups(ctx context.Context) ([]string, error) {
	var groups []string
	err := r.db.WithContext(ctx).
		Model(&model.Permission{}).
		Distinct("\"group\"").
		Where("\"group\" != ''").
		Order("\"group\"").
		Pluck("group", &groups).Error
	return groups, err
}

func (r *permissionRepository) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).First(&permission, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &permission, err
}

func (r *permissionRepository) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).
		Where("resource = ? AND action = ?", resource, action).
		First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &permission, err
}

func (r *permissionRepository) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).
		Where("method = ? AND path = ?", method, path).
		First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &permission, err
}

func (r *permissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	return &permission, err
}

func (r *permissionRepository) Create(ctx context.Context, permission *model.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *permissionRepository) CreateBatch(ctx context.Context, permissions []model.Permission) error {
	if len(permissions) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&permissions).Error
}

func (r *permissionRepository) Update(ctx context.Context, permission *model.Permission) error {
	result := r.db.WithContext(ctx).Model(permission).Updates(permission)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Permission{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *permissionRepository) UpsertByMethodPath(ctx context.Context, permission *model.Permission) error {
	// å°è¯æ¥æ¾ç°æè®°å½
	var existing model.Permission
	err := r.db.WithContext(ctx).
		Where("method = ? AND path = ?", permission.Method, permission.Path).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// ä¸å­å¨ï¼åå»ºæ°è®°å½?
		return r.db.WithContext(ctx).Create(permission).Error
	}
	if err != nil {
		return err
	}

	// å­å¨ï¼æ´æ°è®°å½ï¼ä¿ç IDï¼?
	permission.ID = existing.ID
	return r.db.WithContext(ctx).Model(&existing).Updates(permission).Error
}

func (r *permissionRepository) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.group, permissions.method, permissions.path").
		Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.WithContext(ctx).
		Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Order("permissions.group, permissions.method, permissions.path").
		Find(&permissions).Error
	return permissions, err
}

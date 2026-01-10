package menu

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

type repositoryImpl struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) repository.MenuRepository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Create(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *repositoryImpl) Update(ctx context.Context, menu *model.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *repositoryImpl) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Menu{}, id).Error
}

func (r *repositoryImpl) Get(ctx context.Context, id uint64) (*model.Menu, error) {
	var m model.Menu
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *repositoryImpl) List(ctx context.Context, parentID *uint64) ([]model.Menu, error) {
	tx := r.db.WithContext(ctx).Model(&model.Menu{})
	if parentID != nil {
		tx = tx.Where("parent_id = ?", *parentID)
	}
	var menus []model.Menu
	if err := tx.Order("\"order\" ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *repositoryImpl) ListPaged(ctx context.Context, page, pageSize int, parentID *uint64) ([]model.Menu, int64, error) {
	tx := r.db.WithContext(ctx).Model(&model.Menu{})
	if parentID != nil {
		tx = tx.Where("parent_id = ?", *parentID)
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		tx = tx.Offset(offset).Limit(pageSize)
	}

	var menus []model.Menu
	if err := tx.Order("\"order\" ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, 0, err
	}
	return menus, total, nil
}

func (r *repositoryImpl) ListByPermission(ctx context.Context, codes []string) ([]model.Menu, error) {
	if len(codes) == 0 {
		return []model.Menu{}, nil
	}
	var menus []model.Menu
	if err := r.db.WithContext(ctx).
		Where("permission IN ?", codes).
		Order("\"order\" ASC, id ASC").
		Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *repositoryImpl) HasChildren(ctx context.Context, parentID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Menu{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

package serviceitem

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// ServiceItemRepository 服务项目仓储接口（统一管理所有服务类型，包括礼物�?
type ServiceItemRepository interface {
	// 基础CRUD
	Create(ctx context.Context, item *model.ServiceItem) error
	Get(ctx context.Context, id uint64) (*model.ServiceItem, error)
	GetByCode(ctx context.Context, itemCode string) (*model.ServiceItem, error)
	List(ctx context.Context, opts ServiceItemListOptions) ([]model.ServiceItem, int64, error)
	Update(ctx context.Context, item *model.ServiceItem) error
	Delete(ctx context.Context, id uint64) error

	// 批量操作
	BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error
	BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error

	// 特殊查询
	GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error)
	GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error)
}

// ServiceItemListOptions 服务项目查询选项
type ServiceItemListOptions struct {
	Category    *string
	SubCategory *model.ServiceItemSubCategory
	GameID      *uint64
	PlayerID    *uint64
	IsActive    *bool
	Page        int
	PageSize    int
}

type serviceItemRepository struct {
	db *gorm.DB
}

// NewServiceItemRepository 创建服务项目仓储
func NewServiceItemRepository(db *gorm.DB) ServiceItemRepository {
	return &serviceItemRepository{db: db}
}

// Create 创建服务项目
func (r *serviceItemRepository) Create(ctx context.Context, item *model.ServiceItem) error {
	desiredActive := item.IsActive
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return err
	}
	if !desiredActive {
		if err := r.db.WithContext(ctx).
			Model(item).
			Update("is_active", false).Error; err != nil {
			return err
		}
		item.IsActive = false
	}
	return nil
}

// Get 获取服务项目
func (r *serviceItemRepository) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	var item model.ServiceItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// GetByCode 根据编码获取服务项目
func (r *serviceItemRepository) GetByCode(ctx context.Context, itemCode string) (*model.ServiceItem, error) {
	var item model.ServiceItem
	err := r.db.WithContext(ctx).Where("item_code = ?", itemCode).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &item, nil
}

// List 查询服务项目列表
func (r *serviceItemRepository) List(ctx context.Context, opts ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.ServiceItem{})

	// 过滤条件
	if opts.Category != nil {
		query = query.Where("category = ?", *opts.Category)
	}
	if opts.SubCategory != nil {
		query = query.Where("sub_category = ?", *opts.SubCategory)
	}
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}
	if opts.PlayerID != nil {
		query = query.Where("player_id = ?", *opts.PlayerID)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	offset := (opts.Page - 1) * opts.PageSize

	// 查询数据
	var items []model.ServiceItem
	err := query.Order("sort_order ASC, created_at DESC").
		Offset(offset).Limit(opts.PageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Update 更新服务项目
func (r *serviceItemRepository) Update(ctx context.Context, item *model.ServiceItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete 删除服务项目
func (r *serviceItemRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.ServiceItem{}, id).Error
}

// BatchUpdateStatus 批量更新状�?
func (r *serviceItemRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	return r.db.WithContext(ctx).
		Model(&model.ServiceItem{}).
		Where("id IN ?", ids).
		Update("is_active", isActive).Error
}

// BatchUpdatePrice 批量更新价格
func (r *serviceItemRepository) BatchUpdatePrice(ctx context.Context, ids []uint64, basePriceCents int64) error {
	return r.db.WithContext(ctx).
		Model(&model.ServiceItem{}).
		Where("id IN ?", ids).
		Update("base_price_cents", basePriceCents).Error
}

// GetGifts 获取礼物列表
func (r *serviceItemRepository) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	subCat := model.SubCategoryGift
	return r.List(ctx, ServiceItemListOptions{
		SubCategory: &subCat,
		IsActive:    boolPtr(true),
		Page:        page,
		PageSize:    pageSize,
	})
}

// GetGameServices 获取指定游戏的服�?
func (r *serviceItemRepository) GetGameServices(ctx context.Context, gameID uint64, subCategory *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	query := r.db.WithContext(ctx).
		Where("game_id = ? AND is_active = ?", gameID, true)

	if subCategory != nil {
		query = query.Where("sub_category = ?", *subCategory)
	}

	var items []model.ServiceItem
	err := query.Order("sort_order ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func boolPtr(b bool) *bool {
	return &b
}


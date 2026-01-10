package ordergroup

import (
	"context"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository OrderGroup 仓储接口
type Repository interface {
	Create(ctx context.Context, group *model.OrderGroup) error
	Get(ctx context.Context, id uint64) (*model.OrderGroup, error)
	GetByGroupNo(ctx context.Context, groupNo string) (*model.OrderGroup, error)
	GetWithSubOrders(ctx context.Context, id uint64) (*model.OrderGroup, error)
	Update(ctx context.Context, group *model.OrderGroup) error
	UpdateStatus(ctx context.Context, id uint64, status model.OrderGroupStatus) error
	List(ctx context.Context, opts ListOptions) ([]model.OrderGroup, int64, error)
	ListByUser(ctx context.Context, userID uint64, opts ListOptions) ([]model.OrderGroup, int64, error)
}

// ListOptions 列表查询选项
type ListOptions struct {
	UserID   *uint64
	Status   *model.OrderGroupStatus
	GameID   *uint64
	Page     int
	PageSize int
}

type orderGroupRepository struct {
	db *gorm.DB
}

// NewRepository 创建 OrderGroup 仓储
func NewRepository(db *gorm.DB) Repository {
	return &orderGroupRepository{db: db}
}

// Create 创建主订单
func (r *orderGroupRepository) Create(ctx context.Context, group *model.OrderGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// Get 获取主订单
func (r *orderGroupRepository) Get(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	var group model.OrderGroup
	err := r.db.WithContext(ctx).First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

// GetByGroupNo 根据订单号获取主订单
func (r *orderGroupRepository) GetByGroupNo(ctx context.Context, groupNo string) (*model.OrderGroup, error) {
	var group model.OrderGroup
	err := r.db.WithContext(ctx).Where("group_no = ?", groupNo).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

// GetWithSubOrders 获取主订单及其子订单
func (r *orderGroupRepository) GetWithSubOrders(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	var group model.OrderGroup
	err := r.db.WithContext(ctx).
		Preload("SubOrders", func(db *gorm.DB) *gorm.DB {
			return db.Order("hour_index ASC")
		}).
		Preload("SubOrders.Player").
		First(&group, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

// Update 更新主订单
func (r *orderGroupRepository) Update(ctx context.Context, group *model.OrderGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// UpdateStatus 更新主订单状态
func (r *orderGroupRepository) UpdateStatus(ctx context.Context, id uint64, status model.OrderGroupStatus) error {
	return r.db.WithContext(ctx).
		Model(&model.OrderGroup{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// List 查询主订单列表
func (r *orderGroupRepository) List(ctx context.Context, opts ListOptions) ([]model.OrderGroup, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.OrderGroup{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.GameID != nil {
		query = query.Where("game_id = ?", *opts.GameID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 || opts.PageSize > 100 {
		opts.PageSize = 20
	}

	var groups []model.OrderGroup
	err := query.
		Order("created_at DESC").
		Offset((opts.Page - 1) * opts.PageSize).
		Limit(opts.PageSize).
		Find(&groups).Error

	return groups, total, err
}

// ListByUser 查询用户的主订单列表
func (r *orderGroupRepository) ListByUser(ctx context.Context, userID uint64, opts ListOptions) ([]model.OrderGroup, int64, error) {
	opts.UserID = &userID
	return r.List(ctx, opts)
}

package collectionentity

import (
	"context"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"gorm.io/gorm"
)

// CollectionEntityRepository 收款主体仓储接口
// Requirements: 15.3, 16.5
type CollectionEntityRepository interface {
	// Entity CRUD operations
	Create(ctx context.Context, entity *model.CollectionEntity) error
	Get(ctx context.Context, id uint64) (*model.CollectionEntity, error)
	GetByCreditCode(ctx context.Context, creditCode string) (*model.CollectionEntity, error)
	GetDefault(ctx context.Context) (*model.CollectionEntity, error)
	Update(ctx context.Context, entity *model.CollectionEntity) error
	List(ctx context.Context, opts ListOptions) ([]model.CollectionEntity, int64, error)
	ListActive(ctx context.Context) ([]model.CollectionEntity, error) // Requirements: 17.5 - 获取所有活跃主体用于容错
	ToggleStatus(ctx context.Context, id uint64, status model.EntityStatus) error
	SetDefault(ctx context.Context, id uint64) error

	// Payment channel operations
	CreateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error
	GetChannel(ctx context.Context, id uint64) (*model.PaymentChannelConfig, error)
	GetChannelByEntityAndMethod(ctx context.Context, entityID uint64, method model.PaymentMethod) (*model.PaymentChannelConfig, error)
	UpdateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error
	DeleteChannel(ctx context.Context, id uint64) error
	ListChannelsByEntity(ctx context.Context, entityID uint64) ([]model.PaymentChannelConfig, error)

	// History operations
	CreateHistory(ctx context.Context, history *model.CollectionEntityHistory) error
	GetHistory(ctx context.Context, entityID uint64) ([]model.CollectionEntityHistory, error)

	// Statistics
	UpdateCollectionStats(ctx context.Context, entityID uint64, amountCents int64) error
	GetCollectionStats(ctx context.Context, entityID uint64) (totalCents int64, count int64, err error)
}

// ListOptions 收款主体查询选项
type ListOptions struct {
	Status    *model.EntityStatus
	Keyword   string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
}

type collectionEntityRepository struct {
	db *gorm.DB
}

// NewCollectionEntityRepository 创建收款主体仓储
func NewCollectionEntityRepository(db *gorm.DB) CollectionEntityRepository {
	return &collectionEntityRepository{db: db}
}

// Create 创建收款主体
func (r *collectionEntityRepository) Create(ctx context.Context, entity *model.CollectionEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Get 获取收款主体
func (r *collectionEntityRepository) Get(ctx context.Context, id uint64) (*model.CollectionEntity, error) {
	var entity model.CollectionEntity
	err := r.db.WithContext(ctx).
		Preload("PaymentChannels", "enabled = ?", true).
		First(&entity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &entity, nil
}

// GetByCreditCode 根据统一社会信用代码获取收款主体
func (r *collectionEntityRepository) GetByCreditCode(ctx context.Context, creditCode string) (*model.CollectionEntity, error) {
	var entity model.CollectionEntity
	err := r.db.WithContext(ctx).Where("credit_code = ?", creditCode).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &entity, nil
}

// GetDefault 获取默认收款主体
// Requirements: 16.3, 17.4
func (r *collectionEntityRepository) GetDefault(ctx context.Context) (*model.CollectionEntity, error) {
	var entity model.CollectionEntity
	err := r.db.WithContext(ctx).
		Where("is_default = ? AND status = ?", true, model.EntityStatusActive).
		Preload("PaymentChannels", "enabled = ?", true).
		First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &entity, nil
}

// Update 更新收款主体
func (r *collectionEntityRepository) Update(ctx context.Context, entity *model.CollectionEntity) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// List 查询收款主体列表
// Requirements: 15.3
func (r *collectionEntityRepository) List(ctx context.Context, opts ListOptions) ([]model.CollectionEntity, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.CollectionEntity{})

	// 过滤条件
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.Keyword != "" {
		keyword := "%" + opts.Keyword + "%"
		query = query.Where("name LIKE ? OR credit_code LIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderClause := "created_at DESC"
	if opts.SortBy != "" {
		order := "ASC"
		if opts.SortOrder == "desc" {
			order = "DESC"
		}
		switch opts.SortBy {
		case "name":
			orderClause = "name " + order
		case "created_at":
			orderClause = "created_at " + order
		case "total_collection_cents":
			orderClause = "total_collection_cents " + order
		}
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
	var entities []model.CollectionEntity
	err := query.
		Preload("PaymentChannels").
		Order(orderClause).
		Offset(offset).
		Limit(opts.PageSize).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// ToggleStatus 切换收款主体状态
// Requirements: 15.5
func (r *collectionEntityRepository) ToggleStatus(ctx context.Context, id uint64, status model.EntityStatus) error {
	return r.db.WithContext(ctx).Model(&model.CollectionEntity{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// ListActive 获取所有活跃的收款主体
// Requirements: 17.5 - 用于容错处理，当默认主体不可用时查找备用主体
func (r *collectionEntityRepository) ListActive(ctx context.Context) ([]model.CollectionEntity, error) {
	var entities []model.CollectionEntity
	err := r.db.WithContext(ctx).
		Where("status = ?", model.EntityStatusActive).
		Preload("PaymentChannels", "enabled = ?", true).
		Order("is_default DESC, created_at ASC"). // 默认主体优先
		Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// SetDefault 设置默认收款主体
// Requirements: 16.3
func (r *collectionEntityRepository) SetDefault(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先取消所有默认标记
		if err := tx.Model(&model.CollectionEntity{}).
			Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// 设置新的默认主体
		return tx.Model(&model.CollectionEntity{}).
			Where("id = ?", id).
			Update("is_default", true).Error
	})
}

// CreateChannel 创建支付渠道配置
// Requirements: 15.4
func (r *collectionEntityRepository) CreateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error {
	return r.db.WithContext(ctx).Create(channel).Error
}

// GetChannel 获取支付渠道配置
func (r *collectionEntityRepository) GetChannel(ctx context.Context, id uint64) (*model.PaymentChannelConfig, error) {
	var channel model.PaymentChannelConfig
	err := r.db.WithContext(ctx).First(&channel, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &channel, nil
}

// GetChannelByEntityAndMethod 根据收款主体和支付方式获取渠道配置
func (r *collectionEntityRepository) GetChannelByEntityAndMethod(ctx context.Context, entityID uint64, method model.PaymentMethod) (*model.PaymentChannelConfig, error) {
	var channel model.PaymentChannelConfig
	err := r.db.WithContext(ctx).
		Where("collection_entity_id = ? AND channel = ? AND enabled = ?", entityID, method, true).
		Order("priority ASC").
		First(&channel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &channel, nil
}

// UpdateChannel 更新支付渠道配置
func (r *collectionEntityRepository) UpdateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error {
	return r.db.WithContext(ctx).Save(channel).Error
}

// DeleteChannel 删除支付渠道配置
func (r *collectionEntityRepository) DeleteChannel(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.PaymentChannelConfig{}, id).Error
}

// ListChannelsByEntity 获取收款主体的所有支付渠道配置
func (r *collectionEntityRepository) ListChannelsByEntity(ctx context.Context, entityID uint64) ([]model.PaymentChannelConfig, error) {
	var channels []model.PaymentChannelConfig
	err := r.db.WithContext(ctx).
		Where("collection_entity_id = ?", entityID).
		Order("channel ASC, priority ASC").
		Find(&channels).Error
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// CreateHistory 创建收款主体修改历史
func (r *collectionEntityRepository) CreateHistory(ctx context.Context, history *model.CollectionEntityHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// GetHistory 获取收款主体修改历史
func (r *collectionEntityRepository) GetHistory(ctx context.Context, entityID uint64) ([]model.CollectionEntityHistory, error) {
	var histories []model.CollectionEntityHistory
	err := r.db.WithContext(ctx).
		Where("collection_entity_id = ?", entityID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// UpdateCollectionStats 更新收款统计
// Requirements: 18.1
func (r *collectionEntityRepository) UpdateCollectionStats(ctx context.Context, entityID uint64, amountCents int64) error {
	return r.db.WithContext(ctx).Model(&model.CollectionEntity{}).
		Where("id = ?", entityID).
		Updates(map[string]interface{}{
			"total_collection_cents": gorm.Expr("total_collection_cents + ?", amountCents),
			"transaction_count":      gorm.Expr("transaction_count + 1"),
		}).Error
}

// GetCollectionStats 获取收款统计
func (r *collectionEntityRepository) GetCollectionStats(ctx context.Context, entityID uint64) (totalCents int64, count int64, err error) {
	var entity model.CollectionEntity
	err = r.db.WithContext(ctx).
		Select("total_collection_cents, transaction_count").
		First(&entity, entityID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, repository.ErrNotFound
		}
		return 0, 0, err
	}
	return entity.TotalCollectionCents, entity.TransactionCount, nil
}

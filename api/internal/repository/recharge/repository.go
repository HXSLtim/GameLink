package recharge

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 充值仓库
type Repository struct {
	db *gorm.DB
}

// NewRechargeRepository 创建充值仓库
func NewRechargeRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 充值档位 CRUD
// ============================================================================

// OptionListOptions 档位列表查询选项
type OptionListOptions struct {
	Page          int
	PageSize      int
	Keyword       string
	IsActive      *bool
	IsRecommended *bool
	MinVipLevel   *uint64
}

// ListOptions 获取档位列表
func (r *Repository) ListOptions(ctx context.Context, opts OptionListOptions) ([]model.RechargeOption, int64, error) {
	var items []model.RechargeOption
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RechargeOption{})

	if opts.Keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+opts.Keyword+"%", "%"+opts.Keyword+"%")
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}
	if opts.IsRecommended != nil {
		query = query.Where("is_recommended = ?", *opts.IsRecommended)
	}
	if opts.MinVipLevel != nil {
		query = query.Where("min_vip_level IS NULL OR min_vip_level <= ?", *opts.MinVipLevel)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count options: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Order("sort_order ASC, created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list options: %w", err)
	}

	return items, total, nil
}

// GetActiveOptions 获取启用的档位列表（用户端）
func (r *Repository) GetActiveOptions(ctx context.Context, vipLevel *uint64) ([]model.RechargeOption, error) {
	var items []model.RechargeOption
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	if vipLevel != nil {
		query = query.Where("min_vip_level IS NULL OR min_vip_level <= ?", *vipLevel)
	} else {
		query = query.Where("min_vip_level IS NULL")
	}

	if err := query.Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get active options: %w", err)
	}
	return items, nil
}

// GetOptionByID 根据ID获取档位
func (r *Repository) GetOptionByID(ctx context.Context, id uint64) (*model.RechargeOption, error) {
	var item model.RechargeOption
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get option: %w", err)
	}
	return &item, nil
}

// CreateOption 创建档位
func (r *Repository) CreateOption(ctx context.Context, item *model.RechargeOption) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create option: %w", err)
	}
	return nil
}

// UpdateOption 更新档位
func (r *Repository) UpdateOption(ctx context.Context, item *model.RechargeOption) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update option: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteOption 删除档位
func (r *Repository) DeleteOption(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.RechargeOption{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete option: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// IncrementPurchaseCount 增加购买次数
func (r *Repository) IncrementPurchaseCount(ctx context.Context, optionID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.RechargeOption{}).
		Where("id = ?", optionID).
		UpdateColumn("purchase_count", gorm.Expr("purchase_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("increment purchase count: %w", result.Error)
	}
	return nil
}

// BatchUpdateOptionStatus 批量更新档位状态
func (r *Repository) BatchUpdateOptionStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.RechargeOption{}).
		Where("id IN ?", ids).
		Update("is_active", isActive)
	if result.Error != nil {
		return 0, fmt.Errorf("batch update status: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// BatchDeleteOptions 批量删除档位
func (r *Repository) BatchDeleteOptions(ctx context.Context, ids []uint64) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&model.RechargeOption{}, ids)
	if result.Error != nil {
		return 0, fmt.Errorf("batch delete: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ============================================================================
// 充值记录 CRUD
// ============================================================================

// RecordListOptions 记录列表查询选项
type RecordListOptions struct {
	Page           int
	PageSize       int
	UserID         *uint64
	OptionID       *uint64
	Status         *model.RechargeStatus
	PaymentChannel *string
	OrderNo        string
	StartTime      *time.Time
	EndTime        *time.Time
}

// ListRecords 获取充值记录列表
func (r *Repository) ListRecords(ctx context.Context, opts RecordListOptions) ([]model.RechargeRecord, int64, error) {
	var items []model.RechargeRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&model.RechargeRecord{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.OptionID != nil {
		query = query.Where("option_id = ?", *opts.OptionID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.PaymentChannel != nil {
		query = query.Where("payment_channel = ?", *opts.PaymentChannel)
	}
	if opts.OrderNo != "" {
		query = query.Where("order_no ILIKE ? OR merchant_order_no ILIKE ?", "%"+opts.OrderNo+"%", "%"+opts.OrderNo+"%")
	}
	if opts.StartTime != nil {
		query = query.Where("created_at >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		query = query.Where("created_at <= ?", *opts.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count records: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Preload("User").Preload("Option").
		Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}

	return items, total, nil
}

// GetRecordByID 根据ID获取记录
func (r *Repository) GetRecordByID(ctx context.Context, id uint64) (*model.RechargeRecord, error) {
	var item model.RechargeRecord
	if err := r.db.WithContext(ctx).Preload("User").Preload("Option").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get record: %w", err)
	}
	return &item, nil
}

// GetRecordByOrderNo 根据订单号获取记录
func (r *Repository) GetRecordByOrderNo(ctx context.Context, orderNo string) (*model.RechargeRecord, error) {
	var item model.RechargeRecord
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get record by order no: %w", err)
	}
	return &item, nil
}

// GetRecordByMerchantOrderNo 根据商户订单号获取记录
func (r *Repository) GetRecordByMerchantOrderNo(ctx context.Context, merchantOrderNo string) (*model.RechargeRecord, error) {
	var item model.RechargeRecord
	if err := r.db.WithContext(ctx).Where("merchant_order_no = ?", merchantOrderNo).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get record by merchant order no: %w", err)
	}
	return &item, nil
}

// CreateRecord 创建充值记录
func (r *Repository) CreateRecord(ctx context.Context, item *model.RechargeRecord) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create record: %w", err)
	}
	return nil
}

// UpdateRecord 更新充值记录
func (r *Repository) UpdateRecord(ctx context.Context, item *model.RechargeRecord) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update record: %w", result.Error)
	}
	return nil
}

// UpdateRecordStatus 更新记录状态
func (r *Repository) UpdateRecordStatus(ctx context.Context, id uint64, status model.RechargeStatus) error {
	result := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("update status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// MarkAsPaid 标记为已支付
func (r *Repository) MarkAsPaid(ctx context.Context, id uint64, providerTradeNo string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("id = ? AND status = ?", id, model.RechargeStatusPending).
		Updates(map[string]any{
			"status":            model.RechargeStatusPaid,
			"provider_trade_no": providerTradeNo,
			"paid_at":           now,
		})
	if result.Error != nil {
		return fmt.Errorf("mark as paid: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("record not pending or not found")
	}
	return nil
}

// MarkAsRefunded 标记为已退款
func (r *Repository) MarkAsRefunded(ctx context.Context, id uint64, refundAmount int64, reason, providerNo string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("id = ? AND status = ?", id, model.RechargeStatusPaid).
		Updates(map[string]any{
			"status":              model.RechargeStatusRefunded,
			"refunded_at":         now,
			"refund_amount_cents": refundAmount,
			"refund_reason":       reason,
			"refund_provider_no":  providerNo,
		})
	if result.Error != nil {
		return fmt.Errorf("mark as refunded: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("record not paid or not found")
	}
	return nil
}

// MarkCouponIssued 标记优惠券已发放
func (r *Repository) MarkCouponIssued(ctx context.Context, id uint64, couponIDs string) error {
	result := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"coupon_issued": true,
			"coupon_ids":    couponIDs,
		})
	if result.Error != nil {
		return fmt.Errorf("mark coupon issued: %w", result.Error)
	}
	return nil
}

// CountUserPurchases 统计用户购买某档位次数
func (r *Repository) CountUserPurchases(ctx context.Context, userID, optionID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("user_id = ? AND option_id = ? AND status = ?", userID, optionID, model.RechargeStatusPaid).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user purchases: %w", err)
	}
	return count, nil
}

// GetUserRecords 获取用户充值记录
func (r *Repository) GetUserRecords(ctx context.Context, userID uint64, limit int) ([]model.RechargeRecord, error) {
	var items []model.RechargeRecord
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get user records: %w", err)
	}
	return items, nil
}

// GetRechargeStats 获取充值统计
func (r *Repository) GetRechargeStats(ctx context.Context) (map[string]any, error) {
	stats := make(map[string]any)

	// 总充值金额
	var totalAmount int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ?", model.RechargeStatusPaid).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&totalAmount).Error; err != nil {
		return nil, fmt.Errorf("sum total amount: %w", err)
	}
	stats["totalAmountCents"] = totalAmount

	// 总充值笔数
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ?", model.RechargeStatusPaid).
		Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("count total: %w", err)
	}
	stats["totalCount"] = totalCount

	// 今日充值金额
	today := time.Now().Truncate(24 * time.Hour)
	var todayAmount int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ? AND paid_at >= ?", model.RechargeStatusPaid, today).
		Select("COALESCE(SUM(amount_cents), 0)").
		Scan(&todayAmount).Error; err != nil {
		return nil, fmt.Errorf("sum today amount: %w", err)
	}
	stats["todayAmountCents"] = todayAmount

	// 今日充值笔数
	var todayCount int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ? AND paid_at >= ?", model.RechargeStatusPaid, today).
		Count(&todayCount).Error; err != nil {
		return nil, fmt.Errorf("count today: %w", err)
	}
	stats["todayCount"] = todayCount

	// 总退款金额
	var refundAmount int64
	if err := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ?", model.RechargeStatusRefunded).
		Select("COALESCE(SUM(refund_amount_cents), 0)").
		Scan(&refundAmount).Error; err != nil {
		return nil, fmt.Errorf("sum refund amount: %w", err)
	}
	stats["refundAmountCents"] = refundAmount

	return stats, nil
}

// CancelExpiredRecords 取消过期未支付的记录
func (r *Repository) CancelExpiredRecords(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.RechargeRecord{}).
		Where("status = ? AND expire_at IS NOT NULL AND expire_at <= ?", model.RechargeStatusPending, time.Now()).
		Update("status", model.RechargeStatusCanceled)
	if result.Error != nil {
		return 0, fmt.Errorf("cancel expired: %w", result.Error)
	}
	return result.RowsAffected, nil
}

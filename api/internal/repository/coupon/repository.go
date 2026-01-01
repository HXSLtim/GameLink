package coupon

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Repository 优惠券仓库
type Repository struct {
	db *gorm.DB
}

// NewCouponRepository 创建优惠券仓库
func NewCouponRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ============================================================================
// 优惠券模板 CRUD
// ============================================================================

// TemplateListOptions 模板列表查询选项
type TemplateListOptions struct {
	Page     int
	PageSize int
	Keyword  string
	Type     *model.CouponType
	Source   *model.CouponSource
	Scope    *model.CouponScope
	IsActive *bool
}

// ListTemplates 获取模板列表
func (r *Repository) ListTemplates(ctx context.Context, opts TemplateListOptions) ([]model.CouponTemplate, int64, error) {
	var items []model.CouponTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&model.CouponTemplate{})

	if opts.Keyword != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+opts.Keyword+"%", "%"+opts.Keyword+"%")
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.Source != nil {
		query = query.Where("source = ?", *opts.Source)
	}
	if opts.Scope != nil {
		query = query.Where("scope = ?", *opts.Scope)
	}
	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count templates: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list templates: %w", err)
	}

	return items, total, nil
}

// GetTemplateByID 根据ID获取模板
func (r *Repository) GetTemplateByID(ctx context.Context, id uint64) (*model.CouponTemplate, error) {
	var item model.CouponTemplate
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get template: %w", err)
	}
	return &item, nil
}

// GetTemplateByClaimLink 根据领取链接获取模板
func (r *Repository) GetTemplateByClaimLink(ctx context.Context, link string) (*model.CouponTemplate, error) {
	var item model.CouponTemplate
	if err := r.db.WithContext(ctx).Where("claim_link = ?", link).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get template by link: %w", err)
	}
	return &item, nil
}

// CreateTemplate 创建模板
func (r *Repository) CreateTemplate(ctx context.Context, item *model.CouponTemplate) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

// UpdateTemplate 更新模板
func (r *Repository) UpdateTemplate(ctx context.Context, item *model.CouponTemplate) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update template: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteTemplate 删除模板
func (r *Repository) DeleteTemplate(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.CouponTemplate{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete template: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// IncrementClaimedCount 增加已领取数量
func (r *Repository) IncrementClaimedCount(ctx context.Context, templateID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.CouponTemplate{}).
		Where("id = ?", templateID).
		UpdateColumn("claimed_count", gorm.Expr("claimed_count + 1"))
	if result.Error != nil {
		return fmt.Errorf("increment claimed count: %w", result.Error)
	}
	return nil
}

// BatchUpdateTemplateStatus 批量更新模板状态
func (r *Repository) BatchUpdateTemplateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.CouponTemplate{}).
		Where("id IN ?", ids).
		Update("is_active", isActive)
	if result.Error != nil {
		return 0, fmt.Errorf("batch update status: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// BatchDeleteTemplates 批量删除模板
func (r *Repository) BatchDeleteTemplates(ctx context.Context, ids []uint64) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&model.CouponTemplate{}, ids)
	if result.Error != nil {
		return 0, fmt.Errorf("batch delete: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ============================================================================
// 用户优惠券 CRUD
// ============================================================================

// CouponListOptions 优惠券列表查询选项
type CouponListOptions struct {
	Page       int
	PageSize   int
	UserID     *uint64
	TemplateID *uint64
	State      *model.CouponState
	Type       *model.CouponType
	Source     *model.CouponSource
	ExpireSoon bool // 即将过期（7天内）
}

// ListCoupons 获取优惠券列表
func (r *Repository) ListCoupons(ctx context.Context, opts CouponListOptions) ([]model.Coupon, int64, error) {
	var items []model.Coupon
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Coupon{})

	if opts.UserID != nil {
		query = query.Where("user_id = ?", *opts.UserID)
	}
	if opts.TemplateID != nil {
		query = query.Where("template_id = ?", *opts.TemplateID)
	}
	if opts.State != nil {
		query = query.Where("state = ?", *opts.State)
	}
	if opts.Type != nil {
		query = query.Where("type = ?", *opts.Type)
	}
	if opts.Source != nil {
		query = query.Where("source = ?", *opts.Source)
	}
	if opts.ExpireSoon {
		query = query.Where("expire_at <= ? AND expire_at > ?", time.Now().AddDate(0, 0, 7), time.Now())
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count coupons: %w", err)
	}

	offset := (opts.Page - 1) * opts.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(opts.PageSize).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("list coupons: %w", err)
	}

	return items, total, nil
}

// GetCouponByID 根据ID获取优惠券
func (r *Repository) GetCouponByID(ctx context.Context, id uint64) (*model.Coupon, error) {
	var item model.Coupon
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get coupon: %w", err)
	}
	return &item, nil
}

// GetCouponWithTemplate 获取优惠券（含模板）
func (r *Repository) GetCouponWithTemplate(ctx context.Context, id uint64) (*model.Coupon, error) {
	var item model.Coupon
	if err := r.db.WithContext(ctx).Preload("Template").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get coupon with template: %w", err)
	}
	return &item, nil
}

// CreateCoupon 创建优惠券
func (r *Repository) CreateCoupon(ctx context.Context, item *model.Coupon) error {
	if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
		return fmt.Errorf("create coupon: %w", err)
	}
	return nil
}

// UpdateCoupon 更新优惠券
func (r *Repository) UpdateCoupon(ctx context.Context, item *model.Coupon) error {
	result := r.db.WithContext(ctx).Model(item).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update coupon: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// DeleteCoupon 删除优惠券
func (r *Repository) DeleteCoupon(ctx context.Context, id uint64) error {
	result := r.db.WithContext(ctx).Delete(&model.Coupon{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete coupon: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// CountUserCouponsFromTemplate 统计用户从某模板领取的优惠券数量
func (r *Repository) CountUserCouponsFromTemplate(ctx context.Context, userID, templateID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("user_id = ? AND template_id = ?", userID, templateID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user coupons: %w", err)
	}
	return count, nil
}

// GetUserAvailableCoupons 获取用户可用优惠券
func (r *Repository) GetUserAvailableCoupons(ctx context.Context, userID uint64) ([]model.Coupon, error) {
	var items []model.Coupon
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND state = ? AND expire_at > ?", userID, model.CouponStateAvailable, time.Now()).
		Order("expire_at ASC").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("get available coupons: %w", err)
	}
	return items, nil
}

// LockCoupon 锁定优惠券（下单时）
func (r *Repository) LockCoupon(ctx context.Context, couponID, orderID uint64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("id = ? AND state = ? AND expire_at > ?", couponID, model.CouponStateAvailable, now).
		Updates(map[string]interface{}{
			"state":              model.CouponStateLocked,
			"locked_by_order_id": orderID,
			"locked_at":          now,
		})
	if result.Error != nil {
		return fmt.Errorf("lock coupon: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("coupon not available or expired")
	}
	return nil
}

// UnlockCoupon 解锁优惠券（取消订单时）
func (r *Repository) UnlockCoupon(ctx context.Context, couponID uint64) error {
	result := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("id = ? AND state = ?", couponID, model.CouponStateLocked).
		Updates(map[string]interface{}{
			"state":              model.CouponStateAvailable,
			"locked_by_order_id": nil,
			"locked_at":          nil,
		})
	if result.Error != nil {
		return fmt.Errorf("unlock coupon: %w", result.Error)
	}
	return nil
}

// UseCoupon 使用优惠券（支付成功时）
func (r *Repository) UseCoupon(ctx context.Context, couponID, orderID uint64, discountCents int64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("id = ? AND state = ?", couponID, model.CouponStateLocked).
		Updates(map[string]interface{}{
			"state":          model.CouponStateUsed,
			"used_order_id":  orderID,
			"used_at":        now,
			"discount_cents": discountCents,
		})
	if result.Error != nil {
		return fmt.Errorf("use coupon: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("coupon not locked")
	}
	return nil
}

// ExpireOldCoupons 过期旧优惠券（定时任务）
func (r *Repository) ExpireOldCoupons(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("state = ? AND expire_at <= ?", model.CouponStateAvailable, time.Now()).
		Update("state", model.CouponStateExpired)
	if result.Error != nil {
		return 0, fmt.Errorf("expire coupons: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// GetCouponStats 获取优惠券统计
func (r *Repository) GetCouponStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 模板总数
	var totalTemplates int64
	if err := r.db.WithContext(ctx).Model(&model.CouponTemplate{}).Count(&totalTemplates).Error; err != nil {
		return nil, fmt.Errorf("count templates: %w", err)
	}
	stats["totalTemplates"] = totalTemplates

	// 启用的模板数
	var activeTemplates int64
	if err := r.db.WithContext(ctx).Model(&model.CouponTemplate{}).Where("is_active = ?", true).Count(&activeTemplates).Error; err != nil {
		return nil, fmt.Errorf("count active templates: %w", err)
	}
	stats["activeTemplates"] = activeTemplates

	// 优惠券总数
	var totalCoupons int64
	if err := r.db.WithContext(ctx).Model(&model.Coupon{}).Count(&totalCoupons).Error; err != nil {
		return nil, fmt.Errorf("count total coupons: %w", err)
	}
	stats["totalCoupons"] = totalCoupons

	// 按状态统计
	type stateCount struct {
		State model.CouponState
		Count int64
	}
	var stateCounts []stateCount
	if err := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Select("state, count(*) as count").
		Group("state").
		Scan(&stateCounts).Error; err != nil {
		return nil, fmt.Errorf("count by state: %w", err)
	}
	for _, sc := range stateCounts {
		switch sc.State {
		case model.CouponStateAvailable:
			stats["availableCoupons"] = sc.Count
		case model.CouponStateUsed:
			stats["usedCoupons"] = sc.Count
		case model.CouponStateExpired:
			stats["expiredCoupons"] = sc.Count
		}
	}

	// 确保所有字段都有值（即使为0）
	if _, ok := stats["availableCoupons"]; !ok {
		stats["availableCoupons"] = 0
	}
	if _, ok := stats["usedCoupons"]; !ok {
		stats["usedCoupons"] = 0
	}
	if _, ok := stats["expiredCoupons"]; !ok {
		stats["expiredCoupons"] = 0
	}

	// 总折扣金额
	var totalDiscount int64
	if err := r.db.WithContext(ctx).Model(&model.Coupon{}).
		Where("state = ?", model.CouponStateUsed).
		Select("COALESCE(SUM(discount_cents), 0)").
		Scan(&totalDiscount).Error; err != nil {
		return nil, fmt.Errorf("sum discount: %w", err)
	}
	stats["totalDiscountCents"] = totalDiscount

	return stats, nil
}

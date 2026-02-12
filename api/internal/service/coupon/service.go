package coupon

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	couponrepo "gamelink/internal/repository/coupon"
)

// CouponRepository defines the interface for coupon repository operations
type CouponRepository interface {
	ListTemplates(ctx context.Context, opts couponrepo.TemplateListOptions) ([]model.CouponTemplate, int64, error)
	GetTemplateByID(ctx context.Context, id uint64) (*model.CouponTemplate, error)
	GetTemplateByClaimLink(ctx context.Context, link string) (*model.CouponTemplate, error)
	CreateTemplate(ctx context.Context, template *model.CouponTemplate) error
	UpdateTemplate(ctx context.Context, template *model.CouponTemplate) error
	DeleteTemplate(ctx context.Context, id uint64) error
	BatchUpdateTemplateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error)
	BatchDeleteTemplates(ctx context.Context, ids []uint64) (int64, error)
	ListCoupons(ctx context.Context, opts couponrepo.CouponListOptions) ([]model.Coupon, int64, error)
	GetCouponByID(ctx context.Context, id uint64) (*model.Coupon, error)
	GetCouponWithTemplate(ctx context.Context, id uint64) (*model.Coupon, error)
	GetUserAvailableCoupons(ctx context.Context, userID uint64) ([]model.Coupon, error)
	CountUserCouponsFromTemplate(ctx context.Context, userID, templateID uint64) (int64, error)
	CreateCoupon(ctx context.Context, coupon *model.Coupon) error
	IncrementClaimedCount(ctx context.Context, templateID uint64) error
	LockCoupon(ctx context.Context, couponID, orderID uint64) error
	UnlockCoupon(ctx context.Context, couponID uint64) error
	UseCoupon(ctx context.Context, couponID, orderID uint64, discountCents int64) error
	ExpireOldCoupons(ctx context.Context) (int64, error)
	GetCouponStats(ctx context.Context) (map[string]int64, error)
	DeleteCoupon(ctx context.Context, id uint64) error
}

// Service 优惠券业务逻辑层
type Service struct {
	repo CouponRepository
}

// NewCouponService 创建优惠券服务
func NewCouponService(repo CouponRepository) *Service {
	return &Service{repo: repo}
}

// ============================================================================
// 优惠券模板管理
// ============================================================================

// ListTemplates 获取模板列表
func (s *Service) ListTemplates(ctx context.Context, opts couponrepo.TemplateListOptions) ([]model.CouponTemplate, int64, error) {
	templates, total, err := s.repo.ListTemplates(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list templates: %w", err)
	}
	return templates, total, nil
}

// GetTemplate 获取模板详情
func (s *Service) GetTemplate(ctx context.Context, id uint64) (*model.CouponTemplate, error) {
	template, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	return template, nil
}

// GetTemplateByClaimLink 根据领取链接获取模板
func (s *Service) GetTemplateByClaimLink(ctx context.Context, link string) (*model.CouponTemplate, error) {
	template, err := s.repo.GetTemplateByClaimLink(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("get template by link: %w", err)
	}
	return template, nil
}

// CreateTemplate 创建模板
func (s *Service) CreateTemplate(ctx context.Context, template *model.CouponTemplate) error {
	// 验证模板配置
	if err := s.validateTemplate(template); err != nil {
		return err
	}

	if err := s.repo.CreateTemplate(ctx, template); err != nil {
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

// UpdateTemplate 更新模板
func (s *Service) UpdateTemplate(ctx context.Context, template *model.CouponTemplate) error {
	// 检查是否存在
	if _, err := s.repo.GetTemplateByID(ctx, template.ID); err != nil {
		return fmt.Errorf("get template: %w", err)
	}

	// 验证模板配置
	if err := s.validateTemplate(template); err != nil {
		return err
	}

	if err := s.repo.UpdateTemplate(ctx, template); err != nil {
		return fmt.Errorf("update template: %w", err)
	}
	return nil
}

// DeleteTemplate 删除模板
func (s *Service) DeleteTemplate(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteTemplate(ctx, id); err != nil {
		return fmt.Errorf("delete template: %w", err)
	}
	return nil
}

// BatchUpdateTemplateStatus 批量更新模板状态
func (s *Service) BatchUpdateTemplateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	affected, err := s.repo.BatchUpdateTemplateStatus(ctx, ids, isActive)
	if err != nil {
		return 0, fmt.Errorf("batch update status: %w", err)
	}
	return affected, nil
}

// BatchDeleteTemplates 批量删除模板
func (s *Service) BatchDeleteTemplates(ctx context.Context, ids []uint64) (int64, error) {
	affected, err := s.repo.BatchDeleteTemplates(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete: %w", err)
	}
	return affected, nil
}

// validateTemplate 验证模板配置
func (s *Service) validateTemplate(template *model.CouponTemplate) error {
	if template.Name == "" {
		return errors.New("模板名称不能为空")
	}

	switch template.Type {
	case model.CouponTypeDeduct:
		if template.DeductAmountCents <= 0 {
			return errors.New("满减金额必须大于0")
		}
	case model.CouponTypeDiscount:
		if template.DiscountRate <= 0 || template.DiscountRate >= 1 {
			return errors.New("折扣率必须在0-1之间")
		}
	default:
		return errors.New("无效的优惠券类型")
	}

	return nil
}

// ============================================================================
// 用户优惠券管理
// ============================================================================

// ListCoupons 获取优惠券列表
func (s *Service) ListCoupons(ctx context.Context, opts couponrepo.CouponListOptions) ([]model.Coupon, int64, error) {
	coupons, total, err := s.repo.ListCoupons(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list coupons: %w", err)
	}
	return coupons, total, nil
}

// GetCoupon 获取优惠券详情
func (s *Service) GetCoupon(ctx context.Context, id uint64) (*model.Coupon, error) {
	coupon, err := s.repo.GetCouponByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get coupon: %w", err)
	}
	return coupon, nil
}

// GetCouponWithTemplate 获取优惠券详情（含模板）
func (s *Service) GetCouponWithTemplate(ctx context.Context, id uint64) (*model.Coupon, error) {
	coupon, err := s.repo.GetCouponWithTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get coupon with template: %w", err)
	}
	return coupon, nil
}

// GetUserAvailableCoupons 获取用户可用优惠券
func (s *Service) GetUserAvailableCoupons(ctx context.Context, userID uint64) ([]model.Coupon, error) {
	coupons, err := s.repo.GetUserAvailableCoupons(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get available coupons: %w", err)
	}
	return coupons, nil
}

// UserCouponStats 用户优惠券统计
type UserCouponStats struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
	Used      int64 `json:"used"`
	Expired   int64 `json:"expired"`
	Locked    int64 `json:"locked"`
	Deleted   int64 `json:"deleted"`
}

// GetUserCouponStats 获取用户优惠券统计
func (s *Service) GetUserCouponStats(ctx context.Context, userID uint64) (*UserCouponStats, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	countByState := func(state *model.CouponState) (int64, error) {
		opts := couponrepo.CouponListOptions{
			Page:     1,
			PageSize: 1,
			UserID:   &userID,
			State:    state,
		}
		_, total, err := s.repo.ListCoupons(ctx, opts)
		if err != nil {
			return 0, err
		}
		return total, nil
	}

	var (
		availableState = model.CouponStateAvailable
		usedState      = model.CouponStateUsed
		expiredState   = model.CouponStateExpired
		lockedState    = model.CouponStateLocked
		deletedState   = model.CouponStateDeleted
	)

	total, err := countByState(nil)
	if err != nil {
		return nil, fmt.Errorf("count total user coupons: %w", err)
	}
	available, err := countByState(&availableState)
	if err != nil {
		return nil, fmt.Errorf("count available user coupons: %w", err)
	}
	used, err := countByState(&usedState)
	if err != nil {
		return nil, fmt.Errorf("count used user coupons: %w", err)
	}
	expired, err := countByState(&expiredState)
	if err != nil {
		return nil, fmt.Errorf("count expired user coupons: %w", err)
	}
	locked, err := countByState(&lockedState)
	if err != nil {
		return nil, fmt.Errorf("count locked user coupons: %w", err)
	}
	deleted, err := countByState(&deletedState)
	if err != nil {
		return nil, fmt.Errorf("count deleted user coupons: %w", err)
	}

	return &UserCouponStats{
		Total:     total,
		Available: available,
		Used:      used,
		Expired:   expired,
		Locked:    locked,
		Deleted:   deleted,
	}, nil
}

// ClaimCoupon 领取优惠券
func (s *Service) ClaimCoupon(ctx context.Context, userID, templateID uint64) (*model.Coupon, error) {
	// 获取模板
	template, err := s.repo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}

	// 检查模板是否启用
	if !template.IsActive {
		return nil, errors.New("该优惠券已下架")
	}

	// 检查总量限制
	if template.TotalCount > 0 && template.ClaimedCount >= template.TotalCount {
		return nil, errors.New("优惠券已领完")
	}

	// 检查每人限领数量
	if template.PerUserLimit > 0 {
		count, err := s.repo.CountUserCouponsFromTemplate(ctx, userID, templateID)
		if err != nil {
			return nil, fmt.Errorf("count user coupons: %w", err)
		}
		if int(count) >= template.PerUserLimit {
			return nil, errors.New("已达到领取上限")
		}
	}

	// 计算过期时间
	var expireAt time.Time
	if template.ValidityType == "fixed" && template.FixedExpireAt != nil {
		expireAt = *template.FixedExpireAt
	} else {
		expireAt = time.Now().AddDate(0, 0, template.ValidityDays)
	}

	// 创建优惠券
	now := time.Now()
	coupon := &model.Coupon{
		TemplateID:        templateID,
		UserID:            userID,
		State:             model.CouponStateAvailable,
		Name:              template.Name,
		Type:              template.Type,
		Source:            template.Source,
		MinAmountCents:    template.MinAmountCents,
		DeductAmountCents: template.DeductAmountCents,
		DiscountRate:      template.DiscountRate,
		MaxDiscountCents:  template.MaxDiscountCents,
		Scope:             template.Scope,
		GameIDs:           template.GameIDs,
		ItemIDs:           template.ItemIDs,
		ClaimedAt:         &now,
		ExpireAt:          expireAt,
	}

	if err := s.repo.CreateCoupon(ctx, coupon); err != nil {
		return nil, fmt.Errorf("create coupon: %w", err)
	}

	// 增加已领取数量
	if err := s.repo.IncrementClaimedCount(ctx, templateID); err != nil {
		// 不影响主流程，记录日志即可
	}

	return coupon, nil
}

// ClaimCouponByLink 通过链接领取优惠券
func (s *Service) ClaimCouponByLink(ctx context.Context, userID uint64, link string) (*model.Coupon, error) {
	template, err := s.repo.GetTemplateByClaimLink(ctx, link)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("无效的领取链接")
		}
		return nil, fmt.Errorf("get template: %w", err)
	}

	return s.ClaimCoupon(ctx, userID, template.ID)
}

// IssueCoupon 发放优惠券（系统发放，不检查限制）
func (s *Service) IssueCoupon(ctx context.Context, userID, templateID uint64, source model.CouponSource) (*model.Coupon, error) {
	template, err := s.repo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}

	// 计算过期时间
	var expireAt time.Time
	if template.ValidityType == "fixed" && template.FixedExpireAt != nil {
		expireAt = *template.FixedExpireAt
	} else {
		expireAt = time.Now().AddDate(0, 0, template.ValidityDays)
	}

	now := time.Now()
	coupon := &model.Coupon{
		TemplateID:        templateID,
		UserID:            userID,
		State:             model.CouponStateAvailable,
		Name:              template.Name,
		Type:              template.Type,
		Source:            source,
		MinAmountCents:    template.MinAmountCents,
		DeductAmountCents: template.DeductAmountCents,
		DiscountRate:      template.DiscountRate,
		MaxDiscountCents:  template.MaxDiscountCents,
		Scope:             template.Scope,
		GameIDs:           template.GameIDs,
		ItemIDs:           template.ItemIDs,
		ClaimedAt:         &now,
		ExpireAt:          expireAt,
	}

	if err := s.repo.CreateCoupon(ctx, coupon); err != nil {
		return nil, fmt.Errorf("create coupon: %w", err)
	}

	return coupon, nil
}

// LockCoupon 锁定优惠券（下单时）
func (s *Service) LockCoupon(ctx context.Context, couponID, orderID uint64) error {
	if err := s.repo.LockCoupon(ctx, couponID, orderID); err != nil {
		return fmt.Errorf("lock coupon: %w", err)
	}
	return nil
}

// UnlockCoupon 解锁优惠券（取消订单时）
func (s *Service) UnlockCoupon(ctx context.Context, couponID uint64) error {
	if err := s.repo.UnlockCoupon(ctx, couponID); err != nil {
		return fmt.Errorf("unlock coupon: %w", err)
	}
	return nil
}

// UseCoupon 使用优惠券（支付成功时）
func (s *Service) UseCoupon(ctx context.Context, couponID, orderID uint64, discountCents int64) error {
	if err := s.repo.UseCoupon(ctx, couponID, orderID, discountCents); err != nil {
		return fmt.Errorf("use coupon: %w", err)
	}
	return nil
}

// CalculateDiscount 计算优惠券折扣金额
func (s *Service) CalculateDiscount(ctx context.Context, couponID uint64, orderAmountCents int64) (int64, error) {
	coupon, err := s.repo.GetCouponByID(ctx, couponID)
	if err != nil {
		return 0, fmt.Errorf("get coupon: %w", err)
	}

	if !coupon.IsValid() {
		return 0, errors.New("优惠券不可用")
	}

	return coupon.CalculateDiscount(orderAmountCents), nil
}

// ExpireOldCoupons 过期旧优惠券（定时任务）
func (s *Service) ExpireOldCoupons(ctx context.Context) (int64, error) {
	affected, err := s.repo.ExpireOldCoupons(ctx)
	if err != nil {
		return 0, fmt.Errorf("expire coupons: %w", err)
	}
	return affected, nil
}

// GetCouponStats 获取优惠券统计
func (s *Service) GetCouponStats(ctx context.Context) (map[string]int64, error) {
	stats, err := s.repo.GetCouponStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	return stats, nil
}

// DeleteCoupon 删除优惠券（管理员）
func (s *Service) DeleteCoupon(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteCoupon(ctx, id); err != nil {
		return fmt.Errorf("delete coupon: %w", err)
	}
	return nil
}

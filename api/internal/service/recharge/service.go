package recharge

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	rechargerepo "gamelink/internal/repository/recharge"
	walletrepo "gamelink/internal/repository/wallet"
)

// RechargeRepository 定义充值仓库接口
type RechargeRepository interface {
	ListOptions(ctx context.Context, opts rechargerepo.OptionListOptions) ([]model.RechargeOption, int64, error)
	GetActiveOptions(ctx context.Context, vipLevel *uint64) ([]model.RechargeOption, error)
	GetOptionByID(ctx context.Context, id uint64) (*model.RechargeOption, error)
	CreateOption(ctx context.Context, option *model.RechargeOption) error
	UpdateOption(ctx context.Context, option *model.RechargeOption) error
	DeleteOption(ctx context.Context, id uint64) error
	BatchUpdateOptionStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error)
	BatchDeleteOptions(ctx context.Context, ids []uint64) (int64, error)
	IncrementPurchaseCount(ctx context.Context, optionID uint64) error
	ListRecords(ctx context.Context, opts rechargerepo.RecordListOptions) ([]model.RechargeRecord, int64, error)
	GetRecordByID(ctx context.Context, id uint64) (*model.RechargeRecord, error)
	GetRecordByOrderNo(ctx context.Context, orderNo string) (*model.RechargeRecord, error)
	CreateRecord(ctx context.Context, record *model.RechargeRecord) error
	MarkAsPaid(ctx context.Context, id uint64, providerTradeNo string) error
	MarkAsRefunded(ctx context.Context, id uint64, refundAmount int64, reason, providerNo string) error
	MarkCouponIssued(ctx context.Context, id uint64, couponIDs string) error
	CountUserPurchases(ctx context.Context, userID, optionID uint64) (int64, error)
	GetUserRecords(ctx context.Context, userID uint64, limit int) ([]model.RechargeRecord, error)
	GetRechargeStats(ctx context.Context) (map[string]any, error)
	CancelExpiredRecords(ctx context.Context) (int64, error)
}

// CouponService 定义优惠券服务接口
type CouponService interface {
	IssueCoupon(ctx context.Context, userID, templateID uint64, source model.CouponSource) (*model.Coupon, error)
}

// Service 充值业务逻辑层
type Service struct {
	repo       RechargeRepository
	walletRepo walletrepo.Repository
	couponSvc  CouponService
}

// NewRechargeService 创建充值服务
func NewRechargeService(repo RechargeRepository, walletRepo walletrepo.Repository, couponSvc CouponService) *Service {
	return &Service{
		repo:       repo,
		walletRepo: walletRepo,
		couponSvc:  couponSvc,
	}
}

// ============================================================================
// 充值档位管理
// ============================================================================

// ListOptions 获取档位列表
func (s *Service) ListOptions(ctx context.Context, opts rechargerepo.OptionListOptions) ([]model.RechargeOption, int64, error) {
	options, total, err := s.repo.ListOptions(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list options: %w", err)
	}
	return options, total, nil
}

// GetActiveOptions 获取启用的档位列表（用户端）
func (s *Service) GetActiveOptions(ctx context.Context, vipLevel *uint64) ([]model.RechargeOption, error) {
	options, err := s.repo.GetActiveOptions(ctx, vipLevel)
	if err != nil {
		return nil, fmt.Errorf("get active options: %w", err)
	}
	return options, nil
}

// GetOption 获取档位详情
func (s *Service) GetOption(ctx context.Context, id uint64) (*model.RechargeOption, error) {
	option, err := s.repo.GetOptionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get option: %w", err)
	}
	return option, nil
}

// CreateOption 创建档位
func (s *Service) CreateOption(ctx context.Context, option *model.RechargeOption) error {
	if err := s.validateOption(option); err != nil {
		return err
	}

	// 计算实际到账
	option.TotalCents = option.AmountCents + option.BonusCents

	if err := s.repo.CreateOption(ctx, option); err != nil {
		return fmt.Errorf("create option: %w", err)
	}
	return nil
}

// UpdateOption 更新档位
func (s *Service) UpdateOption(ctx context.Context, option *model.RechargeOption) error {
	if _, err := s.repo.GetOptionByID(ctx, option.ID); err != nil {
		return fmt.Errorf("get option: %w", err)
	}

	if err := s.validateOption(option); err != nil {
		return err
	}

	// 计算实际到账
	option.TotalCents = option.AmountCents + option.BonusCents

	if err := s.repo.UpdateOption(ctx, option); err != nil {
		return fmt.Errorf("update option: %w", err)
	}
	return nil
}

// DeleteOption 删除档位
func (s *Service) DeleteOption(ctx context.Context, id uint64) error {
	if err := s.repo.DeleteOption(ctx, id); err != nil {
		return fmt.Errorf("delete option: %w", err)
	}
	return nil
}

// BatchUpdateOptionStatus 批量更新档位状态
func (s *Service) BatchUpdateOptionStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	affected, err := s.repo.BatchUpdateOptionStatus(ctx, ids, isActive)
	if err != nil {
		return 0, fmt.Errorf("batch update status: %w", err)
	}
	return affected, nil
}

// BatchDeleteOptions 批量删除档位
func (s *Service) BatchDeleteOptions(ctx context.Context, ids []uint64) (int64, error) {
	affected, err := s.repo.BatchDeleteOptions(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("batch delete: %w", err)
	}
	return affected, nil
}

// validateOption 验证档位配置
func (s *Service) validateOption(option *model.RechargeOption) error {
	if option.Name == "" {
		return errors.New("档位名称不能为空")
	}
	if option.AmountCents <= 0 {
		return errors.New("充值金额必须大于0")
	}
	if option.BonusCents < 0 {
		return errors.New("赠送金额不能为负数")
	}
	return nil
}

// ============================================================================
// 充值记录管理
// ============================================================================

// ListRecords 获取充值记录列表
func (s *Service) ListRecords(ctx context.Context, opts rechargerepo.RecordListOptions) ([]model.RechargeRecord, int64, error) {
	records, total, err := s.repo.ListRecords(ctx, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("list records: %w", err)
	}
	return records, total, nil
}

// GetRecord 获取充值记录详情
func (s *Service) GetRecord(ctx context.Context, id uint64) (*model.RechargeRecord, error) {
	record, err := s.repo.GetRecordByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get record: %w", err)
	}
	return record, nil
}

// GetRecordByOrderNo 根据订单号获取记录
func (s *Service) GetRecordByOrderNo(ctx context.Context, orderNo string) (*model.RechargeRecord, error) {
	record, err := s.repo.GetRecordByOrderNo(ctx, orderNo)
	if err != nil {
		return nil, fmt.Errorf("get record by order no: %w", err)
	}
	return record, nil
}

// GetUserRecords 获取用户充值记录
func (s *Service) GetUserRecords(ctx context.Context, userID uint64, limit int) ([]model.RechargeRecord, error) {
	if limit <= 0 {
		limit = 20
	}
	records, err := s.repo.GetUserRecords(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get user records: %w", err)
	}
	return records, nil
}

// CreateRechargeOrder 创建充值订单
func (s *Service) CreateRechargeOrder(ctx context.Context, userID, optionID uint64, paymentChannel, paymentMethod, clientIP, userAgent string) (*model.RechargeRecord, error) {
	// 获取档位
	option, err := s.repo.GetOptionByID(ctx, optionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.New("充值档位不存在")
		}
		return nil, fmt.Errorf("get option: %w", err)
	}

	// 检查档位是否启用
	if !option.IsActive {
		return nil, errors.New("该充值档位已下架")
	}

	// 检查总限购
	if option.TotalLimit > 0 && option.PurchaseCount >= option.TotalLimit {
		return nil, errors.New("该档位已售罄")
	}

	// 检查每人限购
	if option.PerUserLimit > 0 {
		count, err := s.repo.CountUserPurchases(ctx, userID, optionID)
		if err != nil {
			return nil, fmt.Errorf("count purchases: %w", err)
		}
		if int(count) >= option.PerUserLimit {
			return nil, errors.New("已达到购买上限")
		}
	}

	// 生成订单号
	orderNo := model.GenerateOrderNo("RC")
	merchantOrderNo := model.GenerateOrderNo("MRC")

	// 计算过期时间（30分钟）
	expireAt := time.Now().Add(30 * time.Minute)

	record := &model.RechargeRecord{
		UserID:          userID,
		OptionID:        &optionID,
		AmountCents:     option.AmountCents,
		BonusCents:      option.BonusCents,
		TotalCents:      option.TotalCents,
		Status:          model.RechargeStatusPending,
		OrderNo:         orderNo,
		MerchantOrderNo: merchantOrderNo,
		PaymentChannel:  paymentChannel,
		PaymentMethod:   paymentMethod,
		ExpireAt:        &expireAt,
		ClientIP:        clientIP,
		UserAgent:       userAgent,
	}

	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return nil, fmt.Errorf("create record: %w", err)
	}

	return record, nil
}

// HandlePaymentCallback 处理支付回调
func (s *Service) HandlePaymentCallback(ctx context.Context, orderNo, providerTradeNo string) error {
	// 获取记录
	record, err := s.repo.GetRecordByOrderNo(ctx, orderNo)
	if err != nil {
		return fmt.Errorf("get record: %w", err)
	}

	// 检查状态
	if record.Status != model.RechargeStatusPending {
		return errors.New("订单状态不正确")
	}

	// 标记为已支付
	if err := s.repo.MarkAsPaid(ctx, record.ID, providerTradeNo); err != nil {
		return fmt.Errorf("mark as paid: %w", err)
	}

	// 增加钱包余额
	if s.walletRepo != nil {
		wallet, err := s.walletRepo.GetByUserID(ctx, record.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				wallet = &model.Wallet{UserID: record.UserID}
			} else {
				fmt.Printf("get wallet error: %v\n", err)
			}
		}
		if wallet != nil {
			wallet.BalanceCents += record.TotalCents
			if err := s.walletRepo.Save(ctx, wallet); err != nil {
				fmt.Printf("save wallet error: %v\n", err)
			}
		}
	}

	// 增加档位购买次数
	if record.OptionID != nil {
		if err := s.repo.IncrementPurchaseCount(ctx, *record.OptionID); err != nil {
			// 记录错误但不影响主流程
			fmt.Printf("increment purchase count error: %v\n", err)
		}

		// 发放优惠券
		if s.couponSvc != nil {
			option, err := s.repo.GetOptionByID(ctx, *record.OptionID)
			if err == nil && option.CouponTemplateID != nil && option.CouponCount > 0 {
				s.issueCoupons(ctx, record.ID, record.UserID, *option.CouponTemplateID, option.CouponCount)
			}
		}
	}

	return nil
}

// issueCoupons 发放优惠券
func (s *Service) issueCoupons(ctx context.Context, recordID, userID, templateID uint64, count int) {
	var couponIDs []uint64
	for i := 0; i < count; i++ {
		coupon, err := s.couponSvc.IssueCoupon(ctx, userID, templateID, model.CouponSourceRecharge)
		if err != nil {
			fmt.Printf("issue coupon error: %v\n", err)
			continue
		}
		couponIDs = append(couponIDs, coupon.ID)
	}

	if len(couponIDs) > 0 {
		// 标记优惠券已发放
		couponIDsJSON := fmt.Sprintf("%v", couponIDs)
		_ = s.repo.MarkCouponIssued(ctx, recordID, couponIDsJSON)
	}
}

// RefundRecord 退款
func (s *Service) RefundRecord(ctx context.Context, id uint64, reason string) error {
	record, err := s.repo.GetRecordByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get record: %w", err)
	}

	if !record.CanRefund() {
		return errors.New("该订单不可退款")
	}

	// Payment channel refund will be implemented when payment service is fully integrated
	// For now, generate a mock refund provider number
	providerNo := "REFUND_" + record.ProviderTradeNo

	// 标记为已退款
	if err := s.repo.MarkAsRefunded(ctx, id, record.AmountCents, reason, providerNo); err != nil {
		return fmt.Errorf("mark as refunded: %w", err)
	}

	// 扣减钱包余额
	if s.walletRepo != nil {
		wallet, err := s.walletRepo.GetByUserID(ctx, record.UserID)
		if err != nil {
			fmt.Printf("get wallet error: %v\n", err)
		} else if wallet != nil {
			wallet.BalanceCents -= record.TotalCents
			if wallet.BalanceCents < 0 {
				wallet.BalanceCents = 0
			}
			if err := s.walletRepo.Save(ctx, wallet); err != nil {
				fmt.Printf("save wallet error: %v\n", err)
			}
		}
	}

	return nil
}

// GetRechargeStats 获取充值统计
func (s *Service) GetRechargeStats(ctx context.Context) (map[string]any, error) {
	stats, err := s.repo.GetRechargeStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get stats: %w", err)
	}
	return stats, nil
}

// CancelExpiredRecords 取消过期未支付的记录（定时任务）
func (s *Service) CancelExpiredRecords(ctx context.Context) (int64, error) {
	affected, err := s.repo.CancelExpiredRecords(ctx)
	if err != nil {
		return 0, fmt.Errorf("cancel expired: %w", err)
	}
	return affected, nil
}

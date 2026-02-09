package referral

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// TriggerService 推荐奖励触发服务
type TriggerService struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewTriggerService 创建推荐奖励触发服务
func NewTriggerService(db *gorm.DB) *TriggerService {
	return &TriggerService{
		db:     db,
		logger: slog.Default(),
	}
}

// CheckAndIssueReward 检查并发放推荐奖励
// condition: registered | first_order | first_recharge
func (s *TriggerService) CheckAndIssueReward(ctx context.Context, userID uint64, condition string) error {
	// 1. 检查推荐系统是否启用
	var enabledConfig struct {
		ConfigValue string
	}
	err := s.db.WithContext(ctx).
		Table("referral_configs").
		Select("config_value").
		Where("config_key = ?", "enabled").
		First(&enabledConfig).Error
	if err != nil || enabledConfig.ConfigValue != "true" {
		return nil // 推荐系统未启用
	}

	// 2. 查找该用户的推荐记录（作为被推荐人）
	var referral struct {
		ID                uint64
		ReferrerID        uint64
		Status            string
		RefereeCondition  string
		RewardType        string
		RewardAmountCents int64
	}
	err = s.db.WithContext(ctx).
		Table("referrals").
		Select("id, referrer_id, status, referee_condition, reward_type, reward_amount_cents").
		Where("referee_id = ?", userID).
		Where("status = ?", "pending").
		First(&referral).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 该用户不是被推荐用户
		}
		return err
	}

	// 3. 检查是否满足触发条件
	if referral.RefereeCondition != condition {
		return nil // 条件不匹配
	}

	// 4. 发放奖励
	now := time.Now()

	// 开启事务
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 4.1 更新推荐记录状态
		if err := tx.Exec(`
			UPDATE referrals 
			SET status = 'completed', 
			    completed_at = ?,
			    rewarded_at = ?,
			    updated_at = ?
			WHERE id = ?
		`, now, now, now, referral.ID).Error; err != nil {
			return err
		}

		// 4.2 发放奖励给推荐人
		if referral.RewardType == "cash" && referral.RewardAmountCents > 0 {
			// 现金奖励：加入推荐人钱包
			if err := tx.Exec(`
				UPDATE wallets 
				SET balance_cents = balance_cents + ?,
				    updated_at = ?
				WHERE user_id = ?
			`, referral.RewardAmountCents, now, referral.ReferrerID).Error; err != nil {
				return err
			}
		} else if referral.RewardType == "coupon" {
			// 优惠券奖励：从配置获取模板并发放
			var rewardTemplateID uint64
			var rewardCount int
			err := tx.Raw(`
				SELECT 
					CAST(JSON_EXTRACT(config_value, '$.template_id') AS UNSIGNED) as template_id,
					CAST(COALESCE(JSON_EXTRACT(config_value, '$.count'), 1) AS UNSIGNED) as count
				FROM referral_configs 
				WHERE config_key = 'user_reward_coupon_config'
			`).Row().Scan(&rewardTemplateID, &rewardCount)
			if err == nil && rewardTemplateID > 0 {
				for i := 0; i < rewardCount; i++ {
					s.issueCouponFromTemplate(ctx, tx, referral.ReferrerID, rewardTemplateID, "referral")
				}
			}
		}

		// 4.3 创建奖励记录
		if err := tx.Exec(`
			INSERT INTO referral_rewards (
				referral_id, user_id, type, amount_cents, status, issued_at, created_at, updated_at
			) VALUES (?, ?, ?, ?, 'issued', ?, ?, ?)
		`, referral.ID, referral.ReferrerID, referral.RewardType, referral.RewardAmountCents,
			now, now, now).Error; err != nil {
			return err
		}

		// 4.4 发送通知给推荐人
		s.createNotification(ctx, tx, referral.ReferrerID, referral.RewardType, referral.RewardAmountCents)

		s.logger.Info("Referral reward issued",
			"referralId", referral.ID,
			"referrerId", referral.ReferrerID,
			"refereeId", userID,
			"condition", condition,
			"rewardType", referral.RewardType,
			"amount", referral.RewardAmountCents)

		return nil
	})
}

// issueCouponFromTemplate 从模板发放优惠券
func (s *TriggerService) issueCouponFromTemplate(ctx context.Context, tx *gorm.DB, userID uint64, templateID uint64, source string) {
	var tpl struct {
		Name              string
		Type              string
		MinAmountCents    int64
		DeductAmountCents int64
		DiscountRate      float64
		MaxDiscountCents  int64
		Scope             string
		GameIDs           string
		ItemIDs           string
		ValidityDays      int
	}

	err := tx.Table("coupon_templates").
		Select("name, type, min_amount_cents, deduct_amount_cents, discount_rate, max_discount_cents, scope, game_ids, item_ids, validity_days").
		Where("id = ?", templateID).
		First(&tpl).Error

	if err != nil {
		return
	}

	expireAt := time.Now().AddDate(0, 0, tpl.ValidityDays)

	tx.Exec(`
		INSERT INTO coupons (
			template_id, user_id, state, name, type, source,
			min_amount_cents, deduct_amount_cents, discount_rate, max_discount_cents,
			scope, game_ids, item_ids, claimed_at, expire_at, created_at, updated_at
		) VALUES (?, ?, 'available', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, NOW(), NOW())
	`, templateID, userID, tpl.Name, tpl.Type, source,
		tpl.MinAmountCents, tpl.DeductAmountCents, tpl.DiscountRate, tpl.MaxDiscountCents,
		tpl.Scope, tpl.GameIDs, tpl.ItemIDs, expireAt)
}

// createNotification 创建推荐奖励通知
func (s *TriggerService) createNotification(ctx context.Context, tx *gorm.DB, userID uint64, rewardType string, amountCents int64) {
	title := "推荐奖励已到账"
	var message string
	if rewardType == "cash" {
		message = "恭喜您获得推荐奖励，已存入您的钱包"
	} else {
		message = "恭喜您获得推荐奖励，优惠券已发放到您的账户"
	}

	tx.Exec(`
		INSERT INTO notification_events (user_id, title, message, priority, channel, reference_type, created_at, updated_at)
		VALUES (?, ?, ?, 'normal', 'in_app', 'referral', NOW(), NOW())
	`, userID, title, message)
}

// OnUserRegistered 用户注册后调用
func (s *TriggerService) OnUserRegistered(ctx context.Context, userID uint64) error {
	return s.CheckAndIssueReward(ctx, userID, "registered")
}

// OnFirstOrderCompleted 首单完成后调用
func (s *TriggerService) OnFirstOrderCompleted(ctx context.Context, userID uint64) error {
	// 先检查是否是首单
	var orderCount int64
	err := s.db.WithContext(ctx).
		Table("orders").
		Where("user_id = ?", userID).
		Where("status = ?", "completed").
		Count(&orderCount).Error

	if err != nil || orderCount > 1 {
		return nil // 不是首单
	}

	return s.CheckAndIssueReward(ctx, userID, "first_order")
}

// OnFirstRechargeCompleted 首次充值后调用
func (s *TriggerService) OnFirstRechargeCompleted(ctx context.Context, userID uint64) error {
	// 先检查是否是首次充值
	var rechargeCount int64
	err := s.db.WithContext(ctx).
		Table("recharge_records").
		Where("user_id = ?", userID).
		Where("status = ?", "paid").
		Count(&rechargeCount).Error

	if err != nil || rechargeCount > 1 {
		return nil // 不是首次充值
	}

	return s.CheckAndIssueReward(ctx, userID, "first_recharge")
}

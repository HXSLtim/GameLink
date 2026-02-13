package scheduler

import (
	"context"
	"log/slog"
	"time"

	"gamelink/pkg/cache"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// BusinessScheduler 综合业务调度器
// 处理：优惠券过期/锁定释放、VIP过期、活动状态、争议SLA等
type BusinessScheduler struct {
	db     *gorm.DB
	lock   cache.DistributedLock
	logger *slog.Logger
	Cron   *cron.Cron
}

// NewBusinessScheduler 创建综合业务调度器
func NewBusinessScheduler(db *gorm.DB, lock cache.DistributedLock) *BusinessScheduler {
	return &BusinessScheduler{
		db:     db,
		lock:   lock,
		logger: slog.Default(),
		Cron:   cron.New(cron.WithSeconds()),
	}
}

// Start 启动调度器
func (s *BusinessScheduler) Start() {
	// 每分钟执行：优惠券锁定超时释放、支付超时订单取消
	if _, err := s.Cron.AddFunc("0 * * * * *", s.releaseCouponLocksWithLock); err != nil {
		s.logger.Error("Failed to add coupon lock release job", "error", err)
	}

	// 每 5 分钟执行：争议 SLA 超时检查
	if _, err := s.Cron.AddFunc("0 */5 * * * *", s.checkDisputeSLAWithLock); err != nil {
		s.logger.Error("Failed to add dispute SLA check job", "error", err)
	}

	// 每小时执行：VIP过期检查、优惠券过期检查
	if _, err := s.Cron.AddFunc("0 0 * * * *", s.checkVipExpireWithLock); err != nil {
		s.logger.Error("Failed to add VIP expire check job", "error", err)
	}
	if _, err := s.Cron.AddFunc("0 5 * * * *", s.checkCouponExpireWithLock); err != nil {
		s.logger.Error("Failed to add coupon expire check job", "error", err)
	}

	// 每 10 分钟执行：活动状态自动切换
	if _, err := s.Cron.AddFunc("0 */10 * * * *", s.updateActivityStatusWithLock); err != nil {
		s.logger.Error("Failed to add activity status update job", "error", err)
	}

	// 每天凌晨 3 点执行：VIP月度券发放检查
	if _, err := s.Cron.AddFunc("0 0 3 * * *", s.issueVipMonthlyCouponsWithLock); err != nil {
		s.logger.Error("Failed to add VIP monthly coupon job", "error", err)
	}

	s.Cron.Start()
	s.logger.Info("Business scheduler started with all jobs")
}

// Stop 停止调度器
func (s *BusinessScheduler) Stop() {
	s.Cron.Stop()
	s.logger.Info("Business scheduler stopped")
}

// ==================== 优惠券锁定超时释放 ====================

func (s *BusinessScheduler) releaseCouponLocksWithLock() {
	s.executeWithLock("scheduler:coupon:release", time.Minute*5, s.releaseCouponLocks)
}

// releaseCouponLocks 释放超时的优惠券锁定
// 业务规则：支付超时（30分钟）后，自动释放优惠券锁定
func (s *BusinessScheduler) releaseCouponLocks(ctx context.Context) {
	// 查找锁定超过 30 分钟的优惠券
	timeout := time.Now().Add(-30 * time.Minute)

	result := s.db.WithContext(ctx).Exec(`
		UPDATE coupons 
		SET state = 'available', 
		    locked_by_order_id = NULL, 
		    locked_at = NULL,
		    updated_at = NOW()
		WHERE state = 'locked' 
		  AND locked_at IS NOT NULL 
		  AND locked_at < ?
	`, timeout)

	if result.Error != nil {
		s.logger.Error("Failed to release coupon locks", "error", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		s.logger.Info("Released expired coupon locks", "count", result.RowsAffected)
	}
}

// ==================== VIP 过期处理 ====================

func (s *BusinessScheduler) checkVipExpireWithLock() {
	s.executeWithLock("scheduler:vip:expire", time.Hour, s.checkVipExpire)
}

// checkVipExpire 检查并处理 VIP 过期
func (s *BusinessScheduler) checkVipExpire(ctx context.Context) {
	now := time.Now()

	// 1. 查找已过期的 VIP 用户
	var expiredUserIDs []uint64
	err := s.db.WithContext(ctx).
		Table("users").
		Select("id").
		Where("vip_level_id IS NOT NULL").
		Where("vip_expire_at IS NOT NULL").
		Where("vip_expire_at < ?", now).
		Pluck("id", &expiredUserIDs).Error

	if err != nil {
		s.logger.Error("Failed to query expired VIP users", "error", err)
		return
	}

	if len(expiredUserIDs) == 0 {
		return
	}

	// 2. 清除 VIP 状态
	result := s.db.WithContext(ctx).Exec(`
		UPDATE users 
		SET vip_level_id = NULL,
		    updated_at = NOW()
		WHERE id IN ?
		  AND vip_expire_at < ?
	`, expiredUserIDs, now)

	if result.Error != nil {
		s.logger.Error("Failed to clear expired VIP status", "error", result.Error)
		return
	}

	s.logger.Info("Cleared expired VIP status", "count", result.RowsAffected)

	// 3. 发送 VIP 过期通知
	for _, userID := range expiredUserIDs {
		s.createNotification(ctx, userID, "vip_expired", "VIP会员已过期", "您的VIP会员已过期，续费可继续享受专属权益")
	}
}

// ==================== 优惠券过期检查（提前通知）====================

func (s *BusinessScheduler) checkCouponExpireWithLock() {
	s.executeWithLock("scheduler:coupon:expire", time.Hour, s.checkCouponExpire)
}

// checkCouponExpire 检查即将过期的优惠券并发送提醒
func (s *BusinessScheduler) checkCouponExpire(ctx context.Context) {
	// 查找 1 天内过期且未通知的优惠券
	expireThreshold := time.Now().Add(24 * time.Hour)

	var couponsToNotify []struct {
		ID     uint64
		UserID uint64
		Name   string
	}

	err := s.db.WithContext(ctx).
		Table("coupons").
		Select("id, user_id, name").
		Where("state = ?", "available").
		Where("expire_at <= ?", expireThreshold).
		Where("expire_at > ?", time.Now()).
		Where("ext_json NOT LIKE ?", "%expire_notified%"). // 未通知过
		Limit(100).
		Scan(&couponsToNotify).Error

	if err != nil {
		s.logger.Error("Failed to query expiring coupons", "error", err)
		return
	}

	for _, coupon := range couponsToNotify {
		// 发送通知
		s.createNotification(ctx, coupon.UserID, "coupon_expire",
			"优惠券即将过期",
			"您的优惠券「"+coupon.Name+"」将于24小时内过期，请尽快使用")

		// 标记已通知
		s.db.WithContext(ctx).Exec(`
			UPDATE coupons
			SET ext_json = jsonb_set(COALESCE(ext_json, '{}'::jsonb), '{expire_notified}', 'true'::jsonb, true)
			WHERE id = ?
		`, coupon.ID)
	}

	if len(couponsToNotify) > 0 {
		s.logger.Info("Sent coupon expire notifications", "count", len(couponsToNotify))
	}
}

// ==================== 争议 SLA 超时检查 ====================

func (s *BusinessScheduler) checkDisputeSLAWithLock() {
	s.executeWithLock("scheduler:dispute:sla", time.Minute*10, s.checkDisputeSLA)
}

// checkDisputeSLA 检查争议处理 SLA 超时
func (s *BusinessScheduler) checkDisputeSLA(ctx context.Context) {
	now := time.Now()

	// SLA 规则：
	// - 30 分钟未响应 → 升级到高级客服
	// - 2 小时未处理 → 升级到主管

	// 1. 30 分钟未响应的争议 → 标记为需升级
	threshold30min := now.Add(-30 * time.Minute)
	result := s.db.WithContext(ctx).Exec(`
		UPDATE order_disputes 
		SET sla_breached = true,
		    sla_breached_at = COALESCE(sla_breached_at, ?),
		    ext_json = jsonb_set(
		        jsonb_set(COALESCE(ext_json, '{}'::jsonb), '{sla_escalated_at}', to_jsonb(CAST(? AS text)), true),
		        '{sla_level}',
		        to_jsonb('high'::text),
		        true
		    ),
		    updated_at = NOW()
		WHERE status = 'pending'
		  AND created_at < ?
		  AND (sla_breached = false OR sla_breached_at IS NULL)
	`, now, now.Format(time.RFC3339), threshold30min)

	if result.Error != nil {
		s.logger.Error("Failed to escalate disputes (30min)", "error", result.Error)
	} else if result.RowsAffected > 0 {
		s.logger.Warn("Escalated disputes due to SLA breach (30min)", "count", result.RowsAffected)
	}

	// 2. 2 小时未处理的争议 → 标记为紧急
	threshold2h := now.Add(-2 * time.Hour)
	result = s.db.WithContext(ctx).Exec(`
		UPDATE order_disputes 
		SET sla_breached = true,
		    sla_breached_at = COALESCE(sla_breached_at, ?),
		    ext_json = jsonb_set(
		        jsonb_set(COALESCE(ext_json, '{}'::jsonb), '{sla_urgent_at}', to_jsonb(CAST(? AS text)), true),
		        '{sla_level}',
		        to_jsonb('urgent'::text),
		        true
		    ),
		    updated_at = NOW()
		WHERE status = 'pending'
		  AND created_at < ?
	`, now, now.Format(time.RFC3339), threshold2h)

	if result.Error != nil {
		s.logger.Error("Failed to escalate disputes (2h)", "error", result.Error)
	} else if result.RowsAffected > 0 {
		s.logger.Error("URGENT: Disputes exceeded 2h SLA", "count", result.RowsAffected)
		// TODO: 发送告警给主管
	}
}

// ==================== 活动状态自动切换 ====================

func (s *BusinessScheduler) updateActivityStatusWithLock() {
	s.executeWithLock("scheduler:activity:status", time.Minute*15, s.updateActivityStatus)
}

// updateActivityStatus 自动更新活动状态
func (s *BusinessScheduler) updateActivityStatus(ctx context.Context) {
	now := time.Now()

	// 1. 预热期 → 活动中 (preheat → active)
	result := s.db.WithContext(ctx).Exec(`
		UPDATE activities 
		SET status = 'active', updated_at = NOW()
		WHERE status = 'preheat'
		  AND start_at <= ?
	`, now)
	if result.RowsAffected > 0 {
		s.logger.Info("Activities started", "count", result.RowsAffected)
	}

	// 2. 活动中 → 已结束 (active → ended)
	result = s.db.WithContext(ctx).Exec(`
		UPDATE activities 
		SET status = 'ended', updated_at = NOW()
		WHERE status = 'active'
		  AND end_at <= ?
	`, now)
	if result.RowsAffected > 0 {
		s.logger.Info("Activities ended", "count", result.RowsAffected)
	}
}

// ==================== VIP 月度券发放 ====================

func (s *BusinessScheduler) issueVipMonthlyCouponsWithLock() {
	s.executeWithLock("scheduler:vip:monthly", time.Hour*2, s.issueVipMonthlyCoupons)
}

// issueVipMonthlyCoupons 发放 VIP 月度优惠券
func (s *BusinessScheduler) issueVipMonthlyCoupons(ctx context.Context) {
	now := time.Now()

	// 查找需要发放月度券的 VIP 用户
	// 条件：有 VIP 等级 + 未过期 + 本月未发放过
	var usersToIssue []struct {
		UserID                  uint64
		VipLevelID              uint64
		MonthlyCouponTemplateID *uint64
		MonthlyCouponCount      int
	}

	err := s.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.vip_level_id, v.monthly_coupon_template_id, v.monthly_coupon_count
		FROM users u
		JOIN vip_levels v ON u.vip_level_id = v.id
		WHERE u.vip_level_id IS NOT NULL
		  AND (u.vip_expire_at IS NULL OR u.vip_expire_at > ?)
		  AND v.monthly_coupon_template_id IS NOT NULL
		  AND v.monthly_coupon_count > 0
		  AND (
		      u.last_monthly_coupon_at IS NULL
		      OR date_trunc('month', u.last_monthly_coupon_at) < date_trunc('month', CAST(? AS timestamptz))
		  )
		LIMIT 500
	`, now, now).Scan(&usersToIssue).Error

	if err != nil {
		s.logger.Error("Failed to query VIP users for monthly coupons", "error", err)
		return
	}

	issuedCount := 0
	for _, user := range usersToIssue {
		if user.MonthlyCouponTemplateID == nil {
			continue
		}

		// 发放优惠券（调用优惠券服务）
		for i := 0; i < user.MonthlyCouponCount; i++ {
			err := s.issueCouponFromTemplate(ctx, user.UserID, *user.MonthlyCouponTemplateID, "vip_monthly")
			if err != nil {
				s.logger.Error("Failed to issue VIP monthly coupon",
					"userId", user.UserID, "templateId", *user.MonthlyCouponTemplateID, "error", err)
				continue
			}
		}

		// 更新最后发放时间
		s.db.WithContext(ctx).Exec(`UPDATE users SET last_monthly_coupon_at = ? WHERE id = ?`, now, user.UserID)
		issuedCount++

		// 发送通知
		s.createNotification(ctx, user.UserID, "vip_monthly_coupon",
			"VIP月度优惠券已到账",
			"尊敬的VIP会员，您的本月专属优惠券已发放，请查收")
	}

	if issuedCount > 0 {
		s.logger.Info("Issued VIP monthly coupons", "userCount", issuedCount)
	}
}

// ==================== 辅助方法 ====================

// executeWithLock 使用分布式锁执行任务
func (s *BusinessScheduler) executeWithLock(lockKey string, ttl time.Duration, fn func(context.Context)) {
	ctx := context.Background()

	locked, err := s.lock.TryLock(ctx, lockKey, ttl, 1, time.Second)
	if err != nil {
		s.logger.Error("Failed to acquire lock", "key", lockKey, "error", err)
		return
	}

	if !locked {
		s.logger.Debug("Another instance is running task, skipping", "key", lockKey)
		return
	}

	defer func() {
		if unlockErr := s.lock.Unlock(ctx, lockKey); unlockErr != nil {
			s.logger.Error("Failed to release lock", "key", lockKey, "error", unlockErr)
		}
	}()

	fn(ctx)
}

// createNotification 创建通知（简化版，实际应调用通知服务）
func (s *BusinessScheduler) createNotification(ctx context.Context, userID uint64, notifType, title, content string) {
	err := s.db.WithContext(ctx).Exec(`
		INSERT INTO user_notifications (user_id, type, channel, title, content, status, created_at, updated_at)
		VALUES (?, ?, 'in_app', ?, ?, 'unread', NOW(), NOW())
	`, userID, notifType, title, content).Error

	if err != nil {
		s.logger.Error("Failed to create notification", "userId", userID, "type", notifType, "error", err)
	}
}

// issueCouponFromTemplate 从模板发放优惠券（简化版）
func (s *BusinessScheduler) issueCouponFromTemplate(ctx context.Context, userID uint64, templateID uint64, source string) error {
	// 查询模板
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

	err := s.db.WithContext(ctx).
		Table("coupon_templates").
		Select("name, type, min_amount_cents, deduct_amount_cents, discount_rate, max_discount_cents, scope, game_ids, item_ids, validity_days").
		Where("id = ?", templateID).
		First(&tpl).Error

	if err != nil {
		return err
	}

	// 计算过期时间
	expireAt := time.Now().AddDate(0, 0, tpl.ValidityDays)

	// 创建优惠券
	return s.db.WithContext(ctx).Exec(`
		INSERT INTO coupons (
			template_id, user_id, state, name, type, source,
			min_amount_cents, deduct_amount_cents, discount_rate, max_discount_cents,
			scope, game_ids, item_ids, claimed_at, expire_at, created_at, updated_at
		) VALUES (?, ?, 'available', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, NOW(), NOW())
	`, templateID, userID, tpl.Name, tpl.Type, source,
		tpl.MinAmountCents, tpl.DeductAmountCents, tpl.DiscountRate, tpl.MaxDiscountCents,
		tpl.Scope, tpl.GameIDs, tpl.ItemIDs, expireAt).Error
}

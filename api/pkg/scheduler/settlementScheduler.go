package scheduler

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gamelink/internal/service/commission"
	"gamelink/pkg/cache"

	"github.com/robfig/cron/v3"
)

// SettlementScheduler 结算调度器
type SettlementScheduler struct {
	commissionSvc interface{}            // Use interface to allow mocking
	lock          cache.DistributedLock  // Redis distributed lock
	logger        *slog.Logger
	Cron          *cron.Cron  // Exported for testing
}

// NewSettlementScheduler 创建结算调度器
func NewSettlementScheduler(commissionSvc *commission.CommissionService, lock cache.DistributedLock) *SettlementScheduler {
	return &SettlementScheduler{
		commissionSvc: commissionSvc,
		lock:          lock,
		logger:        slog.Default(),
		Cron:          cron.New(),
	}
}

// Start 启动调度器
func (s *SettlementScheduler) Start() {
	// 每月1号凌晨2点执行月度结算
	_, err := s.Cron.AddFunc("0 2 1 * *", s.monthlySettlementWithLock)
	if err != nil {
		s.logger.Error("Failed to add monthly settlement job", "error", err)
		return
	}

	s.Cron.Start()
	s.logger.Info("Settlement scheduler started - will run on 1st of each month at 02:00 with distributed lock")
}

// Stop 停止调度器
func (s *SettlementScheduler) Stop() {
	s.Cron.Stop()
	s.logger.Info("Settlement scheduler stopped")
}

// monthlySettlementWithLock 月度结算任务（带分布式锁）
func (s *SettlementScheduler) monthlySettlementWithLock() {
	ctx := context.Background()
	lockKey := "scheduler:settlement:monthly"

	// Try to acquire distributed lock with 1 hour TTL
	// Monthly settlement should complete within 1 hour
	locked, err := s.lock.TryLock(ctx, lockKey, time.Hour, 1, time.Second)
	if err != nil {
		s.logger.Error("Failed to acquire settlement lock", "error", err, "key", lockKey)
		return
	}

	if !locked {
		s.logger.Info("Another instance is running monthly settlement, skipping", "key", lockKey)
		return
	}

	s.logger.Info("Acquired settlement lock, starting monthly settlement", "key", lockKey)

	// Ensure lock is released
	defer func() {
		if unlockErr := s.lock.Unlock(ctx, lockKey); unlockErr != nil {
			s.logger.Error("Failed to release settlement lock", "error", unlockErr, "key", lockKey)
		} else {
			s.logger.Info("Released settlement lock", "key", lockKey)
		}
	}()

	// Execute the actual settlement
	s.monthlySettlement(ctx)
}

// monthlySettlement 月度结算任务
func (s *SettlementScheduler) monthlySettlement(ctx context.Context) {
	// 结算上个月
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	month := lastMonth.Format("2006-01")

	s.logger.Info("Starting monthly settlement", "month", month)

	// Type assertion to get the commission service
	svc, ok := s.commissionSvc.(interface {
		SettleMonth(ctx context.Context, month string) error
	})
	if !ok {
		s.logger.Error("Invalid commission service type")
		return
	}

	err := svc.SettleMonth(ctx, month)
	if err != nil {
		s.logger.Error("Monthly settlement failed", "month", month, "error", err)
		// Alert notification will be sent when notification service is integrated
		return
	}

	s.logger.Info("Monthly settlement completed successfully", "month", month)
	// Success notification will be sent when notification service is integrated
}

// TriggerSettlement 手动触发结算（用于测试和补偿）
func (s *SettlementScheduler) TriggerSettlement(month string) error {
	ctx := context.Background()
	s.logger.Info("Manual trigger for month", "month", month)

	// Type assertion to get the commission service
	svc, ok := s.commissionSvc.(interface {
		SettleMonth(ctx context.Context, month string) error
	})
	if !ok {
		s.logger.Error("Invalid commission service type")
		return errors.New("invalid commission service type")
	}

	err := svc.SettleMonth(ctx, month)
	if err != nil {
		s.logger.Error("Manual settlement failed", "error", err)
		return err
	}

	s.logger.Info("Manual settlement completed successfully")
	return nil
}

// GetNextRunTime 获取下次运行时间
func (s *SettlementScheduler) GetNextRunTime() time.Time {
	entries := s.Cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}

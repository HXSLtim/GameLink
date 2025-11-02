package scheduler

import (
	"context"
	"log"
	"time"

	"gamelink/internal/service/commission"

	"github.com/robfig/cron/v3"
)

// SettlementScheduler 结算调度器
type SettlementScheduler struct {
	commissionSvc *commission.CommissionService
	cron          *cron.Cron
}

// NewSettlementScheduler 创建结算调度器
func NewSettlementScheduler(commissionSvc *commission.CommissionService) *SettlementScheduler {
	return &SettlementScheduler{
		commissionSvc: commissionSvc,
		cron:          cron.New(),
	}
}

// Start 启动调度器
func (s *SettlementScheduler) Start() {
	// 每月1号凌晨2点执行月度结算
	_, err := s.cron.AddFunc("0 2 1 * *", s.monthlySettlement)
	if err != nil {
		log.Printf("Failed to add monthly settlement job: %v", err)
		return
	}

	s.cron.Start()
	log.Println("Settlement scheduler started - will run on 1st of each month at 02:00")
}

// Stop 停止调度器
func (s *SettlementScheduler) Stop() {
	s.cron.Stop()
	log.Println("Settlement scheduler stopped")
}

// monthlySettlement 月度结算任务
func (s *SettlementScheduler) monthlySettlement() {
	ctx := context.Background()

	// 结算上个月
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	month := lastMonth.Format("2006-01")

	log.Printf("[Settlement] Starting monthly settlement for %s", month)

	err := s.commissionSvc.SettleMonth(ctx, month)
	if err != nil {
		log.Printf("[Settlement] ERROR: Monthly settlement failed for %s: %v", month, err)
		// TODO: 发送告警通知
		return
	}

	log.Printf("[Settlement] SUCCESS: Monthly settlement completed for %s", month)
	// TODO: 发送成功通知
}

// TriggerSettlement 手动触发结算（用于测试和补偿）
func (s *SettlementScheduler) TriggerSettlement(month string) error {
	ctx := context.Background()
	log.Printf("[Settlement] Manual trigger for month: %s", month)

	err := s.commissionSvc.SettleMonth(ctx, month)
	if err != nil {
		log.Printf("[Settlement] Manual settlement failed: %v", err)
		return err
	}

	log.Printf("[Settlement] Manual settlement completed successfully")
	return nil
}

// GetNextRunTime 获取下次运行时间
func (s *SettlementScheduler) GetNextRunTime() time.Time {
	entries := s.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next
	}
	return time.Time{}
}


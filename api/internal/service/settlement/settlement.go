// Package settlement provides T+7 settlement and frozen balance unfreezing functionality.
//
// Business Rules:
// - Orders completed 7+ days ago with no disputes can have their frozen income released
// - Frozen income (FrozenCents) is moved to available balance (BalanceCents)
// - CommissionRecord status is updated from "pending" to "settled"
// - This process runs as a scheduled task (cron job)
package settlement

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	"gamelink/internal/repository/common"
	"gamelink/pkg/apierr"
)

// SettlementService handles T+7 frozen balance unfreezing
type SettlementService struct {
	tx          common.TxManager // 事务管理器
	commissions commissionrepo.CommissionRepository
	wallets     repository.WalletRepository
	players     repository.PlayerRepository // To get user ID from player ID
}

// NewSettlementService creates a new settlement service
func NewSettlementService(
	commissions commissionrepo.CommissionRepository,
	wallets repository.WalletRepository,
	players repository.PlayerRepository,
) *SettlementService {
	return &SettlementService{
		commissions: commissions,
		wallets:     wallets,
		players:     players,
	}
}

// SetTxManager injects a transaction manager.
func (s *SettlementService) SetTxManager(tx common.TxManager) {
	s.tx = tx
}

// UnfreezeResult represents the result of an unfreeze operation
type UnfreezeResult struct {
	ProcessedCount int      `json:"processedCount"`
	SuccessCount   int      `json:"successCount"`
	FailedCount    int      `json:"failedCount"`
	TotalUnfrozen  int64    `json:"totalUnfrozen"` // Total cents unfrozen
	FailedIDs      []uint64 `json:"failedIds"`
	Errors         []string `json:"errors"`
}

// ProcessT7Unfreeze processes all commission records that are 7+ days old
// and moves frozen income to available balance.
//
// This should be called by a scheduled task (e.g., daily cron job).
func (s *SettlementService) ProcessT7Unfreeze(ctx context.Context) (*UnfreezeResult, error) {
	result := &UnfreezeResult{
		FailedIDs: make([]uint64, 0),
		Errors:    make([]string, 0),
	}

	// Calculate the cutoff date (7 days ago)
	cutoffDate := time.Now().AddDate(0, 0, -7)

	// Find all pending commission records created before cutoff date
	status := string(model.SettlementStatusPending)
	records, _, err := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		SettlementStatus: &status,
		DateTo:           &cutoffDate,
		Page:             1,
		PageSize:         1000, // Process in batches
	})
	if err != nil {
		return nil, apierr.InternalError("failed to list pending records").WithDetails(err.Error())
	}

	result.ProcessedCount = len(records)

	// Group records by player for batch wallet updates
	playerIncomes := make(map[uint64]int64)
	recordsByPlayer := make(map[uint64][]model.CommissionRecord)

	for _, record := range records {
		playerIncomes[record.PlayerID] += record.PlayerIncomeCents
		recordsByPlayer[record.PlayerID] = append(recordsByPlayer[record.PlayerID], record)
	}

	// Process each player's frozen income in a transaction
	now := time.Now()
	for playerID, incomeToUnfreeze := range playerIncomes {
		// 使用事务确保钱包解冻和佣金记录更新的原子性
		txErr := s.tx.WithTx(ctx, func(r *common.Repos) error {
			// 使用乐观锁解冻钱包余额
			if err := s.unfreezePlayerIncomeWithRepos(ctx, r, playerID, incomeToUnfreeze); err != nil {
				return err
			}

			// 在同一事务中更新佣金记录
			for i := range recordsByPlayer[playerID] {
				record := &recordsByPlayer[playerID][i]
				record.SettlementStatus = model.SettlementStatusSettled
				record.SettledAt = &now

				if err := r.Commissions.UpdateRecord(ctx, record); err != nil {
					return fmt.Errorf("update record %d: %w", record.ID, err)
				}
			}
			return nil
		})

		if txErr != nil {
			for _, record := range recordsByPlayer[playerID] {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, record.ID)
			}
			result.Errors = append(result.Errors, fmt.Sprintf("player %d: %s", playerID, txErr.Error()))
			continue
		}

		for _, record := range recordsByPlayer[playerID] {
			result.SuccessCount++
			result.TotalUnfrozen += record.PlayerIncomeCents
		}
	}

	return result, nil
}

// unfreezePlayerIncome moves frozen income to available balance for a player (non-transactional).
func (s *SettlementService) unfreezePlayerIncome(ctx context.Context, playerID uint64, amount int64) error {
	return s.unfreezePlayerIncomeWithRepos(ctx, s.tx.Repos(), playerID, amount)
}

// unfreezePlayerIncomeWithRepos moves frozen income to available balance using optimistic locking.
func (s *SettlementService) unfreezePlayerIncomeWithRepos(ctx context.Context, r *common.Repos, playerID uint64, amount int64) error {
	if amount <= 0 {
		return nil
	}

	// Get user ID from player
	var userID uint64
	if s.players != nil {
		player, err := s.players.Get(ctx, playerID)
		if err != nil {
			return fmt.Errorf("get player: %w", err)
		}
		userID = player.UserID
	} else {
		// Fallback: assume playerID == userID (for testing)
		userID = playerID
	}

	// 使用乐观锁确保钱包更新的并发安全
	wallet, err := r.Wallets.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	// Validate frozen balance is sufficient
	if wallet.FrozenCents < amount {
		return fmt.Errorf("insufficient frozen balance: have %d, need %d", wallet.FrozenCents, amount)
	}

	// Move from frozen to available
	wallet.FrozenCents -= amount
	wallet.BalanceCents += amount

	// Save wallet with optimistic locking
	err = r.Wallets.SaveWithOptimisticLock(ctx, wallet)
	if err != nil {
		return fmt.Errorf("update wallet: %w", err)
	}

	return nil
}

// UnfreezeByOrderID manually unfreezes income for a specific order
// This is useful for admin operations or dispute resolution
func (s *SettlementService) UnfreezeByOrderID(ctx context.Context, orderID uint64) error {
	record, err := s.commissions.GetRecordByOrderID(ctx, orderID)
	if err != nil {
		return apierr.NotFound("commission record not found")
	}

	if record.SettlementStatus == model.SettlementStatusSettled {
		return apierr.Conflict("already settled")
	}

	// Unfreeze the income
	err = s.unfreezePlayerIncome(ctx, record.PlayerID, record.PlayerIncomeCents)
	if err != nil {
		return apierr.InternalError("failed to unfreeze").WithDetails(err.Error())
	}

	// Update record status
	now := time.Now()
	record.SettlementStatus = model.SettlementStatusSettled
	record.SettledAt = &now

	return s.commissions.UpdateRecord(ctx, record)
}

// GetPendingSettlementStats returns statistics about pending settlements
func (s *SettlementService) GetPendingSettlementStats(ctx context.Context) (*PendingStats, error) {
	status := string(model.SettlementStatusPending)
	records, total, err := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		SettlementStatus: &status,
		Page:             1,
		PageSize:         10000,
	})
	if err != nil {
		return nil, err
	}

	var totalFrozen int64
	cutoffDate := time.Now().AddDate(0, 0, -7)
	var readyToUnfreeze int64
	var readyCount int64

	for _, record := range records {
		totalFrozen += record.PlayerIncomeCents
		if record.CreatedAt.Before(cutoffDate) {
			readyToUnfreeze += record.PlayerIncomeCents
			readyCount++
		}
	}

	return &PendingStats{
		TotalPendingCount:    total,
		TotalFrozenCents:     totalFrozen,
		ReadyToUnfreezeCount: readyCount,
		ReadyToUnfreezeCents: readyToUnfreeze,
	}, nil
}

// PendingStats represents pending settlement statistics
type PendingStats struct {
	TotalPendingCount    int64 `json:"totalPendingCount"`
	TotalFrozenCents     int64 `json:"totalFrozenCents"`
	ReadyToUnfreezeCount int64 `json:"readyToUnfreezeCount"`
	ReadyToUnfreezeCents int64 `json:"readyToUnfreezeCents"`
}

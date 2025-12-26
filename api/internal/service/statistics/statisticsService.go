package statistics

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Service 统计服务
type Service struct {
	db *gorm.DB
}

// NewService 创建统计服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ============================================================================
// 用户统计
// ============================================================================

// UpdateUserStatistics 更新用户统计（全量重算）
func (s *Service) UpdateUserStatistics(ctx context.Context, userID uint64) error {
	stats := &model.UserStatistics{UserID: userID}

	// 订单统计
	var orderStats struct {
		TotalCount     int
		CompletedCount int
		CanceledCount  int
		RefundCount    int
		TotalSpent     int64
		FirstOrderAt   *time.Time
		LastOrderAt    *time.Time
	}

	err := s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END) as canceled_count,
			SUM(CASE WHEN status = 'refunded' THEN 1 ELSE 0 END) as refund_count,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN total_price_cents ELSE 0 END), 0) as total_spent,
			MIN(created_at) as first_order_at,
			MAX(created_at) as last_order_at
		`).
		Where("user_id = ?", userID).
		Scan(&orderStats).Error
	if err != nil {
		return fmt.Errorf("query order stats: %w", err)
	}

	stats.TotalOrderCount = orderStats.TotalCount
	stats.CompletedOrderCount = orderStats.CompletedCount
	stats.CanceledOrderCount = orderStats.CanceledCount
	stats.RefundOrderCount = orderStats.RefundCount
	stats.TotalSpentCents = orderStats.TotalSpent
	stats.FirstOrderAt = orderStats.FirstOrderAt
	stats.LastOrderAt = orderStats.LastOrderAt

	if stats.TotalOrderCount > 0 {
		stats.AvgOrderAmountCents = stats.TotalSpentCents / int64(stats.CompletedOrderCount)
	}

	// 争议统计
	var disputeStats struct {
		TotalCount int
		WonCount   int
		LostCount  int
	}

	err = s.db.WithContext(ctx).Model(&model.OrderDispute{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN resolution = 'refund' OR resolution = 'partial' THEN 1 ELSE 0 END) as won_count,
			SUM(CASE WHEN resolution = 'reject' THEN 1 ELSE 0 END) as lost_count
		`).
		Where("user_id = ?", userID).
		Scan(&disputeStats).Error
	if err != nil {
		return fmt.Errorf("query dispute stats: %w", err)
	}

	stats.DisputeCount = disputeStats.TotalCount
	stats.DisputeWonCount = disputeStats.WonCount
	stats.DisputeLostCount = disputeStats.LostCount

	// 评价统计
	var reviewStats struct {
		ReviewCount int
		AvgScore    float32
	}

	err = s.db.WithContext(ctx).Model(&model.Review{}).
		Select(`
			COUNT(*) as review_count,
			COALESCE(AVG(score), 0) as avg_score
		`).
		Where("user_id = ? AND status = 'approved'", userID).
		Scan(&reviewStats).Error
	if err != nil {
		return fmt.Errorf("query review stats: %w", err)
	}

	stats.ReviewCount = reviewStats.ReviewCount
	stats.AvgReviewScore = reviewStats.AvgScore

	// Upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(stats).Error
}

// ============================================================================
// 陪玩师统计
// ============================================================================

// UpdatePlayerStatistics 更新陪玩师统计（全量重算）
func (s *Service) UpdatePlayerStatistics(ctx context.Context, playerID uint64) error {
	stats := &model.PlayerStatistics{PlayerID: playerID}

	// 订单统计
	var orderStats struct {
		TotalCount       int
		CompletedCount   int
		CanceledCount    int
		RefundCount      int
		TotalEarnings    int64
		TotalCommission  int64
		TotalServiceMins int
		FirstOrderAt     *time.Time
		LastOrderAt      *time.Time
	}

	err := s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END) as canceled_count,
			SUM(CASE WHEN status = 'refunded' THEN 1 ELSE 0 END) as refund_count,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN player_income_cents ELSE 0 END), 0) as total_earnings,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN commission_cents ELSE 0 END), 0) as total_commission,
			COALESCE(SUM(CASE WHEN status = 'completed' AND completed_at IS NOT NULL AND started_at IS NOT NULL 
				THEN EXTRACT(EPOCH FROM (completed_at - started_at))/60 ELSE 0 END), 0) as total_service_mins,
			MIN(created_at) as first_order_at,
			MAX(created_at) as last_order_at
		`).
		Where("player_id = ?", playerID).
		Scan(&orderStats).Error
	if err != nil {
		return fmt.Errorf("query player order stats: %w", err)
	}

	stats.TotalOrderCount = orderStats.TotalCount
	stats.CompletedOrderCount = orderStats.CompletedCount
	stats.CanceledOrderCount = orderStats.CanceledCount
	stats.RefundOrderCount = orderStats.RefundCount
	stats.TotalEarningsCents = orderStats.TotalEarnings
	stats.TotalCommissionCents = orderStats.TotalCommission
	stats.TotalServiceMinutes = orderStats.TotalServiceMins
	stats.FirstOrderAt = orderStats.FirstOrderAt
	stats.LastOrderAt = orderStats.LastOrderAt

	if stats.CompletedOrderCount > 0 {
		stats.AvgOrderAmountCents = stats.TotalEarningsCents / int64(stats.CompletedOrderCount)
	}

	// 客户统计（去重用户数 + 回头客）
	var customerStats struct {
		TotalCustomers  int
		RepeatCustomers int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(DISTINCT user_id) as total_customers,
			COUNT(*) FILTER (WHERE order_count > 1) as repeat_customers
		FROM (
			SELECT user_id, COUNT(*) as order_count
			FROM orders
			WHERE player_id = ? AND status = 'completed'
			GROUP BY user_id
		) sub
	`, playerID).Scan(&customerStats).Error
	if err != nil {
		return fmt.Errorf("query customer stats: %w", err)
	}

	stats.TotalCustomerCount = customerStats.TotalCustomers
	stats.RepeatCustomerCount = customerStats.RepeatCustomers
	if stats.TotalCustomerCount > 0 {
		stats.RepeatOrderRate = float32(stats.RepeatCustomerCount) / float32(stats.TotalCustomerCount)
	}

	// 争议统计
	var disputeStats struct {
		TotalCount int
		WonCount   int
		LostCount  int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) as total_count,
			SUM(CASE WHEN resolution = 'reject' THEN 1 ELSE 0 END) as won_count,
			SUM(CASE WHEN resolution IN ('refund', 'partial') THEN 1 ELSE 0 END) as lost_count
		FROM order_disputes d
		JOIN orders o ON d.order_id = o.id
		WHERE o.player_id = ?
	`, playerID).Scan(&disputeStats).Error
	if err != nil {
		return fmt.Errorf("query player dispute stats: %w", err)
	}

	stats.DisputeCount = disputeStats.TotalCount
	stats.DisputeWonCount = disputeStats.WonCount
	stats.DisputeLostCount = disputeStats.LostCount

	// 礼物统计
	var giftStats struct {
		GiftCount  int
		GiftAmount int64
	}

	err = s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as gift_count,
			COALESCE(SUM(total_price_cents), 0) as gift_amount
		`).
		Where("recipient_player_id = ? AND status = 'completed'", playerID).
		Scan(&giftStats).Error
	if err != nil {
		return fmt.Errorf("query gift stats: %w", err)
	}

	stats.GiftReceivedCount = giftStats.GiftCount
	stats.GiftReceivedAmountCents = giftStats.GiftAmount

	// 提现统计
	var withdrawStats struct {
		TotalWithdraw   int64
		PendingWithdraw int64
	}

	err = s.db.WithContext(ctx).Model(&model.Withdraw{}).
		Select(`
			COALESCE(SUM(CASE WHEN status = 'completed' THEN amount_cents ELSE 0 END), 0) as total_withdraw,
			COALESCE(SUM(CASE WHEN status IN ('pending', 'approved') THEN amount_cents ELSE 0 END), 0) as pending_withdraw
		`).
		Where("player_id = ?", playerID).
		Scan(&withdrawStats).Error
	if err != nil {
		return fmt.Errorf("query withdraw stats: %w", err)
	}

	stats.TotalWithdrawCents = withdrawStats.TotalWithdraw
	stats.PendingWithdrawCents = withdrawStats.PendingWithdraw

	// Upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "player_id"}},
		UpdateAll: true,
	}).Create(stats).Error
}

// ============================================================================
// 服务项目统计
// ============================================================================

// UpdateServiceItemStatistics 更新服务项目统计
func (s *Service) UpdateServiceItemStatistics(ctx context.Context, itemID uint64) error {
	stats := &model.ServiceItemStatistics{ServiceItemID: itemID}

	// 销售统计
	var salesStats struct {
		TotalCount      int
		TotalAmount     int64
		TotalCommission int64
		RefundCount     int
		RefundAmount    int64
		FirstSoldAt     *time.Time
		LastSoldAt      *time.Time
	}

	err := s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as total_count,
			COALESCE(SUM(total_price_cents), 0) as total_amount,
			COALESCE(SUM(commission_cents), 0) as total_commission,
			SUM(CASE WHEN status = 'refunded' THEN 1 ELSE 0 END) as refund_count,
			COALESCE(SUM(CASE WHEN status = 'refunded' THEN refund_amount_cents ELSE 0 END), 0) as refund_amount,
			MIN(created_at) as first_sold_at,
			MAX(created_at) as last_sold_at
		`).
		Where("item_id = ?", itemID).
		Scan(&salesStats).Error
	if err != nil {
		return fmt.Errorf("query service item sales stats: %w", err)
	}

	stats.TotalSalesCount = salesStats.TotalCount
	stats.TotalSalesAmountCents = salesStats.TotalAmount
	stats.TotalCommissionCents = salesStats.TotalCommission
	stats.RefundCount = salesStats.RefundCount
	stats.RefundAmountCents = salesStats.RefundAmount
	stats.FirstSoldAt = salesStats.FirstSoldAt
	stats.LastSoldAt = salesStats.LastSoldAt

	if stats.TotalSalesCount > 0 {
		stats.RefundRate = float32(stats.RefundCount) / float32(stats.TotalSalesCount)
	}

	// 评价统计
	var reviewStats struct {
		ReviewCount   int
		AvgRating     float32
		FiveStarCount int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) as review_count,
			COALESCE(AVG(r.score), 0) as avg_rating,
			SUM(CASE WHEN r.score = 5 THEN 1 ELSE 0 END) as five_star_count
		FROM reviews r
		JOIN orders o ON r.order_id = o.id
		WHERE o.item_id = ? AND r.status = 'approved'
	`, itemID).Scan(&reviewStats).Error
	if err != nil {
		return fmt.Errorf("query service item review stats: %w", err)
	}

	stats.ReviewCount = reviewStats.ReviewCount
	stats.AvgRating = reviewStats.AvgRating
	stats.FiveStarCount = reviewStats.FiveStarCount

	// Upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "service_item_id"}},
		UpdateAll: true,
	}).Create(stats).Error
}

// ============================================================================
// 游戏统计
// ============================================================================

// UpdateGameStatistics 更新游戏统计
func (s *Service) UpdateGameStatistics(ctx context.Context, gameID uint64) error {
	stats := &model.GameStatistics{GameID: gameID}

	// 订单统计
	var orderStats struct {
		TotalCount      int
		TotalGMV        int64
		TotalCommission int64
	}

	err := s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as total_count,
			COALESCE(SUM(total_price_cents), 0) as total_gmv,
			COALESCE(SUM(commission_cents), 0) as total_commission
		`).
		Where("game_id = ? AND status = 'completed'", gameID).
		Scan(&orderStats).Error
	if err != nil {
		return fmt.Errorf("query game order stats: %w", err)
	}

	stats.TotalOrderCount = orderStats.TotalCount
	stats.TotalGMVCents = orderStats.TotalGMV
	stats.TotalCommissionCents = orderStats.TotalCommission

	if stats.TotalOrderCount > 0 {
		stats.AvgOrderAmountCents = stats.TotalGMVCents / int64(stats.TotalOrderCount)
	}

	// 陪玩师统计
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var playerStats struct {
		TotalPlayers  int
		ActivePlayers int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(DISTINCT p.id) as total_players,
			COUNT(DISTINCT CASE WHEN o.created_at > ? THEN p.id END) as active_players
		FROM players p
		LEFT JOIN orders o ON o.player_id = p.id AND o.game_id = ?
		WHERE p.main_game_id = ?
	`, sevenDaysAgo, gameID, gameID).Scan(&playerStats).Error
	if err != nil {
		return fmt.Errorf("query game player stats: %w", err)
	}

	stats.TotalPlayerCount = playerStats.TotalPlayers
	stats.ActivePlayerCount = playerStats.ActivePlayers

	// 用户统计
	var userStats struct {
		TotalUsers  int
		ActiveUsers int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(DISTINCT user_id) as total_users,
			COUNT(DISTINCT CASE WHEN created_at > ? THEN user_id END) as active_users
		FROM orders
		WHERE game_id = ?
	`, sevenDaysAgo, gameID).Scan(&userStats).Error
	if err != nil {
		return fmt.Errorf("query game user stats: %w", err)
	}

	stats.TotalUserCount = userStats.TotalUsers
	stats.ActiveUserCount = userStats.ActiveUsers

	// 评价统计
	var avgRating float32
	err = s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(AVG(r.score), 0)
		FROM reviews r
		JOIN orders o ON r.order_id = o.id
		WHERE o.game_id = ? AND r.status = 'approved'
	`, gameID).Scan(&avgRating).Error
	if err != nil {
		return fmt.Errorf("query game review stats: %w", err)
	}

	stats.AvgRating = avgRating

	// Upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "game_id"}},
		UpdateAll: true,
	}).Create(stats).Error
}

// ============================================================================
// 平台日统计
// ============================================================================

// UpdatePlatformDailyStatistics 更新平台日统计
func (s *Service) UpdatePlatformDailyStatistics(ctx context.Context, date time.Time) error {
	// 标准化日期（去掉时分秒）
	statDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nextDate := statDate.AddDate(0, 0, 1)

	stats := &model.PlatformStatistics{StatDate: statDate}

	// 订单统计
	var orderStats struct {
		TotalCount      int
		CompletedCount  int
		CanceledCount   int
		TotalGMV        int64
		TotalCommission int64
		RefundAmount    int64
	}

	err := s.db.WithContext(ctx).Model(&model.Order{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'canceled' THEN 1 ELSE 0 END) as canceled_count,
			COALESCE(SUM(total_price_cents), 0) as total_gmv,
			COALESCE(SUM(commission_cents), 0) as total_commission,
			COALESCE(SUM(refund_amount_cents), 0) as refund_amount
		`).
		Where("created_at >= ? AND created_at < ?", statDate, nextDate).
		Scan(&orderStats).Error
	if err != nil {
		return fmt.Errorf("query daily order stats: %w", err)
	}

	stats.DailyOrderCount = orderStats.TotalCount
	stats.DailyCompletedCount = orderStats.CompletedCount
	stats.DailyCanceledCount = orderStats.CanceledCount
	stats.DailyGMVCents = orderStats.TotalGMV
	stats.DailyCommissionCents = orderStats.TotalCommission
	stats.DailyRefundAmountCents = orderStats.RefundAmount

	// 用户统计
	var userStats struct {
		NewUsers    int
		ActiveUsers int
		PayingUsers int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			(SELECT COUNT(*) FROM users WHERE created_at >= ? AND created_at < ?) as new_users,
			(SELECT COUNT(DISTINCT user_id) FROM orders WHERE created_at >= ? AND created_at < ?) as active_users,
			(SELECT COUNT(DISTINCT user_id) FROM orders WHERE created_at >= ? AND created_at < ? AND status = 'completed') as paying_users
	`, statDate, nextDate, statDate, nextDate, statDate, nextDate).Scan(&userStats).Error
	if err != nil {
		return fmt.Errorf("query daily user stats: %w", err)
	}

	stats.DailyNewUserCount = userStats.NewUsers
	stats.DailyActiveUserCount = userStats.ActiveUsers
	stats.DailyPayingUserCount = userStats.PayingUsers

	// 陪玩师统计
	var playerStats struct {
		NewPlayers    int
		ActivePlayers int
	}

	err = s.db.WithContext(ctx).Raw(`
		SELECT 
			(SELECT COUNT(*) FROM players WHERE created_at >= ? AND created_at < ?) as new_players,
			(SELECT COUNT(DISTINCT player_id) FROM orders WHERE created_at >= ? AND created_at < ? AND player_id IS NOT NULL) as active_players
	`, statDate, nextDate, statDate, nextDate).Scan(&playerStats).Error
	if err != nil {
		return fmt.Errorf("query daily player stats: %w", err)
	}

	stats.DailyNewPlayerCount = playerStats.NewPlayers
	stats.DailyActivePlayerCount = playerStats.ActivePlayers

	// 提现统计
	var withdrawStats struct {
		TotalAmount int64
		TotalCount  int
	}

	err = s.db.WithContext(ctx).Model(&model.Withdraw{}).
		Select(`
			COALESCE(SUM(amount_cents), 0) as total_amount,
			COUNT(*) as total_count
		`).
		Where("created_at >= ? AND created_at < ?", statDate, nextDate).
		Scan(&withdrawStats).Error
	if err != nil {
		return fmt.Errorf("query daily withdraw stats: %w", err)
	}

	stats.DailyWithdrawCents = withdrawStats.TotalAmount
	stats.DailyWithdrawCount = withdrawStats.TotalCount

	// 争议统计
	var disputeStats struct {
		TotalCount     int
		ResolvedCount  int
		SLABreachCount int
	}

	err = s.db.WithContext(ctx).Model(&model.OrderDispute{}).
		Select(`
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'resolved' THEN 1 ELSE 0 END) as resolved_count,
			SUM(CASE WHEN sla_breached = true THEN 1 ELSE 0 END) as sla_breach_count
		`).
		Where("created_at >= ? AND created_at < ?", statDate, nextDate).
		Scan(&disputeStats).Error
	if err != nil {
		return fmt.Errorf("query daily dispute stats: %w", err)
	}

	stats.DailyDisputeCount = disputeStats.TotalCount
	stats.DailyResolvedCount = disputeStats.ResolvedCount
	stats.DailySLABreachCount = disputeStats.SLABreachCount

	// Upsert
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "stat_date"}},
		UpdateAll: true,
	}).Create(stats).Error
}

// ============================================================================
// 批量更新
// ============================================================================

// UpdateAllUserStatistics 更新所有用户统计
func (s *Service) UpdateAllUserStatistics(ctx context.Context) error {
	var userIDs []uint64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Pluck("id", &userIDs).Error; err != nil {
		return fmt.Errorf("get user ids: %w", err)
	}

	for _, userID := range userIDs {
		if err := s.UpdateUserStatistics(ctx, userID); err != nil {
			return fmt.Errorf("update user %d stats: %w", userID, err)
		}
	}
	return nil
}

// UpdateAllPlayerStatistics 更新所有陪玩师统计
func (s *Service) UpdateAllPlayerStatistics(ctx context.Context) error {
	var playerIDs []uint64
	if err := s.db.WithContext(ctx).Model(&model.Player{}).Pluck("id", &playerIDs).Error; err != nil {
		return fmt.Errorf("get player ids: %w", err)
	}

	for _, playerID := range playerIDs {
		if err := s.UpdatePlayerStatistics(ctx, playerID); err != nil {
			return fmt.Errorf("update player %d stats: %w", playerID, err)
		}
	}
	return nil
}

// UpdateAllServiceItemStatistics 更新所有服务项目统计
func (s *Service) UpdateAllServiceItemStatistics(ctx context.Context) error {
	var itemIDs []uint64
	if err := s.db.WithContext(ctx).Model(&model.ServiceItem{}).Pluck("id", &itemIDs).Error; err != nil {
		return fmt.Errorf("get service item ids: %w", err)
	}

	for _, itemID := range itemIDs {
		if err := s.UpdateServiceItemStatistics(ctx, itemID); err != nil {
			return fmt.Errorf("update service item %d stats: %w", itemID, err)
		}
	}
	return nil
}

// UpdateAllGameStatistics 更新所有游戏统计
func (s *Service) UpdateAllGameStatistics(ctx context.Context) error {
	var gameIDs []uint64
	if err := s.db.WithContext(ctx).Model(&model.Game{}).Pluck("id", &gameIDs).Error; err != nil {
		return fmt.Errorf("get game ids: %w", err)
	}

	for _, gameID := range gameIDs {
		if err := s.UpdateGameStatistics(ctx, gameID); err != nil {
			return fmt.Errorf("update game %d stats: %w", gameID, err)
		}
	}
	return nil
}

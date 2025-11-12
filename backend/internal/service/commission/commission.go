package commission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
)

var (
	// ErrNotFound 记录不存在
	ErrNotFound = repository.ErrNotFound
	// ErrValidation 表示输入校验失败
	ErrValidation = errors.New("validation failed")
	// ErrAlreadyRecorded 抽成已记录
	ErrAlreadyRecorded = errors.New("commission already recorded")
	// ErrAlreadySettled 已经结算
	ErrAlreadySettled = errors.New("already settled")
)

// CommissionService 抽成服务
type CommissionService struct {
	commissions commissionrepo.CommissionRepository
	orders      repository.OrderRepository
	players     repository.PlayerRepository
}

// NewCommissionService 创建抽成服务
func NewCommissionService(
	commissions commissionrepo.CommissionRepository,
	orders repository.OrderRepository,
	players repository.PlayerRepository,
) *CommissionService {
	return &CommissionService{
		commissions: commissions,
		orders:      orders,
		players:     players,
	}
}

// CalculateCommission 计算订单抽成（便捷方法：通过orderID）
func (s *CommissionService) CalculateCommission(ctx context.Context, orderID uint64) (*CommissionCalculation, error) {
	// 获取订单
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 调用完整的三层抽成计算
	return s.CalculateOrderCommission(ctx, order)
}

// RecordCommission 记录订单抽成
func (s *CommissionService) RecordCommission(ctx context.Context, orderID uint64) error {
	// 1. 检查是否已记录
	existing, _ := s.commissions.GetRecordByOrderID(ctx, orderID)
	if existing != nil {
		return ErrAlreadyRecorded
	}

	// 2. 计算抽成
	calc, err := s.CalculateCommission(ctx, orderID)
	if err != nil {
		return err
	}

	// 3. 获取订单信息
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return err
	}

	// 4. 创建抽成记录
	now := time.Now()
	playerID := order.GetPlayerID()
	if playerID == 0 {
		return errors.New("order has no player assigned")
	}
	
	record := &model.CommissionRecord{
		OrderID:            orderID,
		PlayerID:           playerID,
		TotalAmountCents:   calc.TotalAmountCents,
		CommissionRate:     calc.CommissionRate,
		CommissionCents:    calc.CommissionCents,
		PlayerIncomeCents:  calc.PlayerIncomeCents,
		SettlementStatus:   "pending",
		SettlementMonth:    now.Format("2006-01"),
	}

	return s.commissions.CreateRecord(ctx, record)
}

// PlayerMonthStats 玩家月度统计
type PlayerMonthStats struct {
	PlayerID             uint64
	OrderCount           int64
	TotalAmountCents     int64
	TotalCommissionCents int64
	TotalIncomeCents     int64
}

// SettleMonth 月度结算
func (s *CommissionService) SettleMonth(ctx context.Context, month string) error {
	// 1. 检查是否已经结算过
	settlements, _, err := s.commissions.ListSettlements(ctx, commissionrepo.SettlementListOptions{
		SettlementMonth: &month,
		Page:            1,
		PageSize:        1,
	})
	if err != nil {
		return err
	}
	if len(settlements) > 0 {
		return ErrAlreadySettled
	}

	// 2. 获取该月所有待结算记录
	status := "pending"
	records, _, err := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		SettlementMonth:  &month,
		SettlementStatus: &status,
		Page:             1,
		PageSize:         10000,
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records to settle for month %s", month)
	}

	// 3. 按陪玩师分组统计
	playerStats := make(map[uint64]*PlayerMonthStats)
	for _, record := range records {
		stats, exists := playerStats[record.PlayerID]
		if !exists {
			stats = &PlayerMonthStats{PlayerID: record.PlayerID}
			playerStats[record.PlayerID] = stats
		}
		stats.OrderCount++
		stats.TotalAmountCents += record.TotalAmountCents
		stats.TotalCommissionCents += record.CommissionCents
		stats.TotalIncomeCents += record.PlayerIncomeCents
	}

	// 4. 为每个陪玩师创建月度结算记录
	for _, stats := range playerStats {
		settlement := &model.MonthlySettlement{
			PlayerID:             stats.PlayerID,
			SettlementMonth:      month,
			TotalOrderCount:      stats.OrderCount,
			TotalAmountCents:     stats.TotalAmountCents,
			TotalCommissionCents: stats.TotalCommissionCents,
			TotalIncomeCents:     stats.TotalIncomeCents,
			BonusCents:           0, // 奖金在排名系统中计算
			FinalIncomeCents:     stats.TotalIncomeCents,
			Status:               "pending",
		}

		err := s.commissions.CreateSettlement(ctx, settlement)
		if err != nil {
			return fmt.Errorf("failed to create settlement for player %d: %w", stats.PlayerID, err)
		}
	}

	// 5. 更新抽成记录状态
	now := time.Now()
	for _, record := range records {
		record.SettlementStatus = "settled"
		record.SettledAt = &now
		if err := s.commissions.UpdateRecord(ctx, &record); err != nil {
			return fmt.Errorf("failed to update record %d: %w", record.ID, err)
		}
	}

	return nil
}

// GetPlayerCommissionSummary 获取玩家抽成汇总
func (s *CommissionService) GetPlayerCommissionSummary(ctx context.Context, playerID uint64, month string) (*CommissionSummaryResponse, error) {
	// 获取月度收入
	income, err := s.commissions.GetPlayerMonthlyIncome(ctx, playerID, month)
	if err != nil {
		return nil, err
	}

	// 获取抽成记录
	records, total, err := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		PlayerID: &playerID,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, err
	}

	var totalCommission int64
	var totalIncome int64
	if len(records) > 0 {
		// 统计所有抽成记录
	allRecords, _, _ := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
			PlayerID: &playerID,
			Page:     1,
			PageSize: 10000,
		})
		for _, r := range allRecords {
			totalCommission += r.CommissionCents
			totalIncome += r.PlayerIncomeCents
		}
	}

	return &CommissionSummaryResponse{
		MonthlyIncome:   income,
		TotalCommission: totalCommission,
		TotalIncome:     totalIncome,
		TotalOrders:     total,
	}, nil
}

// CommissionSummaryResponse 抽成汇总响应
type CommissionSummaryResponse struct {
	MonthlyIncome   int64 `json:"monthlyIncome"`
	TotalCommission int64 `json:"totalCommission"`
	TotalIncome     int64 `json:"totalIncome"`
	TotalOrders     int64 `json:"totalOrders"`
}

// GetCommissionRecords 获取抽成记录列表
func (s *CommissionService) GetCommissionRecords(ctx context.Context, playerID uint64, page, pageSize int) (*CommissionRecordListResponse, error) {
	records, total, err := s.commissions.ListRecords(ctx, commissionrepo.CommissionRecordListOptions{
		PlayerID: &playerID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	recordDTOs := make([]CommissionRecordDTO, 0, len(records))
	for _, r := range records {
		recordDTOs = append(recordDTOs, CommissionRecordDTO{
			ID:                 r.ID,
			OrderID:            r.OrderID,
			TotalAmountCents:   r.TotalAmountCents,
			CommissionRate:     r.CommissionRate,
			CommissionCents:    r.CommissionCents,
			PlayerIncomeCents:  r.PlayerIncomeCents,
			SettlementStatus:   r.SettlementStatus,
			SettlementMonth:    r.SettlementMonth,
			CreatedAt:          r.CreatedAt,
		})
	}

	return &CommissionRecordListResponse{
		Records: recordDTOs,
		Total:   total,
	}, nil
}

// CommissionRecordDTO 抽成记录DTO
type CommissionRecordDTO struct {
	ID                 uint64    `json:"id"`
	OrderID            uint64    `json:"orderId"`
	TotalAmountCents   int64     `json:"totalAmountCents"`
	CommissionRate     int       `json:"commissionRate"`
	CommissionCents    int64     `json:"commissionCents"`
	PlayerIncomeCents  int64     `json:"playerIncomeCents"`
	SettlementStatus   string    `json:"settlementStatus"`
	SettlementMonth    string    `json:"settlementMonth"`
	CreatedAt          time.Time `json:"createdAt"`
}

// CommissionRecordListResponse 抽成记录列表响应
type CommissionRecordListResponse struct {
	Records []CommissionRecordDTO `json:"records"`
	Total   int64                 `json:"total"`
}

// GetMonthlySettlements 获取月度结算列表
func (s *CommissionService) GetMonthlySettlements(ctx context.Context, playerID uint64, page, pageSize int) (*SettlementListResponse, error) {
	settlements, total, err := s.commissions.ListSettlements(ctx, commissionrepo.SettlementListOptions{
		PlayerID: &playerID,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	settlementDTOs := make([]SettlementDTO, 0, len(settlements))
	for _, s := range settlements {
		settlementDTOs = append(settlementDTOs, SettlementDTO{
			ID:                   s.ID,
			SettlementMonth:      s.SettlementMonth,
			TotalOrderCount:      s.TotalOrderCount,
			TotalAmountCents:     s.TotalAmountCents,
			TotalCommissionCents: s.TotalCommissionCents,
			TotalIncomeCents:     s.TotalIncomeCents,
			BonusCents:           s.BonusCents,
			FinalIncomeCents:     s.FinalIncomeCents,
			Status:               s.Status,
			CreatedAt:            s.CreatedAt,
			SettledAt:            s.SettledAt,
		})
	}

	return &SettlementListResponse{
		Settlements: settlementDTOs,
		Total:       total,
	}, nil
}

// SettlementDTO 结算DTO
type SettlementDTO struct {
	ID                   uint64     `json:"id"`
	SettlementMonth      string     `json:"settlementMonth"`
	TotalOrderCount      int64      `json:"totalOrderCount"`
	TotalAmountCents     int64      `json:"totalAmountCents"`
	TotalCommissionCents int64      `json:"totalCommissionCents"`
	TotalIncomeCents     int64      `json:"totalIncomeCents"`
	BonusCents           int64      `json:"bonusCents"`
	FinalIncomeCents     int64      `json:"finalIncomeCents"`
	Status               string     `json:"status"`
	CreatedAt            time.Time  `json:"createdAt"`
	SettledAt            *time.Time `json:"settledAt"`
}

// SettlementListResponse 结算列表响应
type SettlementListResponse struct {
	Settlements []SettlementDTO `json:"settlements"`
	Total       int64           `json:"total"`
}

// CreateCommissionRule 创建抽成规则（管理员）
func (s *CommissionService) CreateCommissionRule(ctx context.Context, req CreateCommissionRuleRequest) (*model.CommissionRule, error) {
	// 验证抽成比例
	if req.Rate < 0 || req.Rate > 100 {
		return nil, fmt.Errorf("commission rate must be between 0 and 100")
	}

	rule := &model.CommissionRule{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Rate:        req.Rate,
		IsActive:    true,
		GameID:      req.GameID,
		PlayerID:    req.PlayerID,
		ServiceType: req.ServiceType,
	}

	if err := s.commissions.CreateRule(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// CreateCommissionRuleRequest 创建抽成规则请求
type CreateCommissionRuleRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`
	Description string  `json:"description"`
	Type        string  `json:"type" binding:"required,oneof=default special gift"`
	Rate        int     `json:"rate" binding:"required,min=0,max=100"`
	GameID      *uint64 `json:"gameId"`
	PlayerID    *uint64 `json:"playerId"`
	ServiceType *string `json:"serviceType"`
}

// UpdateCommissionRule 更新抽成规则（管理员）
func (s *CommissionService) UpdateCommissionRule(ctx context.Context, id uint64, req UpdateCommissionRuleRequest) error {
	rule, err := s.commissions.GetRule(ctx, id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = *req.Description
	}
	if req.Rate != nil {
		if *req.Rate < 0 || *req.Rate > 100 {
			return fmt.Errorf("commission rate must be between 0 and 100")
		}
		rule.Rate = *req.Rate
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	return s.commissions.UpdateRule(ctx, rule)
}

// UpdateCommissionRuleRequest 更新抽成规则请求
type UpdateCommissionRuleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Rate        *int    `json:"rate"`
	IsActive    *bool   `json:"isActive"`
}

// GetPlatformStats 获取平台统计（管理员）
func (s *CommissionService) GetPlatformStats(ctx context.Context, month string) (*PlatformStatsResponse, error) {
	stats, err := s.commissions.GetMonthlyStats(ctx, month)
	if err != nil {
		return nil, err
	}

	return &PlatformStatsResponse{
		Month:             month,
		TotalOrders:       stats.TotalOrders,
		TotalIncome:       stats.TotalIncome,
		TotalCommission:   stats.TotalCommission,
		TotalPlayerIncome: stats.TotalPlayerIncome,
	}, nil
}

// PlatformStatsResponse 平台统计响应
type PlatformStatsResponse struct {
	Month             string `json:"month"`
	TotalOrders       int64  `json:"totalOrders"`
	TotalIncome       int64  `json:"totalIncome"`
	TotalCommission   int64  `json:"totalCommission"`
	TotalPlayerIncome int64  `json:"totalPlayerIncome"`
}

// ============================================================================
// 三层抽成计算逻辑
// ============================================================================

// CalculateOrderCommission 计算订单的实际抽成（三层取最低）
//
// 抽成来源：
// 1. 服务项目抽成 (service_items.commission_rate)
// 2. 陪玩师专属抽成 (commission_rules WHERE player_id = ?)
// 3. 排名抽成 (基于上月排名，前N名享受优惠)
//
// 实际抽成 = MIN(服务项目抽成, 陪玩师抽成, 排名抽成)
//
// 注意：礼物订单不参与排名优惠
func (s *CommissionService) CalculateOrderCommission(ctx context.Context, order *model.Order) (*CommissionCalculation, error) {
	var candidateRates []CommissionCandidate

	// 1. 获取服务项目抽成（基础抽成）
	serviceItem, err := s.getServiceItemForOrder(ctx, order.ItemID)
	if err == nil && serviceItem != nil {
		candidateRates = append(candidateRates, CommissionCandidate{
			Source: "服务项目",
			Rate:   int(serviceItem.CommissionRate * 100),
			Detail: serviceItem.Name,
		})
	}

	// 2. 查找陪玩师专属抽成规则
	playerID := order.GetPlayerID()
	if playerID > 0 {
		playerRule, err := s.commissions.GetRuleForOrder(ctx, order.GameID, order.PlayerID, nil)
		if err == nil && playerRule != nil {
			candidateRates = append(candidateRates, CommissionCandidate{
				Source: "陪玩师专属",
				Rate:   playerRule.Rate,
				Detail: playerRule.Name,
			})
		}
	}

	// 3. 查找排名抽成（基于上月排名）
	// 注意：礼物订单不参与排名优惠
	if !order.IsGiftOrder() && playerID > 0 {
		rankingRate, rankingDetail := s.getRankingCommissionRate(ctx, playerID)
		if rankingRate > 0 && rankingRate < 100 {
			candidateRates = append(candidateRates, CommissionCandidate{
				Source: "排名优惠",
				Rate:   rankingRate,
				Detail: rankingDetail,
			})
		}
	}

	// 如果没有任何规则，使用默认20%
	if len(candidateRates) == 0 {
		defaultRule, err := s.commissions.GetDefaultRule(ctx)
		defaultRate := 20 // 默认20%
		defaultDetail := "平台默认20%抽成"
		
		if err == nil && defaultRule != nil {
			defaultRate = defaultRule.Rate
			defaultDetail = defaultRule.Name
		}
		
		candidateRates = append(candidateRates, CommissionCandidate{
			Source: "默认规则",
			Rate:   defaultRate,
			Detail: defaultDetail,
		})
	}

	// 取最低抽成比例
	finalRate := selectLowestRate(candidateRates)
	totalAmount := order.TotalPriceCents
	commissionCents := totalAmount * int64(finalRate.Rate) / 100
	playerIncome := totalAmount - commissionCents

	return &CommissionCalculation{
		OrderID:           order.ID,
		TotalAmountCents:  totalAmount,
		CommissionRate:    finalRate.Rate,
		CommissionCents:   commissionCents,
		PlayerIncomeCents: playerIncome,
		AppliedRule:       finalRate.Source,
		AppliedRuleDetail: finalRate.Detail,
		CandidateRates:    candidateRates,
	}, nil
}

// CommissionCandidate 抽成候选项
type CommissionCandidate struct {
	Source string `json:"source"` // 来源：服务项目/陪玩师专属/排名优惠/默认规则
	Rate   int    `json:"rate"`   // 抽成比例
	Detail string `json:"detail"` // 详细说明
}

// CommissionCalculation 抽成计算结果
type CommissionCalculation struct {
	OrderID           uint64                `json:"orderId"`
	TotalAmountCents  int64                 `json:"totalAmountCents"`
	CommissionRate    int                   `json:"commissionRate"`    // 实际使用的抽成比例
	CommissionCents   int64                 `json:"commissionCents"`   // 平台抽成
	PlayerIncomeCents int64                 `json:"playerIncomeCents"` // 陪玩师收入
	AppliedRule       string                `json:"appliedRule"`       // 实际应用的规则
	AppliedRuleDetail string                `json:"appliedRuleDetail"` // 规则详情
	CandidateRates    []CommissionCandidate `json:"candidateRates"`    // 所有候选抽成
}

// getRankingCommissionRate 获取陪玩师的排名抽成比例
func (s *CommissionService) getRankingCommissionRate(ctx context.Context, playerID uint64) (int, string) {
	// TODO: 查询陪玩师上月排名
	// 1. 获取上月（例如：当前是2月，查询11月的排名）
	// lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	// 2. 查询该陪玩师在上月的排名
	// rankings, err := s.rankings.GetPlayerRankings(ctx, playerID, lastMonth)

	// 3. 查找对应的排名抽成配置
	// for _, ranking := range rankings {
	//     config := s.findRankingCommissionConfig(ctx, ranking.RankingType, lastMonth)
	//     if config != nil {
	//         // 解析JSON规则
	//         var rules []model.RankingCommissionRule
	//         json.Unmarshal([]byte(config.RulesJSON), &rules)
	//
	//         // 查找该排名对应的抽成
	//         for _, rule := range rules {
	//             if ranking.Rank >= rule.RankStart && ranking.Rank <= rule.RankEnd {
	//                 return rule.CommissionRate, fmt.Sprintf("%s第%d名", config.Name, ranking.Rank)
	//             }
	//         }
	//     }
	// }

    return 0, "" // 暂时返回0，等待排名系统完整实现
}

// getServiceItemForOrder 获取订单的服务项
func (s *CommissionService) getServiceItemForOrder(ctx context.Context, itemID uint64) (*model.ServiceItem, error) {
	// TODO: 需要注入 ServiceItemRepository
	// return s.serviceItems.Get(ctx, itemID)
	return nil, nil // 暂时返回nil
}

// selectLowestRate 选择最低抽成比例
func selectLowestRate(candidates []CommissionCandidate) CommissionCandidate {
	if len(candidates) == 0 {
		return CommissionCandidate{
			Source: "默认规则",
			Rate:   20,
			Detail: "平台默认20%抽成",
		}
	}

	lowest := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.Rate < lowest.Rate {
			lowest = candidate
		}
	}

	return lowest
}

// ParseRankingCommissionRules 解析排名抽成规则JSON
func ParseRankingCommissionRules(rulesJSON string) ([]model.RankingCommissionRule, error) {
	var rules []model.RankingCommissionRule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// FindCommissionRateForRank 根据排名查找对应的抽成比例
func FindCommissionRateForRank(rules []model.RankingCommissionRule, rank int) int {
	for _, rule := range rules {
		if rank >= rule.RankStart && rank <= rule.RankEnd {
			return rule.CommissionRate
		}
	}
    return 0 // 不在任何规则范围内
}

// ValidateRankingRules 验证排名规则的合法性
func ValidateRankingRules(rules []model.RankingCommissionRule) error {
	for _, rule := range rules {
		if rule.RankStart < 1 || rule.RankEnd < rule.RankStart {
			return ErrValidation
		}
		if rule.CommissionRate < 0 || rule.CommissionRate > 100 {
			return ErrValidation
		}
	}

	// 检查是否有重叠
	for i := 0; i < len(rules); i++ {
		for j := i + 1; j < len(rules); j++ {
			if rangesOverlap(rules[i].RankStart, rules[i].RankEnd, rules[j].RankStart, rules[j].RankEnd) {
				return ErrValidation
			}
		}
	}

	return nil
}

// rangesOverlap 检查两个范围是否重叠
func rangesOverlap(start1, end1, start2, end2 int) bool {
	return math.Max(float64(start1), float64(start2)) <= math.Min(float64(end1), float64(end2))
}


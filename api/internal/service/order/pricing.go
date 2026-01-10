package order

import (
	"context"

	"gamelink/internal/model"
)

// calculateOrderPricing 计算订单价格和抽成相关金额
// 使用 CommissionRule 系统获取动态抽成率
func (s *OrderService) calculateOrderPricing(player *model.Player, req CreateOrderRequest) (totalPrice int64, commissionCents int64, playerIncomeCents int64) {
	// 从陪玩师时薪计算价格
	hourlyRate := player.HourlyRateCents
	totalPrice = int64(float32(hourlyRate) * req.DurationHours)

	// 获取抽成率（使用 CommissionRule 系统）
	commissionRate := s.getCommissionRate(context.Background(), &req.GameID, &player.ID)
	commissionCents = totalPrice * int64(commissionRate) / 100
	playerIncomeCents = totalPrice - commissionCents

	return
}

// getCommissionRate 获取抽成率
// 优先级：陪玩师个人规则 > 游戏规则 > 默认规则
func (s *OrderService) getCommissionRate(ctx context.Context, gameID *uint64, playerID *uint64) int {
	const defaultRate = 20 // 默认 20%

	if s.commissions == nil {
		return defaultRate
	}

	// 尝试获取适用的抽成规则
	rule, err := s.commissions.GetRuleForOrder(ctx, gameID, playerID, nil)
	if err == nil && rule != nil {
		return rule.Rate
	}

	// 尝试获取默认规则
	defaultRule, err := s.commissions.GetDefaultRule(ctx)
	if err == nil && defaultRule != nil {
		return defaultRule.Rate
	}

	return defaultRate
}

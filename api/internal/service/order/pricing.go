package order

import "gamelink/internal/model"

// calculateOrderPricing 计算订单价格和抽成相关金额
func (s *OrderService) calculateOrderPricing(player *model.Player, req CreateOrderRequest) (totalPrice int64, commissionCents int64, playerIncomeCents int64) {
	// 从陪玩师时薪计算价格（保持原有简化版本逻辑）
	hourlyRate := player.HourlyRateCents
	totalPrice = int64(float32(hourlyRate) * req.DurationHours)
	// 默认抽成20%
	commissionRate := int64(20)
	commissionCents = totalPrice * commissionRate / 100
	playerIncomeCents = totalPrice - commissionCents

	return
}

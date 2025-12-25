package order

import (
	"time"

	"gamelink/internal/model"
)

// buildOrderForCreation 根据请求和定价结果构建待持久化的订单实体
func (s *OrderService) buildOrderForCreation(userID uint64, req CreateOrderRequest, totalPrice, commissionCents, playerIncomeCents int64) *model.Order {
	// 计算结束时间（保持原有实现方式）
	scheduledEnd := req.ScheduledStart.Add(time.Duration(req.DurationHours * float32(time.Hour)))

	playerID := req.PlayerID
	gameID := req.GameID

	return &model.Order{
		OrderNo:           model.GenerateEscortOrderNo(),
		UserID:            userID,
		ItemID:            1, // TODO: 后续从 service_items 选择对应的服务项
		PlayerID:          &playerID,
		GameID:            &gameID,
		Quantity:          1,
		UnitPriceCents:    totalPrice,
		TotalPriceCents:   totalPrice,
		CommissionCents:   commissionCents,
		PlayerIncomeCents: playerIncomeCents,
		Currency:          model.CurrencyCNY,
		Status:            model.OrderStatusPending,
		Title:             req.Title,
		Description:       req.Description,
		ScheduledStart:    req.ScheduledStart,
		ScheduledEnd:      &scheduledEnd,
	}
}

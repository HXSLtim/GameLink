package order

import (
	"math"
	"time"

	"gamelink/internal/model"
)

// buildOrderForCreation 根据请求和定价结果构建待持久化的订单实体
func (s *OrderService) buildOrderForCreation(userID uint64, req CreateOrderRequest, totalPrice, commissionCents, playerIncomeCents int64) *model.Order {
	// 计算结束时间（保持原有实现方式）
	scheduledEnd := req.ScheduledStart.Add(time.Duration(req.DurationHours * float32(time.Hour)))

	playerID := req.PlayerID
	gameID := req.GameID

	// 使用请求中的 ServiceID，如果未提供则默认为 0（允许 NULL）
	var itemID uint64
	if req.ServiceID != nil {
		itemID = *req.ServiceID
	}

	return &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:           model.GenerateEscortOrderNo(),
		UserID:            userID,
		ItemID:            itemID,
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
		OrderConfig:       "{}",
	}
}

// buildOrderGroupWithSubOrders 创建主订单和拆分的子订单
// 将 N 小时的订单拆分成 N 个子订单，每个子订单 1 小时
func (s *OrderService) buildOrderGroupWithSubOrders(
	userID uint64,
	req CreateOrderRequest,
	hourlyPrice, commissionPerHour, playerIncomePerHour int64,
) (*model.OrderGroup, []*model.Order) {
	// 计算总小时数（向上取整）
	totalHours := int(math.Ceil(float64(req.DurationHours)))
	if totalHours < 1 {
		totalHours = 1
	}

	// 计算总价
	totalPrice := hourlyPrice * int64(totalHours)

	// 计算结束时间
	scheduledEnd := req.ScheduledStart.Add(time.Duration(totalHours) * time.Hour)

	playerID := req.PlayerID
	gameID := req.GameID

	var itemID uint64
	if req.ServiceID != nil {
		itemID = *req.ServiceID
	}

	// 创建主订单
	group := &model.OrderGroup{
		Base: model.Base{
			ExtJSON: "{}",
		},
		GroupNo:         model.GenerateGroupOrderNo(),
		UserID:          userID,
		GameID:          gameID,
		ItemID:          itemID,
		OriginalPlayer:  playerID,
		TotalPriceCents: totalPrice,
		TotalHours:      totalHours,
		CompletedHours:  0,
		Status:          model.OrderGroupStatusPending,
		Title:           req.Title,
		Description:     req.Description,
		ScheduledStart:  req.ScheduledStart,
		ScheduledEnd:    &scheduledEnd,
		Currency:        model.CurrencyCNY,
	}

	// 创建子订单（每小时一个）
	subOrders := make([]*model.Order, 0, totalHours)
	for i := 0; i < totalHours; i++ {
		hourStart := req.ScheduledStart.Add(time.Duration(i) * time.Hour)
		hourEnd := hourStart.Add(time.Hour)

		subOrder := &model.Order{
			Base: model.Base{
				ExtJSON: "{}",
			},
			OrderNo:           model.GenerateEscortOrderNo(),
			UserID:            userID,
			ItemID:            itemID,
			PlayerID:          &playerID,
			GameID:            &gameID,
			Quantity:          1,
			UnitPriceCents:    hourlyPrice,
			TotalPriceCents:   hourlyPrice,
			CommissionCents:   commissionPerHour,
			PlayerIncomeCents: playerIncomePerHour,
			Currency:          model.CurrencyCNY,
			Status:            model.OrderStatusPending,
			Title:             req.Title,
			Description:       req.Description,
			ScheduledStart:    &hourStart,
			ScheduledEnd:      &hourEnd,
			OrderConfig:       "{}",
			// 拆分相关字段
			HourIndex:   i + 1, // 从 1 开始
			IsSubOrder:  true,
			CanTransfer: true,
		}
		subOrders = append(subOrders, subOrder)
	}

	return group, subOrders
}

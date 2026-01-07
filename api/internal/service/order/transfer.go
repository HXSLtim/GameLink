package order

import (
	"context"
	"fmt"

	"gamelink/internal/model"
	"gamelink/pkg/apierr"
)

// TransferSubOrderRequest 转单请求
type TransferSubOrderRequest struct {
	SubOrderID   uint64 `json:"subOrderId" binding:"required"`   // 要转的子订单ID
	NewPlayerID  uint64 `json:"newPlayerId" binding:"required"`  // 新陪玩师ID
	TransferNote string `json:"transferNote"`                    // 转单备注
}

// TransferSubOrderResponse 转单响应
type TransferSubOrderResponse struct {
	Success        bool   `json:"success"`
	NewSubOrderID  uint64 `json:"newSubOrderId"`  // 新子订单ID
	Message        string `json:"message"`
}

// TransferSubOrder 转单 - 将子订单转给另一个陪玩师
// 场景：陪玩 A 打了一半打不了了，将剩余时间转给陪玩 B
func (s *OrderService) TransferSubOrder(ctx context.Context, operatorID uint64, req TransferSubOrderRequest) (*TransferSubOrderResponse, error) {
	// 1. 获取原子订单
	subOrder, err := s.orders.Get(ctx, req.SubOrderID)
	if err != nil {
		return nil, apierr.NotFound("子订单不存在")
	}

	// 2. 验证是否可转单
	if !subOrder.IsSubOrder {
		return nil, apierr.BadRequest("只有子订单可以转单")
	}
	if !subOrder.CanTransfer {
		return nil, apierr.BadRequest("该订单不可转单")
	}
	if subOrder.Status == model.OrderStatusCompleted {
		return nil, apierr.BadRequest("已完成的订单不能转单")
	}
	if subOrder.Status == model.OrderStatusCanceled || subOrder.Status == model.OrderStatusRefunded {
		return nil, apierr.BadRequest("已取消/退款的订单不能转单")
	}

	// 3. 验证新陪玩师
	newPlayer, err := s.players.Get(ctx, req.NewPlayerID)
	if err != nil {
		return nil, apierr.NotFound("新陪玩师不存在")
	}
	if newPlayer.VerificationStatus != "verified" {
		return nil, apierr.BadRequest("新陪玩师未通过认证")
	}

	// 4. 不能转给同一个陪玩师
	if subOrder.PlayerID != nil && *subOrder.PlayerID == req.NewPlayerID {
		return nil, apierr.BadRequest("不能转给同一个陪玩师")
	}

	// 5. 创建新的子订单（复制原订单信息，更换陪玩师）
	newSubOrder := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:           model.GenerateEscortOrderNo(),
		UserID:            subOrder.UserID,
		ItemID:            subOrder.ItemID,
		PlayerID:          &req.NewPlayerID,
		GameID:            subOrder.GameID,
		GroupID:           subOrder.GroupID,
		Quantity:          subOrder.Quantity,
		UnitPriceCents:    subOrder.UnitPriceCents,
		TotalPriceCents:   subOrder.TotalPriceCents,
		CommissionCents:   subOrder.CommissionCents,
		PlayerIncomeCents: subOrder.PlayerIncomeCents,
		Currency:          subOrder.Currency,
		Status:            model.OrderStatusPending, // 新订单待确认
		Title:             subOrder.Title,
		Description:       subOrder.Description,
		ScheduledStart:    subOrder.ScheduledStart,
		ScheduledEnd:      subOrder.ScheduledEnd,
		OrderConfig:       subOrder.OrderConfig,
		// 拆分相关
		HourIndex:    subOrder.HourIndex,
		IsSubOrder:   true,
		CanTransfer:  true,
		TransferFrom: &subOrder.ID,
		TransferNote: req.TransferNote,
	}

	// 6. 更新原订单状态
	subOrder.Status = model.OrderStatusCanceled
	subOrder.CanTransfer = false
	subOrder.TransferTo = &newSubOrder.ID
	subOrder.CancelReason = fmt.Sprintf("转单给陪玩师 %d: %s", req.NewPlayerID, req.TransferNote)

	// 7. 保存新订单
	if err := s.orders.Create(ctx, newSubOrder); err != nil {
		return nil, apierr.InternalError("创建新订单失败").WithDetails(err.Error())
	}

	// 8. 更新原订单的 TransferTo
	subOrder.TransferTo = &newSubOrder.ID
	if err := s.orders.Update(ctx, subOrder); err != nil {
		return nil, apierr.InternalError("更新原订单失败").WithDetails(err.Error())
	}

	// 9. 更新主订单状态
	if subOrder.GroupID != nil && s.orderGroups != nil {
		group, err := s.orderGroups.GetWithSubOrders(ctx, *subOrder.GroupID)
		if err == nil {
			group.UpdateStatusFromSubOrders(group.SubOrders)
			_ = s.orderGroups.Update(ctx, group)
		}
	}

	return &TransferSubOrderResponse{
		Success:       true,
		NewSubOrderID: newSubOrder.ID,
		Message:       "转单成功",
	}, nil
}

// BatchTransferSubOrders 批量转单 - 将多个子订单转给另一个陪玩师
// 场景：陪玩 A 打了 1 小时后无法继续，将剩余 2 小时都转给陪玩 B
type BatchTransferRequest struct {
	SubOrderIDs  []uint64 `json:"subOrderIds" binding:"required,min=1"` // 要转的子订单ID列表
	NewPlayerID  uint64   `json:"newPlayerId" binding:"required"`       // 新陪玩师ID
	TransferNote string   `json:"transferNote"`                         // 转单备注
}

type BatchTransferResponse struct {
	SuccessCount int      `json:"successCount"`
	FailedCount  int      `json:"failedCount"`
	NewOrderIDs  []uint64 `json:"newOrderIds"`
	Errors       []string `json:"errors,omitempty"`
}

// BatchTransferSubOrders 批量转单
func (s *OrderService) BatchTransferSubOrders(ctx context.Context, operatorID uint64, req BatchTransferRequest) (*BatchTransferResponse, error) {
	resp := &BatchTransferResponse{
		NewOrderIDs: make([]uint64, 0),
		Errors:      make([]string, 0),
	}

	for _, subOrderID := range req.SubOrderIDs {
		result, err := s.TransferSubOrder(ctx, operatorID, TransferSubOrderRequest{
			SubOrderID:   subOrderID,
			NewPlayerID:  req.NewPlayerID,
			TransferNote: req.TransferNote,
		})
		if err != nil {
			resp.FailedCount++
			resp.Errors = append(resp.Errors, fmt.Sprintf("订单 %d: %s", subOrderID, err.Error()))
		} else {
			resp.SuccessCount++
			resp.NewOrderIDs = append(resp.NewOrderIDs, result.NewSubOrderID)
		}
	}

	return resp, nil
}

// GetTransferableSubOrders 获取可转单的子订单列表
func (s *OrderService) GetTransferableSubOrders(ctx context.Context, groupID uint64) ([]*model.Order, error) {
	if s.orderGroups == nil {
		return nil, apierr.InternalError("orderGroups repository not configured")
	}

	group, err := s.orderGroups.GetWithSubOrders(ctx, groupID)
	if err != nil {
		return nil, err
	}

	transferable := make([]*model.Order, 0)
	for i := range group.SubOrders {
		sub := &group.SubOrders[i]
		if sub.CanTransfer && sub.Status != model.OrderStatusCompleted &&
			sub.Status != model.OrderStatusCanceled && sub.Status != model.OrderStatusRefunded {
			transferable = append(transferable, sub)
		}
	}

	return transferable, nil
}

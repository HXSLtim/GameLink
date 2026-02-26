package order

import (
	"context"
	"fmt"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/common"
	"gamelink/pkg/apierr"
)

// TransferSubOrderRequest 转单请求
type TransferSubOrderRequest struct {
	SubOrderID       uint64 `json:"subOrderId" binding:"required"`  // 要转的子订单ID
	NewPlayerID      uint64 `json:"newPlayerId" binding:"required"` // 新陪玩师ID
	TransferNote     string `json:"transferNote"`                   // 转单备注
	CompletedMinutes int    `json:"completedMinutes"`               // 原陪玩师已完成的分钟数（0表示未开始）
}

// TransferSubOrderResponse 转单响应
type TransferSubOrderResponse struct {
	Success              bool   `json:"success"`
	NewSubOrderID        uint64 `json:"newSubOrderId"`        // 新子订单ID
	OriginalPlayerIncome int64  `json:"originalPlayerIncome"` // 原陪玩师应得收入（分）
	NewPlayerIncome      int64  `json:"newPlayerIncome"`      // 新陪玩师应得收入（分）
	Message              string `json:"message"`
}

// TransferSubOrder 转单 - 将子订单转给另一个陪玩师
// 场景：陪玩 A 打了一半打不了了，将剩余时间转给陪玩 B
//
// 收入归属规则：
// 1. 如果原陪玩师未开始服务（CompletedMinutes=0），全部收入归新陪玩师
// 2. 如果原陪玩师已开始服务，按已完成时间比例分配收入
// 3. 平台抽成只计算一次，不重复扣除
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

	// 5. 计算收入分配
	// 每个子订单代表1小时（60分钟）
	totalMinutes := 60
	completedMinutes := req.CompletedMinutes
	if completedMinutes < 0 {
		completedMinutes = 0
	}
	if completedMinutes > totalMinutes {
		completedMinutes = totalMinutes
	}
	remainingMinutes := totalMinutes - completedMinutes

	// 计算原陪玩师和新陪玩师的收入分配
	// 注意：抽成已经在原订单中计算过，不需要重复计算
	originalPlayerIncome := int64(0)
	newPlayerIncome := subOrder.PlayerIncomeCents

	if completedMinutes > 0 {
		// 按比例分配陪玩师收入
		originalPlayerIncome = subOrder.PlayerIncomeCents * int64(completedMinutes) / int64(totalMinutes)
		newPlayerIncome = subOrder.PlayerIncomeCents - originalPlayerIncome
	}

	// 6. 创建新的子订单（复制原订单信息，更换陪玩师）
	now := time.Now()
	newSubOrder := &model.Order{
		Base: model.Base{
			ExtJSON: "{}",
		},
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          subOrder.UserID,
		ItemID:          subOrder.ItemID,
		PlayerID:        &req.NewPlayerID,
		GameID:          subOrder.GameID,
		GroupID:         subOrder.GroupID,
		Quantity:        subOrder.Quantity,
		UnitPriceCents:  subOrder.UnitPriceCents,
		TotalPriceCents: subOrder.TotalPriceCents,
		// 关键修复：新订单的抽成为0（抽成已在原订单计算），收入为剩余部分
		CommissionCents:   0, // 抽成不重复计算
		PlayerIncomeCents: newPlayerIncome,
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
		TransferNote: fmt.Sprintf("%s (剩余%d分钟)", req.TransferNote, remainingMinutes),
	}

	// 7. 更新原订单状态和收入
	subOrder.Status = model.OrderStatusCanceled
	subOrder.CanTransfer = false
	subOrder.PlayerIncomeCents = originalPlayerIncome // 更新为实际应得收入
	subOrder.CancelReason = fmt.Sprintf("转单给陪玩师 %d (已完成%d分钟): %s", req.NewPlayerID, completedMinutes, req.TransferNote)
	subOrder.CompletedAt = &now

	if s.tx == nil {
		return nil, apierr.InternalError("transaction manager not configured")
	}

	// 8-10. 使用事务确保新订单创建、原订单更新、主订单状态更新的原子性
	err = s.tx.WithTx(ctx, func(r *common.Repos) error {
		// 8. 保存新订单
		if err := r.Orders.Create(ctx, newSubOrder); err != nil {
			return fmt.Errorf("创建新订单失败: %w", err)
		}

		// 9. 更新原订单的 TransferTo
		subOrder.TransferTo = &newSubOrder.ID
		if err := r.Orders.Update(ctx, subOrder); err != nil {
			return fmt.Errorf("更新原订单失败: %w", err)
		}

		// 10. 更新主订单状态
		if subOrder.GroupID != nil {
			group, err := r.OrderGroups.GetWithSubOrders(ctx, *subOrder.GroupID)
			if err == nil {
				group.UpdateStatusFromSubOrders(group.SubOrders)
				_ = r.OrderGroups.Update(ctx, group)
			}
		}

		return nil
	})
	if err != nil {
		return nil, apierr.InternalError(err.Error())
	}

	return &TransferSubOrderResponse{
		Success:              true,
		NewSubOrderID:        newSubOrder.ID,
		OriginalPlayerIncome: originalPlayerIncome,
		NewPlayerIncome:      newPlayerIncome,
		Message:              fmt.Sprintf("转单成功，原陪玩师收入 %d 分，新陪玩师收入 %d 分", originalPlayerIncome, newPlayerIncome),
	}, nil
}

// BatchTransferSubOrders 批量转单 - 将多个子订单转给另一个陪玩师
// 场景：陪玩 A 打了 1 小时后无法继续，将剩余 2 小时都转给陪玩 B
type BatchTransferRequest struct {
	SubOrderIDs      []uint64 `json:"subOrderIds" binding:"required,min=1"` // 要转的子订单ID列表
	NewPlayerID      uint64   `json:"newPlayerId" binding:"required"`       // 新陪玩师ID
	TransferNote     string   `json:"transferNote"`                         // 转单备注
	CompletedMinutes int      `json:"completedMinutes"`                     // 第一个订单已完成的分钟数（后续订单视为未开始）
}

type BatchTransferResponse struct {
	SuccessCount         int      `json:"successCount"`
	FailedCount          int      `json:"failedCount"`
	NewOrderIDs          []uint64 `json:"newOrderIds"`
	TotalOriginalIncome  int64    `json:"totalOriginalIncome"`  // 原陪玩师总收入
	TotalNewPlayerIncome int64    `json:"totalNewPlayerIncome"` // 新陪玩师总收入
	Errors               []string `json:"errors,omitempty"`
}

// BatchTransferSubOrders 批量转单
// 注意：只有第一个订单使用 CompletedMinutes，后续订单视为未开始（CompletedMinutes=0）
func (s *OrderService) BatchTransferSubOrders(ctx context.Context, operatorID uint64, req BatchTransferRequest) (*BatchTransferResponse, error) {
	resp := &BatchTransferResponse{
		NewOrderIDs: make([]uint64, 0),
		Errors:      make([]string, 0),
	}

	for i, subOrderID := range req.SubOrderIDs {
		// 只有第一个订单使用已完成分钟数，后续订单视为未开始
		completedMinutes := 0
		if i == 0 {
			completedMinutes = req.CompletedMinutes
		}

		result, err := s.TransferSubOrder(ctx, operatorID, TransferSubOrderRequest{
			SubOrderID:       subOrderID,
			NewPlayerID:      req.NewPlayerID,
			TransferNote:     req.TransferNote,
			CompletedMinutes: completedMinutes,
		})
		if err != nil {
			resp.FailedCount++
			resp.Errors = append(resp.Errors, fmt.Sprintf("订单 %d: %s", subOrderID, err.Error()))
		} else {
			resp.SuccessCount++
			resp.NewOrderIDs = append(resp.NewOrderIDs, result.NewSubOrderID)
			resp.TotalOriginalIncome += result.OriginalPlayerIncome
			resp.TotalNewPlayerIncome += result.NewPlayerIncome
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

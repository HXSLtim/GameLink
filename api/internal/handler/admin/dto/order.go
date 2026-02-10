package dto

import (
	"time"

	"gamelink/internal/model"
)

// ==================== Response DTOs ====================

// OrderResponse 订单响应 DTO
type OrderResponse struct {
	ID              uint64            `json:"id"`
	OrderNo         string            `json:"orderNo"`
	UserID          uint64            `json:"userId"`
	PlayerID        *uint64           `json:"playerId,omitempty"`
	ItemID          uint64            `json:"itemId"`
	GameID          *uint64           `json:"gameId,omitempty"`
	Status          model.OrderStatus `json:"status"`
	Title           string            `json:"title,omitempty"`
	Description     string            `json:"description,omitempty"`

	// 价格
	Quantity          int    `json:"quantity"`
	UnitPriceCents    int64  `json:"unitPriceCents"`
	TotalPriceCents   int64  `json:"totalPriceCents"`
	CommissionCents   int64  `json:"commissionCents"`
	PlayerIncomeCents int64  `json:"playerIncomeCents"`

	// 时间
	ScheduledStart *time.Time `json:"scheduledStart,omitempty"`
	ScheduledEnd   *time.Time `json:"scheduledEnd,omitempty"`
	StartedAt      *time.Time `json:"startedAt,omitempty"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`

	// 取消/退款
	CancelReason      string     `json:"cancelReason,omitempty"`
	RefundAmountCents int64      `json:"refundAmountCents,omitempty"`
	RefundedAt        *time.Time `json:"refundedAt,omitempty"`

	// 争议
	HasDispute bool `json:"hasDispute"`

	// 简要关联信息
	UserBrief   *UserBrief   `json:"user,omitempty"`
	PlayerBrief *PlayerBrief `json:"player,omitempty"`
}

// UserBrief 用户简要信息（嵌入其他 DTO 用）
type UserBrief struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Items      []OrderResponse `json:"items"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalPages int             `json:"totalPages"`
}

// ==================== 转换函数 ====================

// ToOrderResponse 将 model.Order 转换为 OrderResponse
func ToOrderResponse(order *model.Order) *OrderResponse {
	if order == nil {
		return nil
	}

	resp := &OrderResponse{
		ID:                order.ID,
		OrderNo:           order.OrderNo,
		UserID:            order.UserID,
		PlayerID:          order.PlayerID,
		ItemID:            order.ItemID,
		GameID:            order.GameID,
		Status:            order.Status,
		Title:             order.Title,
		Description:       order.Description,
		Quantity:          order.Quantity,
		UnitPriceCents:    order.UnitPriceCents,
		TotalPriceCents:   order.TotalPriceCents,
		CommissionCents:   order.CommissionCents,
		PlayerIncomeCents: order.PlayerIncomeCents,
		ScheduledStart:    order.ScheduledStart,
		ScheduledEnd:      order.ScheduledEnd,
		StartedAt:         order.StartedAt,
		CompletedAt:       order.CompletedAt,
		CreatedAt:         order.CreatedAt,
		CancelReason:      order.CancelReason,
		RefundAmountCents: order.RefundAmountCents,
		RefundedAt:        order.RefundedAt,
		HasDispute:        order.HasDispute,
	}

	// 嵌入用户简要信息
	if order.User.ID > 0 {
		resp.UserBrief = &UserBrief{
			ID:        order.User.ID,
			Name:      order.User.Name,
			AvatarURL: order.User.AvatarURL,
		}
	}

	// 嵌入陪玩师简要信息
	if order.Player != nil {
		resp.PlayerBrief = &PlayerBrief{
			ID:       order.Player.ID,
			Nickname: order.Player.Nickname,
			UserID:   order.Player.UserID,
		}
	}

	return resp
}

// ToOrderResponseList 批量转换
func ToOrderResponseList(orders []model.Order) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))
	for i := range orders {
		if resp := ToOrderResponse(&orders[i]); resp != nil {
			responses = append(responses, *resp)
		}
	}
	return responses
}

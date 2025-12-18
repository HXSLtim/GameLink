package statistics

import (
	"context"
	"log"

	"gamelink/internal/model"
)

// EventHooks 统计事件钩子
type EventHooks struct {
	svc *Service
}

// NewEventHooks 创建事件钩子
func NewEventHooks(svc *Service) *EventHooks {
	return &EventHooks{svc: svc}
}

// OnOrderCompleted 订单完成时触发
func (h *EventHooks) OnOrderCompleted(ctx context.Context, order *model.Order) {
	go func() {
		bgCtx := context.Background()

		// 更新用户统计
		if err := h.svc.UpdateUserStatistics(bgCtx, order.UserID); err != nil {
			log.Printf("[Statistics] update user %d stats failed: %v", order.UserID, err)
		}

		// 更新陪玩师统计
		if order.PlayerID != nil {
			if err := h.svc.UpdatePlayerStatistics(bgCtx, *order.PlayerID); err != nil {
				log.Printf("[Statistics] update player %d stats failed: %v", *order.PlayerID, err)
			}
		}

		// 更新服务项目统计
		if err := h.svc.UpdateServiceItemStatistics(bgCtx, order.ItemID); err != nil {
			log.Printf("[Statistics] update service item %d stats failed: %v", order.ItemID, err)
		}

		// 更新游戏统计
		if order.GameID != nil {
			if err := h.svc.UpdateGameStatistics(bgCtx, *order.GameID); err != nil {
				log.Printf("[Statistics] update game %d stats failed: %v", *order.GameID, err)
			}
		}
	}()
}

// OnOrderCanceled 订单取消时触发
func (h *EventHooks) OnOrderCanceled(ctx context.Context, order *model.Order) {
	go func() {
		bgCtx := context.Background()

		if err := h.svc.UpdateUserStatistics(bgCtx, order.UserID); err != nil {
			log.Printf("[Statistics] update user %d stats failed: %v", order.UserID, err)
		}

		if order.PlayerID != nil {
			if err := h.svc.UpdatePlayerStatistics(bgCtx, *order.PlayerID); err != nil {
				log.Printf("[Statistics] update player %d stats failed: %v", *order.PlayerID, err)
			}
		}
	}()
}

// OnOrderRefunded 订单退款时触发
func (h *EventHooks) OnOrderRefunded(ctx context.Context, order *model.Order) {
	go func() {
		bgCtx := context.Background()

		if err := h.svc.UpdateUserStatistics(bgCtx, order.UserID); err != nil {
			log.Printf("[Statistics] update user %d stats failed: %v", order.UserID, err)
		}

		if order.PlayerID != nil {
			if err := h.svc.UpdatePlayerStatistics(bgCtx, *order.PlayerID); err != nil {
				log.Printf("[Statistics] update player %d stats failed: %v", *order.PlayerID, err)
			}
		}

		if err := h.svc.UpdateServiceItemStatistics(bgCtx, order.ItemID); err != nil {
			log.Printf("[Statistics] update service item %d stats failed: %v", order.ItemID, err)
		}
	}()
}

// OnDisputeResolved 争议解决时触发
func (h *EventHooks) OnDisputeResolved(ctx context.Context, dispute *model.OrderDispute, order *model.Order) {
	go func() {
		bgCtx := context.Background()

		if err := h.svc.UpdateUserStatistics(bgCtx, dispute.UserID); err != nil {
			log.Printf("[Statistics] update user %d stats failed: %v", dispute.UserID, err)
		}

		if order != nil && order.PlayerID != nil {
			if err := h.svc.UpdatePlayerStatistics(bgCtx, *order.PlayerID); err != nil {
				log.Printf("[Statistics] update player %d stats failed: %v", *order.PlayerID, err)
			}
		}
	}()
}

// OnReviewCreated 评价创建时触发
func (h *EventHooks) OnReviewCreated(ctx context.Context, review *model.Review) {
	go func() {
		bgCtx := context.Background()

		if err := h.svc.UpdateUserStatistics(bgCtx, review.UserID); err != nil {
			log.Printf("[Statistics] update user %d stats failed: %v", review.UserID, err)
		}

		if err := h.svc.UpdatePlayerStatistics(bgCtx, review.PlayerID); err != nil {
			log.Printf("[Statistics] update player %d stats failed: %v", review.PlayerID, err)
		}
	}()
}

// OnWithdrawCompleted 提现完成时触发
func (h *EventHooks) OnWithdrawCompleted(ctx context.Context, withdraw *model.Withdraw) {
	go func() {
		bgCtx := context.Background()

		if err := h.svc.UpdatePlayerStatistics(bgCtx, withdraw.PlayerID); err != nil {
			log.Printf("[Statistics] update player %d stats failed: %v", withdraw.PlayerID, err)
		}
	}()
}

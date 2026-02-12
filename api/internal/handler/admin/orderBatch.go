package admin

import (
	"context"
	"time"

	_ "gamelink/internal/model" // Imported for Swagger annotations

	"github.com/gin-gonic/gin"

	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
	"gamelink/pkg/logging"
)

// BatchCancelOrders 批量取消订单
// @Summary      批量取消订单
// @Description  取消多个待处理或已确认的订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchCancelOrdersRequest  true  "订单ID列表和取消原因"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/cancel [post]
func (h *OrderHandler) BatchCancelOrders(c *gin.Context) {
	var req BatchCancelOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}

	result, err := h.svc.BatchCancelOrders(contextWithActor(c), req.OrderIDs, req.Reason, req.Note)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch cancel orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchConfirmOrders 批量确认订单
// @Summary      批量确认订单
// @Description  确认多个待处理的订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchConfirmOrdersRequest  true  "订单ID列表"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/confirm [post]
func (h *OrderHandler) BatchConfirmOrders(c *gin.Context) {
	var req BatchConfirmOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}

	result, err := h.svc.BatchConfirmOrders(contextWithActor(c), req.OrderIDs, req.Note)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch confirm orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchCompleteOrders 批量完成订单
// @Summary      批量完成订单
// @Description  将多个进行中的订单标记为已完成
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchCompleteOrdersRequest  true  "订单ID列表"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/complete [post]
func (h *OrderHandler) BatchCompleteOrders(c *gin.Context) {
	var req BatchCompleteOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}

	result, err := h.svc.BatchCompleteOrders(contextWithActor(c), req.OrderIDs, req.Note)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch complete orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchRefundOrders 批量退款订单
// @Summary      批量退款订单
// @Description  对多个订单执行退款操作
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchRefundOrdersRequest  true  "订单ID列表和退款信息"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/refund [post]
func (h *OrderHandler) BatchRefundOrders(c *gin.Context) {
	var req BatchRefundOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}
	if req.Reason == "" {
		respondAPIError(c, apierr.BadRequest("reason is required"))
		return
	}

	result, err := h.svc.BatchRefundOrders(contextWithActor(c), req.OrderIDs, adminservice.BatchRefundInput{
		Reason:      req.Reason,
		AmountCents: req.AmountCents,
		Note:        req.Note,
		RefundedAt:  parseTimePtr(req.RefundedAt),
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch refund orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchDeleteOrders 批量删除订单
// @Summary      批量删除订单
// @Description  软删除多个订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchDeleteOrdersRequest  true  "订单ID列表"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/delete [post]
func (h *OrderHandler) BatchDeleteOrders(c *gin.Context) {
	var req BatchDeleteOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}

	result, err := h.svc.BatchDeleteOrders(contextWithActor(c), req.OrderIDs)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch delete orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchUpdateOrderStatus 批量更新订单状态
// @Summary      批量更新订单状态
// @Description  更新多个订单的状态
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchUpdateOrderStatusRequest  true  "订单ID列表和新状态"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/status [put]
func (h *OrderHandler) BatchUpdateOrderStatus(c *gin.Context) {
	var req BatchUpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}
	if req.Status == "" {
		respondAPIError(c, apierr.BadRequest("status is required"))
		return
	}

	result, err := h.svc.BatchUpdateOrderStatus(contextWithActor(c), req.OrderIDs, adminservice.BatchUpdateStatusInput{
		Status:       normalizeOrderStatus(req.Status),
		Note:         req.Note,
		StartedAt:    parseTimePtr(req.StartedAt),
		CompletedAt:  parseTimePtr(req.CompletedAt),
		CancelReason: req.CancelReason,
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch update order status failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// BatchAssignOrders 批量指派订单
// @Summary      批量指派订单
// @Description  将多个订单指派给指定的陪玩师
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  BatchAssignOrdersRequest  true  "订单ID列表和陪玩师ID"
// @Success      200  {object}  BatchOperationResponse
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders/batch/assign [post]
func (h *OrderHandler) BatchAssignOrders(c *gin.Context) {
	var req BatchAssignOrdersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	if len(req.OrderIDs) == 0 {
		respondAPIError(c, apierr.BadRequest("order_ids is required"))
		return
	}
	if len(req.OrderIDs) > 100 {
		respondAPIError(c, apierr.BadRequest("maximum 100 orders per batch"))
		return
	}
	if req.PlayerID == 0 {
		respondAPIError(c, apierr.BadRequest("player_id is required"))
		return
	}

	result, err := h.svc.BatchAssignOrders(contextWithActor(c), req.OrderIDs, req.PlayerID)
	if err != nil {
		respondAPIError(c, apierr.InternalError("batch assign orders failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, result)
}

// Batch Operation Request/Response DTOs

// BatchCancelOrdersRequest 批量取消订单请求
type BatchCancelOrdersRequest struct {
	OrderIDs []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	Reason   string   `json:"reason" binding:"required"`
	Note     string   `json:"note"`
}

// BatchConfirmOrdersRequest 批量确认订单请求
type BatchConfirmOrdersRequest struct {
	OrderIDs []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	Note     string   `json:"note"`
}

// BatchCompleteOrdersRequest 批量完成订单请求
type BatchCompleteOrdersRequest struct {
	OrderIDs []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	Note     string   `json:"note"`
}

// BatchRefundOrdersRequest 批量退款订单请求
type BatchRefundOrdersRequest struct {
	OrderIDs    []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	Reason      string   `json:"reason" binding:"required"`
	AmountCents *int64   `json:"amount_cents,omitempty"` // nil 表示全额退款
	Note        string   `json:"note"`
	RefundedAt  *string  `json:"refunded_at,omitempty"`
}

// BatchDeleteOrdersRequest 批量删除订单请求
type BatchDeleteOrdersRequest struct {
	OrderIDs []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
}

// BatchUpdateOrderStatusRequest 批量更新订单状态请求
type BatchUpdateOrderStatusRequest struct {
	OrderIDs     []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	Status       string   `json:"status" binding:"required"`
	Note         string   `json:"note"`
	StartedAt    *string  `json:"started_at,omitempty"`
	CompletedAt  *string  `json:"completed_at,omitempty"`
	CancelReason string   `json:"cancel_reason,omitempty"`
}

// BatchAssignOrdersRequest 批量指派订单请求
type BatchAssignOrdersRequest struct {
	OrderIDs []uint64 `json:"order_ids" binding:"required,min=1,max=100"`
	PlayerID uint64   `json:"player_id" binding:"required"`
}

func contextWithActor(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if actorUserID := getUserIDFromContext(c); actorUserID != 0 {
		return logging.WithActorUserID(ctx, actorUserID)
	}
	return ctx
}

// Helper function to parse time pointer
func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}

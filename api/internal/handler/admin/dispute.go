package admin

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	orderservice "gamelink/internal/service/order"
	apierr "gamelink/pkg/apierr"
)

// OrderDispute 订单纠纷模型（类型别名）
type OrderDispute = model.OrderDispute

// DisputeHandler handles order dispute related endpoints
type DisputeHandler struct {
	svc *orderservice.DisputeService
}

// NewDisputeHandler creates a new dispute handler
func NewDisputeHandler(svc *orderservice.DisputeService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

// GetDisputeDetail retrieves dispute details
// @Summary      获取纠纷详情
// @Description  根据 ID 获取单个订单纠纷的详细信
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "纠纷ID"
// @Success      200  {object}  model.APIResponse[OrderDispute]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/disputes/{id} [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
	disputeID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	dispute, err := h.svc.GetDisputeDetail(c.Request.Context(), disputeID)
	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeNotFound) {
			respondError(c, apierr.NotFound(apierr.ErrDisputeNotFound))
			return
		}
		respondError(c, err)
		return
	}

	respondSuccess(c, dispute)
}

// ListPendingDisputes lists disputes pending assignment
// @Summary      列出待处理纠纷
// @Description  获取状态为待处理的订单纠纷列表，支持分页
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query  int  false  "页码"  default(1)
// @Param        pageSize  query  int  false  "每页数量"    default(20)
// @Success      200  {object}  model.APIResponse[[]OrderDispute]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/disputes/pending [get]
func (h *DisputeHandler) ListPendingDisputes(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	disputes, total, err := h.svc.ListPendingDisputes(c.Request.Context(), page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}

	type ListResponse struct {
		Disputes []model.OrderDispute `json:"disputes"`
		Total    int64                `json:"total"`
		Page     int                  `json:"page"`
		PageSize int                  `json:"pageSize"`
	}
	respondSuccess(c, ListResponse{
		Disputes: disputes,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ListDisputes lists all disputes with optional status filter
// @Summary      列出纠纷列表
// @Description  获取订单纠纷列表，支持状态筛选和分页
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "页码"  default(1)
// @Param        pageSize  query  int     false  "每页数量"    default(20)
// @Param        status    query  string  false  "状态筛选"
// @Param        orderNo   query  string  false  "订单号"
// @Success      200  {object}  model.APIResponse[[]OrderDispute]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/disputes [get]
func (h *DisputeHandler) ListDisputes(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("pageSize"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	status := c.Query("status")
	orderNo := c.Query("orderNo")

	disputes, total, err := h.svc.ListDisputes(c.Request.Context(), orderservice.ListDisputesRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		OrderNo:  orderNo,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	type ListResponse struct {
		Disputes []model.OrderDispute `json:"disputes"`
		Total    int64                `json:"total"`
		Page     int                  `json:"page"`
		PageSize int                  `json:"pageSize"`
	}
	respondSuccess(c, ListResponse{
		Disputes: disputes,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetDisputeStats returns dispute statistics
// @Summary      获取纠纷统计
// @Description  获取各状态纠纷数量统计
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.APIResponse[map[string]int64]
// @Failure      500  {object}  model.ErrorResponse
// @Router       /admin/disputes/stats [get]
func (h *DisputeHandler) GetDisputeStats(c *gin.Context) {
	stats, err := h.svc.GetDisputeStats(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	respondSuccess(c, stats)
}

// AssignDisputePayload represents the request to assign a dispute
type AssignDisputePayload struct {
	AssignedServiceID uint64  `json:"assignedServiceId" binding:"required"`
	OriginalServiceID *uint64 `json:"originalServiceId"`
}

// AssignDispute assigns a dispute to a customer service representative
// @Summary      分配纠纷
// @Description  将一个订单纠纷分配给指定的客服人员处
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "纠纷ID"
// @Param        request  body  AssignDisputePayload   true  "分配信息"
// @Success      200  {object}  model.APIResponse[string]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/disputes/{id}/assign [post]
func (h *DisputeHandler) AssignDispute(c *gin.Context) {
	disputeID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload AssignDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondBadRequest(c, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context (set by auth middleware)
	actorUserID, exists := c.Get("userID")
	if !exists {
		respondError(c, apierr.Unauthorized(apierr.ErrUserIDNotInContext))
		return
	}

	err := h.svc.AssignDispute(c.Request.Context(), orderservice.AssignDisputeRequest{
		DisputeID:         disputeID,
		AssignedServiceID: payload.AssignedServiceID,
		OriginalServiceID: payload.OriginalServiceID,
		ActorUserID:       actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeValidation) {
			respondBadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orderservice.ErrDisputeInvalidStatus) {
			respondError(c, apierr.Conflict(err.Error()))
			return
		}
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]string{
		"message": "Dispute assigned successfully",
	})
}

// RollbackAssignmentPayload represents the request to rollback an assignment
type RollbackAssignmentPayload struct {
	RollbackReason string `json:"rollbackReason" binding:"required"`
}

// RollbackAssignment rolls back a dispute assignment
// @Summary      回滚分配
// @Description  撤销一个订单纠纷的分配，使其回到待处理状
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "纠纷ID"
// @Param        request  body  RollbackAssignmentPayload   true  "回滚信息"
// @Success      200  {object}  model.APIResponse[string]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/disputes/{id}/rollback [post]
func (h *DisputeHandler) RollbackAssignment(c *gin.Context) {
	disputeID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload RollbackAssignmentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondBadRequest(c, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	actorUserID, exists := c.Get("userID")
	if !exists {
		respondError(c, apierr.Unauthorized(apierr.ErrUserIDNotInContext))
		return
	}

	err := h.svc.RollbackDisputeAssignment(c.Request.Context(), orderservice.RollbackDisputeRequest{
		DisputeID:      disputeID,
		RollbackReason: payload.RollbackReason,
		ActorUserID:    actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeValidation) {
			respondBadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orderservice.ErrDisputeInvalidStatus) {
			respondError(c, apierr.Conflict(err.Error()))
			return
		}
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]string{
		"message": "Assignment rolled back successfully",
	})
}

// ResolveDisputePayload represents the request to resolve a dispute
type ResolveDisputePayload struct {
	Resolution    string `json:"resolution" binding:"required,oneof=refund partial reassign reject"`
	ResolveRemark string `json:"resolveRemark" binding:"required"`
}

// ResolveDispute resolves a dispute with a decision
// @Summary      解决纠纷
// @Description  对一个订单纠纷做出最终处理决定，例如退款、重新分配等
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "纠纷ID"
// @Param        request  body  ResolveDisputePayload   true  "处理结果信息"
// @Success      200  {object}  model.APIResponse[string]
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/disputes/{id}/resolve [post]
func (h *DisputeHandler) ResolveDispute(c *gin.Context) {
	disputeID, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload ResolveDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondBadRequest(c, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	actorUserID, exists := c.Get("userID")
	if !exists {
		respondError(c, apierr.Unauthorized(apierr.ErrUserIDNotInContext))
		return
	}

	resolution := model.DisputeResolution(payload.Resolution)

	err := h.svc.ResolveDispute(c.Request.Context(), orderservice.ResolveDisputeRequest{
		DisputeID:     disputeID,
		Resolution:    resolution,
		ResolveRemark: payload.ResolveRemark,
		ActorUserID:   actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeValidation) {
			respondBadRequest(c, err.Error())
			return
		}
		if errors.Is(err, orderservice.ErrDisputeInvalidStatus) {
			respondError(c, apierr.Conflict(err.Error()))
			return
		}
		respondError(c, err)
		return
	}

	respondSuccess(c, map[string]string{
		"message": fmt.Sprintf("Dispute resolved with %s decision", resolution),
	})
}

package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/service/assignment"
)

// DisputeHandler handles order dispute related endpoints
type DisputeHandler struct {
	svc *assignment.AssignmentService
}

// NewDisputeHandler creates a new dispute handler
func NewDisputeHandler(svc *assignment.AssignmentService) *DisputeHandler {
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
// @Success      200  {object}  model.APIResponse[model.OrderDispute]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/disputes/{id} [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
	disputeID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	dispute, err := h.svc.GetDisputeDetail(c.Request.Context(), disputeID)
	if err != nil {
		if errors.Is(err, assignment.ErrNotFound) {
			writeJSONError(c, http.StatusNotFound, apierr.ErrDisputeNotFound)
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[model.OrderDispute]{
		Success: true,
		Code:    http.StatusOK,
		Data:    dispute,
	})
}

// ListPendingDisputes lists disputes pending assignment
// @Summary      列出待处理纠
// @Description  获取状态为待处理的订单纠纷列表，支持分
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query  int  false  "页码"  default(1)
// @Param        pageSize  query  int  false  "每页数量"    default(20)
// @Success      200  {object}  model.APIResponse[[]model.OrderDispute]
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
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	type ListResponse struct {
		Disputes []model.OrderDispute `json:"disputes"`
		Total    int64                `json:"total"`
		Page     int                  `json:"page"`
		PageSize int                  `json:"pageSize"`
	}
	writeJSON(c, http.StatusOK, model.APIResponse[ListResponse]{
		Success: true,
		Code:    http.StatusOK,
		Data: ListResponse{
			Disputes: disputes,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// AssignDisputePayload represents the request to assign a dispute
type AssignDisputePayload struct {
	AssignedToUserID uint64 `json:"assignedToUserId" binding:"required"`
	Source           string `json:"source" binding:"required,oneof=system manual"`
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
	disputeID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload AssignDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context (set by auth middleware)
	actorUserID, exists := c.Get("userID")
	if !exists {
		writeJSONError(c, http.StatusUnauthorized, apierr.ErrUserIDNotInContext)
		return
	}

	source := model.AssignmentSource(payload.Source)
	if source != model.AssignmentSourceSystem && source != model.AssignmentSourceManual {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidAssignmentSource)
		return
	}

	err = h.svc.AssignDispute(c.Request.Context(), assignment.AssignDisputeRequest{
		DisputeID:        disputeID,
		AssignedToUserID: payload.AssignedToUserID,
		Source:           source,
		ActorUserID:      actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, assignment.ErrValidation) {
			writeJSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, assignment.ErrInvalidStatus) {
			writeJSONError(c, http.StatusConflict, err.Error())
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[map[string]string]{
		Success: true,
		Code:    http.StatusOK,
		Data: map[string]string{
			"message": "Dispute assigned successfully",
		},
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
	disputeID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload RollbackAssignmentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	actorUserID, exists := c.Get("userID")
	if !exists {
		writeJSONError(c, http.StatusUnauthorized, apierr.ErrUserIDNotInContext)
		return
	}

	err = h.svc.RollbackAssignment(c.Request.Context(), assignment.RollbackAssignmentRequest{
		DisputeID:      disputeID,
		RollbackReason: payload.RollbackReason,
		ActorUserID:    actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, assignment.ErrValidation) {
			writeJSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, assignment.ErrInvalidStatus) {
			writeJSONError(c, http.StatusConflict, err.Error())
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[map[string]string]{
		Success: true,
		Code:    http.StatusOK,
		Data: map[string]string{
			"message": "Assignment rolled back successfully",
		},
	})
}

// ResolveDisputePayload represents the request to resolve a dispute
type ResolveDisputePayload struct {
	Resolution       string `json:"resolution" binding:"required,oneof=refund partial reassign reject"`
	ResolutionAmount int64  `json:"resolutionAmount"`
	ResolutionNotes  string `json:"resolutionNotes" binding:"required"`
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
	disputeID, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload ResolveDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	// Get actor user ID from context
	actorUserID, exists := c.Get("userID")
	if !exists {
		writeJSONError(c, http.StatusUnauthorized, apierr.ErrUserIDNotInContext)
		return
	}

	resolution := model.DisputeResolution(payload.Resolution)

	err = h.svc.ResolveDispute(c.Request.Context(), assignment.ResolveDisputeRequest{
		DisputeID:        disputeID,
		Resolution:       resolution,
		ResolutionAmount: payload.ResolutionAmount,
		ResolutionNotes:  payload.ResolutionNotes,
		ActorUserID:      actorUserID.(uint64),
	})

	if err != nil {
		if errors.Is(err, assignment.ErrValidation) {
			writeJSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, assignment.ErrInvalidStatus) {
			writeJSONError(c, http.StatusConflict, err.Error())
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[map[string]string]{
		Success: true,
		Code:    http.StatusOK,
		Data: map[string]string{
			"message": fmt.Sprintf("Dispute resolved with %s decision", resolution),
		},
	})
}

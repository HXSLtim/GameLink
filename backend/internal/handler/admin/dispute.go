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
// @Summary      Get Dispute Detail
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "Dispute ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/disputes [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
	disputeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid dispute ID")
		return
	}

	dispute, err := h.svc.GetDisputeDetail(c.Request.Context(), disputeID)
	if err != nil {
		if errors.Is(err, assignment.ErrNotFound) {
			writeJSONError(c, http.StatusNotFound, "Dispute not found")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.OrderDispute]{
		Success: true,
		Code:    http.StatusOK,
		Data:    dispute,
	})
}

// ListPendingDisputes lists disputes pending assignment
// @Summary      List Pending Disputes
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page      query  int  false  "Page number"  default(1)
// @Param        pageSize  query  int  false  "Page size"    default(20)
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/orders/pending-assign [get]
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
		Total    int64               `json:"total"`
		Page     int                 `json:"page"`
		PageSize int                 `json:"pageSize"`
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
// @Summary      Assign Dispute
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  uint64                 true  "Dispute ID"
// @Param        request  body  AssignDisputePayload   true  "Assignment info"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/assign [post]
func (h *DisputeHandler) AssignDispute(c *gin.Context) {
	disputeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid dispute ID")
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
		writeJSONError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	source := model.AssignmentSource(payload.Source)
	if source != model.AssignmentSourceSystem && source != model.AssignmentSourceManual {
		writeJSONError(c, http.StatusBadRequest, "Invalid assignment source")
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
// @Summary      Rollback Assignment
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  uint64                      true  "Dispute ID"
// @Param        request  body  RollbackAssignmentPayload   true  "Rollback info"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/assign/cancel [post]
func (h *DisputeHandler) RollbackAssignment(c *gin.Context) {
	disputeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid dispute ID")
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
		writeJSONError(c, http.StatusUnauthorized, "User ID not found in context")
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
// @Summary      Resolve Dispute
// @Tags         Admin/Disputes
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  uint64                  true  "Dispute ID"
// @Param        request  body  ResolveDisputePayload   true  "Resolution info"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/mediate [post]
func (h *DisputeHandler) ResolveDispute(c *gin.Context) {
	disputeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "Invalid dispute ID")
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
		writeJSONError(c, http.StatusUnauthorized, "User ID not found in context")
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

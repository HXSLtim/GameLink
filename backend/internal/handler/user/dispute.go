package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/service/assignment"
)

// DisputeHandler handles order dispute related endpoints for users
type DisputeHandler struct {
	svc *assignment.AssignmentService
}

// NewDisputeHandler creates a new dispute handler
func NewDisputeHandler(svc *assignment.AssignmentService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

// InitiateDisputePayload represents the request to initiate a dispute
type InitiateDisputePayload struct {
	OrderID      uint64   `json:"orderId" binding:"required"`
	Reason       string   `json:"reason" binding:"required,max=255"`
	Description  string   `json:"description" binding:"max=2000"`
	EvidenceURLs []string `json:"evidenceUrls" binding:"max=9"`
}

// InitiateDispute creates a new dispute for an order
// @Summary      Initiate Dispute
// @Tags         User/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  InitiateDisputePayload  true  "Dispute info"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /user/orders/{id}/dispute [post]
func (h *DisputeHandler) InitiateDispute(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		writeJSONError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	var payload InitiateDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	// Validate evidence URLs count
	if len(payload.EvidenceURLs) > 9 {
		writeJSONError(c, http.StatusBadRequest, "Maximum 9 evidence URLs allowed")
		return
	}

	// Validate evidence URLs are not empty
	for _, url := range payload.EvidenceURLs {
		if url == "" {
			writeJSONError(c, http.StatusBadRequest, "Evidence URLs cannot be empty")
			return
		}
	}

	resp, err := h.svc.InitiateDispute(c.Request.Context(), assignment.InitiateDisputeRequest{
		OrderID:      payload.OrderID,
		UserID:       userID.(uint64),
		Reason:       payload.Reason,
		Description:  payload.Description,
		EvidenceURLs: payload.EvidenceURLs,
	})

	if err != nil {
		if errors.Is(err, assignment.ErrValidation) {
			writeJSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, assignment.ErrUnauthorized) {
			writeJSONError(c, http.StatusForbidden, "You can only initiate disputes for your own orders")
			return
		}
		if errors.Is(err, assignment.ErrCannotInitiateDispute) {
			writeJSONError(c, http.StatusConflict, "Cannot initiate dispute for this order")
			return
		}
		if errors.Is(err, assignment.ErrDisputeExists) {
			writeJSONError(c, http.StatusConflict, "A dispute already exists for this order")
			return
		}
		if errors.Is(err, assignment.ErrOrderNotFound) {
			writeJSONError(c, http.StatusNotFound, "Order not found")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	type InitiateDisputeResponse struct {
		DisputeID   uint64 `json:"disputeId"`
		TraceID     string `json:"traceId"`
		SLADeadline string `json:"slaDeadline"`
	}

	writeJSON(c, http.StatusCreated, model.APIResponse[InitiateDisputeResponse]{
		Success: true,
		Code:    http.StatusCreated,
		Data: InitiateDisputeResponse{
			DisputeID:   resp.DisputeID,
			TraceID:     resp.TraceID,
			SLADeadline: resp.SLADeadline.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}

// GetDisputeDetail retrieves dispute details for a user
// @Summary      Get Dispute Detail
// @Tags         User/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "Dispute ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /user/orders/{id}/disputes [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		writeJSONError(c, http.StatusUnauthorized, "User ID not found in context")
		return
	}

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

	// Verify user owns this dispute
	if dispute.UserID != userID.(uint64) {
		writeJSONError(c, http.StatusForbidden, "You can only view your own disputes")
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.OrderDispute]{
		Success: true,
		Code:    http.StatusOK,
		Data:    dispute,
	})
}

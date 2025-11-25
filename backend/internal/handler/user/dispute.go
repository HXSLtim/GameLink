package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gamelink/internal/apierr"
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

// InitiateDisputeResponse represents the response for initiating a dispute
type InitiateDisputeResponse = assignment.InitiateDisputeResponse

// InitiateDispute creates a new dispute for an order
// @Summary      Initiate Dispute
// @Description  发起订单纠纷
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  InitiateDisputePayload  true  "纠纷信息"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  apierr.APIError
// @Failure      401      {object}  apierr.APIError
// @Failure      403      {object}  apierr.APIError
// @Failure      404      {object}  apierr.APIError
// @Failure      409      {object}  apierr.APIError
// @Failure      500      {object}  apierr.APIError
// @Router       /user/orders/{id}/dispute [post]
func (h *DisputeHandler) InitiateDispute(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		respondAPIError(c, apierr.Unauthorized("用户ID不在上下文中"))
		return
	}

	var payload InitiateDisputePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}

	// Validate evidence URLs count
	if len(payload.EvidenceURLs) > 9 {
		respondAPIError(c, apierr.BadRequest("最多允许上传9个证据链接"))
		return
	}

	// Validate evidence URLs are not empty
	for _, url := range payload.EvidenceURLs {
		if url == "" {
			respondAPIError(c, apierr.BadRequest("证据链接不能为空"))
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
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		if errors.Is(err, assignment.ErrCannotInitiateDispute) {
			respondAPIError(c, apierr.Conflict("当前订单状态不能发起纠纷"))
			return
		}
		if errors.Is(err, assignment.ErrDisputeExists) {
			respondAPIError(c, apierr.Conflict("该订单已存在纠纷"))
			return
		}
		if errors.Is(err, assignment.ErrOrderNotFound) {
			respondAPIError(c, apierr.NotFound("订单不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("发起纠纷失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "纠纷发起成功", assignment.InitiateDisputeResponse{
		DisputeID:   resp.DisputeID,
		TraceID:     resp.TraceID,
		SLADeadline: resp.SLADeadline,
	})
}

// GetDisputeDetail retrieves dispute details for a user
// @Summary      Get Dispute Detail
// @Description  获取订单纠纷详情
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  uint64  true  "纠纷ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  apierr.APIError
// @Failure      401  {object}  apierr.APIError
// @Failure      403  {object}  apierr.APIError
// @Failure      404  {object}  apierr.APIError
// @Failure      500  {object}  apierr.APIError
// @Router       /user/orders/{id}/disputes [get]
func (h *DisputeHandler) GetDisputeDetail(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		respondAPIError(c, apierr.Unauthorized("用户ID不在上下文中"))
		return
	}

	disputeID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	dispute, err := h.svc.GetDisputeDetail(c.Request.Context(), disputeID)
	if err != nil {
		if errors.Is(err, assignment.ErrNotFound) {
			respondAPIError(c, apierr.NotFound("纠纷不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("获取纠纷详情失败").WithDetails(err.Error()))
		return
	}

	// Verify user owns this dispute
	if dispute.UserID != userID.(uint64) {
		respondAPIError(c, apierr.Forbidden("您只能查看自己的纠纷"))
		return
	}

	respondSuccess(c, "OK", dispute)
}

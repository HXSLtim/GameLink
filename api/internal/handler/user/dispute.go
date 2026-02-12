package user

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	orderservice "gamelink/internal/service/order"
	"gamelink/pkg/apierr"
)

// DisputeHandler handles order dispute related endpoints for users
type DisputeHandler struct {
	svc *orderservice.DisputeService
}

// NewDisputeHandler creates a new dispute handler
func NewDisputeHandler(svc *orderservice.DisputeService) *DisputeHandler {
	return &DisputeHandler{svc: svc}
}

// InitiateDisputePayload represents the request to initiate a dispute
type InitiateDisputePayload struct {
	OrderID      uint64   `json:"orderId"`
	Type         string   `json:"type,omitempty"`
	Reason       string   `json:"reason" binding:"required,max=255"`
	Description  string   `json:"description" binding:"max=2000"`
	EvidenceURLs []string `json:"evidenceUrls" binding:"max=9"`
}

// InitiateDisputeResponse represents the response for initiating a dispute
type InitiateDisputeResponse = orderservice.InitiateDisputeResponse

type refundPayload struct {
	Reason       string   `json:"reason" binding:"required,max=255"`
	Description  string   `json:"description" binding:"max=2000"`
	EvidenceURLs []string `json:"evidenceUrls" binding:"max=9"`
}

type refundStatusResponse struct {
	ID           uint64 `json:"id"`
	Status       string `json:"status"`
	Reason       string `json:"reason"`
	Amount       int64  `json:"amount"`
	CreatedAt    string `json:"createdAt"`
	ProcessedAt  string `json:"processedAt,omitempty"`
	RejectReason string `json:"rejectReason,omitempty"`
}

// RegisterDisputeRoutes 注册用户端订单争议路由
func RegisterDisputeRoutes(router gin.IRouter, svc *orderservice.DisputeService, authMiddleware gin.HandlerFunc) {
	h := NewDisputeHandler(svc)
	group := router.Group("/orders")
	group.Use(authMiddleware)
	group.POST("/:id/dispute", h.InitiateDispute)
	group.GET("/:id/disputes", h.GetDisputeDetail)
	group.POST("/:id/refund", h.RequestRefund)
	group.GET("/:id/refund", h.GetRefundStatus)
}

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
	userID := getUserIDFromContext(c)
	if userID == 0 {
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

	orderIDFromPath, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}
	if payload.OrderID == 0 {
		payload.OrderID = orderIDFromPath
	}
	if payload.OrderID != orderIDFromPath {
		respondAPIError(c, apierr.BadRequest("orderId 与路径参数不一致"))
		return
	}

	disputeType := model.DisputeType(payload.Type)
	if disputeType == "" {
		disputeType = model.DisputeTypeServiceQuality
	}

	resp, err := h.svc.InitiateDispute(c.Request.Context(), orderservice.InitiateDisputeRequest{
		OrderID:       payload.OrderID,
		InitiatorID:   userID,
		InitiatorType: "user",
		Type:          disputeType,
		Reason:        payload.Reason,
		EvidenceText:  payload.Description,
		EvidenceURLs:  payload.EvidenceURLs,
	})

	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeValidation) {
			respondAPIError(c, apierr.BadRequest(err.Error()))
			return
		}
		if errors.Is(err, orderservice.ErrCannotInitiateDispute) {
			respondAPIError(c, apierr.Conflict("当前订单状态不能发起纠纷"))
			return
		}
		if errors.Is(err, orderservice.ErrDisputeExists) {
			respondAPIError(c, apierr.Conflict("该订单已存在纠纷"))
			return
		}
		if errors.Is(err, orderservice.ErrOrderNotFound) {
			respondAPIError(c, apierr.NotFound("订单不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("发起纠纷失败").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, "纠纷发起成功", orderservice.InitiateDisputeResponse{
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
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("用户ID不在上下文中"))
		return
	}

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	dispute, err := h.svc.GetDisputeByOrderID(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeNotFound) {
			respondAPIError(c, apierr.NotFound("纠纷不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("获取纠纷详情失败").WithDetails(err.Error()))
		return
	}

	// Verify user owns this dispute
	if dispute.InitiatorID != userID {
		respondAPIError(c, apierr.Forbidden("您只能查看自己的纠纷"))
		return
	}

	respondSuccess(c, "OK", dispute)
}

// RequestRefund creates a refund request for an order (compat endpoint).
// @Summary      Request Refund
// @Description  用户申请退款（兼容接口，底层发起订单纠纷）
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  uint64         true  "订单ID"
// @Param        request  body  refundPayload  true  "退款申请信息"
// @Success      200      {object} map[string]interface{}
// @Failure      400      {object} apierr.APIError
// @Failure      401      {object} apierr.APIError
// @Failure      404      {object} apierr.APIError
// @Failure      409      {object} apierr.APIError
// @Failure      500      {object} apierr.APIError
// @Router       /user/orders/{id}/refund [post]
func (h *DisputeHandler) RequestRefund(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("用户ID不在上下文中"))
		return
	}

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	var payload refundPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidJSONPayload).WithDetails(err.Error()))
		return
	}
	if len(payload.EvidenceURLs) > 9 {
		respondAPIError(c, apierr.BadRequest("最多允许上传9个证据链接"))
		return
	}

	resp, err := h.svc.InitiateDispute(c.Request.Context(), orderservice.InitiateDisputeRequest{
		OrderID:       orderID,
		InitiatorID:   userID,
		InitiatorType: model.DisputeInitiatorUser,
		Type:          model.DisputeTypeServiceQuality,
		Reason:        payload.Reason,
		EvidenceText:  payload.Description,
		EvidenceURLs:  payload.EvidenceURLs,
	})
	if err != nil {
		switch {
		case errors.Is(err, orderservice.ErrDisputeValidation):
			respondAPIError(c, apierr.BadRequest(err.Error()))
		case errors.Is(err, orderservice.ErrCannotInitiateDispute):
			respondAPIError(c, apierr.Conflict("当前订单状态不能申请退款"))
		case errors.Is(err, orderservice.ErrDisputeExists):
			respondAPIError(c, apierr.Conflict("该订单已提交退款/争议申请"))
		case errors.Is(err, orderservice.ErrOrderNotFound):
			respondAPIError(c, apierr.NotFound("订单不存在"))
		default:
			respondAPIError(c, apierr.InternalError("申请退款失败").WithDetails(err.Error()))
		}
		return
	}

	respondSuccess(c, "退款申请已提交", gin.H{
		"id":          resp.DisputeID,
		"status":      "pending",
		"reason":      payload.Reason,
		"createdAt":   nowRFC3339(),
		"traceId":     resp.TraceID,
		"slaDeadline": resp.SLADeadline,
	})
}

// GetRefundStatus returns refund request status for an order (compat endpoint).
// @Summary      Get Refund Status
// @Description  查询订单退款申请状态（兼容接口）
// @Tags         User - Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  uint64  true  "订单ID"
// @Success      200  {object} refundStatusResponse
// @Failure      400  {object} apierr.APIError
// @Failure      401  {object} apierr.APIError
// @Failure      404  {object} apierr.APIError
// @Failure      500  {object} apierr.APIError
// @Router       /user/orders/{id}/refund [get]
func (h *DisputeHandler) GetRefundStatus(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		respondAPIError(c, apierr.Unauthorized("用户ID不在上下文中"))
		return
	}

	orderID, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest(apierr.ErrInvalidID))
		return
	}

	dispute, err := h.svc.GetDisputeByOrderID(c.Request.Context(), orderID)
	if err != nil {
		if errors.Is(err, orderservice.ErrDisputeNotFound) {
			respondAPIError(c, apierr.NotFound("退款申请不存在"))
			return
		}
		respondAPIError(c, apierr.InternalError("获取退款状态失败").WithDetails(err.Error()))
		return
	}

	if dispute.InitiatorID != userID {
		respondAPIError(c, apierr.Forbidden("您只能查看自己的退款申请"))
		return
	}

	status := mapDisputeStatusToRefundStatus(dispute)
	resp := refundStatusResponse{
		ID:        dispute.ID,
		Status:    status,
		Reason:    dispute.Reason,
		Amount:    0,
		CreatedAt: dispute.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if dispute.ResolvedAt != nil {
		resp.ProcessedAt = dispute.ResolvedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	if status == "rejected" {
		resp.RejectReason = dispute.ResolveRemark
	}

	respondSuccess(c, "OK", resp)
}

func mapDisputeStatusToRefundStatus(dispute *model.OrderDispute) string {
	switch dispute.Status {
	case model.DisputeStatusPending:
		return "pending"
	case model.DisputeStatusAssigned, model.DisputeStatusMediating:
		return "processing"
	case model.DisputeStatusResolved:
		if dispute.Resolution == model.ResolutionRefund || dispute.Resolution == model.ResolutionPartial {
			return "refunded"
		}
		return "rejected"
	case model.DisputeStatusRejected, model.DisputeStatusCanceled:
		return "rejected"
	default:
		return "pending"
	}
}

func nowRFC3339() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
}

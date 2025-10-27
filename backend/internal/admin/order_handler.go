package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/service"
)

// OrderHandler 管理订单相关接口。
type OrderHandler struct {
	svc *service.AdminService
}

// NewOrderHandler 创建 Handler。
func NewOrderHandler(svc *service.AdminService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// ListOrders returns a paginated list of orders with filters.
func (h *OrderHandler) ListOrders(c *gin.Context) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page")
		return
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page_size")
		return
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.OrderStatus, 0, len(statusTokens))
	for _, token := range statusTokens {
		statuses = append(statuses, normalizeOrderStatus(token))
	}

	userID, err := queryUint64Ptr(c, "user_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid user_id")
		return
	}
	playerID, err := queryUint64Ptr(c, "player_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid player_id")
		return
	}
	gameID, err := queryUint64Ptr(c, "game_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid game_id")
		return
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid date_from")
		return
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid date_to")
		return
	}

	opts := repository.OrderListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
		UserID:   userID,
		PlayerID: playerID,
		GameID:   gameID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
	}

	orders, pagination, err := h.svc.ListOrders(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Order]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       orders,
		Pagination: pagination,
	})
}

// GetOrder returns a single order by id.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "order not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    order,
	})
}

// UpdateOrder updates order fields such as status and schedule.
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var payload UpdateOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	scheduledStart, err := parseRFC3339Ptr(payload.ScheduledStart)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid scheduled_start format")
		return
	}
	scheduledEnd, err := parseRFC3339Ptr(payload.ScheduledEnd)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid scheduled_end format")
		return
	}

	input := service.UpdateOrderInput{
		Status:         normalizeOrderStatus(payload.Status),
		PriceCents:     payload.PriceCents,
		Currency:       model.Currency(strings.ToUpper(strings.TrimSpace(payload.Currency))),
		ScheduledStart: scheduledStart,
		ScheduledEnd:   scheduledEnd,
		CancelReason:   payload.CancelReason,
	}

	order, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if errors.Is(err, service.ErrValidation) {
		writeJSONError(c, http.StatusBadRequest, "invalid order payload")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "order not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    order,
	})
}

// DeleteOrder deletes an order by id.
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteOrder(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "order not found")
		return
	} else if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// UpdateOrderPayload defines the request body for updating an order.
type UpdateOrderPayload struct {
	Status         string  `json:"status" binding:"required"`
	PriceCents     int64   `json:"price_cents" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`
	ScheduledStart *string `json:"scheduled_start"`
	ScheduledEnd   *string `json:"scheduled_end"`
	CancelReason   string  `json:"cancel_reason"`
}

// PaymentHandler 管理支付记录。
type PaymentHandler struct {
	svc *service.AdminService
}

// NewPaymentHandler 创建 Handler。
func NewPaymentHandler(svc *service.AdminService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// ListPayments returns a paginated list of payments with filters.
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	page, err := queryIntDefault(c, "page", 1)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page")
		return
	}
	pageSize, err := queryIntDefault(c, "page_size", 20)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid page_size")
		return
	}

	statusTokens := parseCSVParams(c.QueryArray("status"))
	statuses := make([]model.PaymentStatus, 0, len(statusTokens))
	for _, token := range statusTokens {
		statuses = append(statuses, model.PaymentStatus(strings.ToLower(token)))
	}

	methodTokens := parseCSVParams(c.QueryArray("method"))
	methods := make([]model.PaymentMethod, 0, len(methodTokens))
	for _, token := range methodTokens {
		methods = append(methods, model.PaymentMethod(strings.ToLower(token)))
	}

	userID, err := queryUint64Ptr(c, "user_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid user_id")
		return
	}
	orderID, err := queryUint64Ptr(c, "order_id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid order_id")
		return
	}
	dateFrom, err := queryTimePtr(c, "date_from")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid date_from")
		return
	}
	dateTo, err := queryTimePtr(c, "date_to")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid date_to")
		return
	}

	opts := repository.PaymentListOptions{
		Page:     page,
		PageSize: pageSize,
		Statuses: statuses,
		Methods:  methods,
		UserID:   userID,
		OrderID:  orderID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	payments, pagination, err := h.svc.ListPayments(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Payment]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       payments,
		Pagination: pagination,
	})
}

// GetPayment returns a single payment by id.
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "payment not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Payment]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    payment,
	})
}

// UpdatePayment updates payment fields such as status and provider info.
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var payload UpdatePaymentPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid json payload")
		return
	}

	paidAt, err := parseRFC3339Ptr(payload.PaidAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid paid_at format")
		return
	}
	refundedAt, err := parseRFC3339Ptr(payload.RefundedAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid refunded_at format")
		return
	}

	input := service.UpdatePaymentInput{
		Status:          model.PaymentStatus(strings.TrimSpace(payload.Status)),
		ProviderTradeNo: payload.ProviderTradeNo,
		ProviderRaw:     payload.ProviderRaw,
		PaidAt:          paidAt,
		RefundedAt:      refundedAt,
	}

	payment, err := h.svc.UpdatePayment(c.Request.Context(), id, input)
	if errors.Is(err, service.ErrValidation) {
		writeJSONError(c, http.StatusBadRequest, "invalid payment payload")
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "payment not found")
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Payment]{
		Success: true,
		Code:    http.StatusOK,
		Message: "updated",
		Data:    payment,
	})
}

// DeletePayment deletes a payment record by id.
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeletePayment(c.Request.Context(), id); errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, "payment not found")
		return
	} else if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// UpdatePaymentPayload defines the request body for updating a payment.
type UpdatePaymentPayload struct {
	Status          string          `json:"status" binding:"required"`
	ProviderTradeNo string          `json:"provider_trade_no"`
	ProviderRaw     json.RawMessage `json:"provider_raw"`
	PaidAt          *string         `json:"paid_at"`
	RefundedAt      *string         `json:"refunded_at"`
}

func parseRFC3339Ptr(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*value))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

package admin

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	apierr "gamelink/internal/handler"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
	"strconv"
)

// OrderHandler 管理订单相关接口�?type OrderHandler struct {
	svc *service.AdminService
}

// NewOrderHandler 创建 Handler�?func NewOrderHandler(svc *service.AdminService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// CreateOrder
// @Summary      创建订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateOrderPayload  true  "订单信息"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var p CreateOrderPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	start, err := parseRFC3339Ptr(p.ScheduledStart)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidScheduledStart)
		return
	}
	end, err := parseRFC3339Ptr(p.ScheduledEnd)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidScheduledEnd)
		return
	}
	var playerID *uint64
	if p.PlayerID != nil {
		playerID = p.PlayerID
	}
	order, err := h.svc.CreateOrder(c.Request.Context(), service.CreateOrderInput{
		UserID:         p.UserID,
		PlayerID:       playerID,
		GameID:         p.GameID,
		Title:          p.Title,
		Description:    p.Description,
		PriceCents:     p.PriceCents,
		Currency:       model.Currency(strings.ToUpper(strings.TrimSpace(p.Currency))),
		ScheduledStart: start,
		ScheduledEnd:   end,
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusCreated, model.APIResponse[*model.Order]{Success: true, Code: http.StatusCreated, Message: "created", Data: order})
}

// AssignOrder
// @Summary      指派订单的陪玩师
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "订单ID"
// @Param        request  body  AssignOrderPayload  true  "指派信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/assign [post]
func (h *OrderHandler) AssignOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var p AssignOrderPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	order, err := h.svc.AssignOrder(c.Request.Context(), id, p.PlayerID)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: order})
}

// ConfirmOrder 确认订单�?// @Summary      确认订单
// @Description  将订单状态从 pending 置为 confirmed
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/confirm [post]
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
			return
		}
	}
	order, err := h.svc.ConfirmOrder(c.Request.Context(), id, payload.Note)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: order})
}

// StartOrder 开始服务�?// @Summary      开始服�?// @Description  将订单状态从 confirmed 置为 in_progress，并记录开始时�?// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/start [post]
func (h *OrderHandler) StartOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
			return
		}
	}
	order, err := h.svc.StartOrder(c.Request.Context(), id, payload.Note)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: order})
}

// CompleteOrder 完成订单�?// @Summary      完成订单
// @Description  将订单状态从 in_progress 置为 completed，并记录完成时间
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
			return
		}
	}
	order, err := h.svc.CompleteOrder(c.Request.Context(), id, payload.Note)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: order})
}

// ListOrders
// @Summary      列出订单
// @Description  根据状�?用户/玩家/游戏和时间范围筛选，支持分页
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        page        query  int     false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status      query  []string  false  "订单状态，可多�?
// @Param        userId     query     int       false  "用户ID"
// @Param        player_id   query  int     false  "玩家ID"
// @Param        gameId     query     int       false  "游戏ID"
// @Param        dateFrom   query     string    false  "开始时�?
// @Param        dateTo     query     string    false  "结束时间"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/orders [get]
//
// ListOrders returns a paginated list of orders with filters.
func (h *OrderHandler) ListOrders(c *gin.Context) {
	opts, ok := buildOrderListOptions(c)
	if !ok {
		return
	}

	orders, pagination, err := h.svc.ListOrders(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	orders = ensureSlice(orders)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Order]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       orders,
		Pagination: pagination,
	})
}

// GetOrder
// @Summary      获取订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        id   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id} [get]
//
// GetOrder returns a single order by id.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, apierr.ErrOrderNotFound)
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

// RefundOrder 处理订单退款�?// @Summary      订单退�?// @Description  将订单状态标记为 refunded，并记录退款金额与原因，同时关联支付退�?// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "订单ID"
// @Param        request  body  orderRefundPayload   true  "退款信�?
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/refund [post]
func (h *OrderHandler) RefundOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload orderRefundPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	order, err := h.svc.RefundOrder(c.Request.Context(), id, service.RefundOrderInput{
		Reason:      payload.Reason,
		AmountCents: payload.AmountCents,
		Note:        payload.Note,
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: order})
}

// GetOrderTimeline 返回订单时间线�?// @Summary      获取订单时间�?// @Description  汇总订单的状态变更、操作日志、支付事件等信息
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/timeline [get]
func (h *OrderHandler) GetOrderTimeline(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	items, err := h.svc.GetOrderTimeline(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]service.OrderTimelineItem]{Success: true, Code: http.StatusOK, Message: "OK", Data: ensureSlice(items)})
}

// ListOrderPayments 返回订单关联的支付记录�?// @Summary      获取订单支付记录
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/payments [get]
func (h *OrderHandler) ListOrderPayments(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	items, err := h.svc.GetOrderPayments(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	payments := ensureSlice(items)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Payment]{Success: true, Code: http.StatusOK, Message: "OK", Data: payments})
}

// ListOrderRefunds 返回订单的退款记录�?// @Summary      获取订单退款记�?// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/refunds [get]
func (h *OrderHandler) ListOrderRefunds(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	items, err := h.svc.GetOrderRefunds(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]service.OrderRefundItem]{Success: true, Code: http.StatusOK, Message: "OK", Data: ensureSlice(items)})
}

// ListOrderReviews 返回订单评价列表�?// @Summary      获取订单评价列表
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/reviews [get]
func (h *OrderHandler) ListOrderReviews(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	items, err := h.svc.GetOrderReviews(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Review]{Success: true, Code: http.StatusOK, Message: "OK", Data: ensureSlice(items)})
}

// UpdateOrder
// @Summary      更新订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                true  "订单ID"
// @Param        request  body  UpdateOrderPayload true  "订单信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id} [put]
//
// UpdateOrder updates order fields such as status and schedule.
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload UpdateOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	scheduledStart, err := parseRFC3339Ptr(payload.ScheduledStart)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidScheduledStart)
		return
	}
	scheduledEnd, err := parseRFC3339Ptr(payload.ScheduledEnd)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidScheduledEnd)
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
	if errors.Is(err, service.ErrOrderInvalidTransition) {
		_ = c.Error(service.ErrOrderInvalidTransition)
		return
	}
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
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

// DeleteOrder
// @Summary      删除订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        id   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id} [delete]
//
// DeleteOrder deletes an order by id.
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	err = h.svc.DeleteOrder(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// ListOrderLogs
// @Summary      获取订单操作日志
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int  true  "订单ID"
// @Param        page       query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action     query  string false "动作过滤" Enums(create,assign_player,update_status,cancel,delete)
// @Param        actor_user_id query int false "操作者用户ID"
// @Param        dateFrom   query     string    false  "开始时�?
// @Param        dateTo     query     string    false  "结束时间"
// @Param        export     query  string false "导出格式" Enums(csv)
// @Param        fields     query  string false "导出列（逗号分隔），默认：id,entity_type,entity_id,actor_user_id,action,reason,metadata,created_at"
// @Param        header_lang query string false "列头语言" Enums(en,zh)
// @Success      200  {object}  map[string]any
// @Router       /admin/orders/{id}/logs [get]
func (h *OrderHandler) ListOrderLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "order", id, opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "order", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.OperationLog]{Success: true, Code: http.StatusOK, Message: "OK", Data: items, Pagination: p})
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

// CreateOrderPayload defines payload for creating an order.
type CreateOrderPayload struct {
	UserID         uint64  `json:"user_id" binding:"required"`
	PlayerID       *uint64 `json:"player_id"`
	GameID         uint64  `json:"game_id" binding:"required"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	PriceCents     int64   `json:"price_cents" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`
	ScheduledStart *string `json:"scheduled_start"`
	ScheduledEnd   *string `json:"scheduled_end"`
}

// AssignOrderPayload defines player assignment.
type AssignOrderPayload struct {
	PlayerID uint64 `json:"player_id" binding:"required"`
}

// PaymentHandler 管理支付记录�?type PaymentHandler struct {
	svc *service.AdminService
}

// NewPaymentHandler 创建 Handler�?func NewPaymentHandler(svc *service.AdminService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// CreatePayment
// @Summary      创建支付记录
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreatePaymentPayload  true  "支付信息"
// @Success      201  {object}  map[string]any
// @Failure      400  {object}  map[string]any
// @Router       /admin/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var p CreatePaymentPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	pay, err := h.svc.CreatePayment(c.Request.Context(), service.CreatePaymentInput{
		OrderID:     p.OrderID,
		UserID:      p.UserID,
		Method:      model.PaymentMethod(strings.ToLower(strings.TrimSpace(p.Method))),
		AmountCents: p.AmountCents,
		Currency:    model.Currency(strings.ToUpper(strings.TrimSpace(p.Currency))),
		ProviderRaw: p.ProviderRaw,
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusCreated, model.APIResponse[*model.Payment]{Success: true, Code: http.StatusCreated, Message: "created", Data: pay})
}

// CapturePayment
// @Summary      确认支付入账
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                     true  "支付ID"
// @Param        request  body  CapturePaymentPayload   true  "入账信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/payments/{id}/capture [post]
func (h *PaymentHandler) CapturePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var p CapturePaymentPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}
	paidAt, err := parseRFC3339Ptr(p.PaidAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPaidAt)
		return
	}
	pay, err := h.svc.CapturePayment(c.Request.Context(), id, service.CapturePaymentInput{
		ProviderTradeNo: p.ProviderTradeNo,
		ProviderRaw:     p.ProviderRaw,
		PaidAt:          paidAt,
	})
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Payment]{Success: true, Code: http.StatusOK, Message: "updated", Data: pay})
}

// ListPayments
// @Summary      列出支付
// @Description  根据状�?方法/用户/订单和时间范围筛选，支持分页
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        page        query  int       false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status      query  []string  false  "支付状�?
// @Param        method      query  []string  false  "支付方式"
// @Param        userId     query     int       false  "用户ID"
// @Param        orderId     query     int       false  "订单ID"
// @Param        dateFrom   query     string    false  "开始时�?
// @Param        dateTo     query     string    false  "结束时间"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /admin/payments [get]
//
// ListPayments returns a paginated list of payments with filters.
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	opts, ok := buildPaymentListOptions(c)
	if !ok {
		return
	}

	payments, pagination, err := h.svc.ListPayments(c.Request.Context(), opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	payments = ensureSlice(payments)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.Payment]{
		Success:    true,
		Code:       http.StatusOK,
		Message:    "OK",
		Data:       payments,
		Pagination: pagination,
	})
}

// GetPayment
// @Summary      获取支付
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        id   path  int  true  "支付ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/payments/{id} [get]
//
// GetPayment returns a single payment by id.
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		writeJSONError(c, http.StatusNotFound, apierr.ErrPaymentNotFound)
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

// UpdatePayment
// @Summary      更新支付
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "支付ID"
// @Param        request  body  UpdatePaymentPayload true  "支付信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/payments/{id} [put]
//
// UpdatePayment updates payment fields such as status and provider info.
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}

	var payload UpdatePaymentPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	paidAt, err := parseRFC3339Ptr(payload.PaidAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidPaidAt)
		return
	}
	refundedAt, err := parseRFC3339Ptr(payload.RefundedAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidRefundedAt)
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
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
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

// DeletePayment
// @Summary      删除支付
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        id   path  int  true  "支付ID"
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/payments/{id} [delete]
//
// DeletePayment deletes a payment record by id.
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	err = h.svc.DeletePayment(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[any]{
		Success: true,
		Code:    http.StatusOK,
		Message: "deleted",
	})
}

// ListPaymentLogs
// @Summary      获取支付操作日志
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id         path   int  true  "支付ID"
// @Param        page       query  int  false "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        action     query  string false "动作过滤" Enums(create,capture,update_status,refund,delete)
// @Param        actor_user_id query int false "操作者用户ID"
// @Param        dateFrom   query     string    false  "开始时�?
// @Param        dateTo     query     string    false  "结束时间"
// @Param        export     query  string false "导出格式" Enums(csv)
// @Param        fields     query  string false "导出列（逗号分隔），默认：id,entity_type,entity_id,actor_user_id,action,reason,metadata,created_at"
// @Param        header_lang query string false "列头语言" Enums(en,zh)
// @Success      200  {object}  map[string]any
// @Router       /admin/payments/{id}/logs [get]
func (h *PaymentHandler) ListPaymentLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateFrom)
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		writeJSONError(c, 400, apierr.ErrInvalidDateTo)
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "payment", id, opts)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "payment", id, items)
		return
	}
	items = ensureSlice(items)
	writeJSON(c, http.StatusOK, model.APIResponse[[]model.OperationLog]{Success: true, Code: http.StatusOK, Message: "OK", Data: items, Pagination: p})
}

// exportOperationLogsCSV writes operation logs as CSV attachment.
func exportOperationLogsCSV(c *gin.Context, entity string, entityID uint64, items []model.OperationLog) {
	// default columns
	allowed := []string{"id", "entity_type", "entity_id", "actor_user_id", "action", "reason", "metadata", "created_at"}
	// parse fields
	rawFields := strings.TrimSpace(c.Query("fields"))
	fields := allowed
	if rawFields != "" {
		req := parseCSVParams([]string{rawFields})
		// validate and keep order
		pick := make([]string, 0, len(req))
		for _, f := range req {
			for _, a := range allowed {
				if f == a {
					pick = append(pick, f)
					break
				}
			}
		}
		if len(pick) > 0 {
			fields = pick
		}
	}

	// header i18n
	lang := strings.ToLower(strings.TrimSpace(c.Query("header_lang")))
	headerMapEn := map[string]string{
		"id": "id", "entity_type": "entity_type", "entity_id": "entity_id", "actor_user_id": "actor_user_id",
		"action": "action", "reason": "reason", "metadata": "metadata", "created_at": "created_at",
	}
	headerMapZh := map[string]string{
		"id": "编号", "entity_type": "实体", "entity_id": "实体ID", "actor_user_id": "操作人ID",
		"action": "动作", "reason": "原因", "metadata": "元数�?, "created_at": "创建时间",
	}
	var header []string
	for _, f := range fields {
		if lang == "zh" {
			header = append(header, headerMapZh[f])
		} else {
			header = append(header, headerMapEn[f])
		}
	}

	filename := entity + "_" + strconv.FormatUint(entityID, 10) + "_logs.csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	// excel-friendly BOM when requested or zh header
	bom := strings.EqualFold(strings.TrimSpace(c.Query("bom")), "true") || lang == "zh"
	if bom {
		_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	w := csv.NewWriter(c.Writer)
	_ = w.Write(header)
	// timezone
	tz := strings.TrimSpace(c.Query("tz"))
	var loc *time.Location
	if tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	for _, it := range items {
		row := make([]string, 0, len(fields))
		for _, f := range fields {
			switch f {
			case "id":
				row = append(row, strconv.FormatUint(it.ID, 10))
			case "entity_type":
				row = append(row, it.EntityType)
			case "entity_id":
				row = append(row, strconv.FormatUint(it.EntityID, 10))
			case "actor_user_id":
				if it.ActorUserID != nil {
					row = append(row, strconv.FormatUint(*it.ActorUserID, 10))
				} else {
					row = append(row, "")
				}
			case "action":
				row = append(row, it.Action)
			case "reason":
				row = append(row, it.Reason)
			case "metadata":
				row = append(row, string(it.MetadataJSON))
			case "created_at":
				t := it.CreatedAt
				if loc != nil {
					t = t.In(loc)
				}
				row = append(row, t.Format(time.RFC3339))
			default:
				row = append(row, "")
			}
		}
		_ = w.Write(row)
	}
	w.Flush()
}

// UpdatePaymentPayload defines the request body for updating a payment.
type UpdatePaymentPayload struct {
	Status          string          `json:"status" binding:"required"`
	ProviderTradeNo string          `json:"provider_trade_no"`
	ProviderRaw     json.RawMessage `json:"provider_raw,omitempty" swaggertype:"string" example:"{\"result\":\"update\"}"`
	PaidAt          *string         `json:"paid_at,omitempty" example:"2025-10-28T10:00:00Z"`
	RefundedAt      *string         `json:"refunded_at,omitempty" example:"2025-10-28T12:00:00Z"`
}

// CreatePaymentPayload defines create payment body.
type CreatePaymentPayload struct {
	OrderID     uint64          `json:"order_id" binding:"required"`
	UserID      uint64          `json:"user_id" binding:"required"`
	Method      string          `json:"method" binding:"required"`
	AmountCents int64           `json:"amount_cents" binding:"required"`
	Currency    string          `json:"currency" binding:"required"`
	ProviderRaw json.RawMessage `json:"provider_raw,omitempty" swaggertype:"string" example:"{\"result\":\"success\"}"`
}

// CapturePaymentPayload defines capture info.
type CapturePaymentPayload struct {
	ProviderTradeNo string          `json:"provider_trade_no"`
	ProviderRaw     json.RawMessage `json:"provider_raw,omitempty" swaggertype:"string" example:"{\"result\":\"captured\"}"`
	PaidAt          *string         `json:"paid_at" example:"2025-10-28T10:00:00Z"`
}

// RefundPayment
// @Summary      退款处�?// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "支付ID"
// @Param        request  body  RefundPaymentPayload   false "退款信�?
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/payments/{id}/refund [post]
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload RefundPaymentPayload
	// optional body
	if c.Request.Body != nil {
		_ = c.ShouldBindJSON(&payload)
	}
	refundedAt, err := parseRFC3339Ptr(payload.RefundedAt)
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidRefundedAt)
		return
	}
	if refundedAt == nil {
		now := time.Now().UTC()
		refundedAt = &now
	}

	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Only allow refund from paid
	input := service.UpdatePaymentInput{
		Status:          model.PaymentStatusRefunded,
		ProviderTradeNo: payload.ProviderTradeNo,
		ProviderRaw:     payload.ProviderRaw,
		PaidAt:          payment.PaidAt,
		RefundedAt:      refundedAt,
	}
	updated, err := h.svc.UpdatePayment(c.Request.Context(), id, input)
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, http.StatusOK, model.APIResponse[*model.Payment]{Success: true, Code: http.StatusOK, Message: "updated", Data: updated})
}

// RefundPaymentPayload defines optional refund fields.
type RefundPaymentPayload struct {
	RefundedAt      *string         `json:"refunded_at,omitempty" example:"2025-10-28T12:00:00Z"`
	ProviderTradeNo string          `json:"provider_trade_no,omitempty"`
	ProviderRaw     json.RawMessage `json:"provider_raw,omitempty" swaggertype:"string" example:"{\"result\":\"refunded\"}"`
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

// ReviewOrder
// @Summary      审核订单（通过/拒绝�?// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "订单ID"
// @Param        request  body  ReviewOrderPayload     true  "审核信息"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/review [post]
func (h *OrderHandler) ReviewOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload ReviewOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	next := model.OrderStatusConfirmed
	cancelReason := ""
	if !payload.Approved {
		next = model.OrderStatusCanceled
		cancelReason = strings.TrimSpace(payload.Reason)
	}

	input := service.UpdateOrderInput{
		Status:         next,
		PriceCents:     order.TotalPriceCents,
		Currency:       order.Currency,
		ScheduledStart: order.ScheduledStart,
		ScheduledEnd:   order.ScheduledEnd,
		CancelReason:   cancelReason,
	}
	updated, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if errors.Is(err, service.ErrOrderInvalidTransition) {
		_ = c.Error(service.ErrOrderInvalidTransition)
		return
	}
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: updated})
}

// CancelOrder
// @Summary      取消订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "订单ID"
// @Param        request  body  CancelOrderPayload   true  "取消原因"
// @Success      200  {object}  map[string]any
// @Failure      404  {object}  map[string]any
// @Router       /admin/orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidID)
		return
	}
	var payload CancelOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		writeJSONError(c, http.StatusBadRequest, apierr.ErrInvalidJSONPayload)
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	input := service.UpdateOrderInput{
		Status:         model.OrderStatusCanceled,
		PriceCents:     order.TotalPriceCents,
		Currency:       order.Currency,
		ScheduledStart: order.ScheduledStart,
		ScheduledEnd:   order.ScheduledEnd,
		CancelReason:   strings.TrimSpace(payload.Reason),
	}
	updated, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if errors.Is(err, service.ErrOrderInvalidTransition) {
		_ = c.Error(service.ErrOrderInvalidTransition)
		return
	}
	if errors.Is(err, service.ErrValidation) {
		_ = c.Error(service.ErrValidation)
		return
	}
	if errors.Is(err, service.ErrNotFound) {
		_ = c.Error(service.ErrNotFound)
		return
	}
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(c, http.StatusOK, model.APIResponse[*model.Order]{Success: true, Code: http.StatusOK, Message: "updated", Data: updated})
}

// ReviewOrderPayload defines approval decision.
type ReviewOrderPayload struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// CancelOrderPayload defines cancel reason.
type CancelOrderPayload struct {
	Reason string `json:"reason"`
}

type orderNotePayload struct {
	Note string `json:"note"`
}

type orderRefundPayload struct {
	Reason      string `json:"reason" binding:"required"`
	AmountCents *int64 `json:"amount_cents,omitempty"`
	Note        string `json:"note"`
}

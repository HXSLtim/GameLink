package admin

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/apierr"
)

// Order 订单模型（类型别名）
type Order = model.Order

// Payment 支付模型（类型别名）
type Payment = model.Payment

// OrderHandler 管理订单相关接口
type OrderHandler struct {
	svc *adminservice.AdminService
}

// NewOrderHandler 创建 Handler
func NewOrderHandler(svc *adminservice.AdminService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// CreateOrder
// @Summary      创建订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreateOrderPayload  true  "订单信息"
// @Success      201  {object}  model.APIResponse[Order]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var p CreateOrderPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}
	start, err := parseRFC3339Ptr(p.ScheduledStart)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid scheduled start time"))
		return
	}
	end, err := parseRFC3339Ptr(p.ScheduledEnd)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid scheduled end time"))
		return
	}
	var playerID *uint64
	if p.PlayerID != nil {
		playerID = p.PlayerID
	}
	order, err := h.svc.CreateOrder(c.Request.Context(), adminservice.CreateOrderInput{
		UserID:          p.UserID,
		PlayerID:        playerID,
		GameID:          p.GameID,
		ItemID:          p.ItemID,
		Title:           p.Title,
		Description:     p.Description,
		TotalPriceCents: p.TotalPriceCents,
		Currency:        model.Currency(strings.ToUpper(strings.TrimSpace(p.Currency))),
		ScheduledStart:  start,
		ScheduledEnd:    end,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("create order failed").WithDetails(err.Error()))
		return
	}
	respondCreated(c, *order)
}

// AssignOrder
// @Summary      指派订单的陪玩师
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                 true  "订单ID"
// @Param        request  body  AssignOrderPayload  true  "指派信息"
// @Success      200  {object}  model.APIResponse[Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/assign [post]
func (h *OrderHandler) AssignOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var p AssignOrderPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}
	order, err := h.svc.AssignOrder(c.Request.Context(), id, p.PlayerID)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, apierr.BadRequest("validation failed").WithDetails(err.Error()))
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order or player not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("assign order failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, *order)
}

// ConfirmOrder 确认订单// @Summary      确认订单
// @Description  将订单状态从 pending 置为 confirmed
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  model.APIResponse[[]Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/confirm [post]
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
			return
		}
	}
	order, err := h.svc.ConfirmOrder(c.Request.Context(), id, payload.Note)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("confirm order failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, *order)
}

// @Description  API endpoint// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/start [post]
func (h *OrderHandler) StartOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
			return
		}
	}
	order, err := h.svc.StartOrder(c.Request.Context(), id, payload.Note)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("start order failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, *order)
}

// CompleteOrder 完成订单// @Summary      完成订单
// @Description  将订单状态从 in_progress 置为 completed，并记录完成时间
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int               true  "订单ID"
// @Param        request  body  orderNotePayload  false "备注（可选）"
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/complete [post]
func (h *OrderHandler) CompleteOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload orderNotePayload
	if c.Request.ContentLength > 0 {
		if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
			respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
			return
		}
	}
	order, err := h.svc.CompleteOrder(c.Request.Context(), id, payload.Note)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("complete order failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, *order)
}

// ListOrders
// @Summary      列出订单
// @Description  API endpoint// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        page        query  int     false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status         query    []string     false  "Status filter"// @Param        userId     query     int       false  "用户ID"
// @Param        player_id   query  int     false  "玩家ID"
// @Param        gameId     query     int       false  "游戏ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.Order]
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
		respondAPIError(c, apierr.InternalError("list orders failed").WithDetails(err.Error()))
		return
	}
	orders = ensureSlice(orders)
	respondList(c, orders, pagination)
}

// GetOrder
// @Summary      获取订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        id   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[[]adminservice.OrderTimelineItem]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id} [get]
//
// GetOrder returns a single order by id.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("get order failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, *order)
}

// @Description  API endpoint// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "订单ID"
// @Param        request        body     orderRefundPayload true   "Request body"// @Success      200  {object}  model.APIResponse[[]Payment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/refund [post]
func (h *OrderHandler) RefundOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload orderRefundPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}
	order, err := h.svc.RefundOrder(c.Request.Context(), id, adminservice.RefundOrderInput{
		Reason:      payload.Reason,
		AmountCents: payload.AmountCents,
		Note:        payload.Note,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("refund order failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, *order)
}

// @Description  API endpoint// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  model.APIResponse[[]adminservice.OrderRefundItem]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/timeline [get]
func (h *OrderHandler) GetOrderTimeline(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.GetOrderTimeline(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("get order timeline failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, ensureSlice(items))
}

// ListOrderPayments 返回订单关联的支付记录// @Summary      获取订单支付记录
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  model.APIResponse[[]Review]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/payments [get]
func (h *OrderHandler) ListOrderPayments(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.GetOrderPayments(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("get order payments failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, ensureSlice(items))
}

// ListOrderRefunds 返回订单的退款记录// @Summary      获取订单退款记// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/refunds [get]
func (h *OrderHandler) ListOrderRefunds(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.GetOrderRefunds(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("get order refunds failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, ensureSlice(items))
}

// ListOrderReviews 返回订单评价列表// @Summary      获取订单评价列表
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "订单ID"
// @Success      200  {object}  model.SuccessResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/reviews [get]
func (h *OrderHandler) ListOrderReviews(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	items, err := h.svc.GetOrderReviews(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("get order reviews failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, ensureSlice(items))
}

// UpdateOrder
// @Summary      更新订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                true  "订单ID"
// @Param        request  body  UpdateOrderPayload true  "订单信息"
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id} [put]
//
// UpdateOrder updates order fields such as status and schedule.
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid order ID"))
		return
	}

	var payload UpdateOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
		return
	}

	scheduledStart, err := parseRFC3339Ptr(payload.ScheduledStart)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid scheduled start time"))
		return
	}
	scheduledEnd, err := parseRFC3339Ptr(payload.ScheduledEnd)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid scheduled end time"))
		return
	}

	input := adminservice.UpdateOrderInput{
		Status:          normalizeOrderStatus(payload.Status),
		TotalPriceCents: payload.TotalPriceCents,
		Currency:        model.Currency(strings.ToUpper(strings.TrimSpace(payload.Currency))),
		ScheduledStart:  scheduledStart,
		ScheduledEnd:    scheduledEnd,
		CancelReason:    payload.CancelReason,
	}

	order, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, apierr.BadRequest("validation failed").WithDetails(err.Error()))
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		if errors.Is(err, adminservice.ErrOrderInvalidTransition) {
			respondAPIError(c, apierr.BadRequest("invalid order status transition").WithDetails(err.Error()))
			return
		}
		respondAPIError(c, apierr.InternalError("update order failed").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, *order)
}

// DeleteOrder
// @Summary      删除订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Param        id   path  int  true  "订单ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[Payment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id} [delete]
//
// DeleteOrder deletes an order by id.
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid order ID"))
		return
	}
	err = h.svc.DeleteOrder(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, apierr.NotFound("order not found"))
			return
		}
		respondAPIError(c, apierr.InternalError("delete order failed").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
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
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export     query  string false "导出格式" Enums(csv)
// @Param        fields     query  string false "Export columns (comma separated)"
// @Param        header_lang query string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[[]model.Payment]
// @Router       /admin/orders/{id}/logs [get]
func (h *OrderHandler) ListOrderLogs(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid order ID"))
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid user ID"))
		return
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid date from"))
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid date to"))
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "order", id, opts)
	if err != nil {
		respondAPIError(c, apierr.InternalError("list order logs failed").WithDetails(err.Error()))
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "order", id, items)
		return
	}
	respondList(c, ensureSlice(items), p)
}

// UpdateOrderPayload defines the request body for updating an order.
type UpdateOrderPayload struct {
	Status          string  `json:"status" binding:"required"`
	TotalPriceCents int64   `json:"total_price_cents" binding:"required"`
	Currency        string  `json:"currency" binding:"required"`
	ScheduledStart  *string `json:"scheduled_start"`
	ScheduledEnd    *string `json:"scheduled_end"`
	CancelReason    string  `json:"cancel_reason"`
}

// CreateOrderPayload defines payload for creating an order.
type CreateOrderPayload struct {
	UserID          uint64  `json:"user_id" binding:"required"`
	PlayerID        *uint64 `json:"player_id"`
	GameID          uint64  `json:"game_id" binding:"required"`
	ItemID          uint64  `json:"item_id" binding:"required"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TotalPriceCents int64   `json:"total_price_cents" binding:"required"`
	Currency        string  `json:"currency" binding:"required"`
	ScheduledStart  *string `json:"scheduled_start"`
	ScheduledEnd    *string `json:"scheduled_end"`
}

// AssignOrderPayload defines player assignment.
type AssignOrderPayload struct {
	PlayerID uint64 `json:"player_id" binding:"required"`
}

// PaymentHandler 管理支付记录
type PaymentHandler struct {
	svc *adminservice.AdminService
}

// NewPaymentHandler 创建 Handler
func NewPaymentHandler(svc *adminservice.AdminService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// CreatePayment
// @Summary      创建支付记录
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  CreatePaymentPayload  true  "支付信息"
// @Success      201  {object}  model.APIResponse[Payment]
// @Failure      400  {object}  model.ErrorResponse
// @Router       /admin/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var p CreatePaymentPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}
	pay, err := h.svc.CreatePayment(c.Request.Context(), adminservice.CreatePaymentInput{
		OrderID:     p.OrderID,
		UserID:      p.UserID,
		Method:      model.PaymentMethod(strings.ToLower(strings.TrimSpace(p.Method))),
		AmountCents: p.AmountCents,
		Currency:    model.Currency(strings.ToUpper(strings.TrimSpace(p.Currency))),
		ProviderRaw: p.ProviderRaw,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("create payment failed").WithDetails(err.Error()))
		return
	}
	respondCreated(c, pay)
}

// CapturePayment
// @Summary      确认支付入账
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                     true  "支付ID"
// @Param        request  body  CapturePaymentPayload   true  "入账信息"
// @Success      200  {object}  model.APIResponse[model.Payment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/payments/{id}/capture [post]
func (h *PaymentHandler) CapturePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid order ID"))
		return
	}
	var p CapturePaymentPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}
	paidAt, err := parseRFC3339Ptr(p.PaidAt)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid paid at time"))
		return
	}
	pay, err := h.svc.CapturePayment(c.Request.Context(), id, adminservice.CapturePaymentInput{
		ProviderTradeNo: p.ProviderTradeNo,
		ProviderRaw:     p.ProviderRaw,
		PaidAt:          paidAt,
	})
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("capture payment failed").WithDetails(err.Error()))
		return
	}
	respondUpdated(c, pay)
}

// ListPayments
// @Summary      列出支付
// @Description  API endpoint// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        page        query  int       false  "页码"
// @Param        pageSize   query     int       false  "每页数量"
// @Param        status         query    []string     false  "Status filter"// @Param        method      query  []string  false  "支付方式"
// @Param        userId     query     int       false  "用户ID"
// @Param        orderId     query     int       false  "订单ID"
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.Payment]
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
		respondAPIError(c, apierr.InternalError("list payments failed").WithDetails(err.Error()))
		return
	}
	respondList(c, ensureSlice(payments), pagination)
}

// GetPayment
// @Summary      获取支付详情
// @Description  获取支付记录的完整详细信息，包括关联订单信息、用户信息和支付时间线
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        id   path  int  true  "支付ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.Payment]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/payments/{id} [get]
//
// GetPayment returns a single payment by id with related order and user info.
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid payment ID"))
		return
	}
	payment, err := h.svc.GetPaymentWithRelations(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get payment failed").WithDetails(err.Error()))
		return
	}
	respondSuccess(c, payment)
}

// UpdatePayment
// @Summary      更新支付
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "支付ID"
// @Param        request  body  UpdatePaymentPayload true  "支付信息"
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/payments/{id} [put]
//
// UpdatePayment updates payment fields such as status and provider info.
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid order ID"))
		return
	}

	var payload UpdatePaymentPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
		return
	}

	paidAt, err := parseRFC3339Ptr(payload.PaidAt)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid paid at time"))
		return
	}
	refundedAt, err := parseRFC3339Ptr(payload.RefundedAt)
	if err != nil {
		respondAPIError(c, apierr.BadRequest("invalid refunded at time"))
		return
	}

	input := adminservice.UpdatePaymentInput{
		Status:          model.PaymentStatus(strings.TrimSpace(payload.Status)),
		ProviderTradeNo: payload.ProviderTradeNo,
		ProviderRaw:     payload.ProviderRaw,
		PaidAt:          paidAt,
		RefundedAt:      refundedAt,
	}
	payment, err := h.svc.UpdatePayment(c.Request.Context(), id, input)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("update payment failed").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, payment)
}

// DeletePayment
// @Summary      删除支付
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Param        id   path  int  true  "支付ID"
// @Produce      json
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/payments/{id} [delete]
//
// DeletePayment deletes a payment record by id.
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	err := h.svc.DeletePayment(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("delete payment failed").WithDetails(err.Error()))
		return
	}

	respondDeleted(c)
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
// @Param        dateFrom       query    string       false  "Start date (YYYY-MM-DD)"// @Param        dateTo     query     string    false  "End date (YYYY-MM-DD)"
// @Param        export     query  string false "导出格式" Enums(csv)
// @Param        fields     query  string false "Export columns (comma separated)"
// @Param        header_lang query string false "列头语言" Enums(en,zh)
// @Success      200  {object}  model.APIResponse[model.Order]
// @Router       /admin/payments/{id}/logs [get]
func (h *PaymentHandler) ListPaymentLogs(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	page, pageSize, ok := parsePagination(c)
	if !ok {
		return
	}
	var actorID *uint64
	if v, err := queryUint64Ptr(c, "actor_user_id"); err == nil {
		actorID = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid user ID"))
		return
	}
	var dateFrom, dateTo *time.Time
	if v, err := queryTimePtr(c, "date_from"); err == nil {
		dateFrom = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid date from"))
		return
	}
	if v, err := queryTimePtr(c, "date_to"); err == nil {
		dateTo = v
	} else {
		respondAPIError(c, apierr.BadRequest("invalid date to"))
		return
	}
	opts := repository.OperationLogListOptions{Page: page, PageSize: pageSize, Action: strings.TrimSpace(c.Query("action")), ActorUserID: actorID, DateFrom: dateFrom, DateTo: dateTo}
	items, p, err := h.svc.ListOperationLogs(c.Request.Context(), "payment", id, opts)
	if err != nil {
		respondAPIError(c, apierr.InternalError("list payment logs failed").WithDetails(err.Error()))
		return
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("export")), "csv") {
		exportOperationLogsCSV(c, "payment", id, items)
		return
	}
	respondList(c, ensureSlice(items), p)
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

// RefundPayment processes a refund request with amount validation.
// @Summary      发起退款
// @Description  处理退款请求，验证退款金额不超过剩余可退款金额
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "支付ID"
// @Param        request  body  RefundPaymentPayload   true  "退款信息"
// @Success      200  {object}  model.APIResponse[model.Payment]
// @Failure      400  {object}  model.ErrorResponse "退款金额无效或超过可退款金额"
// @Failure      404  {object}  model.ErrorResponse "支付记录不存在"
// @Router       /admin/payments/{id}/refund [post]
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	var payload RefundPaymentPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(err.Error()))
		return
	}

	// Get the payment to validate refund amount
	payment, err := h.svc.GetPayment(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get payment failed").WithDetails(err.Error()))
		return
	}

	// Validate refund amount using model validation
	// Requirements: 9.1, 9.2, 9.3
	if err := payment.ValidateRefundAmount(payload.AmountCents); err != nil {
		if refundErr, ok := err.(*model.RefundValidationError); ok {
			respondAPIError(c, apierr.BadRequest(refundErr.Message).WithDetails(refundErr.Code))
			return
		}
		respondAPIError(c, apierr.BadRequest("invalid refund amount"))
		return
	}

	// Get operator ID from context if available
	var operatorID *uint64
	if uid := c.GetUint64("user_id"); uid != 0 {
		operatorID = &uid
	}

	// Process the refund
	refundedAt, _ := parseRFC3339Ptr(payload.RefundedAt)
	if refundedAt == nil {
		now := time.Now().UTC()
		refundedAt = &now
	}

	// Update payment with refund amount
	payment.RefundedAmountCents += payload.AmountCents
	payment.RefundedAt = refundedAt
	payment.ProviderTradeNo = payload.ProviderTradeNo
	payment.ProviderRaw = payload.ProviderRaw

	// Check if fully refunded
	if payment.IsFullyRefunded() {
		payment.Status = model.PaymentStatusRefunded
	}

	// Use the existing UpdatePayment with proper status transition
	input := adminservice.UpdatePaymentInput{
		Status:              payment.Status,
		ProviderTradeNo:     payment.ProviderTradeNo,
		ProviderRaw:         payment.ProviderRaw,
		PaidAt:              payment.PaidAt,
		RefundedAt:          refundedAt,
		RefundedAmountCents: &payment.RefundedAmountCents,
	}

	updated, err := h.svc.UpdatePaymentWithRefund(c.Request.Context(), id, input, payload.AmountCents, payload.Reason, operatorID)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("refund processing failed").WithDetails(err.Error()))
		return
	}

	respondSuccessWithMsg(c, "refund processed", updated)
}

// RefundPaymentPayload defines refund request fields.
// Requirements: 2.1, 9.1, 9.2, 9.3
type RefundPaymentPayload struct {
	AmountCents     int64           `json:"amount_cents" binding:"required,gt=0" example:"1000"`
	Reason          string          `json:"reason" binding:"required" example:"Customer requested refund"`
	Note            string          `json:"note,omitempty" example:"Internal note"`
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

// GetRefundHistory returns the refund history for a payment.
// @Summary      获取退款历史
// @Description  获取支付记录的所有退款操作历史
// @Tags         Admin/Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id   path  int  true  "支付ID"
// @Success      200  {object}  model.APIResponse[[]model.OperationLog]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/payments/{id}/refunds [get]
func (h *PaymentHandler) GetRefundHistory(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}

	// Verify payment exists
	_, err := h.svc.GetPayment(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get payment failed").WithDetails(err.Error()))
		return
	}

	// Get refund-related operation logs for this payment
	logs, _, err := h.svc.GetPaymentLogs(c.Request.Context(), id, repository.OperationLogListOptions{
		Page:     1,
		PageSize: 100,
		Action:   string(model.OpActionRefund),
	})
	if err != nil {
		respondAPIError(c, apierr.InternalError("get refund history failed").WithDetails(err.Error()))
		return
	}

	respondSuccess(c, logs)
}

// ReviewOrder
// @Summary      审核订单（通过/拒绝// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                    true  "订单ID"
// @Param        request  body  ReviewOrderPayload     true  "审核信息"
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/review [post]
func (h *OrderHandler) ReviewOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload ReviewOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get order failed").WithDetails(err.Error()))
		return
	}

	next := model.OrderStatusConfirmed
	cancelReason := ""
	if !payload.Approved {
		next = model.OrderStatusCanceled
		cancelReason = strings.TrimSpace(payload.Reason)
	}

	input := adminservice.UpdateOrderInput{
		Status:          next,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    cancelReason,
	}
	updated, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("review order failed").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, *updated)
}

// CancelOrder
// @Summary      取消订单
// @Tags         Admin/Orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path  int                  true  "订单ID"
// @Param        request  body  CancelOrderPayload   true  "取消原因"
// @Success      200  {object}  model.APIResponse[model.Order]
// @Failure      404  {object}  model.ErrorResponse
// @Router       /admin/orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, ok := ParseIDAndRespond(c, "id")
	if !ok {
		return
	}
	var payload CancelOrderPayload
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		respondAPIError(c, apierr.BadRequest("invalid request payload").WithDetails(bindErr.Error()))
		return
	}

	order, err := h.svc.GetOrder(c.Request.Context(), id)
	if err != nil {
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("get order failed").WithDetails(err.Error()))
		return
	}

	input := adminservice.UpdateOrderInput{
		Status:          model.OrderStatusCanceled,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    strings.TrimSpace(payload.Reason),
	}
	updated, err := h.svc.UpdateOrder(c.Request.Context(), id, input)
	if err != nil {
		if apierr.IsValidationError(err) {
			respondAPIError(c, err)
			return
		}
		if apierr.IsNotFound(err) {
			respondAPIError(c, err)
			return
		}
		respondAPIError(c, apierr.InternalError("cancel order failed").WithDetails(err.Error()))
		return
	}

	respondUpdated(c, *updated)
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

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"log/slog"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/apierr"
	"gamelink/pkg/logging"
)

// --- Order management ---

// CreateOrderInput 创建订单请求。
type CreateOrderInput struct {
	UserID          uint64
	PlayerID        *uint64
	GameID          uint64
	ItemID          uint64 // 服务项目ID (必填)
	Title           string
	Description     string
	TotalPriceCents int64
	Currency        model.Currency
	ScheduledStart  *time.Time
	ScheduledEnd    *time.Time
}

// CreateOrder 新建订单，默认状态为 pending。
func (s *AdminService) CreateOrder(ctx context.Context, in CreateOrderInput) (*model.Order, error) {
	// ✅ 数据安全修复: 验证所有必填字段，包括ItemID
	if in.UserID == 0 || in.GameID == 0 || in.ItemID == 0 || in.TotalPriceCents < 0 || !model.IsValidCurrency(in.Currency) {
		return nil, ErrValidation
	}
	if in.ScheduledStart != nil && in.ScheduledEnd != nil && in.ScheduledEnd.Before(*in.ScheduledStart) {
		return nil, ErrValidation
	}

	// 验证服务项目是否存在
	serviceItem, err := s.serviceItems.Get(ctx, in.ItemID)
	if err != nil {
		return nil, apierr.BadRequest("服务项目不存在")
	}

	// 验证服务项目是否激活
	if !serviceItem.IsActive {
		return nil, apierr.BadRequest("服务项目已停用")
	}

	// 可选: 验证服务项目与游戏的关联性
	if serviceItem.GameID != nil && *serviceItem.GameID != in.GameID {
		return nil, apierr.BadRequest("服务项目与游戏不匹配")
	}

	// 验证陪玩师是否存在
	if in.PlayerID != nil && *in.PlayerID != 0 {
		if _, err := s.players.Get(ctx, *in.PlayerID); err != nil {
			return nil, err
		}
	}

	gameID := in.GameID
	order := &model.Order{
		OrderNo:         model.GenerateEscortOrderNo(),
		UserID:          in.UserID,
		ItemID:          in.ItemID, // ✅ 修复: 使用传入的ItemID而不是硬编码
		GameID:          &gameID,
		Quantity:        1,
		UnitPriceCents:  in.TotalPriceCents,
		TotalPriceCents: in.TotalPriceCents,
		Currency:        in.Currency,
		Status:          model.OrderStatusPending,
		Title:           strings.TrimSpace(in.Title),
		Description:     strings.TrimSpace(in.Description),
		ScheduledStart:  in.ScheduledStart,
		ScheduledEnd:    in.ScheduledEnd,
	}
	if in.PlayerID != nil {
		order.PlayerID = in.PlayerID
	}
	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionCreate), map[string]any{"status": order.Status})
	return order, nil
}

// AssignOrder 指派陪玩师。
func (s *AdminService) AssignOrder(ctx context.Context, id uint64, playerID uint64) (*model.Order, error) {
	if playerID == 0 {
		return nil, ErrValidation
	}
	if _, err := s.players.Get(ctx, playerID); err != nil {
		return nil, err
	}
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// 不允许在完成/取消/退款后指派
	switch order.Status {
	case model.OrderStatusCompleted, model.OrderStatusCanceled, model.OrderStatusRefunded:
		return nil, ErrValidation
	}
	order.SetPlayerID(playerID)
	if err := s.orders.Update(ctx, order); err != nil {
		return nil, WrapError(err, "update order")
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(model.OpActionAssignPlayer), map[string]any{"player_id": playerID})
	return order, nil
}

// UpdateOrderInput 用于更新订单状态。
type UpdateOrderInput struct {
	Status            model.OrderStatus
	TotalPriceCents   int64
	Currency          model.Currency
	ScheduledStart    *time.Time
	ScheduledEnd      *time.Time
	CancelReason      string
	StartedAt         *time.Time
	CompletedAt       *time.Time
	RefundAmountCents *int64
	RefundReason      string
	RefundedAt        *time.Time
	Note              string
}

// RefundOrderInput 描述退款请求。
type RefundOrderInput struct {
	Reason      string
	AmountCents *int64
	Note        string
}

// OrderTimelineItem 组合订单历史时间线。
type OrderTimelineItem struct {
	ID           uint64         `json:"id"`
	OrderID      uint64         `json:"order_id"`
	PaymentID    *uint64        `json:"payment_id,omitempty"`
	EventType    string         `json:"event_type"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	Operator     string         `json:"operator,omitempty"`
	OperatorRole string         `json:"operator_role,omitempty"`
	OperatorID   *uint64        `json:"operator_id,omitempty"`
	StatusBefore string         `json:"status_before,omitempty"`
	StatusAfter  string         `json:"status_after,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// OrderRefundItem 描述订单退款记录。
type OrderRefundItem struct {
	ID          uint64     `json:"id"`
	OrderID     uint64     `json:"order_id"`
	PaymentID   uint64     `json:"payment_id"`
	AmountCents int64      `json:"amount_cents"`
	Reason      string     `json:"reason,omitempty"`
	Status      string     `json:"status"`
	Method      string     `json:"refund_method"`
	Note        string     `json:"note,omitempty"`
	RefundedAt  *time.Time `json:"refunded_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ListOrders 列出订单。
func (s *AdminService) ListOrders(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, *model.Pagination, error) {
	normalized := opts
	normalized.Page = repository.NormalizePage(opts.Page)
	normalized.PageSize = repository.NormalizePageSize(opts.PageSize)

	orders, total, err := s.orders.List(ctx, normalized)
	if err != nil {
		return nil, nil, err
	}

	pagination := buildPagination(normalized.Page, normalized.PageSize, total)
	return orders, &pagination, nil
}

// GetOrder 获取订单详情。
func (s *AdminService) GetOrder(ctx context.Context, id uint64) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}
	return order, nil
}

// UpdateOrder 更新订单信息。
func (s *AdminService) UpdateOrder(ctx context.Context, id uint64, input UpdateOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}

	if !isValidOrderStatus(input.Status) {
		return nil, ErrValidation
	}
	if !model.IsValidCurrency(input.Currency) {
		return nil, ErrValidation
	}
	if input.TotalPriceCents < 0 {
		return nil, ErrValidation
	}
	if input.ScheduledStart != nil && input.ScheduledEnd != nil && input.ScheduledEnd.Before(*input.ScheduledStart) {
		return nil, ErrValidation
	}

	// state machine guard
	if !isAllowedOrderTransition(order.Status, input.Status) {
		return nil, ErrOrderInvalidTransition
	}
	if input.Status == model.OrderStatusCompleted && order.Status != model.OrderStatusCompleted {
		hasPaid, err := s.hasPaidPayment(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		if !hasPaid {
			return nil, apierr.BadRequest("order must be paid before completion")
		}
	}

	prevStatus := order.Status
	if input.Status == model.OrderStatusCanceled {
		if err := s.syncCancellationPayments(ctx, order); err != nil {
			return nil, err
		}
	}

	order.Status = input.Status
	order.TotalPriceCents = input.TotalPriceCents
	order.Currency = input.Currency
	order.ScheduledStart = input.ScheduledStart
	order.ScheduledEnd = input.ScheduledEnd
	order.CancelReason = strings.TrimSpace(input.CancelReason)
	if input.StartedAt != nil {
		order.StartedAt = input.StartedAt
	}
	if input.CompletedAt != nil {
		order.CompletedAt = input.CompletedAt
	}
	if input.RefundAmountCents != nil {
		order.RefundAmountCents = *input.RefundAmountCents
	}
	if input.RefundReason != "" || input.RefundAmountCents != nil {
		order.RefundReason = strings.TrimSpace(input.RefundReason)
	}
	if input.RefundedAt != nil {
		order.RefundedAt = input.RefundedAt
	}

	if err := s.orders.Update(ctx, order); err != nil {
		return nil, err
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	action := model.OpActionUpdateStatus
	switch order.Status {
	case model.OrderStatusCanceled:
		action = model.OpActionCancel
	case model.OrderStatusRefunded:
		action = model.OpActionRefund
	default:
		switch {
		case prevStatus == model.OrderStatusPending && order.Status == model.OrderStatusConfirmed:
			action = model.OpActionConfirm
		case prevStatus == model.OrderStatusConfirmed && order.Status == model.OrderStatusInProgress:
			action = model.OpActionStart
		case prevStatus == model.OrderStatusInProgress && order.Status == model.OrderStatusCompleted:
			action = model.OpActionComplete
		}
	}
	meta := map[string]any{
		"status":      order.Status,
		"from_status": prevStatus,
	}
	if order.CancelReason != "" {
		meta["reason"] = order.CancelReason
	}
	if input.Note != "" {
		meta["note"] = strings.TrimSpace(input.Note)
	}
	if order.StartedAt != nil {
		meta["started_at"] = order.StartedAt.Format(time.RFC3339)
	}
	if order.CompletedAt != nil {
		meta["completed_at"] = order.CompletedAt.Format(time.RFC3339)
	}
	if input.RefundAmountCents != nil {
		meta["refund_amount_cents"] = order.RefundAmountCents
	}
	if order.RefundReason != "" {
		meta["refund_reason"] = order.RefundReason
	}
	if order.RefundedAt != nil {
		meta["refunded_at"] = order.RefundedAt.Format(time.RFC3339)
	}
	s.appendLogAsync(ctx, string(model.OpEntityOrder), order.ID, string(action), meta)
	if rid, ok := logging.RequestIDFromContext(ctx); ok {
		slog.Info("order_status_changed", slog.Uint64("order_id", order.ID), slog.String("status", string(order.Status)), slog.String("request_id", rid))
	} else {
		slog.Info("order_status_changed", slog.Uint64("order_id", order.ID), slog.String("status", string(order.Status)))
	}
	return order, nil
}

func (s *AdminService) syncCancellationPayments(ctx context.Context, order *model.Order) error {
	payments, err := s.payments.GetByOrderID(ctx, order.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	var refundedAmount int64

	for _, payment := range payments {
		if payment.Status != model.PaymentStatusPaid {
			continue
		}
		payment.Status = model.PaymentStatusRefunded
		payment.RefundedAmountCents = payment.AmountCents
		payment.RefundedAt = &now
		if err := s.payments.Update(ctx, &payment); err != nil {
			return err
		}
		refundedAmount += payment.AmountCents
	}

	if refundedAmount > 0 {
		order.RefundAmountCents = refundedAmount
		order.RefundedAt = &now
		if strings.TrimSpace(order.RefundReason) == "" {
			order.RefundReason = "order canceled"
		}
	}

	return nil
}

// ConfirmOrder 将订单从 pending 确认到 confirmed。
func (s *AdminService) ConfirmOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       order.StartedAt,
		CompletedAt:     order.CompletedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
	})
}

// StartOrder 将订单置为进行中，并记录实际开始时间。
func (s *AdminService) StartOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	note = strings.TrimSpace(note)
	startedAt := time.Now().UTC()
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:          model.OrderStatusInProgress,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       &startedAt,
		CompletedAt:     order.CompletedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
	})
}

// CompleteOrder 完成订单服务，并记录完成时间。
func (s *AdminService) CompleteOrder(ctx context.Context, id uint64, note string) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	hasPaid, err := s.hasPaidPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	if !hasPaid {
		return nil, apierr.BadRequest("order must be paid before completion")
	}
	note = strings.TrimSpace(note)
	completedAt := time.Now().UTC()
	return s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:          model.OrderStatusCompleted,
		TotalPriceCents: order.TotalPriceCents,
		Currency:        order.Currency,
		ScheduledStart:  order.ScheduledStart,
		ScheduledEnd:    order.ScheduledEnd,
		CancelReason:    order.CancelReason,
		StartedAt:       order.StartedAt,
		CompletedAt:     &completedAt,
		RefundReason:    order.RefundReason,
		RefundedAt:      order.RefundedAt,
		Note:            note,
	})
}

// RefundOrder 执行退款并记录退款信息。
func (s *AdminService) RefundOrder(ctx context.Context, id uint64, input RefundOrderInput) (*model.Order, error) {
	order, err := s.orders.Get(ctx, id)
	if err != nil {
		return nil, WrapError(err, "get order")
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, apierr.BadRequest("reason is required")
	}
	switch order.Status {
	case model.OrderStatusCompleted, model.OrderStatusInProgress, model.OrderStatusConfirmed, model.OrderStatusRefunded:
		// allowed
	default:
		return nil, apierr.BadRequest("invalid order status for refund")
	}

	// 先汇总当前支付退款信息，避免订单字段与支付字段不一致
	payments, err := s.listPaymentsByOrder(ctx, id)
	if err != nil {
		return nil, WrapError(err, "list payments by order")
	}

	var paymentRefundedCents int64
	var paymentRemainingCents int64
	for _, pay := range payments {
		if pay.Status != model.PaymentStatusPaid && pay.Status != model.PaymentStatusRefunded {
			continue
		}
		refunded := clampRefundedAmount(pay.RefundedAmountCents, pay.AmountCents)
		paymentRefundedCents += refunded
		paymentRemainingCents += maxInt64(pay.AmountCents-refunded, 0)
	}

	if len(payments) == 0 && order.RefundAmountCents == 0 {
		return nil, apierr.BadRequest("no payment record for refund")
	}

	baseRefunded := maxInt64(order.RefundAmountCents, paymentRefundedCents)
	remainingOrderRefundable := order.TotalPriceCents - baseRefunded
	if remainingOrderRefundable <= 0 {
		return nil, apierr.BadRequest("order already fully refunded")
	}

	amount := remainingOrderRefundable
	if input.AmountCents != nil {
		if *input.AmountCents <= 0 || *input.AmountCents > remainingOrderRefundable {
			return nil, apierr.BadRequest("invalid refund amount")
		}
		amount = *input.AmountCents
	}

	refundedAt := time.Now().UTC()
	note := strings.TrimSpace(input.Note)

	// 优先把退款分摊到支付记录（最多分摊到可退款余额）
	remainingToAllocate := minInt64(amount, paymentRemainingCents)
	for _, pay := range payments {
		if remainingToAllocate <= 0 {
			break
		}
		if pay.Status != model.PaymentStatusPaid && pay.Status != model.PaymentStatusRefunded {
			continue
		}
		currentRefunded := clampRefundedAmount(pay.RefundedAmountCents, pay.AmountCents)
		paymentRemaining := pay.AmountCents - currentRefunded
		if paymentRemaining <= 0 {
			continue
		}

		allocate := minInt64(paymentRemaining, remainingToAllocate)
		newRefundedTotal := currentRefunded + allocate
		nextStatus := model.PaymentStatusPaid
		if newRefundedTotal >= pay.AmountCents {
			nextStatus = model.PaymentStatusRefunded
		}

		_, updErr := s.UpdatePayment(ctx, pay.ID, UpdatePaymentInput{
			Status:              nextStatus,
			ProviderTradeNo:     pay.ProviderTradeNo,
			ProviderRaw:         pay.ProviderRaw,
			PaidAt:              pay.PaidAt,
			RefundedAt:          &refundedAt,
			RefundedAmountCents: &newRefundedTotal,
		})
		if updErr != nil && !errors.Is(updErr, ErrValidation) {
			return nil, WrapError(updErr, "update payment")
		}

		if s.refunds != nil {
			var operatorID *uint64
			if uid, ok := logging.ActorUserIDFromContext(ctx); ok {
				id := uid
				operatorID = &id
			}
			record := &model.RefundRecord{
				PaymentID:       pay.ID,
				OrderID:         id,
				UserID:          pay.UserID,
				AmountCents:     allocate,
				Reason:          strings.TrimSpace(reason),
				Status:          model.RefundStatusProcessed,
				ProviderTradeNo: strings.TrimSpace(pay.ProviderTradeNo),
				OperatorID:      operatorID,
				RefundedAt:      &refundedAt,
				Note:            note,
			}
			_ = s.refunds.Create(ctx, record)
		}
		remainingToAllocate -= allocate
	}

	updatedRefundedTotal := baseRefunded + amount
	if updatedRefundedTotal > order.TotalPriceCents {
		updatedRefundedTotal = order.TotalPriceCents
	}

	nextOrderStatus := order.Status
	if updatedRefundedTotal >= order.TotalPriceCents {
		nextOrderStatus = model.OrderStatusRefunded
	}

	updatedOrder, err := s.UpdateOrder(ctx, id, UpdateOrderInput{
		Status:            nextOrderStatus,
		TotalPriceCents:   order.TotalPriceCents,
		Currency:          order.Currency,
		ScheduledStart:    order.ScheduledStart,
		ScheduledEnd:      order.ScheduledEnd,
		CancelReason:      order.CancelReason,
		StartedAt:         order.StartedAt,
		CompletedAt:       order.CompletedAt,
		RefundAmountCents: &updatedRefundedTotal,
		RefundReason:      reason,
		RefundedAt:        &refundedAt,
		Note:              note,
	})
	if err != nil {
		return nil, WrapError(err, "update order")
	}

	return updatedOrder, nil
}

// GetOrderPayments 返回订单下的所有支付记录。
func (s *AdminService) GetOrderPayments(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	return s.listPaymentsByOrder(ctx, orderID)
}

// GetOrderRefunds 汇总订单退款记录（基于支付信息与订单字段）。
func (s *AdminService) GetOrderRefunds(ctx context.Context, orderID uint64) ([]OrderRefundItem, error) {
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if s.refunds != nil {
		records, err := s.refunds.ListByOrderID(ctx, orderID)
		if err == nil && len(records) > 0 {
			paymentByID := make(map[uint64]model.Payment, len(payments))
			for _, p := range payments {
				paymentByID[p.ID] = p
			}
			result := make([]OrderRefundItem, 0, len(records))
			for _, rec := range records {
				method := "unknown"
				status := mapRefundRecordStatus(rec.Status)
				if pay, ok := paymentByID[rec.PaymentID]; ok {
					method = string(pay.Method)
					if rec.Status == model.RefundStatusProcessed && pay.Status != model.PaymentStatusRefunded {
						status = "partial"
					}
				}
				refundedAt := rec.RefundedAt
				result = append(result, OrderRefundItem{
					ID:          rec.ID,
					OrderID:     orderID,
					PaymentID:   rec.PaymentID,
					AmountCents: rec.AmountCents,
					Method:      method,
					Status:      status,
					RefundedAt:  refundedAt,
					CreatedAt:   rec.CreatedAt,
					Reason:      rec.Reason,
					Note:        rec.Note,
				})
			}
			sort.Slice(result, func(i, j int) bool { return result[i].CreatedAt.Before(result[j].CreatedAt) })
			return result, nil
		}
	}

	result := make([]OrderRefundItem, 0)
	var aggregatedRefundedCents int64
	for _, pay := range payments {
		refundedAmount := clampRefundedAmount(pay.RefundedAmountCents, pay.AmountCents)
		if refundedAmount <= 0 {
			continue
		}
		aggregatedRefundedCents += refundedAmount
		status := "partial"
		if pay.IsFullyRefunded() {
			status = mapRefundStatus(pay.Status)
		}
		item := OrderRefundItem{
			ID:          pay.ID,
			OrderID:     orderID,
			PaymentID:   pay.ID,
			AmountCents: refundedAmount,
			Method:      string(pay.Method),
			Status:      status,
			RefundedAt:  pay.RefundedAt,
			CreatedAt:   pay.CreatedAt,
			Reason:      order.RefundReason,
			Note:        order.RefundReason,
		}
		result = append(result, item)
	}

	// 如果订单存在“人工退款差额”（未映射到支付记录），补充一条摘要
	manualRefundCents := order.RefundAmountCents - aggregatedRefundedCents
	if manualRefundCents > 0 {
		createdAt := order.UpdatedAt
		if order.RefundedAt != nil {
			createdAt = *order.RefundedAt
		}
		status := "partial"
		if order.RefundAmountCents >= order.TotalPriceCents {
			status = "success"
		}
		item := OrderRefundItem{
			ID:          orderID*10 + 1,
			OrderID:     orderID,
			PaymentID:   0,
			AmountCents: manualRefundCents,
			Method:      "manual",
			Status:      status,
			Reason:      order.RefundReason,
			RefundedAt:  order.RefundedAt,
			CreatedAt:   createdAt,
			Note:        order.RefundReason,
		}
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func mapRefundRecordStatus(status model.RefundStatus) string {
	switch status {
	case model.RefundStatusProcessed:
		return "success"
	case model.RefundStatusPending:
		return "pending"
	case model.RefundStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func clampRefundedAmount(refunded, total int64) int64 {
	if refunded < 0 {
		return 0
	}
	if total >= 0 && refunded > total {
		return total
	}
	return refunded
}

// GetOrderReviews 返回订单相关的全部评价。
func (s *AdminService) GetOrderReviews(ctx context.Context, orderID uint64) ([]model.Review, error) {
	reviews := make([]model.Review, 0)
	page := 1
	orderIDCopy := orderID
	for {
		opts := repository.ReviewListOptions{
			Page:     page,
			PageSize: 200,
			OrderID:  &orderIDCopy,
		}
		items, pagination, err := s.ListReviews(ctx, opts)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, items...)
		if pagination == nil || !pagination.HasNext {
			break
		}
		page++
	}
	return reviews, nil
}

// GetOrderTimeline 汇总订单的状态流转与关键事件。
func (s *AdminService) GetOrderTimeline(ctx context.Context, orderID uint64) ([]OrderTimelineItem, error) {
	order, err := s.orders.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	logs, err := s.collectOperationLogs(ctx, string(model.OpEntityOrder), orderID)
	if err != nil {
		return nil, err
	}

	userCache := make(map[uint64]*model.User)
	items := make([]OrderTimelineItem, 0, len(logs))
	for _, logEntry := range logs {
		meta := map[string]any{}
		if len(logEntry.MetadataJSON) > 0 {
			_ = json.Unmarshal(logEntry.MetadataJSON, &meta)
		}
		item := OrderTimelineItem{
			ID:        logEntry.ID,
			OrderID:   orderID,
			EventType: mapTimelineEventType(logEntry.Action),
			Title:     mapTimelineTitle(logEntry.Action),
			Metadata:  meta,
			CreatedAt: logEntry.CreatedAt,
		}
		if note, ok := meta["note"].(string); ok && strings.TrimSpace(note) != "" {
			item.Description = strings.TrimSpace(note)
		} else if reason, ok := meta["reason"].(string); ok && strings.TrimSpace(reason) != "" {
			item.Description = strings.TrimSpace(reason)
		}
		if before, ok := meta["from_status"].(string); ok {
			item.StatusBefore = before
		}
		if after, ok := meta["status"].(string); ok {
			item.StatusAfter = after
		}
		if logEntry.ActorUserID != nil {
			if user := s.resolveUser(ctx, userCache, *logEntry.ActorUserID); user != nil {
				item.Operator = user.Name
				item.OperatorRole = string(user.Role)
				id := user.ID
				item.OperatorID = &id
			}
		}
		items = append(items, item)
	}

	// 追加支付关键事件
	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}
	for _, pay := range payments {
		if pay.PaidAt != nil {
			item := OrderTimelineItem{
				ID:        pay.ID*10 + 1,
				OrderID:   orderID,
				PaymentID: ptrUint64(pay.ID),
				EventType: "action",
				Title:     "支付确认",
				Metadata: map[string]any{
					"payment_status": pay.Status,
					"payment_method": pay.Method,
					"amount_cents":   pay.AmountCents,
				},
				CreatedAt: *pay.PaidAt,
			}
			items = append(items, item)
		}
		if pay.RefundedAt != nil {
			item := OrderTimelineItem{
				ID:          pay.ID*10 + 2,
				OrderID:     orderID,
				PaymentID:   ptrUint64(pay.ID),
				EventType:   "status_change",
				Title:       "支付退款",
				Description: strings.TrimSpace(order.RefundReason),
				Metadata: map[string]any{
					"payment_status": pay.Status,
					"payment_method": pay.Method,
					"amount_cents":   pay.AmountCents,
				},
				CreatedAt:   *pay.RefundedAt,
				StatusAfter: string(model.OrderStatusRefunded),
			}
			items = append(items, item)
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items, nil
}

// DeleteOrder 删除订单（软删）。
func (s *AdminService) DeleteOrder(ctx context.Context, id uint64) error {
	if err := s.orders.Delete(ctx, id); err != nil {
		return WrapError(err, "delete order")
	}
	s.invalidateCache(ctx, cacheKeyOrders)
	s.appendLogAsync(ctx, string(model.OpEntityOrder), id, string(model.OpActionDelete), nil)
	return nil
}

func (s *AdminService) hasPaidPayment(ctx context.Context, orderID uint64) (bool, error) {
	payments, err := s.listPaymentsByOrder(ctx, orderID)
	if err != nil {
		return false, err
	}
	for _, payment := range payments {
		if payment.Status == model.PaymentStatusPaid {
			return true, nil
		}
	}
	return false, nil
}

func isValidOrderStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusPending, model.OrderStatusConfirmed, model.OrderStatusInProgress,
		model.OrderStatusCompleted, model.OrderStatusCanceled, model.OrderStatusRefunded:
		return true
	default:
		return false
	}
}

func isAllowedOrderTransition(prev, next model.OrderStatus) bool {
	if prev == next {
		return true
	}
	switch prev {
	case model.OrderStatusPending:
		return next == model.OrderStatusConfirmed || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusConfirmed:
		return next == model.OrderStatusInProgress || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusInProgress:
		return next == model.OrderStatusCompleted || next == model.OrderStatusCanceled || next == model.OrderStatusRefunded
	case model.OrderStatusCompleted:
		return next == model.OrderStatusRefunded
	case model.OrderStatusCanceled, model.OrderStatusRefunded:
		return false
	default:
		return false
	}
}

func mapRefundStatus(status model.PaymentStatus) string {
	switch status {
	case model.PaymentStatusRefunded:
		return "success"
	case model.PaymentStatusPending:
		return "pending"
	case model.PaymentStatusFailed:
		return "failed"
	default:
		return strings.ToLower(string(status))
	}
}

func mapTimelineEventType(action string) string {
	switch action {
	case string(model.OpActionCreate):
		return "system"
	case string(model.OpActionAssignPlayer):
		return "action"
	case string(model.OpActionConfirm), string(model.OpActionStart), string(model.OpActionComplete),
		string(model.OpActionUpdateStatus), string(model.OpActionCancel), string(model.OpActionRefund):
		return "status_change"
	default:
		return "action"
	}
}

func mapTimelineTitle(action string) string {
	switch action {
	case string(model.OpActionCreate):
		return "订单创建"
	case string(model.OpActionAssignPlayer):
		return "指派陪玩师"
	case string(model.OpActionConfirm):
		return "订单确认"
	case string(model.OpActionStart):
		return "开始服务"
	case string(model.OpActionComplete):
		return "完成订单"
	case string(model.OpActionCancel):
		return "订单取消"
	case string(model.OpActionRefund):
		return "订单退款"
	case string(model.OpActionUpdateStatus):
		return "状态更新"
	default:
		return strings.ReplaceAll(action, "_", " ")
	}
}

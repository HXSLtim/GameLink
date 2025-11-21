package payment

import (
	"context"
	"time"
	"gamelink/internal/model"
	"encoding/json"
	"errors"
	"testing"
	"gamelink/internal/repository"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPaymentRepository struct {
	payments   map[uint64]*model.Payment
	createHook func(ctx context.Context, payment *model.Payment) error
	listHook   func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error)
	getHook    func(ctx context.Context, id uint64) (*model.Payment, error)
	updateHook func(ctx context.Context, payment *model.Payment) error
	deleteHook func(ctx context.Context, id uint64) error
}

func newMockPaymentRepository() *mockPaymentRepository {
	return &mockPaymentRepository{
		payments: make(map[uint64]*model.Payment),
	}
}

func (m *mockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	if m.createHook != nil {
		return m.createHook(ctx, payment)
	}
	payment.ID = uint64(len(m.payments) + 1)
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if m.listHook != nil {
		return m.listHook(ctx, opts)
	}
	var result []model.Payment
	for _, p := range m.payments {
		if opts.OrderID != nil && p.OrderID != *opts.OrderID {
			continue
		}
		result = append(result, *p)
	}
	return result, int64(len(result)), nil
}

func (m *mockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	if m.getHook != nil {
		return m.getHook(ctx, id)
	}
	if payment, ok := m.payments[id]; ok {
		return payment, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	if m.updateHook != nil {
		return m.updateHook(ctx, payment)
	}
	if _, ok := m.payments[payment.ID]; !ok {
		return repository.ErrNotFound
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	if m.deleteHook != nil {
		return m.deleteHook(ctx, id)
	}
	delete(m.payments, id)
	return nil
}

type mockOrderRepository struct {
	orders     map[uint64]*model.Order
	createHook func(ctx context.Context, order *model.Order) error
	listHook   func(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error)
	getHook    func(ctx context.Context, id uint64) (*model.Order, error)
	updateHook func(ctx context.Context, order *model.Order) error
	deleteHook func(ctx context.Context, id uint64) error
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	if m.createHook != nil {
		return m.createHook(ctx, order)
	}
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	if m.listHook != nil {
		return m.listHook(ctx, opts)
	}
	return []model.Order{}, 0, nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if m.getHook != nil {
		return m.getHook(ctx, id)
	}
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if m.updateHook != nil {
		return m.updateHook(ctx, order)
	}
	if _, ok := m.orders[order.ID]; !ok {
		return repository.ErrNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	if m.deleteHook != nil {
		return m.deleteHook(ctx, id)
	}
	delete(m.orders, id)
	return nil
}

func TestCreatePayment(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建测试订单
	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		UnitPriceCents:  10000,
		Quantity:        1,
		ItemID:          1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试创建支付
	resp, err := svc.CreatePayment(context.Background(), 1, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.PaymentID == 0 {
		t.Error("expected payment ID, got 0")
	}

	if resp.PayInfo == nil {
		t.Error("expected pay info, got nil")
	}

	// 验证支付信息包含必要字段
	if _, ok := resp.PayInfo["paymentId"]; !ok {
		t.Error("expected paymentId in pay info")
	}

	if method, ok := resp.PayInfo["method"]; !ok || method != "wechat" {
		t.Error("expected method 'wechat' in pay info")
	}
}

func TestGetPaymentStatus(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建测试订单和支付
	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		UnitPriceCents:  10000,
		Quantity:        1,
		ItemID:          1,
	}
	_ = orderRepo.Create(context.Background(), order)

	resp, _ := svc.CreatePayment(context.Background(), 1, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodAlipay,
	})

	// 测试查询支付状态
	status, err := svc.GetPaymentStatus(context.Background(), resp.PaymentID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status == nil {
		t.Fatal("expected status, got nil")
	}

	if status.PaymentID != resp.PaymentID {
		t.Errorf("expected payment ID %d, got %d", resp.PaymentID, status.PaymentID)
	}

	if status.OrderID != order.ID {
		t.Errorf("expected order ID %d, got %d", order.ID, status.OrderID)
	}

	// Mock 支付应该自动成功
	if status.Status != model.PaymentStatusPaid {
		t.Errorf("expected status paid, got %s", status.Status)
	}
}

func TestCancelPayment(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建待支付的支付记录
	payment := &model.Payment{
		OrderID:     1,
		UserID:      1,
		Status:      model.PaymentStatusPending,
		AmountCents: 10000,
	}
	_ = paymentRepo.Create(context.Background(), payment)

	// 测试取消支付
	err := svc.CancelPayment(context.Background(), 1, payment.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证支付状态已更新
	updated, _ := paymentRepo.Get(context.Background(), payment.ID)
	if updated.Status != model.PaymentStatusFailed {
		t.Errorf("expected status failed, got %s", updated.Status)
	}
}

func TestCreatePaymentInvalidOrderStatus(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建已支付的订单
	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 10000,
		UnitPriceCents:  10000,
		Quantity:        1,
		ItemID:          1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试为已确认订单创建支付（应该失败）
	_, err := svc.CreatePayment(context.Background(), 1, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})

	if err != ErrInvalidOrderStatus {
		t.Errorf("expected ErrInvalidOrderStatus, got %v", err)
	}
}

func TestCreatePaymentUnauthorized(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建其他用户的订单
	order := &model.Order{
		UserID:          2,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000,
		UnitPriceCents:  10000,
		Quantity:        1,
		ItemID:          1,
	}
	_ = orderRepo.Create(context.Background(), order)

	// 测试用户1为用户2的订单创建支付（应该失败）
	_, err := svc.CreatePayment(context.Background(), 1, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRefundPayment(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 先创建订单
	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 10000,
		Currency:        "CNY",
	}
	_ = orderRepo.Create(context.Background(), order)

	// 创建已支付的支付记录
	now := time.Now()
	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      1,
		Status:      model.PaymentStatusPaid,
		AmountCents: 10000,
		PaidAt:      &now,
	}
	_ = paymentRepo.Create(context.Background(), payment)

	// 测试退款
	err := svc.RefundPayment(context.Background(), payment.ID, "用户取消订单")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证退款状态
	refunded, err := paymentRepo.Get(context.Background(), payment.ID)
	if err != nil {
		t.Fatalf("failed to get refunded payment: %v", err)
	}

	if refunded.Status != model.PaymentStatusRefunded {
		t.Errorf("expected status refunded, got %s", refunded.Status)
	}

	if refunded.RefundedAt == nil {
		t.Error("expected refunded_at to be set")
	}
}

func TestHandlePaymentCallback(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建订单和待支付的支付记录
	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 10000,
		Currency:        "CNY",
	}
	_ = orderRepo.Create(context.Background(), order)

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      1,
		Status:      model.PaymentStatusPending,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
	}
	_ = paymentRepo.Create(context.Background(), payment)

	// 模拟支付回调数据
	callbackData := map[string]interface{}{
		"payment_id":   float64(payment.ID),
		"status":       "paid",
		"trade_no":     "MOCK123456",
		"amount_cents": int64(10000),
	}

	err := svc.HandlePaymentCallback(context.Background(), "alipay", callbackData)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 验证支付状态已更新
	updated, err := paymentRepo.Get(context.Background(), payment.ID)
	if err != nil {
		t.Fatalf("failed to get updated payment: %v", err)
	}

	if updated.Status != model.PaymentStatusPaid {
		t.Errorf("expected status paid, got %s", updated.Status)
	}

	// 验证订单状态已更新
	updatedOrder, err := orderRepo.Get(context.Background(), order.ID)
	if err != nil {
		t.Fatalf("failed to get updated order: %v", err)
	}

	if updatedOrder.Status != model.OrderStatusConfirmed {
		t.Errorf("expected order status confirmed, got %s", updatedOrder.Status)
	}
}

func TestCreatePayment_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("创建支付成功", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建待支付订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotZero(t, resp.PaymentID)
		assert.NotEmpty(t, resp.PayInfo)
	})

	t.Run("订单不存在应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		req := CreatePaymentRequest{
			OrderID: 999,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("无权限支付他人订单", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建用户1的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		// 用户2尝试支付
		resp, err := svc.CreatePayment(ctx, 2, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("非pending状态订单不能支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建已确认的订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrInvalidOrderStatus, err)
	})

	t.Run("订单已支付不能重复支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		// 创建已支付的支付记录
		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, ErrOrderAlreadyPaid, err)
	})

	t.Run("支付金额为0", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建0元订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 0,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		// 当前实现允许0元支付，但实际业务可能需要验证
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("支付金额为极大值", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建大额订单
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000000, // 100,000元
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func TestGetPaymentStatus_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("查询支付状态成功", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		resp, err := svc.GetPaymentStatus(ctx, payment.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, payment.ID, resp.PaymentID)
		assert.Equal(t, model.PaymentStatusPaid, resp.Status)
		assert.NotNil(t, resp.PaidAt)
	})

	t.Run("查询不存在的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		resp, err := svc.GetPaymentStatus(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("查询pending状态的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodAlipay,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		resp, err := svc.GetPaymentStatus(ctx, payment.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, model.PaymentStatusPending, resp.Status)
		assert.Nil(t, resp.PaidAt)
	})
}

func TestCancelPayment_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("取消pending状态的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		err := svc.CancelPayment(ctx, 1, payment.ID)

		assert.NoError(t, err)

		// 验证状态已更新
		updated, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusFailed, updated.Status)
	})

	t.Run("无权限取消他人支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		// 用户2尝试取消用户1的支付
		err := svc.CancelPayment(ctx, 2, payment.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})

	t.Run("不能取消已支付的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		err := svc.CancelPayment(ctx, 1, payment.ID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot cancel payment")
	})

	t.Run("取消不存在的支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		err := svc.CancelPayment(ctx, 1, 999)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestHandlePaymentCallback_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("成功处理支付回调", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建订单和支付
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		// 模拟支付回调
		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
			"trade_no":   "wx_trade_123",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.NoError(t, err)

		// 验证支付状态已更新
		updatedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusPaid, updatedPayment.Status)
		assert.NotNil(t, updatedPayment.PaidAt)
	})

	t.Run("重复回调应该幂等", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		// 创建已支付的支付
		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 10000,
		}
		orderRepo.Create(ctx, order)

		paidAt := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &paidAt,
		}
		paymentRepo.Create(ctx, payment)

		// 重复回调
		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		// 应该成功返回，不报错
		assert.NoError(t, err)
	})

	t.Run("缺少payment_id应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		callbackData := map[string]interface{}{
			"status": "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing payment_id")
	})

	t.Run("支付不存在应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		callbackData := map[string]interface{}{
			"payment_id": float64(999),
			"status":     "success",
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", callbackData)

		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("支付方式不匹配应该失败", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		payment := &model.Payment{
			OrderID:     1,
			UserID:      1,
			Method:      model.PaymentMethodWeChat,
			AmountCents: 10000,
			Status:      model.PaymentStatusPending,
		}
		paymentRepo.Create(ctx, payment)

		callbackData := map[string]interface{}{
			"payment_id": float64(payment.ID),
			"status":     "success",
		}

		// 使用alipay回调wechat支付
		err := svc.HandlePaymentCallback(ctx, "alipay", callbackData)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provider mismatch")
	})
}

func TestPaymentMethods(t *testing.T) {
	ctx := context.Background()

	t.Run("微信支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.PayInfo, "prepay_id")
		assert.Contains(t, resp.PayInfo, "code_url")
	})

	t.Run("支付宝支付", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          1,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 10000,
			Currency:        model.CurrencyCNY,
		}
		orderRepo.Create(ctx, order)

		req := CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		}

		resp, err := svc.CreatePayment(ctx, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.PayInfo, "trade_no")
		assert.Contains(t, resp.PayInfo, "qr_code")
	})
}

func TestCreatePayment_DetectsExistingPaidRecord(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusPending,
		TotalPriceCents: 8888,
		Currency:        "CNY",
	}
	_ = orderRepo.Create(context.Background(), order)

	// 预先写入一条已支付记录以触发 ErrOrderAlreadyPaid
	_ = paymentRepo.Create(context.Background(), &model.Payment{
		OrderID: order.ID,
		UserID:  1,
		Status:  model.PaymentStatusPaid,
	})

	_, err := svc.CreatePayment(context.Background(), 1, CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})

	if err != ErrOrderAlreadyPaid {
		t.Fatalf("expected ErrOrderAlreadyPaid, got %v", err)
	}
}

func TestCancelPayment_InvalidScenarios(t *testing.T) {
	ctx := context.Background()
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	t.Run("unauthorized user", func(t *testing.T) {
		payment := &model.Payment{
			OrderID:  1,
			UserID:   2,
			Status:   model.PaymentStatusPending,
			Currency: "CNY",
		}
		_ = paymentRepo.Create(ctx, payment)

		if err := svc.CancelPayment(ctx, 1, payment.ID); err == nil {
			t.Fatal("expected unauthorized error")
		}
	})

	t.Run("non pending status", func(t *testing.T) {
		payment := &model.Payment{
			OrderID: 2,
			UserID:  2,
			Status:  model.PaymentStatusPaid,
		}
		_ = paymentRepo.Create(ctx, payment)

		err := svc.CancelPayment(ctx, 2, payment.ID)
		if err == nil || !strings.Contains(err.Error(), "cannot cancel payment") {
			t.Fatalf("expected cannot cancel payment error, got %v", err)
		}
	})
}

func TestHandlePaymentCallback_Errors(t *testing.T) {
	ctx := context.Background()
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 12345,
	}
	_ = orderRepo.Create(ctx, order)

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      1,
		Status:      model.PaymentStatusPending,
		Method:      model.PaymentMethodWeChat,
		AmountCents: 12345,
	}
	_ = paymentRepo.Create(ctx, payment)

	t.Run("provider mismatch", func(t *testing.T) {
		err := svc.HandlePaymentCallback(ctx, "alipay", map[string]interface{}{
			"payment_id":   float64(payment.ID),
			"amount_cents": int64(payment.AmountCents),
		})
		if err == nil || !strings.Contains(err.Error(), "provider mismatch") {
			t.Fatalf("expected provider mismatch error, got %v", err)
		}
		stored, _ := paymentRepo.Get(ctx, payment.ID)
		if stored.Status != model.PaymentStatusPending {
			t.Fatalf("payment status should remain pending on error, got %s", stored.Status)
		}
	})

	t.Run("amount mismatch", func(t *testing.T) {
		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id":   float64(payment.ID),
			"amount_cents": int64(payment.AmountCents + 1),
		})
		if err == nil || !strings.Contains(err.Error(), "amount mismatch") {
			t.Fatalf("expected amount mismatch error, got %v", err)
		}
	})

	t.Run("missing payment id", func(t *testing.T) {
		if err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{}); err == nil {
			t.Fatal("expected error for missing payment_id")
		}
	})

	t.Run("duplicate callback returns nil", func(t *testing.T) {
		stored, _ := paymentRepo.Get(ctx, payment.ID)
		stored.Status = model.PaymentStatusPaid
		_ = paymentRepo.Update(ctx, stored)

		if err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": float64(payment.ID),
		}); err != nil {
			t.Fatalf("expected nil error when payment already processed, got %v", err)
		}
	})
}

func TestRefundPayment_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	order := &model.Order{
		UserID:          1,
		Status:          model.OrderStatusConfirmed,
		TotalPriceCents: 5000,
	}
	_ = orderRepo.Create(ctx, order)

	payment := &model.Payment{
		OrderID:     order.ID,
		UserID:      1,
		Status:      model.PaymentStatusPending,
		AmountCents: 5000,
	}
	_ = paymentRepo.Create(ctx, payment)

	err := svc.RefundPayment(ctx, payment.ID, "not paid yet")
	if err == nil || !strings.Contains(err.Error(), "payment status must be paid") {
		t.Fatalf("expected payment status error, got %v", err)
	}
}

func TestPaymentService_CreatePayment_ErrorBranches(t *testing.T) {
	ctx := context.Background()

	t.Run("list error bypasses duplicate check", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		paymentRepo.listHook = func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return nil, 0, errors.New("db unavailable")
		}
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          7,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 3210,
			Currency:        model.CurrencyCNY,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		resp, err := svc.CreatePayment(ctx, order.UserID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("create error surfaces", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		paymentRepo.createHook = func(ctx context.Context, payment *model.Payment) error {
			return errors.New("insert failed")
		}
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          9,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 555,
			Currency:        model.CurrencyCNY,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		resp, err := svc.CreatePayment(ctx, order.UserID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.True(t, strings.Contains(err.Error(), "insert failed"))
	})

	t.Run("existing pending payment allows retry", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          10,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 4321,
			Currency:        model.CurrencyCNY,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		// Existing pending payment should not block creation.
		existing := &model.Payment{
			OrderID: order.ID,
			UserID:  order.UserID,
			Status:  model.PaymentStatusPending,
		}
		require.NoError(t, paymentRepo.Create(ctx, existing))

		resp, err := svc.CreatePayment(ctx, order.UserID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodWeChat,
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Greater(t, resp.PaymentID, existing.ID)
	})

	t.Run("mock payment success failure surfaces", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		paymentRepo.getHook = func(ctx context.Context, id uint64) (*model.Payment, error) {
			return nil, errors.New("mock success failed")
		}
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          11,
			Status:          model.OrderStatusPending,
			TotalPriceCents: 9876,
			Currency:        model.CurrencyCNY,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		resp, err := svc.CreatePayment(ctx, order.UserID, CreatePaymentRequest{
			OrderID: order.ID,
			Method:  model.PaymentMethodAlipay,
		})

		require.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "mock success failed")
	})
}

func TestPaymentService_mockPaymentSuccess_ErrorScenarios(t *testing.T) {
	ctx := context.Background()

	t.Run("payment get error", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		paymentRepo.getHook = func(ctx context.Context, id uint64) (*model.Payment, error) {
			return nil, repository.ErrNotFound
		}
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		err := svc.mockPaymentSuccess(ctx, 1, &model.Order{})
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("payment update error", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{Status: model.OrderStatusPending}
		require.NoError(t, orderRepo.Create(ctx, order))

		payment := &model.Payment{OrderID: order.ID, Status: model.PaymentStatusPending}
		require.NoError(t, paymentRepo.Create(ctx, payment))

		paymentRepo.updateHook = func(ctx context.Context, payment *model.Payment) error {
			return errors.New("update payment failed")
		}

		err := svc.mockPaymentSuccess(ctx, payment.ID, order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update payment failed")
	})

	t.Run("order update error", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{Status: model.OrderStatusPending}
		require.NoError(t, orderRepo.Create(ctx, order))

		payment := &model.Payment{OrderID: order.ID, Status: model.PaymentStatusPending}
		require.NoError(t, paymentRepo.Create(ctx, payment))

		orderRepo.updateHook = func(ctx context.Context, order *model.Order) error {
			return errors.New("update order failed")
		}

		err := svc.mockPaymentSuccess(ctx, payment.ID, order)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update order failed")
	})
}

func TestPaymentService_HandlePaymentCallback_CompletesBranches(t *testing.T) {
	ctx := context.Background()

	baseOrder := func(status model.OrderStatus) (*mockPaymentRepository, *mockOrderRepository, *PaymentService, *model.Order, *model.Payment) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          42,
			Status:          status,
			TotalPriceCents: 2000,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      order.UserID,
			Method:      model.PaymentMethodWeChat,
			AmountCents: order.TotalPriceCents,
			Status:      model.PaymentStatusPending,
		}
		require.NoError(t, paymentRepo.Create(ctx, payment))

		return paymentRepo, orderRepo, svc, order, payment
	}

	t.Run("success with generated trade number", func(t *testing.T) {
		paymentRepo, orderRepo, svc, order, payment := baseOrder(model.OrderStatusPending)

		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": payment.ID,
		})
		require.NoError(t, err)

		storedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusPaid, storedPayment.Status)
		assert.NotEmpty(t, storedPayment.ProviderTradeNo)
		assert.True(t, strings.HasPrefix(storedPayment.ProviderTradeNo, "wechat_"))

		storedOrder, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusConfirmed, storedOrder.Status)
	})

	t.Run("success with provided trade number and amount", func(t *testing.T) {
		paymentRepo, _, svc, order, payment := baseOrder(model.OrderStatusPending)

		amount := int64(order.TotalPriceCents)
		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id":   payment.ID,
			"amount_cents": amount,
			"trade_no":     "custom_trade",
		})
		require.NoError(t, err)

		storedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, "custom_trade", storedPayment.ProviderTradeNo)
		assert.Equal(t, amount, storedPayment.AmountCents)
	})

	t.Run("payment retrieval failure", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": uint64(999),
		})
		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("payment update failure bubbles up", func(t *testing.T) {
		paymentRepo, _, svc, _, payment := baseOrder(model.OrderStatusPending)
		paymentRepo.updateHook = func(ctx context.Context, payment *model.Payment) error {
			return errors.New("persist payment error")
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": payment.ID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "persist payment error")
	})

	t.Run("order update failure bubbles up", func(t *testing.T) {
		_, orderRepo, svc, _, payment := baseOrder(model.OrderStatusPending)
		orderRepo.updateHook = func(ctx context.Context, order *model.Order) error {
			return errors.New("persist order error")
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": payment.ID,
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "persist order error")
	})

	t.Run("order retrieval failure", func(t *testing.T) {
		_, orderRepo, svc, _, payment := baseOrder(model.OrderStatusPending)
		orderRepo.getHook = func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		}

		err := svc.HandlePaymentCallback(ctx, "wechat", map[string]interface{}{
			"payment_id": payment.ID,
		})
		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func TestPaymentService_RefundPayment_FullCoverage(t *testing.T) {
	ctx := context.Background()

	newSetup := func(method model.PaymentMethod) (*mockPaymentRepository, *mockOrderRepository, *PaymentService, *model.Payment, *model.Order) {
		paymentRepo := newMockPaymentRepository()
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		order := &model.Order{
			UserID:          88,
			Status:          model.OrderStatusConfirmed,
			TotalPriceCents: 900,
		}
		require.NoError(t, orderRepo.Create(ctx, order))

		now := time.Now()
		payment := &model.Payment{
			OrderID:     order.ID,
			UserID:      order.UserID,
			Method:      method,
			AmountCents: order.TotalPriceCents,
			Status:      model.PaymentStatusPaid,
			PaidAt:      &now,
		}
		require.NoError(t, paymentRepo.Create(ctx, payment))

		return paymentRepo, orderRepo, svc, payment, order
	}

	t.Run("wechat provider", func(t *testing.T) {
		paymentRepo, orderRepo, svc, payment, order := newSetup(model.PaymentMethodWeChat)

		err := svc.RefundPayment(ctx, payment.ID, "duplicate")
		require.NoError(t, err)

		storedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.Equal(t, model.PaymentStatusRefunded, storedPayment.Status)
		assert.True(t, strings.HasPrefix(storedPayment.ProviderTradeNo, "wx_refund_"))
		assert.NotNil(t, storedPayment.RefundedAt)

		storedOrder, _ := orderRepo.Get(ctx, order.ID)
		assert.Equal(t, model.OrderStatusRefunded, storedOrder.Status)
		assert.Equal(t, "duplicate", storedOrder.RefundReason)
	})

	t.Run("payment retrieval failure", func(t *testing.T) {
		paymentRepo := newMockPaymentRepository()
		paymentRepo.getHook = func(ctx context.Context, id uint64) (*model.Payment, error) {
			return nil, repository.ErrNotFound
		}
		orderRepo := newMockOrderRepository()
		svc := NewPaymentService(paymentRepo, orderRepo)

		err := svc.RefundPayment(ctx, 999, "missing")
		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("alipay provider", func(t *testing.T) {
		paymentRepo, _, svc, payment, _ := newSetup(model.PaymentMethodAlipay)

		err := svc.RefundPayment(ctx, payment.ID, "quality")
		require.NoError(t, err)

		storedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.True(t, strings.HasPrefix(storedPayment.ProviderTradeNo, "ali_refund_"))
	})

	t.Run("generic provider fallback", func(t *testing.T) {
		paymentRepo, _, svc, payment, _ := newSetup(model.PaymentMethod("bank_transfer"))

		err := svc.RefundPayment(ctx, payment.ID, "manual")
		require.NoError(t, err)

		storedPayment, _ := paymentRepo.Get(ctx, payment.ID)
		assert.True(t, strings.HasPrefix(storedPayment.ProviderTradeNo, "refund_"))
	})

	t.Run("provider refund failure", func(t *testing.T) {
		_, _, svc, payment, _ := newSetup(model.PaymentMethodWeChat)
		svc.providers[model.PaymentMethodWeChat] = failingProvider{}

		err := svc.RefundPayment(ctx, payment.ID, "issue")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "provider failure")
	})

	t.Run("payment update failure", func(t *testing.T) {
		paymentRepo, _, svc, payment, _ := newSetup(model.PaymentMethodWeChat)
		paymentRepo.updateHook = func(ctx context.Context, payment *model.Payment) error {
			return errors.New("update payment error")
		}

		err := svc.RefundPayment(ctx, payment.ID, "fail")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update payment error")
	})

	t.Run("order update failure", func(t *testing.T) {
		_, orderRepo, svc, payment, _ := newSetup(model.PaymentMethodWeChat)
		orderRepo.updateHook = func(ctx context.Context, order *model.Order) error {
			return errors.New("update order error")
		}

		err := svc.RefundPayment(ctx, payment.ID, "fail")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "update order error")
	})

	t.Run("order retrieval failure", func(t *testing.T) {
		_, orderRepo, svc, payment, _ := newSetup(model.PaymentMethodWeChat)
		orderRepo.getHook = func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		}

		err := svc.RefundPayment(ctx, payment.ID, "missing order")
		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})
}

func (failingProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	return "", nil, time.Time{}, errors.New("provider failure")
}

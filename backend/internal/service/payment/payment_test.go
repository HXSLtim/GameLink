package payment

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Mock repositories
type mockPaymentRepository struct {
	payments map[uint64]*model.Payment
}

func newMockPaymentRepository() *mockPaymentRepository {
	return &mockPaymentRepository{
		payments: make(map[uint64]*model.Payment),
	}
}

func (m *mockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	payment.ID = uint64(len(m.payments) + 1)
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
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
	if payment, ok := m.payments[id]; ok {
		return payment, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	if _, ok := m.payments[payment.ID]; !ok {
		return repository.ErrNotFound
	}
	m.payments[payment.ID] = payment
	return nil
}

func (m *mockPaymentRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.payments, id)
	return nil
}

type mockOrderRepository struct {
	orders map[uint64]*model.Order
}

func newMockOrderRepository() *mockOrderRepository {
	return &mockOrderRepository{
		orders: make(map[uint64]*model.Order),
	}
}

func (m *mockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	order.ID = uint64(len(m.orders) + 1)
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	return []model.Order{}, 0, nil
}

func (m *mockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if order, ok := m.orders[id]; ok {
		return order, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if _, ok := m.orders[order.ID]; !ok {
		return repository.ErrNotFound
	}
	m.orders[order.ID] = order
	return nil
}

func (m *mockOrderRepository) Delete(ctx context.Context, id uint64) error {
	delete(m.orders, id)
	return nil
}

func TestCreatePayment(t *testing.T) {
	paymentRepo := newMockPaymentRepository()
	orderRepo := newMockOrderRepository()
	svc := NewPaymentService(paymentRepo, orderRepo)

	// 创建测试订单
	order := &model.Order{
		UserID:     1,
		Status:     model.OrderStatusPending,
		PriceCents: 10000,
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
		UserID:     1,
		Status:     model.OrderStatusPending,
		PriceCents: 10000,
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
		UserID:     1,
		Status:     model.OrderStatusConfirmed,
		PriceCents: 10000,
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
		UserID:     2,
		Status:     model.OrderStatusPending,
		PriceCents: 10000,
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

	// 创建已支付的支付记录
	now := time.Now()
	payment := &model.Payment{
		OrderID:     1,
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
	refunded, _ := paymentRepo.Get(context.Background(), payment.ID)
	if refunded.Status != model.PaymentStatusRefunded {
		t.Errorf("expected status refunded, got %s", refunded.Status)
	}

	if refunded.RefundedAt == nil {
		t.Error("expected refunded_at to be set")
	}
}

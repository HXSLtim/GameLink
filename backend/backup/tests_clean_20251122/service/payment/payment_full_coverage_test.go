package payment

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

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

type failingProvider struct{}

func (failingProvider) Refund(ctx context.Context, p *model.Payment, reason string) (string, json.RawMessage, time.Time, error) {
	return "", nil, time.Time{}, errors.New("provider failure")
}

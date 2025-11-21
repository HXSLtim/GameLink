package payment

import (
	"context"
	"strings"
	"testing"

	"gamelink/internal/model"
)

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

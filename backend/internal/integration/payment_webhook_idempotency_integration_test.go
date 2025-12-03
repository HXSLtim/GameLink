package integration

import (
	"testing"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/pkg/testutil"
)

// 支付回调幂等 & 提供商校验：首次回调成功更新支付/订单，再次回调无副作用，提供商不匹配返回错误
func TestPaymentCallbackIdempotencyAndProviderCheck(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePaymentModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	// 准备订单（待支付）
	order := &model.Order{
		UserID:          1,
		Title:           "Callback order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 5000,
		Currency:        model.CurrencyCNY,
	}
	if err := orderRepo.Create(ctx(), order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// 准备支付（待支付）
	p := &model.Payment{
		OrderID:     order.ID,
		UserID:      order.UserID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: order.TotalPriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}
	if err := paymentRepo.Create(ctx(), p); err != nil {
		t.Fatalf("seed payment: %v", err)
	}

	cb := map[string]interface{}{
		"payment_id":   p.ID,
		"amount_cents": p.AmountCents,
	}

	// 首次回调：应变更为已支付/已确认
	if err := svc.HandlePaymentCallback(ctx(), "wechat", cb); err != nil {
		t.Fatalf("callback failed: %v", err)
	}
	updatedPayment, _ := paymentRepo.Get(ctx(), p.ID)
	if updatedPayment.Status != model.PaymentStatusPaid {
		t.Fatalf("expected payment paid, got %s", updatedPayment.Status)
	}
	updatedOrder, _ := orderRepo.Get(ctx(), order.ID)
	if updatedOrder.Status != model.OrderStatusConfirmed {
		t.Fatalf("expected order confirmed, got %s", updatedOrder.Status)
	}

	// 再次回调：幂等，无错误，无额外变化
	if err := svc.HandlePaymentCallback(ctx(), "wechat", cb); err != nil {
		t.Fatalf("second callback should be idempotent, got %v", err)
	}

	// 提供商不匹配：报错且保持待支付
	other := &model.Payment{
		OrderID:     order.ID,
		UserID:      order.UserID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: order.TotalPriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}
	_ = paymentRepo.Create(ctx(), other)
	err := svc.HandlePaymentCallback(ctx(), "wechat", map[string]interface{}{
		"payment_id":   other.ID,
		"amount_cents": other.AmountCents,
	})
	if err == nil {
		t.Fatalf("expected provider mismatch error, got nil")
	}
	otherPayment, _ := paymentRepo.Get(ctx(), other.ID)
	if otherPayment.Status != model.PaymentStatusPending {
		t.Fatalf("expected payment remain pending on provider mismatch, got %s", otherPayment.Status)
	}

	// 金额不匹配：报错并保持 pending
	amountMismatch := &model.Payment{
		OrderID:     order.ID,
		UserID:      order.UserID,
		Method:      model.PaymentMethodWeChat,
		AmountCents: order.TotalPriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}
	_ = paymentRepo.Create(ctx(), amountMismatch)
	err = svc.HandlePaymentCallback(ctx(), "wechat", map[string]interface{}{
		"payment_id":   amountMismatch.ID,
		"amount_cents": amountMismatch.AmountCents + 1, // wrong amount
	})
	if err == nil {
		t.Fatalf("expected amount mismatch error, got nil")
	}
	mismatchPayment, _ := paymentRepo.Get(ctx(), amountMismatch.ID)
	if mismatchPayment.Status != model.PaymentStatusPending {
		t.Fatalf("expected payment remain pending on amount mismatch, got %s", mismatchPayment.Status)
	}
}

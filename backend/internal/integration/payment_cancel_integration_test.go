package integration

import (
	"testing"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/pkg/testutil"
)

// 待支付取消后可重新发起支付
func TestCancelPendingPaymentAndRepay(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePaymentModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	// seed order pending
	order := &model.Order{
		UserID:          10,
		Title:           "Cancelable order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 4000,
		Currency:        model.CurrencyCNY,
	}
	if err := orderRepo.Create(ctx(), order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// seed pending payment
	p := &model.Payment{
		OrderID:     order.ID,
		UserID:      order.UserID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: order.TotalPriceCents,
		Currency:    order.Currency,
		Status:      model.PaymentStatusPending,
	}
	if err := paymentRepo.Create(ctx(), p); err != nil {
		t.Fatalf("seed payment: %v", err)
	}

	// cancel pending
	if err := svc.CancelPayment(ctx(), order.UserID, p.ID); err != nil {
		t.Fatalf("cancel payment: %v", err)
	}
	canceled, _ := paymentRepo.Get(ctx(), p.ID)
	if canceled.Status != model.PaymentStatusFailed {
		t.Fatalf("expected failed status after cancel, got %s", canceled.Status)
	}

	// create new payment and ensure paid + order confirmed
	resp, err := svc.CreatePayment(ctx(), order.UserID, paymentsvc.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})
	if err != nil {
		t.Fatalf("create payment after cancel: %v", err)
	}
	newPayment, _ := paymentRepo.Get(ctx(), resp.PaymentID)
	if newPayment.Status != model.PaymentStatusPaid {
		t.Fatalf("expected new payment paid, got %s", newPayment.Status)
	}
	updatedOrder, _ := orderRepo.Get(ctx(), order.ID)
	if updatedOrder.Status != model.OrderStatusConfirmed {
		t.Fatalf("expected order confirmed, got %s", updatedOrder.Status)
	}
}

// 已支付取消应失败且状态保持 paid
func TestCancelPaidPaymentShouldFail(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePaymentModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	order := &model.Order{
		UserID:          11,
		Title:           "Paid order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 3000,
		Currency:        model.CurrencyCNY,
	}
	_ = orderRepo.Create(ctx(), order)

	// create payment -> auto paid
	resp, err := svc.CreatePayment(ctx(), order.UserID, paymentsvc.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWeChat,
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}

	// try cancel should error
	if err := svc.CancelPayment(ctx(), order.UserID, resp.PaymentID); err == nil {
		t.Fatalf("expected cancel error for paid payment")
	}
	p, _ := paymentRepo.Get(ctx(), resp.PaymentID)
	if p.Status != model.PaymentStatusPaid {
		t.Fatalf("payment status changed after failed cancel: %s", p.Status)
	}
}

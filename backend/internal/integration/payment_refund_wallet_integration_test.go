package integration

import (
	"testing"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	userrepo "gamelink/internal/repository/user"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/pkg/testutil"
)

// 退款回滚到钱包余额
func TestRefundCreditsWalletBalance(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	walletRepo := userrepo.NewWalletRepository(db)

	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)
	svc.SetWalletRepository(walletRepo)

	order := &model.Order{
		UserID:          21,
		Title:           "Refund to wallet",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 7000,
		Currency:        model.CurrencyCNY,
	}
	if err := orderRepo.Create(ctx(), order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// 支付 -> 自动 paid + confirmed
	payResp, err := svc.CreatePayment(ctx(), order.UserID, paymentsvc.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodAlipay,
	})
	if err != nil {
		t.Fatalf("create payment: %v", err)
	}

	// 退款
	if err := svc.RefundPayment(ctx(), payResp.PaymentID, "user request"); err != nil {
		t.Fatalf("refund payment: %v", err)
	}

	// 钱包余额应增加 7000
	w, err := walletRepo.GetByUserID(ctx(), order.UserID)
	if err != nil {
		t.Fatalf("get wallet: %v", err)
	}
	if w.BalanceCents != 7000 {
		t.Fatalf("expected wallet balance 7000, got %d", w.BalanceCents)
	}
}

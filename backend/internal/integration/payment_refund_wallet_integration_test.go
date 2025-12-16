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

// 钱包支付退款回滚到钱包余额
func TestRefundCreditsWalletBalance(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	walletRepo := userrepo.NewWalletRepository(db)

	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)
	svc.SetWalletRepository(walletRepo)

	// 创建用户钱包（有足够余额）
	wallet := &model.Wallet{
		UserID:       21,
		BalanceCents: 10000, // 100元
	}
	if err := db.Create(wallet).Error; err != nil {
		t.Fatalf("seed wallet: %v", err)
	}

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

	// 钱包支付 -> 自动 paid + confirmed
	payResp, err := svc.CreatePayment(ctx(), order.UserID, paymentsvc.CreatePaymentRequest{
		OrderID: order.ID,
		Method:  model.PaymentMethodWallet,
	})
	if err != nil {
		t.Fatalf("create wallet payment: %v", err)
	}

	// 验证钱包余额减少
	w, err := walletRepo.GetByUserID(ctx(), order.UserID)
	if err != nil {
		t.Fatalf("get wallet after payment: %v", err)
	}
	if w.BalanceCents != 3000 {
		t.Fatalf("expected wallet balance 3000 after payment, got %d", w.BalanceCents)
	}

	// 退款
	if err := svc.RefundPayment(ctx(), payResp.PaymentID, "user request"); err != nil {
		t.Fatalf("refund payment: %v", err)
	}

	// 钱包余额应恢复到 10000
	w, err = walletRepo.GetByUserID(ctx(), order.UserID)
	if err != nil {
		t.Fatalf("get wallet after refund: %v", err)
	}
	if w.BalanceCents != 10000 {
		t.Fatalf("expected wallet balance 10000 after refund, got %d", w.BalanceCents)
	}
}

// 组合支付退款测试
func TestCombinedPaymentRefund(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateWalletModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	walletRepo := userrepo.NewWalletRepository(db)

	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)
	svc.SetWalletRepository(walletRepo)

	// 创建用户钱包
	wallet := &model.Wallet{
		UserID:       22,
		BalanceCents: 5000, // 50元
	}
	if err := db.Create(wallet).Error; err != nil {
		t.Fatalf("seed wallet: %v", err)
	}

	order := &model.Order{
		UserID:          22,
		Title:           "Combined payment order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 10000, // 100元
		Currency:        model.CurrencyCNY,
	}
	if err := orderRepo.Create(ctx(), order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// 组合支付：3000钱包 + 7000微信
	payResp, err := svc.CreatePayment(ctx(), order.UserID, paymentsvc.CreatePaymentRequest{
		OrderID:           order.ID,
		Method:            model.PaymentMethodCombined,
		WalletAmountCents: 3000,
		ThirdPartyMethod:  model.PaymentMethodWeChat,
	})
	if err != nil {
		t.Fatalf("create combined payment: %v", err)
	}

	// 验证钱包余额减少
	w, err := walletRepo.GetByUserID(ctx(), order.UserID)
	if err != nil {
		t.Fatalf("get wallet after payment: %v", err)
	}
	if w.BalanceCents != 2000 {
		t.Fatalf("expected wallet balance 2000 after payment, got %d", w.BalanceCents)
	}

	// 退款
	if err := svc.RefundPayment(ctx(), payResp.PaymentID, "user request"); err != nil {
		t.Fatalf("refund payment: %v", err)
	}

	// 钱包余额应恢复到 5000（只退回钱包部分）
	w, err = walletRepo.GetByUserID(ctx(), order.UserID)
	if err != nil {
		t.Fatalf("get wallet after refund: %v", err)
	}
	if w.BalanceCents != 5000 {
		t.Fatalf("expected wallet balance 5000 after refund, got %d", w.BalanceCents)
	}
}

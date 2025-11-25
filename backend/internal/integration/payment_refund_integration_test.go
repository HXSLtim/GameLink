package integration

import (
	"testing"

	"gamelink/internal/model"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/internal/testutil"
)

// 支付退款：已支付订单退款成功，重复退款返回错误且状态不变
func TestPaymentRefundFlow(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePaymentModels(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	svc := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	// 创建订单（pending）
	order := &model.Order{
		UserID:          1,
		Title:           "Refund order",
		Status:          model.OrderStatusPending,
		TotalPriceCents: 2000,
		Currency:        model.CurrencyCNY,
	}
	if err := orderRepo.Create(ctx(), order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// 创建支付（自动置为 paid + 订单 confirmed）
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

	refundedPayment, _ := paymentRepo.Get(ctx(), payResp.PaymentID)
	if refundedPayment.Status != model.PaymentStatusRefunded {
		t.Fatalf("expected refunded payment, got %s", refundedPayment.Status)
	}
	refundedOrder, _ := orderRepo.Get(ctx(), order.ID)
	if refundedOrder.Status != model.OrderStatusRefunded {
		t.Fatalf("expected refunded order, got %s", refundedOrder.Status)
	}

	// 重复退款应报错且状态保持 refunded
	if err := svc.RefundPayment(ctx(), payResp.PaymentID, "again"); err == nil {
		t.Fatalf("expected error on second refund, got nil")
	}
	refundedPayment2, _ := paymentRepo.Get(ctx(), payResp.PaymentID)
	if refundedPayment2.Status != model.PaymentStatusRefunded {
		t.Fatalf("payment status changed on second refund: %s", refundedPayment2.Status)
	}
}

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/pkg/cache"
	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/user"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/serviceitem"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/testutil"
)

// 管理端支付退款：管理员对已支付记录执行退款
func TestAdminPaymentRefund(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateOrderModels(t, db)

	ctx := context.Background()
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)

	adminSvc := adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, serviceItemRepo, memCache)
	paymentHandler := adminhandler.NewPaymentHandler(adminSvc)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.POST("/payments/:id/refund", paymentHandler.RefundPayment)

	// seed paid payment
	orderID := createOrderForAdmin(t, ctx, orderRepo, seed, model.OrderStatusCompleted, ptrTime(time.Now().Add(-time.Hour)))
	paymentModel := &model.Payment{
		OrderID:     orderID,
		UserID:      seed.userID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: 10000,
		Currency:    model.CurrencyCNY,
		Status:      model.PaymentStatusPaid,
		PaidAt:      ptrTime(time.Now().Add(-50 * time.Minute)),
	}
	if err := paymentRepo.Create(ctx, paymentModel); err != nil {
		t.Fatalf("seed payment: %v", err)
	}

	refundAt := time.Now().UTC().Format(time.RFC3339)
	payload := map[string]interface{}{
		"refunded_at":       refundAt,
		"provider_trade_no": "trade-no-1",
		"provider_raw":      `{"state":"refunded"}`,
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/payments/"+uintToStr(paymentModel.ID)+"/refund", payload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("refund payment status=%d body=%s", resp.Code, resp.Body.String())
	}

	var parsed apiResp[model.Payment]
	if err := json.Unmarshal(resp.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("parse refund resp: %v", err)
	}
	if parsed.Data.Status != model.PaymentStatusRefunded {
		t.Fatalf("expected refunded status, got %s", parsed.Data.Status)
	}
	if parsed.Data.RefundedAt == nil {
		t.Fatalf("expected refundedAt set")
	}
}

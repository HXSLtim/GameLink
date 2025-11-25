package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/role"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/testutil"
	"gamelink/internal/cache"
)

// 管理端订单取消/退款链路：管理员取消待付款订单，管理员为已完成订单发起退款
func TestAdminOrderCancelAndRefund(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateOrderModels(t, db)

	ctx := context.Background()
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	roleRepo := role.NewRoleRepository(db)
	memCache := cache.NewMemory()

	adminSvc := adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, memCache)
	orderHandler := adminhandler.NewOrderHandler(adminSvc)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.POST("/orders/:id/cancel", orderHandler.CancelOrder)
	adminGroup.POST("/orders/:id/refund", orderHandler.RefundOrder)

	// 待付款订单 -> 管理员取消
	pendingOrderID := createOrderForAdmin(t, ctx, orderRepo, seed, model.OrderStatusPending, nil)
	cancelResp := doJSON(router, http.MethodPost, "/api/v1/admin/orders/"+uintToStr(pendingOrderID)+"/cancel", map[string]string{
		"reason": "admin cancel",
	}, "")
	if cancelResp.Code != http.StatusOK {
		t.Fatalf("admin cancel order status=%d body=%s", cancelResp.Code, cancelResp.Body.String())
	}
	var cancelParsed apiResp[model.Order]
	if err := json.Unmarshal(cancelResp.Body.Bytes(), &cancelParsed); err != nil {
		t.Fatalf("parse cancel resp: %v", err)
	}
	if cancelParsed.Data.Status != model.OrderStatusCanceled {
		t.Fatalf("expected canceled status, got %s", cancelParsed.Data.Status)
	}

	// 已完成订单 -> 管理员部分退款
	completedAt := time.Now().Add(-time.Hour).UTC()
	refundOrderID := createOrderForAdmin(t, ctx, orderRepo, seed, model.OrderStatusCompleted, &completedAt)
	refundAmount := int64(6000)
	refundResp := doJSON(router, http.MethodPost, "/api/v1/admin/orders/"+uintToStr(refundOrderID)+"/refund", map[string]interface{}{
		"reason":       "partial refund",
		"amount_cents": refundAmount,
		"note":         "test refund",
	}, "")
	if refundResp.Code != http.StatusOK {
		t.Fatalf("admin refund order status=%d body=%s", refundResp.Code, refundResp.Body.String())
	}
	var refundParsed apiResp[model.Order]
	if err := json.Unmarshal(refundResp.Body.Bytes(), &refundParsed); err != nil {
		t.Fatalf("parse refund resp: %v", err)
	}
	if refundParsed.Data.Status != model.OrderStatusRefunded {
		t.Fatalf("expected refunded status, got %s", refundParsed.Data.Status)
	}
	if refundParsed.Data.RefundAmountCents != refundAmount {
		t.Fatalf("expected refund amount %d got %d", refundAmount, refundParsed.Data.RefundAmountCents)
	}
	if refundParsed.Data.RefundReason == "" {
		t.Fatalf("expected refund reason set")
	}
	if refundParsed.Data.RefundedAt == nil {
		t.Fatalf("expected refundedAt set")
	}
}

func createOrderForAdmin(t *testing.T, ctx context.Context, orderRepo repoiface.OrderRepository, seed orderSeed, status model.OrderStatus, completedAt *time.Time) uint64 {
	t.Helper()
	order := &model.Order{
		UserID:          seed.userID,
		ItemID:          seed.gameID,
		Status:          status,
		Title:           "Admin Order",
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
	}
	order.SetPlayerID(seed.playerID)
	order.SetGameID(seed.gameID)
	if completedAt != nil {
		order.CompletedAt = completedAt
	}
	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("seed admin order: %v", err)
	}
	return order.ID
}

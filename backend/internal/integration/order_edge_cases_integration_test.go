package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"context"

	"github.com/gin-gonic/gin"

	playerhandler "gamelink/internal/handler/player"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/internal/testutil"
)

// 支付/订单边缘：重复支付幂等 & 退款后再次支付被拒
func TestOrderPaymentEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateOrderModels(t, db)
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)

	orderService := ordersvc.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	paymentService := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	playerAuth := fakeAuthMiddleware(seed.playerUserID)
	userhandler.RegisterOrderRoutes(api, orderService, userAuth)
	userhandler.RegisterPaymentRoutes(api, paymentService, userAuth)
	playerhandler.RegisterOrderRoutes(api, orderService, playerAuth)

	// 创建订单
	scheduledStart := time.Now().Add(30 * time.Minute).UTC()
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Edge Order",
		"scheduledStart": scheduledStart.Format(time.RFC3339),
		"durationHours":  1.0,
	}
	w := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
	if w.Code != http.StatusOK {
		t.Fatalf("create order status=%d body=%s", w.Code, w.Body.String())
	}
	var orderResp apiResp[ordersvc.CreateOrderResponse]
	_ = json.Unmarshal(w.Body.Bytes(), &orderResp)
	orderID := orderResp.Data.OrderID

	// 第一次支付成功
	payPayload := map[string]interface{}{"orderId": orderID, "method": "alipay"}
	payResp1 := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp1.Code != http.StatusOK {
		t.Fatalf("first payment status=%d body=%s", payResp1.Code, payResp1.Body.String())
	}

	// 再次支付同一订单 -> 应拒绝（非 pending 状态）
	payResp2 := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp2.Code == http.StatusOK {
		t.Fatalf("expected duplicate payment rejected, got 200")
	}

	// 玩家完成订单 -> 管理端退款后，再次支付应拒绝（订单已退款）
	_ = doJSON(router, http.MethodPost, "/api/v1/player/orders/"+uintToStr(orderID)+"/accept", nil, "")
	_ = doJSON(router, http.MethodPut, "/api/v1/player/orders/"+uintToStr(orderID)+"/complete", nil, "")
	// 模拟退款：直接更新支付状态
	pList, _, _ := paymentRepo.List(context.Background(), repository.PaymentListOptions{OrderID: &orderID})
	if len(pList) > 0 {
		p := pList[0]
		p.Status = model.PaymentStatusRefunded
		_ = paymentRepo.Update(context.Background(), &p)
	}
	payResp3 := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp3.Code == http.StatusOK {
		t.Fatalf("expected payment rejected after refund, got 200")
	}
}

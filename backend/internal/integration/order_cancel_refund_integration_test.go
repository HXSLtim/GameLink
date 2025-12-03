package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	playerhandler "gamelink/internal/handler/player"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/user"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/pkg/testutil"
)

// 用户支付后取消订单 -> 自动退款 + 状态退款
func TestOrderCancelAfterPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	migrateOrderModels(t, db)
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)

	orderService := ordersvc.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	paymentService := paymentsvc.NewPaymentService(paymentRepo, orderRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	playerAuth := fakeAuthMiddleware(seed.playerUserID)
	userGroup := api.Group("/user")
	playerGroup := api.Group("/player")
	userhandler.RegisterOrderRoutes(userGroup, orderService, userAuth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentService, userAuth)
	playerhandler.RegisterOrderRoutes(playerGroup, orderService, playerAuth)

	// 创建订单
	scheduledStart := time.Now().Add(30 * time.Minute).UTC()
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Cancel Flow",
		"scheduledStart": scheduledStart.Format(time.RFC3339),
		"durationHours":  2.0,
	}
	createResp := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
	if createResp.Code != http.StatusOK {
		t.Fatalf("create order status=%d body=%s", createResp.Code, createResp.Body.String())
	}
	var createParsed apiResp[ordersvc.CreateOrderResponse]
	if err := json.Unmarshal(createResp.Body.Bytes(), &createParsed); err != nil {
		t.Fatalf("parse create order: %v", err)
	}
	orderID := createParsed.Data.OrderID
	priceCents := createParsed.Data.PriceCents

	// 支付 -> mock 自动已支付（订单变为 confirmed + payment paid）
	payPayload := map[string]interface{}{
		"orderId": orderID,
		"method":  "alipay",
	}
	payResp := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp.Code != http.StatusOK {
		t.Fatalf("create payment status=%d body=%s", payResp.Code, payResp.Body.String())
	}

	// 用户取消（已支付订单应自动走退款并置为 refunded）
	cancelPayload := map[string]interface{}{
		"reason": "change mind",
	}
	cancelResp := doJSON(router, http.MethodPut, "/api/v1/user/orders/"+uintToStr(orderID)+"/cancel", cancelPayload, "")
	if cancelResp.Code != http.StatusOK {
		t.Fatalf("cancel order status=%d body=%s", cancelResp.Code, cancelResp.Body.String())
	}

	// 查询详情，验证退款状态与金额
	detailResp := doJSON(router, http.MethodGet, "/api/v1/user/orders/"+uintToStr(orderID), nil, "")
	if detailResp.Code != http.StatusOK {
		t.Fatalf("order detail status=%d body=%s", detailResp.Code, detailResp.Body.String())
	}
	var detailParsed apiResp[ordersvc.OrderDetailResponse]
	if err := json.Unmarshal(detailResp.Body.Bytes(), &detailParsed); err != nil {
		t.Fatalf("parse order detail: %v", err)
	}
	if detailParsed.Data.Order.Status != model.OrderStatusRefunded {
		t.Fatalf("expected refunded status, got %s", detailParsed.Data.Order.Status)
	}
	if detailParsed.Data.Order.RefundAmount != priceCents {
		t.Fatalf("refund amount mismatch: expect %d got %d", priceCents, detailParsed.Data.Order.RefundAmount)
	}
	if detailParsed.Data.Order.RefundReason == "" {
		t.Fatalf("expected refund reason to be set")
	}
}

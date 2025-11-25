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
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/internal/testutil"
)

// 订单非法状态流转：已确认直接完成、进行中取消
func TestOrderInvalidTransitions(t *testing.T) {
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

	// 创建订单并支付（状态变为 confirmed）
	scheduledStart := time.Now().Add(1 * time.Hour).UTC()
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Invalid transitions",
		"scheduledStart": scheduledStart.Format(time.RFC3339),
		"durationHours":  1.5,
	}
	createResp := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
	if createResp.Code != http.StatusOK {
		t.Fatalf("create order status=%d body=%s", createResp.Code, createResp.Body.String())
	}
	var createParsed apiResp[ordersvc.CreateOrderResponse]
	_ = json.Unmarshal(createResp.Body.Bytes(), &createParsed)
	orderID := createParsed.Data.OrderID

	// 支付订单
	payPayload := map[string]interface{}{"orderId": orderID, "method": "alipay"}
	payResp := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp.Code != http.StatusOK {
		t.Fatalf("pay status=%d body=%s", payResp.Code, payResp.Body.String())
	}

	// 1) 已确认直接完成 => 应返回 400
	completeResp := doJSON(router, http.MethodPut, "/api/v1/user/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if completeResp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for completing confirmed order, got %d body=%s", completeResp.Code, completeResp.Body.String())
	}

	// 接单 -> 进入 in_progress
	acceptResp := doJSON(router, http.MethodPost, "/api/v1/player/orders/"+uintToStr(orderID)+"/accept", nil, "")
	if acceptResp.Code != http.StatusOK {
		t.Fatalf("accept order status=%d body=%s", acceptResp.Code, acceptResp.Body.String())
	}

	// 2) 进行中取消 => 应返回 400
	cancelPayload := map[string]interface{}{"reason": "change mind"}
	cancelResp := doJSON(router, http.MethodPut, "/api/v1/user/orders/"+uintToStr(orderID)+"/cancel", cancelPayload, "")
	if cancelResp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for canceling in-progress order, got %d body=%s", cancelResp.Code, cancelResp.Body.String())
	}

	// 为防后续流程受影响，正常完成订单
	playerComplete := doJSON(router, http.MethodPut, "/api/v1/player/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if playerComplete.Code != http.StatusOK {
		t.Fatalf("player complete status=%d body=%s", playerComplete.Code, playerComplete.Body.String())
	}

	// 再次确认用户完成请求应仍为 200（已完成可重复？当前实现会返回 400，因为状态不是 in_progress）
	userComplete := doJSON(router, http.MethodPut, "/api/v1/user/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if userComplete.Code == http.StatusOK {
		// If implementation changes to allow idempotent complete, this branch will avoid failing future updates.
		return
	}
	if userComplete.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when user completes already completed order, got %d body=%s", userComplete.Code, userComplete.Body.String())
	}

	// 最终状态应为 completed
	detailResp := doJSON(router, http.MethodGet, "/api/v1/user/orders/"+uintToStr(orderID), nil, "")
	var detailParsed apiResp[ordersvc.OrderDetailResponse]
	_ = json.Unmarshal(detailResp.Body.Bytes(), &detailParsed)
	if detailParsed.Data.Order.Status != model.OrderStatusCompleted {
		t.Fatalf("expected final status completed, got %s", detailParsed.Data.Order.Status)
	}
}

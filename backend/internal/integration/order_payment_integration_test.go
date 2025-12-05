package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	playerhandler "gamelink/internal/handler/player"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"

	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/pkg/testutil"
)

// 用户创建订单 -> 支付（自动确认）-> 陪玩接单完成 -> 查询详情/列表
func TestOrderPaymentFlow(t *testing.T) {
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

	// 1) 创建订单
	scheduledStart := time.Now().Add(30 * time.Minute).UTC()
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Test Order",
		"scheduledStart": scheduledStart.Format(time.RFC3339),
		"durationHours":  2.0,
	}
	w := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
	if w.Code != http.StatusOK {
		t.Fatalf("create order status=%d body=%s", w.Code, w.Body.String())
	}
	var orderResp apiResp[ordersvc.CreateOrderResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &orderResp); err != nil {
		t.Fatalf("parse create order resp: %v", err)
	}
	if !orderResp.Success || orderResp.Data.OrderID == 0 {
		t.Fatalf("unexpected create order response: %+v", orderResp)
	}
	orderID := orderResp.Data.OrderID

	// 2) 创建支付（内部自动标记订单确认）
	payPayload := map[string]interface{}{
		"orderId": orderID,
		"method":  "alipay",
	}
	payResp := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp.Code != http.StatusOK {
		t.Fatalf("create payment status=%d body=%s", payResp.Code, payResp.Body.String())
	}
	var payParsed apiResp[paymentsvc.CreatePaymentResponse]
	if err := json.Unmarshal(payResp.Body.Bytes(), &payParsed); err != nil {
		t.Fatalf("parse payment resp: %v", err)
	}
	if !payParsed.Success || payParsed.Data.PaymentID == 0 {
		t.Fatalf("unexpected payment response: %+v", payParsed)
	}

	// 3) 玩家接单 -> 完成订单
	acceptResp := doJSON(router, http.MethodPost, "/api/v1/player/orders/"+uintToStr(orderID)+"/accept", nil, "")
	if acceptResp.Code != http.StatusOK {
		t.Fatalf("accept order status=%d body=%s", acceptResp.Code, acceptResp.Body.String())
	}
	completeResp := doJSON(router, http.MethodPut, "/api/v1/player/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if completeResp.Code != http.StatusOK {
		t.Fatalf("complete order status=%d body=%s", completeResp.Code, completeResp.Body.String())
	}

	// 4) 查询订单详情，状态应为 completed
	detailResp := doJSON(router, http.MethodGet, "/api/v1/user/orders/"+uintToStr(orderID), nil, "")
	if detailResp.Code != http.StatusOK {
		t.Fatalf("order detail status=%d body=%s", detailResp.Code, detailResp.Body.String())
	}
	var detailParsed apiResp[ordersvc.OrderDetailResponse]
	if err := json.Unmarshal(detailResp.Body.Bytes(), &detailParsed); err != nil {
		t.Fatalf("parse order detail: %v", err)
	}
	if detailParsed.Data.Order.Status != model.OrderStatusCompleted {
		t.Fatalf("expected status completed, got %s", detailParsed.Data.Order.Status)
	}

	// 5) 查询我的订单列表应包含该订单
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/orders", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("order list status=%d body=%s", listResp.Code, listResp.Body.String())
	}
	var listParsed apiResp[ordersvc.MyOrderListResponse]
	if err := json.Unmarshal(listResp.Body.Bytes(), &listParsed); err != nil {
		t.Fatalf("parse order list: %v", err)
	}
	if listParsed.Data.Total == 0 || len(listParsed.Data.Orders) == 0 {
		t.Fatalf("order list empty, resp=%+v", listParsed)
	}
}

func migrateOrderModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.Review{},
	); err != nil {
		t.Fatalf("migrate order models: %v", err)
	}
}

type orderSeed struct {
	userID       uint64
	playerUserID uint64
	playerID     uint64
	gameID       uint64
}

func seedOrderData(t *testing.T, db *gorm.DB) orderSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)

	userModel := &model.User{
		Name:         "User A",
		Email:        "user@example.com",
		Phone:        "10000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, userModel); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	playerUser := &model.User{
		Name:         "Player U",
		Email:        "player@example.com",
		Phone:        "10000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RolePlayer,
	}
	if err := userRepo.Create(ctx, playerUser); err != nil {
		t.Fatalf("seed player user: %v", err)
	}

	playerModel := &model.Player{
		UserID:          playerUser.ID,
		Nickname:        "Pro",
		HourlyRateCents: 5000,
	}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	gameModel := &model.Game{
		Key:  "lol",
		Name: "League of Legends",
	}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	return orderSeed{
		userID:       userModel.ID,
		playerUserID: playerUser.ID,
		playerID:     playerModel.ID,
		gameID:       gameModel.ID,
	}
}

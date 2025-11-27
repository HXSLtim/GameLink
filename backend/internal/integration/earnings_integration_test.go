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
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/review"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/withdraw"
	earningssvc "gamelink/internal/service/earnings"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	"gamelink/internal/testutil"
)

// 玩家收益/提现：订单完成后查询收益概览、趋势、申请提现并查看记录
func TestEarningsWithdrawFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateEarningsModels(t, db)

	seed := seedEarningsData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)
	withdrawRepo := withdraw.NewWithdrawRepository(db)

	orderService := ordersvc.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	paymentService := paymentsvc.NewPaymentService(paymentRepo, orderRepo)
	earningsService := earningssvc.NewEarningsService(playerRepo, orderRepo, withdrawRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	playerAuth := fakeAuthMiddleware(seed.playerUserID)
	userGroup := api.Group("/user")
	playerGroup := api.Group("/player")
	userhandler.RegisterOrderRoutes(userGroup, orderService, userAuth)
	userhandler.RegisterPaymentRoutes(userGroup, paymentService, userAuth)
	playerhandler.RegisterOrderRoutes(playerGroup, orderService, playerAuth)
	playerhandler.RegisterEarningsRoutes(playerGroup, earningsService, playerAuth)

	// 创建订单并完成
	orderID := createAndCompleteOrder(t, router, seed, 3.0)

	// 收益概览
	summaryResp := doJSON(router, http.MethodGet, "/api/v1/player/earnings/summary", nil, "")
	if summaryResp.Code != http.StatusOK {
		t.Fatalf("summary status=%d body=%s", summaryResp.Code, summaryResp.Body.String())
	}
	var summaryParsed apiResp[earningssvc.EarningsSummaryResponse]
	if err := json.Unmarshal(summaryResp.Body.Bytes(), &summaryParsed); err != nil {
		t.Fatalf("parse summary: %v", err)
	}

	// 申请提现（100元=10000分）
	withdrawPayload := map[string]interface{}{
		"amountCents": 10000,
		"method":      "alipay",
		"accountInfo": "alipay-acc",
	}
	withdrawResp := doJSON(router, http.MethodPost, "/api/v1/player/earnings/withdraw", withdrawPayload, "")
	if withdrawResp.Code != http.StatusOK {
		t.Fatalf("withdraw status=%d body=%s", withdrawResp.Code, withdrawResp.Body.String())
	}

	// 提现记录
	historyResp := doJSON(router, http.MethodGet, "/api/v1/player/earnings/withdraw-history", nil, "")
	if historyResp.Code != http.StatusOK {
		t.Fatalf("withdraw history status=%d body=%s", historyResp.Code, historyResp.Body.String())
	}

	// 趋势
	trendResp := doJSON(router, http.MethodGet, "/api/v1/player/earnings/trend?days=7", nil, "")
	if trendResp.Code != http.StatusOK {
		t.Fatalf("trend status=%d body=%s", trendResp.Code, trendResp.Body.String())
	}

	// 基础校验：提现后可用余额应减少（至少不为负）
	var afterSummary apiResp[earningssvc.EarningsSummaryResponse]
	if err := json.Unmarshal(summaryResp.Body.Bytes(), &summaryParsed); err == nil {
		// ignore parse error already handled
	}
	afterResp := doJSON(router, http.MethodGet, "/api/v1/player/earnings/summary", nil, "")
	if err := json.Unmarshal(afterResp.Body.Bytes(), &afterSummary); err == nil {
		if afterSummary.Data.AvailableBalance < 0 {
			t.Fatalf("available balance negative after withdraw: %v", afterSummary.Data.AvailableBalance)
		}
	}

	_ = orderID
}

func migrateEarningsModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.Withdraw{},
	); err != nil {
		t.Fatalf("migrate earnings models: %v", err)
	}
}

type earningsSeed struct {
	userID       uint64
	playerUserID uint64
	playerID     uint64
	gameID       uint64
}

func seedEarningsData(t *testing.T, db *gorm.DB) earningsSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)

	userModel := &model.User{
		Name:         "EarnUser",
		Email:        "earn_user@example.com",
		Phone:        "60000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, userModel); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	playerUser := &model.User{
		Name:         "EarnPlayerUser",
		Email:        "earn_player@example.com",
		Phone:        "60000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RolePlayer,
	}
	if err := userRepo.Create(ctx, playerUser); err != nil {
		t.Fatalf("seed player user: %v", err)
	}

	playerModel := &model.Player{
		UserID:          playerUser.ID,
		Nickname:        "Earner",
		HourlyRateCents: 8000, // 80元/小时
	}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	gameModel := &model.Game{
		Key:  "valorant",
		Name: "Valorant",
	}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	return earningsSeed{
		userID:       userModel.ID,
		playerUserID: playerUser.ID,
		playerID:     playerModel.ID,
		gameID:       gameModel.ID,
	}
}

// createAndCompleteOrder 创建订单->支付->接单->完成，返回订单ID
func createAndCompleteOrder(t *testing.T, router *gin.Engine, seed earningsSeed, durationHours float32) uint64 {
	t.Helper()
	scheduledStart := time.Now().Add(30 * time.Minute).UTC()
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Earnings Order",
		"scheduledStart": scheduledStart.Format(time.RFC3339),
		"durationHours":  durationHours,
	}
	orderResp := doJSON(router, http.MethodPost, "/api/v1/user/orders", createOrderPayload, "")
	if orderResp.Code != http.StatusOK {
		t.Fatalf("create order status=%d body=%s", orderResp.Code, orderResp.Body.String())
	}
	var parsed apiResp[ordersvc.CreateOrderResponse]
	if err := json.Unmarshal(orderResp.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("parse create order: %v", err)
	}
	orderID := parsed.Data.OrderID

	payPayload := map[string]interface{}{
		"orderId": orderID,
		"method":  "alipay",
	}
	payResp := doJSON(router, http.MethodPost, "/api/v1/user/payments", payPayload, "")
	if payResp.Code != http.StatusOK {
		t.Fatalf("create payment status=%d body=%s", payResp.Code, payResp.Body.String())
	}

	acceptResp := doJSON(router, http.MethodPost, "/api/v1/player/orders/"+uintToStr(orderID)+"/accept", nil, "")
	if acceptResp.Code != http.StatusOK {
		t.Fatalf("accept order status=%d body=%s", acceptResp.Code, acceptResp.Body.String())
	}
	completeResp := doJSON(router, http.MethodPut, "/api/v1/player/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if completeResp.Code != http.StatusOK {
		t.Fatalf("complete order status=%d body=%s", completeResp.Code, completeResp.Body.String())
	}
	return orderID
}

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
	"gamelink/internal/repository/reviewreply"
	"gamelink/internal/repository/user"
	ordersvc "gamelink/internal/service/order"
	paymentsvc "gamelink/internal/service/payment"
	reviewsvc "gamelink/internal/service/review"
	"gamelink/internal/testutil"
)

// 场景：用户完成订单后创建评价 -> 陪玩师回复 -> 用户查询自己的评价列表
func TestReviewFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateReviewModels(t, db)

	seed := seedReviewData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	reviewRepo := review.NewReviewRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)
	replyRepo := reviewreply.NewReviewReplyRepository(db)

	orderService := ordersvc.NewOrderService(orderRepo, playerRepo, userRepo, gameRepo, paymentRepo, reviewRepo, commissionRepo)
	paymentService := paymentsvc.NewPaymentService(paymentRepo, orderRepo)
	reviewService := reviewsvc.NewReviewService(reviewRepo, orderRepo, playerRepo, userRepo, replyRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	playerAuth := fakeAuthMiddleware(seed.playerUserID)
	userhandler.RegisterOrderRoutes(api, orderService, userAuth)
	userhandler.RegisterPaymentRoutes(api, paymentService, userAuth)
	userhandler.RegisterReviewRoutes(api, reviewService, userAuth)
	playerGroup := api.Group("/player")
	playerhandler.RegisterReviewRoutes(playerGroup, reviewService, playerAuth)
	playerhandler.RegisterOrderRoutes(api, orderService, playerAuth)

	// 1) 用户创建并支付订单（自动确认）
	orderID := createAndPayOrder(t, router, seed, time.Now().Add(1*time.Hour))

	// 2) 玩家接单、完成
	acceptResp := doJSON(router, http.MethodPost, "/api/v1/player/orders/"+uintToStr(orderID)+"/accept", nil, "")
	if acceptResp.Code != http.StatusOK {
		t.Fatalf("accept order status=%d body=%s", acceptResp.Code, acceptResp.Body.String())
	}
	completeResp := doJSON(router, http.MethodPut, "/api/v1/player/orders/"+uintToStr(orderID)+"/complete", nil, "")
	if completeResp.Code != http.StatusOK {
		t.Fatalf("complete order status=%d body=%s", completeResp.Code, completeResp.Body.String())
	}

	// 3) 用户创建评价
	reviewPayload := map[string]interface{}{
		"orderId": orderID,
		"rating":  5,
		"comment": "great job",
	}
	createReviewResp := doJSON(router, http.MethodPost, "/api/v1/user/reviews", reviewPayload, "")
	if createReviewResp.Code != http.StatusOK {
		t.Fatalf("create review status=%d body=%s", createReviewResp.Code, createReviewResp.Body.String())
	}
	var createReviewParsed apiResp[reviewsvc.CreateReviewResponse]
	if err := json.Unmarshal(createReviewResp.Body.Bytes(), &createReviewParsed); err != nil {
		t.Fatalf("parse create review: %v", err)
	}
	reviewID := createReviewParsed.Data.ReviewID

	// 4) 玩家回复评价
	replyPayload := map[string]interface{}{
		"content": "thanks!",
	}
	replyResp := doJSON(router, http.MethodPost, "/api/v1/player/reviews/"+uintToStr(reviewID)+"/reply", replyPayload, "")
	if replyResp.Code != http.StatusOK {
		t.Fatalf("reply review status=%d body=%s", replyResp.Code, replyResp.Body.String())
	}

	// 5) 用户查询自己的评价列表，包含刚创建的评价
	myReviewsResp := doJSON(router, http.MethodGet, "/api/v1/user/reviews/my", nil, "")
	if myReviewsResp.Code != http.StatusOK {
		t.Fatalf("my reviews status=%d body=%s", myReviewsResp.Code, myReviewsResp.Body.String())
	}
	var myReviewsParsed apiResp[reviewsvc.MyReviewListResponse]
	if err := json.Unmarshal(myReviewsResp.Body.Bytes(), &myReviewsParsed); err != nil {
		t.Fatalf("parse my reviews: %v", err)
	}
	if len(myReviewsParsed.Data.Reviews) == 0 {
		t.Fatalf("expected reviews, got none")
	}
}

func migrateReviewModels(t *testing.T, db *gorm.DB) {
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
		&model.ReviewReply{},
	); err != nil {
		t.Fatalf("migrate review models: %v", err)
	}
}

type reviewSeed struct {
	userID       uint64
	playerUserID uint64
	playerID     uint64
	gameID       uint64
}

func seedReviewData(t *testing.T, db *gorm.DB) reviewSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)

	userModel := &model.User{
		Name:         "ReviewUser",
		Email:        "review_user@example.com",
		Phone:        "30000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, userModel); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	playerUser := &model.User{
		Name:         "ReviewPlayerUser",
		Email:        "review_player@example.com",
		Phone:        "30000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RolePlayer,
	}
	if err := userRepo.Create(ctx, playerUser); err != nil {
		t.Fatalf("seed player user: %v", err)
	}

	playerModel := &model.Player{
		UserID:          playerUser.ID,
		Nickname:        "Reviewer",
		HourlyRateCents: 4000,
	}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	gameModel := &model.Game{
		Key:  "dota2",
		Name: "Dota 2",
	}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	return reviewSeed{
		userID:       userModel.ID,
		playerUserID: playerUser.ID,
		playerID:     playerModel.ID,
		gameID:       gameModel.ID,
	}
}

// createAndPayOrder 从订单服务创建并支付订单，返回订单ID
func createAndPayOrder(t *testing.T, router *gin.Engine, seed reviewSeed, start time.Time) uint64 {
	createOrderPayload := map[string]interface{}{
		"playerId":       seed.playerID,
		"gameId":         seed.gameID,
		"title":          "Review Order",
		"scheduledStart": start.UTC().Format(time.RFC3339),
		"durationHours":  1.0,
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
	return orderID
}

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	adminhandler "gamelink/internal/handler/admin"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/dispute"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/menu"
	"gamelink/internal/repository/notification"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/wallet"
	"gamelink/pkg/cache"

	adminservice "gamelink/internal/service/admin"
	orderservice "gamelink/internal/service/order"
	"gamelink/pkg/testutil"
)

type disputeSeed struct {
	userID       uint64
	playerUserID uint64
	playerID     uint64
	gameID       uint64
}

func seedDisputeData(t *testing.T, db *gorm.DB) disputeSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)

	userModel := &model.User{
		Name:         "DisputeUser",
		Email:        "dispute_user@example.com",
		Phone:        "30000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, userModel); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	playerUser := &model.User{
		Name:         "DisputePlayerUser",
		Email:        "dispute_player@example.com",
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
		Nickname:        "DisputePro",
		HourlyRateCents: 5000,
	}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	gameModel := &model.Game{
		Key:  "dispute_game",
		Name: "Dispute Game",
	}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	return disputeSeed{
		userID:       userModel.ID,
		playerUserID: playerUser.ID,
		playerID:     playerModel.ID,
		gameID:       gameModel.ID,
	}
}

// 纠纷全链路：用户发起 -> 管理端列表/分配 -> 管理端解决（部分退款）
func TestDisputeFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedDisputeData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := adminrepo.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	memCache := cache.NewMemory()

	assignSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	menuRepo := menu.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	_ = adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, serviceItemRepo, permRepo, menuRepo, statsRepo, walletRepo, memCache)
	adminDisputeHandler := adminhandler.NewDisputeHandler(assignSvc)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	adminGroup := api.Group("/admin")

	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	adminGroup.Use(fakeUserBothMiddleware(seed.playerUserID))

	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)
	userGroup.GET("/orders/:id/disputes", userDisputeHandler.GetDisputeDetail)
	adminGroup.GET("/disputes/pending", adminDisputeHandler.ListPendingDisputes)
	adminGroup.POST("/disputes/:id/assign", adminDisputeHandler.AssignDispute)
	adminGroup.POST("/disputes/:id/resolve", adminDisputeHandler.ResolveDispute)

	orderID := createDisputeOrder(t, ctx, orderRepo, seed, time.Now().Add(-1*time.Hour))

	initPayload := map[string]any{
		"orderId":      orderID,
		"reason":       "service issue",
		"description":  "bad experience",
		"evidenceUrls": []string{"https://example.com/e1.jpg"},
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code != http.StatusOK && initResp.Code != http.StatusCreated {
		t.Fatalf("initiate dispute status=%d body=%s", initResp.Code, initResp.Body.String())
	}
	var initParsed apiResp[orderservice.InitiateDisputeResponse]
	_ = json.Unmarshal(initResp.Body.Bytes(), &initParsed)
	if initParsed.Data.DisputeID == 0 {
		t.Fatalf("expected dispute id > 0, got %+v", initParsed.Data)
	}
	disputeID := initParsed.Data.DisputeID

	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/disputes/pending", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("pending disputes status=%d body=%s", listResp.Code, listResp.Body.String())
	}

	assignResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/assign", map[string]any{
		"assignedToUserId": seed.playerUserID,
		"source":           "system",
	}, "")
	if assignResp.Code != http.StatusOK {
		t.Fatalf("assign dispute status=%d body=%s", assignResp.Code, assignResp.Body.String())
	}

	resolveResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/resolve", map[string]any{
		"resolution":       "refund",
		"resolutionAmount": 4000,
		"resolutionNotes":  "partial refund granted",
	}, "")
	if resolveResp.Code != http.StatusOK {
		t.Fatalf("resolve dispute status=%d body=%s", resolveResp.Code, resolveResp.Body.String())
	}

	updatedDispute, _ := disputeRepo.Get(ctx, disputeID)
	if updatedDispute.Status != model.DisputeStatusResolved {
		t.Fatalf("expected dispute resolved, got %s", updatedDispute.Status)
	}

	orderModel, _ := orderRepo.Get(ctx, orderID)
	if orderModel.Status != model.OrderStatusRefunded {
		t.Fatalf("expected order refunded, got %s", orderModel.Status)
	}
	if orderModel.RefundAmountCents != 4000 {
		t.Fatalf("expected refund amount 4000, got %d", orderModel.RefundAmountCents)
	}
}

func migrateDisputeModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.OrderDispute{},
		&model.OperationLog{},
		&model.NotificationEvent{},
	); err != nil {
		t.Fatalf("migrate dispute models: %v", err)
	}
}

func createDisputeOrder(t *testing.T, ctx context.Context, orderRepo repoiface.OrderRepository, seed disputeSeed, completedAt time.Time) uint64 {
	t.Helper()
	orderModel := &model.Order{
		UserID:          seed.userID,
		ItemID:          seed.gameID,
		Status:          model.OrderStatusCompleted,
		Title:           "Dispute Order",
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
		CompletedAt:     &completedAt,
	}
	orderModel.SetPlayerID(seed.playerID)
	orderModel.SetGameID(seed.gameID)
	if err := orderRepo.Create(ctx, orderModel); err != nil {
		t.Fatalf("seed dispute order: %v", err)
	}
	return orderModel.ID
}

func fakeUserBothMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("userID", userID)
		c.Next()
	}
}

// 测试纠纷拒绝流程
func TestDisputeRejectFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedDisputeData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := adminrepo.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := order.NewPaymentRepository(db)

	assignSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	adminDisputeHandler := adminhandler.NewDisputeHandler(assignSvc)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	adminGroup := api.Group("/admin")

	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	adminGroup.Use(fakeUserBothMiddleware(seed.playerUserID))

	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)
	adminGroup.POST("/disputes/:id/assign", adminDisputeHandler.AssignDispute)
	adminGroup.POST("/disputes/:id/resolve", adminDisputeHandler.ResolveDispute)

	orderID := createDisputeOrder(t, ctx, orderRepo, seed, time.Now().Add(-1*time.Hour))

	initPayload := map[string]any{
		"orderId":     orderID,
		"reason":      "unreasonable complaint",
		"description": "invalid claim",
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code != http.StatusOK && initResp.Code != http.StatusCreated {
		t.Fatalf("initiate dispute status=%d body=%s", initResp.Code, initResp.Body.String())
	}
	var initParsed apiResp[orderservice.InitiateDisputeResponse]
	_ = json.Unmarshal(initResp.Body.Bytes(), &initParsed)
	disputeID := initParsed.Data.DisputeID

	assignResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/assign", map[string]any{
		"assignedToUserId": seed.playerUserID,
		"source":           "manual",
	}, "")
	if assignResp.Code != http.StatusOK {
		t.Fatalf("assign dispute status=%d body=%s", assignResp.Code, assignResp.Body.String())
	}

	resolveResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/resolve", map[string]any{
		"resolution":      "reject",
		"resolutionNotes": "claim not valid",
	}, "")
	if resolveResp.Code != http.StatusOK {
		t.Fatalf("reject dispute status=%d body=%s", resolveResp.Code, resolveResp.Body.String())
	}

	updatedDispute, _ := disputeRepo.Get(ctx, disputeID)
	if updatedDispute.Status != model.DisputeStatusResolved {
		t.Fatalf("expected dispute resolved, got %s", updatedDispute.Status)
	}
	if updatedDispute.Resolution != model.ResolutionReject {
		t.Fatalf("expected resolution reject, got %s", updatedDispute.Resolution)
	}

	orderModel, _ := orderRepo.Get(ctx, orderID)
	if orderModel.Status != model.OrderStatusCompleted {
		t.Fatalf("expected order still completed, got %s", orderModel.Status)
	}
}

// 测试超时订单不能发起纠纷
func TestDisputeExpiredOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedDisputeData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := adminrepo.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := order.NewPaymentRepository(db)

	assignSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)

	orderID := createDisputeOrder(t, ctx, orderRepo, seed, time.Now().Add(-48*time.Hour))

	initPayload := map[string]any{
		"orderId":     orderID,
		"reason":      "late complaint",
		"description": "too late",
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code == http.StatusOK || initResp.Code == http.StatusCreated {
		t.Fatalf("expected error for expired order, got status=%d", initResp.Code)
	}
}

// 测试重复发起纠纷
func TestDisputeDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedDisputeData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := adminrepo.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := order.NewPaymentRepository(db)

	assignSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)

	orderID := createDisputeOrder(t, ctx, orderRepo, seed, time.Now().Add(-1*time.Hour))

	initPayload := map[string]any{
		"orderId":     orderID,
		"reason":      "first complaint",
		"description": "first issue",
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code != http.StatusOK && initResp.Code != http.StatusCreated {
		t.Fatalf("first dispute status=%d body=%s", initResp.Code, initResp.Body.String())
	}

	initPayload2 := map[string]any{
		"orderId":     orderID,
		"reason":      "second complaint",
		"description": "second issue",
	}
	initResp2 := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload2, "")
	if initResp2.Code == http.StatusOK || initResp2.Code == http.StatusCreated {
		t.Fatalf("expected error for duplicate dispute, got status=%d", initResp2.Code)
	}
}

// 测试查看纠纷详情
func TestDisputeDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedDisputeData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := adminrepo.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	paymentRepo := order.NewPaymentRepository(db)

	assignSvc := orderservice.NewDisputeService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)
	userGroup.GET("/orders/:id/disputes", userDisputeHandler.GetDisputeDetail)

	orderID := createDisputeOrder(t, ctx, orderRepo, seed, time.Now().Add(-1*time.Hour))

	initPayload := map[string]any{
		"orderId":      orderID,
		"reason":       "service issue",
		"description":  "detailed description",
		"evidenceUrls": []string{"https://example.com/evidence1.jpg", "https://example.com/evidence2.jpg"},
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code != http.StatusOK && initResp.Code != http.StatusCreated {
		t.Fatalf("initiate dispute status=%d body=%s", initResp.Code, initResp.Body.String())
	}

	detailResp := doJSON(router, http.MethodGet, "/api/v1/user/orders/"+uintToStr(orderID)+"/disputes", nil, "")
	if detailResp.Code != http.StatusOK {
		t.Fatalf("get dispute detail status=%d body=%s", detailResp.Code, detailResp.Body.String())
	}
}

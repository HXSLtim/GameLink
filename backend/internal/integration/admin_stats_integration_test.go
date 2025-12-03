package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/user"
	statsrepo "gamelink/internal/repository/stats"
	"gamelink/internal/service/admin"
	"gamelink/pkg/testutil"
)

// 管理端统计：dashboard + revenue/user/orders
func TestAdminStatsEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateStatsModels(t, db)

	seedStatsData(t, db)

	statsSvc := stats.NewStatsService(statsrepo.NewStatsRepository(db))
	statsHandler := adminhandler.NewStatsHandler(statsSvc)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.GET("/stats/dashboard", statsHandler.Dashboard)
	adminGroup.GET("/stats/revenue-trend", statsHandler.RevenueTrend)
	adminGroup.GET("/stats/user-growth", statsHandler.UserGrowth)
	adminGroup.GET("/stats/orders", statsHandler.OrdersSummary)

	// Dashboard
	dashboardResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/dashboard", nil, "")
	if dashboardResp.Code != http.StatusOK {
		t.Fatalf("dashboard status=%d body=%s", dashboardResp.Code, dashboardResp.Body.String())
	}
	var dashParsed apiResp[map[string]interface{}]
	_ = json.Unmarshal(dashboardResp.Body.Bytes(), &dashParsed)

	// Revenue trend
	revenueResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/revenue-trend?days=7", nil, "")
	if revenueResp.Code != http.StatusOK {
		t.Fatalf("revenue status=%d body=%s", revenueResp.Code, revenueResp.Body.String())
	}

	// User growth
	userGrowthResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/user-growth?days=7", nil, "")
	if userGrowthResp.Code != http.StatusOK {
		t.Fatalf("user growth status=%d body=%s", userGrowthResp.Code, userGrowthResp.Body.String())
	}

	// Orders summary
	ordersResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/orders", nil, "")
	if ordersResp.Code != http.StatusOK {
		t.Fatalf("orders summary status=%d body=%s", ordersResp.Code, ordersResp.Body.String())
	}
}

// 顶级陪玩&审计统计
func TestAdminStatsTopPlayersAndAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateStatsModels(t, db)

	// seed players with ratings
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	opRepo := adminrepo.NewOperationLogRepository(db)

	u1 := &model.User{Name: "P1", Email: "p1@example.com", Phone: "18000000011", PasswordHash: "x", Role: model.RolePlayer, Status: model.UserStatusActive}
	u2 := &model.User{Name: "P2", Email: "p2@example.com", Phone: "18000000022", PasswordHash: "x", Role: model.RolePlayer, Status: model.UserStatusActive}
	_ = userRepo.Create(context.Background(), u1)
	_ = userRepo.Create(context.Background(), u2)
	p1 := &model.Player{UserID: u1.ID, Nickname: "A", RatingAverage: 4.9, RatingCount: 10}
	p2 := &model.Player{UserID: u2.ID, Nickname: "B", RatingAverage: 4.7, RatingCount: 5}
	_ = playerRepo.Create(context.Background(), p1)
	_ = playerRepo.Create(context.Background(), p2)

	// seed operation logs for audit stats
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	_ = opRepo.Append(context.Background(), &model.OperationLog{Base: model.Base{CreatedAt: yesterday}, EntityType: "order", EntityID: 1, Action: "create"})
	_ = opRepo.Append(context.Background(), &model.OperationLog{Base: model.Base{CreatedAt: now}, EntityType: "order", EntityID: 2, Action: "create"})

	// router
	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminhandler.RegisterStatsAnalysisRoutes(adminGroup, db)

	// top players
	topResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/top-players?limit=1", nil, "")
	if topResp.Code != http.StatusOK {
		t.Fatalf("top players status=%d body=%s", topResp.Code, topResp.Body.String())
	}

	// audit overview
	overviewResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/audit/overview", nil, "")
	if overviewResp.Code != http.StatusOK {
		t.Fatalf("audit overview status=%d body=%s", overviewResp.Code, overviewResp.Body.String())
	}

	// audit trend
	trendResp := doJSON(router, http.MethodGet, "/api/v1/admin/stats/audit/trend?entity=order&action=create", nil, "")
	if trendResp.Code != http.StatusOK {
		t.Fatalf("audit trend status=%d body=%s", trendResp.Code, trendResp.Body.String())
	}
}

func migrateStatsModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate stats models: %v", err)
	}
}

func seedStatsData(t *testing.T, db *gorm.DB) {
	t.Helper()
	ctx := context.Background()
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	orderRepo := orderimpl.NewOrderRepository(db)
	paymentRepo := order.NewPaymentRepository(db)

	// users
	userA := &model.User{Name: "A", Email: "a@example.com", Phone: "10000000001", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	userB := &model.User{Name: "B", Email: "b@example.com", Phone: "10000000002", PasswordHash: "x", Role: model.RoleUser, Status: model.UserStatusActive}
	if err := userRepo.Create(ctx, userA); err != nil {
		t.Fatalf("seed userA: %v", err)
	}
	if err := userRepo.Create(ctx, userB); err != nil {
		t.Fatalf("seed userB: %v", err)
	}

	// player
	playerModel := &model.Player{UserID: userA.ID, Nickname: "Pro", HourlyRateCents: 5000}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	// game
	gameModel := &model.Game{Key: "game", Name: "Game"}
	if err := gameRepo.Create(ctx, gameModel); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	// orders
	now := time.Now().UTC()
	order1 := &model.Order{
		UserID:          userA.ID,
		ItemID:          gameModel.ID,
		Status:          model.OrderStatusCompleted,
		Title:           "Finished",
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
		CompletedAt:     &now,
	}
	order1.SetPlayerID(playerModel.ID)
	order1.SetGameID(gameModel.ID)
	if err := orderRepo.Create(ctx, order1); err != nil {
		t.Fatalf("seed order1: %v", err)
	}

	order2 := &model.Order{
		UserID:          userB.ID,
		ItemID:          gameModel.ID,
		Status:          model.OrderStatusPending,
		Title:           "Pending",
		UnitPriceCents:  8000,
		TotalPriceCents: 8000,
	}
	order2.SetPlayerID(playerModel.ID)
	order2.SetGameID(gameModel.ID)
	if err := orderRepo.Create(ctx, order2); err != nil {
		t.Fatalf("seed order2: %v", err)
	}

	// payment for order1 (paid)
	paidAt := now.Add(-2 * time.Hour)
	paymentModel := &model.Payment{
		OrderID:     order1.ID,
		UserID:      userA.ID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: order1.TotalPriceCents,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &paidAt,
	}
	if err := paymentRepo.Create(ctx, paymentModel); err != nil {
		t.Fatalf("seed payment: %v", err)
	}
}

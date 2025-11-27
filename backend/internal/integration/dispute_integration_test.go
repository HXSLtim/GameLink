package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"gamelink/internal/cache"
	adminhandler "gamelink/internal/handler/admin"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/dispute"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/internal/repository/notification"
	operationlog "gamelink/internal/repository/operation_log"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/role"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/service/assignment"
	"gamelink/internal/testutil"
)

// 纠纷全链路：用户发起 -> 管理端列表/分配 -> 管理端解决（部分退款）
func TestDisputeFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateDisputeModels(t, db)

	ctx := context.Background()
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := payment.NewPaymentRepository(db)
	roleRepo := role.NewRoleRepository(db)
	disputeRepo := dispute.NewDisputeRepository(db)
	opLogRepo := operationlog.NewOperationLogRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	memCache := cache.NewMemory()

	assignSvc := assignment.NewAssignmentService(disputeRepo, orderRepo, userRepo, opLogRepo, notificationRepo, paymentRepo)
	_ = adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, memCache) // retain parity with admin handler deps
	adminDisputeHandler := adminhandler.NewDisputeHandler(assignSvc)
	userDisputeHandler := userhandler.NewDisputeHandler(assignSvc)

	router := gin.New()
	api := router.Group("/api/v1")
	userGroup := api.Group("/user")
	adminGroup := api.Group("/admin")

	userGroup.Use(fakeUserBothMiddleware(seed.userID))
	adminGroup.Use(fakeUserBothMiddleware(seed.playerUserID)) // actor user id for admin ops

	userGroup.POST("/orders/:id/dispute", userDisputeHandler.InitiateDispute)
	userGroup.GET("/orders/:id/disputes", userDisputeHandler.GetDisputeDetail)
	adminGroup.GET("/disputes/pending", adminDisputeHandler.ListPendingDisputes)
	adminGroup.POST("/disputes/:id/assign", adminDisputeHandler.AssignDispute)
	adminGroup.POST("/disputes/:id/resolve", adminDisputeHandler.ResolveDispute)

	// seed completed order within 24h (可发起纠纷)
	orderID := createCompletedOrder(t, ctx, orderRepo, seed, time.Now().Add(-1*time.Hour))

	// 用户发起纠纷
	initPayload := map[string]interface{}{
		"orderId":      orderID,
		"reason":       "service issue",
		"description":  "bad experience",
		"evidenceUrls": []string{"https://example.com/e1.jpg"},
	}
	initResp := doJSON(router, http.MethodPost, "/api/v1/user/orders/"+uintToStr(orderID)+"/dispute", initPayload, "")
	if initResp.Code != http.StatusOK && initResp.Code != http.StatusCreated {
		t.Fatalf("initiate dispute status=%d body=%s", initResp.Code, initResp.Body.String())
	}
	var initParsed apiResp[assignment.InitiateDisputeResponse]
	_ = json.Unmarshal(initResp.Body.Bytes(), &initParsed)
	if initParsed.Data.DisputeID == 0 {
		t.Fatalf("expected dispute id > 0, got %+v", initParsed.Data)
	}
	disputeID := initParsed.Data.DisputeID

	// 管理端查询待分配列表
	listResp := doJSON(router, http.MethodGet, "/api/v1/admin/disputes/pending", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("pending disputes status=%d body=%s", listResp.Code, listResp.Body.String())
	}

	// 管理端分配纠纷
	assignResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/assign", map[string]interface{}{
		"assignedToUserId": seed.playerUserID,
		"source":           "system",
	}, "")
	if assignResp.Code != http.StatusOK {
		t.Fatalf("assign dispute status=%d body=%s", assignResp.Code, assignResp.Body.String())
	}

	// 管理端解决纠纷（部分退款）
	resolveResp := doJSON(router, http.MethodPost, "/api/v1/admin/disputes/"+uintToStr(disputeID)+"/resolve", map[string]interface{}{
		"resolution":       "refund",
		"resolutionAmount": 4000,
		"resolutionNotes":  "partial refund granted",
	}, "")
	if resolveResp.Code != http.StatusOK {
		t.Fatalf("resolve dispute status=%d body=%s", resolveResp.Code, resolveResp.Body.String())
	}

	// 验证纠纷状态 & 订单退款状态
	updatedDispute, _ := disputeRepo.Get(ctx, disputeID)
	if updatedDispute.Status != model.DisputeStatusResolved {
		t.Fatalf("expected dispute resolved, got %s", updatedDispute.Status)
	}

	order, _ := orderRepo.Get(ctx, orderID)
	if order.Status != model.OrderStatusRefunded {
		t.Fatalf("expected order refunded, got %s", order.Status)
	}
	if order.RefundAmountCents != 4000 {
		t.Fatalf("expected refund amount 4000, got %d", order.RefundAmountCents)
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

func createCompletedOrder(t *testing.T, ctx context.Context, orderRepo repoiface.OrderRepository, seed orderSeed, completedAt time.Time) uint64 {
	t.Helper()
	order := &model.Order{
		UserID:          seed.userID,
		ItemID:          seed.gameID,
		Status:          model.OrderStatusCompleted,
		Title:           "Dispute Order",
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
		CompletedAt:     &completedAt,
	}
	order.SetPlayerID(seed.playerID)
	order.SetGameID(seed.gameID)
	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("seed dispute order: %v", err)
	}
	return order.ID
}

// fakeUserBothMiddleware sets both user_id and userID for compatibility with handlers
func fakeUserBothMiddleware(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("userID", userID)
		c.Next()
	}
}

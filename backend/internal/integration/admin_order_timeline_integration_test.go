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
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/common"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/menu"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/internal/repository/wallet"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// 管理端订单时间线：操作日志+支付事件汇总
func TestAdminOrderTimeline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateTimelineModels(t, db)

	ctx := context.Background()
	seed := seedOrderData(t, db)

	orderRepo := orderimpl.NewOrderRepository(db)
	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	gameRepo := game.NewGameRepository(db)
	paymentRepo := order.NewPaymentRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	opRepo := adminrepo.NewOperationLogRepository(db)
	memCache := cache.NewMemory()
	serviceItemRepo := serviceitem.NewServiceItemRepository(db)

	permRepo := permission.NewPermissionRepository(db)
	menuRepo := menu.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	walletRepo := wallet.NewWalletRepository(db)
	adminSvc := adminservice.NewAdminService(gameRepo, userRepo, playerRepo, orderRepo, paymentRepo, roleRepo, serviceItemRepo, permRepo, menuRepo, statsRepo, walletRepo, memCache)
	adminSvc.SetTxManager(common.NewUnitOfWork(db))
	orderHandler := adminhandler.NewOrderHandler(adminSvc)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.GET("/orders/:id/timeline", orderHandler.GetOrderTimeline)

	// 创建订单
	order := &model.Order{
		UserID:          seed.userID,
		ItemID:          seed.gameID,
		Status:          model.OrderStatusPending,
		Title:           "Timeline Order",
		UnitPriceCents:  10000,
		TotalPriceCents: 10000,
	}
	order.SetPlayerID(seed.playerID)
	order.SetGameID(seed.gameID)
	if err := orderRepo.Create(ctx, order); err != nil {
		t.Fatalf("seed order: %v", err)
	}

	// 插入操作日志：状态 pending -> confirmed
	meta, _ := json.Marshal(map[string]any{
		"from_status": string(model.OrderStatusPending),
		"status":      string(model.OrderStatusConfirmed),
	})
	if err := opRepo.Append(ctx, &model.OperationLog{
		EntityType:   string(model.OpEntityOrder),
		EntityID:     order.ID,
		Action:       string(model.OpActionUpdateStatus),
		MetadataJSON: meta,
	}); err != nil {
		t.Fatalf("seed op log: %v", err)
	}

	// 插入支付记录（已支付）
	now := time.Now().Add(-10 * time.Minute).UTC()
	pay := &model.Payment{
		OrderID:     order.ID,
		UserID:      seed.userID,
		Method:      model.PaymentMethodAlipay,
		AmountCents: order.TotalPriceCents,
		Status:      model.PaymentStatusPaid,
		PaidAt:      &now,
	}
	if err := paymentRepo.Create(ctx, pay); err != nil {
		t.Fatalf("seed payment: %v", err)
	}

	// 查询时间线
	resp := doJSON(router, http.MethodGet, "/api/v1/admin/orders/"+uintToStr(order.ID)+"/timeline", nil, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("timeline status=%d body=%s", resp.Code, resp.Body.String())
	}
	var parsed apiResp[[]adminservice.OrderTimelineItem]
	if err := json.Unmarshal(resp.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("parse timeline: %v", err)
	}
	if len(parsed.Data) == 0 {
		t.Fatalf("expected timeline items, got 0")
	}
}

func migrateTimelineModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate timeline models: %v", err)
	}
}

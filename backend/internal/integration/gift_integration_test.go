package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	playerhandler "gamelink/internal/handler/player"
	userhandler "gamelink/internal/handler/user"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/user"

	giftsvc "gamelink/internal/service/gift"
	itemsvc "gamelink/internal/service/item"
	"gamelink/pkg/testutil"
)

// 礼物业务集成：列举礼物->发送礼物->陪玩师查看收到的礼物与统计
func TestGiftFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateGiftModels(t, db)

	seed := seedGiftData(t, db)

	// repositories & services
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	itemRepo := serviceitem.NewServiceItemRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)

	itemService := itemsvc.NewServiceItemService(itemRepo, nil, playerRepo)
	giftService := giftsvc.NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	// player gift 路由假设 user_id == playerID，因此这里用 playerID
	playerGiftAuth := fakeAuthMiddleware(seed.playerID)
	userGroup := api.Group("/user")
	playerGroup := api.Group("/player")
	userhandler.RegisterGiftRoutes(userGroup, giftService, itemService, userAuth)
	playerhandler.RegisterGiftRoutes(playerGroup, giftService, playerGiftAuth)

	// 1) 列出礼物
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/gifts", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list gifts status=%d body=%s", listResp.Code, listResp.Body.String())
	}

	// 2) 发送礼物
	sendPayload := map[string]interface{}{
		"playerId":   seed.playerID,
		"giftItemId": seed.giftItemID,
		"quantity":   3,
		"message":    "nice",
	}
	sendResp := doJSON(router, http.MethodPost, "/api/v1/user/gifts/send", sendPayload, "")
	if sendResp.Code != http.StatusOK {
		t.Fatalf("send gift status=%d body=%s", sendResp.Code, sendResp.Body.String())
	}

	// 3) 陪玩端查看收到的礼物
	receivedResp := doJSON(router, http.MethodGet, "/api/v1/player/gifts/received", nil, "")
	if receivedResp.Code != http.StatusOK {
		t.Fatalf("received gifts status=%d body=%s", receivedResp.Code, receivedResp.Body.String())
	}
	var receivedParsed apiResp[giftsvc.ReceivedGiftsResponse]
	if err := json.Unmarshal(receivedResp.Body.Bytes(), &receivedParsed); err != nil {
		t.Fatalf("parse received gifts: %v", err)
	}
	if len(receivedParsed.Data.Gifts) == 0 {
		t.Fatalf("expected received gifts, got none")
	}

	// 4) 礼物统计
	statsResp := doJSON(router, http.MethodGet, "/api/v1/player/gifts/stats", nil, "")
	if statsResp.Code != http.StatusOK {
		t.Fatalf("gift stats status=%d body=%s", statsResp.Code, statsResp.Body.String())
	}
	var statsParsed apiResp[giftsvc.GiftStatsResponse]
	if err := json.Unmarshal(statsResp.Body.Bytes(), &statsParsed); err != nil {
		t.Fatalf("parse gift stats: %v", err)
	}
	if statsParsed.Data.TotalGiftOrders == 0 || statsParsed.Data.TotalGiftsReceived == 0 {
		t.Fatalf("unexpected stats: %+v", statsParsed.Data)
	}
}

// 下架的礼物不应出现在用户礼物列表
func TestGiftListExcludeInactive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateGiftModels(t, db)

	seed := seedGiftData(t, db)
	itemRepo := serviceitem.NewServiceItemRepository(db)

	// 下架礼物
	_ = itemRepo.BatchUpdateStatus(ctx(), []uint64{seed.giftItemID}, false)

	// 路由
	orderRepo := orderimpl.NewOrderRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	itemRepo = serviceitem.NewServiceItemRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)
	itemService := itemsvc.NewServiceItemService(itemRepo, nil, playerRepo)
	giftService := giftsvc.NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

	router := gin.New()
	api := router.Group("/api/v1")
	userAuth := fakeAuthMiddleware(seed.userID)
	userGroup := api.Group("/user")
	userhandler.RegisterGiftRoutes(userGroup, giftService, itemService, userAuth)

	// 列出礼物应为空
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/gifts", nil, "")
	var listParsed apiResp[itemsvc.ServiceItemListResponse]
	_ = json.Unmarshal(listResp.Body.Bytes(), &listParsed)
	if listParsed.Data.Total != 0 || len(listParsed.Data.Items) != 0 {
		t.Fatalf("expected no active gifts, got %+v", listParsed.Data)
	}
}

func migrateGiftModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.ServiceItem{},
		&model.Order{},
		&model.CommissionRecord{},
	); err != nil {
		t.Fatalf("migrate gift models: %v", err)
	}
}

type giftSeed struct {
	userID     uint64
	playerID   uint64
	giftItemID uint64
}

func seedGiftData(t *testing.T, db *gorm.DB) giftSeed {
	t.Helper()
	ctx := context.Background()

	userRepo := user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	itemRepo := serviceitem.NewServiceItemRepository(db)

	userModel := &model.User{
		Name:         "GiftUser",
		Email:        "gift_user@example.com",
		Phone:        "20000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, userModel); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	playerUser := &model.User{
		Name:         "GiftPlayerUser",
		Email:        "gift_player@example.com",
		Phone:        "20000000002",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RolePlayer,
	}
	if err := userRepo.Create(ctx, playerUser); err != nil {
		t.Fatalf("seed player user: %v", err)
	}

	playerModel := &model.Player{
		UserID:          playerUser.ID,
		Nickname:        "GiftPro",
		HourlyRateCents: 3000,
	}
	if err := playerRepo.Create(ctx, playerModel); err != nil {
		t.Fatalf("seed player: %v", err)
	}

	// 创建礼物 ServiceItem
	itemService := itemsvc.NewServiceItemService(itemRepo, nil, playerRepo)
	giftItem, err := itemService.CreateServiceItem(ctx, itemsvc.CreateServiceItemRequest{
		ItemCode:       "GIFT-001",
		Name:           "Heart",
		SubCategory:    model.SubCategoryGift,
		BasePriceCents: 1000,
		CommissionRate: 0.2,
		IconURL:        "https://example.com/heart.png",
	})
	if err != nil {
		t.Fatalf("seed gift item: %v", err)
	}

	return giftSeed{
		userID:     userModel.ID,
		playerID:   playerModel.ID,
		giftItemID: giftItem.ID,
	}
}

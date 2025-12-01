package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/service/item"
	"gamelink/internal/testutil"
)

// 服务项管理：创建礼物与陪玩服务、查询过滤、更新/批量上下架与调价、删除后数量减少
func TestAdminServiceItemCRUDAndBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateServiceItemModels(t, db)

	gameRepo := game.NewGameRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	itemRepo := serviceitem.NewServiceItemRepository(db)
	svc := item.NewServiceItemService(itemRepo, gameRepo, playerRepo)

	// seed game
	g := &model.Game{Key: "lol", Name: "League"}
	if err := gameRepo.Create(ctx(), g); err != nil {
		t.Fatalf("seed game: %v", err)
	}

	router := gin.New()
	api := router.Group("/api/v1/admin")
	adminhandler.RegisterServiceItemRoutes(api, svc)

	// 1) 创建礼物
	giftPayload := map[string]interface{}{
		"itemCode":       "gift_rose",
		"name":           "Rose Gift",
		"subCategory":    "gift",
		"basePriceCents": 199,
		"serviceHours":   0,
		"minUsers":       1,
		"maxPlayers":     1,
		"commissionRate": 0.2,
		"sortOrder":      1,
	}
	giftResp := doJSON(router, http.MethodPost, "/api/v1/admin/service-items", giftPayload, "")
	if giftResp.Code != http.StatusOK {
		t.Fatalf("create gift status=%d body=%s", giftResp.Code, giftResp.Body.String())
	}
	var giftParsed apiResp[model.ServiceItem]
	_ = json.Unmarshal(giftResp.Body.Bytes(), &giftParsed)
	giftID := giftParsed.Data.ID

	// 2) 创建陪玩服务
	servicePayload := map[string]interface{}{
		"itemCode":       "escort_basic",
		"name":           "Basic Escort",
		"subCategory":    "solo",
		"gameId":         g.ID,
		"basePriceCents": 5000,
		"serviceHours":   2,
		"minUsers":       1,
		"maxPlayers":     1,
		"commissionRate": 0.2,
		"sortOrder":      2,
	}
	svcResp := doJSON(router, http.MethodPost, "/api/v1/admin/service-items", servicePayload, "")
	if svcResp.Code != http.StatusOK {
		t.Fatalf("create escort status=%d body=%s", svcResp.Code, svcResp.Body.String())
	}
	var escortParsed apiResp[model.ServiceItem]
	_ = json.Unmarshal(svcResp.Body.Bytes(), &escortParsed)
	escortID := escortParsed.Data.ID

	// 3) 礼物过滤查询
	listGiftResp := doJSON(router, http.MethodGet, "/api/v1/admin/service-items?subCategory=gift", nil, "")
	if listGiftResp.Code != http.StatusOK {
		t.Fatalf("list gift status=%d body=%s", listGiftResp.Code, listGiftResp.Body.String())
	}
	var listGiftParsed apiResp[item.ServiceItemListResponse]
	_ = json.Unmarshal(listGiftResp.Body.Bytes(), &listGiftParsed)
	if listGiftParsed.Data.Total != 1 {
		t.Fatalf("expected 1 gift, got %+v", listGiftParsed.Data)
	}

	// 4) 更新礼物价格
	updatePayload := map[string]interface{}{
		"basePriceCents": 299,
	}
	updResp := doJSON(router, http.MethodPut, "/api/v1/admin/service-items/"+uintToStr(giftID), updatePayload, "")
	if updResp.Code != http.StatusOK {
		t.Fatalf("update gift status=%d body=%s", updResp.Code, updResp.Body.String())
	}

	// 5) 批量下架两个服务
	batchStatus := map[string]interface{}{
		"ids":      []uint64{giftID, escortID},
		"isActive": false,
	}
	statusResp := doJSON(router, http.MethodPost, "/api/v1/admin/service-items/batch-update-status", batchStatus, "")
	if statusResp.Code != http.StatusOK {
		t.Fatalf("batch status status=%d body=%s", statusResp.Code, statusResp.Body.String())
	}

	// 6) 批量调价
	batchPrice := map[string]interface{}{
		"ids":            []uint64{giftID, escortID},
		"basePriceCents": 888,
	}
	priceResp := doJSON(router, http.MethodPost, "/api/v1/admin/service-items/batch-update-price", batchPrice, "")
	if priceResp.Code != http.StatusOK {
		t.Fatalf("batch price status=%d body=%s", priceResp.Code, priceResp.Body.String())
	}

	// 7) 删除礼物
	delResp := doJSON(router, http.MethodDelete, "/api/v1/admin/service-items/"+uintToStr(giftID), nil, "")
	if delResp.Code != http.StatusOK {
		t.Fatalf("delete gift status=%d body=%s", delResp.Code, delResp.Body.String())
	}

	// 8) 列表应只剩 1 个（escort）
	listAllResp := doJSON(router, http.MethodGet, "/api/v1/admin/service-items", nil, "")
	if listAllResp.Code != http.StatusOK {
		t.Fatalf("list all status=%d body=%s", listAllResp.Code, listAllResp.Body.String())
	}
	var listAllParsed apiResp[item.ServiceItemListResponse]
	_ = json.Unmarshal(listAllResp.Body.Bytes(), &listAllParsed)
	if listAllParsed.Data.Total != 1 || len(listAllParsed.Data.Items) != 1 || listAllParsed.Data.Items[0].ID != escortID {
		t.Fatalf("expected only escort remain, got %+v", listAllParsed.Data)
	}
}

// 礼物 service_hours 非 0 应校验失败
func TestServiceItemGiftInvalidServiceHours(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateServiceItemModels(t, db)

	gameRepo := game.NewGameRepository(db)
	playerRepo := player.NewPlayerRepository(db)
	itemRepo := serviceitem.NewServiceItemRepository(db)
	svc := item.NewServiceItemService(itemRepo, gameRepo, playerRepo)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	adminhandler.RegisterServiceItemRoutes(api, svc)

	payload := map[string]interface{}{
		"itemCode":       "gift_bad",
		"name":           "Bad Gift",
		"subCategory":    "gift",
		"basePriceCents": 100,
		"serviceHours":   1, // invalid for gift
		"minUsers":       1,
		"maxPlayers":     1,
		"commissionRate": 0.2,
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/service-items", payload, "")
	if resp.Code == http.StatusOK {
		t.Fatalf("expected failure for gift service_hours != 0, got 200")
	}
}

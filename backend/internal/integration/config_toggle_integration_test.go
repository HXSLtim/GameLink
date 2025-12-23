package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	"gamelink/internal/repository/commission"
	"gamelink/internal/repository/game"
	orderimpl "gamelink/internal/repository/implementations"
	"gamelink/internal/repository/order"
	"gamelink/internal/repository/ranking"
	"gamelink/internal/repository/user"

	commissionservice "gamelink/internal/service/commission"
	rankingservice "gamelink/internal/service/ranking"
	"gamelink/pkg/testutil"
)

// 配置禁用/更新：抽成规则、排名抽成配置
func TestConfigUpdateAndDisable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateConfigModels(t, db)

	ctx := context.Background()
	orderRepo := orderimpl.NewOrderRepository(db)
	_ = user.NewUserRepository(db)
	playerRepo := user.NewPlayerRepository(db)
	_ = game.NewGameRepository(db)
	_ = order.NewPaymentRepository(db)
	commissionRepo := commission.NewCommissionRepository(db)
	rankingRepo := ranking.NewRankingCommissionRepository(db)

	commissionSvc := commissionservice.NewCommissionService(commissionRepo, orderRepo, playerRepo)
	rankingSvc := rankingservice.NewRankingService(nil, rankingRepo, orderRepo)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminhandler.RegisterCommissionRoutes(adminGroup, commissionSvc, nil)
	adminhandler.RegisterRankingCommissionRoutes(adminGroup, rankingRepo)

	// 创建抽成规则
	createRulePayload := map[string]interface{}{
		"name":     "default",
		"type":     "default",
		"rate":     20,
		"isActive": true,
	}
	ruleResp := doJSON(router, http.MethodPost, "/api/v1/admin/commission/rules", createRulePayload, "")
	if ruleResp.Code != http.StatusCreated && ruleResp.Code != http.StatusOK {
		t.Fatalf("create rule status=%d body=%s", ruleResp.Code, ruleResp.Body.String())
	}
	var ruleParsed apiResp[model.CommissionRule]
	_ = json.Unmarshal(ruleResp.Body.Bytes(), &ruleParsed)
	ruleID := ruleParsed.Data.ID

	// 更新并禁用规则
	updateRulePayload := map[string]interface{}{
		"rate":     15,
		"isActive": false,
	}
	updateResp := doJSON(router, http.MethodPut, "/api/v1/admin/commission/rules/"+uintToStr(ruleID), updateRulePayload, "")
	if updateResp.Code != http.StatusOK {
		t.Fatalf("update rule status=%d body=%s", updateResp.Code, updateResp.Body.String())
	}
	updatedRule, _ := commissionRepo.GetRule(ctx, ruleID)
	if updatedRule.Rate != 15 || updatedRule.IsActive {
		t.Fatalf("rule not updated, got rate=%d active=%v", updatedRule.Rate, updatedRule.IsActive)
	}

	// 创建排名抽成配置
	createRankingPayload := map[string]interface{}{
		"name":        "cfg1",
		"rankingType": "income",
		"month":       "2025-01",
		"rules": []map[string]interface{}{
			{"rankStart": 1, "rankEnd": 3, "commissionRate": 5},
		},
	}
	rcResp := doJSON(router, http.MethodPost, "/api/v1/admin/ranking-commission/configs", createRankingPayload, "")
	if rcResp.Code != http.StatusCreated && rcResp.Code != http.StatusOK {
		t.Fatalf("create ranking cfg status=%d body=%s", rcResp.Code, rcResp.Body.String())
	}
	var rcParsed apiResp[model.RankingCommissionConfig]
	_ = json.Unmarshal(rcResp.Body.Bytes(), &rcParsed)
	cfgID := rcParsed.Data.ID

	// 更新禁用 ranking 配置
	updateCfgPayload := map[string]interface{}{
		"description": "updated desc",
		"isActive":    false,
	}
	cfgResp := doJSON(router, http.MethodPut, "/api/v1/admin/ranking-commission/configs/"+uintToStr(cfgID), updateCfgPayload, "")
	if cfgResp.Code != http.StatusOK {
		t.Fatalf("update cfg status=%d body=%s", cfgResp.Code, cfgResp.Body.String())
	}
	cfg, err := rankingRepo.GetConfig(ctx, cfgID)
	if err != nil {
		t.Fatalf("get cfg: %v", err)
	}
	if cfg.IsActive {
		t.Fatalf("expected cfg disabled")
	}
	if cfg.Description != "updated desc" {
		t.Fatalf("cfg desc not updated")
	}

	// 禁用配置后计算排名不应报错
	if err := rankingSvc.CalculateMonthlyRankings(ctx, "2025-01"); err != nil {
		t.Fatalf("calculate rankings after disable cfg: %v", err)
	}
}

func migrateConfigModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Player{},
		&model.Game{},
		&model.Order{},
		&model.Payment{},
		&model.CommissionRule{},
		&model.CommissionRecord{},
		&model.RankingCommissionConfig{},
	); err != nil {
		t.Fatalf("migrate config models: %v", err)
	}
}

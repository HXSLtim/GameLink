package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	notificationrepo "gamelink/internal/repository/content"
	userrepo "gamelink/internal/repository/user"
	userservice "gamelink/internal/service/user"
	"gamelink/pkg/testutil"
)

func fakeAdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userId", "1")        // userBatch.go expects string "userId"
		c.Set("user_id", uint64(1)) // standard convention
		c.Next()
	}
}

// 批量增加积分：创建多个用户 -> 批量增加积分 -> 验证余额
func TestAdminBatchAddPoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateBatchModels(t, db)

	// repos & service
	userRepo := userrepo.NewUserRepository(db)
	tagRepo := userrepo.NewUserTagRepository(db)
	notifyRepo := notificationrepo.NewNotificationRepository(db)
	batchSvc := userservice.NewBatchOperationService(db, userRepo, tagRepo, notifyRepo)

	router := gin.New()
	adminGroup := router.Group("/api/v1/admin")
	adminGroup.Use(fakeAdminAuthMiddleware()) // 模拟管理员ID=1
	// register batch routes
	adminhandler.RegisterBatchRoutes(adminGroup, batchSvc)

	// 创建测试用户
	ctx := context.Background()
	users := make([]*model.User, 3)
	for i := 0; i < 3; i++ {
		u := &model.User{
			Name:         "BatchUser" + uintToStr(uint64(i)),
			Email:        "batch" + uintToStr(uint64(i)) + "@example.com",
			Phone:        "1900000000" + uintToStr(uint64(i)),
			PasswordHash: "x",
			Status:       model.UserStatusActive,
			Role:         model.RoleUser,
		}
		if err := userRepo.Create(ctx, u); err != nil {
			t.Fatalf("create user %d: %v", i, err)
		}
		users[i] = u
		// seed wallet
		wallet := &model.Wallet{
			UserID:       u.ID,
			BalanceCents: 0,
			FrozenCents:  0,
		}
		if err := db.Create(wallet).Error; err != nil {
			t.Fatalf("create wallet %d: %v", i, err)
		}
	}

	// 批量增加积分
	payload := map[string]interface{}{
		"target":  "users",
		"userIds": []uint64{users[0].ID, users[1].ID, users[2].ID},
		"cents":   1000,
		"reason":  "test batch add",
		"type":    "admin",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/points", payload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("batch add points status=%d body=%s", resp.Code, resp.Body.String())
	}

	var batchResp adminhandler.BatchResponse
	if err := json.Unmarshal(resp.Body.Bytes(), &batchResp); err != nil {
		t.Fatalf("parse batch resp: %v", err)
	}
	if !batchResp.Success || batchResp.SuccessCount != 3 {
		fmt.Printf("BATCH_FAILURE: %s DETAILS: %+v\n", batchResp.Message, batchResp)
		t.FailNow()
	}

	// 验证每个用户的余额
	for _, u := range users {
		var wallet model.Wallet
		if err := db.Where("user_id = ?", u.ID).First(&wallet).Error; err != nil {
			t.Fatalf("get wallet for user %d: %v", u.ID, err)
		}
		if wallet.BalanceCents != 1000 {
			t.Fatalf("user %d balance mismatch: expected 1000, got %d", u.ID, wallet.BalanceCents)
		}
	}
}

func migrateBatchModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.NotificationEvent{},
		&model.OperationLog{},
	); err != nil {
		t.Fatalf("migrate batch models: %v", err)
	}
}

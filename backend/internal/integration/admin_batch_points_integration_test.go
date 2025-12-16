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

// 批量发送通知
func TestAdminBatchSendNotification(t *testing.T) {
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
	adminGroup.Use(fakeAdminAuthMiddleware())
	adminhandler.RegisterBatchRoutes(adminGroup, batchSvc)

	// 创建测试用户
	ctx := context.Background()
	users := make([]*model.User, 3)
	for i := 0; i < 3; i++ {
		u := &model.User{
			Name:         "NotifyUser" + uintToStr(uint64(i)),
			Email:        "notify" + uintToStr(uint64(i)) + "@example.com",
			Phone:        "1800000000" + uintToStr(uint64(i)),
			PasswordHash: "x",
			Status:       model.UserStatusActive,
			Role:         model.RoleUser,
		}
		if err := userRepo.Create(ctx, u); err != nil {
			t.Fatalf("create user %d: %v", i, err)
		}
		users[i] = u
	}

	// 批量发送通知
	payload := map[string]interface{}{
		"target":  "users",
		"userIds": []uint64{users[0].ID, users[1].ID, users[2].ID},
		"title":   "Batch Notification",
		"content": "Hello Batch Users",
		"type":    "system",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/notification", payload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("batch notify status=%d body=%s", resp.Code, resp.Body.String())
	}

	// 验证通知是否创建
	var count int64
	if err := db.Model(&model.NotificationEvent{}).Where("title = ?", "Batch Notification").Count(&count).Error; err != nil {
		t.Fatalf("count notifications: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 notifications, got %d", count)
	}
}

// 批量更新状态和角色
func TestAdminBatchUpdateStatusAndRole(t *testing.T) {
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
	adminGroup.Use(fakeAdminAuthMiddleware())
	adminhandler.RegisterBatchRoutes(adminGroup, batchSvc)

	// 创建测试用户
	ctx := context.Background()
	users := make([]*model.User, 2)
	for i := 0; i < 2; i++ {
		u := &model.User{
			Name:         "StatusUser" + uintToStr(uint64(i)),
			Email:        "status" + uintToStr(uint64(i)) + "@example.com",
			Phone:        "1700000000" + uintToStr(uint64(i)),
			PasswordHash: "x",
			Status:       model.UserStatusActive, // 初始Active
			Role:         model.RoleUser,         // 初始User
		}
		if err := userRepo.Create(ctx, u); err != nil {
			t.Fatalf("create user %d: %v", i, err)
		}
		users[i] = u
	}

	// 1. 批量更新状态 -> Banned
	statusPayload := map[string]interface{}{
		"userIds": []uint64{users[0].ID, users[1].ID},
		"status":  "banned",
		"reason":  "violation",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/status", statusPayload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("batch update status code=%d body=%s", resp.Code, resp.Body.String())
	}

	// Verify Status
	for _, u := range users {
		var updated model.User
		if err := db.First(&updated, u.ID).Error; err != nil {
			t.Fatalf("get user %d: %v", u.ID, err)
		}
		if updated.Status != model.UserStatusBanned {
			t.Fatalf("user %d status expected Banned, got %s", u.ID, updated.Status)
		}
	}

	// 2. 批量更新角色 -> Player
	rolePayload := map[string]interface{}{
		"userIds": []uint64{users[0].ID, users[1].ID},
		"role":    "player",
	}
	respRole := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/role", rolePayload, "")
	if respRole.Code != http.StatusOK {
		t.Fatalf("batch update role code=%d body=%s", respRole.Code, respRole.Body.String())
	}

	// Verify Role
	for _, u := range users {
		var updated model.User
		if err := db.First(&updated, u.ID).Error; err != nil {
			t.Fatalf("get user %d: %v", u.ID, err)
		}
		if updated.Role != model.RolePlayer {
			t.Fatalf("user %d role expected Player, got %s", u.ID, updated.Role)
		}
	}
}

// 批量删除用户
func TestAdminBatchDeleteUser(t *testing.T) {
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
	adminGroup.Use(fakeAdminAuthMiddleware())
	adminhandler.RegisterBatchRoutes(adminGroup, batchSvc)

	// 创建测试用户
	ctx := context.Background()
	u := &model.User{
		Name:         "DelUser",
		Email:        "del@example.com",
		Phone:        "16000000001",
		PasswordHash: "x",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// 批量删除
	delPayload := map[string]interface{}{
		"userIds": []uint64{u.ID},
		"reason":  "cleanup",
	}
	resp := doJSON(router, http.MethodPost, "/api/v1/admin/users/batch/delete", delPayload, "")
	if resp.Code != http.StatusOK {
		t.Fatalf("batch delete status=%d body=%s", resp.Code, resp.Body.String())
	}

	// Verify Soft Delete (should not be found with normal find)
	var found model.User
	err := db.First(&found, u.ID).Error
	if err != gorm.ErrRecordNotFound {
		t.Fatalf("expected ErrRecordNotFound, got %v (found: %+v)", err, found)
	}

	// Verify Unscoped (should be found)
	if err := db.Unscoped().First(&found, u.ID).Error; err != nil {
		t.Fatalf("expected found in unscoped, got %v", err)
	}
	if found.DeletedAt.Time.IsZero() {
		t.Fatalf("expected DeletedAt to be set")
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

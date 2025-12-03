package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	notificationhandler "gamelink/internal/handler/notification"
	"gamelink/internal/model"
	notificationrepo "gamelink/internal/repository/content"
	userrepo "gamelink/internal/repository/user"
	notificationservice "gamelink/internal/service/notification"
	"gamelink/pkg/testutil"
)

// 场景：种子通知 -> 列表 -> 未读数 -> 标记已读 -> 未读数减少
func TestNotificationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateNotificationModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	notifyRepo := notificationrepo.NewNotificationRepository(db)

	user := &model.User{
		Name:         "NotifyUser",
		Email:        "notify@example.com",
		Phone:        "40000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	// 种子两条未读通知
	ev1 := &model.NotificationEvent{UserID: user.ID, Title: "A", Message: "m1", Priority: model.NotificationPriorityNormal}
	ev2 := &model.NotificationEvent{UserID: user.ID, Title: "B", Message: "m2", Priority: model.NotificationPriorityHigh}
	if err := db.WithContext(ctx).Create(ev1).Error; err != nil {
		t.Fatalf("seed notif1: %v", err)
	}
	if err := db.WithContext(ctx).Create(ev2).Error; err != nil {
		t.Fatalf("seed notif2: %v", err)
	}

	svc := notificationservice.NewService(notifyRepo)
	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(user.ID)
	notificationhandler.RegisterRoutes(api, svc, auth)

	// 列表
	listResp := doJSON(router, http.MethodGet, "/api/v1/notifications", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list notif status=%d body=%s", listResp.Code, listResp.Body.String())
	}
	var listParsed apiResp[*notificationservice.ListResponse]
	if err := json.Unmarshal(listResp.Body.Bytes(), &listParsed); err != nil {
		t.Fatalf("parse list: %v", err)
	}
	if len(listParsed.Data.Items) != 2 || listParsed.Data.UnreadCount != 2 {
		t.Fatalf("unexpected list resp: %+v", listParsed.Data)
	}

	// 未读数
	unreadResp := doJSON(router, http.MethodGet, "/api/v1/notifications/unread-count", nil, "")
	if unreadResp.Code != http.StatusOK {
		t.Fatalf("unread count status=%d body=%s", unreadResp.Code, unreadResp.Body.String())
	}

	// 标记已读
	readPayload := map[string]interface{}{
		"ids": []uint64{ev1.ID},
	}
	markResp := doJSON(router, http.MethodPost, "/api/v1/notifications/read", readPayload, "")
	if markResp.Code != http.StatusOK {
		t.Fatalf("mark read status=%d body=%s", markResp.Code, markResp.Body.String())
	}

	// 再查未读数应为1
	listResp2 := doJSON(router, http.MethodGet, "/api/v1/notifications?unread=true", nil, "")
	if listResp2.Code != http.StatusOK {
		t.Fatalf("list unread status=%d body=%s", listResp2.Code, listResp2.Body.String())
	}
	if err := json.Unmarshal(listResp2.Body.Bytes(), &listParsed); err != nil {
		t.Fatalf("parse list2: %v", err)
	}
	if listParsed.Data.UnreadCount != 1 {
		t.Fatalf("expected unread=1, got %+v", listParsed.Data)
	}
}

func migrateNotificationModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.NotificationEvent{},
	); err != nil {
		t.Fatalf("migrate notification models: %v", err)
	}
}

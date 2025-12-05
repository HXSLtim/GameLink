package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	userfeed "gamelink/internal/handler/user"
	"gamelink/internal/model"
	feedrepo "gamelink/internal/repository/content"
	userrepo "gamelink/internal/repository/user"
	feedservice "gamelink/internal/service/content"
	"gamelink/pkg/testutil"
)

// 动态：发布 -> 列表 -> 举报
func TestFeedFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateFeedModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	feedRepo := feedrepo.NewFeedRepository(db)

	user := &model.User{
		Name:         "FeedUser",
		Email:        "feed@example.com",
		Phone:        "70000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	svc := feedservice.NewFeedService(feedRepo, nil)

	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(user.ID)
	userGroup := api.Group("/user")
	userfeed.RegisterFeedRoutes(userGroup, svc, auth)

	// 发布动态
	createPayload := map[string]interface{}{
		"content":    "hello feed",
		"visibility": "public",
		"images": []map[string]interface{}{
			{"url": "https://example.com/a.jpg", "width": 100, "height": 100, "sizeBytes": 1024, "order": 1},
		},
	}
	createResp := doJSON(router, http.MethodPost, "/api/v1/user/feeds", createPayload, "")
	if createResp.Code != http.StatusOK {
		t.Fatalf("create feed status=%d body=%s", createResp.Code, createResp.Body.String())
	}
	var createParsed apiResp[feedservice.FeedView]
	if err := json.Unmarshal(createResp.Body.Bytes(), &createParsed); err != nil {
		t.Fatalf("parse create feed: %v", err)
	}
	feedID := createParsed.Data.ID

	// 列表
	listResp := doJSON(router, http.MethodGet, "/api/v1/user/feeds", nil, "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list feeds status=%d body=%s", listResp.Code, listResp.Body.String())
	}

	// 举报
	reportPayload := map[string]interface{}{
		"reason": "abuse report",
	}
	reportResp := doJSON(router, http.MethodPost, "/api/v1/user/feeds/"+uintToStr(feedID)+"/report", reportPayload, "")
	if reportResp.Code != http.StatusOK {
		t.Fatalf("report feed status=%d body=%s", reportResp.Code, reportResp.Body.String())
	}
}

func migrateFeedModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.Feed{},
		&model.FeedImage{},
		&model.FeedReport{},
	); err != nil {
		t.Fatalf("migrate feed models: %v", err)
	}
}

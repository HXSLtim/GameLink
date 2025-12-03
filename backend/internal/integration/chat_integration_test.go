package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/pkg/cache"
	userchat "gamelink/internal/handler/user"
	"gamelink/internal/model"
	chatrepo "gamelink/internal/repository/chat"
	userrepo "gamelink/internal/repository/user"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/testutil"
)

// IM 集成：列群组 -> 发消息 -> 拉取消息 -> 举报消息
func TestChatFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateChatModels(t, db)

	ctx := context.Background()
	userRepo := userrepo.NewUserRepository(db)
	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)

	user := &model.User{
		Name:         "ChatUser",
		Email:        "chat@example.com",
		Phone:        "50000000001",
		PasswordHash: "hashed",
		Status:       model.UserStatusActive,
		Role:         model.RoleUser,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	group := &model.ChatGroup{
		GroupName: "OrderChat",
		GroupType: model.ChatGroupTypeOrder,
		CreatedBy: user.ID,
		IsActive:  true,
	}
	if err := db.WithContext(ctx).Create(group).Error; err != nil {
		t.Fatalf("seed group: %v", err)
	}
	member := &model.ChatGroupMember{
		GroupID:  group.ID,
		UserID:   user.ID,
		JoinedAt: time.Now(),
		IsActive: true,
	}
	if err := db.WithContext(ctx).Create(member).Error; err != nil {
		t.Fatalf("seed member: %v", err)
	}

	svc := chatservice.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, cache.NewMemory())
	router := gin.New()
	api := router.Group("/api/v1")
	auth := fakeAuthMiddleware(user.ID)
	userchat.RegisterChatRoutes(api, svc, auth)

	// 1) 列群组
	groupsResp := doJSON(router, http.MethodGet, "/api/v1/chat/groups", nil, "")
	if groupsResp.Code != http.StatusOK {
		t.Fatalf("list groups status=%d body=%s", groupsResp.Code, groupsResp.Body.String())
	}

	// 2) 发消息
	sendPayload := map[string]interface{}{
		"content": "hello world",
	}
	sendResp := doJSON(router, http.MethodPost, "/api/v1/chat/groups/"+uintToStr(group.ID)+"/messages", sendPayload, "")
	if sendResp.Code != http.StatusCreated {
		t.Fatalf("send message status=%d body=%s", sendResp.Code, sendResp.Body.String())
	}
	var sendParsed apiResp[model.ChatMessage]
	if err := json.Unmarshal(sendResp.Body.Bytes(), &sendParsed); err != nil {
		t.Fatalf("parse send resp: %v", err)
	}
	messageID := sendParsed.Data.ID

	// 3) 拉取消息
	msgResp := doJSON(router, http.MethodGet, "/api/v1/chat/groups/"+uintToStr(group.ID)+"/messages", nil, "")
	if msgResp.Code != http.StatusOK {
		t.Fatalf("list messages status=%d body=%s", msgResp.Code, msgResp.Body.String())
	}

	// 4) 举报消息
	reportPayload := map[string]interface{}{
		"reason": "test",
	}
	reportResp := doJSON(router, http.MethodPost, "/api/v1/chat/messages/"+uintToStr(messageID)+"/report", reportPayload, "")
	if reportResp.Code != http.StatusOK {
		t.Fatalf("report message status=%d body=%s", reportResp.Code, reportResp.Body.String())
	}
}

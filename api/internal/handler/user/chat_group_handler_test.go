package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	chatrepo "gamelink/internal/repository/chat"
	userrepo "gamelink/internal/repository/user"
	chatservice "gamelink/internal/service/chat"
	"gamelink/pkg/cache"
)

type userChatGroupTestContext struct {
	Router  *gin.Engine
	DB      *gorm.DB
	User    *model.User
	Target  *model.User
	Service *chatservice.ChatService
	Members repository.ChatMemberRepository
}

func setupUserChatGroupTest(t *testing.T) *userChatGroupTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)
	router := testutil.SetupGinTest(t)

	user := testutil.CreateAdminUser(t, db, model.RoleUser)
	target := testutil.CreateAdminUser(t, db, model.RolePlayer)

	groupRepo := chatrepo.NewChatGroupRepository(db)
	memberRepo := chatrepo.NewChatMemberRepository(db)
	messageRepo := chatrepo.NewChatMessageRepository(db)
	reportRepo := chatrepo.NewChatReportRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	memCache := cache.NewMemory()

	svc := chatservice.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, userRepo, memCache)

	return &userChatGroupTestContext{
		Router:  router,
		DB:      db,
		User:    user,
		Target:  target,
		Service: svc,
		Members: memberRepo,
	}
}

func (ctx *userChatGroupTestContext) registerRoutes() {
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", ctx.User.ID)
		c.Next()
	}
	group := ctx.Router.Group("/user")
	RegisterChatRoutes(group, ctx.Service, authMiddleware)
}

func TestUserChat_CreateAndGetGroup(t *testing.T) {
	ctx := setupUserChatGroupTest(t)
	ctx.registerRoutes()

	payload := map[string]interface{}{
		"targetUserId": ctx.Target.ID,
		"groupType":    "private",
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/user/chat/groups", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)
	testutil.AssertSuccess(t, w)

	var createResp model.APIResponse[model.ChatGroup]
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	assert.Equal(t, model.ChatGroupTypePrivate, createResp.Data.GroupType)
	assert.Len(t, createResp.Data.Members, 2)

	getPath := testutil.BuildPath("/user/chat/groups/:id", map[string]string{
		"id": testutil.Uint64ToStr(createResp.Data.ID),
	})
	getResp := testutil.MakeRequest(t, ctx.Router, http.MethodGet, getPath, nil)
	testutil.AssertSuccess(t, getResp)
}

func TestUserChat_JoinReadLeave(t *testing.T) {
	ctx := setupUserChatGroupTest(t)
	ctx.registerRoutes()

	group := &model.ChatGroup{
		GroupName: "Public Room",
		GroupType: model.ChatGroupTypePublic,
		CreatedBy: ctx.Target.ID,
		IsActive:  true,
	}
	require.NoError(t, ctx.DB.Create(group).Error)

	joinPath := testutil.BuildPath("/user/chat/groups/:id/join", map[string]string{
		"id": testutil.Uint64ToStr(group.ID),
	})
	joinResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, joinPath, nil)
	testutil.AssertSuccess(t, joinResp)

	member, err := ctx.Members.Get(context.Background(), group.ID, ctx.User.ID)
	require.NoError(t, err)
	assert.True(t, member.IsActive)

	readPayload := map[string]interface{}{
		"messageId": 123,
	}
	readPath := testutil.BuildPath("/user/chat/groups/:id/read", map[string]string{
		"id": testutil.Uint64ToStr(group.ID),
	})
	readResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, readPath, readPayload)
	testutil.AssertSuccess(t, readResp)

	member, err = ctx.Members.Get(context.Background(), group.ID, ctx.User.ID)
	require.NoError(t, err)
	require.NotNil(t, member.LastReadMessageID)
	assert.Equal(t, uint64(123), *member.LastReadMessageID)
	require.NotNil(t, member.LastReadAt)
	assert.WithinDuration(t, time.Now(), *member.LastReadAt, time.Minute)

	leavePath := testutil.BuildPath("/user/chat/groups/:id/leave", map[string]string{
		"id": testutil.Uint64ToStr(group.ID),
	})
	leaveResp := testutil.MakeRequest(t, ctx.Router, http.MethodPost, leavePath, nil)
	testutil.AssertSuccess(t, leaveResp)

	member, err = ctx.Members.Get(context.Background(), group.ID, ctx.User.ID)
	require.NoError(t, err)
	assert.False(t, member.IsActive)
}

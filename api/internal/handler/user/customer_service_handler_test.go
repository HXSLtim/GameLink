package user

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	conversationservice "gamelink/internal/service/conversation"
)

type csFakeRepo struct {
	group    *model.ChatGroup
	messages []model.ChatMessage
}

func (r *csFakeRepo) Create(ctx context.Context, group *model.ChatGroup, members []*model.ChatGroupMember) error {
	_ = ctx
	if group.ID == 0 {
		group.ID = 1
	}
	group.Members = make([]model.ChatGroupMember, 0, len(members))
	for _, member := range members {
		if member == nil {
			continue
		}
		member.GroupID = group.ID
		group.Members = append(group.Members, *member)
	}
	group.IsActive = true
	r.group = group
	return nil
}

func (r *csFakeRepo) Get(ctx context.Context, conversationID uint64) (*model.ChatGroup, error) {
	_ = ctx
	if r.group == nil || r.group.ID != conversationID {
		return nil, assert.AnError
	}
	return r.group, nil
}

func (r *csFakeRepo) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	_ = ctx
	_ = page
	_ = pageSize
	if r.group == nil {
		return []model.ChatGroup{}, 0, nil
	}
	for _, member := range r.group.Members {
		if member.UserID == userID {
			return []model.ChatGroup{*r.group}, 1, nil
		}
	}
	return []model.ChatGroup{}, 0, nil
}

func (r *csFakeRepo) Update(ctx context.Context, group *model.ChatGroup) error {
	_ = ctx
	r.group = group
	return nil
}

func (r *csFakeRepo) Delete(ctx context.Context, conversationID uint64) error {
	_ = ctx
	if r.group != nil && r.group.ID == conversationID {
		r.group.IsActive = false
	}
	return nil
}

func (r *csFakeRepo) FindActivePrivateByParticipants(ctx context.Context, userID uint64, agentIDs []uint64) (*model.ChatGroup, error) {
	_ = ctx
	if r.group == nil || !r.group.IsActive {
		return nil, nil
	}
	hasUser := false
	hasAgent := false
	agentSet := make(map[uint64]struct{}, len(agentIDs))
	for _, agentID := range agentIDs {
		agentSet[agentID] = struct{}{}
	}
	for _, member := range r.group.Members {
		if member.UserID == userID {
			hasUser = true
		}
		if _, ok := agentSet[member.UserID]; ok {
			hasAgent = true
		}
	}
	if hasUser && hasAgent {
		return r.group, nil
	}
	return nil, nil
}

func (r *csFakeRepo) ListMessages(ctx context.Context, groupID uint64, page, pageSize int, beforeID *uint64) ([]model.ChatMessage, int64, error) {
	_ = ctx
	_ = page
	_ = pageSize
	_ = beforeID
	filtered := make([]model.ChatMessage, 0, len(r.messages))
	for _, message := range r.messages {
		if message.GroupID == groupID {
			filtered = append(filtered, message)
		}
	}
	return filtered, int64(len(filtered)), nil
}

func (r *csFakeRepo) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	_ = ctx
	message.ID = uint64(len(r.messages) + 1)
	message.CreatedAt = time.Now()
	r.messages = append(r.messages, *message)
	return nil
}

func (r *csFakeRepo) CountActiveByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error) {
	_ = ctx
	result := make(map[uint64]int64, len(agentIDs))
	for _, agentID := range agentIDs {
		result[agentID] = 0
	}
	return result, nil
}

func (r *csFakeRepo) IsMember(ctx context.Context, groupID, userID uint64) (bool, error) {
	_ = ctx
	if r.group == nil || r.group.ID != groupID {
		return false, nil
	}
	for _, member := range r.group.Members {
		if member.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

type csFakeUsers struct {
	usersByID    map[uint64]*model.User
	usersByEmail map[string]*model.User
}

func (u *csFakeUsers) Get(ctx context.Context, id uint64) (*model.User, error) {
	_ = ctx
	user, ok := u.usersByID[id]
	if !ok {
		return nil, assert.AnError
	}
	return user, nil
}

func (u *csFakeUsers) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	_ = ctx
	user, ok := u.usersByEmail[email]
	if !ok {
		return nil, assert.AnError
	}
	return user, nil
}

type csFakeRoles struct {
	role  *model.RoleModel
	users []uint64
}

func (r *csFakeRoles) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	_ = ctx
	_ = slug
	return r.role, nil
}

func (r *csFakeRoles) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	_ = ctx
	_ = roleID
	return r.users, nil
}

func setupCustomerServiceRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	repo := &csFakeRepo{}
	users := &csFakeUsers{
		usersByID: map[uint64]*model.User{
			2001: {
				Base:   model.Base{ID: 2001},
				Name:   "客服A",
				Status: model.UserStatusActive,
			},
		},
		usersByEmail: map[string]*model.User{
			"cs.agent@gamelink.com": {
				Base:   model.Base{ID: 2001},
				Name:   "客服A",
				Status: model.UserStatusActive,
			},
		},
	}
	roles := &csFakeRoles{
		role:  &model.RoleModel{Base: model.Base{ID: 1}},
		users: []uint64{2001},
	}

	svc := conversationservice.NewService(repo, users, roles, nil)

	auth := func(c *gin.Context) {
		c.Set("user_id", uint64(1001))
		c.Next()
	}

	router := gin.New()
	userGroup := router.Group("/user")
	RegisterCustomerServiceRoutes(userGroup, svc, auth)
	return router
}

func TestCustomerService_LegacyEndpoints(t *testing.T) {
	router := setupCustomerServiceRouter(t)

	sessionResp := testutil.MakeRequest(t, router, http.MethodGet, "/user/customer-service/session", nil)
	assert.Equal(t, http.StatusOK, sessionResp.Code)

	sendResp := testutil.MakeRequest(t, router, http.MethodPost, "/user/customer-service/messages", map[string]string{
		"content": "你好，客服",
	})
	assert.Equal(t, http.StatusCreated, sendResp.Code)

	listResp := testutil.MakeRequest(t, router, http.MethodGet, "/user/customer-service/messages?page=1&pageSize=20", nil)
	assert.Equal(t, http.StatusOK, listResp.Code)
}

func TestCustomerService_ConversationEndpoints(t *testing.T) {
	router := setupCustomerServiceRouter(t)

	createResp := testutil.MakeRequest(t, router, http.MethodPost, "/user/customer-service/conversations", map[string]string{
		"content": "我要咨询订单问题",
	})
	assert.Equal(t, http.StatusCreated, createResp.Code)

	listResp := testutil.MakeRequest(t, router, http.MethodGet, "/user/customer-service/conversations?page=1&pageSize=20", nil)
	assert.Equal(t, http.StatusOK, listResp.Code)

	messageListResp := testutil.MakeRequest(t, router, http.MethodGet, "/user/customer-service/conversations/1/messages?page=1&pageSize=20", nil)
	assert.Equal(t, http.StatusOK, messageListResp.Code)

	sendResp := testutil.MakeRequest(t, router, http.MethodPost, "/user/customer-service/conversations/1/messages", map[string]string{
		"content": "继续补充描述",
	})
	assert.Equal(t, http.StatusCreated, sendResp.Code)

	closeResp := testutil.MakeRequest(t, router, http.MethodDelete, "/user/customer-service/conversations/1", nil)
	assert.Equal(t, http.StatusOK, closeResp.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(closeResp.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusOK), body["code"])
}

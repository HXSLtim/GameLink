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
	"gorm.io/gorm"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	chatservice "gamelink/internal/service/chat"
)

// ---- Mock Repositories ----

type mockChatGroupRepo struct {
	groups map[uint64]*model.ChatGroup
}

func newMockChatGroupRepo() *mockChatGroupRepo {
	return &mockChatGroupRepo{
		groups: make(map[uint64]*model.ChatGroup),
	}
}

func (m *mockChatGroupRepo) Create(ctx context.Context, group *model.ChatGroup) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockChatGroupRepo) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	if g, ok := m.groups[id]; ok {
		return g, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockChatGroupRepo) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	return nil, repository.ErrNotFound
}

func (m *mockChatGroupRepo) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	var result []model.ChatGroup
	for _, g := range m.groups {
		result = append(result, *g)
	}
	return result, int64(len(result)), nil
}

func (m *mockChatGroupRepo) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	return nil, 0, nil
}

func (m *mockChatGroupRepo) Update(ctx context.Context, group *model.ChatGroup) error {
	m.groups[group.ID] = group
	return nil
}

func (m *mockChatGroupRepo) Deactivate(ctx context.Context, id uint64) error {
	if g, ok := m.groups[id]; ok {
		g.IsActive = false
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockChatGroupRepo) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
	return nil, nil
}

func (m *mockChatGroupRepo) DeleteByIDs(ctx context.Context, ids []uint64) error {
	for _, id := range ids {
		delete(m.groups, id)
	}
	return nil
}

type mockChatMemberRepo struct {
	members map[string]*model.ChatGroupMember
}

func newMockChatMemberRepo() *mockChatMemberRepo {
	return &mockChatMemberRepo{
		members: make(map[string]*model.ChatGroupMember),
	}
}

func (m *mockChatMemberRepo) Add(ctx context.Context, member *model.ChatGroupMember) error {
	key := memberKey(member.GroupID, member.UserID)
	m.members[key] = member
	return nil
}

func (m *mockChatMemberRepo) AddBatch(ctx context.Context, members []*model.ChatGroupMember) error {
	for _, member := range members {
		_ = m.Add(ctx, member)
	}
	return nil
}

func (m *mockChatMemberRepo) Update(ctx context.Context, member *model.ChatGroupMember) error {
	key := memberKey(member.GroupID, member.UserID)
	m.members[key] = member
	return nil
}

func (m *mockChatMemberRepo) Remove(ctx context.Context, groupID, userID uint64) error {
	key := memberKey(groupID, userID)
	delete(m.members, key)
	return nil
}

func (m *mockChatMemberRepo) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
	key := memberKey(groupID, userID)
	if member, ok := m.members[key]; ok {
		return member, nil
	}
	return nil, repository.ErrNotFound
}

func memberKey(groupID, userID uint64) string {
	return "g" + string(rune(groupID)) + "u" + string(rune(userID))
}

type mockChatMessageRepo struct {
	messages map[uint64]*model.ChatMessage
	nextID   uint64
}

func newMockChatMessageRepo() *mockChatMessageRepo {
	return &mockChatMessageRepo{
		messages: make(map[uint64]*model.ChatMessage),
		nextID:   1,
	}
}

func (m *mockChatMessageRepo) Create(ctx context.Context, message *model.ChatMessage) error {
	message.ID = m.nextID
	m.nextID++
	m.messages[message.ID] = message
	return nil
}

func (m *mockChatMessageRepo) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error {
	for _, msg := range messages {
		_ = m.Create(ctx, msg)
	}
	return nil
}

func (m *mockChatMessageRepo) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) {
	var result []model.ChatMessage
	for _, msg := range m.messages {
		if msg.GroupID == opts.GroupID {
			result = append(result, *msg)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockChatMessageRepo) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) {
	if msg, ok := m.messages[id]; ok {
		return msg, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockChatMessageRepo) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error {
	if msg, ok := m.messages[id]; ok {
		now := time.Now()
		msg.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockChatMessageRepo) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) {
	return nil, 0, nil
}

func (m *mockChatMessageRepo) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error {
	if msg, ok := m.messages[id]; ok {
		msg.AuditStatus = status
		return nil
	}
	return repository.ErrNotFound
}

func (m *mockChatMessageRepo) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
	for id, msg := range m.messages {
		for _, gid := range groupIDs {
			if msg.GroupID == gid {
				delete(m.messages, id)
			}
		}
	}
	return nil
}

type mockChatReportRepo struct {
	reports map[uint64]*model.ChatReport
	nextID  uint64
}

func newMockChatReportRepo() *mockChatReportRepo {
	return &mockChatReportRepo{
		reports: make(map[uint64]*model.ChatReport),
		nextID:  1,
	}
}

func (m *mockChatReportRepo) Create(ctx context.Context, report *model.ChatReport) error {
	report.ID = m.nextID
	m.nextID++
	m.reports[report.ID] = report
	return nil
}

func (m *mockChatReportRepo) Get(ctx context.Context, id uint64) (*model.ChatReport, error) {
	if r, ok := m.reports[id]; ok {
		return r, nil
	}
	return nil, repository.ErrNotFound
}

func (m *mockChatReportRepo) Update(ctx context.Context, report *model.ChatReport) error {
	m.reports[report.ID] = report
	return nil
}

func (m *mockChatReportRepo) List(ctx context.Context, opts repository.ChatReportListOptions) ([]model.ChatReport, int64, error) {
	var result []model.ChatReport
	for _, r := range m.reports {
		result = append(result, *r)
	}
	return result, int64(len(result)), nil
}

// ---- Test Helpers ----

func setupChatTest() (*chatservice.ChatService, *mockChatGroupRepo, *mockChatMemberRepo, *mockChatMessageRepo, *mockChatReportRepo) {
	groupRepo := newMockChatGroupRepo()
	memberRepo := newMockChatMemberRepo()
	messageRepo := newMockChatMessageRepo()
	reportRepo := newMockChatReportRepo()
	cacheRepo := cache.NewMemory()

	svc := chatservice.NewChatService(groupRepo, memberRepo, messageRepo, reportRepo, cacheRepo)
	return svc, groupRepo, memberRepo, messageRepo, reportRepo
}

func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Set("user_id", uint64(1))
	return c, w
}

// ---- Tests ----

func TestListChatGroupsHandler_Success(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=20", nil)

	listChatGroupsHandler(c, svc)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestListChatGroupsHandler_DefaultPagination(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/", nil)

	listChatGroupsHandler(c, svc)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestListChatMessagesHandler_Success(t *testing.T) {
	svc, groupRepo, memberRepo, messageRepo, _ := setupChatTest()

	// Setup test data
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})
	messageRepo.Create(context.Background(), &model.ChatMessage{
		GroupID:     1,
		SenderID:    1,
		Content:     "Hello",
		MessageType: model.ChatMessageTypeText,
		AuditStatus: model.ChatMessageAuditApproved,
	})

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=50", nil)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	listChatMessagesHandler(c, svc)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestListChatMessagesHandler_InvalidGroupID(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	listChatMessagesHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListChatMessagesHandler_NotMember(t *testing.T) {
	svc, groupRepo, _, _, _ := setupChatTest()

	// Setup test data - group exists but user is not a member
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=50", nil)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	listChatMessagesHandler(c, svc)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestListChatMessagesHandler_InactiveGroup(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data - group is inactive
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  false,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: false,
	})

	c, w := createTestContext()
	c.Request = httptest.NewRequest("GET", "/?page=1&pageSize=50", nil)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	listChatMessagesHandler(c, svc)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSendChatMessageHandler_Success(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})

	c, w := createTestContext()
	payload := sendMessageRequest{
		Content:     "Hello World",
		MessageType: "text",
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp model.APIResponse[*model.ChatMessage]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Data)
}

func TestSendChatMessageHandler_InvalidGroupID(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	payload := sendMessageRequest{Content: "Hello"}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendChatMessageHandler_InvalidJSON(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendChatMessageHandler_NotMember(t *testing.T) {
	svc, groupRepo, _, _, _ := setupChatTest()

	// Setup test data - group exists but user is not a member
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)

	c, w := createTestContext()
	payload := sendMessageRequest{Content: "Hello"}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestSendChatMessageHandler_InactiveGroup(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data - group is inactive
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  false,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})

	c, w := createTestContext()
	payload := sendMessageRequest{Content: "Hello"}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusGone, w.Code)
}

func TestSendChatMessageHandler_EmptyContent(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})

	c, w := createTestContext()
	payload := sendMessageRequest{Content: ""}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendChatMessageHandler_UnsupportedMessageType(t *testing.T) {
	svc, groupRepo, memberRepo, _, _ := setupChatTest()

	// Setup test data
	group := &model.ChatGroup{
		Base:      model.Base{ID: 1},
		GroupName: "Test Group",
		IsActive:  true,
		GroupType: model.ChatGroupTypeOrder,
	}
	groupRepo.Create(context.Background(), group)
	memberRepo.Add(context.Background(), &model.ChatGroupMember{
		GroupID:  1,
		UserID:   1,
		IsActive: true,
	})

	c, w := createTestContext()
	payload := sendMessageRequest{
		Content:     "Hello",
		MessageType: "unsupported",
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	sendChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportChatMessageHandler_Success(t *testing.T) {
	svc, _, _, messageRepo, reportRepo := setupChatTest()

	// Setup test data
	messageRepo.Create(context.Background(), &model.ChatMessage{
		Base:        model.Base{ID: 1, DeletedAt: gorm.DeletedAt{}},
		GroupID:     1,
		Content:     "Offensive content",
		AuditStatus: model.ChatMessageAuditApproved,
	})

	c, w := createTestContext()
	payload := reportMessageRequest{
		Reason:   "offensive",
		Evidence: "screenshot",
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	reportChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.APIResponse[any]
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify report was created
	reports, _, _ := reportRepo.List(context.Background(), repository.ChatReportListOptions{})
	assert.Equal(t, 1, len(reports))
}

func TestReportChatMessageHandler_InvalidMessageID(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	payload := reportMessageRequest{Reason: "offensive"}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	reportChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReportChatMessageHandler_InvalidJSON(t *testing.T) {
	svc, _, _, _, _ := setupChatTest()

	c, w := createTestContext()
	c.Request = httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	reportChatMessageHandler(c, svc)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestParseUintFromParam_Success(t *testing.T) {
	c, _ := createTestContext()
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	result, err := parseUintFromParam(c, "id")

	assert.NoError(t, err)
	assert.Equal(t, uint64(123), result)
}

func TestParseUintFromParam_InvalidValue(t *testing.T) {
	c, _ := createTestContext()
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	_, err := parseUintFromParam(c, "id")

	assert.Error(t, err)
}

func TestParseUintFromParam_NegativeValue(t *testing.T) {
	c, _ := createTestContext()
	c.Params = []gin.Param{{Key: "id", Value: "-1"}}

	_, err := parseUintFromParam(c, "id")

	assert.Error(t, err)
}

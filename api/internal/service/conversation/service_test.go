package conversation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	repoerr "gamelink/internal/repository"
)

type mockConversationRepository struct {
	findConversationResult *model.ChatGroup
	findConversationErr    error
	countByAgentResult     map[uint64]int64
	countByAgentErr        error
	createErr              error
	getResult              *model.ChatGroup
	getErr                 error
	createMessageErr       error
	listByUserResult       []model.ChatGroup
	listByUserTotal        int64
	listByUserErr          error
	listMessagesResult     []model.ChatMessage
	listMessagesTotal      int64
	listMessagesErr        error
	isMember               bool
	isMemberErr            error
	deleteErr              error

	lastCreatedGroup   *model.ChatGroup
	lastCreatedMembers []*model.ChatGroupMember
	lastCreatedMessage *model.ChatMessage
}

func (m *mockConversationRepository) Create(ctx context.Context, group *model.ChatGroup, members []*model.ChatGroupMember) error {
	if m.createErr != nil {
		return m.createErr
	}
	if group.ID == 0 {
		group.ID = 88
	}
	m.lastCreatedGroup = group
	cloned := make([]*model.ChatGroupMember, 0, len(members))
	for _, item := range members {
		copied := *item
		copied.GroupID = group.ID
		cloned = append(cloned, &copied)
	}
	m.lastCreatedMembers = cloned
	return nil
}

func (m *mockConversationRepository) Get(ctx context.Context, conversationID uint64) (*model.ChatGroup, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.getResult != nil {
		return m.getResult, nil
	}
	if m.lastCreatedGroup == nil {
		return nil, repoerr.ErrNotFound
	}
	groupCopy := *m.lastCreatedGroup
	groupCopy.Members = make([]model.ChatGroupMember, 0, len(m.lastCreatedMembers))
	for _, member := range m.lastCreatedMembers {
		groupCopy.Members = append(groupCopy.Members, *member)
	}
	return &groupCopy, nil
}

func (m *mockConversationRepository) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	if m.listByUserErr != nil {
		return nil, 0, m.listByUserErr
	}
	return m.listByUserResult, m.listByUserTotal, nil
}

func (m *mockConversationRepository) Update(ctx context.Context, group *model.ChatGroup) error {
	m.lastCreatedGroup = group
	return nil
}

func (m *mockConversationRepository) Delete(ctx context.Context, conversationID uint64) error {
	return m.deleteErr
}

func (m *mockConversationRepository) FindActivePrivateByParticipants(ctx context.Context, userID uint64, agentIDs []uint64) (*model.ChatGroup, error) {
	if m.findConversationErr != nil {
		return nil, m.findConversationErr
	}
	return m.findConversationResult, nil
}

func (m *mockConversationRepository) ListMessages(ctx context.Context, groupID uint64, page, pageSize int, beforeID *uint64) ([]model.ChatMessage, int64, error) {
	if m.listMessagesErr != nil {
		return nil, 0, m.listMessagesErr
	}
	return m.listMessagesResult, m.listMessagesTotal, nil
}

func (m *mockConversationRepository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	if m.createMessageErr != nil {
		return m.createMessageErr
	}
	if message.ID == 0 {
		message.ID = 9001
	}
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}
	m.lastCreatedMessage = message
	return nil
}

func (m *mockConversationRepository) CountActiveByAgentIDs(ctx context.Context, agentIDs []uint64) (map[uint64]int64, error) {
	if m.countByAgentErr != nil {
		return nil, m.countByAgentErr
	}
	if m.countByAgentResult == nil {
		return map[uint64]int64{}, nil
	}
	return m.countByAgentResult, nil
}

func (m *mockConversationRepository) IsMember(ctx context.Context, groupID, userID uint64) (bool, error) {
	return m.isMember, m.isMemberErr
}

type mockUserStore struct {
	usersByID    map[uint64]*model.User
	usersByEmail map[string]*model.User
}

func (m *mockUserStore) Get(ctx context.Context, id uint64) (*model.User, error) {
	if user, ok := m.usersByID[id]; ok {
		return user, nil
	}
	return nil, repoerr.ErrNotFound
}

func (m *mockUserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if user, ok := m.usersByEmail[email]; ok {
		return user, nil
	}
	return nil, repoerr.ErrNotFound
}

type mockRoleStore struct {
	rolesBySlug  map[string]*model.RoleModel
	usersByRole  map[uint64][]uint64
	getBySlugErr error
}

func (m *mockRoleStore) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	if m.getBySlugErr != nil {
		return nil, m.getBySlugErr
	}
	role, ok := m.rolesBySlug[slug]
	if !ok {
		return nil, repoerr.ErrNotFound
	}
	return role, nil
}

func (m *mockRoleStore) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	return m.usersByRole[roleID], nil
}

type mockHubPublisher struct {
	broadcastToGroupCalls int
	broadcastToUserCalls  int
	lastGroupID           uint64
	lastGroupPayload      []byte
	lastUserID            uint64
	lastUserPayload       []byte
}

func (m *mockHubPublisher) BroadcastToGroup(message []byte, groupID uint64) {
	m.broadcastToGroupCalls++
	m.lastGroupID = groupID
	m.lastGroupPayload = message
}

func (m *mockHubPublisher) BroadcastToUser(message []byte, userID uint64) {
	m.broadcastToUserCalls++
	m.lastUserID = userID
	m.lastUserPayload = message
}

func TestService_EnsureSessionSelectsLeastLoadedAgent(t *testing.T) {
	repo := &mockConversationRepository{
		countByAgentResult: map[uint64]int64{
			2001: 5,
			2002: 2,
		},
	}

	roleStore := &mockRoleStore{
		rolesBySlug: map[string]*model.RoleModel{
			string(model.RoleSlugCSAgent):  {Base: model.Base{ID: 1}},
			string(model.RoleSlugCSLeader): {Base: model.Base{ID: 2}},
		},
		usersByRole: map[uint64][]uint64{
			1: {2001, 2002},
			2: {},
		},
	}

	userStore := &mockUserStore{
		usersByID: map[uint64]*model.User{
			2001: {Base: model.Base{ID: 2001}, Status: model.UserStatusActive, Name: "客服1"},
			2002: {Base: model.Base{ID: 2002}, Status: model.UserStatusActive, Name: "客服2"},
		},
	}

	svc := NewService(repo, userStore, roleStore, nil)

	group, agent, err := svc.EnsureSession(context.Background(), 1001)
	require.NoError(t, err)
	require.NotNil(t, group)
	require.NotNil(t, agent)
	assert.Equal(t, uint64(2002), agent.ID)
	require.Len(t, repo.lastCreatedMembers, 2)
	assert.Equal(t, uint64(2002), repo.lastCreatedMembers[1].UserID)
}

func TestService_SendMessagePersistsAndBroadcasts(t *testing.T) {
	hub := &mockHubPublisher{}

	repo := &mockConversationRepository{
		isMember: true,
		getResult: &model.ChatGroup{
			Base:      model.Base{ID: 88},
			GroupType: model.ChatGroupTypePrivate,
			IsActive:  true,
			Members: []model.ChatGroupMember{
				{UserID: 1001, IsActive: true},
				{UserID: 2001, IsActive: true},
			},
		},
	}

	svc := NewService(repo, &mockUserStore{}, &mockRoleStore{}, hub)

	message, err := svc.SendMessage(context.Background(), 1001, 88, "  hello  ")
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, "hello", message.Content)
	assert.Equal(t, uint64(88), repo.lastCreatedMessage.GroupID)
	assert.Equal(t, 1, hub.broadcastToGroupCalls)
	assert.Equal(t, uint64(88), hub.lastGroupID)
	assert.NotEmpty(t, hub.lastGroupPayload)
	assert.Equal(t, 2, hub.broadcastToUserCalls)
}

func TestService_ListConversationsIncludesLastMessage(t *testing.T) {
	now := time.Now()
	repo := &mockConversationRepository{
		listByUserResult: []model.ChatGroup{
			{
				Base:      model.Base{ID: 88, CreatedAt: now, UpdatedAt: now},
				GroupType: model.ChatGroupTypePrivate,
				IsActive:  true,
				Members: []model.ChatGroupMember{
					{UserID: 1001, IsActive: true},
					{UserID: 2001, IsActive: true},
				},
			},
		},
		listByUserTotal: 1,
		listMessagesResult: []model.ChatMessage{
			{
				Base:        model.Base{ID: 9001, CreatedAt: now},
				GroupID:     88,
				SenderID:    2001,
				Content:     "您好，请问有什么可以帮助您",
				MessageType: model.ChatMessageTypeText,
			},
		},
		listMessagesTotal: 1,
		isMember:          true,
	}

	userStore := &mockUserStore{
		usersByID: map[uint64]*model.User{
			2001: {Base: model.Base{ID: 2001}, Status: model.UserStatusActive, Nickname: "在线客服"},
		},
	}

	svc := NewService(repo, userStore, &mockRoleStore{}, nil)

	conversations, total, err := svc.ListConversations(context.Background(), 1001, 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, conversations, 1)
	assert.Equal(t, "您好，请问有什么可以帮助您", conversations[0].LastMessage)
	assert.Equal(t, uint64(2001), conversations[0].AgentID)
}

func TestService_EnsureSessionFailsWithoutAgents(t *testing.T) {
	svc := NewService(
		&mockConversationRepository{},
		&mockUserStore{},
		&mockRoleStore{getBySlugErr: errors.New("role lookup failed")},
		nil,
	)

	_, _, err := svc.EnsureSession(context.Background(), 1001)
	assert.Error(t, err)
}

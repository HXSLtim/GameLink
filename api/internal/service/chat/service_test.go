package chat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ============================================================================
// Mock Implementations
// ============================================================================

// MockChatGroupRepository is a mock implementation of ChatGroupRepository
type MockChatGroupRepository struct {
	mock.Mock
}

func (m *MockChatGroupRepository) Create(ctx context.Context, group *model.ChatGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockChatGroupRepository) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	args := m.Called(ctx, userID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatGroupRepository) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	args := m.Called(ctx, groupID, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatGroupMember), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatGroupRepository) Update(ctx context.Context, group *model.ChatGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockChatGroupRepository) Deactivate(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChatGroupRepository) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
	args := m.Called(ctx, cutoff, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) DeleteByIDs(ctx context.Context, ids []uint64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockChatGroupRepository) GetWithRelations(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) GetByRelatedTeamID(ctx context.Context, teamID uint64) (*model.ChatGroup, error) {
	args := m.Called(ctx, teamID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) GetByRelatedLFGID(ctx context.Context, lfgID uint64) (*model.ChatGroup, error) {
	args := m.Called(ctx, lfgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) GetByVoiceRoomID(ctx context.Context, voiceRoomID string) (*model.ChatGroup, error) {
	args := m.Called(ctx, voiceRoomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) UpdateRoomStatus(ctx context.Context, id uint64, status model.ChatGroupStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockChatGroupRepository) ListGameRooms(ctx context.Context, opts repository.GameRoomListOptions) ([]model.ChatGroup, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatGroupRepository) ListPublicRooms(ctx context.Context, gameID *uint64, page, pageSize int) ([]model.ChatGroup, int64, error) {
	args := m.Called(ctx, gameID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatGroupRepository) ListByHostUserID(ctx context.Context, hostUserID uint64, status *model.ChatGroupStatus) ([]model.ChatGroup, error) {
	args := m.Called(ctx, hostUserID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ChatGroup), args.Error(1)
}

func (m *MockChatGroupRepository) IncrementMemberCount(ctx context.Context, groupID uint64) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

func (m *MockChatGroupRepository) DecrementMemberCount(ctx context.Context, groupID uint64) error {
	args := m.Called(ctx, groupID)
	return args.Error(0)
}

func (m *MockChatGroupRepository) CountByRoomStatus(ctx context.Context) (map[model.ChatGroupStatus]int64, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[model.ChatGroupStatus]int64), args.Error(1)
}

func (m *MockChatGroupRepository) CountActiveRooms(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockChatMemberRepository is a mock implementation of ChatMemberRepository
type MockChatMemberRepository struct {
	mock.Mock
}

func (m *MockChatMemberRepository) Add(ctx context.Context, member *model.ChatGroupMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockChatMemberRepository) AddBatch(ctx context.Context, members []*model.ChatGroupMember) error {
	args := m.Called(ctx, members)
	return args.Error(0)
}

func (m *MockChatMemberRepository) Update(ctx context.Context, member *model.ChatGroupMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MockChatMemberRepository) Remove(ctx context.Context, groupID, userID uint64) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *MockChatMemberRepository) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroupMember), args.Error(1)
}

// MockChatMessageRepository is a mock implementation of ChatMessageRepository
type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) Create(ctx context.Context, message *model.ChatMessage) error {
	args := m.Called(ctx, message)
	if args.Error(0) == nil && message.ID == 0 {
		message.ID = 1
	}
	return args.Error(0)
}

func (m *MockChatMessageRepository) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error {
	args := m.Called(ctx, messages)
	return args.Error(0)
}

func (m *MockChatMessageRepository) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockChatMessageRepository) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error {
	args := m.Called(ctx, id, status, moderatorID, reason)
	return args.Error(0)
}

func (m *MockChatMessageRepository) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
	args := m.Called(ctx, groupIDs)
	return args.Error(0)
}

// MockChatReportRepository is a mock implementation of ChatReportRepository
type MockChatReportRepository struct {
	mock.Mock
}

func (m *MockChatReportRepository) Create(ctx context.Context, report *model.ChatReport) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockChatReportRepository) Get(ctx context.Context, id uint64) (*model.ChatReport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatReport), args.Error(1)
}

func (m *MockChatReportRepository) Update(ctx context.Context, report *model.ChatReport) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockChatReportRepository) List(ctx context.Context, opts repository.ChatReportListOptions) ([]model.ChatReport, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.ChatReport), args.Get(1).(int64), args.Error(2)
}

// MockCache is a mock implementation of cache.Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, key string) (string, bool, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache) GetRedisClient() interface{} {
	return nil
}

func (m *MockCache) Close(ctx context.Context) error {
	return nil
}

// ============================================================================
// ChatService Tests
// ============================================================================

func TestNewChatService(t *testing.T) {
	groups := &MockChatGroupRepository{}
	members := &MockChatMemberRepository{}
	messages := &MockChatMessageRepository{}
	reports := &MockChatReportRepository{}
	cache := &MockCache{}

	svc := NewChatService(groups, members, messages, reports, nil, cache)

	assert.NotNil(t, svc)
	assert.Equal(t, groups, svc.groups)
	assert.Equal(t, members, svc.members)
	assert.Equal(t, messages, svc.messages)
	assert.Equal(t, reports, svc.reports)
	assert.Equal(t, cache, svc.cache)
}

func TestChatService_ListUserGroups(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint64
		page         int
		pageSize     int
		setupMock    func(*MockChatGroupRepository)
		expectError  bool
		expectGroups int
	}{
		{
			name:     "successful list groups",
			userID:   1,
			page:     1,
			pageSize: 10,
			setupMock: func(groups *MockChatGroupRepository) {
				chatGroups := []model.ChatGroup{
					{Base: model.Base{ID: 1}, GroupName: "Order Chat 1", GroupType: model.ChatGroupTypeOrder, IsActive: true},
					{Base: model.Base{ID: 2}, GroupName: "Order Chat 2", GroupType: model.ChatGroupTypeOrder, IsActive: true},
				}
				groups.On("ListByUser", mock.Anything, uint64(1), mock.MatchedBy(func(opts repository.ChatGroupListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 10
				})).Return(chatGroups, int64(2), nil)
			},
			expectError:  false,
			expectGroups: 2,
		},
		{
			name:     "default pagination values",
			userID:   1,
			page:     0,
			pageSize: 0,
			setupMock: func(groups *MockChatGroupRepository) {
				groups.On("ListByUser", mock.Anything, uint64(1), mock.MatchedBy(func(opts repository.ChatGroupListOptions) bool {
					return opts.Page == 1 && opts.PageSize == 20
				})).Return([]model.ChatGroup{}, int64(0), nil)
			},
			expectError:  false,
			expectGroups: 0,
		},
		{
			name:     "page size exceeds max",
			userID:   1,
			page:     1,
			pageSize: 200,
			setupMock: func(groups *MockChatGroupRepository) {
				groups.On("ListByUser", mock.Anything, uint64(1), mock.MatchedBy(func(opts repository.ChatGroupListOptions) bool {
					return opts.PageSize == 20
				})).Return([]model.ChatGroup{}, int64(0), nil)
			},
			expectError:  false,
			expectGroups: 0,
		},
		{
			name:     "repository error",
			userID:   1,
			page:     1,
			pageSize: 10,
			setupMock: func(groups *MockChatGroupRepository) {
				groups.On("ListByUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(groups)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			result, total, err := svc.ListUserGroups(context.Background(), tt.userID, tt.page, tt.pageSize)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectGroups)
				assert.Equal(t, int64(tt.expectGroups), total)
			}
		})
	}
}

func TestChatService_EnsureMembership(t *testing.T) {
	tests := []struct {
		name        string
		groupID     uint64
		userID      uint64
		setupMock   func(*MockChatMemberRepository)
		expectError bool
	}{
		{
			name:    "active member",
			groupID: 1,
			userID:  1,
			setupMock: func(members *MockChatMemberRepository) {
				member := &model.ChatGroupMember{
					GroupID:  1,
					UserID:   1,
					IsActive: true,
				}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)
			},
			expectError: false,
		},
		{
			name:    "inactive member",
			groupID: 1,
			userID:  2,
			setupMock: func(members *MockChatMemberRepository) {
				member := &model.ChatGroupMember{
					GroupID:  1,
					UserID:   2,
					IsActive: false,
				}
				members.On("Get", mock.Anything, uint64(1), uint64(2)).Return(member, nil)
			},
			expectError: true,
		},
		{
			name:    "member not found",
			groupID: 1,
			userID:  3,
			setupMock: func(members *MockChatMemberRepository) {
				members.On("Get", mock.Anything, uint64(1), uint64(3)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(members)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			member, err := svc.EnsureMembership(context.Background(), tt.groupID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, member)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, member)
			}
		})
	}
}

func TestChatService_SendMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       SendMessageInput
		setupMock   func(*MockChatGroupRepository, *MockChatMemberRepository, *MockChatMessageRepository, *MockCache)
		expectError bool
	}{
		{
			name: "successful send message in order group",
			input: SendMessageInput{
				GroupID:     1,
				SenderID:    1,
				Content:     "Hello!",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
				member := &model.ChatGroupMember{GroupID: 1, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 1}, GroupType: model.ChatGroupTypeOrder, IsActive: true}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)

				messages.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "send message in public group with throttle",
			input: SendMessageInput{
				GroupID:     2,
				SenderID:    1,
				Content:     "Public message",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
				member := &model.ChatGroupMember{GroupID: 2, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(2), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 2}, GroupType: model.ChatGroupTypePublic, IsActive: true}
				groups.On("Get", mock.Anything, uint64(2)).Return(group, nil)

				cache.On("Get", mock.Anything, mock.Anything).Return("", false, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				messages.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "throttled in public group",
			input: SendMessageInput{
				GroupID:     2,
				SenderID:    1,
				Content:     "Throttled message",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
				member := &model.ChatGroupMember{GroupID: 2, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(2), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 2}, GroupType: model.ChatGroupTypePublic, IsActive: true}
				groups.On("Get", mock.Anything, uint64(2)).Return(group, nil)

				cache.On("Get", mock.Anything, mock.Anything).Return("1", true, nil)
			},
			expectError: true,
		},
		{
			name: "empty content",
			input: SendMessageInput{
				GroupID:     1,
				SenderID:    1,
				Content:     "",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
			},
			expectError: true,
		},
		{
			name: "inactive group",
			input: SendMessageInput{
				GroupID:     1,
				SenderID:    1,
				Content:     "Hello",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
				member := &model.ChatGroupMember{GroupID: 1, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 1}, GroupType: model.ChatGroupTypeOrder, IsActive: false}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)
			},
			expectError: true,
		},
		{
			name: "not a member",
			input: SendMessageInput{
				GroupID:     1,
				SenderID:    1,
				Content:     "Hello",
				MessageType: model.ChatMessageTypeText,
			},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository, cache *MockCache) {
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(groups, members, messages, cache)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			msg, err := svc.SendMessage(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, msg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
			}
		})
	}
}

func TestChatService_JoinGroup(t *testing.T) {
	tests := []struct {
		name        string
		groupID     uint64
		userID      uint64
		nickname    string
		setupMock   func(*MockChatGroupRepository, *MockChatMemberRepository)
		expectError bool
	}{
		{
			name:     "join new group",
			groupID:  1,
			userID:   1,
			nickname: "TestUser",
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository) {
				group := &model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(nil, repository.ErrNotFound)
				members.On("Add", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "reactivate existing member",
			groupID:  1,
			userID:   2,
			nickname: "ReactivatedUser",
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository) {
				group := &model.ChatGroup{Base: model.Base{ID: 1}, IsActive: true}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)
				member := &model.ChatGroupMember{GroupID: 1, UserID: 2, IsActive: false}
				members.On("Get", mock.Anything, uint64(1), uint64(2)).Return(member, nil)
				members.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:     "join inactive group",
			groupID:  1,
			userID:   1,
			nickname: "TestUser",
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository) {
				group := &model.ChatGroup{Base: model.Base{ID: 1}, IsActive: false}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(groups, members)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			err := svc.JoinGroup(context.Background(), tt.groupID, tt.userID, tt.nickname)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatService_LeaveGroup(t *testing.T) {
	tests := []struct {
		name        string
		groupID     uint64
		userID      uint64
		setupMock   func(*MockChatMemberRepository)
		expectError bool
	}{
		{
			name:    "successful leave",
			groupID: 1,
			userID:  1,
			setupMock: func(members *MockChatMemberRepository) {
				member := &model.ChatGroupMember{GroupID: 1, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)
				members.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:    "not a member",
			groupID: 1,
			userID:  2,
			setupMock: func(members *MockChatMemberRepository) {
				members.On("Get", mock.Anything, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(members)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			err := svc.LeaveGroup(context.Background(), tt.groupID, tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatService_MarkRead(t *testing.T) {
	tests := []struct {
		name        string
		groupID     uint64
		userID      uint64
		messageID   uint64
		setupMock   func(*MockChatMemberRepository)
		expectError bool
	}{
		{
			name:      "successful mark read",
			groupID:   1,
			userID:    1,
			messageID: 100,
			setupMock: func(members *MockChatMemberRepository) {
				member := &model.ChatGroupMember{GroupID: 1, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)
				members.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "not a member",
			groupID:   1,
			userID:    2,
			messageID: 100,
			setupMock: func(members *MockChatMemberRepository) {
				members.On("Get", mock.Anything, uint64(1), uint64(2)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(members)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			err := svc.MarkRead(context.Background(), tt.groupID, tt.userID, tt.messageID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatService_ApproveMessage(t *testing.T) {
	groups := &MockChatGroupRepository{}
	members := &MockChatMemberRepository{}
	messages := &MockChatMessageRepository{}
	reports := &MockChatReportRepository{}
	cache := &MockCache{}

	moderatorID := uint64(100)
	messages.On("UpdateAuditStatus", mock.Anything, uint64(1), model.ChatMessageAuditApproved, &moderatorID, "").Return(nil)

	svc := NewChatService(groups, members, messages, reports, nil, cache)
	err := svc.ApproveMessage(context.Background(), 1, moderatorID)

	assert.NoError(t, err)
}

func TestChatService_RejectMessage(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{name: "with reason", reason: "inappropriate content"},
		{name: "empty reason defaults to rejected", reason: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			moderatorID := uint64(100)
			expectedReason := tt.reason
			if expectedReason == "" {
				expectedReason = "rejected"
			}
			messages.On("UpdateAuditStatus", mock.Anything, uint64(1), model.ChatMessageAuditRejected, &moderatorID, expectedReason).Return(nil)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			err := svc.RejectMessage(context.Background(), 1, moderatorID, tt.reason)

			assert.NoError(t, err)
		})
	}
}

func TestChatService_ReportMessage(t *testing.T) {
	tests := []struct {
		name        string
		reporterID  uint64
		messageID   uint64
		reason      string
		evidence    string
		setupMock   func(*MockChatReportRepository)
		nilReports  bool
		expectError bool
	}{
		{
			name:       "successful report",
			reporterID: 1,
			messageID:  100,
			reason:     "spam",
			evidence:   "screenshot.png",
			setupMock: func(reports *MockChatReportRepository) {
				reports.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "empty reason defaults to unspecified",
			reporterID: 1,
			messageID:  100,
			reason:     "",
			evidence:   "",
			setupMock: func(reports *MockChatReportRepository) {
				reports.On("Create", mock.Anything, mock.MatchedBy(func(r *model.ChatReport) bool {
					return r.Reason == "unspecified"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "nil reports repository",
			reporterID:  1,
			messageID:   100,
			reason:      "spam",
			nilReports:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			cache := &MockCache{}

			var svc *ChatService
			if tt.nilReports {
				// Pass actual nil interface, not typed nil
				svc = NewChatService(groups, members, messages, nil, nil, cache)
			} else {
				reports := &MockChatReportRepository{}
				tt.setupMock(reports)
				svc = NewChatService(groups, members, messages, reports, nil, cache)
			}

			err := svc.ReportMessage(context.Background(), tt.reporterID, tt.messageID, tt.reason, tt.evidence)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatService_ListMessages(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint64
		groupID        uint64
		opts           ListMessagesOptions
		setupMock      func(*MockChatGroupRepository, *MockChatMemberRepository, *MockChatMessageRepository)
		expectError    bool
		expectMessages int
	}{
		{
			name:    "successful list messages in order group",
			userID:  1,
			groupID: 1,
			opts:    ListMessagesOptions{Page: 1, PageSize: 20},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository) {
				member := &model.ChatGroupMember{GroupID: 1, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 1}, GroupType: model.ChatGroupTypeOrder}
				groups.On("Get", mock.Anything, uint64(1)).Return(group, nil)

				msgs := []model.ChatMessage{
					{Base: model.Base{ID: 1}, Content: "Hello"},
					{Base: model.Base{ID: 2}, Content: "World"},
				}
				messages.On("ListByGroup", mock.Anything, mock.Anything).Return(msgs, int64(2), nil)
			},
			expectError:    false,
			expectMessages: 2,
		},
		{
			name:    "list messages in public group filters by audit status",
			userID:  1,
			groupID: 2,
			opts:    ListMessagesOptions{Page: 1, PageSize: 20},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository) {
				member := &model.ChatGroupMember{GroupID: 2, UserID: 1, IsActive: true}
				members.On("Get", mock.Anything, uint64(2), uint64(1)).Return(member, nil)

				group := &model.ChatGroup{Base: model.Base{ID: 2}, GroupType: model.ChatGroupTypePublic}
				groups.On("Get", mock.Anything, uint64(2)).Return(group, nil)

				msgs := []model.ChatMessage{
					{Base: model.Base{ID: 1}, Content: "Approved message", AuditStatus: model.ChatMessageAuditApproved},
				}
				messages.On("ListByGroup", mock.Anything, mock.MatchedBy(func(opts repository.ChatMessageListOptions) bool {
					return len(opts.AuditStatuses) == 1 && opts.AuditStatuses[0] == model.ChatMessageAuditApproved
				})).Return(msgs, int64(1), nil)
			},
			expectError:    false,
			expectMessages: 1,
		},
		{
			name:    "not a member",
			userID:  1,
			groupID: 1,
			opts:    ListMessagesOptions{Page: 1, PageSize: 20},
			setupMock: func(groups *MockChatGroupRepository, members *MockChatMemberRepository, messages *MockChatMessageRepository) {
				members.On("Get", mock.Anything, uint64(1), uint64(1)).Return(nil, repository.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			members := &MockChatMemberRepository{}
			messages := &MockChatMessageRepository{}
			reports := &MockChatReportRepository{}
			cache := &MockCache{}

			tt.setupMock(groups, members, messages)

			svc := NewChatService(groups, members, messages, reports, nil, cache)
			msgs, total, err := svc.ListMessages(context.Background(), tt.userID, tt.groupID, tt.opts)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, msgs, tt.expectMessages)
				assert.Equal(t, int64(tt.expectMessages), total)
			}
		})
	}
}

// ============================================================================
// CleanupService Tests
// ============================================================================

func TestNewCleanupService(t *testing.T) {
	groups := &MockChatGroupRepository{}

	svc := NewCleanupService(groups)

	assert.NotNil(t, svc)
	assert.Equal(t, groups, svc.groups)
}

func TestCleanupService_CleanupInactiveOrderGroups(t *testing.T) {
	tests := []struct {
		name        string
		cutoff      time.Time
		limit       int
		setupMock   func(*MockChatGroupRepository)
		expectError bool
		expectCount int
	}{
		{
			name:   "successful cleanup",
			cutoff: time.Now().Add(-24 * time.Hour),
			limit:  100,
			setupMock: func(groups *MockChatGroupRepository) {
				deactivatedGroups := []model.ChatGroup{
					{Base: model.Base{ID: 1}},
					{Base: model.Base{ID: 2}},
					{Base: model.Base{ID: 3}},
				}
				groups.On("ListDeactivatedBefore", mock.Anything, mock.Anything, 100).Return(deactivatedGroups, nil)
				groups.On("DeleteByIDs", mock.Anything, []uint64{1, 2, 3}).Return(nil)
			},
			expectError: false,
			expectCount: 3,
		},
		{
			name:   "no groups to cleanup",
			cutoff: time.Now().Add(-24 * time.Hour),
			limit:  100,
			setupMock: func(groups *MockChatGroupRepository) {
				groups.On("ListDeactivatedBefore", mock.Anything, mock.Anything, 100).Return([]model.ChatGroup{}, nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:   "list error",
			cutoff: time.Now().Add(-24 * time.Hour),
			limit:  100,
			setupMock: func(groups *MockChatGroupRepository) {
				groups.On("ListDeactivatedBefore", mock.Anything, mock.Anything, 100).Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
		{
			name:   "delete error",
			cutoff: time.Now().Add(-24 * time.Hour),
			limit:  100,
			setupMock: func(groups *MockChatGroupRepository) {
				deactivatedGroups := []model.ChatGroup{
					{Base: model.Base{ID: 1}},
				}
				groups.On("ListDeactivatedBefore", mock.Anything, mock.Anything, 100).Return(deactivatedGroups, nil)
				groups.On("DeleteByIDs", mock.Anything, []uint64{1}).Return(errors.New("delete error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := &MockChatGroupRepository{}
			tt.setupMock(groups)

			svc := NewCleanupService(groups)
			count, err := svc.CleanupInactiveOrderGroups(context.Background(), tt.cutoff, tt.limit)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectCount, count)
			}
		})
	}
}

// ============================================================================
// Structure Tests
// ============================================================================

func TestSendMessageInput_Structure(t *testing.T) {
	replyID := uint64(10)
	input := SendMessageInput{
		GroupID:     1,
		SenderID:    2,
		Content:     "Hello World",
		MessageType: model.ChatMessageTypeText,
		ReplyToID:   &replyID,
		ImageURL:    "https://example.com/image.png",
	}

	assert.Equal(t, uint64(1), input.GroupID)
	assert.Equal(t, uint64(2), input.SenderID)
	assert.Equal(t, "Hello World", input.Content)
	assert.Equal(t, model.ChatMessageTypeText, input.MessageType)
	assert.Equal(t, uint64(10), *input.ReplyToID)
	assert.Equal(t, "https://example.com/image.png", input.ImageURL)
}

func TestListMessagesOptions_Structure(t *testing.T) {
	beforeID := uint64(100)
	afterID := uint64(50)
	opts := ListMessagesOptions{
		Page:     1,
		PageSize: 20,
		BeforeID: &beforeID,
		AfterID:  &afterID,
	}

	assert.Equal(t, 1, opts.Page)
	assert.Equal(t, 20, opts.PageSize)
	assert.Equal(t, uint64(100), *opts.BeforeID)
	assert.Equal(t, uint64(50), *opts.AfterID)
}

func TestChatErrors(t *testing.T) {
	assert.NotNil(t, ErrNotFound)
	assert.NotNil(t, ErrNotMember)
	assert.NotNil(t, ErrInactiveGroup)
	assert.NotNil(t, ErrMessageTooLarge)
	assert.NotNil(t, ErrThrottled)
}

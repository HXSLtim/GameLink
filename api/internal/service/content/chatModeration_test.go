package content

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockChatMessageRepository is a mock for ChatMessageRepository
type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) Create(ctx context.Context, message *model.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockChatMessageRepository) CreateBatch(ctx context.Context, messages []*model.ChatMessage) error {
	args := m.Called(ctx, messages)
	return args.Error(0)
}

func (m *MockChatMessageRepository) Get(ctx context.Context, id uint64) (*model.ChatMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) ListByGroup(ctx context.Context, opts repository.ChatMessageListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) ListForModeration(ctx context.Context, opts repository.ChatMessageModerationListOptions) ([]model.ChatMessage, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ChatMessage), args.Get(1).(int64), args.Error(2)
}

func (m *MockChatMessageRepository) MarkDeleted(ctx context.Context, id uint64, deletedBy uint64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockChatMessageRepository) UpdateAuditStatus(ctx context.Context, id uint64, status model.ChatMessageAuditStatus, moderatorID *uint64, reason string) error {
	args := m.Called(ctx, id, status, moderatorID, reason)
	return args.Error(0)
}

func (m *MockChatMessageRepository) DeleteByGroupIDs(ctx context.Context, groupIDs []uint64) error {
	args := m.Called(ctx, groupIDs)
	return args.Error(0)
}

// MockChatMemberRepository is a mock for ChatMemberRepository
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

func (m *MockChatMemberRepository) Get(ctx context.Context, groupID, userID uint64) (*model.ChatGroupMember, error) {
	args := m.Called(ctx, groupID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatGroupMember), args.Error(1)
}

func (m *MockChatMemberRepository) Remove(ctx context.Context, groupID, userID uint64) error {
	args := m.Called(ctx, groupID, userID)
	return args.Error(0)
}

func (m *MockChatMemberRepository) Update(ctx context.Context, member *model.ChatGroupMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func TestNewChatModerationService(t *testing.T) {
	messageRepo := &MockChatMessageRepository{}
	memberRepo := &MockChatMemberRepository{}
	opLogRepo := &MockOperationLogRepository{}

	svc := NewChatModerationService(messageRepo, memberRepo, nil, opLogRepo)

	require.NotNil(t, svc)
	assert.Equal(t, messageRepo, svc.messageRepo)
	assert.Equal(t, memberRepo, svc.memberRepo)
	assert.Nil(t, svc.sensitiveWord)
	assert.Equal(t, opLogRepo, svc.opLogRepo)
}

func TestChatModerationService_ListMessages(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with default pagination", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		messages := []model.ChatMessage{
			{Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now}, GroupID: 100, SenderID: 1, Content: "Hello", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditPending},
			{Base: model.Base{ID: 2, CreatedAt: now, UpdatedAt: now}, GroupID: 100, SenderID: 2, Content: "Hi", MessageType: model.ChatMessageTypeText, AuditStatus: model.ChatMessageAuditApproved},
		}
		messageRepo.On("ListForModeration", ctx, mock.MatchedBy(func(opts repository.ChatMessageModerationListOptions) bool {
			return opts.Page == 1 && opts.PageSize == 20
		})).Return(messages, int64(2), nil)

		req := ListMessagesRequest{Page: 0, PageSize: 0}
		result, err := svc.ListMessages(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
	})

	t.Run("success with filters", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		groupID := uint64(100)
		senderID := uint64(1)
		status := model.ChatMessageAuditRejected
		messages := []model.ChatMessage{
			{Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now}, GroupID: 100, SenderID: 1, Content: "Bad", AuditStatus: model.ChatMessageAuditRejected},
		}
		messageRepo.On("ListForModeration", ctx, mock.MatchedBy(func(opts repository.ChatMessageModerationListOptions) bool {
			return opts.Page == 1 && opts.PageSize == 10 && opts.GroupID != nil && *opts.GroupID == 100
		})).Return(messages, int64(1), nil)

		req := ListMessagesRequest{
			Page:        1,
			PageSize:    10,
			GroupID:     &groupID,
			SenderID:    &senderID,
			AuditStatus: &status,
		}
		result, err := svc.ListMessages(ctx, req)

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, int64(1), result.Total)
	})

	t.Run("page size defaults to 20 when invalid", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		messageRepo.On("ListForModeration", ctx, mock.MatchedBy(func(opts repository.ChatMessageModerationListOptions) bool {
			return opts.PageSize == 20 // Should default to 20
		})).Return([]model.ChatMessage{}, int64(0), nil)

		req := ListMessagesRequest{Page: 1, PageSize: 200}
		_, err := svc.ListMessages(ctx, req)

		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		messageRepo.On("ListForModeration", ctx, mock.AnythingOfType("repository.ChatMessageModerationListOptions")).
			Return([]model.ChatMessage{}, int64(0), errors.New("db error"))

		req := ListMessagesRequest{Page: 1, PageSize: 10}
		_, err := svc.ListMessages(ctx, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_DeleteMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, opLogRepo)

		moderatorID := uint64(10)
		messageRepo.On("MarkDeleted", ctx, uint64(1), moderatorID).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.DeleteMessage(ctx, 1, moderatorID, "spam")

		require.NoError(t, err)
		messageRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		messageRepo.On("MarkDeleted", ctx, uint64(1), uint64(10)).
			Return(errors.New("db error"))

		err := svc.DeleteMessage(ctx, 1, 10, "spam")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_ApproveMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, opLogRepo)

		moderatorID := uint64(10)
		messageRepo.On("UpdateAuditStatus", ctx, uint64(1), model.ChatMessageAuditApproved, &moderatorID, "").Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.ApproveMessage(ctx, 1, moderatorID)

		require.NoError(t, err)
		messageRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		moderatorID := uint64(10)
		messageRepo.On("UpdateAuditStatus", ctx, uint64(1), model.ChatMessageAuditApproved, &moderatorID, "").
			Return(errors.New("db error"))

		err := svc.ApproveMessage(ctx, 1, moderatorID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_RejectMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, opLogRepo)

		moderatorID := uint64(10)
		messageRepo.On("UpdateAuditStatus", ctx, uint64(1), model.ChatMessageAuditRejected, &moderatorID, "violation").Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.RejectMessage(ctx, 1, moderatorID, "violation")

		require.NoError(t, err)
		messageRepo.AssertExpectations(t)
	})

	t.Run("empty reason returns validation error", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		err := svc.RejectMessage(ctx, 1, 10, "")

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("repo error", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		svc := NewChatModerationService(messageRepo, nil, nil, nil)

		moderatorID := uint64(10)
		messageRepo.On("UpdateAuditStatus", ctx, uint64(1), model.ChatMessageAuditRejected, &moderatorID, "reason").
			Return(errors.New("db error"))

		err := svc.RejectMessage(ctx, 1, moderatorID, "reason")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_MuteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		memberRepo := &MockChatMemberRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewChatModerationService(messageRepo, memberRepo, nil, opLogRepo)

		member := &model.ChatGroupMember{Base: model.Base{ID: 1}, GroupID: 100, UserID: 5}
		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(member, nil)
		memberRepo.On("Update", ctx, mock.MatchedBy(func(m *model.ChatGroupMember) bool {
			return m.IsMuted && m.MutedUntil != nil
		})).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		req := MuteUserRequest{
			GroupID:  100,
			UserID:   5,
			Duration: 30,
			Reason:   "spam",
		}
		err := svc.MuteUser(ctx, req, 10)

		require.NoError(t, err)
		memberRepo.AssertExpectations(t)
	})

	t.Run("duration validation error", func(t *testing.T) {
		svc := NewChatModerationService(nil, nil, nil, nil)

		req := MuteUserRequest{
			GroupID:  100,
			UserID:   5,
			Duration: 0,
			Reason:   "test",
		}
		err := svc.MuteUser(ctx, req, 10)

		require.Error(t, err)
		assert.Equal(t, ErrAdminValidation, err)
	})

	t.Run("member not found", func(t *testing.T) {
		memberRepo := &MockChatMemberRepository{}
		svc := NewChatModerationService(nil, memberRepo, nil, nil)

		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(nil, repository.ErrNotFound)

		req := MuteUserRequest{
			GroupID:  100,
			UserID:   5,
			Duration: 30,
		}
		err := svc.MuteUser(ctx, req, 10)

		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("update error", func(t *testing.T) {
		memberRepo := &MockChatMemberRepository{}
		svc := NewChatModerationService(nil, memberRepo, nil, nil)

		member := &model.ChatGroupMember{Base: model.Base{ID: 1}, GroupID: 100, UserID: 5}
		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(member, nil)
		memberRepo.On("Update", ctx, mock.AnythingOfType("*model.ChatGroupMember")).Return(errors.New("db error"))

		req := MuteUserRequest{
			GroupID:  100,
			UserID:   5,
			Duration: 30,
		}
		err := svc.MuteUser(ctx, req, 10)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_UnmuteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		messageRepo := &MockChatMessageRepository{}
		memberRepo := &MockChatMemberRepository{}
		opLogRepo := &MockOperationLogRepository{}
		svc := NewChatModerationService(messageRepo, memberRepo, nil, opLogRepo)

		mutedUntil := time.Now().Add(1 * time.Hour)
		member := &model.ChatGroupMember{
			Base:       model.Base{ID: 1},
			GroupID:    100,
			UserID:     5,
			IsMuted:    true,
			MutedUntil: &mutedUntil,
			MutedBy:    pointerTo(uint64(10)),
			MuteReason: "spam",
		}
		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(member, nil)
		memberRepo.On("Update", ctx, mock.MatchedBy(func(m *model.ChatGroupMember) bool {
			return !m.IsMuted && m.MutedUntil == nil
		})).Return(nil)
		opLogRepo.On("Append", ctx, mock.AnythingOfType("*model.OperationLog")).Return(nil)

		err := svc.UnmuteUser(ctx, 100, 5, 10)

		require.NoError(t, err)
		memberRepo.AssertExpectations(t)
	})

	t.Run("member not found", func(t *testing.T) {
		memberRepo := &MockChatMemberRepository{}
		svc := NewChatModerationService(nil, memberRepo, nil, nil)

		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(nil, repository.ErrNotFound)

		err := svc.UnmuteUser(ctx, 100, 5, 10)

		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("update error", func(t *testing.T) {
		memberRepo := &MockChatMemberRepository{}
		svc := NewChatModerationService(nil, memberRepo, nil, nil)

		member := &model.ChatGroupMember{Base: model.Base{ID: 1}, GroupID: 100, UserID: 5, IsMuted: true}
		memberRepo.On("Get", ctx, uint64(100), uint64(5)).Return(member, nil)
		memberRepo.On("Update", ctx, mock.AnythingOfType("*model.ChatGroupMember")).Return(errors.New("db error"))

		err := svc.UnmuteUser(ctx, 100, 5, 10)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestChatModerationService_ToMessageDTO(t *testing.T) {
	now := time.Now()

	t.Run("basic message", func(t *testing.T) {
		svc := &ChatModerationService{}
		msg := &model.ChatMessage{
			Base:        model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			GroupID:     100,
			SenderID:    5,
			Content:     "Hello world",
			MessageType: model.ChatMessageTypeText,
			AuditStatus: model.ChatMessageAuditPending,
		}

		dto := svc.toMessageDTO(msg)

		assert.Equal(t, uint64(1), dto.ID)
		assert.Equal(t, uint64(100), dto.GroupID)
		assert.Equal(t, uint64(5), dto.SenderID)
		assert.Equal(t, "Hello world", dto.Content)
		assert.Equal(t, model.ChatMessageTypeText, dto.MessageType)
		assert.Equal(t, model.ChatMessageAuditPending, dto.AuditStatus)
		assert.NotEmpty(t, dto.CreatedAt)
	})

	t.Run("message with rejection", func(t *testing.T) {
		svc := &ChatModerationService{}
		msg := &model.ChatMessage{
			Base:        model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			Content:     "Spam",
			AuditStatus: model.ChatMessageAuditRejected,
			RejectReason: "violation",
		}

		dto := svc.toMessageDTO(msg)

		assert.Equal(t, "Spam", dto.Content)
		assert.Equal(t, model.ChatMessageAuditRejected, dto.AuditStatus)
		assert.Equal(t, "violation", dto.RejectReason)
	})
}

// Helper function
func pointerTo[T any](v T) *T {
	return &v
}

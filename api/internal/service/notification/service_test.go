package notification

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

// MockNotificationRepository is a mock implementation of NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.NotificationEvent), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	args := m.Called(ctx, userID, ids)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllRead(ctx context.Context, userID uint64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) Create(ctx context.Context, event *model.NotificationEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(ctx context.Context, userID uint64, id uint64) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestNewService(t *testing.T) {
	repo := &MockNotificationRepository{}

	svc := NewService(repo)

	assert.NotNil(t, svc)
	assert.Equal(t, repo, svc.repo)
}

func TestService_List(t *testing.T) {
	now := time.Now()
	readAt := now.Add(-time.Hour)
	refID := uint64(100)

	tests := []struct {
		name          string
		userID        uint64
		req           ListRequest
		setupMock     func(*MockNotificationRepository)
		expectError   bool
		expectItems   int
		expectTotal   int64
		expectUnread  int64
	}{
		{
			name:   "successful list notifications",
			userID: 1,
			req: ListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(repo *MockNotificationRepository) {
				notifications := []model.NotificationEvent{
					{
						Base:          model.Base{ID: 1, CreatedAt: now},
						UserID:        1,
						Title:         "Order Update",
						Message:       "Your order has been accepted",
						Priority:      model.NotificationPriorityNormal,
						Channel:       "web",
						ReferenceType: "order",
						ReferenceID:   &refID,
						ReadAt:        nil,
					},
					{
						Base:          model.Base{ID: 2, CreatedAt: now.Add(-time.Hour)},
						UserID:        1,
						Title:         "System Notice",
						Message:       "Welcome to GameLink",
						Priority:      model.NotificationPriorityLow,
						Channel:       "web",
						ReferenceType: "system",
						ReferenceID:   nil,
						ReadAt:        &readAt,
					},
				}
				repo.On("ListByUser", mock.Anything, mock.MatchedBy(func(opts repository.NotificationListOptions) bool {
					return opts.UserID == 1 && opts.Page == 1 && opts.PageSize == 10
				})).Return(notifications, int64(2), nil)
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(1), nil)
			},
			expectError:  false,
			expectItems:  2,
			expectTotal:  2,
			expectUnread: 1,
		},
		{
			name:   "list with unread only filter",
			userID: 1,
			req: ListRequest{
				Page:       1,
				PageSize:   10,
				UnreadOnly: true,
			},
			setupMock: func(repo *MockNotificationRepository) {
				notifications := []model.NotificationEvent{
					{
						Base:          model.Base{ID: 1, CreatedAt: now},
						UserID:        1,
						Title:         "Unread Notification",
						Message:       "This is unread",
						Priority:      model.NotificationPriorityHigh,
						Channel:       "push",
						ReferenceType: "order",
						ReferenceID:   &refID,
						ReadAt:        nil,
					},
				}
				repo.On("ListByUser", mock.Anything, mock.MatchedBy(func(opts repository.NotificationListOptions) bool {
					return opts.UserID == 1 && opts.Unread != nil && *opts.Unread == true
				})).Return(notifications, int64(1), nil)
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(1), nil)
			},
			expectError:  false,
			expectItems:  1,
			expectTotal:  1,
			expectUnread: 1,
		},
		{
			name:   "list with priority filter",
			userID: 1,
			req: ListRequest{
				Page:       1,
				PageSize:   10,
				Priorities: []model.NotificationPriority{model.NotificationPriorityHigh},
			},
			setupMock: func(repo *MockNotificationRepository) {
				notifications := []model.NotificationEvent{
					{
						Base:          model.Base{ID: 1, CreatedAt: now},
						UserID:        1,
						Title:         "High Priority",
						Message:       "Urgent notification",
						Priority:      model.NotificationPriorityHigh,
						Channel:       "push",
						ReferenceType: "dispute",
						ReferenceID:   &refID,
						ReadAt:        nil,
					},
				}
				repo.On("ListByUser", mock.Anything, mock.MatchedBy(func(opts repository.NotificationListOptions) bool {
					return opts.UserID == 1 && len(opts.Priority) == 1 && opts.Priority[0] == model.NotificationPriorityHigh
				})).Return(notifications, int64(1), nil)
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(1), nil)
			},
			expectError:  false,
			expectItems:  1,
			expectTotal:  1,
			expectUnread: 1,
		},
		{
			name:   "empty list",
			userID: 2,
			req: ListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("ListByUser", mock.Anything, mock.MatchedBy(func(opts repository.NotificationListOptions) bool {
					return opts.UserID == 2
				})).Return([]model.NotificationEvent{}, int64(0), nil)
				repo.On("CountUnread", mock.Anything, uint64(2)).Return(int64(0), nil)
			},
			expectError:  false,
			expectItems:  0,
			expectTotal:  0,
			expectUnread: 0,
		},
		{
			name:   "repository error on list",
			userID: 1,
			req: ListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("ListByUser", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
			},
			expectError: true,
		},
		{
			name:   "repository error on count unread",
			userID: 1,
			req: ListRequest{
				Page:     1,
				PageSize: 10,
			},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("ListByUser", mock.Anything, mock.Anything).Return([]model.NotificationEvent{}, int64(0), nil)
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(0), errors.New("count error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockNotificationRepository{}
			tt.setupMock(repo)

			svc := NewService(repo)
			resp, err := svc.List(context.Background(), tt.userID, tt.req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Items, tt.expectItems)
				assert.Equal(t, tt.expectTotal, resp.Total)
				assert.Equal(t, tt.expectUnread, resp.UnreadCount)
				assert.Equal(t, tt.req.Page, resp.Page)
				assert.Equal(t, tt.req.PageSize, resp.PageSize)
			}
		})
	}
}

func TestService_MarkRead(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		ids         []uint64
		setupMock   func(*MockNotificationRepository)
		expectError bool
	}{
		{
			name:   "successful mark read",
			userID: 1,
			ids:    []uint64{1, 2, 3},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("MarkRead", mock.Anything, uint64(1), []uint64{1, 2, 3}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "mark single notification read",
			userID: 1,
			ids:    []uint64{1},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("MarkRead", mock.Anything, uint64(1), []uint64{1}).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "repository error",
			userID: 1,
			ids:    []uint64{1},
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("MarkRead", mock.Anything, uint64(1), []uint64{1}).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockNotificationRepository{}
			tt.setupMock(repo)

			svc := NewService(repo)
			err := svc.MarkRead(context.Background(), tt.userID, tt.ids)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetUnreadCount(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		setupMock   func(*MockNotificationRepository)
		expectError bool
		expectCount int64
	}{
		{
			name:   "successful get unread count",
			userID: 1,
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(5), nil)
			},
			expectError: false,
			expectCount: 5,
		},
		{
			name:   "zero unread count",
			userID: 2,
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("CountUnread", mock.Anything, uint64(2)).Return(int64(0), nil)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:   "repository error",
			userID: 1,
			setupMock: func(repo *MockNotificationRepository) {
				repo.On("CountUnread", mock.Anything, uint64(1)).Return(int64(0), errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockNotificationRepository{}
			tt.setupMock(repo)

			svc := NewService(repo)
			count, err := svc.GetUnreadCount(context.Background(), tt.userID)

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

func TestListRequest_Structure(t *testing.T) {
	req := ListRequest{
		Page:       1,
		PageSize:   20,
		UnreadOnly: true,
		Priorities: []model.NotificationPriority{model.NotificationPriorityHigh, model.NotificationPriorityNormal},
	}

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 20, req.PageSize)
	assert.True(t, req.UnreadOnly)
	assert.Len(t, req.Priorities, 2)
}

func TestNotificationView_Structure(t *testing.T) {
	now := time.Now()
	readAt := now.Add(-time.Hour)
	refID := uint64(123)

	view := NotificationView{
		ID:            1,
		Title:         "Test Notification",
		Message:       "This is a test message",
		Priority:      model.NotificationPriorityHigh,
		Channel:       "push",
		ReferenceType: "order",
		ReferenceID:   &refID,
		ReadAt:        &readAt,
		CreatedAt:     now,
	}

	assert.Equal(t, uint64(1), view.ID)
	assert.Equal(t, "Test Notification", view.Title)
	assert.Equal(t, "This is a test message", view.Message)
	assert.Equal(t, model.NotificationPriorityHigh, view.Priority)
	assert.Equal(t, "push", view.Channel)
	assert.Equal(t, "order", view.ReferenceType)
	assert.Equal(t, uint64(123), *view.ReferenceID)
	assert.NotNil(t, view.ReadAt)
	assert.Equal(t, now, view.CreatedAt)
}

func TestListResponse_Structure(t *testing.T) {
	resp := ListResponse{
		Items: []NotificationView{
			{ID: 1, Title: "Notification 1"},
			{ID: 2, Title: "Notification 2"},
		},
		Page:        1,
		PageSize:    10,
		Total:       2,
		UnreadCount: 1,
	}

	assert.Len(t, resp.Items, 2)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, int64(1), resp.UnreadCount)
}

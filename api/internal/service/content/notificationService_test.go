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

// MockNotificationRepository is a mock implementation of NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	args := m.Called(ctx, opts)
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

func TestNewNotificationService(t *testing.T) {
	repo := &MockNotificationRepository{}
	svc := NewNotificationService(repo)
	require.NotNil(t, svc)
	assert.Equal(t, repo, svc.repo)
}

func TestNotificationService_ListNotifications(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with notifications", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		readAt := now.Add(-1 * time.Hour)
		refID := uint64(100)
		notifications := []model.NotificationEvent{
			{
				Base:          model.Base{ID: 1, CreatedAt: now},
				UserID:        1,
				Title:         "订单完成",
				Message:       "您的订单已完成",
				Priority:      model.NotificationPriorityNormal,
				Channel:       "web",
				ReferenceType: "order",
				ReferenceID:   &refID,
				ReadAt:        &readAt,
			},
			{
				Base:          model.Base{ID: 2, CreatedAt: now},
				UserID:        1,
				Title:         "新消息",
				Message:       "您有新消息",
				Priority:      model.NotificationPriorityHigh,
				Channel:       "push",
				ReferenceType: "chat",
			},
		}

		repo.On("ListByUser", ctx, mock.AnythingOfType("repository.NotificationListOptions")).Return(notifications, int64(2), nil)
		repo.On("CountUnread", ctx, uint64(1)).Return(int64(1), nil)

		req := NotificationListRequest{Page: 1, PageSize: 10}
		result, err := svc.ListNotifications(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, int64(1), result.UnreadCount)
		assert.True(t, result.Items[0].IsRead)
		assert.False(t, result.Items[1].IsRead)
	})

	t.Run("success with unread only filter", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		notifications := []model.NotificationEvent{
			{Base: model.Base{ID: 1, CreatedAt: now}, Title: "Unread"},
		}

		repo.On("ListByUser", ctx, mock.MatchedBy(func(opts repository.NotificationListOptions) bool {
			return opts.Unread != nil && *opts.Unread == true
		})).Return(notifications, int64(1), nil)
		repo.On("CountUnread", ctx, uint64(1)).Return(int64(1), nil)

		req := NotificationListRequest{Page: 1, PageSize: 10, UnreadOnly: true}
		result, err := svc.ListNotifications(ctx, 1, req)

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
	})

	t.Run("repo list error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("ListByUser", ctx, mock.AnythingOfType("repository.NotificationListOptions")).Return([]model.NotificationEvent{}, int64(0), errors.New("db error"))

		req := NotificationListRequest{Page: 1, PageSize: 10}
		_, err := svc.ListNotifications(ctx, 1, req)

		require.Error(t, err)
	})

	t.Run("repo count error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("ListByUser", ctx, mock.AnythingOfType("repository.NotificationListOptions")).Return([]model.NotificationEvent{}, int64(0), nil)
		repo.On("CountUnread", ctx, uint64(1)).Return(int64(0), errors.New("count error"))

		req := NotificationListRequest{Page: 1, PageSize: 10}
		_, err := svc.ListNotifications(ctx, 1, req)

		require.Error(t, err)
	})
}

func TestNotificationService_MarkNotificationsRead(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("MarkRead", ctx, uint64(1), []uint64{1, 2, 3}).Return(nil)

		err := svc.MarkNotificationsRead(ctx, 1, []uint64{1, 2, 3})

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("MarkRead", ctx, uint64(1), []uint64{1}).Return(errors.New("db error"))

		err := svc.MarkNotificationsRead(ctx, 1, []uint64{1})

		require.Error(t, err)
	})
}

func TestNotificationService_MarkAllNotificationsRead(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("MarkAllRead", ctx, uint64(1)).Return(nil)

		err := svc.MarkAllNotificationsRead(ctx, 1)

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("MarkAllRead", ctx, uint64(1)).Return(errors.New("db error"))

		err := svc.MarkAllNotificationsRead(ctx, 1)

		require.Error(t, err)
	})
}

func TestNotificationService_GetUnreadNotificationCount(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("CountUnread", ctx, uint64(1)).Return(int64(5), nil)

		count, err := svc.GetUnreadNotificationCount(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("CountUnread", ctx, uint64(1)).Return(int64(0), errors.New("db error"))

		_, err := svc.GetUnreadNotificationCount(ctx, 1)

		require.Error(t, err)
	})
}

func TestNotificationService_DeleteNotification(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("Delete", ctx, uint64(1), uint64(100)).Return(nil)

		err := svc.DeleteNotification(ctx, 1, 100)

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockNotificationRepository{}
		svc := NewNotificationService(repo)

		repo.On("Delete", ctx, uint64(1), uint64(100)).Return(errors.New("db error"))

		err := svc.DeleteNotification(ctx, 1, 100)

		require.Error(t, err)
	})
}

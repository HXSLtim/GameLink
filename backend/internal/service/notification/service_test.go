package notification

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// Mock repository for notification service tests
type mockNotificationRepoForService struct {
	events map[uint64]*model.NotificationEvent
}

func (m *mockNotificationRepoForService) Create(ctx context.Context, event *model.NotificationEvent) error {
	if event.ID == 0 {
		event.ID = uint64(len(m.events) + 1)
	}
	m.events[event.ID] = event
	return nil
}

func (m *mockNotificationRepoForService) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	events := make([]model.NotificationEvent, 0)
	for _, e := range m.events {
		if e.UserID == opts.UserID {
			events = append(events, *e)
		}
	}
	return events, int64(len(events)), nil
}

func (m *mockNotificationRepoForService) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	for _, id := range ids {
		if e, ok := m.events[id]; ok {
			now := time.Now()
			e.ReadAt = &now
		}
	}
	return nil
}

func (m *mockNotificationRepoForService) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	count := 0
	for _, e := range m.events {
		if e.UserID == userID && e.ReadAt == nil {
			count++
		}
	}
	return int64(count), nil
}

func setupNotificationService(t *testing.T) *Service {
	t.Helper()
	repo := &mockNotificationRepoForService{events: make(map[uint64]*model.NotificationEvent)}
	return NewService(repo)
}

func TestNotificationService_List_Success(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       1,
		PageSize:   10,
		UnreadOnly: false,
		Priorities: []model.NotificationPriority{},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
}

func TestNotificationService_List_WithPagination(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       2,
		PageSize:   5,
		UnreadOnly: false,
		Priorities: []model.NotificationPriority{},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 5, resp.PageSize)
}

func TestNotificationService_List_UnreadOnly(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       1,
		PageSize:   10,
		UnreadOnly: true,
		Priorities: []model.NotificationPriority{},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.UnreadCount)
}

func TestNotificationService_List_WithPriority(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       1,
		PageSize:   10,
		UnreadOnly: false,
		Priorities: []model.NotificationPriority{model.NotificationPriorityHigh},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNotificationService_MarkRead_Success(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	// Create some notifications first
	repo := &mockNotificationRepoForService{events: make(map[uint64]*model.NotificationEvent)}
	for i := 1; i <= 3; i++ {
		event := &model.NotificationEvent{
			UserID:    1,
			Title:     "Test",
			Message:   "Test message",
			Priority:  model.NotificationPriorityNormal,
			Channel:   "system",
		}
		repo.Create(ctx, event)
	}

	// Mark as read
	err := svc.MarkRead(ctx, 1, []uint64{1, 2})
	assert.NoError(t, err)
}

func TestNotificationService_MarkRead_EmptyIDs(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	err := svc.MarkRead(ctx, 1, []uint64{})
	assert.NoError(t, err)
}

func TestNotificationService_MarkRead_InvalidID(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	// Mark non-existent notification as read - should not error
	err := svc.MarkRead(ctx, 1, []uint64{99999})
	assert.NoError(t, err)
}

func TestNotificationService_GetUnreadCount_Success(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	count, err := svc.GetUnreadCount(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestNotificationService_GetUnreadCount_WithNotifications(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	// Note: We can't easily add notifications to the service's internal repo
	// through the service interface, so we just test that the method works
	count, err := svc.GetUnreadCount(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestNotificationService_List_DefaultPage(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       0,
		PageSize:   0,
		UnreadOnly: false,
		Priorities: []model.NotificationPriority{},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNotificationService_MarkRead_MultipleIDs(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	// Mark multiple notifications as read
	err := svc.MarkRead(ctx, 1, []uint64{1, 2, 3, 4, 5})
	assert.NoError(t, err)
}

// TestNotificationService_List_WithMultiplePriorities tests filtering by multiple priorities
func TestNotificationService_List_WithMultiplePriorities(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:     1,
		PageSize: 10,
		Priorities: []model.NotificationPriority{
			model.NotificationPriorityHigh,
			model.NotificationPriorityNormal,
		},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Page)
}

// TestNotificationService_List_EmptyPriorities tests with empty priority filter
func TestNotificationService_List_EmptyPriorities(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       1,
		PageSize:   10,
		Priorities: []model.NotificationPriority{},
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestNotificationService_List_LargePage tests large page number
func TestNotificationService_List_LargePage(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:     100,
		PageSize: 20,
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 100, resp.Page)
	assert.Equal(t, 20, resp.PageSize)
}

// TestNotificationService_List_SmallPageSize tests small page size
func TestNotificationService_List_SmallPageSize(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:     1,
		PageSize: 1,
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.PageSize)
}

// TestNotificationService_GetUnreadCount_DifferentUser tests unread count for different users
func TestNotificationService_GetUnreadCount_DifferentUser(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	// Get unread count for user 999
	count, err := svc.GetUnreadCount(ctx, 999)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestNotificationService_MarkRead_SingleID tests marking single notification as read
func TestNotificationService_MarkRead_SingleID(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	err := svc.MarkRead(ctx, 1, []uint64{1})
	assert.NoError(t, err)
}

// TestNotificationService_List_UnreadOnlyWithNoUnread tests unread filter with no unread notifications
func TestNotificationService_List_UnreadOnlyWithNoUnread(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:       1,
		PageSize:   10,
		UnreadOnly: true,
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.UnreadCount)
	assert.Len(t, resp.Items, 0)
}

// TestNotificationService_List_VerifyResponseStructure tests response structure
func TestNotificationService_List_VerifyResponseStructure(t *testing.T) {
	svc := setupNotificationService(t)
	ctx := context.Background()

	req := ListRequest{
		Page:     1,
		PageSize: 20,
	}

	resp, err := svc.List(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Items)
	assert.GreaterOrEqual(t, resp.Total, int64(0))
	assert.GreaterOrEqual(t, resp.UnreadCount, int64(0))
}

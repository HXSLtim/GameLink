// Package integration provides integration tests for notification service.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	notificationrepo "gamelink/internal/repository/notification"
	"gamelink/internal/service/notification"
)

// setupNotificationService creates a NotificationService with real database repository.
func setupNotificationService(t *testing.T) (*notification.Service, context.Context) {
	t.Helper()
	db := SetupTestDB(t)
	ctx := context.Background()

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)
	return svc, ctx
}

// createTestNotificationEvent creates a NotificationEvent directly in the database.
func createTestNotificationEvent(t *testing.T, db interface {
	Create(interface{}) interface{ Error() error }
}, userID uint64, title, message string, priority model.NotificationPriority) *model.NotificationEvent {
	t.Helper()
	event := &model.NotificationEvent{
		Base: model.Base{
			ExtJSON: "{}",
		},
		UserID:        userID,
		Title:         title,
		Message:       message,
		Channel:       "web",
		Priority:      priority,
		ReferenceType: "test",
	}
	// Use the global db from SetupTestDB
	return event
}

// TestNotificationService_List tests listing notifications.
func TestNotificationService_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create notifications directly
	for i := 0; i < 5; i++ {
		event := &model.NotificationEvent{
			Base: model.Base{
				ExtJSON: "{}",
			},
			UserID:        user.ID,
			Title:         "Notification " + string(rune('A'+i)),
			Message:       "Message content",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	// Create service with fresh repo
	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// List notifications
	req := notification.ListRequest{
		Page:     1,
		PageSize: 10,
	}
	resp, err := svc.List(ctx, user.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(5), resp.Total)
	assert.Len(t, resp.Items, 5)
	assert.Equal(t, int64(5), resp.UnreadCount)
}

// TestNotificationService_List_UnreadOnly tests listing only unread notifications.
func TestNotificationService_List_UnreadOnly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create 3 unread and 2 read notifications
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Unread " + string(rune('A'+i)),
			Message:       "Unread message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	now := time.Now()
	for i := 0; i < 2; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Read " + string(rune('A'+i)),
			Message:       "Read message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
			ReadAt:        &now,
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// List only unread
	req := notification.ListRequest{
		Page:       1,
		PageSize:   10,
		UnreadOnly: true,
	}
	resp, err := svc.List(ctx, user.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Items, 3)
}

// TestNotificationService_List_ByPriority tests filtering by priority.
func TestNotificationService_List_ByPriority(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create notifications with different priorities
	priorities := []model.NotificationPriority{
		model.NotificationPriorityLow,
		model.NotificationPriorityNormal,
		model.NotificationPriorityHigh,
	}
	for _, p := range priorities {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Priority " + string(p),
			Message:       "Message",
			Channel:       "web",
			Priority:      p,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// Filter by high priority only
	req := notification.ListRequest{
		Page:       1,
		PageSize:   10,
		Priorities: []model.NotificationPriority{model.NotificationPriorityHigh},
	}
	resp, err := svc.List(ctx, user.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(1), resp.Total)
}

// TestNotificationService_MarkRead tests marking notifications as read.
func TestNotificationService_MarkRead(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create notifications
	var ids []uint64
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Notification " + string(rune('A'+i)),
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
		ids = append(ids, event.ID)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// Mark first two as read
	err := svc.MarkRead(ctx, user.ID, ids[:2])
	require.NoError(t, err)

	// Verify unread count
	count, err := svc.GetUnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// TestNotificationService_MarkRead_AllNotifications tests marking all as read.
func TestNotificationService_MarkRead_AllNotifications(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create notifications
	for i := 0; i < 5; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Notification " + string(rune('A'+i)),
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// Mark all as read (empty slice marks all)
	err := svc.MarkRead(ctx, user.ID, []uint64{})
	require.NoError(t, err)

	// Verify all are read
	count, err := svc.GetUnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestNotificationService_GetUnreadCount tests getting unread count.
func TestNotificationService_GetUnreadCount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create 5 unread notifications
	for i := 0; i < 5; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Notification " + string(rune('A'+i)),
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	count, err := svc.GetUnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}

// TestNotificationService_GetUnreadCount_NoNotifications tests count with no notifications.
func TestNotificationService_GetUnreadCount_NoNotifications(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	count, err := svc.GetUnreadCount(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestNotificationService_List_Pagination tests pagination.
func TestNotificationService_List_Pagination(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user := CreateUniqueTestUser(t, db, "notif_user")

	// Create 15 notifications
	for i := 0; i < 15; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user.ID,
			Title:         "Notification " + string(rune('A'+i)),
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// Get first page
	req := notification.ListRequest{
		Page:     1,
		PageSize: 10,
	}
	resp, err := svc.List(ctx, user.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(15), resp.Total)
	assert.Len(t, resp.Items, 10)

	// Get second page
	req.Page = 2
	resp, err = svc.List(ctx, user.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(15), resp.Total)
	assert.Len(t, resp.Items, 5)
}

// TestNotificationService_List_DifferentUsers tests isolation between users.
func TestNotificationService_List_DifferentUsers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	_, ctx := setupNotificationService(t)

	user1 := CreateUniqueTestUser(t, db, "notif_user1")
	user2 := CreateUniqueTestUser(t, db, "notif_user2")

	// Create notifications for user1
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user1.ID,
			Title:         "User1 Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	// Create notifications for user2
	for i := 0; i < 5; i++ {
		event := &model.NotificationEvent{
			Base:          model.Base{ExtJSON: "{}"},
			UserID:        user2.ID,
			Title:         "User2 Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "test",
		}
		err := db.Create(event).Error
		require.NoError(t, err)
	}

	repo := notificationrepo.NewNotificationRepository(db)
	svc := notification.NewService(repo)

	// User1 should only see their notifications
	req := notification.ListRequest{Page: 1, PageSize: 10}
	resp1, err := svc.List(ctx, user1.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(3), resp1.Total)

	// User2 should only see their notifications
	resp2, err := svc.List(ctx, user2.ID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(5), resp2.Total)
}

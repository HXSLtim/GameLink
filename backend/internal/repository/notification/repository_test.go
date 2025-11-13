package notification

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func setupNotificationTest(t *testing.T) repository.NotificationRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.NotificationEvent{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return NewNotificationRepository(db)
}

func TestNotificationRepository_Create(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	event := &model.NotificationEvent{
		UserID:        1,
		Title:         "Test Notification",
		Message:       "This is a test",
		Channel:       "web",
		Priority:      model.NotificationPriorityNormal,
		ReferenceType: "order",
	}

	err := repo.Create(ctx, event)
	assert.NoError(t, err)
	assert.NotZero(t, event.ID)
}

func TestNotificationRepository_ListByUser(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	// Create multiple notifications for different users
	for i := 0; i < 5; i++ {
		event := &model.NotificationEvent{
			UserID:        uint64(i%2 + 1), // Alternate between user 1 and 2
			Title:         "Notification " + string(rune(i)),
			Message:       "Message " + string(rune(i)),
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	}

	// List notifications for user 1
	events, total, err := repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   1,
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total) // 3 notifications for user 1
	assert.Len(t, events, 3)

	// List notifications for user 2
	events, total, err = repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   2,
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // 2 notifications for user 2
	assert.Len(t, events, 2)
}

func TestNotificationRepository_MarkRead(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	// Create notifications
	var ids []uint64
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			UserID:        1,
			Title:         "Notification " + string(rune(i)),
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
		ids = append(ids, event.ID)
	}

	// Mark as read
	err := repo.MarkRead(ctx, 1, ids[:2])
	assert.NoError(t, err)

	// Verify
	events, _, err := repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   1,
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	readCount := 0
	for _, e := range events {
		if e.ReadAt != nil {
			readCount++
		}
	}
	assert.Equal(t, 2, readCount)
}

func TestNotificationRepository_CountUnread(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	// Create notifications
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			UserID:        1,
			Title:         "Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	}

	// Count unread
	count, err := repo.CountUnread(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestNotificationRepository_ListByUserWithPagination(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	// Create many notifications
	for i := 0; i < 25; i++ {
		event := &model.NotificationEvent{
			UserID:        1,
			Title:         "Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityNormal,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	}

	// List with pagination
	events, total, err := repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   1,
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(25), total)
	assert.Len(t, events, 10)

	// Get second page
	events, total, err = repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   1,
		Page:     2,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(25), total)
	assert.Len(t, events, 10)
}

func TestNotificationRepository_ListByUserWithPriority(t *testing.T) {
	repo := setupNotificationTest(t)
	ctx := context.Background()

	// Create notifications with different priorities
	for i := 0; i < 3; i++ {
		event := &model.NotificationEvent{
			UserID:        1,
			Title:         "Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityHigh,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	}

	for i := 0; i < 2; i++ {
		event := &model.NotificationEvent{
			UserID:        1,
			Title:         "Notification",
			Message:       "Message",
			Channel:       "web",
			Priority:      model.NotificationPriorityLow,
			ReferenceType: "order",
		}
		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	}

	// List high priority
	events, total, err := repo.ListByUser(ctx, repository.NotificationListOptions{
		UserID:   1,
		Priority: []model.NotificationPriority{model.NotificationPriorityHigh},
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, events, 3)
}

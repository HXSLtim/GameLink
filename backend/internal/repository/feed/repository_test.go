package feed

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func setupFeedTest(t *testing.T) repository.FeedRepository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&model.Feed{}, &model.FeedImage{}, &model.FeedReport{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return NewFeedRepository(db)
}

func TestFeedRepository_Create(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	feed := &model.Feed{
		AuthorID:           1,
		Content:            "Hello world",
		Visibility:         model.FeedVisibilityPublic,
		ModerationStatus:   model.FeedModerationPending,
	}

	err := repo.Create(ctx, feed)
	assert.NoError(t, err)
	assert.NotZero(t, feed.ID)
}

func TestFeedRepository_Get(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create a feed
	feed := &model.Feed{
		AuthorID:           1,
		Content:            "Hello world",
		Visibility:         model.FeedVisibilityPublic,
		ModerationStatus:   model.FeedModerationPending,
	}
	err := repo.Create(ctx, feed)
	assert.NoError(t, err)

	// Get the feed
	retrieved, err := repo.Get(ctx, feed.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, feed.AuthorID, retrieved.AuthorID)
	assert.Equal(t, feed.Content, retrieved.Content)

	// Get non-existent feed
	_, err = repo.Get(ctx, 99999)
	assert.Equal(t, repository.ErrNotFound, err)
}

func TestFeedRepository_List(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create multiple feeds
	for i := 0; i < 5; i++ {
		feed := &model.Feed{
			AuthorID:           uint64(i + 1),
			Content:            "Content " + string(rune(i)),
			Visibility:         model.FeedVisibilityPublic,
			ModerationStatus:   model.FeedModerationPending,
		}
		err := repo.Create(ctx, feed)
		assert.NoError(t, err)
	}

	// List all feeds
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit: 10,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 5)

	// List with limit
	feeds, err = repo.List(ctx, repository.FeedListOptions{
		Limit: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 2)

	// List with author filter
	feeds, err = repo.List(ctx, repository.FeedListOptions{
		AuthorID: ptrUint64(1),
		Limit:    10,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)

	// List with visibility filter
	feeds, err = repo.List(ctx, repository.FeedListOptions{
		Visibility: []model.FeedVisibility{model.FeedVisibilityPublic},
		Limit:      10,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 5)

	// List approved only
	feeds, err = repo.List(ctx, repository.FeedListOptions{
		OnlyApproved: true,
		Limit:        10,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 0)
}

func TestFeedRepository_UpdateModeration(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create a feed
	feed := &model.Feed{
		AuthorID:           1,
		Content:            "Hello world",
		Visibility:         model.FeedVisibilityPublic,
		ModerationStatus:   model.FeedModerationPending,
	}
	err := repo.Create(ctx, feed)
	assert.NoError(t, err)

	// Update moderation status (auto)
	err = repo.UpdateModeration(ctx, feed.ID, model.FeedModerationApproved, "Auto approved", false)
	assert.NoError(t, err)

	// Verify update
	retrieved, err := repo.Get(ctx, feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.FeedModerationApproved, retrieved.ModerationStatus)
	assert.Equal(t, "Auto approved", retrieved.ModerationNote)
	assert.NotNil(t, retrieved.AutoModeratedAt)

	// Update moderation status (manual)
	err = repo.UpdateModeration(ctx, feed.ID, model.FeedModerationRejected, "Manual rejected", true)
	assert.NoError(t, err)

	// Verify update
	retrieved, err = repo.Get(ctx, feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.FeedModerationRejected, retrieved.ModerationStatus)
	assert.Equal(t, "Manual rejected", retrieved.ModerationNote)
	assert.NotNil(t, retrieved.ManualModeratedAt)
}

func TestFeedRepository_CreateReport(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create a feed first
	feed := &model.Feed{
		AuthorID:           1,
		Content:            "Hello world",
		Visibility:         model.FeedVisibilityPublic,
		ModerationStatus:   model.FeedModerationPending,
	}
	err := repo.Create(ctx, feed)
	assert.NoError(t, err)

	// Create a report
	report := &model.FeedReport{
		FeedID:   feed.ID,
		Reporter: 2,
		Reason:   "Inappropriate content",
		Status:   "pending",
	}

	err = repo.CreateReport(ctx, report)
	assert.NoError(t, err)
	assert.NotZero(t, report.ID)
}

func TestFeedRepository_ListDefaultPageSize(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create more feeds than default page size
	for i := 0; i < 25; i++ {
		feed := &model.Feed{
			AuthorID:           1,
			Content:            "Content",
			Visibility:         model.FeedVisibilityPublic,
			ModerationStatus:   model.FeedModerationPending,
		}
		err := repo.Create(ctx, feed)
		assert.NoError(t, err)
	}

	// List with default limit (should be 20)
	feeds, err := repo.List(ctx, repository.FeedListOptions{})
	assert.NoError(t, err)
	assert.Len(t, feeds, 20)
}

func TestFeedRepository_ListMaxPageSize(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create more feeds than max page size
	for i := 0; i < 60; i++ {
		feed := &model.Feed{
			AuthorID:           1,
			Content:            "Content",
			Visibility:         model.FeedVisibilityPublic,
			ModerationStatus:   model.FeedModerationPending,
		}
		err := repo.Create(ctx, feed)
		assert.NoError(t, err)
	}

	// List with limit > max (should be capped at 50)
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit: 100,
	})
	assert.NoError(t, err)
	assert.Len(t, feeds, 50)
}

func TestFeedRepository_ListCursorPagination(t *testing.T) {
	repo := setupFeedTest(t)
	ctx := context.Background()

	// Create feeds
	var feedIDs []uint64
	for i := 0; i < 5; i++ {
		feed := &model.Feed{
			AuthorID:           1,
			Content:            "Content",
			Visibility:         model.FeedVisibilityPublic,
			ModerationStatus:   model.FeedModerationPending,
		}
		err := repo.Create(ctx, feed)
		assert.NoError(t, err)
		feedIDs = append(feedIDs, feed.ID)
	}

	// List with cursor
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		CursorBefore: &feedIDs[2],
		Limit:        10,
	})
	assert.NoError(t, err)
	// Should return feeds with ID < feedIDs[2]
	assert.True(t, len(feeds) <= 2)
}

// Helper function
func ptrUint64(v uint64) *uint64 {
	return &v
}

// Package integration provides integration tests for Feed/Content module.
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/feed"
)

// ============================================================================
// Feed/Content CRUD Tests
// ============================================================================

// TestFeedRepository_Create tests creating a new feed
func TestFeedRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_author")

	feed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "This is a test feed content",
		Visibility: model.FeedVisibilityPublic,
	}

	err := repo.Create(ctx, feed)
	require.NoError(t, err)
	assert.NotZero(t, feed.ID)
	assert.NotZero(t, feed.CreatedAt)
}

// TestFeedRepository_Get tests retrieving a feed by ID
func TestFeedRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_get_author")

	feed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "Test feed for get",
		Visibility: model.FeedVisibilityPublic,
	}
	require.NoError(t, repo.Create(ctx, feed))

	got, err := repo.Get(ctx, feed.ID)
	require.NoError(t, err)
	assert.Equal(t, feed.Content, got.Content)
	assert.Equal(t, feed.AuthorID, got.AuthorID)
	assert.Equal(t, model.FeedVisibilityPublic, got.Visibility)
}

// TestFeedRepository_Get_NonExistent tests getting non-existent feed
func TestFeedRepository_Get_NonExistent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	_, err := repo.Get(ctx, 99999)
	assert.Error(t, err)
}

// TestFeedRepository_Update tests updating a feed
func TestFeedRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_update_author")

	feed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "Original content",
		Visibility: model.FeedVisibilityPublic,
	}
	require.NoError(t, repo.Create(ctx, feed))

	// Update content
	feed.Content = "Updated content"
	err := repo.Update(ctx, feed)
	require.NoError(t, err)

	got, err := repo.Get(ctx, feed.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated content", got.Content)
}

// TestFeedRepository_Delete tests deleting a feed
func TestFeedRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_delete_author")

	feed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "Feed to delete",
		Visibility: model.FeedVisibilityPublic,
	}
	require.NoError(t, repo.Create(ctx, feed))

	err := repo.Delete(ctx, feed.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, feed.ID)
	assert.Error(t, err)
}

// TestFeedRepository_List tests listing feeds with cursor pagination
func TestFeedRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_list_author")

	// Create multiple feeds
	for i := 0; i < 5; i++ {
		f := &model.Feed{
			AuthorID:          user.ID,
			Content:           fmt.Sprintf("Feed %d", i),
			Visibility:        model.FeedVisibilityPublic,
			ModerationStatus:  model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// List feeds
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        10,
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(feeds), 5)
}

// TestFeedRepository_List_WithCursor tests cursor-based pagination
func TestFeedRepository_List_WithCursor(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_cursor_author")

	// Create feeds
	var feedIDs []uint64
	for i := 0; i < 5; i++ {
		f := &model.Feed{
			AuthorID:         user.ID,
			Content:          fmt.Sprintf("Cursor feed %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, f))
		feedIDs = append(feedIDs, f.ID)
	}

	// Get first page
	feeds1, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        3,
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.Len(t, feeds1, 3)

	// Get second page with cursor
	cursor := feeds1[len(feeds1)-1].ID
	feeds2, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        3,
		CursorBefore: &cursor,
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(feeds2), 2)
}

// TestFeedRepository_List_ByAuthor tests filtering by author
func TestFeedRepository_List_ByAuthor(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user1 := CreateUniqueTestUser(t, db, "feed_author1")
	user2 := CreateUniqueTestUser(t, db, "feed_author2")

	// Create feeds for user1
	for i := 0; i < 3; i++ {
		f := &model.Feed{
			AuthorID:         user1.ID,
			Content:          fmt.Sprintf("User1 feed %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// Create feeds for user2
	for i := 0; i < 2; i++ {
		f := &model.Feed{
			AuthorID:         user2.ID,
			Content:          fmt.Sprintf("User2 feed %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// List only user1's feeds
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        10,
		AuthorID:     &user1.ID,
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(feeds), 3)
	for _, f := range feeds {
		assert.Equal(t, user1.ID, f.AuthorID)
	}
}

// TestFeedRepository_List_ByVisibility tests filtering by visibility
func TestFeedRepository_List_ByVisibility(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_visibility_author")

	// Create feeds with different visibility
	visibilities := []model.FeedVisibility{
		model.FeedVisibilityPublic,
		model.FeedVisibilityFollowers,
		model.FeedVisibilityPrivate,
	}

	for _, vis := range visibilities {
		f := &model.Feed{
			AuthorID:         user.ID,
			Content:          fmt.Sprintf("Feed with visibility %s", vis),
			Visibility:       vis,
			ModerationStatus: model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// List only public feeds
	publicVis := model.FeedVisibilityPublic
	feeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        10,
		Visibility:   []model.FeedVisibility{publicVis},
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(feeds), 1)
}

// TestFeedRepository_ListPaged tests paginated listing
func TestFeedRepository_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_paged_author")

	// Create multiple feeds
	for i := 0; i < 15; i++ {
		f := &model.Feed{
			AuthorID:   user.ID,
			Content:    fmt.Sprintf("Paged feed %d", i),
			Visibility: model.FeedVisibilityPublic,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// Get first page
	feeds1, total, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, feeds1, 10)

	// Get second page
	feeds2, _, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		Page:     2,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.Len(t, feeds2, 5)
}

// TestFeedRepository_ListPaged_ByKeyword tests searching by keyword
func TestFeedRepository_ListPaged_ByKeyword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_keyword_author")

	// Create feeds with specific keyword
	uniqueKeyword := fmt.Sprintf("unique_%d", time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		f := &model.Feed{
			AuthorID:   user.ID,
			Content:    fmt.Sprintf("Feed with %s keyword %d", uniqueKeyword, i),
			Visibility: model.FeedVisibilityPublic,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// Search by keyword
	feeds, total, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		Keyword:  uniqueKeyword,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	for _, f := range feeds {
		assert.Contains(t, f.Content, uniqueKeyword)
	}
}

// TestFeedRepository_ListPaged_ByModerationStatus tests filtering by moderation status
func TestFeedRepository_ListPaged_ByModerationStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_mod_status_author")

	// Create feeds with different moderation statuses
	statuses := []model.FeedModerationStatus{
		model.FeedModerationPending,
		model.FeedModerationApproved,
		model.FeedModerationRejected,
	}

	for _, status := range statuses {
		f := &model.Feed{
			AuthorID:          user.ID,
			Content:           fmt.Sprintf("Feed with status %s", status),
			Visibility:        model.FeedVisibilityPublic,
			ModerationStatus:  status,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// List only pending feeds
	pendingStatus := model.FeedModerationPending
	feeds, total, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		ModerationStatus: &pendingStatus,
		Page:             1,
		PageSize:         10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	for _, f := range feeds {
		assert.Equal(t, model.FeedModerationPending, f.ModerationStatus)
	}
}

// TestFeedRepository_UpdateModeration tests updating moderation status
func TestFeedRepository_UpdateModeration(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_mod_update_author")
	moderator := CreateUniqueTestUser(t, db, "feed_moderator")

	feed := &model.Feed{
		AuthorID:          user.ID,
		Content:           "Feed to moderate",
		Visibility:        model.FeedVisibilityPublic,
		ModerationStatus:  model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, feed))

	// Approve feed
	err := repo.UpdateModeration(ctx, feed.ID, model.FeedModerationApproved, "Approved by moderator", &moderator.ID)
	require.NoError(t, err)

	got, err := repo.Get(ctx, feed.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationApproved, got.ModerationStatus)
	assert.Equal(t, "Approved by moderator", got.ModerationNote)
	assert.NotNil(t, got.ModeratedBy)
	assert.NotNil(t, got.ManualModeratedAt)
}

// TestFeedRepository_UpdateModeration_AutoModeration tests auto-moderation (without moderator)
func TestFeedRepository_UpdateModeration_AutoModeration(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_auto_mod_author")

	feed := &model.Feed{
		AuthorID:          user.ID,
		Content:           "Feed for auto moderation",
		Visibility:        model.FeedVisibilityPublic,
		ModerationStatus:  model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, feed))

	// Auto-moderate (reject without moderator)
	err := repo.UpdateModeration(ctx, feed.ID, model.FeedModerationRejected, "Auto-rejected by system", nil)
	require.NoError(t, err)

	got, err := repo.Get(ctx, feed.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationRejected, got.ModerationStatus)
	assert.NotNil(t, got.AutoModeratedAt)
}

// TestFeedRepository_BatchUpdateModeration tests batch moderation update
func TestFeedRepository_BatchUpdateModeration(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_batch_mod_author")
	moderator := CreateUniqueTestUser(t, db, "feed_batch_moderator")

	// Create feeds pending moderation
	var feedIDs []uint64
	for i := 0; i < 3; i++ {
		f := &model.Feed{
			AuthorID:          user.ID,
			Content:           fmt.Sprintf("Pending feed %d", i),
			Visibility:        model.FeedVisibilityPublic,
			ModerationStatus:  model.FeedModerationPending,
		}
		require.NoError(t, repo.Create(ctx, f))
		feedIDs = append(feedIDs, f.ID)
	}

	// Batch approve
	err := repo.BatchUpdateModeration(ctx, feedIDs, model.FeedModerationApproved, "Batch approved", &moderator.ID)
	require.NoError(t, err)

	// Verify all are approved
	for _, id := range feedIDs {
		feed, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationApproved, feed.ModerationStatus)
	}
}

// TestFeedRepository_CreateReport tests creating a feed report
func TestFeedRepository_CreateReport(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "feed_report_author")
	reporter := CreateUniqueTestUser(t, db, "feed_reporter")

	feed := &model.Feed{
		AuthorID:          author.ID,
		Content:           "Feed to be reported",
		Visibility:        model.FeedVisibilityPublic,
		ModerationStatus:  model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, feed))

	report := &model.FeedReport{
		FeedID:   feed.ID,
		Reporter: reporter.ID,
		Reason:   "Inappropriate content",
		Status:   "pending",
	}

	err := repo.CreateReport(ctx, report)
	require.NoError(t, err)
	assert.NotZero(t, report.ID)
}

// TestFeedRepository_ListReports tests listing feed reports
func TestFeedRepository_ListReports(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "feed_list_report_author")
	reporter := CreateUniqueTestUser(t, db, "feed_list_reporter")

	feed := &model.Feed{
		AuthorID:          author.ID,
		Content:           "Feed with reports",
		Visibility:        model.FeedVisibilityPublic,
		ModerationStatus:  model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, feed))

	// Create multiple reports
	for i := 0; i < 3; i++ {
		report := &model.FeedReport{
			FeedID:   feed.ID,
			Reporter: reporter.ID,
			Reason:   fmt.Sprintf("Report reason %d", i),
			Status:   "pending",
		}
		require.NoError(t, repo.CreateReport(ctx, report))
	}

	// List reports for this feed
	reports, total, err := repo.ListReports(ctx, repository.FeedReportListOptions{
		FeedID:   &feed.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	_ = reports // reports is used for assertion
}

// TestFeedRepository_CountByStatus tests counting feeds by status
func TestFeedRepository_CountByStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_count_author")

	// Create feeds with different statuses
	statuses := []model.FeedModerationStatus{
		model.FeedModerationPending,
		model.FeedModerationApproved,
		model.FeedModerationApproved,
		model.FeedModerationRejected,
	}

	for _, status := range statuses {
		f := &model.Feed{
			AuthorID:          user.ID,
			Content:           fmt.Sprintf("Count test feed %s", status),
			Visibility:        model.FeedVisibilityPublic,
			ModerationStatus:  status,
		}
		require.NoError(t, repo.Create(ctx, f))
	}

	// Get counts
	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, counts[model.FeedModerationApproved], int64(2))
	assert.GreaterOrEqual(t, counts[model.FeedModerationPending], int64(1))
	assert.GreaterOrEqual(t, counts[model.FeedModerationRejected], int64(1))
}

// TestFeedRepository_WithImages tests creating feeds with images
func TestFeedRepository_WithImages(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_images_author")

	feed := &model.Feed{
		AuthorID:          user.ID,
		Content:           "Feed with images",
		Visibility:        model.FeedVisibilityPublic,
		ModerationStatus:  model.FeedModerationApproved,
		Images: []model.FeedImage{
			{
				URL:       "https://example.com/image1.jpg",
				Width:     800,
				Height:    600,
				SizeBytes: 102400,
				Order:     0,
			},
			{
				URL:       "https://example.com/image2.jpg",
				Width:     1024,
				Height:    768,
				SizeBytes: 204800,
				Order:     1,
			},
		},
	}

	err := repo.Create(ctx, feed)
	require.NoError(t, err)

	// Get feed with images
	got, err := repo.Get(ctx, feed.ID)
	require.NoError(t, err)
	assert.Len(t, got.Images, 2)
	assert.Equal(t, "https://example.com/image1.jpg", got.Images[0].URL)
}

// TestFeedRepository_AllVisibilityTypes tests all visibility types
func TestFeedRepository_AllVisibilityTypes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_all_vis_author")

	visibilities := []model.FeedVisibility{
		model.FeedVisibilityPublic,
		model.FeedVisibilityFollowers,
		model.FeedVisibilityPrivate,
	}

	for _, vis := range visibilities {
		feed := &model.Feed{
			AuthorID:          user.ID,
			Content:           fmt.Sprintf("Feed with visibility %s", vis),
			Visibility:        vis,
			ModerationStatus:  model.FeedModerationApproved,
		}
		require.NoError(t, repo.Create(ctx, feed))
		assert.NotZero(t, feed.ID)
	}
}

// TestFeedRepository_ListPaged_DateRange tests filtering by date range
func TestFeedRepository_ListPaged_DateRange(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "feed_date_author")

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Create feed yesterday
	oldFeed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "Old feed",
		Visibility: model.FeedVisibilityPublic,
	}
	oldFeed.CreatedAt = yesterday
	require.NoError(t, repo.Create(ctx, oldFeed))

	// Create feed today
	newFeed := &model.Feed{
		AuthorID:   user.ID,
		Content:    "New feed",
		Visibility: model.FeedVisibilityPublic,
	}
	require.NoError(t, repo.Create(ctx, newFeed))

	// List feeds from today onwards
	today := now.Truncate(24 * time.Hour)
	feeds, total, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		DateFrom: &today,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(1))
	_ = feeds // feeds is used for assertion
}

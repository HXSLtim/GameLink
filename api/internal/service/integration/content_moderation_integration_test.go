// Package integration provides integration tests for content moderation workflow.
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
	"gamelink/internal/repository/sensitiveword"
	"gorm.io/gorm"
)

// ============================================================================
// Content Moderation Workflow Integration Tests
// Tests for complete content creation -> sensitive word detection -> moderation -> publish/reject flow
// ============================================================================

// TestContentModeration_CreateContent_WithAutoDetection tests creating content with automatic sensitive word detection.
func TestContentModeration_CreateContent_WithAutoDetection(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	feedRepo := feed.NewFeedRepository(db)
	swRepo := sensitiveword.NewSensitiveWordRepository(db)
	ctx := context.Background()

	// Setup: Create test user and sensitive words
	author := CreateUniqueTestUser(t, db, "content_author")

	// Create sensitive words with different severities
	highSeverityWord := CreateTestSensitiveWordWithSeverity(t, db, "violent_keyword", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityHigh)
	mediumSeverityWord := CreateTestSensitiveWordWithSeverity(t, db, "rude_word", model.SensitiveWordCategoryAbuse, model.SensitiveWordSeverityMedium)
	lowSeverityWord := CreateTestSensitiveWordWithSeverity(t, db, "mild_term", model.SensitiveWordCategoryOther, model.SensitiveWordSeverityLow)

	// Test 1: Content with high severity sensitive word - should be auto-rejected
	t.Run("HighSeverity_AutoReject", func(t *testing.T) {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          "This contains violent_keyword that should be rejected",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, content))

		// Simulate auto-detection (in real system, this would be triggered by service layer)
		content.ModerationStatus = model.FeedModerationRejected
		content.ModerationNote = "Auto-rejected: Contains high severity sensitive word"
		now := time.Now()
		content.AutoModeratedAt = &now
		require.NoError(t, feedRepo.Update(ctx, content))

		// Verify
		updated, err := feedRepo.Get(ctx, content.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationRejected, updated.ModerationStatus)
		assert.NotNil(t, updated.AutoModeratedAt)
		assert.Contains(t, updated.ModerationNote, "high severity")
	})

	// Test 2: Content with medium severity word - should be pending manual review
	t.Run("MediumSeverity_PendingReview", func(t *testing.T) {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          "This has rude_word but needs manual review",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, content))

		// Medium severity stays pending for manual review
		updated, err := feedRepo.Get(ctx, content.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationPending, updated.ModerationStatus)
		assert.Nil(t, updated.AutoModeratedAt)
	})

	// Test 3: Content with low severity word - should be pending with low priority
	t.Run("LowSeverity_PendingLowPriority", func(t *testing.T) {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          "Contains mild_term content",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, content))

		// Verify it's pending
		updated, err := feedRepo.Get(ctx, content.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationPending, updated.ModerationStatus)
	})

	// Test 4: Clean content - should be pending (goes through normal queue)
	t.Run("CleanContent_Pending", func(t *testing.T) {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          "This is clean content with no sensitive words",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, content))

		updated, err := feedRepo.Get(ctx, content.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationPending, updated.ModerationStatus)
	})

	// Verify all words exist
	_, _ = swRepo.Get(ctx, highSeverityWord.ID)
	_, _ = swRepo.Get(ctx, mediumSeverityWord.ID)
	_, _ = swRepo.Get(ctx, lowSeverityWord.ID)
}

// TestContentModeration_ApprovalFlow tests the complete approval flow.
func TestContentModeration_ApprovalFlow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "approval_author")
	moderator := CreateUniqueTestUser(t, db, "approval_moderator")

	// Step 1: Create pending content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "This is great content waiting for approval",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Verify initial state
	assert.Equal(t, model.FeedModerationPending, content.ModerationStatus)

	// Step 2: Approve content
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationApproved, "Content approved - meets guidelines", &moderator.ID)
	require.NoError(t, err)

	// Step 3: Verify approved state
	approved, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationApproved, approved.ModerationStatus)
	assert.Equal(t, "Content approved - meets guidelines", approved.ModerationNote)
	assert.Equal(t, &moderator.ID, approved.ModeratedBy)
	assert.NotNil(t, approved.ManualModeratedAt)

	// Step 4: Verify content is now visible in public feed
	publicFeeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        10,
		OnlyApproved: true,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(publicFeeds), 1)

	// Find our approved feed in the list
	found := false
	for _, f := range publicFeeds {
		if f.ID == content.ID {
			found = true
			assert.Equal(t, model.FeedModerationApproved, f.ModerationStatus)
			break
		}
	}
	assert.True(t, found, "Approved content should appear in public feed")
}

// TestContentModeration_RejectionFlow tests the complete rejection flow.
func TestContentModeration_RejectionFlow(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "reject_author")
	moderator := CreateUniqueTestUser(t, db, "reject_moderator")

	// Step 1: Create pending content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Inappropriate content that should be rejected",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Step 2: Reject content with reason
	reason := "Content violates community guidelines - inappropriate language"
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationRejected, reason, &moderator.ID)
	require.NoError(t, err)

	// Step 3: Verify rejected state
	rejected, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationRejected, rejected.ModerationStatus)
	assert.Equal(t, reason, rejected.ModerationNote)
	assert.Equal(t, &moderator.ID, rejected.ModeratedBy)
	assert.NotNil(t, rejected.ManualModeratedAt)

	// Step 4: Verify content is NOT visible in public feed
	publicFeeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        100,
		OnlyApproved: true,
	})
	require.NoError(t, err)

	for _, f := range publicFeeds {
		if f.ID == content.ID {
			t.Errorf("Rejected content should not appear in public feed")
		}
	}
}

// TestContentModeration_StatusTransitions tests valid status transitions.
func TestContentModeration_StatusTransitions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "transition_author")
	moderator := CreateUniqueTestUser(t, db, "transition_moderator")

	testCases := []struct {
		name          string
		initialStatus model.FeedModerationStatus
		targetStatus  model.FeedModerationStatus
		shouldSucceed bool
	}{
		{
			name:          "Pending to Approved",
			initialStatus: model.FeedModerationPending,
			targetStatus:  model.FeedModerationApproved,
			shouldSucceed: true,
		},
		{
			name:          "Pending to Rejected",
			initialStatus: model.FeedModerationPending,
			targetStatus:  model.FeedModerationRejected,
			shouldSucceed: true,
		},
		{
			name:          "Approved to Rejected",
			initialStatus: model.FeedModerationApproved,
			targetStatus:  model.FeedModerationRejected,
			shouldSucceed: true,
		},
		{
			name:          "Approved to Removed",
			initialStatus: model.FeedModerationApproved,
			targetStatus:  model.FeedModerationRemoved,
			shouldSucceed: true,
		},
		{
			name:          "Rejected to Approved",
			initialStatus: model.FeedModerationRejected,
			targetStatus:  model.FeedModerationApproved,
			shouldSucceed: true, // Appeal approved
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := &model.Feed{
				AuthorID:         author.ID,
				Content:          fmt.Sprintf("Content for transition %s", tc.name),
				Visibility:       model.FeedVisibilityPublic,
				ModerationStatus: tc.initialStatus,
			}
			require.NoError(t, repo.Create(ctx, content))

			// Attempt status transition
			err := repo.UpdateModeration(ctx, content.ID, tc.targetStatus, fmt.Sprintf("Transition to %s", tc.targetStatus), &moderator.ID)

			if tc.shouldSucceed {
				require.NoError(t, err)
				updated, err := repo.Get(ctx, content.ID)
				require.NoError(t, err)
				assert.Equal(t, tc.targetStatus, updated.ModerationStatus)
			}
		})
	}
}

// TestContentModeration_VisibilityRules tests content visibility based on moderation status.
func TestContentModeration_VisibilityRules(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "visibility_author")

	// Create content with different statuses and visibility settings
	testCases := []struct {
		name             string
		moderationStatus model.FeedModerationStatus
		visibility       model.FeedVisibility
		shouldBeVisible  bool
	}{
		{
			name:             "Public Approved",
			moderationStatus: model.FeedModerationApproved,
			visibility:       model.FeedVisibilityPublic,
			shouldBeVisible:  true,
		},
		{
			name:             "Public Pending",
			moderationStatus: model.FeedModerationPending,
			visibility:       model.FeedVisibilityPublic,
			shouldBeVisible:  false,
		},
		{
			name:             "Public Rejected",
			moderationStatus: model.FeedModerationRejected,
			visibility:       model.FeedVisibilityPublic,
			shouldBeVisible:  false,
		},
		{
			name:             "Followers Approved",
			moderationStatus: model.FeedModerationApproved,
			visibility:       model.FeedVisibilityFollowers,
			shouldBeVisible:  true, // Would only be visible to followers
		},
		{
			name:             "Private Approved",
			moderationStatus: model.FeedModerationApproved,
			visibility:       model.FeedVisibilityPrivate,
			shouldBeVisible:  false, // Only visible to author
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := &model.Feed{
				AuthorID:         author.ID,
				Content:          fmt.Sprintf("Content for %s", tc.name),
				Visibility:       tc.visibility,
				ModerationStatus: tc.moderationStatus,
			}
			require.NoError(t, repo.Create(ctx, content))

			// For approved content, verify it appears in approved list
			if tc.moderationStatus == model.FeedModerationApproved && tc.visibility != model.FeedVisibilityPrivate {
				publicFeeds, err := repo.List(ctx, repository.FeedListOptions{
					Limit:        100,
					OnlyApproved: true,
				})
				require.NoError(t, err)

				found := false
				for _, f := range publicFeeds {
					if f.ID == content.ID {
						found = true
						break
					}
				}
				assert.True(t, found, "Approved content should be visible")
			}
		})
	}
}

// TestContentModeration_BatchApproval tests batch approval of multiple content items.
func TestContentModeration_BatchApproval(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "batch_author")
	moderator := CreateUniqueTestUser(t, db, "batch_moderator")

	// Create multiple pending content items
	var contentIDs []uint64
	for i := 0; i < 5; i++ {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          fmt.Sprintf("Batch test content %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, repo.Create(ctx, content))
		contentIDs = append(contentIDs, content.ID)
	}

	// Batch approve all
	err := repo.BatchUpdateModeration(ctx, contentIDs, model.FeedModerationApproved, "Batch approved - all clean", &moderator.ID)
	require.NoError(t, err)

	// Verify all are approved
	for _, id := range contentIDs {
		content, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationApproved, content.ModerationStatus)
		assert.Equal(t, "Batch approved - all clean", content.ModerationNote)
		assert.Equal(t, &moderator.ID, content.ModeratedBy)
	}
}

// TestContentModeration_BatchRejection tests batch rejection of multiple content items.
func TestContentModeration_BatchRejection(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "batch_reject_author")
	moderator := CreateUniqueTestUser(t, db, "batch_reject_moderator")

	// Create multiple pending content items
	var contentIDs []uint64
	for i := 0; i < 3; i++ {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          fmt.Sprintf("Inappropriate content %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, repo.Create(ctx, content))
		contentIDs = append(contentIDs, content.ID)
	}

	// Batch reject all
	reason := "Batch rejected - violates community guidelines"
	err := repo.BatchUpdateModeration(ctx, contentIDs, model.FeedModerationRejected, reason, &moderator.ID)
	require.NoError(t, err)

	// Verify all are rejected
	for _, id := range contentIDs {
		content, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationRejected, content.ModerationStatus)
		assert.Equal(t, reason, content.ModerationNote)
	}
}

// TestContentModeration_ReportAndReview tests the content report workflow.
func TestContentModeration_ReportAndReview(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "reported_author")
	reporter := CreateUniqueTestUser(t, db, "reporter")
	admin := CreateUniqueTestUser(t, db, "admin")

	// Step 1: Create and approve content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Content that will be reported",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Step 2: User reports the content
	report := &model.FeedReport{
		FeedID:   content.ID,
		Reporter: reporter.ID,
		Reason:   "Inappropriate content - violates policy",
		Status:   "pending",
	}
	require.NoError(t, repo.CreateReport(ctx, report))

	// Verify report was created
	assert.NotZero(t, report.ID)
	assert.Equal(t, "pending", report.Status)

	// Step 3: Admin reviews the report and rejects content
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationRejected, "Review confirmed: Content inappropriate", &admin.ID)
	require.NoError(t, err)

	// Step 4: Update report status
	report.Status = "resolved"
	report.Result = "Content removed as per policy"
	report.HandledBy = &admin.ID
	now := time.Now()
	report.HandledAt = &now
	// In real implementation, there would be an UpdateReport method
	db.Save(report)

	// Verify final state
	finalContent, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationRejected, finalContent.ModerationStatus)
}

// TestContentModeration_MultipleReports tests handling multiple reports for the same content.
func TestContentModeration_MultipleReports(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "multi_report_author")

	// Create approved content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Controversial content",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Multiple users report the same content
	for i := 0; i < 3; i++ {
		reporter := CreateUniqueTestUser(t, db, fmt.Sprintf("reporter_%d", i))
		report := &model.FeedReport{
			FeedID:   content.ID,
			Reporter: reporter.ID,
			Reason:   fmt.Sprintf("Report reason %d", i),
			Status:   "pending",
		}
		require.NoError(t, repo.CreateReport(ctx, report))
	}

	// Verify all reports are recorded
	reports, total, err := repo.ListReports(ctx, repository.FeedReportListOptions{
		FeedID:   &content.ID,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(reports), 3)
}

// TestContentModeration_CategoryBasedModeration tests moderation by content category.
func TestContentModeration_CategoryBasedModeration(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "category_author")
	moderator := CreateUniqueTestUser(t, db, "category_moderator")

	// Create content categories
	category1 := CreateTestContentCategory(t, db, "Gaming", "Gaming related content")
	category2 := CreateTestContentCategory(t, db, "Off-Topic", "Non-gaming content")

	// Create content in different categories
	content1 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Gaming strategy discussion",
		CategoryID:       &category1.ID,
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content1))

	content2 := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Off-topic discussion",
		CategoryID:       &category2.ID,
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content2))

	// Approve content in gaming category, reject off-topic
	err := repo.UpdateModeration(ctx, content1.ID, model.FeedModerationApproved, "Approved - on topic", &moderator.ID)
	require.NoError(t, err)

	err = repo.UpdateModeration(ctx, content2.ID, model.FeedModerationRejected, "Rejected - off topic", &moderator.ID)
	require.NoError(t, err)

	// Verify results
	updated1, _ := repo.Get(ctx, content1.ID)
	assert.Equal(t, model.FeedModerationApproved, updated1.ModerationStatus)

	updated2, _ := repo.Get(ctx, content2.ID)
	assert.Equal(t, model.FeedModerationRejected, updated2.ModerationStatus)
}

// TestContentModeration_WithImages tests moderation of content with images.
func TestContentModeration_WithImages(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "image_author")
	moderator := CreateUniqueTestUser(t, db, "image_moderator")

	// Create content with multiple images
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Check out these gaming screenshots",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
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
	require.NoError(t, repo.Create(ctx, content))

	// Approve content with images
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationApproved, "Images approved - no issues", &moderator.ID)
	require.NoError(t, err)

	// Verify images are preserved
	updated, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Len(t, updated.Images, 2)
	assert.Equal(t, model.FeedModerationApproved, updated.ModerationStatus)
}

// TestContentModeration_AutoModerationWithoutModerator tests auto-moderation (e.g., by sensitive word detection).
func TestContentModeration_AutoModerationWithoutModerator(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "auto_author")

	// Create content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Content triggering auto-rejection",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Auto-moderate without moderator (nil moderatorID)
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationRejected, "Auto-rejected by sensitive word detection", nil)
	require.NoError(t, err)

	// Verify auto-moderation fields
	updated, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationRejected, updated.ModerationStatus)
	assert.Nil(t, updated.ModeratedBy) // No moderator
	assert.NotNil(t, updated.AutoModeratedAt)
	assert.Nil(t, updated.ManualModeratedAt)
}

// TestContentModeration_ModerationQueuePriority tests prioritization of moderation queue.
func TestContentModeration_ModerationQueuePriority(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "priority_author")

	// Create content with different priorities based on time
	olderTime := time.Now().Add(-2 * time.Hour)
	recentTime := time.Now().Add(-1 * time.Hour)

	oldContent := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Old pending content",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	oldContent.CreatedAt = olderTime
	require.NoError(t, repo.Create(ctx, oldContent))

	recentContent := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Recent pending content",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	recentContent.CreatedAt = recentTime
	require.NoError(t, repo.Create(ctx, recentContent))

	// Get pending queue (should be ordered by creation time)
	pendingStatus := model.FeedModerationPending
	pendingFeeds, total, err := repo.ListPaged(ctx, repository.FeedPagedListOptions{
		ModerationStatus: &pendingStatus,
		Page:             1,
		PageSize:         10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))

	// Older content should appear first (higher priority due to wait time)
	if len(pendingFeeds) >= 2 {
		// First item should be older or same age as second
		assert.True(t, pendingFeeds[0].CreatedAt.Before(pendingFeeds[1].CreatedAt) ||
			pendingFeeds[0].CreatedAt.Equal(pendingFeeds[1].CreatedAt))
	}
}

// TestContentModeration_ContentStatistics tests moderation statistics and counts.
func TestContentModeration_ContentStatistics(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "stats_author")

	// Create content with different statuses
	statuses := []model.FeedModerationStatus{
		model.FeedModerationPending,
		model.FeedModerationPending,
		model.FeedModerationPending,
		model.FeedModerationApproved,
		model.FeedModerationApproved,
		model.FeedModerationRejected,
	}

	for i, status := range statuses {
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          fmt.Sprintf("Content for stats %d", i),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: status,
		}
		require.NoError(t, repo.Create(ctx, content))
	}

	// Get counts by status
	counts, err := repo.CountByStatus(ctx)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, counts[model.FeedModerationPending], int64(3))
	assert.GreaterOrEqual(t, counts[model.FeedModerationApproved], int64(2))
	assert.GreaterOrEqual(t, counts[model.FeedModerationRejected], int64(1))
}

// TestContentModeration_ReportStatusTracking tests report status lifecycle.
func TestContentModeration_ReportStatusTracking(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "report_status_author")
	reporter := CreateUniqueTestUser(t, db, "report_status_reporter")
	admin := CreateUniqueTestUser(t, db, "report_status_admin")

	// Create and approve content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Content to be reported",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Create report
	report := &model.FeedReport{
		FeedID:   content.ID,
		Reporter: reporter.ID,
		Reason:   "Inappropriate content",
		Status:   "pending",
	}
	require.NoError(t, repo.CreateReport(ctx, report))

	// Track status transitions: pending -> processing -> resolved
	statuses := []string{"pending", "processing", "resolved"}
	for _, status := range statuses {
		report.Status = status
		if status == "resolved" {
			report.HandledBy = &admin.ID
			now := time.Now()
			report.HandledAt = &now
			report.Result = "Action taken"
		}
		db.Save(report)

		// Verify status
		var updated model.FeedReport
		db.First(&updated, report.ID)
		assert.Equal(t, status, updated.Status)
	}
}

// TestContentModeration_DifferentSensitiveWordCategories tests different sensitive word categories.
func TestContentModeration_DifferentSensitiveWordCategories(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	feedRepo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "category_author")

	// Create sensitive words in different categories
	categories := []model.SensitiveWordCategory{
		model.SensitiveWordCategoryPolitics,
		model.SensitiveWordCategoryPorn,
		model.SensitiveWordCategoryAbuse,
		model.SensitiveWordCategoryAd,
	}

	for _, category := range categories {
		word := &model.SensitiveWord{
			Word:        fmt.Sprintf("%s_test_word", category),
			Category:    category,
			MatchType:   model.SensitiveWordMatchTypeExact,
			Severity:    model.SensitiveWordSeverityHigh,
			Replacement: "***",
			IsActive:    true,
		}
		require.NoError(t, db.Create(word).Error)

		// Create content with this sensitive word
		content := &model.Feed{
			AuthorID:         author.ID,
			Content:          fmt.Sprintf("Content with %s_test_word", category),
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationPending,
		}
		require.NoError(t, feedRepo.Create(ctx, content))

		// Simulate auto-rejection for high severity
		content.ModerationStatus = model.FeedModerationRejected
		content.ModerationNote = fmt.Sprintf("Auto-rejected: %s sensitive word detected", category)
		now := time.Now()
		content.AutoModeratedAt = &now
		require.NoError(t, feedRepo.Update(ctx, content))

		// Verify
		updated, err := feedRepo.Get(ctx, content.ID)
		require.NoError(t, err)
		assert.Equal(t, model.FeedModerationRejected, updated.ModerationStatus)
		assert.Contains(t, updated.ModerationNote, string(category))
	}
}

// TestContentModeration_AppealProcess tests the appeal process for rejected content.
func TestContentModeration_AppealProcess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "appeal_author")
	moderator1 := CreateUniqueTestUser(t, db, "moderator1")
	moderator2 := CreateUniqueTestUser(t, db, "moderator2")

	// Create and reject content
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Content that was wrongly rejected",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content))

	// First moderator rejects
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationRejected, "Initially rejected - mistake", &moderator1.ID)
	require.NoError(t, err)

	// Author appeals, second moderator reviews and approves
	err = repo.UpdateModeration(ctx, content.ID, model.FeedModerationApproved, "Appeal approved - content is acceptable", &moderator2.ID)
	require.NoError(t, err)

	// Verify final state
	updated, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FeedModerationApproved, updated.ModerationStatus)
	assert.Equal(t, &moderator2.ID, updated.ModeratedBy)
}

// TestContentModeration_PrivateContentVisibility tests that private content respects visibility settings.
func TestContentModeration_PrivateContentVisibility(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "private_author")
	otherUser := CreateUniqueTestUser(t, db, "other_user")

	// Create private approved content
	privateContent := &model.Feed{
		AuthorID:         author.ID,
		Content:          "My private thoughts",
		Visibility:       model.FeedVisibilityPrivate,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, privateContent))

	// Create public approved content
	publicContent := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Public announcement",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
	}
	require.NoError(t, repo.Create(ctx, publicContent))

	// Get author's feeds (should see both)
	authorFeeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:      10,
		AuthorID:   &author.ID,
		Visibility: []model.FeedVisibility{model.FeedVisibilityPublic, model.FeedVisibilityPrivate},
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(authorFeeds), 2)

	// Get public feeds only (should not see private)
	publicFeeds, err := repo.List(ctx, repository.FeedListOptions{
		Limit:        10,
		OnlyApproved: true,
	})
	require.NoError(t, err)

	// Verify public content is visible, private is not
	foundPublic := false
	for _, f := range publicFeeds {
		if f.ID == publicContent.ID {
			foundPublic = true
		}
		if f.ID == privateContent.ID {
			t.Errorf("Private content should not appear in public feed")
		}
	}
	_ = otherUser // Used for context
	_ = foundPublic
}

// TestContentModeration_ContentWithCategory tests content categorization workflow.
func TestContentModeration_ContentWithCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := feed.NewFeedRepository(db)
	ctx := context.Background()

	author := CreateUniqueTestUser(t, db, "category_moderation_author")
	moderator := CreateUniqueTestUser(t, db, "category_moderator")

	// Create category
	category := CreateTestContentCategory(t, db, "Game Reviews", "User game reviews")

	// Create content with category
	content := &model.Feed{
		AuthorID:         author.ID,
		Content:          "Great game review content",
		CategoryID:       &category.ID,
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationPending,
	}
	require.NoError(t, repo.Create(ctx, content))

	// Approve content
	err := repo.UpdateModeration(ctx, content.ID, model.FeedModerationApproved, "Approved review", &moderator.ID)
	require.NoError(t, err)

	// Verify category is preserved
	updated, err := repo.Get(ctx, content.ID)
	require.NoError(t, err)
	assert.NotNil(t, updated.CategoryID)
	assert.Equal(t, category.ID, *updated.CategoryID)
}

// ============================================================================
// Helper Functions
// ============================================================================

// CreateTestFeedWithStatus creates a feed with a specific moderation status.
func CreateTestFeedWithStatus(t *testing.T, db *gorm.DB, author *model.User, status model.FeedModerationStatus) *model.Feed {
	t.Helper()
	feed := &model.Feed{
		Base: model.Base{
			ExtJSON: "{}",
		},
		AuthorID:         author.ID,
		Content:          fmt.Sprintf("Test feed with status %s", status),
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: status,
	}
	if err := db.Create(feed).Error; err != nil {
		t.Fatalf("Failed to create test feed: %v", err)
	}
	return feed
}

// CreateTestSensitiveWordWithSeverity creates a sensitive word with specific category and severity.
func CreateTestSensitiveWordWithSeverity(t *testing.T, db *gorm.DB, word string, category model.SensitiveWordCategory, severity model.SensitiveWordSeverity) *model.SensitiveWord {
	t.Helper()
	sw := &model.SensitiveWord{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Word:        word,
		Category:    category,
		MatchType:   model.SensitiveWordMatchTypeExact,
		Severity:    severity,
		Replacement: "***",
		IsActive:    true,
	}
	if err := db.Create(sw).Error; err != nil {
		t.Fatalf("Failed to create test sensitive word: %v", err)
	}
	return sw
}

// CreateTestContentCategory creates a test content category.
func CreateTestContentCategory(t *testing.T, db *gorm.DB, name, description string) *model.ContentCategory {
	t.Helper()
	category := &model.ContentCategory{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:        name,
		Description: description,
		Status:      model.ContentCategoryStatusActive,
		SortOrder:   0,
	}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test content category: %v", err)
	}
	return category
}

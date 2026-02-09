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

// MockFeedModerationEngine is a mock for FeedModerationEngine
type MockFeedModerationEngine struct {
	mock.Mock
}

func (m *MockFeedModerationEngine) Evaluate(ctx context.Context, input FeedModerationInput) (FeedModerationResult, error) {
	args := m.Called(ctx, input)
	// Handle nil result case for error returns
	if args.Get(0) == nil {
		return FeedModerationResult{}, args.Error(1)
	}
	return args.Get(0).(FeedModerationResult), args.Error(1)
}

func TestNewFeedService(t *testing.T) {
	repo := &MockFeedRepository{}
	moderation := &MockFeedModerationEngine{}

	t.Run("with all dependencies", func(t *testing.T) {
		svc := NewFeedService(repo, moderation)
		require.NotNil(t, svc)
		assert.Equal(t, repo, svc.repo)
		assert.Equal(t, moderation, svc.moderation)
	})

	t.Run("with nil moderation", func(t *testing.T) {
		svc := NewFeedService(repo, nil)
		require.NotNil(t, svc)
		assert.NotNil(t, svc.moderation)
	})
}

func TestFeedService_CreateFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success - text only", func(t *testing.T) {
		repo := &MockFeedRepository{}
		moderation := &MockFeedModerationEngine{}
		svc := NewFeedService(repo, moderation)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		moderation.On("Evaluate", ctx, mock.AnythingOfType("content.FeedModerationInput")).
			Return(FeedModerationResult{
				Decision: FeedModerationDecisionApprove,
				Reason:   "auto-approved",
			}, nil)
		repo.On("UpdateModeration", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := CreateFeedRequest{
			Content:    "Hello world",
			Visibility: model.FeedVisibilityPublic,
			Images:     []FeedImageInput{},
		}

		result, err := svc.CreateFeed(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Hello world", result.Content)
		assert.Equal(t, model.FeedVisibilityPublic, result.Visibility)
		assert.Equal(t, string(model.FeedModerationApproved), result.ModerationStatus)
		repo.AssertExpectations(t)
		moderation.AssertExpectations(t)
	})

	t.Run("success - with images", func(t *testing.T) {
		repo := &MockFeedRepository{}
		moderation := &MockFeedModerationEngine{}
		svc := NewFeedService(repo, moderation)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		moderation.On("Evaluate", ctx, mock.AnythingOfType("content.FeedModerationInput")).
			Return(FeedModerationResult{
				Decision: FeedModerationDecisionManual,
				Reason:   "needs review",
			}, nil)
		repo.On("UpdateModeration", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := CreateFeedRequest{
			Content:    "Check out this photo",
			Visibility: model.FeedVisibilityFollowers,
			Images: []FeedImageInput{
				{URL: "https://example.com/photo.jpg", Width: 800, Height: 600, SizeBytes: 1024 * 1024},
			},
		}

		result, err := svc.CreateFeed(ctx, 100, req)

		require.NoError(t, err)
		assert.Equal(t, "Check out this photo", result.Content)
	})

	t.Run("invalid visibility", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibility("invalid"),
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "visibility 不支持")
	})

	t.Run("content too long", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		// Create content > 1000 runes
		longContent := string(make([]rune, 1001))
		req := CreateFeedRequest{
			Content:    longContent,
			Visibility: model.FeedVisibilityPublic,
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "动态内容验证失败")
	})

	t.Run("too many images", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		images := make([]FeedImageInput, 10)
		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibilityPublic,
			Images:     images,
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "图片数量超过限制")
	})

	t.Run("image size exceeds limit", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibilityPublic,
			Images: []FeedImageInput{
				{URL: "https://example.com/large.jpg", SizeBytes: 11 * 1024 * 1024},
			},
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "超过10MB")
	})

	t.Run("empty image URL", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibilityPublic,
			Images: []FeedImageInput{
				{URL: "   ", SizeBytes: 1024},
			},
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "URL为空")
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		moderation := &MockFeedModerationEngine{}
		svc := NewFeedService(repo, moderation)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(errors.New("db error"))

		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibilityPublic,
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("moderation error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		moderation := &MockFeedModerationEngine{}
		svc := NewFeedService(repo, moderation)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		moderation.On("Evaluate", ctx, mock.AnythingOfType("content.FeedModerationInput")).
			Return(nil, errors.New("moderation service error"))

		req := CreateFeedRequest{
			Content:    "Test",
			Visibility: model.FeedVisibilityPublic,
		}

		_, err := svc.CreateFeed(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "moderation service error")
	})
}

func TestFeedService_ListFeeds(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success with cursor", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		feeds := []model.Feed{
			{Base: model.Base{ID: 100, CreatedAt: now, UpdatedAt: now}, AuthorID: 1, Content: "Feed 1", ModerationStatus: model.FeedModerationApproved},
			{Base: model.Base{ID: 99, CreatedAt: now, UpdatedAt: now}, AuthorID: 2, Content: "Feed 2", ModerationStatus: model.FeedModerationApproved},
		}
		repo.On("List", ctx, mock.MatchedBy(func(opts repository.FeedListOptions) bool {
			return opts.OnlyApproved == true && opts.Limit == 20
		})).Return(feeds, nil)

		req := UserListFeedsRequest{Cursor: "", Limit: 20}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, "99", result.NextCursor)
		repo.AssertExpectations(t)
	})

	t.Run("success with cursor pagination", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		feeds := []model.Feed{
			{Base: model.Base{ID: 50, CreatedAt: now, UpdatedAt: now}, AuthorID: 1, Content: "Feed 1", ModerationStatus: model.FeedModerationApproved},
		}
		repo.On("List", ctx, mock.MatchedBy(func(opts repository.FeedListOptions) bool {
			return opts.CursorBefore != nil && *opts.CursorBefore == 100
		})).Return(feeds, nil)

		req := UserListFeedsRequest{Cursor: "100", Limit: 10}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		assert.Len(t, result.Items, 1)
		assert.Equal(t, "50", result.NextCursor)
	})

	t.Run("invalid cursor", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		req := UserListFeedsRequest{Cursor: "invalid", Limit: 10}
		_, err := svc.ListFeeds(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cursor 无效")
	})

	t.Run("empty result", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return([]model.Feed{}, nil)

		req := UserListFeedsRequest{Limit: 10}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		assert.Len(t, result.Items, 0)
		assert.Empty(t, result.NextCursor)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return(nil, errors.New("db error"))

		req := UserListFeedsRequest{Limit: 10}
		_, err := svc.ListFeeds(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestFeedService_ReportFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("Get", ctx, uint64(100)).Return(&model.Feed{Base: model.Base{ID: 100}}, nil)
		repo.On("CreateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(nil)

		err := svc.ReportFeed(ctx, 1, 100, "test report content")

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("report reason too long", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		longReason := string(make([]rune, 501))
		err := svc.ReportFeed(ctx, 1, 100, longReason)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "验证失败")
	})

	t.Run("feed not found", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		err := svc.ReportFeed(ctx, 1, 999, "test")

		require.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("create report error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("Get", ctx, uint64(100)).Return(&model.Feed{Base: model.Base{ID: 100}}, nil)
		repo.On("CreateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(errors.New("db error"))

		err := svc.ReportFeed(ctx, 1, 100, "test")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestValidateFeedVisibility(t *testing.T) {
	t.Run("empty visibility", func(t *testing.T) {
		err := validateFeedVisibility("")
		assert.NoError(t, err)
	})

	t.Run("public visibility", func(t *testing.T) {
		err := validateFeedVisibility(model.FeedVisibilityPublic)
		assert.NoError(t, err)
	})

	t.Run("followers visibility", func(t *testing.T) {
		err := validateFeedVisibility(model.FeedVisibilityFollowers)
		assert.NoError(t, err)
	})

	t.Run("private visibility", func(t *testing.T) {
		err := validateFeedVisibility(model.FeedVisibilityPrivate)
		assert.NoError(t, err)
	})

	t.Run("invalid visibility", func(t *testing.T) {
		err := validateFeedVisibility(model.FeedVisibility("unknown"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "visibility 不支持")
	})
}

func TestToFeedView(t *testing.T) {
	now := time.Now()

	t.Run("basic feed", func(t *testing.T) {
		feed := &model.Feed{
			Base:             model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			AuthorID:         100,
			Content:          "Test content",
			Visibility:       model.FeedVisibilityPublic,
			ModerationStatus: model.FeedModerationApproved,
			ModerationNote:   "approved",
			Images:           []model.FeedImage{},
		}

		view := toFeedView(feed)

		assert.Equal(t, uint64(1), view.ID)
		assert.Equal(t, uint64(100), view.AuthorID)
		assert.Equal(t, "Test content", view.Content)
		assert.Equal(t, model.FeedVisibilityPublic, view.Visibility)
		assert.Equal(t, "approved", view.ModerationStatus)
		assert.Empty(t, view.Images)
	})

	t.Run("feed with images", func(t *testing.T) {
		feed := &model.Feed{
			Base: model.Base{ID: 1, CreatedAt: now, UpdatedAt: now},
			Images: []model.FeedImage{
				{Base: model.Base{ID: 1}, URL: "https://example.com/1.jpg", Width: 800, Height: 600, SizeBytes: 1024, Order: 0},
				{Base: model.Base{ID: 2}, URL: "https://example.com/2.jpg", Width: 1024, Height: 768, SizeBytes: 2048, Order: 1},
			},
		}

		view := toFeedView(feed)

		assert.Len(t, view.Images, 2)
		assert.Equal(t, "https://example.com/1.jpg", view.Images[0].URL)
		assert.Equal(t, 800, view.Images[0].Width)
		assert.Equal(t, 0, view.Images[0].Order)
		assert.Equal(t, int64(2048), view.Images[1].SizeBytes)
		assert.Equal(t, 1, view.Images[1].Order)
	})
}

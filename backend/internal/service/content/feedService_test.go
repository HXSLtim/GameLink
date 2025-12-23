package content

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
)

// mockFeedModerationEngine is a mock implementation of FeedModerationEngine
type mockFeedModerationEngine struct {
	mock.Mock
}

func (m *mockFeedModerationEngine) Evaluate(ctx context.Context, input FeedModerationInput) (FeedModerationResult, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(FeedModerationResult), args.Error(1)
}

func TestNewFeedService(t *testing.T) {
	repo := &MockFeedRepository{}

	t.Run("with nil moderation engine uses default", func(t *testing.T) {
		svc := NewFeedService(repo, nil)
		require.NotNil(t, svc)
		require.NotNil(t, svc.moderation)
	})

	t.Run("with custom moderation engine", func(t *testing.T) {
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)
		require.NotNil(t, svc)
		assert.Equal(t, modEngine, svc.moderation)
	})
}

func TestFeedService_CreateFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success with approved content", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		repo.On("UpdateModeration", ctx, uint64(0), model.FeedModerationApproved, "auto approved", (*uint64)(nil)).Return(nil)
		modEngine.On("Evaluate", ctx, mock.AnythingOfType("FeedModerationInput")).Return(FeedModerationResult{
			Decision: FeedModerationDecisionApprove,
			Reason:   "auto approved",
		}, nil)

		req := CreateFeedRequest{
			Content:    "测试动态内容",
			Visibility: model.FeedVisibilityPublic,
		}

		result, err := svc.CreateFeed(ctx, 1, req)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "测试动态内容", result.Content)

		repo.AssertExpectations(t)
		modEngine.AssertExpectations(t)
	})

	t.Run("success with rejected content", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		repo.On("UpdateModeration", ctx, uint64(0), model.FeedModerationRejected, "内容不当", (*uint64)(nil)).Return(nil)
		modEngine.On("Evaluate", ctx, mock.AnythingOfType("FeedModerationInput")).Return(FeedModerationResult{
			Decision: FeedModerationDecisionReject,
			Reason:   "内容不当",
		}, nil)

		req := CreateFeedRequest{
			Content:    "测试内容",
			Visibility: model.FeedVisibilityPublic,
		}

		result, err := svc.CreateFeed(ctx, 1, req)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, string(model.FeedModerationRejected), result.ModerationStatus)

		repo.AssertExpectations(t)
		modEngine.AssertExpectations(t)
	})

	t.Run("success with manual review", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		repo.On("UpdateModeration", ctx, uint64(0), model.FeedModerationPending, "需要人工审核", (*uint64)(nil)).Return(nil)
		modEngine.On("Evaluate", ctx, mock.AnythingOfType("FeedModerationInput")).Return(FeedModerationResult{
			Decision: FeedModerationDecisionManual,
			Reason:   "需要人工审核",
		}, nil)

		req := CreateFeedRequest{
			Content:    "可疑内容",
			Visibility: model.FeedVisibilityPublic,
		}

		result, err := svc.CreateFeed(ctx, 1, req)
		require.NoError(t, err)
		require.NotNil(t, result)

		repo.AssertExpectations(t)
		modEngine.AssertExpectations(t)
	})

	t.Run("invalid visibility", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: "invalid_visibility",
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "可见性")
	})

	t.Run("too many images", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		images := make([]FeedImageInput, 10)
		for i := range images {
			images[i] = FeedImageInput{URL: "https://example.com/img.jpg"}
		}

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: model.FeedVisibilityPublic,
			Images:     images,
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "图片数量")
	})

	t.Run("image too large", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: model.FeedVisibilityPublic,
			Images: []FeedImageInput{
				{URL: "https://example.com/img.jpg", SizeBytes: 20 * 1024 * 1024},
			},
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "10MB")
	})

	t.Run("empty image URL", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: model.FeedVisibilityPublic,
			Images: []FeedImageInput{
				{URL: "", SizeBytes: 1024},
			},
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "URL为空")
	})

	t.Run("repo create error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(errors.New("db error"))

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: model.FeedVisibilityPublic,
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("moderation error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		modEngine := &mockFeedModerationEngine{}
		svc := NewFeedService(repo, modEngine)

		repo.On("Create", ctx, mock.AnythingOfType("*model.Feed")).Return(nil)
		modEngine.On("Evaluate", ctx, mock.AnythingOfType("FeedModerationInput")).Return(FeedModerationResult{}, errors.New("moderation error"))

		req := CreateFeedRequest{
			Content:    "测试",
			Visibility: model.FeedVisibilityPublic,
		}

		_, err := svc.CreateFeed(ctx, 1, req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "moderation error")
	})
}

func TestFeedService_ListFeeds(t *testing.T) {
	ctx := context.Background()

	t.Run("success without cursor", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		feeds := []model.Feed{
			{Base: model.Base{ID: 1}, AuthorID: 1, Content: "Feed 1"},
			{Base: model.Base{ID: 2}, AuthorID: 2, Content: "Feed 2"},
		}
		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return(feeds, nil)

		req := UserListFeedsRequest{Limit: 10}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 2)
		assert.Equal(t, "2", result.NextCursor)
	})

	t.Run("success with cursor", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		feeds := []model.Feed{
			{Base: model.Base{ID: 3}, AuthorID: 1, Content: "Feed 3"},
		}
		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return(feeds, nil)

		req := UserListFeedsRequest{Cursor: "2", Limit: 10}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 1)
	})

	t.Run("invalid cursor", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		req := UserListFeedsRequest{Cursor: "invalid", Limit: 10}
		_, err := svc.ListFeeds(ctx, 1, req)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cursor")
	})

	t.Run("empty result", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return([]model.Feed{}, nil)

		req := UserListFeedsRequest{Limit: 10}
		result, err := svc.ListFeeds(ctx, 1, req)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Items)
		assert.Empty(t, result.NextCursor)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("List", ctx, mock.AnythingOfType("repository.FeedListOptions")).Return(nil, errors.New("db error"))

		req := UserListFeedsRequest{Limit: 10}
		_, err := svc.ListFeeds(ctx, 1, req)

		require.Error(t, err)
	})
}

func TestFeedService_ReportFeed(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		feed := &model.Feed{Base: model.Base{ID: 1}, AuthorID: 1, Content: "Test"}
		repo.On("Get", ctx, uint64(1)).Return(feed, nil)
		repo.On("CreateReport", ctx, mock.AnythingOfType("*model.FeedReport")).Return(nil)

		err := svc.ReportFeed(ctx, 2, 1, "内容不当")
		require.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("feed not found", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		repo.On("Get", ctx, uint64(999)).Return(nil, errors.New("not found"))

		err := svc.ReportFeed(ctx, 2, 999, "内容不当")
		require.Error(t, err)
	})

	t.Run("reason too long", func(t *testing.T) {
		repo := &MockFeedRepository{}
		svc := NewFeedService(repo, nil)

		longReason := ""
		for i := 0; i < 600; i++ {
			longReason += "a"
		}

		err := svc.ReportFeed(ctx, 2, 1, longReason)
		require.Error(t, err)
	})
}

func TestValidateFeedVisibility(t *testing.T) {
	tests := []struct {
		name       string
		visibility model.FeedVisibility
		wantErr    bool
	}{
		{"empty is valid", "", false},
		{"public is valid", model.FeedVisibilityPublic, false},
		{"followers is valid", model.FeedVisibilityFollowers, false},
		{"private is valid", model.FeedVisibilityPrivate, false},
		{"invalid visibility", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeedVisibility(tt.visibility)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToFeedView(t *testing.T) {
	feed := &model.Feed{
		Base:             model.Base{ID: 1},
		AuthorID:         100,
		Content:          "Test content",
		Visibility:       model.FeedVisibilityPublic,
		ModerationStatus: model.FeedModerationApproved,
		ModerationNote:   "auto approved",
		Images: []model.FeedImage{
			{URL: "https://example.com/1.jpg", Width: 800, Height: 600, SizeBytes: 1024, Order: 0},
			{URL: "https://example.com/2.jpg", Width: 1024, Height: 768, SizeBytes: 2048, Order: 1},
		},
	}

	view := toFeedView(feed)

	assert.Equal(t, uint64(1), view.ID)
	assert.Equal(t, uint64(100), view.AuthorID)
	assert.Equal(t, "Test content", view.Content)
	assert.Equal(t, model.FeedVisibilityPublic, view.Visibility)
	assert.Equal(t, string(model.FeedModerationApproved), view.ModerationStatus)
	assert.Equal(t, "auto approved", view.ModerationNote)
	assert.Len(t, view.Images, 2)
	assert.Equal(t, "https://example.com/1.jpg", view.Images[0].URL)
	assert.Equal(t, 800, view.Images[0].Width)
}

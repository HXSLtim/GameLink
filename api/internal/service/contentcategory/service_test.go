package contentcategory

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

// MockContentCategoryRepository is a mock implementation of ContentCategoryRepository
type MockContentCategoryRepository struct {
	mock.Mock
}

func (m *MockContentCategoryRepository) Create(ctx context.Context, category *model.ContentCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockContentCategoryRepository) Get(ctx context.Context, id uint64) (*model.ContentCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContentCategory), args.Error(1)
}

func (m *MockContentCategoryRepository) GetByName(ctx context.Context, name string) (*model.ContentCategory, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ContentCategory), args.Error(1)
}

func (m *MockContentCategoryRepository) List(ctx context.Context, opts repository.ContentCategoryListOptions) ([]model.ContentCategory, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ContentCategory), args.Get(1).(int64), args.Error(2)
}

func (m *MockContentCategoryRepository) Update(ctx context.Context, category *model.ContentCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockContentCategoryRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockContentCategoryRepository) GetFeedCount(ctx context.Context, categoryID uint64) (int64, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockContentCategoryRepository) MigrateFeeds(ctx context.Context, fromCategoryID, toCategoryID uint64) error {
	args := m.Called(ctx, fromCategoryID, toCategoryID)
	return args.Error(0)
}

// Helper function to create a test category
func createTestCategory(id uint64, name string) *model.ContentCategory {
	now := time.Now()
	return &model.ContentCategory{
		Base: model.Base{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		Description: "Test description",
		SortOrder:   1,
		Status:      model.ContentCategoryStatusActive,
		IconURL:     "https://example.com/icon.png",
	}
}

func TestNewContentCategoryService(t *testing.T) {
	t.Run("creates service successfully", func(t *testing.T) {
		repo := &MockContentCategoryRepository{}
		svc := NewContentCategoryService(repo)

		require.NotNil(t, svc)
		assert.Equal(t, repo, svc.repo)
	})
}

// ==================== Create Tests ====================

func TestContentCategoryService_Create_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name:        "游戏攻略",
		Description: "分享游戏攻略和技巧",
		SortOrder:   1,
		IconURL:     "https://example.com/icon.png",
	}

	repo.On("Create", ctx, mock.AnythingOfType("*model.ContentCategory")).Return(nil)

	result, err := svc.Create(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "游戏攻略", result.Name)
	assert.Equal(t, "分享游戏攻略和技巧", result.Description)
	assert.Equal(t, 1, result.SortOrder)
	assert.Equal(t, "https://example.com/icon.png", result.IconURL)
	assert.Equal(t, model.ContentCategoryStatusActive, result.Status)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Create_EmptyName(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name: "",
	}

	_, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrValidation, err)
	repo.AssertNotCalled(t, "Create")
}

func TestContentCategoryService_Create_WhitespaceName(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name: "   ",
	}

	_, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrValidation, err)
}

func TestContentCategoryService_Create_NameTrimmed(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name: "  游戏攻略  ",
	}

	repo.On("Create", ctx, mock.MatchedBy(func(c *model.ContentCategory) bool {
		return c.Name == "游戏攻略"
	})).Return(nil)

	result, err := svc.Create(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "游戏攻略", result.Name)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Create_DuplicateName(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name: "游戏攻略",
	}

	repo.On("Create", ctx, mock.AnythingOfType("*model.ContentCategory")).Return(errors.New("分类名称已存在"))

	_, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrDuplicate, err)
}

func TestContentCategoryService_Create_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := CreateRequest{
		Name: "游戏攻略",
	}

	repo.On("Create", ctx, mock.AnythingOfType("*model.ContentCategory")).Return(errors.New("database error"))

	_, err := svc.Create(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

// ==================== Get Tests ====================

func TestContentCategoryService_Get_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(10), nil)

	result, err := svc.Get(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, "游戏攻略", result.Name)
	assert.Equal(t, int64(10), result.FeedCount)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Get_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	repo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	_, err := svc.Get(ctx, 999)

	require.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
	repo.AssertNotCalled(t, "GetFeedCount")
}

func TestContentCategoryService_Get_FeedCountError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(0), errors.New("feed count error"))

	result, err := svc.Get(ctx, 1)

	// Feed count error should not prevent returning the category
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint64(1), result.ID)
	assert.Equal(t, int64(0), result.FeedCount) // Should default to 0 on error
}

// ==================== List Tests ====================

func TestContentCategoryService_List_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	categories := []model.ContentCategory{
		*createTestCategory(1, "游戏攻略"),
		*createTestCategory(2, "游戏分享"),
	}

	req := ListRequest{
		Page:     1,
		PageSize: 10,
	}

	repo.On("List", ctx, mock.MatchedBy(func(opts repository.ContentCategoryListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 10
	})).Return(categories, int64(2), nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(5), nil)
	repo.On("GetFeedCount", ctx, uint64(2)).Return(int64(8), nil)

	result, err := svc.List(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, "游戏攻略", result.Items[0].Name)
	assert.Equal(t, int64(5), result.Items[0].FeedCount)
	assert.Equal(t, "游戏分享", result.Items[1].Name)
	assert.Equal(t, int64(8), result.Items[1].FeedCount)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_List_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := ListRequest{
		Page:     0,
		PageSize: 0,
	}

	repo.On("List", ctx, mock.MatchedBy(func(opts repository.ContentCategoryListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return([]model.ContentCategory{}, int64(0), nil)

	_, err := svc.List(ctx, req)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_List_PageSizeCapped(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := ListRequest{
		Page:     1,
		PageSize: 200, // When PageSize > 100, it defaults to 20
	}

	repo.On("List", ctx, mock.MatchedBy(func(opts repository.ContentCategoryListOptions) bool {
		return opts.PageSize == 20
	})).Return([]model.ContentCategory{}, int64(0), nil)
	repo.On("GetFeedCount", ctx, mock.Anything).Return(int64(0), nil).Maybe()

	_, err := svc.List(ctx, req)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_List_WithFilters(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	status := model.ContentCategoryStatusActive
	req := ListRequest{
		Page:     1,
		PageSize: 10,
		Keyword:  "游戏",
		Status:   &status,
	}

	repo.On("List", ctx, mock.MatchedBy(func(opts repository.ContentCategoryListOptions) bool {
		return opts.Keyword == "游戏" && opts.Status != nil && *opts.Status == model.ContentCategoryStatusActive
	})).Return([]model.ContentCategory{}, int64(0), nil)
	repo.On("GetFeedCount", ctx, mock.Anything).Return(int64(0), nil).Maybe()

	_, err := svc.List(ctx, req)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_List_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := ListRequest{
		Page:     1,
		PageSize: 10,
	}

	repo.On("List", ctx, mock.AnythingOfType("repository.ContentCategoryListOptions")).Return([]model.ContentCategory{}, int64(0), errors.New("database error"))

	_, err := svc.List(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestContentCategoryService_List_EmptyResult(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := ListRequest{
		Page:     1,
		PageSize: 10,
	}

	repo.On("List", ctx, mock.AnythingOfType("repository.ContentCategoryListOptions")).Return([]model.ContentCategory{}, int64(0), nil)

	result, err := svc.List(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Items, 0)
	assert.Equal(t, int64(0), result.Total)
}

// ==================== Update Tests ====================

func TestContentCategoryService_Update_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "旧名称")

	req := UpdateRequest{
		Name:        "新名称",
		Description: "新描述",
		SortOrder:   5,
		Status:      model.ContentCategoryStatusInactive,
		IconURL:     "https://example.com/new-icon.png",
	}

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(c *model.ContentCategory) bool {
		return c.Name == "新名称" &&
			c.Description == "新描述" &&
			c.SortOrder == 5 &&
			c.Status == model.ContentCategoryStatusInactive &&
			c.IconURL == "https://example.com/new-icon.png"
	})).Return(nil)

	err := svc.Update(ctx, 1, req)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Update_EmptyName(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := UpdateRequest{
		Name: "",
	}

	err := svc.Update(ctx, 1, req)

	require.Error(t, err)
	assert.Equal(t, ErrValidation, err)
	repo.AssertNotCalled(t, "Get")
	repo.AssertNotCalled(t, "Update")
}

func TestContentCategoryService_Update_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	req := UpdateRequest{
		Name: "新名称",
	}

	repo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.Update(ctx, 999, req)

	require.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
	repo.AssertNotCalled(t, "Update")
}

func TestContentCategoryService_Update_DuplicateName(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	req := UpdateRequest{
		Name: "已存在的名称",
	}

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*model.ContentCategory")).Return(errors.New("分类名称已存在"))

	err := svc.Update(ctx, 1, req)

	require.Error(t, err)
	assert.Equal(t, ErrDuplicate, err)
}

func TestContentCategoryService_Update_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	req := UpdateRequest{
		Name:   "新名称",
		Status: model.ContentCategoryStatus("invalid"),
	}

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("Update", ctx, mock.MatchedBy(func(c *model.ContentCategory) bool {
		// Status should remain unchanged when invalid
		return c.Status == model.ContentCategoryStatusActive
	})).Return(nil)

	err := svc.Update(ctx, 1, req)

	require.NoError(t, err)
}

func TestContentCategoryService_Update_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	req := UpdateRequest{
		Name: "新名称",
	}

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*model.ContentCategory")).Return(errors.New("database error"))

	err := svc.Update(ctx, 1, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

// ==================== Delete Tests ====================

func TestContentCategoryService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(0), nil)
	repo.On("Delete", ctx, uint64(1)).Return(nil)

	err := svc.Delete(ctx, 1, nil)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	repo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.Delete(ctx, 999, nil)

	require.Error(t, err)
	assert.Equal(t, repository.ErrNotFound, err)
	repo.AssertNotCalled(t, "GetFeedCount")
	repo.AssertNotCalled(t, "Delete")
}

func TestContentCategoryService_Delete_HasFeedsWithoutMigration(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(10), nil)

	err := svc.Delete(ctx, 1, nil)

	require.Error(t, err)
	assert.Equal(t, ErrHasFeeds, err)
	repo.AssertNotCalled(t, "Delete")
}

func TestContentCategoryService_Delete_WithMigration(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")
	targetCategoryID := uint64(2)

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(10), nil)
	repo.On("MigrateFeeds", ctx, uint64(1), targetCategoryID).Return(nil)
	repo.On("Delete", ctx, uint64(1)).Return(nil)

	err := svc.Delete(ctx, 1, &targetCategoryID)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestContentCategoryService_Delete_MigrationError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")
	targetCategoryID := uint64(2)

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(10), nil)
	repo.On("MigrateFeeds", ctx, uint64(1), targetCategoryID).Return(errors.New("migration failed"))

	err := svc.Delete(ctx, 1, &targetCategoryID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "migration failed")
	repo.AssertNotCalled(t, "Delete")
}

func TestContentCategoryService_Delete_RepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := &MockContentCategoryRepository{}
	svc := NewContentCategoryService(repo)

	category := createTestCategory(1, "游戏攻略")

	repo.On("Get", ctx, uint64(1)).Return(category, nil)
	repo.On("GetFeedCount", ctx, uint64(1)).Return(int64(0), nil)
	repo.On("Delete", ctx, uint64(1)).Return(errors.New("database error"))

	err := svc.Delete(ctx, 1, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

// ==================== toDTO Tests ====================

func TestContentCategoryService_toDTO(t *testing.T) {
	svc := &ContentCategoryService{}
	category := createTestCategory(1, "游戏攻略")

	t.Run("full DTO with feed count", func(t *testing.T) {
		dto := svc.toDTO(category, 15)

		require.NotNil(t, dto)
		assert.Equal(t, uint64(1), dto.ID)
		assert.Equal(t, "游戏攻略", dto.Name)
		assert.Equal(t, "Test description", dto.Description)
		assert.Equal(t, 1, dto.SortOrder)
		assert.Equal(t, model.ContentCategoryStatusActive, dto.Status)
		assert.Equal(t, "https://example.com/icon.png", dto.IconURL)
		assert.Equal(t, int64(15), dto.FeedCount)
		assert.NotEmpty(t, dto.CreatedAt)
		assert.NotEmpty(t, dto.UpdatedAt)
	})

	t.Run("DTO without feed count", func(t *testing.T) {
		dto := svc.toDTO(category, 0)

		require.NotNil(t, dto)
		assert.Equal(t, int64(0), dto.FeedCount)
	})
}

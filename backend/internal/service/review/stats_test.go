package review

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockReviewRepository is a mock implementation of ReviewRepository
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Review), args.Get(1).(int64), args.Error(2)
}

func (m *MockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Review), args.Error(1)
}

func (m *MockReviewRepository) Create(ctx context.Context, review *model.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockReviewRepository) Update(ctx context.Context, review *model.Review) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockReviewRepository) UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error {
	args := m.Called(ctx, id, status, rejectionReason)
	return args.Error(0)
}

func (m *MockReviewRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error {
	args := m.Called(ctx, ids, status, rejectionReason)
	return args.Error(0)
}

func (m *MockReviewRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReviewRepository) GetStats(ctx context.Context) (repository.ReviewStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(repository.ReviewStats), args.Error(1)
}

func (m *MockReviewRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	args := m.Called(ctx, days)
	return args.Get(0).([]repository.DateValue), args.Error(1)
}

func (m *MockReviewRepository) GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]repository.PlayerReviewStats), args.Error(1)
}

func (m *MockReviewRepository) GetTopPlayersByRating(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]repository.PlayerReviewStats), args.Error(1)
}

func (m *MockReviewRepository) GetGameStats(ctx context.Context) ([]repository.GameReviewStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.GameReviewStats), args.Error(1)
}

func TestGetReviewStats(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedStats := repository.ReviewStats{
		TotalReviews:  100,
		AverageRating: 4.5,
		RatingDistribution: map[int]int64{
			1: 5,
			2: 10,
			3: 15,
			4: 30,
			5: 40,
		},
	}

	mockRepo.On("GetStats", mock.Anything).Return(expectedStats, nil)

	result, err := service.GetReviewStats(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.TotalReviews)
	assert.Equal(t, 4.5, result.AverageRating)
	assert.Equal(t, int64(40), result.RatingDistribution[5])

	mockRepo.AssertExpectations(t)
}

func TestGetReviewTrend(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedTrend := []repository.DateValue{
		{Date: "2024-01-01", Value: 10},
		{Date: "2024-01-02", Value: 15},
		{Date: "2024-01-03", Value: 20},
	}

	mockRepo.On("GetTrend", mock.Anything, 30).Return(expectedTrend, nil)

	result, err := service.GetReviewTrend(context.Background(), 30)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Trend, 3)
	assert.Equal(t, "2024-01-01", result.Trend[0].Date)
	assert.Equal(t, int64(10), result.Trend[0].Value)

	mockRepo.AssertExpectations(t)
}

func TestGetTopPlayers_ByReviewCount(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedPlayers := []repository.PlayerReviewStats{
		{PlayerID: 1, PlayerName: "Player1", ReviewCount: 100, AverageRating: 4.5},
		{PlayerID: 2, PlayerName: "Player2", ReviewCount: 90, AverageRating: 4.3},
	}

	mockRepo.On("GetTopPlayersByReviewCount", mock.Anything, 10).Return(expectedPlayers, nil)

	result, err := service.GetTopPlayers(context.Background(), 10, "count")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Players, 2)
	assert.Equal(t, "Player1", result.Players[0].PlayerName)
	assert.Equal(t, int64(100), result.Players[0].ReviewCount)

	mockRepo.AssertExpectations(t)
}

func TestGetTopPlayers_ByRating(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedPlayers := []repository.PlayerReviewStats{
		{PlayerID: 1, PlayerName: "Player1", ReviewCount: 50, AverageRating: 4.9},
		{PlayerID: 2, PlayerName: "Player2", ReviewCount: 60, AverageRating: 4.8},
	}

	mockRepo.On("GetTopPlayersByRating", mock.Anything, 10).Return(expectedPlayers, nil)

	result, err := service.GetTopPlayers(context.Background(), 10, "rating")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Players, 2)
	assert.Equal(t, "Player1", result.Players[0].PlayerName)
	assert.Equal(t, 4.9, result.Players[0].AverageRating)

	mockRepo.AssertExpectations(t)
}

func TestGetGameStats(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedGames := []repository.GameReviewStats{
		{GameID: 1, GameName: "Game1", ReviewCount: 200, AverageRating: 4.5},
		{GameID: 2, GameName: "Game2", ReviewCount: 150, AverageRating: 4.3},
	}

	mockRepo.On("GetGameStats", mock.Anything).Return(expectedGames, nil)

	result, err := service.GetGameStats(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Games, 2)
	assert.Equal(t, "Game1", result.Games[0].GameName)
	assert.Equal(t, int64(200), result.Games[0].ReviewCount)

	mockRepo.AssertExpectations(t)
}

func TestGetReviewTrend_DefaultDays(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedTrend := []repository.DateValue{
		{Date: "2024-01-01", Value: 10},
	}

	// When days is 0 or negative, it should default to 30
	mockRepo.On("GetTrend", mock.Anything, 30).Return(expectedTrend, nil)

	result, err := service.GetReviewTrend(context.Background(), 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestGetTopPlayers_DefaultLimit(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewStatsService(mockRepo)

	expectedPlayers := []repository.PlayerReviewStats{
		{PlayerID: 1, PlayerName: "Player1", ReviewCount: 100, AverageRating: 4.5},
	}

	// When limit is 0 or negative, it should default to 10
	mockRepo.On("GetTopPlayersByReviewCount", mock.Anything, 10).Return(expectedPlayers, nil)

	result, err := service.GetTopPlayers(context.Background(), 0, "count")

	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

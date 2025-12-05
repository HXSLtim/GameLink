package review

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
)

// mockSettingsRepo is a mock implementation of ReviewDisplaySettingsRepository
type mockSettingsRepo struct {
	settings *model.ReviewDisplaySettings
}

func (m *mockSettingsRepo) Get(ctx context.Context) (*model.ReviewDisplaySettings, error) {
	if m.settings == nil {
		return model.DefaultReviewDisplaySettings(), nil
	}
	return m.settings, nil
}

func (m *mockSettingsRepo) Save(ctx context.Context, settings *model.ReviewDisplaySettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	m.settings = settings
	return nil
}

func TestSettingsService_GetSettings(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	settings, err := svc.GetSettings(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, model.ReviewSortByTime, settings.SortBy)
}

func TestSettingsService_UpdateSettings(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	newSettings := &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByScore,
		MinScore:      3,
		ShowAnonymous: false,
		PageSize:      15,
	}

	result, err := svc.UpdateSettings(ctx, newSettings)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewSortByScore, result.SortBy)
	assert.Equal(t, 3, result.MinScore)
	assert.False(t, result.ShowAnonymous)
	assert.Equal(t, 15, result.PageSize)
}

func TestSettingsService_UpdateSettingsPartial(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	// Only update sortBy
	sortBy := model.ReviewSortByLikes
	input := UpdateSettingsInput{
		SortBy: &sortBy,
	}

	result, err := svc.UpdateSettingsPartial(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewSortByLikes, result.SortBy)
	// Other fields should remain default
	assert.Equal(t, 1, result.MinScore)
	assert.True(t, result.ShowAnonymous)
	assert.Equal(t, 20, result.PageSize)
}

func TestSettingsService_UpdateSettingsPartial_AllFields(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	sortBy := model.ReviewSortByScore
	minScore := 4
	showAnonymous := false
	pageSize := 50

	input := UpdateSettingsInput{
		SortBy:        &sortBy,
		MinScore:      &minScore,
		ShowAnonymous: &showAnonymous,
		PageSize:      &pageSize,
	}

	result, err := svc.UpdateSettingsPartial(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, model.ReviewSortByScore, result.SortBy)
	assert.Equal(t, 4, result.MinScore)
	assert.False(t, result.ShowAnonymous)
	assert.Equal(t, 50, result.PageSize)
}

func TestSettingsService_UpdateSettings_ValidationError(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	// Invalid sortBy
	invalidSettings := &model.ReviewDisplaySettings{
		SortBy:        "invalid",
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      20,
	}

	_, err := svc.UpdateSettings(ctx, invalidSettings)
	assert.Error(t, err)
}

func TestSettingsService_UpdateSettingsPartial_ValidationError(t *testing.T) {
	repo := &mockSettingsRepo{}
	svc := NewSettingsService(repo)
	ctx := context.Background()

	// Invalid minScore
	minScore := 10
	input := UpdateSettingsInput{
		MinScore: &minScore,
	}

	_, err := svc.UpdateSettingsPartial(ctx, input)
	assert.Error(t, err)
}

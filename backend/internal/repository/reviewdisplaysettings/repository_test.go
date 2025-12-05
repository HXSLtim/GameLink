package reviewdisplaysettings

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.ReviewDisplaySettings{})
	require.NoError(t, err)

	return db
}

func TestRepository_Get_ReturnsDefaultWhenEmpty(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	settings, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, uint64(1), settings.ID)
	assert.Equal(t, model.ReviewSortByTime, settings.SortBy)
	assert.Equal(t, 1, settings.MinScore)
	assert.True(t, settings.ShowAnonymous)
	assert.Equal(t, 20, settings.PageSize)
}

func TestRepository_Save_CreatesNewSettings(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	settings := &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByScore,
		MinScore:      3,
		ShowAnonymous: true,
		PageSize:      10,
	}

	err := repo.Save(ctx, settings)
	require.NoError(t, err)

	// Verify the settings were saved
	saved, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), saved.ID)
	assert.Equal(t, model.ReviewSortByScore, saved.SortBy)
	assert.Equal(t, 3, saved.MinScore)
	assert.True(t, saved.ShowAnonymous)
	assert.Equal(t, 10, saved.PageSize)
}

func TestRepository_Save_UpdatesExistingSettings(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	// Create initial settings
	initial := &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByTime,
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      20,
	}
	err := repo.Save(ctx, initial)
	require.NoError(t, err)

	// Update settings
	updated := &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByLikes,
		MinScore:      4,
		ShowAnonymous: false,
		PageSize:      50,
	}
	err = repo.Save(ctx, updated)
	require.NoError(t, err)

	// Verify only one record exists
	var count int64
	db.Model(&model.ReviewDisplaySettings{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// Verify the settings were updated
	saved, err := repo.Get(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), saved.ID)
	assert.Equal(t, model.ReviewSortByLikes, saved.SortBy)
	assert.Equal(t, 4, saved.MinScore)
	assert.False(t, saved.ShowAnonymous)
	assert.Equal(t, 50, saved.PageSize)
}

func TestRepository_Save_ValidationError(t *testing.T) {
	db := setupTestDB(t)
	repo := New(db)
	ctx := context.Background()

	// Invalid sortBy
	settings := &model.ReviewDisplaySettings{
		SortBy:        "invalid",
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      20,
	}
	err := repo.Save(ctx, settings)
	assert.Error(t, err)

	// Invalid minScore (too low)
	settings = &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByTime,
		MinScore:      0,
		ShowAnonymous: true,
		PageSize:      20,
	}
	err = repo.Save(ctx, settings)
	assert.Error(t, err)

	// Invalid minScore (too high)
	settings = &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByTime,
		MinScore:      6,
		ShowAnonymous: true,
		PageSize:      20,
	}
	err = repo.Save(ctx, settings)
	assert.Error(t, err)

	// Invalid pageSize (too low)
	settings = &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByTime,
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      0,
	}
	err = repo.Save(ctx, settings)
	assert.Error(t, err)

	// Invalid pageSize (too high)
	settings = &model.ReviewDisplaySettings{
		SortBy:        model.ReviewSortByTime,
		MinScore:      1,
		ShowAnonymous: true,
		PageSize:      101,
	}
	err = repo.Save(ctx, settings)
	assert.Error(t, err)
}

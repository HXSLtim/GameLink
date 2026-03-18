package game

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/pkg/cache"
)

func TestGameRepository_CacheConfiguration(t *testing.T) {
	// Verify cache constants
	assert.Equal(t, "games:list:all", cacheKeyGamesList)
	assert.Equal(t, 1*time.Hour, cacheTTLGames)
}

func TestGameRepository_CacheHelpers(t *testing.T) {
	cacheClient := cache.NewMemory()
	repo := &gormGameRepository{
		db:    nil, // Not needed for cache helper tests
		cache: cacheClient,
	}

	ctx := context.Background()

	// Test cacheGamesList
	games := []model.Game{
		{
			Key:         "test-game-1",
			Name:        "Test Game 1",
			Category:    "MOBA",
			IconURL:     "https://example.com/icon1.png",
			Description: "Test game 1",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Key:         "test-game-2",
			Name:        "Test Game 2",
			Category:    "FPS",
			IconURL:     "https://example.com/icon2.png",
			Description: "Test game 2",
			IsActive:    true,
			SortOrder:   2,
		},
	}

	// Cache the games list
	repo.cacheGamesList(ctx, games)

	// Verify cache exists
	cached, ok, err := cacheClient.Get(ctx, cacheKeyGamesList)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.NotEmpty(t, cached)

	// Test invalidateCache
	repo.invalidateCache(ctx)

	// Verify cache is cleared
	_, ok, err = cacheClient.Get(ctx, cacheKeyGamesList)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestGameRepository_CacheHelpersWithNilCache(t *testing.T) {
	repo := &gormGameRepository{
		db:    nil,
		cache: nil, // No cache
	}

	ctx := context.Background()

	games := []model.Game{
		{Key: "test-game", Name: "Test Game"},
	}

	// Should not panic with nil cache
	repo.cacheGamesList(ctx, games)
	repo.invalidateCache(ctx)
}

func TestGameRepository_ConstructorWithCache(t *testing.T) {
	cacheClient := cache.NewMemory()
	repo := NewGameRepositoryWithCache(nil, cacheClient)

	// Verify cache is set
	gormRepo, ok := repo.(*gormGameRepository)
	assert.True(t, ok)
	assert.NotNil(t, gormRepo.cache)
}

func TestGameRepository_ConstructorWithoutCache(t *testing.T) {
	repo := NewGameRepository(nil)

	// Verify cache is nil
	gormRepo, ok := repo.(*gormGameRepository)
	assert.True(t, ok)
	assert.Nil(t, gormRepo.cache)
}

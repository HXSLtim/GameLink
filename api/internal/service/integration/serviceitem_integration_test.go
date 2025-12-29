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
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/serviceitem"
)

// ============================================================================
// ServiceItem CRUD Tests
// ============================================================================

func TestServiceItemRepository_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("ITEM_%d", time.Now().UnixNano()%1000000),
		Name:           "Test Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 5000,
		ServiceHours:   1,
		CommissionRate: 0.20,
		IsActive:       true,
	}

	err := repo.Create(ctx, item)
	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestServiceItemRepository_Get(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("GET_ITEM_%d", time.Now().UnixNano()%1000000),
		Name:           "Get Test Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 3000,
		IsActive:       true,
	}
	require.NoError(t, repo.Create(ctx, item))

	got, err := repo.Get(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Test Service", got.Name)
	assert.Equal(t, int64(3000), got.BasePriceCents)
}

func TestServiceItemRepository_GetByCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	code := fmt.Sprintf("CODE_%d", time.Now().UnixNano()%1000000)
	item := &model.ServiceItem{
		ItemCode:       code,
		Name:           "Code Test Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 4000,
		IsActive:       true,
	}
	require.NoError(t, repo.Create(ctx, item))

	got, err := repo.GetByCode(ctx, code)
	require.NoError(t, err)
	assert.Equal(t, item.ID, got.ID)
}

func TestServiceItemRepository_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("UPDATE_ITEM_%d", time.Now().UnixNano()%1000000),
		Name:           "Update Test",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 2000,
		IsActive:       true,
	}
	require.NoError(t, repo.Create(ctx, item))

	// Update
	item.Name = "Updated Service"
	item.BasePriceCents = 6000
	err := repo.Update(ctx, item)
	require.NoError(t, err)

	// Verify
	got, err := repo.Get(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Service", got.Name)
	assert.Equal(t, int64(6000), got.BasePriceCents)
}

func TestServiceItemRepository_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("DELETE_ITEM_%d", time.Now().UnixNano()%1000000),
		Name:           "Delete Test",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 1000,
		IsActive:       true,
	}
	require.NoError(t, repo.Create(ctx, item))

	err := repo.Delete(ctx, item.ID)
	require.NoError(t, err)

	_, err = repo.Get(ctx, item.ID)
	assert.Error(t, err)
}

func TestServiceItemRepository_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	// Create items
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("LIST_ITEM_%d_%d", i, time.Now().UnixNano()%1000000),
			Name:           fmt.Sprintf("List Service %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: int64(1000 * (i + 1)),
			IsActive:       true,
		}
		require.NoError(t, repo.Create(ctx, item))
	}

	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(items), 3)
}

func TestServiceItemRepository_ListByCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	// Create items with specific category
	for i := 0; i < 2; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("CAT_ITEM_%d_%d", i, time.Now().UnixNano()%1000000),
			Name:           fmt.Sprintf("Category Service %d", i),
			Category:       "special_cat",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 2000,
			IsActive:       true,
		}
		require.NoError(t, repo.Create(ctx, item))
	}

	category := "special_cat"
	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{
		Category: &category,
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, item := range items {
		assert.Equal(t, "special_cat", item.Category)
	}
}

func TestServiceItemRepository_BatchUpdateStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	var ids []uint64
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BS%d%d", i, time.Now().UnixNano()%1000000),
			Name:           fmt.Sprintf("Batch Status %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 1000,
			IsActive:       true,
		}
		require.NoError(t, repo.Create(ctx, item))
		ids = append(ids, item.ID)
	}

	// Batch update to inactive
	err := repo.BatchUpdateStatus(ctx, ids, false)
	require.NoError(t, err)

	// Verify
	for _, id := range ids {
		got, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.False(t, got.IsActive)
	}
}

func TestServiceItemRepository_GetGifts(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()

	// Create gift items
	for i := 0; i < 2; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("GIFT_%d_%d", i, time.Now().UnixNano()%1000000),
			Name:           fmt.Sprintf("Gift %d", i),
			Category:       "gift",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 500,
			IsActive:       true,
		}
		require.NoError(t, repo.Create(ctx, item))
	}

	gifts, total, err := repo.GetGifts(ctx, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(2))
	for _, gift := range gifts {
		assert.Equal(t, model.SubCategoryGift, gift.SubCategory)
	}
}

func TestServiceItemRepository_GetGameServices(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	itemRepo := serviceitem.NewServiceItemRepository(db)
	gameRepo := game.NewGameRepository(db)
	ctx := context.Background()

	// Create a game first
	gameObj := &model.Game{
		Key:      fmt.Sprintf("game_svc_%d", time.Now().UnixNano()%1000000),
		Name:     "Game for Services",
		Category: "moba",
	}
	require.NoError(t, gameRepo.Create(ctx, gameObj))

	// Create service items for this game
	for i := 0; i < 2; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("GAME_SVC_%d_%d", i, time.Now().UnixNano()%1000000),
			Name:           fmt.Sprintf("Game Service %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			GameID:         &gameObj.ID,
			BasePriceCents: 3000,
			IsActive:       true,
		}
		require.NoError(t, itemRepo.Create(ctx, item))
	}

	services, err := itemRepo.GetGameServices(ctx, gameObj.ID, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(services), 2)
}

// Package integration provides integration tests for ServiceItem module.
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
	"gamelink/internal/repository/serviceitem"
	"gorm.io/gorm"
)

// ============================================================================
// Helper Functions for ServiceItem Tests
// ============================================================================

// CreateTestServiceItemWithDetails creates a test service item with full details.
func CreateTestServiceItemWithDetails(t *testing.T, db *gorm.DB, game *model.Game, category *model.GameCategory, name string, itemCode string, priceCents int64, subCategory model.ServiceItemSubCategory) *model.ServiceItem {
	t.Helper()
	var gameID, categoryID *uint64
	if game != nil {
		gameID = &game.ID
	}
	if category != nil {
		categoryID = &category.ID
	}
	if itemCode == "" {
		itemCode = fmt.Sprintf("IT%d", time.Now().UnixNano()%1000000000)
	}
	if name == "" {
		name = "Test Service Item"
	}
	item := &model.ServiceItem{
		ItemCode:        itemCode,
		Name:            name,
		Description:     "Test service item description for " + name,
		Category:        "escort",
		SubCategory:     subCategory,
		GameID:          gameID,
		CategoryID:      categoryID,
		BasePriceCents:  priceCents,
		ServiceHours:    1,
		CommissionRate:  0.20,
		MinUsers:        1,
		MaxPlayers:      1,
		RequiredPlayers: 1,
		MaxPerOrder:     0,
		UsageLimitType:  model.UsageLimitNone,
		UsageLimitCount: 0,
		IsActive:        true,
		SortOrder:       0,
		Tags:            "[]",
		IconURL:         "https://example.com/icon.png",
		RankLevel:       "",
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test service item: %v", err)
	}
	return item
}

// CreateTestServiceItemWithVIP creates a service item with VIP pricing.
func CreateTestServiceItemWithVIP(t *testing.T, db *gorm.DB, game *model.Game, name string, basePrice, vipPrice int64) *model.ServiceItem {
	t.Helper()
	var gameID *uint64
	if game != nil {
		gameID = &game.ID
	}
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("VIP%d", time.Now().UnixNano()%1000000000),
		Name:           name,
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         gameID,
		BasePriceCents: basePrice,
		VipPriceCents:  &vipPrice,
		CommissionRate: 0.20,
		IsActive:       true,
		Tags:           "[]",
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test service item with VIP: %v", err)
	}
	return item
}

// CreateTestServiceItemWithUsageLimit creates a service item with usage limits.
func CreateTestServiceItemWithUsageLimit(t *testing.T, db *gorm.DB, game *model.Game, name string, limitType model.UsageLimitType, limitCount int) *model.ServiceItem {
	t.Helper()
	var gameID *uint64
	if game != nil {
		gameID = &game.ID
	}
	item := &model.ServiceItem{
		ItemCode:        fmt.Sprintf("UL%d", time.Now().UnixNano()%1000000000),
		Name:            name,
		Category:        "escort",
		SubCategory:     model.SubCategorySolo,
		GameID:          gameID,
		BasePriceCents:  10000,
		CommissionRate:  0.20,
		UsageLimitType:  limitType,
		UsageLimitCount: limitCount,
		MaxPerOrder:     1,
		IsActive:        true,
		Tags:            "[]",
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test service item with usage limit: %v", err)
	}
	return item
}

// CreateTestServiceItemWithStatus creates a service item with specific active status.
func CreateTestServiceItemWithStatus(t *testing.T, db *gorm.DB, game *model.Game, name string, isActive bool) *model.ServiceItem {
	t.Helper()
	var gameID *uint64
	if game != nil {
		gameID = &game.ID
	}
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("ST%d", time.Now().UnixNano()%1000000000),
		Name:           name,
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         gameID,
		BasePriceCents: 10000,
		CommissionRate: 0.20,
		IsActive:       isActive,
		Tags:           "[]",
	}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("Failed to create test service item with status: %v", err)
	}
	return item
}

func TestServiceItemRepository_Create_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "TestGame")
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("TEST%d", time.Now().UnixNano()),
		Name:           "Test Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		GameID:         &game.ID,
		BasePriceCents: 10000,
		CommissionRate: 0.20,
		IsActive:       true,
		Tags:           "[]",
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	assert.NotZero(t, item.ID)
}

func TestServiceItemRepository_Get_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "GetGame")
	item := CreateTestServiceItem(t, db, game, "Get Test Service", 10000)
	got, err := repo.Get(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, "Get Test Service", got.Name)
	assert.Equal(t, int64(10000), got.BasePriceCents)
}

func TestServiceItemRepository_GetByCode_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	itemCode := fmt.Sprintf("CODE%d", time.Now().UnixNano())
	item := &model.ServiceItem{
		ItemCode:       itemCode,
		Name:           "Code Test Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 15000,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	got, err := repo.GetByCode(ctx, itemCode)
	require.NoError(t, err)
	assert.Equal(t, "Code Test Service", got.Name)
	assert.Equal(t, itemCode, got.ItemCode)
}

func TestServiceItemRepository_Update_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "UpdateGame")
	item := CreateTestServiceItem(t, db, game, "Update Test", 10000)
	item.Name = "Updated Service Name"
	item.BasePriceCents = 20000
	item.Description = "Updated description"
	err := repo.Update(ctx, item)
	require.NoError(t, err)
	got, err := repo.Get(ctx, item.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Service Name", got.Name)
	assert.Equal(t, int64(20000), got.BasePriceCents)
	assert.Equal(t, "Updated description", got.Description)
}

func TestServiceItemRepository_Delete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "DeleteGame")
	item := CreateTestServiceItem(t, db, game, "Delete Test", 10000)
	err := repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.Get(ctx, item.ID)
	assert.Error(t, err)
}

func TestServiceItemRepository_List_All(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "ListGame")
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("LIST%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("List Item %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			GameID:         &game.ID,
			BasePriceCents: 10000 + int64(i*1000),
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
	}
	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(items), 3)
}

func TestServiceItemRepository_BatchDelete_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	var ids []uint64
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BDEL%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Batch Delete %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
		ids = append(ids, item.ID)
	}
	affected, err := repo.BatchDelete(ctx, ids)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)
	for _, id := range ids {
		_, err := repo.Get(ctx, id)
		assert.Error(t, err)
	}
}

func TestServiceItemRepository_BatchUpdateStatus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	var ids []uint64
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BSTAT%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Batch Status %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			IsActive:       true,
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
		ids = append(ids, item.ID)
	}
	err := repo.BatchUpdateStatus(ctx, ids, false)
	require.NoError(t, err)
	for _, id := range ids {
		item, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.False(t, item.IsActive)
	}
}

func TestServiceItemRepository_BatchUpdatePrice_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	var ids []uint64
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BPRICE%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Batch Price %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000 + int64(i*1000),
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
		ids = append(ids, item.ID)
	}
	newPrice := int64(25000)
	err := repo.BatchUpdatePrice(ctx, ids, newPrice)
	require.NoError(t, err)
	for _, id := range ids {
		item, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newPrice, item.BasePriceCents)
	}
}

func TestServiceItemRepository_BatchUpdateCommission_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	var ids []uint64
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BCOMM%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Batch Commission %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			CommissionRate: 0.15 + float64(i)*0.05,
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
		ids = append(ids, item.ID)
	}
	newRate := 0.25
	err := repo.BatchUpdateCommission(ctx, ids, newRate)
	require.NoError(t, err)
	for _, id := range ids {
		item, err := repo.Get(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, newRate, item.CommissionRate)
	}
}

func TestServiceItem_GameCategoryAssociation_CreateWithCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	category := CreateTestCategory(t, db, "Action")
	game := CreateTestGameWithCategoryID(t, db, "action_game", "Action Game", &category.ID)
	item := CreateTestServiceItemWithDetails(t, db, game, category, "Action Service", "ACT001", 10000, model.SubCategorySolo)
	assert.Equal(t, game.ID, *item.GameID)
	assert.Equal(t, category.ID, *item.CategoryID)
}

func TestServiceItem_GameCategoryAssociation_QueryByCategory(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	category1 := CreateTestCategory(t, db, "RPG")
	category2 := CreateTestCategory(t, db, "FPS")
	game1 := CreateTestGameWithCategoryID(t, db, "rpg_game", "RPG Game", &category1.ID)
	game2 := CreateTestGameWithCategoryID(t, db, "fps_game", "FPS Game", &category2.ID)
	CreateTestServiceItemWithDetails(t, db, game1, category1, "RPG Service 1", "RPG001", 10000, model.SubCategorySolo)
	CreateTestServiceItemWithDetails(t, db, game1, category1, "RPG Service 2", "RPG002", 15000, model.SubCategorySolo)
	CreateTestServiceItemWithDetails(t, db, game2, category2, "FPS Service", "FPS001", 12000, model.SubCategorySolo)
	var items []model.ServiceItem
	err := db.Where("category_id = ?", category1.ID).Find(&items).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2)
}

func TestServiceItem_GameAssociation_QueryByGame(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	game1 := CreateTestGame(t, db, "GameA")
	game2 := CreateTestGame(t, db, "GameB")
	CreateTestServiceItemWithDetails(t, db, game1, nil, "GameA Service 1", "GA001", 10000, model.SubCategorySolo)
	CreateTestServiceItemWithDetails(t, db, game1, nil, "GameA Service 2", "GA002", 15000, model.SubCategoryTeam)
	CreateTestServiceItemWithDetails(t, db, game2, nil, "GameB Service", "GB001", 12000, model.SubCategorySolo)
	var items []model.ServiceItem
	err := db.Where("game_id = ?", game1.ID).Find(&items).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(items), 2)
	for _, item := range items {
		assert.Equal(t, game1.ID, *item.GameID)
	}
}

func TestServiceItem_PriceManagement_BasePrice(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Base Price Test", "BP001", 25000, model.SubCategorySolo)
	assert.Equal(t, int64(25000), item.BasePriceCents)
}

func TestServiceItem_PriceManagement_VipPrice(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	basePrice := int64(30000)
	vipPrice := int64(24000)
	item := CreateTestServiceItemWithVIP(t, db, nil, "VIP Price Test", basePrice, vipPrice)
	assert.Equal(t, basePrice, item.BasePriceCents)
	assert.NotNil(t, item.VipPriceCents)
	assert.Equal(t, vipPrice, *item.VipPriceCents)
}

func TestServiceItem_PriceManagement_NoVipPrice(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItem(t, db, nil, "No VIP Test", 10000)
	assert.Nil(t, item.VipPriceCents)
}

func TestServiceItem_RequiredPlayers_Solo(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Solo Service", "SOLO001", 10000, model.SubCategorySolo)
	assert.Equal(t, 1, item.RequiredPlayers)
	assert.Equal(t, 1, item.MaxPlayers)
}

func TestServiceItem_RequiredPlayers_Team(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Team Service", "TEAM001", 50000, model.SubCategoryTeam)
	item.RequiredPlayers = 3
	item.MaxPlayers = 5
	db.Save(item)
	assert.Equal(t, 3, item.RequiredPlayers)
	assert.Equal(t, 5, item.MaxPlayers)
}

func TestServiceItem_MaxPerOrder_Unlimited(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItem(t, db, nil, "Unlimited Order", 20000)
	assert.Equal(t, 0, item.MaxPerOrder)
}

func TestServiceItem_MaxPerOrder_Limited(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "Limited Order", model.UsageLimitNone, 0)
	item.MaxPerOrder = 5
	db.Save(item)
	assert.Equal(t, 5, item.MaxPerOrder)
}

func TestServiceItem_UsageLimit_None(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "No Limit", model.UsageLimitNone, 0)
	assert.Equal(t, model.UsageLimitNone, item.UsageLimitType)
	assert.Equal(t, 0, item.UsageLimitCount)
}

func TestServiceItem_UsageLimit_Once(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "Once Only", model.UsageLimitOnce, 1)
	assert.Equal(t, model.UsageLimitOnce, item.UsageLimitType)
	assert.Equal(t, 1, item.UsageLimitCount)
}

func TestServiceItem_UsageLimit_Daily(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "Daily Limit", model.UsageLimitDaily, 5)
	assert.Equal(t, model.UsageLimitDaily, item.UsageLimitType)
	assert.Equal(t, 5, item.UsageLimitCount)
}

func TestServiceItem_UsageLimit_Weekly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "Weekly Limit", model.UsageLimitWeekly, 10)
	assert.Equal(t, model.UsageLimitWeekly, item.UsageLimitType)
	assert.Equal(t, 10, item.UsageLimitCount)
}

func TestServiceItem_UsageLimit_Monthly(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithUsageLimit(t, db, nil, "Monthly Limit", model.UsageLimitMonthly, 20)
	assert.Equal(t, model.UsageLimitMonthly, item.UsageLimitType)
	assert.Equal(t, 20, item.UsageLimitCount)
}

func TestServiceItem_CommissionRate_Default(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItem(t, db, nil, "Default Commission", 20000)
	assert.Equal(t, 0.20, item.CommissionRate)
}

func TestServiceItem_CommissionRate_Custom(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("COMM%d", time.Now().UnixNano()),
		Name:           "Custom Commission",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		CommissionRate: 0.25,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	assert.Equal(t, 0.25, item.CommissionRate)
}

func TestServiceItem_CommissionRate_Minimum(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("MINC%d", time.Now().UnixNano()),
		Name:           "Min Commission",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		CommissionRate: 0.0,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	assert.Equal(t, 0.0, item.CommissionRate)
}

func TestServiceItem_CommissionRate_Maximum(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("MAXC%d", time.Now().UnixNano()),
		Name:           "Max Commission",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		CommissionRate: 1.0,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	assert.Equal(t, 1.0, item.CommissionRate)
}

func TestServiceItem_CalculateCommission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("CALC%d", time.Now().UnixNano()),
		Name:           "Calculate Test",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		CommissionRate: 0.20,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	platformCommission, playerIncome := item.CalculateCommission(1)
	assert.Equal(t, int64(2000), platformCommission)
	assert.Equal(t, int64(8000), playerIncome)
	platformCommission, playerIncome = item.CalculateCommission(2)
	assert.Equal(t, int64(4000), platformCommission)
	assert.Equal(t, int64(16000), playerIncome)
}

func TestServiceItem_SubCategory_Solo(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Solo Escort", "SOLO1", 10000, model.SubCategorySolo)
	assert.Equal(t, model.SubCategorySolo, item.SubCategory)
	assert.False(t, item.IsGift())
}

func TestServiceItem_SubCategory_Team(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Team Escort", "TEAM1", 50000, model.SubCategoryTeam)
	assert.Equal(t, model.SubCategoryTeam, item.SubCategory)
	assert.False(t, item.IsGift())
}

func TestServiceItem_SubCategory_Gift(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItemWithDetails(t, db, nil, nil, "Gift Item", "GIFT1", 5000, model.SubCategoryGift)
	assert.Equal(t, model.SubCategoryGift, item.SubCategory)
	assert.True(t, item.IsGift())
	assert.Equal(t, 0, item.ServiceHours)
}

func TestServiceItem_SubCategory_QueryGifts(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("GIFT%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Gift %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 5000 + int64(i*1000),
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
	}
	CreateTestServiceItemWithDetails(t, db, nil, nil, "Solo Service", "SOLO1", 10000, model.SubCategorySolo)
	CreateTestServiceItemWithDetails(t, db, nil, nil, "Team Service", "TEAM1", 50000, model.SubCategoryTeam)
	gifts, total, err := repo.GetGifts(ctx, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, int64(3))
	assert.GreaterOrEqual(t, len(gifts), 3)
	for _, gift := range gifts {
		assert.True(t, gift.IsGift())
		assert.Equal(t, model.SubCategoryGift, gift.SubCategory)
	}
}

func TestServiceItem_SubCategory_QueryGameServices(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	game := CreateTestGame(t, db, "GameForServices")
	for i := 0; i < 2; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("SOLO%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Solo %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			GameID:         &game.ID,
			BasePriceCents: 10000,
			IsActive:       true,
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
	}
	for i := 0; i < 2; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("TEAM%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Team %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			GameID:         &game.ID,
			BasePriceCents: 50000,
			IsActive:       true,
			Tags:           "[]",
		}
		require.NoError(t, db.Create(item).Error)
	}
	services, err := repo.GetGameServices(ctx, game.ID, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(services), 4)
	soloSubCat := model.SubCategorySolo
	soloServices, err := repo.GetGameServices(ctx, game.ID, &soloSubCat)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(soloServices), 2)
	for _, svc := range soloServices {
		assert.Equal(t, model.SubCategorySolo, svc.SubCategory)
	}
}

func TestServiceItem_SortOrder_Default(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItem(t, db, nil, "Default Sort", 20000)
	assert.Equal(t, 0, item.SortOrder)
}

func TestServiceItem_SortOrder_Custom(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("SORT%d", time.Now().UnixNano()),
		Name:           "Custom Sort",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		SortOrder:      100,
		Tags:           "[]",
	}
	require.NoError(t, db.Create(item).Error)
	assert.Equal(t, 100, item.SortOrder)
}

func TestServiceItem_Timestamps(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	beforeCreate := time.Now()
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("TIME%d", time.Now().UnixNano()),
		Name:           "Timestamp Test",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		Tags:           "[]",
	}
	require.NoError(t, repo.Create(ctx, item))
	afterCreate := time.Now()
	assert.True(t, item.CreatedAt.After(beforeCreate) || item.CreatedAt.Equal(beforeCreate))
	assert.True(t, item.CreatedAt.Before(afterCreate) || item.CreatedAt.Equal(afterCreate))
	beforeUpdate := time.Now()
	item.Name = "Updated Name"
	err := repo.Update(ctx, item)
	require.NoError(t, err)
	afterUpdate := time.Now()
	fetched, _ := repo.Get(ctx, item.ID)
	assert.True(t, fetched.UpdatedAt.After(beforeUpdate) || fetched.UpdatedAt.Equal(beforeUpdate))
	assert.True(t, fetched.UpdatedAt.Before(afterUpdate) || fetched.UpdatedAt.Equal(afterUpdate))
}

func TestServiceItem_SoftDelete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("SOFT%d", time.Now().UnixNano()),
		Name:           "Soft Delete Test",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		Tags:           "[]",
	}
	require.NoError(t, repo.Create(ctx, item))
	err := repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.Get(ctx, item.ID)
	assert.Error(t, err)
	// ServiceItem uses hard delete - verify record is completely removed
	var item2 model.ServiceItem
	err = db.First(&item2, item.ID).Error
	assert.Error(t, err) // Record should not exist
}

func TestServiceItem_Pagination_FirstPage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	for i := 0; i < 25; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("PAGE1%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Page1 Item %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			Tags:           "[]",
		}
		require.NoError(t, repo.Create(ctx, item))
	}
	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(25), total)
	assert.Len(t, items, 10)
}

func TestServiceItem_Pagination_LastPage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	for i := 0; i < 25; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("PAGE3%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Page3 Item %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			Tags:           "[]",
		}
		require.NoError(t, repo.Create(ctx, item))
	}
	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{Page: 3, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(25), total)
	assert.LessOrEqual(t, len(items), 5)
}

func TestServiceItem_Pagination_BeyondLastPage(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		item := &model.ServiceItem{
			ItemCode:       fmt.Sprintf("BEYOND%d_%d", i, time.Now().UnixNano()),
			Name:           fmt.Sprintf("Beyond Item %d", i),
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			Tags:           "[]",
		}
		require.NoError(t, repo.Create(ctx, item))
	}
	items, total, err := repo.List(ctx, repository.ServiceItemListOptions{Page: 5, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(10), total)
	assert.Len(t, items, 0)
}

func TestServiceItem_Tags_EmptyArray(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	item := CreateTestServiceItem(t, db, nil, "Empty Tags", 20000)
	assert.Equal(t, "[]", item.Tags)
}

func TestServiceItem_Tags_WithValues(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	tags := `["popular", "premium", "verified"]`
	item := &model.ServiceItem{
		ItemCode:       fmt.Sprintf("TAG%d", time.Now().UnixNano()),
		Name:           "Tagged Service",
		Category:       "escort",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
		Tags:           tags,
	}
	require.NoError(t, db.Create(item).Error)
	assert.Equal(t, tags, item.Tags)
}

func TestServiceItem_IsActiveToggle_ActiveToInactive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	item := CreateTestServiceItemWithStatus(t, db, nil, "Toggle Test", true)
	assert.True(t, item.IsActive)
	err := repo.BatchUpdateStatus(ctx, []uint64{item.ID}, false)
	require.NoError(t, err)
	fetched, _ := repo.Get(ctx, item.ID)
	assert.False(t, fetched.IsActive)
}

func TestServiceItem_IsActiveToggle_InactiveToActive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := serviceitem.NewServiceItemRepository(db)
	ctx := context.Background()
	item := CreateTestServiceItemWithStatus(t, db, nil, "Toggle Test 2", false)
	assert.False(t, item.IsActive)
	err := repo.BatchUpdateStatus(ctx, []uint64{item.ID}, true)
	require.NoError(t, err)
	fetched, _ := repo.Get(ctx, item.ID)
	assert.True(t, fetched.IsActive)
}

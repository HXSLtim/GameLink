package serviceitem

import (
	"context"
	"testing"

	"gamelink/internal/model"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(&model.ServiceItem{})
	require.NoError(t, err)

	return db
}

func TestServiceItemRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	t.Run("创建单人护航服务", func(t *testing.T) {
		item := &model.ServiceItem{
			ItemCode:       "ESCORT_RANK_DIAMOND",
			Name:           "钻石段位护航",
			Category:       "escort",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 50000,
			ServiceHours:   1,
			CommissionRate: 0.20,
			IsActive:       true,
		}

		err := repo.Create(ctx, item)
		assert.NoError(t, err)
		assert.NotZero(t, item.ID)
	})

	t.Run("创建礼物", func(t *testing.T) {
		gift := &model.ServiceItem{
			ItemCode:       "ESCORT_GIFT_ROSE",
			Name:           "玫瑰花",
			Category:       "escort",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   0, // 礼物service_hours为0
			CommissionRate: 0.20,
			IsActive:       true,
		}

		err := repo.Create(ctx, gift)
		assert.NoError(t, err)
		assert.NotZero(t, gift.ID)
		assert.True(t, gift.IsGift())
	})

	t.Run("创建团队护航", func(t *testing.T) {
		team := &model.ServiceItem{
			ItemCode:       "ESCORT_TEAM_RANKED",
			Name:           "团队排位",
			Category:       "escort",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 80000,
			ServiceHours:   2,
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     3,
			IsActive:       true,
		}

		err := repo.Create(ctx, team)
		assert.NoError(t, err)
		assert.NotZero(t, team.ID)
	})
}

func TestServiceItemRepository_GetGifts(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建测试数据
	gift1 := &model.ServiceItem{
		ItemCode:       "GIFT_ROSE",
		Name:           "玫瑰",
		SubCategory:    model.SubCategoryGift,
		BasePriceCents: 10000,
		ServiceHours:   0,
		IsActive:       true,
	}
	repo.Create(ctx, gift1)

	escort := &model.ServiceItem{
		ItemCode:       "ESCORT_SOLO",
		Name:           "护航",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 50000,
		ServiceHours:   1,
		IsActive:       true,
	}
	repo.Create(ctx, escort)

	gift2 := &model.ServiceItem{
		ItemCode:       "GIFT_CHOCOLATE",
		Name:           "巧克力",
		SubCategory:    model.SubCategoryGift,
		BasePriceCents: 5000,
		ServiceHours:   0,
		IsActive:       true,
	}
	repo.Create(ctx, gift2)

	// 测试：只获取礼物
	gifts, total, err := repo.GetGifts(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, gifts, 2)

	// 验证都是礼物
	for i, gift := range gifts {
		t.Logf("Gift %d: ItemCode=%s, SubCategory=%s, ServiceHours=%d",
			i, gift.ItemCode, gift.SubCategory, gift.ServiceHours)
		assert.Equal(t, model.SubCategoryGift, gift.SubCategory)
		assert.Equal(t, 0, gift.ServiceHours, "Gift %s should have ServiceHours=0", gift.ItemCode)
	}
}

func TestServiceItemRepository_List_FilterBySubCategory(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建各种类型的服务
	types := []model.ServiceItemSubCategory{
		model.SubCategorySolo,
		model.SubCategoryTeam,
		model.SubCategoryGift,
	}

	for i, subCat := range types {
		item := &model.ServiceItem{
			ItemCode:       string(subCat) + "_TEST",
			Name:           string(subCat),
			SubCategory:    subCat,
			BasePriceCents: int64((i + 1) * 10000),
			ServiceHours:   0,
			IsActive:       true,
		}
		if subCat != model.SubCategoryGift {
			item.ServiceHours = 1
		}
		repo.Create(ctx, item)
	}

	// 测试：只获取solo类型
	soloType := model.SubCategorySolo
	soloItems, total, err := repo.List(ctx, ServiceItemListOptions{
		SubCategory: &soloType,
		Page:        1,
		PageSize:    10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, model.SubCategorySolo, soloItems[0].SubCategory)

	// 测试：获取所有类型
	allItems, total, err := repo.List(ctx, ServiceItemListOptions{
		Page:     1,
		PageSize: 10,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, allItems, 3)
}

func TestServiceItem_CalculateCommission(t *testing.T) {
	item := &model.ServiceItem{
		BasePriceCents: 10000,
		CommissionRate: 0.20, // 20%
	}

	// 测试：购买1个
	commission, playerIncome := item.CalculateCommission(1)
	assert.Equal(t, int64(2000), commission)   // 20%
	assert.Equal(t, int64(8000), playerIncome) // 80%

	// 测试：购买3个
	commission, playerIncome = item.CalculateCommission(3)
	assert.Equal(t, int64(6000), commission)    // 30000 * 20%
	assert.Equal(t, int64(24000), playerIncome) // 30000 * 80%
}

func TestServiceItemRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建测试数据
	item := &model.ServiceItem{
		ItemCode:       "TEST_ITEM",
		Name:           "测试服务",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)

	// 测试Get
	result, err := repo.Get(ctx, item.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TEST_ITEM", result.ItemCode)
}

func TestServiceItemRepository_GetByCode(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建测试数据
	item := &model.ServiceItem{
		ItemCode:       "UNIQUE_CODE",
		Name:           "测试服务",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)

	// 测试GetByCode
	result, err := repo.GetByCode(ctx, "UNIQUE_CODE")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "测试服务", result.Name)
}

func TestServiceItemRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建测试数据
	item := &model.ServiceItem{
		ItemCode:       "ITEM_UPDATE",
		Name:           "Original Name",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)

	// 更新
	item.Name = "Updated Name"
	item.BasePriceCents = 15000
	err = repo.Update(ctx, item)
	assert.NoError(t, err)

	// 验证更新
	updated, err := repo.Get(ctx, item.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, int64(15000), updated.BasePriceCents)
}

func TestServiceItemRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建测试数据
	item := &model.ServiceItem{
		ItemCode:       "TO_DELETE",
		Name:           "将被删除",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 10000,
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)

	// 删除
	err = repo.Delete(ctx, item.ID)
	assert.NoError(t, err)

	// 验证已删除（软删除）
	_, err = repo.Get(ctx, item.ID)
	assert.Error(t, err)
}

func TestServiceItemRepository_GetGameServices(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	gameID := uint64(1)

	// 创建测试数据
	items := []model.ServiceItem{
		{ItemCode: "G1_SOLO", Name: "单人服务", GameID: &gameID, SubCategory: model.SubCategorySolo, BasePriceCents: 10000, IsActive: true},
		{ItemCode: "G1_TEAM", Name: "组队服务", GameID: &gameID, SubCategory: model.SubCategoryTeam, BasePriceCents: 20000, IsActive: true},
		{ItemCode: "G2_SOLO", Name: "其他游戏", GameID: nil, SubCategory: model.SubCategorySolo, BasePriceCents: 10000, IsActive: true},
	}
	for _, item := range items {
		err := repo.Create(ctx, &item)
		require.NoError(t, err)
	}

	// 测试获取特定游戏的服务
	results, err := repo.GetGameServices(ctx, gameID, nil)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	// 测试按SubCategory过滤
	subCat := model.SubCategorySolo
	results2, err := repo.GetGameServices(ctx, gameID, &subCat)
	assert.NoError(t, err)
	assert.Len(t, results2, 1)
	assert.Equal(t, "单人服务", results2[0].Name)
}

func TestServiceItemRepository_BatchOperations(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	// 创建多个服务
	ids := make([]uint64, 0)
	for i := 0; i < 3; i++ {
		item := &model.ServiceItem{
			ItemCode:       "TEST_" + string(rune(i+'0')),
			Name:           "Test Item",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			IsActive:       true,
		}
		repo.Create(ctx, item)
		ids = append(ids, item.ID)
	}

	// 测试：批量更新状态
	err := repo.BatchUpdateStatus(ctx, ids, false)
	assert.NoError(t, err)

	// 验证
	for _, id := range ids {
		item, _ := repo.Get(ctx, id)
		assert.False(t, item.IsActive)
	}

	// 测试：批量更新价格
	err = repo.BatchUpdatePrice(ctx, ids, 15000)
	assert.NoError(t, err)

	// 验证
	for _, id := range ids {
		item, _ := repo.Get(ctx, id)
		assert.Equal(t, int64(15000), item.BasePriceCents)
	}
}

// TestServiceItemRepository_List_WithFilters 测试带多种过滤条件的列表查询
func TestServiceItemRepository_List_WithFilters(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	gameID1 := uint64(1)
	gameID2 := uint64(2)

	// 创建测试数据
	items := []model.ServiceItem{
		{ItemCode: "G1_ACTIVE", Name: "Active Item 1", GameID: &gameID1, BasePriceCents: 10000, IsActive: true},
		{ItemCode: "G1_INACTIVE", Name: "Inactive Item", GameID: &gameID1, BasePriceCents: 20000, IsActive: false},
		{ItemCode: "G2_ACTIVE", Name: "Active Item 2", GameID: &gameID2, BasePriceCents: 15000, IsActive: true},
		{ItemCode: "NO_GAME", Name: "No Game Item", GameID: nil, BasePriceCents: 5000, IsActive: true},
	}
	for _, item := range items {
		err := repo.Create(ctx, &item)
		require.NoError(t, err)
	}

	t.Run("按游戏ID过滤", func(t *testing.T) {
		results, total, err := repo.List(ctx, ServiceItemListOptions{
			GameID:   &gameID1,
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), total) // 游戏1有2个服务
		for _, item := range results {
			assert.NotNil(t, item.GameID)
			assert.Equal(t, gameID1, *item.GameID)
		}
	})

	t.Run("只获取激活的服务", func(t *testing.T) {
		isActive := true
		results, total, err := repo.List(ctx, ServiceItemListOptions{
			IsActive: &isActive,
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.True(t, total >= 3) // 至少有3个激活的服务
		for _, item := range results {
			assert.True(t, item.IsActive)
		}
	})

	t.Run("组合过滤条件", func(t *testing.T) {
		isActive := true
		results, total, err := repo.List(ctx, ServiceItemListOptions{
			GameID:   &gameID1,
			IsActive: &isActive,
			Page:     1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total) // 游戏1只有1个激活的服务
		assert.True(t, results[0].IsActive)
		assert.Equal(t, gameID1, *results[0].GameID)
	})
}

// TestServiceItemRepository_GetGameServices_EdgeCases 测试GetGameServices的边界条件
func TestServiceItemRepository_GetGameServices_EdgeCases(t *testing.T) {
	db := setupTestDB(t)
	repo := NewServiceItemRepository(db)
	ctx := context.Background()

	t.Run("查询不存在的游戏ID", func(t *testing.T) {
		nonExistentGameID := uint64(99999)
		results, err := repo.GetGameServices(ctx, nonExistentGameID, nil)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("查询存在但无服务的游戏", func(t *testing.T) {
		// 创建一个游戏ID但没有对应的服务
		emptyGameID := uint64(888)
		results, err := repo.GetGameServices(ctx, emptyGameID, nil)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("按SubCategory过滤不存在的类型", func(t *testing.T) {
		gameID := uint64(1)
		// 创建solo类型的服务
		item := &model.ServiceItem{
			ItemCode:       "SOLO_ONLY",
			Name:           "Solo Service",
			GameID:         &gameID,
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			IsActive:       true,
		}
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		// 查询team类型（不存在）
		teamType := model.SubCategoryTeam
		results, err := repo.GetGameServices(ctx, gameID, &teamType)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("查询已禁用的服务", func(t *testing.T) {
		gameID := uint64(2)
		// 创建禁用的服务
		item := &model.ServiceItem{
			ItemCode:       "INACTIVE_SERVICE",
			Name:           "Inactive",
			GameID:         &gameID,
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			IsActive:       false, // 禁用
		}
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		// GetGameServices应该只返回激活的服务
		results, err := repo.GetGameServices(ctx, gameID, nil)
		assert.NoError(t, err)
		// 如果实现中过滤了IsActive，则结果应该为空或只包含激活的
		for _, result := range results {
			assert.True(t, result.IsActive, "GetGameServices should only return active services")
		}
	})
}

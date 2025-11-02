package repository

import (
	"context"
	"testing"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
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
	for _, gift := range gifts {
		assert.Equal(t, model.SubCategoryGift, gift.SubCategory)
		assert.Equal(t, 0, gift.ServiceHours)
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
	assert.Equal(t, int64(2000), commission)      // 20%
	assert.Equal(t, int64(8000), playerIncome)    // 80%

	// 测试：购买3个
	commission, playerIncome = item.CalculateCommission(3)
	assert.Equal(t, int64(6000), commission)      // 30000 * 20%
	assert.Equal(t, int64(24000), playerIncome)   // 30000 * 80%
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


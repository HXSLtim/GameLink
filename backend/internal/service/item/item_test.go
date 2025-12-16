package item

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/pkg/testutil"
)

func setupItemService(t *testing.T) (*ServiceItemService, context.Context) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.ServiceItem{},
		&model.Game{},
		&model.Player{},
		&model.User{},
	)

	itemRepo := serviceitem.NewServiceItemRepository(db)
	gameRepo := game.NewGameRepository(db)
	playerRepo := player.NewPlayerRepository(db)

	service := NewServiceItemService(itemRepo, gameRepo, playerRepo)
	return service, context.Background()
}

func TestServiceItemService_CreateServiceItem(t *testing.T) {
	service, ctx := setupItemService(t)

	t.Run("create solo service item", func(t *testing.T) {
		req := CreateServiceItemRequest{
			ItemCode:       "SOLO001",
			Name:           "单人护航服务",
			Description:    "专业单人护航",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 5000,
			ServiceHours:   2,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := service.CreateServiceItem(ctx, req)
		require.NoError(t, err)
		assert.NotZero(t, item.ID)
		assert.Equal(t, "SOLO001", item.ItemCode)
		assert.Equal(t, "单人护航服务", item.Name)
		assert.Equal(t, model.SubCategorySolo, item.SubCategory)
		assert.True(t, item.IsActive)
	})

	t.Run("create team service item", func(t *testing.T) {
		req := CreateServiceItemRequest{
			ItemCode:       "TEAM001",
			Name:           "团队护航服务",
			Description:    "专业团队护航",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 10000,
			ServiceHours:   3,
			CommissionRate: 0.25,
			MinUsers:       2,
			MaxPlayers:     5,
		}

		item, err := service.CreateServiceItem(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, model.SubCategoryTeam, item.SubCategory)
		assert.Equal(t, 5, item.MaxPlayers)
	})

	t.Run("create gift item", func(t *testing.T) {
		req := CreateServiceItemRequest{
			ItemCode:       "GIFT001",
			Name:           "玫瑰花",
			Description:    "送给喜欢的陪玩师",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 100,
			ServiceHours:   0, // 礼物必须为0
			CommissionRate: 0.3,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := service.CreateServiceItem(ctx, req)
		require.NoError(t, err)
		assert.True(t, item.IsGift())
		assert.Equal(t, 0, item.ServiceHours)
	})

	t.Run("gift with non-zero service hours should fail", func(t *testing.T) {
		req := CreateServiceItemRequest{
			ItemCode:       "GIFT002",
			Name:           "错误礼物",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 100,
			ServiceHours:   1, // 礼物不能有服务时长
			CommissionRate: 0.3,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		_, err := service.CreateServiceItem(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "礼物类项目的服务时长必须为0")
	})
}

func TestServiceItemService_UpdateServiceItem(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	item, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode:       "UPDATE001",
		Name:           "原始名称",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 5000,
		ServiceHours:   2,
		CommissionRate: 0.2,
		MinUsers:       1,
		MaxPlayers:     1,
	})
	require.NoError(t, err)

	t.Run("update name and description", func(t *testing.T) {
		newName := "更新后的名称"
		newDesc := "更新后的描述"
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			Name:        &newName,
			Description: &newDesc,
		})
		require.NoError(t, err)

		updated, err := service.GetServiceItem(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, newName, updated.Name)
		assert.Equal(t, newDesc, updated.Description)
	})

	t.Run("update price", func(t *testing.T) {
		newPrice := int64(8000)
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			BasePriceCents: &newPrice,
		})
		require.NoError(t, err)

		updated, err := service.GetServiceItem(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, newPrice, updated.BasePriceCents)
	})

	t.Run("update with negative price should fail", func(t *testing.T) {
		negativePrice := int64(-100)
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			BasePriceCents: &negativePrice,
		})
		assert.Error(t, err)
	})

	t.Run("update commission rate", func(t *testing.T) {
		newRate := 0.15
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			CommissionRate: &newRate,
		})
		require.NoError(t, err)

		updated, err := service.GetServiceItem(ctx, item.ID)
		require.NoError(t, err)
		assert.Equal(t, newRate, updated.CommissionRate)
	})

	t.Run("update with invalid commission rate should fail", func(t *testing.T) {
		invalidRate := 1.5 // 超过1
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			CommissionRate: &invalidRate,
		})
		assert.Error(t, err)
	})

	t.Run("update non-existent item should fail", func(t *testing.T) {
		newName := "不存在"
		err := service.UpdateServiceItem(ctx, 99999, UpdateServiceItemRequest{
			Name: &newName,
		})
		assert.Error(t, err)
	})
}

func TestServiceItemService_DeleteServiceItem(t *testing.T) {
	service, ctx := setupItemService(t)

	t.Run("delete existing item", func(t *testing.T) {
		item, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
			ItemCode:       "DELETE001",
			Name:           "待删除",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 5000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		})
		require.NoError(t, err)

		err = service.DeleteServiceItem(ctx, item.ID)
		require.NoError(t, err)

		_, err = service.GetServiceItem(ctx, item.ID)
		assert.Error(t, err)
	})

	t.Run("delete non-existent item should fail", func(t *testing.T) {
		err := service.DeleteServiceItem(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestServiceItemService_ListServiceItems(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	items := []CreateServiceItemRequest{
		{ItemCode: "LIST001", Name: "服务1", SubCategory: model.SubCategorySolo, BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1},
		{ItemCode: "LIST002", Name: "服务2", SubCategory: model.SubCategoryTeam, BasePriceCents: 10000, ServiceHours: 2, CommissionRate: 0.25, MinUsers: 2, MaxPlayers: 5},
		{ItemCode: "LIST003", Name: "礼物1", SubCategory: model.SubCategoryGift, BasePriceCents: 100, ServiceHours: 0, CommissionRate: 0.3, MinUsers: 1, MaxPlayers: 1},
	}
	for _, req := range items {
		_, err := service.CreateServiceItem(ctx, req)
		require.NoError(t, err)
	}

	t.Run("list all items", func(t *testing.T) {
		resp, err := service.ListServiceItems(ctx, ListServiceItemsRequest{
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Items, 3)
	})

	t.Run("list by sub category", func(t *testing.T) {
		subCat := model.SubCategoryGift
		resp, err := service.ListServiceItems(ctx, ListServiceItemsRequest{
			Page:        1,
			PageSize:    10,
			SubCategory: &subCat,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "gift", resp.Items[0].SubCategory)
	})

	t.Run("list with pagination", func(t *testing.T) {
		resp, err := service.ListServiceItems(ctx, ListServiceItemsRequest{
			Page:     1,
			PageSize: 2,
		})
		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
		assert.Len(t, resp.Items, 2)
	})
}

func TestServiceItemService_GetGiftList(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	_, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "GIFTLIST001", Name: "礼物A", SubCategory: model.SubCategoryGift,
		BasePriceCents: 100, ServiceHours: 0, CommissionRate: 0.3, MinUsers: 1, MaxPlayers: 1,
	})
	require.NoError(t, err)

	_, err = service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "GIFTLIST002", Name: "礼物B", SubCategory: model.SubCategoryGift,
		BasePriceCents: 200, ServiceHours: 0, CommissionRate: 0.3, MinUsers: 1, MaxPlayers: 1,
	})
	require.NoError(t, err)

	// 创建非礼物项目
	_, err = service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "NOTGIFT001", Name: "服务", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})
	require.NoError(t, err)

	t.Run("get gift list", func(t *testing.T) {
		resp, err := service.GetGiftList(ctx, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
		for _, item := range resp.Items {
			assert.Equal(t, "gift", item.SubCategory)
		}
	})
}

func TestServiceItemService_BatchUpdateStatus(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	item1, _ := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "BATCH001", Name: "批量1", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})
	item2, _ := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "BATCH002", Name: "批量2", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})

	t.Run("batch update status to inactive", func(t *testing.T) {
		err := service.BatchUpdateStatus(ctx, BatchUpdateStatusRequest{
			IDs:      []uint64{item1.ID, item2.ID},
			IsActive: false,
		})
		require.NoError(t, err)

		updated1, _ := service.GetServiceItem(ctx, item1.ID)
		updated2, _ := service.GetServiceItem(ctx, item2.ID)
		assert.False(t, updated1.IsActive)
		assert.False(t, updated2.IsActive)
	})

	t.Run("batch update with empty IDs should fail", func(t *testing.T) {
		err := service.BatchUpdateStatus(ctx, BatchUpdateStatusRequest{
			IDs:      []uint64{},
			IsActive: true,
		})
		assert.Error(t, err)
	})
}

func TestServiceItemService_BatchUpdatePrice(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	item1, _ := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "PRICE001", Name: "价格1", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})
	item2, _ := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "PRICE002", Name: "价格2", SubCategory: model.SubCategorySolo,
		BasePriceCents: 6000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})

	t.Run("batch update price", func(t *testing.T) {
		err := service.BatchUpdatePrice(ctx, BatchUpdatePriceRequest{
			IDs:            []uint64{item1.ID, item2.ID},
			BasePriceCents: 8000,
		})
		require.NoError(t, err)

		updated1, _ := service.GetServiceItem(ctx, item1.ID)
		updated2, _ := service.GetServiceItem(ctx, item2.ID)
		assert.Equal(t, int64(8000), updated1.BasePriceCents)
		assert.Equal(t, int64(8000), updated2.BasePriceCents)
	})

	t.Run("batch update price with empty IDs should fail", func(t *testing.T) {
		err := service.BatchUpdatePrice(ctx, BatchUpdatePriceRequest{
			IDs:            []uint64{},
			BasePriceCents: 9000,
		})
		assert.Error(t, err)
	})
}

func TestServiceItemService_CreateWithInvalidGameID(t *testing.T) {
	service, ctx := setupItemService(t)

	t.Run("create with invalid game ID", func(t *testing.T) {
		invalidGameID := uint64(99999)
		req := CreateServiceItemRequest{
			ItemCode:       "INVALID001",
			Name:           "无效游戏",
			SubCategory:    model.SubCategorySolo,
			GameID:         &invalidGameID,
			BasePriceCents: 5000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		_, err := service.CreateServiceItem(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "游戏ID无效")
	})
}

func TestServiceItemService_CreateWithInvalidPlayerID(t *testing.T) {
	service, ctx := setupItemService(t)

	t.Run("create with invalid player ID", func(t *testing.T) {
		invalidPlayerID := uint64(99999)
		req := CreateServiceItemRequest{
			ItemCode:       "INVALID002",
			Name:           "无效陪玩师",
			SubCategory:    model.SubCategorySolo,
			PlayerID:       &invalidPlayerID,
			BasePriceCents: 5000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		_, err := service.CreateServiceItem(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "陪玩师ID无效")
	})
}

func TestServiceItemService_UpdateGiftServiceHours(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建礼物项目
	gift, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode:       "GIFTUPDATE001",
		Name:           "礼物更新测试",
		SubCategory:    model.SubCategoryGift,
		BasePriceCents: 100,
		ServiceHours:   0,
		CommissionRate: 0.3,
		MinUsers:       1,
		MaxPlayers:     1,
	})
	require.NoError(t, err)

	t.Run("update gift with non-zero service hours should fail", func(t *testing.T) {
		nonZeroHours := 1
		err := service.UpdateServiceItem(ctx, gift.ID, UpdateServiceItemRequest{
			ServiceHours: &nonZeroHours,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "礼物类项目的服务时长必须为0")
	})

	t.Run("update gift with zero service hours should succeed", func(t *testing.T) {
		zeroHours := 0
		err := service.UpdateServiceItem(ctx, gift.ID, UpdateServiceItemRequest{
			ServiceHours: &zeroHours,
		})
		require.NoError(t, err)
	})
}

func TestServiceItemService_UpdateOptionalFields(t *testing.T) {
	service, ctx := setupItemService(t)

	item, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode:       "OPTIONAL001",
		Name:           "可选字段测试",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 5000,
		ServiceHours:   1,
		CommissionRate: 0.2,
		MinUsers:       1,
		MaxPlayers:     1,
	})
	require.NoError(t, err)

	t.Run("update rank level", func(t *testing.T) {
		rankLevel := "钻石"
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			RankLevel: &rankLevel,
		})
		require.NoError(t, err)

		updated, _ := service.GetServiceItem(ctx, item.ID)
		assert.Equal(t, "钻石", updated.RankLevel)
	})

	t.Run("update tags", func(t *testing.T) {
		tags := "热门,推荐"
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			Tags: &tags,
		})
		require.NoError(t, err)

		updated, _ := service.GetServiceItem(ctx, item.ID)
		assert.Equal(t, "热门,推荐", updated.Tags)
	})

	t.Run("update icon URL", func(t *testing.T) {
		iconURL := "https://example.com/icon.png"
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			IconURL: &iconURL,
		})
		require.NoError(t, err)

		updated, _ := service.GetServiceItem(ctx, item.ID)
		assert.Equal(t, "https://example.com/icon.png", updated.IconURL)
	})

	t.Run("update sort order", func(t *testing.T) {
		sortOrder := 10
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			SortOrder: &sortOrder,
		})
		require.NoError(t, err)

		updated, _ := service.GetServiceItem(ctx, item.ID)
		assert.Equal(t, 10, updated.SortOrder)
	})

	t.Run("update is active", func(t *testing.T) {
		isActive := false
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			IsActive: &isActive,
		})
		require.NoError(t, err)

		updated, _ := service.GetServiceItem(ctx, item.ID)
		assert.False(t, updated.IsActive)
	})
}

func TestServiceItemService_UpdateCommissionRateEdgeCases(t *testing.T) {
	service, ctx := setupItemService(t)

	item, err := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode:       "COMMISSION001",
		Name:           "佣金测试",
		SubCategory:    model.SubCategorySolo,
		BasePriceCents: 5000,
		ServiceHours:   1,
		CommissionRate: 0.2,
		MinUsers:       1,
		MaxPlayers:     1,
	})
	require.NoError(t, err)

	t.Run("update with negative commission rate should fail", func(t *testing.T) {
		negativeRate := -0.1
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			CommissionRate: &negativeRate,
		})
		assert.Error(t, err)
	})

	t.Run("update with zero commission rate should succeed", func(t *testing.T) {
		zeroRate := 0.0
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			CommissionRate: &zeroRate,
		})
		require.NoError(t, err)
	})

	t.Run("update with max commission rate should succeed", func(t *testing.T) {
		maxRate := 1.0
		err := service.UpdateServiceItem(ctx, item.ID, UpdateServiceItemRequest{
			CommissionRate: &maxRate,
		})
		require.NoError(t, err)
	})
}

func TestServiceItemService_ListWithFilters(t *testing.T) {
	service, ctx := setupItemService(t)

	// 创建测试数据
	active := true
	inactive := false

	_, _ = service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "FILTER001", Name: "活跃服务", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})

	item2, _ := service.CreateServiceItem(ctx, CreateServiceItemRequest{
		ItemCode: "FILTER002", Name: "非活跃服务", SubCategory: model.SubCategorySolo,
		BasePriceCents: 5000, ServiceHours: 1, CommissionRate: 0.2, MinUsers: 1, MaxPlayers: 1,
	})
	// 设置为非活跃
	_ = service.UpdateServiceItem(ctx, item2.ID, UpdateServiceItemRequest{IsActive: &inactive})

	t.Run("list active items only", func(t *testing.T) {
		resp, err := service.ListServiceItems(ctx, ListServiceItemsRequest{
			Page:     1,
			PageSize: 10,
			IsActive: &active,
		})
		require.NoError(t, err)
		for _, item := range resp.Items {
			assert.True(t, item.IsActive)
		}
	})

	t.Run("list inactive items only", func(t *testing.T) {
		resp, err := service.ListServiceItems(ctx, ListServiceItemsRequest{
			Page:     1,
			PageSize: 10,
			IsActive: &inactive,
		})
		require.NoError(t, err)
		for _, item := range resp.Items {
			assert.False(t, item.IsActive)
		}
	})
}

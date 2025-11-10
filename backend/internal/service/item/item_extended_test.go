package item

import (
	"context"
	"errors"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateServiceItem_EdgeCases 测试创建服务项目的边界情况
func TestCreateServiceItem_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("抽成率为0应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			return item.CommissionRate == 0.0
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "ZERO_COMMISSION",
			Name:           "零抽成服务",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 0.0, // 零抽成
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, 0.0, item.CommissionRate)
	})

	t.Run("抽成率为1应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			return item.CommissionRate == 1.0
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "FULL_COMMISSION",
			Name:           "全额抽成服务",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 1.0, // 100%抽成
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, 1.0, item.CommissionRate)
	})

	t.Run("价格为0应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			return item.BasePriceCents == 0
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "FREE_SERVICE",
			Name:           "免费服务",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 0, // 免费
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, int64(0), item.BasePriceCents)
	})

	t.Run("服务时长为0的非礼物项目应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			return item.ServiceHours == 0 && item.SubCategory != model.SubCategoryGift
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "INSTANT_SERVICE",
			Name:           "即时服务",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			ServiceHours:   0, // 非礼物但服务时长为0
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, item)
	})

	t.Run("玩家ID有效时应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		playerID := uint64(100)
		player := &model.Player{
			Base:     model.Base{ID: playerID},
			Nickname: "测试玩家",
		}

		playerRepo.On("Get", ctx, playerID).Return(player, nil)
		itemRepo.On("Create", ctx, mock.Anything).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "PLAYER_SERVICE",
			Name:           "玩家专属服务",
			SubCategory:    model.SubCategorySolo,
			PlayerID:       &playerID,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, playerID, *item.PlayerID)
	})

	t.Run("玩家ID无效时应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		playerID := uint64(999)
		playerRepo.On("Get", ctx, playerID).Return(nil, repository.ErrNotFound)

		req := CreateServiceItemRequest{
			ItemCode:       "INVALID_PLAYER",
			Name:           "无效玩家服务",
			SubCategory:    model.SubCategorySolo,
			PlayerID:       &playerID,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "invalid player_id")
	})

	t.Run("数据库创建失败应该返回错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		dbErr := errors.New("database connection failed")
		itemRepo.On("Create", ctx, mock.Anything).Return(dbErr)

		req := CreateServiceItemRequest{
			ItemCode:       "DB_ERROR",
			Name:           "数据库错误",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 0.2,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Equal(t, dbErr, err)
	})
}

// TestUpdateServiceItem_EdgeCases 测试更新服务项目的边界情况
func TestUpdateServiceItem_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("更新价格为负数应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:             1,
			ItemCode:       "TEST",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)

		negativePrice := int64(-100)
		req := UpdateServiceItemRequest{
			BasePriceCents: &negativePrice,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "base price must be >= 0")
	})

	t.Run("更新抽成率超过1应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:             1,
			ItemCode:       "TEST",
			SubCategory:    model.SubCategorySolo,
			CommissionRate: 0.2,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)

		invalidRate := 1.5
		req := UpdateServiceItemRequest{
			CommissionRate: &invalidRate,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commission rate must be between 0 and 1")
	})

	t.Run("更新抽成率为负数应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:             1,
			ItemCode:       "TEST",
			SubCategory:    model.SubCategorySolo,
			CommissionRate: 0.2,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)

		negativeRate := -0.1
		req := UpdateServiceItemRequest{
			CommissionRate: &negativeRate,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "commission rate must be between 0 and 1")
	})

	t.Run("更新礼物的服务时长为非0应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:           1,
			ItemCode:     "GIFT_ROSE",
			SubCategory:  model.SubCategoryGift,
			ServiceHours: 0,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)

		invalidHours := 1
		req := UpdateServiceItemRequest{
			ServiceHours: &invalidHours,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "gift items must have service_hours = 0")
	})

	t.Run("更新礼物的服务时长为0应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:           1,
			ItemCode:     "GIFT_ROSE",
			SubCategory:  model.SubCategoryGift,
			ServiceHours: 0,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
		itemRepo.On("Update", ctx, mock.Anything).Return(nil)

		validHours := 0
		req := UpdateServiceItemRequest{
			ServiceHours: &validHours,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.NoError(t, err)
	})

	t.Run("更新不存在的项目应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		name := "New Name"
		req := UpdateServiceItemRequest{
			Name: &name,
		}

		err := svc.UpdateServiceItem(ctx, 999, req)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("只更新部分字段应该成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:             1,
			ItemCode:       "TEST",
			Name:           "Original Name",
			Description:    "Original Desc",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
			CommissionRate: 0.2,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
		itemRepo.On("Update", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			// 只有Name被更新
			return item.Name == "New Name" &&
				item.Description == "Original Desc" &&
				item.BasePriceCents == 10000
		})).Return(nil)

		newName := "New Name"
		req := UpdateServiceItemRequest{
			Name: &newName,
			// 其他字段不更新
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.NoError(t, err)
	})

	t.Run("数据库更新失败应该返回错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:          1,
			ItemCode:    "TEST",
			SubCategory: model.SubCategorySolo,
		}

		dbErr := errors.New("database update failed")
		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
		itemRepo.On("Update", ctx, mock.Anything).Return(dbErr)

		name := "New Name"
		req := UpdateServiceItemRequest{
			Name: &name,
		}

		err := svc.UpdateServiceItem(ctx, 1, req)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}

// TestDeleteServiceItem_EdgeCases 测试删除服务项目的边界情况
func TestDeleteServiceItem_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("删除不存在的项目应该失败", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		err := svc.DeleteServiceItem(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
	})

	t.Run("数据库删除失败应该返回错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		existingItem := &model.ServiceItem{
			ID:       1,
			ItemCode: "TEST",
		}

		dbErr := errors.New("foreign key constraint failed")
		itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
		itemRepo.On("Delete", ctx, uint64(1)).Return(dbErr)

		err := svc.DeleteServiceItem(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}

// TestBatchOperations_EdgeCases 测试批量操作的边界情况
func TestBatchOperations_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("批量更新状态-空ID列表", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		req := BatchUpdateStatusRequest{
			IDs:      []uint64{},
			IsActive: true,
		}

		err := svc.BatchUpdateStatus(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no item ids provided")
	})

	t.Run("批量更新状态-nil ID列表", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		req := BatchUpdateStatusRequest{
			IDs:      nil,
			IsActive: true,
		}

		err := svc.BatchUpdateStatus(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no item ids provided")
	})

	t.Run("批量更新价格-空ID列表", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		req := BatchUpdatePriceRequest{
			IDs:            []uint64{},
			BasePriceCents: 10000,
		}

		err := svc.BatchUpdatePrice(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no item ids provided")
	})

	t.Run("批量更新状态-单个ID", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		ids := []uint64{1}
		itemRepo.On("BatchUpdateStatus", ctx, ids, true).Return(nil)

		req := BatchUpdateStatusRequest{
			IDs:      ids,
			IsActive: true,
		}

		err := svc.BatchUpdateStatus(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("批量更新价格-大量ID", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		// 生成100个ID
		ids := make([]uint64, 100)
		for i := range ids {
			ids[i] = uint64(i + 1)
		}

		itemRepo.On("BatchUpdatePrice", ctx, ids, int64(50000)).Return(nil)

		req := BatchUpdatePriceRequest{
			IDs:            ids,
			BasePriceCents: 50000,
		}

		err := svc.BatchUpdatePrice(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("批量更新状态-数据库错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		ids := []uint64{1, 2, 3}
		dbErr := errors.New("database error")
		itemRepo.On("BatchUpdateStatus", ctx, ids, false).Return(dbErr)

		req := BatchUpdateStatusRequest{
			IDs:      ids,
			IsActive: false,
		}

		err := svc.BatchUpdateStatus(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})

	t.Run("批量更新价格-数据库错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		ids := []uint64{1, 2}
		dbErr := errors.New("constraint violation")
		itemRepo.On("BatchUpdatePrice", ctx, ids, int64(10000)).Return(dbErr)

		req := BatchUpdatePriceRequest{
			IDs:            ids,
			BasePriceCents: 10000,
		}

		err := svc.BatchUpdatePrice(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
	})
}

// TestListServiceItems_EdgeCases 测试列表查询的边界情况
func TestListServiceItems_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("空结果集应该返回空列表", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("List", ctx, mock.Anything).Return([]model.ServiceItem{}, int64(0), nil)

		req := ListServiceItemsRequest{
			Page:     1,
			PageSize: 20,
		}

		result, err := svc.ListServiceItems(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Items)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("数据库查询失败应该返回错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		dbErr := errors.New("database query failed")
		itemRepo.On("List", ctx, mock.Anything).Return([]model.ServiceItem{}, int64(0), dbErr)

		req := ListServiceItemsRequest{
			Page:     1,
			PageSize: 20,
		}

		result, err := svc.ListServiceItems(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, dbErr, err)
	})
}

// TestGetGiftList_EdgeCases 测试礼物列表的边界情况
func TestGetGiftList_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("空礼物列表应该返回空结果", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("GetGifts", ctx, 1, 20).Return([]model.ServiceItem{}, int64(0), nil)

		result, err := svc.GetGiftList(ctx, 1, 20)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Items)
		assert.Equal(t, int64(0), result.Total)
	})

	t.Run("数据库错误应该返回错误", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		dbErr := errors.New("connection timeout")
		itemRepo.On("GetGifts", ctx, 1, 20).Return([]model.ServiceItem{}, int64(0), dbErr)

		result, err := svc.GetGiftList(ctx, 1, 20)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, dbErr, err)
	})

	t.Run("大页码应该正常处理", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)
		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("GetGifts", ctx, 1000, 20).Return([]model.ServiceItem{}, int64(0), nil)

		result, err := svc.GetGiftList(ctx, 1000, 20)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

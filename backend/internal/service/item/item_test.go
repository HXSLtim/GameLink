package item

import (
	"gamelink/internal/model"
	"gamelink/internal/repository"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	"testing"
	"errors"
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type MockServiceItemRepo struct {
	mock.Mock
}

func (m *MockServiceItemRepo) Create(ctx context.Context, item *model.ServiceItem) error {
	args := m.Called(ctx, item)
	if args.Get(0) != nil {
		item.ID = 1
	}
	return args.Error(0)
}

func (m *MockServiceItemRepo) Get(ctx context.Context, id uint64) (*model.ServiceItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServiceItem), args.Error(1)
}

func (m *MockServiceItemRepo) GetByCode(ctx context.Context, code string) (*model.ServiceItem, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ServiceItem), args.Error(1)
}

func (m *MockServiceItemRepo) List(ctx context.Context, opts serviceitemrepo.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ServiceItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockServiceItemRepo) Update(ctx context.Context, item *model.ServiceItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockServiceItemRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockServiceItemRepo) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) error {
	args := m.Called(ctx, ids, isActive)
	return args.Error(0)
}

func (m *MockServiceItemRepo) BatchUpdatePrice(ctx context.Context, ids []uint64, price int64) error {
	args := m.Called(ctx, ids, price)
	return args.Error(0)
}

func (m *MockServiceItemRepo) GetGifts(ctx context.Context, page, pageSize int) ([]model.ServiceItem, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.ServiceItem), args.Get(1).(int64), args.Error(2)
}

func (m *MockServiceItemRepo) GetGameServices(ctx context.Context, gameID uint64, subCat *model.ServiceItemSubCategory) ([]model.ServiceItem, error) {
	args := m.Called(ctx, gameID, subCat)
	return args.Get(0).([]model.ServiceItem), args.Error(1)
}

type MockGameRepo struct {
	mock.Mock
}

func (m *MockGameRepo) Get(ctx context.Context, id uint64) (*model.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Game), args.Error(1)
}

func (m *MockGameRepo) Create(ctx context.Context, game *model.Game) error { return nil }
func (m *MockGameRepo) Update(ctx context.Context, game *model.Game) error { return nil }

func (m *MockGameRepo) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *MockGameRepo) List(ctx context.Context) ([]model.Game, error)     { return nil, nil }

func (m *MockGameRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return nil, 0, nil
}

func (m *MockGameRepo) GetByKey(ctx context.Context, key string) (*model.Game, error) {
	return nil, nil
}

type MockPlayerRepo struct {
	mock.Mock
}

func (m *MockPlayerRepo) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepo) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *MockPlayerRepo) Update(ctx context.Context, player *model.Player) error { return nil }

func (m *MockPlayerRepo) Delete(ctx context.Context, id uint64) error            { return nil }
func (m *MockPlayerRepo) List(ctx context.Context) ([]model.Player, error)       { return nil, nil }

func (m *MockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *MockPlayerRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func TestServiceItemService_CreateServiceItem(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	t.Run("创建护航服务成功", func(t *testing.T) {
		game := &model.Game{
			Base: model.Base{
				ID: 1,
			},
			Name: "王者荣耀",
		}
		gameID := uint64(1)

		gameRepo.On("Get", ctx, uint64(1)).Return(game, nil)
		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			assert.Equal(t, "ESCORT_RANK_DIAMOND", item.ItemCode)
			assert.Equal(t, model.SubCategorySolo, item.SubCategory)
			assert.Equal(t, "escort", item.Category)
			assert.Equal(t, int64(50000), item.BasePriceCents)
			assert.Equal(t, 1, item.ServiceHours)
			assert.Equal(t, 0.20, item.CommissionRate)
			assert.True(t, item.IsActive)
			return true
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "ESCORT_RANK_DIAMOND",
			Name:           "钻石段位护航",
			SubCategory:    model.SubCategorySolo,
			GameID:         &gameID,
			BasePriceCents: 50000,
			ServiceHours:   1,
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "ESCORT_RANK_DIAMOND", item.ItemCode)
	})

	t.Run("创建礼物成功", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Create", ctx, mock.MatchedBy(func(item *model.ServiceItem) bool {
			assert.Equal(t, "GIFT_ROSE", item.ItemCode)
			assert.Equal(t, model.SubCategoryGift, item.SubCategory)
			assert.Equal(t, 0, item.ServiceHours) // 礼物必须为0
			assert.Equal(t, 0.20, item.CommissionRate)
			return true
		})).Return(nil)

		req := CreateServiceItemRequest{
			ItemCode:       "GIFT_ROSE",
			Name:           "玫瑰花",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   0, // 礼物为0
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.True(t, item.IsGift())
	})

	t.Run("礼物service_hours必须为0", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		req := CreateServiceItemRequest{
			ItemCode:       "GIFT_INVALID",
			Name:           "无效礼物",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   1, // ❌ 礼物不能有服务时长
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, item)
		assert.Contains(t, err.Error(), "gift items must have service_hours = 0")
	})

	t.Run("游戏不存在", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		gameID := uint64(999)
		gameRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		req := CreateServiceItemRequest{
			ItemCode:       "TEST",
			Name:           "Test",
			SubCategory:    model.SubCategorySolo,
			GameID:         &gameID,
			BasePriceCents: 10000,
			ServiceHours:   1,
			CommissionRate: 0.20,
			MinUsers:       1,
			MaxPlayers:     1,
		}

		item, err := svc.CreateServiceItem(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, item)
	})
}

func TestServiceItemService_GetGiftList(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	gifts := []model.ServiceItem{
		{
			ID:             1,
			ItemCode:       "GIFT_ROSE",
			Name:           "玫瑰",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   0,
		},
		{
			ID:             2,
			ItemCode:       "GIFT_CHOCOLATE",
			Name:           "巧克力",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 5000,
			ServiceHours:   0,
		},
	}

	itemRepo.On("GetGifts", ctx, 1, 20).Return(gifts, int64(2), nil)

	// 获取礼物列表
	resp, err := svc.GetGiftList(ctx, 1, 20)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, int64(2), resp.Total)

	// 验证都是礼物
	for _, item := range resp.Items {
		assert.Equal(t, "gift", item.SubCategory)
	}
}

func TestServiceItemService_BatchOperations(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	t.Run("批量更新状态", func(t *testing.T) {
		ids := []uint64{1, 2, 3}
		req := BatchUpdateStatusRequest{
			IDs:      ids,
			IsActive: false,
		}

		itemRepo.On("BatchUpdateStatus", ctx, ids, false).Return(nil)

		err := svc.BatchUpdateStatus(ctx, req)

		assert.NoError(t, err)
		itemRepo.AssertExpectations(t)
	})

	t.Run("批量更新价格", func(t *testing.T) {
		ids := []uint64{1, 2, 3}
		req := BatchUpdatePriceRequest{
			IDs:            ids,
			BasePriceCents: 15000,
		}

		itemRepo.On("BatchUpdatePrice", ctx, ids, int64(15000)).Return(nil)

		err := svc.BatchUpdatePrice(ctx, req)

		assert.NoError(t, err)
		itemRepo.AssertExpectations(t)
	})

	t.Run("IDs为空应该报错", func(t *testing.T) {
		req := BatchUpdateStatusRequest{
			IDs:      []uint64{},
			IsActive: true,
		}

		err := svc.BatchUpdateStatus(ctx, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no item ids provided")
	})
}

func TestServiceItemService_GetServiceItem(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	t.Run("成功获取服务项目", func(t *testing.T) {
		item := &model.ServiceItem{
			ID:             1,
			ItemCode:       "TEST_ITEM",
			Name:           "测试服务",
			BasePriceCents: 10000,
			SubCategory:    model.SubCategorySolo,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(item, nil)

		result, err := svc.GetServiceItem(ctx, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "TEST_ITEM", result.ItemCode)
	})

	t.Run("服务项目不存在", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		gameRepo := new(MockGameRepo)
		playerRepo := new(MockPlayerRepo)

		svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

		itemRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		result, err := svc.GetServiceItem(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestServiceItemService_UpdateServiceItem(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	existingItem := &model.ServiceItem{
		ID:             1,
		ItemCode:       "TEST_ITEM",
		SubCategory:    model.SubCategorySolo,
		ServiceHours:   1,
		BasePriceCents: 10000,
	}

	itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*model.ServiceItem")).Return(nil)

	name := "Updated Name"
	price := int64(15000)
	req := UpdateServiceItemRequest{
		Name:           &name,
		BasePriceCents: &price,
	}

	err := svc.UpdateServiceItem(ctx, 1, req)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestServiceItemService_DeleteServiceItem(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	existingItem := &model.ServiceItem{
		ID:       1,
		ItemCode: "TEST_ITEM",
	}

	itemRepo.On("Get", ctx, uint64(1)).Return(existingItem, nil)
	itemRepo.On("Delete", ctx, uint64(1)).Return(nil)

	err := svc.DeleteServiceItem(ctx, 1)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestServiceItemService_ListServiceItems(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	items := []model.ServiceItem{
		{
			ID:             1,
			ItemCode:       "ITEM1",
			Name:           "服务1",
			SubCategory:    model.SubCategorySolo,
			BasePriceCents: 10000,
		},
		{
			ID:             2,
			ItemCode:       "ITEM2",
			Name:           "服务2",
			SubCategory:    model.SubCategoryTeam,
			BasePriceCents: 20000,
		},
	}

	opts := serviceitemrepo.ServiceItemListOptions{
		Page:     1,
		PageSize: 20,
	}

	itemRepo.On("List", ctx, opts).Return(items, int64(2), nil)

	req := ListServiceItemsRequest{
		Page:     1,
		PageSize: 20,
	}

	result, err := svc.ListServiceItems(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, int64(2), result.Total)
}

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

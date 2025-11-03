package serviceitem

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockServiceItemRepo) List(ctx context.Context, opts repository.ServiceItemListOptions) ([]model.ServiceItem, int64, error) {
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
func (m *MockGameRepo) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockGameRepo) List(ctx context.Context) ([]model.Game, error) { return nil, nil }
func (m *MockGameRepo) GetByKey(ctx context.Context, key string) (*model.Game, error) { return nil, nil }

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
func (m *MockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func TestServiceItemService_CreateServiceItem(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	gameRepo := new(MockGameRepo)
	playerRepo := new(MockPlayerRepo)

	svc := NewServiceItemService(itemRepo, gameRepo, playerRepo)

	t.Run("创建护航服务成功", func(t *testing.T) {
		game := &model.Game{
			ID:   1,
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



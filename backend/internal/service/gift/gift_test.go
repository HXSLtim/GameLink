package gift

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock Repositories
type MockServiceItemRepo struct {
	mock.Mock
}

func (m *MockServiceItemRepo) Create(ctx context.Context, item *model.ServiceItem) error {
	args := m.Called(ctx, item)
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

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	if args.Get(0) != nil {
		// 模拟自动生成ID
		order.ID = 1001
	}
	return args.Error(0)
}

func (m *MockOrderRepo) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepo) List(ctx context.Context, opts repository.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
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

func (m *MockPlayerRepo) Create(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepo) Update(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockPlayerRepo) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

type MockCommissionRepo struct {
	mock.Mock
}

func (m *MockCommissionRepo) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockCommissionRepo) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommissionRecord), args.Error(1)
}

// Mock所有其他必需的方法
func (m *MockCommissionRepo) CreateRule(ctx context.Context, rule *model.CommissionRule) error { return nil }
func (m *MockCommissionRepo) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) { return nil, nil }
func (m *MockCommissionRepo) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) { return nil, nil }
func (m *MockCommissionRepo) GetRuleForOrder(ctx context.Context, gameID, playerID *uint64, serviceType *string) (*model.CommissionRule, error) { return nil, nil }
func (m *MockCommissionRepo) ListRules(ctx context.Context, opts repository.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) { return nil, 0, nil }
func (m *MockCommissionRepo) UpdateRule(ctx context.Context, rule *model.CommissionRule) error { return nil }
func (m *MockCommissionRepo) DeleteRule(ctx context.Context, id uint64) error { return nil }
func (m *MockCommissionRepo) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) { return nil, nil }
func (m *MockCommissionRepo) ListRecords(ctx context.Context, opts repository.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) { return nil, 0, nil }
func (m *MockCommissionRepo) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error { return nil }
func (m *MockCommissionRepo) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error { return nil }
func (m *MockCommissionRepo) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) { return nil, nil }
func (m *MockCommissionRepo) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) { return nil, nil }
func (m *MockCommissionRepo) ListSettlements(ctx context.Context, opts repository.SettlementListOptions) ([]model.MonthlySettlement, int64, error) { return nil, 0, nil }
func (m *MockCommissionRepo) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error { return nil }
func (m *MockCommissionRepo) GetMonthlyStats(ctx context.Context, month string) (*repository.MonthlyStats, error) { return nil, nil }
func (m *MockCommissionRepo) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) { return 0, nil }

func TestGiftService_SendGift(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	commissionRepo := new(MockCommissionRepo)

	svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

	t.Run("成功赠送礼物", func(t *testing.T) {
		// Mock礼物项目
		giftItem := &model.ServiceItem{
			ID:             1,
			ItemCode:       "GIFT_ROSE",
			Name:           "玫瑰花",
			SubCategory:    model.SubCategoryGift,
			BasePriceCents: 10000,
			ServiceHours:   0,
			CommissionRate: 0.20,
			IsActive:       true,
		}

		// Mock陪玩师
		player := &model.Player{
			ID:       5,
			Nickname: "测试陪玩师",
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(giftItem, nil)
		playerRepo.On("Get", ctx, uint64(5)).Return(player, nil)
		orderRepo.On("Create", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
		orderRepo.On("Update", ctx, mock.AnythingOfType("*model.Order")).Return(nil)
		commissionRepo.On("CreateRecord", ctx, mock.AnythingOfType("*model.CommissionRecord")).Return(nil)

		// 赠送礼物
		req := SendGiftRequest{
			PlayerID:    5,
			GiftItemID:  1,
			Quantity:    3,
			Message:     "感谢你！",
			IsAnonymous: false,
		}

		resp, err := svc.SendGift(ctx, 100, req)

		// 验证
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint64(5), resp.PlayerID)
		assert.Equal(t, "玫瑰花", resp.GiftName)
		assert.Equal(t, 3, resp.Quantity)
		assert.Equal(t, int64(30000), resp.TotalPrice) // 10000 * 3
		assert.NotNil(t, resp.DeliveredAt)

		// 验证订单创建的参数
		orderRepo.AssertCalled(t, "Create", ctx, mock.MatchedBy(func(order *model.Order) bool {
			assert.Equal(t, uint64(1), order.ItemID)
			assert.Equal(t, 3, order.Quantity)
			assert.Equal(t, int64(10000), order.UnitPriceCents)
			assert.Equal(t, int64(30000), order.TotalPriceCents)
			assert.Equal(t, int64(6000), order.CommissionCents)    // 20%
			assert.Equal(t, int64(24000), order.PlayerIncomeCents) // 80%
			assert.Equal(t, "感谢你！", order.GiftMessage)
			assert.False(t, order.IsAnonymous)
			assert.NotNil(t, order.RecipientPlayerID)
			assert.Equal(t, uint64(5), *order.RecipientPlayerID)
			return true
		}))
	})

	t.Run("礼物项目不存在", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		commissionRepo := new(MockCommissionRepo)

		svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

		itemRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		req := SendGiftRequest{
			PlayerID:   5,
			GiftItemID: 999,
			Quantity:   1,
		}

		resp, err := svc.SendGift(ctx, 100, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("不是礼物类型", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		commissionRepo := new(MockCommissionRepo)

		svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

		// Mock护航服务（不是礼物）
		escortItem := &model.ServiceItem{
			ID:          1,
			SubCategory: model.SubCategorySolo,
			IsActive:    true,
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(escortItem, nil)

		req := SendGiftRequest{
			PlayerID:   5,
			GiftItemID: 1,
			Quantity:   1,
		}

		resp, err := svc.SendGift(ctx, 100, req)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidGiftItem, err)
		assert.Nil(t, resp)
	})

	t.Run("礼物未激活", func(t *testing.T) {
		itemRepo := new(MockServiceItemRepo)
		orderRepo := new(MockOrderRepo)
		playerRepo := new(MockPlayerRepo)
		commissionRepo := new(MockCommissionRepo)

		svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

		// Mock未激活的礼物
		giftItem := &model.ServiceItem{
			ID:          1,
			SubCategory: model.SubCategoryGift,
			IsActive:    false, // 未激活
		}

		itemRepo.On("Get", ctx, uint64(1)).Return(giftItem, nil)

		req := SendGiftRequest{
			PlayerID:   5,
			GiftItemID: 1,
			Quantity:   1,
		}

		resp, err := svc.SendGift(ctx, 100, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not active")
		assert.Nil(t, resp)
	})
}

func TestOrder_IsGiftOrder(t *testing.T) {
	t.Run("是礼物订单", func(t *testing.T) {
		recipientID := uint64(5)
		order := &model.Order{
			RecipientPlayerID: &recipientID,
		}
		assert.True(t, order.IsGiftOrder())
	})

	t.Run("不是礼物订单", func(t *testing.T) {
		order := &model.Order{
			RecipientPlayerID: nil,
		}
		assert.False(t, order.IsGiftOrder())
	})

	t.Run("RecipientPlayerID为0", func(t *testing.T) {
		recipientID := uint64(0)
		order := &model.Order{
			RecipientPlayerID: &recipientID,
		}
		assert.False(t, order.IsGiftOrder())
	})
}

func TestGiftService_GetPlayerReceivedGifts(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	commissionRepo := new(MockCommissionRepo)

	svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

	// Mock订单数据：包含礼物订单和护航订单
	recipientID := uint64(5)
	playerID := uint64(5)
	deliveredAt := time.Now()

	giftOrder := model.Order{
		ID:                1001,
		OrderNo:           "GIFT001",
		ItemID:            1,
		RecipientPlayerID: &recipientID, // 礼物订单
		Quantity:          2,
		TotalPriceCents:   20000,
		PlayerIncomeCents: 16000,
		GiftMessage:       "谢谢！",
		IsAnonymous:       false,
		DeliveredAt:       &deliveredAt,
		CreatedAt:         time.Now(),
	}

	escortOrder := model.Order{
		ID:                1002,
		OrderNo:           "ESC001",
		ItemID:            2,
		PlayerID:          &playerID, // 护航订单（不是礼物）
		RecipientPlayerID: nil,
		TotalPriceCents:   50000,
	}

	orders := []model.Order{giftOrder, escortOrder}

	orderRepo.On("List", ctx, mock.AnythingOfType("repository.OrderListOptions")).
		Return(orders, int64(2), nil)

	// Mock礼物项目
	giftItem := &model.ServiceItem{
		ID:      1,
		Name:    "玫瑰花",
		IconURL: "/rose.png",
	}
	itemRepo.On("Get", ctx, uint64(1)).Return(giftItem, nil)

	// 获取收到的礼物
	resp, err := svc.GetPlayerReceivedGifts(ctx, 5, 1, 10)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Gifts, 1) // 只有1个礼物订单
	assert.Equal(t, uint64(1001), resp.Gifts[0].OrderID)
	assert.Equal(t, "玫瑰花", resp.Gifts[0].GiftName)
	assert.Equal(t, 2, resp.Gifts[0].Quantity)
	assert.Equal(t, int64(20000), resp.Gifts[0].TotalPrice)
	assert.Equal(t, int64(16000), resp.Gifts[0].Income)
	assert.Equal(t, "谢谢！", resp.Gifts[0].Message)
}

func TestGiftService_GetGiftStats(t *testing.T) {
	ctx := context.Background()

	itemRepo := new(MockServiceItemRepo)
	orderRepo := new(MockOrderRepo)
	playerRepo := new(MockPlayerRepo)
	commissionRepo := new(MockCommissionRepo)

	svc := NewGiftService(itemRepo, orderRepo, playerRepo, commissionRepo)

	// Mock数据：3个礼物订单
	recipientID := uint64(5)
	orders := []model.Order{
		{
			ID:                1001,
			ItemID:            1,
			RecipientPlayerID: &recipientID,
			Quantity:          2,
			PlayerIncomeCents: 16000,
			Status:            model.OrderStatusCompleted,
		},
		{
			ID:                1002,
			ItemID:            1,
			RecipientPlayerID: &recipientID,
			Quantity:          1,
			PlayerIncomeCents: 8000,
			Status:            model.OrderStatusCompleted,
		},
		{
			ID:                1003,
			ItemID:            2,
			RecipientPlayerID: &recipientID,
			Quantity:          5,
			PlayerIncomeCents: 40000,
			Status:            model.OrderStatusCompleted,
		},
	}

	orderRepo.On("List", ctx, mock.AnythingOfType("repository.OrderListOptions")).
		Return(orders, int64(3), nil)

	// 获取礼物统计
	resp, err := svc.GetGiftStats(ctx, 5)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(8), resp.TotalGiftsReceived)      // 2+1+5
	assert.Equal(t, int64(64000), resp.TotalGiftIncome)     // 16000+8000+40000
	assert.Equal(t, int64(3), resp.TotalGiftOrders)
}


package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/ordergroup"
)

// MockOrderGroupRepository 主订单仓储 Mock (for creation tests)
type mockOrderGroupRepo struct {
	createGroup       func(ctx context.Context, group *model.OrderGroup) error
	getGroup          func(ctx context.Context, id uint64) (*model.OrderGroup, error)
	getByGroupNo      func(ctx context.Context, groupNo string) (*model.OrderGroup, error)
	getWithSubOrders  func(ctx context.Context, id uint64) (*model.OrderGroup, error)
	updateGroup       func(ctx context.Context, group *model.OrderGroup) error
	updateStatus      func(ctx context.Context, id uint64, status model.OrderGroupStatus) error
}

func (m *mockOrderGroupRepo) Create(ctx context.Context, group *model.OrderGroup) error {
	if m.createGroup != nil {
		return m.createGroup(ctx, group)
	}
	return nil
}

func (m *mockOrderGroupRepo) Get(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	if m.getGroup != nil {
		return m.getGroup(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderGroupRepo) GetByGroupNo(ctx context.Context, groupNo string) (*model.OrderGroup, error) {
	if m.getByGroupNo != nil {
		return m.getByGroupNo(ctx, groupNo)
	}
	return nil, nil
}

func (m *mockOrderGroupRepo) GetWithSubOrders(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	if m.getWithSubOrders != nil {
		return m.getWithSubOrders(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderGroupRepo) Update(ctx context.Context, group *model.OrderGroup) error {
	if m.updateGroup != nil {
		return m.updateGroup(ctx, group)
	}
	return nil
}

func (m *mockOrderGroupRepo) UpdateStatus(ctx context.Context, id uint64, status model.OrderGroupStatus) error {
	if m.updateStatus != nil {
		return m.updateStatus(ctx, id, status)
	}
	return nil
}

func (m *mockOrderGroupRepo) List(ctx context.Context, opts ordergroup.ListOptions) ([]model.OrderGroup, int64, error) {
	return nil, 0, nil
}

func (m *mockOrderGroupRepo) ListByUser(ctx context.Context, userID uint64, opts ordergroup.ListOptions) ([]model.OrderGroup, int64, error) {
	return nil, 0, nil
}

// TestBuildOrderGroupWithSubOrders_Success 测试订单拆分构建
func TestBuildOrderGroupWithSubOrders_Success(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	serviceID := uint64(10)
	scheduledStart := time.Now().Add(time.Hour)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		ServiceID:      &serviceID,
		Title:          "3小时陪玩",
		Description:    "王者荣耀陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	}

	hourlyPrice := int64(5000)
	commissionPerHour := int64(1000)
	playerIncomePerHour := int64(4000)

	group, subOrders := service.buildOrderGroupWithSubOrders(userID, req, hourlyPrice, commissionPerHour, playerIncomePerHour)

	// 验证主订单
	assert.NotNil(t, group)
	assert.Equal(t, userID, group.UserID)
	assert.Equal(t, gameID, group.GameID)
	assert.Equal(t, serviceID, group.ItemID)
	assert.Equal(t, playerID, group.OriginalPlayer)
	assert.Equal(t, int64(15000), group.TotalPriceCents) // 5000 * 3
	assert.Equal(t, 3, group.TotalHours)
	assert.Equal(t, 0, group.CompletedHours)
	assert.Equal(t, model.OrderGroupStatusPending, group.Status)
	assert.Equal(t, "3小时陪玩", group.Title)
	assert.NotEmpty(t, group.GroupNo)
	assert.True(t, len(group.GroupNo) > 0 && group.GroupNo[0] == 'G')

	// 验证子订单
	assert.Len(t, subOrders, 3)

	for i, sub := range subOrders {
		assert.Equal(t, userID, sub.UserID)
		assert.Equal(t, serviceID, sub.ItemID)
		assert.Equal(t, playerID, *sub.PlayerID)
		assert.Equal(t, gameID, *sub.GameID)
		assert.Equal(t, hourlyPrice, sub.TotalPriceCents)
		assert.Equal(t, commissionPerHour, sub.CommissionCents)
		assert.Equal(t, playerIncomePerHour, sub.PlayerIncomeCents)
		assert.Equal(t, i+1, sub.HourIndex)
		assert.True(t, sub.IsSubOrder)
		assert.True(t, sub.CanTransfer)
		assert.Equal(t, model.OrderStatusPending, sub.Status)
		assert.NotEmpty(t, sub.OrderNo)

		// 验证时间
		expectedStart := scheduledStart.Add(time.Duration(i) * time.Hour)
		expectedEnd := expectedStart.Add(time.Hour)
		assert.Equal(t, expectedStart.Unix(), sub.ScheduledStart.Unix())
		assert.Equal(t, expectedEnd.Unix(), sub.ScheduledEnd.Unix())
	}
}

// TestBuildOrderGroupWithSubOrders_SingleHour 测试单小时订单拆分
func TestBuildOrderGroupWithSubOrders_SingleHour(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "1小时陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	group, subOrders := service.buildOrderGroupWithSubOrders(userID, req, 5000, 1000, 4000)

	assert.NotNil(t, group)
	assert.Equal(t, 1, group.TotalHours)
	assert.Len(t, subOrders, 1)
	assert.Equal(t, 1, subOrders[0].HourIndex)
}

// TestBuildOrderGroupWithSubOrders_FractionalHours 测试小数时长向上取整
func TestBuildOrderGroupWithSubOrders_FractionalHours(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "2.5小时陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  2.5, // 向上取整为 3
	}

	group, subOrders := service.buildOrderGroupWithSubOrders(userID, req, 5000, 1000, 4000)

	assert.NotNil(t, group)
	assert.Equal(t, 3, group.TotalHours)
	assert.Len(t, subOrders, 3)
}

// TestBuildOrderGroupWithSubOrders_ZeroHours 测试零时长默认为1
func TestBuildOrderGroupWithSubOrders_ZeroHours(t *testing.T) {
	service := &OrderService{}

	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "测试",
		ScheduledStart: &scheduledStart,
		DurationHours:  0.3, // 小于1，向上取整为1
	}

	group, subOrders := service.buildOrderGroupWithSubOrders(userID, req, 5000, 1000, 4000)

	assert.NotNil(t, group)
	assert.Equal(t, 1, group.TotalHours)
	assert.Len(t, subOrders, 1)
}

// TestCreateOrderWithSplit_Success 测试拆分订单创建
func TestCreateOrderWithSplit_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	var createdGroup *model.OrderGroup
	var createdSubOrders []*model.Order

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = uint64(len(createdSubOrders) + 10)
			createdSubOrders = append(createdSubOrders, order)
			return nil
		},
	}

	orderGroups := &mockOrderGroupRepo{
		createGroup: func(ctx context.Context, group *model.OrderGroup) error {
			group.ID = 1
			createdGroup = group
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, games, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "3小时陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  3,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsSplit)
	assert.Equal(t, 3, resp.SubOrderCount)
	assert.Equal(t, 3, resp.TotalHours)
	assert.Equal(t, uint64(1), resp.OrderID)
	assert.NotEmpty(t, resp.GroupNo)

	// 验证创建的主订单
	assert.NotNil(t, createdGroup)
	assert.Equal(t, 3, createdGroup.TotalHours)

	// 验证创建的子订单
	assert.Len(t, createdSubOrders, 3)
	for i, sub := range createdSubOrders {
		assert.Equal(t, uint64(1), *sub.GroupID)
		assert.Equal(t, i+1, sub.HourIndex)
		assert.True(t, sub.IsSubOrder)
	}
}

// TestCreateOrder_NoSplitForShortDuration 测试短时长不拆分
func TestCreateOrder_NoSplitForShortDuration(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	var createdOrder *model.Order

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
			createdOrder = order
			return nil
		},
	}

	orderGroups := &mockOrderGroupRepo{}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, games, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "1小时陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  1, // 1小时不拆分
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsSplit)
	assert.Equal(t, 0, resp.SubOrderCount)
	assert.Equal(t, uint64(123), resp.OrderID)

	// 验证创建的是普通订单
	assert.NotNil(t, createdOrder)
	assert.False(t, createdOrder.IsSubOrder)
	assert.Nil(t, createdOrder.GroupID)
}

// TestCreateOrder_NoSplitWithoutOrderGroupRepo 测试没有 orderGroups 仓储时不拆分
func TestCreateOrder_NoSplitWithoutOrderGroupRepo(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	var createdOrder *model.Order

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
			createdOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	// 不设置 orderGroups 仓储
	service := NewOrderService(orders, players, &MockUserRepository{}, games, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "3小时陪玩",
		ScheduledStart: &scheduledStart,
		DurationHours:  3, // 即使3小时，没有仓储也不拆分
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsSplit)

	// 验证创建的是普通订单
	assert.NotNil(t, createdOrder)
	assert.False(t, createdOrder.IsSubOrder)
}

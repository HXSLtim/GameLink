package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/ordergroup"
)

// MockOrderGroupRepository 主订单仓储 Mock
type MockOrderGroupRepository struct {
	createGroup       func(ctx context.Context, group *model.OrderGroup) error
	getGroup          func(ctx context.Context, id uint64) (*model.OrderGroup, error)
	getByGroupNo      func(ctx context.Context, groupNo string) (*model.OrderGroup, error)
	getWithSubOrders  func(ctx context.Context, id uint64) (*model.OrderGroup, error)
	updateGroup       func(ctx context.Context, group *model.OrderGroup) error
	updateStatus      func(ctx context.Context, id uint64, status model.OrderGroupStatus) error
}

func (m *MockOrderGroupRepository) Create(ctx context.Context, group *model.OrderGroup) error {
	if m.createGroup != nil {
		return m.createGroup(ctx, group)
	}
	return nil
}

func (m *MockOrderGroupRepository) Get(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	if m.getGroup != nil {
		return m.getGroup(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderGroupRepository) GetByGroupNo(ctx context.Context, groupNo string) (*model.OrderGroup, error) {
	if m.getByGroupNo != nil {
		return m.getByGroupNo(ctx, groupNo)
	}
	return nil, nil
}

func (m *MockOrderGroupRepository) GetWithSubOrders(ctx context.Context, id uint64) (*model.OrderGroup, error) {
	if m.getWithSubOrders != nil {
		return m.getWithSubOrders(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderGroupRepository) Update(ctx context.Context, group *model.OrderGroup) error {
	if m.updateGroup != nil {
		return m.updateGroup(ctx, group)
	}
	return nil
}

func (m *MockOrderGroupRepository) UpdateStatus(ctx context.Context, id uint64, status model.OrderGroupStatus) error {
	if m.updateStatus != nil {
		return m.updateStatus(ctx, id, status)
	}
	return nil
}

func (m *MockOrderGroupRepository) List(ctx context.Context, opts ordergroup.ListOptions) ([]model.OrderGroup, int64, error) {
	return nil, 0, nil
}

func (m *MockOrderGroupRepository) ListByUser(ctx context.Context, userID uint64, opts ordergroup.ListOptions) ([]model.OrderGroup, int64, error) {
	return nil, 0, nil
}

// createTestSubOrder 创建测试子订单
func createTestSubOrder(id uint64, userID uint64, groupID uint64, hourIndex int, status model.OrderStatus) *model.Order {
	playerID := uint64(100)
	gameID := uint64(1)
	now := time.Now()
	hourStart := now.Add(time.Duration(hourIndex-1) * time.Hour)
	hourEnd := hourStart.Add(time.Hour)

	return &model.Order{
		Base:              model.Base{ID: id},
		OrderNo:           model.GenerateEscortOrderNo(),
		UserID:            userID,
		ItemID:            1,
		PlayerID:          &playerID,
		GameID:            &gameID,
		GroupID:           &groupID,
		Quantity:          1,
		UnitPriceCents:    5000,
		TotalPriceCents:   5000,
		CommissionCents:   1000,
		PlayerIncomeCents: 4000,
		Currency:          model.CurrencyCNY,
		Status:            status,
		Title:             "Test Sub Order",
		ScheduledStart:    &hourStart,
		ScheduledEnd:      &hourEnd,
		OrderConfig:       "{}",
		HourIndex:         hourIndex,
		IsSubOrder:        true,
		CanTransfer:       true,
	}
}

// TestTransferSubOrder_Success 测试成功转单
func TestTransferSubOrder_Success(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusPending)

	var createdOrder *model.Order
	var updatedOrder *model.Order

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == subOrderID {
				return subOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			createdOrder = order
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				UserID:             id + 100,
				Nickname:           "NewPlayer",
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{
				Base:       model.Base{ID: groupID},
				TotalHours: 3,
				SubOrders:  []model.Order{*subOrder},
			}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error {
			return nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := TransferSubOrderRequest{
		SubOrderID:   subOrderID,
		NewPlayerID:  newPlayerID,
		TransferNote: "陪玩A有事",
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, uint64(100), resp.NewSubOrderID)

	// 验证新订单
	assert.NotNil(t, createdOrder)
	assert.Equal(t, newPlayerID, *createdOrder.PlayerID)
	assert.Equal(t, subOrder.HourIndex, createdOrder.HourIndex)
	assert.True(t, createdOrder.IsSubOrder)
	assert.Equal(t, &subOrderID, createdOrder.TransferFrom)

	// 验证原订单更新
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusCanceled, updatedOrder.Status)
	assert.False(t, updatedOrder.CanTransfer)
}

// TestTransferSubOrder_NotSubOrder 测试非子订单不能转单
func TestTransferSubOrder_NotSubOrder(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	orderID := uint64(10)

	normalOrder := createTestOrder(orderID, 1, model.OrderStatusPending)
	normalOrder.IsSubOrder = false

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return normalOrder, nil
		},
	}

	service := NewOrderService(orders, &MockPlayerRepository{}, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := TransferSubOrderRequest{
		SubOrderID:  orderID,
		NewPlayerID: 200,
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "只有子订单可以转单")
}

// TestTransferSubOrder_CannotTransfer 测试不可转单的订单
func TestTransferSubOrder_CannotTransfer(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusPending)
	subOrder.CanTransfer = false

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return subOrder, nil
		},
	}

	service := NewOrderService(orders, &MockPlayerRepository{}, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := TransferSubOrderRequest{
		SubOrderID:  subOrderID,
		NewPlayerID: 200,
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "该订单不可转单")
}

// TestTransferSubOrder_CompletedOrder 测试已完成订单不能转单
func TestTransferSubOrder_CompletedOrder(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusCompleted)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return subOrder, nil
		},
	}

	service := NewOrderService(orders, &MockPlayerRepository{}, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := TransferSubOrderRequest{
		SubOrderID:  subOrderID,
		NewPlayerID: 200,
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "已完成的订单不能转单")
}

// TestTransferSubOrder_SamePlayer 测试不能转给同一个陪玩师
func TestTransferSubOrder_SamePlayer(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	playerID := uint64(100)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusPending)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return subOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := TransferSubOrderRequest{
		SubOrderID:  subOrderID,
		NewPlayerID: playerID, // 同一个陪玩师
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "不能转给同一个陪玩师")
}

// TestTransferSubOrder_UnverifiedPlayer 测试新陪玩师未认证
func TestTransferSubOrder_UnverifiedPlayer(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusPending)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return subOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationPending, // 未认证
			}, nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})

	req := TransferSubOrderRequest{
		SubOrderID:  subOrderID,
		NewPlayerID: newPlayerID,
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "新陪玩师未通过认证")
}

// TestBatchTransferSubOrders_Success 测试批量转单成功
func TestBatchTransferSubOrders_Success(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder1 := createTestSubOrder(10, 1, groupID, 2, model.OrderStatusPending)
	subOrder2 := createTestSubOrder(11, 1, groupID, 3, model.OrderStatusPending)

	orderMap := map[uint64]*model.Order{
		10: subOrder1,
		11: subOrder2,
	}

	newOrderID := uint64(100)
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if order, ok := orderMap[id]; ok {
				return order, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = newOrderID
			newOrderID++
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{
				Base:       model.Base{ID: groupID},
				TotalHours: 3,
				SubOrders:  []model.Order{*subOrder1, *subOrder2},
			}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error {
			return nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := BatchTransferRequest{
		SubOrderIDs:  []uint64{10, 11},
		NewPlayerID:  newPlayerID,
		TransferNote: "批量转单测试",
	}

	resp, err := service.BatchTransferSubOrders(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.SuccessCount)
	assert.Equal(t, 0, resp.FailedCount)
	assert.Len(t, resp.NewOrderIDs, 2)
}

// TestBatchTransferSubOrders_PartialFailure 测试批量转单部分失败
func TestBatchTransferSubOrders_PartialFailure(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder1 := createTestSubOrder(10, 1, groupID, 2, model.OrderStatusPending)
	subOrder2 := createTestSubOrder(11, 1, groupID, 3, model.OrderStatusCompleted) // 已完成，不能转

	orderMap := map[uint64]*model.Order{
		10: subOrder1,
		11: subOrder2,
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if order, ok := orderMap[id]; ok {
				return order, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{
				Base:       model.Base{ID: groupID},
				TotalHours: 3,
				SubOrders:  []model.Order{*subOrder1, *subOrder2},
			}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error {
			return nil
		},
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := BatchTransferRequest{
		SubOrderIDs:  []uint64{10, 11},
		NewPlayerID:  newPlayerID,
		TransferNote: "批量转单测试",
	}

	resp, err := service.BatchTransferSubOrders(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.SuccessCount)
	assert.Equal(t, 1, resp.FailedCount)
	assert.Len(t, resp.NewOrderIDs, 1)
	assert.Len(t, resp.Errors, 1)
}

// TestGetTransferableSubOrders_Success 测试获取可转单子订单
func TestGetTransferableSubOrders_Success(t *testing.T) {
	ctx := context.Background()
	groupID := uint64(1)

	subOrder1 := createTestSubOrder(10, 1, groupID, 1, model.OrderStatusCompleted)
	subOrder1.CanTransfer = false
	subOrder2 := createTestSubOrder(11, 1, groupID, 2, model.OrderStatusPending)
	subOrder3 := createTestSubOrder(12, 1, groupID, 3, model.OrderStatusPending)

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{
				Base:       model.Base{ID: groupID},
				TotalHours: 3,
				SubOrders:  []model.Order{*subOrder1, *subOrder2, *subOrder3},
			}, nil
		},
	}

	service := NewOrderService(&MockOrderRepository{}, &MockPlayerRepository{}, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	transferable, err := service.GetTransferableSubOrders(ctx, groupID)

	require.NoError(t, err)
	assert.Len(t, transferable, 2) // 只有 subOrder2 和 subOrder3 可转
}


// TestTransferSubOrder_IncomeAttribution_NoServiceStarted 测试转单收入归属 - 未开始服务
func TestTransferSubOrder_IncomeAttribution_NoServiceStarted(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	// 原订单：总价5000分，抽成1000分，陪玩师收入4000分
	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusPending)

	var createdOrder *model.Order
	var updatedOrder *model.Order

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == subOrderID {
				return subOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			createdOrder = order
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{Base: model.Base{ID: groupID}, SubOrders: []model.Order{*subOrder}}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error { return nil },
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := TransferSubOrderRequest{
		SubOrderID:       subOrderID,
		NewPlayerID:      newPlayerID,
		TransferNote:     "未开始服务转单",
		CompletedMinutes: 0, // 未开始服务
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证收入归属：原陪玩师0，新陪玩师全部
	assert.Equal(t, int64(0), resp.OriginalPlayerIncome)
	assert.Equal(t, int64(4000), resp.NewPlayerIncome)

	// 验证新订单：抽成为0（不重复计算），收入为全部
	assert.Equal(t, int64(0), createdOrder.CommissionCents)
	assert.Equal(t, int64(4000), createdOrder.PlayerIncomeCents)

	// 验证原订单：收入更新为0
	assert.Equal(t, int64(0), updatedOrder.PlayerIncomeCents)
}

// TestTransferSubOrder_IncomeAttribution_HalfCompleted 测试转单收入归属 - 完成一半
func TestTransferSubOrder_IncomeAttribution_HalfCompleted(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusInProgress)

	var createdOrder *model.Order
	var updatedOrder *model.Order

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == subOrderID {
				return subOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			createdOrder = order
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{Base: model.Base{ID: groupID}, SubOrders: []model.Order{*subOrder}}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error { return nil },
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := TransferSubOrderRequest{
		SubOrderID:       subOrderID,
		NewPlayerID:      newPlayerID,
		TransferNote:     "完成一半转单",
		CompletedMinutes: 30, // 完成30分钟（一半）
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	// 验证收入归属：按比例分配 4000 * 30/60 = 2000
	assert.Equal(t, int64(2000), resp.OriginalPlayerIncome)
	assert.Equal(t, int64(2000), resp.NewPlayerIncome)

	// 验证新订单
	assert.Equal(t, int64(0), createdOrder.CommissionCents)
	assert.Equal(t, int64(2000), createdOrder.PlayerIncomeCents)

	// 验证原订单
	assert.Equal(t, int64(2000), updatedOrder.PlayerIncomeCents)
}

// TestTransferSubOrder_IncomeAttribution_MostlyCompleted 测试转单收入归属 - 完成大部分
func TestTransferSubOrder_IncomeAttribution_MostlyCompleted(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusInProgress)

	var createdOrder *model.Order
	var updatedOrder *model.Order

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == subOrderID {
				return subOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			createdOrder = order
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{Base: model.Base{ID: groupID}, SubOrders: []model.Order{*subOrder}}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error { return nil },
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := TransferSubOrderRequest{
		SubOrderID:       subOrderID,
		NewPlayerID:      newPlayerID,
		TransferNote:     "完成45分钟转单",
		CompletedMinutes: 45, // 完成45分钟
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)

	// 验证收入归属：4000 * 45/60 = 3000, 4000 - 3000 = 1000
	assert.Equal(t, int64(3000), resp.OriginalPlayerIncome)
	assert.Equal(t, int64(1000), resp.NewPlayerIncome)

	// 验证新订单
	assert.Equal(t, int64(0), createdOrder.CommissionCents)
	assert.Equal(t, int64(1000), createdOrder.PlayerIncomeCents)

	// 验证原订单
	assert.Equal(t, int64(3000), updatedOrder.PlayerIncomeCents)
}

// TestTransferSubOrder_IncomeAttribution_InvalidMinutes 测试转单收入归属 - 无效分钟数
func TestTransferSubOrder_IncomeAttribution_InvalidMinutes(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	subOrderID := uint64(10)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	subOrder := createTestSubOrder(subOrderID, 1, groupID, 2, model.OrderStatusInProgress)

	var createdOrder *model.Order

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if id == subOrderID {
				return subOrder, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 100
			createdOrder = order
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error { return nil },
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{Base: model.Base{ID: groupID}, SubOrders: []model.Order{*subOrder}}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error { return nil },
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	// 测试负数分钟数 - 应该被修正为0
	req := TransferSubOrderRequest{
		SubOrderID:       subOrderID,
		NewPlayerID:      newPlayerID,
		CompletedMinutes: -10,
	}

	resp, err := service.TransferSubOrder(ctx, operatorID, req)
	require.NoError(t, err)
	assert.Equal(t, int64(0), resp.OriginalPlayerIncome)
	assert.Equal(t, int64(4000), resp.NewPlayerIncome)

	// 测试超过60分钟 - 应该被修正为60
	subOrder2 := createTestSubOrder(subOrderID+1, 1, groupID, 3, model.OrderStatusInProgress)
	orders.getOrder = func(ctx context.Context, id uint64) (*model.Order, error) {
		return subOrder2, nil
	}

	req2 := TransferSubOrderRequest{
		SubOrderID:       subOrderID + 1,
		NewPlayerID:      newPlayerID,
		CompletedMinutes: 100, // 超过60分钟
	}

	resp2, err := service.TransferSubOrder(ctx, operatorID, req2)
	require.NoError(t, err)
	// 修正为60分钟，原陪玩师获得全部收入
	assert.Equal(t, int64(4000), resp2.OriginalPlayerIncome)
	assert.Equal(t, int64(0), resp2.NewPlayerIncome)
	assert.Equal(t, int64(0), createdOrder.PlayerIncomeCents)
}

// TestBatchTransferSubOrders_IncomeAttribution 测试批量转单收入归属
func TestBatchTransferSubOrders_IncomeAttribution(t *testing.T) {
	ctx := context.Background()
	operatorID := uint64(1)
	groupID := uint64(1)
	newPlayerID := uint64(200)

	// 第一个订单完成20分钟，第二个订单未开始
	subOrder1 := createTestSubOrder(10, 1, groupID, 2, model.OrderStatusInProgress)
	subOrder2 := createTestSubOrder(11, 1, groupID, 3, model.OrderStatusPending)

	orderMap := map[uint64]*model.Order{
		10: subOrder1,
		11: subOrder2,
	}

	newOrderID := uint64(100)
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			if order, ok := orderMap[id]; ok {
				return order, nil
			}
			return nil, repository.ErrNotFound
		},
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = newOrderID
			newOrderID++
			return nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error { return nil },
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return &model.Player{
				Base:               model.Base{ID: id},
				VerificationStatus: model.VerificationVerified,
			}, nil
		},
	}

	orderGroups := &MockOrderGroupRepository{
		getWithSubOrders: func(ctx context.Context, id uint64) (*model.OrderGroup, error) {
			return &model.OrderGroup{
				Base:       model.Base{ID: groupID},
				SubOrders:  []model.Order{*subOrder1, *subOrder2},
			}, nil
		},
		updateGroup: func(ctx context.Context, group *model.OrderGroup) error { return nil },
	}

	service := NewOrderService(orders, players, &MockUserRepository{}, &MockGameRepository{}, &MockPaymentRepository{}, &MockReviewRepository{}, &MockCommissionRepository{})
	service.SetOrderGroupRepository(orderGroups)

	req := BatchTransferRequest{
		SubOrderIDs:      []uint64{10, 11},
		NewPlayerID:      newPlayerID,
		TransferNote:     "批量转单",
		CompletedMinutes: 20, // 第一个订单完成20分钟
	}

	resp, err := service.BatchTransferSubOrders(ctx, operatorID, req)

	require.NoError(t, err)
	assert.Equal(t, 2, resp.SuccessCount)

	// 第一个订单：4000 * 20/60 ≈ 1333 原陪玩师，2667 新陪玩师
	// 第二个订单：0 原陪玩师，4000 新陪玩师
	// 总计：原陪玩师 1333，新陪玩师 6667
	assert.Equal(t, int64(1333), resp.TotalOriginalIncome)
	assert.Equal(t, int64(6667), resp.TotalNewPlayerIncome)
}

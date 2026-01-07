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

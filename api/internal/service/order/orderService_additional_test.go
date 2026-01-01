package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/cache"
)

// MockDistributedLock for testing distributed lock functionality
type MockDistributedLock struct {
	lock    func(ctx context.Context, key string, ttl time.Duration) (bool, error)
	tryLock func(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error)
	unlock  func(ctx context.Context, key string) error
}

func (m *MockDistributedLock) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if m.lock != nil {
		return m.lock(ctx, key, ttl)
	}
	return true, nil
}

func (m *MockDistributedLock) TryLock(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error) {
	if m.tryLock != nil {
		return m.tryLock(ctx, key, ttl, retry, interval)
	}
	return true, nil
}

func (m *MockDistributedLock) Unlock(ctx context.Context, key string) error {
	if m.unlock != nil {
		return m.unlock(ctx, key)
	}
	return nil
}

// MockChatGroupRepository for testing chat deactivation
type MockChatGroupRepository struct {
	getByRelatedOrderID func(ctx context.Context, orderID uint64) (*model.ChatGroup, error)
	deactivate          func(ctx context.Context, groupID uint64) error
}

func (m *MockChatGroupRepository) Create(ctx context.Context, group *model.ChatGroup) error {
	return nil
}

func (m *MockChatGroupRepository) Get(ctx context.Context, id uint64) (*model.ChatGroup, error) {
	return nil, nil
}

func (m *MockChatGroupRepository) GetByRelatedOrderID(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
	if m.getByRelatedOrderID != nil {
		return m.getByRelatedOrderID(ctx, orderID)
	}
	return nil, nil
}

func (m *MockChatGroupRepository) ListByUser(ctx context.Context, userID uint64, opts repository.ChatGroupListOptions) ([]model.ChatGroup, int64, error) {
	return nil, 0, nil
}

func (m *MockChatGroupRepository) ListMembers(ctx context.Context, groupID uint64, opts repository.ChatGroupMemberListOptions) ([]model.ChatGroupMember, int64, error) {
	return nil, 0, nil
}

func (m *MockChatGroupRepository) Update(ctx context.Context, group *model.ChatGroup) error {
	return nil
}

func (m *MockChatGroupRepository) Deactivate(ctx context.Context, id uint64) error {
	if m.deactivate != nil {
		return m.deactivate(ctx, id)
	}
	return nil
}

func (m *MockChatGroupRepository) ListDeactivatedBefore(ctx context.Context, cutoff time.Time, limit int) ([]model.ChatGroup, error) {
	return nil, nil
}

func (m *MockChatGroupRepository) DeleteByIDs(ctx context.Context, ids []uint64) error {
	return nil
}

// TestOrderService_CreateOrder_WithDistributedLock tests order creation with distributed lock
func TestOrderService_CreateOrder_WithDistributedLock(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{
		createOrder: func(ctx context.Context, order *model.Order) error {
			order.ID = 123
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Inject distributed lock
	mockLock := &MockDistributedLock{
		tryLock: func(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error) {
			return true, nil
		},
		unlock: func(ctx context.Context, key string) error {
			return nil
		},
	}
	service.SetDistributedLock(mockLock)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.OrderID)
}

// TestOrderService_CreateOrder_LockConflict tests order creation with lock conflict
func TestOrderService_CreateOrder_LockConflict(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Inject distributed lock that returns false (locked)
	mockLock := &MockDistributedLock{
		tryLock: func(ctx context.Context, key string, ttl time.Duration, retry int, interval time.Duration) (bool, error) {
			return false, nil // Already locked
		},
	}
	service.SetDistributedLock(mockLock)

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Test Order",
		ScheduledStart: &scheduledStart,
		DurationHours:  1,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

// TestOrderService_SetDistributedLock tests setting distributed lock
func TestOrderService_SetDistributedLock(t *testing.T) {
	service := &OrderService{}

	var mockLock cache.DistributedLock = &MockDistributedLock{}
	service.SetDistributedLock(mockLock)

	assert.Equal(t, mockLock, service.distributedLock)
}

// TestOrderService_SetChatGroupRepository tests setting chat group repository
func TestOrderService_SetChatGroupRepository(t *testing.T) {
	service := &OrderService{}

	mockChatGroups := &MockChatGroupRepository{}
	service.SetChatGroupRepository(mockChatGroups)

	assert.Equal(t, mockChatGroups, service.chatGroups)
}

// TestOrderService_deactivateOrderChat_NoRepository tests deactivation when no chat repository
func TestOrderService_deactivateOrderChat_NoRepository(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	service := &OrderService{}

	// Should not panic when chatGroups is nil
	service.deactivateOrderChat(ctx, orderID)
}

// TestOrderService_deactivateOrderChat_NoGroup tests deactivation when no group found
func TestOrderService_deactivateOrderChat_NoGroup(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return nil, nil // No group found
		},
	}

	service := &OrderService{
		chatGroups: mockChatGroups,
	}

	// Should not error when no group found
	service.deactivateOrderChat(ctx, orderID)
}

// TestOrderService_deactivateOrderChat_Error tests deactivation when repository returns error
func TestOrderService_deactivateOrderChat_Error(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return nil, assert.AnError
		},
	}

	service := &OrderService{
		chatGroups: mockChatGroups,
	}

	// Should not error when repository errors
	service.deactivateOrderChat(ctx, orderID)
}

// TestOrderService_deactivateOrderChat_NotActive tests deactivation when group is not active
func TestOrderService_deactivateOrderChat_NotActive(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return &model.ChatGroup{
				Base:           model.Base{ID: 1},
				IsActive:       false,
				GroupType:      model.ChatGroupTypeOrder,
				RelatedOrderID: &orderID,
			}, nil
		},
	}

	service := &OrderService{
		chatGroups: mockChatGroups,
	}

	// Should not try to deactivate inactive group
	service.deactivateOrderChat(ctx, orderID)
}

// TestOrderService_deactivateOrderChat_WrongType tests deactivation when group is not order type
func TestOrderService_deactivateOrderChat_WrongType(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return &model.ChatGroup{
				Base:           model.Base{ID: 1},
				IsActive:       true,
				GroupType:      model.ChatGroupTypePublic, // Wrong type
				RelatedOrderID: &orderID,
			}, nil
		},
	}

	service := &OrderService{
		chatGroups: mockChatGroups,
	}

	// Should not deactivate non-order groups
	service.deactivateOrderChat(ctx, orderID)
}

// TestOrderService_deactivateOrderChat_Success tests successful chat deactivation
func TestOrderService_deactivateOrderChat_Success(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	var deactivatedGroupID uint64
	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return &model.ChatGroup{
				Base:           model.Base{ID: 123},
				IsActive:       true,
				GroupType:      model.ChatGroupTypeOrder,
				RelatedOrderID: &orderID,
			}, nil
		},
		deactivate: func(ctx context.Context, groupID uint64) error {
			deactivatedGroupID = groupID
			return nil
		},
	}

	service := &OrderService{
		chatGroups: mockChatGroups,
	}

	service.deactivateOrderChat(ctx, orderID)

	assert.Equal(t, uint64(123), deactivatedGroupID)
}

// TestOrderService_CancelOrder_WithChatDeactivation tests cancellation with chat deactivation
func TestOrderService_CancelOrder_WithChatDeactivation(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	var deactivatedGroupID uint64
	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return &model.ChatGroup{
				Base:           model.Base{ID: 123},
				IsActive:       true,
				GroupType:      model.ChatGroupTypeOrder,
				RelatedOrderID: &orderID,
			}, nil
		},
		deactivate: func(ctx context.Context, groupID uint64) error {
			deactivatedGroupID = groupID
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)
	service.SetChatGroupRepository(mockChatGroups)

	req := CancelOrderRequest{
		Reason: "Test",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.Equal(t, uint64(123), deactivatedGroupID)
}

// TestOrderService_CompleteOrder_WithChatDeactivation tests completion with chat deactivation
func TestOrderService_CompleteOrder_WithChatDeactivation(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusInProgress)

	var deactivatedGroupID uint64
	mockChatGroups := &MockChatGroupRepository{
		getByRelatedOrderID: func(ctx context.Context, orderID uint64) (*model.ChatGroup, error) {
			return &model.ChatGroup{
				Base:           model.Base{ID: 123},
				IsActive:       true,
				GroupType:      model.ChatGroupTypeOrder,
				RelatedOrderID: &orderID,
			}, nil
		},
		deactivate: func(ctx context.Context, groupID uint64) error {
			deactivatedGroupID = groupID
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			return nil
		},
	}

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 20}, nil
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			return nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)
	service.SetChatGroupRepository(mockChatGroups)

	err := service.CompleteOrder(ctx, userID, orderID)

	require.NoError(t, err)
	assert.Equal(t, uint64(123), deactivatedGroupID)
}

// TestOrderService_recordCommissionAsync_WithDefaultRule tests commission recording with default rule
func TestOrderService_recordCommissionAsync_WithDefaultRule(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	var createdRecord *model.CommissionRecord
	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return nil, nil // No specific rule
		},
		getDefaultRule: func(ctx context.Context) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 25}, nil // Default 25%
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			createdRecord = record
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.recordCommissionAsync(ctx, orderID)

	require.NoError(t, err)
	assert.NotNil(t, createdRecord)
	assert.Equal(t, 25, createdRecord.CommissionRate)
	assert.Equal(t, int64(1250), createdRecord.CommissionCents) // 25% of 5000
}

// TestOrderService_recordCommissionAsync_DefaultRuleError tests commission recording when default rule fails
func TestOrderService_recordCommissionAsync_DefaultRuleError(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	var createdRecord *model.CommissionRecord
	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return nil, nil
		},
		getDefaultRule: func(ctx context.Context) (*model.CommissionRule, error) {
			return nil, assert.AnError
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			createdRecord = record
			return nil
		},
	}

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	// Should still succeed with hardcoded default rate of 20%
	err := service.recordCommissionAsync(ctx, orderID)

	// The function should not error even if getting default rule fails
	// because it falls back to hardcoded default
	require.NoError(t, err)
	assert.NotNil(t, createdRecord)
	assert.Equal(t, 20, createdRecord.CommissionRate) // Should use hardcoded default
}

// TestOrderService_GetMyOrders_WithReviewCheck tests order listing with review status check
func TestOrderService_GetMyOrders_WithReviewCheck(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrder := createTestOrder(1, userID, model.OrderStatusCompleted)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return []model.Order{*testOrder}, 1, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	// Mock existing review
	reviews := &MockReviewRepository{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				{Base: model.Base{ID: 1}, OrderID: 1, UserID: userID},
			}, 1, nil
		},
	}

	payments := &MockPaymentRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := MyOrderListRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Orders))
	// CanReview should be false because review exists
	assert.False(t, resp.Orders[0].CanReview)
}

// TestOrderService_GetOrderDetail_WithReview tests order detail with review
func TestOrderService_GetOrderDetail_WithReview(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusCompleted)
	now := time.Now()
	testOrder.CompletedAt = &now

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	paymentTime := time.Now()
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, Method: model.PaymentMethodAlipay, AmountCents: 5000, Status: model.PaymentStatusPaid, PaidAt: &paymentTime},
			}, 1, nil
		},
	}

	// Mock review
	reviews := &MockReviewRepository{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				{Base: model.Base{ID: 1}, OrderID: orderID, UserID: userID, Score: 5, Content: "Great service!"},
			}, 1, nil
		},
	}

	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Review)
	assert.Equal(t, 5, resp.Review.Rating)
	assert.Equal(t, "Great service!", resp.Review.Comment)
}

// TestOrderService_GetOrderDetail_NoPlayerInfo tests order detail when player info not available
func TestOrderService_GetOrderDetail_NoPlayerInfo(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return nil, repository.ErrNotFound
		},
	}

	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Player) // No player info available
}

// TestOrderService_AcceptOrder_UpdateError tests accept order with update error
func TestOrderService_AcceptOrder_UpdateError(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	orderID := uint64(1)

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(100, playerUserID), nil
		},
	}

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, id uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			return false, assert.AnError
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.AcceptOrder(ctx, playerUserID, orderID)

	assert.Error(t, err)
}

// TestOrderService_CompleteOrderByPlayer_InvalidStatus tests player completion with invalid status
func TestOrderService_CompleteOrderByPlayer_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	// Order is pending, not in progress
	testOrder := createTestOrder(orderID, 1, model.OrderStatusPending)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_GetAvailableOrders_CalculateDuration tests available orders with duration calculation
func TestOrderService_GetAvailableOrders_CalculateDuration(t *testing.T) {
	ctx := context.Background()

	scheduledStart := time.Now()
	scheduledEnd := scheduledStart.Add(2 * time.Hour)

	testOrder := createTestOrder(1, 1, model.OrderStatusConfirmed)
	testOrder.ScheduledStart = &scheduledStart
	testOrder.ScheduledEnd = &scheduledEnd

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return []model.Order{*testOrder}, 1, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		Page:     1,
		PageSize: 10,
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, len(ordersResp))
	assert.Equal(t, int64(1), total)
	// Duration should be approximately 2 hours
	assert.GreaterOrEqual(t, ordersResp[0].DurationHours, float32(1.9))
	assert.LessOrEqual(t, ordersResp[0].DurationHours, float32(2.1))
}

// TestOrderService_GetAvailableOrders_NoDurationInfo tests available orders without duration info
func TestOrderService_GetAvailableOrders_NoDurationInfo(t *testing.T) {
	ctx := context.Background()

	testOrder := createTestOrder(1, 1, model.OrderStatusConfirmed)
	testOrder.ScheduledStart = nil
	testOrder.ScheduledEnd = nil

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return []model.Order{*testOrder}, 1, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}

	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		Page:     1,
		PageSize: 10,
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 1, len(ordersResp))
	assert.Equal(t, int64(1), total)
	assert.Equal(t, float32(0), ordersResp[0].DurationHours)
}

// TestOrderService_GetMyOrders_PageSizeTooLarge tests pagination with page size exceeding maximum
func TestOrderService_GetMyOrders_PageSizeTooLarge(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	var capturedPageSize int
	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			capturedPageSize = opts.PageSize
			return []model.Order{}, 0, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := MyOrderListRequest{
		Page:     1,
		PageSize: 200, // Exceeds max of 100, should default to 20
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	// Should default to 20 when out of valid range
	assert.Equal(t, 20, capturedPageSize)
}

// TestOrderService_GetAvailableOrders_PageSizeTooLarge tests available orders pagination
func TestOrderService_GetAvailableOrders_PageSizeTooLarge(t *testing.T) {
	ctx := context.Background()

	var capturedPageSize int
	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			capturedPageSize = opts.PageSize
			return []model.Order{}, 0, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	req := AvailableOrdersRequest{
		Page:     1,
		PageSize: 200, // Exceeds max of 100, should default to 20
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, ordersResp)
	// Should default to 20 when out of valid range
	assert.Equal(t, 20, capturedPageSize)
	assert.Equal(t, int64(0), total)
}

// TestOrderService_toOrderCardDTO_WithReviewCheck tests order card DTO with review check
func TestOrderService_toOrderCardDTO_WithReviewCheck(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrder := createTestOrder(1, userID, model.OrderStatusCompleted)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}
	payments := &MockPaymentRepository{}
	// Mock review exists
	reviews := &MockReviewRepository{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{
				{Base: model.Base{ID: 1}, OrderID: 1, UserID: userID},
			}, 1, nil
		},
	}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, testOrder, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.False(t, card.CanReview) // Already reviewed
}

// TestOrderService_toOrderCardDTO_NoReviewYet tests order card DTO when not yet reviewed
func TestOrderService_toOrderCardDTO_NoReviewYet(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrder := createTestOrder(1, userID, model.OrderStatusCompleted)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return createTestPlayer(id, 200), nil
		},
	}
	users := &MockUserRepository{
		getUser: func(ctx context.Context, id uint64) (*model.User, error) {
			return createTestUser(id), nil
		},
	}
	games := &MockGameRepository{
		getGame: func(ctx context.Context, id uint64) (*model.Game, error) {
			return createTestGame(id), nil
		},
	}
	payments := &MockPaymentRepository{}
	// No reviews yet
	reviews := &MockReviewRepository{
		listReviews: func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
			return []model.Review{}, 0, nil
		},
	}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	card, err := service.toOrderCardDTO(ctx, testOrder, userID)

	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.True(t, card.CanReview) // Can review
}

// TestOrderService_buildOrderTimeline_CompletedOrder tests timeline for completed order
func TestOrderService_buildOrderTimeline_CompletedOrder(t *testing.T) {
	testOrder := createTestOrder(1, 1, model.OrderStatusCompleted)
	now := time.Now()
	testOrder.StartedAt = &now
	testOrder.CompletedAt = &now

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	timeline := service.buildOrderTimeline(testOrder)

	assert.GreaterOrEqual(t, len(timeline), 3)
	completedFound := false
	for _, item := range timeline {
		if item.Status == string(model.OrderStatusCompleted) {
			completedFound = true
		}
	}
	assert.True(t, completedFound, "Completed status should be in timeline")
}

// TestOrderService_buildOrderTimeline_InProgressOrder tests timeline for in-progress order
func TestOrderService_buildOrderTimeline_InProgressOrder(t *testing.T) {
	testOrder := createTestOrder(1, 1, model.OrderStatusInProgress)
	now := time.Now()
	testOrder.StartedAt = &now

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	timeline := service.buildOrderTimeline(testOrder)

	assert.GreaterOrEqual(t, len(timeline), 2)
	inProgressFound := false
	for _, item := range timeline {
		if item.Status == string(model.OrderStatusInProgress) {
			inProgressFound = true
		}
	}
	assert.True(t, inProgressFound, "InProgress status should be in timeline")
}

// TestOrderService_buildOrderTimeline_PendingOrder tests timeline for pending order
func TestOrderService_buildOrderTimeline_PendingOrder(t *testing.T) {
	testOrder := createTestOrder(1, 1, model.OrderStatusPending)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(orders, players, users, games, payments, reviews, commissions)

	timeline := service.buildOrderTimeline(testOrder)

	assert.Equal(t, 1, len(timeline)) // Only created event
	assert.Equal(t, string(model.OrderStatusPending), timeline[0].Status)
}

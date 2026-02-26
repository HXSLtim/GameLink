package order

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	commissionrepo "gamelink/internal/repository/commission"
	repoiface "gamelink/internal/repository/interfaces"
)

// Mock implementations for testing
type MockOrderRepository struct {
	createOrder         func(ctx context.Context, order *model.Order) error
	getOrder            func(ctx context.Context, id uint64) (*model.Order, error)
	getOrdersByIDs      func(ctx context.Context, ids []uint64) ([]model.Order, error)
	updateOrder         func(ctx context.Context, order *model.Order) error
	updateWithCondition func(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error)
	listOrders          func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error)
	deleteOrder         func(ctx context.Context, id uint64) error
}

func (m *MockOrderRepository) Create(ctx context.Context, order *model.Order) error {
	if m.createOrder != nil {
		return m.createOrder(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) Get(ctx context.Context, id uint64) (*model.Order, error) {
	if m.getOrder != nil {
		return m.getOrder(ctx, id)
	}
	return nil, nil
}

func (m *MockOrderRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Order, error) {
	if m.getOrdersByIDs != nil {
		return m.getOrdersByIDs(ctx, ids)
	}
	return []model.Order{}, nil
}

func (m *MockOrderRepository) Update(ctx context.Context, order *model.Order) error {
	if m.updateOrder != nil {
		return m.updateOrder(ctx, order)
	}
	return nil
}

func (m *MockOrderRepository) UpdateWithCondition(ctx context.Context, orderID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
	if m.updateWithCondition != nil {
		return m.updateWithCondition(ctx, orderID, expectedStatus, updates)
	}
	return true, nil
}

func (m *MockOrderRepository) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	if m.listOrders != nil {
		return m.listOrders(ctx, opts)
	}
	return []model.Order{}, 0, nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint64) error {
	if m.deleteOrder != nil {
		return m.deleteOrder(ctx, id)
	}
	return nil
}

type MockPlayerRepository struct {
	getPlayer       func(ctx context.Context, id uint64) (*model.Player, error)
	getPlayerByUser func(ctx context.Context, userID uint64) (*model.Player, error)
	updatePlayer    func(ctx context.Context, player *model.Player) error
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	if m.getPlayer != nil {
		return m.getPlayer(ctx, id)
	}
	return nil, nil
}

func (m *MockPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	if m.getPlayerByUser != nil {
		return m.getPlayerByUser(ctx, userID)
	}
	return nil, nil
}

func (m *MockPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	if m.updatePlayer != nil {
		return m.updatePlayer(ctx, player)
	}
	return nil
}

func (m *MockPlayerRepository) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (m *MockPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockPlayerRepository) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *MockPlayerRepository) Delete(ctx context.Context, id uint64) error            { return nil }
func (m *MockPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	return 0, nil
}
func (m *MockPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	return 0, nil
}
func (m *MockPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	return 0, nil
}
func (m *MockPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return 0, nil
}
func (m *MockPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	return nil, nil
}

func (m *MockPlayerRepository) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, int64(0), nil
}

type MockUserRepository struct {
	getUser func(ctx context.Context, id uint64) (*model.User, error)
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	if m.getUser != nil {
		return m.getUser(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *MockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}
func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	return nil, nil
}
func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error { return nil }
func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
	return nil
}

func (m *MockUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	return nil, nil
}

type MockGameRepository struct {
	getGame func(ctx context.Context, id uint64) (*model.Game, error)
}

func (m *MockGameRepository) Get(ctx context.Context, id uint64) (*model.Game, error) {
	if m.getGame != nil {
		return m.getGame(ctx, id)
	}
	return nil, nil
}

func (m *MockGameRepository) List(ctx context.Context) ([]model.Game, error) { return nil, nil }
func (m *MockGameRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (m *MockGameRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string) ([]model.Game, int64, error) {
	return nil, 0, nil
}
func (m *MockGameRepository) Create(ctx context.Context, game *model.Game) error { return nil }
func (m *MockGameRepository) Update(ctx context.Context, game *model.Game) error { return nil }
func (m *MockGameRepository) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *MockGameRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return 0, nil
}
func (m *MockGameRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	return 0, nil
}
func (m *MockGameRepository) BatchUpdateSortOrder(ctx context.Context, updates map[uint64]int) (int64, error) {
	return 0, nil
}
func (m *MockGameRepository) BatchUpdateCategory(ctx context.Context, ids []uint64, category string) (int64, error) {
	return 0, nil
}
func (m *MockGameRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Game, error) {
	return nil, nil
}

type MockPaymentRepository struct {
	listPayments  func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error)
	updatePayment func(ctx context.Context, payment *model.Payment) error
}

func (m *MockPaymentRepository) List(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
	if m.listPayments != nil {
		return m.listPayments(ctx, opts)
	}
	return nil, 0, nil
}

func (m *MockPaymentRepository) Create(ctx context.Context, payment *model.Payment) error { return nil }
func (m *MockPaymentRepository) Get(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, nil
}
func (m *MockPaymentRepository) GetWithRelations(ctx context.Context, id uint64) (*model.Payment, error) {
	return nil, nil
}
func (m *MockPaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	if m.updatePayment != nil {
		return m.updatePayment(ctx, payment)
	}
	return nil
}
func (m *MockPaymentRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockPaymentRepository) GetByOrderID(ctx context.Context, orderID uint64) ([]model.Payment, error) {
	return nil, nil
}
func (m *MockPaymentRepository) GetByRequestID(ctx context.Context, requestID string) (*model.Payment, error) {
	return nil, nil
}

type MockReviewRepository struct {
	listReviews func(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error)
}

func (m *MockReviewRepository) List(ctx context.Context, opts repository.ReviewListOptions) ([]model.Review, int64, error) {
	if m.listReviews != nil {
		return m.listReviews(ctx, opts)
	}
	return nil, 0, nil
}

func (m *MockReviewRepository) Get(ctx context.Context, id uint64) (*model.Review, error) {
	return nil, nil
}
func (m *MockReviewRepository) Create(ctx context.Context, review *model.Review) error { return nil }
func (m *MockReviewRepository) Update(ctx context.Context, review *model.Review) error { return nil }
func (m *MockReviewRepository) UpdateStatus(ctx context.Context, id uint64, status model.ReviewStatus, rejectionReason string) error {
	return nil
}
func (m *MockReviewRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.ReviewStatus, rejectionReason string) error {
	return nil
}
func (m *MockReviewRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockReviewRepository) ListPending(ctx context.Context, page, pageSize int) ([]model.Review, int64, error) {
	return nil, 0, nil
}
func (m *MockReviewRepository) GetStats(ctx context.Context) (repository.ReviewStats, error) {
	return repository.ReviewStats{}, nil
}
func (m *MockReviewRepository) GetTrend(ctx context.Context, days int) ([]repository.DateValue, error) {
	return nil, nil
}
func (m *MockReviewRepository) GetTopPlayersByReviewCount(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	return nil, nil
}
func (m *MockReviewRepository) GetTopPlayersByRating(ctx context.Context, limit int) ([]repository.PlayerReviewStats, error) {
	return nil, nil
}
func (m *MockReviewRepository) GetGameStats(ctx context.Context) ([]repository.GameReviewStats, error) {
	return nil, nil
}

type MockCommissionRepository struct {
	getRule            func(ctx context.Context, id uint64) (*model.CommissionRule, error)
	getRuleForOrder    func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error)
	getDefaultRule     func(ctx context.Context) (*model.CommissionRule, error)
	getRecordByOrderID func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error)
	createRecord       func(ctx context.Context, record *model.CommissionRecord) error
}

func (m *MockCommissionRepository) GetRule(ctx context.Context, id uint64) (*model.CommissionRule, error) {
	if m.getRule != nil {
		return m.getRule(ctx, id)
	}
	return nil, nil
}

func (m *MockCommissionRepository) GetRuleForOrder(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
	if m.getRuleForOrder != nil {
		return m.getRuleForOrder(ctx, gameID, playerID, serviceType)
	}
	return nil, nil
}

func (m *MockCommissionRepository) GetDefaultRule(ctx context.Context) (*model.CommissionRule, error) {
	if m.getDefaultRule != nil {
		return m.getDefaultRule(ctx)
	}
	return &model.CommissionRule{Rate: 20}, nil
}

func (m *MockCommissionRepository) GetRecordByOrderID(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
	if m.getRecordByOrderID != nil {
		return m.getRecordByOrderID(ctx, orderID)
	}
	return nil, nil
}

func (m *MockCommissionRepository) CreateRecord(ctx context.Context, record *model.CommissionRecord) error {
	if m.createRecord != nil {
		return m.createRecord(ctx, record)
	}
	return nil
}

func (m *MockCommissionRepository) CreateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *MockCommissionRepository) ListRules(ctx context.Context, opts commissionrepo.CommissionRuleListOptions) ([]model.CommissionRule, int64, error) {
	return nil, 0, nil
}
func (m *MockCommissionRepository) UpdateRule(ctx context.Context, rule *model.CommissionRule) error {
	return nil
}
func (m *MockCommissionRepository) DeleteRule(ctx context.Context, id uint64) error { return nil }
func (m *MockCommissionRepository) GetRecord(ctx context.Context, id uint64) (*model.CommissionRecord, error) {
	return nil, nil
}
func (m *MockCommissionRepository) ListRecords(ctx context.Context, opts commissionrepo.CommissionRecordListOptions) ([]model.CommissionRecord, int64, error) {
	return nil, 0, nil
}
func (m *MockCommissionRepository) UpdateRecord(ctx context.Context, record *model.CommissionRecord) error {
	return nil
}
func (m *MockCommissionRepository) CreateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (m *MockCommissionRepository) GetSettlement(ctx context.Context, id uint64) (*model.MonthlySettlement, error) {
	return nil, nil
}
func (m *MockCommissionRepository) GetSettlementByPlayerMonth(ctx context.Context, playerID uint64, month string) (*model.MonthlySettlement, error) {
	return nil, nil
}
func (m *MockCommissionRepository) ListSettlements(ctx context.Context, opts commissionrepo.SettlementListOptions) ([]model.MonthlySettlement, int64, error) {
	return nil, 0, nil
}
func (m *MockCommissionRepository) UpdateSettlement(ctx context.Context, settlement *model.MonthlySettlement) error {
	return nil
}
func (m *MockCommissionRepository) GetPlayerMonthlyIncome(ctx context.Context, playerID uint64, month string) (int64, error) {
	return 0, nil
}
func (m *MockCommissionRepository) GetMonthlyStats(ctx context.Context, month string) (*commissionrepo.MonthlyStats, error) {
	return &commissionrepo.MonthlyStats{}, nil
}

// Helper function to create test order
func createTestOrder(id uint64, userID uint64, status model.OrderStatus) *model.Order {
	playerID := uint64(100)
	gameID := uint64(1)
	now := time.Now()
	scheduledStart := now.Add(time.Hour)
	scheduledEnd := scheduledStart.Add(time.Hour)

	return &model.Order{
		Base:              model.Base{ID: id},
		OrderNo:           model.GenerateEscortOrderNo(),
		UserID:            userID,
		ItemID:            1,
		PlayerID:          &playerID,
		GameID:            &gameID,
		Quantity:          1,
		UnitPriceCents:    5000,
		TotalPriceCents:   5000,
		CommissionCents:   1000,
		PlayerIncomeCents: 4000,
		Currency:          model.CurrencyCNY,
		Status:            status,
		Title:             "Test Order",
		Description:       "Test Description",
		ScheduledStart:    &scheduledStart,
		ScheduledEnd:      &scheduledEnd,
		OrderConfig:       "{}",
	}
}

func createTestPlayer(id uint64, userID uint64) *model.Player {
	return &model.Player{
		Base:               model.Base{ID: id},
		UserID:             userID,
		Nickname:           "TestPlayer",
		Rank:               "Gold",
		HourlyRateCents:    5000,
		VerificationStatus: model.VerificationVerified,
	}
}

func createTestUser(id uint64) *model.User {
	return &model.User{
		Base:      model.Base{ID: id},
		Name:      "Test User",
		AvatarURL: "https://example.com/avatar.jpg",
	}
}

func createTestGame(id uint64) *model.Game {
	return &model.Game{
		Base: model.Base{ID: id},
		Name: "Test Game",
	}
}

// TestOrderService_CreateOrder_Success tests successful order creation
func TestOrderService_CreateOrder_Success(t *testing.T) {
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

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := CreateOrderRequest{
		PlayerID:       playerID,
		GameID:         gameID,
		Title:          "Test Order",
		Description:    "Test Description",
		ScheduledStart: &scheduledStart,
		DurationHours:  2,
	}

	resp, err := service.CreateOrder(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.OrderID)
	assert.Equal(t, int64(10000), resp.PriceCents) // 5000 * 2 hours
	assert.True(t, resp.NeedPayment)
}

// TestOrderService_CreateOrder_PlayerNotFound tests order creation when player doesn't exist
func TestOrderService_CreateOrder_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(999)
	gameID := uint64(1)
	scheduledStart := time.Now().Add(time.Hour)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return nil, repository.ErrNotFound
		},
	}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

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
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestOrderService_CreateOrder_GameNotFound tests order creation when game doesn't exist
func TestOrderService_CreateOrder_GameNotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	playerID := uint64(100)
	gameID := uint64(999)
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
			return nil, repository.ErrNotFound
		},
	}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

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
	assert.Equal(t, repository.ErrNotFound, err)
}

// TestOrderService_GetMyOrders_Success tests successful retrieval of user's orders
func TestOrderService_GetMyOrders_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrders := []model.Order{
		*createTestOrder(1, userID, model.OrderStatusPending),
		*createTestOrder(2, userID, model.OrderStatusConfirmed),
	}

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			return testOrders, 2, nil
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

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := MyOrderListRequest{
		Page:     1,
		PageSize: 10,
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Orders))
	assert.Equal(t, int64(2), resp.Total)
}

// TestOrderService_GetMyOrders_WithStatusFilter tests order retrieval with status filter
func TestOrderService_GetMyOrders_WithStatusFilter(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	testOrders := []model.Order{
		*createTestOrder(1, userID, model.OrderStatusPending),
	}

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			assert.Equal(t, []model.OrderStatus{model.OrderStatusPending}, opts.Statuses)
			return testOrders, 1, nil
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

	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := MyOrderListRequest{
		Status:   "pending",
		Page:     1,
		PageSize: 10,
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Orders))
}

// TestOrderService_GetMyOrders_DefaultPagination tests default pagination parameters
func TestOrderService_GetMyOrders_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			assert.Equal(t, 1, opts.Page)
			assert.Equal(t, 20, opts.PageSize)
			return []model.Order{}, 0, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := MyOrderListRequest{
		Page:     0, // Invalid, should default to 1
		PageSize: 0, // Invalid, should default to 20
	}

	resp, err := service.GetMyOrders(ctx, userID, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestOrderService_GetOrderDetail_Success tests successful order detail retrieval
func TestOrderService_GetOrderDetail_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusCompleted)
	now := time.Now()
	testOrder.CompletedAt = &now
	// Set preloaded relations expected by GetOrderDetail
	testPlayer := createTestPlayer(100, 200)
	testUser := createTestUser(200)
	testPlayer.User = testUser
	testOrder.Player = testPlayer

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

	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, orderID, resp.Order.ID)
	assert.NotNil(t, resp.Player)
	assert.NotNil(t, resp.Payment)
	assert.NotEmpty(t, resp.Timeline)
}

// TestOrderService_GetOrderDetail_NotFound tests order detail when order doesn't exist
func TestOrderService_GetOrderDetail_NotFound(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(999)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return nil, repository.ErrNotFound
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrNotFound, err)
}

// TestOrderService_GetOrderDetail_Unauthorized tests order detail when user is not authorized
func TestOrderService_GetOrderDetail_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	orderUserID := uint64(999)

	testOrder := createTestOrder(orderID, orderUserID, model.OrderStatusCompleted)

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
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	resp, err := service.GetOrderDetail(ctx, userID, orderID)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestOrderService_CancelOrder_Success tests successful order cancellation
func TestOrderService_CancelOrder_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

	// CancelOrder now uses UpdateWithCondition (atomic update)
	var capturedUpdates map[string]any
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateWithCondition: func(ctx context.Context, oID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			capturedUpdates = updates
			return true, nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := CancelOrderRequest{
		Reason: "Changed mind",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.NotNil(t, capturedUpdates)
	assert.Equal(t, model.OrderStatusCanceled, capturedUpdates["status"])
	assert.Equal(t, "Changed mind", capturedUpdates["cancel_reason"])
}

// TestOrderService_CancelOrder_WithRefund tests order cancellation with refund
func TestOrderService_CancelOrder_WithRefund(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusConfirmed)

	// CancelOrder now uses UpdateWithCondition (atomic update) for the fallback refund path
	var capturedUpdates map[string]any
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateWithCondition: func(ctx context.Context, oID uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			capturedUpdates = updates
			return true, nil
		},
	}

	paymentTime := time.Now()
	var updatedPayment *model.Payment
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{Base: model.Base{ID: 1}, Method: model.PaymentMethodAlipay, AmountCents: 5000, Status: model.PaymentStatusPaid, PaidAt: &paymentTime},
			}, 1, nil
		},
		updatePayment: func(ctx context.Context, payment *model.Payment) error {
			clone := *payment
			updatedPayment = &clone
			return nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := CancelOrderRequest{
		Reason: "Need to cancel",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	require.NoError(t, err)
	assert.NotNil(t, capturedUpdates)
	assert.Equal(t, model.OrderStatusRefunded, capturedUpdates["status"])
	assert.Equal(t, int64(5000), capturedUpdates["refund_amount_cents"])
	assert.Equal(t, "用户取消订单", capturedUpdates["refund_reason"])
	assert.NotNil(t, capturedUpdates["refunded_at"])
	require.NotNil(t, updatedPayment)
	assert.Equal(t, model.PaymentStatusRefunded, updatedPayment.Status)
	assert.Equal(t, int64(5000), updatedPayment.RefundedAmountCents)
	assert.NotNil(t, updatedPayment.RefundedAt)
}

// TestOrderService_CancelOrder_InvalidStatus tests order cancellation with invalid status
func TestOrderService_CancelOrder_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusCompleted)

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
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := CancelOrderRequest{
		Reason: "Test",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_CancelOrder_Unauthorized tests order cancellation by unauthorized user
func TestOrderService_CancelOrder_Unauthorized(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)
	orderUserID := uint64(999)

	testOrder := createTestOrder(orderID, orderUserID, model.OrderStatusPending)

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
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := CancelOrderRequest{
		Reason: "Test",
	}

	err := service.CancelOrder(ctx, userID, orderID, req)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestOrderService_CompleteOrder_Success tests successful order completion by user
func TestOrderService_CompleteOrder_Success(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusInProgress)

	var updatedOrder *model.Order
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	commissionRecorded := false
	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil // No existing record
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 20}, nil
		},
		createRecord: func(ctx context.Context, record *model.CommissionRecord) error {
			commissionRecorded = true
			return nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{OrderID: orderID, Status: model.PaymentStatusPaid},
			}, 1, nil
		},
	}
	reviews := &MockReviewRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrder(ctx, userID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.CompletedAt)
	assert.True(t, commissionRecorded)
}

// TestOrderService_CompleteOrder_InvalidStatus tests order completion with invalid status
func TestOrderService_CompleteOrder_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusPending)

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
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrder(ctx, userID, orderID)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_CompleteOrder_UnpaidOrder tests completion fails when order has no paid payment
func TestOrderService_CompleteOrder_UnpaidOrder(t *testing.T) {
	ctx := context.Background()
	userID := uint64(1)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, userID, model.OrderStatusInProgress)

	updateCalled := false
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updateCalled = true
			return nil
		},
	}

	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{OrderID: orderID, Status: model.PaymentStatusPending},
			}, 1, nil
		},
	}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrder(ctx, userID, orderID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "order must be paid before completion")
	assert.False(t, updateCalled)
}

// TestOrderService_AcceptOrder_Success tests successful order acceptance by player
func TestOrderService_AcceptOrder_Success(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, id uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			assert.Equal(t, orderID, id)
			assert.Equal(t, model.OrderStatusConfirmed, expectedStatus)
			assert.Equal(t, playerID, updates["player_id"])
			assert.Equal(t, model.OrderStatusInProgress, updates["status"])
			return true, nil
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

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.AcceptOrder(ctx, playerUserID, orderID)

	require.NoError(t, err)
}

// TestOrderService_AcceptOrder_PlayerNotFound tests order acceptance when player doesn't exist
func TestOrderService_AcceptOrder_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(999)
	orderID := uint64(1)

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return nil, repository.ErrNotFound
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.AcceptOrder(ctx, playerUserID, orderID)

	assert.Error(t, err)
}

// TestOrderService_AcceptOrder_AlreadyAccepted tests order acceptance when already accepted
func TestOrderService_AcceptOrder_AlreadyAccepted(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	orders := &MockOrderRepository{
		updateWithCondition: func(ctx context.Context, id uint64, expectedStatus model.OrderStatus, updates map[string]any) (bool, error) {
			return false, nil // Already taken
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

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.AcceptOrder(ctx, playerUserID, orderID)

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTransition, err)
}

// TestOrderService_CompleteOrderByPlayer_Success tests successful order completion by player
func TestOrderService_CompleteOrderByPlayer_Success(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	var updatedOrder *model.Order
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updatedOrder = order
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{OrderID: orderID, Status: model.PaymentStatusPaid},
			}, 1, nil
		},
	}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	require.NoError(t, err)
	assert.NotNil(t, updatedOrder)
	assert.Equal(t, model.OrderStatusCompleted, updatedOrder.Status)
	assert.NotNil(t, updatedOrder.CompletedAt)
}

// TestOrderService_CompleteOrderByPlayer_UnpaidOrder tests player completion fails when order has no paid payment
func TestOrderService_CompleteOrderByPlayer_UnpaidOrder(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	playerID := uint64(100)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	updateCalled := false
	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
		updateOrder: func(ctx context.Context, order *model.Order) error {
			updateCalled = true
			return nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(playerID, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{
		listPayments: func(ctx context.Context, opts repository.PaymentListOptions) ([]model.Payment, int64, error) {
			return []model.Payment{
				{OrderID: orderID, Status: model.PaymentStatusPending},
			}, 1, nil
		},
	}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "order must be paid before completion")
	assert.False(t, updateCalled)
}

// TestOrderService_CompleteOrderByPlayer_Unauthorized tests order completion by unauthorized player
func TestOrderService_CompleteOrderByPlayer_Unauthorized(t *testing.T) {
	ctx := context.Background()
	playerUserID := uint64(200)
	wrongPlayerID := uint64(999)
	orderID := uint64(1)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusInProgress)

	orders := &MockOrderRepository{
		getOrder: func(ctx context.Context, id uint64) (*model.Order, error) {
			return testOrder, nil
		},
	}

	players := &MockPlayerRepository{
		getPlayerByUser: func(ctx context.Context, userID uint64) (*model.Player, error) {
			return createTestPlayer(wrongPlayerID, playerUserID), nil
		},
	}

	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}
	commissions := &MockCommissionRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.CompleteOrderByPlayer(ctx, playerUserID, orderID)

	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

// TestOrderService_GetAvailableOrders_Success tests successful retrieval of available orders
func TestOrderService_GetAvailableOrders_Success(t *testing.T) {
	ctx := context.Background()
	gameID := uint64(1)

	testOrders := []model.Order{
		*createTestOrder(1, 1, model.OrderStatusConfirmed),
		*createTestOrder(2, 2, model.OrderStatusConfirmed),
	}

	orders := &MockOrderRepository{
		listOrders: func(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
			assert.Equal(t, []model.OrderStatus{model.OrderStatusConfirmed}, opts.Statuses)
			return testOrders, 2, nil
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

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	req := AvailableOrdersRequest{
		GameID:   &gameID,
		Page:     1,
		PageSize: 10,
	}

	ordersResp, total, err := service.GetAvailableOrders(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 2, len(ordersResp))
	assert.Equal(t, int64(2), total)
}

// TestOrderService_recordCommissionAsync tests commission recording
func TestOrderService_recordCommissionAsync(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)
	playerID := uint64(100)

	testOrder := createTestOrder(orderID, 1, model.OrderStatusCompleted)

	var createdRecord *model.CommissionRecord
	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return nil, nil // No existing record
		},
		getRuleForOrder: func(ctx context.Context, gameID *uint64, playerID *uint64, serviceType *string) (*model.CommissionRule, error) {
			return &model.CommissionRule{Rate: 20}, nil
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

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.recordCommissionAsync(ctx, orderID)

	require.NoError(t, err)
	assert.NotNil(t, createdRecord)
	assert.Equal(t, orderID, createdRecord.OrderID)
	assert.Equal(t, playerID, createdRecord.PlayerID)
	assert.Equal(t, int64(5000), createdRecord.TotalAmountCents)
	assert.Equal(t, int64(1000), createdRecord.CommissionCents) // 20% of 5000
	assert.Equal(t, int64(4000), createdRecord.PlayerIncomeCents)
	assert.Equal(t, 20, createdRecord.CommissionRate)
}

// TestOrderService_recordCommissionAsync_AlreadyRecorded tests commission recording when already recorded
func TestOrderService_recordCommissionAsync_AlreadyRecorded(t *testing.T) {
	ctx := context.Background()
	orderID := uint64(1)

	existingRecord := &model.CommissionRecord{
		ID:               1,
		OrderID:          orderID,
		TotalAmountCents: 5000,
	}

	commissions := &MockCommissionRepository{
		getRecordByOrderID: func(ctx context.Context, orderID uint64) (*model.CommissionRecord, error) {
			return existingRecord, nil
		},
	}

	orders := &MockOrderRepository{}
	players := &MockPlayerRepository{}
	users := &MockUserRepository{}
	games := &MockGameRepository{}
	payments := &MockPaymentRepository{}
	reviews := &MockReviewRepository{}

	service := NewOrderService(OrderDeps{Orders: orders, Players: players, Users: users, Games: games, Payments: payments, Reviews: reviews, Commissions: commissions})

	err := service.recordCommissionAsync(ctx, orderID)

	require.NoError(t, err) // Should not error, just skip
}

// TestOrderService_calculateOrderPricing tests order pricing calculation
func TestOrderService_calculateOrderPricing(t *testing.T) {
	player := createTestPlayer(100, 200)
	player.HourlyRateCents = 5000

	req := CreateOrderRequest{
		DurationHours: 2,
	}

	service := &OrderService{}

	totalPrice, commissionCents, playerIncomeCents := service.calculateOrderPricing(player, req)

	assert.Equal(t, int64(10000), totalPrice)       // 5000 * 2 hours
	assert.Equal(t, int64(2000), commissionCents)   // 20% of 10000
	assert.Equal(t, int64(8000), playerIncomeCents) // 10000 - 2000
}

// TestOrderService_buildOrderForCreation tests order building
func TestOrderService_buildOrderForCreation(t *testing.T) {
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
		Title:          "Test Order",
		Description:    "Test Description",
		ScheduledStart: &scheduledStart,
		DurationHours:  2,
	}

	order := service.buildOrderForCreation(userID, req, 10000, 2000, 8000)

	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, playerID, *order.PlayerID)
	assert.Equal(t, gameID, *order.GameID)
	assert.Equal(t, serviceID, order.ItemID)
	assert.Equal(t, "Test Order", order.Title)
	assert.Equal(t, "Test Description", order.Description)
	assert.Equal(t, int64(10000), order.TotalPriceCents)
	assert.Equal(t, int64(2000), order.CommissionCents)
	assert.Equal(t, int64(8000), order.PlayerIncomeCents)
	assert.Equal(t, model.OrderStatusPending, order.Status)
	assert.NotNil(t, order.ScheduledEnd)
}

// TestOrderService_validateCreateOrder_Success tests successful order validation
func TestOrderService_validateCreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	playerID := uint64(100)
	gameID := uint64(1)

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

	service := &OrderService{
		players: players,
		games:   games,
	}

	req := CreateOrderRequest{
		PlayerID: playerID,
		GameID:   gameID,
	}

	player, err := service.validateCreateOrder(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, playerID, player.ID)
}

// TestOrderService_validateCreateOrder_PlayerNotFound tests validation with non-existent player
func TestOrderService_validateCreateOrder_PlayerNotFound(t *testing.T) {
	ctx := context.Background()

	players := &MockPlayerRepository{
		getPlayer: func(ctx context.Context, id uint64) (*model.Player, error) {
			return nil, repository.ErrNotFound
		},
	}

	games := &MockGameRepository{}

	service := &OrderService{
		players: players,
		games:   games,
	}

	req := CreateOrderRequest{
		PlayerID: 999,
		GameID:   1,
	}

	player, err := service.validateCreateOrder(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, player)
}

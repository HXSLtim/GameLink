package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	repoiface "gamelink/internal/repository/interfaces"
	withdrawrepo "gamelink/internal/repository/withdraw"
)

// ============================================================================
// Mock Implementations for Earnings Service
// ============================================================================

// MockEarningsPlayerRepository is a mock implementation of PlayerRepository
type MockEarningsPlayerRepository struct {
	mock.Mock
}

func (m *MockEarningsPlayerRepository) Create(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockEarningsPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockEarningsPlayerRepository) Update(ctx context.Context, player *model.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockEarningsPlayerRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEarningsPlayerRepository) List(ctx context.Context) ([]model.Player, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockEarningsPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockEarningsPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, page, pageSize, keyword, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

func (m *MockEarningsPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockEarningsPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Player), args.Error(1)
}

func (m *MockEarningsPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	args := m.Called(ctx, ids, rank)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEarningsPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	args := m.Called(ctx, ids, rateCents)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEarningsPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	args := m.Called(ctx, ids, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEarningsPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockEarningsPlayerRepository) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	args := m.Called(ctx, limit, status)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Player), args.Get(1).(int64), args.Error(2)
}

// MockOrderQuery is a mock implementation of OrderQuery
type MockOrderQuery struct {
	mock.Mock
}

func (m *MockOrderQuery) List(ctx context.Context, opts repoiface.OrderListOptions) ([]model.Order, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderQuery) Get(ctx context.Context, id uint64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderQuery) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

// MockWithdrawRepository is a mock implementation of WithdrawRepository
type MockWithdrawRepository struct {
	mock.Mock
}

func (m *MockWithdrawRepository) Create(ctx context.Context, withdraw *model.Withdraw) error {
	args := m.Called(ctx, withdraw)
	if args.Error(0) == nil && withdraw.ID == 0 {
		withdraw.ID = 1
	}
	return args.Error(0)
}

func (m *MockWithdrawRepository) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Withdraw), args.Error(1)
}

func (m *MockWithdrawRepository) Update(ctx context.Context, withdraw *model.Withdraw) error {
	args := m.Called(ctx, withdraw)
	return args.Error(0)
}

func (m *MockWithdrawRepository) List(ctx context.Context, opts withdrawrepo.WithdrawListOptions) ([]model.Withdraw, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Withdraw), args.Get(1).(int64), args.Error(2)
}

func (m *MockWithdrawRepository) GetPlayerBalance(ctx context.Context, playerID uint64) (*withdrawrepo.PlayerBalance, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*withdrawrepo.PlayerBalance), args.Error(1)
}

func (m *MockWithdrawRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Withdraw, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Withdraw), args.Error(1)
}

func (m *MockWithdrawRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.WithdrawStatus, processedBy *uint64, processedAt *time.Time, reason string) ([]uint64, []withdrawrepo.BatchOperationError, error) {
	args := m.Called(ctx, ids, status, processedBy, processedAt, reason)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]uint64), args.Get(1).([]withdrawrepo.BatchOperationError), args.Error(2)
}

func (m *MockWithdrawRepository) BatchComplete(ctx context.Context, ids []uint64, adminUserID uint64, completedAt time.Time) ([]uint64, []withdrawrepo.BatchOperationError, error) {
	args := m.Called(ctx, ids, adminUserID, completedAt)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]uint64), args.Get(1).([]withdrawrepo.BatchOperationError), args.Error(2)
}

func (m *MockWithdrawRepository) ListByCompany(ctx context.Context, opts withdrawrepo.WithdrawByCompanyOptions) ([]model.Withdraw, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.Withdraw), args.Get(1).(int64), args.Error(2)
}

func (m *MockWithdrawRepository) GetRoutingStats(ctx context.Context, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStatsResponse, error) {
	args := m.Called(ctx, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WithdrawRoutingStatsResponse), args.Error(1)
}

func (m *MockWithdrawRepository) GetRoutingStatsByCompany(ctx context.Context, companyID uint64, dateFrom, dateTo *time.Time) (*model.WithdrawRoutingStats, error) {
	args := m.Called(ctx, companyID, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WithdrawRoutingStats), args.Error(1)
}

func (m *MockWithdrawRepository) CreateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockWithdrawRepository) GetSalaryPaymentRecordByWithdrawID(ctx context.Context, withdrawID uint64) (*model.SalaryPaymentRecord, error) {
	args := m.Called(ctx, withdrawID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SalaryPaymentRecord), args.Error(1)
}

func (m *MockWithdrawRepository) UpdateSalaryPaymentRecord(ctx context.Context, record *model.SalaryPaymentRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestNewEarningsService(t *testing.T) {
	players := &MockEarningsPlayerRepository{}
	orders := &MockOrderQuery{}
	withdraws := &MockWithdrawRepository{}

	svc := NewEarningsService(players, orders, withdraws)

	assert.NotNil(t, svc)
}

func TestEarningsService_GetEarningsSummary(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		setupMock   func(*MockEarningsPlayerRepository, *MockOrderQuery, *MockWithdrawRepository)
		expectError bool
	}{
		{
			name:   "successful get earnings summary",
			userID: 1,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// Mock order queries for today, month, and total
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{{TotalPriceCents: 10000}}, int64(1), nil,
				)

				// Mock balance
				withdraws.On("GetPlayerBalance", mock.Anything, uint64(1)).Return(&withdrawrepo.PlayerBalance{
					TotalEarnings:    100000,
					AvailableBalance: 80000,
					PendingBalance:   20000,
					WithdrawTotal:    0,
				}, nil)
			},
			expectError: false,
		},
		{
			name:   "player not found",
			userID: 999,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockEarningsPlayerRepository{}
			orders := &MockOrderQuery{}
			withdraws := &MockWithdrawRepository{}

			tt.setupMock(players, orders, withdraws)

			svc := NewEarningsService(players, orders, withdraws)
			summary, err := svc.GetEarningsSummary(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, summary)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, summary)
			}
		})
	}
}

func TestEarningsService_GetEarningsTrend(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		days        int
		setupMock   func(*MockEarningsPlayerRepository, *MockOrderQuery, *MockWithdrawRepository)
		expectError bool
		expectDays  int
	}{
		{
			name:   "successful get 7 day trend",
			userID: 1,
			days:   7,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// Mock order queries for each day
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{{TotalPriceCents: 10000}}, int64(1), nil,
				)
			},
			expectError: false,
			expectDays:  7,
		},
		{
			name:   "days less than 7 defaults to 7",
			userID: 1,
			days:   3,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(0), nil,
				)
			},
			expectError: false,
			expectDays:  7,
		},
		{
			name:   "days more than 90 caps at 90",
			userID: 1,
			days:   100,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(0), nil,
				)
			},
			expectError: false,
			expectDays:  90,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockEarningsPlayerRepository{}
			orders := &MockOrderQuery{}
			withdraws := &MockWithdrawRepository{}

			tt.setupMock(players, orders, withdraws)

			svc := NewEarningsService(players, orders, withdraws)
			trend, err := svc.GetEarningsTrend(context.Background(), tt.userID, tt.days)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, trend)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, trend)
				assert.Len(t, trend.Trend, tt.expectDays)
			}
		})
	}
}

func TestFormatCents(t *testing.T) {
	tests := []struct {
		cents    int64
		expected string
	}{
		{10000, "100.00元"},
		{12345, "123.45元"},
		{100, "1.00元"},
		{0, "0.00元"},
		{1, "0.01元"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatCents(tt.cents)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEarningsService_RequestWithdraw(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint64
		req           WithdrawRequest
		setupMock     func(*MockEarningsPlayerRepository, *MockOrderQuery, *MockWithdrawRepository)
		expectError   bool
		errorContains string
	}{
		{
			name:   "successful withdraw request",
			userID: 1,
			req: WithdrawRequest{
				AmountCents: 10000, // 100元
				Method:      "wechat",
				AccountInfo: "test_account",
			},
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// No pending withdraws
				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{}, int64(0), nil,
				)

				// Mock balance check
				withdraws.On("GetPlayerBalance", mock.Anything, uint64(1)).Return(&withdrawrepo.PlayerBalance{
					TotalEarnings:    100000,
					AvailableBalance: 80000,
					PendingBalance:   20000,
					WithdrawTotal:    0,
				}, nil)

				// Mock order queries for earnings summary
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{{TotalPriceCents: 10000}}, int64(1), nil,
				)

				// Create withdraw
				withdraws.On("Create", mock.Anything, mock.AnythingOfType("*model.Withdraw")).Return(nil)
			},
			expectError: false,
		},
		{
			name:   "amount too small",
			userID: 1,
			req: WithdrawRequest{
				AmountCents: 5000, // 50元 < 100元 minimum
				Method:      "wechat",
				AccountInfo: "test_account",
			},
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)
			},
			expectError:   true,
			errorContains: "提现金额不能少于100元",
		},
		{
			name:   "player not found",
			userID: 999,
			req: WithdrawRequest{
				AmountCents: 10000,
				Method:      "wechat",
				AccountInfo: "test_account",
			},
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError:   true,
			errorContains: "not found",
		},
		{
			name:   "has pending withdraw",
			userID: 1,
			req: WithdrawRequest{
				AmountCents: 10000,
				Method:      "wechat",
				AccountInfo: "test_account",
			},
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// Has pending withdraw
				pendingWithdraw := model.Withdraw{Status: model.WithdrawStatusPending}
				pendingWithdraw.ID = 1
				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{pendingWithdraw}, int64(1), nil,
				)
			},
			expectError:   true,
			errorContains: "pending withdrawal",
		},
		{
			name:   "insufficient balance",
			userID: 1,
			req: WithdrawRequest{
				AmountCents: 100000, // 1000元
				Method:      "wechat",
				AccountInfo: "test_account",
			},
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// No pending withdraws
				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{}, int64(0), nil,
				)

				// Low balance
				withdraws.On("GetPlayerBalance", mock.Anything, uint64(1)).Return(&withdrawrepo.PlayerBalance{
					TotalEarnings:    50000,
					AvailableBalance: 50000, // Only 500元 available
					PendingBalance:   0,
					WithdrawTotal:    0,
				}, nil)

				// Mock order queries
				orders.On("List", mock.Anything, mock.AnythingOfType("interfaces.OrderListOptions")).Return(
					[]model.Order{}, int64(0), nil,
				)
			},
			expectError:   true,
			errorContains: "insufficient balance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockEarningsPlayerRepository{}
			orders := &MockOrderQuery{}
			withdraws := &MockWithdrawRepository{}

			tt.setupMock(players, orders, withdraws)

			svc := NewEarningsService(players, orders, withdraws)
			resp, err := svc.RequestWithdraw(context.Background(), tt.userID, tt.req)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotZero(t, resp.WithdrawID)
			}
		})
	}
}

func TestEarningsService_GetWithdrawHistory(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint64
		page        int
		pageSize    int
		setupMock   func(*MockEarningsPlayerRepository, *MockOrderQuery, *MockWithdrawRepository)
		expectError bool
		expectCount int
	}{
		{
			name:     "successful get withdraw history",
			userID:   1,
			page:     1,
			pageSize: 10,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				// Return withdraw history
				w1 := model.Withdraw{
					PlayerID:    1,
					AmountCents: 10000,
					Method:      model.WithdrawMethodWeChat,
					Status:      model.WithdrawStatusCompleted,
				}
				w1.ID = 1
				w2 := model.Withdraw{
					PlayerID:    1,
					AmountCents: 20000,
					Method:      model.WithdrawMethodAlipay,
					Status:      model.WithdrawStatusPending,
				}
				w2.ID = 2
				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{w1, w2}, int64(2), nil,
				)
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name:     "empty withdraw history",
			userID:   1,
			page:     1,
			pageSize: 10,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{}, int64(0), nil,
				)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:     "player not found",
			userID:   999,
			page:     1,
			pageSize: 10,
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{}, int64(0), nil)
			},
			expectError: true,
			expectCount: 0,
		},
		{
			name:     "default pagination values",
			userID:   1,
			page:     0, // Should default to 1
			pageSize: 0, // Should default to 20
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{}, int64(0), nil,
				)
			},
			expectError: false,
			expectCount: 0,
		},
		{
			name:     "page size exceeds max",
			userID:   1,
			page:     1,
			pageSize: 200, // Should cap at 100
			setupMock: func(players *MockEarningsPlayerRepository, orders *MockOrderQuery, withdraws *MockWithdrawRepository) {
				player := model.Player{UserID: 1}
				player.ID = 1
				players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

				withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
					[]model.Withdraw{}, int64(0), nil,
				)
			},
			expectError: false,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			players := &MockEarningsPlayerRepository{}
			orders := &MockOrderQuery{}
			withdraws := &MockWithdrawRepository{}

			tt.setupMock(players, orders, withdraws)

			svc := NewEarningsService(players, orders, withdraws)
			resp, err := svc.GetWithdrawHistory(context.Background(), tt.userID, tt.page, tt.pageSize)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Len(t, resp.Records, tt.expectCount)
			}
		})
	}
}

func TestWithdrawRequest_Structure(t *testing.T) {
	req := WithdrawRequest{
		AmountCents: 10000,
		Method:      "wechat",
		AccountInfo: "test_account_info",
	}

	assert.Equal(t, int64(10000), req.AmountCents)
	assert.Equal(t, "wechat", req.Method)
	assert.Equal(t, "test_account_info", req.AccountInfo)
}

func TestWithdrawResponse_Structure(t *testing.T) {
	resp := WithdrawResponse{
		WithdrawID: 123,
		Status:     "pending",
	}

	assert.Equal(t, uint64(123), resp.WithdrawID)
	assert.Equal(t, "pending", resp.Status)
}

func TestWithdrawRecordDTO_Structure(t *testing.T) {
	now := time.Now()
	record := WithdrawRecordDTO{
		ID:          1,
		AmountCents: 10000,
		Method:      "wechat",
		Status:      "completed",
		CreatedAt:   now,
		ProcessedAt: &now,
	}

	assert.Equal(t, uint64(1), record.ID)
	assert.Equal(t, int64(10000), record.AmountCents)
	assert.Equal(t, "wechat", record.Method)
	assert.Equal(t, "completed", record.Status)
	assert.Equal(t, now, record.CreatedAt)
	assert.NotNil(t, record.ProcessedAt)
}

func TestWithdrawHistoryResponse_Structure(t *testing.T) {
	resp := WithdrawHistoryResponse{
		Records: []WithdrawRecordDTO{
			{ID: 1, AmountCents: 10000},
			{ID: 2, AmountCents: 20000},
		},
		Total: 2,
	}

	assert.Len(t, resp.Records, 2)
	assert.Equal(t, int64(2), resp.Total)
}

func TestEarningsService_GetWithdrawHistory_RepositoryError(t *testing.T) {
	players := &MockEarningsPlayerRepository{}
	orders := &MockOrderQuery{}
	withdraws := &MockWithdrawRepository{}

	player := model.Player{UserID: 1}
	player.ID = 1
	players.On("ListPaged", mock.Anything, 1, 100).Return([]model.Player{player}, int64(1), nil)

	// Return error from withdraw list
	withdraws.On("List", mock.Anything, mock.AnythingOfType("withdraw.WithdrawListOptions")).Return(
		nil, int64(0), errors.New("database error"),
	)

	svc := NewEarningsService(players, orders, withdraws)
	resp, err := svc.GetWithdrawHistory(context.Background(), 1, 1, 10)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestStrPtr(t *testing.T) {
	s := "test"
	ptr := strPtr(s)

	assert.NotNil(t, ptr)
	assert.Equal(t, "test", *ptr)
}

func TestEarningsSummaryResponse_Structure(t *testing.T) {
	summary := EarningsSummaryResponse{
		TodayEarnings:    10000,
		MonthEarnings:    100000,
		TotalEarnings:    1000000,
		AvailableBalance: 800000,
		PendingBalance:   200000,
		WithdrawTotal:    500000,
	}

	assert.Equal(t, int64(10000), summary.TodayEarnings)
	assert.Equal(t, int64(100000), summary.MonthEarnings)
	assert.Equal(t, int64(1000000), summary.TotalEarnings)
	assert.Equal(t, int64(800000), summary.AvailableBalance)
	assert.Equal(t, int64(200000), summary.PendingBalance)
	assert.Equal(t, int64(500000), summary.WithdrawTotal)
}

func TestEarningsTrendResponse_Structure(t *testing.T) {
	trend := EarningsTrendResponse{
		Trend: []DailyEarningDTO{
			{Date: "2025-01-01", Earnings: 10000, OrderCount: 5},
			{Date: "2025-01-02", Earnings: 20000, OrderCount: 10},
		},
	}

	assert.Len(t, trend.Trend, 2)
	assert.Equal(t, "2025-01-01", trend.Trend[0].Date)
	assert.Equal(t, int64(10000), trend.Trend[0].Earnings)
	assert.Equal(t, 5, trend.Trend[0].OrderCount)
}

func TestDailyEarningDTO_Structure(t *testing.T) {
	daily := DailyEarningDTO{
		Date:       "2025-01-15",
		Earnings:   50000,
		OrderCount: 25,
	}

	assert.Equal(t, "2025-01-15", daily.Date)
	assert.Equal(t, int64(50000), daily.Earnings)
	assert.Equal(t, 25, daily.OrderCount)
}

package withdraw

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/settlementcompany"
	"gamelink/internal/repository/withdraw"
	"gamelink/pkg/apierr"
)

// Test helpers

func newTestWithdraw(id uint64, playerID uint64, amountCents int64, status model.WithdrawStatus) *model.Withdraw {
	return &model.Withdraw{
		ID:          id,
		PlayerID:    playerID,
		UserID:      playerID,
		AmountCents: amountCents,
		Method:      model.WithdrawMethodBank,
		AccountInfo: "6222021234567890123",
		Status:      status,
	}
}

func newTestSettlementCompany(id uint64, name string, status model.CompanyStatus) *model.SettlementCompany {
	return &model.SettlementCompany{
		Base:        model.Base{ID: id},
		Name:        name,
		CreditCode:  "91110000MA1234567X",
		BankName:    "Test Bank",
		BankAccount: "1234567890",
		Status:      status,
	}
}

// Mock implementations using testify/mock

type MockWithdrawRepository struct {
	mock.Mock
}

func (m *MockWithdrawRepository) Create(ctx context.Context, w *model.Withdraw) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWithdrawRepository) Get(ctx context.Context, id uint64) (*model.Withdraw, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Withdraw), args.Error(1)
}

func (m *MockWithdrawRepository) Update(ctx context.Context, w *model.Withdraw) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *MockWithdrawRepository) List(ctx context.Context, opts withdraw.WithdrawListOptions) ([]model.Withdraw, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Withdraw), args.Get(1).(int64), args.Error(2)
}

func (m *MockWithdrawRepository) GetPlayerBalance(ctx context.Context, playerID uint64) (*withdraw.PlayerBalance, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*withdraw.PlayerBalance), args.Error(1)
}

func (m *MockWithdrawRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Withdraw, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.Withdraw), args.Error(1)
}

func (m *MockWithdrawRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.WithdrawStatus, processedBy *uint64, processedAt *time.Time, reason string) ([]uint64, []withdraw.BatchOperationError, error) {
	args := m.Called(ctx, ids, status, processedBy, processedAt, reason)
	return args.Get(0).([]uint64), args.Get(1).([]withdraw.BatchOperationError), args.Error(2)
}

func (m *MockWithdrawRepository) BatchComplete(ctx context.Context, ids []uint64, adminUserID uint64, completedAt time.Time) ([]uint64, []withdraw.BatchOperationError, error) {
	args := m.Called(ctx, ids, adminUserID, completedAt)
	return args.Get(0).([]uint64), args.Get(1).([]withdraw.BatchOperationError), args.Error(2)
}

func (m *MockWithdrawRepository) ListByCompany(ctx context.Context, opts withdraw.WithdrawByCompanyOptions) ([]model.Withdraw, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
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

type MockSettlementCompanyRepository struct {
	mock.Mock
}

func (m *MockSettlementCompanyRepository) Create(ctx context.Context, company *model.SettlementCompany) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) Get(ctx context.Context, id uint64) (*model.SettlementCompany, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SettlementCompany), args.Error(1)
}

func (m *MockSettlementCompanyRepository) GetByCreditCode(ctx context.Context, creditCode string) (*model.SettlementCompany, error) {
	args := m.Called(ctx, creditCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.SettlementCompany), args.Error(1)
}

func (m *MockSettlementCompanyRepository) Update(ctx context.Context, company *model.SettlementCompany) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) List(ctx context.Context, opts settlementcompany.ListOptions) ([]model.SettlementCompany, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.SettlementCompany), args.Get(1).(int64), args.Error(2)
}

func (m *MockSettlementCompanyRepository) ToggleStatus(ctx context.Context, id uint64, status model.CompanyStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.CompanyStatus) error {
	args := m.Called(ctx, ids, status)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) BatchDelete(ctx context.Context, ids []uint64) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) GetByIDsWithPlayerCount(ctx context.Context, ids []uint64) ([]model.SettlementCompany, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.SettlementCompany), args.Error(1)
}

func (m *MockSettlementCompanyRepository) AssignPlayer(ctx context.Context, assignment *model.PlayerCompanyAssignment) error {
	args := m.Called(ctx, assignment)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) GetCurrentAssignment(ctx context.Context, playerID uint64) (*model.PlayerCompanyAssignment, error) {
	args := m.Called(ctx, playerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PlayerCompanyAssignment), args.Error(1)
}

func (m *MockSettlementCompanyRepository) GetAssignmentHistory(ctx context.Context, playerID uint64) ([]model.PlayerCompanyAssignment, error) {
	args := m.Called(ctx, playerID)
	return args.Get(0).([]model.PlayerCompanyAssignment), args.Error(1)
}

func (m *MockSettlementCompanyRepository) EndCurrentAssignment(ctx context.Context, playerID uint64, endDate time.Time) error {
	args := m.Called(ctx, playerID, endDate)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) BatchAssignPlayers(ctx context.Context, assignments []model.PlayerCompanyAssignment) error {
	args := m.Called(ctx, assignments)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) CreateHistory(ctx context.Context, history *model.SettlementCompanyHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) GetHistory(ctx context.Context, companyID uint64) ([]model.SettlementCompanyHistory, error) {
	args := m.Called(ctx, companyID)
	return args.Get(0).([]model.SettlementCompanyHistory), args.Error(1)
}

func (m *MockSettlementCompanyRepository) GetPlayerCount(ctx context.Context, companyID uint64) (int, error) {
	args := m.Called(ctx, companyID)
	return args.Int(0), args.Error(1)
}

func (m *MockSettlementCompanyRepository) UpdatePlayerCount(ctx context.Context, companyID uint64) error {
	args := m.Called(ctx, companyID)
	return args.Error(0)
}

// ============================================================================
// WithdrawRoutingService Tests
// ============================================================================

// TestNewWithdrawRoutingService tests the constructor
func TestNewWithdrawRoutingService(t *testing.T) {
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)

	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	assert.NotNil(t, service)
	assert.Equal(t, withdrawRepo, service.withdrawRepo)
	assert.Equal(t, settlementRepo, service.settlementRepo)
}

// TestRouteWithdrawal_Success tests successful withdrawal routing
func TestRouteWithdrawal_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	playerID := uint64(123)
	companyID := uint64(1)
	withdraw := newTestWithdraw(1, playerID, 10000, model.WithdrawStatusPending)
	company := newTestSettlementCompany(companyID, "Test Company", model.CompanyStatusActive)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: companyID}, nil)

	settlementRepo.On("Get", ctx, companyID).Return(company, nil)

	result, err := service.RouteWithdrawal(ctx, withdraw)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, companyID, *withdraw.SettlementCompanyID)
	assert.Equal(t, company.Name, withdraw.SettlementCompanyName)
	assert.Equal(t, company.BankAccount, withdraw.PaymentBankAccount)

	settlementRepo.AssertExpectations(t)
}

// TestRouteWithdrawal_NoSettlementCompany tests error when player has no settlement company
func TestRouteWithdrawal_NoSettlementCompany(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	playerID := uint64(123)
	withdraw := newTestWithdraw(1, playerID, 10000, model.WithdrawStatusPending)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		nil, repository.ErrNotFound)

	result, err := service.RouteWithdrawal(ctx, withdraw)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, apierr.IsNotFound(err) || errors.Is(err, ErrNoSettlementCompany))

	settlementRepo.AssertExpectations(t)
}

// TestRouteWithdrawal_CompanyInactive tests error when settlement company is inactive
func TestRouteWithdrawal_CompanyInactive(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	playerID := uint64(123)
	companyID := uint64(1)
	withdraw := newTestWithdraw(1, playerID, 10000, model.WithdrawStatusPending)
	company := newTestSettlementCompany(companyID, "Inactive Company", model.CompanyStatusInactive)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: companyID}, nil)

	settlementRepo.On("Get", ctx, companyID).Return(company, nil)

	result, err := service.RouteWithdrawal(ctx, withdraw)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, ErrCompanyInactive))

	settlementRepo.AssertExpectations(t)
}

// TestRouteWithdrawal_CompanyNotFound tests error when settlement company is not found
func TestRouteWithdrawal_CompanyNotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	playerID := uint64(123)
	companyID := uint64(999)
	withdraw := newTestWithdraw(1, playerID, 10000, model.WithdrawStatusPending)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: companyID}, nil)

	settlementRepo.On("Get", ctx, companyID).Return(nil, repository.ErrNotFound)

	result, err := service.RouteWithdrawal(ctx, withdraw)

	assert.Error(t, err)
	assert.Nil(t, result)

	settlementRepo.AssertExpectations(t)
}

// TestProcessWithdrawalRouting_Success tests successful processing of withdrawal routing
func TestProcessWithdrawalRouting_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(1)
	playerID := uint64(123)
	companyID := uint64(1)
	taxDeductedCents := int64(500)

	withdraw := newTestWithdraw(withdrawID, playerID, 10000, model.WithdrawStatusApproved)
	company := newTestSettlementCompany(companyID, "Test Company", model.CompanyStatusActive)

	withdrawRepo.On("Get", ctx, withdrawID).Return(withdraw, nil)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: companyID}, nil)

	settlementRepo.On("Get", ctx, companyID).Return(company, nil)

	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	record, err := service.ProcessWithdrawalRouting(ctx, withdrawID, taxDeductedCents)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, withdrawID, record.WithdrawID)
	assert.Equal(t, playerID, record.PlayerID)
	assert.Equal(t, companyID, record.SettlementCompanyID)
	assert.Equal(t, company.Name, record.SettlementCompanyName)
	assert.Equal(t, int64(10000), record.AmountCents)
	assert.Equal(t, taxDeductedCents, record.TaxDeductedCents)
	assert.Equal(t, int64(9500), record.ActualAmountCents)

	withdrawRepo.AssertExpectations(t)
	settlementRepo.AssertExpectations(t)
}

// TestProcessWithdrawalRouting_WithdrawalNotFound tests error when withdrawal not found
func TestProcessWithdrawalRouting_WithdrawalNotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(999)

	withdrawRepo.On("Get", ctx, withdrawID).Return(nil, repository.ErrNotFound)

	record, err := service.ProcessWithdrawalRouting(ctx, withdrawID, 500)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.True(t, errors.Is(err, ErrNotFound))

	withdrawRepo.AssertExpectations(t)
}

// TestProcessWithdrawalRouting_InvalidStatus tests error when withdrawal is not approved
func TestProcessWithdrawalRouting_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(1)
	withdraw := newTestWithdraw(withdrawID, 123, 10000, model.WithdrawStatusPending)

	withdrawRepo.On("Get", ctx, withdrawID).Return(withdraw, nil)

	record, err := service.ProcessWithdrawalRouting(ctx, withdrawID, 500)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "must be approved")

	withdrawRepo.AssertExpectations(t)
}

// TestProcessWithdrawalRouting_WithExistingRouting tests processing when withdrawal already has routing
func TestProcessWithdrawalRouting_WithExistingRouting(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(1)
	playerID := uint64(123)
	companyID := uint64(1)
	taxDeductedCents := int64(500)

	withdraw := newTestWithdraw(withdrawID, playerID, 10000, model.WithdrawStatusApproved)
	withdraw.SettlementCompanyID = &companyID
	withdraw.SettlementCompanyName = "Test Company"

	withdrawRepo.On("Get", ctx, withdrawID).Return(withdraw, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	record, err := service.ProcessWithdrawalRouting(ctx, withdrawID, taxDeductedCents)

	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, taxDeductedCents, withdraw.TaxDeductedCents)

	withdrawRepo.AssertExpectations(t)
}

// TestProcessWithdrawalRouting_UpdateFailed tests error when update fails
func TestProcessWithdrawalRouting_UpdateFailed(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(1)
	withdraw := newTestWithdraw(withdrawID, 123, 10000, model.WithdrawStatusApproved)

	withdrawRepo.On("Get", ctx, withdrawID).Return(withdraw, nil)

	settlementRepo.On("GetCurrentAssignment", ctx, mock.Anything).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: 1}, nil)

	settlementRepo.On("Get", ctx, mock.Anything).Return(
		newTestSettlementCompany(1, "Test Company", model.CompanyStatusActive), nil)

	withdrawRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

	record, err := service.ProcessWithdrawalRouting(ctx, withdrawID, 500)

	assert.Error(t, err)
	assert.Nil(t, record)
	assert.Contains(t, err.Error(), "failed to update withdrawal")
}

// TestCompleteWithdrawalPayment_Success tests successful completion of withdrawal payment
func TestCompleteWithdrawalPayment_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawID := uint64(1)
	bankTransactionNo := "BANK_TXN_12345"
	withdraw := newTestWithdraw(withdrawID, 123, 10000, model.WithdrawStatusApproved)

	withdrawRepo.On("Get", ctx, withdrawID).Return(withdraw, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := service.CompleteWithdrawalPayment(ctx, withdrawID, bankTransactionNo)

	assert.NoError(t, err)
	assert.Equal(t, bankTransactionNo, withdraw.BankTransactionNo)
	assert.Equal(t, model.WithdrawStatusCompleted, withdraw.Status)
	assert.NotNil(t, withdraw.PaidAt)
	assert.NotNil(t, withdraw.CompletedAt)

	withdrawRepo.AssertExpectations(t)
}

// TestCompleteWithdrawalPayment_NotFound tests error when withdrawal not found
func TestCompleteWithdrawalPayment_NotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(nil, repository.ErrNotFound)

	err := service.CompleteWithdrawalPayment(ctx, 999, "BANK_TXN_123")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

// TestCompleteWithdrawalPayment_UpdateFailed tests error when update fails
func TestCompleteWithdrawalPayment_UpdateFailed(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdraw := newTestWithdraw(1, 123, 10000, model.WithdrawStatusApproved)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(withdraw, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

	err := service.CompleteWithdrawalPayment(ctx, 1, "BANK_TXN_123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update withdrawal")
}

// TestGetWithdrawal_Success tests successful retrieval of withdrawal
func TestGetWithdrawal_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	expectedWithdraw := newTestWithdraw(1, 123, 10000, model.WithdrawStatusPending)

	withdrawRepo.On("Get", ctx, uint64(1)).Return(expectedWithdraw, nil)

	result, err := service.GetWithdrawal(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedWithdraw, result)

	withdrawRepo.AssertExpectations(t)
}

// TestGetWithdrawal_NotFound tests error when withdrawal not found
func TestGetWithdrawal_NotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(nil, repository.ErrNotFound)

	result, err := service.GetWithdrawal(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, ErrNotFound))
}

// TestGetWithdrawal_InternalError tests error when repository returns internal error
func TestGetWithdrawal_InternalError(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(nil, errors.New("internal database error"))

	result, err := service.GetWithdrawal(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get withdrawal")
}

// TestGetPlayerCurrentCompany_Success tests successful retrieval of player's current company
func TestGetPlayerCurrentCompany_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	playerID := uint64(123)
	companyID := uint64(1)
	company := newTestSettlementCompany(companyID, "Test Company", model.CompanyStatusActive)

	settlementRepo.On("GetCurrentAssignment", ctx, playerID).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: companyID}, nil)

	settlementRepo.On("Get", ctx, companyID).Return(company, nil)

	result, err := service.GetPlayerCurrentCompany(ctx, playerID)

	assert.NoError(t, err)
	assert.Equal(t, company, result)

	settlementRepo.AssertExpectations(t)
}

// TestGetPlayerCurrentCompany_NoAssignment tests error when player has no assignment
func TestGetPlayerCurrentCompany_NoAssignment(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	settlementRepo.On("GetCurrentAssignment", ctx, mock.Anything).Return(
		nil, repository.ErrNotFound)

	result, err := service.GetPlayerCurrentCompany(ctx, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, ErrNoSettlementCompany))
}

// TestGetPlayerCurrentCompany_CompanyNotFound tests error when company not found
func TestGetPlayerCurrentCompany_CompanyNotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	settlementRepo.On("GetCurrentAssignment", ctx, mock.Anything).Return(
		&model.PlayerCompanyAssignment{SettlementCompanyID: 999}, nil)

	settlementRepo.On("Get", ctx, mock.Anything).Return(nil, errors.New("database error"))

	result, err := service.GetPlayerCurrentCompany(ctx, 123)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get settlement company")
}

// ============================================================================
// WithdrawRoutingStatsService Tests
// ============================================================================

// TestNewWithdrawRoutingStatsService tests the constructor
func TestNewWithdrawRoutingStatsService(t *testing.T) {
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	assert.NotNil(t, service)
	assert.Equal(t, withdrawRepo, service.withdrawRepo)
	assert.Equal(t, settlementRepo, service.settlementRepo)
}

// TestGetRoutingStats_Success tests successful retrieval of routing stats
func TestGetRoutingStats_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	dateFrom := time.Now().AddDate(-1, 0, 0)
	dateTo := time.Now()

	expectedStats := &model.WithdrawRoutingStatsResponse{
		TotalWithdrawals:       100,
		TotalAmountCents:       1000000,
		TotalTaxDeductedCents:  50000,
		TotalActualAmountCents: 950000,
		ByCompany: []model.WithdrawRoutingStats{
			{
				SettlementCompanyID:    1,
				SettlementCompanyName:  "Company A",
				TotalWithdrawals:       60,
				TotalAmountCents:       600000,
				TotalTaxDeductedCents:  30000,
				TotalActualAmountCents: 570000,
			},
		},
	}

	withdrawRepo.On("GetRoutingStats", ctx, &dateFrom, &dateTo).Return(expectedStats, nil)

	result, err := service.GetRoutingStats(ctx, &model.WithdrawRoutingStatsRequest{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedStats, result)

	withdrawRepo.AssertExpectations(t)
}

// TestGetRoutingStats_Error tests error when repository fails
func TestGetRoutingStats_Error(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	withdrawRepo.On("GetRoutingStats", ctx, mock.Anything, mock.Anything).Return(
		nil, errors.New("database error"))

	result, err := service.GetRoutingStats(ctx, &model.WithdrawRoutingStatsRequest{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get routing stats")
}

// TestListWithdrawalsByCompany_Success tests successful listing of withdrawals by company
func TestListWithdrawalsByCompany_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	companyID := uint64(1)
	dateFrom := time.Now().AddDate(-1, 0, 0)
	dateTo := time.Now()

	withdrawals := []model.Withdraw{
		*newTestWithdraw(1, 123, 10000, model.WithdrawStatusCompleted),
		*newTestWithdraw(2, 124, 20000, model.WithdrawStatusCompleted),
	}

	withdrawRepo.On("ListByCompany", ctx, mock.Anything).Return(withdrawals, int64(2), nil)

	req := &model.ListWithdrawsByCompanyRequest{
		SettlementCompanyID: &companyID,
		DateFrom:            &dateFrom,
		DateTo:              &dateTo,
		Page:                1,
		PageSize:            20,
	}

	result, err := service.ListWithdrawalsByCompany(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
	assert.Equal(t, withdrawals, result.Withdraws)
}

// TestListWithdrawalsByCompany_Error tests error when repository fails
func TestListWithdrawalsByCompany_Error(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	withdrawRepo.On("ListByCompany", ctx, mock.Anything).Return(
		nil, int64(0), errors.New("database error"))

	_, err := service.ListWithdrawalsByCompany(ctx, &model.ListWithdrawsByCompanyRequest{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list withdrawals by company")
}

// TestGenerateRoutingReport_Monthly tests monthly report generation
func TestGenerateRoutingReport_Monthly(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	withdrawRepo.On("GetRoutingStats", ctx, mock.Anything, mock.Anything).Return(
		&model.WithdrawRoutingStatsResponse{
			TotalWithdrawals:       10,
			TotalAmountCents:       100000,
			TotalTaxDeductedCents:  5000,
			TotalActualAmountCents: 95000,
			ByCompany:              []model.WithdrawRoutingStats{},
		}, nil)

	req := &model.WithdrawRoutingReportRequest{
		ReportType: "monthly",
		Year:       2025,
		Month:      1,
	}

	result, err := service.GenerateRoutingReport(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "monthly", result.ReportType)
	assert.Equal(t, 2025, result.Year)
	assert.Equal(t, 1, result.Month)
	assert.Equal(t, int64(10), result.TotalWithdrawals)
}

// TestGenerateRoutingReport_InvalidReportType tests error for invalid report type
func TestGenerateRoutingReport_InvalidReportType(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	req := &model.WithdrawRoutingReportRequest{
		ReportType: "invalid",
		Year:       2025,
	}

	result, err := service.GenerateRoutingReport(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid report type")
}

// TestGenerateRoutingReport_InvalidMonth tests error for invalid month
func TestGenerateRoutingReport_InvalidMonth(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	req := &model.WithdrawRoutingReportRequest{
		ReportType: "monthly",
		Year:       2025,
		Month:      13, // Invalid
	}

	result, err := service.GenerateRoutingReport(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid month")
}

// TestGenerateRoutingReport_InvalidQuarter tests error for invalid quarter
func TestGenerateRoutingReport_InvalidQuarter(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	req := &model.WithdrawRoutingReportRequest{
		ReportType: "quarterly",
		Year:       2025,
		Quarter:    5, // Invalid
	}

	result, err := service.GenerateRoutingReport(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid quarter")
}

// TestGetCompanyWithdrawalStats_Success tests successful retrieval of company withdrawal stats
func TestGetCompanyWithdrawalStats_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	companyID := uint64(1)
	dateFrom := time.Now().AddDate(-1, 0, 0)
	dateTo := time.Now()

	expectedStats := &model.WithdrawRoutingStats{
		SettlementCompanyID:    companyID,
		SettlementCompanyName:  "Test Company",
		TotalWithdrawals:       50,
		TotalAmountCents:       500000,
		TotalTaxDeductedCents:  25000,
		TotalActualAmountCents: 475000,
		AverageAmountCents:     10000,
		Percentage:             100.0,
	}

	settlementRepo.On("Get", ctx, companyID).Return(
		newTestSettlementCompany(companyID, "Test Company", model.CompanyStatusActive), nil)

	withdrawRepo.On("GetRoutingStatsByCompany", ctx, companyID, &dateFrom, &dateTo).Return(
		expectedStats, nil)

	result, err := service.GetCompanyWithdrawalStats(ctx, companyID, &dateFrom, &dateTo)

	assert.NoError(t, err)
	assert.Equal(t, expectedStats, result)
}

// TestGetCompanyWithdrawalStats_CompanyNotFound tests error when company not found
func TestGetCompanyWithdrawalStats_CompanyNotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingStatsService(withdrawRepo, settlementRepo)

	settlementRepo.On("Get", ctx, mock.Anything).Return(nil, repository.ErrNotFound)

	result, err := service.GetCompanyWithdrawalStats(ctx, 999, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "settlement company not found")
}

// ============================================================================
// Batch Operations Tests
// ============================================================================

// TestBatchApprove_Success tests successful batch approval
func TestBatchApprove_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	adminUserID := uint64(1)
	remark := "Approved after review"

	withdraw1 := newTestWithdraw(1, 123, 10000, model.WithdrawStatusPending)
	withdraw2 := newTestWithdraw(2, 124, 20000, model.WithdrawStatusPending)

	withdrawRepo.On("Get", ctx, uint64(1)).Return(withdraw1, nil)
	withdrawRepo.On("Get", ctx, uint64(2)).Return(withdraw2, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	req := &BatchApproveRequest{
		WithdrawIDs: []uint64{1, 2},
		ProcessedBy: adminUserID,
		Remark:      remark,
	}

	result, err := service.BatchApprove(ctx, req, adminUserID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 2)
}

// TestBatchApprove_EmptyIDs tests error when withdrawal IDs are empty
func TestBatchApprove_EmptyIDs(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	req := &BatchApproveRequest{
		WithdrawIDs: []uint64{},
		ProcessedBy: 1,
	}

	result, err := service.BatchApprove(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestBatchApprove_TooManyIDs tests error when too many withdrawal IDs provided
func TestBatchApprove_TooManyIDs(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	ids := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		ids[i] = uint64(i + 1)
	}

	req := &BatchApproveRequest{
		WithdrawIDs: ids,
		ProcessedBy: 1,
	}

	result, err := service.BatchApprove(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 100 withdrawals")
}

// TestBatchApprove_NotFound tests handling of not found withdrawals
func TestBatchApprove_NotFound(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(nil, repository.ErrNotFound)

	req := &BatchApproveRequest{
		WithdrawIDs: []uint64{999},
		ProcessedBy: 1,
	}

	result, err := service.BatchApprove(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)
	assert.Contains(t, result.FailedItems[0].Message, "not found")
}

// TestBatchApprove_InvalidStatus tests handling of invalid status withdrawals
func TestBatchApprove_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(
		newTestWithdraw(1, 123, 10000, model.WithdrawStatusApproved), nil)

	req := &BatchApproveRequest{
		WithdrawIDs: []uint64{1},
		ProcessedBy: 1,
	}

	result, err := service.BatchApprove(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedItems[0].Message, "cannot approve withdrawal with status")
}

// TestBatchApprove_UpdateFailed tests handling of update failures
func TestBatchApprove_UpdateFailed(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(
		newTestWithdraw(1, 123, 10000, model.WithdrawStatusPending), nil)

	withdrawRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

	req := &BatchApproveRequest{
		WithdrawIDs: []uint64{1},
		ProcessedBy: 1,
	}

	result, err := service.BatchApprove(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedItems[0].Message, "update failed")
}

// TestBatchReject_Success tests successful batch rejection
func TestBatchReject_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	adminUserID := uint64(1)
	reason := "Invalid account information"

	withdraw1 := newTestWithdraw(1, 123, 10000, model.WithdrawStatusPending)
	withdraw2 := newTestWithdraw(2, 124, 20000, model.WithdrawStatusPending)

	withdrawRepo.On("Get", ctx, uint64(1)).Return(withdraw1, nil)
	withdrawRepo.On("Get", ctx, uint64(2)).Return(withdraw2, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	req := &BatchRejectRequest{
		WithdrawIDs: []uint64{1, 2},
		ProcessedBy: adminUserID,
		Reason:      reason,
	}

	result, err := service.BatchReject(ctx, req, adminUserID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 2)
}

// TestBatchReject_EmptyIDs tests error when withdrawal IDs are empty
func TestBatchReject_EmptyIDs(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	req := &BatchRejectRequest{
		WithdrawIDs: []uint64{},
		ProcessedBy: 1,
		Reason:      "Test",
	}

	result, err := service.BatchReject(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestBatchReject_InvalidStatus tests handling of invalid status withdrawals
func TestBatchReject_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(
		newTestWithdraw(1, 123, 10000, model.WithdrawStatusCompleted), nil)

	req := &BatchRejectRequest{
		WithdrawIDs: []uint64{1},
		ProcessedBy: 1,
		Reason:      "Test",
	}

	result, err := service.BatchReject(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedItems[0].Message, "cannot reject withdrawal with status")
}

// TestBatchComplete_Success tests successful batch completion
func TestBatchComplete_Success(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	adminUserID := uint64(1)

	withdraw1 := newTestWithdraw(1, 123, 10000, model.WithdrawStatusApproved)
	withdraw2 := newTestWithdraw(2, 124, 20000, model.WithdrawStatusApproved)

	withdrawRepo.On("Get", ctx, uint64(1)).Return(withdraw1, nil)
	withdrawRepo.On("Get", ctx, uint64(2)).Return(withdraw2, nil)
	withdrawRepo.On("Update", ctx, mock.Anything).Return(nil)

	req := &BatchCompleteRequest{
		WithdrawIDs: []uint64{1, 2},
	}

	result, err := service.BatchComplete(ctx, req, adminUserID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessIDs, 2)
}

// TestBatchComplete_EmptyIDs tests error when withdrawal IDs are empty
func TestBatchComplete_EmptyIDs(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	req := &BatchCompleteRequest{
		WithdrawIDs: []uint64{},
	}

	result, err := service.BatchComplete(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "withdrawal IDs are required")
}

// TestBatchComplete_InvalidStatus tests handling of invalid status withdrawals
func TestBatchComplete_InvalidStatus(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	withdrawRepo.On("Get", ctx, mock.Anything).Return(
		newTestWithdraw(1, 123, 10000, model.WithdrawStatusPending), nil)

	req := &BatchCompleteRequest{
		WithdrawIDs: []uint64{1},
	}

	result, err := service.BatchComplete(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.FailedItems[0].Message, "cannot complete withdrawal with status")
}

// TestBatchComplete_WithProcessedBy tests setting ProcessedBy when not set
func TestBatchComplete_WithProcessedBy(t *testing.T) {
	ctx := context.Background()
	withdrawRepo := new(MockWithdrawRepository)
	settlementRepo := new(MockSettlementCompanyRepository)
	service := NewWithdrawRoutingService(withdrawRepo, settlementRepo)

	adminUserID := uint64(1)

	withdraw := newTestWithdraw(1, 123, 10000, model.WithdrawStatusApproved)
	withdraw.ProcessedBy = nil // Not set

	withdrawRepo.On("Get", ctx, mock.Anything).Return(withdraw, nil)
	withdrawRepo.On("Update", ctx, mock.MatchedBy(func(w *model.Withdraw) bool {
		return w.ProcessedBy != nil && *w.ProcessedBy == adminUserID
	})).Return(nil)

	req := &BatchCompleteRequest{
		WithdrawIDs: []uint64{1},
	}

	result, err := service.BatchComplete(ctx, req, adminUserID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.SuccessCount)
	assert.NotNil(t, withdraw.ProcessedBy)
	assert.Equal(t, adminUserID, *withdraw.ProcessedBy)
}

// ============================================================================
// Test Withdraw Model Methods
// ============================================================================

// TestWithdraw_CalculateActualAmount tests the actual amount calculation
func TestWithdraw_CalculateActualAmount(t *testing.T) {
	tests := []struct {
		name             string
		amountCents      int64
		taxDeductedCents int64
		expected         int64
	}{
		{
			name:             "No tax",
			amountCents:      10000,
			taxDeductedCents: 0,
			expected:         10000,
		},
		{
			name:             "With tax",
			amountCents:      10000,
			taxDeductedCents: 500,
			expected:         9500,
		},
		{
			name:             "Full tax",
			amountCents:      10000,
			taxDeductedCents: 10000,
			expected:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withdraw := &model.Withdraw{
				AmountCents:      tt.amountCents,
				TaxDeductedCents: tt.taxDeductedCents,
			}
			result := withdraw.CalculateActualAmount()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestWithdraw_SetRoutingInfo tests setting routing information
func TestWithdraw_SetRoutingInfo(t *testing.T) {
	company := &model.SettlementCompany{
		Base:        model.Base{ID: 1},
		Name:        "Test Company",
		BankAccount: "1234567890",
	}

	withdraw := &model.Withdraw{}
	withdraw.SetRoutingInfo(company)

	assert.NotNil(t, withdraw.SettlementCompanyID)
	assert.Equal(t, uint64(1), *withdraw.SettlementCompanyID)
	assert.Equal(t, "Test Company", withdraw.SettlementCompanyName)
	assert.Equal(t, "1234567890", withdraw.PaymentBankAccount)
}

// TestWithdraw_SetRoutingInfo_NilCompany tests handling of nil company
func TestWithdraw_SetRoutingInfo_NilCompany(t *testing.T) {
	withdraw := &model.Withdraw{}
	withdraw.SetRoutingInfo(nil)

	assert.Nil(t, withdraw.SettlementCompanyID)
	assert.Empty(t, withdraw.SettlementCompanyName)
	assert.Empty(t, withdraw.PaymentBankAccount)
}

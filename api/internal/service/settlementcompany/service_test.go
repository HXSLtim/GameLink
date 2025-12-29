package settlementcompany

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/settlementcompany"
)

// MockSettlementCompanyRepository is a mock implementation
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

func (m *MockSettlementCompanyRepository) List(ctx context.Context, opts settlementcompany.ListOptions) ([]model.SettlementCompany, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.SettlementCompany), args.Get(1).(int64), args.Error(2)
}

func (m *MockSettlementCompanyRepository) ToggleStatus(ctx context.Context, id uint64, status model.CompanyStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
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
	return args.Get(0).(int), args.Error(1)
}

func (m *MockSettlementCompanyRepository) UpdatePlayerCount(ctx context.Context, companyID uint64) error {
	args := m.Called(ctx, companyID)
	return args.Error(0)
}

func (m *MockSettlementCompanyRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
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
	if args.Get(0) == nil {
		return []model.SettlementCompany{}, args.Error(1)
	}
	return args.Get(0).([]model.SettlementCompany), args.Error(1)
}

// MockPlayerRepository is a mock implementation
type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockPlayerRepository) Create(_ context.Context, _ *model.Player) error {
	return nil
}

func (m *MockPlayerRepository) Update(_ context.Context, _ *model.Player) error {
	return nil
}

func (m *MockPlayerRepository) Delete(_ context.Context, _ uint64) error {
	return nil
}

func (m *MockPlayerRepository) List(_ context.Context) ([]model.Player, error) {
	return nil, nil
}

func (m *MockPlayerRepository) ListPaged(_ context.Context, _, _ int) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *MockPlayerRepository) ListPagedWithFilter(_ context.Context, _, _ int, _ string, _ *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}

func (m *MockPlayerRepository) GetByUserID(_ context.Context, _ uint64) (*model.Player, error) {
	return nil, nil
}

func (m *MockPlayerRepository) BatchUpdateStatus(_ context.Context, _ []uint64, _ model.VerificationStatus) (int64, error) {
	return 0, nil
}

func (m *MockPlayerRepository) BatchUpdateRank(_ context.Context, _ []uint64, _ string) (int64, error) {
	return 0, nil
}

func (m *MockPlayerRepository) BatchDelete(_ context.Context, _ []uint64) (int64, error) {
	return 0, nil
}

func (m *MockPlayerRepository) BatchUpdateHourlyRate(_ context.Context, _ []uint64, _ int64) (int64, error) {
	return 0, nil
}

// Test cases

func TestCreateCompany_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.CreateSettlementCompanyRequest{
		Name:       "Test Company",
		CreditCode: "91110000100000000A", // Valid 18-char credit code
		BankName:   "Test Bank",
	}

	mockRepo.On("GetByCreditCode", ctx, req.CreditCode).Return(nil, repository.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(nil)

	result, err := svc.CreateCompany(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Company", result.Name)
	assert.Equal(t, model.CompanyStatusActive, result.Status)
}

func TestCreateCompany_InvalidCreditCode(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.CreateSettlementCompanyRequest{
		Name:       "Test Company",
		CreditCode: "invalid", // Invalid credit code
	}

	result, err := svc.CreateCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid credit code")
}

func TestCreateCompany_CreditCodeExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.CreateSettlementCompanyRequest{
		Name:       "Test Company",
		CreditCode: "91110000100000000A",
	}

	existingCompany := &model.SettlementCompany{Name: "Existing"}
	mockRepo.On("GetByCreditCode", ctx, req.CreditCode).Return(existingCompany, nil)

	result, err := svc.CreateCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
}

func TestGetCompany_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base: model.Base{ID: 1},
		Name: "Test Company",
	}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)

	result, err := svc.GetCompany(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Test Company", result.Name)
}

func TestGetCompany_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetCompany(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateCompany_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base: model.Base{ID: 1},
		Name: "Old Name",
	}
	newName := "New Name"
	req := &model.UpdateSettlementCompanyRequest{Name: &newName}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(nil)

	result, err := svc.UpdateCompany(ctx, 1, req, 1)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
}

func TestUpdateCompany_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	newName := "New Name"
	req := &model.UpdateSettlementCompanyRequest{Name: &newName}
	result, err := svc.UpdateCompany(ctx, 999, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListCompanies_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	companies := []model.SettlementCompany{
		{Base: model.Base{ID: 1}, Name: "Company 1"},
		{Base: model.Base{ID: 2}, Name: "Company 2"},
	}
	mockRepo.On("List", ctx, mock.AnythingOfType("settlementcompany.ListOptions")).Return(companies, int64(2), nil)

	req := &model.ListSettlementCompaniesRequest{Page: 1, PageSize: 10}
	result, err := svc.ListCompanies(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Companies, 2)
}

func TestToggleCompanyStatus_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(nil)
	mockRepo.On("ToggleStatus", ctx, uint64(1), model.CompanyStatusInactive).Return(nil)

	err := svc.ToggleCompanyStatus(ctx, 1, false, 1)

	assert.NoError(t, err)
}

func TestToggleCompanyStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.ToggleCompanyStatus(ctx, 999, false, 1)

	assert.Error(t, err)
}

func TestGetCompanyHistory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}}
	histories := []model.SettlementCompanyHistory{
		{FieldName: "name", OldValue: "old", NewValue: "new"},
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("GetHistory", ctx, uint64(1)).Return(histories, nil)

	result, err := svc.GetCompanyHistory(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetCompanyHistory_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetCompanyHistory(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignPlayerToCompany_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}
	player := &model.Player{Base: model.Base{ID: 1}}

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return(player, nil)
	mockRepo.On("AssignPlayer", ctx, mock.AnythingOfType("*model.PlayerCompanyAssignment")).Return(nil)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(1), result.PlayerID)
}

func TestAssignPlayerToCompany_CompanyNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 999,
	}

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignPlayerToCompany_CompanyInactive(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusInactive,
	}

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "inactive")
}

func TestAssignPlayerToCompany_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            999,
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "player")
}

func TestBatchAssignPlayers_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}
	player1 := &model.Player{Base: model.Base{ID: 1}}
	player2 := &model.Player{Base: model.Base{ID: 2}}

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return(player1, nil)
	mockPlayerRepo.On("Get", ctx, uint64(2)).Return(player2, nil)
	mockRepo.On("BatchAssignPlayers", ctx, mock.AnythingOfType("[]model.PlayerCompanyAssignment")).Return(nil)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestBatchAssignPlayers_CompanyInactive(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusInactive,
	}

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestGetCurrentAssignment_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	assignment := &model.PlayerCompanyAssignment{
		PlayerID:            1,
		SettlementCompanyID: 1,
		IsCurrent:           true,
	}
	mockRepo.On("GetCurrentAssignment", ctx, uint64(1)).Return(assignment, nil)

	result, err := svc.GetCurrentAssignment(ctx, 1)

	assert.NoError(t, err)
	assert.True(t, result.IsCurrent)
}

func TestGetCurrentAssignment_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("GetCurrentAssignment", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetCurrentAssignment(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAssignmentHistory_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	player := &model.Player{Base: model.Base{ID: 1}}
	assignments := []model.PlayerCompanyAssignment{
		{PlayerID: 1, SettlementCompanyID: 1},
		{PlayerID: 1, SettlementCompanyID: 2},
	}

	mockPlayerRepo.On("Get", ctx, uint64(1)).Return(player, nil)
	mockRepo.On("GetAssignmentHistory", ctx, uint64(1)).Return(assignments, nil)

	result, err := svc.GetAssignmentHistory(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Assignments, 2)
}

func TestGetAssignmentHistory_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	mockPlayerRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetAssignmentHistory(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEndCurrentAssignment_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	assignment := &model.PlayerCompanyAssignment{PlayerID: 1, IsCurrent: true}
	endDate := time.Now()

	mockRepo.On("GetCurrentAssignment", ctx, uint64(1)).Return(assignment, nil)
	mockRepo.On("EndCurrentAssignment", ctx, uint64(1), endDate).Return(nil)

	err := svc.EndCurrentAssignment(ctx, 1, endDate)

	assert.NoError(t, err)
}

func TestEndCurrentAssignment_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("GetCurrentAssignment", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.EndCurrentAssignment(ctx, 999, time.Now())

	assert.Error(t, err)
}

func TestUpdateCompany_AllFieldChanges(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:              model.Base{ID: 1},
		Name:              "Old Name",
		TaxRegistrationNo: "OldTax",
		BankName:          "Old Bank",
		BankAccount:       "OldAccount",
		BankBranch:        "Old Branch",
		ContactName:       "Old Contact",
		ContactPhone:      "OldPhone",
		Address:           "Old Address",
	}

	newName := "New Name"
	newTax := "NewTax"
	newBank := "New Bank"
	newAccount := "NewAccount"
	newBranch := "New Branch"
	newContact := "New Contact"
	newPhone := "NewPhone"
	newAddress := "New Address"

	req := &model.UpdateSettlementCompanyRequest{
		Name:              &newName,
		TaxRegistrationNo: &newTax,
		BankName:          &newBank,
		BankAccount:       &newAccount,
		BankBranch:        &newBranch,
		ContactName:       &newContact,
		ContactPhone:      &newPhone,
		Address:           &newAddress,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(nil)

	result, err := svc.UpdateCompany(ctx, 1, req, 1)

	assert.NoError(t, err)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newTax, result.TaxRegistrationNo)
	assert.Equal(t, newBank, result.BankName)
	assert.Equal(t, newAccount, result.BankAccount)
	assert.Equal(t, newBranch, result.BankBranch)
	assert.Equal(t, newContact, result.ContactName)
	assert.Equal(t, newPhone, result.ContactPhone)
	assert.Equal(t, newAddress, result.Address)
}

func TestListCompanies_WithStatusFilter(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	companies := []model.SettlementCompany{
		{Base: model.Base{ID: 1}, Name: "Active Company", Status: model.CompanyStatusActive},
	}
	mockRepo.On("List", ctx, mock.MatchedBy(func(opts settlementcompany.ListOptions) bool {
		return opts.Status != nil && *opts.Status == model.CompanyStatusActive
	})).Return(companies, int64(1), nil)

	req := &model.ListSettlementCompaniesRequest{
		Page:     1,
		PageSize: 10,
		Status:   model.CompanyStatusActive,
	}
	result, err := svc.ListCompanies(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}

func TestListCompanies_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	companies := []model.SettlementCompany{}
	mockRepo.On("List", ctx, mock.MatchedBy(func(opts settlementcompany.ListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return(companies, int64(0), nil)

	req := &model.ListSettlementCompaniesRequest{
		Page:     0,
		PageSize: 0,
	}
	result, err := svc.ListCompanies(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
}

func TestListCompanies_PageSizeCapped(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	companies := []model.SettlementCompany{}
	mockRepo.On("List", ctx, mock.MatchedBy(func(opts settlementcompany.ListOptions) bool {
		return opts.PageSize == 20
	})).Return(companies, int64(0), nil)

	req := &model.ListSettlementCompaniesRequest{
		Page:     1,
		PageSize: 200,
	}
	result, err := svc.ListCompanies(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 20, result.PageSize)
}

func TestToggleCompanyStatus_Enable(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusInactive,
	}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(nil)
	mockRepo.On("ToggleStatus", ctx, uint64(1), model.CompanyStatusActive).Return(nil)

	err := svc.ToggleCompanyStatus(ctx, 1, true, 1)

	assert.NoError(t, err)
}

func TestToggleCompanyStatus_NoChange(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("ToggleStatus", ctx, uint64(1), model.CompanyStatusActive).Return(nil)

	err := svc.ToggleCompanyStatus(ctx, 1, true, 1)

	assert.NoError(t, err)
}

func TestBatchAssignPlayers_CompanyNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 999,
	}

	mockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestBatchAssignPlayers_PlayerNotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}
	player1 := &model.Player{Base: model.Base{ID: 1}}

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 999},
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return(player1, nil)
	mockPlayerRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.Contains(t, err.Error(), "player")
}

func TestAssignPlayerToCompany_WithoutPlayerRepo(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("AssignPlayer", ctx, mock.AnythingOfType("*model.PlayerCompanyAssignment")).Return(nil)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBatchAssignPlayers_WithoutPlayerRepo(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Base:   model.Base{ID: 1},
		Status: model.CompanyStatusActive,
	}

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("BatchAssignPlayers", ctx, mock.AnythingOfType("[]model.PlayerCompanyAssignment")).Return(nil)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestGetAssignmentHistory_WithoutPlayerRepo(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	assignments := []model.PlayerCompanyAssignment{
		{PlayerID: 1, SettlementCompanyID: 1},
	}

	mockRepo.On("GetAssignmentHistory", ctx, uint64(1)).Return(assignments, nil)

	result, err := svc.GetAssignmentHistory(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}

func TestDetectChanges_NoChanges(t *testing.T) {
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{
		Name:         "Test",
		BankName:     "Bank",
		BankAccount:  "Account",
		BankBranch:   "Branch",
		ContactName:  "Contact",
		ContactPhone: "Phone",
		Address:      "Address",
	}

	sameName := "Test"
	req := &model.UpdateSettlementCompanyRequest{
		Name: &sameName,
	}

	changes := svc.detectChanges(company, req)

	assert.Empty(t, changes)
}

// ============================================================================
// Additional Tests for Coverage Improvement
// ============================================================================

func TestCreateCompany_GetByCreditCodeError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.CreateSettlementCompanyRequest{
		Name:       "Test Company",
		CreditCode: "91110000100000000A",
	}

	mockRepo.On("GetByCreditCode", ctx, req.CreditCode).Return((*model.SettlementCompany)(nil), assert.AnError)

	result, err := svc.CreateCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateCompany_CreateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.CreateSettlementCompanyRequest{
		Name:       "Test Company",
		CreditCode: "91110000100000000A",
	}

	mockRepo.On("GetByCreditCode", ctx, req.CreditCode).Return((*model.SettlementCompany)(nil), repository.ErrNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(assert.AnError)

	result, err := svc.CreateCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetCompany_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	result, err := svc.GetCompany(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateCompany_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	newName := "New Name"
	req := &model.UpdateSettlementCompanyRequest{Name: &newName}
	result, err := svc.UpdateCompany(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateCompany_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Name: "Old"}
	newName := "New"
	req := &model.UpdateSettlementCompanyRequest{Name: &newName}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(assert.AnError)

	result, err := svc.UpdateCompany(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListCompanies_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("List", ctx, mock.Anything).Return([]model.SettlementCompany{}, int64(0), assert.AnError)

	req := &model.ListSettlementCompaniesRequest{Page: 1, PageSize: 10}
	result, err := svc.ListCompanies(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestToggleCompanyStatus_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	err := svc.ToggleCompanyStatus(ctx, 1, true, 1)

	assert.Error(t, err)
}

func TestToggleCompanyStatus_ToggleError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRepo.On("ToggleStatus", ctx, uint64(1), model.CompanyStatusInactive).Return(assert.AnError)

	err := svc.ToggleCompanyStatus(ctx, 1, false, 1)

	assert.Error(t, err)
}

func TestGetCompanyHistory_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	result, err := svc.GetCompanyHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetCompanyHistory_HistoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("GetHistory", ctx, uint64(1)).Return([]model.SettlementCompanyHistory{}, assert.AnError)

	result, err := svc.GetCompanyHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignPlayerToCompany_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignPlayerToCompany_PlayerGetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return((*model.Player)(nil), assert.AnError)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAssignPlayerToCompany_AssignError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	req := &model.AssignPlayerToCompanyRequest{
		PlayerID:            1,
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("AssignPlayer", ctx, mock.AnythingOfType("*model.PlayerCompanyAssignment")).Return(assert.AnError)

	result, err := svc.AssignPlayerToCompany(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBatchAssignPlayers_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return((*model.SettlementCompany)(nil), assert.AnError)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestBatchAssignPlayers_PlayerGetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1},
		SettlementCompanyID: 1,
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return((*model.Player)(nil), assert.AnError)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestBatchAssignPlayers_BatchAssignError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	req := &model.BatchAssignPlayersRequest{
		PlayerIDs:           []uint64{1, 2},
		SettlementCompanyID: 1,
		EffectiveDate:       time.Now(),
	}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("BatchAssignPlayers", ctx, mock.AnythingOfType("[]model.PlayerCompanyAssignment")).Return(assert.AnError)

	count, err := svc.BatchAssignPlayers(ctx, req, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

func TestGetCurrentAssignment_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("GetCurrentAssignment", ctx, uint64(1)).Return((*model.PlayerCompanyAssignment)(nil), assert.AnError)

	result, err := svc.GetCurrentAssignment(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAssignmentHistory_PlayerGetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	mockPlayerRepo.On("Get", ctx, uint64(1)).Return((*model.Player)(nil), assert.AnError)

	result, err := svc.GetAssignmentHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAssignmentHistory_HistoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	mockPlayerRepo := new(MockPlayerRepository)
	svc := NewSettlementCompanyService(mockRepo, mockPlayerRepo)

	player := &model.Player{Base: model.Base{ID: 1}}
	mockPlayerRepo.On("Get", ctx, uint64(1)).Return(player, nil)
	mockRepo.On("GetAssignmentHistory", ctx, uint64(1)).Return([]model.PlayerCompanyAssignment{}, assert.AnError)

	result, err := svc.GetAssignmentHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEndCurrentAssignment_GetError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	mockRepo.On("GetCurrentAssignment", ctx, uint64(1)).Return((*model.PlayerCompanyAssignment)(nil), assert.AnError)

	err := svc.EndCurrentAssignment(ctx, 1, time.Now())

	assert.Error(t, err)
}

func TestEndCurrentAssignment_EndError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	assignment := &model.PlayerCompanyAssignment{PlayerID: 1, IsCurrent: true}
	endDate := time.Now()

	mockRepo.On("GetCurrentAssignment", ctx, uint64(1)).Return(assignment, nil)
	mockRepo.On("EndCurrentAssignment", ctx, uint64(1), endDate).Return(assert.AnError)

	err := svc.EndCurrentAssignment(ctx, 1, endDate)

	assert.Error(t, err)
}

func TestUpdateCompany_HistoryCreateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Name: "Old"}
	newName := "New"
	req := &model.UpdateSettlementCompanyRequest{Name: &newName}

	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.SettlementCompany")).Return(nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(assert.AnError)

	result, err := svc.UpdateCompany(ctx, 1, req, 1)

	// History error should not affect main flow
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestToggleCompanyStatus_HistoryCreateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSettlementCompanyRepository)
	svc := NewSettlementCompanyService(mockRepo, nil)

	company := &model.SettlementCompany{Base: model.Base{ID: 1}, Status: model.CompanyStatusActive}
	mockRepo.On("Get", ctx, uint64(1)).Return(company, nil)
	mockRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.SettlementCompanyHistory")).Return(assert.AnError)
	mockRepo.On("ToggleStatus", ctx, uint64(1), model.CompanyStatusInactive).Return(nil)

	err := svc.ToggleCompanyStatus(ctx, 1, false, 1)

	// History error should not affect main flow
	assert.NoError(t, err)
}

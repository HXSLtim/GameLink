package routingrule

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/collectionentity"
	"gamelink/internal/repository/routingrule"
)

// MockRoutingRuleRepository is a mock implementation of RoutingRuleRepository
type MockRoutingRuleRepository struct {
	mock.Mock
}

func (m *MockRoutingRuleRepository) Create(ctx context.Context, rule *model.RoutingRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) Get(ctx context.Context, id uint64) (*model.RoutingRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoutingRule), args.Error(1)
}

func (m *MockRoutingRuleRepository) Update(ctx context.Context, rule *model.RoutingRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) List(ctx context.Context, opts routingrule.ListOptions) ([]model.RoutingRule, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RoutingRule), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoutingRuleRepository) ToggleStatus(ctx context.Context, id uint64, status model.RuleStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) ListActiveByPriority(ctx context.Context) ([]model.RoutingRule, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.RoutingRule), args.Error(1)
}

func (m *MockRoutingRuleRepository) CreateHistory(ctx context.Context, history *model.RoutingRuleHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) GetHistory(ctx context.Context, ruleID uint64) ([]model.RoutingRuleHistory, error) {
	args := m.Called(ctx, ruleID)
	return args.Get(0).([]model.RoutingRuleHistory), args.Error(1)
}

func (m *MockRoutingRuleRepository) CreateRoutingLog(ctx context.Context, log *model.RoutingLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockRoutingRuleRepository) GetRoutingLogByPayment(ctx context.Context, paymentID uint64) (*model.RoutingLog, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoutingLog), args.Error(1)
}

func (m *MockRoutingRuleRepository) ListRoutingLogs(ctx context.Context, opts routingrule.RoutingLogListOptions) ([]model.RoutingLog, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RoutingLog), args.Get(1).(int64), args.Error(2)
}

// MockCollectionEntityRepository is a mock implementation of CollectionEntityRepository
type MockCollectionEntityRepository struct {
	mock.Mock
}

func (m *MockCollectionEntityRepository) Create(ctx context.Context, entity *model.CollectionEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) Get(ctx context.Context, id uint64) (*model.CollectionEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CollectionEntity), args.Error(1)
}

func (m *MockCollectionEntityRepository) GetByCreditCode(ctx context.Context, creditCode string) (*model.CollectionEntity, error) {
	args := m.Called(ctx, creditCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CollectionEntity), args.Error(1)
}

func (m *MockCollectionEntityRepository) GetDefault(ctx context.Context) (*model.CollectionEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CollectionEntity), args.Error(1)
}

func (m *MockCollectionEntityRepository) Update(ctx context.Context, entity *model.CollectionEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) List(ctx context.Context, opts collectionentity.ListOptions) ([]model.CollectionEntity, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CollectionEntity), args.Get(1).(int64), args.Error(2)
}

func (m *MockCollectionEntityRepository) ListActive(ctx context.Context) ([]model.CollectionEntity, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.CollectionEntity), args.Error(1)
}

func (m *MockCollectionEntityRepository) ToggleStatus(ctx context.Context, id uint64, status model.EntityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) SetDefault(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) CreateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) GetChannel(ctx context.Context, id uint64) (*model.PaymentChannelConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaymentChannelConfig), args.Error(1)
}

func (m *MockCollectionEntityRepository) GetChannelByEntityAndMethod(ctx context.Context, entityID uint64, method model.PaymentMethod) (*model.PaymentChannelConfig, error) {
	args := m.Called(ctx, entityID, method)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaymentChannelConfig), args.Error(1)
}

func (m *MockCollectionEntityRepository) UpdateChannel(ctx context.Context, channel *model.PaymentChannelConfig) error {
	args := m.Called(ctx, channel)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) DeleteChannel(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) ListChannelsByEntity(ctx context.Context, entityID uint64) ([]model.PaymentChannelConfig, error) {
	args := m.Called(ctx, entityID)
	return args.Get(0).([]model.PaymentChannelConfig), args.Error(1)
}

func (m *MockCollectionEntityRepository) CreateHistory(ctx context.Context, history *model.CollectionEntityHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) GetHistory(ctx context.Context, entityID uint64) ([]model.CollectionEntityHistory, error) {
	args := m.Called(ctx, entityID)
	return args.Get(0).([]model.CollectionEntityHistory), args.Error(1)
}

func (m *MockCollectionEntityRepository) UpdateCollectionStats(ctx context.Context, entityID uint64, amountCents int64) error {
	args := m.Called(ctx, entityID, amountCents)
	return args.Error(0)
}

func (m *MockCollectionEntityRepository) GetCollectionStats(ctx context.Context, entityID uint64) (int64, int64, error) {
	args := m.Called(ctx, entityID)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

// Helper function to create JSON raw message
func mustJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

// Test cases

func TestCreateRule_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Name:   "Test Entity",
		Status: model.EntityStatusActive,
	}

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
		Description:    "Test description",
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockRuleRepo.On("Create", ctx, mock.AnythingOfType("*model.RoutingRule")).Return(nil)
	mockRuleRepo.On("Get", ctx, mock.AnythingOfType("uint64")).Return(&model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Test Rule",
		Priority:       1,
		TargetEntityID: 1,
	}, nil)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockEntityRepo.AssertExpectations(t)
	mockRuleRepo.AssertExpectations(t)
}

func TestCreateRule_EntityNotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 999,
	}

	mockEntityRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestCreateRule_EntityInactive(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Name:   "Test Entity",
		Status: model.EntityStatusInactive,
	}

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "inactive")
}

func TestCreateRule_InvalidConditions(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}

	// Empty conditions
	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     []model.RoutingCondition{},
		TargetEntityID: 1,
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

func TestGetRule_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{
		Base:     model.Base{ID: 1},
		Name:     "Test Rule",
		Priority: 1,
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)

	result, err := svc.GetRule(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Test Rule", result.Name)
}

func TestGetRule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetRule(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteRule_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("Delete", ctx, uint64(1)).Return(nil)

	err := svc.DeleteRule(ctx, 1)

	assert.NoError(t, err)
}

func TestDeleteRule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.DeleteRule(ctx, 999)

	assert.Error(t, err)
}

func TestListRules_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Name: "Rule 1"},
		{Base: model.Base{ID: 2}, Name: "Rule 2"},
	}

	mockRuleRepo.On("List", ctx, mock.AnythingOfType("routingrule.ListOptions")).Return(rules, int64(2), nil)

	req := &model.ListRoutingRulesRequest{Page: 1, PageSize: 10}
	result, err := svc.ListRules(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.Total)
	assert.Len(t, result.Rules, 2)
}

func TestToggleRuleStatus_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Status: model.RuleStatusActive}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.RoutingRuleHistory")).Return(nil)
	mockRuleRepo.On("ToggleStatus", ctx, uint64(1), model.RuleStatusInactive).Return(nil)

	err := svc.ToggleRuleStatus(ctx, 1, false, 1)

	assert.NoError(t, err)
}

func TestGetRuleHistory_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}}
	histories := []model.RoutingRuleHistory{
		{FieldName: "name", OldValue: "old", NewValue: "new"},
	}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("GetHistory", ctx, uint64(1)).Return(histories, nil)

	result, err := svc.GetRuleHistory(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSetDefaultEntity_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockEntityRepo.On("SetDefault", ctx, uint64(1)).Return(nil)

	err := svc.SetDefaultEntity(ctx, 1, 1)

	assert.NoError(t, err)
}

func TestSetDefaultEntity_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.SetDefaultEntity(ctx, 999, 1)

	assert.Error(t, err)
}

func TestSetDefaultEntity_Inactive(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusInactive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	err := svc.SetDefaultEntity(ctx, 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "inactive")
}

func TestGetDefaultEntity_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}
	mockEntityRepo.On("GetDefault", ctx).Return(entity, nil)

	result, err := svc.GetDefaultEntity(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "Default", result.Name)
}

func TestGetDefaultEntity_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("GetDefault", ctx).Return(nil, repository.ErrNotFound)

	result, err := svc.GetDefaultEntity(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListActiveRulesByPriority_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Priority: 1},
		{Base: model.Base{ID: 2}, Priority: 2},
	}
	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)

	result, err := svc.ListActiveRulesByPriority(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCreateRoutingLog_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	log := &model.RoutingLog{PaymentID: 1, OrderID: 1}
	mockRuleRepo.On("CreateRoutingLog", ctx, log).Return(nil)

	err := svc.CreateRoutingLog(ctx, log)

	assert.NoError(t, err)
}

func TestGetRoutingLogByPayment_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	log := &model.RoutingLog{PaymentID: 1, OrderID: 1}
	mockRuleRepo.On("GetRoutingLogByPayment", ctx, uint64(1)).Return(log, nil)

	result, err := svc.GetRoutingLogByPayment(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), result.PaymentID)
}

func TestGetRoutingLogByPayment_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("GetRoutingLogByPayment", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetRoutingLogByPayment(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListRoutingLogs_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	logs := []model.RoutingLog{{PaymentID: 1}, {PaymentID: 2}}
	mockRuleRepo.On("ListRoutingLogs", ctx, mock.AnythingOfType("routingrule.RoutingLogListOptions")).Return(logs, int64(2), nil)

	opts := routingrule.RoutingLogListOptions{Page: 1, PageSize: 10}
	result, total, err := svc.ListRoutingLogs(ctx, opts)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestReorderPriorities_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule1 := &model.RoutingRule{Base: model.Base{ID: 1}, Priority: 2}
	rule2 := &model.RoutingRule{Base: model.Base{ID: 2}, Priority: 1}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule1, nil)
	mockRuleRepo.On("Get", ctx, uint64(2)).Return(rule2, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.RoutingRuleHistory")).Return(nil)
	mockRuleRepo.On("Update", ctx, mock.AnythingOfType("*model.RoutingRule")).Return(nil)

	err := svc.ReorderPriorities(ctx, []uint64{1, 2}, 1)

	assert.NoError(t, err)
}

func TestValidateConditions_InvalidField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	conditions := []model.RoutingCondition{
		{Field: "invalid_field", Operator: model.ConditionOperatorEquals, Value: mustJSON("test")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

func TestValidateConditions_InvalidOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: "invalid_op", Value: mustJSON("test")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

func TestValidateConditions_EmptyValue(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: nil},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid routing condition")
}

func TestMatchCollectionEntity_NoRulesUsesDefault(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:      model.Base{ID: 1},
		Name:      "Default Entity",
		IsDefault: true,
		PaymentChannels: []model.PaymentChannelConfig{
			{MerchantNo: "MERCHANT_001"},
		},
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
	assert.Equal(t, uint64(1), result.CollectionEntityID)
}

func TestUpdateRule_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Old Name",
		Priority:       1,
		TargetEntityID: 1,
		Description:    "Old description",
		Conditions:     mustJSON([]model.RoutingCondition{{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")}}),
	}

	newName := "New Name"
	newPriority := 2
	newDesc := "New description"
	newEntityID := uint64(2)
	newConditions := []model.RoutingCondition{
		{Field: model.ConditionFieldServiceType, Operator: model.ConditionOperatorEquals, Value: mustJSON("escort")},
	}

	req := &model.UpdateRoutingRuleRequest{
		Name:           &newName,
		Priority:       &newPriority,
		Description:    &newDesc,
		TargetEntityID: &newEntityID,
		Conditions:     &newConditions,
	}

	newEntity := &model.CollectionEntity{Base: model.Base{ID: 2}, Status: model.EntityStatusActive}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil).Once()
	mockEntityRepo.On("Get", ctx, uint64(2)).Return(newEntity, nil)
	mockRuleRepo.On("Update", ctx, mock.AnythingOfType("*model.RoutingRule")).Return(nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.AnythingOfType("*model.RoutingRuleHistory")).Return(nil)
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(&model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           newName,
		Priority:       newPriority,
		TargetEntityID: newEntityID,
	}, nil).Once()

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
}

func TestUpdateRule_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	newName := "New Name"
	req := &model.UpdateRoutingRuleRequest{Name: &newName}

	result, err := svc.UpdateRule(ctx, 999, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_TargetEntityNotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{Base: model.Base{ID: 1}, TargetEntityID: 1}
	newEntityID := uint64(999)
	req := &model.UpdateRoutingRuleRequest{TargetEntityID: &newEntityID}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockEntityRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_TargetEntityInactive(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{Base: model.Base{ID: 1}, TargetEntityID: 1}
	newEntityID := uint64(2)
	req := &model.UpdateRoutingRuleRequest{TargetEntityID: &newEntityID}

	inactiveEntity := &model.CollectionEntity{Base: model.Base{ID: 2}, Status: model.EntityStatusInactive}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockEntityRepo.On("Get", ctx, uint64(2)).Return(inactiveEntity, nil)

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "inactive")
}

func TestMatchCollectionEntity_MatchesRule(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{
			Base:           model.Base{ID: 1},
			Name:           "LOL Rule",
			Priority:       1,
			Conditions:     conditionsJSON,
			TargetEntityID: 10,
			Status:         model.RuleStatusActive,
		},
	}

	targetEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 10},
		Name:   "LOL Entity",
		Status: model.EntityStatusActive,
		PaymentChannels: []model.PaymentChannelConfig{
			{MerchantNo: "LOL_MERCHANT"},
		},
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
	assert.Equal(t, uint64(10), result.CollectionEntityID)
	assert.Equal(t, "LOL Entity", result.EntityName)
	assert.Equal(t, "LOL_MERCHANT", result.MerchantNo)
	assert.NotNil(t, result.MatchedRuleID)
	assert.Equal(t, uint64(1), *result.MatchedRuleID)
}

func TestMatchCollectionEntity_InOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorIn, Value: mustJSON([]string{"LOL", "DOTA2", "CSGO"})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "DOTA2"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_NotEqualsOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorNotEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "DOTA2"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_NotInOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorNotIn, Value: mustJSON([]string{"LOL", "DOTA2"})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "CSGO"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_GreaterThanOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(10000))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 15000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_LessThanOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorLessThan, Value: mustJSON(int64(10000))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 5000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_BetweenOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorBetween, Value: mustJSON([]int64{5000, 15000})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 10000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_RegionField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldRegion, Operator: model.ConditionOperatorEquals, Value: mustJSON("CN")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{Region: "CN"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_ServiceTypeField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldServiceType, Operator: model.ConditionOperatorEquals, Value: mustJSON("escort")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{ServiceType: "escort"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_MultipleConditions(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(10000))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL", AmountCents: 15000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_NoDefaultEntity(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(nil, repository.ErrNotFound)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMatchCollectionEntity_SkipsInactiveEntity(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	inactiveEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusInactive}
	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(inactiveEntity, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestToggleRuleStatus_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	err := svc.ToggleRuleStatus(ctx, 999, true, 1)

	assert.Error(t, err)
}

func TestGetRuleHistory_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

	result, err := svc.GetRuleHistory(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTestRouting(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.TestRouting(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_IntInOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorIn, Value: mustJSON([]int64{5000, 10000, 15000})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 10000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestMatchCollectionEntity_IntEqualsOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorEquals, Value: mustJSON(int64(10000))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 10000}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

// ============================================================================
// Additional Tests for Coverage Improvement
// ============================================================================

func TestCreateRule_EntityGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return((*model.CollectionEntity)(nil), assert.AnError)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetRule_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	result, err := svc.GetRule(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDeleteRule_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	err := svc.DeleteRule(ctx, 1)

	assert.Error(t, err)
}

func TestDeleteRule_DeleteError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("Delete", ctx, uint64(1)).Return(assert.AnError)

	err := svc.DeleteRule(ctx, 1)

	assert.Error(t, err)
}

func TestListRules_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("List", ctx, mock.Anything).Return([]model.RoutingRule{}, int64(0), assert.AnError)

	req := &model.ListRoutingRulesRequest{Page: 1, PageSize: 10}
	result, err := svc.ListRules(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestToggleRuleStatus_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	err := svc.ToggleRuleStatus(ctx, 1, true, 1)

	assert.Error(t, err)
}

func TestToggleRuleStatus_ToggleError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Status: model.RuleStatusActive}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("ToggleStatus", ctx, uint64(1), model.RuleStatusInactive).Return(assert.AnError)

	err := svc.ToggleRuleStatus(ctx, 1, false, 1)

	assert.Error(t, err)
}

func TestGetRuleHistory_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	result, err := svc.GetRuleHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetRuleHistory_HistoryError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("GetHistory", ctx, uint64(1)).Return([]model.RoutingRuleHistory{}, assert.AnError)

	result, err := svc.GetRuleHistory(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSetDefaultEntity_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("Get", ctx, uint64(1)).Return((*model.CollectionEntity)(nil), assert.AnError)

	err := svc.SetDefaultEntity(ctx, 1, 1)

	assert.Error(t, err)
}

func TestSetDefaultEntity_SetDefaultError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockEntityRepo.On("SetDefault", ctx, uint64(1)).Return(assert.AnError)

	err := svc.SetDefaultEntity(ctx, 1, 1)

	assert.Error(t, err)
}

func TestGetDefaultEntity_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), assert.AnError)

	result, err := svc.GetDefaultEntity(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListActiveRulesByPriority_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, assert.AnError)

	result, err := svc.ListActiveRulesByPriority(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateRoutingLog_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	log := &model.RoutingLog{PaymentID: 1}
	mockRuleRepo.On("CreateRoutingLog", ctx, log).Return(assert.AnError)

	err := svc.CreateRoutingLog(ctx, log)

	assert.Error(t, err)
}

func TestGetRoutingLogByPayment_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("GetRoutingLogByPayment", ctx, uint64(1)).Return((*model.RoutingLog)(nil), assert.AnError)

	result, err := svc.GetRoutingLogByPayment(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListRoutingLogs_Error(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListRoutingLogs", ctx, mock.Anything).Return([]model.RoutingLog{}, int64(0), assert.AnError)

	opts := routingrule.RoutingLogListOptions{Page: 1, PageSize: 10}
	result, total, err := svc.ListRoutingLogs(ctx, opts)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
}

func TestReorderPriorities_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	err := svc.ReorderPriorities(ctx, []uint64{1}, 1)

	assert.Error(t, err)
}

func TestReorderPriorities_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Priority: 5}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(assert.AnError)

	err := svc.ReorderPriorities(ctx, []uint64{1}, 1)

	assert.Error(t, err)
}

func TestReorderPriorities_SkipsNotFound(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(999)).Return((*model.RoutingRule)(nil), repository.ErrNotFound)

	err := svc.ReorderPriorities(ctx, []uint64{999}, 1)

	assert.NoError(t, err)
}

func TestUpdateRule_GetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("Get", ctx, uint64(1)).Return((*model.RoutingRule)(nil), assert.AnError)

	newName := "New Name"
	req := &model.UpdateRoutingRuleRequest{Name: &newName}

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_EntityGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{Base: model.Base{ID: 1}, TargetEntityID: 1}
	newEntityID := uint64(2)
	req := &model.UpdateRoutingRuleRequest{TargetEntityID: &newEntityID}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)
	mockEntityRepo.On("Get", ctx, uint64(2)).Return((*model.CollectionEntity)(nil), assert.AnError)

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_InvalidConditions(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{Base: model.Base{ID: 1}, TargetEntityID: 1}
	emptyConditions := []model.RoutingCondition{}
	req := &model.UpdateRoutingRuleRequest{Conditions: &emptyConditions}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil)

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{Base: model.Base{ID: 1}, Name: "Old", TargetEntityID: 1}
	newName := "New"
	req := &model.UpdateRoutingRuleRequest{Name: &newName}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil).Once()
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(assert.AnError)

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMatchCollectionEntity_ListActiveError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, assert.AnError)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMatchCollectionEntity_GetDefaultError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), assert.AnError)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestListRules_WithStatus(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{{Base: model.Base{ID: 1}}}
	mockRuleRepo.On("List", ctx, mock.Anything).Return(rules, int64(1), nil)

	req := &model.ListRoutingRulesRequest{Page: 1, PageSize: 10, Status: "active"}
	result, err := svc.ListRules(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}

func TestToggleRuleStatus_NoStatusChange(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Status: model.RuleStatusActive}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("ToggleStatus", ctx, uint64(1), model.RuleStatusActive).Return(nil)

	err := svc.ToggleRuleStatus(ctx, 1, true, 1)

	assert.NoError(t, err)
}

func TestMatchCollectionEntity_EntityGetError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return((*model.CollectionEntity)(nil), assert.AnError)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_ConditionNotMatched(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "DOTA2"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_InvalidConditionsJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: []byte("invalid json"), TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCondition_UnknownField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: "unknown_field", Operator: model.ConditionOperatorEquals, Value: mustJSON("test")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCondition_UnknownOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: "unknown_op", Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchEquals_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to string
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: json.RawMessage(`{"not":"a string"}`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchIn_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to []string
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorIn, Value: json.RawMessage(`{"not":"an array"}`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchGreaterThan_NonNumericField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchLessThan_NonNumericField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorLessThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchBetween_NonNumericField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorBetween, Value: mustJSON([]int64{100, 200})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchBetween_InvalidValueCount(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorBetween, Value: mustJSON([]int64{100})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchGreaterThan_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: json.RawMessage(`"not a number"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchLessThan_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorLessThan, Value: json.RawMessage(`"not a number"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchBetween_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to []int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorBetween, Value: json.RawMessage(`"not an array"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchIn_IntNotInList(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorIn, Value: mustJSON([]int64{100, 200, 300})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchIn_StringNotInList(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorIn, Value: mustJSON([]string{"LOL", "DOTA2"})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "CSGO"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchIn_IntInvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to []int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorIn, Value: json.RawMessage(`{"not":"an array"}`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchEquals_IntInvalidJSON(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Use truly invalid JSON that cannot be unmarshaled to int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorEquals, Value: json.RawMessage(`"not a number"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 150}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestReorderPriorities_NoPriorityChange(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Priority: 1}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)

	err := svc.ReorderPriorities(ctx, []uint64{1}, 1)

	assert.NoError(t, err)
}

// ============================================================================
// Additional Tests for 90%+ Coverage - Part 2
// ============================================================================

func TestCreateRule_CreateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockRuleRepo.On("Create", ctx, mock.Anything).Return(assert.AnError)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdateRule_HistoryCreateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Old Name",
		Priority:       1,
		TargetEntityID: 1,
		Description:    "Old description",
	}

	newName := "New Name"
	req := &model.UpdateRoutingRuleRequest{Name: &newName}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil).Once()
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(assert.AnError) // History error should not affect main flow
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(&model.RoutingRule{
		Base: model.Base{ID: 1},
		Name: newName,
	}, nil).Once()

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestToggleRuleStatus_HistoryCreateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Status: model.RuleStatusActive}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(assert.AnError) // History error should not affect main flow
	mockRuleRepo.On("ToggleStatus", ctx, uint64(1), model.RuleStatusInactive).Return(nil)

	err := svc.ToggleRuleStatus(ctx, 1, false, 1)

	assert.NoError(t, err)
}

func TestReorderPriorities_HistoryCreateError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rule := &model.RoutingRule{Base: model.Base{ID: 1}, Priority: 5}
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(rule, nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(assert.AnError) // History error should not affect main flow
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(nil)

	err := svc.ReorderPriorities(ctx, []uint64{1}, 1)

	assert.NoError(t, err)
}

func TestMatchEquals_UnsupportedType(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Test with a field that returns an unsupported type (this is hard to trigger directly)
	// We'll test the matchEquals function indirectly through matchCondition
	// by using a condition that would result in an unsupported type comparison

	// Since all fields map to string or int64, we need to test the error path differently
	// Let's test with invalid JSON that causes unmarshal error for int64
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorEquals, Value: json.RawMessage(`{"invalid":"json"}`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 100}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchIn_UnsupportedType(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Test matchIn with invalid JSON for int64 array
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorIn, Value: json.RawMessage(`"not an array"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 100}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestValidateConditions_AllValidFields(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}

	// Test all valid field types
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
		{Field: model.ConditionFieldServiceType, Operator: model.ConditionOperatorEquals, Value: mustJSON("escort")},
		{Field: model.ConditionFieldRegion, Operator: model.ConditionOperatorEquals, Value: mustJSON("CN")},
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(100))},
	}

	req := &model.CreateRoutingRuleRequest{
		Name:           "Test Rule",
		Priority:       1,
		Conditions:     conditions,
		TargetEntityID: 1,
	}

	mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
	mockRuleRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Get", ctx, mock.Anything).Return(&model.RoutingRule{Base: model.Base{ID: 1}}, nil)

	result, err := svc.CreateRule(ctx, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	_ = svc // Use svc to avoid unused variable error
}

func TestValidateConditions_AllValidOperators(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)
	_ = svc // suppress unused variable warning - subtests create their own svc

	entity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}

	// Test all valid operators
	testCases := []struct {
		name     string
		operator model.ConditionOperator
		value    json.RawMessage
	}{
		{"equals", model.ConditionOperatorEquals, mustJSON("LOL")},
		{"notEquals", model.ConditionOperatorNotEquals, mustJSON("LOL")},
		{"in", model.ConditionOperatorIn, mustJSON([]string{"LOL", "DOTA2"})},
		{"notIn", model.ConditionOperatorNotIn, mustJSON([]string{"LOL", "DOTA2"})},
		{"greaterThan", model.ConditionOperatorGreaterThan, mustJSON(int64(100))},
		{"lessThan", model.ConditionOperatorLessThan, mustJSON(int64(100))},
		{"between", model.ConditionOperatorBetween, mustJSON([]int64{100, 200})},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRuleRepo := new(MockRoutingRuleRepository)
			mockEntityRepo := new(MockCollectionEntityRepository)
			svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

			field := model.ConditionFieldGameType
			if tc.operator == model.ConditionOperatorGreaterThan || tc.operator == model.ConditionOperatorLessThan || tc.operator == model.ConditionOperatorBetween {
				field = model.ConditionFieldOrderAmount
			}

			conditions := []model.RoutingCondition{
				{Field: field, Operator: tc.operator, Value: tc.value},
			}

			req := &model.CreateRoutingRuleRequest{
				Name:           "Test Rule",
				Priority:       1,
				Conditions:     conditions,
				TargetEntityID: 1,
			}

			mockEntityRepo.On("Get", ctx, uint64(1)).Return(entity, nil)
			mockRuleRepo.On("Create", ctx, mock.Anything).Return(nil)
			mockRuleRepo.On("Get", ctx, mock.Anything).Return(&model.RoutingRule{Base: model.Base{ID: 1}}, nil)

			result, err := svc.CreateRule(ctx, req, 1)

			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestDetectChanges_AllFields(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	existingRule := &model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           "Old Name",
		Priority:       1,
		TargetEntityID: 1,
		Description:    "Old description",
	}

	newName := "New Name"
	newPriority := 2
	newEntityID := uint64(2)
	newDesc := "New description"

	req := &model.UpdateRoutingRuleRequest{
		Name:           &newName,
		Priority:       &newPriority,
		TargetEntityID: &newEntityID,
		Description:    &newDesc,
	}

	newEntity := &model.CollectionEntity{Base: model.Base{ID: 2}, Status: model.EntityStatusActive}

	mockRuleRepo.On("Get", ctx, uint64(1)).Return(existingRule, nil).Once()
	mockEntityRepo.On("Get", ctx, uint64(2)).Return(newEntity, nil)
	mockRuleRepo.On("Update", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("CreateHistory", ctx, mock.Anything).Return(nil)
	mockRuleRepo.On("Get", ctx, uint64(1)).Return(&model.RoutingRule{
		Base:           model.Base{ID: 1},
		Name:           newName,
		Priority:       newPriority,
		TargetEntityID: newEntityID,
		Description:    newDesc,
	}, nil).Once()

	result, err := svc.UpdateRule(ctx, 1, req, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Verify CreateHistory was called 4 times (once for each changed field)
	mockRuleRepo.AssertNumberOfCalls(t, "CreateHistory", 4)
}

func TestMatchCollectionEntity_NoPaymentChannels(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	// Entity without payment channels
	targetEntity := &model.CollectionEntity{
		Base:            model.Base{ID: 10},
		Name:            "LOL Entity",
		Status:          model.EntityStatusActive,
		PaymentChannels: []model.PaymentChannelConfig{}, // Empty
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
	assert.Equal(t, "", result.MerchantNo) // No merchant number
}

func TestMatchCollectionEntity_DefaultEntityNoPaymentChannels(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	// Default entity without payment channels
	defaultEntity := &model.CollectionEntity{
		Base:            model.Base{ID: 1},
		Name:            "Default",
		IsDefault:       true,
		PaymentChannels: []model.PaymentChannelConfig{}, // Empty
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"}
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
	assert.Equal(t, "", result.MerchantNo) // No merchant number
}

func TestListRules_DefaultPagination(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{{Base: model.Base{ID: 1}}}
	mockRuleRepo.On("List", ctx, mock.MatchedBy(func(opts routingrule.ListOptions) bool {
		return opts.Page == 1 && opts.PageSize == 20
	})).Return(rules, int64(1), nil)

	// Test with invalid page and pageSize
	req := &model.ListRoutingRulesRequest{Page: 0, PageSize: 0}
	result, err := svc.ListRules(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 20, result.PageSize)
}

func TestListRules_MaxPageSize(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	rules := []model.RoutingRule{{Base: model.Base{ID: 1}}}
	mockRuleRepo.On("List", ctx, mock.MatchedBy(func(opts routingrule.ListOptions) bool {
		return opts.PageSize == 20 // Should be capped at 20 (or 100?)
	})).Return(rules, int64(1), nil)

	// Test with page size over limit
	req := &model.ListRoutingRulesRequest{Page: 1, PageSize: 200}
	result, err := svc.ListRules(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestMatchCollectionEntity_CaseInsensitiveStringMatch(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("lol")}, // lowercase
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"} // uppercase
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault) // Should match due to case-insensitive comparison
}

func TestMatchCollectionEntity_CaseInsensitiveInMatch(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorIn, Value: mustJSON([]string{"lol", "dota2"})}, // lowercase
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"} // uppercase
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault) // Should match due to case-insensitive comparison
}

func TestMatchCollectionEntity_BetweenBoundaryValues(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)
	_ = svc // suppress unused variable warning - subtests create their own svc

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorBetween, Value: mustJSON([]int64{100, 200})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}

	// Test lower boundary
	t.Run("lower boundary", func(t *testing.T) {
		mockRuleRepo := new(MockRoutingRuleRepository)
		mockEntityRepo := new(MockCollectionEntityRepository)
		svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

		mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
		mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

		req := &model.RoutingTestRequest{AmountCents: 100} // Exactly at lower boundary
		result, err := svc.MatchCollectionEntity(ctx, req)

		assert.NoError(t, err)
		assert.False(t, result.IsDefault)
	})

	// Test upper boundary
	t.Run("upper boundary", func(t *testing.T) {
		mockRuleRepo := new(MockRoutingRuleRepository)
		mockEntityRepo := new(MockCollectionEntityRepository)
		svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

		mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
		mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)

		req := &model.RoutingTestRequest{AmountCents: 200} // Exactly at upper boundary
		result, err := svc.MatchCollectionEntity(ctx, req)

		assert.NoError(t, err)
		assert.False(t, result.IsDefault)
	})

	// Test outside boundaries
	t.Run("below lower boundary", func(t *testing.T) {
		mockRuleRepo := new(MockRoutingRuleRepository)
		mockEntityRepo := new(MockCollectionEntityRepository)
		svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

		defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

		mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
		mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

		req := &model.RoutingTestRequest{AmountCents: 99} // Below lower boundary
		result, err := svc.MatchCollectionEntity(ctx, req)

		assert.NoError(t, err)
		assert.True(t, result.IsDefault)
	})
}

func TestMatchCollectionEntity_GreaterThanNotMatched(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 100} // Equal, not greater
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_LessThanNotMatched(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorLessThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{AmountCents: 100} // Equal, not less
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_NotEqualsNotMatched(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorNotEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"} // Equal, so notEquals should not match
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestMatchCollectionEntity_NotInNotMatched(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	svc := NewRoutingRuleService(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorNotIn, Value: mustJSON([]string{"LOL", "DOTA2"})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", IsDefault: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)

	req := &model.RoutingTestRequest{GameType: "LOL"} // In list, so notIn should not match
	result, err := svc.MatchCollectionEntity(ctx, req)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

// ============================================================================
// RoutingEngine Tests
// ============================================================================

func TestNewRoutingEngine(t *testing.T) {
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)

	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	assert.NotNil(t, engine)
}

func TestRoutingEngine_RoutePayment_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 10},
		Name:   "LOL Entity",
		Status: model.EntityStatusActive,
	}

	channel := &model.PaymentChannelConfig{
		MerchantNo: "MERCHANT001",
		Enabled:    true,
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{
		GameType: "LOL",
		Method:   model.PaymentMethodWeChat,
	}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint64(10), result.CollectionEntityID)
	assert.Equal(t, "MERCHANT001", result.MerchantNo)
	assert.False(t, result.IsDefault)
}

func TestRoutingEngine_RoutePayment_FallbackToDefault(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:      model.Base{ID: 1},
		Name:      "Default Entity",
		Status:    model.EntityStatusActive,
		IsDefault: true,
	}

	channel := &model.PaymentChannelConfig{
		MerchantNo: "DEFAULT_MERCHANT",
		Enabled:    true,
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodAlipay).Return(channel, nil)

	routingCtx := &RoutingContext{
		GameType: "UNKNOWN",
		Method:   model.PaymentMethodAlipay,
	}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.IsDefault)
	assert.True(t, result.IsFallback)
	assert.Equal(t, "DEFAULT_MERCHANT", result.MerchantNo)
}

func TestRoutingEngine_RoutePayment_NoDefaultEntity(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), repository.ErrNotFound)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_RoutePayment_ListRulesError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, assert.AnError)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_RoutePayment_EntityInactive(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	inactiveEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 10},
		Status: model.EntityStatusInactive,
	}

	defaultEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Name:   "Default",
		Status: model.EntityStatusActive,
	}

	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(inactiveEntity, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_RoutePayment_NoPaymentChannel(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}
	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Name: "Default", Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return((*model.PaymentChannelConfig)(nil), repository.ErrNotFound)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_FallbackToDefault_InactiveDefault(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	inactiveDefault := &model.CollectionEntity{
		Base:      model.Base{ID: 1},
		Status:    model.EntityStatusInactive,
		IsDefault: true,
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(inactiveDefault, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_FallbackToDefault_NoChannelUseAny(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:      model.Base{ID: 1},
		Status:    model.EntityStatusActive,
		IsDefault: true,
		PaymentChannels: []model.PaymentChannelConfig{
			{MerchantNo: "ANY_MERCHANT", Enabled: true},
		},
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return((*model.PaymentChannelConfig)(nil), repository.ErrNotFound)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.Equal(t, "ANY_MERCHANT", result.MerchantNo)
}

func TestRoutingEngine_FallbackToDefault_NoChannelAtAll(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:            model.Base{ID: 1},
		Status:          model.EntityStatusActive,
		IsDefault:       true,
		PaymentChannels: []model.PaymentChannelConfig{},
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return((*model.PaymentChannelConfig)(nil), repository.ErrNotFound)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_MatchRule_EmptyConditions(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Rule with empty conditions
	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: json.RawMessage(`[]`), TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_MatchCondition_AllFields(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
		{Field: model.ConditionFieldServiceType, Operator: model.ConditionOperatorEquals, Value: mustJSON("escort")},
		{Field: model.ConditionFieldRegion, Operator: model.ConditionOperatorEquals, Value: mustJSON("CN")},
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "MERCHANT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{
		GameType:    "LOL",
		ServiceType: "escort",
		Region:      "CN",
		AmountCents: 200,
		Method:      model.PaymentMethodWeChat,
	}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestRoutingEngine_MatchCondition_UnknownField(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: "unknown_field", Operator: model.ConditionOperatorEquals, Value: mustJSON("value")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_EvaluateCondition_AllOperators(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		operator    model.ConditionOperator
		field       model.ConditionField
		value       json.RawMessage
		routingCtx  *RoutingContext
		shouldMatch bool
	}{
		{"equals_match", model.ConditionOperatorEquals, model.ConditionFieldGameType, mustJSON("LOL"), &RoutingContext{GameType: "LOL"}, true},
		{"equals_no_match", model.ConditionOperatorEquals, model.ConditionFieldGameType, mustJSON("DOTA2"), &RoutingContext{GameType: "LOL"}, false},
		{"notEquals_match", model.ConditionOperatorNotEquals, model.ConditionFieldGameType, mustJSON("DOTA2"), &RoutingContext{GameType: "LOL"}, true},
		{"in_match", model.ConditionOperatorIn, model.ConditionFieldGameType, mustJSON([]string{"LOL", "DOTA2"}), &RoutingContext{GameType: "LOL"}, true},
		{"notIn_match", model.ConditionOperatorNotIn, model.ConditionFieldGameType, mustJSON([]string{"DOTA2", "CS"}), &RoutingContext{GameType: "LOL"}, true},
		{"greaterThan_match", model.ConditionOperatorGreaterThan, model.ConditionFieldOrderAmount, mustJSON(int64(100)), &RoutingContext{AmountCents: 200}, true},
		{"lessThan_match", model.ConditionOperatorLessThan, model.ConditionFieldOrderAmount, mustJSON(int64(100)), &RoutingContext{AmountCents: 50}, true},
		{"between_match", model.ConditionOperatorBetween, model.ConditionFieldOrderAmount, mustJSON([]int64{100, 200}), &RoutingContext{AmountCents: 150}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRuleRepo := new(MockRoutingRuleRepository)
			mockEntityRepo := new(MockCollectionEntityRepository)
			engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

			conditions := []model.RoutingCondition{
				{Field: tc.field, Operator: tc.operator, Value: tc.value},
			}
			conditionsJSON, _ := json.Marshal(conditions)

			rules := []model.RoutingRule{
				{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
			}

			targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}
			defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
			channel := &model.PaymentChannelConfig{MerchantNo: "MERCHANT", Enabled: true}

			mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
			if tc.shouldMatch {
				mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
				mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return(channel, nil)
			} else {
				mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
				mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)
			}

			tc.routingCtx.Method = model.PaymentMethodWeChat
			result, err := engine.RoutePayment(ctx, tc.routingCtx)

			assert.NoError(t, err)
			assert.Equal(t, !tc.shouldMatch, result.IsDefault)
		})
	}
}

func TestRoutingEngine_EvaluateCondition_UnknownOperator(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: "unknown_operator", Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_MatchEquals_UnsupportedType(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Invalid JSON for int64 field
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorEquals, Value: json.RawMessage(`"not a number"`)},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{AmountCents: 100, Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_MatchIn_IntArray(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorIn, Value: mustJSON([]int64{100, 200, 300})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "MERCHANT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{AmountCents: 200, Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.False(t, result.IsDefault)
}

func TestRoutingEngine_MatchBetween_InvalidValues(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// Between with only 1 value (invalid)
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldOrderAmount, Operator: model.ConditionOperatorBetween, Value: mustJSON([]int64{100})},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{AmountCents: 150, Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_MatchGreaterThan_NonNumeric(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// GreaterThan on string field (invalid)
	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorGreaterThan, Value: mustJSON(int64(100))},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_CreateRoutingLog(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	ruleID := uint64(1)
	result := &RoutingResult{
		CollectionEntityID: 10,
		EntityName:         "Test Entity",
		MerchantNo:         "MERCHANT001",
		MatchedRuleID:      &ruleID,
		MatchedRuleName:    "Test Rule",
		IsDefault:          false,
		MatchDetails: []model.RoutingCondition{
			{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
		},
	}

	mockRuleRepo.On("CreateRoutingLog", ctx, mock.Anything).Return(nil)

	err := engine.CreateRoutingLog(ctx, 100, 200, result)

	assert.NoError(t, err)
	mockRuleRepo.AssertCalled(t, "CreateRoutingLog", ctx, mock.Anything)
}

func TestRoutingEngine_GetRoutingLogByPayment(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	log := &model.RoutingLog{
		Base:               model.Base{ID: 1},
		PaymentID:          100,
		CollectionEntityID: 10,
		MerchantNo:         "MERCHANT001",
	}

	mockRuleRepo.On("GetRoutingLogByPayment", ctx, uint64(100)).Return(log, nil)

	result, err := engine.GetRoutingLogByPayment(ctx, 100)

	assert.NoError(t, err)
	assert.Equal(t, uint64(100), result.PaymentID)
}

func TestRoutingEngine_RoutePaymentWithFallback_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}
	channel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(channel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePaymentWithFallback(ctx, routingCtx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestRoutingEngine_RoutePaymentWithFallback_FindAnyEntity(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	// First call fails
	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil).Once()
	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), repository.ErrNotFound).Once()

	// Fallback to find any entity
	entities := []model.CollectionEntity{
		{
			Base:   model.Base{ID: 5},
			Name:   "Any Entity",
			Status: model.EntityStatusActive,
			PaymentChannels: []model.PaymentChannelConfig{
				{MerchantNo: "ANY_MERCHANT", Enabled: true},
			},
		},
	}
	mockEntityRepo.On("ListActive", ctx).Return(entities, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(5), model.PaymentMethodWeChat).Return((*model.PaymentChannelConfig)(nil), repository.ErrNotFound)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePaymentWithFallback(ctx, routingCtx)

	assert.NoError(t, err)
	assert.Equal(t, "ANY_MERCHANT", result.MerchantNo)
	assert.True(t, result.IsFallback)
}

func TestRoutingEngine_FindAnyAvailableEntity_NoEntities(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), repository.ErrNotFound)
	mockEntityRepo.On("ListActive", ctx).Return([]model.CollectionEntity{}, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePaymentWithFallback(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_FindAnyAvailableEntity_ListError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), repository.ErrNotFound)
	mockEntityRepo.On("ListActive", ctx).Return([]model.CollectionEntity{}, assert.AnError)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePaymentWithFallback(ctx, routingCtx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoutingEngine_ValidateRoutingConfiguration_Success(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}
	channels := []model.PaymentChannelConfig{
		{MerchantNo: "MERCHANT001", Enabled: true},
	}

	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("ListChannelsByEntity", ctx, uint64(1)).Return(channels, nil)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.NoError(t, err)
}

func TestRoutingEngine_ValidateRoutingConfiguration_NoDefault(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), repository.ErrNotFound)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.Error(t, err)
}

func TestRoutingEngine_ValidateRoutingConfiguration_InactiveDefault(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	inactiveDefault := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusInactive,
	}

	mockEntityRepo.On("GetDefault", ctx).Return(inactiveDefault, nil)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.Error(t, err)
}

func TestRoutingEngine_ValidateRoutingConfiguration_NoEnabledChannel(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}
	channels := []model.PaymentChannelConfig{
		{MerchantNo: "MERCHANT001", Enabled: false},
	}

	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("ListChannelsByEntity", ctx, uint64(1)).Return(channels, nil)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.Error(t, err)
}

func TestRoutingEngine_ValidateRoutingConfiguration_ListChannelsError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:   model.Base{ID: 1},
		Status: model.EntityStatusActive,
	}

	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("ListChannelsByEntity", ctx, uint64(1)).Return([]model.PaymentChannelConfig{}, assert.AnError)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.Error(t, err)
}

func TestRoutingEngine_ValidateRoutingConfiguration_GetDefaultError(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	mockEntityRepo.On("GetDefault", ctx).Return((*model.CollectionEntity)(nil), assert.AnError)

	err := engine.ValidateRoutingConfiguration(ctx)

	assert.Error(t, err)
}

func TestRoutingEngine_GetMerchantNo_DisabledChannel(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	conditions := []model.RoutingCondition{
		{Field: model.ConditionFieldGameType, Operator: model.ConditionOperatorEquals, Value: mustJSON("LOL")},
	}
	conditionsJSON, _ := json.Marshal(conditions)

	rules := []model.RoutingRule{
		{Base: model.Base{ID: 1}, Conditions: conditionsJSON, TargetEntityID: 10, Status: model.RuleStatusActive},
	}

	targetEntity := &model.CollectionEntity{Base: model.Base{ID: 10}, Status: model.EntityStatusActive}
	disabledChannel := &model.PaymentChannelConfig{MerchantNo: "MERCHANT", Enabled: false}
	defaultEntity := &model.CollectionEntity{Base: model.Base{ID: 1}, Status: model.EntityStatusActive}
	enabledChannel := &model.PaymentChannelConfig{MerchantNo: "DEFAULT", Enabled: true}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return(rules, nil)
	mockEntityRepo.On("Get", ctx, uint64(10)).Return(targetEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(10), model.PaymentMethodWeChat).Return(disabledChannel, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return(enabledChannel, nil)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.True(t, result.IsDefault)
}

func TestRoutingEngine_GetAnyMerchantNo_DisabledChannels(t *testing.T) {
	ctx := context.Background()
	mockRuleRepo := new(MockRoutingRuleRepository)
	mockEntityRepo := new(MockCollectionEntityRepository)
	engine := NewRoutingEngine(mockRuleRepo, mockEntityRepo)

	defaultEntity := &model.CollectionEntity{
		Base:      model.Base{ID: 1},
		Status:    model.EntityStatusActive,
		IsDefault: true,
		PaymentChannels: []model.PaymentChannelConfig{
			{MerchantNo: "DISABLED1", Enabled: false},
			{MerchantNo: "", Enabled: true}, // Empty merchant no
			{MerchantNo: "ENABLED", Enabled: true},
		},
	}

	mockRuleRepo.On("ListActiveByPriority", ctx).Return([]model.RoutingRule{}, nil)
	mockEntityRepo.On("GetDefault", ctx).Return(defaultEntity, nil)
	mockEntityRepo.On("GetChannelByEntityAndMethod", ctx, uint64(1), model.PaymentMethodWeChat).Return((*model.PaymentChannelConfig)(nil), repository.ErrNotFound)

	routingCtx := &RoutingContext{GameType: "LOL", Method: model.PaymentMethodWeChat}
	result, err := engine.RoutePayment(ctx, routingCtx)

	assert.NoError(t, err)
	assert.Equal(t, "ENABLED", result.MerchantNo)
}

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

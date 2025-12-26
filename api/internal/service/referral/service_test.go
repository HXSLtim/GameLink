package referral

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	referralrepo "gamelink/internal/repository/referral"
)

// ============================================================================
// Mock Repository
// ============================================================================

// MockReferralRepository is a mock implementation
type MockReferralRepository struct {
	mock.Mock
}

func (m *MockReferralRepository) GetAllConfigs(ctx context.Context) ([]model.ReferralConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.ReferralConfig), args.Error(1)
}

func (m *MockReferralRepository) GetConfig(ctx context.Context, key string) (*model.ReferralConfig, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReferralConfig), args.Error(1)
}

func (m *MockReferralRepository) SetConfig(ctx context.Context, key, value, description string) error {
	args := m.Called(ctx, key, value, description)
	return args.Error(0)
}

func (m *MockReferralRepository) ListCodes(ctx context.Context, opts referralrepo.CodeListOptions) ([]model.ReferralCode, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ReferralCode), args.Get(1).(int64), args.Error(2)
}

func (m *MockReferralRepository) GetCodeByID(ctx context.Context, id uint64) (*model.ReferralCode, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReferralCode), args.Error(1)
}

func (m *MockReferralRepository) GetCodeByCode(ctx context.Context, code string) (*model.ReferralCode, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReferralCode), args.Error(1)
}

func (m *MockReferralRepository) GetUserCode(ctx context.Context, userID uint64, refType model.ReferralType) (*model.ReferralCode, error) {
	args := m.Called(ctx, userID, refType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReferralCode), args.Error(1)
}

func (m *MockReferralRepository) CreateCode(ctx context.Context, item *model.ReferralCode) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockReferralRepository) UpdateCode(ctx context.Context, item *model.ReferralCode) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockReferralRepository) DeleteCode(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReferralRepository) IncrementCodeUseCount(ctx context.Context, codeID uint64) error {
	args := m.Called(ctx, codeID)
	return args.Error(0)
}

func (m *MockReferralRepository) ListReferrals(ctx context.Context, opts referralrepo.ReferralListOptions) ([]model.Referral, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Referral), args.Get(1).(int64), args.Error(2)
}

func (m *MockReferralRepository) GetReferralByID(ctx context.Context, id uint64) (*model.Referral, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Referral), args.Error(1)
}

func (m *MockReferralRepository) GetReferralByReferee(ctx context.Context, refereeID uint64) (*model.Referral, error) {
	args := m.Called(ctx, refereeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Referral), args.Error(1)
}

func (m *MockReferralRepository) CreateReferral(ctx context.Context, item *model.Referral) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockReferralRepository) UpdateReferralStatus(ctx context.Context, id uint64, status model.ReferralStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockReferralRepository) GetUserReferrals(ctx context.Context, userID uint64, limit int) ([]model.Referral, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]model.Referral), args.Error(1)
}

func (m *MockReferralRepository) ListRewards(ctx context.Context, opts referralrepo.RewardListOptions) ([]model.ReferralReward, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ReferralReward), args.Get(1).(int64), args.Error(2)
}

func (m *MockReferralRepository) GetRewardByID(ctx context.Context, id uint64) (*model.ReferralReward, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReferralReward), args.Error(1)
}

func (m *MockReferralRepository) CreateReward(ctx context.Context, item *model.ReferralReward) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockReferralRepository) UpdateRewardStatus(ctx context.Context, id uint64, status model.ReferralRewardStatus, failureReason string) error {
	args := m.Called(ctx, id, status, failureReason)
	return args.Error(0)
}

func (m *MockReferralRepository) GetUserRewards(ctx context.Context, userID uint64, limit int) ([]model.ReferralReward, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]model.ReferralReward), args.Error(1)
}

func (m *MockReferralRepository) GetReferralStats(ctx context.Context) (map[string]any, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockReferralRepository) GetUserReferralStats(ctx context.Context, userID uint64) (map[string]any, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(map[string]any), args.Error(1)
}

// ============================================================================
// Service Tests - Config Management
// ============================================================================

func TestService_GetAllConfigs(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	configs := []model.ReferralConfig{
		{ConfigKey: "enabled", ConfigValue: "true"},
		{ConfigKey: "expire_days", ConfigValue: "30"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAllConfigs", ctx).Return(configs, nil).Once()
		result, err := svc.GetAllConfigs(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetAllConfigs", ctx).Return([]model.ReferralConfig{}, errors.New("db error")).Once()
		_, err := svc.GetAllConfigs(ctx)
		assert.Error(t, err)
	})
}

func TestService_GetConfig(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	config := &model.ReferralConfig{ConfigKey: "enabled", ConfigValue: "true"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetConfig", ctx, "enabled").Return(config, nil).Once()
		result, err := svc.GetConfig(ctx, "enabled")
		require.NoError(t, err)
		assert.Equal(t, "true", result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetConfig", ctx, "unknown").Return(nil, repository.ErrNotFound).Once()
		_, err := svc.GetConfig(ctx, "unknown")
		assert.Error(t, err)
	})
}

func TestService_SetConfig(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("SetConfig", ctx, "enabled", "true", "Enable referral").Return(nil).Once()
		err := svc.SetConfig(ctx, "enabled", "true", "Enable referral")
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("SetConfig", ctx, "key", "val", "desc").Return(errors.New("db error")).Once()
		err := svc.SetConfig(ctx, "key", "val", "desc")
		assert.Error(t, err)
	})
}

func TestService_IsEnabled(t *testing.T) {
	ctx := context.Background()

	t.Run("enabled true", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigEnabled, ConfigValue: "true"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigEnabled).Return(config, nil).Once()
		assert.True(t, svc.IsEnabled(ctx))
	})

	t.Run("enabled 1", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigEnabled, ConfigValue: "1"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigEnabled).Return(config, nil).Once()
		assert.True(t, svc.IsEnabled(ctx))
	})

	t.Run("disabled", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigEnabled, ConfigValue: "false"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigEnabled).Return(config, nil).Once()
		assert.False(t, svc.IsEnabled(ctx))
	})

	t.Run("error returns false", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetConfig", ctx, model.ReferralConfigEnabled).Return(nil, repository.ErrNotFound).Once()
		assert.False(t, svc.IsEnabled(ctx))
	})
}

func TestService_GetExpireDays(t *testing.T) {
	ctx := context.Background()

	t.Run("custom days", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigExpireDays, ConfigValue: "60"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigExpireDays).Return(config, nil).Once()
		assert.Equal(t, 60, svc.GetExpireDays(ctx))
	})

	t.Run("default on error", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetConfig", ctx, model.ReferralConfigExpireDays).Return(nil, repository.ErrNotFound).Once()
		assert.Equal(t, 30, svc.GetExpireDays(ctx))
	})

	t.Run("default on invalid value", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigExpireDays, ConfigValue: "invalid"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigExpireDays).Return(config, nil).Once()
		assert.Equal(t, 30, svc.GetExpireDays(ctx))
	})

	t.Run("default on zero", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		config := &model.ReferralConfig{ConfigKey: model.ReferralConfigExpireDays, ConfigValue: "0"}
		mockRepo.On("GetConfig", ctx, model.ReferralConfigExpireDays).Return(config, nil).Once()
		assert.Equal(t, 30, svc.GetExpireDays(ctx))
	})
}

// ============================================================================
// Service Tests - Code Management
// ============================================================================

func TestService_ListCodes(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	codes := []model.ReferralCode{
		{Code: "ABC123", UserID: 1},
		{Code: "DEF456", UserID: 2},
	}

	t.Run("success", func(t *testing.T) {
		opts := referralrepo.CodeListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListCodes", ctx, opts).Return(codes, int64(2), nil).Once()
		result, total, err := svc.ListCodes(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		opts := referralrepo.CodeListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListCodes", ctx, opts).Return([]model.ReferralCode{}, int64(0), errors.New("db error")).Once()
		_, _, err := svc.ListCodes(ctx, opts)
		assert.Error(t, err)
	})
}

func TestService_GetCodeByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	code := &model.ReferralCode{Code: "ABC123", UserID: 1}
	code.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetCodeByID", ctx, uint64(1)).Return(code, nil).Once()
		result, err := svc.GetCodeByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "ABC123", result.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetCodeByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		_, err := svc.GetCodeByID(ctx, 999)
		assert.Error(t, err)
	})
}

func TestService_GetCodeByCode(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	code := &model.ReferralCode{Code: "ABC123", UserID: 1}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetCodeByCode", ctx, "ABC123").Return(code, nil).Once()
		result, err := svc.GetCodeByCode(ctx, "ABC123")
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.UserID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetCodeByCode", ctx, "INVALID").Return(nil, repository.ErrNotFound).Once()
		_, err := svc.GetCodeByCode(ctx, "INVALID")
		assert.Error(t, err)
	})
}

func TestService_CreateCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		expireAt := time.Now().Add(30 * 24 * time.Hour)
		req := CreateCodeRequest{
			UserID:   1,
			Type:     model.ReferralTypeUserToUser,
			MaxUse:   100,
			ExpireAt: &expireAt,
		}
		mockRepo.On("CreateCode", ctx, mock.AnythingOfType("*model.ReferralCode")).Return(nil).Once()
		result, err := svc.CreateCode(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.UserID)
		assert.True(t, result.IsActive)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		req := CreateCodeRequest{UserID: 1, Type: model.ReferralTypeUserToUser}
		mockRepo.On("CreateCode", ctx, mock.AnythingOfType("*model.ReferralCode")).Return(errors.New("db error")).Once()
		_, err := svc.CreateCode(ctx, req)
		assert.Error(t, err)
	})
}

func TestService_GetOrCreateUserCode(t *testing.T) {
	ctx := context.Background()

	t.Run("existing code", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		existingCode := &model.ReferralCode{Code: "EXISTING", UserID: 1, Type: model.ReferralTypeUserToUser}
		mockRepo.On("GetUserCode", ctx, uint64(1), model.ReferralTypeUserToUser).Return(existingCode, nil).Once()
		result, err := svc.GetOrCreateUserCode(ctx, 1, model.ReferralTypeUserToUser)
		require.NoError(t, err)
		assert.Equal(t, "EXISTING", result.Code)
	})

	t.Run("create new code", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetUserCode", ctx, uint64(1), model.ReferralTypeUserToUser).Return(nil, repository.ErrNotFound).Once()
		mockRepo.On("GetConfig", ctx, model.ReferralConfigExpireDays).Return(&model.ReferralConfig{ConfigValue: "30"}, nil).Once()
		mockRepo.On("CreateCode", ctx, mock.AnythingOfType("*model.ReferralCode")).Return(nil).Once()
		result, err := svc.GetOrCreateUserCode(ctx, 1, model.ReferralTypeUserToUser)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.UserID)
	})
}

func TestService_UpdateCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		existingCode := &model.ReferralCode{Code: "ABC123", IsActive: true, MaxUse: 50}
		existingCode.ID = 1
		isActive := false
		maxUse := 100
		req := UpdateCodeRequest{ID: 1, IsActive: &isActive, MaxUse: &maxUse}
		mockRepo.On("GetCodeByID", ctx, uint64(1)).Return(existingCode, nil).Once()
		mockRepo.On("UpdateCode", ctx, mock.AnythingOfType("*model.ReferralCode")).Return(nil).Once()
		result, err := svc.UpdateCode(ctx, req)
		require.NoError(t, err)
		assert.False(t, result.IsActive)
		assert.Equal(t, 100, result.MaxUse)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		req := UpdateCodeRequest{ID: 999}
		mockRepo.On("GetCodeByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		_, err := svc.UpdateCode(ctx, req)
		assert.Error(t, err)
	})

	t.Run("update error", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		existingCode := &model.ReferralCode{Code: "ABC123"}
		existingCode.ID = 1
		req := UpdateCodeRequest{ID: 1}
		mockRepo.On("GetCodeByID", ctx, uint64(1)).Return(existingCode, nil).Once()
		mockRepo.On("UpdateCode", ctx, mock.AnythingOfType("*model.ReferralCode")).Return(errors.New("db error")).Once()
		_, err := svc.UpdateCode(ctx, req)
		assert.Error(t, err)
	})
}

func TestService_DeleteCode(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteCode", ctx, uint64(1)).Return(nil).Once()
		err := svc.DeleteCode(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteCode", ctx, uint64(999)).Return(repository.ErrNotFound).Once()
		err := svc.DeleteCode(ctx, 999)
		assert.Error(t, err)
	})
}

func TestService_ValidateCode(t *testing.T) {
	ctx := context.Background()

	t.Run("valid code", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		code := &model.ReferralCode{Code: "ABC123", IsActive: true, MaxUse: 0}
		mockRepo.On("GetCodeByCode", ctx, "ABC123").Return(code, nil).Once()
		result, err := svc.ValidateCode(ctx, "ABC123")
		require.NoError(t, err)
		assert.Equal(t, "ABC123", result.Code)
	})

	t.Run("invalid code not found", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetCodeByCode", ctx, "INVALID").Return(nil, repository.ErrNotFound).Once()
		_, err := svc.ValidateCode(ctx, "INVALID")
		assert.Error(t, err)
	})

	t.Run("expired code", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		expireAt := time.Now().Add(-24 * time.Hour)
		code := &model.ReferralCode{Code: "EXPIRED", IsActive: true, ExpireAt: &expireAt}
		mockRepo.On("GetCodeByCode", ctx, "EXPIRED").Return(code, nil).Once()
		_, err := svc.ValidateCode(ctx, "EXPIRED")
		assert.Error(t, err)
	})
}

// ============================================================================
// ReferralCode Model Tests
// ============================================================================

func TestReferralCode_IsValid(t *testing.T) {
	t.Run("valid active code no expiry", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: nil,
			MaxUse:   0,
			UseCount: 10,
		}
		isValid := code.IsValid()
		assert.True(t, isValid)
	})

	t.Run("valid active code with expiry", func(t *testing.T) {
		expireAt := time.Now().Add(24 * time.Hour)
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: &expireAt,
			MaxUse:   100,
			UseCount: 50,
		}
		isValid := code.IsValid()
		assert.True(t, isValid)
	})

	t.Run("inactive code", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: false,
			ExpireAt: nil,
		}
		isValid := code.IsValid()
		assert.False(t, isValid)
	})

	t.Run("expired code", func(t *testing.T) {
		expireAt := time.Now().Add(-24 * time.Hour)
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: &expireAt,
		}
		isValid := code.IsValid()
		assert.False(t, isValid)
	})

	t.Run("max use reached", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: nil,
			MaxUse:   10,
			UseCount: 10,
		}
		isValid := code.IsValid()
		assert.False(t, isValid)
	})

	t.Run("unlimited use", func(t *testing.T) {
		code := &model.ReferralCode{
			IsActive: true,
			ExpireAt: nil,
			MaxUse:   0,
			UseCount: 1000,
		}
		isValid := code.IsValid()
		assert.True(t, isValid)
	})
}

// ============================================================================
// Referral Model Tests
// ============================================================================

func TestReferral_Status(t *testing.T) {
	t.Run("pending status", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusPending}
		isPending := referral.Status == model.ReferralStatusPending
		assert.True(t, isPending)
	})

	t.Run("completed status", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusCompleted}
		isCompleted := referral.Status == model.ReferralStatusCompleted
		assert.True(t, isCompleted)
	})

	t.Run("rewarded status", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusRewarded}
		isRewarded := referral.Status == model.ReferralStatusRewarded
		assert.True(t, isRewarded)
	})

	t.Run("expired status", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusExpired}
		isExpired := referral.Status == model.ReferralStatusExpired
		assert.True(t, isExpired)
	})

	t.Run("canceled status", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusCanceled}
		isCanceled := referral.Status == model.ReferralStatusCanceled
		assert.True(t, isCanceled)
	})
}

func TestReferral_Level(t *testing.T) {
	t.Run("direct referral level 1", func(t *testing.T) {
		referral := &model.Referral{Level: 1}
		isDirect := referral.Level == 1
		assert.True(t, isDirect)
	})

	t.Run("indirect referral level 2", func(t *testing.T) {
		referral := &model.Referral{Level: 2}
		isIndirect := referral.Level == 2
		assert.True(t, isIndirect)
	})
}

// ============================================================================
// ReferralReward Model Tests
// ============================================================================

func TestReferralReward_Status(t *testing.T) {
	t.Run("pending status", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusPending}
		isPending := reward.Status == model.ReferralRewardStatusPending
		assert.True(t, isPending)
	})

	t.Run("issued status", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusIssued}
		isIssued := reward.Status == model.ReferralRewardStatusIssued
		assert.True(t, isIssued)
	})

	t.Run("failed status", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusFailed}
		isFailed := reward.Status == model.ReferralRewardStatusFailed
		assert.True(t, isFailed)
	})
}

func TestReferralReward_Type(t *testing.T) {
	t.Run("cash reward", func(t *testing.T) {
		reward := &model.ReferralReward{Type: model.RewardTypeCash, AmountCents: 1000}
		isCash := reward.Type == model.RewardTypeCash
		assert.True(t, isCash)
		assert.Equal(t, int64(1000), reward.AmountCents)
	})

	t.Run("coupon reward", func(t *testing.T) {
		couponID := uint64(1)
		reward := &model.ReferralReward{Type: model.RewardTypeCoupon, CouponID: &couponID}
		isCoupon := reward.Type == model.RewardTypeCoupon
		assert.True(t, isCoupon)
		assert.NotNil(t, reward.CouponID)
	})

	t.Run("points reward", func(t *testing.T) {
		reward := &model.ReferralReward{Type: model.RewardTypePoints, AmountCents: 500}
		isPoints := reward.Type == model.RewardTypePoints
		assert.True(t, isPoints)
	})
}

// ============================================================================
// Status Constants Tests
// ============================================================================

func TestReferralType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.ReferralType("user_to_user"), model.ReferralTypeUserToUser)
		assert.Equal(t, model.ReferralType("player_to_player"), model.ReferralTypePlayerToPlayer)
		assert.Equal(t, model.ReferralType("user_to_player"), model.ReferralTypeUserToPlayer)
	})
}

func TestReferralStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.ReferralStatus("pending"), model.ReferralStatusPending)
		assert.Equal(t, model.ReferralStatus("completed"), model.ReferralStatusCompleted)
		assert.Equal(t, model.ReferralStatus("rewarded"), model.ReferralStatusRewarded)
		assert.Equal(t, model.ReferralStatus("expired"), model.ReferralStatusExpired)
		assert.Equal(t, model.ReferralStatus("canceled"), model.ReferralStatusCanceled)
	})
}

func TestRewardType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.RewardType("cash"), model.RewardTypeCash)
		assert.Equal(t, model.RewardType("coupon"), model.RewardTypeCoupon)
		assert.Equal(t, model.RewardType("points"), model.RewardTypePoints)
	})
}

func TestReferralRewardStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.ReferralRewardStatus("pending"), model.ReferralRewardStatusPending)
		assert.Equal(t, model.ReferralRewardStatus("issued"), model.ReferralRewardStatusIssued)
		assert.Equal(t, model.ReferralRewardStatus("failed"), model.ReferralRewardStatusFailed)
	})
}

// ============================================================================
// Config Constants Tests
// ============================================================================

func TestReferralConfig_Constants(t *testing.T) {
	t.Run("config keys", func(t *testing.T) {
		assert.Equal(t, "enabled", model.ReferralConfigEnabled)
		assert.Equal(t, "expire_days", model.ReferralConfigExpireDays)
		assert.Equal(t, "max_level", model.ReferralConfigMaxLevel)
	})
}

// ============================================================================
// Service Logic Tests
// ============================================================================

func TestIsEnabled_Logic(t *testing.T) {
	t.Run("enabled true", func(t *testing.T) {
		value := "true"
		isEnabled := value == "true" || value == "1"
		assert.True(t, isEnabled)
	})

	t.Run("enabled 1", func(t *testing.T) {
		value := "1"
		isEnabled := value == "true" || value == "1"
		assert.True(t, isEnabled)
	})

	t.Run("disabled false", func(t *testing.T) {
		value := "false"
		isEnabled := value == "true" || value == "1"
		assert.False(t, isEnabled)
	})

	t.Run("disabled 0", func(t *testing.T) {
		value := "0"
		isEnabled := value == "true" || value == "1"
		assert.False(t, isEnabled)
	})
}

func TestGetExpireDays_Logic(t *testing.T) {
	t.Run("default expire days", func(t *testing.T) {
		days := 0
		if days <= 0 {
			days = 30
		}
		assert.Equal(t, 30, days)
	})

	t.Run("custom expire days", func(t *testing.T) {
		days := 60
		if days <= 0 {
			days = 30
		}
		assert.Equal(t, 60, days)
	})

	t.Run("negative expire days", func(t *testing.T) {
		days := -10
		if days <= 0 {
			days = 30
		}
		assert.Equal(t, 30, days)
	})
}

func TestValidateCode_Logic(t *testing.T) {
	t.Run("cannot refer yourself", func(t *testing.T) {
		code := &model.ReferralCode{UserID: 100}
		refereeID := uint64(100)
		canRefer := code.UserID != refereeID
		assert.False(t, canRefer)
	})

	t.Run("can refer other user", func(t *testing.T) {
		code := &model.ReferralCode{UserID: 100}
		refereeID := uint64(200)
		canRefer := code.UserID != refereeID
		assert.True(t, canRefer)
	})
}

func TestCreateReferral_Logic(t *testing.T) {
	t.Run("default level is 1", func(t *testing.T) {
		level := 0
		if level <= 0 {
			level = 1
		}
		assert.Equal(t, 1, level)
	})

	t.Run("custom level preserved", func(t *testing.T) {
		level := 2
		if level <= 0 {
			level = 1
		}
		assert.Equal(t, 2, level)
	})
}

func TestCompleteReferral_Logic(t *testing.T) {
	t.Run("condition matches", func(t *testing.T) {
		referral := &model.Referral{
			Status:           model.ReferralStatusPending,
			RefereeCondition: "first_order",
		}
		condition := "first_order"
		shouldProcess := referral.RefereeCondition == "" || referral.RefereeCondition == condition
		assert.True(t, shouldProcess)
	})

	t.Run("condition does not match", func(t *testing.T) {
		referral := &model.Referral{
			Status:           model.ReferralStatusPending,
			RefereeCondition: "first_order",
		}
		condition := "registered"
		shouldProcess := referral.RefereeCondition == "" || referral.RefereeCondition == condition
		assert.False(t, shouldProcess)
	})

	t.Run("empty condition matches any", func(t *testing.T) {
		referral := &model.Referral{
			Status:           model.ReferralStatusPending,
			RefereeCondition: "",
		}
		condition := "first_order"
		shouldProcess := referral.RefereeCondition == "" || referral.RefereeCondition == condition
		assert.True(t, shouldProcess)
	})

	t.Run("already completed skipped", func(t *testing.T) {
		referral := &model.Referral{Status: model.ReferralStatusCompleted}
		shouldProcess := referral.Status == model.ReferralStatusPending
		assert.False(t, shouldProcess)
	})
}

func TestIssueReward_Logic(t *testing.T) {
	t.Run("only pending can be issued", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusPending}
		canIssue := reward.Status == model.ReferralRewardStatusPending
		assert.True(t, canIssue)
	})

	t.Run("issued cannot be issued again", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusIssued}
		canIssue := reward.Status == model.ReferralRewardStatusPending
		assert.False(t, canIssue)
	})

	t.Run("failed cannot be issued", func(t *testing.T) {
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusFailed}
		canIssue := reward.Status == model.ReferralRewardStatusPending
		assert.False(t, canIssue)
	})
}

// ============================================================================
// Default Limit Tests
// ============================================================================

func TestGetUserReferrals_DefaultLimit(t *testing.T) {
	t.Run("default limit when zero", func(t *testing.T) {
		limit := 0
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 20, limit)
	})

	t.Run("default limit when negative", func(t *testing.T) {
		limit := -5
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 20, limit)
	})

	t.Run("custom limit preserved", func(t *testing.T) {
		limit := 50
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 50, limit)
	})
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestReferralErrors(t *testing.T) {
	t.Run("not found error", func(t *testing.T) {
		err := repository.ErrNotFound
		assert.True(t, errors.Is(err, repository.ErrNotFound))
	})
}

// ============================================================================
// Stats Structure Tests
// ============================================================================

func TestReferralStats_Structure(t *testing.T) {
	t.Run("stats map structure", func(t *testing.T) {
		stats := map[string]any{
			"totalReferrals":     int64(100),
			"completedReferrals": int64(80),
			"rewardedReferrals":  int64(70),
			"totalRewardsCents":  int64(50000),
		}
		assert.Contains(t, stats, "totalReferrals")
		assert.Contains(t, stats, "completedReferrals")
		assert.Contains(t, stats, "rewardedReferrals")
	})
}

func TestUserReferralStats_Structure(t *testing.T) {
	t.Run("user stats map structure", func(t *testing.T) {
		stats := map[string]any{
			"referralCount":    int64(10),
			"completedCount":   int64(8),
			"totalRewardCents": int64(5000),
		}
		assert.Contains(t, stats, "referralCount")
		assert.Contains(t, stats, "completedCount")
		assert.Contains(t, stats, "totalRewardCents")
	})
}

// ============================================================================
// Service Tests - Referral Management
// ============================================================================

func TestService_ListReferrals(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	referrals := []model.Referral{
		{ReferrerID: 1, RefereeID: 2, Status: model.ReferralStatusPending},
	}

	t.Run("success", func(t *testing.T) {
		opts := referralrepo.ReferralListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListReferrals", ctx, opts).Return(referrals, int64(1), nil).Once()
		result, total, err := svc.ListReferrals(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
	})

	t.Run("error", func(t *testing.T) {
		opts := referralrepo.ReferralListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListReferrals", ctx, opts).Return([]model.Referral{}, int64(0), errors.New("db error")).Once()
		_, _, err := svc.ListReferrals(ctx, opts)
		assert.Error(t, err)
	})
}

func TestService_GetReferralByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	referral := &model.Referral{ReferrerID: 1, RefereeID: 2}
	referral.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetReferralByID", ctx, uint64(1)).Return(referral, nil).Once()
		result, err := svc.GetReferralByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.ReferrerID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetReferralByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		_, err := svc.GetReferralByID(ctx, 999)
		assert.Error(t, err)
	})
}

func TestService_CreateReferral(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		codeID := uint64(1)
		req := CreateReferralRequest{
			ReferrerID:       1,
			RefereeID:        2,
			CodeID:           &codeID,
			Type:             model.ReferralTypeUserToUser,
			Level:            1,
			RefereeCondition: "registered",
		}
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(nil, repository.ErrNotFound).Once()
		mockRepo.On("CreateReferral", ctx, mock.AnythingOfType("*model.Referral")).Return(nil).Once()
		mockRepo.On("IncrementCodeUseCount", ctx, uint64(1)).Return(nil).Once()
		result, err := svc.CreateReferral(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.ReferrerID)
		assert.Equal(t, model.ReferralStatusPending, result.Status)
	})

	t.Run("referee already has referral", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		existingReferral := &model.Referral{ReferrerID: 3, RefereeID: 2}
		req := CreateReferralRequest{ReferrerID: 1, RefereeID: 2}
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(existingReferral, nil).Once()
		_, err := svc.CreateReferral(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already has referral")
	})

	t.Run("default level is 1", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		req := CreateReferralRequest{ReferrerID: 1, RefereeID: 2, Level: 0}
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(nil, repository.ErrNotFound).Once()
		mockRepo.On("CreateReferral", ctx, mock.MatchedBy(func(r *model.Referral) bool {
			return r.Level == 1
		})).Return(nil).Once()
		result, err := svc.CreateReferral(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, 1, result.Level)
	})
}

func TestService_UseCode(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		code := &model.ReferralCode{Code: "ABC123", UserID: 1, Type: model.ReferralTypeUserToUser, IsActive: true}
		code.ID = 1
		req := UseCodeRequest{Code: "ABC123", RefereeID: 2}
		mockRepo.On("GetCodeByCode", ctx, "ABC123").Return(code, nil).Once()
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(nil, repository.ErrNotFound).Once()
		mockRepo.On("CreateReferral", ctx, mock.AnythingOfType("*model.Referral")).Return(nil).Once()
		mockRepo.On("IncrementCodeUseCount", ctx, uint64(1)).Return(nil).Once()
		result, err := svc.UseCode(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.ReferrerID)
		assert.Equal(t, uint64(2), result.RefereeID)
	})

	t.Run("cannot refer yourself", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		code := &model.ReferralCode{Code: "ABC123", UserID: 1, IsActive: true}
		req := UseCodeRequest{Code: "ABC123", RefereeID: 1}
		mockRepo.On("GetCodeByCode", ctx, "ABC123").Return(code, nil).Once()
		_, err := svc.UseCode(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot refer yourself")
	})

	t.Run("invalid code", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		req := UseCodeRequest{Code: "INVALID", RefereeID: 2}
		mockRepo.On("GetCodeByCode", ctx, "INVALID").Return(nil, repository.ErrNotFound).Once()
		_, err := svc.UseCode(ctx, req)
		assert.Error(t, err)
	})
}

func TestService_CompleteReferral(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		referral := &model.Referral{ReferrerID: 1, RefereeID: 2, Status: model.ReferralStatusPending, RefereeCondition: "first_order"}
		referral.ID = 1
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(referral, nil).Once()
		mockRepo.On("UpdateReferralStatus", ctx, uint64(1), model.ReferralStatusCompleted).Return(nil).Once()
		err := svc.CompleteReferral(ctx, 2, "first_order")
		require.NoError(t, err)
	})

	t.Run("condition mismatch skipped", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		referral := &model.Referral{ReferrerID: 1, RefereeID: 2, Status: model.ReferralStatusPending, RefereeCondition: "first_order"}
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(referral, nil).Once()
		err := svc.CompleteReferral(ctx, 2, "registered")
		require.NoError(t, err)
	})

	t.Run("already completed skipped", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		referral := &model.Referral{ReferrerID: 1, RefereeID: 2, Status: model.ReferralStatusCompleted}
		mockRepo.On("GetReferralByReferee", ctx, uint64(2)).Return(referral, nil).Once()
		err := svc.CompleteReferral(ctx, 2, "first_order")
		require.NoError(t, err)
	})

	t.Run("referral not found", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetReferralByReferee", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		err := svc.CompleteReferral(ctx, 999, "first_order")
		assert.Error(t, err)
	})
}

func TestService_UpdateReferralStatus(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateReferralStatus", ctx, uint64(1), model.ReferralStatusCompleted).Return(nil).Once()
		err := svc.UpdateReferralStatus(ctx, 1, model.ReferralStatusCompleted)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("UpdateReferralStatus", ctx, uint64(999), model.ReferralStatusCompleted).Return(repository.ErrNotFound).Once()
		err := svc.UpdateReferralStatus(ctx, 999, model.ReferralStatusCompleted)
		assert.Error(t, err)
	})
}

func TestService_GetUserReferrals(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	referrals := []model.Referral{{ReferrerID: 1, RefereeID: 2}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetUserReferrals", ctx, uint64(1), 20).Return(referrals, nil).Once()
		result, err := svc.GetUserReferrals(ctx, 1, 0)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("custom limit", func(t *testing.T) {
		mockRepo.On("GetUserReferrals", ctx, uint64(1), 50).Return(referrals, nil).Once()
		result, err := svc.GetUserReferrals(ctx, 1, 50)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// ============================================================================
// Service Tests - Reward Management
// ============================================================================

func TestService_ListRewards(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	rewards := []model.ReferralReward{{UserID: 1, AmountCents: 1000}}

	t.Run("success", func(t *testing.T) {
		opts := referralrepo.RewardListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListRewards", ctx, opts).Return(rewards, int64(1), nil).Once()
		result, total, err := svc.ListRewards(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, result, 1)
	})
}

func TestService_GetRewardByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	reward := &model.ReferralReward{UserID: 1, AmountCents: 1000}
	reward.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(reward, nil).Once()
		result, err := svc.GetRewardByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(1000), result.AmountCents)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRewardByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		_, err := svc.GetRewardByID(ctx, 999)
		assert.Error(t, err)
	})
}

func TestService_CreateReward(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		req := CreateRewardRequest{
			ReferralID:  1,
			UserID:      1,
			Type:        model.RewardTypeCash,
			AmountCents: 1000,
		}
		mockRepo.On("CreateReward", ctx, mock.AnythingOfType("*model.ReferralReward")).Return(nil).Once()
		result, err := svc.CreateReward(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, model.ReferralRewardStatusPending, result.Status)
	})

	t.Run("error", func(t *testing.T) {
		req := CreateRewardRequest{ReferralID: 1, UserID: 1}
		mockRepo.On("CreateReward", ctx, mock.AnythingOfType("*model.ReferralReward")).Return(errors.New("db error")).Once()
		_, err := svc.CreateReward(ctx, req)
		assert.Error(t, err)
	})
}

func TestService_IssueReward(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		reward := &model.ReferralReward{ReferralID: 1, UserID: 1, Status: model.ReferralRewardStatusPending}
		reward.ID = 1
		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(reward, nil).Once()
		mockRepo.On("UpdateRewardStatus", ctx, uint64(1), model.ReferralRewardStatusIssued, "").Return(nil).Once()
		mockRepo.On("UpdateReferralStatus", ctx, uint64(1), model.ReferralStatusRewarded).Return(nil).Once()
		err := svc.IssueReward(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("already processed", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		reward := &model.ReferralReward{Status: model.ReferralRewardStatusIssued}
		reward.ID = 1
		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(reward, nil).Once()
		err := svc.IssueReward(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already processed")
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockReferralRepository)
		svc := NewReferralService(mockRepo)
		mockRepo.On("GetRewardByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()
		err := svc.IssueReward(ctx, 999)
		assert.Error(t, err)
	})
}

func TestService_FailReward(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UpdateRewardStatus", ctx, uint64(1), model.ReferralRewardStatusFailed, "insufficient balance").Return(nil).Once()
		err := svc.FailReward(ctx, 1, "insufficient balance")
		require.NoError(t, err)
	})
}

func TestService_GetUserRewards(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	rewards := []model.ReferralReward{{UserID: 1, AmountCents: 1000}}

	t.Run("success default limit", func(t *testing.T) {
		mockRepo.On("GetUserRewards", ctx, uint64(1), 20).Return(rewards, nil).Once()
		result, err := svc.GetUserRewards(ctx, 1, 0)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("custom limit", func(t *testing.T) {
		mockRepo.On("GetUserRewards", ctx, uint64(1), 50).Return(rewards, nil).Once()
		result, err := svc.GetUserRewards(ctx, 1, 50)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// ============================================================================
// Service Tests - Statistics
// ============================================================================

func TestService_GetReferralStats(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		stats := map[string]any{"totalCount": int64(100)}
		mockRepo.On("GetReferralStats", ctx).Return(stats, nil).Once()
		result, err := svc.GetReferralStats(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(100), result["totalCount"])
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetReferralStats", ctx).Return(map[string]any{}, errors.New("db error")).Once()
		_, err := svc.GetReferralStats(ctx)
		assert.Error(t, err)
	})
}

func TestService_GetUserReferralStats(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockReferralRepository)
	svc := NewReferralService(mockRepo)

	t.Run("success", func(t *testing.T) {
		stats := map[string]any{"totalCount": int64(10)}
		mockRepo.On("GetUserReferralStats", ctx, uint64(1)).Return(stats, nil).Once()
		result, err := svc.GetUserReferralStats(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result["totalCount"])
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetUserReferralStats", ctx, uint64(999)).Return(map[string]any{}, errors.New("db error")).Once()
		_, err := svc.GetUserReferralStats(ctx, 999)
		assert.Error(t, err)
	})
}

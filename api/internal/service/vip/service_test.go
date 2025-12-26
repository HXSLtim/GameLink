package vip

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// MockVipRepository is a mock implementation of VipRepository
type MockVipRepository struct {
	mock.Mock
}

func (m *MockVipRepository) CreateLevel(ctx context.Context, level *model.VipLevel) error {
	args := m.Called(ctx, level)
	return args.Error(0)
}

func (m *MockVipRepository) GetLevel(ctx context.Context, id uint64) (*model.VipLevel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) GetLevelBySlug(ctx context.Context, slug string) (*model.VipLevel, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) GetDefaultLevel(ctx context.Context) (*model.VipLevel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) ListLevels(ctx context.Context) ([]model.VipLevel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) ListActiveLevels(ctx context.Context) ([]model.VipLevel, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) ListLevelsPaged(ctx context.Context, opts repository.VipLevelListOptions) ([]model.VipLevel, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.VipLevel), args.Get(1).(int64), args.Error(2)
}

func (m *MockVipRepository) UpdateLevel(ctx context.Context, level *model.VipLevel) error {
	args := m.Called(ctx, level)
	return args.Error(0)
}

func (m *MockVipRepository) DeleteLevel(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVipRepository) SetDefaultLevel(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVipRepository) GetLevelByExp(ctx context.Context, exp int64) (*model.VipLevel, error) {
	args := m.Called(ctx, exp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VipLevel), args.Error(1)
}

func (m *MockVipRepository) BatchUpdateLevelStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	args := m.Called(ctx, ids, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVipRepository) BatchDeleteLevels(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVipRepository) GetConfig(ctx context.Context, key string) (*model.VipConfig, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VipConfig), args.Error(1)
}

func (m *MockVipRepository) ListConfigs(ctx context.Context) ([]model.VipConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.VipConfig), args.Error(1)
}

func (m *MockVipRepository) SaveConfig(ctx context.Context, config *model.VipConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockVipRepository) DeleteConfig(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestNewVipService(t *testing.T) {
	repo := &MockVipRepository{}
	svc := NewVipService(repo)
	require.NotNil(t, svc)
	assert.Equal(t, repo, svc.repo)
}

func TestService_CreateLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		level := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		repo.On("GetLevelBySlug", ctx, "vip1").Return(nil, repository.ErrNotFound)
		repo.On("CreateLevel", ctx, level).Return(nil)

		err := svc.CreateLevel(ctx, level)
		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("slug already exists", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		existing := &model.VipLevel{Slug: "vip1"}
		level := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		repo.On("GetLevelBySlug", ctx, "vip1").Return(existing, nil)

		err := svc.CreateLevel(ctx, level)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("repo error on check", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		level := &model.VipLevel{Slug: "vip1"}
		repo.On("GetLevelBySlug", ctx, "vip1").Return(nil, errors.New("db error"))

		err := svc.CreateLevel(ctx, level)
		require.Error(t, err)
	})
}

func TestService_GetLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		expected := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		expected.ID = 1
		repo.On("GetLevel", ctx, uint64(1)).Return(expected, nil)

		result, err := svc.GetLevel(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetLevel", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		_, err := svc.GetLevel(ctx, 999)
		require.Error(t, err)
	})
}

func TestService_ListLevels(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		levels := []model.VipLevel{
			{Slug: "vip1", Title: "VIP 1"},
			{Slug: "vip2", Title: "VIP 2"},
		}
		repo.On("ListLevels", ctx).Return(levels, nil)

		result, err := svc.ListLevels(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

func TestService_UpdateLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success same slug", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		existing := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		existing.ID = 1
		level := &model.VipLevel{Slug: "vip1", Title: "VIP 1 Updated"}
		level.ID = 1

		repo.On("GetLevel", ctx, uint64(1)).Return(existing, nil)
		repo.On("UpdateLevel", ctx, level).Return(nil)

		err := svc.UpdateLevel(ctx, level)
		require.NoError(t, err)
	})

	t.Run("success different slug", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		existing := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		existing.ID = 1
		level := &model.VipLevel{Slug: "vip2", Title: "VIP 2"}
		level.ID = 1

		repo.On("GetLevel", ctx, uint64(1)).Return(existing, nil)
		repo.On("GetLevelBySlug", ctx, "vip2").Return(nil, repository.ErrNotFound)
		repo.On("UpdateLevel", ctx, level).Return(nil)

		err := svc.UpdateLevel(ctx, level)
		require.NoError(t, err)
	})

	t.Run("slug conflict", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		existing := &model.VipLevel{Slug: "vip1"}
		existing.ID = 1
		other := &model.VipLevel{Slug: "vip2"}
		other.ID = 2
		level := &model.VipLevel{Slug: "vip2"}
		level.ID = 1

		repo.On("GetLevel", ctx, uint64(1)).Return(existing, nil)
		repo.On("GetLevelBySlug", ctx, "vip2").Return(other, nil)

		err := svc.UpdateLevel(ctx, level)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestService_DeleteLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("DeleteLevel", ctx, uint64(1)).Return(nil)

		err := svc.DeleteLevel(ctx, 1)
		require.NoError(t, err)
	})
}

func TestService_GetConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "test_key", ConfigValue: "test_value"}
		repo.On("GetConfig", ctx, "test_key").Return(config, nil)

		result, err := svc.GetConfig(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, "test_value", result.ConfigValue)
	})
}

func TestService_GetConfigInt64(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "threshold", ConfigValue: "10000"}
		repo.On("GetConfig", ctx, "threshold").Return(config, nil)

		result, err := svc.GetConfigInt64(ctx, "threshold")
		require.NoError(t, err)
		assert.Equal(t, int64(10000), result)
	})

	t.Run("invalid value", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "threshold", ConfigValue: "invalid"}
		repo.On("GetConfig", ctx, "threshold").Return(config, nil)

		_, err := svc.GetConfigInt64(ctx, "threshold")
		require.Error(t, err)
	})
}

func TestService_SetConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("create new config", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetConfig", ctx, "new_key").Return(nil, repository.ErrNotFound)
		repo.On("SaveConfig", ctx, mock.AnythingOfType("*model.VipConfig")).Return(nil)

		err := svc.SetConfig(ctx, "new_key", "new_value", "description")
		require.NoError(t, err)
	})

	t.Run("update existing config", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		existing := &model.VipConfig{ConfigKey: "existing_key", ConfigValue: "old_value"}
		repo.On("GetConfig", ctx, "existing_key").Return(existing, nil)
		repo.On("SaveConfig", ctx, mock.AnythingOfType("*model.VipConfig")).Return(nil)

		err := svc.SetConfig(ctx, "existing_key", "new_value", "")
		require.NoError(t, err)
	})
}

func TestService_CheckVipUnlock(t *testing.T) {
	ctx := context.Background()

	t.Run("unlock by consume", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		consumeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByConsume, ConfigValue: "10000"}
		rechargeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByRecharge, ConfigValue: "5000"}
		repo.On("GetConfig", ctx, model.VipConfigUnlockByConsume).Return(consumeConfig, nil)
		repo.On("GetConfig", ctx, model.VipConfigUnlockByRecharge).Return(rechargeConfig, nil)

		unlocked, err := svc.CheckVipUnlock(ctx, 15000, 0)
		require.NoError(t, err)
		assert.True(t, unlocked)
	})

	t.Run("unlock by recharge", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		consumeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByConsume, ConfigValue: "10000"}
		rechargeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByRecharge, ConfigValue: "5000"}
		repo.On("GetConfig", ctx, model.VipConfigUnlockByConsume).Return(consumeConfig, nil)
		repo.On("GetConfig", ctx, model.VipConfigUnlockByRecharge).Return(rechargeConfig, nil)

		unlocked, err := svc.CheckVipUnlock(ctx, 0, 6000)
		require.NoError(t, err)
		assert.True(t, unlocked)
	})

	t.Run("not unlocked", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		consumeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByConsume, ConfigValue: "10000"}
		rechargeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByRecharge, ConfigValue: "5000"}
		repo.On("GetConfig", ctx, model.VipConfigUnlockByConsume).Return(consumeConfig, nil)
		repo.On("GetConfig", ctx, model.VipConfigUnlockByRecharge).Return(rechargeConfig, nil)

		unlocked, err := svc.CheckVipUnlock(ctx, 5000, 2000)
		require.NoError(t, err)
		assert.False(t, unlocked)
	})
}

func TestService_BatchUpdateLevelStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("BatchUpdateLevelStatus", ctx, []uint64{1, 2, 3}, true).Return(int64(3), nil)

		affected, err := svc.BatchUpdateLevelStatus(ctx, []uint64{1, 2, 3}, true)
		require.NoError(t, err)
		assert.Equal(t, int64(3), affected)
	})
}

func TestService_BatchDeleteLevels(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("BatchDeleteLevels", ctx, []uint64{1, 2}).Return(int64(2), nil)

		affected, err := svc.BatchDeleteLevels(ctx, []uint64{1, 2})
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)
	})
}

func TestService_GetLevelBySlug(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		expected := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
		repo.On("GetLevelBySlug", ctx, "vip1").Return(expected, nil)

		result, err := svc.GetLevelBySlug(ctx, "vip1")
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetLevelBySlug", ctx, "nonexistent").Return(nil, repository.ErrNotFound)

		_, err := svc.GetLevelBySlug(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestService_GetDefaultLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		expected := &model.VipLevel{Slug: "vip1", Title: "VIP 1", IsDefault: true}
		repo.On("GetDefaultLevel", ctx).Return(expected, nil)

		result, err := svc.GetDefaultLevel(ctx)
		require.NoError(t, err)
		assert.True(t, result.IsDefault)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetDefaultLevel", ctx).Return(nil, repository.ErrNotFound)

		_, err := svc.GetDefaultLevel(ctx)
		require.Error(t, err)
	})
}

func TestService_ListActiveLevels(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		levels := []model.VipLevel{
			{Slug: "vip1", Title: "VIP 1", IsActive: true},
		}
		repo.On("ListActiveLevels", ctx).Return(levels, nil)

		result, err := svc.ListActiveLevels(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("ListActiveLevels", ctx).Return([]model.VipLevel{}, errors.New("db error"))

		_, err := svc.ListActiveLevels(ctx)
		require.Error(t, err)
	})
}

func TestService_ListLevelsPaged(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		opts := repository.VipLevelListOptions{Page: 1, PageSize: 10}
		levels := []model.VipLevel{{Slug: "vip1"}}
		repo.On("ListLevelsPaged", ctx, opts).Return(levels, int64(1), nil)

		result, total, err := svc.ListLevelsPaged(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		opts := repository.VipLevelListOptions{Page: 1, PageSize: 10}
		repo.On("ListLevelsPaged", ctx, opts).Return([]model.VipLevel{}, int64(0), errors.New("db error"))

		_, _, err := svc.ListLevelsPaged(ctx, opts)
		require.Error(t, err)
	})
}

func TestService_SetDefaultLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("SetDefaultLevel", ctx, uint64(1)).Return(nil)

		err := svc.SetDefaultLevel(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("SetDefaultLevel", ctx, uint64(999)).Return(errors.New("not found"))

		err := svc.SetDefaultLevel(ctx, 999)
		require.Error(t, err)
	})
}

func TestService_GetLevelByExp(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		expected := &model.VipLevel{Slug: "vip2", ExpRequired: 10000}
		repo.On("GetLevelByExp", ctx, int64(15000)).Return(expected, nil)

		result, err := svc.GetLevelByExp(ctx, 15000)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetLevelByExp", ctx, int64(0)).Return(nil, repository.ErrNotFound)

		_, err := svc.GetLevelByExp(ctx, 0)
		require.Error(t, err)
	})
}

func TestService_GetConfigValue(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "test_key", ConfigValue: "test_value"}
		repo.On("GetConfig", ctx, "test_key").Return(config, nil)

		result, err := svc.GetConfigValue(ctx, "test_key")
		require.NoError(t, err)
		assert.Equal(t, "test_value", result)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetConfig", ctx, "nonexistent").Return(nil, repository.ErrNotFound)

		_, err := svc.GetConfigValue(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestService_ListConfigs(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		configs := []model.VipConfig{
			{ConfigKey: "key1", ConfigValue: "value1"},
			{ConfigKey: "key2", ConfigValue: "value2"},
		}
		repo.On("ListConfigs", ctx).Return(configs, nil)

		result, err := svc.ListConfigs(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("ListConfigs", ctx).Return([]model.VipConfig{}, errors.New("db error"))

		_, err := svc.ListConfigs(ctx)
		require.Error(t, err)
	})
}

func TestService_SaveConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "key", ConfigValue: "value"}
		repo.On("SaveConfig", ctx, config).Return(nil)

		err := svc.SaveConfig(ctx, config)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: "key", ConfigValue: "value"}
		repo.On("SaveConfig", ctx, config).Return(errors.New("db error"))

		err := svc.SaveConfig(ctx, config)
		require.Error(t, err)
	})
}

func TestService_DeleteConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("DeleteConfig", ctx, "key").Return(nil)

		err := svc.DeleteConfig(ctx, "key")
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("DeleteConfig", ctx, "nonexistent").Return(errors.New("not found"))

		err := svc.DeleteConfig(ctx, "nonexistent")
		require.Error(t, err)
	})
}

func TestService_GetUnlockThreshold(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		consumeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByConsume, ConfigValue: "10000"}
		rechargeConfig := &model.VipConfig{ConfigKey: model.VipConfigUnlockByRecharge, ConfigValue: "5000"}
		repo.On("GetConfig", ctx, model.VipConfigUnlockByConsume).Return(consumeConfig, nil)
		repo.On("GetConfig", ctx, model.VipConfigUnlockByRecharge).Return(rechargeConfig, nil)

		consume, recharge, err := svc.GetUnlockThreshold(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), consume)
		assert.Equal(t, int64(5000), recharge)
	})

	t.Run("config not found returns zero", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetConfig", ctx, model.VipConfigUnlockByConsume).Return(nil, repository.ErrNotFound)
		repo.On("GetConfig", ctx, model.VipConfigUnlockByRecharge).Return(nil, repository.ErrNotFound)

		consume, recharge, err := svc.GetUnlockThreshold(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), consume)
		assert.Equal(t, int64(0), recharge)
	})
}

func TestService_CalculateVipLevel(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		expected := &model.VipLevel{Slug: "vip2", ExpRequired: 10000}
		repo.On("GetLevelByExp", ctx, int64(15000)).Return(expected, nil)

		result, err := svc.CalculateVipLevel(ctx, 15000)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestService_GetVipExpireDays(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: model.VipConfigExpireDays, ConfigValue: "365"}
		repo.On("GetConfig", ctx, model.VipConfigExpireDays).Return(config, nil)

		days, err := svc.GetVipExpireDays(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(365), days)
	})

	t.Run("not found returns zero (permanent)", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		repo.On("GetConfig", ctx, model.VipConfigExpireDays).Return(nil, repository.ErrNotFound)

		days, err := svc.GetVipExpireDays(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), days)
	})

	t.Run("invalid value error", func(t *testing.T) {
		repo := &MockVipRepository{}
		svc := NewVipService(repo)

		config := &model.VipConfig{ConfigKey: model.VipConfigExpireDays, ConfigValue: "invalid"}
		repo.On("GetConfig", ctx, model.VipConfigExpireDays).Return(config, nil)

		_, err := svc.GetVipExpireDays(ctx)
		require.Error(t, err)
	})
}

func TestService_UpdateLevel_GetLevelError(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	level := &model.VipLevel{Slug: "vip1"}
	level.ID = 1
	repo.On("GetLevel", ctx, uint64(1)).Return(nil, errors.New("db error"))

	err := svc.UpdateLevel(ctx, level)
	require.Error(t, err)
}

func TestService_UpdateLevel_CheckSlugError(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	existing := &model.VipLevel{Slug: "vip1"}
	existing.ID = 1
	level := &model.VipLevel{Slug: "vip2"}
	level.ID = 1

	repo.On("GetLevel", ctx, uint64(1)).Return(existing, nil)
	repo.On("GetLevelBySlug", ctx, "vip2").Return(nil, errors.New("db error"))

	err := svc.UpdateLevel(ctx, level)
	require.Error(t, err)
}

func TestService_CreateLevel_RepoError(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	level := &model.VipLevel{Slug: "vip1", Title: "VIP 1"}
	repo.On("GetLevelBySlug", ctx, "vip1").Return(nil, repository.ErrNotFound)
	repo.On("CreateLevel", ctx, level).Return(errors.New("db error"))

	err := svc.CreateLevel(ctx, level)
	require.Error(t, err)
}

func TestService_ListLevels_Error(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	repo.On("ListLevels", ctx).Return([]model.VipLevel{}, errors.New("db error"))

	_, err := svc.ListLevels(ctx)
	require.Error(t, err)
}

func TestService_DeleteLevel_Error(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	repo.On("DeleteLevel", ctx, uint64(1)).Return(errors.New("db error"))

	err := svc.DeleteLevel(ctx, 1)
	require.Error(t, err)
}

func TestService_UpdateLevel_UpdateError(t *testing.T) {
	ctx := context.Background()

	repo := &MockVipRepository{}
	svc := NewVipService(repo)

	existing := &model.VipLevel{Slug: "vip1"}
	existing.ID = 1
	level := &model.VipLevel{Slug: "vip1", Title: "Updated"}
	level.ID = 1

	repo.On("GetLevel", ctx, uint64(1)).Return(existing, nil)
	repo.On("UpdateLevel", ctx, level).Return(errors.New("db error"))

	err := svc.UpdateLevel(ctx, level)
	require.Error(t, err)
}

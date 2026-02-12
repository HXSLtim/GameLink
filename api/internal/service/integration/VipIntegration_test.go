package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/vip"
	vipservice "gamelink/internal/service/vip"
)

// setupVipService creates a VIP service for testing.
func setupVipService(t *testing.T) *vipservice.Service {
	t.Helper()
	db := SetupTestDB(t)
	repo := vip.NewVipRepository(db)
	return vipservice.NewVipService(repo)
}

// ============================================================================
// VIP Level Tests
// ============================================================================

func TestVipService_CreateLevel(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level := &model.VipLevel{
		Slug:          "vip1",
		Title:         "VIP 1",
		ExpRequired:   10000,
		OrderDiscount: 0.98,
		IsActive:      true,
		Benefits:      "{}",
	}

	err := svc.CreateLevel(ctx, level)
	require.NoError(t, err)
	assert.NotZero(t, level.ID)

	// Verify
	got, err := svc.GetLevel(ctx, level.ID)
	require.NoError(t, err)
	assert.Equal(t, "vip1", got.Slug)
	assert.Equal(t, "VIP 1", got.Title)
	assert.Equal(t, int64(10000), got.ExpRequired)
	assert.Equal(t, 0.98, got.OrderDiscount)
}

func TestVipService_CreateLevel_DuplicateSlug(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level1 := &model.VipLevel{
		Slug:     "vip_dup",
		Title:    "VIP Dup",
		IsActive: true,
		Benefits: "{}",
	}
	err := svc.CreateLevel(ctx, level1)
	require.NoError(t, err)

	level2 := &model.VipLevel{
		Slug:     "vip_dup",
		Title:    "VIP Dup 2",
		IsActive: true,
		Benefits: "{}",
	}
	err = svc.CreateLevel(ctx, level2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestVipService_GetLevelBySlug(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level := &model.VipLevel{
		Slug:     "vip_slug_test",
		Title:    "VIP Slug Test",
		IsActive: true,
		Benefits: "{}",
	}
	err := svc.CreateLevel(ctx, level)
	require.NoError(t, err)

	got, err := svc.GetLevelBySlug(ctx, "vip_slug_test")
	require.NoError(t, err)
	assert.Equal(t, level.ID, got.ID)
	assert.Equal(t, "VIP Slug Test", got.Title)
}

func TestVipService_UpdateLevel(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level := &model.VipLevel{
		Slug:          "vip_update",
		Title:         "VIP Update",
		ExpRequired:   5000,
		OrderDiscount: 0.95,
		IsActive:      true,
		Benefits:      "{}",
	}
	err := svc.CreateLevel(ctx, level)
	require.NoError(t, err)

	// Update
	level.Title = "VIP Updated"
	level.ExpRequired = 8000
	level.OrderDiscount = 0.90
	err = svc.UpdateLevel(ctx, level)
	require.NoError(t, err)

	// Verify
	got, err := svc.GetLevel(ctx, level.ID)
	require.NoError(t, err)
	assert.Equal(t, "VIP Updated", got.Title)
	assert.Equal(t, int64(8000), got.ExpRequired)
	assert.Equal(t, 0.90, got.OrderDiscount)
}

func TestVipService_DeleteLevel(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level := &model.VipLevel{
		Slug:     "vip_delete",
		Title:    "VIP Delete",
		IsActive: true,
		Benefits: "{}",
	}
	err := svc.CreateLevel(ctx, level)
	require.NoError(t, err)

	err = svc.DeleteLevel(ctx, level.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetLevel(ctx, level.ID)
	assert.Error(t, err)
}

func TestVipService_ListActiveLevels(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	// Create active and inactive levels
	active := &model.VipLevel{
		Slug:     "vip_active",
		Title:    "VIP Active",
		IsActive: true,
		Benefits: "{}",
	}
	inactive := &model.VipLevel{
		Slug:     "vip_inactive",
		Title:    "VIP Inactive",
		IsActive: false,
		Benefits: "{}",
	}
	require.NoError(t, svc.CreateLevel(ctx, active))
	require.NoError(t, svc.CreateLevel(ctx, inactive))

	levels, err := svc.ListActiveLevels(ctx)
	require.NoError(t, err)

	// Should only contain active levels
	for _, l := range levels {
		assert.True(t, l.IsActive)
	}
}

func TestVipService_SetDefaultLevel(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	level1 := &model.VipLevel{
		Slug:      "vip_default1",
		Title:     "VIP Default 1",
		IsActive:  true,
		IsDefault: true,
		Benefits:  "{}",
	}
	level2 := &model.VipLevel{
		Slug:      "vip_default2",
		Title:     "VIP Default 2",
		IsActive:  true,
		IsDefault: false,
		Benefits:  "{}",
	}
	require.NoError(t, svc.CreateLevel(ctx, level1))
	require.NoError(t, svc.CreateLevel(ctx, level2))

	// Set level2 as default
	err := svc.SetDefaultLevel(ctx, level2.ID)
	require.NoError(t, err)

	// Verify level2 is now default
	got, err := svc.GetDefaultLevel(ctx)
	require.NoError(t, err)
	assert.Equal(t, level2.ID, got.ID)
}

func TestVipService_GetLevelByExp(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	// Create levels with different exp requirements
	level1 := &model.VipLevel{
		Slug:        "vip_exp1",
		Title:       "VIP Bronze",
		ExpRequired: 0,
		IsActive:    true,
		Benefits:    "{}",
	}
	level2 := &model.VipLevel{
		Slug:        "vip_exp2",
		Title:       "VIP Silver",
		ExpRequired: 10000,
		IsActive:    true,
		Benefits:    "{}",
	}
	level3 := &model.VipLevel{
		Slug:        "vip_exp3",
		Title:       "VIP Gold",
		ExpRequired: 50000,
		IsActive:    true,
		Benefits:    "{}",
	}
	require.NoError(t, svc.CreateLevel(ctx, level1))
	require.NoError(t, svc.CreateLevel(ctx, level2))
	require.NoError(t, svc.CreateLevel(ctx, level3))

	// Test exp = 25000 should get Silver (10000 <= 25000 < 50000)
	got, err := svc.GetLevelByExp(ctx, 25000)
	require.NoError(t, err)
	assert.Equal(t, "VIP Silver", got.Title)

	// Test exp = 100000 should get Gold
	got, err = svc.GetLevelByExp(ctx, 100000)
	require.NoError(t, err)
	assert.Equal(t, "VIP Gold", got.Title)
}

// ============================================================================
// VIP Config Tests
// ============================================================================

func TestVipService_SetConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	err := svc.SetConfig(ctx, "test_key", "test_value", "Test description")
	require.NoError(t, err)

	// Verify
	value, err := svc.GetConfigValue(ctx, "test_key")
	require.NoError(t, err)
	assert.Equal(t, "test_value", value)
}

func TestVipService_GetConfigInt64(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	err := svc.SetConfig(ctx, "int_config", "12345", "Integer config")
	require.NoError(t, err)

	value, err := svc.GetConfigInt64(ctx, "int_config")
	require.NoError(t, err)
	assert.Equal(t, int64(12345), value)
}

func TestVipService_CheckVipUnlock(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	// Set unlock thresholds
	err := svc.SetConfig(ctx, model.VipConfigUnlockByConsume, "10000", "Consume threshold")
	require.NoError(t, err)
	err = svc.SetConfig(ctx, model.VipConfigUnlockByRecharge, "5000", "Recharge threshold")
	require.NoError(t, err)

	// Test: below both thresholds
	unlocked, err := svc.CheckVipUnlock(ctx, 3000, 2000)
	require.NoError(t, err)
	assert.False(t, unlocked)

	// Test: meets consume threshold
	unlocked, err = svc.CheckVipUnlock(ctx, 15000, 2000)
	require.NoError(t, err)
	assert.True(t, unlocked)

	// Test: meets recharge threshold
	unlocked, err = svc.CheckVipUnlock(ctx, 3000, 6000)
	require.NoError(t, err)
	assert.True(t, unlocked)
}

func TestVipService_ListConfigs(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	// Create multiple configs
	require.NoError(t, svc.SetConfig(ctx, "config_a", "value_a", "Config A"))
	require.NoError(t, svc.SetConfig(ctx, "config_b", "value_b", "Config B"))

	configs, err := svc.ListConfigs(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(configs), 2)
}

func TestVipService_DeleteConfig(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupVipService(t)
	ctx := context.Background()

	err := svc.SetConfig(ctx, "delete_config", "value", "To be deleted")
	require.NoError(t, err)

	err = svc.DeleteConfig(ctx, "delete_config")
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetConfigValue(ctx, "delete_config")
	assert.Error(t, err)
}

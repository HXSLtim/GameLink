package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/recharge"
	rechargeservice "gamelink/internal/service/recharge"
)

// setupRechargeService creates a Recharge service for testing.
func setupRechargeService(t *testing.T) *rechargeservice.Service {
	t.Helper()
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	// Create service without wallet and coupon dependencies for basic tests
	return rechargeservice.NewRechargeService(repo, nil, nil)
}

// ============================================================================
// Recharge Option Tests
// ============================================================================

func TestRechargeService_CreateOption(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupRechargeService(t)
	ctx := context.Background()

	option := &model.RechargeOption{
		Name:        "Basic Recharge",
		AmountCents: 10000, // ¥100
		BonusCents:  1000,  // ¥10 bonus
		Description: "Basic recharge option",
		SortOrder:   1,
		IsActive:    true,
	}

	err := svc.CreateOption(ctx, option)
	require.NoError(t, err)
	assert.NotZero(t, option.ID)
	assert.Equal(t, int64(11000), option.TotalCents) // Auto-calculated

	// Verify
	got, err := svc.GetOption(ctx, option.ID)
	require.NoError(t, err)
	assert.Equal(t, "Basic Recharge", got.Name)
	assert.Equal(t, int64(10000), got.AmountCents)
	assert.Equal(t, int64(1000), got.BonusCents)
	assert.Equal(t, int64(11000), got.TotalCents)
}

func TestRechargeService_CreateOption_InvalidAmount(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupRechargeService(t)
	ctx := context.Background()

	option := &model.RechargeOption{
		Name:        "Invalid Option",
		AmountCents: 0, // Invalid
		IsActive:    true,
	}

	err := svc.CreateOption(ctx, option)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "金额必须大于0")
}

func TestRechargeService_UpdateOption(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupRechargeService(t)
	ctx := context.Background()

	option := &model.RechargeOption{
		Name:        "Update Test",
		AmountCents: 5000,
		BonusCents:  500,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// Update
	option.Name = "Updated Name"
	option.AmountCents = 6000
	option.BonusCents = 600
	err := svc.UpdateOption(ctx, option)
	require.NoError(t, err)

	// Verify
	got, err := svc.GetOption(ctx, option.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, int64(6000), got.AmountCents)
	assert.Equal(t, int64(6600), got.TotalCents)
}

func TestRechargeService_DeleteOption(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupRechargeService(t)
	ctx := context.Background()

	option := &model.RechargeOption{
		Name:        "Delete Test",
		AmountCents: 1000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	err := svc.DeleteOption(ctx, option.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetOption(ctx, option.ID)
	assert.Error(t, err)
}

func TestRechargeService_GetActiveOptions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	// Create active options
	for i := 1; i <= 3; i++ {
		option := &model.RechargeOption{
			Name:        "Active Option",
			AmountCents: int64(i * 1000),
			IsActive:    true,
			SortOrder:   i,
		}
		require.NoError(t, svc.CreateOption(ctx, option))
	}

	// Create inactive option
	inactiveOption := &model.RechargeOption{
		Name:        "Inactive Option",
		AmountCents: 5000,
		IsActive:    true, // Create as active first
	}
	require.NoError(t, svc.CreateOption(ctx, inactiveOption))
	// Update to inactive
	require.NoError(t, db.Model(inactiveOption).Update("is_active", false).Error)

	// Get active options
	options, err := svc.GetActiveOptions(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, options, 3)
}

func TestRechargeService_BatchUpdateOptionStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupRechargeService(t)
	ctx := context.Background()

	var ids []uint64
	for i := 0; i < 3; i++ {
		option := &model.RechargeOption{
			Name:        "Batch Test",
			AmountCents: int64((i + 1) * 1000),
			IsActive:    true,
		}
		require.NoError(t, svc.CreateOption(ctx, option))
		ids = append(ids, option.ID)
	}

	// Batch disable
	affected, err := svc.BatchUpdateOptionStatus(ctx, ids, false)
	require.NoError(t, err)
	assert.Equal(t, int64(3), affected)

	// Verify all disabled
	options, err := svc.GetActiveOptions(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, options, 0)
}

// ============================================================================
// Recharge Order Tests
// ============================================================================

func TestRechargeService_CreateRechargeOrder(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_order")

	option := &model.RechargeOption{
		Name:        "Order Test",
		AmountCents: 10000,
		BonusCents:  1000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// Create recharge order
	record, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	require.NoError(t, err)
	assert.NotZero(t, record.ID)
	assert.Equal(t, user.ID, record.UserID)
	assert.Equal(t, model.RechargeStatusPending, record.Status)
	assert.Equal(t, int64(10000), record.AmountCents)
	assert.Equal(t, int64(1000), record.BonusCents)
	assert.Equal(t, int64(11000), record.TotalCents)
	assert.NotEmpty(t, record.OrderNo)
}

func TestRechargeService_CreateRechargeOrder_InactiveOption(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_inactive")

	option := &model.RechargeOption{
		Name:        "Inactive Test",
		AmountCents: 10000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))
	// Update to inactive
	require.NoError(t, db.Model(option).Update("is_active", false).Error)

	_, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已下架")
}

func TestRechargeService_CreateRechargeOrder_SoldOut(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_soldout")

	option := &model.RechargeOption{
		Name:          "Sold Out Test",
		AmountCents:   10000,
		IsActive:      true,
		TotalLimit:    1,
		PurchaseCount: 1, // Already sold out
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	_, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已售罄")
}

func TestRechargeService_CreateRechargeOrder_PerUserLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_limit")

	option := &model.RechargeOption{
		Name:         "Limit Test",
		AmountCents:  10000,
		IsActive:     true,
		PerUserLimit: 1,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// First order - should succeed
	record1, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	require.NoError(t, err)

	// Simulate payment success
	require.NoError(t, repo.MarkAsPaid(ctx, record1.ID, "TRADE123"))

	// Second order - should fail
	_, err = svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "购买上限")
}

func TestRechargeService_GetUserRecords(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_records")

	option := &model.RechargeOption{
		Name:        "Records Test",
		AmountCents: 10000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// Create multiple orders
	for i := 0; i < 3; i++ {
		_, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
		require.NoError(t, err)
	}

	records, err := svc.GetUserRecords(ctx, user.ID, 10)
	require.NoError(t, err)
	assert.Len(t, records, 3)
}

func TestRechargeService_CancelExpiredRecords(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_expire")

	option := &model.RechargeOption{
		Name:        "Expire Test",
		AmountCents: 10000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// Create order
	record, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	require.NoError(t, err)

	// Manually set expire time to past
	require.NoError(t, db.Exec("UPDATE recharge_records SET expire_at = NOW() - INTERVAL '1 hour' WHERE id = ?", record.ID).Error)

	// Run cancel job
	affected, err := svc.CancelExpiredRecords(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, affected, int64(1))

	// Verify canceled
	canceled, err := svc.GetRecord(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, model.RechargeStatusCanceled, canceled.Status)
}

func TestRechargeService_GetRechargeStats(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := recharge.NewRechargeRepository(db)
	svc := rechargeservice.NewRechargeService(repo, nil, nil)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "recharge_stats")

	option := &model.RechargeOption{
		Name:        "Stats Test",
		AmountCents: 10000,
		BonusCents:  1000,
		IsActive:    true,
	}
	require.NoError(t, svc.CreateOption(ctx, option))

	// Create and pay order
	record, err := svc.CreateRechargeOrder(ctx, user.ID, option.ID, "wechat", "wechat_h5", "127.0.0.1", "TestAgent")
	require.NoError(t, err)
	require.NoError(t, repo.MarkAsPaid(ctx, record.ID, "TRADE_STATS"))

	// Get stats
	stats, err := svc.GetRechargeStats(ctx)
	require.NoError(t, err)
	assert.NotNil(t, stats["totalAmountCents"])
	assert.NotNil(t, stats["totalCount"])
}

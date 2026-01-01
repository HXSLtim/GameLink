package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository/coupon"
	couponservice "gamelink/internal/service/coupon"
)

// setupCouponService creates a Coupon service for testing.
func setupCouponService(t *testing.T) *couponservice.Service {
	t.Helper()
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	return couponservice.NewCouponService(repo)
}

// ============================================================================
// Coupon Template Tests
// ============================================================================

func TestCouponService_CreateTemplate(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupCouponService(t)
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:              "Test Coupon",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 1000,
		MinAmountCents:    5000,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		TotalCount:        100,
		PerUserLimit:      1,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}

	err := svc.CreateTemplate(ctx, template)
	require.NoError(t, err)
	assert.NotZero(t, template.ID)

	// Verify
	got, err := svc.GetTemplate(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, "Test Coupon", got.Name)
	assert.Equal(t, model.CouponTypeDeduct, got.Type)
	assert.Equal(t, int64(1000), got.DeductAmountCents)
}

func TestCouponService_CreateTemplate_Discount(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupCouponService(t)
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:             "Discount Coupon",
		Type:             model.CouponTypeDiscount,
		Source:           model.CouponSourceManual,
		DiscountRate:     0.9,
		MaxDiscountCents: 5000,
		MinAmountCents:   10000,
		Scope:            model.CouponScopeAll,
		ValidityType:     "days",
		ValidityDays:     15,
		IsActive:         true,
		GameIDs:          "[]",
		ItemIDs:          "[]",
	}

	err := svc.CreateTemplate(ctx, template)
	require.NoError(t, err)

	got, err := svc.GetTemplate(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CouponTypeDiscount, got.Type)
	assert.Equal(t, 0.9, got.DiscountRate)
}

func TestCouponService_CreateTemplate_InvalidType(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupCouponService(t)
	ctx := context.Background()

	// Deduct type without amount
	template := &model.CouponTemplate{
		Name:              "Invalid Coupon",
		Type:              model.CouponTypeDeduct,
		DeductAmountCents: 0, // Invalid
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}

	err := svc.CreateTemplate(ctx, template)
	assert.Error(t, err)
}

func TestCouponService_UpdateTemplate(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupCouponService(t)
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:              "Update Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      7,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	// Update
	template.Name = "Updated Name"
	template.DeductAmountCents = 800
	err := svc.UpdateTemplate(ctx, template)
	require.NoError(t, err)

	// Verify
	got, err := svc.GetTemplate(ctx, template.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, int64(800), got.DeductAmountCents)
}

func TestCouponService_DeleteTemplate(t *testing.T) {
	SkipIfNoTestDB(t)
	svc := setupCouponService(t)
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:              "Delete Test",
		Type:              model.CouponTypeDeduct,
		DeductAmountCents: 100,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      7,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	err := svc.DeleteTemplate(ctx, template.ID)
	require.NoError(t, err)

	// Verify deleted
	_, err = svc.GetTemplate(ctx, template.ID)
	assert.Error(t, err)
}

// ============================================================================
// Coupon Claim Tests
// ============================================================================

func TestCouponService_ClaimCoupon(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	// Create user
	user := CreateUniqueTestUser(t, db, "coupon_claim")

	// Create template
	template := &model.CouponTemplate{
		Name:              "Claim Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 1000,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		TotalCount:        100,
		PerUserLimit:      2,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	// Claim coupon
	claimed, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)
	assert.NotZero(t, claimed.ID)
	assert.Equal(t, user.ID, claimed.UserID)
	assert.Equal(t, model.CouponStateAvailable, claimed.State)
}

func TestCouponService_ClaimCoupon_ExceedLimit(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_limit")

	template := &model.CouponTemplate{
		Name:              "Limit Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		TotalCount:        100,
		PerUserLimit:      1, // Only 1 per user
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	// First claim - should succeed
	_, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)

	// Second claim - should fail
	_, err = svc.ClaimCoupon(ctx, user.ID, template.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "领取上限")
}

func TestCouponService_ClaimCoupon_Inactive(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_inactive")

	// Create active template first, then update to inactive
	template := &model.CouponTemplate{
		Name:              "Inactive Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, db.Create(template).Error)

	// Update to inactive
	require.NoError(t, db.Model(template).Update("is_active", false).Error)

	_, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "已下架")
}

// ============================================================================
// Coupon Usage Tests
// ============================================================================

func TestCouponService_LockAndUseCoupon(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_use")

	template := &model.CouponTemplate{
		Name:              "Use Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 1000,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	// Claim
	claimed, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)

	// Lock for order
	orderID := uint64(12345)
	err = svc.LockCoupon(ctx, claimed.ID, orderID)
	require.NoError(t, err)

	// Verify locked
	locked, err := svc.GetCoupon(ctx, claimed.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CouponStateLocked, locked.State)

	// Use coupon
	err = svc.UseCoupon(ctx, claimed.ID, orderID, 1000)
	require.NoError(t, err)

	// Verify used
	used, err := svc.GetCoupon(ctx, claimed.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CouponStateUsed, used.State)
	assert.Equal(t, int64(1000), used.DiscountCents)
}

func TestCouponService_UnlockCoupon(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_unlock")

	template := &model.CouponTemplate{
		Name:              "Unlock Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 500,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	claimed, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)

	// Lock
	err = svc.LockCoupon(ctx, claimed.ID, 99999)
	require.NoError(t, err)

	// Unlock (cancel order)
	err = svc.UnlockCoupon(ctx, claimed.ID)
	require.NoError(t, err)

	// Verify available again
	unlocked, err := svc.GetCoupon(ctx, claimed.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CouponStateAvailable, unlocked.State)
}

func TestCouponService_CalculateDiscount_Deduct(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_calc_deduct")

	template := &model.CouponTemplate{
		Name:              "Deduct Calc",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 1000, // ¥10
		MinAmountCents:    5000, // Min ¥50
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      30,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	claimed, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)

	// Order ¥100, should get ¥10 discount
	discount, err := svc.CalculateDiscount(ctx, claimed.ID, 10000)
	require.NoError(t, err)
	assert.Equal(t, int64(1000), discount)
}

func TestCouponService_CalculateDiscount_Discount(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_calc_discount")

	template := &model.CouponTemplate{
		Name:             "Discount Calc",
		Type:             model.CouponTypeDiscount,
		Source:           model.CouponSourceManual,
		DiscountRate:     0.9,   // 90% = 10% off
		MaxDiscountCents: 2000,  // Max ¥20
		MinAmountCents:   10000, // Min ¥100
		Scope:            model.CouponScopeAll,
		ValidityType:     "days",
		ValidityDays:     30,
		IsActive:         true,
		GameIDs:          "[]",
		ItemIDs:          "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	claimed, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
	require.NoError(t, err)

	// Order ¥150, 10% off = ¥15 discount (allow for floating point variance)
	discount, err := svc.CalculateDiscount(ctx, claimed.ID, 15000)
	require.NoError(t, err)
	assert.InDelta(t, int64(1500), discount, 1) // Allow 1 cent variance

	// Order ¥300, 10% off = ¥30, but max is ¥20
	discount, err = svc.CalculateDiscount(ctx, claimed.ID, 30000)
	require.NoError(t, err)
	assert.Equal(t, int64(2000), discount) // Capped at max
}

func TestCouponService_GetUserAvailableCoupons(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_available")

	// Create and claim multiple coupons
	for i := 0; i < 3; i++ {
		template := &model.CouponTemplate{
			Name:              fmt.Sprintf("Available Test %d", i),
			Type:              model.CouponTypeDeduct,
			Source:            model.CouponSourceManual,
			DeductAmountCents: int64(100 * (i + 1)),
			Scope:             model.CouponScopeAll,
			ValidityType:      "days",
			ValidityDays:      30,
			IsActive:          true,
			GameIDs:           "[]",
			ItemIDs:           "[]",
			ClaimLink:         fmt.Sprintf("available_test_%d_%d", user.ID, i), // Unique claim link
		}
		require.NoError(t, svc.CreateTemplate(ctx, template))
		_, err := svc.ClaimCoupon(ctx, user.ID, template.ID)
		require.NoError(t, err)
	}

	coupons, err := svc.GetUserAvailableCoupons(ctx, user.ID)
	require.NoError(t, err)
	assert.Len(t, coupons, 3)
}

func TestCouponService_IssueCoupon(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_issue")

	template := &model.CouponTemplate{
		Name:              "Issue Test",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 2000,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      7,
		PerUserLimit:      1, // Limit 1
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, svc.CreateTemplate(ctx, template))

	// Issue bypasses per-user limit
	issued, err := svc.IssueCoupon(ctx, user.ID, template.ID, model.CouponSourceVip)
	require.NoError(t, err)
	assert.Equal(t, model.CouponSourceVip, issued.Source)

	// Can issue again (no limit check)
	issued2, err := svc.IssueCoupon(ctx, user.ID, template.ID, model.CouponSourceActivity)
	require.NoError(t, err)
	assert.NotEqual(t, issued.ID, issued2.ID)
}

func TestCouponService_ExpireOldCoupons(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	repo := coupon.NewCouponRepository(db)
	svc := couponservice.NewCouponService(repo)
	ctx := context.Background()

	user := CreateUniqueTestUser(t, db, "coupon_expire")

	// Create template first (required for foreign key)
	template := &model.CouponTemplate{
		Name:              "Expire Test Template",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 100,
		Scope:             model.CouponScopeAll,
		ValidityType:      "days",
		ValidityDays:      1,
		IsActive:          true,
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, db.Create(template).Error)

	// Create expired coupon directly with valid TemplateID
	expiredCoupon := &model.Coupon{
		TemplateID:        template.ID,
		UserID:            user.ID,
		State:             model.CouponStateAvailable,
		Name:              "Expired Coupon",
		Type:              model.CouponTypeDeduct,
		Source:            model.CouponSourceManual,
		DeductAmountCents: 100,
		Scope:             model.CouponScopeAll,
		ExpireAt:          time.Now().Add(-24 * time.Hour), // Expired yesterday
		GameIDs:           "[]",
		ItemIDs:           "[]",
	}
	require.NoError(t, db.Create(expiredCoupon).Error)

	// Run expire job
	affected, err := svc.ExpireOldCoupons(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, affected, int64(1))

	// Verify expired
	expired, err := svc.GetCoupon(ctx, expiredCoupon.ID)
	require.NoError(t, err)
	assert.Equal(t, model.CouponStateExpired, expired.State)
}

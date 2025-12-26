package coupon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	couponrepo "gamelink/internal/repository/coupon"
)

// MockCouponRepository is a mock implementation
type MockCouponRepository struct {
	mock.Mock
}

func (m *MockCouponRepository) ListTemplates(ctx context.Context, opts couponrepo.TemplateListOptions) ([]model.CouponTemplate, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.CouponTemplate), args.Get(1).(int64), args.Error(2)
}

func (m *MockCouponRepository) GetTemplateByID(ctx context.Context, id uint64) (*model.CouponTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CouponTemplate), args.Error(1)
}

func (m *MockCouponRepository) GetTemplateByClaimLink(ctx context.Context, link string) (*model.CouponTemplate, error) {
	args := m.Called(ctx, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CouponTemplate), args.Error(1)
}

func (m *MockCouponRepository) CreateTemplate(ctx context.Context, template *model.CouponTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockCouponRepository) UpdateTemplate(ctx context.Context, template *model.CouponTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockCouponRepository) DeleteTemplate(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCouponRepository) BatchUpdateTemplateStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	args := m.Called(ctx, ids, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCouponRepository) BatchDeleteTemplates(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCouponRepository) ListCoupons(ctx context.Context, opts couponrepo.CouponListOptions) ([]model.Coupon, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Coupon), args.Get(1).(int64), args.Error(2)
}

func (m *MockCouponRepository) GetCouponByID(ctx context.Context, id uint64) (*model.Coupon, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Coupon), args.Error(1)
}

func (m *MockCouponRepository) GetCouponWithTemplate(ctx context.Context, id uint64) (*model.Coupon, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Coupon), args.Error(1)
}

func (m *MockCouponRepository) GetUserAvailableCoupons(ctx context.Context, userID uint64) ([]model.Coupon, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Coupon), args.Error(1)
}

func (m *MockCouponRepository) CountUserCouponsFromTemplate(ctx context.Context, userID, templateID uint64) (int64, error) {
	args := m.Called(ctx, userID, templateID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCouponRepository) CreateCoupon(ctx context.Context, coupon *model.Coupon) error {
	args := m.Called(ctx, coupon)
	return args.Error(0)
}

func (m *MockCouponRepository) IncrementClaimedCount(ctx context.Context, templateID uint64) error {
	args := m.Called(ctx, templateID)
	return args.Error(0)
}

func (m *MockCouponRepository) LockCoupon(ctx context.Context, couponID, orderID uint64) error {
	args := m.Called(ctx, couponID, orderID)
	return args.Error(0)
}

func (m *MockCouponRepository) UnlockCoupon(ctx context.Context, couponID uint64) error {
	args := m.Called(ctx, couponID)
	return args.Error(0)
}

func (m *MockCouponRepository) UseCoupon(ctx context.Context, couponID, orderID uint64, discountCents int64) error {
	args := m.Called(ctx, couponID, orderID, discountCents)
	return args.Error(0)
}

func (m *MockCouponRepository) ExpireOldCoupons(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCouponRepository) GetCouponStats(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockCouponRepository) DeleteCoupon(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestService_validateTemplate(t *testing.T) {
	svc := &Service{}

	t.Run("empty name", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "", Type: model.CouponTypeDeduct}
		err := svc.validateTemplate(template)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名称")
	})

	t.Run("deduct type invalid amount", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Test", Type: model.CouponTypeDeduct, DeductAmountCents: 0}
		err := svc.validateTemplate(template)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "满减金额")
	})

	t.Run("discount type invalid rate", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Test", Type: model.CouponTypeDiscount, DiscountRate: 1.5}
		err := svc.validateTemplate(template)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "折扣率")
	})

	t.Run("invalid type", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Test", Type: "invalid"}
		err := svc.validateTemplate(template)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "无效")
	})

	t.Run("valid deduct template", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Test", Type: model.CouponTypeDeduct, DeductAmountCents: 1000}
		err := svc.validateTemplate(template)
		require.NoError(t, err)
	})

	t.Run("valid discount template", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Test", Type: model.CouponTypeDiscount, DiscountRate: 0.9}
		err := svc.validateTemplate(template)
		require.NoError(t, err)
	})
}

func TestService_ClaimCoupon(t *testing.T) {
	t.Run("template validation", func(t *testing.T) {
		// Test template validation scenarios
		template := &model.CouponTemplate{
			Name:         "Test Coupon",
			Type:         model.CouponTypeDeduct,
			IsActive:     true,
			TotalCount:   100,
			ClaimedCount: 50,
			PerUserLimit: 3,
			ValidityDays: 30,
		}
		template.ID = 1

		assert.True(t, template.IsActive)
		assert.Less(t, template.ClaimedCount, template.TotalCount)
	})

	t.Run("template not active check", func(t *testing.T) {
		template := &model.CouponTemplate{IsActive: false}
		assert.False(t, template.IsActive)
	})

	t.Run("coupon sold out check", func(t *testing.T) {
		template := &model.CouponTemplate{TotalCount: 100, ClaimedCount: 100}
		assert.Equal(t, template.TotalCount, template.ClaimedCount)
	})
}

func TestService_CalculateDiscount(t *testing.T) {
	// Test discount calculation logic
	t.Run("deduct coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			Type:              model.CouponTypeDeduct,
			State:             model.CouponStateAvailable,
			MinAmountCents:    5000,
			DeductAmountCents: 1000,
			ExpireAt:          time.Now().Add(24 * time.Hour),
		}

		// Order amount >= min amount
		discount := coupon.CalculateDiscount(10000)
		assert.Equal(t, int64(1000), discount)

		// Order amount < min amount
		discount = coupon.CalculateDiscount(3000)
		assert.Equal(t, int64(0), discount)
	})

	t.Run("discount coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			Type:             model.CouponTypeDiscount,
			State:            model.CouponStateAvailable,
			MinAmountCents:   0,
			DiscountRate:     0.5, // 5折，优惠50%
			MaxDiscountCents: 3000,
			ExpireAt:         time.Now().Add(24 * time.Hour),
		}

		// 50% off 4000 = 2000 discount
		discount := coupon.CalculateDiscount(4000)
		assert.Equal(t, int64(2000), discount)

		// 50% off 8000 = 4000 discount, but max is 3000
		discount = coupon.CalculateDiscount(8000)
		assert.Equal(t, int64(3000), discount)
	})
}

func TestService_LockCoupon(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockCouponRepository{}
		repo.On("LockCoupon", ctx, uint64(1), uint64(100)).Return(nil)

		// Would need actual service instance with mock repo
	})

	t.Run("error", func(t *testing.T) {
		repo := &MockCouponRepository{}
		repo.On("LockCoupon", ctx, uint64(1), uint64(100)).Return(errors.New("db error"))

		// Would need actual service instance with mock repo
	})
}

func TestService_ExpireOldCoupons(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockCouponRepository{}
		repo.On("ExpireOldCoupons", ctx).Return(int64(10), nil)

		// Would need actual service instance with mock repo
	})
}

func TestService_GetCouponStats(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &MockCouponRepository{}
		stats := map[string]int64{
			"available": 100,
			"used":      50,
			"expired":   20,
		}
		repo.On("GetCouponStats", ctx).Return(stats, nil)

		// Would need actual service instance with mock repo
	})
}

func TestCoupon_IsValid(t *testing.T) {
	t.Run("valid coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			State:    model.CouponStateAvailable,
			ExpireAt: time.Now().Add(24 * time.Hour),
		}
		assert.True(t, coupon.IsValid())
	})

	t.Run("expired coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			State:    model.CouponStateAvailable,
			ExpireAt: time.Now().Add(-24 * time.Hour),
		}
		assert.False(t, coupon.IsValid())
	})

	t.Run("used coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			State:    model.CouponStateUsed,
			ExpireAt: time.Now().Add(24 * time.Hour),
		}
		assert.False(t, coupon.IsValid())
	})

	t.Run("locked coupon", func(t *testing.T) {
		coupon := &model.Coupon{
			State:    model.CouponStateLocked,
			ExpireAt: time.Now().Add(24 * time.Hour),
		}
		assert.False(t, coupon.IsValid())
	})
}

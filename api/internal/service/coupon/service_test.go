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
	"gamelink/internal/repository"
	couponrepo "gamelink/internal/repository/coupon"
)

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

func (m *MockCouponRepository) CreateTemplate(ctx context.Context, t *model.CouponTemplate) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockCouponRepository) UpdateTemplate(ctx context.Context, t *model.CouponTemplate) error {
	return m.Called(ctx, t).Error(0)
}

func (m *MockCouponRepository) DeleteTemplate(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
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

func (m *MockCouponRepository) CreateCoupon(ctx context.Context, c *model.Coupon) error {
	return m.Called(ctx, c).Error(0)
}

func (m *MockCouponRepository) IncrementClaimedCount(ctx context.Context, templateID uint64) error {
	return m.Called(ctx, templateID).Error(0)
}

func (m *MockCouponRepository) LockCoupon(ctx context.Context, couponID, orderID uint64) error {
	return m.Called(ctx, couponID, orderID).Error(0)
}

func (m *MockCouponRepository) UnlockCoupon(ctx context.Context, couponID uint64) error {
	return m.Called(ctx, couponID).Error(0)
}

func (m *MockCouponRepository) UseCoupon(ctx context.Context, couponID, orderID uint64, discountCents int64) error {
	return m.Called(ctx, couponID, orderID, discountCents).Error(0)
}

func (m *MockCouponRepository) ExpireOldCoupons(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCouponRepository) GetCouponStats(ctx context.Context) (map[string]int64, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(map[string]int64), args.Error(1)
}

func (m *MockCouponRepository) DeleteCoupon(ctx context.Context, id uint64) error {
	return m.Called(ctx, id).Error(0)
}

func newTestService() (*Service, *MockCouponRepository) {
	repo := new(MockCouponRepository)
	return &Service{repo: repo}, repo
}

// ============================================================================
// Template Tests
// ============================================================================

func TestListTemplates(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	templates := []model.CouponTemplate{
		{Name: "Test1", Type: model.CouponTypeDeduct},
		{Name: "Test2", Type: model.CouponTypeDiscount},
	}

	t.Run("success", func(t *testing.T) {
		opts := couponrepo.TemplateListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListTemplates", ctx, opts).Return(templates, int64(2), nil).Once()

		result, total, err := svc.ListTemplates(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		opts := couponrepo.TemplateListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListTemplates", ctx, opts).Return([]model.CouponTemplate{}, int64(0), errors.New("db error")).Once()

		_, _, err := svc.ListTemplates(ctx, opts)
		assert.Error(t, err)
	})
}

func TestGetTemplate(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	template := &model.CouponTemplate{Name: "Test", Type: model.CouponTypeDeduct}
	template.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(template, nil).Once()

		result, err := svc.GetTemplate(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "Test", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetTemplateByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetTemplate(ctx, 999)
		assert.Error(t, err)
	})
}

func TestGetTemplateByClaimLink(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	template := &model.CouponTemplate{Name: "Test", ClaimLink: "abc123"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetTemplateByClaimLink", ctx, "abc123").Return(template, nil).Once()

		result, err := svc.GetTemplateByClaimLink(ctx, "abc123")
		require.NoError(t, err)
		assert.Equal(t, "abc123", result.ClaimLink)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetTemplateByClaimLink", ctx, "invalid").Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetTemplateByClaimLink(ctx, "invalid")
		assert.Error(t, err)
	})
}

func TestCreateTemplate(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success deduct", func(t *testing.T) {
		template := &model.CouponTemplate{
			Name:              "满减券",
			Type:              model.CouponTypeDeduct,
			DeductAmountCents: 1000,
		}
		mockRepo.On("CreateTemplate", ctx, template).Return(nil).Once()

		err := svc.CreateTemplate(ctx, template)
		require.NoError(t, err)
	})

	t.Run("success discount", func(t *testing.T) {
		template := &model.CouponTemplate{
			Name:         "折扣券",
			Type:         model.CouponTypeDiscount,
			DiscountRate: 0.9,
		}
		mockRepo.On("CreateTemplate", ctx, template).Return(nil).Once()

		err := svc.CreateTemplate(ctx, template)
		require.NoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "", Type: model.CouponTypeDeduct}
		err := svc.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "名称不能为空")
	})

	t.Run("invalid deduct amount", func(t *testing.T) {
		template := &model.CouponTemplate{
			Name:              "满减券",
			Type:              model.CouponTypeDeduct,
			DeductAmountCents: 0,
		}
		err := svc.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "满减金额必须大于0")
	})

	t.Run("invalid discount rate", func(t *testing.T) {
		template := &model.CouponTemplate{
			Name:         "折扣券",
			Type:         model.CouponTypeDiscount,
			DiscountRate: 1.5,
		}
		err := svc.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "折扣率必须在0-1之间")
	})

	t.Run("invalid type", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "无效券", Type: "invalid"}
		err := svc.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的优惠券类型")
	})
}

func TestUpdateTemplate(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	existing := &model.CouponTemplate{Name: "Old", Type: model.CouponTypeDeduct, DeductAmountCents: 500}
	existing.ID = 1

	t.Run("success", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Updated", Type: model.CouponTypeDeduct, DeductAmountCents: 1000}
		template.ID = 1
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(existing, nil).Once()
		mockRepo.On("UpdateTemplate", ctx, template).Return(nil).Once()

		err := svc.UpdateTemplate(ctx, template)
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "Updated", Type: model.CouponTypeDeduct, DeductAmountCents: 1000}
		template.ID = 999
		mockRepo.On("GetTemplateByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		err := svc.UpdateTemplate(ctx, template)
		assert.Error(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		template := &model.CouponTemplate{Name: "", Type: model.CouponTypeDeduct}
		template.ID = 1
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(existing, nil).Once()

		err := svc.UpdateTemplate(ctx, template)
		assert.Error(t, err)
	})
}

func TestDeleteTemplate(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteTemplate", ctx, uint64(1)).Return(nil).Once()
		err := svc.DeleteTemplate(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteTemplate", ctx, uint64(999)).Return(errors.New("not found")).Once()
		err := svc.DeleteTemplate(ctx, 999)
		assert.Error(t, err)
	})
}

func TestBatchUpdateTemplateStatus(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ids := []uint64{1, 2, 3}
		mockRepo.On("BatchUpdateTemplateStatus", ctx, ids, true).Return(int64(3), nil).Once()

		affected, err := svc.BatchUpdateTemplateStatus(ctx, ids, true)
		require.NoError(t, err)
		assert.Equal(t, int64(3), affected)
	})

	t.Run("error", func(t *testing.T) {
		ids := []uint64{1, 2}
		mockRepo.On("BatchUpdateTemplateStatus", ctx, ids, false).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.BatchUpdateTemplateStatus(ctx, ids, false)
		assert.Error(t, err)
	})
}

func TestBatchDeleteTemplates(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ids := []uint64{1, 2}
		mockRepo.On("BatchDeleteTemplates", ctx, ids).Return(int64(2), nil).Once()

		affected, err := svc.BatchDeleteTemplates(ctx, ids)
		require.NoError(t, err)
		assert.Equal(t, int64(2), affected)
	})

	t.Run("error", func(t *testing.T) {
		ids := []uint64{1}
		mockRepo.On("BatchDeleteTemplates", ctx, ids).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.BatchDeleteTemplates(ctx, ids)
		assert.Error(t, err)
	})
}

// ============================================================================
// Coupon Tests
// ============================================================================

func TestListCoupons(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	coupons := []model.Coupon{
		{UserID: 1, State: model.CouponStateAvailable},
		{UserID: 1, State: model.CouponStateUsed},
	}

	t.Run("success", func(t *testing.T) {
		userID := uint64(1)
		opts := couponrepo.CouponListOptions{UserID: &userID, Page: 1, PageSize: 10}
		mockRepo.On("ListCoupons", ctx, opts).Return(coupons, int64(2), nil).Once()

		result, total, err := svc.ListCoupons(ctx, opts)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, result, 2)
	})

	t.Run("error", func(t *testing.T) {
		opts := couponrepo.CouponListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListCoupons", ctx, opts).Return([]model.Coupon{}, int64(0), errors.New("db error")).Once()

		_, _, err := svc.ListCoupons(ctx, opts)
		assert.Error(t, err)
	})
}

func TestGetCoupon(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	coupon := &model.Coupon{UserID: 1, State: model.CouponStateAvailable}
	coupon.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetCouponByID", ctx, uint64(1)).Return(coupon, nil).Once()

		result, err := svc.GetCoupon(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.UserID)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetCouponByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetCoupon(ctx, 999)
		assert.Error(t, err)
	})
}

func TestGetCouponWithTemplate(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	coupon := &model.Coupon{UserID: 1, TemplateID: 1}
	coupon.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetCouponWithTemplate", ctx, uint64(1)).Return(coupon, nil).Once()

		result, err := svc.GetCouponWithTemplate(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.TemplateID)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetCouponWithTemplate", ctx, uint64(999)).Return(nil, errors.New("not found")).Once()

		_, err := svc.GetCouponWithTemplate(ctx, 999)
		assert.Error(t, err)
	})
}

func TestGetUserAvailableCoupons(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	coupons := []model.Coupon{{UserID: 1, State: model.CouponStateAvailable}}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetUserAvailableCoupons", ctx, uint64(1)).Return(coupons, nil).Once()

		result, err := svc.GetUserAvailableCoupons(ctx, 1)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetUserAvailableCoupons", ctx, uint64(999)).Return([]model.Coupon{}, errors.New("db error")).Once()

		_, err := svc.GetUserAvailableCoupons(ctx, 999)
		assert.Error(t, err)
	})
}

func TestClaimCoupon(t *testing.T) {
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:              "满减券",
		Type:              model.CouponTypeDeduct,
		DeductAmountCents: 1000,
		IsActive:          true,
		TotalCount:        100,
		ClaimedCount:      10,
		PerUserLimit:      3,
		ValidityType:      "days",
		ValidityDays:      30,
	}
	template.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(template, nil).Once()
		mockRepo.On("CountUserCouponsFromTemplate", ctx, uint64(100), uint64(1)).Return(int64(0), nil).Once()
		mockRepo.On("CreateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil).Once()
		mockRepo.On("IncrementClaimedCount", ctx, uint64(1)).Return(nil).Once()

		coupon, err := svc.ClaimCoupon(ctx, 100, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), coupon.UserID)
		assert.Equal(t, model.CouponStateAvailable, coupon.State)
	})

	t.Run("template not active", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		inactiveTemplate := &model.CouponTemplate{IsActive: false}
		inactiveTemplate.ID = 2
		mockRepo.On("GetTemplateByID", ctx, uint64(2)).Return(inactiveTemplate, nil).Once()

		_, err := svc.ClaimCoupon(ctx, 100, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已下架")
	})

	t.Run("sold out", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		soldOutTemplate := &model.CouponTemplate{IsActive: true, TotalCount: 10, ClaimedCount: 10}
		soldOutTemplate.ID = 3
		mockRepo.On("GetTemplateByID", ctx, uint64(3)).Return(soldOutTemplate, nil).Once()

		_, err := svc.ClaimCoupon(ctx, 100, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已领完")
	})

	t.Run("user limit reached", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(template, nil).Once()
		mockRepo.On("CountUserCouponsFromTemplate", ctx, uint64(100), uint64(1)).Return(int64(3), nil).Once()

		_, err := svc.ClaimCoupon(ctx, 100, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "领取上限")
	})

	t.Run("fixed expire date", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		fixedExpire := time.Now().AddDate(0, 1, 0)
		fixedTemplate := &model.CouponTemplate{
			Name:              "固定过期券",
			Type:              model.CouponTypeDeduct,
			DeductAmountCents: 500,
			IsActive:          true,
			ValidityType:      "fixed",
			FixedExpireAt:     &fixedExpire,
		}
		fixedTemplate.ID = 4
		mockRepo.On("GetTemplateByID", ctx, uint64(4)).Return(fixedTemplate, nil).Once()
		mockRepo.On("CreateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil).Once()
		mockRepo.On("IncrementClaimedCount", ctx, uint64(4)).Return(nil).Once()

		coupon, err := svc.ClaimCoupon(ctx, 100, 4)
		require.NoError(t, err)
		assert.Equal(t, fixedExpire.Unix(), coupon.ExpireAt.Unix())
	})
}

func TestClaimCouponByLink(t *testing.T) {
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:              "链接券",
		Type:              model.CouponTypeDeduct,
		DeductAmountCents: 500,
		IsActive:          true,
		ClaimLink:         "abc123",
		ValidityType:      "days",
		ValidityDays:      7,
	}
	template.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByClaimLink", ctx, "abc123").Return(template, nil).Once()
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(template, nil).Once()
		mockRepo.On("CreateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil).Once()
		mockRepo.On("IncrementClaimedCount", ctx, uint64(1)).Return(nil).Once()

		coupon, err := svc.ClaimCouponByLink(ctx, 100, "abc123")
		require.NoError(t, err)
		assert.Equal(t, uint64(100), coupon.UserID)
	})

	t.Run("invalid link", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByClaimLink", ctx, "invalid").Return(nil, repository.ErrNotFound).Once()

		_, err := svc.ClaimCouponByLink(ctx, 100, "invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的领取链接")
	})
}

func TestIssueCoupon(t *testing.T) {
	ctx := context.Background()

	template := &model.CouponTemplate{
		Name:         "系统发放券",
		Type:         model.CouponTypeDiscount,
		DiscountRate: 0.9,
		ValidityType: "days",
		ValidityDays: 14,
	}
	template.ID = 1

	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByID", ctx, uint64(1)).Return(template, nil).Once()
		mockRepo.On("CreateCoupon", ctx, mock.AnythingOfType("*model.Coupon")).Return(nil).Once()

		coupon, err := svc.IssueCoupon(ctx, 100, 1, model.CouponSourceVip)
		require.NoError(t, err)
		assert.Equal(t, model.CouponSourceVip, coupon.Source)
	})

	t.Run("template not found", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		mockRepo.On("GetTemplateByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.IssueCoupon(ctx, 100, 999, model.CouponSourceManual)
		assert.Error(t, err)
	})
}

func TestLockCoupon(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("LockCoupon", ctx, uint64(1), uint64(100)).Return(nil).Once()
		err := svc.LockCoupon(ctx, 1, 100)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("LockCoupon", ctx, uint64(999), uint64(100)).Return(errors.New("not found")).Once()
		err := svc.LockCoupon(ctx, 999, 100)
		assert.Error(t, err)
	})
}

func TestUnlockCoupon(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UnlockCoupon", ctx, uint64(1)).Return(nil).Once()
		err := svc.UnlockCoupon(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("UnlockCoupon", ctx, uint64(999)).Return(errors.New("not found")).Once()
		err := svc.UnlockCoupon(ctx, 999)
		assert.Error(t, err)
	})
}

func TestUseCoupon(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("UseCoupon", ctx, uint64(1), uint64(100), int64(1000)).Return(nil).Once()
		err := svc.UseCoupon(ctx, 1, 100, 1000)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("UseCoupon", ctx, uint64(999), uint64(100), int64(500)).Return(errors.New("not found")).Once()
		err := svc.UseCoupon(ctx, 999, 100, 500)
		assert.Error(t, err)
	})
}

func TestCalculateDiscount(t *testing.T) {
	ctx := context.Background()

	t.Run("deduct coupon", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		coupon := &model.Coupon{
			Type:              model.CouponTypeDeduct,
			State:             model.CouponStateAvailable,
			MinAmountCents:    5000,
			DeductAmountCents: 1000,
			ExpireAt:          time.Now().AddDate(0, 0, 7),
		}
		coupon.ID = 1
		mockRepo.On("GetCouponByID", ctx, uint64(1)).Return(coupon, nil).Once()

		discount, err := svc.CalculateDiscount(ctx, 1, 10000)
		require.NoError(t, err)
		assert.Equal(t, int64(1000), discount)
	})

	t.Run("discount coupon", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		coupon := &model.Coupon{
			Type:             model.CouponTypeDiscount,
			State:            model.CouponStateAvailable,
			DiscountRate:     0.9,
			MaxDiscountCents: 2000,
			ExpireAt:         time.Now().AddDate(0, 0, 7),
		}
		coupon.ID = 2
		mockRepo.On("GetCouponByID", ctx, uint64(2)).Return(coupon, nil).Once()

		discount, err := svc.CalculateDiscount(ctx, 2, 10000)
		require.NoError(t, err)
		// 10000 * (1 - 0.9) = 1000, 但浮点计算可能有精度问题
		assert.True(t, discount >= 999 && discount <= 1001)
	})

	t.Run("invalid coupon", func(t *testing.T) {
		mockRepo := new(MockCouponRepository)
		svc := NewCouponService(mockRepo)
		coupon := &model.Coupon{
			State:    model.CouponStateUsed,
			ExpireAt: time.Now().AddDate(0, 0, -1),
		}
		coupon.ID = 3
		mockRepo.On("GetCouponByID", ctx, uint64(3)).Return(coupon, nil).Once()

		_, err := svc.CalculateDiscount(ctx, 3, 10000)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不可用")
	})
}

func TestExpireOldCoupons(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("ExpireOldCoupons", ctx).Return(int64(5), nil).Once()

		affected, err := svc.ExpireOldCoupons(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(5), affected)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("ExpireOldCoupons", ctx).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.ExpireOldCoupons(ctx)
		assert.Error(t, err)
	})
}

func TestGetCouponStats(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		stats := map[string]int64{
			"totalTemplates":     10,
			"activeTemplates":    8,
			"totalCoupons":       100,
			"availableCoupons":   50,
			"usedCoupons":        30,
			"expiredCoupons":     20,
			"totalDiscountCents": 50000,
		}
		mockRepo.On("GetCouponStats", ctx).Return(stats, nil).Once()

		result, err := svc.GetCouponStats(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result["totalTemplates"])
		assert.Equal(t, int64(100), result["totalCoupons"])
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetCouponStats", ctx).Return(nil, errors.New("db error")).Once()

		_, err := svc.GetCouponStats(ctx)
		assert.Error(t, err)
	})
}

func TestDeleteCoupon(t *testing.T) {
	mockRepo := new(MockCouponRepository)
	svc := NewCouponService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteCoupon", ctx, uint64(1)).Return(nil).Once()
		err := svc.DeleteCoupon(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteCoupon", ctx, uint64(999)).Return(errors.New("not found")).Once()
		err := svc.DeleteCoupon(ctx, 999)
		assert.Error(t, err)
	})
}

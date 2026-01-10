package recharge

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	rechargerepo "gamelink/internal/repository/recharge"
)

// ============================================================================
// Mock Implementations
// ============================================================================

type MockRechargeRepo struct {
	mock.Mock
}

func (m *MockRechargeRepo) ListOptions(ctx context.Context, opts rechargerepo.OptionListOptions) ([]model.RechargeOption, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RechargeOption), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetActiveOptions(ctx context.Context, vipLevel *uint64) ([]model.RechargeOption, error) {
	args := m.Called(ctx, vipLevel)
	return args.Get(0).([]model.RechargeOption), args.Error(1)
}

func (m *MockRechargeRepo) GetOptionByID(ctx context.Context, id uint64) (*model.RechargeOption, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeOption), args.Error(1)
}

func (m *MockRechargeRepo) CreateOption(ctx context.Context, option *model.RechargeOption) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockRechargeRepo) UpdateOption(ctx context.Context, option *model.RechargeOption) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockRechargeRepo) DeleteOption(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRechargeRepo) BatchUpdateOptionStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	args := m.Called(ctx, ids, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepo) BatchDeleteOptions(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepo) IncrementPurchaseCount(ctx context.Context, optionID uint64) error {
	args := m.Called(ctx, optionID)
	return args.Error(0)
}

func (m *MockRechargeRepo) ListRecords(ctx context.Context, opts rechargerepo.RecordListOptions) ([]model.RechargeRecord, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RechargeRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepo) GetRecordByID(ctx context.Context, id uint64) (*model.RechargeRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepo) GetRecordByOrderNo(ctx context.Context, orderNo string) (*model.RechargeRecord, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepo) CreateRecord(ctx context.Context, record *model.RechargeRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRechargeRepo) MarkAsPaid(ctx context.Context, id uint64, providerTradeNo string) error {
	args := m.Called(ctx, id, providerTradeNo)
	return args.Error(0)
}

func (m *MockRechargeRepo) MarkAsRefunded(ctx context.Context, id uint64, refundAmount int64, reason, providerNo string) error {
	args := m.Called(ctx, id, refundAmount, reason, providerNo)
	return args.Error(0)
}

func (m *MockRechargeRepo) MarkCouponIssued(ctx context.Context, id uint64, couponIDs string) error {
	args := m.Called(ctx, id, couponIDs)
	return args.Error(0)
}

func (m *MockRechargeRepo) CountUserPurchases(ctx context.Context, userID, optionID uint64) (int64, error) {
	args := m.Called(ctx, userID, optionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepo) GetUserRecords(ctx context.Context, userID uint64, limit int) ([]model.RechargeRecord, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepo) GetRechargeStats(ctx context.Context) (map[string]any, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockRechargeRepo) CancelExpiredRecords(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockWalletRepo struct {
	mock.Mock
}

func (m *MockWalletRepo) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepo) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepo) SaveWithOptimisticLock(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepo) UpdateBalanceWithLock(ctx context.Context, userID uint64, delta int64, maxRetries int) (*model.Wallet, error) {
	args := m.Called(ctx, userID, delta, maxRetries)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

type MockCouponService struct {
	mock.Mock
}

func (m *MockCouponService) IssueCoupon(ctx context.Context, userID, templateID uint64, source model.CouponSource) (*model.Coupon, error) {
	args := m.Called(ctx, userID, templateID, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Coupon), args.Error(1)
}

// ============================================================================
// Test Helper
// ============================================================================

func setupTestService() (*Service, *MockRechargeRepo, *MockWalletRepo, *MockCouponService) {
	mockRepo := new(MockRechargeRepo)
	mockWalletRepo := new(MockWalletRepo)
	mockCouponSvc := new(MockCouponService)
	svc := NewRechargeService(mockRepo, mockWalletRepo, mockCouponSvc)
	return svc, mockRepo, mockWalletRepo, mockCouponSvc
}

// ============================================================================
// Option Tests
// ============================================================================

func TestListOptions(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		opts := rechargerepo.OptionListOptions{Page: 1, PageSize: 10}
		expected := []model.RechargeOption{{Name: "Option1"}, {Name: "Option2"}}
		mockRepo.On("ListOptions", ctx, opts).Return(expected, int64(2), nil).Once()

		result, total, err := svc.ListOptions(ctx, opts)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		opts := rechargerepo.OptionListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListOptions", ctx, opts).Return([]model.RechargeOption{}, int64(0), errors.New("db error")).Once()

		_, _, err := svc.ListOptions(ctx, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "list options")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetActiveOptions(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []model.RechargeOption{{Name: "Active1", IsActive: true}}
		mockRepo.On("GetActiveOptions", ctx, (*uint64)(nil)).Return(expected, nil).Once()

		result, err := svc.GetActiveOptions(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("with vip level", func(t *testing.T) {
		vipLevel := uint64(2)
		expected := []model.RechargeOption{{Name: "VIP Option"}}
		mockRepo.On("GetActiveOptions", ctx, &vipLevel).Return(expected, nil).Once()

		result, err := svc.GetActiveOptions(ctx, &vipLevel)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetActiveOptions", ctx, (*uint64)(nil)).Return([]model.RechargeOption{}, errors.New("db error")).Once()

		_, err := svc.GetActiveOptions(ctx, nil)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetOption(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &model.RechargeOption{Name: "Test Option"}
		expected.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(expected, nil).Once()

		result, err := svc.GetOption(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetOptionByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetOption(ctx, 999)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateOption(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "New Option",
			AmountCents: 1000,
			BonusCents:  100,
		}
		mockRepo.On("CreateOption", ctx, mock.AnythingOfType("*model.RechargeOption")).Return(nil).Once()

		err := svc.CreateOption(ctx, option)
		assert.NoError(t, err)
		assert.Equal(t, int64(1100), option.TotalCents)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "",
			AmountCents: 1000,
		}

		err := svc.CreateOption(ctx, option)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "名称不能为空")
	})

	t.Run("invalid amount", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "Test",
			AmountCents: 0,
		}

		err := svc.CreateOption(ctx, option)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "金额必须大于0")
	})

	t.Run("negative bonus", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "Test",
			AmountCents: 1000,
			BonusCents:  -100,
		}

		err := svc.CreateOption(ctx, option)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "赠送金额不能为负数")
	})

	t.Run("repo error", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "Test",
			AmountCents: 1000,
		}
		mockRepo.On("CreateOption", ctx, mock.AnythingOfType("*model.RechargeOption")).Return(errors.New("db error")).Once()

		err := svc.CreateOption(ctx, option)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create option")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateOption(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existing := &model.RechargeOption{Name: "Existing"}
		existing.ID = 1
		option := &model.RechargeOption{
			Name:        "Updated",
			AmountCents: 2000,
			BonusCents:  200,
		}
		option.ID = 1

		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(existing, nil).Once()
		mockRepo.On("UpdateOption", ctx, mock.AnythingOfType("*model.RechargeOption")).Return(nil).Once()

		err := svc.UpdateOption(ctx, option)
		assert.NoError(t, err)
		assert.Equal(t, int64(2200), option.TotalCents)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		option := &model.RechargeOption{Name: "Test", AmountCents: 1000}
		option.ID = 999
		mockRepo.On("GetOptionByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		err := svc.UpdateOption(ctx, option)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		existing := &model.RechargeOption{Name: "Existing"}
		existing.ID = 1
		option := &model.RechargeOption{Name: "", AmountCents: 1000}
		option.ID = 1

		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(existing, nil).Once()

		err := svc.UpdateOption(ctx, option)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "名称不能为空")
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteOption(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteOption", ctx, uint64(1)).Return(nil).Once()

		err := svc.DeleteOption(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("DeleteOption", ctx, uint64(999)).Return(errors.New("not found")).Once()

		err := svc.DeleteOption(ctx, 999)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBatchUpdateOptionStatus(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ids := []uint64{1, 2, 3}
		mockRepo.On("BatchUpdateOptionStatus", ctx, ids, true).Return(int64(3), nil).Once()

		affected, err := svc.BatchUpdateOptionStatus(ctx, ids, true)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), affected)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		ids := []uint64{1, 2}
		mockRepo.On("BatchUpdateOptionStatus", ctx, ids, false).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.BatchUpdateOptionStatus(ctx, ids, false)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBatchDeleteOptions(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		ids := []uint64{1, 2}
		mockRepo.On("BatchDeleteOptions", ctx, ids).Return(int64(2), nil).Once()

		affected, err := svc.BatchDeleteOptions(ctx, ids)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), affected)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		ids := []uint64{1}
		mockRepo.On("BatchDeleteOptions", ctx, ids).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.BatchDeleteOptions(ctx, ids)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Record Tests
// ============================================================================

func TestListRecords(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		opts := rechargerepo.RecordListOptions{Page: 1, PageSize: 10}
		expected := []model.RechargeRecord{{OrderNo: "RC001"}, {OrderNo: "RC002"}}
		mockRepo.On("ListRecords", ctx, opts).Return(expected, int64(2), nil).Once()

		result, total, err := svc.ListRecords(ctx, opts)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		opts := rechargerepo.RecordListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListRecords", ctx, opts).Return([]model.RechargeRecord{}, int64(0), errors.New("db error")).Once()

		_, _, err := svc.ListRecords(ctx, opts)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecord(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &model.RechargeRecord{OrderNo: "RC001"}
		expected.ID = 1
		mockRepo.On("GetRecordByID", ctx, uint64(1)).Return(expected, nil).Once()

		result, err := svc.GetRecord(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRecordByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetRecord(ctx, 999)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecordByOrderNo(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := &model.RechargeRecord{OrderNo: "RC001"}
		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(expected, nil).Once()

		result, err := svc.GetRecordByOrderNo(ctx, "RC001")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRecordByOrderNo", ctx, "INVALID").Return(nil, repository.ErrNotFound).Once()

		_, err := svc.GetRecordByOrderNo(ctx, "INVALID")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserRecords(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []model.RechargeRecord{{OrderNo: "RC001"}}
		mockRepo.On("GetUserRecords", ctx, uint64(1), 10).Return(expected, nil).Once()

		result, err := svc.GetUserRecords(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("default limit", func(t *testing.T) {
		expected := []model.RechargeRecord{{OrderNo: "RC001"}}
		mockRepo.On("GetUserRecords", ctx, uint64(1), 20).Return(expected, nil).Once()

		result, err := svc.GetUserRecords(ctx, 1, 0)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetUserRecords", ctx, uint64(1), 10).Return([]model.RechargeRecord{}, errors.New("db error")).Once()

		_, err := svc.GetUserRecords(ctx, 1, 10)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateRechargeOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:        "Test Option",
			AmountCents: 1000,
			BonusCents:  100,
			TotalCents:  1100,
			IsActive:    true,
		}
		option.ID = 1

		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()
		mockRepo.On("CreateRecord", ctx, mock.AnythingOfType("*model.RechargeRecord")).Return(nil).Once()

		record, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "127.0.0.1", "Mozilla")
		assert.NoError(t, err)
		assert.NotNil(t, record)
		assert.Equal(t, uint64(100), record.UserID)
		assert.Equal(t, int64(1000), record.AmountCents)
		assert.Equal(t, int64(100), record.BonusCents)
		assert.Equal(t, int64(1100), record.TotalCents)
		assert.Equal(t, model.RechargeStatusPending, record.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("option not found", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		mockRepo.On("GetOptionByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 999, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")
		mockRepo.AssertExpectations(t)
	})

	t.Run("option inactive", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:     "Inactive",
			IsActive: false,
		}
		option.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已下架")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sold out", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:          "Sold Out",
			IsActive:      true,
			TotalLimit:    10,
			PurchaseCount: 10,
		}
		option.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "售罄")
		mockRepo.AssertExpectations(t)
	})

	t.Run("per user limit exceeded", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:         "Limited",
			IsActive:     true,
			PerUserLimit: 2,
		}
		option.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()
		mockRepo.On("CountUserPurchases", ctx, uint64(100), uint64(1)).Return(int64(2), nil).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "购买上限")
		mockRepo.AssertExpectations(t)
	})

	t.Run("count purchases error", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:         "Limited",
			IsActive:     true,
			PerUserLimit: 2,
		}
		option.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()
		mockRepo.On("CountUserPurchases", ctx, uint64(100), uint64(1)).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "count purchases")
		mockRepo.AssertExpectations(t)
	})

	t.Run("create record error", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		option := &model.RechargeOption{
			Name:        "Test",
			AmountCents: 1000,
			IsActive:    true,
		}
		option.ID = 1
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()
		mockRepo.On("CreateRecord", ctx, mock.AnythingOfType("*model.RechargeRecord")).Return(errors.New("db error")).Once()

		_, err := svc.CreateRechargeOrder(ctx, 100, 1, "wechat", "wechat_h5", "", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create record")
		mockRepo.AssertExpectations(t)
	})
}

func TestHandlePaymentCallback(t *testing.T) {
	ctx := context.Background()

	t.Run("success without wallet", func(t *testing.T) {
		mockRepo := new(MockRechargeRepo)
		svc := NewRechargeService(mockRepo, nil, nil)

		record := &model.RechargeRecord{
			UserID:      100,
			AmountCents: 1000,
			TotalCents:  1100,
			Status:      model.RechargeStatusPending,
			OrderNo:     "RC001",
		}
		record.ID = 1

		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()
		mockRepo.On("MarkAsPaid", ctx, uint64(1), "TRADE123").Return(nil).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with wallet", func(t *testing.T) {
		svc, mockRepo, mockWalletRepo, _ := setupTestService()

		record := &model.RechargeRecord{
			UserID:      100,
			AmountCents: 1000,
			TotalCents:  1100,
			Status:      model.RechargeStatusPending,
			OrderNo:     "RC001",
		}
		record.ID = 1

		wallet := &model.Wallet{UserID: 100, BalanceCents: 500}

		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()
		mockRepo.On("MarkAsPaid", ctx, uint64(1), "TRADE123").Return(nil).Once()
		mockWalletRepo.On("GetByUserID", ctx, uint64(100)).Return(wallet, nil).Once()
		mockWalletRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.NoError(t, err)
		assert.Equal(t, int64(1600), wallet.BalanceCents) // 500 + 1100
		mockRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("success with wallet not found creates new", func(t *testing.T) {
		svc, mockRepo, mockWalletRepo, _ := setupTestService()

		record := &model.RechargeRecord{
			UserID:      100,
			AmountCents: 1000,
			TotalCents:  1100,
			Status:      model.RechargeStatusPending,
			OrderNo:     "RC001",
		}
		record.ID = 1

		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()
		mockRepo.On("MarkAsPaid", ctx, uint64(1), "TRADE123").Return(nil).Once()
		mockWalletRepo.On("GetByUserID", ctx, uint64(100)).Return(nil, repository.ErrNotFound).Once()
		mockWalletRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("success with coupon issuance", func(t *testing.T) {
		svc, mockRepo, mockWalletRepo, mockCouponSvc := setupTestService()

		optionID := uint64(1)
		templateID := uint64(10)
		record := &model.RechargeRecord{
			UserID:      100,
			OptionID:    &optionID,
			AmountCents: 1000,
			TotalCents:  1100,
			Status:      model.RechargeStatusPending,
			OrderNo:     "RC001",
		}
		record.ID = 1

		option := &model.RechargeOption{
			CouponTemplateID: &templateID,
			CouponCount:      2,
		}
		option.ID = 1

		coupon := &model.Coupon{}
		coupon.ID = 100

		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()
		mockRepo.On("MarkAsPaid", ctx, uint64(1), "TRADE123").Return(nil).Once()
		mockWalletRepo.On("GetByUserID", ctx, uint64(100)).Return(nil, repository.ErrNotFound).Once()
		mockWalletRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Once()
		mockRepo.On("IncrementPurchaseCount", ctx, uint64(1)).Return(nil).Once()
		mockRepo.On("GetOptionByID", ctx, uint64(1)).Return(option, nil).Once()
		mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(10), model.CouponSourceRecharge).Return(coupon, nil).Twice()
		mockRepo.On("MarkCouponIssued", ctx, uint64(1), mock.AnythingOfType("string")).Return(nil).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
		mockCouponSvc.AssertExpectations(t)
	})

	t.Run("record not found", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		mockRepo.On("GetRecordByOrderNo", ctx, "INVALID").Return(nil, repository.ErrNotFound).Once()

		err := svc.HandlePaymentCallback(ctx, "INVALID", "TRADE123")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("wrong status", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		record := &model.RechargeRecord{
			Status:  model.RechargeStatusPaid,
			OrderNo: "RC001",
		}
		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "状态不正确")
		mockRepo.AssertExpectations(t)
	})

	t.Run("mark as paid error", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		record := &model.RechargeRecord{
			Status:  model.RechargeStatusPending,
			OrderNo: "RC001",
		}
		record.ID = 1
		mockRepo.On("GetRecordByOrderNo", ctx, "RC001").Return(record, nil).Once()
		mockRepo.On("MarkAsPaid", ctx, uint64(1), "TRADE123").Return(errors.New("db error")).Once()

		err := svc.HandlePaymentCallback(ctx, "RC001", "TRADE123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mark as paid")
		mockRepo.AssertExpectations(t)
	})
}

func TestRefundRecord(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, mockRepo, mockWalletRepo, _ := setupTestService()

		record := &model.RechargeRecord{
			UserID:          100,
			AmountCents:     1000,
			TotalCents:      1100,
			Status:          model.RechargeStatusPaid,
			ProviderTradeNo: "TRADE123",
		}
		record.ID = 1

		wallet := &model.Wallet{UserID: 100, BalanceCents: 2000}

		mockRepo.On("GetRecordByID", ctx, uint64(1)).Return(record, nil).Once()
		mockRepo.On("MarkAsRefunded", ctx, uint64(1), int64(1000), "用户申请", "REFUND_TRADE123").Return(nil).Once()
		mockWalletRepo.On("GetByUserID", ctx, uint64(100)).Return(wallet, nil).Once()
		mockWalletRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Once()

		err := svc.RefundRecord(ctx, 1, "用户申请")
		assert.NoError(t, err)
		assert.Equal(t, int64(900), wallet.BalanceCents) // 2000 - 1100
		mockRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("success with insufficient balance", func(t *testing.T) {
		svc, mockRepo, mockWalletRepo, _ := setupTestService()

		record := &model.RechargeRecord{
			UserID:          100,
			AmountCents:     1000,
			TotalCents:      1100,
			Status:          model.RechargeStatusPaid,
			ProviderTradeNo: "TRADE123",
		}
		record.ID = 1

		wallet := &model.Wallet{UserID: 100, BalanceCents: 500}

		mockRepo.On("GetRecordByID", ctx, uint64(1)).Return(record, nil).Once()
		mockRepo.On("MarkAsRefunded", ctx, uint64(1), int64(1000), "用户申请", "REFUND_TRADE123").Return(nil).Once()
		mockWalletRepo.On("GetByUserID", ctx, uint64(100)).Return(wallet, nil).Once()
		mockWalletRepo.On("Save", ctx, mock.AnythingOfType("*model.Wallet")).Return(nil).Once()

		err := svc.RefundRecord(ctx, 1, "用户申请")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), wallet.BalanceCents) // Should not go negative
		mockRepo.AssertExpectations(t)
		mockWalletRepo.AssertExpectations(t)
	})

	t.Run("record not found", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		mockRepo.On("GetRecordByID", ctx, uint64(999)).Return(nil, repository.ErrNotFound).Once()

		err := svc.RefundRecord(ctx, 999, "reason")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("cannot refund", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		record := &model.RechargeRecord{
			Status: model.RechargeStatusPending, // Not paid, cannot refund
		}
		record.ID = 1
		mockRepo.On("GetRecordByID", ctx, uint64(1)).Return(record, nil).Once()

		err := svc.RefundRecord(ctx, 1, "reason")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不可退款")
		mockRepo.AssertExpectations(t)
	})

	t.Run("mark as refunded error", func(t *testing.T) {
		svc, mockRepo, _, _ := setupTestService()
		record := &model.RechargeRecord{
			AmountCents:     1000,
			Status:          model.RechargeStatusPaid,
			ProviderTradeNo: "TRADE123",
		}
		record.ID = 1
		mockRepo.On("GetRecordByID", ctx, uint64(1)).Return(record, nil).Once()
		mockRepo.On("MarkAsRefunded", ctx, uint64(1), int64(1000), "reason", "REFUND_TRADE123").Return(errors.New("db error")).Once()

		err := svc.RefundRecord(ctx, 1, "reason")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mark as refunded")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRechargeStats(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := map[string]any{
			"totalAmountCents": int64(100000),
			"totalCount":       int64(50),
		}
		mockRepo.On("GetRechargeStats", ctx).Return(expected, nil).Once()

		result, err := svc.GetRechargeStats(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetRechargeStats", ctx).Return(nil, errors.New("db error")).Once()

		_, err := svc.GetRechargeStats(ctx)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCancelExpiredRecords(t *testing.T) {
	svc, mockRepo, _, _ := setupTestService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CancelExpiredRecords", ctx).Return(int64(5), nil).Once()

		affected, err := svc.CancelExpiredRecords(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), affected)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("CancelExpiredRecords", ctx).Return(int64(0), errors.New("db error")).Once()

		_, err := svc.CancelExpiredRecords(ctx)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

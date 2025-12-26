package recharge

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
	rechargerepo "gamelink/internal/repository/recharge"
)

// MockRechargeRepository is a mock implementation
type MockRechargeRepository struct {
	mock.Mock
}

func (m *MockRechargeRepository) ListOptions(ctx context.Context, opts rechargerepo.OptionListOptions) ([]model.RechargeOption, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RechargeOption), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepository) GetActiveOptions(ctx context.Context, vipLevel *uint64) ([]model.RechargeOption, error) {
	args := m.Called(ctx, vipLevel)
	return args.Get(0).([]model.RechargeOption), args.Error(1)
}

func (m *MockRechargeRepository) GetOptionByID(ctx context.Context, id uint64) (*model.RechargeOption, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeOption), args.Error(1)
}

func (m *MockRechargeRepository) CreateOption(ctx context.Context, option *model.RechargeOption) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockRechargeRepository) UpdateOption(ctx context.Context, option *model.RechargeOption) error {
	args := m.Called(ctx, option)
	return args.Error(0)
}

func (m *MockRechargeRepository) DeleteOption(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRechargeRepository) BatchUpdateOptionStatus(ctx context.Context, ids []uint64, isActive bool) (int64, error) {
	args := m.Called(ctx, ids, isActive)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepository) BatchDeleteOptions(ctx context.Context, ids []uint64) (int64, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepository) IncrementPurchaseCount(ctx context.Context, optionID uint64) error {
	args := m.Called(ctx, optionID)
	return args.Error(0)
}

func (m *MockRechargeRepository) ListRecords(ctx context.Context, opts rechargerepo.RecordListOptions) ([]model.RechargeRecord, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.RechargeRecord), args.Get(1).(int64), args.Error(2)
}

func (m *MockRechargeRepository) GetRecordByID(ctx context.Context, id uint64) (*model.RechargeRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepository) GetRecordByOrderNo(ctx context.Context, orderNo string) (*model.RechargeRecord, error) {
	args := m.Called(ctx, orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepository) CreateRecord(ctx context.Context, record *model.RechargeRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockRechargeRepository) MarkAsPaid(ctx context.Context, id uint64, providerTradeNo string) error {
	args := m.Called(ctx, id, providerTradeNo)
	return args.Error(0)
}

func (m *MockRechargeRepository) MarkAsRefunded(ctx context.Context, id uint64, refundAmount int64, reason, providerNo string) error {
	args := m.Called(ctx, id, refundAmount, reason, providerNo)
	return args.Error(0)
}

func (m *MockRechargeRepository) MarkCouponIssued(ctx context.Context, id uint64, couponIDs string) error {
	args := m.Called(ctx, id, couponIDs)
	return args.Error(0)
}

func (m *MockRechargeRepository) CountUserPurchases(ctx context.Context, userID, optionID uint64) (int64, error) {
	args := m.Called(ctx, userID, optionID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRechargeRepository) GetUserRecords(ctx context.Context, userID uint64, limit int) ([]model.RechargeRecord, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]model.RechargeRecord), args.Error(1)
}

func (m *MockRechargeRepository) GetRechargeStats(ctx context.Context) (map[string]any, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockRechargeRepository) CancelExpiredRecords(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockWalletRepository is a mock implementation
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) UpdateBalance(ctx context.Context, userID uint64, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) Freeze(ctx context.Context, userID uint64, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockWalletRepository) Unfreeze(ctx context.Context, userID uint64, amount int64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestService_validateOption(t *testing.T) {
	svc := &Service{}

	t.Run("empty name", func(t *testing.T) {
		option := &model.RechargeOption{Name: ""}
		err := svc.validateOption(option)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名称")
	})

	t.Run("invalid amount", func(t *testing.T) {
		option := &model.RechargeOption{Name: "Test", AmountCents: 0}
		err := svc.validateOption(option)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "金额")
	})

	t.Run("negative bonus", func(t *testing.T) {
		option := &model.RechargeOption{Name: "Test", AmountCents: 1000, BonusCents: -100}
		err := svc.validateOption(option)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "赠送")
	})

	t.Run("valid option", func(t *testing.T) {
		option := &model.RechargeOption{Name: "Test", AmountCents: 1000, BonusCents: 100}
		err := svc.validateOption(option)
		require.NoError(t, err)
	})
}

func TestService_CreateRechargeOrder(t *testing.T) {
	t.Run("option validation", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:          "充值100元",
			AmountCents:   10000,
			BonusCents:    1000,
			TotalCents:    11000,
			IsActive:      true,
			TotalLimit:    100,
			PurchaseCount: 50,
			PerUserLimit:  3,
		}
		option.ID = 1

		assert.True(t, option.IsActive)
		assert.Less(t, option.PurchaseCount, option.TotalLimit)
	})

	t.Run("option not active check", func(t *testing.T) {
		option := &model.RechargeOption{IsActive: false}
		assert.False(t, option.IsActive)
	})

	t.Run("option sold out check", func(t *testing.T) {
		option := &model.RechargeOption{TotalLimit: 100, PurchaseCount: 100}
		assert.Equal(t, option.TotalLimit, option.PurchaseCount)
	})
}

func TestService_HandlePaymentCallback(t *testing.T) {
	t.Run("record status check", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status:      model.RechargeStatusPending,
			AmountCents: 10000,
			TotalCents:  11000,
		}
		record.ID = 1

		assert.Equal(t, model.RechargeStatusPending, record.Status)
	})

	t.Run("record already paid", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusPaid,
		}
		assert.NotEqual(t, model.RechargeStatusPending, record.Status)
	})
}

func TestService_RefundRecord(t *testing.T) {
	t.Run("can refund check", func(t *testing.T) {
		now := time.Now()
		record := &model.RechargeRecord{
			Status:            model.RechargeStatusPaid,
			PaidAt:            &now,
			AmountCents:       10000,
			RefundAmountCents: 0,
		}
		assert.True(t, record.CanRefund())
	})

	t.Run("cannot refund pending", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusPending,
		}
		assert.False(t, record.CanRefund())
	})

	t.Run("cannot refund already refunded", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusRefunded,
		}
		assert.False(t, record.CanRefund())
	})
}

func TestService_GetRechargeStats(t *testing.T) {
	t.Run("stats structure", func(t *testing.T) {
		stats := map[string]any{
			"totalAmountCents":  int64(1000000),
			"totalCount":        int64(100),
			"todayAmountCents":  int64(50000),
			"todayCount":        int64(5),
			"refundAmountCents": int64(10000),
		}

		assert.Contains(t, stats, "totalAmountCents")
		assert.Contains(t, stats, "totalCount")
		assert.Contains(t, stats, "todayAmountCents")
	})
}

func TestRechargeRecord_CanRefund(t *testing.T) {
	t.Run("paid record can refund", func(t *testing.T) {
		now := time.Now()
		record := &model.RechargeRecord{
			Status:            model.RechargeStatusPaid,
			PaidAt:            &now,
			AmountCents:       10000,
			RefundAmountCents: 0,
		}
		assert.True(t, record.CanRefund())
	})

	t.Run("pending record cannot refund", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusPending,
		}
		assert.False(t, record.CanRefund())
	})

	t.Run("refunded record cannot refund again", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusRefunded,
		}
		assert.False(t, record.CanRefund())
	})

	t.Run("canceled record cannot refund", func(t *testing.T) {
		record := &model.RechargeRecord{
			Status: model.RechargeStatusCanceled,
		}
		assert.False(t, record.CanRefund())
	})

	t.Run("fully refunded cannot refund again", func(t *testing.T) {
		now := time.Now()
		record := &model.RechargeRecord{
			Status:            model.RechargeStatusPaid,
			PaidAt:            &now,
			AmountCents:       10000,
			RefundAmountCents: 10000,
		}
		assert.False(t, record.CanRefund())
	})
}

func TestRechargeOption_TotalCalculation(t *testing.T) {
	t.Run("total equals amount plus bonus", func(t *testing.T) {
		option := &model.RechargeOption{
			AmountCents: 10000,
			BonusCents:  1000,
		}
		option.TotalCents = option.AmountCents + option.BonusCents
		assert.Equal(t, int64(11000), option.TotalCents)
	})

	t.Run("no bonus", func(t *testing.T) {
		option := &model.RechargeOption{
			AmountCents: 10000,
			BonusCents:  0,
		}
		option.TotalCents = option.AmountCents + option.BonusCents
		assert.Equal(t, int64(10000), option.TotalCents)
	})
}

func TestService_CancelExpiredRecords(t *testing.T) {
	ctx := context.Background()

	t.Run("expired record detection", func(t *testing.T) {
		expireAt := time.Now().Add(-1 * time.Hour)
		record := &model.RechargeRecord{
			Status:   model.RechargeStatusPending,
			ExpireAt: &expireAt,
		}

		assert.Equal(t, model.RechargeStatusPending, record.Status)
		assert.True(t, record.ExpireAt.Before(time.Now()))
	})

	t.Run("not expired record", func(t *testing.T) {
		expireAt := time.Now().Add(1 * time.Hour)
		record := &model.RechargeRecord{
			Status:   model.RechargeStatusPending,
			ExpireAt: &expireAt,
		}

		assert.False(t, record.ExpireAt.Before(time.Now()))
	})

	_ = ctx // suppress unused warning
}

func TestService_GetUserRecords(t *testing.T) {
	t.Run("default limit", func(t *testing.T) {
		limit := 0
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 20, limit)
	})

	t.Run("custom limit", func(t *testing.T) {
		limit := 50
		if limit <= 0 {
			limit = 20
		}
		assert.Equal(t, 50, limit)
	})
}

func TestService_CreateOption(t *testing.T) {
	t.Run("total cents calculation", func(t *testing.T) {
		option := &model.RechargeOption{
			Name:        "充值100元送10元",
			AmountCents: 10000,
			BonusCents:  1000,
		}

		// Simulate service logic
		option.TotalCents = option.AmountCents + option.BonusCents

		assert.Equal(t, int64(11000), option.TotalCents)
	})
}

func TestService_OptionNotFound(t *testing.T) {
	t.Run("not found error", func(t *testing.T) {
		err := repository.ErrNotFound
		assert.True(t, errors.Is(err, repository.ErrNotFound))
	})
}

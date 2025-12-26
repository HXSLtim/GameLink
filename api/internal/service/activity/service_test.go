package activity

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
	activityrepo "gamelink/internal/repository/activity"
)

// MockActivityRepository is a mock implementation of ActivityRepository
type MockActivityRepository struct {
	mock.Mock
}

func (m *MockActivityRepository) ListActivities(ctx context.Context, opts activityrepo.ActivityListOptions) ([]model.Activity, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.Activity), args.Get(1).(int64), args.Error(2)
}

func (m *MockActivityRepository) GetActiveActivities(ctx context.Context) ([]model.Activity, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Activity), args.Error(1)
}

func (m *MockActivityRepository) GetVisibleActivities(ctx context.Context) ([]model.Activity, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Activity), args.Error(1)
}

func (m *MockActivityRepository) GetActivityByID(ctx context.Context, id uint64) (*model.Activity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Activity), args.Error(1)
}

func (m *MockActivityRepository) CreateActivity(ctx context.Context, activity *model.Activity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *MockActivityRepository) UpdateActivity(ctx context.Context, activity *model.Activity) error {
	args := m.Called(ctx, activity)
	return args.Error(0)
}

func (m *MockActivityRepository) DeleteActivity(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockActivityRepository) UpdateActivityStatus(ctx context.Context, id uint64, status model.ActivityStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockActivityRepository) GetRewardByID(ctx context.Context, id uint64) (*model.ActivityReward, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ActivityReward), args.Error(1)
}

func (m *MockActivityRepository) GetRewardsByActivityID(ctx context.Context, activityID uint64) ([]model.ActivityReward, error) {
	args := m.Called(ctx, activityID)
	return args.Get(0).([]model.ActivityReward), args.Error(1)
}

func (m *MockActivityRepository) CreateReward(ctx context.Context, reward *model.ActivityReward) error {
	args := m.Called(ctx, reward)
	return args.Error(0)
}

func (m *MockActivityRepository) UpdateReward(ctx context.Context, reward *model.ActivityReward) error {
	args := m.Called(ctx, reward)
	return args.Error(0)
}

func (m *MockActivityRepository) DeleteReward(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockActivityRepository) CountUserParticipations(ctx context.Context, userID, activityID uint64) (int64, error) {
	args := m.Called(ctx, userID, activityID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockActivityRepository) CountTodayParticipations(ctx context.Context, activityID uint64) (int64, error) {
	args := m.Called(ctx, activityID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockActivityRepository) CreateParticipation(ctx context.Context, participation *model.ActivityParticipation) error {
	args := m.Called(ctx, participation)
	return args.Error(0)
}

func (m *MockActivityRepository) GetUserParticipations(ctx context.Context, userID uint64, limit int) ([]model.ActivityParticipation, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]model.ActivityParticipation), args.Error(1)
}

func (m *MockActivityRepository) ListParticipations(ctx context.Context, opts activityrepo.ParticipationListOptions) ([]model.ActivityParticipation, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.ActivityParticipation), args.Get(1).(int64), args.Error(2)
}

func (m *MockActivityRepository) IncrementParticipants(ctx context.Context, activityID uint64) error {
	args := m.Called(ctx, activityID)
	return args.Error(0)
}

func (m *MockActivityRepository) DecrementRewardStock(ctx context.Context, rewardID uint64) error {
	args := m.Called(ctx, rewardID)
	return args.Error(0)
}

func (m *MockActivityRepository) IncrementDailyStats(ctx context.Context, activityID uint64) error {
	args := m.Called(ctx, activityID)
	return args.Error(0)
}

func (m *MockActivityRepository) GetActivityStats(ctx context.Context, activityID uint64) (map[string]any, error) {
	args := m.Called(ctx, activityID)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockActivityRepository) GetAllActivityStats(ctx context.Context) (map[string]any, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *MockActivityRepository) ResetTodayParticipants(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockCouponService is a mock implementation of CouponService
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

// Helper function to create test service
func createTestActivityService() (*Service, *MockActivityRepository, *MockCouponService) {
	mockRepo := new(MockActivityRepository)
	mockCouponSvc := new(MockCouponService)
	svc := &Service{repo: mockRepo, couponSvc: mockCouponSvc}
	return svc, mockRepo, mockCouponSvc
}

// Helper function to create valid activity
func createValidActivity() *model.Activity {
	now := time.Now()
	return &model.Activity{
		Name:    "Test Activity",
		StartAt: now.Add(1 * time.Hour),
		EndAt:   now.Add(24 * time.Hour),
		Status:  model.ActivityStatusDraft,
	}
}

// ============================================================================
// Activity Management Tests
// ============================================================================

func TestService_ListActivities(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		opts := activityrepo.ActivityListOptions{Page: 1, PageSize: 10}
		activities := []model.Activity{{Name: "Activity 1"}, {Name: "Activity 2"}}
		mockRepo.On("ListActivities", ctx, opts).Return(activities, int64(2), nil).Once()

		result, total, err := svc.ListActivities(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		opts := activityrepo.ActivityListOptions{Page: 1, PageSize: 10}
		mockRepo.On("ListActivities", ctx, opts).Return([]model.Activity{}, int64(0), errors.New("db error")).Once()

		_, _, err := svc.ListActivities(ctx, opts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "list activities")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetActiveActivities(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activities := []model.Activity{{Name: "Active Activity"}}
		mockRepo.On("GetActiveActivities", ctx).Return(activities, nil).Once()

		result, err := svc.GetActiveActivities(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetActiveActivities", ctx).Return([]model.Activity{}, errors.New("db error")).Once()

		_, err := svc.GetActiveActivities(ctx)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetVisibleActivities(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activities := []model.Activity{{Name: "Visible Activity"}}
		mockRepo.On("GetVisibleActivities", ctx).Return(activities, nil).Once()

		result, err := svc.GetVisibleActivities(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetVisibleActivities", ctx).Return([]model.Activity{}, errors.New("db error")).Once()

		_, err := svc.GetVisibleActivities(ctx)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetActivity(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activity := &model.Activity{Name: "Test Activity"}
		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()

		result, err := svc.GetActivity(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "Test Activity", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

		_, err := svc.GetActivity(ctx, 999)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CreateActivity(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activity := createValidActivity()
		mockRepo.On("CreateActivity", ctx, activity).Return(nil).Once()

		err := svc.CreateActivity(ctx, activity)
		require.NoError(t, err)
		assert.Equal(t, model.ActivityStatusDraft, activity.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error empty name", func(t *testing.T) {
		activity := createValidActivity()
		activity.Name = ""

		err := svc.CreateActivity(ctx, activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名称")
	})

	t.Run("validation error invalid time", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Name:    "Test",
			StartAt: now,
			EndAt:   now.Add(-1 * time.Hour),
		}

		err := svc.CreateActivity(ctx, activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "晚于")
	})

	t.Run("repo error", func(t *testing.T) {
		activity := createValidActivity()
		mockRepo.On("CreateActivity", ctx, activity).Return(errors.New("db error")).Once()

		err := svc.CreateActivity(ctx, activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "create activity")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateActivity(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		existing := createValidActivity()
		existing.Status = model.ActivityStatusDraft
		activity := createValidActivity()
		activity.ID = 1

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(existing, nil).Once()
		mockRepo.On("UpdateActivity", ctx, activity).Return(nil).Once()

		err := svc.UpdateActivity(ctx, activity)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ended activity cannot be updated", func(t *testing.T) {
		existing := createValidActivity()
		existing.Status = model.ActivityStatusEnded
		activity := createValidActivity()
		activity.ID = 1

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(existing, nil).Once()

		err := svc.UpdateActivity(ctx, activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "已结束")
		mockRepo.AssertExpectations(t)
	})

	t.Run("canceled activity cannot be updated", func(t *testing.T) {
		existing := createValidActivity()
		existing.Status = model.ActivityStatusCanceled
		activity := createValidActivity()
		activity.ID = 1

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(existing, nil).Once()

		err := svc.UpdateActivity(ctx, activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "已取消")
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		activity := createValidActivity()
		activity.ID = 999

		mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

		err := svc.UpdateActivity(ctx, activity)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeleteActivity(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activity := createValidActivity()
		activity.Status = model.ActivityStatusDraft

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		mockRepo.On("DeleteActivity", ctx, uint64(1)).Return(nil).Once()

		err := svc.DeleteActivity(ctx, 1)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("active activity cannot be deleted", func(t *testing.T) {
		activity := createValidActivity()
		activity.Status = model.ActivityStatusActive

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()

		err := svc.DeleteActivity(ctx, 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "进行中")
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

		err := svc.DeleteActivity(ctx, 999)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateActivityStatus(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("draft to active", func(t *testing.T) {
		activity := createValidActivity()
		activity.Status = model.ActivityStatusDraft

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		mockRepo.On("UpdateActivityStatus", ctx, uint64(1), model.ActivityStatusActive).Return(nil).Once()

		err := svc.UpdateActivityStatus(ctx, 1, model.ActivityStatusActive)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid transition", func(t *testing.T) {
		activity := createValidActivity()
		activity.Status = model.ActivityStatusEnded

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()

		err := svc.UpdateActivityStatus(ctx, 1, model.ActivityStatusActive)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不允许")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetActivityStats(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		stats := map[string]any{"total": int64(10)}
		mockRepo.On("GetActivityStats", ctx, uint64(1)).Return(stats, nil).Once()

		result, err := svc.GetActivityStats(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result["total"])
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Reward Management Tests
// ============================================================================

func TestService_GetReward(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 100}
		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(reward, nil).Once()

		result, err := svc.GetReward(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), result.CouponTemplateID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetRewardByID", ctx, uint64(999)).Return((*model.ActivityReward)(nil), repository.ErrNotFound).Once()

		_, err := svc.GetReward(ctx, 999)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetRewardsByActivityID(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rewards := []model.ActivityReward{{CouponTemplateID: 1}}
		mockRepo.On("GetRewardsByActivityID", ctx, uint64(1)).Return(rewards, nil).Once()

		result, err := svc.GetRewardsByActivityID(ctx, 1)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CreateReward(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activity := createValidActivity()
		reward := &model.ActivityReward{
			ActivityID:       1,
			CouponTemplateID: 1,
			CouponCount:      1,
			Probability:      100,
			TotalStock:       100,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		mockRepo.On("CreateReward", ctx, reward).Return(nil).Once()

		err := svc.CreateReward(ctx, reward)
		require.NoError(t, err)
		assert.Equal(t, 100, reward.RemainingStock)
		mockRepo.AssertExpectations(t)
	})

	t.Run("activity not found", func(t *testing.T) {
		reward := &model.ActivityReward{
			ActivityID:       999,
			CouponTemplateID: 1,
			CouponCount:      1,
			Probability:      100,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

		err := svc.CreateReward(ctx, reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		activity := createValidActivity()
		reward := &model.ActivityReward{
			ActivityID:       1,
			CouponTemplateID: 0,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()

		err := svc.CreateReward(ctx, reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "模板")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateReward(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success with stock adjustment", func(t *testing.T) {
		existing := &model.ActivityReward{
			TotalStock:     100,
			RemainingStock: 50,
		}
		reward := &model.ActivityReward{
			CouponTemplateID: 1,
			CouponCount:      1,
			Probability:      100,
			TotalStock:       150,
		}
		reward.ID = 1

		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(existing, nil).Once()
		mockRepo.On("UpdateReward", ctx, reward).Return(nil).Once()

		err := svc.UpdateReward(ctx, reward)
		require.NoError(t, err)
		assert.Equal(t, 100, reward.RemainingStock)
		mockRepo.AssertExpectations(t)
	})

	t.Run("stock decrease below remaining", func(t *testing.T) {
		existing := &model.ActivityReward{
			TotalStock:     100,
			RemainingStock: 50,
		}
		reward := &model.ActivityReward{
			CouponTemplateID: 1,
			CouponCount:      1,
			Probability:      100,
			TotalStock:       30,
		}
		reward.ID = 1

		mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(existing, nil).Once()
		mockRepo.On("UpdateReward", ctx, reward).Return(nil).Once()

		err := svc.UpdateReward(ctx, reward)
		require.NoError(t, err)
		assert.Equal(t, 0, reward.RemainingStock)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeleteReward(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteReward", ctx, uint64(1)).Return(nil).Once()

		err := svc.DeleteReward(ctx, 1)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Participation Tests
// ============================================================================

func TestService_ParticipateActivity(t *testing.T) {
	svc, mockRepo, mockCouponSvc := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
			// No limits set, so CountUserParticipations and CountTodayParticipations won't be called
		}
		reward := &model.ActivityReward{
			ActivityID:       1,
			CouponTemplateID: 1,
			CouponCount:      1,
			TotalStock:       100,
			RemainingStock:   50,
		}
		coupon := &model.Coupon{UserID: 100}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations or CountTodayParticipations calls since limits are 0
		mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()
		mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(1), model.CouponSourceActivity).Return(coupon, nil).Once()
		mockRepo.On("CreateParticipation", ctx, mock.AnythingOfType("*model.ActivityParticipation")).Return(nil).Once()
		mockRepo.On("IncrementParticipants", ctx, uint64(1)).Return(nil).Once()
		mockRepo.On("DecrementRewardStock", ctx, uint64(10)).Return(nil).Once()
		mockRepo.On("IncrementDailyStats", ctx, uint64(1)).Return(nil).Once()

		result, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.NoError(t, err)
		assert.Equal(t, uint64(100), result.UserID)
		mockRepo.AssertExpectations(t)
		mockCouponSvc.AssertExpectations(t)
	})

	t.Run("success with limits", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:       model.ActivityStatusActive,
			StartAt:      now.Add(-1 * time.Hour),
			EndAt:        now.Add(1 * time.Hour),
			PerUserLimit: 5,
			DailyLimit:   100,
		}
		reward := &model.ActivityReward{
			ActivityID:       1,
			CouponTemplateID: 1,
			CouponCount:      1,
			TotalStock:       100,
			RemainingStock:   50,
		}
		coupon := &model.Coupon{UserID: 100}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		mockRepo.On("CountUserParticipations", ctx, uint64(100), uint64(1)).Return(int64(0), nil).Once()
		mockRepo.On("CountTodayParticipations", ctx, uint64(1)).Return(int64(0), nil).Once()
		mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()
		mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(1), model.CouponSourceActivity).Return(coupon, nil).Once()
		mockRepo.On("CreateParticipation", ctx, mock.AnythingOfType("*model.ActivityParticipation")).Return(nil).Once()
		mockRepo.On("IncrementParticipants", ctx, uint64(1)).Return(nil).Once()
		mockRepo.On("DecrementRewardStock", ctx, uint64(10)).Return(nil).Once()
		mockRepo.On("IncrementDailyStats", ctx, uint64(1)).Return(nil).Once()

		result, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.NoError(t, err)
		assert.Equal(t, uint64(100), result.UserID)
		mockRepo.AssertExpectations(t)
		mockCouponSvc.AssertExpectations(t)
	})

	t.Run("activity not found", func(t *testing.T) {
		mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 999, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")
		mockRepo.AssertExpectations(t)
	})

	t.Run("activity not active", func(t *testing.T) {
		activity := &model.Activity{
			Status:  model.ActivityStatusDraft,
			StartAt: time.Now().Add(1 * time.Hour),
			EndAt:   time.Now().Add(2 * time.Hour),
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "未开始")
		mockRepo.AssertExpectations(t)
	})

	t.Run("per user limit reached", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:       model.ActivityStatusActive,
			StartAt:      now.Add(-1 * time.Hour),
			EndAt:        now.Add(1 * time.Hour),
			PerUserLimit: 3,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		mockRepo.On("CountUserParticipations", ctx, uint64(100), uint64(1)).Return(int64(3), nil).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "上限")
		mockRepo.AssertExpectations(t)
	})

	t.Run("daily limit reached", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:     model.ActivityStatusActive,
			StartAt:    now.Add(-1 * time.Hour),
			EndAt:      now.Add(1 * time.Hour),
			DailyLimit: 100,
			// No PerUserLimit, so CountUserParticipations won't be called
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations since PerUserLimit is 0
		mockRepo.On("CountTodayParticipations", ctx, uint64(1)).Return(int64(100), nil).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名额已满")
		mockRepo.AssertExpectations(t)
	})

	t.Run("total limit reached", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:       model.ActivityStatusActive,
			StartAt:      now.Add(-1 * time.Hour),
			EndAt:        now.Add(1 * time.Hour),
			TotalLimit:   1000,
			TotalClaimed: 1000,
			// No PerUserLimit or DailyLimit, so those checks won't be called
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations or CountTodayParticipations since limits are 0

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名额已满")
		mockRepo.AssertExpectations(t)
	})

	t.Run("reward not found", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
			// No limits set
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations or CountTodayParticipations since limits are 0
		mockRepo.On("GetRewardByID", ctx, uint64(999)).Return((*model.ActivityReward)(nil), repository.ErrNotFound).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 999, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不存在")
		mockRepo.AssertExpectations(t)
	})

	t.Run("reward not belong to activity", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
			// No limits set
		}
		reward := &model.ActivityReward{
			ActivityID: 2,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations or CountTodayParticipations since limits are 0
		mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "不属于")
		mockRepo.AssertExpectations(t)
	})

	t.Run("reward out of stock", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
			// No limits set
		}
		reward := &model.ActivityReward{
			ActivityID:     1,
			TotalStock:     100,
			RemainingStock: 0,
		}

		mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
		// No CountUserParticipations or CountTodayParticipations since limits are 0
		mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()

		_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "库存不足")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetUserParticipations(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success with default limit", func(t *testing.T) {
		participations := []model.ActivityParticipation{{UserID: 100}}
		mockRepo.On("GetUserParticipations", ctx, uint64(100), 20).Return(participations, nil).Once()

		result, err := svc.GetUserParticipations(ctx, 100, 0)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with custom limit", func(t *testing.T) {
		participations := []model.ActivityParticipation{{UserID: 100}}
		mockRepo.On("GetUserParticipations", ctx, uint64(100), 50).Return(participations, nil).Once()

		result, err := svc.GetUserParticipations(ctx, 100, 50)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_ListParticipations(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		opts := activityrepo.ParticipationListOptions{Page: 1, PageSize: 10}
		participations := []model.ActivityParticipation{{UserID: 100}}
		mockRepo.On("ListParticipations", ctx, opts).Return(participations, int64(1), nil).Once()

		result, total, err := svc.ListParticipations(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Statistics Tests
// ============================================================================

func TestService_GetAllActivityStats(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		stats := map[string]any{"total": int64(100)}
		mockRepo.On("GetAllActivityStats", ctx).Return(stats, nil).Once()

		result, err := svc.GetAllActivityStats(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(100), result["total"])
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Scheduled Tasks Tests
// ============================================================================

func TestService_ResetTodayParticipants(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("ResetTodayParticipants", ctx).Return(nil).Once()

		err := svc.ResetTodayParticipants(ctx)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("ResetTodayParticipants", ctx).Return(errors.New("db error")).Once()

		err := svc.ResetTodayParticipants(ctx)
		require.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_AutoUpdateActivityStatus(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	t.Run("success with status transitions", func(t *testing.T) {
		now := time.Now()
		preheatTime := now.Add(-2 * time.Hour)

		activity1 := model.Activity{
			Status:    model.ActivityStatusDraft,
			PreheatAt: &preheatTime,
			StartAt:   now.Add(1 * time.Hour),
			EndAt:     now.Add(2 * time.Hour),
		}
		activity1.ID = 1

		activity2 := model.Activity{
			Status:  model.ActivityStatusPreheat,
			StartAt: now.Add(-1 * time.Hour),
			EndAt:   now.Add(1 * time.Hour),
		}
		activity2.ID = 2

		activity3 := model.Activity{
			Status:  model.ActivityStatusActive,
			StartAt: now.Add(-2 * time.Hour),
			EndAt:   now.Add(-1 * time.Hour),
		}
		activity3.ID = 3

		activities := []model.Activity{activity1, activity2, activity3}

		opts := activityrepo.ActivityListOptions{Page: 1, PageSize: 1000}
		mockRepo.On("ListActivities", ctx, opts).Return(activities, int64(3), nil).Once()
		mockRepo.On("UpdateActivityStatus", ctx, uint64(1), model.ActivityStatusPreheat).Return(nil).Once()
		mockRepo.On("UpdateActivityStatus", ctx, uint64(2), model.ActivityStatusActive).Return(nil).Once()
		mockRepo.On("UpdateActivityStatus", ctx, uint64(3), model.ActivityStatusEnded).Return(nil).Once()

		err := svc.AutoUpdateActivityStatus(ctx)
		require.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestService_validateActivity(t *testing.T) {
	svc := &Service{}

	t.Run("empty name", func(t *testing.T) {
		activity := &model.Activity{Name: ""}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "名称")
	})

	t.Run("empty start time", func(t *testing.T) {
		activity := &model.Activity{Name: "Test"}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "开始时间")
	})

	t.Run("empty end time", func(t *testing.T) {
		activity := &model.Activity{
			Name:    "Test",
			StartAt: time.Now(),
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "结束时间")
	})

	t.Run("end before start", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Name:    "Test",
			StartAt: now,
			EndAt:   now.Add(-1 * time.Hour),
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "晚于")
	})

	t.Run("preheat after start", func(t *testing.T) {
		now := time.Now()
		preheat := now.Add(1 * time.Hour)
		activity := &model.Activity{
			Name:      "Test",
			StartAt:   now,
			EndAt:     now.Add(2 * time.Hour),
			PreheatAt: &preheat,
		}
		err := svc.validateActivity(activity)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "预热")
	})

	t.Run("valid activity", func(t *testing.T) {
		now := time.Now()
		activity := &model.Activity{
			Name:    "Test",
			StartAt: now,
			EndAt:   now.Add(2 * time.Hour),
		}
		err := svc.validateActivity(activity)
		require.NoError(t, err)
	})
}

func TestService_validateStatusTransition(t *testing.T) {
	svc := &Service{}

	t.Run("draft to active", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusActive)
		require.NoError(t, err)
	})

	t.Run("draft to preheat", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusDraft, model.ActivityStatusPreheat)
		require.NoError(t, err)
	})

	t.Run("active to paused", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusActive, model.ActivityStatusPaused)
		require.NoError(t, err)
	})

	t.Run("paused to active", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusPaused, model.ActivityStatusActive)
		require.NoError(t, err)
	})

	t.Run("ended cannot transition", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusEnded, model.ActivityStatusActive)
		require.Error(t, err)
	})

	t.Run("canceled cannot transition", func(t *testing.T) {
		err := svc.validateStatusTransition(model.ActivityStatusCanceled, model.ActivityStatusActive)
		require.Error(t, err)
	})
}

func TestService_validateReward(t *testing.T) {
	svc := &Service{}

	t.Run("empty template id", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "模板")
	})

	t.Run("invalid coupon count", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "数量")
	})

	t.Run("invalid probability low", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 0}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "概率")
	})

	t.Run("invalid probability high", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 101}
		err := svc.validateReward(reward)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "概率")
	})

	t.Run("valid reward", func(t *testing.T) {
		reward := &model.ActivityReward{CouponTemplateID: 1, CouponCount: 1, Probability: 100}
		err := svc.validateReward(reward)
		require.NoError(t, err)
	})
}

// ============================================================================
// Additional Tests for Coverage Improvement
// ============================================================================

func TestNewActivityService(t *testing.T) {
	mockRepo := new(MockActivityRepository)
	mockCouponSvc := new(MockCouponService)

	svc := NewActivityService(mockRepo, mockCouponSvc)
	require.NotNil(t, svc)
	assert.Equal(t, mockRepo, svc.repo)
	assert.Equal(t, mockCouponSvc, svc.couponSvc)
}

func TestService_GetRewardsByActivityID_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("GetRewardsByActivityID", ctx, uint64(1)).Return([]model.ActivityReward{}, errors.New("db error")).Once()

	_, err := svc.GetRewardsByActivityID(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get rewards")
	mockRepo.AssertExpectations(t)
}

func TestService_CreateReward_GetActivityError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	reward := &model.ActivityReward{
		ActivityID:       1,
		CouponTemplateID: 1,
		CouponCount:      1,
		Probability:      100,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return((*model.Activity)(nil), errors.New("db error")).Once()

	err := svc.CreateReward(ctx, reward)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get activity")
	mockRepo.AssertExpectations(t)
}

func TestService_CreateReward_RepoError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	activity := createValidActivity()
	reward := &model.ActivityReward{
		ActivityID:       1,
		CouponTemplateID: 1,
		CouponCount:      1,
		Probability:      100,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("CreateReward", ctx, reward).Return(errors.New("db error")).Once()

	err := svc.CreateReward(ctx, reward)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create reward")
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateReward_NotFound(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	reward := &model.ActivityReward{
		CouponTemplateID: 1,
		CouponCount:      1,
		Probability:      100,
	}
	reward.ID = 999

	mockRepo.On("GetRewardByID", ctx, uint64(999)).Return((*model.ActivityReward)(nil), repository.ErrNotFound).Once()

	err := svc.UpdateReward(ctx, reward)
	require.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateReward_ValidationError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	existing := &model.ActivityReward{
		TotalStock:     100,
		RemainingStock: 50,
	}
	reward := &model.ActivityReward{
		CouponTemplateID: 0, // Invalid
		CouponCount:      1,
		Probability:      100,
	}
	reward.ID = 1

	mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(existing, nil).Once()

	err := svc.UpdateReward(ctx, reward)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "模板")
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateReward_RepoError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	existing := &model.ActivityReward{
		TotalStock:     100,
		RemainingStock: 50,
	}
	reward := &model.ActivityReward{
		CouponTemplateID: 1,
		CouponCount:      1,
		Probability:      100,
		TotalStock:       100,
	}
	reward.ID = 1

	mockRepo.On("GetRewardByID", ctx, uint64(1)).Return(existing, nil).Once()
	mockRepo.On("UpdateReward", ctx, reward).Return(errors.New("db error")).Once()

	err := svc.UpdateReward(ctx, reward)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "update reward")
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteReward_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("DeleteReward", ctx, uint64(1)).Return(errors.New("db error")).Once()

	err := svc.DeleteReward(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete reward")
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateActivityStatus_RepoError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	activity := createValidActivity()
	activity.Status = model.ActivityStatusDraft

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("UpdateActivityStatus", ctx, uint64(1), model.ActivityStatusActive).Return(errors.New("db error")).Once()

	err := svc.UpdateActivityStatus(ctx, 1, model.ActivityStatusActive)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "update status")
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateActivityStatus_NotFound(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("GetActivityByID", ctx, uint64(999)).Return((*model.Activity)(nil), repository.ErrNotFound).Once()

	err := svc.UpdateActivityStatus(ctx, 999, model.ActivityStatusActive)
	require.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_GetActivityStats_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("GetActivityStats", ctx, uint64(1)).Return(map[string]any{}, errors.New("db error")).Once()

	_, err := svc.GetActivityStats(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get activity stats")
	mockRepo.AssertExpectations(t)
}

func TestService_GetAllActivityStats_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("GetAllActivityStats", ctx).Return(map[string]any{}, errors.New("db error")).Once()

	_, err := svc.GetAllActivityStats(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get all activity stats")
	mockRepo.AssertExpectations(t)
}

func TestService_ListParticipations_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	opts := activityrepo.ParticipationListOptions{Page: 1, PageSize: 10}
	mockRepo.On("ListParticipations", ctx, opts).Return([]model.ActivityParticipation{}, int64(0), errors.New("db error")).Once()

	_, _, err := svc.ListParticipations(ctx, opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "list participations")
	mockRepo.AssertExpectations(t)
}

func TestService_GetUserParticipations_Error(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	mockRepo.On("GetUserParticipations", ctx, uint64(100), 20).Return([]model.ActivityParticipation{}, errors.New("db error")).Once()

	_, err := svc.GetUserParticipations(ctx, 100, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get user participations")
	mockRepo.AssertExpectations(t)
}

func TestService_ParticipateActivity_CountUserParticipationsError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:       model.ActivityStatusActive,
		StartAt:      now.Add(-1 * time.Hour),
		EndAt:        now.Add(1 * time.Hour),
		PerUserLimit: 5,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("CountUserParticipations", ctx, uint64(100), uint64(1)).Return(int64(0), errors.New("db error")).Once()

	_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "count participations")
	mockRepo.AssertExpectations(t)
}

func TestService_ParticipateActivity_CountTodayParticipationsError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:     model.ActivityStatusActive,
		StartAt:    now.Add(-1 * time.Hour),
		EndAt:      now.Add(1 * time.Hour),
		DailyLimit: 100,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("CountTodayParticipations", ctx, uint64(1)).Return(int64(0), errors.New("db error")).Once()

	_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "count today participations")
	mockRepo.AssertExpectations(t)
}

func TestService_ParticipateActivity_GetRewardError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-1 * time.Hour),
		EndAt:   now.Add(1 * time.Hour),
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("GetRewardByID", ctx, uint64(10)).Return((*model.ActivityReward)(nil), errors.New("db error")).Once()

	_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "get reward")
	mockRepo.AssertExpectations(t)
}

func TestService_ParticipateActivity_CreateParticipationError(t *testing.T) {
	svc, mockRepo, mockCouponSvc := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-1 * time.Hour),
		EndAt:   now.Add(1 * time.Hour),
	}
	reward := &model.ActivityReward{
		ActivityID:       1,
		CouponTemplateID: 1,
		CouponCount:      1,
		TotalStock:       100,
		RemainingStock:   50,
	}
	coupon := &model.Coupon{UserID: 100}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()
	mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(1), model.CouponSourceActivity).Return(coupon, nil).Once()
	mockRepo.On("CreateParticipation", ctx, mock.AnythingOfType("*model.ActivityParticipation")).Return(errors.New("db error")).Once()

	_, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "create participation")
	mockRepo.AssertExpectations(t)
	mockCouponSvc.AssertExpectations(t)
}

func TestService_ParticipateActivity_CouponIssueFails(t *testing.T) {
	svc, mockRepo, mockCouponSvc := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-1 * time.Hour),
		EndAt:   now.Add(1 * time.Hour),
	}
	reward := &model.ActivityReward{
		ActivityID:       1,
		CouponTemplateID: 1,
		CouponCount:      2,
		TotalStock:       100,
		RemainingStock:   50,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()
	// First coupon fails, second succeeds
	mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(1), model.CouponSourceActivity).Return((*model.Coupon)(nil), errors.New("issue error")).Once()
	mockCouponSvc.On("IssueCoupon", ctx, uint64(100), uint64(1), model.CouponSourceActivity).Return(&model.Coupon{UserID: 100}, nil).Once()
	mockRepo.On("CreateParticipation", ctx, mock.AnythingOfType("*model.ActivityParticipation")).Return(nil).Once()
	mockRepo.On("IncrementParticipants", ctx, uint64(1)).Return(nil).Once()
	mockRepo.On("DecrementRewardStock", ctx, uint64(10)).Return(nil).Once()
	mockRepo.On("IncrementDailyStats", ctx, uint64(1)).Return(nil).Once()

	result, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
	mockCouponSvc.AssertExpectations(t)
}

func TestService_ParticipateActivity_NilCouponService(t *testing.T) {
	mockRepo := new(MockActivityRepository)
	svc := &Service{repo: mockRepo, couponSvc: nil}
	ctx := context.Background()

	now := time.Now()
	activity := &model.Activity{
		Status:  model.ActivityStatusActive,
		StartAt: now.Add(-1 * time.Hour),
		EndAt:   now.Add(1 * time.Hour),
	}
	reward := &model.ActivityReward{
		ActivityID:       1,
		CouponTemplateID: 1,
		CouponCount:      1,
		TotalStock:       100,
		RemainingStock:   50,
	}

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("GetRewardByID", ctx, uint64(10)).Return(reward, nil).Once()
	mockRepo.On("CreateParticipation", ctx, mock.AnythingOfType("*model.ActivityParticipation")).Return(nil).Once()
	mockRepo.On("IncrementParticipants", ctx, uint64(1)).Return(nil).Once()
	mockRepo.On("DecrementRewardStock", ctx, uint64(10)).Return(nil).Once()
	mockRepo.On("IncrementDailyStats", ctx, uint64(1)).Return(nil).Once()

	result, err := svc.ParticipateActivity(ctx, 100, 1, 10, "192.168.1.1")
	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_AutoUpdateActivityStatus_ListError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	opts := activityrepo.ActivityListOptions{Page: 1, PageSize: 1000}
	mockRepo.On("ListActivities", ctx, opts).Return([]model.Activity{}, int64(0), errors.New("db error")).Once()

	err := svc.AutoUpdateActivityStatus(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "list activities")
	mockRepo.AssertExpectations(t)
}

func TestService_AutoUpdateActivityStatus_DraftToActive(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	now := time.Now()
	// Draft activity that should go directly to active (no preheat, past start time)
	activity := model.Activity{
		Status:  model.ActivityStatusDraft,
		StartAt: now.Add(-1 * time.Hour),
		EndAt:   now.Add(1 * time.Hour),
	}
	activity.ID = 1

	activities := []model.Activity{activity}

	opts := activityrepo.ActivityListOptions{Page: 1, PageSize: 1000}
	mockRepo.On("ListActivities", ctx, opts).Return(activities, int64(1), nil).Once()
	mockRepo.On("UpdateActivityStatus", ctx, uint64(1), model.ActivityStatusActive).Return(nil).Once()

	err := svc.AutoUpdateActivityStatus(ctx)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_validateStatusTransition_InvalidFromStatus(t *testing.T) {
	svc := &Service{}

	err := svc.validateStatusTransition(model.ActivityStatus("invalid"), model.ActivityStatusActive)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "无效的当前状态")
}

func TestService_UpdateActivity_RepoError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	existing := createValidActivity()
	existing.Status = model.ActivityStatusDraft
	activity := createValidActivity()
	activity.ID = 1

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(existing, nil).Once()
	mockRepo.On("UpdateActivity", ctx, activity).Return(errors.New("db error")).Once()

	err := svc.UpdateActivity(ctx, activity)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "update activity")
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteActivity_RepoError(t *testing.T) {
	svc, mockRepo, _ := createTestActivityService()
	ctx := context.Background()

	activity := createValidActivity()
	activity.Status = model.ActivityStatusDraft

	mockRepo.On("GetActivityByID", ctx, uint64(1)).Return(activity, nil).Once()
	mockRepo.On("DeleteActivity", ctx, uint64(1)).Return(errors.New("db error")).Once()

	err := svc.DeleteActivity(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "delete activity")
	mockRepo.AssertExpectations(t)
}

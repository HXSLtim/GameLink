package userblock

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ============================================================================
// Mock Repositories
// ============================================================================

// MockUserBlockRepository is a mock implementation of UserBlockRepository
type MockUserBlockRepository struct {
	mock.Mock
}

func (m *MockUserBlockRepository) Create(ctx context.Context, block *model.UserBlock) error {
	args := m.Called(ctx, block)
	return args.Error(0)
}

func (m *MockUserBlockRepository) Get(ctx context.Context, id uint64) (*model.UserBlock, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) GetWithRelations(ctx context.Context, id uint64) (*model.UserBlock, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) GetByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error) {
	args := m.Called(ctx, blockerID, blockedID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) GetActiveByBlockerAndBlocked(ctx context.Context, blockerID, blockedID uint64) (*model.UserBlock, error) {
	args := m.Called(ctx, blockerID, blockedID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) IsBlocked(ctx context.Context, userID1, userID2 uint64) (bool, error) {
	args := m.Called(ctx, userID1, userID2)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserBlockRepository) IsBlockedBy(ctx context.Context, blockerID, blockedID uint64) (bool, error) {
	args := m.Called(ctx, blockerID, blockedID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserBlockRepository) ListByBlockerID(ctx context.Context, blockerID uint64, status *model.BlockStatus) ([]model.UserBlock, error) {
	args := m.Called(ctx, blockerID, status)
	return args.Get(0).([]model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) ListByBlockedID(ctx context.Context, blockedID uint64, status *model.BlockStatus) ([]model.UserBlock, error) {
	args := m.Called(ctx, blockedID, status)
	return args.Get(0).([]model.UserBlock), args.Error(1)
}

func (m *MockUserBlockRepository) ListPaged(ctx context.Context, opts repository.UserBlockListOptions) ([]model.UserBlock, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.UserBlock), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserBlockRepository) Update(ctx context.Context, block *model.UserBlock) error {
	args := m.Called(ctx, block)
	return args.Error(0)
}

func (m *MockUserBlockRepository) UpdateStatus(ctx context.Context, id uint64, status model.BlockStatus, canceledBy *uint64, adminRemark string) error {
	args := m.Called(ctx, id, status, canceledBy, adminRemark)
	return args.Error(0)
}

func (m *MockUserBlockRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserBlockRepository) GetBlockedUserIDs(ctx context.Context, blockerID uint64) ([]uint64, error) {
	args := m.Called(ctx, blockerID)
	return args.Get(0).([]uint64), args.Error(1)
}

func (m *MockUserBlockRepository) GetBlockerUserIDs(ctx context.Context, blockedID uint64) ([]uint64, error) {
	args := m.Called(ctx, blockedID)
	return args.Get(0).([]uint64), args.Error(1)
}

func (m *MockUserBlockRepository) GetAllBlockRelatedUserIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]uint64), args.Error(1)
}

func (m *MockUserBlockRepository) CountByStatus(ctx context.Context) (map[model.BlockStatus]int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[model.BlockStatus]int64), args.Error(1)
}

func (m *MockUserBlockRepository) GetActiveCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) List(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	args := m.Called(ctx, opts)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ============================================================================
// Tests
// ============================================================================

func TestNewUserBlockService(t *testing.T) {
	blockRepo := &MockUserBlockRepository{}
	userRepo := &MockUserRepository{}
	svc := NewUserBlockService(blockRepo, userRepo)
	require.NotNil(t, svc)
	assert.Equal(t, blockRepo, svc.repo)
	assert.Equal(t, userRepo, svc.users)
}

func TestService_Block(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		input := BlockInput{
			BlockerID:   100,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   200,
			BlockedType: model.BlockUserTypePlayer,
			Reason:      "test reason",
		}

		blockedUser := &model.User{Name: "Blocked User"}
		blockedUser.ID = 200

		userRepo.On("Get", ctx, uint64(200)).Return(blockedUser, nil)
		blockRepo.On("GetActiveByBlockerAndBlocked", ctx, uint64(100), uint64(200)).Return(nil, repository.ErrNotFound)
		blockRepo.On("Create", ctx, mock.AnythingOfType("*model.UserBlock")).Return(nil)

		result, err := svc.Block(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, uint64(100), result.BlockerID)
		assert.Equal(t, uint64(200), result.BlockedID)
		assert.Equal(t, model.BlockStatusActive, result.Status)
		blockRepo.AssertExpectations(t)
		userRepo.AssertExpectations(t)
	})

	t.Run("cannot block self", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		input := BlockInput{
			BlockerID: 100,
			BlockedID: 100,
		}

		_, err := svc.Block(ctx, input)
		require.Error(t, err)
		assert.Equal(t, ErrCannotBlockSelf, err)
	})

	t.Run("blocked user not found", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		input := BlockInput{
			BlockerID: 100,
			BlockedID: 999,
		}

		userRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		_, err := svc.Block(ctx, input)
		require.Error(t, err)
		assert.Equal(t, ErrUserNotFound, err)
	})

	t.Run("already blocked", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		input := BlockInput{
			BlockerID: 100,
			BlockedID: 200,
		}

		blockedUser := &model.User{Name: "Blocked User"}
		blockedUser.ID = 200
		existingBlock := &model.UserBlock{BlockerID: 100, BlockedID: 200, Status: model.BlockStatusActive}

		userRepo.On("Get", ctx, uint64(200)).Return(blockedUser, nil)
		blockRepo.On("GetActiveByBlockerAndBlocked", ctx, uint64(100), uint64(200)).Return(existingBlock, nil)

		_, err := svc.Block(ctx, input)
		require.Error(t, err)
		assert.Equal(t, ErrAlreadyBlocked, err)
	})
}

func TestService_Unblock(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		existingBlock := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Status:    model.BlockStatusActive,
		}
		existingBlock.ID = 1

		blockRepo.On("GetActiveByBlockerAndBlocked", ctx, uint64(100), uint64(200)).Return(existingBlock, nil)
		blockRepo.On("Update", ctx, mock.AnythingOfType("*model.UserBlock")).Return(nil)

		err := svc.Unblock(ctx, 100, 200)
		require.NoError(t, err)
		blockRepo.AssertExpectations(t)
	})

	t.Run("not blocked", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("GetActiveByBlockerAndBlocked", ctx, uint64(100), uint64(200)).Return(nil, repository.ErrNotFound)

		err := svc.Unblock(ctx, 100, 200)
		require.Error(t, err)
		assert.Equal(t, ErrNotBlocked, err)
	})
}

func TestService_AdminUnblock(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		existingBlock := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Status:    model.BlockStatusActive,
		}
		existingBlock.ID = 1

		blockRepo.On("Get", ctx, uint64(1)).Return(existingBlock, nil)
		blockRepo.On("Update", ctx, mock.AnythingOfType("*model.UserBlock")).Return(nil)

		err := svc.AdminUnblock(ctx, 1, 999, "admin unblock reason")
		require.NoError(t, err)
		blockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("Get", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		err := svc.AdminUnblock(ctx, 999, 1, "reason")
		require.Error(t, err)
	})

	t.Run("not active", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		existingBlock := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Status:    model.BlockStatusCanceled,
		}
		existingBlock.ID = 1

		blockRepo.On("Get", ctx, uint64(1)).Return(existingBlock, nil)

		err := svc.AdminUnblock(ctx, 1, 999, "reason")
		require.Error(t, err)
		assert.Equal(t, ErrNotBlocked, err)
	})
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		expected := &model.UserBlock{
			BlockerID: 100,
			BlockedID: 200,
			Status:    model.BlockStatusActive,
		}
		expected.ID = 1

		blockRepo.On("GetWithRelations", ctx, uint64(1)).Return(expected, nil)

		result, err := svc.Get(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("not found", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("GetWithRelations", ctx, uint64(999)).Return(nil, repository.ErrNotFound)

		_, err := svc.Get(ctx, 999)
		require.Error(t, err)
	})
}

func TestService_IsBlocked(t *testing.T) {
	ctx := context.Background()

	t.Run("blocked", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("IsBlocked", ctx, uint64(100), uint64(200)).Return(true, nil)

		result, err := svc.IsBlocked(ctx, 100, 200)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("not blocked", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("IsBlocked", ctx, uint64(100), uint64(200)).Return(false, nil)

		result, err := svc.IsBlocked(ctx, 100, 200)
		require.NoError(t, err)
		assert.False(t, result)
	})
}

func TestService_IsBlockedBy(t *testing.T) {
	ctx := context.Background()

	t.Run("blocked by", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("IsBlockedBy", ctx, uint64(100), uint64(200)).Return(true, nil)

		result, err := svc.IsBlockedBy(ctx, 100, 200)
		require.NoError(t, err)
		assert.True(t, result)
	})
}

func TestService_ListByBlocker(t *testing.T) {
	ctx := context.Background()

	t.Run("active only", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blocks := []model.UserBlock{
			{BlockerID: 100, BlockedID: 200, Status: model.BlockStatusActive},
			{BlockerID: 100, BlockedID: 300, Status: model.BlockStatusActive},
		}
		activeStatus := model.BlockStatusActive

		blockRepo.On("ListByBlockerID", ctx, uint64(100), &activeStatus).Return(blocks, nil)

		result, err := svc.ListByBlocker(ctx, 100, true)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("all statuses", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blocks := []model.UserBlock{
			{BlockerID: 100, BlockedID: 200, Status: model.BlockStatusActive},
			{BlockerID: 100, BlockedID: 300, Status: model.BlockStatusCanceled},
		}

		blockRepo.On("ListByBlockerID", ctx, uint64(100), (*model.BlockStatus)(nil)).Return(blocks, nil)

		result, err := svc.ListByBlocker(ctx, 100, false)
		require.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

func TestService_ListByBlocked(t *testing.T) {
	ctx := context.Background()

	t.Run("active only", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blocks := []model.UserBlock{
			{BlockerID: 100, BlockedID: 200, Status: model.BlockStatusActive},
		}
		activeStatus := model.BlockStatusActive

		blockRepo.On("ListByBlockedID", ctx, uint64(200), &activeStatus).Return(blocks, nil)

		result, err := svc.ListByBlocked(ctx, 200, true)
		require.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

func TestService_ListPaged(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blocks := []model.UserBlock{
			{BlockerID: 100, BlockedID: 200},
			{BlockerID: 100, BlockedID: 300},
		}
		opts := repository.UserBlockListOptions{Page: 1, PageSize: 10}

		blockRepo.On("ListPaged", ctx, opts).Return(blocks, int64(2), nil)

		result, pagination, err := svc.ListPaged(ctx, opts)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, pagination.Total)
		assert.Equal(t, 1, pagination.Page)
		assert.Equal(t, 10, pagination.PageSize)
	})
}

func TestService_GetBlockedUserIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		ids := []uint64{200, 300, 400}
		blockRepo.On("GetBlockedUserIDs", ctx, uint64(100)).Return(ids, nil)

		result, err := svc.GetBlockedUserIDs(ctx, 100)
		require.NoError(t, err)
		assert.Equal(t, ids, result)
	})
}

func TestService_GetBlockerUserIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		ids := []uint64{100, 300}
		blockRepo.On("GetBlockerUserIDs", ctx, uint64(200)).Return(ids, nil)

		result, err := svc.GetBlockerUserIDs(ctx, 200)
		require.NoError(t, err)
		assert.Equal(t, ids, result)
	})
}

func TestService_GetAllBlockRelatedUserIDs(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		ids := []uint64{200, 300, 400, 500}
		blockRepo.On("GetAllBlockRelatedUserIDs", ctx, uint64(100)).Return(ids, nil)

		result, err := svc.GetAllBlockRelatedUserIDs(ctx, 100)
		require.NoError(t, err)
		assert.Equal(t, ids, result)
	})
}

func TestService_GetStats(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		stats := map[model.BlockStatus]int64{
			model.BlockStatusActive:        10,
			model.BlockStatusCanceled:      5,
			model.BlockStatusAdminCanceled: 2,
		}
		blockRepo.On("CountByStatus", ctx).Return(stats, nil)

		result, err := svc.GetStats(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(10), result[model.BlockStatusActive])
		assert.Equal(t, int64(5), result[model.BlockStatusCanceled])
	})
}

func TestService_GetActiveCount(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("GetActiveCount", ctx).Return(int64(15), nil)

		result, err := svc.GetActiveCount(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(15), result)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		blockRepo.On("Delete", ctx, uint64(1)).Return(nil)

		err := svc.Delete(ctx, 1)
		require.NoError(t, err)
	})
}

func TestService_BatchUnblock(t *testing.T) {
	ctx := context.Background()

	t.Run("all success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		block1 := &model.UserBlock{Status: model.BlockStatusActive}
		block1.ID = 1
		block2 := &model.UserBlock{Status: model.BlockStatusActive}
		block2.ID = 2

		blockRepo.On("Get", ctx, uint64(1)).Return(block1, nil)
		blockRepo.On("Get", ctx, uint64(2)).Return(block2, nil)
		blockRepo.On("Update", ctx, mock.AnythingOfType("*model.UserBlock")).Return(nil).Times(2)

		count, err := svc.BatchUnblock(ctx, []uint64{1, 2}, 999, "batch unblock")
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("partial success", func(t *testing.T) {
		blockRepo := &MockUserBlockRepository{}
		userRepo := &MockUserRepository{}
		svc := NewUserBlockService(blockRepo, userRepo)

		block1 := &model.UserBlock{Status: model.BlockStatusActive}
		block1.ID = 1

		blockRepo.On("Get", ctx, uint64(1)).Return(block1, nil)
		blockRepo.On("Get", ctx, uint64(2)).Return(nil, repository.ErrNotFound)
		blockRepo.On("Update", ctx, mock.AnythingOfType("*model.UserBlock")).Return(nil)

		count, err := svc.BatchUnblock(ctx, []uint64{1, 2}, 999, "batch unblock")
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

// ============================================================================
// Model Tests
// ============================================================================

func TestBlockStatus_Constants(t *testing.T) {
	t.Run("status values", func(t *testing.T) {
		assert.Equal(t, model.BlockStatus("active"), model.BlockStatusActive)
		assert.Equal(t, model.BlockStatus("canceled"), model.BlockStatusCanceled)
		assert.Equal(t, model.BlockStatus("admin_canceled"), model.BlockStatusAdminCanceled)
	})
}

func TestBlockUserType_Constants(t *testing.T) {
	t.Run("type values", func(t *testing.T) {
		assert.Equal(t, model.BlockUserType("user"), model.BlockUserTypeUser)
		assert.Equal(t, model.BlockUserType("player"), model.BlockUserTypePlayer)
	})
}

func TestUserBlock_Cancel(t *testing.T) {
	t.Run("user cancel", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusActive,
		}
		block.Cancel()

		assert.Equal(t, model.BlockStatusCanceled, block.Status)
		assert.NotNil(t, block.CanceledAt)
	})
}

func TestUserBlock_AdminCancel(t *testing.T) {
	t.Run("admin cancel", func(t *testing.T) {
		block := &model.UserBlock{
			Status: model.BlockStatusActive,
		}
		adminID := uint64(1)
		remark := "admin cancel reason"
		block.AdminCancel(adminID, remark)

		assert.Equal(t, model.BlockStatusAdminCanceled, block.Status)
		assert.NotNil(t, block.CanceledAt)
		assert.NotNil(t, block.CanceledBy)
		assert.Equal(t, adminID, *block.CanceledBy)
		assert.Equal(t, remark, block.AdminRemark)
	})
}

func TestUserBlock_IsActive(t *testing.T) {
	t.Run("active block", func(t *testing.T) {
		block := &model.UserBlock{Status: model.BlockStatusActive}
		assert.True(t, block.IsActive())
	})

	t.Run("canceled block", func(t *testing.T) {
		block := &model.UserBlock{Status: model.BlockStatusCanceled}
		assert.False(t, block.IsActive())
	})

	t.Run("admin canceled block", func(t *testing.T) {
		block := &model.UserBlock{Status: model.BlockStatusAdminCanceled}
		assert.False(t, block.IsActive())
	})
}

func TestUserBlock_Structure(t *testing.T) {
	t.Run("block record", func(t *testing.T) {
		now := time.Now()
		block := &model.UserBlock{
			BlockerID:   100,
			BlockerType: model.BlockUserTypeUser,
			BlockedID:   200,
			BlockedType: model.BlockUserTypePlayer,
			Reason:      "test reason",
			Status:      model.BlockStatusActive,
			BlockedAt:   now,
		}

		assert.Equal(t, uint64(100), block.BlockerID)
		assert.Equal(t, uint64(200), block.BlockedID)
		assert.Equal(t, model.BlockUserTypeUser, block.BlockerType)
		assert.Equal(t, model.BlockUserTypePlayer, block.BlockedType)
		assert.Equal(t, model.BlockStatusActive, block.Status)
	})
}

func TestService_Errors(t *testing.T) {
	t.Run("error types", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.NotNil(t, ErrValidation)
		assert.NotNil(t, ErrAlreadyBlocked)
		assert.NotNil(t, ErrNotBlocked)
		assert.NotNil(t, ErrCannotBlockSelf)
		assert.NotNil(t, ErrUserNotFound)
	})
}

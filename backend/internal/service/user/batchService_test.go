package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	userrepo "gamelink/internal/repository/user"
	"gamelink/pkg/testutil"
)

// mockNotificationRepo is a mock implementation of NotificationRepository
type mockNotificationRepo struct{}

func (m *mockNotificationRepo) ListByUser(ctx context.Context, opts repository.NotificationListOptions) ([]model.NotificationEvent, int64, error) {
	return nil, 0, nil
}

func (m *mockNotificationRepo) MarkRead(ctx context.Context, userID uint64, ids []uint64) error {
	return nil
}

func (m *mockNotificationRepo) MarkAllRead(ctx context.Context, userID uint64) error {
	return nil
}

func (m *mockNotificationRepo) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

func (m *mockNotificationRepo) Create(ctx context.Context, event *model.NotificationEvent) error {
	return nil
}

func (m *mockNotificationRepo) Delete(ctx context.Context, userID uint64, id uint64) error {
	return nil
}

func setupBatchServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.UserTag{},
		&model.Wallet{},
		&model.NotificationEvent{},
	)
	return db
}

func createBatchService(db *gorm.DB) *BatchOperationService {
	userRepo := userrepo.NewUserRepository(db)
	tagRepo := userrepo.NewUserTagRepository(db)
	notificationRepo := &mockNotificationRepo{}
	return NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)
}

func createTestUsersForBatch(t *testing.T, db *gorm.DB, count int) []*model.User {
	t.Helper()
	users := make([]*model.User, count)
	for i := 0; i < count; i++ {
		users[i] = &model.User{
			Phone:        fmt.Sprintf("1380000%04d", i),
			Email:        fmt.Sprintf("batchuser%d@test.com", i),
			Name:         fmt.Sprintf("BatchUser%d", i),
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
			PasswordHash: "hashed",
		}
		require.NoError(t, db.Create(users[i]).Error)
	}
	return users
}

func TestBatchOperationService_BatchUpdateUserRole(t *testing.T) {
	db := setupBatchServiceDB(t)
	defer testutil.CleanDB(t, db)

	svc := createBatchService(db)
	ctx := context.Background()

	// Create test users
	users := createTestUsersForBatch(t, db, 3)
	userIDs := make([]uint64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	t.Run("batch update user role success", func(t *testing.T) {
		success, failed, err := svc.BatchUpdateUserRole(ctx, &BatchUpdateUserRoleRequest{
			UserIDs: userIDs,
			Role:    "player",
		})
		require.NoError(t, err)
		assert.Equal(t, len(userIDs), success)
		assert.Equal(t, 0, failed)

		// Verify role changed
		var user model.User
		require.NoError(t, db.First(&user, users[0].ID).Error)
		assert.Equal(t, model.Role("player"), user.Role)
	})

	t.Run("empty user IDs", func(t *testing.T) {
		_, _, err := svc.BatchUpdateUserRole(ctx, &BatchUpdateUserRoleRequest{
			UserIDs: []uint64{},
			Role:    "player",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "用户ID列表不能为空")
	})

	t.Run("too many users", func(t *testing.T) {
		largeUserIDs := make([]uint64, 1001)
		for i := range largeUserIDs {
			largeUserIDs[i] = uint64(i + 1)
		}
		_, _, err := svc.BatchUpdateUserRole(ctx, &BatchUpdateUserRoleRequest{
			UserIDs: largeUserIDs,
			Role:    "player",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "一次批量操作最多支持1000个用户")
	})
}

func TestBatchOperationService_BatchUpdateUserStatus(t *testing.T) {
	db := setupBatchServiceDB(t)
	defer testutil.CleanDB(t, db)

	svc := createBatchService(db)
	ctx := context.Background()

	// Create test users
	users := createTestUsersForBatch(t, db, 3)
	userIDs := make([]uint64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	t.Run("batch update user status success", func(t *testing.T) {
		success, failed, err := svc.BatchUpdateUserStatus(ctx, &BatchUpdateUserStatusRequest{
			UserIDs: userIDs,
			Status:  "banned",
		}, 1)
		require.NoError(t, err)
		assert.Equal(t, len(userIDs), success)
		assert.Equal(t, 0, failed)

		// Verify status changed
		var user model.User
		require.NoError(t, db.First(&user, users[0].ID).Error)
		assert.Equal(t, model.UserStatus("banned"), user.Status)
	})

	t.Run("too many users", func(t *testing.T) {
		largeUserIDs := make([]uint64, 1001)
		for i := range largeUserIDs {
			largeUserIDs[i] = uint64(i + 1)
		}
		_, _, err := svc.BatchUpdateUserStatus(ctx, &BatchUpdateUserStatusRequest{
			UserIDs: largeUserIDs,
			Status:  "active",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "一次批量操作最多支持1000个用户")
	})
}

func TestBatchOperationService_BatchDeleteUsers(t *testing.T) {
	db := setupBatchServiceDB(t)
	defer testutil.CleanDB(t, db)

	svc := createBatchService(db)
	ctx := context.Background()

	// Create test users
	users := createTestUsersForBatch(t, db, 3)
	userIDs := make([]uint64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	t.Run("batch delete users success", func(t *testing.T) {
		success, failed, err := svc.BatchDeleteUsers(ctx, &BatchDeleteUsersRequest{
			UserIDs: userIDs,
			Reason:  "测试删除",
		}, 1)
		require.NoError(t, err)
		assert.Equal(t, len(userIDs), success)
		assert.Equal(t, 0, failed)

		// Verify users are soft deleted
		var count int64
		db.Model(&model.User{}).Where("id IN ?", userIDs).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("too many users", func(t *testing.T) {
		largeUserIDs := make([]uint64, 1001)
		for i := range largeUserIDs {
			largeUserIDs[i] = uint64(i + 1)
		}
		_, _, err := svc.BatchDeleteUsers(ctx, &BatchDeleteUsersRequest{
			UserIDs: largeUserIDs,
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "一次批量操作最多支持1000个用户")
	})
}

func TestBatchOperationService_BatchAddPoints(t *testing.T) {
	db := setupBatchServiceDB(t)
	defer testutil.CleanDB(t, db)

	svc := createBatchService(db)
	ctx := context.Background()

	// Create test users
	users := createTestUsersForBatch(t, db, 3)
	userIDs := make([]uint64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	t.Run("batch add points to specific users", func(t *testing.T) {
		success, failed, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target:  "users",
			UserIDs: userIDs,
			Cents:   100,
			Reason:  "测试积分",
			Type:    "admin",
		}, 1)
		require.NoError(t, err)
		assert.Equal(t, len(userIDs), success)
		assert.Equal(t, 0, failed)

		// Verify wallet balance
		var wallet model.Wallet
		require.NoError(t, db.Where("user_id = ?", users[0].ID).First(&wallet).Error)
		assert.Equal(t, int64(100), wallet.BalanceCents)
	})

	t.Run("batch add points to all users", func(t *testing.T) {
		success, failed, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target: "all",
			Cents:  50,
			Reason: "全员积分",
			Type:   "activity",
		}, 1)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, success, 0)
		assert.Equal(t, 0, failed)
	})

	t.Run("batch add points by role", func(t *testing.T) {
		success, failed, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target: "role",
			Roles:  []string{"user"},
			Cents:  25,
			Reason: "角色积分",
			Type:   "compensation",
		}, 1)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, success, 0)
		assert.Equal(t, 0, failed)
	})

	t.Run("empty user IDs for users target", func(t *testing.T) {
		_, _, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target:  "users",
			UserIDs: []uint64{},
			Cents:   100,
			Reason:  "测试",
			Type:    "admin",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "userIds不能为空")
	})

	t.Run("empty roles for role target", func(t *testing.T) {
		_, _, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target: "role",
			Roles:  []string{},
			Cents:  100,
			Reason: "测试",
			Type:   "admin",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "roles不能为空")
	})

	t.Run("invalid target", func(t *testing.T) {
		_, _, err := svc.BatchAddPoints(ctx, &BatchAddPointsRequest{
			Target: "invalid",
			Cents:  100,
			Reason: "测试",
			Type:   "admin",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的target类型")
	})
}

func TestBatchOperationService_BatchSendNotification(t *testing.T) {
	db := setupBatchServiceDB(t)
	defer testutil.CleanDB(t, db)

	svc := createBatchService(db)
	ctx := context.Background()

	// Create test users
	users := createTestUsersForBatch(t, db, 3)
	userIDs := make([]uint64, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	t.Run("send notification to specific users", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "users",
			UserIDs: userIDs,
			Title:   "测试通知",
			Content: "这是测试内容",
			Type:    "system",
		}, 1)
		require.NoError(t, err)
	})

	t.Run("send notification to all users", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "all",
			Title:   "全员通知",
			Content: "这是全员通知内容",
			Type:    "marketing",
		}, 1)
		require.NoError(t, err)
	})

	t.Run("send notification by role", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "role",
			Roles:   []string{"user"},
			Title:   "角色通知",
			Content: "这是角色通知内容",
			Type:    "personal",
		}, 1)
		require.NoError(t, err)
	})

	t.Run("empty user IDs for users target", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "users",
			UserIDs: []uint64{},
			Title:   "标题",
			Content: "内容",
			Type:    "system",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "userIds不能为空")
	})

	t.Run("empty roles for role target", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "role",
			Roles:   []string{},
			Title:   "标题",
			Content: "内容",
			Type:    "system",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "roles不能为空")
	})

	t.Run("invalid target", func(t *testing.T) {
		err := svc.BatchSendNotification(ctx, &BatchSendNotificationRequest{
			Target:  "invalid",
			Title:   "标题",
			Content: "内容",
			Type:    "system",
		}, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的target类型")
	})
}

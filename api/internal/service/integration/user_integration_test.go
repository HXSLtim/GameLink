// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	"gamelink/internal/repository/notification"
	userrepository "gamelink/internal/repository/user"
	authservice "gamelink/internal/service/auth"
	"gamelink/internal/service/user"
	pkgauth "gamelink/pkg/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_BanUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test user
	testUser := CreateUniqueTestUser(t, db, "ban_user")
	assert.Equal(t, model.UserStatusActive, testUser.Status)

	// Ban user
	req := &user.BatchUpdateUserStatusRequest{
		UserIDs: []uint64{testUser.ID},
		Status:  "banned",
		Reason:  "违规操作",
	}

	successCount, failedCount, err := svc.BatchUpdateUserStatus(ctx, req, 1) // operatorID = 1
	require.NoError(t, err)
	assert.Equal(t, 1, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify user is banned
	var bannedUser model.User
	err = db.First(&bannedUser, testUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.UserStatusBanned, bannedUser.Status)
}

func TestUserService_UnbanUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create banned user
	testUser := CreateUniqueTestUser(t, db, "unban_user")
	testUser.Status = model.UserStatusBanned
	db.Save(testUser)

	// Unban user
	req := &user.BatchUpdateUserStatusRequest{
		UserIDs: []uint64{testUser.ID},
		Status:  "active",
		Reason:  "解封",
	}

	successCount, failedCount, err := svc.BatchUpdateUserStatus(ctx, req, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify user is active
	var activeUser model.User
	err = db.First(&activeUser, testUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.UserStatusActive, activeUser.Status)
}

func TestUserService_SuspendUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test user
	testUser := CreateUniqueTestUser(t, db, "suspend_user")

	// Suspend user
	req := &user.BatchUpdateUserStatusRequest{
		UserIDs: []uint64{testUser.ID},
		Status:  "suspended",
		Reason:  "暂时冻结",
	}

	successCount, failedCount, err := svc.BatchUpdateUserStatus(ctx, req, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify user is suspended
	var suspendedUser model.User
	err = db.First(&suspendedUser, testUser.ID).Error
	require.NoError(t, err)
	assert.Equal(t, model.UserStatusSuspended, suspendedUser.Status)
}

func TestAuthService_ChangePassword(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepository.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)

	// Create user with password
	oldPassword := "OldPassword123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser := CreateUniqueTestUser(t, db, "changepwd_user")
	testUser.PasswordHash = string(hashedPassword)
	db.Save(testUser)

	// Change password
	newPassword := "NewPassword456"
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	require.NoError(t, err)

	testUser.PasswordHash = string(newHashedPassword)
	err = db.Save(testUser).Error
	require.NoError(t, err)

	// Verify can login with new password
	loginReq := authservice.LoginRequest{
		Username: testUser.Email,
		Password: newPassword,
	}
	svc := authservice.NewAuthService(userRepo, jwtManager)
	_, err = svc.Login(ctx, loginReq)
	require.NoError(t, err)

	// Verify cannot login with old password
	loginReq.Password = oldPassword
	_, err = svc.Login(ctx, loginReq)
	assert.Error(t, err)
}

func TestUserService_BatchBanUsers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create multiple test users
	var userIDs []uint64
	for i := 0; i < 5; i++ {
		u := CreateUniqueTestUser(t, db, "batch_ban_user"+string(rune('0'+i)))
		userIDs = append(userIDs, u.ID)
	}

	// Ban all users
	req := &user.BatchUpdateUserStatusRequest{
		UserIDs: userIDs,
		Status:  "banned",
		Reason:  "批量封禁测试",
	}

	successCount, failedCount, err := svc.BatchUpdateUserStatus(ctx, req, 1)
	require.NoError(t, err)
	assert.Equal(t, 5, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify all users are banned
	for _, userID := range userIDs {
		var u model.User
		err = db.First(&u, userID).Error
		require.NoError(t, err)
		assert.Equal(t, model.UserStatusBanned, u.Status)
	}
}

func TestUserService_BatchAddPoints(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test users
	var userIDs []uint64
	for i := 0; i < 3; i++ {
		u := CreateUniqueTestUser(t, db, "points_user"+string(rune('0'+i)))
		userIDs = append(userIDs, u.ID)
	}

	// Add points to users
	req := &user.BatchAddPointsRequest{
		Target:  "users",
		UserIDs: userIDs,
		Cents:   10000, // 100元
		Reason:  "活动奖励",
		Type:    "activity",
	}

	successCount, failedCount, err := svc.BatchAddPoints(ctx, req, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify users have wallet balance
	for _, userID := range userIDs {
		var wallet model.Wallet
		err = db.Where("user_id = ?", userID).First(&wallet).Error
		require.NoError(t, err)
		assert.Equal(t, int64(10000), wallet.BalanceCents)
	}
}

func TestUserService_BatchAddPoints_ByRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test users with player role
	var playerUsers []uint64
	for i := 0; i < 3; i++ {
		u := CreateUniqueTestUser(t, db, "role_points_user"+string(rune('0'+i)))
		u.Role = model.RolePlayer
		db.Save(u)
		playerUsers = append(playerUsers, u.ID)
	}

	// Add points to all players
	req := &user.BatchAddPointsRequest{
		Target: "role",
		Roles:  []string{"player"},
		Cents:  5000, // 50元
		Reason: "陪玩师补贴",
		Type:   "admin",
	}

	successCount, failedCount, err := svc.BatchAddPoints(ctx, req, 1)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, successCount, 3)
	assert.Equal(t, 0, failedCount)

	// Verify players have wallet balance
	for _, userID := range playerUsers {
		var wallet model.Wallet
		err = db.Where("user_id = ?", userID).First(&wallet).Error
		require.NoError(t, err)
		assert.Equal(t, int64(5000), wallet.BalanceCents)
	}
}

func TestUserService_BatchSendNotification(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test users
	var userIDs []uint64
	for i := 0; i < 3; i++ {
		u := CreateUniqueTestUser(t, db, "notify_user"+string(rune('0'+i)))
		userIDs = append(userIDs, u.ID)
	}

	// Send notification to users
	req := &user.BatchSendNotificationRequest{
		Target:  "users",
		UserIDs: userIDs,
		Title:   "系统通知",
		Content: "这是一条测试通知",
		Type:    "system",
	}

	err := svc.BatchSendNotification(ctx, req, 1)
	require.NoError(t, err)

	// Verify notifications were created
	var count int64
	err = db.Model(&model.NotificationEvent{}).Where("title = ?", "系统通知").Count(&count).Error
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(3))
}

func TestUserService_BatchDeleteUsers(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test users
	var userIDs []uint64
	for i := 0; i < 3; i++ {
		u := CreateUniqueTestUser(t, db, "delete_user"+string(rune('0'+i)))
		userIDs = append(userIDs, u.ID)
	}

	// Delete users (soft delete)
	req := &user.BatchDeleteUsersRequest{
		UserIDs: userIDs,
		Reason:  "测试删除",
	}

	successCount, failedCount, err := svc.BatchDeleteUsers(ctx, req, 1)
	require.NoError(t, err)
	assert.Equal(t, 3, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify users are soft deleted (deleted_at is set)
	for _, userID := range userIDs {
		var u model.User
		err = db.Unscoped().First(&u, userID).Error
		require.NoError(t, err)
		assert.NotNil(t, u.DeletedAt)
	}
}

func TestUserService_BatchUpdateUserRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup repositories
	userRepo := userrepository.NewUserRepository(db)
	tagRepo := userrepository.NewUserTagRepository(db)
	notificationRepo := notification.NewNotificationRepository(db)
	svc := user.NewBatchOperationService(db, userRepo, tagRepo, notificationRepo)

	// Create test users
	var userIDs []uint64
	for i := 0; i < 3; i++ {
		u := CreateUniqueTestUser(t, db, "role_update_user"+string(rune('0'+i)))
		userIDs = append(userIDs, u.ID)
	}

	// Update user roles to player
	req := &user.BatchUpdateUserRoleRequest{
		UserIDs: userIDs,
		Role:    "player",
	}

	successCount, failedCount, err := svc.BatchUpdateUserRole(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, 3, successCount)
	assert.Equal(t, 0, failedCount)

	// Verify users have player role
	for _, userID := range userIDs {
		var u model.User
		err = db.First(&u, userID).Error
		require.NoError(t, err)
		assert.Equal(t, model.RolePlayer, u.Role)
	}
}

func TestAuthService_Login_BannedUserCannotLogin(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepository.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := authservice.NewAuthService(userRepo, jwtManager)

	// Create banned user with password
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "banned_login_user", password)
	testUser.Status = model.UserStatusBanned
	db.Save(testUser)

	// Try to login - should fail
	req := authservice.LoginRequest{
		Username: testUser.Email,
		Password: password,
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, authservice.ErrUserDisabled, err)
}

func TestAuthService_Login_SuspendedUserCannotLogin(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepository.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := authservice.NewAuthService(userRepo, jwtManager)

	// Create suspended user with password
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "suspended_login_user", password)
	testUser.Status = model.UserStatusSuspended
	db.Save(testUser)

	// Try to login - should fail
	req := authservice.LoginRequest{
		Username: testUser.Email,
		Password: password,
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, authservice.ErrUserDisabled, err)
}

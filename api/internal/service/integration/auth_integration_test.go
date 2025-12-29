// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"
	"time"

	"gamelink/internal/model"
	userrepo "gamelink/internal/repository/user"
	"gamelink/internal/service/auth"
	pkgauth "gamelink/pkg/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Register(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Test registration
	req := auth.RegisterRequest{
		Phone:    "13800138001",
		Email:    "test_register@test.com",
		Password: "Test123456",
		Name:     "Test User",
		Role:     model.RoleUser,
	}

	resp, err := svc.Register(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Phone, resp.User.Phone)
	assert.Equal(t, model.UserStatusActive, resp.User.Status)

	// Verify user in database
	var savedUser model.User
	err = db.First(&savedUser, resp.User.ID).Error
	require.NoError(t, err)
	assert.Equal(t, req.Name, savedUser.Name)
	assert.NotEmpty(t, savedUser.PasswordHash)
}

func TestAuthService_Register_DuplicatePhone(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create existing user
	existingUser := CreateUniqueTestUser(t, db, "existing")

	// Try to register with same phone
	req := auth.RegisterRequest{
		Phone:    existingUser.Phone,
		Email:    "new_email@test.com",
		Password: "Test123456",
		Name:     "New User",
		Role:     model.RoleUser,
	}

	_, err := svc.Register(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "手机号已被注册")
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create existing user
	existingUser := CreateUniqueTestUser(t, db, "existing")

	// Try to register with same email
	req := auth.RegisterRequest{
		Phone:    "13900139001",
		Email:    existingUser.Email,
		Password: "Test123456",
		Name:     "New User",
		Role:     model.RoleUser,
	}

	_, err := svc.Register(ctx, req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "邮箱已被注册")
}

func TestAuthService_Login(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create user with password
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "login_user", password)

	// Test login with email
	req := auth.LoginRequest{
		Username: testUser.Email,
		Password: password,
	}

	resp, err := svc.Login(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, testUser.ID, resp.User.ID)
	assert.Equal(t, testUser.Name, resp.User.Name)
}

func TestAuthService_Login_WithPhone(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create user with password
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "login_phone", password)

	// Test login with phone
	req := auth.LoginRequest{
		Username: testUser.Phone,
		Password: password,
	}

	resp, err := svc.Login(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, testUser.ID, resp.User.ID)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create user with password
	testUser := CreateTestUserWithPassword(t, db, "invalid_cred", "CorrectPassword123")

	// Test login with wrong password
	req := auth.LoginRequest{
		Username: testUser.Email,
		Password: "WrongPassword123",
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, auth.ErrInvalidCredentials, err)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Test login with non-existent user
	req := auth.LoginRequest{
		Username: "nonexistent@test.com",
		Password: "Test123456",
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, auth.ErrInvalidCredentials, err)
}

func TestAuthService_Login_DisabledUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create disabled user
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "disabled_user", password)
	testUser.Status = model.UserStatusBanned
	db.Save(testUser)

	// Test login
	req := auth.LoginRequest{
		Username: testUser.Email,
		Password: password,
	}

	_, err := svc.Login(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, auth.ErrUserDisabled, err)
}

func TestAuthService_RefreshToken(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup - use a 25 second token duration so it's immediately in refresh window (< 30s)
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 25*time.Second)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create user and login
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "refresh_user", password)

	loginResp, err := svc.Login(ctx, auth.LoginRequest{
		Username: testUser.Email,
		Password: password,
	})
	require.NoError(t, err)

	// Token has 25 seconds duration, which is < 30 second refresh window
	// So refresh should work immediately
	newToken, err := svc.RefreshToken(ctx, loginResp.Token)
	require.NoError(t, err)
	assert.NotEmpty(t, newToken)

	// Verify the new token is valid by using it to get user info
	// Me expects "Bearer <token>" format
	user, err := svc.Me(ctx, "Bearer "+newToken)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
}

func TestAuthService_Me(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Create user and login
	password := "Test123456"
	testUser := CreateTestUserWithPassword(t, db, "me_user", password)

	loginResp, err := svc.Login(ctx, auth.LoginRequest{
		Username: testUser.Email,
		Password: password,
	})
	require.NoError(t, err)

	// Get current user
	authHeader := "Bearer " + loginResp.Token
	user, err := svc.Me(ctx, authHeader)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Name, user.Name)
}

func TestAuthService_Me_InvalidToken(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Test with invalid token
	_, err := svc.Me(ctx, "Bearer invalid-token")
	assert.Error(t, err)
}

func TestAuthService_Me_MissingHeader(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	userRepo := userrepo.NewUserRepository(db)
	jwtManager := pkgauth.NewJWTManager("test-secret-key-for-integration-tests", 24*time.Hour)
	svc := auth.NewAuthService(userRepo, jwtManager)

	// Test with missing header
	_, err := svc.Me(ctx, "")
	assert.Error(t, err)
}

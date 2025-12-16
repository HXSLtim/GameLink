package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository/user"
	"gamelink/pkg/auth"
	"gamelink/pkg/testutil"
)

func setupAuthTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.User{}, &model.Wallet{})
	return db
}

func createTestUser(t *testing.T, db *gorm.DB, email, phone, password string) *model.User {
	t.Helper()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	u := &model.User{
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Name:         "Test User",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, db.Create(u).Error)
	return u
}

func TestAuthService_Login(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户
	testUser := createTestUser(t, db, "test@example.com", "13800138000", "password123")

	t.Run("通过邮箱登录成功", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), LoginRequest{
			Username: "test@example.com",
			Password: "password123",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, testUser.ID, resp.User.ID)
	})

	t.Run("通过手机号登录成功", func(t *testing.T) {
		resp, err := svc.Login(context.Background(), LoginRequest{
			Username: "13800138000",
			Password: "password123",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, testUser.ID, resp.User.ID)
	})

	t.Run("密码错误", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginRequest{
			Username: "test@example.com",
			Password: "wrongpassword",
		})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("用户不存在", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginRequest{
			Username: "notexist@example.com",
			Password: "password123",
		})
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("空用户名或密码", func(t *testing.T) {
		_, err := svc.Login(context.Background(), LoginRequest{
			Username: "",
			Password: "password123",
		})
		assert.Error(t, err)

		_, err = svc.Login(context.Background(), LoginRequest{
			Username: "test@example.com",
			Password: "",
		})
		assert.Error(t, err)
	})

	t.Run("禁用用户无法登录", func(t *testing.T) {
		// 创建禁用用户
		disabledUser := createTestUser(t, db, "disabled@example.com", "13800138001", "password123")
		disabledUser.Status = model.UserStatusBanned
		require.NoError(t, db.Save(disabledUser).Error)

		_, err := svc.Login(context.Background(), LoginRequest{
			Username: "disabled@example.com",
			Password: "password123",
		})
		assert.ErrorIs(t, err, ErrUserDisabled)
	})
}

func TestAuthService_Register(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	t.Run("注册成功", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "newuser@example.com",
			Phone:    "13900139000",
			Password: "password123",
			Name:     "New User",
			Role:     model.RoleUser,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "newuser@example.com", resp.User.Email)
		assert.Equal(t, "New User", resp.User.Name)
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		// 先注册一个用户
		_, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "existing@example.com",
			Phone:    "13900139001",
			Password: "password123",
			Name:     "Existing User",
		})
		require.NoError(t, err)

		// 尝试用相同邮箱注册
		_, err = svc.Register(context.Background(), RegisterRequest{
			Email:    "existing@example.com",
			Phone:    "13900139002",
			Password: "password123",
			Name:     "Another User",
		})
		assert.Error(t, err)
	})

	t.Run("手机号已存在", func(t *testing.T) {
		// 先注册一个用户
		_, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "phone1@example.com",
			Phone:    "13900139003",
			Password: "password123",
			Name:     "Phone User",
		})
		require.NoError(t, err)

		// 尝试用相同手机号注册
		_, err = svc.Register(context.Background(), RegisterRequest{
			Email:    "phone2@example.com",
			Phone:    "13900139003",
			Password: "password123",
			Name:     "Another Phone User",
		})
		assert.Error(t, err)
	})

	t.Run("缺少必填字段", func(t *testing.T) {
		// 缺少姓名
		_, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "noname@example.com",
			Password: "password123",
		})
		assert.Error(t, err)

		// 缺少邮箱和手机号
		_, err = svc.Register(context.Background(), RegisterRequest{
			Password: "password123",
			Name:     "No Contact",
		})
		assert.Error(t, err)

		// 缺少密码
		_, err = svc.Register(context.Background(), RegisterRequest{
			Email: "nopassword@example.com",
			Name:  "No Password",
		})
		assert.Error(t, err)
	})

	t.Run("密码太短", func(t *testing.T) {
		_, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "shortpwd@example.com",
			Password: "12345",
			Name:     "Short Password",
		})
		assert.Error(t, err)
	})

	t.Run("邮箱格式错误", func(t *testing.T) {
		_, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "invalid-email",
			Password: "password123",
			Name:     "Invalid Email",
		})
		assert.Error(t, err)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户并登录
	createTestUser(t, db, "refresh@example.com", "13800138002", "password123")
	loginResp, err := svc.Login(context.Background(), LoginRequest{
		Username: "refresh@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	t.Run("无效Token刷新失败", func(t *testing.T) {
		_, err := svc.RefreshToken(context.Background(), "invalid-token")
		assert.Error(t, err)
	})

	t.Run("有效Token可以验证", func(t *testing.T) {
		// 验证登录返回的Token是有效的
		claims, err := jwtManager.VerifyToken(loginResp.Token)
		require.NoError(t, err)
		assert.NotNil(t, claims)
	})
}

func TestAuthService_Me(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户并登录
	testUser := createTestUser(t, db, "me@example.com", "13800138003", "password123")
	loginResp, err := svc.Login(context.Background(), LoginRequest{
		Username: "me@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	t.Run("获取当前用户成功", func(t *testing.T) {
		user, err := svc.Me(context.Background(), "Bearer "+loginResp.Token)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, "me@example.com", user.Email)
	})

	t.Run("缺少认证头", func(t *testing.T) {
		_, err := svc.Me(context.Background(), "")
		assert.Error(t, err)
	})

	t.Run("无效Token", func(t *testing.T) {
		_, err := svc.Me(context.Background(), "Bearer invalid-token")
		assert.Error(t, err)
	})
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.org", true},
		{"user+tag@example.com", true},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
		{"test@tempmail.com", false},   // 临时邮箱
		{"test@mailinator.com", false}, // 临时邮箱
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthService_GetUser(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户
	testUser := createTestUser(t, db, "getuser@example.com", "13800138010", "password123")

	t.Run("获取用户成功", func(t *testing.T) {
		u, err := svc.GetUser(context.Background(), testUser.ID)
		require.NoError(t, err)
		assert.Equal(t, testUser.ID, u.ID)
		assert.Equal(t, "getuser@example.com", u.Email)
	})

	t.Run("用户不存在", func(t *testing.T) {
		_, err := svc.GetUser(context.Background(), 99999)
		assert.Error(t, err)
	})
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	t.Run("无效Token刷新失败", func(t *testing.T) {
		_, err := svc.RefreshToken(context.Background(), "invalid-token")
		assert.Error(t, err)
	})
}

func TestAuthService_RefreshToken_DisabledUser(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户并登录
	testUser := createTestUser(t, db, "refreshdisabled@example.com", "13800138012", "password123")
	loginResp, err := svc.Login(context.Background(), LoginRequest{
		Username: "refreshdisabled@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	// 禁用用户
	testUser.Status = model.UserStatusBanned
	require.NoError(t, db.Save(testUser).Error)

	t.Run("禁用用户无法刷新Token", func(t *testing.T) {
		_, err := svc.RefreshToken(context.Background(), loginResp.Token)
		assert.ErrorIs(t, err, ErrUserDisabled)
	})
}

func TestAuthService_Me_DisabledUser(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户并登录
	testUser := createTestUser(t, db, "medisabled@example.com", "13800138013", "password123")
	loginResp, err := svc.Login(context.Background(), LoginRequest{
		Username: "medisabled@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	// 禁用用户
	testUser.Status = model.UserStatusBanned
	require.NoError(t, db.Save(testUser).Error)

	t.Run("禁用用户无法获取信息", func(t *testing.T) {
		_, err := svc.Me(context.Background(), "Bearer "+loginResp.Token)
		assert.ErrorIs(t, err, ErrUserDisabled)
	})
}

func TestAuthService_Register_OnlyPhone(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	t.Run("只用手机号注册成功", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), RegisterRequest{
			Phone:    "13900139010",
			Password: "password123",
			Name:     "Phone Only User",
			Role:     model.RoleUser,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "13900139010", resp.User.Phone)
	})
}

func TestAuthService_Register_DefaultRole(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	t.Run("不指定角色使用默认角色", func(t *testing.T) {
		resp, err := svc.Register(context.Background(), RegisterRequest{
			Email:    "defaultrole@example.com",
			Phone:    "13900139011",
			Password: "password123",
			Name:     "Default Role User",
			// Role 不指定
		})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})
}

func TestIsValidEmail_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"very long email", "a" + string(make([]byte, 130)) + "@example.com", false},
		{"subdomain email", "test@sub.example.com", true},
		{"numeric domain", "test@123.com", true},
		{"guerrillamail", "test@guerrillamail.com", false},
		{"10minutemail", "test@10minutemail.com", false},
		{"subdomain of disposable", "test@sub.tempmail.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateRegisterInput(t *testing.T) {
	tests := []struct {
		name    string
		req     RegisterRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: RegisterRequest{
				Email:    "valid@example.com",
				Phone:    "13800138000",
				Password: "password123",
				Name:     "Valid User",
				Role:     model.RoleUser,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: RegisterRequest{
				Email:    "valid@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "no email or phone",
			req: RegisterRequest{
				Password: "password123",
				Name:     "No Contact",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			req: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Invalid Email",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			req: RegisterRequest{
				Email: "valid@example.com",
				Name:  "No Password",
			},
			wantErr: true,
		},
		{
			name: "short password",
			req: RegisterRequest{
				Email:    "valid@example.com",
				Password: "12345",
				Name:     "Short Password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegisterInput(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewAuthService(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)

	svc := NewAuthService(userRepo, jwtManager)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.userRepo)
	assert.NotNil(t, svc.jwtManager)
}

func TestAuthService_Login_UpdateLastLoginTime(t *testing.T) {
	db := setupAuthTestDB(t)
	defer testutil.CleanDB(t, db)

	userRepo := user.NewUserRepository(db)
	jwtManager := auth.NewJWTManager("test-secret-key-32-chars-long!!", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	// 创建测试用户
	testUser := createTestUser(t, db, "lastlogin@example.com", "13800138020", "password123")
	assert.Nil(t, testUser.LastLoginAt)

	// 登录
	_, err := svc.Login(context.Background(), LoginRequest{
		Username: "lastlogin@example.com",
		Password: "password123",
	})
	require.NoError(t, err)

	// 验证最后登录时间已更新
	var updatedUser model.User
	require.NoError(t, db.First(&updatedUser, testUser.ID).Error)
	assert.NotNil(t, updatedUser.LastLoginAt)
}

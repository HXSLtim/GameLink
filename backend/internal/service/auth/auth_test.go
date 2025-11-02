package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/mocks"
)

// TestNewAuthService 测试构造函数。
func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)

	svc := NewAuthService(userRepo, jwtManager)

	if svc == nil {
		t.Fatal("NewAuthService returned nil")
	}

	if svc.userRepo != userRepo {
		t.Error("userRepo not set correctly")
	}

	if svc.jwtManager != jwtManager {
		t.Error("jwtManager not set correctly")
	}
}

// TestAuthServiceWithNilInputs 测试nil输入处理。
func TestAuthServiceWithNilInputs(t *testing.T) {
	// 虽然不推荐传nil，但我们测试service是否会panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewAuthService panicked with nil inputs: %v", r)
		}
	}()

	_ = NewAuthService(nil, nil)
}

// TestGetUser 测试GetUser方法
func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()
	userID := uint64(123)

	t.Run("成功获取用户", func(t *testing.T) {
		expectedUser := &model.User{
			Base:  model.Base{ID: userID},
			Email: "test@example.com",
			Name:  "Test User",
			Role:  model.RoleUser,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(expectedUser, nil)

		user, err := svc.GetUser(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.ID != userID {
			t.Errorf("Expected user ID %d, got %d", userID, user.ID)
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		userRepo.EXPECT().
			Get(ctx, userID).
			Return(nil, repository.ErrNotFound)

		user, err := svc.GetUser(ctx, userID)
		if err != repository.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	})
}

// TestMe 测试Me方法（验证Token并返回当前用户）
func TestMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()
	userID := uint64(123)

	t.Run("成功验证Token", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(userID, string(model.RoleUser))
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		authHeader := "Bearer " + token

		expectedUser := &model.User{
			Base:   model.Base{ID: userID},
			Email:  "test@example.com",
			Name:   "Test User",
			Role:   model.RoleUser,
			Status: model.UserStatusActive,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(expectedUser, nil)

		user, err := svc.Me(ctx, authHeader)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.ID != userID {
			t.Errorf("Expected user ID %d, got %d", userID, user.ID)
		}
	})

	t.Run("缺少Authorization头", func(t *testing.T) {
		user, err := svc.Me(ctx, "")
		if err == nil {
			t.Error("Expected error for missing authorization header")
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	})

	t.Run("无效的Token格式", func(t *testing.T) {
		user, err := svc.Me(ctx, "InvalidToken")
		if err == nil {
			t.Error("Expected error for invalid token format")
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	})

	t.Run("用户已禁用", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(userID, string(model.RoleUser))
		authHeader := "Bearer " + token

		disabledUser := &model.User{
			Base:   model.Base{ID: userID},
			Email:  "test@example.com",
			Status: model.UserStatusSuspended,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(disabledUser, nil)

		user, err := svc.Me(ctx, authHeader)
		if err != ErrUserDisabled {
			t.Errorf("Expected ErrUserDisabled, got %v", err)
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(userID, string(model.RoleUser))
		authHeader := "Bearer " + token

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(nil, repository.ErrNotFound)

		user, err := svc.Me(ctx, authHeader)
		if err == nil {
			t.Error("Expected error for user not found")
		}

		if user != nil {
			t.Error("Expected nil user")
		}
	})
}

// TestLogin 测试Login方法
func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()

	t.Run("通过邮箱成功登录", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &model.User{
			Base:         model.Base{ID: 1},
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			Name:         "Test User",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		}

		userRepo.EXPECT().
			FindByEmail(ctx, "test@example.com").
			Return(user, nil)

		userRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil)

		req := LoginRequest{
			Username: "test@example.com",
			Password: password,
		}

		resp, err := svc.Login(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Token == "" {
			t.Error("Expected token to be generated")
		}

		if resp.User.ID != user.ID {
			t.Errorf("Expected user ID %d, got %d", user.ID, resp.User.ID)
		}
	})

	t.Run("通过手机号成功登录", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &model.User{
			Base:         model.Base{ID: 1},
			Phone:        "13800138000",
			PasswordHash: string(hashedPassword),
			Name:         "Test User",
			Role:         model.RoleUser,
			Status:       model.UserStatusActive,
		}

		userRepo.EXPECT().
			FindByPhone(ctx, "13800138000").
			Return(user, nil)

		userRepo.EXPECT().
			Update(ctx, gomock.Any()).
			Return(nil)

		req := LoginRequest{
			Username: "13800138000",
			Password: password,
		}

		resp, err := svc.Login(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Token == "" {
			t.Error("Expected token to be generated")
		}
	})

	t.Run("用户名为空", func(t *testing.T) {
		req := LoginRequest{
			Username: "",
			Password: "password123",
		}

		resp, err := svc.Login(ctx, req)
		if err == nil {
			t.Error("Expected error for empty username")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("密码为空", func(t *testing.T) {
		req := LoginRequest{
			Username: "test@example.com",
			Password: "",
		}

		resp, err := svc.Login(ctx, req)
		if err == nil {
			t.Error("Expected error for empty password")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		userRepo.EXPECT().
			FindByEmail(ctx, "notfound@example.com").
			Return(nil, repository.ErrNotFound)

		req := LoginRequest{
			Username: "notfound@example.com",
			Password: "password123",
		}

		resp, err := svc.Login(ctx, req)
		if err != ErrInvalidCredentials {
			t.Errorf("Expected ErrInvalidCredentials, got %v", err)
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("密码错误", func(t *testing.T) {
		password := "correctpassword"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &model.User{
			Base:         model.Base{ID: 1},
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			Status:       model.UserStatusActive,
		}

		userRepo.EXPECT().
			FindByEmail(ctx, "test@example.com").
			Return(user, nil)

		req := LoginRequest{
			Username: "test@example.com",
			Password: "wrongpassword",
		}

		resp, err := svc.Login(ctx, req)
		if err != ErrInvalidCredentials {
			t.Errorf("Expected ErrInvalidCredentials, got %v", err)
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("用户已被禁用", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &model.User{
			Base:         model.Base{ID: 1},
			Email:        "test@example.com",
			PasswordHash: string(hashedPassword),
			Status:       model.UserStatusSuspended,
		}

		userRepo.EXPECT().
			FindByEmail(ctx, "test@example.com").
			Return(user, nil)

		req := LoginRequest{
			Username: "test@example.com",
			Password: password,
		}

		resp, err := svc.Login(ctx, req)
		if err != ErrUserDisabled {
			t.Errorf("Expected ErrUserDisabled, got %v", err)
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})
}

// TestRegister 测试Register方法
func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()

	t.Run("成功注册（通过邮箱）", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
			Role:     model.RoleUser,
		}

		userRepo.EXPECT().
			FindByEmail(ctx, req.Email).
			Return(nil, repository.ErrNotFound)

		userRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *model.User) error {
				user.ID = 123
				return nil
			})

		resp, err := svc.Register(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Token == "" {
			t.Error("Expected token to be generated")
		}

		if resp.User.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, resp.User.Email)
		}
	})

	t.Run("成功注册（通过手机号）", func(t *testing.T) {
		req := RegisterRequest{
			Phone:    "13800138000",
			Password: "password123",
			Name:     "New User",
			Role:     model.RoleUser,
		}

		userRepo.EXPECT().
			FindByPhone(ctx, req.Phone).
			Return(nil, repository.ErrNotFound)

		userRepo.EXPECT().
			Create(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, user *model.User) error {
				user.ID = 123
				return nil
			})

		resp, err := svc.Register(ctx, req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Token == "" {
			t.Error("Expected token to be generated")
		}
	})

	t.Run("姓名为空", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "",
		}

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for empty name")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("邮箱和手机号都为空", func(t *testing.T) {
		req := RegisterRequest{
			Password: "password123",
			Name:     "Test User",
		}

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for missing email and phone")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("密码为空", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "",
			Name:     "Test User",
		}

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for empty password")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("密码太短", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "12345",
			Name:     "Test User",
		}

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for short password")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("邮箱已存在", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "existing@example.com",
			Password: "password123",
			Name:     "Test User",
		}

		existingUser := &model.User{
			Base:  model.Base{ID: 1},
			Email: req.Email,
		}

		userRepo.EXPECT().
			FindByEmail(ctx, req.Email).
			Return(existingUser, nil)

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for existing email")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})

	t.Run("手机号已存在", func(t *testing.T) {
		req := RegisterRequest{
			Phone:    "13800138000",
			Password: "password123",
			Name:     "Test User",
		}

		existingUser := &model.User{
			Base:  model.Base{ID: 1},
			Phone: req.Phone,
		}

		userRepo.EXPECT().
			FindByPhone(ctx, req.Phone).
			Return(existingUser, nil)

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for existing phone")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})
}

// TestRefreshToken 测试RefreshToken方法
func TestRefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()
	userID := uint64(123)

	t.Run("Token还未到刷新时间", func(t *testing.T) {
		// 新生成的Token不能立即刷新（剩余时间>30秒）
		token, _ := jwtManager.GenerateToken(userID, string(model.RoleUser))

		activeUser := &model.User{
			Base:   model.Base{ID: userID},
			Email:  "test@example.com",
			Status: model.UserStatusActive,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(activeUser, nil)

		newToken, err := svc.RefreshToken(ctx, token)
		if err == nil {
			t.Error("Expected error for token not ready to refresh")
		}

		if newToken != "" {
			t.Error("Expected empty token")
		}
	})

	t.Run("成功刷新Token（使用短期Token）", func(t *testing.T) {
		// 创建一个短期有效的JWT Manager (20秒)
		shortJWT := auth.NewJWTManager("test-secret", 20*time.Second)
		shortSvc := NewAuthService(userRepo, shortJWT)

		token, _ := shortJWT.GenerateToken(userID, string(model.RoleUser))

		// 等待一段时间，让Token接近过期（剩余时间<30秒）
		time.Sleep(1 * time.Second)

		activeUser := &model.User{
			Base:   model.Base{ID: userID},
			Email:  "test@example.com",
			Status: model.UserStatusActive,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(activeUser, nil)

		newToken, err := shortSvc.RefreshToken(ctx, token)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if newToken == "" {
			t.Error("Expected new token to be generated")
		}

		if newToken == token {
			t.Error("Expected new token to be different from old token")
		}
	})

	t.Run("无效的Token", func(t *testing.T) {
		newToken, err := svc.RefreshToken(ctx, "invalid-token")
		if err == nil {
			t.Error("Expected error for invalid token")
		}

		if newToken != "" {
			t.Error("Expected empty token")
		}
	})

	t.Run("用户不存在", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(userID, string(model.RoleUser))

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(nil, repository.ErrNotFound)

		newToken, err := svc.RefreshToken(ctx, token)
		if err == nil {
			t.Error("Expected error for user not found")
		}

		if newToken != "" {
			t.Error("Expected empty token")
		}
	})

	t.Run("用户已禁用", func(t *testing.T) {
		token, _ := jwtManager.GenerateToken(userID, string(model.RoleUser))

		disabledUser := &model.User{
			Base:   model.Base{ID: userID},
			Email:  "test@example.com",
			Status: model.UserStatusBanned,
		}

		userRepo.EXPECT().
			Get(ctx, userID).
			Return(disabledUser, nil)

		newToken, err := svc.RefreshToken(ctx, token)
		if err != ErrUserDisabled {
			t.Errorf("Expected ErrUserDisabled, got %v", err)
		}

		if newToken != "" {
			t.Error("Expected empty token")
		}
	})
}

// TestValidateRegisterInput 测试validateRegisterInput函数
func TestValidateRegisterInput(t *testing.T) {
	t.Run("所有字段正确", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
			Role:     model.RoleUser,
		}

		err := validateRegisterInput(req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("姓名为空", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "",
		}

		err := validateRegisterInput(req)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})

	t.Run("邮箱和手机号都为空", func(t *testing.T) {
		req := RegisterRequest{
			Password: "password123",
			Name:     "Test User",
		}

		err := validateRegisterInput(req)
		if err == nil {
			t.Error("Expected error for missing email and phone")
		}
	})

	t.Run("密码为空", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "",
			Name:     "Test User",
		}

		err := validateRegisterInput(req)
		if err == nil {
			t.Error("Expected error for empty password")
		}
	})

	t.Run("密码太短", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "12345",
			Name:     "Test User",
		}

		err := validateRegisterInput(req)
		if err == nil {
			t.Error("Expected error for short password")
		}
	})
}

// TestIsValidEmail 测试isValidEmail函数
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"有效邮箱", "test@example.com", true},
		{"有效邮箱（带加号）", "test+tag@example.com", true},
		{"有效邮箱（子域名）", "test@sub.example.com", true},
		{"空字符串", "", false},
		{"缺少@符号", "testexample.com", false},
		{"缺少域名", "test@", false},
		{"缺少用户名", "@example.com", false},
		{"无效格式", "not-an-email", false},
		{"多个@符号", "test@@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("isValidEmail(%q) = %v, expected %v", tt.email, result, tt.expected)
			}
		})
	}
}

// TestLoginDatabaseError 测试数据库错误处理
func TestLoginDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()

	t.Run("数据库查询错误", func(t *testing.T) {
		userRepo.EXPECT().
			FindByEmail(ctx, "test@example.com").
			Return(nil, errors.New("database connection error"))

		req := LoginRequest{
			Username: "test@example.com",
			Password: "password123",
		}

		resp, err := svc.Login(ctx, req)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})
}

// TestRegisterDatabaseError 测试注册时的数据库错误
func TestRegisterDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	svc := NewAuthService(userRepo, jwtManager)

	ctx := context.Background()

	t.Run("创建用户失败", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "new@example.com",
			Password: "password123",
			Name:     "New User",
		}

		userRepo.EXPECT().
			FindByEmail(ctx, req.Email).
			Return(nil, repository.ErrNotFound)

		userRepo.EXPECT().
			Create(ctx, gomock.Any()).
			Return(errors.New("database insert error"))

		resp, err := svc.Register(ctx, req)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if resp != nil {
			t.Error("Expected nil response")
		}
	})
}

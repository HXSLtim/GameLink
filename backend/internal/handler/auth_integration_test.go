package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gamelink/internal/auth"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	authservice "gamelink/internal/service/auth"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) List(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
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

func (m *MockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
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

func TestLoginHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:        "登录成功",
			requestBody: `{"username":"user@example.com","password":"password123"}`,
			setupMock: func(m *MockUserRepository) {
				// 模拟密码哈希验证 - 使用有效哈希 (password123)
				m.On("FindByEmail", mock.Anything, "user@example.com").Return(&model.User{
					Email:        "user@example.com",
					PasswordHash: "$2a$10$S7cG/Zl76jygWqZtC4LEJ.Zw72aWWF0B3Ki8qAfOjPe9fzRwH3vPa",
					Name:         "Test User",
					Status:       model.UserStatusActive,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "登录成功",
		},
		{
			name:        "无效请求格式",
			requestBody: `{"username":"user@example.com"`, // 缺少闭合括号
			setupMock:   func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "无效的请求格式",
		},
		{
			name:        "用户不存在",
			requestBody: `{"username":"nonexistent@example.com","password":"password123"}`,
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, repository.ErrNotFound).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "用户名或密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			// 创建真实的 AuthService
			jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
			svc := authservice.NewAuthService(mockRepo, jwtManager)

			router := gin.New()
			RegisterAuthRoutes(router, svc)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Contains(t, response["message"].(string), tt.expectedMsg)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRegisterHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockUserRepository)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:        "注册成功",
			requestBody: `{"email":"newuser@example.com","password":"password123","name":"New User"}`,
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "newuser@example.com").Return(nil, repository.ErrNotFound).Once()
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "登录成功",
		},
		{
			name:        "邮箱已存在",
			requestBody: `{"email":"existing@example.com","password":"password123","name":"Test User"}`,
			setupMock: func(m *MockUserRepository) {
				m.On("FindByEmail", mock.Anything, "existing@example.com").Return(&model.User{
					Email: "existing@example.com",
					Name:  "Existing User",
				}, nil).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "注册失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setupMock(mockRepo)

			jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
			svc := authservice.NewAuthService(mockRepo, jwtManager)

			router := gin.New()
			RegisterAuthRoutes(router, svc)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Contains(t, response["message"].(string), tt.expectedMsg)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogoutHandler_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockUserRepository)
	jwtManager := auth.NewJWTManager("test-secret-key-that-is-32-characters-long", 24*time.Hour)
	svc := authservice.NewAuthService(mockRepo, jwtManager)

	router := gin.New()
	RegisterAuthRoutes(router, svc)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "登出成功", response["message"])
}

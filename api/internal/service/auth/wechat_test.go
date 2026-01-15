package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// ============================================================================
// Mock Implementations
// ============================================================================

// MockWeChatUserRepository is a mock implementation of UserRepository for WeChat tests
type MockWeChatUserRepository struct {
	mock.Mock
}

func (m *MockWeChatUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockWeChatUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	args := m.Called(ctx, openID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockWeChatUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	args := m.Called(ctx, unionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockWeChatUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil && user.ID == 0 {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockWeChatUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockWeChatUserRepository) List(ctx context.Context) ([]model.User, error)                   { return nil, nil }
func (m *MockWeChatUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockWeChatUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockWeChatUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) { return 0, nil }
func (m *MockWeChatUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error)      { return nil, nil }
func (m *MockWeChatUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error)     { return nil, nil }
func (m *MockWeChatUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error)     { return nil, nil }
func (m *MockWeChatUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error)    { return nil, nil }
func (m *MockWeChatUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error)    { return nil, nil }
func (m *MockWeChatUserRepository) Delete(ctx context.Context, id uint64) error                          { return nil }
func (m *MockWeChatUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error { return nil }

// MockWeChatPlayerRepository is a mock implementation of PlayerRepository for WeChat tests
type MockWeChatPlayerRepository struct {
	mock.Mock
}

func (m *MockWeChatPlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockWeChatPlayerRepository) Create(ctx context.Context, player *model.Player) error { return nil }
func (m *MockWeChatPlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) { return nil, nil }
func (m *MockWeChatPlayerRepository) Update(ctx context.Context, player *model.Player) error { return nil }
func (m *MockWeChatPlayerRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockWeChatPlayerRepository) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (m *MockWeChatPlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockWeChatPlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockWeChatPlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) { return nil, nil }
func (m *MockWeChatPlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) { return 0, nil }
func (m *MockWeChatPlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) { return 0, nil }
func (m *MockWeChatPlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) { return 0, nil }
func (m *MockWeChatPlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) { return 0, nil }

// ============================================================================
// Tests
// ============================================================================

func TestNewWeChatAuthService(t *testing.T) {
	users := &MockWeChatUserRepository{}
	players := &MockWeChatPlayerRepository{}

	svc := NewWeChatAuthService(users, players)

	assert.NotNil(t, svc)
	assert.Equal(t, users, svc.users)
	assert.Equal(t, players, svc.players)
}

func TestWeChatAuthService_checkIsPlayer(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint64
		setupMock  func(*MockWeChatPlayerRepository)
		expectBool bool
	}{
		{
			name:   "user is verified player",
			userID: 1,
			setupMock: func(m *MockWeChatPlayerRepository) {
				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				m.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectBool: true,
		},
		{
			name:   "user is pending player",
			userID: 2,
			setupMock: func(m *MockWeChatPlayerRepository) {
				player := &model.Player{
					UserID:             2,
					VerificationStatus: model.VerificationPending,
				}
				player.ID = 2
				m.On("GetByUserID", mock.Anything, uint64(2)).Return(player, nil)
			},
			expectBool: false,
		},
		{
			name:   "user is not a player",
			userID: 3,
			setupMock: func(m *MockWeChatPlayerRepository) {
				m.On("GetByUserID", mock.Anything, uint64(3)).Return(nil, repository.ErrNotFound)
			},
			expectBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockWeChatUserRepository{}
			players := &MockWeChatPlayerRepository{}

			tt.setupMock(players)

			svc := NewWeChatAuthService(users, players)
			result := svc.checkIsPlayer(context.Background(), tt.userID)

			assert.Equal(t, tt.expectBool, result)
		})
	}
}

func TestWeChatAuthService_checkIsPlayer_NilPlayers(t *testing.T) {
	users := &MockWeChatUserRepository{}

	svc := NewWeChatAuthService(users, nil)
	result := svc.checkIsPlayer(context.Background(), 1)

	assert.False(t, result)
}

func TestPkcs7Unpad(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "normal padding",
			input:    []byte{1, 2, 3, 4, 5, 3, 3, 3},
			expected: []byte{1, 2, 3, 4, 5},
		},
		{
			name:     "single byte padding",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 1},
			expected: []byte{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "padding larger than length",
			input:    []byte{1, 2, 100},
			expected: []byte{1, 2, 100},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pkcs7Unpad(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeChatLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		req         WeChatLoginRequest
		expectValid bool
	}{
		{
			name: "valid request with code only",
			req: WeChatLoginRequest{
				Code: "test_code",
			},
			expectValid: true,
		},
		{
			name: "valid request with all fields",
			req: WeChatLoginRequest{
				Code:          "test_code",
				EncryptedData: "encrypted_data",
				IV:            "iv_value",
				ReferralCode:  "REF123",
			},
			expectValid: true,
		},
		{
			name: "invalid request without code",
			req: WeChatLoginRequest{
				EncryptedData: "encrypted_data",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.req.Code != ""
			assert.Equal(t, tt.expectValid, isValid)
		})
	}
}

func TestWeChatLoginResponse_Structure(t *testing.T) {
	resp := WeChatLoginResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		ExpiresIn:    86400,
		User: &UserInfo{
			ID:          1,
			Nickname:    "Test User",
			Avatar:      "https://example.com/avatar.jpg",
			Phone:       "13800138000",
			IsPlayer:    true,
			CurrentRole: "player",
		},
	}

	assert.Equal(t, "access_token", resp.AccessToken)
	assert.Equal(t, "refresh_token", resp.RefreshToken)
	assert.Equal(t, int64(86400), resp.ExpiresIn)
	assert.NotNil(t, resp.User)
	assert.Equal(t, uint64(1), resp.User.ID)
	assert.Equal(t, "Test User", resp.User.Nickname)
	assert.True(t, resp.User.IsPlayer)
	assert.Equal(t, "player", resp.User.CurrentRole)
}

func TestWeChatSession_Structure(t *testing.T) {
	session := WeChatSession{
		OpenID:     "openid_123",
		SessionKey: "session_key_456",
		UnionID:    "unionid_789",
		ErrCode:    0,
		ErrMsg:     "",
	}

	assert.Equal(t, "openid_123", session.OpenID)
	assert.Equal(t, "session_key_456", session.SessionKey)
	assert.Equal(t, "unionid_789", session.UnionID)
	assert.Equal(t, 0, session.ErrCode)
}

func TestWeChatPhoneInfo_Structure(t *testing.T) {
	phoneInfo := WeChatPhoneInfo{
		PhoneNumber:     "+8613800138000",
		PurePhoneNumber: "13800138000",
		CountryCode:     "86",
	}

	assert.Equal(t, "+8613800138000", phoneInfo.PhoneNumber)
	assert.Equal(t, "13800138000", phoneInfo.PurePhoneNumber)
	assert.Equal(t, "86", phoneInfo.CountryCode)
}

func TestWeChatErrors(t *testing.T) {
	assert.NotNil(t, ErrWeChatCodeInvalid)
	assert.NotNil(t, ErrWeChatSessionFailed)
	assert.NotNil(t, ErrDecryptFailed)

	assert.Equal(t, "微信登录凭证无效", ErrWeChatCodeInvalid.Error())
	assert.Equal(t, "获取微信会话失败", ErrWeChatSessionFailed.Error())
	assert.Equal(t, "解密数据失败", ErrDecryptFailed.Error())
}

func TestUserInfo_Structure(t *testing.T) {
	userInfo := UserInfo{
		ID:          123,
		Nickname:    "TestNickname",
		Avatar:      "https://example.com/avatar.png",
		Phone:       "13900139000",
		IsPlayer:    false,
		CurrentRole: "user",
	}

	assert.Equal(t, uint64(123), userInfo.ID)
	assert.Equal(t, "TestNickname", userInfo.Nickname)
	assert.Equal(t, "https://example.com/avatar.png", userInfo.Avatar)
	assert.Equal(t, "13900139000", userInfo.Phone)
	assert.False(t, userInfo.IsPlayer)
	assert.Equal(t, "user", userInfo.CurrentRole)
}

// ============================================================================
// Additional Tests for Coverage Improvement
// ============================================================================

func TestWeChatAuthService_findOrCreateWeChatUser(t *testing.T) {
	tests := []struct {
		name        string
		openID      string
		unionID     string
		phone       string
		setupMock   func(*MockWeChatUserRepository)
		expectError bool
		expectNew   bool
	}{
		{
			name:    "find existing user by OpenID",
			openID:  "openid_123",
			unionID: "",
			phone:   "",
			setupMock: func(m *MockWeChatUserRepository) {
				user := &model.User{
					Name:   "Existing User",
					Phone:  "13800138000",
					Status: model.UserStatusActive,
				}
				user.ID = 1
				m.On("GetByWeChatOpenID", mock.Anything, "openid_123").Return(user, nil)
			},
			expectError: false,
			expectNew:   false,
		},
		{
			name:    "find existing user by OpenID and update phone",
			openID:  "openid_123",
			unionID: "",
			phone:   "13900139000",
			setupMock: func(m *MockWeChatUserRepository) {
				user := &model.User{
					Name:   "Existing User",
					Phone:  "",
					Status: model.UserStatusActive,
				}
				user.ID = 1
				m.On("GetByWeChatOpenID", mock.Anything, "openid_123").Return(user, nil)
				m.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Phone == "13900139000"
				})).Return(nil)
			},
			expectError: false,
			expectNew:   false,
		},
		{
			name:    "find existing user by UnionID",
			openID:  "openid_new",
			unionID: "unionid_123",
			phone:   "",
			setupMock: func(m *MockWeChatUserRepository) {
				m.On("GetByWeChatOpenID", mock.Anything, "openid_new").Return(nil, repository.ErrNotFound)
				user := &model.User{
					Name:   "Union User",
					Status: model.UserStatusActive,
				}
				user.ID = 2
				m.On("GetByWeChatUnionID", mock.Anything, "unionid_123").Return(user, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
			expectNew:   false,
		},
		{
			name:    "create new user",
			openID:  "openid_new",
			unionID: "",
			phone:   "13700137000",
			setupMock: func(m *MockWeChatUserRepository) {
				m.On("GetByWeChatOpenID", mock.Anything, "openid_new").Return(nil, repository.ErrNotFound)
				m.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.Phone == "13700137000" && *u.WeChatOpenID == "openid_new"
				})).Return(nil)
			},
			expectError: false,
			expectNew:   true,
		},
		{
			name:    "create new user with UnionID not found",
			openID:  "openid_new",
			unionID: "unionid_new",
			phone:   "",
			setupMock: func(m *MockWeChatUserRepository) {
				m.On("GetByWeChatOpenID", mock.Anything, "openid_new").Return(nil, repository.ErrNotFound)
				m.On("GetByWeChatUnionID", mock.Anything, "unionid_new").Return(nil, repository.ErrNotFound)
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
			expectNew:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockWeChatUserRepository{}
			tt.setupMock(users)

			svc := NewWeChatAuthService(users, nil)
			user, isNew, err := svc.findOrCreateWeChatUser(context.Background(), tt.openID, tt.unionID, tt.phone)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectNew, isNew)
			}
			users.AssertExpectations(t)
		})
	}
}

func TestWeChatAuthService_generateWeChatToken(t *testing.T) {
	users := &MockWeChatUserRepository{}
	players := &MockWeChatPlayerRepository{}

	svc := NewWeChatAuthService(users, players)

	user := &model.User{
		Name:   "Test User",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	user.ID = 1

	token, err := svc.generateWeChatToken(user, false, "user")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestWeChatAuthService_generateWeChatToken_Player(t *testing.T) {
	users := &MockWeChatUserRepository{}
	players := &MockWeChatPlayerRepository{}

	svc := NewWeChatAuthService(users, players)

	user := &model.User{
		Name:   "Player User",
		Role:   model.RolePlayer,
		Status: model.UserStatusActive,
	}
	user.ID = 2

	token, err := svc.generateWeChatToken(user, true, "player")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestWeChatAuthService_generateRefreshToken(t *testing.T) {
	users := &MockWeChatUserRepository{}
	players := &MockWeChatPlayerRepository{}

	svc := NewWeChatAuthService(users, players)

	token, err := svc.generateRefreshToken(1)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestWeChatAuthService_processReferral(t *testing.T) {
	users := &MockWeChatUserRepository{}
	players := &MockWeChatPlayerRepository{}

	svc := NewWeChatAuthService(users, players)

	// processReferral is a stub that returns nil
	err := svc.processReferral(context.Background(), 1, "REF123")

	assert.NoError(t, err)
}

func TestWeChatAuthService_RefreshAccessToken(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockWeChatUserRepository, *MockWeChatPlayerRepository) string
		expectError bool
	}{
		{
			name: "success - refresh token for regular user",
			setupMock: func(users *MockWeChatUserRepository, players *MockWeChatPlayerRepository) string {
				user := &model.User{
					Name:   "Test User",
					Role:   model.RoleUser,
					Status: model.UserStatusActive,
				}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(nil, repository.ErrNotFound)

				// Generate a valid refresh token
				svc := NewWeChatAuthService(users, players)
				token, _ := svc.generateRefreshToken(1)
				return token
			},
			expectError: false,
		},
		{
			name: "success - refresh token for player",
			setupMock: func(users *MockWeChatUserRepository, players *MockWeChatPlayerRepository) string {
				user := &model.User{
					Name:   "Player User",
					Role:   model.RolePlayer,
					Status: model.UserStatusActive,
				}
				user.ID = 2
				users.On("Get", mock.Anything, uint64(2)).Return(user, nil)

				player := &model.Player{
					UserID:             2,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 2
				players.On("GetByUserID", mock.Anything, uint64(2)).Return(player, nil)

				svc := NewWeChatAuthService(users, players)
				token, _ := svc.generateRefreshToken(2)
				return token
			},
			expectError: false,
		},
		{
			name: "error - user not found",
			setupMock: func(users *MockWeChatUserRepository, players *MockWeChatPlayerRepository) string {
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)

				svc := NewWeChatAuthService(users, players)
				token, _ := svc.generateRefreshToken(999)
				return token
			},
			expectError: true,
		},
		{
			name: "error - user banned",
			setupMock: func(users *MockWeChatUserRepository, players *MockWeChatPlayerRepository) string {
				user := &model.User{
					Name:   "Banned User",
					Role:   model.RoleUser,
					Status: model.UserStatusBanned,
				}
				user.ID = 3
				users.On("Get", mock.Anything, uint64(3)).Return(user, nil)

				svc := NewWeChatAuthService(users, players)
				token, _ := svc.generateRefreshToken(3)
				return token
			},
			expectError: true,
		},
		{
			name: "error - invalid token",
			setupMock: func(users *MockWeChatUserRepository, players *MockWeChatPlayerRepository) string {
				return "invalid_token"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockWeChatUserRepository{}
			players := &MockWeChatPlayerRepository{}

			refreshToken := tt.setupMock(users, players)

			svc := NewWeChatAuthService(users, players)
			resp, err := svc.RefreshAccessToken(context.Background(), refreshToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.NotNil(t, resp.User)
			}
		})
	}
}

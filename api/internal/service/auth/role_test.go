package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/pkg/auth"
)

func init() {
	// Set JWT secret key for testing
	auth.SetDefaultSecretKey("test-secret-key-for-role-service-tests")
}

// ============================================================================
// Mock Implementations for Role Service
// ============================================================================

// MockRoleUserRepository is a mock implementation of UserRepository for role tests
type MockRoleUserRepository struct {
	mock.Mock
}

func (m *MockRoleUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRoleUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	args := m.Called(ctx, openID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRoleUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	args := m.Called(ctx, unionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRoleUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil && user.ID == 0 {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockRoleUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRoleUserRepository) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *MockRoleUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockRoleUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *MockRoleUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}
func (m *MockRoleUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	return nil, nil
}
func (m *MockRoleUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockRoleUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockRoleUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *MockRoleUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *MockRoleUserRepository) Delete(ctx context.Context, id uint64) error { return nil }
func (m *MockRoleUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
	return nil
}

// MockRolePlayerRepository is a mock implementation of PlayerRepository for role tests
type MockRolePlayerRepository struct {
	mock.Mock
}

func (m *MockRolePlayerRepository) GetByUserID(ctx context.Context, userID uint64) (*model.Player, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Player), args.Error(1)
}

func (m *MockRolePlayerRepository) Create(ctx context.Context, player *model.Player) error {
	return nil
}
func (m *MockRolePlayerRepository) Get(ctx context.Context, id uint64) (*model.Player, error) {
	return nil, nil
}
func (m *MockRolePlayerRepository) Update(ctx context.Context, player *model.Player) error {
	return nil
}
func (m *MockRolePlayerRepository) Delete(ctx context.Context, id uint64) error      { return nil }
func (m *MockRolePlayerRepository) List(ctx context.Context) ([]model.Player, error) { return nil, nil }
func (m *MockRolePlayerRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockRolePlayerRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}
func (m *MockRolePlayerRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.Player, error) {
	return nil, nil
}
func (m *MockRolePlayerRepository) BatchUpdateRank(ctx context.Context, ids []uint64, rank string) (int64, error) {
	return 0, nil
}
func (m *MockRolePlayerRepository) BatchUpdateHourlyRate(ctx context.Context, ids []uint64, rateCents int64) (int64, error) {
	return 0, nil
}
func (m *MockRolePlayerRepository) BatchUpdateStatus(ctx context.Context, ids []uint64, status model.VerificationStatus) (int64, error) {
	return 0, nil
}
func (m *MockRolePlayerRepository) BatchDelete(ctx context.Context, ids []uint64) (int64, error) {
	return 0, nil
}
func (m *MockRolePlayerRepository) ListFeatured(ctx context.Context, limit int, status *model.VerificationStatus) ([]model.Player, int64, error) {
	return nil, 0, nil
}

// ============================================================================
// Tests
// ============================================================================

func TestNewRoleService(t *testing.T) {
	users := &MockRoleUserRepository{}
	players := &MockRolePlayerRepository{}

	svc := NewRoleService(users, players)

	assert.NotNil(t, svc)
	assert.Equal(t, users, svc.users)
	assert.Equal(t, players, svc.players)
}

func TestRoleService_GetAvailableRoles(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint64
		currentRole    string
		setupMock      func(*MockRoleUserRepository, *MockRolePlayerRepository)
		expectError    bool
		errorContains  string
		expectIsPlayer bool
	}{
		{
			name:        "user with player role available",
			userID:      1,
			currentRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Test User", Status: model.UserStatusActive}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectError:    false,
			expectIsPlayer: true,
		},
		{
			name:        "user without player role",
			userID:      2,
			currentRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Regular User", Status: model.UserStatusActive}
				user.ID = 2
				users.On("Get", mock.Anything, uint64(2)).Return(user, nil)

				players.On("GetByUserID", mock.Anything, uint64(2)).Return(nil, repository.ErrNotFound)
			},
			expectError:    false,
			expectIsPlayer: false,
		},
		{
			name:        "user not found",
			userID:      999,
			currentRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError:   true,
			errorContains: "用户不存在",
		},
		{
			name:        "user is banned",
			userID:      3,
			currentRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Banned User", Status: model.UserStatusBanned}
				user.ID = 3
				users.On("Get", mock.Anything, uint64(3)).Return(user, nil)
			},
			expectError:   true,
			errorContains: "用户已被禁用",
		},
		{
			name:        "player with pending verification",
			userID:      4,
			currentRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Pending Player", Status: model.UserStatusActive}
				user.ID = 4
				users.On("Get", mock.Anything, uint64(4)).Return(user, nil)

				player := &model.Player{
					UserID:             4,
					VerificationStatus: model.VerificationPending,
				}
				player.ID = 4
				players.On("GetByUserID", mock.Anything, uint64(4)).Return(player, nil)
			},
			expectError:    false,
			expectIsPlayer: false, // pending players cannot switch to player role
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockRoleUserRepository{}
			players := &MockRolePlayerRepository{}

			tt.setupMock(users, players)

			svc := NewRoleService(users, players)
			resp, err := svc.GetAvailableRoles(context.Background(), tt.userID, tt.currentRole)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.currentRole, resp.CurrentRole)
				assert.Len(t, resp.Roles, 2)

				// Check user role is always available
				assert.Equal(t, "user", resp.Roles[0].Role)
				assert.True(t, resp.Roles[0].Available)

				// Check player role availability
				assert.Equal(t, "player", resp.Roles[1].Role)
				assert.Equal(t, tt.expectIsPlayer, resp.Roles[1].Available)
			}
		})
	}
}

func TestRoleService_SwitchRole(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint64
		targetRole    string
		setupMock     func(*MockRoleUserRepository, *MockRolePlayerRepository)
		expectError   bool
		errorContains string
	}{
		{
			name:       "switch to user role",
			userID:     1,
			targetRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Test User", Status: model.UserStatusActive, Role: model.RoleUser}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectError: false,
		},
		{
			name:       "switch to player role - verified player",
			userID:     1,
			targetRole: "player",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Test User", Status: model.UserStatusActive, Role: model.RoleUser}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectError: false,
		},
		{
			name:       "switch to player role - not a player",
			userID:     2,
			targetRole: "player",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Regular User", Status: model.UserStatusActive, Role: model.RoleUser}
				user.ID = 2
				users.On("Get", mock.Anything, uint64(2)).Return(user, nil)

				players.On("GetByUserID", mock.Anything, uint64(2)).Return(nil, repository.ErrNotFound)
			},
			expectError:   true,
			errorContains: "您还不是认证陪玩师",
		},
		{
			name:       "switch to player role - pending verification",
			userID:     3,
			targetRole: "player",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Pending Player", Status: model.UserStatusActive, Role: model.RoleUser}
				user.ID = 3
				users.On("Get", mock.Anything, uint64(3)).Return(user, nil)

				player := &model.Player{
					UserID:             3,
					VerificationStatus: model.VerificationPending,
				}
				player.ID = 3
				players.On("GetByUserID", mock.Anything, uint64(3)).Return(player, nil)
			},
			expectError:   true,
			errorContains: "您还不是认证陪玩师",
		},
		{
			name:       "user not found",
			userID:     999,
			targetRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				users.On("Get", mock.Anything, uint64(999)).Return(nil, repository.ErrNotFound)
			},
			expectError:   true,
			errorContains: "用户不存在",
		},
		{
			name:       "user is banned",
			userID:     4,
			targetRole: "user",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Banned User", Status: model.UserStatusBanned, Role: model.RoleUser}
				user.ID = 4
				users.On("Get", mock.Anything, uint64(4)).Return(user, nil)
			},
			expectError:   true,
			errorContains: "用户已被禁用",
		},
		{
			name:       "invalid target role",
			userID:     1,
			targetRole: "admin",
			setupMock: func(users *MockRoleUserRepository, players *MockRolePlayerRepository) {
				user := &model.User{Name: "Test User", Status: model.UserStatusActive, Role: model.RoleUser}
				user.ID = 1
				users.On("Get", mock.Anything, uint64(1)).Return(user, nil)

				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectError:   true,
			errorContains: "无效的目标角色",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockRoleUserRepository{}
			players := &MockRolePlayerRepository{}

			tt.setupMock(users, players)

			svc := NewRoleService(users, players)
			resp, err := svc.SwitchRole(context.Background(), tt.userID, tt.targetRole)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.targetRole, resp.CurrentRole)
				assert.NotEmpty(t, resp.AccessToken)
			}
		})
	}
}

func TestRoleService_checkIsPlayer(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint64
		setupMock  func(*MockRolePlayerRepository)
		expectBool bool
	}{
		{
			name:   "verified player",
			userID: 1,
			setupMock: func(players *MockRolePlayerRepository) {
				player := &model.Player{
					UserID:             1,
					VerificationStatus: model.VerificationVerified,
				}
				player.ID = 1
				players.On("GetByUserID", mock.Anything, uint64(1)).Return(player, nil)
			},
			expectBool: true,
		},
		{
			name:   "pending player",
			userID: 2,
			setupMock: func(players *MockRolePlayerRepository) {
				player := &model.Player{
					UserID:             2,
					VerificationStatus: model.VerificationPending,
				}
				player.ID = 2
				players.On("GetByUserID", mock.Anything, uint64(2)).Return(player, nil)
			},
			expectBool: false,
		},
		{
			name:   "rejected player",
			userID: 3,
			setupMock: func(players *MockRolePlayerRepository) {
				player := &model.Player{
					UserID:             3,
					VerificationStatus: model.VerificationRejected,
				}
				player.ID = 3
				players.On("GetByUserID", mock.Anything, uint64(3)).Return(player, nil)
			},
			expectBool: false,
		},
		{
			name:   "not a player",
			userID: 4,
			setupMock: func(players *MockRolePlayerRepository) {
				players.On("GetByUserID", mock.Anything, uint64(4)).Return(nil, repository.ErrNotFound)
			},
			expectBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &MockRoleUserRepository{}
			players := &MockRolePlayerRepository{}

			tt.setupMock(players)

			svc := NewRoleService(users, players)
			result := svc.checkIsPlayer(context.Background(), tt.userID)

			assert.Equal(t, tt.expectBool, result)
		})
	}
}

func TestRoleService_checkIsPlayer_NilPlayers(t *testing.T) {
	users := &MockRoleUserRepository{}

	svc := NewRoleService(users, nil)
	result := svc.checkIsPlayer(context.Background(), 1)

	assert.False(t, result)
}

func TestRoleSwitchRequest_Structure(t *testing.T) {
	req := RoleSwitchRequest{
		TargetRole: "player",
	}

	assert.Equal(t, "player", req.TargetRole)
}

func TestRoleSwitchResponse_Structure(t *testing.T) {
	resp := RoleSwitchResponse{
		AccessToken: "test_token",
		CurrentRole: "player",
		IsPlayer:    true,
	}

	assert.Equal(t, "test_token", resp.AccessToken)
	assert.Equal(t, "player", resp.CurrentRole)
	assert.True(t, resp.IsPlayer)
}

func TestAvailableRole_Structure(t *testing.T) {
	role := AvailableRole{
		Role:        "player",
		Name:        "陪玩师",
		Description: "接单、管理服务、收益追踪",
		Available:   true,
	}

	assert.Equal(t, "player", role.Role)
	assert.Equal(t, "陪玩师", role.Name)
	assert.Equal(t, "接单、管理服务、收益追踪", role.Description)
	assert.True(t, role.Available)
}

func TestAvailableRolesResponse_Structure(t *testing.T) {
	resp := AvailableRolesResponse{
		CurrentRole: "user",
		Roles: []AvailableRole{
			{Role: "user", Name: "用户", Available: true},
			{Role: "player", Name: "陪玩师", Available: false},
		},
	}

	assert.Equal(t, "user", resp.CurrentRole)
	assert.Len(t, resp.Roles, 2)
}

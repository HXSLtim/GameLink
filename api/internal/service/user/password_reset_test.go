package user

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

// mockUserRepository is a mock implementation of UserRepository for testing
type mockUserRepository struct {
	users         map[uint64]*model.User
	findByEmailFn func(ctx context.Context, email string) (*model.User, error)
	updatePassFn  func(ctx context.Context, userID uint64, newPassword string) error
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepository) UpdatePassword(ctx context.Context, userID uint64, newPassword string) error {
	if m.updatePassFn != nil {
		return m.updatePassFn(ctx, userID, newPassword)
	}
	if u, ok := m.users[userID]; ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost+2)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hashedPassword)
		return nil
	}
	return repository.ErrNotFound
}

// Implement other required methods (stubs)
func (m *mockUserRepository) List(ctx context.Context) ([]model.User, error) { return nil, nil }
func (m *mockUserRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepository) ListWithFilters(ctx context.Context, opts repository.UserListOptions) ([]model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepository) Count(ctx context.Context, opts repository.UserListOptions) (int, error) {
	return 0, nil
}
func (m *mockUserRepository) Get(ctx context.Context, id uint64) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByIDs(ctx context.Context, ids []uint64) ([]model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepository) Delete(ctx context.Context, id uint64) error        { return nil }
func (m *mockUserRepository) GetByWeChatOpenID(ctx context.Context, openID string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByWeChatUnionID(ctx context.Context, unionID string) (*model.User, error) {
	return nil, nil
}

// setupTestRedis creates a miniredis instance for testing
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestPasswordResetService_RequestReset(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupMock     func(*mockUserRepository)
		expectError   bool
		expectSuccess bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {
						Email: "test@example.com",
						Name:  "Test User",
					},
				}
			},
			expectError:   false,
			expectSuccess: true,
		},
		{
			name:        "non-existent email (should still succeed to prevent user enumeration)",
			email:       "nonexistent@example.com",
			setupMock:   func(m *mockUserRepository) {},
			expectError: false,
			// Always returns success to prevent user enumeration
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, redisClient := setupTestRedis(t)
			defer mr.Close()

			mockRepo := &mockUserRepository{}
			tt.setupMock(mockRepo)

			// Create cache wrapper
			testCache := &testCache{client: redisClient}
			logger := slog.Default()
			service := NewPasswordResetService(mockRepo, testCache, logger)

			err := service.RequestReset(context.Background(), tt.email)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// testCache is a simple wrapper for testing
type testCache struct {
	client *redis.Client
}

func (c *testCache) GetClient() *redis.Client {
	return c.client
}

// Implement cache.Cache interface stubs
func (c *testCache) Get(ctx context.Context, key string) (string, bool, error) {
	return "", false, nil
}

func (c *testCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return nil
}

func (c *testCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *testCache) Close(context.Context) error {
	return nil
}

func TestPasswordResetService_ResetPassword(t *testing.T) {
	tests := []struct {
		name         string
		token        string
		newPassword  string
		setupMock    func(*mockUserRepository)
		setupToken   func(*miniredis.Miniredis, string)
		expectError  error
		verifyCalled bool
	}{
		{
			name:        "valid token and strong password",
			token:       "valid-token-123",
			newPassword: "NewSecure123!@#",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				mr.Set(passwordResetKeyPrefix+token, "1")
				mr.FastForward(14 * time.Minute) // Token still valid
			},
			expectError:  nil,
			verifyCalled: true,
		},
		{
			name:        "invalid token",
			token:       "invalid-token",
			newPassword: "NewSecure123!@#",
			setupMock:   func(m *mockUserRepository) {},
			setupToken:  func(mr *miniredis.Miniredis, token string) {},
			expectError: ErrInvalidToken,
		},
		{
			name:        "expired token",
			token:       "expired-token",
			newPassword: "NewSecure123!@#",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				// Don't set the token at all (simulating expired/non-existent)
			},
			expectError: ErrInvalidToken,
		},
		{
			name:        "weak password - too short",
			token:       "valid-token",
			newPassword: "Short1!",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				mr.Set(passwordResetKeyPrefix+token, "1")
			},
			expectError: ErrWeakPassword,
		},
		{
			name:        "weak password - no uppercase",
			token:       "valid-token",
			newPassword: "lowercase123!@#",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				mr.Set(passwordResetKeyPrefix+token, "1")
			},
			expectError: ErrWeakPassword,
		},
		{
			name:        "weak password - no special character",
			token:       "valid-token",
			newPassword: "NoSpecial123",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				mr.Set(passwordResetKeyPrefix+token, "1")
			},
			expectError: ErrWeakPassword,
		},
		{
			name:        "strong password with all requirements",
			token:       "valid-token",
			newPassword: "VeryStr0ng!Pass@2025",
			setupMock: func(m *mockUserRepository) {
				m.users = map[uint64]*model.User{
					1: {Email: "test@example.com"},
				}
			},
			setupToken: func(mr *miniredis.Miniredis, token string) {
				mr.Set(passwordResetKeyPrefix+token, "1")
			},
			expectError:  nil,
			verifyCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, redisClient := setupTestRedis(t)
			defer mr.Close()

			if tt.setupToken != nil {
				tt.setupToken(mr, tt.token)
			}

			mockRepo := &mockUserRepository{}
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			testCache := &testCache{client: redisClient}
			logger := slog.Default()
			service := NewPasswordResetService(mockRepo, testCache, logger)

			err := service.ResetPassword(context.Background(), tt.token, tt.newPassword)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectError) || err.Error() == tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify token is deleted after successful reset
			if err == nil {
				exists := mr.Exists(passwordResetKeyPrefix + tt.token)
				assert.False(t, exists, "reset token should be deleted after use")
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	service := &PasswordResetService{}

	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"valid password", "SecurePass123!@", false},
		{"too short", "Short1!", true},
		{"no uppercase", "lowercase123!@", true},
		{"no lowercase", "UPPERCASE123!@", true},
		{"no number", "NoNumbers!@#", true},
		{"no special", "NoSpecial123", true},
		{"common password", "Password123!", true},
		{"very long but valid", "ThisIsAVeryLongPassword123!@#ButStillSecure", false},
		{"too long (>128 chars)", string(make([]byte, 129)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePasswordStrength(tt.password)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateSecureToken(t *testing.T) {
	// Test that tokens are unique
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := generateSecureToken(32)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.False(t, tokens[token], "token should be unique")
		tokens[token] = true
	}

	// Test token length
	token, err := generateSecureToken(32)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(token), 32, "token should be at least 32 characters")
}

func TestPasswordResetService_ResetPassword_SingleUseToken(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	mockRepo := &mockUserRepository{
		users: map[uint64]*model.User{
			1: {Email: "test@example.com"},
		},
	}

	token := "single-use-token"
	mr.Set(passwordResetKeyPrefix+token, "1")

	testCache := &testCache{client: redisClient}
	logger := slog.Default()
	service := NewPasswordResetService(mockRepo, testCache, logger)

	// First use should succeed
	err := service.ResetPassword(context.Background(), token, "NewSecure123!@#")
	assert.NoError(t, err)

	// Second use should fail (token deleted)
	err = service.ResetPassword(context.Background(), token, "AnotherSecure123!@#")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidToken))
}

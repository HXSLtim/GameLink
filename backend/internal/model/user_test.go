package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoleConstants 验证角色常量
func TestRoleConstants(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{"User role", RoleUser, "user"},
		{"Player role", RolePlayer, "player"},
		{"Admin role", RoleAdmin, "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.role))
		})
	}
}

// TestUserStatusConstants 验证用户状态常量
func TestUserStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected string
	}{
		{"Active status", UserStatusActive, "active"},
		{"Suspended status", UserStatusSuspended, "suspended"},
		{"Banned status", UserStatusBanned, "banned"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

// TestUserCreation 测试创建用户结构体
func TestUserCreation(t *testing.T) {
	now := time.Now()

	t.Run("Create minimal user", func(t *testing.T) {
		user := &User{
			Phone:  "13800138000",
			Email:  "test@example.com",
			Name:   "Test User",
			Role:   RoleUser,
			Status: UserStatusActive,
		}

		assert.Equal(t, "13800138000", user.Phone)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, RoleUser, user.Role)
		assert.Equal(t, UserStatusActive, user.Status)
	})

	t.Run("Create full user with all fields", func(t *testing.T) {
		lastLogin := time.Now()
		user := &User{
			Base: Base{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Phone:        "13800138001",
			Email:        "full@example.com",
			PasswordHash: "hashed_password_12345",
			Name:         "Full Test User",
			AvatarURL:    "https://example.com/avatar.jpg",
			Role:         RoleAdmin,
			Status:       UserStatusActive,
			LastLoginAt:  &lastLogin,
		}

		// Validate Base fields
		assert.Equal(t, uint64(1), user.ID)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)

		// Validate User fields
		assert.Equal(t, "13800138001", user.Phone)
		assert.Equal(t, "full@example.com", user.Email)
		assert.Equal(t, "hashed_password_12345", user.PasswordHash)
		assert.Equal(t, "Full Test User", user.Name)
		assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
		assert.Equal(t, RoleAdmin, user.Role)
		assert.Equal(t, UserStatusActive, user.Status)
		assert.Equal(t, &lastLogin, user.LastLoginAt)
	})

	t.Run("Create user with different roles", func(t *testing.T) {
		roles := []Role{RoleUser, RolePlayer, RoleAdmin}

		for _, role := range roles {
			user := &User{
				Email:  "role@example.com",
				Name:   "Role Test User",
				Role:   role,
				Status: UserStatusActive,
			}
			assert.Equal(t, role, user.Role)
		}
	})

	t.Run("Create user with different statuses", func(t *testing.T) {
		statuses := []UserStatus{UserStatusActive, UserStatusSuspended, UserStatusBanned}

		for _, status := range statuses {
			user := &User{
				Email:  "status@example.com",
				Name:   "Status Test User",
				Role:   RoleUser,
				Status: status,
			}
			assert.Equal(t, status, user.Status)
		}
	})
}

// TestUserJSONSerialization 测试 JSON 序列化和反序列化
func TestUserJSONSerialization(t *testing.T) {
	lastLogin := time.Now()

	t.Run("Marshal full user to JSON", func(t *testing.T) {
		user := &User{
			Base: Base{
				ID:        123,
				CreatedAt: lastLogin.Add(-24 * time.Hour),
				UpdatedAt: lastLogin.Add(-1 * time.Hour),
			},
			Phone:        "13800138001",
			Email:        "json@example.com",
			PasswordHash: "hashed_secret_67890",
			Name:         "JSON User",
			AvatarURL:    "https://cdn.example.com/avatar.png",
			Role:         RolePlayer,
			Status:       UserStatusActive,
			LastLoginAt:  &lastLogin,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)
		assert.NotEmpty(t, data)

		// Verify key fields are present in JSON
		jsonStr := string(data)
		assert.Contains(t, jsonStr, "json@example.com")
		assert.Contains(t, jsonStr, "JSON User")
		assert.Contains(t, jsonStr, "https://cdn.example.com/avatar.png")
		assert.Contains(t, jsonStr, "player")
		assert.Contains(t, jsonStr, "active")

		// PasswordHash should not be in JSON (json:"-")
		assert.NotContains(t, jsonStr, "hashed_secret_67890")
		assert.NotContains(t, jsonStr, "PasswordHash")
		assert.NotContains(t, jsonStr, "password_hash")
	})

	t.Run("Marshal minimal user", func(t *testing.T) {
		user := &User{
			Email:  "minimal@example.com",
			Name:   "Minimal User",
			Role:   RoleUser,
			Status: UserStatusSuspended,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		// Optional fields should be omitted
		assert.NotContains(t, result, "phone")
		assert.NotContains(t, result, "passwordHash")
		assert.NotContains(t, result, "password_hash")
		assert.NotContains(t, result, "lastLoginAt")
	})

	t.Run("Unmarshal user from JSON", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"email": "unmarshal@example.com",
			"name": "Unmarshal User",
			"avatarUrl": "https://example.com/avatar.jpg",
			"role": "admin",
			"status": "active",
			"lastLoginAt": "2023-12-01T10:30:00Z"
		}`

		var user User
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, uint64(456), user.ID)
		assert.Equal(t, "unmarshal@example.com", user.Email)
		assert.Equal(t, "Unmarshal User", user.Name)
		assert.Equal(t, "https://example.com/avatar.jpg", user.AvatarURL)
		assert.Equal(t, RoleAdmin, user.Role)
		assert.Equal(t, UserStatusActive, user.Status)
		assert.NotNil(t, user.LastLoginAt)
	})
}

// TestUserFieldAccess 测试访问所有字段
func TestUserFieldAccess(t *testing.T) {
	user := &User{}

	// Test setting and getting all fields
	t.Run("Set and get all User fields", func(t *testing.T) {
		now := time.Now()
		lastLogin := now.Add(-2 * time.Hour)

		// Set Base fields
		user.ID = 999
		user.CreatedAt = now
		user.UpdatedAt = now

		// Set User fields
		user.Phone = "13900139000"
		user.Email = "field@example.com"
		user.PasswordHash = "secure_hash_999"
		user.Name = "Field Access Test"
		user.AvatarURL = "https://example.com/field.jpg"
		user.Role = RolePlayer
		user.Status = UserStatusActive
		user.LastLoginAt = &lastLogin

		// Verify all fields
		assert.Equal(t, uint64(999), user.ID)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.Equal(t, "13900139000", user.Phone)
		assert.Equal(t, "field@example.com", user.Email)
		assert.Equal(t, "secure_hash_999", user.PasswordHash)
		assert.Equal(t, "Field Access Test", user.Name)
		assert.Equal(t, "https://example.com/field.jpg", user.AvatarURL)
		assert.Equal(t, RolePlayer, user.Role)
		assert.Equal(t, UserStatusActive, user.Status)
		assert.Equal(t, &lastLogin, user.LastLoginAt)
	})

	t.Run("Access nil LastLoginAt", func(t *testing.T) {
		user := &User{
			Email: "nil@example.com",
		}
		assert.Nil(t, user.LastLoginAt)
	})
}

// TestUserFieldInitialization 测试所有字段的初始化
func TestUserFieldInitialization(t *testing.T) {
	now := time.Now()

	t.Run("Initialize with all pointer fields", func(t *testing.T) {
		user := &User{
			Base: Base{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Phone:        "test",
			Email:        "test@test.com",
			PasswordHash: "hash123",
			Name:         "Test User",
			AvatarURL:    "avatar.jpg",
			Role:         RoleAdmin,
			Status:       UserStatusActive,
			LastLoginAt:  &now,
		}

		// Ensure all fields are correctly set
		assert.NotNil(t, user)
		assert.Equal(t, uint64(1), user.ID)
		assert.NotNil(t, &user.CreatedAt)
		assert.NotNil(t, &user.UpdatedAt)
		assert.Equal(t, "test", user.Phone)
		assert.Equal(t, "test@test.com", user.Email)
		assert.Equal(t, "hash123", user.PasswordHash)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "avatar.jpg", user.AvatarURL)
		assert.Equal(t, RoleAdmin, user.Role)
		assert.Equal(t, UserStatusActive, user.Status)
		assert.NotNil(t, user.LastLoginAt)
	})
}

// TestUserRoleValidation 测试用户角色值的合法性
func TestUserRoleValidation(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{"Valid user role", RoleUser, true},
		{"Valid player role", RolePlayer, true},
		{"Valid admin role", RoleAdmin, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Email:  "role@example.com",
				Role:   tt.role,
				Status: UserStatusActive,
			}
			assert.Equal(t, tt.role, user.Role)
			assert.NotEmpty(t, string(user.Role))
		})
	}
}

// TestUserStatusValidation 测试用户状态值的合法性
func TestUserStatusValidation(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{"Valid active status", UserStatusActive, true},
		{"Valid suspended status", UserStatusSuspended, true},
		{"Valid banned status", UserStatusBanned, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Email:  "status@example.com",
				Role:   RoleUser,
				Status: tt.status,
			}
			assert.Equal(t, tt.status, user.Status)
			assert.NotEmpty(t, string(user.Status))
		})
	}
}

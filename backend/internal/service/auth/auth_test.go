package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gamelink/internal/model"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "有效邮箱-标准格式",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱-包含点",
			email:    "user.name@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱-包含加号",
			email:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "有效邮箱-子域名",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "无效邮箱-空字符串",
			email:    "",
			expected: false,
		},
		{
			name:     "无效邮箱-过长",
			email:    "user" + string(make([]byte, 120)) + "@example.com",
			expected: false,
		},
		{
			name:     "无效邮箱-缺少@",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "无效邮箱-缺少域名",
			email:    "user@",
			expected: false,
		},
		{
			name:     "无效邮箱-缺少用户名",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "无效邮箱-特殊字符",
			email:    "user@exa mple.com",
			expected: false,
		},
		{
			name:     "无效邮箱-多个@",
			email:    "user@@example.com",
			expected: false,
		},
		{
			name:     "临时邮箱-tempmail.com",
			email:    "user@tempmail.com",
			expected: false,
		},
		{
			name:     "临时邮箱-10minutemail.com",
			email:    "user@10minutemail.com",
			expected: false,
		},
		{
			name:     "临时邮箱子域名",
			email:    "user@sub.tempmail.com",
			expected: false,
		},
		{
			name:     "有效邮箱-类似临时域名",
			email:    "user@mytempmail.com",
			expected: true,
		},
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
		name        string
		request     RegisterRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效注册-邮箱",
			request: RegisterRequest{
				Email:    "user@example.com",
				Password: "password123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: false,
		},
		{
			name: "有效注册-手机号",
			request: RegisterRequest{
				Phone:    "13812345678",
				Password: "password123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: false,
		},
		{
			name: "有效注册-邮箱和手机号",
			request: RegisterRequest{
				Email:    "user@example.com",
				Phone:    "13812345678",
				Password: "password123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: false,
		},
		{
			name: "无效注册-缺少邮箱和手机号",
			request: RegisterRequest{
				Password: "password123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: true,
			errorMsg:    "email or phone is required",
		},
		{
			name: "无效注册-缺少姓名",
			request: RegisterRequest{
				Email:    "user@example.com",
				Password: "password123",
				Role:     model.RoleUser,
			},
			expectError: true,
			errorMsg:    "name is required",
		},
		{
			name: "无效注册-缺少密码",
			request: RegisterRequest{
				Email: "user@example.com",
				Name:  "Test User",
				Role:  model.RoleUser,
			},
			expectError: true,
			errorMsg:    "password is required",
		},
		{
			name: "无效注册-密码太短",
			request: RegisterRequest{
				Email:    "user@example.com",
				Password: "123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: true,
			errorMsg:    "password must be at least 6 characters",
		},
		{
			name: "无效注册-无效邮箱格式",
			request: RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
				Role:     model.RoleUser,
			},
			expectError: true,
			errorMsg:    "email or phone is required",
		},
		{
			name: "有效注册-默认角色",
			request: RegisterRequest{
				Email:    "user@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegisterInput(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

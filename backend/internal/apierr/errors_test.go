package apierr

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		message     string
		expectCode  int
		expectMsg   string
	}{
		{
			name:       "创建400错误",
			code:       http.StatusBadRequest,
			message:    "请求参数无效",
			expectCode: http.StatusBadRequest,
			expectMsg:  "请求参数无效",
		},
		{
			name:       "创建404错误",
			code:       http.StatusNotFound,
			message:    "资源不存在",
			expectCode: http.StatusNotFound,
			expectMsg:  "资源不存在",
		},
		{
			name:       "创建500错误",
			code:       http.StatusInternalServerError,
			message:    "服务器内部错误",
			expectCode: http.StatusInternalServerError,
			expectMsg:  "服务器内部错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)

			assert.NotNil(t, err)
			assert.Equal(t, tt.expectCode, err.Code)
			assert.Equal(t, tt.expectMsg, err.Message)
			assert.NotZero(t, err.Timestamp)
			assert.Contains(t, err.Error(), tt.message)
		})
	}
}

func TestAPIError_WithDetails(t *testing.T) {
	err := New(http.StatusBadRequest, "验证失败").
		WithDetails("用户名长度必须在3-20个字符之间")

	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "验证失败", err.Message)
	assert.Equal(t, "用户名长度必须在3-20个字符之间", err.Details)
}

func TestAPIError_WithField(t *testing.T) {
	err := New(http.StatusBadRequest, "验证失败").
		WithField("username").
		WithDetails("用户名长度必须在3-20个字符之间")

	assert.Equal(t, "username", err.Field)
	assert.Equal(t, "验证失败", err.Message)
}

func TestAPIError_WithRequestID(t *testing.T) {
	err := New(http.StatusInternalServerError, "服务器内部错误").
		WithRequestID("req-123456")

	assert.Equal(t, "req-123456", err.RequestID)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
}

func TestAPIError_WithExtension(t *testing.T) {
	err := New(http.StatusBadRequest, "验证失败").
		WithExtension("min_length", 3).
		WithExtension("max_length", 20).
		WithExtension("actual_length", 25)

	assert.NotNil(t, err.Extensions)
	assert.Equal(t, 3, err.Extensions["min_length"])
	assert.Equal(t, 20, err.Extensions["max_length"])
	assert.Equal(t, 25, err.Extensions["actual_length"])
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("请求参数无效")

	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "请求参数无效", err.Message)
	assert.NotZero(t, err.Timestamp)
}

func TestUnauthorized(t *testing.T) {
	err := Unauthorized("未授权访问")

	assert.Equal(t, http.StatusUnauthorized, err.Code)
	assert.Equal(t, "未授权访问", err.Message)
}

func TestForbidden(t *testing.T) {
	err := Forbidden("权限不足")

	assert.Equal(t, http.StatusForbidden, err.Code)
	assert.Equal(t, "权限不足", err.Message)
}

func TestNotFound(t *testing.T) {
	err := NotFound("资源不存在")

	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, "资源不存在", err.Message)
}

func TestConflict(t *testing.T) {
	err := Conflict("资源冲突")

	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Equal(t, "资源冲突", err.Message)
}

func TestInternalError(t *testing.T) {
	err := InternalError("服务器内部错误")

	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "服务器内部错误", err.Message)
}

func TestServiceUnavailable(t *testing.T) {
	err := ServiceUnavailable("服务不可用")

	assert.Equal(t, http.StatusServiceUnavailable, err.Code)
	assert.Equal(t, "服务不可用", err.Message)
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("username", "用户名长度必须在3-20个字符之间")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "用户名长度必须在3-20个字符之间", err.Message)
	assert.Equal(t, "username", err.Field)
	assert.NotZero(t, err.Timestamp)
}

func TestValidationError_WithValue(t *testing.T) {
	err := NewValidationError("username", "用户名长度必须在3-20个字符之间")
	err.Value = "ab" // 2个字符，不满足最小长度
	err.Tag = "min"

	assert.Equal(t, "ab", err.Value)
	assert.Equal(t, "min", err.Tag)
}

func TestNewDatabaseError(t *testing.T) {
	dbErr := errors.New("connection timeout")
	err := NewDatabaseError(dbErr)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "数据库操作失败", err.Message)
	assert.Equal(t, "connection timeout", err.Details)
	assert.NotZero(t, err.Timestamp)
}

func TestNewDatabaseError_WithQuery(t *testing.T) {
	dbErr := errors.New("duplicate key value")
	err := NewDatabaseError(dbErr)
	err.Query = "INSERT INTO users (email) VALUES ('test@example.com')"

	assert.Equal(t, "INSERT INTO users (email) VALUES ('test@example.com')", err.Query)
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "404错误",
			err:      NotFound("资源不存在"),
			expected: true,
		},
		{
			name:     "400错误",
			err:      BadRequest("请求无效"),
			expected: false,
		},
		{
			name:     "500错误",
			err:      InternalError("服务器错误"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: false,
		},
		{
			name:     "非APIError",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "401错误",
			err:      Unauthorized("未授权"),
			expected: true,
		},
		{
			name:     "403错误",
			err:      Forbidden("权限不足"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnauthorized(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "403错误",
			err:      Forbidden("权限不足"),
			expected: true,
		},
		{
			name:     "401错误",
			err:      Unauthorized("未授权"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsForbidden(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "400错误",
			err:      BadRequest("请求无效"),
			expected: true,
		},
		{
			name:     "验证错误",
			err:      NewValidationError("username", "用户名无效"),
			expected: true,
		},
		{
			name:     "404错误",
			err:      NotFound("资源不存在"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidationError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAPIError_ErrorMethod(t *testing.T) {
	tests := []struct {
		name          string
		err           *APIError
		expectedError string
	}{
		{
			name:          "400错误",
			err:           BadRequest("请求无效"),
			expectedError: "[400] 请求无效",
		},
		{
			name:          "404错误",
			err:           NotFound("资源不存在"),
			expectedError: "[404] 资源不存在",
		},
		{
			name:          "500错误",
			err:           InternalError("服务器错误"),
			expectedError: "[500] 服务器错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedError, tt.err.Error())
		})
	}
}

func TestAPIError_Timestamp(t *testing.T) {
	before := time.Now().Unix()
	err := BadRequest("请求无效")
	after := time.Now().Unix()

	assert.True(t, err.Timestamp >= before)
	assert.True(t, err.Timestamp <= after)
}

func TestCommonErrors(t *testing.T) {
	// 验证预定义的错误
	assert.NotNil(t, ErrNotFound)
	assert.Equal(t, http.StatusNotFound, ErrNotFound.Code)

	assert.NotNil(t, ErrUnauthorized)
	assert.Equal(t, http.StatusUnauthorized, ErrUnauthorized.Code)

	assert.NotNil(t, ErrForbidden)
	assert.Equal(t, http.StatusForbidden, ErrForbidden.Code)

	assert.NotNil(t, ErrInvalidInput)
	assert.Equal(t, http.StatusBadRequest, ErrInvalidInput.Code)

	assert.NotNil(t, ErrInternal)
	assert.Equal(t, http.StatusInternalServerError, ErrInternal.Code)

	assert.NotNil(t, ErrConflict)
	assert.Equal(t, http.StatusConflict, ErrConflict.Code)

	assert.NotNil(t, ErrTooManyRequests)
	assert.Equal(t, http.StatusTooManyRequests, ErrTooManyRequests.Code)
}

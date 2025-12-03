// Package apierr provides generic API error handling.

package apierr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/**
 * Test_New_CreateAPIError 测试 New 函数创建 API 错误
 */
func Test_New_CreateAPIError(t *testing.T) {
	// Arrange
	code := http.StatusBadRequest
	message := "测试错误消息"
	beforeTime := time.Now().Unix()

	// Act
	err := New(code, message)
	afterTime := time.Now().Unix()

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.GreaterOrEqual(t, err.Timestamp, beforeTime)
	assert.LessOrEqual(t, err.Timestamp, afterTime)
	assert.Empty(t, err.Details)
	assert.Empty(t, err.Field)
	assert.Empty(t, err.RequestID)
	assert.Nil(t, err.Extensions)
}

/**
 * Test_APIError_Error_ReturnsFormattedString 测试 APIError.Error 方法
 */
func Test_APIError_Error_ReturnsFormattedString(t *testing.T) {
	// Arrange
	code := http.StatusNotFound
	message := "资源未找到"
	err := New(code, message)

	// Act
	result := err.Error()

	// Assert
	expected := fmt.Sprintf("[%d] %s", code, message)
	assert.Equal(t, expected, result)
}

/**
 * Test_APIError_WithDetails_AddsDetails 测试链式添加详情
 */
func Test_APIError_WithDetails_AddsDetails(t *testing.T) {
	// Arrange
	err := New(http.StatusBadRequest, "测试错误")
	details := "详细错误信息"

	// Act
	result := err.WithDetails(details)

	// Assert
	assert.Equal(t, details, result.Details)
	assert.Equal(t, err, result) // 确保返回同一个对象
}

/**
 * Test_APIError_WithField_AddsField 测试添加字段信息
 */
func Test_APIError_WithField_AddsField(t *testing.T) {
	// Arrange
	err := New(http.StatusBadRequest, "验证失败")
	field := "username"

	// Act
	result := err.WithField(field)

	// Assert
	assert.Equal(t, field, result.Field)
	assert.Equal(t, err, result)
}

/**
 * Test_APIError_WithRequestID_AddsRequestID 测试添加请求ID
 */
func Test_APIError_WithRequestID_AddsRequestID(t *testing.T) {
	// Arrange
	err := New(http.StatusInternalServerError, "服务器错误")
	requestID := "req-123456"

	// Act
	result := err.WithRequestID(requestID)

	// Assert
	assert.Equal(t, requestID, result.RequestID)
	assert.Equal(t, err, result)
}

/**
 * Test_APIError_WithExtension_AddsExtension 测试添加扩展数据
 */
func Test_APIError_WithExtension_AddsExtension(t *testing.T) {
	// Arrange
	err := New(http.StatusBadRequest, "测试")
	key := "custom_field"
	value := "custom_value"

	// Act
	result := err.WithExtension(key, value)

	// Assert
	require.NotNil(t, result.Extensions)
	assert.Equal(t, value, result.Extensions[key])
	assert.Equal(t, err, result)
}

/**
 * Test_APIError_WithExtension_MultipleExtensions 测试添加多个扩展
 */
func Test_APIError_WithExtension_MultipleExtensions(t *testing.T) {
	// Arrange
	err := New(http.StatusBadRequest, "测试")

	// Act
	result := err.
		WithExtension("key1", "value1").
		WithExtension("key2", 123).
		WithExtension("key3", true)

	// Assert
	require.NotNil(t, result.Extensions)
	assert.Equal(t, "value1", result.Extensions["key1"])
	assert.Equal(t, 123, result.Extensions["key2"])
	assert.Equal(t, true, result.Extensions["key3"])
}

/**
 * Test_APIError_ChainedMethods 测试链式调用
 */
func Test_APIError_ChainedMethods(t *testing.T) {
	// Act
	err := New(http.StatusBadRequest, "验证失败").
		WithDetails("用户名格式不正确").
		WithField("username").
		WithRequestID("req-789").
		WithExtension("attempts", 3)

	// Assert
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "验证失败", err.Message)
	assert.Equal(t, "用户名格式不正确", err.Details)
	assert.Equal(t, "username", err.Field)
	assert.Equal(t, "req-789", err.RequestID)
	assert.Equal(t, 3, err.Extensions["attempts"])
}

/**
 * Test_BadRequest_Creates400Error 测试 BadRequest 函数
 */
func Test_BadRequest_Creates400Error(t *testing.T) {
	// Arrange
	message := "无效的请求"

	// Act
	err := BadRequest(message)

	// Assert
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_Unauthorized_Creates401Error 测试 Unauthorized 函数
 */
func Test_Unauthorized_Creates401Error(t *testing.T) {
	// Arrange
	message := "未授权"

	// Act
	err := Unauthorized(message)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_Forbidden_Creates403Error 测试 Forbidden 函数
 */
func Test_Forbidden_Creates403Error(t *testing.T) {
	// Arrange
	message := "禁止访问"

	// Act
	err := Forbidden(message)

	// Assert
	assert.Equal(t, http.StatusForbidden, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_NotFound_Creates404Error 测试 NotFound 函数
 */
func Test_NotFound_Creates404Error(t *testing.T) {
	// Arrange
	message := "资源不存在"

	// Act
	err := NotFound(message)

	// Assert
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_Conflict_Creates409Error 测试 Conflict 函数
 */
func Test_Conflict_Creates409Error(t *testing.T) {
	// Arrange
	message := "资源冲突"

	// Act
	err := Conflict(message)

	// Assert
	assert.Equal(t, http.StatusConflict, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_InternalError_Creates500Error 测试 InternalError 函数
 */
func Test_InternalError_Creates500Error(t *testing.T) {
	// Arrange
	message := "内部错误"

	// Act
	err := InternalError(message)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_ServiceUnavailable_Creates503Error 测试 ServiceUnavailable 函数
 */
func Test_ServiceUnavailable_Creates503Error(t *testing.T) {
	// Arrange
	message := "服务不可用"

	// Act
	err := ServiceUnavailable(message)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_TooManyRequests_Creates429Error 测试 TooManyRequests 函数
 */
func Test_TooManyRequests_Creates429Error(t *testing.T) {
	// Arrange
	message := "请求过多"

	// Act
	err := TooManyRequests(message)

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, err.Code)
	assert.Equal(t, message, err.Message)
}

/**
 * Test_NewValidationError_CreatesValidationError 测试创建验证错误
 */
func Test_NewValidationError_CreatesValidationError(t *testing.T) {
	// Arrange
	field := "email"
	message := "邮箱格式不正确"
	beforeTime := time.Now().Unix()

	// Act
	err := NewValidationError(field, message)
	afterTime := time.Now().Unix()

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, field, err.Field)
	assert.GreaterOrEqual(t, err.Timestamp, beforeTime)
	assert.LessOrEqual(t, err.Timestamp, afterTime)
}

/**
 * Test_NewDatabaseError_CreatesDatabaseError 测试创建数据库错误
 */
func Test_NewDatabaseError_CreatesDatabaseError(t *testing.T) {
	// Arrange
	originalErr := errors.New("connection timeout")
	beforeTime := time.Now().Unix()

	// Act
	err := NewDatabaseError(originalErr)
	afterTime := time.Now().Unix()

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "数据库操作失败", err.Message)
	assert.Equal(t, "connection timeout", err.Details)
	assert.GreaterOrEqual(t, err.Timestamp, beforeTime)
	assert.LessOrEqual(t, err.Timestamp, afterTime)
}

/**
 * Test_IsNotFound_WithNotFoundError 测试 IsNotFound 检测 404 错误
 */
func Test_IsNotFound_WithNotFoundError(t *testing.T) {
	// Arrange
	err := NotFound("资源不存在")

	// Act
	result := IsNotFound(err)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsNotFound_WithOtherError 测试 IsNotFound 对其他错误返回 false
 */
func Test_IsNotFound_WithOtherError(t *testing.T) {
	// Arrange
	err := BadRequest("错误请求")

	// Act
	result := IsNotFound(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsNotFound_WithNilError 测试 IsNotFound 对 nil 返回 false
 */
func Test_IsNotFound_WithNilError(t *testing.T) {
	// Act
	result := IsNotFound(nil)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsNotFound_WithWrappedError 测试 IsNotFound 检测包装的错误
 */
func Test_IsNotFound_WithWrappedError(t *testing.T) {
	// Arrange
	notFoundErr := NotFound("资源不存在")
	wrappedErr := fmt.Errorf("wrapped: %w", notFoundErr)

	// Act
	result := IsNotFound(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsUnauthorized_WithUnauthorizedError 测试 IsUnauthorized
 */
func Test_IsUnauthorized_WithUnauthorizedError(t *testing.T) {
	// Arrange
	err := Unauthorized("未授权")

	// Act
	result := IsUnauthorized(err)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsUnauthorized_WithOtherError 测试 IsUnauthorized 对其他错误
 */
func Test_IsUnauthorized_WithOtherError(t *testing.T) {
	// Arrange
	err := BadRequest("错误请求")

	// Act
	result := IsUnauthorized(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsUnauthorized_WithNilError 测试 IsUnauthorized 对 nil
 */
func Test_IsUnauthorized_WithNilError(t *testing.T) {
	// Act
	result := IsUnauthorized(nil)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsForbidden_WithForbiddenError 测试 IsForbidden
 */
func Test_IsForbidden_WithForbiddenError(t *testing.T) {
	// Arrange
	err := Forbidden("禁止访问")

	// Act
	result := IsForbidden(err)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsForbidden_WithOtherError 测试 IsForbidden 对其他错误
 */
func Test_IsForbidden_WithOtherError(t *testing.T) {
	// Arrange
	err := BadRequest("错误请求")

	// Act
	result := IsForbidden(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsForbidden_WithNilError 测试 IsForbidden 对 nil
 */
func Test_IsForbidden_WithNilError(t *testing.T) {
	// Act
	result := IsForbidden(nil)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsValidationError_WithValidationError 测试 IsValidationError
 */
func Test_IsValidationError_WithValidationError(t *testing.T) {
	// Arrange
	err := NewValidationError("email", "邮箱格式错误")

	// Act
	result := IsValidationError(err)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsValidationError_WithBadRequestError 测试 IsValidationError 对 BadRequest
 */
func Test_IsValidationError_WithBadRequestError(t *testing.T) {
	// Arrange
	err := BadRequest("错误请求")

	// Act
	result := IsValidationError(err)

	// Assert
	assert.True(t, result) // BadRequest 也应该被认为是验证错误
}

/**
 * Test_IsValidationError_WithOtherError 测试 IsValidationError 对其他错误
 */
func Test_IsValidationError_WithOtherError(t *testing.T) {
	// Arrange
	err := NotFound("未找到")

	// Act
	result := IsValidationError(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsValidationError_WithNilError 测试 IsValidationError 对 nil
 */
func Test_IsValidationError_WithNilError(t *testing.T) {
	// Act
	result := IsValidationError(nil)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsBadRequest_WithBadRequestError 测试 IsBadRequest
 */
func Test_IsBadRequest_WithBadRequestError(t *testing.T) {
	// Arrange
	err := BadRequest("错误请求")

	// Act
	result := IsBadRequest(err)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsBadRequest_WithOtherError 测试 IsBadRequest 对其他错误
 */
func Test_IsBadRequest_WithOtherError(t *testing.T) {
	// Arrange
	err := NotFound("未找到")

	// Act
	result := IsBadRequest(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsBadRequest_WithNilError 测试 IsBadRequest 对 nil
 */
func Test_IsBadRequest_WithNilError(t *testing.T) {
	// Act
	result := IsBadRequest(nil)

	// Assert
	assert.False(t, result)
}

/**
 * Test_PredefinedErrors_AreCorrectlyInitialized 测试预定义错误
 */
func Test_PredefinedErrors_AreCorrectlyInitialized(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		wantCode int
		wantMsg  string
	}{
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			wantCode: http.StatusNotFound,
			wantMsg:  "资源不存在",
		},
		{
			name:     "ErrUnauthorized",
			err:      ErrUnauthorized,
			wantCode: http.StatusUnauthorized,
			wantMsg:  "未授权访问",
		},
		{
			name:     "ErrForbidden",
			err:      ErrForbidden,
			wantCode: http.StatusForbidden,
			wantMsg:  "权限不足",
		},
		{
			name:     "ErrInvalidInput",
			err:      ErrInvalidInput,
			wantCode: http.StatusBadRequest,
			wantMsg:  "输入参数无效",
		},
		{
			name:     "ErrInternal",
			err:      ErrInternal,
			wantCode: http.StatusInternalServerError,
			wantMsg:  "服务器内部错误",
		},
		{
			name:     "ErrConflict",
			err:      ErrConflict,
			wantCode: http.StatusConflict,
			wantMsg:  "资源冲突",
		},
		{
			name:     "ErrTooManyRequests",
			err:      ErrTooManyRequests,
			wantCode: http.StatusTooManyRequests,
			wantMsg:  "请求过于频繁",
		},
		{
			name:     "ErrValidationFailed",
			err:      ErrValidationFailed,
			wantCode: http.StatusBadRequest,
			wantMsg:  "验证失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.err)
			assert.Equal(t, tt.wantCode, tt.err.Code)
			assert.Equal(t, tt.wantMsg, tt.err.Message)
		})
	}
}

/**
 * Test_ErrorMessages_Constants 测试错误消息常量
 */
func Test_ErrorMessages_Constants(t *testing.T) {
	// 确保常量被正确定义
	assert.Equal(t, "invalid json payload", ErrInvalidJSONPayload)
	assert.Equal(t, "invalid id", ErrInvalidID)
	assert.Equal(t, "invalid page", ErrInvalidPage)
	assert.Equal(t, "invalid page_size", ErrInvalidPageSize)
	assert.Equal(t, "invalid user_id", ErrInvalidUserID)
	assert.Equal(t, "invalid order_id", ErrInvalidOrderID)
	assert.Equal(t, "game not found", ErrGameNotFound)
	assert.Equal(t, "player not found", ErrPlayerNotFound)
	assert.Equal(t, "invalid player_id", ErrInvalidPlayerID)
	assert.Equal(t, "invalid game_id", ErrInvalidGameID)
	assert.Equal(t, "invalid email format", ErrInvalidEmailFormat)
	assert.Equal(t, "invalid phone format", ErrInvalidPhoneFormat)
	assert.Equal(t, "user not found", ErrUserNotFound)
	assert.Equal(t, "order not found", ErrOrderNotFound)
	assert.Equal(t, "payment not found", ErrPaymentNotFound)
}

/**
 * Test_IsUnauthorized_WithWrappedError 测试 IsUnauthorized 检测包装的错误
 */
func Test_IsUnauthorized_WithWrappedError(t *testing.T) {
	// Arrange
	unauthorizedErr := Unauthorized("未授权")
	wrappedErr := fmt.Errorf("wrapped: %w", unauthorizedErr)

	// Act
	result := IsUnauthorized(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsForbidden_WithWrappedError 测试 IsForbidden 检测包装的错误
 */
func Test_IsForbidden_WithWrappedError(t *testing.T) {
	// Arrange
	forbiddenErr := Forbidden("禁止访问")
	wrappedErr := fmt.Errorf("wrapped: %w", forbiddenErr)

	// Act
	result := IsForbidden(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsValidationError_WithWrappedValidationError 测试 IsValidationError 检测包装的验证错误
 */
func Test_IsValidationError_WithWrappedValidationError(t *testing.T) {
	// Arrange
	validationErr := NewValidationError("email", "邮箱格式错误")
	wrappedErr := fmt.Errorf("wrapped: %w", validationErr)

	// Act
	result := IsValidationError(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsValidationError_WithWrappedAPIError 测试 IsValidationError 检测包装的API错误
 */
func Test_IsValidationError_WithWrappedAPIError(t *testing.T) {
	// Arrange
	apiErr := BadRequest("错误请求")
	wrappedErr := fmt.Errorf("wrapped: %w", apiErr)

	// Act
	result := IsValidationError(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsBadRequest_WithWrappedError 测试 IsBadRequest 检测包装的错误
 */
func Test_IsBadRequest_WithWrappedError(t *testing.T) {
	// Arrange
	badRequestErr := BadRequest("错误请求")
	wrappedErr := fmt.Errorf("wrapped: %w", badRequestErr)

	// Act
	result := IsBadRequest(wrappedErr)

	// Assert
	assert.True(t, result)
}

/**
 * Test_IsNotFound_WithNonAPIError 测试 IsNotFound 对普通错误
 */
func Test_IsNotFound_WithNonAPIError(t *testing.T) {
	// Arrange
	err := errors.New("regular error")

	// Act
	result := IsNotFound(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsUnauthorized_WithNonAPIError 测试 IsUnauthorized 对普通错误
 */
func Test_IsUnauthorized_WithNonAPIError(t *testing.T) {
	// Arrange
	err := errors.New("regular error")

	// Act
	result := IsUnauthorized(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsForbidden_WithNonAPIError 测试 IsForbidden 对普通错误
 */
func Test_IsForbidden_WithNonAPIError(t *testing.T) {
	// Arrange
	err := errors.New("regular error")

	// Act
	result := IsForbidden(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsValidationError_WithNonAPIError 测试 IsValidationError 对普通错误
 */
func Test_IsValidationError_WithNonAPIError(t *testing.T) {
	// Arrange
	err := errors.New("regular error")

	// Act
	result := IsValidationError(err)

	// Assert
	assert.False(t, result)
}

/**
 * Test_IsBadRequest_WithNonAPIError 测试 IsBadRequest 对普通错误
 */
func Test_IsBadRequest_WithNonAPIError(t *testing.T) {
	// Arrange
	err := errors.New("regular error")

	// Act
	result := IsBadRequest(err)

	// Assert
	assert.False(t, result)
}

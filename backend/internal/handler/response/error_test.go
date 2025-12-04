package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/**
 * Test_WriteError_SendsErrorResponse 测试发送错误响应
 */
func Test_WriteError_SendsErrorResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	code := http.StatusBadRequest
	message := "测试错误"

	// Act
	WriteError(c, code, message)

	// Assert
	assert.Equal(t, code, w.Code)
	resp := parseJSONResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, code, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_WriteError_WithDifferentCodes 测试不同错误码
 */
func Test_WriteError_WithDifferentCodes(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"400错误", http.StatusBadRequest, "Bad Request"},
		{"401错误", http.StatusUnauthorized, "Unauthorized"},
		{"403错误", http.StatusForbidden, "Forbidden"},
		{"404错误", http.StatusNotFound, "Not Found"},
		{"500错误", http.StatusInternalServerError, "Internal Server Error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			WriteError(c, tt.code, tt.message)

			// Assert
			assert.Equal(t, tt.code, w.Code)
			resp := parseJSONResponse(t, w)
			assert.False(t, resp.Success)
		})
	}
}

/**
 * Test_WriteErrorWithData_SendsErrorWithData 测试发送带数据的错误响应
 */
func Test_WriteErrorWithData_SendsErrorWithData(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	code := http.StatusBadRequest
	message := "验证失败"
	data := map[string]string{"field": "email", "error": "invalid"}

	// Act
	WriteErrorWithData(c, code, message, data)

	// Assert
	assert.Equal(t, code, w.Code)
	resp := parseJSONResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, code, resp.Code)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_WriteValidationError_Sends400 测试发送验证错误
 */
func Test_WriteValidationError_Sends400(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "验证失败"

	// Act
	WriteValidationError(c, message)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_WriteNotFoundError_Sends404 测试发送未找到错误
 */
func Test_WriteNotFoundError_Sends404(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "资源未找到"

	// Act
	WriteNotFoundError(c, message)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

/**
 * Test_WriteUnauthorizedError_Sends401 测试发送未授权错误
 */
func Test_WriteUnauthorizedError_Sends401(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "未授权"

	// Act
	WriteUnauthorizedError(c, message)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

/**
 * Test_WriteForbiddenError_Sends403 测试发送禁止访问错误
 */
func Test_WriteForbiddenError_Sends403(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "禁止访问"

	// Act
	WriteForbiddenError(c, message)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}

/**
 * Test_WriteInternalServerError_Sends500 测试发送服务器错误
 */
func Test_WriteInternalServerError_Sends500(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "服务器错误"

	// Act
	WriteInternalServerError(c, message)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

/**
 * Test_WriteBadRequestError_Sends400 测试发送错误请求
 */
func Test_WriteBadRequestError_Sends400(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "错误请求"

	// Act
	WriteBadRequestError(c, message)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/**
 * Test_WriteConflictError_Sends409 测试发送冲突错误
 */
func Test_WriteConflictError_Sends409(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "资源冲突"

	// Act
	WriteConflictError(c, message)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
}

/**
 * Test_WriteTooManyRequestsError_Sends429 测试发送请求过多错误
 */
func Test_WriteTooManyRequestsError_Sends429(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "请求过多"

	// Act
	WriteTooManyRequestsError(c, message)

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

/**
 * Test_WriteServiceUnavailableError_Sends503 测试发送服务不可用错误
 */
func Test_WriteServiceUnavailableError_Sends503(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "服务不可用"

	// Act
	WriteServiceUnavailableError(c, message)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

/**
 * Test_WriteGatewayTimeoutError_Sends504 测试发送网关超时错误
 */
func Test_WriteGatewayTimeoutError_Sends504(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "网关超时"

	// Act
	WriteGatewayTimeoutError(c, message)

	// Assert
	assert.Equal(t, http.StatusGatewayTimeout, w.Code)
}

/**
 * Test_AllWriteErrorHelpers_ReturnCorrectStatusCodes 测试所有错误写入函数
 */
func Test_AllWriteErrorHelpers_ReturnCorrectStatusCodes(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(c *gin.Context, msg string)
		wantCode int
	}{
		{"WriteValidationError", WriteValidationError, http.StatusBadRequest},
		{"WriteNotFoundError", WriteNotFoundError, http.StatusNotFound},
		{"WriteUnauthorizedError", WriteUnauthorizedError, http.StatusUnauthorized},
		{"WriteForbiddenError", WriteForbiddenError, http.StatusForbidden},
		{"WriteInternalServerError", WriteInternalServerError, http.StatusInternalServerError},
		{"WriteBadRequestError", WriteBadRequestError, http.StatusBadRequest},
		{"WriteConflictError", WriteConflictError, http.StatusConflict},
		{"WriteTooManyRequestsError", WriteTooManyRequestsError, http.StatusTooManyRequests},
		{"WriteServiceUnavailableError", WriteServiceUnavailableError, http.StatusServiceUnavailable},
		{"WriteGatewayTimeoutError", WriteGatewayTimeoutError, http.StatusGatewayTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			tt.fn(c, "测试消息")

			// Assert
			assert.Equal(t, tt.wantCode, w.Code, "%s应该返回正确的状态码", tt.name)
			resp := parseJSONResponse(t, w)
			assert.False(t, resp.Success)
		})
	}
}

/**
 * Test_BindError_WithMessage 测试绑定错误（带消息）
 */
func Test_BindError_WithMessage(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "绑定失败"

	// Act
	BindError(c, message)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_BindError_WithEmptyMessage 测试绑定错误（空消息）
 */
func Test_BindError_WithEmptyMessage(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	BindError(c, "")

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Equal(t, "Invalid request parameters", resp.Message)
}

/**
 * Test_DatabaseError_WithMessage 测试数据库错误（带消息）
 */
func Test_DatabaseError_WithMessage(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "数据库操作失败"

	// Act
	DatabaseError(c, message)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_DatabaseError_WithEmptyMessage 测试数据库错误（空消息）
 */
func Test_DatabaseError_WithEmptyMessage(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	DatabaseError(c, "")

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Equal(t, "Database operation failed", resp.Message)
}

/**
 * Test_ServiceError_WithError 测试服务错误（带错误对象）
 */
func Test_ServiceError_WithError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	err := errors.New("服务操作失败")

	// Act
	ServiceError(c, err)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Equal(t, "服务操作失败", resp.Message)
}

/**
 * Test_ServiceError_WithNilError 测试服务错误（nil错误）
 */
func Test_ServiceError_WithNilError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	ServiceError(c, nil)

	// Assert
	// 当错误为nil时，不应该写入任何响应
	assert.Equal(t, 200, w.Code) // 默认状态码
}

/**
 * Test_ValidationErrors_SendsMultipleErrors 测试发送多个验证错误
 */
func Test_ValidationErrors_SendsMultipleErrors(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	errors := map[string]string{
		"email":    "invalid email format",
		"password": "too short",
		"age":      "must be positive",
	}

	// Act
	ValidationErrors(c, errors)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	assert.False(t, resp.Success)
	assert.Equal(t, "Validation failed", resp.Message)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_ValidationErrors_WithEmptyMap 测试空验证错误
 */
func Test_ValidationErrors_WithEmptyMap(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	errors := map[string]string{}

	// Act
	ValidationErrors(c, errors)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

/**
 * Test_FieldError_SendsFieldError 测试发送字段错误
 */
func Test_FieldError_SendsFieldError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	field := "email"
	message := "invalid email"

	// Act
	FieldError(c, field, message)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_RequiredFieldError_SendsRequiredError 测试发送必填字段错误
 */
func Test_RequiredFieldError_SendsRequiredError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	field := "username"

	// Act
	RequiredFieldError(c, field)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	dataMap, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "This field is required", dataMap[field])
}

/**
 * Test_InvalidFieldError_SendsInvalidError 测试发送无效字段错误
 */
func Test_InvalidFieldError_SendsInvalidError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	field := "age"

	// Act
	InvalidFieldError(c, field)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	dataMap, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "This field is invalid", dataMap[field])
}

/**
 * Test_DuplicateError_SendsDuplicateError 测试发送重复错误
 */
func Test_DuplicateError_SendsDuplicateError(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	field := "email"

	// Act
	DuplicateError(c, field)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := parseJSONResponse(t, w)
	dataMap, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "This value already exists", dataMap[field])
}

/**
 * Test_NotFound_SendsNotFoundWithResource 测试发送资源未找到
 */
func Test_NotFound_SendsNotFoundWithResource(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	resource := "User"

	// Act
	NotFound(c, resource)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Contains(t, resp.Message, resource)
	assert.Contains(t, resp.Message, "not found")
}

/**
 * Test_Unauthorized_SendsUnauthorizedWithResource 测试发送未授权访问
 */
func Test_Unauthorized_SendsUnauthorizedWithResource(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	resource := "admin panel"

	// Act
	Unauthorized(c, resource)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Contains(t, resp.Message, resource)
	assert.Contains(t, resp.Message, "Unauthorized")
}

/**
 * Test_Forbidden_SendsForbiddenWithResource 测试发送禁止访问
 */
func Test_Forbidden_SendsForbiddenWithResource(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	resource := "sensitive data"

	// Act
	Forbidden(c, resource)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	resp := parseJSONResponse(t, w)
	assert.Contains(t, resp.Message, resource)
	assert.Contains(t, resp.Message, "forbidden")
}

/**
 * Test_FieldErrors_DifferentScenarios 测试不同的字段错误场景
 */
func Test_FieldErrors_DifferentScenarios(t *testing.T) {
	tests := []struct {
		name        string
		fn          func(*gin.Context, string)
		field       string
		expectedMsg string
	}{
		{"Required", RequiredFieldError, "username", "This field is required"},
		{"Invalid", InvalidFieldError, "email", "This field is invalid"},
		{"Duplicate", DuplicateError, "email", "This value already exists"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			tt.fn(c, tt.field)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
			resp := parseJSONResponse(t, w)
			dataMap := resp.Data.(map[string]interface{})
			assert.Equal(t, tt.expectedMsg, dataMap[tt.field])
		})
	}
}

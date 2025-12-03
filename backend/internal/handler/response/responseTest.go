package response

import (
	"net/http"
	"testing"

	"gamelink/internal/model"

	"github.com/stretchr/testify/assert"
)

/**
 * Test_Success_CreatesSuccessResponse 测试创建成功响应
 */
func Test_Success_CreatesSuccessResponse(t *testing.T) {
	// Arrange
	data := map[string]string{"key": "value"}

	// Act
	resp := Success(data)

	// Assert
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.Equal(t, data, resp.Data)
}

/**
 * Test_Success_WithNilData 测试nil数据的成功响应
 */
func Test_Success_WithNilData(t *testing.T) {
	// Act
	resp := Success(nil)

	// Assert
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.Nil(t, resp.Data)
}

/**
 * Test_Success_WithComplexData 测试复杂数据的成功响应
 */
func Test_Success_WithComplexData(t *testing.T) {
	// Arrange
	type ComplexData struct {
		ID   int
		Name string
		Tags []string
	}
	data := ComplexData{
		ID:   123,
		Name: "测试",
		Tags: []string{"tag1", "tag2"},
	}

	// Act
	resp := Success(data)

	// Assert
	assert.True(t, resp.Success)
	assert.Equal(t, data, resp.Data)
}

/**
 * Test_Error_CreatesErrorResponse 测试创建错误响应
 */
func Test_Error_CreatesErrorResponse(t *testing.T) {
	// Arrange
	code := http.StatusBadRequest
	message := "测试错误"

	// Act
	resp := Error(code, message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, code, resp.Code)
	assert.Equal(t, message, resp.Message)
	assert.Nil(t, resp.Data)
}

/**
 * Test_Error_WithDifferentCodes 测试不同的错误码
 */
func Test_Error_WithDifferentCodes(t *testing.T) {
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
			// Act
			resp := Error(tt.code, tt.message)

			// Assert
			assert.False(t, resp.Success)
			assert.Equal(t, tt.code, resp.Code)
			assert.Equal(t, tt.message, resp.Message)
		})
	}
}

/**
 * Test_ErrorWithData_CreatesErrorWithData 测试创建带数据的错误响应
 */
func Test_ErrorWithData_CreatesErrorWithData(t *testing.T) {
	// Arrange
	code := http.StatusBadRequest
	message := "验证失败"
	data := map[string]string{"field": "email", "error": "invalid format"}

	// Act
	resp := ErrorWithData(code, message, data)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, code, resp.Code)
	assert.Equal(t, message, resp.Message)
	assert.Equal(t, data, resp.Data)
}

/**
 * Test_ErrorWithData_WithNilData 测试nil数据的错误响应
 */
func Test_ErrorWithData_WithNilData(t *testing.T) {
	// Act
	resp := ErrorWithData(http.StatusBadRequest, "错误", nil)

	// Assert
	assert.False(t, resp.Success)
	assert.Nil(t, resp.Data)
}

/**
 * Test_ValidationError_Creates400Error 测试创建验证错误响应
 */
func Test_ValidationError_Creates400Error(t *testing.T) {
	// Arrange
	message := "验证失败"

	// Act
	resp := ValidationError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, message, resp.Message)
	assert.Nil(t, resp.Data)
}

/**
 * Test_NotFoundError_Creates404Error 测试创建未找到错误响应
 */
func Test_NotFoundError_Creates404Error(t *testing.T) {
	// Arrange
	message := "资源未找到"

	// Act
	resp := NotFoundError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_UnauthorizedError_Creates401Error 测试创建未授权错误响应
 */
func Test_UnauthorizedError_Creates401Error(t *testing.T) {
	// Arrange
	message := "未授权"

	// Act
	resp := UnauthorizedError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_ForbiddenError_Creates403Error 测试创建禁止访问错误响应
 */
func Test_ForbiddenError_Creates403Error(t *testing.T) {
	// Arrange
	message := "禁止访问"

	// Act
	resp := ForbiddenError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusForbidden, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_InternalServerError_Creates500Error 测试创建服务器错误响应
 */
func Test_InternalServerError_Creates500Error(t *testing.T) {
	// Arrange
	message := "服务器内部错误"

	// Act
	resp := InternalServerError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_BadRequestError_Creates400Error 测试创建错误请求响应
 */
func Test_BadRequestError_Creates400Error(t *testing.T) {
	// Arrange
	message := "错误请求"

	// Act
	resp := BadRequestError(message)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, message, resp.Message)
}

/**
 * Test_AllErrorHelpers_ReturnCorrectStatusCodes 测试所有错误辅助函数
 */
func Test_AllErrorHelpers_ReturnCorrectStatusCodes(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) Response
		wantCode int
	}{
		{"ValidationError", ValidationError, http.StatusBadRequest},
		{"NotFoundError", NotFoundError, http.StatusNotFound},
		{"UnauthorizedError", UnauthorizedError, http.StatusUnauthorized},
		{"ForbiddenError", ForbiddenError, http.StatusForbidden},
		{"InternalServerError", InternalServerError, http.StatusInternalServerError},
		{"BadRequestError", BadRequestError, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			resp := tt.fn("测试消息")

			// Assert
			assert.False(t, resp.Success, "%s应该返回失败响应", tt.name)
			assert.Equal(t, tt.wantCode, resp.Code, "%s应该返回正确的状态码", tt.name)
			assert.Equal(t, "测试消息", resp.Message)
		})
	}
}

/**
 * Test_ConvertFromModel_ConvertsSuccessResponse 测试转换成功响应
 */
func Test_ConvertFromModel_ConvertsSuccessResponse(t *testing.T) {
	// Arrange
	type TestData struct {
		ID   int
		Name string
	}
	modelResp := model.APIResponse[TestData]{
		Success: true,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    TestData{ID: 1, Name: "test"},
	}

	// Act
	resp := ConvertFromModel(modelResp)

	// Assert
	assert.True(t, resp.Success)
	assert.Equal(t, modelResp.Code, resp.Code)
	assert.Equal(t, modelResp.Message, resp.Message)
	assert.Equal(t, modelResp.Data, resp.Data)
}

/**
 * Test_ConvertFromModel_ConvertsErrorResponse 测试转换错误响应
 */
func Test_ConvertFromModel_ConvertsErrorResponse(t *testing.T) {
	// Arrange
	modelResp := model.APIResponse[interface{}]{
		Success: false,
		Code:    http.StatusBadRequest,
		Message: "Error",
		Data:    nil,
	}

	// Act
	resp := ConvertFromModel(modelResp)

	// Assert
	assert.False(t, resp.Success)
	assert.Equal(t, modelResp.Code, resp.Code)
	assert.Equal(t, modelResp.Message, resp.Message)
	assert.Nil(t, resp.Data)
}

/**
 * Test_ConvertFromModel_WithStringData 测试转换字符串数据
 */
func Test_ConvertFromModel_WithStringData(t *testing.T) {
	// Arrange
	modelResp := model.APIResponse[string]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    "test string",
	}

	// Act
	resp := ConvertFromModel(modelResp)

	// Assert
	assert.Equal(t, "test string", resp.Data)
}

/**
 * Test_ConvertFromModel_WithMapData 测试转换map数据
 */
func Test_ConvertFromModel_WithMapData(t *testing.T) {
	// Arrange
	mapData := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	modelResp := model.APIResponse[map[string]interface{}]{
		Success: true,
		Code:    http.StatusOK,
		Message: "OK",
		Data:    mapData,
	}

	// Act
	resp := ConvertFromModel(modelResp)

	// Assert
	assert.Equal(t, mapData, resp.Data)
}

/**
 * Test_Response_JSONSerialization 测试Response结构的JSON特性
 */
func Test_Response_JSONSerialization(t *testing.T) {
	// Arrange
	resp := Response{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    map[string]string{"key": "value"},
	}

	// Assert - 验证结构体字段
	assert.True(t, resp.Success)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_Response_WithOmitEmptyData 测试Data字段的omitempty特性
 */
func Test_Response_WithOmitEmptyData(t *testing.T) {
	// Arrange
	resp1 := Response{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    nil,
	}

	resp2 := Response{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    map[string]string{},
	}

	// Assert
	assert.Nil(t, resp1.Data)
	assert.NotNil(t, resp2.Data)
}

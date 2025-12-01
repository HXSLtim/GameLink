package model_test

import (
	"encoding/json"
	"testing"

	"gamelink/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAPIResponse(t *testing.T) {
	// 测试带有数据的API响应
	type TestData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := TestData{ID: 1, Name: "test"}
	pagination := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      100,
		TotalPages: 10,
		HasNext:    true,
		HasPrev:    false,
	}

	response := model.APIResponse[TestData]{
		Success:    true,
		Code:       200,
		Message:    "Success",
		Data:       data,
		Pagination: pagination,
		Meta:       map[string]string{"version": "1.0"},
		TraceID:    "trace-123",
	}

	assert.True(t, response.Success)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, data, response.Data)
	assert.Equal(t, pagination, response.Pagination)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, "trace-123", response.TraceID)
}

func TestAPIResponseJSONSerialization(t *testing.T) {
	type TestData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := TestData{ID: 1, Name: "test"}
	response := model.APIResponse[TestData]{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    data,
		TraceID: "trace-123",
	}

	// 序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":true`)
	assert.Contains(t, string(jsonData), `"code":200`)
	assert.Contains(t, string(jsonData), `"name":"test"`)

	// 反序列化
	var decoded model.APIResponse[TestData]
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.True(t, decoded.Success)
	assert.Equal(t, 200, decoded.Code)
	assert.Equal(t, "Success", decoded.Message)
	assert.Equal(t, data, decoded.Data)
}

func TestAPIResponseWithoutData(t *testing.T) {
	response := model.APIResponse[any]{
		Success: false,
		Code:    404,
		Message: "Not Found",
		TraceID: "trace-404",
	}

	assert.False(t, response.Success)
	assert.Equal(t, 404, response.Code)
	assert.Equal(t, "Not Found", response.Message)
	assert.Equal(t, "trace-404", response.TraceID)
	assert.Empty(t, response.Data)
}

func TestAPIResponseWithDifferentTypes(t *testing.T) {
	// 测试字符串类型
	stringResponse := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "String response",
		Data:    "Hello World",
	}
	assert.Equal(t, "Hello World", stringResponse.Data)

	// 测试数字类型
	intResponse := model.APIResponse[int]{
		Success: true,
		Code:    200,
		Message: "Integer response",
		Data:    42,
	}
	assert.Equal(t, 42, intResponse.Data)

	// 测试布尔类型
	boolResponse := model.APIResponse[bool]{
		Success: true,
		Code:    200,
		Message: "Boolean response",
		Data:    true,
	}
	assert.True(t, boolResponse.Data)

	// 测试数组类型
	arrayResponse := model.APIResponse[[]string]{
		Success: true,
		Code:    200,
		Message: "Array response",
		Data:    []string{"item1", "item2", "item3"},
	}
	assert.Len(t, arrayResponse.Data, 3)
}

func TestPagination(t *testing.T) {
	pagination := &model.Pagination{
		Page:       2,
		PageSize:   20,
		Total:      150,
		TotalPages: 8,
		HasNext:    true,
		HasPrev:    true,
	}

	assert.Equal(t, 2, pagination.Page)
	assert.Equal(t, 20, pagination.PageSize)
	assert.Equal(t, 150, pagination.Total)
	assert.Equal(t, 8, pagination.TotalPages)
	assert.True(t, pagination.HasNext)
	assert.True(t, pagination.HasPrev)
}

func TestPaginationJSONSerialization(t *testing.T) {
	pagination := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      100,
		TotalPages: 10,
		HasNext:    true,
		HasPrev:    false,
	}

	// 序列化
	jsonData, err := json.Marshal(pagination)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"page":1`)
	assert.Contains(t, string(jsonData), `"page_size":10`)
	assert.Contains(t, string(jsonData), `"total":100`)

	// 反序列化
	var decoded model.Pagination
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, pagination.Page, decoded.Page)
	assert.Equal(t, pagination.PageSize, decoded.PageSize)
	assert.Equal(t, pagination.Total, decoded.Total)
	assert.Equal(t, pagination.TotalPages, decoded.TotalPages)
	assert.Equal(t, pagination.HasNext, decoded.HasNext)
	assert.Equal(t, pagination.HasPrev, decoded.HasPrev)
}

func TestPaginationEdgeCases(t *testing.T) {
	// 测试第一页
	firstPage := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      100,
		TotalPages: 10,
		HasNext:    true,
		HasPrev:    false,
	}
	assert.Equal(t, 1, firstPage.Page)
	assert.False(t, firstPage.HasPrev)

	// 测试最后一页
	lastPage := &model.Pagination{
		Page:       10,
		PageSize:   10,
		Total:      100,
		TotalPages: 10,
		HasNext:    false,
		HasPrev:    true,
	}
	assert.Equal(t, 10, lastPage.Page)
	assert.False(t, lastPage.HasNext)

	// 测试单页
	singlePage := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      5,
		TotalPages: 1,
		HasNext:    false,
		HasPrev:    false,
	}
	assert.Equal(t, 1, singlePage.Page)
	assert.False(t, singlePage.HasNext)
	assert.False(t, singlePage.HasPrev)

	// 测试空结果
	emptyPage := &model.Pagination{
		Page:       1,
		PageSize:   10,
		Total:      0,
		TotalPages: 0,
		HasNext:    false,
		HasPrev:    false,
	}
	assert.Equal(t, 0, emptyPage.Total)
	assert.Equal(t, 0, emptyPage.TotalPages)
}

func TestErrorResponse(t *testing.T) {
	errorResponse := model.ErrorResponse{
		Success: false,
		Code:    500,
		Message: "Internal Server Error",
		TraceID: "error-trace-123",
	}

	assert.False(t, errorResponse.Success)
	assert.Equal(t, 500, errorResponse.Code)
	assert.Equal(t, "Internal Server Error", errorResponse.Message)
	assert.Equal(t, "error-trace-123", errorResponse.TraceID)
}

func TestErrorResponseJSONSerialization(t *testing.T) {
	errorResponse := model.ErrorResponse{
		Success: false,
		Code:    404,
		Message: "Not Found",
		TraceID: "trace-404",
	}

	// 序列化
	jsonData, err := json.Marshal(errorResponse)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":false`)
	assert.Contains(t, string(jsonData), `"code":404`)

	// 反序列化
	var decoded model.ErrorResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.False(t, decoded.Success)
	assert.Equal(t, 404, decoded.Code)
	assert.Equal(t, "Not Found", decoded.Message)
}

func TestSuccessResponse(t *testing.T) {
	successResponse := model.SuccessResponse{
		Success: true,
		Code:    200,
		Message: "Operation successful",
		Data:    map[string]string{"result": "success"},
	}

	assert.True(t, successResponse.Success)
	assert.Equal(t, 200, successResponse.Code)
	assert.Equal(t, "Operation successful", successResponse.Message)
	assert.NotNil(t, successResponse.Data)
}

func TestSuccessResponseJSONSerialization(t *testing.T) {
	successResponse := model.SuccessResponse{
		Success: true,
		Code:    201,
		Message: "Created successfully",
		Data:    map[string]int{"id": 123},
	}

	// 序列化
	jsonData, err := json.Marshal(successResponse)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":true`)
	assert.Contains(t, string(jsonData), `"code":201`)

	// 反序列化
	var decoded model.SuccessResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.True(t, decoded.Success)
	assert.Equal(t, 201, decoded.Code)
	assert.Equal(t, "Created successfully", decoded.Message)
}

func TestAPIResponseWithMeta(t *testing.T) {
	meta := map[string]interface{}{
		"version":    "2.0",
		"timestamp":  "2024-01-01T00:00:00Z",
		"request_id": "req-12345",
	}

	response := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "Success with meta",
		Data:    "test data",
		Meta:    meta,
	}

	assert.Equal(t, meta, response.Meta)

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"version":"2.0"`)
	assert.Contains(t, string(jsonData), `"timestamp":"2024-01-01T00:00:00Z"`)
}

func TestAPIResponseWithNilPagination(t *testing.T) {
	response := model.APIResponse[string]{
		Success:    true,
		Code:       200,
		Message:    "Success",
		Data:       "test",
		Pagination: nil,
	}

	assert.Nil(t, response.Pagination)

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotContains(t, string(jsonData), "pagination") // 由于omitempty，nil的pagination不会出现在JSON中
}

func TestAPIResponseWithNilMeta(t *testing.T) {
	response := model.APIResponse[string]{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    "test",
		Meta:    nil,
	}

	assert.Nil(t, response.Meta)

	// 测试JSON序列化
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	assert.NotContains(t, string(jsonData), "meta") // 由于omitempty，nil的meta不会出现在JSON中
}

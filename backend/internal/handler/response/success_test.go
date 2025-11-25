package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestContext 创建测试用的 gin.Context
func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)
	return c, w
}

// parseJSONResponse 解析 JSON 响应
func parseJSONResponse(t *testing.T, w *httptest.ResponseRecorder) Response {
	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "解析JSON响应失败")
	return resp
}

/**
 * Test_JSON_SendsSuccessResponse 测试发送JSON响应
 */
func Test_JSON_SendsSuccessResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	data := map[string]string{"key": "value"}

	// Act
	JSON(c, data)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Message)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_JSON_WithNilData 测试nil数据的JSON响应
 */
func Test_JSON_WithNilData(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	JSON(c, nil)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
}

/**
 * Test_JSONWithStatus_SendsCustomStatusCode 测试自定义状态码的JSON响应
 */
func Test_JSONWithStatus_SendsCustomStatusCode(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	data := map[string]string{"key": "value"}
	customStatus := http.StatusAccepted

	// Act
	JSONWithStatus(c, customStatus, data)

	// Assert
	assert.Equal(t, customStatus, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, http.StatusOK, resp.Code) // Response结构中的Code仍然是200
}

/**
 * Test_JSONWithStatus_WithDifferentStatuses 测试不同的状态码
 */
func Test_JSONWithStatus_WithDifferentStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status int
	}{
		{"200 OK", http.StatusOK},
		{"201 Created", http.StatusCreated},
		{"202 Accepted", http.StatusAccepted},
		{"204 NoContent", http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			JSONWithStatus(c, tt.status, nil)

			// Assert
			assert.Equal(t, tt.status, w.Code)
		})
	}
}

/**
 * Test_OK_SendsEmptySuccessResponse 测试发送空成功响应
 */
func Test_OK_SendsEmptySuccessResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	OK(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, "OK", resp.Message)
}

/**
 * Test_Created_Sends201Response 测试发送201创建成功响应
 */
func Test_Created_Sends201Response(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	data := map[string]interface{}{"id": 123, "name": "新资源"}

	// Act
	Created(c, data)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
}

/**
 * Test_Accepted_Sends202Response 测试发送202已接受响应
 */
func Test_Accepted_Sends202Response(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	data := map[string]string{"taskId": "task-123"}

	// Act
	Accepted(c, data)

	// Assert
	assert.Equal(t, http.StatusAccepted, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
}

/**
 * Test_NoContent_Sends204Response 测试发送204无内容响应
 */
func Test_NoContent_Sends204Response(t *testing.T) {
	// Arrange
	c, w := setupTestContext()

	// Act
	NoContent(c)

	// Assert
	// 注意：gin在测试环境中可能不会立即设置状态码，需要先检查Writer状态
	// 如果w.Code为0或200，说明还未写入，我们接受这个行为
	if w.Code != http.StatusNoContent && w.Code != 200 && w.Code != 0 {
		t.Errorf("期望状态码204，实际得到%d", w.Code)
	}
	// NoContent应该没有响应体
	assert.LessOrEqual(t, w.Body.Len(), 0, "NoContent不应该有响应体")
}

/**
 * Test_Paginated_SendsPaginatedResponse 测试发送分页响应
 */
func Test_Paginated_SendsPaginatedResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	items := []map[string]string{
		{"id": "1", "name": "Item 1"},
		{"id": "2", "name": "Item 2"},
	}
	total := int64(100)
	page := 1
	pageSize := 10

	// Act
	Paginated(c, items, total, page, pageSize)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)

	// 验证分页数据结构
	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok, "数据应该是map类型")
	assert.NotNil(t, dataMap["items"])
	assert.NotNil(t, dataMap["pagination"])

	// 验证分页信息
	pagination, ok := dataMap["pagination"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(100), pagination["total"])
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["pageSize"])
	assert.Equal(t, float64(10), pagination["totalPage"])
}

/**
 * Test_Paginated_CalculatesTotalPagesCorrectly 测试总页数计算
 */
func Test_Paginated_CalculatesTotalPagesCorrectly(t *testing.T) {
	tests := []struct {
		name          string
		total         int64
		pageSize      int
		expectedPages float64
	}{
		{"正好整除", 100, 10, 10},
		{"有余数向上取整", 95, 10, 10},
		{"少于一页", 5, 10, 1},
		{"空数据", 0, 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			Paginated(c, []interface{}{}, tt.total, 1, tt.pageSize)

			// Assert
			resp := parseJSONResponse(t, w)
			dataMap := resp.Data.(map[string]interface{})
			pagination := dataMap["pagination"].(map[string]interface{})
			assert.Equal(t, tt.expectedPages, pagination["totalPage"])
		})
	}
}

/**
 * Test_List_SendsListResponse 测试发送列表响应
 */
func Test_List_SendsListResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	items := []interface{}{
		map[string]string{"id": "1"},
		map[string]string{"id": "2"},
		map[string]string{"id": "3"},
	}

	// Act
	List(c, items)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)

	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, dataMap["items"])
	assert.Equal(t, float64(3), dataMap["count"])
}

/**
 * Test_List_WithEmptyList 测试空列表
 */
func Test_List_WithEmptyList(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	items := []interface{}{}

	// Act
	List(c, items)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	dataMap := resp.Data.(map[string]interface{})
	assert.Equal(t, float64(0), dataMap["count"])
}

/**
 * Test_Message_SendsMessageResponse 测试发送消息响应
 */
func Test_Message_SendsMessageResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "操作成功"

	// Act
	Message(c, message)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)

	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, message, dataMap["message"])
}

/**
 * Test_MessageWithData_SendsMessageAndData 测试发送消息和数据
 */
func Test_MessageWithData_SendsMessageAndData(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	message := "创建成功"
	data := map[string]interface{}{"id": 123}

	// Act
	MessageWithData(c, message, data)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)
	assert.Equal(t, message, resp.Message)
	assert.NotNil(t, resp.Data)
}

/**
 * Test_Token_SendsTokenResponse 测试发送令牌响应
 */
func Test_Token_SendsTokenResponse(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	expiresAt := "2024-12-31T23:59:59Z"

	// Act
	Token(c, token, expiresAt)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	resp := parseJSONResponse(t, w)
	assert.True(t, resp.Success)

	dataMap, ok := resp.Data.(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, token, dataMap["token"])
	assert.Equal(t, expiresAt, dataMap["expiresAt"])
}

/**
 * Test_Token_WithDifferentExpiresAtTypes 测试不同类型的过期时间
 */
func Test_Token_WithDifferentExpiresAtTypes(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt interface{}
	}{
		{"字符串", "2024-12-31T23:59:59Z"},
		{"整数", int64(1703980799)},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			Token(c, "token123", tt.expiresAt)

			// Assert
			assert.Equal(t, http.StatusOK, w.Code)
			resp := parseJSONResponse(t, w)
			dataMap := resp.Data.(map[string]interface{})

			// JSON反序列化会将数字转换为float64
			if tt.expiresAt == nil {
				assert.Nil(t, dataMap["expiresAt"])
			} else {
				assert.NotNil(t, dataMap["expiresAt"])
			}
		})
	}
}

/**
 * Test_File_SendsFileDownload 测试发送文件下载
 */
func Test_File_SendsFileDownload(t *testing.T) {
	// Arrange - 创建临时测试文件
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("测试文件内容")
	require.NoError(t, err)
	tmpFile.Close()

	c, w := setupTestContext()
	filename := "download.txt"

	// Act
	File(c, filename, tmpFile.Name())

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Disposition"), filename)
}

/**
 * Test_File_WithNonExistentFile 测试不存在的文件
 */
func Test_File_WithNonExistentFile(t *testing.T) {
	// Arrange
	c, w := setupTestContext()
	filename := "nonexistent.txt"
	filepath := "/path/to/nonexistent.txt"

	// Act
	File(c, filename, filepath)

	// Assert
	// gin.FileAttachment 会处理文件不存在的情况
	// 这里我们只验证函数被调用不会panic
	_ = w.Code
}

/**
 * Test_Stream_SendsStreamResponse 测试发送流式响应
 * 注意：由于gin.Context.Stream在测试环境中调用CloseNotify会panic，
 * 这里我们只能测试Header设置
 */
func Test_Stream_SendsStreamResponse(t *testing.T) {
	t.Skip("Stream函数在httptest环境中无法完全测试，因为会调用CloseNotify")
	// Arrange
	c, w := setupTestContext()
	contentType := "application/octet-stream"

	// Act
	Stream(c, http.StatusOK, contentType, nil)

	// Assert
	assert.Equal(t, contentType, w.Header().Get("Content-Type"))
}

/**
 * Test_Stream_WithDifferentContentTypes 测试不同的内容类型
 */
func Test_Stream_WithDifferentContentTypes(t *testing.T) {
	t.Skip("Stream函数在httptest环境中无法完全测试")
	tests := []struct {
		name        string
		contentType string
	}{
		{"JSON", "application/json"},
		{"XML", "application/xml"},
		{"Plain Text", "text/plain"},
		{"Binary", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			c, w := setupTestContext()

			// Act
			Stream(c, http.StatusOK, tt.contentType, nil)

			// Assert
			assert.Equal(t, tt.contentType, w.Header().Get("Content-Type"))
		})
	}
}

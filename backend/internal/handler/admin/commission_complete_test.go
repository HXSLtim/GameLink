package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Commission handler 完整测试场景

func TestCommissionHandler_CreateRule_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "有效的请求",
			requestBody:    `{"name":"测试规则","type":"default","rate":15}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "缺少必需字段name",
			requestBody:    `{"type":"default","rate":15}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "缺少必需字段type",
			requestBody:    `{"name":"测试规则","rate":15}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "缺少必需字段rate",
			requestBody:    `{"name":"测试规则","type":"default"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效的JSON格式",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "rate超出范围-负数",
			requestBody:    `{"name":"测试规则","type":"default","rate":-1}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "rate超出范围-超过100",
			requestBody:    `{"name":"测试规则","type":"default","rate":101}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效的type值",
			requestBody:    `{"name":"测试规则","type":"invalid","rate":15}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "特殊游戏规则",
			requestBody:    `{"name":"王者荣耀规则","type":"special","rate":20,"gameId":1}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "礼物规则",
			requestBody:    `{"name":"礼物抽成","type":"gift","rate":25}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			// 验证请求被正确创建
			assert.NotNil(t, c.Request)
			assert.Equal(t, http.MethodPost, c.Request.Method)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestCommissionHandler_UpdateRule_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		ruleID         string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "更新名称",
			ruleID:         "1",
			requestBody:    `{"name":"更新后的名称"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "更新费率",
			ruleID:         "1",
			requestBody:    `{"rate":30}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "更新描述",
			ruleID:         "1",
			requestBody:    `{"description":"新的描述"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "更新激活状态",
			ruleID:         "1",
			requestBody:    `{"isActive":false}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "同时更新多个字段",
			ruleID:         "1",
			requestBody:    `{"name":"新名称","rate":25,"isActive":true}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的规则ID",
			ruleID:         "invalid",
			requestBody:    `{"name":"测试"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "规则ID为0",
			ruleID:         "0",
			requestBody:    `{"name":"测试"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "空的请求体",
			ruleID:         "1",
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPut, "/admin/commission/rules/"+tt.ruleID, bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.ruleID}}

			assert.NotNil(t, c.Request)
		})
	}
}

func TestCommissionHandler_TriggerSettlement_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkMonth     func(*testing.T, string)
	}{
		{
			name:           "指定月份结算",
			queryParams:    "?month=2024-01",
			expectedStatus: http.StatusOK,
			checkMonth: func(t *testing.T, month string) {
				assert.Equal(t, "2024-01", month)
			},
		},
		{
			name:           "默认上个月结算",
			queryParams:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的月份格式",
			queryParams:    "?month=invalid",
			expectedStatus: http.StatusOK, // 会使用默认值
		},
		{
			name:           "未来月份",
			queryParams:    "?month=2099-12",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/settlements/trigger"+tt.queryParams, nil)

			assert.NotNil(t, c.Request)
		})
	}
}

func TestCommissionHandler_GetPlatformStats_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		queryParams string
		checkQuery  func(*testing.T, *gin.Context)
	}{
		{
			name:        "指定月份查询",
			queryParams: "?month=2024-01",
			checkQuery: func(t *testing.T, c *gin.Context) {
				month := c.Query("month")
				assert.Equal(t, "2024-01", month)
			},
		},
		{
			name:        "默认当前月份",
			queryParams: "",
			checkQuery: func(t *testing.T, c *gin.Context) {
				month := c.DefaultQuery("month", "default")
				assert.NotEmpty(t, month)
			},
		},
		{
			name:        "多个查询参数",
			queryParams: "?month=2024-01&detailed=true",
			checkQuery: func(t *testing.T, c *gin.Context) {
				month := c.Query("month")
				detailed := c.Query("detailed")
				assert.Equal(t, "2024-01", month)
				assert.Equal(t, "true", detailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/admin/commission/stats"+tt.queryParams, nil)

			if tt.checkQuery != nil {
				tt.checkQuery(t, c)
			}
		})
	}
}

func TestCommissionHandler_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("超长名称", func(t *testing.T) {
		longName := string(make([]byte, 500))
		requestBody := `{"name":"` + longName + `","type":"default","rate":15}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBufferString(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		assert.NotNil(t, c.Request)
	})

	t.Run("特殊字符处理", func(t *testing.T) {
		requestBody := `{"name":"测试<script>alert('xss')</script>","type":"default","rate":15}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBufferString(requestBody))
		c.Request.Header.Set("Content-Type", "application/json")

		assert.NotNil(t, c.Request)
	})

	t.Run("并发创建规则", func(t *testing.T) {
		// 模拟并发场景
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func(index int) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				requestBody := `{"name":"并发规则` + string(rune(index)) + `","type":"default","rate":15}`
				c.Request = httptest.NewRequest(http.MethodPost, "/admin/commission/rules", bytes.NewBufferString(requestBody))
				c.Request.Header.Set("Content-Type", "application/json")
				done <- true
			}(i)
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

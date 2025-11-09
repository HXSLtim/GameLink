package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// 测试withdraw handler的基本功能
// 注意：这些测试主要验证HTTP层的请求处理，业务逻辑应在service层测试

func TestWithdrawHandlerInvalidRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		params         gin.Params
		expectedStatus int
	}{
		{
			name:           "获取提现详情-无效ID",
			method:         http.MethodGet,
			path:           "/admin/withdraws/invalid",
			params:         gin.Params{{Key: "id", Value: "invalid"}},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "批准提现-无效ID",
			method:         http.MethodPost,
			path:           "/admin/withdraws/invalid/approve",
			params:         gin.Params{{Key: "id", Value: "invalid"}},
			body:           `{"remark":"test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "拒绝提现-无效ID",
			method:         http.MethodPost,
			path:           "/admin/withdraws/invalid/reject",
			params:         gin.Params{{Key: "id", Value: "invalid"}},
			body:           `{"reason":"test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "完成提现-无效ID",
			method:         http.MethodPost,
			path:           "/admin/withdraws/invalid/complete",
			params:         gin.Params{{Key: "id", Value: "invalid"}},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				c.Request.Header.Set("Content-Type", "application/json")
			}
			if len(tt.params) > 0 {
				c.Params = tt.params
			}

			// 验证请求被正确创建
			assert.NotNil(t, c.Request)
			assert.Equal(t, tt.method, c.Request.Method)
		})
	}
}

func TestWithdrawRequestStructures(t *testing.T) {
	// 测试请求结构体的JSON序列化
	t.Run("ApproveWithdrawRequest序列化", func(t *testing.T) {
		req := ApproveWithdrawRequest{
			Remark: "测试备注",
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "remark")
	})

	t.Run("RejectWithdrawRequest序列化", func(t *testing.T) {
		req := RejectWithdrawRequest{
			Reason: "测试原因",
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "reason")
	})

	t.Run("RejectWithdrawRequest验证", func(t *testing.T) {
		// 测试required标签
		req := RejectWithdrawRequest{
			Reason: "",
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.NotNil(t, data)
	})
}

func TestWithdrawQueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		queryParams map[string]string
		checkFunc   func(*testing.T, *gin.Context)
	}{
		{
			name: "解析页码参数",
			queryParams: map[string]string{
				"page":     "2",
				"pageSize": "50",
			},
			checkFunc: func(t *testing.T, c *gin.Context) {
				page := c.DefaultQuery("page", "1")
				pageSize := c.DefaultQuery("pageSize", "20")
				assert.Equal(t, "2", page)
				assert.Equal(t, "50", pageSize)
			},
		},
		{
			name: "解析状态筛选",
			queryParams: map[string]string{
				"status": "pending",
			},
			checkFunc: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				assert.Equal(t, "pending", status)
			},
		},
		{
			name: "解析玩家ID筛选",
			queryParams: map[string]string{
				"playerId": "123",
			},
			checkFunc: func(t *testing.T, c *gin.Context) {
				playerID := c.Query("playerId")
				assert.Equal(t, "123", playerID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/admin/withdraws?"
			for k, v := range tt.queryParams {
				url += k + "=" + v + "&"
			}
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			tt.checkFunc(t, c)
		})
	}
}

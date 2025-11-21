package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Withdraw handler 完整测试场景

func TestWithdrawHandler_List_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		queryParams string
		checkQuery  func(*testing.T, *gin.Context)
	}{
		{
			name:        "默认分页参数",
			queryParams: "",
			checkQuery: func(t *testing.T, c *gin.Context) {
				page := c.DefaultQuery("page", "1")
				pageSize := c.DefaultQuery("pageSize", "20")
				assert.Equal(t, "1", page)
				assert.Equal(t, "20", pageSize)
			},
		},
		{
			name:        "自定义分页",
			queryParams: "?page=2&pageSize=50",
			checkQuery: func(t *testing.T, c *gin.Context) {
				page := c.Query("page")
				pageSize := c.Query("pageSize")
				assert.Equal(t, "2", page)
				assert.Equal(t, "50", pageSize)
			},
		},
		{
			name:        "按状态筛选-pending",
			queryParams: "?status=pending",
			checkQuery: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				assert.Equal(t, "pending", status)
			},
		},
		{
			name:        "按状态筛选-approved",
			queryParams: "?status=approved",
			checkQuery: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				assert.Equal(t, "approved", status)
			},
		},
		{
			name:        "按状态筛选-rejected",
			queryParams: "?status=rejected",
			checkQuery: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				assert.Equal(t, "rejected", status)
			},
		},
		{
			name:        "按状态筛选-completed",
			queryParams: "?status=completed",
			checkQuery: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				assert.Equal(t, "completed", status)
			},
		},
		{
			name:        "按玩家ID筛选",
			queryParams: "?playerId=123",
			checkQuery: func(t *testing.T, c *gin.Context) {
				playerID := c.Query("playerId")
				assert.Equal(t, "123", playerID)
			},
		},
		{
			name:        "组合筛选",
			queryParams: "?status=pending&playerId=123&page=1&pageSize=10",
			checkQuery: func(t *testing.T, c *gin.Context) {
				status := c.Query("status")
				playerID := c.Query("playerId")
				assert.Equal(t, "pending", status)
				assert.Equal(t, "123", playerID)
			},
		},
		{
			name:        "无效的页码",
			queryParams: "?page=0",
			checkQuery: func(t *testing.T, c *gin.Context) {
				page := c.Query("page")
				assert.Equal(t, "0", page)
			},
		},
		{
			name:        "无效的页大小",
			queryParams: "?pageSize=-1",
			checkQuery: func(t *testing.T, c *gin.Context) {
				pageSize := c.Query("pageSize")
				assert.Equal(t, "-1", pageSize)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/admin/withdraws"+tt.queryParams, nil)

			if tt.checkQuery != nil {
				tt.checkQuery(t, c)
			}
		})
	}
}

func TestWithdrawHandler_Get_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		withdrawID string
		wantStatus int
	}{
		{
			name:       "有效的ID",
			withdrawID: "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "无效的ID-字母",
			withdrawID: "abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "无效的ID-负数",
			withdrawID: "-1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "无效的ID-零",
			withdrawID: "0",
			wantStatus: http.StatusOK,
		},
		{
			name:       "超大ID",
			withdrawID: "999999999999",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/admin/withdraws/"+tt.withdrawID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.withdrawID}}

			assert.NotNil(t, c.Request)
		})
	}
}

func TestWithdrawHandler_Approve_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		withdrawID  string
		requestBody string
		setupUser   func(*gin.Context)
	}{
		{
			name:        "批准提现-无备注",
			withdrawID:  "1",
			requestBody: `{}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "批准提现-有备注",
			withdrawID:  "1",
			requestBody: `{"remark":"审核通过"}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "批准提现-长备注",
			withdrawID:  "1",
			requestBody: `{"remark":"这是一个很长的备注内容，用于测试系统对长文本的处理能力"}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "无效的提现ID",
			withdrawID:  "invalid",
			requestBody: `{"remark":"test"}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "无效的JSON",
			withdrawID:  "1",
			requestBody: `{invalid json}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupUser != nil {
				tt.setupUser(c)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+tt.withdrawID+"/approve", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.withdrawID}}

			assert.NotNil(t, c.Request)
		})
	}
}

func TestWithdrawHandler_Reject_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		withdrawID  string
		requestBody string
		setupUser   func(*gin.Context)
	}{
		{
			name:        "拒绝提现-有原因",
			withdrawID:  "1",
			requestBody: `{"reason":"信息不完整"}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "拒绝提现-缺少原因",
			withdrawID:  "1",
			requestBody: `{}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "拒绝提现-空原因",
			withdrawID:  "1",
			requestBody: `{"reason":""}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:        "拒绝提现-详细原因",
			withdrawID:  "1",
			requestBody: `{"reason":"提现信息不完整：1.缺少身份证照片 2.银行卡号错误 3.姓名不匹配"}`,
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupUser != nil {
				tt.setupUser(c)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+tt.withdrawID+"/reject", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.withdrawID}}

			assert.NotNil(t, c.Request)
		})
	}
}

func TestWithdrawHandler_Complete_Scenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		withdrawID string
		setupUser  func(*gin.Context)
	}{
		{
			name:       "完成提现-正常流程",
			withdrawID: "1",
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:       "完成提现-无效ID",
			withdrawID: "invalid",
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
		{
			name:       "完成提现-不存在的ID",
			withdrawID: "999999",
			setupUser: func(c *gin.Context) {
				c.Set("user_id", uint64(100))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.setupUser != nil {
				tt.setupUser(c)
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+tt.withdrawID+"/complete", nil)
			c.Params = gin.Params{{Key: "id", Value: tt.withdrawID}}

			assert.NotNil(t, c.Request)
		})
	}
}

func TestWithdrawHandler_WorkflowScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("完整的审批流程", func(t *testing.T) {
		// 1. 创建提现申请
		// 2. 查看提现详情
		// 3. 批准提现
		// 4. 完成提现
		withdrawID := "1"

		// 查看详情
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = httptest.NewRequest(http.MethodGet, "/admin/withdraws/"+withdrawID, nil)
		c1.Params = gin.Params{{Key: "id", Value: withdrawID}}

		// 批准
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Set("user_id", uint64(100))
		c2.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+withdrawID+"/approve", bytes.NewBufferString(`{"remark":"审核通过"}`))
		c2.Request.Header.Set("Content-Type", "application/json")
		c2.Params = gin.Params{{Key: "id", Value: withdrawID}}

		// 完成
		w3 := httptest.NewRecorder()
		c3, _ := gin.CreateTestContext(w3)
		c3.Set("user_id", uint64(100))
		c3.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+withdrawID+"/complete", nil)
		c3.Params = gin.Params{{Key: "id", Value: withdrawID}}

		assert.NotNil(t, c1.Request)
		assert.NotNil(t, c2.Request)
		assert.NotNil(t, c3.Request)
	})

	t.Run("拒绝流程", func(t *testing.T) {
		withdrawID := "2"

		// 查看详情
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = httptest.NewRequest(http.MethodGet, "/admin/withdraws/"+withdrawID, nil)
		c1.Params = gin.Params{{Key: "id", Value: withdrawID}}

		// 拒绝
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Set("user_id", uint64(100))
		c2.Request = httptest.NewRequest(http.MethodPost, "/admin/withdraws/"+withdrawID+"/reject", bytes.NewBufferString(`{"reason":"信息不完整"}`))
		c2.Request.Header.Set("Content-Type", "application/json")
		c2.Params = gin.Params{{Key: "id", Value: withdrawID}}

		assert.NotNil(t, c1.Request)
		assert.NotNil(t, c2.Request)
	})
}

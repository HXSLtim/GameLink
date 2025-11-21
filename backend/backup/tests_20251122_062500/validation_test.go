package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 测试用的请求结构
type testValidationRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required,password"`
	Age      int    `json:"age" validate:"required,min=1,max=150"`
	Gender   string `json:"gender" validate:"oneof=male female other"`
}

func TestValidateJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        interface{}
		expectedStatus int
		expectAbort    bool
	}{
		{
			name: "有效请求",
			request: testValidationRequest{
				Name:     "张三",
				Email:    "test@example.com",
				Phone:    "13812345678",
				Password: "Pass@123",
				Age:      25,
				Gender:   "male",
			},
			expectedStatus: http.StatusOK,
			expectAbort:    false,
		},
		{
			name: "缺少必填字段",
			request: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectAbort:    true,
		},
		{
			name: "无效的邮箱",
			request: testValidationRequest{
				Name:     "张三",
				Email:    "invalid-email",
				Phone:    "13812345678",
				Password: "Pass@123",
				Age:      25,
				Gender:   "male",
			},
			expectedStatus: http.StatusBadRequest,
			expectAbort:    true,
		},
		{
			name: "无效的手机号",
			request: testValidationRequest{
				Name:     "张三",
				Email:    "test@example.com",
				Phone:    "12345",
				Password: "Pass@123",
				Age:      25,
				Gender:   "male",
			},
			expectedStatus: http.StatusBadRequest,
			expectAbort:    true,
		},
		{
			name: "无效的密码",
			request: testValidationRequest{
				Name:     "张三",
				Email:    "test@example.com",
				Phone:    "13812345678",
				Password: "weak",
				Age:      25,
				Gender:   "male",
			},
			expectedStatus: http.StatusBadRequest,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/test", ValidateJSON(&testValidationRequest{}), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestValidateJSON_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/test", ValidateJSON(&testValidationRequest{}), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("期望状态码 %d, 得到 %d", http.StatusBadRequest, w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != false {
		t.Error("期望 success 为 false")
	}
}

func TestGetValidatedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		setData bool
		expect  bool
	}{
		{
			name:    "成功获取验证后的请求",
			setData: true,
			expect:  true,
		},
		{
			name:    "未找到验证后的请求",
			setData: false,
			expect:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			if tt.setData {
				testData := &testValidationRequest{
					Name:     "测试",
					Email:    "test@example.com",
					Phone:    "13812345678",
					Password: "Pass@123",
					Age:      25,
					Gender:   "male",
				}
				c.Set("validated_request", testData)
			}

			var result testValidationRequest
			got := GetValidatedRequest(c, &result)

			if got != tt.expect {
				t.Errorf("期望返回 %v, 得到 %v", tt.expect, got)
			}

			if tt.setData && result.Name != "测试" {
				t.Error("未能正确获取验证后的数据")
			}
		})
	}
}

func TestValidateQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		rules          map[string]string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name: "有效的查询参数",
			rules: map[string]string{
				"name": "required",
				"age":  "",
			},
			queryParams: map[string]string{
				"name": "张三",
				"age":  "25",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "缺少必填参数",
			rules: map[string]string{
				"name": "required",
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "最小长度验证失败",
			rules: map[string]string{
				"name": "min:5",
			},
			queryParams: map[string]string{
				"name": "ab",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "可选参数为空",
			rules: map[string]string{
				"name": "",
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", ValidateQuery(tt.rules), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("期望状态码 %d, 得到 %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{"有效手机号-138", "13812345678", true},
		{"有效手机号-159", "15912345678", true},
		{"有效手机号-186", "18612345678", true},
		{"太短", "1381234567", false},
		{"太长", "138123456789", false},
		{"不以1开头", "23812345678", false},
		{"第二位无效-0", "10812345678", false},
		{"第二位无效-2", "12812345678", false},
		{"包含非数字字符", "1381234567a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			type phoneRequest struct {
				Phone string `json:"phone" validate:"required,phone"`
			}

			router := gin.New()
			router.POST("/test", ValidateJSON(&phoneRequest{}), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			body, _ := json.Marshal(phoneRequest{Phone: tt.phone})
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			isValid := w.Code == http.StatusOK
			if isValid != tt.valid {
				t.Errorf("手机号 %s: 期望有效性 %v, 得到 %v", tt.phone, tt.valid, isValid)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"有效密码", "Pass@123", true},
		{"有效密码-复杂", "Abc123!@#", true},
		{"太短", "Abc@12", false},
		{"无数字", "Password@", false},
		{"无字母", "12345678@", false},
		{"无特殊字符", "Password123", false},
		{"仅数字", "12345678", false},
		{"仅字母", "abcdefgh", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			type passwordRequest struct {
				Password string `json:"password" validate:"required,password"`
			}

			router := gin.New()
			router.POST("/test", ValidateJSON(&passwordRequest{}), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			body, _ := json.Marshal(passwordRequest{Password: tt.password})
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			isValid := w.Code == http.StatusOK
			if isValid != tt.valid {
				t.Errorf("密码 %s: 期望有效性 %v, 得到 %v", tt.password, tt.valid, isValid)
			}
		})
	}
}

func TestGetErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requestData   interface{}
		expectedTags  []string
		expectSuccess bool
	}{
		{
			name: "required错误",
			requestData: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedTags:  []string{"required"},
			expectSuccess: false,
		},
		{
			name: "min错误",
			requestData: testValidationRequest{
				Name:     "a",
				Email:    "test@example.com",
				Phone:    "13812345678",
				Password: "Pass@123",
				Age:      25,
				Gender:   "male",
			},
			expectedTags:  []string{"min"},
			expectSuccess: false,
		},
		{
			name: "email错误",
			requestData: testValidationRequest{
				Name:     "张三",
				Email:    "invalid",
				Phone:    "13812345678",
				Password: "Pass@123",
				Age:      25,
				Gender:   "male",
			},
			expectedTags:  []string{"email"},
			expectSuccess: false,
		},
		{
			name: "oneof错误",
			requestData: testValidationRequest{
				Name:     "张三",
				Email:    "test@example.com",
				Phone:    "13812345678",
				Password: "Pass@123",
				Age:      25,
				Gender:   "invalid",
			},
			expectedTags:  []string{"oneof"},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/test", ValidateJSON(&testValidationRequest{}), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			body, _ := json.Marshal(tt.requestData)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectSuccess {
				if w.Code != http.StatusOK {
					t.Errorf("期望成功，但得到状态码 %d", w.Code)
				}
			} else {
				if w.Code == http.StatusOK {
					t.Error("期望失败，但请求成功")
				}

				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)

				if errors, ok := resp["errors"].([]interface{}); ok && len(errors) > 0 {
					// 验证错误消息包含中文
					t.Logf("验证错误: %v", errors)
				} else {
					t.Error("期望包含验证错误信息")
				}
			}
		})
	}
}

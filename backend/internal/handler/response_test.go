package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/apierr"
)

func TestRespondJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("发送JSON响应", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			RespondJSON(c, http.StatusOK, gin.H{"message": "success", "data": "test"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
		assert.Contains(t, w.Body.String(), "test")
	})

	t.Run("发送不同状态码", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			RespondJSON(c, http.StatusCreated, gin.H{"id": 123})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "123")
	})
}

func TestRespondSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondSuccess(c, "success", gin.H{"user_id": 123, "name": "Test User"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.Equal(t, true, response["success"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(123), data["user_id"])
	assert.Equal(t, "Test User", data["name"])
}

func TestRespondCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/test", func(c *gin.Context) {
		RespondCreated(c, gin.H{"id": 456, "status": "created"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(201), response["code"])
	assert.Equal(t, "created", response["message"])
}

func TestRespondAPIError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("发送错误响应", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			err := apierr.BadRequest("请求参数无效").
				WithDetails("用户名不能为空").
				WithField("username")
			RespondAPIError(c, err)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, float64(400), response["code"])
		assert.Equal(t, "请求参数无效", response["message"])
		assert.Equal(t, "用户名不能为空", response["details"])
		assert.Equal(t, "username", response["field"])
		assert.NotZero(t, response["timestamp"])
	})

	t.Run("错误响应包含RequestID", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("request_id", "test-req-123")
			c.Next()
		})
		router.GET("/test", func(c *gin.Context) {
			err := apierr.NotFound("资源不存在")
			RespondAPIError(c, err)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, "test-req-123", response["requestId"])
	})
}

func TestRespondError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondError(c, http.StatusBadRequest, "请求无效")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(400), response["code"])
	assert.Equal(t, "请求无效", response["message"])
}

func TestRespondBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondBadRequest(c, "参数验证失败")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(400), response["code"])
	assert.Equal(t, "参数验证失败", response["message"])
}

func TestRespondUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondUnauthorized(c, "未授权访问")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(401), response["code"])
	assert.Equal(t, "未授权访问", response["message"])
}

func TestRespondForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondForbidden(c, "权限不足")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(403), response["code"])
	assert.Equal(t, "权限不足", response["message"])
}

func TestRespondNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondNotFound(c, "资源不存在")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(404), response["code"])
	assert.Equal(t, "资源不存在", response["message"])
}

func TestRespondInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondInternalError(c, "服务器内部错误")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(500), response["code"])
	assert.Equal(t, "服务器内部错误", response["message"])
}

func TestRespondValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		RespondValidationError(c, "email", "邮箱格式无效")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, float64(400), response["code"])
	assert.Equal(t, "邮箱格式无效", response["message"])
	assert.Equal(t, "email", response["field"])
}

func TestBindAndValidate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("绑定成功", func(t *testing.T) {
		router := gin.New()
		router.POST("/test", func(c *gin.Context) {
			var req struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			}
			if err := BindAndValidate(c, &req); err != nil {
				return
			}
			RespondSuccess(c, "success", req)
		})

		body := strings.NewReader(`{"name":"Test","email":"test@example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("绑定失败-无效JSON", func(t *testing.T) {
		router := gin.New()
		router.POST("/test", func(c *gin.Context) {
			var req struct {
				Name string `json:"name"`
			}
			BindAndValidate(c, &req)
		})

		body := strings.NewReader(`{"name":invalid}`) // 无效的JSON
		req := httptest.NewRequest(http.MethodPost, "/test", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, float64(400), response["code"])
		assert.Contains(t, response["message"].(string), "无效的请求格式")
	})
}

func TestParseIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		paramValue  string
		expectError bool
		expectedID  uint64
		errorMsg    string
	}{
		{
			name:        "有效ID",
			paramValue:  "123",
			expectError: false,
			expectedID:  123,
		},
		{
			name:        "无效ID-字母",
			paramValue:  "abc",
			expectError: true,
			errorMsg:    "无效的id格式",
		},
		{
			name:        "无效ID-0",
			paramValue:  "0",
			expectError: true,
			errorMsg:    "id不能为0",
		},
		{
			name:        "无效ID-负数",
			paramValue:  "-1",
			expectError: true,
			errorMsg:    "无效的id格式",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test/:id", func(c *gin.Context) {
				id, err := ParseIDParam(c, "id")
				if err != nil {
					RespondAPIError(c, err.(*apierr.APIError))
					return
				}
				RespondSuccess(c, "success", gin.H{"id": id})
			})

			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.paramValue, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response["message"].(string), tt.errorMsg)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				data := response["extensions"].(map[string]interface{})["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.expectedID), data["id"])
			}
		})
	}
}

func TestParseQueryInt(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		queryValue   string
		defaultValue int
		expected     int
		expectError  bool
	}{
		{
			name:         "有效值",
			queryValue:   "10",
			defaultValue: 5,
			expected:     10,
			expectError:  false,
		},
		{
			name:         "使用默认值",
			queryValue:   "",
			defaultValue: 5,
			expected:     5,
			expectError:  false,
		},
		{
			name:         "无效值-字母",
			queryValue:   "abc",
			defaultValue: 5,
			expected:     0,
			expectError:  true,
		},
		{
			name:         "无效值-负数",
			queryValue:   "-1",
			defaultValue: 5,
			expected:     -1,
			expectError:  false, // 负数是有效的int值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				value, err := ParseQueryInt(c, "page", tt.defaultValue)
				if err != nil {
					RespondAPIError(c, err.(*apierr.APIError))
					return
				}
				RespondSuccess(c, "success", gin.H{"value": value})
			})

			url := "/test"
			if tt.queryValue != "" {
				url += "?page=" + tt.queryValue
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				data := response["extensions"].(map[string]interface{})["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.expected), data["value"])
			}
		})
	}
}

func TestParseQueryUint64(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		queryValue  string
		expected    uint64
		expectError bool
	}{
		{
			name:        "有效值",
			queryValue:  "123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "空值返回0",
			queryValue:  "",
			expected:    0,
			expectError: false,
		},
		{
			name:        "无效值-字母",
			queryValue:  "abc",
			expected:    0,
			expectError: true,
		},
		{
			name:        "无效值-负数",
			queryValue:  "-1",
			expected:    0,
			expectError: true, // uint64不能为负数
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/test", func(c *gin.Context) {
				value, err := ParseQueryUint64(c, "id")
				if err != nil {
					RespondAPIError(c, err.(*apierr.APIError))
					return
				}
				RespondSuccess(c, "success", gin.H{"value": value})
			})

			url := "/test"
			if tt.queryValue != "" {
				url += "?id=" + tt.queryValue
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				data := response["extensions"].(map[string]interface{})["data"].(map[string]interface{})
				assert.Equal(t, float64(tt.expected), data["value"])
			}
		})
	}
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("获取存在的RequestID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("request_id", "test-123")

		requestID := GetRequestID(c)
		assert.Equal(t, "test-123", requestID)
	})

	t.Run("获取不存在的RequestID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		requestID := GetRequestID(c)
		assert.Empty(t, requestID)
	})

	t.Run("RequestID类型错误", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("request_id", 123) // int instead of string

		requestID := GetRequestID(c)
		assert.Empty(t, requestID)
	})
}

func TestAddRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("添加RequestID到响应头", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("request_id", "test-req-456")
			AddRequestID(c)
			c.Next()
		})
		router.GET("/test", func(c *gin.Context) {
			RespondSuccess(c, "success", gin.H{"message": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, "test-req-456", w.Header().Get("X-Request-ID"))
	})

	t.Run("生成新的RequestID", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			AddRequestID(c)
			c.Next()
		})
		router.GET("/test", func(c *gin.Context) {
			requestID := GetRequestID(c)
			RespondSuccess(c, "success", gin.H{"request_id": requestID})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)
		assert.Contains(t, requestID, "req-")
	})
}



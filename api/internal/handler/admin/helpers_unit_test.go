// Package admin provides unit tests for admin handler helper functions.
package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
	"gamelink/internal/repository"
	repoiface "gamelink/internal/repository/interfaces"
	"gamelink/pkg/apierr"
)

// ============================================================================
// parsePagination Tests
// ============================================================================

func TestParsePagination_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name          string
		query         string
		wantPage      int
		wantPageSize  int
		shouldSucceed bool
	}{
		{
			name:          "Default pagination",
			query:         "",
			wantPage:      1,
			wantPageSize:  20,
			shouldSucceed: true,
		},
		{
			name:          "Custom page and pageSize",
			query:         "?page=2&page_size=50",
			wantPage:      2,
			wantPageSize:  50,
			shouldSucceed: true,
		},
		{
			name:          "pageSize alternative param",
			query:         "?page=3&pageSize=30",
			wantPage:      3,
			wantPageSize:  30,
			shouldSucceed: true,
		},
		{
			name:          "page_size takes precedence",
			query:         "?page=1&page_size=25&pageSize=30",
			wantPage:      1,
			wantPageSize:  25,
			shouldSucceed: true,
		},
		{
			name:          "Invalid page",
			query:         "?page=abc",
			shouldSucceed: false,
		},
		{
			name:          "Invalid pageSize",
			query:         "?page_size=xyz",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			page, pageSize, ok := parsePagination(c)

			if tt.shouldSucceed {
				assert.True(t, ok)
				assert.Equal(t, tt.wantPage, page)
				assert.Equal(t, tt.wantPageSize, pageSize)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

// ============================================================================
// Query Helper Function Tests
// ============================================================================

func TestQueryUint64Ptr_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name      string
		param     string
		query     string
		wantValue *uint64
	}{
		{
			name:      "Valid number",
			param:     "user_id",
			query:     "?user_id=123",
			wantValue: func() *uint64 { v := uint64(123); return &v }(),
		},
		{
			name:      "Empty value",
			param:     "user_id",
			query:     "",
			wantValue: nil,
		},
		{
			name:      "Zero value",
			param:     "id",
			query:     "?id=0",
			wantValue: func() *uint64 { v := uint64(0); return &v }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, err := queryUint64Ptr(c, tt.param)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, result)
		})
	}
}

func TestQueryUint64Ptr_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test?user_id=abc", nil)
	c.Request = req

	result, err := queryUint64Ptr(c, "user_id")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestQueryTimePtr_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name      string
		query     string
		wantError bool
	}{
		{
			name:      "RFC3339 format",
			query:     "?date=2023-01-01T00:00:00Z",
			wantError: false,
		},
		{
			name:      "Date only format",
			query:     "?date=2023-01-01",
			wantError: false,
		},
		{
			name:      "Datetime format",
			query:     "?date=2023-01-01 12:00:00",
			wantError: false,
		},
		{
			name:      "Unix timestamp",
			query:     "?date=1672531200",
			wantError: false,
		},
		{
			name:      "Empty value",
			query:     "",
			wantError: false,
		},
		{
			name:      "Invalid format",
			query:     "?date=invalid",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, err := queryTimePtr(c, "date")
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.query != "" && !tt.wantError {
					if tt.query != "" {
						assert.NotNil(t, result)
					}
				}
			}
		})
	}
}

// ============================================================================
// buildUserListOptions Tests
// ============================================================================

func TestBuildUserListOptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Basic query",
			query:   "?page=1&page_size=20",
			wantErr: false,
		},
		{
			name:    "With role filter",
			query:   "?role=user",
			wantErr: false,
		},
		{
			name:    "With multiple roles",
			query:   "?role=user&role=player",
			wantErr: false,
		},
		{
			name:    "With status filter",
			query:   "?status=active",
			wantErr: false,
		},
		{
			name:    "With date range",
			query:   "?date_from=2023-01-01&date_to=2023-12-31",
			wantErr: false,
		},
		{
			name:    "With keyword",
			query:   "?keyword=test",
			wantErr: false,
		},
		{
			name:    "Invalid date_from",
			query:   "?date_from=invalid",
			wantErr: true,
		},
		{
			name:    "Invalid date_to",
			query:   "?date_to=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, ok := buildUserListOptions(c)

			if tt.wantErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.IsType(t, repository.UserListOptions{}, result)
			}
		})
	}
}

// ============================================================================
// buildOrderListOptions Tests
// ============================================================================

func TestBuildOrderListOptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Basic query",
			query:   "?page=1&page_size=20",
			wantErr: false,
		},
		{
			name:    "With status filter",
			query:   "?status=pending",
			wantErr: false,
		},
		{
			name:    "With multiple statuses",
			query:   "?status=pending&status=confirmed",
			wantErr: false,
		},
		{
			name:    "With user_id",
			query:   "?user_id=123",
			wantErr: false,
		},
		{
			name:    "With player_id",
			query:   "?player_id=456",
			wantErr: false,
		},
		{
			name:    "With game_id",
			query:   "?game_id=1",
			wantErr: false,
		},
		{
			name:    "With date range",
			query:   "?date_from=2023-01-01&date_to=2023-12-31",
			wantErr: false,
		},
		{
			name:    "Invalid user_id",
			query:   "?user_id=invalid",
			wantErr: true,
		},
		{
			name:    "Invalid date_from",
			query:   "?date_from=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, ok := buildOrderListOptions(c)

			if tt.wantErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.IsType(t, repoiface.OrderListOptions{}, result)
			}
		})
	}
}

// ============================================================================
// buildPaymentListOptions Tests
// ============================================================================

func TestBuildPaymentListOptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Basic query",
			query:   "?page=1&page_size=20",
			wantErr: false,
		},
		{
			name:    "With status filter",
			query:   "?status=pending",
			wantErr: false,
		},
		{
			name:    "With method filter",
			query:   "?method=alipay",
			wantErr: false,
		},
		{
			name:    "With order_id",
			query:   "?order_id=123",
			wantErr: false,
		},
		{
			name:    "With date range",
			query:   "?date_from=2023-01-01&date_to=2023-12-31",
			wantErr: false,
		},
		{
			name:    "Invalid order_id",
			query:   "?order_id=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, ok := buildPaymentListOptions(c)

			if tt.wantErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.IsType(t, repository.PaymentListOptions{}, result)
			}
		})
	}
}

// ============================================================================
// normalizeOrderStatus Tests
// ============================================================================

func TestNormalizeOrderStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  model.OrderStatus
	}{
		{
			name:  "Lowercase pending",
			input: "pending",
			want:  model.OrderStatusPending,
		},
		{
			name:  "Uppercase PENDING",
			input: "PENDING",
			want:  model.OrderStatusPending,
		},
		{
			name:  "Mixed case Confirmed",
			input: "Confirmed",
			want:  model.OrderStatusConfirmed,
		},
		{
			name:  "Legacy spelling cancelled",
			input: "cancelled",
			want:  model.OrderStatusCanceled,
		},
		{
			name:  "Modern spelling canceled",
			input: "canceled",
			want:  model.OrderStatusCanceled,
		},
		{
			name:  "In progress",
			input: "in_progress",
			want:  model.OrderStatusInProgress,
		},
		{
			name:  "Completed",
			input: "completed",
			want:  model.OrderStatusCompleted,
		},
		{
			name:  "With whitespace",
			input: "  pending  ",
			want:  model.OrderStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeOrderStatus(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

// ============================================================================
// parseCSVParams Tests
// ============================================================================

func TestParseCSVParams(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Single value",
			input:    []string{"user"},
			expected: []string{"user"},
		},
		{
			name:     "Comma-separated values",
			input:    []string{"user,player,admin"},
			expected: []string{"user", "player", "admin"},
		},
		{
			name:     "Multiple comma-separated values",
			input:    []string{"user,player", "admin"},
			expected: []string{"user", "player", "admin"},
		},
		{
			name:     "Values with whitespace",
			input:    []string{" user , player "},
			expected: []string{"user", "player"},
		},
		{
			name:     "Empty strings filtered out",
			input:    []string{"user", "", "player"},
			expected: []string{"user", "player"},
		},
		{
			name:     "Multiple commas",
			input:    []string{"user,,player"},
			expected: []string{"user", "player"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCSVParams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================================
// ensureSlice Tests
// ============================================================================

func TestEnsureSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{
			name:  "Non-nil slice",
			input: []int{1, 2, 3},
			want:  []int{1, 2, 3},
		},
		{
			name:  "Nil slice",
			input: nil,
			want:  []int{},
		},
		{
			name:  "Empty slice",
			input: []int{},
			want:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ensureSlice(tt.input)
			assert.Equal(t, tt.want, result)
		})
	}
}

// ============================================================================
// Validation Tests
// ============================================================================

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		valid bool
	}{
		{
			name:  "Valid phone",
			phone: "13800138000",
			valid: true,
		},
		{
			name:  "Valid phone with leading 1",
			phone: "15912345678",
			valid: true,
		},
		{
			name:  "Invalid - too short",
			phone: "123456789",
			valid: false,
		},
		{
			name:  "Invalid - doesn't start with 1",
			phone: "23800138000",
			valid: false,
		},
		{
			name:  "Invalid - contains letters",
			phone: "1380013800a",
			valid: false,
		},
		{
			name:  "Invalid - empty",
			phone: "",
			valid: false,
		},
		{
			name:  "Invalid - with whitespace",
			phone: "138 0013 8000",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPhone(tt.phone)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		valid bool
	}{
		{
			name:  "Valid email",
			email: "test@example.com",
			valid: true,
		},
		{
			name:  "Valid email with subdomain",
			email: "user@mail.example.com",
			valid: true,
		},
		{
			name:  "Invalid - no @",
			email: "testexample.com",
			valid: false,
		},
		{
			name:  "Invalid - no domain",
			email: "test@",
			valid: false,
		},
		{
			name:  "Invalid - empty",
			email: "",
			valid: false,
		},
		{
			name:  "Invalid - multiple @",
			email: "test@@example.com",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmail(tt.email)
			assert.Equal(t, tt.valid, result)
		})
	}
}

// ============================================================================
// buildOperationLogListOptions Tests
// ============================================================================

func TestBuildOperationLogListOptions_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Basic query",
			query:   "?page=1&page_size=20",
			wantErr: false,
		},
		{
			name:    "With action filter",
			query:   "?action=create",
			wantErr: false,
		},
		{
			name:    "With actor_user_id",
			query:   "?actor_user_id=123",
			wantErr: false,
		},
		{
			name:    "With date range",
			query:   "?date_from=2023-01-01&date_to=2023-12-31",
			wantErr: false,
		},
		{
			name:    "Invalid actor_user_id",
			query:   "?actor_user_id=invalid",
			wantErr: true,
		},
		{
			name:    "Invalid date_from",
			query:   "?date_from=invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req, _ := http.NewRequest("GET", "/test"+tt.query, nil)
			c.Request = req

			result, ok := buildOperationLogListOptions(c)

			if tt.wantErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.IsType(t, repository.OperationLogListOptions{}, result)
			}
		})
	}
}

// ============================================================================
// Error Helper Tests
// ============================================================================

func TestQueryUint64PtrAndRespond(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test?user_id=abc", nil)
	c.Request = req

	result, ok := QueryUint64PtrAndRespond(c, "user_id", "invalid user ID")

	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestQueryTimePtrAndRespond(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test?date=invalid", nil)
	c.Request = req

	result, ok := QueryTimePtrAndRespond(c, "date", "invalid date")

	assert.False(t, ok)
	assert.Nil(t, result)
}

// ============================================================================
// Time parsing helper tests
// ============================================================================

func TestParseRFC3339Ptr(t *testing.T) {
	tests := []struct {
		name      string
		input     *string
		wantError bool
	}{
		{
			name:      "Valid RFC3339",
			input:     func() *string { s := "2023-01-01T00:00:00Z"; return &s }(),
			wantError: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			wantError: false,
		},
		{
			name:      "Empty string",
			input:     func() *string { s := ""; return &s }(),
			wantError: false, // Empty string returns nil, nil
		},
		{
			name:      "Invalid format",
			input:     func() *string { s := "invalid"; return &s }(),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseRFC3339Ptr(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.input != nil && *tt.input != "" {
					assert.NotNil(t, result)
					assert.IsType(t, time.Time{}, *result)
				}
			}
		})
	}
}

// ============================================================================
// Response Helper Tests
// ============================================================================

func TestRespondSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondSuccess(c, map[string]string{"message": "test"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test")
}

func TestRespondCreated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondCreated(c, map[string]string{"id": "123"})

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "123")
}

func TestRespondDeleted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondDeleted(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRespondUpdated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondUpdated(c, map[string]string{"status": "updated"})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestRespondBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondBadRequest(c, "bad request")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "bad request")
}

func TestRespondError_ApiErr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, apierr.ErrNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRespondError_RepositoryNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	respondError(c, repository.ErrNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestParseIDAndRespond_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req, _ := http.NewRequest("GET", "/test/123", nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	id, ok := ParseIDAndRespond(c, "id")

	assert.True(t, ok)
	assert.Equal(t, uint64(123), id)
}

func TestParseIDAndRespond_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/test/invalid", nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	id, ok := ParseIDAndRespond(c, "id")

	assert.False(t, ok)
	assert.Equal(t, uint64(0), id)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

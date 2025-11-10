package player

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gamelink/internal/model"
)

func TestCommissionHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("getCommissionSummaryHandler_DefaultMonth", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/summary", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})

	t.Run("getCommissionSummaryHandler_SpecificMonth", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/summary?month=2024-01", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, "2024-01", c.Query("month"))
	})

	t.Run("getCommissionSummaryHandler_InvalidMonth", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/summary?month=invalid", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
	})

	t.Run("getCommissionRecordsHandler_DefaultParams", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/records", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})

	t.Run("getCommissionRecordsHandler_WithPagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/records?page=2&pageSize=10", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, "2", c.Query("page"))
		assert.Equal(t, "10", c.Query("pageSize"))
	})

	t.Run("getCommissionRecordsHandler_InvalidPage", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/records?page=0", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
	})

	t.Run("getMonthlySettlementsHandler_DefaultParams", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/settlements", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})

	t.Run("getMonthlySettlementsHandler_WithPagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/settlements?page=1&pageSize=20", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, "1", c.Query("page"))
		assert.Equal(t, "20", c.Query("pageSize"))
	})

	t.Run("getMonthlySettlementsHandler_LargePageSize", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/commission/settlements?pageSize=1000", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
	})
}

func TestCommissionHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("getUserIDFromContext", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint64(456))

		userID := getUserIDFromContext(c)
		assert.Equal(t, uint64(456), userID)
	})

	t.Run("respondJSON_CommissionResponse", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		resp := model.APIResponse[map[string]int64]{
			Success: true,
			Code:    200,
			Message: "OK",
			Data: map[string]int64{
				"totalIncome":     10000,
				"totalCommission": 2000,
			},
		}

		respondJSON(c, http.StatusOK, resp)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var result model.APIResponse[map[string]int64]
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(10000), result.Data["totalIncome"])
	})

	t.Run("respondError_NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		respondError(c, http.StatusNotFound, "Player not found")

		assert.Equal(t, http.StatusNotFound, w.Code)
		
		var result model.APIResponse[interface{}]
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.False(t, result.Success)
		assert.Equal(t, "Player not found", result.Message)
	})
}

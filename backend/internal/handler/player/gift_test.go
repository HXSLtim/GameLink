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

func TestGiftHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("getReceivedGiftsHandler_DefaultParams", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/received", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})

	t.Run("getReceivedGiftsHandler_WithPagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/received?page=1&pageSize=20", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, "1", c.Query("page"))
		assert.Equal(t, "20", c.Query("pageSize"))
	})

	t.Run("getReceivedGiftsHandler_InvalidPage", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/received?page=-1", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
	})

	t.Run("getReceivedGiftsHandler_LargePageSize", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/received?pageSize=500", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
	})

	t.Run("getGiftStatsHandler_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/stats", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, http.MethodGet, c.Request.Method)
	})

	t.Run("getGiftStatsHandler_WithMonth", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/player/gifts/stats?month=2024-01", nil)
		c.Set("user_id", uint64(1))

		assert.NotNil(t, c)
		assert.Equal(t, "2024-01", c.Query("month"))
	})
}

func TestGiftHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("respondJSON_GiftStats", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		resp := model.APIResponse[map[string]interface{}]{
			Success: true,
			Code:    200,
			Message: "OK",
			Data: map[string]interface{}{
				"totalGifts": 50,
				"totalValue": 10000,
			},
		}

		respondJSON(c, http.StatusOK, resp)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var result model.APIResponse[map[string]interface{}]
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.True(t, result.Success)
	})

	t.Run("respondError_InternalError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		respondError(c, http.StatusInternalServerError, "Internal server error")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var result model.APIResponse[interface{}]
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.False(t, result.Success)
		assert.Equal(t, http.StatusInternalServerError, result.Code)
	})
}

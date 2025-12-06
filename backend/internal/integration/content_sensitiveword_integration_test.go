package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
	"gamelink/internal/model"
	sensitivewordrepo "gamelink/internal/repository/sensitiveword"
	"gamelink/internal/service/sensitiveword"
	"gamelink/pkg/testutil"
)

// TestSensitiveWordCRUDFlow tests the complete sensitive word CRUD workflow
// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5
func TestSensitiveWordCRUDFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateSensitiveWordModels(t, db)

	// Setup service and handler
	repo := sensitivewordrepo.NewSensitiveWordRepository(db)
	svc := sensitiveword.NewSensitiveWordService(repo)
	handler := adminhandler.NewSensitiveWordHandler(svc)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeSensitiveWordAdminAuthMiddleware(1)
	api.Use(auth)
	api.GET("/sensitive-words", handler.ListSensitiveWords)
	api.POST("/sensitive-words", handler.AddSensitiveWord)
	api.PUT("/sensitive-words/:id", handler.UpdateSensitiveWord)
	api.DELETE("/sensitive-words/:id", handler.DeleteSensitiveWord)
	api.POST("/reviews/detect-sensitive", handler.DetectSensitiveWords)

	var createdWordID uint64

	t.Run("Add sensitive word", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "敏感词1",
			"category": "political",
			"severity": "high",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusCreated, resp.Code)

		var result apiResp[sensitiveword.SensitiveWordDTO]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, "敏感词1", result.Data.Word)
		assert.Equal(t, model.SensitiveWordCategoryPolitical, result.Data.Category)
		assert.Equal(t, model.SensitiveWordSeverityHigh, result.Data.Severity)
		createdWordID = result.Data.ID
	})

	t.Run("Add another sensitive word", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "广告词",
			"category": "advertising",
			"severity": "low",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("Add duplicate word should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "敏感词1",
			"category": "political",
			"severity": "high",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("List sensitive words", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/sensitive-words", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.ListSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(2), result.Data.Total)
	})

	t.Run("Filter by category", func(t *testing.T) {
		resp := doJSON(router, http.MethodGet, "/api/v1/admin/sensitive-words?category=political", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.ListSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.Data.Total)
	})

	t.Run("Update sensitive word", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "敏感词1更新",
			"category": "violent",
			"severity": "medium",
		}
		resp := doJSON(router, http.MethodPut, "/api/v1/admin/sensitive-words/"+uintToStr(createdWordID), payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Delete sensitive word", func(t *testing.T) {
		resp := doJSON(router, http.MethodDelete, "/api/v1/admin/sensitive-words/"+uintToStr(createdWordID), nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		// Verify deletion
		listResp := doJSON(router, http.MethodGet, "/api/v1/admin/sensitive-words", nil, "")
		var result apiResp[sensitiveword.ListSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(listResp.Body.Bytes(), &result))
		assert.Equal(t, int64(1), result.Data.Total)
	})
}

// TestSensitiveWordDetectionFlow tests the sensitive word detection functionality
// Validates: Requirements 1.5, 4.5
func TestSensitiveWordDetectionFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateSensitiveWordModels(t, db)

	// Setup service and handler
	repo := sensitivewordrepo.NewSensitiveWordRepository(db)
	svc := sensitiveword.NewSensitiveWordService(repo)
	handler := adminhandler.NewSensitiveWordHandler(svc)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeSensitiveWordAdminAuthMiddleware(1)
	api.Use(auth)
	api.POST("/sensitive-words", handler.AddSensitiveWord)
	api.POST("/reviews/detect-sensitive", handler.DetectSensitiveWords)

	// Add test sensitive words
	words := []map[string]interface{}{
		{"word": "违禁词", "category": "political", "severity": "high"},
		{"word": "广告", "category": "advertising", "severity": "low"},
		{"word": "暴力", "category": "violent", "severity": "medium"},
	}
	for _, w := range words {
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", w, "")
		require.Equal(t, http.StatusCreated, resp.Code)
	}

	t.Run("Detect sensitive words in content", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "这是一段包含违禁词和广告的内容",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/detect-sensitive", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.DetectSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.True(t, result.Data.HasSensitiveWords)
		assert.Len(t, result.Data.DetectedWords, 2)
	})

	t.Run("No sensitive words in clean content", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "这是一段正常的内容",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/detect-sensitive", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.DetectSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.False(t, result.Data.HasSensitiveWords)
		assert.Empty(t, result.Data.DetectedWords)
	})

	t.Run("Empty content returns no sensitive words", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/detect-sensitive", payload, "")
		// Empty content should fail validation
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Detect multiple occurrences", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "广告广告广告",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/detect-sensitive", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.DetectSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.True(t, result.Data.HasSensitiveWords)
		// Should detect "广告" with multiple positions
		assert.Len(t, result.Data.DetectedWords, 1)
		assert.Len(t, result.Data.DetectedWords[0].Positions, 3)
	})

	t.Run("Highlighted content shows markers", func(t *testing.T) {
		payload := map[string]interface{}{
			"content": "包含暴力内容",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/reviews/detect-sensitive", payload, "")
		assert.Equal(t, http.StatusOK, resp.Code)

		var result apiResp[sensitiveword.DetectSensitiveWordsResponse]
		require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &result))
		assert.True(t, result.Success)
		assert.Contains(t, result.Data.HighlightedContent, "**")
	})
}

// TestSensitiveWordValidationFlow tests input validation
// Validates: Requirements 4.2, 4.3
func TestSensitiveWordValidationFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateSensitiveWordModels(t, db)

	// Setup service and handler
	repo := sensitivewordrepo.NewSensitiveWordRepository(db)
	svc := sensitiveword.NewSensitiveWordService(repo)
	handler := adminhandler.NewSensitiveWordHandler(svc)

	router := gin.New()
	api := router.Group("/api/v1/admin")
	auth := fakeSensitiveWordAdminAuthMiddleware(1)
	api.Use(auth)
	api.POST("/sensitive-words", handler.AddSensitiveWord)

	t.Run("Invalid category should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "测试词",
			"category": "invalid_category",
			"severity": "high",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Invalid severity should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "测试词",
			"category": "political",
			"severity": "invalid_severity",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Empty word should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"word":     "",
			"category": "political",
			"severity": "high",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Missing required fields should fail", func(t *testing.T) {
		payload := map[string]interface{}{
			"word": "测试词",
		}
		resp := doJSON(router, http.MethodPost, "/api/v1/admin/sensitive-words", payload, "")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func fakeSensitiveWordAdminAuthMiddleware(adminID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", adminID)
		c.Set("userID", adminID)
		c.Set("admin_id", adminID)
		c.Set("adminID", adminID)
		c.Next()
	}
}

func migrateSensitiveWordModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.SensitiveWord{},
	); err != nil {
		t.Fatalf("migrate sensitive word models: %v", err)
	}
}

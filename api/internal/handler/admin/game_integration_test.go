// Package admin provides integration tests for game admin handler.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"
)

// TestGameHandler_ListGames_Success tests successful game listing with pagination.
func TestGameHandler_ListGames_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test games
	_ = integration.CreateTestGame(t, helper.db, "League of Legends")
	_ = integration.CreateTestGame(t, helper.db, "Dota 2")
	_ = integration.CreateTestGame(t, helper.db, "Heroes of the Storm")

	// Test listing games
	w := helper.MakeRequest("GET", "/admin/games?page=1&page_size=2", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID    uint64 `json:"id"`
			Name  string `json:"name"`
			Slug  string `json:"slug"`
			Image string `json:"image"`
		} `json:"data"`
		Pagination struct {
			Page       int `json:"page"`
			PageSize   int `json:"page_size"`
			Total      int `json:"total"`
			TotalPages int `json:"total_pages"`
		} `json:"pagination"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify response structure
	assert.True(t, response.Success)
	assert.Equal(t, 2, len(response.Data)) // page_size=2
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 2, response.Pagination.PageSize)
	assert.Equal(t, 3, response.Pagination.Total)
}

// TestGameHandler_ListGames_WithKeyword tests game listing with keyword filter.
func TestGameHandler_ListGames_WithKeyword(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test games
	integration.CreateTestGame(t, helper.db, "League of Legends")
	integration.CreateTestGame(t, helper.db, "Dota 2")
	integration.CreateTestGame(t, helper.db, "Legends of Runeterra")

	// Test searching with keyword "Legends"
	w := helper.MakeRequest("GET", "/admin/games?keyword=Legends", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    []struct {
			ID   uint64 `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, 2, len(response.Data)) // "League of Legends" and "Legends of Runeterra"
}

// TestGameHandler_GetGame_Success tests successful game retrieval.
func TestGameHandler_GetGame_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test game
	game := integration.CreateTestGame(t, helper.db, "League of Legends")

	// Test getting the game
	w := helper.MakeRequest("GET", fmt.Sprintf("/admin/games/%d", game.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "League of Legends", data["name"])
	assert.Equal(t, "League of Legends", data["key"])
}

// TestGameHandler_GetGame_NotFound tests getting a non-existent game.
func TestGameHandler_GetGame_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Test getting non-existent game
	w := helper.MakeRequest("GET", "/admin/games/999999", nil, true)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
}

// TestGameHandler_CreateGame_Success tests successful game creation.
func TestGameHandler_CreateGame_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"key":         "lol",
		"name":        "League of Legends",
		"category":    "moba",
		"icon_url":    "https://example.com/lol.png",
		"description": "A popular MOBA game",
	}

	w := helper.MakeRequest("POST", "/admin/games", payload, true)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "League of Legends", data["name"])
	assert.Equal(t, "lol", data["key"])
	assert.Equal(t, "moba", data["category"])

	// Verify game exists in database
	var count int64
	helper.db.Model(&model.Game{}).Where("key = ?", "lol").Count(&count)
	assert.Equal(t, int64(1), count)
}

// TestGameHandler_CreateGame_ValidationError tests game creation with validation errors.
func TestGameHandler_CreateGame_ValidationError(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	testCases := []struct {
		name        string
		payload     map[string]interface{}
		expectError bool
	}{
		{
			name:        "missing required fields",
			payload:     map[string]interface{}{},
			expectError: true,
		},
		{
			name: "missing key",
			payload: map[string]interface{}{
				"name": "Test Game",
			},
			expectError: true,
		},
		{
			name: "missing name",
			payload: map[string]interface{}{
				"key": "test",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := helper.MakeRequest("POST", "/admin/games", tc.payload, true)
			if tc.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			}
		})
	}
}

// TestGameHandler_CreateGame_DuplicateKey tests game creation with duplicate key.
func TestGameHandler_CreateGame_DuplicateKey(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create first game
	integration.CreateTestGame(t, helper.db, "lol")

	// Try to create duplicate
	payload := map[string]interface{}{
		"key":  "lol",
		"name": "League of Legends Duplicate",
	}

	w := helper.MakeRequest("POST", "/admin/games", payload, true)
	// Should return error (conflict or bad request depending on service implementation)
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

// TestGameHandler_UpdateGame_Success tests successful game update.
func TestGameHandler_UpdateGame_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game := integration.CreateTestGame(t, helper.db, "League of Legends")

	payload := map[string]interface{}{
		"key":         "lol-updated",
		"name":        "League of Legends Updated",
		"category":    "moba-updated",
		"icon_url":    "https://example.com/lol-updated.png",
		"description": "Updated description",
	}

	w := helper.MakeRequest("PUT", fmt.Sprintf("/admin/games/%d", game.ID), payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "League of Legends Updated", data["name"])
	assert.Equal(t, "lol-updated", data["key"])
	assert.Equal(t, "moba-updated", data["category"])

	// Verify update in database
	var updatedGame model.Game
	err = helper.db.First(&updatedGame, game.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "League of Legends Updated", updatedGame.Name)
	assert.Equal(t, "lol-updated", updatedGame.Key)
}

// TestGameHandler_UpdateGame_NotFound tests updating a non-existent game.
func TestGameHandler_UpdateGame_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"key":  "test",
		"name": "Test Game",
	}

	w := helper.MakeRequest("PUT", "/admin/games/999999", payload, true)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// TestGameHandler_DeleteGame_Success tests successful game deletion.
func TestGameHandler_DeleteGame_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game := integration.CreateTestGame(t, helper.db, "League of Legends")

	// Delete the game
	w := helper.MakeRequest("DELETE", fmt.Sprintf("/admin/games/%d", game.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify game is deleted from database
	var deletedGame model.Game
	err = helper.db.First(&deletedGame, game.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// TestGameHandler_DeleteGame_NotFound tests deleting a non-existent game.
func TestGameHandler_DeleteGame_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	w := helper.MakeRequest("DELETE", "/admin/games/999999", nil, true)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// TestGameHandler_BatchDeleteGames_Success tests successful batch deletion.
func TestGameHandler_BatchDeleteGames_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game1 := integration.CreateTestGame(t, helper.db, "Game 1")
	game2 := integration.CreateTestGame(t, helper.db, "Game 2")
	game3 := integration.CreateTestGame(t, helper.db, "Game 3")

	payload := map[string]interface{}{
		"gameIds": []string{
			fmt.Sprintf("%d", game1.ID),
			fmt.Sprintf("%d", game2.ID),
			fmt.Sprintf("%d", game3.ID),
		},
	}

	w := helper.MakeRequest("POST", "/admin/games/batch/delete", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	deletedCount := int64(data["deleted"].(float64))
	assert.Equal(t, int64(3), deletedCount)

	// Verify all games are deleted
	var count int64
	helper.db.Model(&model.Game{}).Where("id IN ?", []uint64{game1.ID, game2.ID, game3.ID}).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestGameHandler_BatchDeleteGames_InvalidID tests batch deletion with invalid IDs.
func TestGameHandler_BatchDeleteGames_InvalidID(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"gameIds": []string{"invalid", "not-a-number"},
	}

	w := helper.MakeRequest("POST", "/admin/games/batch/delete", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// TestGameHandler_BatchDeleteGames_EmptyList tests batch deletion with empty list.
func TestGameHandler_BatchDeleteGames_EmptyList(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"gameIds": []string{},
	}

	w := helper.MakeRequest("POST", "/admin/games/batch/delete", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGameHandler_ListGameLogs_Success tests successful game logs listing.
func TestGameHandler_ListGameLogs_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game := integration.CreateTestGame(t, helper.db, "League of Legends")

	// Create some operation logs
	log1 := &model.OperationLog{
		EntityType:  "game",
		EntityID:    game.ID,
		Action:      "create",
		ActorUserID: &helper.adminUserID,
		Reason:      "Test log entry 1",
	}
	log2 := &model.OperationLog{
		EntityType:  "game",
		EntityID:    game.ID,
		Action:      "update",
		ActorUserID: &helper.adminUserID,
		Reason:      "Test log entry 2",
	}
	require.NoError(t, helper.db.Create(log1).Error)
	require.NoError(t, helper.db.Create(log2).Error)

	// Test listing logs
	w := helper.MakeRequest("GET", fmt.Sprintf("/admin/games/%d/logs?page=1&page_size=10", game.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	items := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 2)
}

// TestGameHandler_ListGameLogs_WithFilters tests game logs listing with filters.
func TestGameHandler_ListGameLogs_WithFilters(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game := integration.CreateTestGame(t, helper.db, "League of Legends")

	// Create logs with different actions
	createLog := &model.OperationLog{
		EntityType:  "game",
		EntityID:    game.ID,
		Action:      "create",
		ActorUserID: &helper.adminUserID,
	}
	deleteLog := &model.OperationLog{
		EntityType:  "game",
		EntityID:    game.ID,
		Action:      "delete",
		ActorUserID: &helper.adminUserID,
	}
	require.NoError(t, helper.db.Create(createLog).Error)
	require.NoError(t, helper.db.Create(deleteLog).Error)

	// Filter by action "create"
	w := helper.MakeRequest("GET", fmt.Sprintf("/admin/games/%d/logs?action=create", game.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	items := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 1)

	// Verify filtered items have action="create"
	for _, item := range items {
		logItem := item.(map[string]interface{})
		assert.Equal(t, "create", logItem["action"])
	}
}

// TestGameHandler_ListGameLogs_GameNotFound tests logs listing for non-existent game.
func TestGameHandler_ListGameLogs_GameNotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	w := helper.MakeRequest("GET", "/admin/games/999999/logs", nil, true)
	// Should return either empty list or not found depending on implementation
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, w.Code)
}

// TestGameHandler_ListGames_EmptyDatabase tests game listing when database is empty.
func TestGameHandler_ListGames_EmptyDatabase(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Don't create any games, query directly
	w := helper.MakeRequest("GET", "/admin/games", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// API returns data as array directly, not as object with items
	items := response["data"].([]interface{})
	assert.Equal(t, 0, len(items))

	pagination := response["pagination"].(map[string]interface{})
	assert.Equal(t, float64(0), pagination["total"])
}

// TestGameHandler_ListGames_PaginationTests tests various pagination scenarios.
func TestGameHandler_ListGames_PaginationTests(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create 25 games
	for i := 1; i <= 25; i++ {
		integration.CreateTestGame(t, helper.db, fmt.Sprintf("Game %d", i))
	}

	testCases := []struct {
		name          string
		page          int
		pageSize      int
		expectedItems int
		expectedTotal int
	}{
		{
			name:          "first page",
			page:          1,
			pageSize:      10,
			expectedItems: 10,
			expectedTotal: 25,
		},
		{
			name:          "second page",
			page:          2,
			pageSize:      10,
			expectedItems: 10,
			expectedTotal: 25,
		},
		{
			name:          "last partial page",
			page:          3,
			pageSize:      10,
			expectedItems: 5,
			expectedTotal: 25,
		},
		{
			name:          "beyond last page",
			page:          10,
			pageSize:      10,
			expectedItems: 0,
			expectedTotal: 25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := helper.MakeRequest("GET", fmt.Sprintf("/admin/games?page=%d&page_size=%d", tc.page, tc.pageSize), nil, true)
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// API returns data as array directly, not as object with items
			items := response["data"].([]interface{})
			pagination := response["pagination"].(map[string]interface{})

			assert.Equal(t, tc.expectedItems, len(items))
			assert.Equal(t, float64(tc.expectedTotal), pagination["total"])
			assert.Equal(t, float64(tc.page), pagination["page"])
			assert.Equal(t, float64(tc.pageSize), pagination["page_size"])
		})
	}
}

// TestGameHandler_UpdateGame_ValidationTests tests update with various validation scenarios.
func TestGameHandler_UpdateGame_ValidationTests(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	game := integration.CreateTestGame(t, helper.db, "Original Name")

	testCases := []struct {
		name       string
		payload    map[string]interface{}
		expectCode int
	}{
		{
			name: "empty payload (with existing values)",
			payload: map[string]interface{}{
				"key":  game.Key,
				"name": game.Name,
			},
			expectCode: http.StatusOK,
		},
		{
			name: "partial update - only name (should fail - key is required)",
			payload: map[string]interface{}{
				"name": "Updated Name Only",
			},
			expectCode: http.StatusBadRequest, // Key is required by service layer
		},
		{
			name: "partial update - only key (should fail - name is required)",
			payload: map[string]interface{}{
				"key": "updated-key",
			},
			expectCode: http.StatusBadRequest, // Name is required by service layer
		},
		{
			name: "update all fields",
			payload: map[string]interface{}{
				"key":         "updated-key",
				"name":        "Updated Name",
				"category":    "updated-category",
				"icon_url":    "https://example.com/updated.png",
				"description": "Updated description",
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := helper.MakeRequest("PUT", fmt.Sprintf("/admin/games/%d", game.ID), tc.payload, true)
			assert.Equal(t, tc.expectCode, w.Code)
		})
	}
}

// TestGameHandler_IDParameterParsing tests ID parameter parsing edge cases.
func TestGameHandler_IDParameterParsing(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	testCases := []struct {
		name       string
		endpoint   string
		expectCode int
	}{
		{
			name:       "zero ID",
			endpoint:   "/admin/games/0",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "negative ID",
			endpoint:   "/admin/games/-1",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "non-numeric ID",
			endpoint:   "/admin/games/abc",
			expectCode: http.StatusBadRequest,
		},
		{
			name:       "float ID",
			endpoint:   "/admin/games/1.5",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := helper.MakeRequest("GET", tc.endpoint, nil, true)
			assert.Equal(t, tc.expectCode, w.Code)
		})
	}
}

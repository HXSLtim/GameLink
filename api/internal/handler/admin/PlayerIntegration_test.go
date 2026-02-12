// Package admin provides integration tests for admin handlers.
package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/internal/service/integration"
)

// TestPlayerHandler_ListPlayers_Success tests listing players with pagination.
func TestPlayerHandler_ListPlayers_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test players
	for i := 0; i < 5; i++ {
		user := integration.CreateUniqueTestUser(t, helper.db, "list_player")
		player := integration.CreateTestPlayer(t, helper.db, user)
		player.VerificationStatus = model.VerificationVerified
		helper.db.Save(player)
	}

	// Test listing players
	w := helper.MakeRequest("GET", "/admin/players?page=1&pageSize=10", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 5)
}

// TestPlayerHandler_ListPlayers_WithFilter tests listing players with status filter.
func TestPlayerHandler_ListPlayers_WithFilter(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create verified players
	for i := 0; i < 3; i++ {
		user := integration.CreateUniqueTestUser(t, helper.db, "verified_player")
		player := integration.CreateTestPlayer(t, helper.db, user)
		player.VerificationStatus = model.VerificationVerified
		helper.db.Save(player)
	}

	// Create pending players
	for i := 0; i < 2; i++ {
		user := integration.CreateUniqueTestUser(t, helper.db, "pending_player")
		player := integration.CreateTestPlayer(t, helper.db, user)
		player.VerificationStatus = model.VerificationPending
		helper.db.Save(player)
	}

	// Test filtering by verified status
	w := helper.MakeRequest("GET", "/admin/players?status=verified&page=1&pageSize=10", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.Equal(t, 3, len(items))
}

// TestPlayerHandler_ListPlayers_WithKeyword tests listing players with keyword search.
func TestPlayerHandler_ListPlayers_WithKeyword(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create player with specific nickname
	user := integration.CreateUniqueTestUser(t, helper.db, "keyword_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	player.Nickname = "DiamondPlayer123"
	helper.db.Save(player)

	// Create another player
	user2 := integration.CreateUniqueTestUser(t, helper.db, "other_player")
	player2 := integration.CreateTestPlayer(t, helper.db, user2)
	player2.Nickname = "GoldPlayer456"
	helper.db.Save(player2)

	// Test keyword search
	w := helper.MakeRequest("GET", "/admin/players?keyword=Diamond&page=1&pageSize=10", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	assert.GreaterOrEqual(t, len(items), 1)
}

// TestPlayerHandler_CreatePlayer_Success tests creating a new player.
func TestPlayerHandler_CreatePlayer_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test user and game
	user := integration.CreateUniqueTestUser(t, helper.db, "new_player_user")
	game := integration.CreateTestGame(t, helper.db, "test_game")

	payload := map[string]interface{}{
		"user_id":             user.ID,
		"nickname":            "TestPlayer",
		"bio":                 "Test bio",
		"rank":                "Diamond",
		"hourly_rate_cents":   5000,
		"main_game_id":        game.ID,
		"verification_status": "pending",
	}

	w := helper.MakeRequest("POST", "/admin/players", payload, true)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotZero(t, data["id"])
	assert.Equal(t, "TestPlayer", data["nickname"])
	assert.Equal(t, "pending", data["verificationStatus"])
}

// TestPlayerHandler_CreatePlayer_ValidationError tests validation errors.
func TestPlayerHandler_CreatePlayer_ValidationError(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"nickname": "TestPlayer",
	}

	w := helper.MakeRequest("POST", "/admin/players", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Invalid verification status
	payload2 := map[string]interface{}{
		"user_id":             uint64(123),
		"verification_status": "invalid_status",
	}

	w2 := helper.MakeRequest("POST", "/admin/players", payload2, true)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// TestPlayerHandler_GetPlayer_Success tests getting a player by ID.
func TestPlayerHandler_GetPlayer_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "get_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	player.Nickname = "GetTestPlayer"
	player.HourlyRateCents = 6000
	helper.db.Save(player)

	w := helper.MakeRequest("GET", "/admin/players/"+uint64ToString(player.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, player.ID, uint64(data["id"].(float64)))
	assert.Equal(t, "GetTestPlayer", data["nickname"])
	assert.Equal(t, int64(6000), int64(data["hourlyRateCents"].(float64)))
}

// TestPlayerHandler_GetPlayer_NotFound tests getting a non-existent player.
func TestPlayerHandler_GetPlayer_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	w := helper.MakeRequest("GET", "/admin/players/99999999", nil, true)
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response["success"].(bool))
}

// TestPlayerHandler_GetPlayer_InvalidID tests getting with invalid ID.
func TestPlayerHandler_GetPlayer_InvalidID(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	w := helper.MakeRequest("GET", "/admin/players/invalid", nil, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPlayerHandler_UpdatePlayer_Success tests updating a player.
func TestPlayerHandler_UpdatePlayer_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "update_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	game := integration.CreateTestGame(t, helper.db, "update_game")

	payload := map[string]interface{}{
		"nickname":            "UpdatedPlayer",
		"bio":                 "Updated bio",
		"rank":                "Master",
		"hourly_rate_cents":   8000,
		"main_game_id":        game.ID,
		"verification_status": "verified",
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID), payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "UpdatedPlayer", data["nickname"])
	assert.Equal(t, "Updated bio", data["bio"])
	assert.Equal(t, int64(8000), int64(data["hourlyRateCents"].(float64)))
}

// TestPlayerHandler_UpdatePlayer_NotFound tests updating a non-existent player.
func TestPlayerHandler_UpdatePlayer_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"nickname":            "Test",
		"verification_status": "verified",
	}

	w := helper.MakeRequest("PUT", "/admin/players/99999999", payload, true)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestPlayerHandler_DeletePlayer_Success tests deleting a player.
func TestPlayerHandler_DeletePlayer_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "delete_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	w := helper.MakeRequest("DELETE", "/admin/players/"+uint64ToString(player.ID), nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	// Verify player is deleted
	var count int64
	helper.db.Model(&model.Player{}).Where("id = ?", player.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

// TestPlayerHandler_DeletePlayer_NotFound tests deleting a non-existent player.
func TestPlayerHandler_DeletePlayer_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	w := helper.MakeRequest("DELETE", "/admin/players/99999999", nil, true)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestPlayerHandler_UpdatePlayerVerification_Success tests updating player verification status.
func TestPlayerHandler_UpdatePlayerVerification_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "verify_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	player.VerificationStatus = model.VerificationPending
	helper.db.Save(player)

	payload := map[string]interface{}{
		"verification_status": "verified",
		"remark":              "Approved by admin",
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/verification", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "verified", data["verificationStatus"])

	// Verify in database
	var updatedPlayer model.Player
	helper.db.First(&updatedPlayer, player.ID)
	assert.Equal(t, model.VerificationVerified, updatedPlayer.VerificationStatus)
	assert.NotNil(t, updatedPlayer.VerifiedBy)
}

// TestPlayerHandler_UpdatePlayerVerification_Reject tests rejecting a player.
func TestPlayerHandler_UpdatePlayerVerification_Reject(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "reject_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	player.VerificationStatus = model.VerificationPending
	helper.db.Save(player)

	payload := map[string]interface{}{
		"verification_status": "rejected",
		"remark":              "Incomplete profile",
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/verification", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "rejected", data["verificationStatus"])
}

// TestPlayerHandler_UpdatePlayerVerification_InvalidStatus tests invalid verification status.
func TestPlayerHandler_UpdatePlayerVerification_InvalidStatus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "invalid_verify_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	payload := map[string]interface{}{
		"verification_status": "invalid_status",
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/verification", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPlayerHandler_UpdatePlayerGames_Success tests updating player's main game.
func TestPlayerHandler_UpdatePlayerGames_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player and games
	user := integration.CreateUniqueTestUser(t, helper.db, "game_update_player")
	player := integration.CreateTestPlayer(t, helper.db, user)
	game1 := integration.CreateTestGame(t, helper.db, "game1")
	game2 := integration.CreateTestGame(t, helper.db, "game2")

	// Set initial game
	player.MainGameID = game1.ID
	helper.db.Save(player)

	// Update to new game
	payload := map[string]interface{}{
		"main_game_id": game2.ID,
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/games", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, game2.ID, uint64(data["mainGameId"].(float64)))
}

// TestPlayerHandler_UpdatePlayerGames_NotFound tests updating games for non-existent player.
func TestPlayerHandler_UpdatePlayerGames_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"main_game_id": uint64(123),
	}

	w := helper.MakeRequest("PUT", "/admin/players/99999999/games", payload, true)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestPlayerHandler_UpdatePlayerSkillTags_Success tests updating player skill tags.
func TestPlayerHandler_UpdatePlayerSkillTags_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "tags_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	payload := map[string]interface{}{
		"tags": []string{"friendly", "patient", "skilled"},
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/skill-tags", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

// TestPlayerHandler_UpdatePlayerSkillTags_EmptyTags tests updating with empty tags.
func TestPlayerHandler_UpdatePlayerSkillTags_EmptyTags(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "empty_tags_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	payload := map[string]interface{}{
		"tags": []string{},
	}

	w := helper.MakeRequest("PUT", "/admin/players/"+uint64ToString(player.ID)+"/skill-tags", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

// TestPlayerHandler_UpdatePlayerSkillTags_NotFound tests updating tags for non-existent player.
func TestPlayerHandler_UpdatePlayerSkillTags_NotFound(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"tags": []string{"test"},
	}

	w := helper.MakeRequest("PUT", "/admin/players/99999999/skill-tags", payload, true)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestPlayerHandler_ListPlayerLogs_Success tests listing player operation logs.
func TestPlayerHandler_ListPlayerLogs_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "logs_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	w := helper.MakeRequest("GET", "/admin/players/"+uint64ToString(player.ID)+"/logs?page=1&pageSize=10", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	items := data["items"].([]interface{})
	// May be empty or have some logs
	assert.NotNil(t, items)
}

// TestPlayerHandler_ListPlayerLogs_WithActionFilter tests filtering logs by action.
func TestPlayerHandler_ListPlayerLogs_WithActionFilter(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test player
	user := integration.CreateUniqueTestUser(t, helper.db, "action_filter_player")
	player := integration.CreateTestPlayer(t, helper.db, user)

	w := helper.MakeRequest("GET", "/admin/players/"+uint64ToString(player.ID)+"/logs?action=update&page=1&pageSize=10", nil, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

// TestPlayerHandler_BatchUpdatePlayerStatus_Success tests batch updating player status.
func TestPlayerHandler_BatchUpdatePlayerStatus_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test players
	var playerIDs []string
	for i := 0; i < 3; i++ {
		user := integration.CreateUniqueTestUser(t, helper.db, "batch_status_player")
		player := integration.CreateTestPlayer(t, helper.db, user)
		player.VerificationStatus = model.VerificationPending
		helper.db.Save(player)
		playerIDs = append(playerIDs, uint64ToString(player.ID))
	}

	payload := map[string]interface{}{
		"playerIds": playerIDs,
		"status":    "verified",
	}

	w := helper.MakeRequest("PUT", "/admin/players/batch/status", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["updated"])

	// Verify all players are verified
	for _, idStr := range playerIDs {
		var player model.Player
		id := parseStringToUint64(t, idStr)
		helper.db.First(&player, id)
		assert.Equal(t, model.VerificationVerified, player.VerificationStatus)
	}
}

// TestPlayerHandler_BatchUpdatePlayerStatus_InvalidID tests batch update with invalid ID.
func TestPlayerHandler_BatchUpdatePlayerStatus_InvalidID(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"invalid_id"},
		"status":    "verified",
	}

	w := helper.MakeRequest("PUT", "/admin/players/batch/status", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPlayerHandler_BatchUpdatePlayerStatus_InvalidStatus tests batch update with invalid status.
func TestPlayerHandler_BatchUpdatePlayerStatus_InvalidStatus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"123"},
		"status":    "invalid_status",
	}

	w := helper.MakeRequest("PUT", "/admin/players/batch/status", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPlayerHandler_BatchDeletePlayers_Success tests batch deleting players.
func TestPlayerHandler_BatchDeletePlayers_Success(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	// Create test players
	var playerIDs []string
	for i := 0; i < 3; i++ {
		user := integration.CreateUniqueTestUser(t, helper.db, "batch_delete_player")
		player := integration.CreateTestPlayer(t, helper.db, user)
		playerIDs = append(playerIDs, uint64ToString(player.ID))
	}

	payload := map[string]interface{}{
		"playerIds": playerIDs,
	}

	w := helper.MakeRequest("POST", "/admin/players/batch/delete", payload, true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["deleted"])

	// Verify all players are deleted
	for _, idStr := range playerIDs {
		var count int64
		id := parseStringToUint64(t, idStr)
		helper.db.Model(&model.Player{}).Where("id = ?", id).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

// TestPlayerHandler_BatchDeletePlayers_InvalidID tests batch delete with invalid ID.
func TestPlayerHandler_BatchDeletePlayers_InvalidID(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"invalid_id"},
	}

	w := helper.MakeRequest("POST", "/admin/players/batch/delete", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestPlayerHandler_BatchDeletePlayers_EmptyList tests batch delete with empty list.
func TestPlayerHandler_BatchDeletePlayers_EmptyList(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	helper := SetupAdminTest(t)
	helper.RegisterRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{},
	}

	w := helper.MakeRequest("POST", "/admin/players/batch/delete", payload, true)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Helper functions

// uint64ToString converts uint64 to string for URL parameters.
func uint64ToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// parseStringToUint64 parses string to uint64 for testing.
func parseStringToUint64(t *testing.T, s string) uint64 {
	t.Helper()
	id, err := strconv.ParseUint(s, 10, 64)
	require.NoError(t, err)
	return id
}

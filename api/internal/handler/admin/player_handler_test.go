// Package admin provides unit tests for player handlers.
package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// PlayerTestContext provides test context for player handler tests.
type PlayerTestContext struct {
	Router     *gin.Engine
	Handler    *PlayerHandler
	Service    *adminservice.AdminService
	DB         *gorm.DB
	AdminUser  *model.User
	AdminToken string
	TestUser   *model.User
	TestGame   *model.Game
}

// SetupPlayerTest initializes test environment for player handler tests.
func SetupPlayerTest(t *testing.T) *PlayerTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	orders := implementations.NewOrderRepository(db)
	payments := payment.NewPaymentRepository(db)
	roles := admin.NewRoleRepository(db)
	serviceItems := serviceitem.NewServiceItemRepository(db)
	permissions := admin.NewPermissionRepository(db)
	menus := admin.NewMenuRepository(db)
	statsRepo := stats.NewStatsRepository(db)
	gameCategories := gamecategory.NewGameCategoryRepository(db)
	c := cache.NewMemory()

	// Create admin service
	svc := adminservice.NewAdminService(
		games, users, players, orders, payments,
		roles, serviceItems, permissions, menus, statsRepo, nil, gameCategories, c,
	)

	// Setup router
	router := testutil.SetupGinTest(t)
	handler := NewPlayerHandler(svc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	// Create test data
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
	testGame := testutil.CreateTestGame(t, db)

	return &PlayerTestContext{
		Router:     router,
		Handler:    handler,
		Service:    svc,
		DB:         db,
		AdminUser:  adminUser,
		AdminToken: adminToken,
		TestUser:   testUser,
		TestGame:   testGame,
	}
}

// RegisterPlayerRoutes registers player routes for testing.
func (ctx *PlayerTestContext) RegisterPlayerRoutes() {
	group := ctx.Router.Group("/admin/players")
	{
		group.GET("", ctx.Handler.ListPlayers)
		group.POST("", ctx.Handler.CreatePlayer)
		group.GET("/:id", ctx.Handler.GetPlayer)
		group.PUT("/:id", ctx.Handler.UpdatePlayer)
		group.DELETE("/:id", ctx.Handler.DeletePlayer)
		group.PUT("/:id/verification", ctx.Handler.UpdatePlayerVerification)
		group.PUT("/:id/games", ctx.Handler.UpdatePlayerGames)
		group.PUT("/:id/skill-tags", ctx.Handler.UpdatePlayerSkillTags)
		group.PUT("/batch/status", ctx.Handler.BatchUpdatePlayerStatus)
		group.POST("/batch/delete", ctx.Handler.BatchDeletePlayers)
		group.GET("/:id/logs", ctx.Handler.ListPlayerLogs)
	}
}

// ============================================================================
// CreatePlayer Tests
// ============================================================================

func TestPlayerHandler_Unit_CreatePlayer_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"user_id":              ctx.TestUser.ID,
		"nickname":             "Test Player",
		"bio":                  "Test player bio",
		"rank":                 "bronze",
		"hourly_rate_cents":    5000,
		"main_game_id":         ctx.TestGame.ID,
		"verification_status":  "pending",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/players", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "Test Player", data["nickname"])
}

func TestPlayerHandler_Unit_CreatePlayer_ValidationError(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"nickname": "Test Player",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/players", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPlayerHandler_Unit_CreateUser_UserNotFound(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"user_id":              999999,
		"nickname":             "Test Player",
		"verification_status":  "pending",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/players", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// ListPlayers Tests
// ============================================================================

func TestPlayerHandler_Unit_ListPlayers_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	// Create test players
	for i := 0; i < 5; i++ {
		user := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
		testutil.CreateTestPlayer(t, ctx.DB, user.ID)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/players", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestPlayerHandler_Unit_ListPlayers_WithPagination(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	// Create test players
	for i := 0; i < 25; i++ {
		user := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
		testutil.CreateTestPlayer(t, ctx.DB, user.ID)
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/players?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestPlayerHandler_Unit_ListPlayers_WithKeywordFilter(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	// Create player with specific nickname
	player := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)
	ctx.DB.Model(&player).Update("nickname", "SuperPlayer")

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/players?keyword=Super", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPlayerHandler_Unit_ListPlayers_WithStatusFilter(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	// Create players with different statuses
	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player1 := testutil.CreateTestPlayer(t, ctx.DB, user1.ID)
	ctx.DB.Model(&player1).Update("verification_status", model.VerificationPending)

	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player2 := testutil.CreateTestPlayer(t, ctx.DB, user2.ID)
	ctx.DB.Model(&player2).Update("verification_status", model.VerificationVerified)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/players?status=pending", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

// ============================================================================
// GetPlayer Tests
// ============================================================================

func TestPlayerHandler_Unit_GetPlayer_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	path := fmt.Sprintf("/admin/players/%d", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(testPlayer.ID), data["id"])
	assert.Equal(t, testPlayer.Nickname, data["nickname"])
}

func TestPlayerHandler_Unit_GetPlayer_NotFound(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/players/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPlayerHandler_Unit_GetPlayer_InvalidID(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/players/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdatePlayer Tests
// ============================================================================

func TestPlayerHandler_Unit_UpdatePlayer_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"nickname":             "Updated Nickname",
		"bio":                  "Updated bio",
		"rank":                 "silver",
		"hourly_rate_cents":    6000,
		"main_game_id":         ctx.TestGame.ID,
		"verification_status":  "verified",
	}

	path := fmt.Sprintf("/admin/players/%d", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedPlayer model.Player
	ctx.DB.First(&updatedPlayer, testPlayer.ID)
	assert.Equal(t, "Updated Nickname", updatedPlayer.Nickname)
}

func TestPlayerHandler_Unit_UpdatePlayer_NotFound(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"nickname":             "Updated Nickname",
		"verification_status":  "verified",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/players/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// DeletePlayer Tests
// ============================================================================

func TestPlayerHandler_Unit_DeletePlayer_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	path := fmt.Sprintf("/admin/players/%d", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify deletion
	var count int64
	ctx.DB.Model(&model.Player{}).Where("id = ?", testPlayer.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestPlayerHandler_Unit_DeletePlayer_NotFound(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/players/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// UpdatePlayerVerification Tests
// ============================================================================

func TestPlayerHandler_Unit_UpdatePlayerVerification_Verify(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"verification_status": "verified",
		"remark":              "Player verified successfully",
	}

	path := fmt.Sprintf("/admin/players/%d/verification", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedPlayer model.Player
	ctx.DB.First(&updatedPlayer, testPlayer.ID)
	assert.Equal(t, model.VerificationVerified, updatedPlayer.VerificationStatus)
}

func TestPlayerHandler_Unit_UpdatePlayerVerification_Reject(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"verification_status": "rejected",
		"remark":              "Invalid information",
	}

	path := fmt.Sprintf("/admin/players/%d/verification", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedPlayer model.Player
	ctx.DB.First(&updatedPlayer, testPlayer.ID)
	assert.Equal(t, model.VerificationRejected, updatedPlayer.VerificationStatus)
}

func TestPlayerHandler_Unit_UpdatePlayerVerification_InvalidStatus(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"verification_status": "invalid_status",
	}

	path := fmt.Sprintf("/admin/players/%d/verification", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdatePlayerGames Tests
// ============================================================================

func TestPlayerHandler_Unit_UpdatePlayerGames_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)
	newGame := testutil.CreateTestGame(t, ctx.DB)

	payload := map[string]interface{}{
		"main_game_id": newGame.ID,
	}

	path := fmt.Sprintf("/admin/players/%d/games", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update
	var updatedPlayer model.Player
	ctx.DB.First(&updatedPlayer, testPlayer.ID)
	assert.Equal(t, newGame.ID, updatedPlayer.MainGameID)
}

// ============================================================================
// UpdatePlayerSkillTags Tests
// ============================================================================

func TestPlayerHandler_Unit_UpdatePlayerSkillTags_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"tags": []string{"friendly", "skilled", "patient"},
	}

	path := fmt.Sprintf("/admin/players/%d/skill-tags", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

func TestPlayerHandler_Unit_UpdatePlayerSkillTags_EmptyTags(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.TestUser.ID)

	payload := map[string]interface{}{
		"tags": []string{},
	}

	path := fmt.Sprintf("/admin/players/%d/skill-tags", testPlayer.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPlayerHandler_Unit_UpdatePlayerSkillTags_PlayerNotFound(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"tags": []string{"friendly", "skilled"},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/players/999999/skill-tags", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// BatchUpdatePlayerStatus Tests
// ============================================================================

func TestPlayerHandler_Unit_BatchUpdatePlayerStatus_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player1 := testutil.CreateTestPlayer(t, ctx.DB, user1.ID)

	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player2 := testutil.CreateTestPlayer(t, ctx.DB, user2.ID)

	payload := map[string]interface{}{
		"playerIds": []string{fmt.Sprintf("%d", player1.ID), fmt.Sprintf("%d", player2.ID)},
		"status":    "verified",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/players/batch/status", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["updated"])
}

func TestPlayerHandler_Unit_BatchUpdatePlayerStatus_InvalidPlayerID(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"invalid_id"},
		"status":    "verified",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/players/batch/status", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPlayerHandler_Unit_BatchUpdatePlayerStatus_InvalidStatus(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"1"},
		"status":    "invalid_status",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/players/batch/status", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// BatchDeletePlayers Tests
// ============================================================================

func TestPlayerHandler_Unit_BatchDeletePlayers_Success(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player1 := testutil.CreateTestPlayer(t, ctx.DB, user1.ID)

	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	player2 := testutil.CreateTestPlayer(t, ctx.DB, user2.ID)

	payload := map[string]interface{}{
		"playerIds": []string{fmt.Sprintf("%d", player1.ID), fmt.Sprintf("%d", player2.ID)},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/players/batch/delete", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["deleted"])
}

func TestPlayerHandler_Unit_BatchDeletePlayers_InvalidID(t *testing.T) {
	ctx := SetupPlayerTest(t)
	ctx.RegisterPlayerRoutes()

	payload := map[string]interface{}{
		"playerIds": []string{"invalid_id"},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/players/batch/delete", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

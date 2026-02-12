// Package public provides unit tests for public handlers.
package public

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	gamerepo "gamelink/internal/repository/game"
	gamecategoryrepo "gamelink/internal/repository/gamecategory"
	playerrepo "gamelink/internal/repository/player"
	playerservicerepo "gamelink/internal/repository/playerservice"
	serviceitemrepo "gamelink/internal/repository/serviceitem"
	userrepo "gamelink/internal/repository/user"
)

// ============================================================================
// Test Setup
// ============================================================================

// PublicTestContext provides test context for public handler tests.
type PublicTestContext struct {
	Router     *gin.Engine
	DB         *gorm.DB
	TestUser   *model.User
	TestPlayer *model.Player
	TestGame   *model.Game
	TestItem   *model.ServiceItem
}

// SetupPublicTest initializes test environment for public handler tests.
func SetupPublicTest(t *testing.T) *PublicTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create test data
	testUser := testutil.CreateAdminUser(t, db, model.RolePlayer)
	testGame := testutil.CreateTestGame(t, db)
	testItem := testutil.CreateTestServiceItem(t, db, testGame.ID, 5000, true)
	testPlayer := testutil.CreateTestPlayer(t, db, testUser.ID)

	// Update player to verified status
	db.Model(testPlayer).Update("verification_status", model.VerificationVerified)

	return &PublicTestContext{
		Router:     router,
		DB:         db,
		TestUser:   testUser,
		TestPlayer: testPlayer,
		TestGame:   testGame,
		TestItem:   testItem,
	}
}

// RegisterPublicRoutes registers public routes for testing.
func (ctx *PublicTestContext) RegisterPublicRoutes() {
	publicGroup := ctx.Router.Group("/public")

	userRepo := userrepo.NewUserRepository(ctx.DB)
	playerRepo := playerrepo.NewPlayerRepository(ctx.DB)
	playerServiceRepo := playerservicerepo.NewPlayerServiceRepository(ctx.DB)
	gameRepo := gamerepo.NewGameRepository(ctx.DB)
	gameCategoryRepo := gamecategoryrepo.NewGameCategoryRepository(ctx.DB)
	serviceItemRepo := serviceitemrepo.NewServiceItemRepository(ctx.DB)

	// Register player routes
	playerHandler := NewPlayerHandler(playerRepo, playerServiceRepo, userRepo)
	playerHandler.RegisterRoutes(publicGroup)

	// Register game routes
	gameHandler := NewGameHandler(gameRepo, gameCategoryRepo)
	gameHandler.RegisterRoutes(publicGroup)

	// Register service item routes
	serviceItemHandler := NewServiceItemHandler(serviceItemRepo)
	serviceItemHandler.RegisterRoutes(publicGroup)

	// Register search routes
	searchHandler := NewSearchHandler(playerRepo, gameRepo, userRepo)
	RegisterSearchRoutes(publicGroup, searchHandler)
}

// makeRequest helper for making requests
func (ctx *PublicTestContext) makeRequest(t *testing.T, method, path string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)
	return w
}

// ============================================================================
// Player List Tests
// ============================================================================

func TestPublicHandler_ListPlayers_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/players")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_ListPlayers_WithPagination(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/players?page=1&pageSize=10")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_GetPlayerDetail_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	path := "/public/players/" + testutil.UintToString(ctx.TestPlayer.ID)
	w := ctx.makeRequest(t, "GET", path)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_GetPlayerDetail_NotFound(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/players/999999")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================================================
// Game List Tests
// ============================================================================

func TestPublicHandler_ListGames_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/games")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_GetGameDetail_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	path := "/public/games/" + testutil.UintToString(ctx.TestGame.ID)
	w := ctx.makeRequest(t, "GET", path)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_GetGameDetail_NotFound(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/games/999999")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============================================================================
// Service Item Tests
// ============================================================================

func TestPublicHandler_ListServiceItems_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/service-items")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_GetServiceItemDetail_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	path := "/public/service-items/" + testutil.UintToString(ctx.TestItem.ID)
	w := ctx.makeRequest(t, "GET", path)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

// ============================================================================
// Search Tests
// ============================================================================

func TestPublicHandler_Search_Success(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/search?q=test")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "players")
	assert.Contains(t, data, "games")
}

func TestPublicHandler_Search_PlayerOnly(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/search?q=test&type=player")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_Search_GameOnly(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/search?q=test&type=game")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestPublicHandler_Search_MissingQuery(t *testing.T) {
	ctx := SetupPublicTest(t)
	ctx.RegisterPublicRoutes()

	w := ctx.makeRequest(t, "GET", "/public/search")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

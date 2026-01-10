// Package user provides unit tests for favorite handlers.
package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/testutil"
	"gamelink/internal/model"
	favoriterepo "gamelink/internal/repository/favorite"
	playerrepo "gamelink/internal/repository/player"
)

// ============================================================================
// Test Setup
// ============================================================================

// FavoriteTestContext provides test context for favorite handler tests.
type FavoriteTestContext struct {
	Router     *gin.Engine
	DB         *gorm.DB
	TestUser   *model.User
	TestPlayer *model.Player
	AuthToken  string
}

// SetupFavoriteTest initializes test environment for favorite handler tests.
func SetupFavoriteTest(t *testing.T) *FavoriteTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create test user
	testUser := testutil.CreateAdminUser(t, db, model.RoleUser)
	authToken := fmt.Sprintf("Bearer test-token-%d", testUser.ID)

	// Create player user and player
	playerUser := testutil.CreateAdminUser(t, db, model.RolePlayer)
	testPlayer := testutil.CreateTestPlayer(t, db, playerUser.ID)

	return &FavoriteTestContext{
		Router:     router,
		DB:         db,
		TestUser:   testUser,
		TestPlayer: testPlayer,
		AuthToken:  authToken,
	}
}

// RegisterFavoriteRoutes registers favorite routes for testing.
func (ctx *FavoriteTestContext) RegisterFavoriteRoutes() {
	favoriteRepo := favoriterepo.NewRepository(ctx.DB)
	playerRepo := playerrepo.NewPlayerRepository(ctx.DB)
	handler := NewFavoriteHandler(favoriteRepo, playerRepo)

	group := ctx.Router.Group("/user/favorites/players")
	{
		group.GET("", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			handler.listFavorites(c)
		})
		group.POST("/:id", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			handler.addFavorite(c)
		})
		group.DELETE("/:id", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			handler.removeFavorite(c)
		})
		group.GET("/:id/check", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			handler.checkFavorite(c)
		})
	}
}

// makeRequest helper for making authenticated requests
func (ctx *FavoriteTestContext) makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	req, err := http.NewRequest(method, path, reqBody)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.AuthToken)

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)
	return w
}

// ============================================================================
// Add Favorite Tests
// ============================================================================

func TestFavoriteHandler_AddFavorite_Success(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	w := ctx.makeRequest(t, "POST", path, nil)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestFavoriteHandler_AddFavorite_AlreadyExists(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	// Add favorite first
	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	ctx.makeRequest(t, "POST", path, nil)

	// Try to add again
	w := ctx.makeRequest(t, "POST", path, nil)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestFavoriteHandler_AddFavorite_PlayerNotFound(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "POST", "/user/favorites/players/999999", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFavoriteHandler_AddFavorite_InvalidID(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "POST", "/user/favorites/players/invalid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============================================================================
// List Favorites Tests
// ============================================================================

func TestFavoriteHandler_ListFavorites_Success(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	// Add a favorite first
	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	ctx.makeRequest(t, "POST", path, nil)

	// List favorites
	w := ctx.makeRequest(t, "GET", "/user/favorites/players", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestFavoriteHandler_ListFavorites_Empty(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "GET", "/user/favorites/players", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))
}

func TestFavoriteHandler_ListFavorites_WithPagination(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "GET", "/user/favorites/players?page=1&pageSize=10", nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============================================================================
// Remove Favorite Tests
// ============================================================================

func TestFavoriteHandler_RemoveFavorite_Success(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	// Add favorite first
	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	ctx.makeRequest(t, "POST", path, nil)

	// Remove favorite
	w := ctx.makeRequest(t, "DELETE", path, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFavoriteHandler_RemoveFavorite_NotFound(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	w := ctx.makeRequest(t, "DELETE", path, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestFavoriteHandler_RemoveFavorite_InvalidID(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "DELETE", "/user/favorites/players/invalid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============================================================================
// Check Favorite Tests
// ============================================================================

func TestFavoriteHandler_CheckFavorite_IsFavorite(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	// Add favorite first
	path := fmt.Sprintf("/user/favorites/players/%d", ctx.TestPlayer.ID)
	ctx.makeRequest(t, "POST", path, nil)

	// Check favorite
	checkPath := fmt.Sprintf("/user/favorites/players/%d/check", ctx.TestPlayer.ID)
	w := ctx.makeRequest(t, "GET", checkPath, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.True(t, data["isFavorite"].(bool))
}

func TestFavoriteHandler_CheckFavorite_NotFavorite(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	checkPath := fmt.Sprintf("/user/favorites/players/%d/check", ctx.TestPlayer.ID)
	w := ctx.makeRequest(t, "GET", checkPath, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.False(t, data["isFavorite"].(bool))
}

func TestFavoriteHandler_CheckFavorite_InvalidID(t *testing.T) {
	ctx := SetupFavoriteTest(t)
	ctx.RegisterFavoriteRoutes()

	w := ctx.makeRequest(t, "GET", "/user/favorites/players/invalid/check", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Package player provides unit tests for status handlers.
package player

import (
	"bytes"
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
	orderrepo "gamelink/internal/repository/implementations"
	playerrepo "gamelink/internal/repository/player"
	reviewrepo "gamelink/internal/repository/review"
	userrepo "gamelink/internal/repository/user"
	playerservice "gamelink/internal/service/player"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// StatusTestContext provides test context for status handler tests.
type StatusTestContext struct {
	Router     *gin.Engine
	DB         *gorm.DB
	Service    *playerservice.PlayerService
	TestUser   *model.User
	TestPlayer *model.Player
}

// SetupStatusTest initializes test environment for status handler tests.
func SetupStatusTest(t *testing.T) *StatusTestContext {
	t.Helper()

	db := testutil.SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create repositories
	playerRepo := playerrepo.NewPlayerRepository(db)
	userRepo := userrepo.NewUserRepository(db)
	gameRepo := gamerepo.NewGameRepository(db)
	orderRepo := orderrepo.NewOrderRepository(db)
	reviewRepo := reviewrepo.NewReviewRepository(db)
	playerTagRepo := userrepo.NewPlayerTagRepository(db)

	// Create memory cache
	memCache := cache.NewMemory()

	// Create player service
	playerSvc := playerservice.NewPlayerService(
		playerRepo, userRepo, gameRepo, orderRepo, reviewRepo, playerTagRepo, memCache,
	)

	// Create test user (player user)
	testUser := testutil.CreateAdminUser(t, db, model.RolePlayer)
	testPlayer := testutil.CreateTestPlayer(t, db, testUser.ID)

	return &StatusTestContext{
		Router:     router,
		DB:         db,
		Service:    playerSvc,
		TestUser:   testUser,
		TestPlayer: testPlayer,
	}
}

// RegisterStatusRoutes registers status routes for testing.
func (ctx *StatusTestContext) RegisterStatusRoutes() {
	group := ctx.Router.Group("/player/online-status")
	{
		group.GET("", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			getOnlineStatusHandler(c, ctx.Service)
		})
		group.PUT("", func(c *gin.Context) {
			c.Set("user_id", ctx.TestUser.ID)
			updateOnlineStatusHandler(c, ctx.Service)
		})
	}
}

// makeRequest helper for making requests
func (ctx *StatusTestContext) makeRequest(t *testing.T, method, path string, body interface{}) *httptest.ResponseRecorder {
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

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)
	return w
}

// ============================================================================
// Get Online Status Tests
// ============================================================================

func TestStatusHandler_GetOnlineStatus_Success(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	w := ctx.makeRequest(t, "GET", "/player/online-status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Contains(t, data, "online")
}

func TestStatusHandler_GetOnlineStatus_DefaultOffline(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	w := ctx.makeRequest(t, "GET", "/player/online-status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.False(t, data["online"].(bool))
}

// ============================================================================
// Update Online Status Tests
// ============================================================================

func TestStatusHandler_UpdateOnlineStatus_SetOnline(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	payload := map[string]interface{}{
		"online": true,
	}

	w := ctx.makeRequest(t, "PUT", "/player/online-status", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.True(t, data["online"].(bool))
}

func TestStatusHandler_UpdateOnlineStatus_SetOffline(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	// First set online
	ctx.makeRequest(t, "PUT", "/player/online-status", map[string]interface{}{"online": true})

	// Then set offline
	payload := map[string]interface{}{
		"online": false,
	}

	w := ctx.makeRequest(t, "PUT", "/player/online-status", payload)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.False(t, data["online"].(bool))
}

func TestStatusHandler_UpdateOnlineStatus_VerifyPersistence(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	// Set online
	ctx.makeRequest(t, "PUT", "/player/online-status", map[string]interface{}{"online": true})

	// Verify status persisted
	w := ctx.makeRequest(t, "GET", "/player/online-status", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.True(t, data["online"].(bool))
}

func TestStatusHandler_UpdateOnlineStatus_InvalidBody(t *testing.T) {
	ctx := SetupStatusTest(t)
	ctx.RegisterStatusRoutes()

	// Send invalid JSON
	reqBody := bytes.NewBufferString("{invalid json}")
	req, err := http.NewRequest("PUT", "/player/online-status", reqBody)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ============================================================================
// Non-Player User Tests
// ============================================================================

func TestStatusHandler_GetOnlineStatus_NonPlayerUser(t *testing.T) {
	ctx := SetupStatusTest(t)

	// Create a non-player user
	nonPlayerUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	// Register routes with non-player user
	group := ctx.Router.Group("/player/online-status")
	{
		group.GET("", func(c *gin.Context) {
			c.Set("user_id", nonPlayerUser.ID)
			getOnlineStatusHandler(c, ctx.Service)
		})
	}

	w := ctx.makeRequest(t, "GET", "/player/online-status", nil)

	// Should return 404 because user is not a player
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStatusHandler_UpdateOnlineStatus_NonPlayerUser(t *testing.T) {
	ctx := SetupStatusTest(t)

	// Create a non-player user
	nonPlayerUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	// Register routes with non-player user
	group := ctx.Router.Group("/player/online-status")
	{
		group.PUT("", func(c *gin.Context) {
			c.Set("user_id", nonPlayerUser.ID)
			updateOnlineStatusHandler(c, ctx.Service)
		})
	}

	payload := map[string]interface{}{
		"online": true,
	}

	w := ctx.makeRequest(t, "PUT", "/player/online-status", payload)

	// Should return 404 because user is not a player
	assert.Equal(t, http.StatusNotFound, w.Code)
}

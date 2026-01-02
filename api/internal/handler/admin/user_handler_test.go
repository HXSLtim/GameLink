// Package admin provides unit tests for user handlers.
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
	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// UserTestContext provides test context for user handler tests.
type UserTestContext struct {
	Router     *gin.Engine
	Handler    *UserHandler
	Service    *adminservice.AdminService
	DB         *gorm.DB
	AdminUser  *model.User
	AdminToken string
}

// SetupUserTest initializes test environment for user handler tests.
func SetupUserTest(t *testing.T) *UserTestContext {
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
	handler := NewUserHandler(svc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	return &UserTestContext{
		Router:     router,
		Handler:    handler,
		Service:    svc,
		DB:         db,
		AdminUser:  adminUser,
		AdminToken: adminToken,
	}
}

// RegisterUserRoutes registers user routes for testing.
func (ctx *UserTestContext) RegisterUserRoutes() {
	group := ctx.Router.Group("/admin/users")
	{
		group.GET("", ctx.Handler.ListUsers)
		group.GET("/stats", ctx.Handler.GetUserStats)
		group.POST("", ctx.Handler.CreateUser)
		group.POST("/with-player", ctx.Handler.CreateUserWithPlayer)
		group.POST("/batch-delete", ctx.Handler.BatchDeleteUsers)
		group.GET("/:id", ctx.Handler.GetUser)
		group.PUT("/:id", ctx.Handler.UpdateUser)
		group.DELETE("/:id", ctx.Handler.DeleteUser)
		group.PUT("/:id/status", ctx.Handler.UpdateUserStatus)
		group.PUT("/:id/role", ctx.Handler.UpdateUserRole)
		group.GET("/:id/orders", ctx.Handler.ListUserOrders)
		group.GET("/:id/logs", ctx.Handler.ListUserLogs)
		group.GET("/:id/login-history", ctx.Handler.ListUserLoginHistory)
	}
}

// ============================================================================
// CreateUser Tests
// ============================================================================

func TestUserHandler_Unit_CreateUser_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"phone":      "13800138001",
		"email":      "test@example.com",
		"password":   "password123",
		"name":       "Test User",
		"avatar_url": "https://example.com/avatar.jpg",
		"role":       "user",
		"status":     "active",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "Test User", data["name"])
	assert.Equal(t, "user", data["role"])
}

func TestUserHandler_Unit_CreateUser_ValidationError(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"name": "Test User",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserHandler_Unit_CreateUser_InvalidEmail(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"phone":    "13800138001",
		"email":    "invalid-email",
		"password": "password123",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
	testutil.AssertErrorMessage(t, w, "email")
}

func TestUserHandler_Unit_CreateUser_InvalidPhone(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"phone":    "12345",
		"password": "password123",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
	testutil.AssertErrorMessage(t, w, "phone")
}

func TestUserHandler_Unit_CreateUser_ShortPassword(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"phone":    "13800138001",
		"password": "12345",
		"name":     "Test User",
		"role":     "user",
		"status":   "active",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// ListUsers Tests
// ============================================================================

func TestUserHandler_Unit_ListUsers_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create test users
	for i := 0; i < 5; i++ {
		testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/users", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestUserHandler_Unit_ListUsers_WithPagination(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create test users
	for i := 0; i < 25; i++ {
		testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/users?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestUserHandler_Unit_ListUsers_WithRoleFilter(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create users with different roles
	testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testutil.CreateAdminUser(t, ctx.DB, model.RolePlayer)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/users?role=user", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	for _, item := range items {
		user := item.(map[string]interface{})
		assert.Equal(t, "user", user["role"])
	}
}

func TestUserHandler_Unit_ListUsers_WithStatusFilter(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create users with different statuses
	testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	inactiveUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	ctx.DB.Model(&inactiveUser).Update("status", model.UserStatusSuspended)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/users?status=active", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestUserHandler_Unit_ListUsers_WithKeywordFilter(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create users with specific names
	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	ctx.DB.Model(&user1).Update("name", "John Doe")

	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	ctx.DB.Model(&user2).Update("name", "Jane Smith")

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/users?keyword=John", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

// ============================================================================
// GetUser Tests
// ============================================================================

func TestUserHandler_Unit_GetUser_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	path := fmt.Sprintf("/admin/users/%d", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(testUser.ID), data["id"])
}

func TestUserHandler_Unit_GetUser_NotFound(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/users/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestUserHandler_Unit_GetUser_InvalidID(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/users/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdateUser Tests
// ============================================================================

func TestUserHandler_Unit_UpdateUser_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"name":       "Updated Name",
		"avatar_url": "https://example.com/new-avatar.jpg",
		"role":       "user",
		"status":     "active",
	}

	path := fmt.Sprintf("/admin/users/%d", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedUser model.User
	ctx.DB.First(&updatedUser, testUser.ID)
	assert.Equal(t, "Updated Name", updatedUser.Name)
}

func TestUserHandler_Unit_UpdateUser_InvalidEmail(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"email":  "invalid-email",
		"name":   "Updated Name",
		"role":   "user",
		"status": "active",
	}

	path := fmt.Sprintf("/admin/users/%d", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestUserHandler_Unit_UpdateUser_UpdatePassword(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	newPassword := "newPassword123"

	payload := map[string]interface{}{
		"name":     testUser.Name,
		"password": &newPassword,
		"role":     "user",
		"status":   "active",
	}

	path := fmt.Sprintf("/admin/users/%d", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

// ============================================================================
// DeleteUser Tests
// ============================================================================

func TestUserHandler_Unit_DeleteUser_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	path := fmt.Sprintf("/admin/users/%d", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify deletion
	var count int64
	ctx.DB.Model(&model.User{}).Where("id = ?", testUser.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestUserHandler_Unit_DeleteUser_NotFound(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/users/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// BatchDeleteUsers Tests
// ============================================================================

func TestUserHandler_Unit_BatchDeleteUsers_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	user2 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"ids": []uint64{user1.ID, user2.ID},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users/batch-delete", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["deleted"])
}

func TestUserHandler_Unit_BatchDeleteUsers_PartialFailure(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	user1 := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"ids": []uint64{user1.ID, 999999},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users/batch-delete", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["deleted"])
	assert.Equal(t, float64(1), data["failed"])
}

func TestUserHandler_Unit_BatchDeleteUsers_EmptyIDs(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"ids": []uint64{},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users/batch-delete", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdateUserStatus Tests
// ============================================================================

func TestUserHandler_Unit_UpdateUserStatus_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"status": "inactive",
	}

	path := fmt.Sprintf("/admin/users/%d/status", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify status update
	var updatedUser model.User
	ctx.DB.First(&updatedUser, testUser.ID)
	assert.Equal(t, model.UserStatusSuspended, updatedUser.Status)
}

func TestUserHandler_Unit_UpdateUserStatus_InvalidStatus(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"status": "invalid_status",
	}

	path := fmt.Sprintf("/admin/users/%d/status", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdateUserRole Tests
// ============================================================================

func TestUserHandler_Unit_UpdateUserRole_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"role": "player",
	}

	path := fmt.Sprintf("/admin/users/%d/role", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify role update
	var updatedUser model.User
	ctx.DB.First(&updatedUser, testUser.ID)
	assert.Equal(t, model.RolePlayer, updatedUser.Role)
}

func TestUserHandler_Unit_UpdateUserRole_InvalidRole(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)

	payload := map[string]interface{}{
		"role": "invalid_role",
	}

	path := fmt.Sprintf("/admin/users/%d/role", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// CreateUserWithPlayer Tests
// ============================================================================

func TestUserHandler_Unit_CreateUserWithPlayer_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	payload := map[string]interface{}{
		"phone":    "13800138002",
		"email":    "player@example.com",
		"password": "password123",
		"name":     "Player User",
		"role":     "player",
		"status":   "active",
		"player": map[string]interface{}{
			"nickname":            "Test Player",
			"bio":                 "Player bio",
			"hourly_rate_cents":   5000,
			"main_game_id":        1,
			"verification_status": "pending",
		},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users/with-player", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["user"])
	assert.NotNil(t, data["player"])
}

func TestUserHandler_Unit_CreateUserWithPlayer_ValidationError(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Missing player info
	payload := map[string]interface{}{
		"phone":    "13800138002",
		"password": "password123",
		"name":     "Player User",
		"role":     "player",
		"status":   "active",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/users/with-player", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// GetUserStats Tests
// ============================================================================

func TestUserHandler_Unit_GetUserStats_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create test users
	testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testutil.CreateAdminUser(t, ctx.DB, model.RolePlayer)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/users/stats", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["total_users"])
	assert.NotNil(t, data["role_distribution"])
}

// ============================================================================
// ListUserOrders Tests
// ============================================================================

func TestUserHandler_Unit_ListUserOrders_Success(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	// Create test game
	testGame := testutil.CreateTestGame(t, ctx.DB)

	// Create test player
	testPlayer := testutil.CreateTestPlayer(t, ctx.DB, ctx.AdminUser.ID)

	// Create test user and order
	testUser := testutil.CreateAdminUser(t, ctx.DB, model.RoleUser)
	testutil.CreateTestOrder(t, ctx.DB, testUser.ID, testPlayer.ID, testGame.ID, model.OrderStatusPending)

	path := fmt.Sprintf("/admin/users/%d/orders", testUser.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestUserHandler_Unit_ListUserOrders_UserNotFound(t *testing.T) {
	ctx := SetupUserTest(t)
	ctx.RegisterUserRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/users/999999/orders", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

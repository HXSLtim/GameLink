// Package admin provides integration tests for admin handlers.
package admin

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

	"gamelink/internal/model"
	adminservice "gamelink/internal/service/admin"
	"gamelink/internal/service/integration"
	"gamelink/pkg/cache"

	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/common"
	"gamelink/internal/repository/game"
	"gamelink/internal/repository/gamecategory"
	"gamelink/internal/repository/implementations"
	"gamelink/internal/repository/payment"
	"gamelink/internal/repository/player"
	"gamelink/internal/repository/serviceitem"
	"gamelink/internal/repository/stats"
	"gamelink/internal/repository/user"
)

// AdminTestHelper provides utilities for admin handler testing.
type AdminTestHelper struct {
	router      *gin.Engine
	svc         *adminservice.AdminService
	db          *gorm.DB
	adminToken  string
	adminUserID uint64
}

// SetupAdminTest initializes the test environment for admin handlers.
func SetupAdminTest(t *testing.T) *AdminTestHelper {
	t.Helper()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test database
	db := integration.SetupTestDB(t)

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
	wallets := user.NewWalletRepository(db)
	c := cache.NewMemory()

	// Create Unit of Work for transaction support
	uow := common.NewUnitOfWork(db)

	// Create admin service
	svc := adminservice.NewAdminService(adminservice.AdminDeps{
		Games: games, Users: users, Players: players, Orders: orders, Payments: payments,
		Roles: roles, ServiceItems: serviceItems, Permissions: permissions, Menus: menus,
		Stats: statsRepo, Wallets: wallets, GameCategories: gameCategories, Cache: c,
	})

	// Set TxManager for operation logs support
	svc.SetTxManager(uow)

	// Create test router
	router := gin.New()

	// Create super admin user for authentication
	adminUser := &model.User{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Phone:  "13800138000",
		Name:   "Test Admin",
		Role:   model.RoleAdmin,
		Status: model.UserStatusActive,
	}
	require.NoError(t, db.Create(adminUser).Error)
	// Update role to superAdmin for testing
	require.NoError(t, db.Model(adminUser).Update("role", string(model.RoleSlugSuperAdmin)).Error)

	// Generate test token (simplified - in real scenario use JWT)
	adminToken := fmt.Sprintf("Bearer admin-test-token-%d", adminUser.ID)

	return &AdminTestHelper{
		router:      router,
		svc:         svc,
		db:          db,
		adminToken:  adminToken,
		adminUserID: adminUser.ID,
	}
}

// RegisterRoutes registers admin routes for testing.
func (h *AdminTestHelper) RegisterRoutes() {
	// Register routes directly without permission middleware for testing
	userHandler := NewUserHandler(h.svc)
	gameHandler := NewGameHandler(h.svc)

	// User routes
	h.router.GET("/admin/users", userHandler.ListUsers)
	h.router.GET("/admin/users/stats", userHandler.GetUserStats)
	h.router.POST("/admin/users", userHandler.CreateUser)
	h.router.POST("/admin/users/with-player", userHandler.CreateUserWithPlayer)
	h.router.POST("/admin/users/batch-delete", userHandler.BatchDeleteUsers)
	h.router.GET("/admin/users/:id", userHandler.GetUser)
	h.router.PUT("/admin/users/:id", userHandler.UpdateUser)
	h.router.DELETE("/admin/users/:id", userHandler.DeleteUser)
	h.router.PUT("/admin/users/:id/status", userHandler.UpdateUserStatus)
	h.router.PUT("/admin/users/:id/role", userHandler.UpdateUserRole)
	h.router.GET("/admin/users/:id/orders", userHandler.ListUserOrders)
	h.router.GET("/admin/users/:id/logs", userHandler.ListUserLogs)
	h.router.GET("/admin/users/:id/login-history", userHandler.ListUserLoginHistory)

	// Game routes
	h.router.GET("/admin/games", gameHandler.ListGames)
	h.router.POST("/admin/games", gameHandler.CreateGame)
	h.router.GET("/admin/games/:id", gameHandler.GetGame)
	h.router.PUT("/admin/games/:id", gameHandler.UpdateGame)
	h.router.DELETE("/admin/games/:id", gameHandler.DeleteGame)
	h.router.POST("/admin/games/batch/delete", gameHandler.BatchDeleteGames)
	h.router.POST("/admin/games/batch/status", gameHandler.BatchUpdateGamesStatus)
	h.router.POST("/admin/games/batch/category", gameHandler.BatchUpdateGamesCategory)
	h.router.GET("/admin/games/:id/logs", gameHandler.ListGameLogs)
}

// MakeRequest performs an HTTP request with authentication.
func (h *AdminTestHelper) MakeRequest(method, path string, body interface{}, auth bool) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal request body: %v", err))
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if auth {
		req.Header.Set("Authorization", h.adminToken)
	}

	w := httptest.NewRecorder()
	h.router.ServeHTTP(w, req)
	return w
}

// AssertJSONResponse asserts the response has the expected status and contains expected data.
func (h *AdminTestHelper) AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedData map[string]interface{}) {
	t.Helper()

	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	if expectedData != nil {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check response structure
		assert.True(t, response["success"].(bool))

		// Check expected fields
		for key, expectedValue := range expectedData {
			actualValue, exists := response[key]
			assert.True(t, exists, "Key %s should exist in response", key)
			assert.Equal(t, expectedValue, actualValue, "Key %s should match", key)
		}
	}
}

// mockPermissionMiddleware is a mock for testing without real permission checks.
type mockPermissionMiddleware struct{}

func (m *mockPermissionMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth in tests
		c.Next()
	}
}

func (m *mockPermissionMiddleware) RequirePermission(method string, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip permission checks in tests
		c.Next()
	}
}

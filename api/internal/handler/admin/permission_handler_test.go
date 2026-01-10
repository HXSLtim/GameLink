// Package admin provides unit tests for permission handlers.
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
	"gamelink/internal/service/integration"
	"gamelink/pkg/cache"
)

// ============================================================================
// Test Setup
// ============================================================================

// PermissionTestContext provides test context for permission handler tests.
type PermissionTestContext struct {
	Router       *gin.Engine
	Handler      *PermissionHandler
	Service      *adminservice.AdminService
	PermissionSvc *adminservice.PermissionService
	RoleSvc      *adminservice.RoleService
	DB           *gorm.DB
	AdminUser    *model.User
	AdminToken   string
}

// SetupPermissionTest initializes test environment for permission handler tests.
func SetupPermissionTest(t *testing.T) *PermissionTestContext {
	t.Helper()

	// Setup test database
	db := testutil.SetupTestDB(t)

	// Create repositories
	games := game.NewGameRepository(db)
	users := user.NewUserRepository(db)
	players := player.NewPlayerRepository(db)
	ordersRepo := implementations.NewOrderRepository(db)
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
		games, users, players, ordersRepo, payments,
		roles, serviceItems, permissions, menus, statsRepo, nil, gameCategories, c,
	)

	// Create permission and role services
	permissionSvc := adminservice.NewPermissionService(permissions, c)
	roleSvc := adminservice.NewRoleService(roles, c)

	// Setup router
	router := testutil.SetupGinTest(t)
	handler := NewPermissionHandlerWithRoleService(permissionSvc, roleSvc)

	// Create super admin user
	adminUser := testutil.CreateSuperAdmin(t, db)
	adminToken := testutil.GenerateTestToken(adminUser.ID)

	return &PermissionTestContext{
		Router:        router,
		Handler:       handler,
		Service:       svc,
		PermissionSvc: permissionSvc,
		RoleSvc:       roleSvc,
		DB:            db,
		AdminUser:     adminUser,
		AdminToken:    adminToken,
	}
}

// RegisterPermissionRoutes registers permission routes for testing.
func (ctx *PermissionTestContext) RegisterPermissionRoutes() {
	group := ctx.Router.Group("/admin/permissions")
	{
		group.GET("", ctx.Handler.ListPermissions)
		group.POST("", ctx.Handler.CreatePermission)
		group.GET("/groups", ctx.Handler.GetPermissionGroups)
		group.GET("/tree", ctx.Handler.GetPermissionTree)
		group.GET("/tree/grouped", ctx.Handler.GetPermissionTreeByGroup)
		group.GET("/me", ctx.Handler.GetCurrentUserPermissions)
		group.POST("/batch/delete", ctx.Handler.BatchDeletePermissions)
		group.DELETE("/batch", ctx.Handler.BatchDelete)
		group.GET("/:id", ctx.Handler.GetPermission)
		group.PUT("/:id", ctx.Handler.UpdatePermission)
		group.PATCH("/:id", ctx.Handler.PatchPermission)
		group.DELETE("/:id", ctx.Handler.DeletePermission)
	}

	// Role permissions routes
	ctx.Router.GET("/admin/roles/:id/permissions", ctx.Handler.GetRolePermissions)
	ctx.Router.GET("/admin/users/:id/permissions", ctx.Handler.GetUserPermissions)
}

// ============================================================================
// Test Data Helpers
// ============================================================================

// CreateTestPermission creates a test permission.
func (ctx *PermissionTestContext) CreateTestPermission(t *testing.T, method model.HTTPMethod, path, code string) *model.Permission {
	t.Helper()

	permission := &model.Permission{
		Base:        model.Base{ExtJSON: "{}"},
		Method:      method,
		Path:        path,
		Code:        code,
		Group:       "test",
		Description: "Test permission",
		SortOrder:   0,
		IsSystem:    false,
	}

	require.NoError(t, ctx.DB.Create(permission).Error, "Failed to create test permission")
	return permission
}

// ============================================================================
// ListPermissions Tests
// ============================================================================

func TestPermissionHandler_ListPermissions_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Create test permissions
	for i := 0; i < 5; i++ {
		ctx.CreateTestPermission(t, model.HTTPMethodGET, fmt.Sprintf("/api/test%d", i), fmt.Sprintf("test.read%d", i))
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 5)
	assert.Equal(t, float64(1), pagination["page"])
}

func TestPermissionHandler_ListPermissions_WithPagination(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Create test permissions
	for i := 0; i < 25; i++ {
		ctx.CreateTestPermission(t, model.HTTPMethodGET, fmt.Sprintf("/api/test%d", i), fmt.Sprintf("test.read%d", i))
	}

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/permissions?page=1&page_size=10", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, pagination := testutil.GetResponseList(t, w)
	assert.LessOrEqual(t, len(items), 10)
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(10), pagination["page_size"])
}

func TestPermissionHandler_ListPermissions_WithKeywordFilter(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")
	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/orders", "orders.read")

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/permissions?keyword=users", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPermissionHandler_ListPermissions_WithMethodFilter(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	ctx.CreateTestPermission(t, model.HTTPMethodPOST, "/api/test", "test.create")

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/permissions?method=GET", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPermissionHandler_ListPermissions_WithGroupFilter(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/permissions?group=test", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

func TestPermissionHandler_ListPermissions_WithIsSystemFilter(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	ctx.DB.Model(&permission).Update("is_system", true)

	w := testutil.MakeRequest(t, ctx.Router, "GET", "/admin/permissions?is_system=true", nil, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)

	items, _ := testutil.GetResponseList(t, w)
	assert.GreaterOrEqual(t, len(items), 1)
}

// ============================================================================
// GetPermission Tests
// ============================================================================

func TestPermissionHandler_GetPermission_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(permission.ID), data["id"])
	assert.Equal(t, "users.read", data["code"])
}

func TestPermissionHandler_GetPermission_NotFound(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPermissionHandler_GetPermission_InvalidID(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/invalid", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// CreatePermission Tests
// ============================================================================

func TestPermissionHandler_CreatePermission_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	payload := map[string]interface{}{
		"method":      "GET",
		"path":        "/api/test/new",
		"code":        "test.new.read",
		"group":       "test",
		"description": "Test permission description",
		"sortOrder":   10,
		"isSystem":    false,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w, http.StatusCreated)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "test.new.read", data["code"])
	assert.Equal(t, "GET", data["method"])
	assert.Equal(t, "/api/test/new", data["path"])
}

func TestPermissionHandler_CreatePermission_ValidationError(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Missing required fields
	payload := map[string]interface{}{
		"method": "GET",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPermissionHandler_CreatePermission_InvalidMethod(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	payload := map[string]interface{}{
		"method":      "INVALID",
		"path":        "/api/test",
		"code":        "test.read",
		"group":       "test",
		"description": "Test",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPermissionHandler_CreatePermission_PathTooLong(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	longPath := string(make([]byte, 256))
	for i := range longPath {
		longPath = longPath[:i] + "a"
	}

	payload := map[string]interface{}{
		"method":      "GET",
		"path":        longPath,
		"code":        "test.read",
		"group":       "test",
		"description": "Test",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// UpdatePermission Tests
// ============================================================================

func TestPermissionHandler_UpdatePermission_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	payload := map[string]interface{}{
		"group":       "updated_group",
		"description": "Updated description",
		"sortOrder":   20,
	}

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedPermission model.Permission
	ctx.DB.First(&updatedPermission, permission.ID)
	assert.Equal(t, "updated_group", updatedPermission.Group)
	assert.Equal(t, "Updated description", updatedPermission.Description)
	assert.Equal(t, 20, updatedPermission.SortOrder)
	// Method and code should remain unchanged
	assert.Equal(t, model.HTTPMethodGET, updatedPermission.Method)
	assert.Equal(t, "/api/test", updatedPermission.Path)
	assert.Equal(t, "test.read", updatedPermission.Code)
}

func TestPermissionHandler_UpdatePermission_NotFound(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	payload := map[string]interface{}{
		"group":       "updated_group",
		"description": "Updated description",
		"sortOrder":   20,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", "/admin/permissions/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPermissionHandler_UpdatePermission_InvalidGroupLength(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	longGroup := string(make([]byte, 65))
	for i := range longGroup {
		longGroup = longGroup[:i] + "a"
	}

	payload := map[string]interface{}{
		"group": longGroup,
	}

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PUT", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

// ============================================================================
// PatchPermission Tests
// ============================================================================

func TestPermissionHandler_PatchPermission_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	payload := map[string]interface{}{
		"group":       "patched_group",
		"description": "Patched description",
	}

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PATCH", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedPermission model.Permission
	ctx.DB.First(&updatedPermission, permission.ID)
	assert.Equal(t, "patched_group", updatedPermission.Group)
	assert.Equal(t, "Patched description", updatedPermission.Description)
}

func TestPermissionHandler_PatchPermission_OnlyCode(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	payload := map[string]interface{}{
		"code": "patched.code",
	}

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PATCH", path, ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	// Verify update in DB
	var updatedPermission model.Permission
	ctx.DB.First(&updatedPermission, permission.ID)
	assert.Equal(t, "patched.code", updatedPermission.Code)
}

func TestPermissionHandler_PatchPermission_NoFields(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	payload := map[string]interface{}{}

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PATCH", path, ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPermissionHandler_PatchPermission_NotFound(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	payload := map[string]interface{}{
		"description": "Patched description",
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "PATCH", "/admin/permissions/999999", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusNotFound)
}

// ============================================================================
// DeletePermission Tests
// ============================================================================

func TestPermissionHandler_DeletePermission_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	path := fmt.Sprintf("/admin/permissions/%d", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)

	// Verify soft deletion
	var count int64
	ctx.DB.Unscoped().Model(&model.Permission{}).Where("id = ?", permission.ID).Count(&count)
	assert.Equal(t, int64(1), count) // Record still exists (soft delete)

	var deletedPermission model.Permission
	err := ctx.DB.First(&deletedPermission, permission.ID).Error
	assert.Error(t, err) // Should not find active record
}

func TestPermissionHandler_DeletePermission_NotFound(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/permissions/999999", ctx.AdminToken, nil)
	testutil.AssertError(t, w, http.StatusNotFound)
}

func TestPermissionHandler_DeletePermission_ForceDelete(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	// Create a role and assign permission
	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")
	integration.AssignPermissionToRole(t, ctx.DB, role.ID, permission.ID)

	path := fmt.Sprintf("/admin/permissions/%d?force=true", permission.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", path, ctx.AdminToken, nil)
	testutil.AssertDeleted(t, w)
}

// ============================================================================
// GetPermissionGroups Tests
// ============================================================================

func TestPermissionHandler_GetPermissionGroups_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Create permissions with different groups
	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")
	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/orders", "orders.read")
	ctx.CreateTestPermission(t, model.HTTPMethodPOST, "/api/users", "users.create")

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/groups", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// GetPermissionTree Tests
// ============================================================================

func TestPermissionHandler_GetPermissionTree_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Create parent permission
	parent := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")

	// Create child permission
	child := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users/:id", "users.detail")
	ctx.DB.Model(&child).Update("parent_id", parent.ID)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/tree", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

func TestPermissionHandler_GetPermissionTree_Empty(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/tree", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.NotNil(t, data)
}

// ============================================================================
// GetPermissionTreeByGroup Tests
// ============================================================================

func TestPermissionHandler_GetPermissionTreeByGroup_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/users", "users.read")
	ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/orders", "orders.read")

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/tree/grouped", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// GetCurrentUserPermissions Tests
// ============================================================================

func TestPermissionHandler_GetCurrentUserPermissions_SuperAdmin(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/me", ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, []interface{}{"*"}, data)
}

func TestPermissionHandler_GetCurrentUserPermissions_RegularUser(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	// Create regular user
	regularUser := integration.CreateTestUser(t, ctx.DB, "regularuser")
	regularToken := testutil.GenerateTestToken(regularUser.ID)

	// Create a role and assign permissions
	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")
	integration.AssignPermissionToRole(t, ctx.DB, role.ID, permission.ID)
	integration.AssignRoleToUser(t, ctx.DB, regularUser.ID, role.ID)

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", "/admin/permissions/me", regularToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// GetRolePermissions Tests
// ============================================================================

func TestPermissionHandler_GetRolePermissions_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")
	integration.AssignPermissionToRole(t, ctx.DB, role.ID, permission.ID)

	path := fmt.Sprintf("/admin/roles/%d/permissions", role.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

func TestPermissionHandler_GetRolePermissions_Empty(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")

	path := fmt.Sprintf("/admin/roles/%d/permissions", role.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.Equal(t, 0, len(data))
}

// ============================================================================
// GetUserPermissions Tests
// ============================================================================

func TestPermissionHandler_GetUserPermissions_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permission := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")
	user := integration.CreateTestUser(t, ctx.DB, "testuser")

	integration.AssignPermissionToRole(t, ctx.DB, role.ID, permission.ID)
	integration.AssignRoleToUser(t, ctx.DB, user.ID, role.ID)

	path := fmt.Sprintf("/admin/users/%d/permissions", user.ID)
	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "GET", path, ctx.AdminToken, nil)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
}

// ============================================================================
// BatchDeletePermissions Tests
// ============================================================================

func TestPermissionHandler_BatchDeletePermissions_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	var permissionIDs []uint64
	for i := 0; i < 3; i++ {
		perm := ctx.CreateTestPermission(t, model.HTTPMethodGET, fmt.Sprintf("/api/test%d", i), fmt.Sprintf("test.read%d", i))
		permissionIDs = append(permissionIDs, perm.ID)
	}

	payload := map[string]interface{}{
		"permission_ids": permissionIDs,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions/batch/delete", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(len(permissionIDs)), data["deleted"])
	assert.Equal(t, float64(0), data["failed"])
}

func TestPermissionHandler_BatchDeletePermissions_EmptyList(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	payload := map[string]interface{}{
		"permission_ids": []uint64{},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions/batch/delete", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPermissionHandler_BatchDeletePermissions_ExceedsLimit(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	permissionIDs := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		permissionIDs[i] = uint64(i + 1)
	}

	payload := map[string]interface{}{
		"permission_ids": permissionIDs,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions/batch/delete", ctx.AdminToken, payload)
	testutil.AssertError(t, w, http.StatusBadRequest)
}

func TestPermissionHandler_BatchDeletePermissions_PartialFailure(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	perm := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")

	payload := map[string]interface{}{
		"permission_ids": []uint64{perm.ID, 999999},
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "POST", "/admin/permissions/batch/delete", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["deleted"])
	assert.Equal(t, float64(1), data["failed"])
}

// ============================================================================
// BatchDelete (Legacy) Tests
// ============================================================================

func TestPermissionHandler_BatchDelete_Success(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	var permissionIDs []uint64
	for i := 0; i < 3; i++ {
		perm := ctx.CreateTestPermission(t, model.HTTPMethodGET, fmt.Sprintf("/api/test%d", i), fmt.Sprintf("test.read%d", i))
		permissionIDs = append(permissionIDs, perm.ID)
	}

	payload := map[string]interface{}{
		"permission_ids": permissionIDs,
	}

	w := testutil.MakeAuthenticatedRequest(t, ctx.Router, "DELETE", "/admin/permissions/batch", ctx.AdminToken, payload)
	testutil.AssertSuccess(t, w)
}

func TestPermissionHandler_BatchDelete_WithForce(t *testing.T) {
	ctx := SetupPermissionTest(t)
	ctx.RegisterPermissionRoutes()

	perm := ctx.CreateTestPermission(t, model.HTTPMethodGET, "/api/test", "test.read")
	role := integration.CreateTestRole(t, ctx.DB, "test-role", "Test Role")
	integration.AssignPermissionToRole(t, ctx.DB, role.ID, perm.ID)

	payload := map[string]interface{}{
		"permission_ids": []uint64{perm.ID},
	}

	w := testutil.MakeRequest(t, ctx.Router, "DELETE", "/admin/permissions/batch?force=true", payload, testutil.WithAuth(ctx.AdminToken))
	testutil.AssertSuccess(t, w)
}

package admin

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

	"gamelink/internal/model"
	roleperm "gamelink/internal/repository/admin"
	permadmin "gamelink/internal/service/admin"
	"gamelink/internal/service/integration"
	"gamelink/pkg/cache"
)

// setupRoleTest creates test environment for role handler tests
func setupRoleTest(t *testing.T) (*gin.Engine, *RoleHandler, *gorm.DB) {
	t.Helper()
	db := integration.SetupTestDB(t)

	// Create repositories
	roleRepo := roleperm.NewRoleRepository(db)

	// Create cache
	cacheClient := cache.NewMemory()

	// Create services
	roleSvc := permadmin.NewRoleService(roleRepo, cacheClient)

	// Create handler
	handler := NewRoleHandler(roleSvc)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register routes
	router.GET("/admin/roles", handler.ListRoles)
	router.GET("/admin/roles/:id", handler.GetRole)
	router.POST("/admin/roles", handler.CreateRole)
	router.PUT("/admin/roles/:id", handler.UpdateRole)
	router.DELETE("/admin/roles/:id", handler.DeleteRole)
	router.GET("/admin/roles/:id/permissions", handler.GetRolePermissionIDs)
	router.PUT("/admin/roles/:id/permissions/batch", handler.AssignPermissions)
	router.POST("/admin/roles/:id/permissions/:pid", handler.AddPermissionToRole)
	router.DELETE("/admin/roles/:id/permissions/:pid", handler.RemovePermissionFromRole)
	router.POST("/admin/users/roles", handler.AssignRolesToUser)
	router.GET("/admin/users/:id/roles", handler.GetUserRoles)
	router.PUT("/admin/users/:id/roles", handler.UpdateUserRoles)
	router.PUT("/admin/users/roles/batch", handler.BatchAssignRolesToUsers)
	router.DELETE("/admin/roles/batch", handler.BatchDeleteRoles)
	router.PUT("/admin/roles/batch/permissions", handler.BatchAssignPermissionsToRoles)

	return router, handler, db
}

// TestRoleHandler_ListRoles tests listing roles
func TestRoleHandler_ListRoles(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test roles
	roles := []model.RoleModel{
		{Slug: "admin", Name: "Admin", Description: "Administrator", IsSystem: true},
		{Slug: "user", Name: "User", Description: "Regular user", IsSystem: false},
		{Slug: "player", Name: "Player", Description: "Game player", IsSystem: false},
	}
	for _, role := range roles {
		require.NoError(t, db.Create(&role).Error)
	}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:           "list all roles",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"data"`)
			},
		},
		{
			name:           "list with pagination",
			queryParams:    "?page=1&pageSize=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"page":1`)
				assert.Contains(t, body, `"page_size":2`)
				assert.Contains(t, body, `"total":3`)
			},
		},
		{
			name:           "filter by keyword",
			queryParams:    "?keyword=admin",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "admin")
			},
		},
		{
			name:           "filter system roles",
			queryParams:    "?isSystem=true",
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/admin/roles"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}

			items := response["items"].([]interface{})
			assert.Equal(t, tt.expectedCount, len(items))
		})
	}
}

// TestRoleHandler_GetRole tests getting a role by ID
func TestRoleHandler_GetRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	tests := []struct {
		name           string
		roleID         string
		queryParams    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:           "get role without permissions",
			roleID:         "1",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"slug":"test-role"`)
				assert.Contains(t, body, `"name":"Test Role"`)
			},
		},
		{
			name:           "role not found",
			roleID:         "999",
			queryParams:    "",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "not found")
			},
		},
		{
			name:           "invalid role ID",
			roleID:         "invalid",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/admin/roles/"+tt.roleID+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestRoleHandler_CreateRole tests creating a new role
func TestRoleHandler_CreateRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, _ := setupRoleTest(t)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name: "create role successfully",
			requestBody: `{
				"slug": "moderator",
				"name": "Moderator",
				"description": "Content moderator"
			}`,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"slug":"moderator"`)
				assert.Contains(t, body, `"name":"Moderator"`)
			},
		},
		{
			name: "validation error - empty slug",
			requestBody: `{
				"slug": "",
				"name": "Test"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty name",
			requestBody: `{
				"slug": "test",
				"name": ""
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{invalid json}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/admin/roles", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestRoleHandler_UpdateRole tests updating a role
func TestRoleHandler_UpdateRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	tests := []struct {
		name           string
		roleID         string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:   "update role successfully",
			roleID: "1",
			requestBody: `{
				"name": "Updated Role",
				"description": "Updated description"
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"name":"Updated Role"`)
				assert.Contains(t, body, `"description":"Updated description"`)
			},
		},
		{
			name:   "role not found",
			roleID: "999",
			requestBody: `{
				"name": "Updated"
			}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid JSON",
			roleID:         "1",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/admin/roles/"+tt.roleID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestRoleHandler_DeleteRole tests deleting a role
func TestRoleHandler_DeleteRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	tests := []struct {
		name           string
		roleID         string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:           "delete role successfully",
			roleID:         "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "deleted")
			},
		},
		{
			name:           "role not found",
			roleID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid role ID",
			roleID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/admin/roles/"+tt.roleID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestRoleHandler_AssignPermissions tests assigning permissions to a role
func TestRoleHandler_AssignPermissions(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	// Create test permissions
	perms := []model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/admin/test1", Code: "admin.test1", Group: "admin"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/admin/test2", Code: "admin.test2", Group: "admin"},
	}
	for _, perm := range perms {
		require.NoError(t, db.Create(&perm).Error)
	}

	tests := []struct {
		name           string
		roleID         string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:   "assign permissions successfully",
			roleID: "1",
			requestBody: `{
				"permissionIds": [1, 2]
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "success")
			},
		},
		{
			name:   "role not found",
			roleID: "999",
			requestBody: `{
				"permissionIds": [1]
			}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid JSON",
			roleID:         "1",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/admin/roles/"+tt.roleID+"/permissions/batch", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestRoleHandler_AddPermissionToRole tests adding a single permission to a role
func TestRoleHandler_AddPermissionToRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	// Create test permission
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/test",
		Code:        "admin.test",
		Group:       "admin",
		Description: "Test permission",
	}
	require.NoError(t, db.Create(perm).Error)

	tests := []struct {
		name           string
		roleID         string
		permID         string
		expectedStatus int
	}{
		{
			name:           "add permission successfully",
			roleID:         "1",
			permID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "role not found",
			roleID:         "999",
			permID:         "1",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid permission ID",
			roleID:         "1",
			permID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/admin/roles/"+tt.roleID+"/permissions/"+tt.permID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestRoleHandler_RemovePermissionFromRole tests removing a permission from a role
func TestRoleHandler_RemovePermissionFromRole(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test role
	role := &model.RoleModel{
		Slug:        "test-role",
		Name:        "Test Role",
		Description: "Test description",
		IsSystem:    false,
	}
	require.NoError(t, db.Create(role).Error)

	// Create test permission
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/test",
		Code:        "admin.test",
		Group:       "admin",
		Description: "Test permission",
	}
	require.NoError(t, db.Create(perm).Error)

	// Assign permission to role
	require.NoError(t, db.Create(&model.RolePermission{
		RoleID:       role.ID,
		PermissionID: perm.ID,
	}).Error)

	tests := []struct {
		name           string
		roleID         string
		permID         string
		expectedStatus int
	}{
		{
			name:           "remove permission successfully",
			roleID:         "1",
			permID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "role not found",
			roleID:         "999",
			permID:         "1",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/admin/roles/"+tt.roleID+"/permissions/"+tt.permID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestRoleHandler_BatchDeleteRoles tests batch deleting roles
func TestRoleHandler_BatchDeleteRoles(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupRoleTest(t)

	// Create test roles
	roles := []model.RoleModel{
		{Slug: "role1", Name: "Role 1", Description: "Test 1", IsSystem: false},
		{Slug: "role2", Name: "Role 2", Description: "Test 2", IsSystem: false},
	}
	for _, role := range roles {
		require.NoError(t, db.Create(&role).Error)
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name: "batch delete successfully",
			requestBody: `{
				"ids": [1, 2]
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "successCount")
			},
		},
		{
			name:           "empty ids",
			requestBody:    `{"ids": []}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/admin/roles/batch", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

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

// setupMenuTest creates test environment for menu handler tests
func setupMenuTest(t *testing.T) (*gin.Engine, *MenuHandler, *gorm.DB) {
	t.Helper()
	db := integration.SetupTestDB(t)

	// Create repositories
	menuRepo := roleperm.NewMenuRepository(db)
	permRepo := roleperm.NewPermissionRepository(db)
	roleRepo := roleperm.NewRoleRepository(db)

	// Create cache
	cacheClient := cache.NewMemory()

	// Create services
	permSvc := permadmin.NewPermissionService(permRepo, cacheClient)
	roleSvc := permadmin.NewRoleService(roleRepo, cacheClient)
	menuSvc := permadmin.NewMenuService(menuRepo)

	// Create handler with role service
	handler := NewMenuHandlerWithRoleService(menuSvc, permSvc, roleSvc)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register routes
	router.GET("/admin/menus", handler.List)
	router.GET("/admin/menus/:id", handler.Get)
	router.POST("/admin/menus", handler.Create)
	router.PUT("/admin/menus/:id", handler.Update)
	router.DELETE("/admin/menus/:id", handler.Delete)
	router.GET("/admin/menus/me", handler.ListMyMenus)
	router.DELETE("/admin/menus/batch", handler.BatchDelete)
	router.PUT("/admin/menus/batch/status", handler.BatchUpdateStatus)
	router.PUT("/admin/menus/batch/sort", handler.BatchUpdateSort)

	return router, handler, db
}

// TestMenuHandler_ListMenus tests listing menus
func TestMenuHandler_ListMenus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menus
	parentMenu := &model.Menu{
		Name:   "Parent Menu",
		Path:   "/parent",
		Icon:   "parent",
		Order:  1,
		Hidden: false,
	}
	require.NoError(t, db.Create(parentMenu).Error)

	childMenu := &model.Menu{
		Name:     "Child Menu",
		Path:     "/child",
		Icon:     "child",
		Order:    2,
		Hidden:   false,
		ParentID: &parentMenu.ID,
	}
	require.NoError(t, db.Create(childMenu).Error)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount   int
	}{
		{
			name:           "list all menus",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:   2,
		},
		{
			name:           "list with pagination",
			queryParams:    "?page=1&pageSize=1",
			expectedStatus: http.StatusOK,
			expectedCount:   1,
		},
		{
			name:           "list by parent ID",
			queryParams:    "?parentId=1",
			expectedStatus: http.StatusOK,
			expectedCount:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/admin/menus"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Check structure
			if tt.queryParams == "" || tt.queryParams == "?parentId=1" {
				// Non-paginated response
				assert.Contains(t, response, "data")
			} else {
				// Paginated response
				assert.Contains(t, response, "items")
				items := response["items"].([]interface{})
				assert.Equal(t, tt.expectedCount, len(items))
			}
		})
	}
}

// TestMenuHandler_GetMenu tests getting a menu by ID
func TestMenuHandler_GetMenu(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menu
	menu := &model.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Icon:     "test",
		Order: 1,
		Hidden:   false,
	}
	require.NoError(t, db.Create(menu).Error)

	tests := []struct {
		name           string
		menuID         string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:           "get menu successfully",
			menuID:         "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"name":"Test Menu"`)
				assert.Contains(t, body, `"path":"/test"`)
			},
		},
		{
			name:           "menu not found",
			menuID:         "999",
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "not found")
			},
		},
		{
			name:           "invalid menu ID",
			menuID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/admin/menus/"+tt.menuID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestMenuHandler_CreateMenu tests creating a new menu
func TestMenuHandler_CreateMenu(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, _ := setupMenuTest(t)

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name: "create menu successfully",
			requestBody: `{
				"name": "New Menu",
				"path": "/new",
				"icon": "new-icon",
				"order": 1,
				"hidden": false
			}`,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"name":"New Menu"`)
				assert.Contains(t, body, `"path":"/new"`)
			},
		},
		{
			name: "create menu with parent",
			requestBody: `{
				"name": "Child Menu",
				"path": "/child",
				"parentId": 1
			}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "validation error - empty name",
			requestBody: `{
				"name": "",
				"path": "/test"
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
			req, _ := http.NewRequest("POST", "/admin/menus", bytes.NewBufferString(tt.requestBody))
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

// TestMenuHandler_UpdateMenu tests updating a menu
func TestMenuHandler_UpdateMenu(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menu
	menu := &model.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Icon:     "test",
		Order: 1,
		Hidden:   false,
	}
	require.NoError(t, db.Create(menu).Error)

	tests := []struct {
		name           string
		menuID         string
		requestBody    string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:   "update menu successfully",
			menuID: "1",
			requestBody: `{
				"name": "Updated Menu",
				"path": "/updated",
				"icon": "updated-icon",
				"order": 2
			}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, `"name":"Updated Menu"`)
				assert.Contains(t, body, `"path":"/updated"`)
			},
		},
		{
			name:   "menu not found",
			menuID: "999",
			requestBody: `{
				"name": "Updated"
			}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid JSON",
			menuID:         "1",
			requestBody:    `{invalid}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("PUT", "/admin/menus/"+tt.menuID, bytes.NewBufferString(tt.requestBody))
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

// TestMenuHandler_DeleteMenu tests deleting a menu
func TestMenuHandler_DeleteMenu(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menu
	menu := &model.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Icon:     "test",
		Order: 1,
		Hidden:   false,
	}
	require.NoError(t, db.Create(menu).Error)

	tests := []struct {
		name           string
		menuID         string
		expectedStatus int
		checkResponse  func(*testing.T, string)
	}{
		{
			name:           "delete menu successfully",
			menuID:         "1",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				assert.Contains(t, body, "deleted")
			},
		},
		{
			name:           "menu not found",
			menuID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid menu ID",
			menuID:         "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/admin/menus/"+tt.menuID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}
		})
	}
}

// TestMenuHandler_ListMyMenus tests listing menus for current user
func TestMenuHandler_ListMyMenus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	_, handler, db := setupMenuTest(t)

	// Create test menus
	menus := []model.Menu{
		{Name: "Menu 1", Path: "/menu1", Icon: "menu1", Order: 1, Hidden: false},
		{Name: "Menu 2", Path: "/menu2", Icon: "menu2", Order: 2, Hidden: false},
	}
	for _, menu := range menus {
		require.NoError(t, db.Create(&menu).Error)
	}

	// Create test user
	user := &model.User{
		Name:   "Test User",
		Phone:  "13800138000",
		Role:   model.RoleAdmin,
		Status: model.UserStatusActive,
	}
	require.NoError(t, db.Create(user).Error)

	// Create super admin role
	superAdminRole := &model.RoleModel{
		Slug:        string(model.RoleSlugSuperAdmin),
		Name:        "Super Admin",
		Description: "Super administrator",
		IsSystem:    true,
	}
	require.NoError(t, db.Create(superAdminRole).Error)

	// Assign super admin role to user
	require.NoError(t, db.Create(&model.UserRole{
		UserID: user.ID,
		RoleID: superAdminRole.ID,
	}).Error)

	// Create permissions for menus
	for _, menu := range menus {
		perm := &model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/admin/menus",
			Code:        "admin.menus." + menu.Path,
			Group:       "admin.menus",
			Description: "Access " + menu.Name,
		}
		require.NoError(t, db.Create(perm).Error)

		// Assign permission to role
		require.NoError(t, db.Create(&model.RolePermission{
			RoleID:       superAdminRole.ID,
			PermissionID: perm.ID,
		}).Error)
	}

	// Setup context with user ID
	gin.SetMode(gin.TestMode)
	testRouter := gin.New()
	testRouter.GET("/admin/menus/me", func(c *gin.Context) {
		c.Set("userID", uint64(user.ID))
		c.Set("userRole", string(model.RoleAdmin))
		handler.ListMyMenus(c)
	})

	req, _ := http.NewRequest("GET", "/admin/menus/me", nil)
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Menu 1")
	assert.Contains(t, w.Body.String(), "Menu 2")
}

// TestMenuHandler_BatchDelete tests batch deleting menus
func TestMenuHandler_BatchDelete(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menus
	menus := []model.Menu{
		{Name: "Menu 1", Path: "/menu1", Icon: "menu1", Order: 1, Hidden: false},
		{Name: "Menu 2", Path: "/menu2", Icon: "menu2", Order: 2, Hidden: false},
	}
	for _, menu := range menus {
		require.NoError(t, db.Create(&menu).Error)
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
			req, _ := http.NewRequest("DELETE", "/admin/menus/batch", bytes.NewBufferString(tt.requestBody))
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

// TestMenuHandler_BatchUpdateStatus tests batch updating menu status
func TestMenuHandler_BatchUpdateStatus(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menus
	menus := []model.Menu{
		{Name: "Menu 1", Path: "/menu1", Icon: "menu1", Order: 1, Hidden: false},
		{Name: "Menu 2", Path: "/menu2", Icon: "menu2", Order: 2, Hidden: false},
	}
	for _, menu := range menus {
		require.NoError(t, db.Create(&menu).Error)
	}

	requestBody := `{
		"updates": [
			{"id": 1, "hidden": true},
			{"id": 2, "hidden": true}
		]
	}`

	req, _ := http.NewRequest("PUT", "/admin/menus/batch/status", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "successCount")
}

// TestMenuHandler_BatchUpdateSort tests batch updating menu sort order
func TestMenuHandler_BatchUpdateSort(t *testing.T) {
	integration.SkipIfNoTestDB(t)
	router, _, db := setupMenuTest(t)

	// Create test menus
	menus := []model.Menu{
		{Name: "Menu 1", Path: "/menu1", Icon: "menu1", Order: 1, Hidden: false},
		{Name: "Menu 2", Path: "/menu2", Icon: "menu2", Order: 2, Hidden: false},
	}
	for _, menu := range menus {
		require.NoError(t, db.Create(&menu).Error)
	}

	requestBody := `{
		"updates": [
			{"id": 1, "sortOrder": 10},
			{"id": 2, "sortOrder": 20}
		]
	}`

	req, _ := http.NewRequest("PUT", "/admin/menus/batch/sort", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "successCount")
}

// TestBuildMenuTree tests the buildMenuTree helper function
func TestBuildMenuTree(t *testing.T) {
	tests := []struct {
		name     string
		menus    []model.Menu
		expected int
	}{
		{
			name: "flat list becomes tree",
			menus: []model.Menu{
				{Base: model.Base{ID: 1}, Name: "Parent", Path: "/parent"},
				{Base: model.Base{ID: 2}, Name: "Child", Path: "/child", ParentID: uint64Ptr(1)},
			},
			expected: 1, // One root node
		},
		{
			name: "multiple root nodes",
			menus: []model.Menu{
				{Base: model.Base{ID: 1}, Name: "Root1", Path: "/root1"},
				{Base: model.Base{ID: 2}, Name: "Root2", Path: "/root2"},
			},
			expected: 2, // Two root nodes
		},
		{
			name:     "empty list",
			menus:    []model.Menu{},
			expected: 0,
		},
		{
			name: "nested children",
			menus: []model.Menu{
				{Base: model.Base{ID: 1}, Name: "Parent", Path: "/parent"},
				{Base: model.Base{ID: 2}, Name: "Child", Path: "/child", ParentID: uint64Ptr(1)},
				{Base: model.Base{ID: 3}, Name: "Grandchild", Path: "/grandchild", ParentID: uint64Ptr(2)},
			},
			expected: 1, // One root node with nested children
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildMenuTree(tt.menus)
			assert.Equal(t, tt.expected, len(result), "should have correct number of root nodes")

			// Verify tree structure
			if len(tt.menus) > 0 && tt.menus[0].ParentID == nil {
				// First menu is a root, should be in result
				found := false
				for _, root := range result {
					if root.ID == tt.menus[0].ID {
						found = true
						break
					}
				}
				assert.True(t, found, "root menu should be in result")
			}
		})
	}
}

// Helper function
func uint64Ptr(v uint64) *uint64 {
	return &v
}

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/user"
	adminservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// TestUserPermissionQuery tests user permission query functionality
// Requirements: 4.1 - WHEN 用户请求受保护的 API THEN 系统 SHALL 验证用户是否拥有该 API 对应的权限
func TestUserPermissionQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionVerificationModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()

	// Create test users with isolated transaction
	testUser := &model.User{
		Name:         "TestUser",
		Email:        "testuser@example.com",
		Phone:        "18800000101",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), testUser))

	// Create a role with specific permissions
	testRole := &model.RoleModel{
		Slug:     "test-role",
		Name:     "Test Role",
		IsSystem: false,
	}
	require.NoError(t, roleRepo.Create(ctx(), testRole))

	// Create permissions
	readPerm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/test-resource",
		Code:        "test.resource.read",
		Group:       "test",
		Description: "Read test resource",
	}
	writePerm := &model.Permission{
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/admin/test-resource",
		Code:        "test.resource.write",
		Group:       "test",
		Description: "Write test resource",
	}
	require.NoError(t, permRepo.Create(ctx(), readPerm))
	require.NoError(t, permRepo.Create(ctx(), writePerm))

	// Assign only read permission to role
	require.NoError(t, roleRepo.AssignPermissions(ctx(), testRole.ID, []uint64{readPerm.ID}))

	// Assign role to user
	require.NoError(t, roleRepo.AssignToUser(ctx(), testUser.ID, []uint64{testRole.ID}))

	// Create services
	permSvc := adminservice.NewPermissionService(permRepo, memCache)
	roleSvc := adminservice.NewRoleService(roleRepo, memCache)

	// Test: Query user permissions
	t.Run("ListPermissionsByUserID returns assigned permissions", func(t *testing.T) {
		permissions, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, permissions, 1)
		assert.Equal(t, readPerm.Code, permissions[0].Code)
	})

	t.Run("CheckUserHasPermission returns true for assigned permission", func(t *testing.T) {
		hasPermission, err := permSvc.CheckUserHasPermission(ctx(), testUser.ID, model.HTTPMethodGET, "/api/v1/admin/test-resource")
		require.NoError(t, err)
		assert.True(t, hasPermission)
	})

	t.Run("CheckUserHasPermission returns false for unassigned permission", func(t *testing.T) {
		hasPermission, err := permSvc.CheckUserHasPermission(ctx(), testUser.ID, model.HTTPMethodPOST, "/api/v1/admin/test-resource")
		require.NoError(t, err)
		assert.False(t, hasPermission)
	})

	t.Run("ListRolesByUserID returns assigned roles", func(t *testing.T) {
		roles, err := roleSvc.ListRolesByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
		assert.Equal(t, testRole.Slug, roles[0].Slug)
	})
}

// TestAPIPermissionVerification tests API permission verification with middleware
// Requirements: 4.1, 4.2 - API permission verification for authorized and unauthorized scenarios
func TestAPIPermissionVerification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionVerificationModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()

	// Create users
	authorizedUser := &model.User{
		Name:         "AuthorizedUser",
		Email:        "authorized@example.com",
		Phone:        "18800000201",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	unauthorizedUser := &model.User{
		Name:         "UnauthorizedUser",
		Email:        "unauthorized@example.com",
		Phone:        "18800000202",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), authorizedUser))
	require.NoError(t, userRepo.Create(ctx(), unauthorizedUser))

	// Create role with permission
	authorizedRole := &model.RoleModel{
		Slug:     "authorized-role",
		Name:     "Authorized Role",
		IsSystem: false,
	}
	require.NoError(t, roleRepo.Create(ctx(), authorizedRole))

	// Create protected permission
	protectedPerm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/protected",
		Code:        "admin.protected.read",
		Group:       "admin",
		Description: "Access protected resource",
	}
	require.NoError(t, permRepo.Create(ctx(), protectedPerm))

	// Assign permission to role and role to authorized user
	require.NoError(t, roleRepo.AssignPermissions(ctx(), authorizedRole.ID, []uint64{protectedPerm.ID}))
	require.NoError(t, roleRepo.AssignToUser(ctx(), authorizedUser.ID, []uint64{authorizedRole.ID}))

	// Create services and middleware
	permSvc := adminservice.NewPermissionService(permRepo, memCache)
	roleSvc := adminservice.NewRoleService(roleRepo, memCache)
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	protectedHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}

	t.Run("Authorized user can access protected API", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected",
			setUserID(authorizedUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/protected"),
			protectedHandler,
		)

		resp := doJSON(router, http.MethodGet, "/protected", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Unauthorized user receives 403 Forbidden", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected",
			setUserID(unauthorizedUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/protected"),
			protectedHandler,
		)

		resp := doJSON(router, http.MethodGet, "/protected", nil, "")
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("API returns standardized error response for unauthorized access", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected",
			setUserID(unauthorizedUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/protected"),
			protectedHandler,
		)

		resp := doJSON(router, http.MethodGet, "/protected", nil, "")
		assert.Equal(t, http.StatusForbidden, resp.Code)

		var response map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		require.NoError(t, err)
		// Verify error response structure
		assert.Contains(t, response, "success")
		assert.Equal(t, false, response["success"])
	})
}

// TestSuperAdminBypass tests that super admin bypasses all permission checks
// Requirements: 4.4 - WHEN 超级管理员请求任何 API THEN 系统 SHALL 跳过权限检查直接放行
func TestSuperAdminBypass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionVerificationModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()

	// Create super admin user
	superAdminUser := &model.User{
		Name:         "SuperAdmin",
		Email:        "superadmin@example.com",
		Phone:        "18800000301",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), superAdminUser))

	// Create regular user without any permissions
	regularUser := &model.User{
		Name:         "RegularUser",
		Email:        "regular@example.com",
		Phone:        "18800000302",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), regularUser))

	// Create super admin role
	superAdminRole := &model.RoleModel{
		Slug:     string(model.RoleSlugSuperAdmin),
		Name:     "Super Admin",
		IsSystem: true,
	}
	require.NoError(t, roleRepo.Create(ctx(), superAdminRole))

	// Assign super admin role to super admin user
	require.NoError(t, roleRepo.AssignToUser(ctx(), superAdminUser.ID, []uint64{superAdminRole.ID}))

	// Create a permission that is NOT assigned to anyone
	restrictedPerm := &model.Permission{
		Method:      model.HTTPMethodDELETE,
		Path:        "/api/v1/admin/restricted",
		Code:        "admin.restricted.delete",
		Group:       "admin",
		Description: "Delete restricted resource",
	}
	require.NoError(t, permRepo.Create(ctx(), restrictedPerm))

	// Create services and middleware
	permSvc := adminservice.NewPermissionService(permRepo, memCache)
	roleSvc := adminservice.NewRoleService(roleRepo, memCache)
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	restrictedHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "restricted action performed"})
	}

	t.Run("Super admin bypasses permission check", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/restricted",
			setUserID(superAdminUser.ID),
			pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/restricted"),
			restrictedHandler,
		)

		resp := doJSON(router, http.MethodDelete, "/restricted", nil, "")
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Regular user without permission is blocked", func(t *testing.T) {
		router := gin.New()
		router.DELETE("/restricted",
			setUserID(regularUser.ID),
			pm.RequirePermission(model.HTTPMethodDELETE, "/api/v1/admin/restricted"),
			restrictedHandler,
		)

		resp := doJSON(router, http.MethodDelete, "/restricted", nil, "")
		assert.Equal(t, http.StatusForbidden, resp.Code)
	})

	t.Run("CheckUserIsSuperAdmin returns true for super admin", func(t *testing.T) {
		isSuperAdmin, err := roleSvc.CheckUserIsSuperAdmin(ctx(), superAdminUser.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)
	})

	t.Run("CheckUserIsSuperAdmin returns false for regular user", func(t *testing.T) {
		isSuperAdmin, err := roleSvc.CheckUserIsSuperAdmin(ctx(), regularUser.ID)
		require.NoError(t, err)
		assert.False(t, isSuperAdmin)
	})
}

// TestPermissionCacheInvalidation tests that permission cache is properly invalidated
// Requirements: 5.1, 5.2 - Cache invalidation when user roles or role permissions change
func TestPermissionCacheInvalidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionVerificationModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()

	// Create test user
	testUser := &model.User{
		Name:         "CacheTestUser",
		Email:        "cachetest@example.com",
		Phone:        "18800000401",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), testUser))

	// Create role
	testRole := &model.RoleModel{
		Slug:     "cache-test-role",
		Name:     "Cache Test Role",
		IsSystem: false,
	}
	require.NoError(t, roleRepo.Create(ctx(), testRole))

	// Create permissions
	perm1 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/cache-test-1",
		Code:        "cache.test.read1",
		Group:       "cache",
		Description: "Cache test permission 1",
	}
	perm2 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/cache-test-2",
		Code:        "cache.test.read2",
		Group:       "cache",
		Description: "Cache test permission 2",
	}
	require.NoError(t, permRepo.Create(ctx(), perm1))
	require.NoError(t, permRepo.Create(ctx(), perm2))

	// Assign perm1 to role and role to user
	require.NoError(t, roleRepo.AssignPermissions(ctx(), testRole.ID, []uint64{perm1.ID}))
	require.NoError(t, roleRepo.AssignToUser(ctx(), testUser.ID, []uint64{testRole.ID}))

	// Create services
	permSvc := adminservice.NewPermissionService(permRepo, memCache)
	roleSvc := adminservice.NewRoleService(roleRepo, memCache)

	t.Run("User permissions are cached after first query", func(t *testing.T) {
		// First query - should hit database and cache
		permissions1, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, permissions1, 1)
		assert.Equal(t, perm1.Code, permissions1[0].Code)

		// Second query - should hit cache
		permissions2, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, permissions2, 1)
		assert.Equal(t, perm1.Code, permissions2[0].Code)
	})

	t.Run("Cache is invalidated when role permissions change", func(t *testing.T) {
		// Initial state - user has perm1
		permissions, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, permissions, 1)

		// Add perm2 to role (this should invalidate cache)
		err = roleSvc.AssignPermissionsToRole(ctx(), testRole.ID, []uint64{perm1.ID, perm2.ID})
		require.NoError(t, err)

		// Query again - should reflect new permissions
		permissions, err = permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Len(t, permissions, 2)
	})

	t.Run("Cache is invalidated when user roles change", func(t *testing.T) {
		// Create another role with different permission
		anotherRole := &model.RoleModel{
			Slug:     "another-cache-role",
			Name:     "Another Cache Role",
			IsSystem: false,
		}
		require.NoError(t, roleRepo.Create(ctx(), anotherRole))

		perm3 := &model.Permission{
			Method:      model.HTTPMethodPOST,
			Path:        "/api/v1/admin/cache-test-3",
			Code:        "cache.test.write3",
			Group:       "cache",
			Description: "Cache test permission 3",
		}
		require.NoError(t, permRepo.Create(ctx(), perm3))
		require.NoError(t, roleRepo.AssignPermissions(ctx(), anotherRole.ID, []uint64{perm3.ID}))

		// Get current permissions count
		permsBefore, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		countBefore := len(permsBefore)

		// Assign another role to user (this should invalidate cache)
		err = roleSvc.AssignRolesToUser(ctx(), testUser.ID, []uint64{testRole.ID, anotherRole.ID})
		require.NoError(t, err)

		// Query again - should reflect new role's permissions
		permsAfter, err := permSvc.ListPermissionsByUserID(ctx(), testUser.ID)
		require.NoError(t, err)
		assert.Greater(t, len(permsAfter), countBefore)
	})
}

// TestMultiRolePermissionMerge tests that permissions from multiple roles are merged correctly
// Requirements: 10.1 - WHEN 用户拥有多个角色 THEN 系统 SHALL 合并所有角色的权限（并集）
func TestMultiRolePermissionMerge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionVerificationModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()

	// Create test user
	multiRoleUser := &model.User{
		Name:         "MultiRoleUser",
		Email:        "multirole@example.com",
		Phone:        "18800000501",
		PasswordHash: "x",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
	}
	require.NoError(t, userRepo.Create(ctx(), multiRoleUser))

	// Create multiple roles
	role1 := &model.RoleModel{
		Slug:     "multi-role-1",
		Name:     "Multi Role 1",
		IsSystem: false,
		Priority: 10,
	}
	role2 := &model.RoleModel{
		Slug:     "multi-role-2",
		Name:     "Multi Role 2",
		IsSystem: false,
		Priority: 20,
	}
	role3 := &model.RoleModel{
		Slug:     "multi-role-3",
		Name:     "Multi Role 3",
		IsSystem: false,
		Priority: 5,
	}
	require.NoError(t, roleRepo.Create(ctx(), role1))
	require.NoError(t, roleRepo.Create(ctx(), role2))
	require.NoError(t, roleRepo.Create(ctx(), role3))

	// Create permissions for each role
	perm1 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/resource-a",
		Code:        "multi.resource.a",
		Group:       "multi",
		Description: "Resource A permission",
	}
	perm2 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/resource-b",
		Code:        "multi.resource.b",
		Group:       "multi",
		Description: "Resource B permission",
	}
	perm3 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/resource-c",
		Code:        "multi.resource.c",
		Group:       "multi",
		Description: "Resource C permission",
	}
	sharedPerm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/shared-resource",
		Code:        "multi.resource.shared",
		Group:       "multi",
		Description: "Shared permission across roles",
	}
	require.NoError(t, permRepo.Create(ctx(), perm1))
	require.NoError(t, permRepo.Create(ctx(), perm2))
	require.NoError(t, permRepo.Create(ctx(), perm3))
	require.NoError(t, permRepo.Create(ctx(), sharedPerm))

	// Assign permissions to roles (with some overlap)
	require.NoError(t, roleRepo.AssignPermissions(ctx(), role1.ID, []uint64{perm1.ID, sharedPerm.ID}))
	require.NoError(t, roleRepo.AssignPermissions(ctx(), role2.ID, []uint64{perm2.ID, sharedPerm.ID}))
	require.NoError(t, roleRepo.AssignPermissions(ctx(), role3.ID, []uint64{perm3.ID}))

	// Assign all roles to user
	require.NoError(t, roleRepo.AssignToUser(ctx(), multiRoleUser.ID, []uint64{role1.ID, role2.ID, role3.ID}))

	// Create services
	permSvc := adminservice.NewPermissionService(permRepo, memCache)
	roleSvc := adminservice.NewRoleService(roleRepo, memCache)
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	t.Run("User has all permissions from all assigned roles", func(t *testing.T) {
		permissions, err := permSvc.ListPermissionsByUserID(ctx(), multiRoleUser.ID)
		require.NoError(t, err)

		// Should have 4 unique permissions (perm1, perm2, perm3, sharedPerm)
		assert.Len(t, permissions, 4)

		// Verify all permission codes are present
		codes := make(map[string]bool)
		for _, p := range permissions {
			codes[p.Code] = true
		}
		assert.True(t, codes["multi.resource.a"])
		assert.True(t, codes["multi.resource.b"])
		assert.True(t, codes["multi.resource.c"])
		assert.True(t, codes["multi.resource.shared"])
	})

	t.Run("Shared permissions are deduplicated", func(t *testing.T) {
		permissions, err := permSvc.ListPermissionsByUserID(ctx(), multiRoleUser.ID)
		require.NoError(t, err)

		// Count occurrences of shared permission
		sharedCount := 0
		for _, p := range permissions {
			if p.Code == "multi.resource.shared" {
				sharedCount++
			}
		}
		assert.Equal(t, 1, sharedCount, "Shared permission should appear only once")
	})

	t.Run("User can access resources from any assigned role", func(t *testing.T) {
		handler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		}

		// Test access to resource A (from role1)
		router1 := gin.New()
		router1.GET("/resource-a",
			setUserID(multiRoleUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/resource-a"),
			handler,
		)
		resp1 := doJSON(router1, http.MethodGet, "/resource-a", nil, "")
		assert.Equal(t, http.StatusOK, resp1.Code)

		// Test access to resource B (from role2)
		router2 := gin.New()
		router2.GET("/resource-b",
			setUserID(multiRoleUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/resource-b"),
			handler,
		)
		resp2 := doJSON(router2, http.MethodGet, "/resource-b", nil, "")
		assert.Equal(t, http.StatusOK, resp2.Code)

		// Test access to resource C (from role3)
		router3 := gin.New()
		router3.GET("/resource-c",
			setUserID(multiRoleUser.ID),
			pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/resource-c"),
			handler,
		)
		resp3 := doJSON(router3, http.MethodGet, "/resource-c", nil, "")
		assert.Equal(t, http.StatusOK, resp3.Code)
	})

	t.Run("ListRolesByUserID returns all assigned roles", func(t *testing.T) {
		roles, err := roleSvc.ListRolesByUserID(ctx(), multiRoleUser.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 3)
	})
}

// migratePermissionVerificationModels migrates all models needed for permission verification tests
func migratePermissionVerificationModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
	); err != nil {
		t.Fatalf("migrate permission verification models: %v", err)
	}
}

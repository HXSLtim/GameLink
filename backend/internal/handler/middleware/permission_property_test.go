package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/permission"
	"gamelink/internal/repository/user"
	permissionservice "gamelink/internal/service/admin"
	roleservice "gamelink/internal/service/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// TestPermissionVerificationConsistency tests Property 8: Permission Verification Consistency
// **Feature: review-management-module, Property 8: 权限验证一致性**
// **Validates: Requirements 10.1, 10.2, 10.3, 10.4**
// For any operation requiring permissions, the system must verify user permissions first,
// and when verification fails, must return 403 error without executing any business logic
func TestPermissionVerificationConsistency(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// Property: Permission check must occur before business logic execution
	properties.Property("permission check occurs before handler execution", prop.ForAll(
		func(hasPermission bool) bool {
			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			// Create test user and role
			testRole := &model.RoleModel{Slug: "test_role", Name: "Test Role", IsSystem: false}
			_ = roleRepo.Create(context.Background(), testRole)

			testUser := &model.User{
				Name:         "Test User",
				Email:        "test@example.com",
				Phone:        "19900000001",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			_ = userRepo.Create(context.Background(), testUser)
			_ = roleRepo.AssignToUser(context.Background(), testUser.ID, []uint64{testRole.ID})

			// Create test permission
			testPerm := &model.Permission{
				Method:      model.HTTPMethodGET,
				Path:        "/api/v1/admin/test",
				Code:        "test:view",
				Group:       "/admin/test",
				Description: "Test permission",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			// Assign permission to role if hasPermission is true
			if hasPermission {
				_ = roleRepo.AssignPermissions(context.Background(), testRole.ID, []uint64{testPerm.ID})
			}

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			// Track if handler was executed
			handlerExecuted := false
			handler := func(c *gin.Context) {
				handlerExecuted = true
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			router.GET("/test", setUserID(testUser.ID), pm.RequirePermission(testPerm.Method, testPerm.Path), handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Property: If user has permission, handler should execute and return 200
			// If user lacks permission, handler should NOT execute and return 403
			if hasPermission {
				return w.Code == http.StatusOK && handlerExecuted
			}
			return w.Code == http.StatusForbidden && !handlerExecuted
		},
		gen.Bool(),
	))

	// Property: Unauthorized access returns 403 and does not execute business logic
	properties.Property("unauthorized access returns 403 without executing logic", prop.ForAll(
		func(permissionCode string) bool {
			// Skip empty permission codes
			if permissionCode == "" {
				return true
			}

			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			// Create user without any permissions
			testUser := &model.User{
				Name:         "Unprivileged User",
				Email:        "unprivileged@example.com",
				Phone:        "19900000002",
				PasswordHash: "x",
				Role:         model.RoleUser,
			}
			_ = userRepo.Create(context.Background(), testUser)

			// Create permission
			testPerm := &model.Permission{
				Method:      model.HTTPMethodGET,
				Path:        "/api/v1/admin/test",
				Code:        permissionCode,
				Group:       "/admin/test",
				Description: "Test permission",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			// Track if handler was executed
			handlerExecuted := false
			handler := func(c *gin.Context) {
				handlerExecuted = true
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			router.GET("/test", setUserID(testUser.ID), pm.RequirePermission(testPerm.Method, testPerm.Path), handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Property: Must return 403 and handler must not execute
			return w.Code == http.StatusForbidden && !handlerExecuted
		},
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) < 50 }),
	))

	// Property: Super admin bypasses permission checks
	properties.Property("super admin has access to all operations", prop.ForAll(
		func(method model.HTTPMethod, path string) bool {
			// Skip invalid paths
			if path == "" || len(path) > 100 {
				return true
			}

			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			// Create super admin role
			superRole := &model.RoleModel{
				Slug:     string(model.RoleSlugSuperAdmin),
				Name:     "Super Admin",
				IsSystem: true,
			}
			_ = roleRepo.Create(context.Background(), superRole)

			// Create super admin user
			superUser := &model.User{
				Name:         "Super Admin",
				Email:        "super@example.com",
				Phone:        "19900000003",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			_ = userRepo.Create(context.Background(), superUser)
			_ = roleRepo.AssignToUser(context.Background(), superUser.ID, []uint64{superRole.ID})

			// Create a permission (super admin should have access even without explicit assignment)
			testPerm := &model.Permission{
				Method:      method,
				Path:        path,
				Code:        "test:operation",
				Group:       "/admin/test",
				Description: "Test operation",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			handlerExecuted := false
			handler := func(c *gin.Context) {
				handlerExecuted = true
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			
			// Register route based on method
			switch method {
			case model.HTTPMethodGET:
				router.GET("/test", setUserID(superUser.ID), pm.RequirePermission(method, path), handler)
			case model.HTTPMethodPOST:
				router.POST("/test", setUserID(superUser.ID), pm.RequirePermission(method, path), handler)
			case model.HTTPMethodPUT:
				router.PUT("/test", setUserID(superUser.ID), pm.RequirePermission(method, path), handler)
			case model.HTTPMethodDELETE:
				router.DELETE("/test", setUserID(superUser.ID), pm.RequirePermission(method, path), handler)
			default:
				return true // Skip unsupported methods
			}

			req := httptest.NewRequest(string(method), "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Property: Super admin should always have access
			return w.Code == http.StatusOK && handlerExecuted
		},
		gen.OneConstOf(
			model.HTTPMethodGET,
			model.HTTPMethodPOST,
			model.HTTPMethodPUT,
			model.HTTPMethodDELETE,
		),
		gen.AlphaString().SuchThat(func(s string) bool { return len(s) > 0 && len(s) < 100 }),
	))

	// Property: Permission verification is consistent across different HTTP methods
	properties.Property("permission verification is consistent across HTTP methods", prop.ForAll(
		func(method model.HTTPMethod, hasPermission bool) bool {
			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			testRole := &model.RoleModel{Slug: "test_role", Name: "Test Role", IsSystem: false}
			_ = roleRepo.Create(context.Background(), testRole)

			testUser := &model.User{
				Name:         "Test User",
				Email:        "test@example.com",
				Phone:        "19900000004",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			_ = userRepo.Create(context.Background(), testUser)
			_ = roleRepo.AssignToUser(context.Background(), testUser.ID, []uint64{testRole.ID})

			testPerm := &model.Permission{
				Method:      method,
				Path:        "/api/v1/admin/test",
				Code:        "test:operation",
				Group:       "/admin/test",
				Description: "Test operation",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			if hasPermission {
				_ = roleRepo.AssignPermissions(context.Background(), testRole.ID, []uint64{testPerm.ID})
			}

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			handlerExecuted := false
			handler := func(c *gin.Context) {
				handlerExecuted = true
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			
			switch method {
			case model.HTTPMethodGET:
				router.GET("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodPOST:
				router.POST("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodPUT:
				router.PUT("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodDELETE:
				router.DELETE("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			default:
				return true // Skip unsupported methods
			}

			req := httptest.NewRequest(string(method), "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Property: Behavior should be consistent regardless of HTTP method
			if hasPermission {
				return w.Code == http.StatusOK && handlerExecuted
			}
			return w.Code == http.StatusForbidden && !handlerExecuted
		},
		gen.OneConstOf(
			model.HTTPMethodGET,
			model.HTTPMethodPOST,
			model.HTTPMethodPUT,
			model.HTTPMethodDELETE,
		),
		gen.Bool(),
	))

	// Property: Multiple permission checks are independent
	properties.Property("multiple permission checks are independent", prop.ForAll(
		func(hasPerm1, hasPerm2 bool) bool {
			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			testRole := &model.RoleModel{Slug: "test_role", Name: "Test Role", IsSystem: false}
			_ = roleRepo.Create(context.Background(), testRole)

			testUser := &model.User{
				Name:         "Test User",
				Email:        "test@example.com",
				Phone:        "19900000005",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			_ = userRepo.Create(context.Background(), testUser)
			_ = roleRepo.AssignToUser(context.Background(), testUser.ID, []uint64{testRole.ID})

			// Create two different permissions
			perm1 := &model.Permission{
				Method:      model.HTTPMethodGET,
				Path:        "/api/v1/admin/test1",
				Code:        "test:view1",
				Group:       "/admin/test",
				Description: "Test permission 1",
			}
			perm2 := &model.Permission{
				Method:      model.HTTPMethodGET,
				Path:        "/api/v1/admin/test2",
				Code:        "test:view2",
				Group:       "/admin/test",
				Description: "Test permission 2",
			}
			_ = permRepo.Create(context.Background(), perm1)
			_ = permRepo.Create(context.Background(), perm2)

			// Assign permissions based on flags
			var permIDs []uint64
			if hasPerm1 {
				permIDs = append(permIDs, perm1.ID)
			}
			if hasPerm2 {
				permIDs = append(permIDs, perm2.ID)
			}
			if len(permIDs) > 0 {
				_ = roleRepo.AssignPermissions(context.Background(), testRole.ID, permIDs)
			}

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			handler := func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			router.GET("/test1", setUserID(testUser.ID), pm.RequirePermission(perm1.Method, perm1.Path), handler)
			router.GET("/test2", setUserID(testUser.ID), pm.RequirePermission(perm2.Method, perm2.Path), handler)

			// Test first permission
			req1 := httptest.NewRequest("GET", "/test1", nil)
			w1 := httptest.NewRecorder()
			router.ServeHTTP(w1, req1)

			// Test second permission
			req2 := httptest.NewRequest("GET", "/test2", nil)
			w2 := httptest.NewRecorder()
			router.ServeHTTP(w2, req2)

			// Property: Each permission check should be independent
			expectedCode1 := http.StatusForbidden
			if hasPerm1 {
				expectedCode1 = http.StatusOK
			}
			expectedCode2 := http.StatusForbidden
			if hasPerm2 {
				expectedCode2 = http.StatusOK
			}

			return w1.Code == expectedCode1 && w2.Code == expectedCode2
		},
		gen.Bool(),
		gen.Bool(),
	))

	// Property: Permission check failure does not affect subsequent requests
	properties.Property("permission check failure does not affect subsequent requests", prop.ForAll(
		func(firstHasPermission, secondHasPermission bool) bool {
			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			// Create two users with different permissions
			role1 := &model.RoleModel{Slug: "role1", Name: "Role 1", IsSystem: false}
			role2 := &model.RoleModel{Slug: "role2", Name: "Role 2", IsSystem: false}
			_ = roleRepo.Create(context.Background(), role1)
			_ = roleRepo.Create(context.Background(), role2)

			user1 := &model.User{
				Name:         "User 1",
				Email:        "user1@example.com",
				Phone:        "19900000006",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			user2 := &model.User{
				Name:         "User 2",
				Email:        "user2@example.com",
				Phone:        "19900000007",
				PasswordHash: "x",
				Role:         model.RoleAdmin,
			}
			_ = userRepo.Create(context.Background(), user1)
			_ = userRepo.Create(context.Background(), user2)
			_ = roleRepo.AssignToUser(context.Background(), user1.ID, []uint64{role1.ID})
			_ = roleRepo.AssignToUser(context.Background(), user2.ID, []uint64{role2.ID})

			testPerm := &model.Permission{
				Method:      model.HTTPMethodGET,
				Path:        "/api/v1/admin/test",
				Code:        "test:view",
				Group:       "/admin/test",
				Description: "Test permission",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			if firstHasPermission {
				_ = roleRepo.AssignPermissions(context.Background(), role1.ID, []uint64{testPerm.ID})
			}
			if secondHasPermission {
				_ = roleRepo.AssignPermissions(context.Background(), role2.ID, []uint64{testPerm.ID})
			}

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			handler := func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			}

			// First request from user1
			router1 := gin.New()
			router1.GET("/test", setUserID(user1.ID), pm.RequirePermission(testPerm.Method, testPerm.Path), handler)
			req1 := httptest.NewRequest("GET", "/test", nil)
			w1 := httptest.NewRecorder()
			router1.ServeHTTP(w1, req1)

			// Second request from user2
			router2 := gin.New()
			router2.GET("/test", setUserID(user2.ID), pm.RequirePermission(testPerm.Method, testPerm.Path), handler)
			req2 := httptest.NewRequest("GET", "/test", nil)
			w2 := httptest.NewRecorder()
			router2.ServeHTTP(w2, req2)

			// Property: Each request should be evaluated independently
			expectedCode1 := http.StatusForbidden
			if firstHasPermission {
				expectedCode1 = http.StatusOK
			}
			expectedCode2 := http.StatusForbidden
			if secondHasPermission {
				expectedCode2 = http.StatusOK
			}

			return w1.Code == expectedCode1 && w2.Code == expectedCode2
		},
		gen.Bool(),
		gen.Bool(),
	))

	// Property: Permission verification always returns 403 for unauthorized access
	properties.Property("unauthorized access always returns 403", prop.ForAll(
		func(method model.HTTPMethod) bool {
			gin.SetMode(gin.TestMode)
			db := testutil.NewMemoryDB(t)
			defer testutil.CleanDB(t, db)
			migratePermissionModels(t, db)

			userRepo := user.NewUserRepository(db)
			roleRepo := adminrepo.NewRoleRepository(db)
			permRepo := permission.NewPermissionRepository(db)

			// Create user without permissions
			testUser := &model.User{
				Name:         "Unprivileged User",
				Email:        "unprivileged@example.com",
				Phone:        "19900000008",
				PasswordHash: "x",
				Role:         model.RoleUser,
			}
			_ = userRepo.Create(context.Background(), testUser)

			testPerm := &model.Permission{
				Method:      method,
				Path:        "/api/v1/admin/test",
				Code:        "test:operation",
				Group:       "/admin/test",
				Description: "Test operation",
			}
			_ = permRepo.Create(context.Background(), testPerm)

			permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
			roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
			pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

			handler := func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			}

			router := gin.New()
			
			switch method {
			case model.HTTPMethodGET:
				router.GET("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodPOST:
				router.POST("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodPUT:
				router.PUT("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			case model.HTTPMethodDELETE:
				router.DELETE("/test", setUserID(testUser.ID), pm.RequirePermission(method, testPerm.Path), handler)
			default:
				return true
			}

			req := httptest.NewRequest(string(method), "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Property: Must always return 403 for unauthorized access
			return w.Code == http.StatusForbidden
		},
		gen.OneConstOf(
			model.HTTPMethodGET,
			model.HTTPMethodPOST,
			model.HTTPMethodPUT,
			model.HTTPMethodDELETE,
		),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper function to set user ID in context
func setUserID(userID uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.UserIDKey, userID)
		c.Next()
	}
}

// Helper function to migrate permission models
func migratePermissionModels(t *testing.T, db interface{}) {
	type migrator interface {
		AutoMigrate(dst ...interface{}) error
	}
	
	if m, ok := db.(migrator); ok {
		if err := m.AutoMigrate(
			&model.User{},
			&model.RoleModel{},
			&model.Permission{},
			&model.UserRole{},
			&model.RolePermission{},
		); err != nil {
			t.Fatalf("Failed to migrate models: %v", err)
		}
	}
}

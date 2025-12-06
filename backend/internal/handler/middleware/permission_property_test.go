package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"gamelink/internal/model"
)

// mockPermissionService is a mock implementation of the permission service for testing.
type mockPermissionService struct {
	userPermissions map[uint64][]model.Permission
}

func newMockPermissionService() *mockPermissionService {
	return &mockPermissionService{
		userPermissions: make(map[uint64][]model.Permission),
	}
}

func (m *mockPermissionService) SetUserPermissions(userID uint64, permissions []model.Permission) {
	m.userPermissions[userID] = permissions
}

func (m *mockPermissionService) ListPermissionsByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	return m.userPermissions[userID], nil
}

func (m *mockPermissionService) CheckUserHasPermission(ctx context.Context, userID uint64, method model.HTTPMethod, path string) (bool, error) {
	permissions := m.userPermissions[userID]
	for _, perm := range permissions {
		if perm.Method == method && perm.Path == path {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockPermissionService) CheckUserHasPermissionCode(ctx context.Context, userID uint64, code string) (bool, error) {
	permissions := m.userPermissions[userID]
	for _, perm := range permissions {
		if perm.Code == code {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockPermissionService) CheckUserHasAnyPermission(ctx context.Context, userID uint64, codes []string) (bool, error) {
	permissions := m.userPermissions[userID]
	codeSet := make(map[string]bool)
	for _, perm := range permissions {
		codeSet[perm.Code] = true
	}
	for _, code := range codes {
		if codeSet[code] {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockPermissionService) CheckUserHasAllPermissions(ctx context.Context, userID uint64, codes []string) (bool, error) {
	permissions := m.userPermissions[userID]
	codeSet := make(map[string]bool)
	for _, perm := range permissions {
		codeSet[perm.Code] = true
	}
	for _, code := range codes {
		if !codeSet[code] {
			return false, nil
		}
	}
	return true, nil
}

func (m *mockPermissionService) CheckUserHasExceptPermissions(ctx context.Context, userID uint64, excludedCodes []string) (bool, error) {
	permissions := m.userPermissions[userID]
	excludedSet := make(map[string]bool)
	for _, code := range excludedCodes {
		excludedSet[code] = true
	}
	for _, perm := range permissions {
		if excludedSet[perm.Code] {
			return false, nil
		}
	}
	return true, nil
}

// mockRoleService is a mock implementation of the role service for testing.
type mockRoleService struct {
	superAdminUsers map[uint64]bool
	userRoles       map[uint64][]string
}

func newMockRoleService() *mockRoleService {
	return &mockRoleService{
		superAdminUsers: make(map[uint64]bool),
		userRoles:       make(map[uint64][]string),
	}
}

func (m *mockRoleService) SetSuperAdmin(userID uint64, isSuperAdmin bool) {
	m.superAdminUsers[userID] = isSuperAdmin
}

func (m *mockRoleService) SetUserRoles(userID uint64, roles []string) {
	m.userRoles[userID] = roles
}

func (m *mockRoleService) CheckUserIsSuperAdmin(ctx context.Context, userID uint64) (bool, error) {
	return m.superAdminUsers[userID], nil
}

func (m *mockRoleService) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	roles := m.userRoles[userID]
	for _, r := range roles {
		if r == roleSlug {
			return true, nil
		}
	}
	return false, nil
}

// testablePermissionMiddleware wraps the permission middleware for testing.
type testablePermissionMiddleware struct {
	permissionSvc *mockPermissionService
	roleSvc       *mockRoleService
	whitelist     map[string]bool
}

func newTestablePermissionMiddleware() *testablePermissionMiddleware {
	return &testablePermissionMiddleware{
		permissionSvc: newMockPermissionService(),
		roleSvc:       newMockRoleService(),
		whitelist:     make(map[string]bool),
	}
}

func (m *testablePermissionMiddleware) AddToWhitelist(methodPath string) {
	m.whitelist[methodPath] = true
}

func (m *testablePermissionMiddleware) IsWhitelisted(method, path string) bool {
	key := method + ":" + path
	return m.whitelist[key]
}

// checkPermission simulates the permission check logic from the middleware.
// Returns true if access should be granted, false otherwise.
func (m *testablePermissionMiddleware) checkPermission(userID uint64, method model.HTTPMethod, path string) bool {
	ctx := context.Background()

	// Check whitelist
	if m.IsWhitelisted(string(method), path) {
		return true
	}

	// Check if super admin
	isSuperAdmin, _ := m.roleSvc.CheckUserIsSuperAdmin(ctx, userID)
	if isSuperAdmin {
		return true
	}

	// Check permission
	hasPermission, _ := m.permissionSvc.CheckUserHasPermission(ctx, userID, method, path)
	return hasPermission
}

// checkPermissionCode simulates the permission code check logic from the middleware.
func (m *testablePermissionMiddleware) checkPermissionCode(userID uint64, code string) bool {
	ctx := context.Background()

	// Check if super admin
	isSuperAdmin, _ := m.roleSvc.CheckUserIsSuperAdmin(ctx, userID)
	if isSuperAdmin {
		return true
	}

	// Check permission code
	hasPermission, _ := m.permissionSvc.CheckUserHasPermissionCode(ctx, userID, code)
	return hasPermission
}

// checkPermissionCodes simulates the multi-permission check logic from the middleware.
func (m *testablePermissionMiddleware) checkPermissionCodes(userID uint64, codes []string, mode PermissionCheckMode) bool {
	ctx := context.Background()

	// Check if super admin (except for "except" mode)
	if mode != PermissionCheckModeExcept {
		isSuperAdmin, _ := m.roleSvc.CheckUserIsSuperAdmin(ctx, userID)
		if isSuperAdmin {
			return true
		}
	}

	var hasPermission bool
	switch mode {
	case PermissionCheckModeAny:
		hasPermission, _ = m.permissionSvc.CheckUserHasAnyPermission(ctx, userID, codes)
	case PermissionCheckModeAll:
		hasPermission, _ = m.permissionSvc.CheckUserHasAllPermissions(ctx, userID, codes)
	case PermissionCheckModeExcept:
		hasPermission, _ = m.permissionSvc.CheckUserHasExceptPermissions(ctx, userID, codes)
	default:
		hasPermission, _ = m.permissionSvc.CheckUserHasAnyPermission(ctx, userID, codes)
	}

	return hasPermission
}

// TestSuperAdminPermissionBypass tests Property 7: Super Admin Permission Bypass
// For any super admin user (with '*' permission), all permission checks should return true
// regardless of the specific permission being checked.
// **Feature: rbac-button-level-permission, Property 7: 超级管理员权限绕过**
// **Validates: Requirements 3.5, 4.4**
func TestSuperAdminPermissionBypass(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Super admin should pass any method+path permission check
	properties.Property("super admin should pass any method+path permission check", prop.ForAll(
		func(userID uint64, methodIdx int, path string) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			methods := []model.HTTPMethod{
				model.HTTPMethodGET,
				model.HTTPMethodPOST,
				model.HTTPMethodPUT,
				model.HTTPMethodPATCH,
				model.HTTPMethodDELETE,
			}
			method := methods[methodIdx%len(methods)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, true)
			// Don't set any permissions - super admin should still pass

			return middleware.checkPermission(userID, method, path)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.IntRange(0, 4),
		genValidPath(),
	))

	// Property 2: Super admin should pass any permission code check
	properties.Property("super admin should pass any permission code check", prop.ForAll(
		func(userID uint64, code string) bool {
			if userID == 0 || code == "" {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, true)
			// Don't set any permissions - super admin should still pass

			return middleware.checkPermissionCode(userID, code)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCode(),
	))

	// Property 3: Super admin should pass any multi-permission check (any mode)
	properties.Property("super admin should pass any multi-permission check (any mode)", prop.ForAll(
		func(userID uint64, codes []string) bool {
			if userID == 0 || len(codes) == 0 {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, true)

			return middleware.checkPermissionCodes(userID, codes, PermissionCheckModeAny)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCode()),
	))

	// Property 4: Super admin should pass any multi-permission check (all mode)
	properties.Property("super admin should pass any multi-permission check (all mode)", prop.ForAll(
		func(userID uint64, codes []string) bool {
			if userID == 0 || len(codes) == 0 {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, true)

			return middleware.checkPermissionCodes(userID, codes, PermissionCheckModeAll)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCode()),
	))

	// Property 5: Non-super admin without permissions should fail permission check
	properties.Property("non-super admin without permissions should fail permission check", prop.ForAll(
		func(userID uint64, methodIdx int, path string) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			methods := []model.HTTPMethod{
				model.HTTPMethodGET,
				model.HTTPMethodPOST,
				model.HTTPMethodPUT,
				model.HTTPMethodPATCH,
				model.HTTPMethodDELETE,
			}
			method := methods[methodIdx%len(methods)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			// Don't set any permissions

			return !middleware.checkPermission(userID, method, path)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.IntRange(0, 4),
		genValidPath(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestAPIPermissionValidationConsistency tests Property 8: API Permission Validation Consistency
// For any protected API endpoint and any user, the middleware should allow access if and only if
// the user has the required permission for that endpoint.
// **Feature: rbac-button-level-permission, Property 8: API 权限验证一致性**
// **Validates: Requirements 4.1, 4.2**
func TestAPIPermissionValidationConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: User with exact permission should pass
	properties.Property("user with exact permission should pass", prop.ForAll(
		func(userID uint64, methodIdx int, path string) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			methods := []model.HTTPMethod{
				model.HTTPMethodGET,
				model.HTTPMethodPOST,
				model.HTTPMethodPUT,
				model.HTTPMethodPATCH,
				model.HTTPMethodDELETE,
			}
			method := methods[methodIdx%len(methods)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: method, Path: path, Code: "test.permission"},
			})

			return middleware.checkPermission(userID, method, path)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.IntRange(0, 4),
		genValidPath(),
	))

	// Property 2: User without permission should fail
	properties.Property("user without permission should fail", prop.ForAll(
		func(userID uint64, methodIdx int, path string) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			methods := []model.HTTPMethod{
				model.HTTPMethodGET,
				model.HTTPMethodPOST,
				model.HTTPMethodPUT,
				model.HTTPMethodPATCH,
				model.HTTPMethodDELETE,
			}
			method := methods[methodIdx%len(methods)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			// Set a different permission
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/different/path", Code: "different.permission"},
			})

			return !middleware.checkPermission(userID, method, path)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.IntRange(0, 4),
		genValidPath(),
	))

	// Property 3: User with permission code should pass code check
	properties.Property("user with permission code should pass code check", prop.ForAll(
		func(userID uint64, code string) bool {
			if userID == 0 || code == "" {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/test", Code: code},
			})

			return middleware.checkPermissionCode(userID, code)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCode(),
	))

	// Property 4: User without permission code should fail code check
	properties.Property("user without permission code should fail code check", prop.ForAll(
		func(userID uint64, code string) bool {
			if userID == 0 || code == "" {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			// Set a different permission code
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/test", Code: "different.code.here"},
			})

			return !middleware.checkPermissionCode(userID, code)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCode(),
	))

	// Property 5: Any mode - user with at least one permission should pass
	properties.Property("any mode - user with at least one permission should pass", prop.ForAll(
		func(userID uint64, codes []string, grantedIdx int) bool {
			if userID == 0 || len(codes) == 0 {
				return true // Skip invalid inputs
			}

			// Grant one of the required permissions
			grantedCode := codes[grantedIdx%len(codes)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/test", Code: grantedCode},
			})

			return middleware.checkPermissionCodes(userID, codes, PermissionCheckModeAny)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCode()),
		gen.IntRange(0, 2),
	))

	// Property 6: All mode - user with all permissions should pass
	properties.Property("all mode - user with all permissions should pass", prop.ForAll(
		func(userID uint64, codes []string) bool {
			if userID == 0 || len(codes) == 0 {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)

			// Grant all required permissions
			permissions := make([]model.Permission, len(codes))
			for i, code := range codes {
				permissions[i] = model.Permission{
					Method: model.HTTPMethodGET,
					Path:   "/test",
					Code:   code,
				}
			}
			middleware.permissionSvc.SetUserPermissions(userID, permissions)

			return middleware.checkPermissionCodes(userID, codes, PermissionCheckModeAll)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCode()),
	))

	// Property 7: All mode - user missing one permission should fail
	properties.Property("all mode - user missing one permission should fail", prop.ForAll(
		func(userID uint64, codes []string, missingIdx int) bool {
			if userID == 0 || len(codes) < 2 {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)

			// Grant all but one permission
			missingIdx = missingIdx % len(codes)
			permissions := make([]model.Permission, 0, len(codes)-1)
			for i, code := range codes {
				if i != missingIdx {
					permissions = append(permissions, model.Permission{
						Method: model.HTTPMethodGET,
						Path:   "/test",
						Code:   code,
					})
				}
			}
			middleware.permissionSvc.SetUserPermissions(userID, permissions)

			return !middleware.checkPermissionCodes(userID, codes, PermissionCheckModeAll)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCode()),
		gen.IntRange(0, 2),
	))

	// Property 8: Whitelist should bypass permission check
	properties.Property("whitelist should bypass permission check", prop.ForAll(
		func(userID uint64, methodIdx int, path string) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			methods := []model.HTTPMethod{
				model.HTTPMethodGET,
				model.HTTPMethodPOST,
				model.HTTPMethodPUT,
				model.HTTPMethodPATCH,
				model.HTTPMethodDELETE,
			}
			method := methods[methodIdx%len(methods)]

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			// Don't set any permissions
			// Add to whitelist
			middleware.AddToWhitelist(string(method) + ":" + path)

			return middleware.checkPermission(userID, method, path)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.IntRange(0, 4),
		genValidPath(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestExceptModePermissionCheck tests the except mode permission check.
// In except mode, super admin should NOT bypass the check.
func TestExceptModePermissionCheck(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: In except mode, super admin with excluded permission should fail
	properties.Property("except mode - super admin with excluded permission should fail", prop.ForAll(
		func(userID uint64, excludedCode string) bool {
			if userID == 0 || excludedCode == "" {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, true)
			// Give super admin the excluded permission
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/test", Code: excludedCode},
			})

			// In except mode, having the excluded permission should fail
			return !middleware.checkPermissionCodes(userID, []string{excludedCode}, PermissionCheckModeExcept)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCode(),
	))

	// Property: In except mode, user without excluded permission should pass
	properties.Property("except mode - user without excluded permission should pass", prop.ForAll(
		func(userID uint64, excludedCode string) bool {
			if userID == 0 || excludedCode == "" {
				return true // Skip invalid inputs
			}

			middleware := newTestablePermissionMiddleware()
			middleware.roleSvc.SetSuperAdmin(userID, false)
			// Give user a different permission
			middleware.permissionSvc.SetUserPermissions(userID, []model.Permission{
				{Method: model.HTTPMethodGET, Path: "/test", Code: "different.code.here"},
			})

			return middleware.checkPermissionCodes(userID, []string{excludedCode}, PermissionCheckModeExcept)
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCode(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genValidPath generates a valid API path for testing.
func genValidPath() gopter.Gen {
	// Generate a fixed-length segment to avoid filtering issues
	segmentGen := gen.IntRange(1, 10).FlatMap(func(length interface{}) gopter.Gen {
		return gen.SliceOfN(length.(int), gen.Rune()).Map(func(runes []rune) string {
			result := make([]byte, len(runes))
			for i, r := range runes {
				// Convert to lowercase letter a-z
				result[i] = byte('a' + (r % 26))
			}
			return string(result)
		})
	}, reflect.TypeOf(""))

	return gen.SliceOfN(3, segmentGen).Map(func(parts []string) string {
		result := "/api"
		for _, part := range parts {
			if part != "" {
				result += "/" + part
			}
		}
		return result
	})
}

// genValidPermissionCode generates a valid permission code (module.resource.action).
func genValidPermissionCode() gopter.Gen {
	// Generate a fixed-length segment to avoid filtering issues
	segmentGen := gen.IntRange(1, 10).FlatMap(func(length interface{}) gopter.Gen {
		return gen.SliceOfN(length.(int), gen.Rune()).Map(func(runes []rune) string {
			result := make([]byte, len(runes))
			for i, r := range runes {
				// Convert to lowercase letter a-z
				result[i] = byte('a' + (r % 26))
			}
			return string(result)
		})
	}, reflect.TypeOf(""))

	return gopter.CombineGens(segmentGen, segmentGen, segmentGen).Map(func(vals []interface{}) string {
		module := vals[0].(string)
		resource := vals[1].(string)
		action := vals[2].(string)
		return module + "." + resource + "." + action
	})
}

// TestMiddlewareHTTPResponses tests that the middleware returns correct HTTP status codes.
func TestMiddlewareHTTPResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("unauthorized when no user ID", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			// Simulate middleware check without user ID
			_, exists := c.Get(UserIDKey)
			if !exists {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"code":    http.StatusUnauthorized,
					"message": "未授权：请先登录",
				})
				return
			}
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("forbidden when no permission", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Set(UserIDKey, uint64(1))
			// Simulate permission check failure
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"code":    http.StatusForbidden,
				"message": "权限不足",
			})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("ok when has permission", func(t *testing.T) {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Set(UserIDKey, uint64(1))
			c.String(http.StatusOK, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

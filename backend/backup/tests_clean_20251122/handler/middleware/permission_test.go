package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"gamelink/internal/auth"
	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	permissionservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"
)

// ---- Local stubs to build real services ----

type testCache struct {
	values map[string]string
}

func newTestCache() *testCache {
	return &testCache{values: make(map[string]string)}
}

func (c *testCache) Get(_ context.Context, key string) (string, bool, error) {
	val, ok := c.values[key]
	return val, ok, nil
}

func (c *testCache) Set(_ context.Context, key, value string, _ time.Duration) error {
	c.values[key] = value
	return nil
}

func (c *testCache) Delete(_ context.Context, key string) error {
	delete(c.values, key)
	return nil
}

func (c *testCache) Close(context.Context) error { return nil }

var _ cache.Cache = (*testCache)(nil)

type stubPermissionRepository struct {
	repository.PermissionRepository
	permsByUser   map[uint64][]model.Permission
	listByUserErr error
}

func (r *stubPermissionRepository) List(ctx context.Context) ([]model.Permission, error) {
	return nil, nil
}

func (r *stubPermissionRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.Permission, int64, error) {
	return nil, 0, nil
}

func (r *stubPermissionRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword, method, group string) ([]model.Permission, int64, error) {
	return nil, 0, nil
}

func (r *stubPermissionRepository) ListByGroup(ctx context.Context) (map[string][]model.Permission, error) {
	return map[string][]model.Permission{}, nil
}

func (r *stubPermissionRepository) ListGroups(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (r *stubPermissionRepository) Get(ctx context.Context, id uint64) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (r *stubPermissionRepository) GetByResource(ctx context.Context, resource, action string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (r *stubPermissionRepository) GetByMethodAndPath(ctx context.Context, method, path string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (r *stubPermissionRepository) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	return nil, repository.ErrNotFound
}

func (r *stubPermissionRepository) Create(ctx context.Context, perm *model.Permission) error {
	return nil
}

func (r *stubPermissionRepository) CreateBatch(ctx context.Context, perms []model.Permission) error {
	return nil
}

func (r *stubPermissionRepository) Update(ctx context.Context, perm *model.Permission) error {
	return nil
}

func (r *stubPermissionRepository) UpsertByMethodPath(ctx context.Context, perm *model.Permission) error {
	return nil
}

func (r *stubPermissionRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (r *stubPermissionRepository) ListByRoleID(ctx context.Context, roleID uint64) ([]model.Permission, error) {
	return nil, nil
}

func (r *stubPermissionRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	if r.listByUserErr != nil {
		return nil, r.listByUserErr
	}
	if perms, ok := r.permsByUser[userID]; ok {
		return perms, nil
	}
	return nil, nil
}

type stubRoleRepository struct {
	repository.RoleRepository
	rolesByUser map[uint64]map[string]bool
	checkErr    error
}

func (r *stubRoleRepository) List(ctx context.Context) ([]model.RoleModel, error) { return nil, nil }
func (r *stubRoleRepository) ListPaged(ctx context.Context, page, pageSize int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (r *stubRoleRepository) ListPagedWithFilter(ctx context.Context, page, pageSize int, keyword string, isSystem *bool) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (r *stubRoleRepository) ListWithPermissions(ctx context.Context) ([]model.RoleModel, error) {
	return nil, nil
}
func (r *stubRoleRepository) Get(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (r *stubRoleRepository) GetWithPermissions(ctx context.Context, id uint64) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (r *stubRoleRepository) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	return nil, repository.ErrNotFound
}
func (r *stubRoleRepository) Create(ctx context.Context, role *model.RoleModel) error { return nil }
func (r *stubRoleRepository) Update(ctx context.Context, role *model.RoleModel) error { return nil }
func (r *stubRoleRepository) Delete(ctx context.Context, id uint64) error             { return nil }
func (r *stubRoleRepository) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (r *stubRoleRepository) AddPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (r *stubRoleRepository) RemovePermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	return nil
}
func (r *stubRoleRepository) AssignToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (r *stubRoleRepository) RemoveFromUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return nil
}
func (r *stubRoleRepository) ListByUserID(ctx context.Context, userID uint64) ([]model.RoleModel, error) {
	return nil, nil
}

func (r *stubRoleRepository) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	if r.checkErr != nil {
		return false, r.checkErr
	}
	if roles, ok := r.rolesByUser[userID]; ok {
		return roles[roleSlug], nil
	}
	return false, nil
}

var _ repository.RoleRepository = (*stubRoleRepository)(nil)
var _ repository.PermissionRepository = (*stubPermissionRepository)(nil)

func newMiddlewareWithServices(perms map[uint64][]model.Permission, roles map[uint64]map[string]bool, permErr error, roleErr error) *PermissionMiddleware {
	permRepo := &stubPermissionRepository{
		permsByUser:   perms,
		listByUserErr: permErr,
	}
	roleRepo := &stubRoleRepository{
		rolesByUser: roles,
		checkErr:    roleErr,
	}

	permSvc := permissionservice.NewPermissionService(permRepo, newTestCache())
	roleSvc := roleservice.NewRoleService(roleRepo, newTestCache())

	return &PermissionMiddleware{
		jwtManager:    auth.NewJWTManager("middleware-test-secret", 24*time.Hour),
		permissionSvc: permSvc,
		roleSvc:       roleSvc,
	}
}

// ---- Fake Services for Permission Tests ----

type fakePermissionService struct {
	permissions map[uint64]map[string]bool // userID -> permission -> hasPermission
}

func newFakePermissionService() *fakePermissionService {
	return &fakePermissionService{
		permissions: map[uint64]map[string]bool{
			1: { // Admin user
				"GET:/admin/users":    true,
				"POST:/admin/users":   true,
				"DELETE:/admin/users": true,
			},
			2: { // Regular user
				"GET:/user/profile": true,
			},
			3: {}, // User with no permissions
		},
	}
}

func (f *fakePermissionService) CheckUserHasPermission(ctx context.Context, userID uint64, method model.HTTPMethod, path string) (bool, error) {
	key := string(method) + ":" + path
	if perms, ok := f.permissions[userID]; ok {
		return perms[key], nil
	}
	return false, nil
}

type fakeRoleService struct {
	superAdmins map[uint64]bool
	userRoles   map[uint64]map[string]bool // userID -> roleSlug -> hasRole
}

func newFakeRoleService() *fakeRoleService {
	return &fakeRoleService{
		superAdmins: map[uint64]bool{
			999: true, // User 999 is super admin
		},
		userRoles: map[uint64]map[string]bool{
			1: {"admin": true},
			2: {"user": true},
			3: {"player": true},
		},
	}
}

func (f *fakeRoleService) CheckUserIsSuperAdmin(ctx context.Context, userID uint64) (bool, error) {
	return f.superAdmins[userID], nil
}

func (f *fakeRoleService) CheckUserHasRole(ctx context.Context, userID uint64, roleSlug string) (bool, error) {
	if roles, ok := f.userRoles[userID]; ok {
		return roles[roleSlug], nil
	}
	return false, nil
}

// ---- Helper Functions ----

func setupPermissionMiddleware() *PermissionMiddleware {
	jwtManager := auth.NewJWTManager("test-secret-key-for-testing-only", 24*time.Hour)
	_ = newFakePermissionService()
	_ = newFakeRoleService()

	// Type assertion to satisfy interface requirements
	var permService *permissionservice.PermissionService
	var roleService *roleservice.RoleService

	// Use unsafe type conversion for testing
	permService = (*permissionservice.PermissionService)(nil)
	roleService = (*roleservice.RoleService)(nil)

	return &PermissionMiddleware{
		jwtManager:    jwtManager,
		permissionSvc: permService,
		roleSvc:       roleService,
	}
}

func generateTestToken(jwtManager *auth.JWTManager, userID uint64, role string) string {
	token, _ := jwtManager.GenerateToken(userID, role)
	return token
}

// ---- Tests for RequireAuth ----

func TestRequireAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		userID, _ := c.Get(UserIDKey)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	token := generateTestToken(jwtManager, 123, "user")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequireAuth_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create JWT manager with very short expiry
	jwtManager := auth.NewJWTManager("test-secret", 1*time.Nanosecond)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	token := generateTestToken(jwtManager, 123, "user")
	time.Sleep(10 * time.Millisecond) // Wait for token to expire

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestRequireAuth_MalformedAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

// ---- Tests for RequireRole ----

func TestRequireRole_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/admin", middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	token := generateTestToken(jwtManager, 1, "admin")

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequireRole_WrongRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/admin", middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	token := generateTestToken(jwtManager, 2, "user")

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected status 403, got %d", w.Code)
	}
}

func TestRequireRole_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/admin", middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

// ---- Tests for RequireAnyRole ----

func TestRequireAnyRole_AllowsWhenAnyRoleMatches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roles := map[uint64]map[string]bool{
		1: {
			"admin": true,
		},
	}
	middleware := newMiddlewareWithServices(nil, roles, nil, nil)

	router := gin.New()
	router.GET("/resource", middleware.RequireAnyRole("admin", "editor"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := generateTestToken(middleware.jwtManager, 1, "user")
	req := httptest.NewRequest(http.MethodGet, "/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequireAnyRole_ForbiddenWhenNoRolesMatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roles := map[uint64]map[string]bool{
		2: {
			"user": true,
		},
	}
	middleware := newMiddlewareWithServices(nil, roles, nil, nil)

	router := gin.New()
	router.GET("/resource", middleware.RequireAnyRole("admin", "editor"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	token := generateTestToken(middleware.jwtManager, 2, "user")
	req := httptest.NewRequest(http.MethodGet, "/resource", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// ---- Tests for RequirePermission ----

func TestRequirePermission_NoUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	_ = newFakePermissionService()
	_ = newFakeRoleService()
	middleware := &PermissionMiddleware{
		permissionSvc: (*permissionservice.PermissionService)(nil),
		roleSvc:       (*roleservice.RoleService)(nil),
	}

	router := gin.New()
	router.GET("/test", middleware.RequirePermission(model.HTTPMethodGET, "/test"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No user_id in context
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", w.Code)
	}
}

func TestRequirePermission_AllowsWithPermission(t *testing.T) {
	gin.SetMode(gin.TestMode)

	perms := map[uint64][]model.Permission{
		10: {
			{Method: model.HTTPMethodGET, Path: "/secure"},
		},
	}
	middleware := newMiddlewareWithServices(perms, nil, nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(10))
		c.Next()
	})
	router.GET("/secure", middleware.RequirePermission(model.HTTPMethodGET, "/secure"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestRequirePermission_ForbiddenWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := newMiddlewareWithServices(nil, nil, nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(20))
		c.Next()
	})
	router.GET("/secure", middleware.RequirePermission(model.HTTPMethodGET, "/secure"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when permission missing, got %d", w.Code)
	}
}

func TestRequirePermission_SuperAdminBypass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roleMatrix := map[uint64]map[string]bool{
		77: {string(model.RoleSlugSuperAdmin): true},
	}
	middleware := newMiddlewareWithServices(nil, roleMatrix, nil, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(77))
		c.Next()
	})
	router.POST("/secure", middleware.RequirePermission(model.HTTPMethodPOST, "/secure"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	req := httptest.NewRequest(http.MethodPost, "/secure", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected super admin bypass, got %d", w.Code)
	}
}

func TestRequirePermission_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	middleware := newMiddlewareWithServices(nil, nil, errors.New("boom"), nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(UserIDKey, uint64(55))
		c.Next()
	})
	router.GET("/secure", middleware.RequirePermission(model.HTTPMethodGET, "/secure"), func(c *gin.Context) {})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when permission service errors, got %d", w.Code)
	}
}

// ---- Tests for Context Keys ----

func TestContextKeys(t *testing.T) {
	if UserIDKey != "user_id" {
		t.Errorf("Expected UserIDKey to be 'user_id', got '%s'", UserIDKey)
	}

	if UserRoleKey != "user_role" {
		t.Errorf("Expected UserRoleKey to be 'user_role', got '%s'", UserRoleKey)
	}

	if UserPermissionsKey != "user_permissions" {
		t.Errorf("Expected UserPermissionsKey to be 'user_permissions', got '%s'", UserPermissionsKey)
	}
}

// ---- Integration Tests ----

func TestPermissionMiddleware_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()

	// Public endpoint
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public"})
	})

	// Protected endpoint
	router.GET("/protected", middleware.RequireAuth(), func(c *gin.Context) {
		userID, _ := c.Get(UserIDKey)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Admin endpoint
	router.GET("/admin", middleware.RequireRole("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin"})
	})

	// Test public endpoint
	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Public endpoint failed: %d", w.Code)
	}

	// Test protected endpoint without auth
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Protected endpoint should return 401, got %d", w.Code)
	}

	// Test protected endpoint with auth
	token := generateTestToken(jwtManager, 123, "user")
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Protected endpoint with auth should return 200, got %d", w.Code)
	}

	// Test admin endpoint with user role
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("Admin endpoint with user role should return 403, got %d", w.Code)
	}

	// Test admin endpoint with admin role
	adminToken := generateTestToken(jwtManager, 1, "admin")
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Admin endpoint with admin role should return 200, got %d", w.Code)
	}
}

// ---- Security Tests ----

func TestSecurity_TokenReuse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	middleware := &PermissionMiddleware{
		jwtManager: jwtManager,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	token := generateTestToken(jwtManager, 123, "user")

	// First request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("First request failed: %d", w.Code)
	}

	// Second request with same token (should still work)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Token reuse should work: %d", w.Code)
	}
}

func TestSecurity_DifferentSecretKey(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager1 := auth.NewJWTManager("secret-1", 24*time.Hour)
	jwtManager2 := auth.NewJWTManager("secret-2", 24*time.Hour)

	middleware := &PermissionMiddleware{
		jwtManager: jwtManager2,
	}

	router := gin.New()
	router.GET("/test", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Generate token with different secret
	token := generateTestToken(jwtManager1, 123, "user")

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Token with different secret should be rejected, got %d", w.Code)
	}
}

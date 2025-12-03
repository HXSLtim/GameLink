package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gamelink/pkg/cache"
	"gamelink/internal/handler/middleware"
	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/internal/repository/user"
	permissionservice "gamelink/internal/service/admin"
	roleservice "gamelink/internal/service/admin"
	"gamelink/pkg/testutil"
)

// 权限中间件：超级管理员放行、具备权限放行、无权限拒绝
func TestPermissionMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)

	// seed roles
	superRole := &model.RoleModel{Slug: string(model.RoleSlugSuperAdmin), Name: "Super", IsSystem: true}
	adminRole := &model.RoleModel{Slug: string(model.RoleSlugAdmin), Name: "Admin", IsSystem: true}
	_ = roleRepo.Create(ctx(), superRole)
	_ = roleRepo.Create(ctx(), adminRole)

	// seed permission
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/secure",
		Code:        "secure:get",
		Group:       "/secure",
		Description: "secure endpoint",
	}
	_ = permRepo.Create(ctx(), perm)
	// grant to admin role
	_ = roleRepo.AssignPermissions(ctx(), adminRole.ID, []uint64{perm.ID})

	// seed users
	superUser := &model.User{Name: "Super", Email: "s@example.com", Phone: "19900000001", PasswordHash: "x", Role: model.RoleAdmin}
	adminUser := &model.User{Name: "Admin", Email: "a@example.com", Phone: "19900000002", PasswordHash: "x", Role: model.RoleAdmin}
	plainUser := &model.User{Name: "User", Email: "u@example.com", Phone: "19900000003", PasswordHash: "x", Role: model.RoleUser}
	_ = userRepo.Create(ctx(), superUser)
	_ = userRepo.Create(ctx(), adminUser)
	_ = userRepo.Create(ctx(), plainUser)

	// bind roles
	_ = roleRepo.AssignToUser(ctx(), superUser.ID, []uint64{superRole.ID})
	_ = roleRepo.AssignToUser(ctx(), adminUser.ID, []uint64{adminRole.ID})

	permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
	roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	securePath := "/api/v1/secure"
	handler := func(c *gin.Context) { c.String(http.StatusOK, "ok") }

	testCases := []struct {
		name   string
		userID uint64
		expect int
	}{
		{"super_admin_pass", superUser.ID, http.StatusOK},
		{"admin_with_perm_pass", adminUser.ID, http.StatusOK},
		{"no_perm_forbidden", plainUser.ID, http.StatusForbidden},
	}

	for _, tc := range testCases {
		router := gin.New()
		router.GET("/secure", setUserID(tc.userID), pm.RequirePermission(model.HTTPMethodGET, securePath), handler)
		req := httptest.NewRequest(http.MethodGet, "/secure", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != tc.expect {
			t.Fatalf("%s expected %d got %d, body=%s", tc.name, tc.expect, w.Code, w.Body.String())
		}
	}
}

func setUserID(id uint64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.UserIDKey, id)
		c.Next()
	}
}

func migratePermissionModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
	); err != nil {
		t.Fatalf("migrate permission models: %v", err)
	}
}

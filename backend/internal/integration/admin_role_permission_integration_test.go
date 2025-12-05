package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	adminhandler "gamelink/internal/handler/admin"
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

// 管理端角色/权限：创建角色、创建权限、分配权限给角色、分配角色给用户，并验证权限中间件放行/拦截
func TestAdminRolePermissionSync(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migrateRBACModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)

	// 种子：超级管理员 + 普通管理员 + 无权限用户
	superUser := &model.User{Name: "Super", Email: "s@example.com", Phone: "18800000001", PasswordHash: "x", Role: model.RoleAdmin, Status: model.UserStatusActive}
	adminUser := &model.User{Name: "Admin", Email: "a@example.com", Phone: "18800000002", PasswordHash: "x", Role: model.RoleAdmin, Status: model.UserStatusActive}
	plainUser := &model.User{Name: "Plain", Email: "p@example.com", Phone: "18800000003", PasswordHash: "x", Role: model.RoleAdmin, Status: model.UserStatusActive}
	_ = userRepo.Create(ctx(), superUser)
	_ = userRepo.Create(ctx(), adminUser)
	_ = userRepo.Create(ctx(), plainUser)

	// 系统内置超级管理员角色
	superRole := &model.RoleModel{Slug: string(model.RoleSlugSuperAdmin), Name: "Super Admin", IsSystem: true}
	_ = roleRepo.Create(ctx(), superRole)
	_ = roleRepo.AssignToUser(ctx(), superUser.ID, []uint64{superRole.ID})

	permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
	roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	roleHandler := adminhandler.NewRoleHandler(roleSvc)
	permHandler := adminhandler.NewPermissionHandler(permSvc)

	router := gin.New()
	admin := router.Group("/api/v1/admin")
	admin.Use(setUserID(superUser.ID)) // 超级管理员调用管理端 API

	// 角色/权限管理接口（与实际路由保持一致的权限配置）
	admin.POST("/permissions", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/permissions"), permHandler.CreatePermission)
	admin.POST("/roles", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles"), roleHandler.CreateRole)
	admin.PUT("/roles/:id/permissions", pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id/permissions"), roleHandler.AssignPermissions)
	admin.POST("/roles/assign-user", pm.RequirePermission(model.HTTPMethodPOST, "/api/v1/admin/roles/assign-user"), roleHandler.AssignRolesToUser)
	admin.GET("/users/:user_id/roles", pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/users/:id/roles"), roleHandler.GetUserRoles)

	// 1) 创建权限
	permPayload := map[string]interface{}{
		"method":      "GET",
		"path":        "/api/v1/admin/secure-data",
		"code":        "secure.data.read",
		"group":       "secure",
		"description": "secure data read",
	}
	permResp := doJSON(router, http.MethodPost, "/api/v1/admin/permissions", permPayload, "")
	if permResp.Code != http.StatusCreated {
		t.Fatalf("create permission status=%d body=%s", permResp.Code, permResp.Body.String())
	}
	var permParsed apiResp[adminhandler.Permission]
	_ = json.Unmarshal(permResp.Body.Bytes(), &permParsed)

	// 2) 创建角色
	rolePayload := map[string]interface{}{
		"slug": "auditor",
		"name": "Auditor",
	}
	roleResp := doJSON(router, http.MethodPost, "/api/v1/admin/roles", rolePayload, "")
	if roleResp.Code != http.StatusCreated {
		t.Fatalf("create role status=%d body=%s", roleResp.Code, roleResp.Body.String())
	}
	var roleParsed apiResp[model.RoleModel]
	_ = json.Unmarshal(roleResp.Body.Bytes(), &roleParsed)
	roleID := roleParsed.Data.ID

	// 3) 分配权限给角色
	assignPayload := map[string]interface{}{"permissionIds": []uint64{permParsed.Data.ID}}
	assignResp := doJSON(router, http.MethodPut, "/api/v1/admin/roles/"+uintToStr(roleID)+"/permissions", assignPayload, "")
	if assignResp.Code != http.StatusOK {
		t.Fatalf("assign perm status=%d body=%s", assignResp.Code, assignResp.Body.String())
	}

	// 4) 分配角色给普通管理员
	assignUserPayload := map[string]interface{}{
		"userId":  adminUser.ID,
		"roleIds": []uint64{roleID},
	}
	assignUserResp := doJSON(router, http.MethodPost, "/api/v1/admin/roles/assign-user", assignUserPayload, "")
	if assignUserResp.Code != http.StatusOK {
		t.Fatalf("assign role to user status=%d body=%s", assignUserResp.Code, assignUserResp.Body.String())
	}

	// 5) 查询用户角色应该包含新角色
	userRolesResp := doJSON(router, http.MethodGet, "/api/v1/admin/users/"+uintToStr(adminUser.ID)+"/roles", nil, "")
	if userRolesResp.Code != http.StatusOK {
		t.Fatalf("get user roles status=%d body=%s", userRolesResp.Code, userRolesResp.Body.String())
	}

	// 6) 验证权限中间件效果：具备权限放行，无权限拦截
	secure := router.Group("/api/v1/admin")
	secure.GET("/secure-data",
		setUserID(adminUser.ID),
		pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/secure-data"),
		func(c *gin.Context) { c.String(http.StatusOK, "ok") },
	)

	secure.GET("/secure-data/forbidden",
		setUserID(plainUser.ID),
		pm.RequirePermission(model.HTTPMethodGET, "/api/v1/admin/secure-data"),
		func(c *gin.Context) { c.String(http.StatusOK, "nope") },
	)

	okResp := doJSON(router, http.MethodGet, "/api/v1/admin/secure-data", nil, "")
	if okResp.Code != http.StatusOK {
		t.Fatalf("expected 200 for permitted user, got %d body=%s", okResp.Code, okResp.Body.String())
	}

	forbiddenResp := doJSON(router, http.MethodGet, "/api/v1/admin/secure-data/forbidden", nil, "")
	if forbiddenResp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for user without perm, got %d body=%s", forbiddenResp.Code, forbiddenResp.Body.String())
	}
}

func migrateRBACModels(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.AutoMigrate(
		&model.User{},
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.UserRole{},
	); err != nil {
		t.Fatalf("migrate rbac models: %v", err)
	}
}

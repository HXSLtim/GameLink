package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

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

// TestReviewPermissions 测试评价管理权限
func TestReviewPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)
	migratePermissionModels(t, db)

	userRepo := user.NewUserRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := permission.NewPermissionRepository(db)

	// 创建角色
	superRole := &model.RoleModel{Slug: string(model.RoleSlugSuperAdmin), Name: "Super Admin", IsSystem: true}
	reviewerRole := &model.RoleModel{Slug: "reviewer", Name: "Reviewer", IsSystem: false}
	_ = roleRepo.Create(ctx(), superRole)
	_ = roleRepo.Create(ctx(), reviewerRole)

	// 创建评价管理权限
	reviewViewPerm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/reviews",
		Code:        "review:view",
		Group:       "/admin/reviews",
		Description: "查看评价列表",
	}
	reviewApprovePerm := &model.Permission{
		Method:      model.HTTPMethodPUT,
		Path:        "/api/v1/admin/reviews/:id/approve",
		Code:        "review:approve",
		Group:       "/admin/reviews",
		Description: "批准评价",
	}
	reviewDeletePerm := &model.Permission{
		Method:      model.HTTPMethodDELETE,
		Path:        "/api/v1/admin/reviews/:id",
		Code:        "review:delete",
		Group:       "/admin/reviews",
		Description: "删除评价",
	}
	reviewManagePerm := &model.Permission{
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/admin/sensitive-words",
		Code:        "review:manage",
		Group:       "/admin/reviews",
		Description: "添加敏感词",
	}

	_ = permRepo.Create(ctx(), reviewViewPerm)
	_ = permRepo.Create(ctx(), reviewApprovePerm)
	_ = permRepo.Create(ctx(), reviewDeletePerm)
	_ = permRepo.Create(ctx(), reviewManagePerm)

	// 为 reviewer 角色分配 review:view 和 review:approve 权限
	_ = roleRepo.AssignPermissions(ctx(), reviewerRole.ID, []uint64{reviewViewPerm.ID, reviewApprovePerm.ID})

	// 创建用户
	superUser := &model.User{Name: "Super", Email: "super@example.com", Phone: "19900000001", PasswordHash: "x", Role: model.RoleAdmin}
	reviewerUser := &model.User{Name: "Reviewer", Email: "reviewer@example.com", Phone: "19900000002", PasswordHash: "x", Role: model.RoleAdmin}
	normalUser := &model.User{Name: "Normal", Email: "normal@example.com", Phone: "19900000003", PasswordHash: "x", Role: model.RoleUser}

	_ = userRepo.Create(ctx(), superUser)
	_ = userRepo.Create(ctx(), reviewerUser)
	_ = userRepo.Create(ctx(), normalUser)

	// 绑定角色
	_ = roleRepo.AssignToUser(ctx(), superUser.ID, []uint64{superRole.ID})
	_ = roleRepo.AssignToUser(ctx(), reviewerUser.ID, []uint64{reviewerRole.ID})

	permSvc := permissionservice.NewPermissionService(permRepo, cache.NewMemory())
	roleSvc := roleservice.NewRoleService(roleRepo, cache.NewMemory())
	pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

	handler := func(c *gin.Context) { c.String(http.StatusOK, "ok") }

	testCases := []struct {
		name       string
		userID     uint64
		method     model.HTTPMethod
		path       string
		routePath  string
		expectCode int
	}{
		// 超级管理员可以访问所有功能
		{"super_admin_view", superUser.ID, model.HTTPMethodGET, "/api/v1/admin/reviews", "/reviews", http.StatusOK},
		{"super_admin_approve", superUser.ID, model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/approve", "/reviews/:id/approve", http.StatusOK},
		{"super_admin_delete", superUser.ID, model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id", "/reviews/:id", http.StatusOK},
		{"super_admin_manage", superUser.ID, model.HTTPMethodPOST, "/api/v1/admin/sensitive-words", "/sensitive-words", http.StatusOK},

		// Reviewer 可以查看和批准，但不能删除和管理敏感词
		{"reviewer_view", reviewerUser.ID, model.HTTPMethodGET, "/api/v1/admin/reviews", "/reviews", http.StatusOK},
		{"reviewer_approve", reviewerUser.ID, model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/approve", "/reviews/:id/approve", http.StatusOK},
		{"reviewer_delete_forbidden", reviewerUser.ID, model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id", "/reviews/:id", http.StatusForbidden},
		{"reviewer_manage_forbidden", reviewerUser.ID, model.HTTPMethodPOST, "/api/v1/admin/sensitive-words", "/sensitive-words", http.StatusForbidden},

		// 普通用户无权访问任何评价管理功能
		{"normal_view_forbidden", normalUser.ID, model.HTTPMethodGET, "/api/v1/admin/reviews", "/reviews", http.StatusForbidden},
		{"normal_approve_forbidden", normalUser.ID, model.HTTPMethodPUT, "/api/v1/admin/reviews/:id/approve", "/reviews/:id/approve", http.StatusForbidden},
		{"normal_delete_forbidden", normalUser.ID, model.HTTPMethodDELETE, "/api/v1/admin/reviews/:id", "/reviews/:id", http.StatusForbidden},
		{"normal_manage_forbidden", normalUser.ID, model.HTTPMethodPOST, "/api/v1/admin/sensitive-words", "/sensitive-words", http.StatusForbidden},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			
			// 根据 HTTP 方法注册路由
			switch tc.method {
			case model.HTTPMethodGET:
				router.GET(tc.routePath, setUserID(tc.userID), pm.RequirePermission(tc.method, tc.path), handler)
			case model.HTTPMethodPOST:
				router.POST(tc.routePath, setUserID(tc.userID), pm.RequirePermission(tc.method, tc.path), handler)
			case model.HTTPMethodPUT:
				router.PUT(tc.routePath, setUserID(tc.userID), pm.RequirePermission(tc.method, tc.path), handler)
			case model.HTTPMethodDELETE:
				router.DELETE(tc.routePath, setUserID(tc.userID), pm.RequirePermission(tc.method, tc.path), handler)
			}

			req := httptest.NewRequest(string(tc.method), tc.routePath, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectCode {
				t.Errorf("%s: expected status %d, got %d, body=%s", tc.name, tc.expectCode, w.Code, w.Body.String())
			}
		})
	}
}

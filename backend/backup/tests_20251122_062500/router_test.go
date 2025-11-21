package router

import (
	"context"
	"errors"
	"testing"
	"time"

	"gamelink/internal/config"
	"gamelink/internal/model"
	permissionservice "gamelink/internal/service/permission"
	roleservice "gamelink/internal/service/role"

	"github.com/gin-gonic/gin"
)

// TestResolveGinMode 测试 resolveGinMode 函数的各种场景
func TestResolveGinMode(t *testing.T) {
	tests := []struct {
		name     string
		ginMode  string
		appEnv   string
		expected string
	}{
		{
			name:     "默认调试模式",
			ginMode:  "",
			appEnv:   "",
			expected: gin.DebugMode,
		},
		{
			name:     "生产环境",
			appEnv:   "production",
			expected: gin.ReleaseMode,
		},
		{
			name:     "GIN_MODE 环境变量覆盖",
			ginMode:  "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理环境变量
			t.Setenv("GIN_MODE", tt.ginMode)
			t.Setenv("APP_ENV", tt.appEnv)

			if got := resolveGinMode(); got != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

// dummyCache 是一个用于测试的缓存实现
type dummyCache struct{}

func (dummyCache) Get(context.Context, string) (string, bool, error)        { return "", false, nil }
func (dummyCache) Set(context.Context, string, string, time.Duration) error { return nil }
func (dummyCache) Delete(context.Context, string) error                     { return nil }
func (dummyCache) Close(context.Context) error                              { return nil }

// fakePermRepo 是一个用于测试的权限仓库实现
type fakePermRepo struct{ perms []model.Permission }

func (f *fakePermRepo) List(context.Context) ([]model.Permission, error) { return f.perms, nil }
func (f *fakePermRepo) ListPaged(context.Context, int, int) ([]model.Permission, int64, error) {
	return nil, 0, nil
}
func (f *fakePermRepo) ListPagedWithFilter(context.Context, int, int, string, string, string) ([]model.Permission, int64, error) {
	return nil, 0, nil
}
func (f *fakePermRepo) ListByGroup(context.Context) (map[string][]model.Permission, error) {
	return map[string][]model.Permission{}, nil
}
func (f *fakePermRepo) ListGroups(context.Context) ([]string, error) { return nil, nil }
func (f *fakePermRepo) Get(context.Context, uint64) (*model.Permission, error) {
	return nil, errors.New("not found")
}
func (f *fakePermRepo) GetByResource(context.Context, string, string) (*model.Permission, error) {
	return nil, errors.New("not found")
}
func (f *fakePermRepo) GetByCode(context.Context, string) (*model.Permission, error) {
	return nil, errors.New("not found")
}
func (f *fakePermRepo) GetByMethodAndPath(context.Context, string, string) (*model.Permission, error) {
	return nil, errors.New("not found")
}
func (f *fakePermRepo) Create(context.Context, *model.Permission) error             { return nil }
func (f *fakePermRepo) Update(context.Context, *model.Permission) error             { return nil }
func (f *fakePermRepo) UpsertByMethodPath(context.Context, *model.Permission) error { return nil }
func (f *fakePermRepo) Delete(context.Context, uint64) error                        { return nil }
func (f *fakePermRepo) ListByRoleID(context.Context, uint64) ([]model.Permission, error) {
	return nil, nil
}
func (f *fakePermRepo) ListByUserID(context.Context, uint64) ([]model.Permission, error) {
	return nil, nil
}

// fakeRoleRepo 是一个用于测试的角色仓库实现
type fakeRoleRepo struct {
	bySlug  map[string]*model.RoleModel
	assigns map[uint64][]uint64
}

func (f *fakeRoleRepo) List(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) ListPaged(context.Context, int, int) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRoleRepo) ListPagedWithFilter(context.Context, int, int, string, *bool) ([]model.RoleModel, int64, error) {
	return nil, 0, nil
}
func (f *fakeRoleRepo) ListWithPermissions(context.Context) ([]model.RoleModel, error) {
	return nil, nil
}
func (f *fakeRoleRepo) Get(context.Context, uint64) (*model.RoleModel, error) {
	return nil, errors.New("not found")
}
func (f *fakeRoleRepo) GetWithPermissions(context.Context, uint64) (*model.RoleModel, error) {
	return nil, errors.New("not found")
}
func (f *fakeRoleRepo) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
	r := f.bySlug[slug]
	if r == nil {
		return nil, errors.New("not found")
	}
	return r, nil
}
func (f *fakeRoleRepo) Create(context.Context, *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Update(context.Context, *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Delete(context.Context, uint64) error           { return nil }
func (f *fakeRoleRepo) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	if f.assigns == nil {
		f.assigns = map[uint64][]uint64{}
	}
	f.assigns[roleID] = append([]uint64{}, permissionIDs...)
	return nil
}
func (f *fakeRoleRepo) AddPermissions(context.Context, uint64, []uint64) error    { return nil }
func (f *fakeRoleRepo) RemovePermissions(context.Context, uint64, []uint64) error { return nil }
func (f *fakeRoleRepo) AssignToUser(context.Context, uint64, []uint64) error      { return nil }
func (f *fakeRoleRepo) RemoveFromUser(context.Context, uint64, []uint64) error    { return nil }
func (f *fakeRoleRepo) ListByUserID(context.Context, uint64) ([]model.RoleModel, error) {
	return nil, nil
}
func (f *fakeRoleRepo) CheckUserHasRole(context.Context, uint64, string) (bool, error) {
	return false, nil
}

// TestAssignDefaultRolePermissions_AllRoles 测试为所有默认角色分配权限
func TestAssignDefaultRolePermissions_AllRoles(t *testing.T) {
	permRepo := &fakePermRepo{perms: []model.Permission{{Base: model.Base{ID: 11}}, {Base: model.Base{ID: 22}}}}
	roleRepo := &fakeRoleRepo{bySlug: map[string]*model.RoleModel{
		string(model.RoleSlugSuperAdmin): {Base: model.Base{ID: 1}, Slug: string(model.RoleSlugSuperAdmin)},
		string(model.RoleSlugAdmin):      {Base: model.Base{ID: 2}, Slug: string(model.RoleSlugAdmin)},
	}}
	permSvc := permissionservice.NewPermissionService(permRepo, dummyCache{})
	roleSvc := roleservice.NewRoleService(roleRepo, dummyCache{})
	err := assignDefaultRolePermissions(context.Background(), roleSvc, permSvc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roleRepo.assigns[1]) != 2 || len(roleRepo.assigns[2]) != 2 {
		t.Fatalf("expected 2 permissions assigned to both roles")
	}
}

// TestAssignDefaultRolePermissions_MissingSuperAdmin 测试当超级管理员角色不存在的情况
func TestAssignDefaultRolePermissions_MissingSuperAdmin(t *testing.T) {
	permRepo := &fakePermRepo{perms: []model.Permission{{Base: model.Base{ID: 99}}}}
	roleRepo := &fakeRoleRepo{bySlug: map[string]*model.RoleModel{
		string(model.RoleSlugAdmin): {Base: model.Base{ID: 2}, Slug: string(model.RoleSlugAdmin)},
	}}
	permSvc := permissionservice.NewPermissionService(permRepo, dummyCache{})
	roleSvc := roleservice.NewRoleService(roleRepo, dummyCache{})
	err := assignDefaultRolePermissions(context.Background(), roleSvc, permSvc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roleRepo.assigns[2]) != 1 {
		t.Fatalf("expected admin to receive permissions")
	}
	if _, ok := roleRepo.assigns[1]; ok {
		t.Fatalf("unexpected assignment to super_admin")
	}
}

// TestSyncAPIPermissions_DisabledInProduction 测试生产环境且未显式开启时不同步权限（覆盖条件分支）。
func TestSyncAPIPermissions_DisabledInProduction(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("SYNC_API_PERMISSIONS", "")

	r := &Router{}
	permSvc := permissionservice.NewPermissionService(&fakePermRepo{}, dummyCache{})
	roleSvc := roleservice.NewRoleService(&fakeRoleRepo{}, dummyCache{})

	// 条件不满足时不应尝试访问 engine，调用不应 panic。
	r.syncAPIPermissions(permSvc, roleSvc)
}

// TestRegisterSwaggerRoutes 根据配置开关保证 swagger 路由按预期注册/禁用。
func TestRegisterSwaggerRoutes_EnableAndDisable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 关闭 swagger 时不应注册 /swagger/*any 路由
	{
		r := &Router{
			engine: gin.New(),
			cfg:    config.AppConfig{EnableSwagger: false},
		}
		r.registerSwaggerRoutes()
		for _, rt := range r.engine.Routes() {
			if rt.Path == "/swagger/*any" {
				t.Fatalf("swagger 被配置为禁用时不应注册 /swagger/*any 路由")
			}
		}
	}

	// 开启 swagger 时应注册 /swagger/*any 路由
	{
		r := &Router{
			engine: gin.New(),
			cfg:    config.AppConfig{EnableSwagger: true},
		}
		r.registerSwaggerRoutes()
		found := false
		for _, rt := range r.engine.Routes() {
			if rt.Path == "/swagger/*any" {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("启用 swagger 时应注册 /swagger/*any 路由")
		}
	}
}

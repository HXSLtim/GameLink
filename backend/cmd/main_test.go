package main

import (
    "testing"

    "github.com/gin-gonic/gin"
    "context"
    "errors"
    "time"
    "gamelink/internal/model"
    permissionservice "gamelink/internal/service/permission"
    roleservice "gamelink/internal/service/role"
)

func TestResolveGinMode(t *testing.T) {
    t.Setenv("GIN_MODE", "")
    t.Setenv("APP_ENV", "")
    if got := resolveGinMode(); got != gin.DebugMode {
        t.Fatalf("expected debug mode, got %s", got)
    }

	t.Setenv("APP_ENV", "production")
	if got := resolveGinMode(); got != gin.ReleaseMode {
		t.Fatalf("expected release mode, got %s", got)
	}

	t.Setenv("GIN_MODE", "test")
    if got := resolveGinMode(); got != "test" {
        t.Fatalf("expected env override to 'test', got %s", got)
    }
}

type dummyCache struct{}
func (dummyCache) Get(context.Context, string) (string, bool, error) { return "", false, nil }
func (dummyCache) Set(context.Context, string, string, time.Duration) error { return nil }
func (dummyCache) Delete(context.Context, string) error { return nil }
func (dummyCache) Close(context.Context) error { return nil }

type fakePermRepo struct{ perms []model.Permission }
func (f *fakePermRepo) List(context.Context) ([]model.Permission, error) { return f.perms, nil }
func (f *fakePermRepo) ListPaged(context.Context, int, int) ([]model.Permission, int64, error) { return nil, 0, nil }
func (f *fakePermRepo) ListPagedWithFilter(context.Context, int, int, string, string, string) ([]model.Permission, int64, error) { return nil, 0, nil }
func (f *fakePermRepo) ListByGroup(context.Context) (map[string][]model.Permission, error) { return map[string][]model.Permission{}, nil }
func (f *fakePermRepo) ListGroups(context.Context) ([]string, error) { return nil, nil }
func (f *fakePermRepo) Get(context.Context, uint64) (*model.Permission, error) { return nil, errors.New("not found") }
func (f *fakePermRepo) GetByResource(context.Context, string, string) (*model.Permission, error) { return nil, errors.New("not found") }
func (f *fakePermRepo) GetByCode(context.Context, string) (*model.Permission, error) { return nil, errors.New("not found") }
func (f *fakePermRepo) GetByMethodAndPath(context.Context, string, string) (*model.Permission, error) { return nil, errors.New("not found") }
func (f *fakePermRepo) Create(context.Context, *model.Permission) error { return nil }
func (f *fakePermRepo) Update(context.Context, *model.Permission) error { return nil }
func (f *fakePermRepo) UpsertByMethodPath(context.Context, *model.Permission) error { return nil }
func (f *fakePermRepo) Delete(context.Context, uint64) error { return nil }
func (f *fakePermRepo) ListByRoleID(context.Context, uint64) ([]model.Permission, error) { return nil, nil }
func (f *fakePermRepo) ListByUserID(context.Context, uint64) ([]model.Permission, error) { return nil, nil }

type fakeRoleRepo struct{
    bySlug map[string]*model.RoleModel
    assigns map[uint64][]uint64
}
func (f *fakeRoleRepo) List(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) ListPaged(context.Context, int, int) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (f *fakeRoleRepo) ListPagedWithFilter(context.Context, int, int, string, *bool) ([]model.RoleModel, int64, error) { return nil, 0, nil }
func (f *fakeRoleRepo) ListWithPermissions(context.Context) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) Get(context.Context, uint64) (*model.RoleModel, error) { return nil, errors.New("not found") }
func (f *fakeRoleRepo) GetWithPermissions(context.Context, uint64) (*model.RoleModel, error) { return nil, errors.New("not found") }
func (f *fakeRoleRepo) GetBySlug(ctx context.Context, slug string) (*model.RoleModel, error) {
    r := f.bySlug[slug]
    if r == nil { return nil, errors.New("not found") }
    return r, nil
}
func (f *fakeRoleRepo) Create(context.Context, *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Update(context.Context, *model.RoleModel) error { return nil }
func (f *fakeRoleRepo) Delete(context.Context, uint64) error { return nil }
func (f *fakeRoleRepo) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
    if f.assigns == nil { f.assigns = map[uint64][]uint64{} }
    f.assigns[roleID] = append([]uint64{}, permissionIDs...)
    return nil
}
func (f *fakeRoleRepo) AddPermissions(context.Context, uint64, []uint64) error { return nil }
func (f *fakeRoleRepo) RemovePermissions(context.Context, uint64, []uint64) error { return nil }
func (f *fakeRoleRepo) AssignToUser(context.Context, uint64, []uint64) error { return nil }
func (f *fakeRoleRepo) RemoveFromUser(context.Context, uint64, []uint64) error { return nil }
func (f *fakeRoleRepo) ListByUserID(context.Context, uint64) ([]model.RoleModel, error) { return nil, nil }
func (f *fakeRoleRepo) CheckUserHasRole(context.Context, uint64, string) (bool, error) { return false, nil }

func TestAssignDefaultRolePermissions_AllRoles(t *testing.T) {
    permRepo := &fakePermRepo{perms: []model.Permission{{Base: model.Base{ID: 11}}, {Base: model.Base{ID: 22}}}}
    roleRepo := &fakeRoleRepo{bySlug: map[string]*model.RoleModel{
        string(model.RoleSlugSuperAdmin): {Base: model.Base{ID: 1}, Slug: string(model.RoleSlugSuperAdmin)},
        string(model.RoleSlugAdmin):      {Base: model.Base{ID: 2}, Slug: string(model.RoleSlugAdmin)},
    }}
    permSvc := permissionservice.NewPermissionService(permRepo, dummyCache{})
    roleSvc := roleservice.NewRoleService(roleRepo, dummyCache{})
    err := assignDefaultRolePermissions(context.Background(), roleSvc, permSvc)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if len(roleRepo.assigns[1]) != 2 || len(roleRepo.assigns[2]) != 2 { t.Fatalf("expected 2 permissions assigned to both roles") }
}

func TestAssignDefaultRolePermissions_MissingSuperAdmin(t *testing.T) {
    permRepo := &fakePermRepo{perms: []model.Permission{{Base: model.Base{ID: 99}}}}
    roleRepo := &fakeRoleRepo{bySlug: map[string]*model.RoleModel{
        string(model.RoleSlugAdmin): {Base: model.Base{ID: 2}, Slug: string(model.RoleSlugAdmin)},
    }}
    permSvc := permissionservice.NewPermissionService(permRepo, dummyCache{})
    roleSvc := roleservice.NewRoleService(roleRepo, dummyCache{})
    err := assignDefaultRolePermissions(context.Background(), roleSvc, permSvc)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if len(roleRepo.assigns[2]) != 1 { t.Fatalf("expected admin to receive permissions") }
    if _, ok := roleRepo.assigns[1]; ok { t.Fatalf("unexpected assignment to super_admin") }
}

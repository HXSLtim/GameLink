package role

import (
    "context"
    "testing"

    "github.com/golang/mock/gomock"

    "gamelink/internal/cache"
    "gamelink/internal/model"
    "gamelink/internal/repository"
    "gamelink/internal/repository/mocks"
)

func TestAssignPermissionsToRole_InvalidatesCaches(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    _ = c.Set(ctx, "rbac:role_permissions:7", "x", 0)
    repo.EXPECT().AssignPermissions(ctx, uint64(7), []uint64{1, 2}).Return(nil)
    if err := svc.AssignPermissionsToRole(ctx, 7, []uint64{1, 2}); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
    if _, ok, _ := c.Get(ctx, "rbac:role_permissions:7"); ok { t.Fatalf("role permissions cache not invalidated") }
}

func TestAssignRolesToUser_InvalidatesCaches(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles:user:99", "x", 0)
    _ = c.Set(ctx, "rbac:user_permissions:99", "x", 0)
    repo.EXPECT().AssignToUser(ctx, uint64(99), []uint64{2}).Return(nil)
    if err := svc.AssignRolesToUser(ctx, 99, []uint64{2}); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles:user:99"); ok { t.Fatalf("user roles cache not invalidated") }
    if _, ok, _ := c.Get(ctx, "rbac:user_permissions:99"); ok { t.Fatalf("user permissions cache not invalidated") }
}

func TestAddPermissionsToRole_InvalidatesCaches(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    _ = c.Set(ctx, "rbac:role_permissions:8", "x", 0)
    repo.EXPECT().AddPermissions(ctx, uint64(8), []uint64{3}).Return(nil)
    if err := svc.AddPermissionsToRole(ctx, 8, []uint64{3}); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
    if _, ok, _ := c.Get(ctx, "rbac:role_permissions:8"); ok { t.Fatalf("role permissions cache not invalidated") }
}

func TestRemovePermissionsFromRole_InvalidatesCaches(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    _ = c.Set(ctx, "rbac:role_permissions:9", "x", 0)
    repo.EXPECT().RemovePermissions(ctx, uint64(9), []uint64{4}).Return(nil)
    if err := svc.RemovePermissionsFromRole(ctx, 9, []uint64{4}); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
    if _, ok, _ := c.Get(ctx, "rbac:role_permissions:9"); ok { t.Fatalf("role permissions cache not invalidated") }
}

func TestCreateRole_InvalidatesRolesCache(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    r := &model.RoleModel{Slug: "r1", Name: "r1"}
    repo.EXPECT().GetBySlug(ctx, "r1").Return(nil, repository.ErrNotFound)
    repo.EXPECT().Create(ctx, r).Return(nil)
    if err := svc.CreateRole(ctx, r); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
}

func TestUpdateRole_InvalidatesRolesCache(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    r := &model.RoleModel{Base: model.Base{ID: 10}, Slug: "r10", Name: "n"}
    repo.EXPECT().Get(ctx, uint64(10)).Return(&model.RoleModel{Base: model.Base{ID: 10}, Slug: "r10", Name: "n"}, nil)
    repo.EXPECT().Update(ctx, r).Return(nil)
    if err := svc.UpdateRole(ctx, r); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
}

func TestDeleteRole_InvalidatesRolesCache(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    repo := mocks.NewMockRoleRepository(ctrl)
    c := cache.NewMemory()
    svc := NewRoleService(repo, c)
    ctx := context.Background()
    _ = c.Set(ctx, "admin:roles", "x", 0)
    repo.EXPECT().Delete(ctx, uint64(11)).Return(nil)
    if err := svc.DeleteRole(ctx, 11); err != nil { t.Fatalf("%v", err) }
    if _, ok, _ := c.Get(ctx, "admin:roles"); ok { t.Fatalf("roles cache not invalidated") }
}

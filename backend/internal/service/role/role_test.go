package role

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"gamelink/internal/cache"
	"gamelink/internal/model"
	"gamelink/internal/repository"
	"gamelink/internal/repository/mocks"
)

// TestNewRoleService 测试构造函数。
func TestNewRoleService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()

	svc := NewRoleService(repo, mockCache)

	if svc == nil {
		t.Fatal("NewRoleService returned nil")
	}

	if svc.roles != repo {
		t.Error("roles repository not set correctly")
	}

	if svc.cache != mockCache {
		t.Error("cache not set correctly")
	}
}

// TestListRoles 测试获取所有角色。
func TestListRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功获取角色列表", func(t *testing.T) {
		expectedRoles := []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
			{Base: model.Base{ID: 2}, Slug: "user", Name: "用户"},
		}

		repo.EXPECT().
			List(ctx).
			Return(expectedRoles, nil)

		roles, err := svc.ListRoles(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles) != 2 {
			t.Errorf("Expected 2 roles, got %d", len(roles))
		}
	})

	t.Run("数据库错误", func(t *testing.T) {
		repo.EXPECT().
			List(ctx).
			Return(nil, errors.New("database error"))

		roles, err := svc.ListRoles(ctx)
		if err == nil {
			t.Error("Expected error for database failure")
		}

		if roles != nil {
			t.Error("Expected nil roles")
		}
	})
}

// TestListRolesPaged 测试分页获取角色。
func TestListRolesPaged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功分页获取", func(t *testing.T) {
		expectedRoles := []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		}

		repo.EXPECT().
			ListPaged(ctx, 1, 20).
			Return(expectedRoles, int64(1), nil)

		roles, total, err := svc.ListRolesPaged(ctx, 1, 20)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}

		if total != 1 {
			t.Errorf("Expected total 1, got %d", total)
		}
	})

	t.Run("自动修正无效页码", func(t *testing.T) {
		repo.EXPECT().
			ListPaged(ctx, 1, 20).
			Return([]model.RoleModel{}, int64(0), nil)

		_, _, err := svc.ListRolesPaged(ctx, 0, 20)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("自动修正无效页大小", func(t *testing.T) {
		repo.EXPECT().
			ListPaged(ctx, 1, 20).
			Return([]model.RoleModel{}, int64(0), nil)

		_, _, err := svc.ListRolesPaged(ctx, 1, 0)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("限制最大页大小", func(t *testing.T) {
		repo.EXPECT().
			ListPaged(ctx, 1, 20).
			Return([]model.RoleModel{}, int64(0), nil)

		_, _, err := svc.ListRolesPaged(ctx, 1, 200)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

// TestListRolesWithPermissions 测试获取角色（带权限）。
func TestListRolesWithPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功获取角色列表（带权限）", func(t *testing.T) {
		expectedRoles := []model.RoleModel{
			{
				Base: model.Base{ID: 1},
				Slug: "admin",
				Name: "管理员",
				Permissions: []model.Permission{
					{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/test"},
				},
			},
		}

		repo.EXPECT().
			ListWithPermissions(ctx).
			Return(expectedRoles, nil)

		roles, err := svc.ListRolesWithPermissions(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}

		if len(roles[0].Permissions) != 1 {
			t.Errorf("Expected 1 permission, got %d", len(roles[0].Permissions))
		}
	})
}

// TestGetRole 测试根据ID获取角色。
func TestGetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功获取角色", func(t *testing.T) {
		expectedRole := &model.RoleModel{
			Base: model.Base{ID: 1},
			Slug: "admin",
			Name: "管理员",
		}

		repo.EXPECT().
			Get(ctx, uint64(1)).
			Return(expectedRole, nil)

		role, err := svc.GetRole(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if role.ID != 1 {
			t.Errorf("Expected role ID 1, got %d", role.ID)
		}
	})

	t.Run("角色不存在", func(t *testing.T) {
		repo.EXPECT().
			Get(ctx, uint64(999)).
			Return(nil, repository.ErrNotFound)

		role, err := svc.GetRole(ctx, 999)
		if err != repository.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}

		if role != nil {
			t.Error("Expected nil role")
		}
	})
}

// TestGetRoleWithPermissions 测试根据ID获取角色（带权限）。
func TestGetRoleWithPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功获取角色（带权限）", func(t *testing.T) {
		expectedRole := &model.RoleModel{
			Base: model.Base{ID: 1},
			Slug: "admin",
			Name: "管理员",
			Permissions: []model.Permission{
				{Base: model.Base{ID: 1}, Method: model.HTTPMethodGET, Path: "/api/test"},
			},
		}

		repo.EXPECT().
			GetWithPermissions(ctx, uint64(1)).
			Return(expectedRole, nil)

		role, err := svc.GetRoleWithPermissions(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(role.Permissions) != 1 {
			t.Errorf("Expected 1 permission, got %d", len(role.Permissions))
		}
	})
}

// TestGetRoleBySlug 测试根据Slug获取角色。
func TestGetRoleBySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功根据Slug获取角色", func(t *testing.T) {
		expectedRole := &model.RoleModel{
			Base: model.Base{ID: 1},
			Slug: "admin",
			Name: "管理员",
		}

		repo.EXPECT().
			GetBySlug(ctx, "admin").
			Return(expectedRole, nil)

		role, err := svc.GetRoleBySlug(ctx, "admin")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if role.Slug != "admin" {
			t.Errorf("Expected slug 'admin', got '%s'", role.Slug)
		}
	})

	t.Run("Slug不存在", func(t *testing.T) {
		repo.EXPECT().
			GetBySlug(ctx, "nonexistent").
			Return(nil, repository.ErrNotFound)

		role, err := svc.GetRoleBySlug(ctx, "nonexistent")
		if err != repository.ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}

		if role != nil {
			t.Error("Expected nil role")
		}
	})
}

// TestCreateRole 测试创建角色。
func TestCreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功创建角色", func(t *testing.T) {
		newRole := &model.RoleModel{
			Slug: "new-role",
			Name: "新角色",
		}

		repo.EXPECT().
			GetBySlug(ctx, "new-role").
			Return(nil, repository.ErrNotFound)

		repo.EXPECT().
			Create(ctx, newRole).
			Return(nil)

		err := svc.CreateRole(ctx, newRole)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Slug为空", func(t *testing.T) {
		newRole := &model.RoleModel{
			Name: "测试角色",
		}

		err := svc.CreateRole(ctx, newRole)
		if err == nil {
			t.Error("Expected error for empty slug")
		}
	})

	t.Run("Name为空", func(t *testing.T) {
		newRole := &model.RoleModel{
			Slug: "test-role",
		}

		err := svc.CreateRole(ctx, newRole)
		if err == nil {
			t.Error("Expected error for empty name")
		}
	})

	t.Run("Slug已存在", func(t *testing.T) {
		newRole := &model.RoleModel{
			Slug: "existing-role",
			Name: "现有角色",
		}

		existingRole := &model.RoleModel{
			Base: model.Base{ID: 1},
			Slug: "existing-role",
			Name: "现有角色",
		}

		repo.EXPECT().
			GetBySlug(ctx, "existing-role").
			Return(existingRole, nil)

		err := svc.CreateRole(ctx, newRole)
		if err == nil {
			t.Error("Expected error for existing slug")
		}
	})
}

// TestUpdateRole 测试更新角色。
func TestUpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功更新普通角色", func(t *testing.T) {
		role := &model.RoleModel{
			Base:        model.Base{ID: 1},
			Slug:        "test-role",
			Name:        "更新后的角色",
			Description: "更新后的描述",
		}

		existingRole := &model.RoleModel{
			Base:     model.Base{ID: 1},
			Slug:     "test-role",
			Name:     "测试角色",
			IsSystem: false,
		}

		repo.EXPECT().
			Get(ctx, uint64(1)).
			Return(existingRole, nil)

		repo.EXPECT().
			Update(ctx, role).
			Return(nil)

		err := svc.UpdateRole(ctx, role)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("更新系统角色（仅允许更新描述）", func(t *testing.T) {
		role := &model.RoleModel{
			Base:        model.Base{ID: 1},
			Slug:        "super-admin",
			Name:        "新名称",
			Description: "新描述",
		}

		existingRole := &model.RoleModel{
			Base:        model.Base{ID: 1},
			Slug:        "super-admin",
			Name:        "超级管理员",
			Description: "旧描述",
			IsSystem:    true,
		}

		repo.EXPECT().
			Get(ctx, uint64(1)).
			Return(existingRole, nil)

		repo.EXPECT().
			Update(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, r *model.RoleModel) error {
				// 验证系统角色只更新了描述
				if r.Name != "超级管理员" {
					t.Error("System role name should not be changed")
				}
				if r.Description != "新描述" {
					t.Error("System role description should be updated")
				}
				return nil
			})

		err := svc.UpdateRole(ctx, role)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("ID为空", func(t *testing.T) {
		role := &model.RoleModel{
			Slug: "test-role",
			Name: "测试角色",
		}

		err := svc.UpdateRole(ctx, role)
		if err == nil {
			t.Error("Expected error for missing ID")
		}
	})

	t.Run("角色不存在", func(t *testing.T) {
		role := &model.RoleModel{
			Base: model.Base{ID: 999},
			Slug: "test-role",
			Name: "测试角色",
		}

		repo.EXPECT().
			Get(ctx, uint64(999)).
			Return(nil, repository.ErrNotFound)

		err := svc.UpdateRole(ctx, role)
		if err == nil {
			t.Error("Expected error for role not found")
		}
	})
}

// TestDeleteRole 测试删除角色。
func TestDeleteRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功删除角色", func(t *testing.T) {
		repo.EXPECT().
			Delete(ctx, uint64(1)).
			Return(nil)

		err := svc.DeleteRole(ctx, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("删除不存在的角色", func(t *testing.T) {
		repo.EXPECT().
			Delete(ctx, uint64(999)).
			Return(repository.ErrNotFound)

		err := svc.DeleteRole(ctx, 999)
		if err == nil {
			t.Error("Expected error for role not found")
		}
	})
}

// TestAssignPermissionsToRole 测试分配权限到角色。
func TestAssignPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功分配权限", func(t *testing.T) {
		roleID := uint64(1)
		permissionIDs := []uint64{1, 2, 3}

		repo.EXPECT().
			AssignPermissions(ctx, roleID, permissionIDs).
			Return(nil)

		err := svc.AssignPermissionsToRole(ctx, roleID, permissionIDs)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("分配权限失败", func(t *testing.T) {
		roleID := uint64(999)
		permissionIDs := []uint64{1, 2, 3}

		repo.EXPECT().
			AssignPermissions(ctx, roleID, permissionIDs).
			Return(errors.New("role not found"))

		err := svc.AssignPermissionsToRole(ctx, roleID, permissionIDs)
		if err == nil {
			t.Error("Expected error for role not found")
		}
	})
}

// TestAddPermissionsToRole 测试添加权限到角色。
func TestAddPermissionsToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功添加权限", func(t *testing.T) {
		roleID := uint64(1)
		permissionIDs := []uint64{4, 5}

		repo.EXPECT().
			AddPermissions(ctx, roleID, permissionIDs).
			Return(nil)

		err := svc.AddPermissionsToRole(ctx, roleID, permissionIDs)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// TestRemovePermissionsFromRole 测试移除角色的权限。
func TestRemovePermissionsFromRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功移除权限", func(t *testing.T) {
		roleID := uint64(1)
		permissionIDs := []uint64{1, 2}

		repo.EXPECT().
			RemovePermissions(ctx, roleID, permissionIDs).
			Return(nil)

		err := svc.RemovePermissionsFromRole(ctx, roleID, permissionIDs)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// TestListRolesByUserID 测试获取用户的角色。
func TestListRolesByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功获取用户角色（从数据库）", func(t *testing.T) {
		userID := uint64(123)
		expectedRoles := []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "user", Name: "用户"},
		}

		repo.EXPECT().
			ListByUserID(ctx, userID).
			Return(expectedRoles, nil)

		roles, err := svc.ListRolesByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
	})

	t.Run("成功获取用户角色（从缓存）", func(t *testing.T) {
		userID := uint64(456)

		// 第一次调用，从数据库获取并缓存
		expectedRoles := []model.RoleModel{
			{Base: model.Base{ID: 1}, Slug: "admin", Name: "管理员"},
		}

		repo.EXPECT().
			ListByUserID(ctx, userID).
			Return(expectedRoles, nil).
			Times(1)

		// 第一次调用
		roles1, err := svc.ListRolesByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// 第二次调用应该从缓存获取，不会调用repository
		roles2, err := svc.ListRolesByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles1) != len(roles2) {
			t.Error("Cached roles should match original roles")
		}
	})

	t.Run("用户没有角色", func(t *testing.T) {
		userID := uint64(789)

		repo.EXPECT().
			ListByUserID(ctx, userID).
			Return([]model.RoleModel{}, nil)

		roles, err := svc.ListRolesByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(roles) != 0 {
			t.Errorf("Expected 0 roles, got %d", len(roles))
		}
	})
}

// TestAssignRolesToUser 测试分配角色给用户。
func TestAssignRolesToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功分配角色", func(t *testing.T) {
		userID := uint64(123)
		roleIDs := []uint64{1, 2}

		repo.EXPECT().
			AssignToUser(ctx, userID, roleIDs).
			Return(nil)

		err := svc.AssignRolesToUser(ctx, userID, roleIDs)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("分配失败", func(t *testing.T) {
		userID := uint64(999)
		roleIDs := []uint64{1, 2}

		repo.EXPECT().
			AssignToUser(ctx, userID, roleIDs).
			Return(errors.New("user not found"))

		err := svc.AssignRolesToUser(ctx, userID, roleIDs)
		if err == nil {
			t.Error("Expected error for user not found")
		}
	})
}

// TestRemoveRolesFromUser 测试移除用户的角色。
func TestRemoveRolesFromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("成功移除角色", func(t *testing.T) {
		userID := uint64(123)
		roleIDs := []uint64{1}

		repo.EXPECT().
			RemoveFromUser(ctx, userID, roleIDs).
			Return(nil)

		err := svc.RemoveRolesFromUser(ctx, userID, roleIDs)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})
}

// TestCheckUserHasRole 测试检查用户是否有角色。
func TestCheckUserHasRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("用户拥有角色", func(t *testing.T) {
		userID := uint64(123)
		roleSlug := "admin"

		repo.EXPECT().
			CheckUserHasRole(ctx, userID, roleSlug).
			Return(true, nil)

		hasRole, err := svc.CheckUserHasRole(ctx, userID, roleSlug)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !hasRole {
			t.Error("Expected user to have role")
		}
	})

	t.Run("用户没有角色", func(t *testing.T) {
		userID := uint64(123)
		roleSlug := "super-admin"

		repo.EXPECT().
			CheckUserHasRole(ctx, userID, roleSlug).
			Return(false, nil)

		hasRole, err := svc.CheckUserHasRole(ctx, userID, roleSlug)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if hasRole {
			t.Error("Expected user not to have role")
		}
	})
}

// TestCheckUserIsSuperAdmin 测试检查用户是否为超级管理员。
func TestCheckUserIsSuperAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRoleRepository(ctrl)
	mockCache := cache.NewMemory()
	svc := NewRoleService(repo, mockCache)

	ctx := context.Background()

	t.Run("用户是超级管理员", func(t *testing.T) {
		userID := uint64(1)

		repo.EXPECT().
			CheckUserHasRole(ctx, userID, string(model.RoleSlugSuperAdmin)).
			Return(true, nil)

		isSuperAdmin, err := svc.CheckUserIsSuperAdmin(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !isSuperAdmin {
			t.Error("Expected user to be super admin")
		}
	})

	t.Run("用户不是超级管理员", func(t *testing.T) {
		userID := uint64(123)

		repo.EXPECT().
			CheckUserHasRole(ctx, userID, string(model.RoleSlugSuperAdmin)).
			Return(false, nil)

		isSuperAdmin, err := svc.CheckUserIsSuperAdmin(ctx, userID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if isSuperAdmin {
			t.Error("Expected user not to be super admin")
		}
	})
}

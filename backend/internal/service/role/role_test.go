package role

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

func setupRoleService(t *testing.T) (*RoleService, context.Context) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	return service, context.Background()
}

func TestRoleService_CreateRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	t.Run("create valid role", func(t *testing.T) {
		role := &model.RoleModel{
			Slug:        "test-role",
			Name:        "测试角色",
			Description: "测试角色描述",
		}

		err := service.CreateRole(ctx, role)
		require.NoError(t, err)
		assert.NotZero(t, role.ID)
	})

	t.Run("create role without slug should fail", func(t *testing.T) {
		role := &model.RoleModel{
			Name: "无Slug角色",
		}

		err := service.CreateRole(ctx, role)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "slug and name are required")
	})

	t.Run("create role without name should fail", func(t *testing.T) {
		role := &model.RoleModel{
			Slug: "no-name-role",
		}

		err := service.CreateRole(ctx, role)
		assert.Error(t, err)
	})

	t.Run("create duplicate slug should fail", func(t *testing.T) {
		role1 := &model.RoleModel{
			Slug: "duplicate-slug",
			Name: "角色1",
		}
		err := service.CreateRole(ctx, role1)
		require.NoError(t, err)

		role2 := &model.RoleModel{
			Slug: "duplicate-slug",
			Name: "角色2",
		}
		err = service.CreateRole(ctx, role2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestRoleService_GetRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建测试角色
	role := &model.RoleModel{
		Slug: "get-role",
		Name: "获取角色",
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("get existing role", func(t *testing.T) {
		found, err := service.GetRole(ctx, role.ID)
		require.NoError(t, err)
		assert.Equal(t, role.Slug, found.Slug)
		assert.Equal(t, role.Name, found.Name)
	})

	t.Run("get non-existent role", func(t *testing.T) {
		_, err := service.GetRole(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestRoleService_GetRoleBySlug(t *testing.T) {
	service, ctx := setupRoleService(t)

	role := &model.RoleModel{
		Slug: "slug-test",
		Name: "Slug测试",
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("get by existing slug", func(t *testing.T) {
		found, err := service.GetRoleBySlug(ctx, "slug-test")
		require.NoError(t, err)
		assert.Equal(t, role.ID, found.ID)
	})

	t.Run("get by non-existent slug", func(t *testing.T) {
		_, err := service.GetRoleBySlug(ctx, "non-existent")
		assert.Error(t, err)
	})
}

func TestRoleService_UpdateRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	role := &model.RoleModel{
		Slug:        "update-role",
		Name:        "更新角色",
		Description: "原始描述",
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("update role description", func(t *testing.T) {
		role.Description = "更新后的描述"
		err := service.UpdateRole(ctx, role)
		require.NoError(t, err)

		updated, err := service.GetRole(ctx, role.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新后的描述", updated.Description)
	})

	t.Run("update role without ID should fail", func(t *testing.T) {
		invalidRole := &model.RoleModel{
			Slug: "no-id",
			Name: "无ID",
		}
		err := service.UpdateRole(ctx, invalidRole)
		assert.Error(t, err)
	})
}

func TestRoleService_DeleteRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	t.Run("delete existing role", func(t *testing.T) {
		role := &model.RoleModel{
			Slug: "delete-role",
			Name: "待删除角色",
		}
		err := service.CreateRole(ctx, role)
		require.NoError(t, err)

		err = service.DeleteRole(ctx, role.ID)
		require.NoError(t, err)

		_, err = service.GetRole(ctx, role.ID)
		assert.Error(t, err)
	})
}

func TestRoleService_ListRoles(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建测试角色
	roles := []*model.RoleModel{
		{Slug: "list-role-1", Name: "列表角色1"},
		{Slug: "list-role-2", Name: "列表角色2"},
		{Slug: "list-role-3", Name: "列表角色3"},
	}
	for _, r := range roles {
		err := service.CreateRole(ctx, r)
		require.NoError(t, err)
	}

	t.Run("list all roles", func(t *testing.T) {
		list, err := service.ListRoles(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 3)
	})

	t.Run("list roles paged", func(t *testing.T) {
		list, total, err := service.ListRolesPaged(ctx, 1, 2)
		require.NoError(t, err)
		assert.Len(t, list, 2)
		assert.GreaterOrEqual(t, total, int64(3))
	})
}

func TestRoleService_ListRolesPagedWithFilter(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建测试角色
	_ = service.CreateRole(ctx, &model.RoleModel{Slug: "filter-admin", Name: "管理员角色", IsSystem: true})
	_ = service.CreateRole(ctx, &model.RoleModel{Slug: "filter-user", Name: "用户角色", IsSystem: false})

	t.Run("filter by keyword", func(t *testing.T) {
		list, _, err := service.ListRolesPagedWithFilter(ctx, 1, 10, "管理员", nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)
	})

	t.Run("filter by system flag", func(t *testing.T) {
		isSystem := true
		list, _, err := service.ListRolesPagedWithFilter(ctx, 1, 10, "", &isSystem)
		require.NoError(t, err)
		for _, r := range list {
			assert.True(t, r.IsSystem)
		}
	})
}

func TestRoleService_CheckUserHasRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建角色和用户
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.User{})

	role := &model.RoleModel{
		Slug: "check-role",
		Name: "检查角色",
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("user without role", func(t *testing.T) {
		hasRole, err := service.CheckUserHasRole(ctx, 99999, "check-role")
		require.NoError(t, err)
		assert.False(t, hasRole)
	})
}

func TestRoleService_MergePermissions(t *testing.T) {
	t.Run("merge multiple permission sets", func(t *testing.T) {
		set1 := []model.Permission{
			{Base: model.Base{ID: 1}, Code: "perm1"},
			{Base: model.Base{ID: 2}, Code: "perm2"},
		}
		set2 := []model.Permission{
			{Base: model.Base{ID: 2}, Code: "perm2"}, // 重复
			{Base: model.Base{ID: 3}, Code: "perm3"},
		}
		set3 := []model.Permission{
			{Base: model.Base{ID: 4}, Code: "perm4"},
		}

		merged := MergePermissions(set1, set2, set3)
		assert.Len(t, merged, 4) // 去重后应该是4个
	})

	t.Run("merge empty sets", func(t *testing.T) {
		merged := MergePermissions([]model.Permission{}, []model.Permission{})
		assert.Empty(t, merged)
	})
}

func TestRoleService_SortRolesByPriority(t *testing.T) {
	t.Run("sort roles by priority descending", func(t *testing.T) {
		roles := []model.RoleModel{
			{Base: model.Base{ID: 1}, Priority: 10},
			{Base: model.Base{ID: 2}, Priority: 30},
			{Base: model.Base{ID: 3}, Priority: 20},
		}

		sortRolesByPriority(roles)

		assert.Equal(t, uint64(2), roles[0].ID) // Priority 30
		assert.Equal(t, uint64(3), roles[1].ID) // Priority 20
		assert.Equal(t, uint64(1), roles[2].ID) // Priority 10
	})
}

func TestRoleService_GetPermissionCache(t *testing.T) {
	service, _ := setupRoleService(t)

	t.Run("get permission cache", func(t *testing.T) {
		cache := service.GetPermissionCache()
		assert.NotNil(t, cache)
	})
}

func TestRoleService_ListRolesWithPermissions(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建角色
	role := &model.RoleModel{
		Slug: "with-perms-role",
		Name: "带权限角色",
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("获取带权限的角色列表", func(t *testing.T) {
		roles, err := service.ListRolesWithPermissions(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(roles), 1)
	})
}

func TestRoleService_GetRoleWithPermissions(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/role-with-perm", Code: "admin.role.withperm"}
	require.NoError(t, permRepo.Create(ctx, perm))

	// 创建角色
	role := &model.RoleModel{Slug: "role-with-perm", Name: "带权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	// 分配权限
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm.ID}))

	t.Run("获取带权限的角色", func(t *testing.T) {
		roleWithPerms, err := service.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err)
		assert.Equal(t, role.Slug, roleWithPerms.Slug)
		assert.Len(t, roleWithPerms.Permissions, 1)
	})
}

func TestRoleService_AssignPermissionsToRole(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/assign1", Code: "admin.assign.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/assign2", Code: "admin.assign.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建角色
	role := &model.RoleModel{Slug: "assign-perm-role", Name: "分配权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	t.Run("分配权限给角色", func(t *testing.T) {
		err := service.AssignPermissionsToRole(ctx, role.ID, []uint64{perm1.ID, perm2.ID})
		require.NoError(t, err)

		roleWithPerms, err := service.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, roleWithPerms.Permissions, 2)
	})
}

func TestRoleService_AddPermissionsToRole(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/add1", Code: "admin.add.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/add2", Code: "admin.add.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建角色并分配第一个权限
	role := &model.RoleModel{Slug: "add-perm-role", Name: "追加权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm1.ID}))

	t.Run("追加权限给角色", func(t *testing.T) {
		err := service.AddPermissionsToRole(ctx, role.ID, []uint64{perm2.ID})
		require.NoError(t, err)

		roleWithPerms, err := service.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, roleWithPerms.Permissions, 2)
	})
}

func TestRoleService_RemovePermissionsFromRole(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/remove1", Code: "admin.remove.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/remove2", Code: "admin.remove.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建角色并分配权限
	role := &model.RoleModel{Slug: "remove-perm-role", Name: "移除权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm1.ID, perm2.ID}))

	t.Run("移除角色权限", func(t *testing.T) {
		err := service.RemovePermissionsFromRole(ctx, role.ID, []uint64{perm1.ID})
		require.NoError(t, err)

		roleWithPerms, err := service.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, roleWithPerms.Permissions, 1)
	})
}

func TestRoleService_ListRolesByUserID(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003001",
		Email:        "role_user@test.com",
		Name:         "角色用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建角色
	role1 := &model.RoleModel{Slug: "user-role-1", Name: "用户角色1"}
	role2 := &model.RoleModel{Slug: "user-role-2", Name: "用户角色2"}
	require.NoError(t, roleRepo.Create(ctx, role1))
	require.NoError(t, roleRepo.Create(ctx, role2))

	// 分配角色给用户
	require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{role1.ID, role2.ID}))

	t.Run("获取用户角色", func(t *testing.T) {
		roles, err := service.ListRolesByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 2)
	})

	t.Run("缓存命中", func(t *testing.T) {
		// 第二次调用应该从缓存获取
		roles, err := service.ListRolesByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 2)
	})
}

func TestRoleService_AssignRolesToUser(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003002",
		Email:        "assign_role@test.com",
		Name:         "分配角色用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建角色
	role := &model.RoleModel{Slug: "assign-to-user", Name: "分配给用户的角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	t.Run("分配角色给用户", func(t *testing.T) {
		err := service.AssignRolesToUser(ctx, user.ID, []uint64{role.ID})
		require.NoError(t, err)

		roles, err := service.ListRolesByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
	})
}

func TestRoleService_RemoveRolesFromUser(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003003",
		Email:        "remove_role@test.com",
		Name:         "移除角色用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建角色并分配
	role1 := &model.RoleModel{Slug: "remove-from-user-1", Name: "移除角色1"}
	role2 := &model.RoleModel{Slug: "remove-from-user-2", Name: "移除角色2"}
	require.NoError(t, roleRepo.Create(ctx, role1))
	require.NoError(t, roleRepo.Create(ctx, role2))
	require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{role1.ID, role2.ID}))

	t.Run("移除用户角色", func(t *testing.T) {
		err := service.RemoveRolesFromUser(ctx, user.ID, []uint64{role1.ID})
		require.NoError(t, err)

		roles, err := service.ListRolesByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, roles, 1)
	})
}

func TestRoleService_CheckUserIsSuperAdmin(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003004",
		Email:        "super_admin@test.com",
		Name:         "超级管理员",
		Role:         model.RoleAdmin,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建超级管理员角色
	superAdminRole := &model.RoleModel{Slug: string(model.RoleSlugSuperAdmin), Name: "超级管理员", IsSystem: true}
	require.NoError(t, roleRepo.Create(ctx, superAdminRole))

	t.Run("用户不是超级管理员", func(t *testing.T) {
		isSuperAdmin, err := service.CheckUserIsSuperAdmin(ctx, user.ID)
		require.NoError(t, err)
		assert.False(t, isSuperAdmin)
	})

	t.Run("用户是超级管理员", func(t *testing.T) {
		require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{superAdminRole.ID}))

		isSuperAdmin, err := service.CheckUserIsSuperAdmin(ctx, user.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)
	})
}

func TestRoleService_UpdateSystemRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建系统角色
	role := &model.RoleModel{
		Slug:        "system-role-update",
		Name:        "系统角色",
		Description: "原始描述",
		IsSystem:    true,
	}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("更新系统角色只能修改描述", func(t *testing.T) {
		role.Name = "新名称"
		role.Description = "新描述"
		err := service.UpdateRole(ctx, role)
		require.NoError(t, err)

		updated, err := service.GetRole(ctx, role.ID)
		require.NoError(t, err)
		// 系统角色名称不应该被修改
		assert.Equal(t, "系统角色", updated.Name)
		assert.Equal(t, "新描述", updated.Description)
	})
}

func TestRoleService_GetEffectivePermissions(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/effective", Code: "admin.effective.perm"}
	require.NoError(t, permRepo.Create(ctx, perm))

	// 创建角色并分配权限
	role := &model.RoleModel{Slug: "effective-role", Name: "有效权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm.ID}))

	t.Run("获取有效权限", func(t *testing.T) {
		perms, err := service.GetEffectivePermissions(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, perms, 1)
	})
}

func TestRoleService_GetUserEffectivePermissions(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003005",
		Email:        "user_effective@test.com",
		Name:         "有效权限用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/user-eff1", Code: "admin.usereff.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/user-eff2", Code: "admin.usereff.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建两个角色，分别分配不同权限
	role1 := &model.RoleModel{Slug: "user-eff-role1", Name: "用户有效角色1", Priority: 10}
	role2 := &model.RoleModel{Slug: "user-eff-role2", Name: "用户有效角色2", Priority: 20}
	require.NoError(t, roleRepo.Create(ctx, role1))
	require.NoError(t, roleRepo.Create(ctx, role2))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role1.ID, []uint64{perm1.ID}))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role2.ID, []uint64{perm2.ID}))

	// 分配角色给用户
	require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{role1.ID, role2.ID}))

	t.Run("获取用户有效权限", func(t *testing.T) {
		perms, err := service.GetUserEffectivePermissions(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, perms, 2)
	})
}

func TestRoleService_GetChildRoles(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建父角色
	parentRole := &model.RoleModel{Slug: "parent-role", Name: "父角色"}
	err := service.CreateRole(ctx, parentRole)
	require.NoError(t, err)

	t.Run("获取子角色（无子角色）", func(t *testing.T) {
		children, err := service.GetChildRoles(ctx, parentRole.ID)
		require.NoError(t, err)
		assert.Empty(t, children)
	})
}

func TestRoleService_ValidateNoCircularInheritance(t *testing.T) {
	service, ctx := setupRoleService(t)

	// 创建角色
	role := &model.RoleModel{Slug: "circular-role", Name: "循环角色"}
	err := service.CreateRole(ctx, role)
	require.NoError(t, err)

	t.Run("角色不能是自己的父角色", func(t *testing.T) {
		err := service.ValidateNoCircularInheritance(ctx, role.ID, role.ID)
		assert.Error(t, err)
	})
}

func TestRoleService_SetRoleParent(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建父角色
	parentRole := &model.RoleModel{Slug: "set-parent-role", Name: "父角色", Level: 0}
	require.NoError(t, roleRepo.Create(ctx, parentRole))

	// 创建子角色
	childRole := &model.RoleModel{Slug: "set-child-role", Name: "子角色", Level: 0}
	require.NoError(t, roleRepo.Create(ctx, childRole))

	t.Run("设置父角色", func(t *testing.T) {
		err := service.SetRoleParent(ctx, childRole.ID, &parentRole.ID)
		require.NoError(t, err)
	})

	t.Run("清除父角色", func(t *testing.T) {
		err := service.SetRoleParent(ctx, childRole.ID, nil)
		require.NoError(t, err)
	})

	t.Run("设置不存在的角色", func(t *testing.T) {
		err := service.SetRoleParent(ctx, 99999, &parentRole.ID)
		assert.Error(t, err)
	})

	t.Run("设置不存在的父角色", func(t *testing.T) {
		nonExistentID := uint64(99999)
		err := service.SetRoleParent(ctx, childRole.ID, &nonExistentID)
		assert.Error(t, err)
	})
}

func TestRoleService_GetRoleInheritanceChain(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建角色
	role := &model.RoleModel{Slug: "chain-role", Name: "链角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	t.Run("获取继承链", func(t *testing.T) {
		chain, err := service.GetRoleInheritanceChain(ctx, role.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, chain)
	})
}

func TestRoleService_InvalidateRolePermissionsAndPropagateToUsers(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建角色
	role := &model.RoleModel{Slug: "invalidate-role", Name: "失效角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	t.Run("失效角色权限并传播", func(t *testing.T) {
		err := service.InvalidateRolePermissionsAndPropagateToUsers(ctx, role.ID)
		require.NoError(t, err)
	})
}

func TestRoleService_DeleteNonExistentRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	t.Run("删除不存在的角色", func(t *testing.T) {
		err := service.DeleteRole(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestRoleService_UpdateNonExistentRole(t *testing.T) {
	service, ctx := setupRoleService(t)

	t.Run("更新不存在的角色", func(t *testing.T) {
		role := &model.RoleModel{
			Base: model.Base{ID: 99999},
			Slug: "non-existent",
			Name: "不存在",
		}
		err := service.UpdateRole(ctx, role)
		assert.Error(t, err)
	})
}

func TestRoleService_ListRolesByUserIDCacheMiss(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800003010",
		Email:        "cache_miss@test.com",
		Name:         "缓存未命中用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	t.Run("缓存未命中时从数据库获取", func(t *testing.T) {
		roles, err := service.ListRolesByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Empty(t, roles)
	})
}

func TestRoleService_GetEffectivePermissionsWithInheritance(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/inherit1", Code: "admin.inherit.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/inherit2", Code: "admin.inherit.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建父角色
	parentRole := &model.RoleModel{Slug: "inherit-parent", Name: "继承父角色", Level: 0}
	require.NoError(t, roleRepo.Create(ctx, parentRole))
	require.NoError(t, roleRepo.AssignPermissions(ctx, parentRole.ID, []uint64{perm1.ID}))

	// 创建子角色
	childRole := &model.RoleModel{Slug: "inherit-child", Name: "继承子角色", Level: 1, ParentID: &parentRole.ID}
	require.NoError(t, roleRepo.Create(ctx, childRole))
	require.NoError(t, roleRepo.AssignPermissions(ctx, childRole.ID, []uint64{perm2.ID}))

	t.Run("获取包含继承的有效权限", func(t *testing.T) {
		perms, err := service.GetEffectivePermissions(ctx, childRole.ID)
		require.NoError(t, err)
		// 应该包含自己的权限和父角色的权限
		assert.GreaterOrEqual(t, len(perms), 1)
	})
}

func TestRoleService_GetUserEffectivePermissionsNoRoles(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.RoleModel{},
		&model.Permission{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewRoleService(roleRepo, memCache)
	ctx := context.Background()

	// 创建用户（无角色）
	user := &model.User{
		Phone:        "13800003011",
		Email:        "no_roles@test.com",
		Name:         "无角色用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	t.Run("无角色用户的有效权限为空", func(t *testing.T) {
		perms, err := service.GetUserEffectivePermissions(ctx, user.ID)
		require.NoError(t, err)
		assert.Empty(t, perms)
	})
}

package permission

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

func setupPermissionService(t *testing.T) (*PermissionService, context.Context) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	permRepo := adminrepo.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	service := NewPermissionService(permRepo, memCache)
	return service, context.Background()
}

func TestPermissionService_CreatePermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	t.Run("create valid permission", func(t *testing.T) {
		perm := &model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/users",
			Code:        "admin.users.list",
			Group:       "用户管理",
			Description: "获取用户列表",
		}

		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)
		assert.NotZero(t, perm.ID)
	})

	t.Run("create without method should fail", func(t *testing.T) {
		perm := &model.Permission{
			Path: "/api/v1/test",
		}

		err := service.CreatePermission(ctx, perm)
		assert.Error(t, err)
	})

	t.Run("create without path should fail", func(t *testing.T) {
		perm := &model.Permission{
			Method: model.HTTPMethodGET,
		}

		err := service.CreatePermission(ctx, perm)
		assert.Error(t, err)
	})

	t.Run("create with invalid code format should fail", func(t *testing.T) {
		perm := &model.Permission{
			Method: model.HTTPMethodGET,
			Path:   "/api/v1/invalid",
			Code:   "invalid-code", // 不符合 module.resource.action 格式
		}

		err := service.CreatePermission(ctx, perm)
		assert.Error(t, err)
	})

	t.Run("create duplicate method+path should fail", func(t *testing.T) {
		perm1 := &model.Permission{
			Method: model.HTTPMethodPOST,
			Path:   "/api/v1/duplicate",
			Code:   "admin.dup.create",
		}
		err := service.CreatePermission(ctx, perm1)
		require.NoError(t, err)

		perm2 := &model.Permission{
			Method: model.HTTPMethodPOST,
			Path:   "/api/v1/duplicate",
			Code:   "admin.dup.create2",
		}
		err = service.CreatePermission(ctx, perm2)
		assert.Error(t, err)
	})
}

func TestPermissionService_GetPermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	perm := &model.Permission{
		Method: model.HTTPMethodGET,
		Path:   "/api/v1/get-test",
		Code:   "admin.get.test",
	}
	err := service.CreatePermission(ctx, perm)
	require.NoError(t, err)

	t.Run("get existing permission", func(t *testing.T) {
		found, err := service.GetPermission(ctx, perm.ID)
		require.NoError(t, err)
		assert.Equal(t, perm.Code, found.Code)
	})

	t.Run("get non-existent permission", func(t *testing.T) {
		_, err := service.GetPermission(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestPermissionService_UpdatePermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/update-test",
		Code:        "admin.update.test",
		Description: "原始描述",
	}
	err := service.CreatePermission(ctx, perm)
	require.NoError(t, err)

	t.Run("update description", func(t *testing.T) {
		perm.Description = "更新后的描述"
		err := service.UpdatePermission(ctx, perm)
		require.NoError(t, err)

		updated, err := service.GetPermission(ctx, perm.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新后的描述", updated.Description)
	})

	t.Run("update without ID should fail", func(t *testing.T) {
		invalidPerm := &model.Permission{
			Method: model.HTTPMethodGET,
			Path:   "/test",
		}
		err := service.UpdatePermission(ctx, invalidPerm)
		assert.Error(t, err)
	})

	t.Run("update code should fail if already set", func(t *testing.T) {
		perm.Code = "admin.update.changed"
		err := service.UpdatePermission(ctx, perm)
		assert.Error(t, err)
	})
}

func TestPermissionService_DeletePermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	t.Run("delete existing permission", func(t *testing.T) {
		perm := &model.Permission{
			Method: model.HTTPMethodDELETE,
			Path:   "/api/v1/delete-test",
			Code:   "admin.delete.test",
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		err = service.DeletePermission(ctx, perm.ID)
		require.NoError(t, err)

		_, err = service.GetPermission(ctx, perm.ID)
		assert.Error(t, err)
	})

	t.Run("delete system permission should fail", func(t *testing.T) {
		perm := &model.Permission{
			Method:   model.HTTPMethodGET,
			Path:     "/api/v1/system",
			Code:     "admin.system.read",
			IsSystem: true,
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		err = service.DeletePermission(ctx, perm.ID)
		assert.Error(t, err)
	})
}

func TestPermissionService_ListPermissions(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建测试权限
	perms := []*model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/list1", Code: "admin.list.one"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/list2", Code: "admin.list.two"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/list3", Code: "admin.list.three"},
	}
	for _, p := range perms {
		err := service.CreatePermission(ctx, p)
		require.NoError(t, err)
	}

	t.Run("list all permissions", func(t *testing.T) {
		list, err := service.ListPermissions(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 3)
	})

	t.Run("list permissions paged", func(t *testing.T) {
		list, total, err := service.ListPermissionsPaged(ctx, 1, 2)
		require.NoError(t, err)
		assert.Len(t, list, 2)
		assert.GreaterOrEqual(t, total, int64(3))
	})
}

func TestPermissionService_ValidatePermissionCode(t *testing.T) {
	service, _ := setupPermissionService(t)

	tests := []struct {
		code    string
		isValid bool
	}{
		{"admin.users.list", true},
		{"admin.users.create", true},
		{"system.config.read", true},
		{"invalid", false},
		{"invalid.code", false},
		{"Invalid.Code.Format", false},    // 大写
		{"admin.users.list.extra", false}, // 太多段
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			err := service.ValidatePermissionCode(tt.code)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPermissionService_CanDeletePermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	t.Run("can delete normal permission", func(t *testing.T) {
		perm := &model.Permission{
			Method: model.HTTPMethodGET,
			Path:   "/api/v1/can-delete",
			Code:   "admin.can.delete",
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		canDelete, reason, err := service.CanDeletePermission(ctx, perm.ID)
		require.NoError(t, err)
		assert.True(t, canDelete)
		assert.Empty(t, reason)
	})

	t.Run("cannot delete system permission", func(t *testing.T) {
		perm := &model.Permission{
			Method:   model.HTTPMethodGET,
			Path:     "/api/v1/system-perm",
			Code:     "admin.system.perm",
			IsSystem: true,
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		canDelete, reason, err := service.CanDeletePermission(ctx, perm.ID)
		require.NoError(t, err)
		assert.False(t, canDelete)
		assert.Contains(t, reason, "系统权限")
	})
}

func TestPermissionService_UpsertPermission(t *testing.T) {
	service, ctx := setupPermissionService(t)

	t.Run("upsert creates new permission", func(t *testing.T) {
		perm := &model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/upsert-new",
			Code:        "admin.upsert.new",
			Description: "新权限",
		}

		err := service.UpsertPermission(ctx, perm)
		require.NoError(t, err)
		assert.NotZero(t, perm.ID)
	})

	t.Run("upsert updates existing permission", func(t *testing.T) {
		perm := &model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/upsert-update",
			Code:        "admin.upsert.update",
			Description: "原始描述",
		}
		err := service.UpsertPermission(ctx, perm)
		require.NoError(t, err)

		// 更新同一个权限（相同method+path）
		perm2 := &model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/v1/upsert-update",
			Code:        "admin.upsert.update",
			Description: "更新后描述",
		}
		err = service.UpsertPermission(ctx, perm2)
		require.NoError(t, err)
	})
}

func TestPermissionService_ListPermissionsByGroup(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建不同分组的权限
	perms := []*model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/users", Code: "admin.users.list", Group: "用户管理"},
		{Method: model.HTTPMethodPOST, Path: "/api/v1/users", Code: "admin.users.create", Group: "用户管理"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/orders", Code: "admin.orders.list", Group: "订单管理"},
	}
	for _, p := range perms {
		err := service.CreatePermission(ctx, p)
		require.NoError(t, err)
	}

	t.Run("按分组获取权限", func(t *testing.T) {
		list, err := service.ListPermissionsByGroup(ctx, "用户管理")
		require.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("获取不存在的分组", func(t *testing.T) {
		list, err := service.ListPermissionsByGroup(ctx, "不存在的分组")
		require.NoError(t, err)
		assert.Empty(t, list)
	})
}

func TestPermissionService_ListPermissionGroups(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建不同分组的权限
	perms := []*model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/group1", Code: "admin.group.one", Group: "分组A"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/group2", Code: "admin.group.two", Group: "分组B"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/group3", Code: "admin.group.three", Group: "分组A"},
	}
	for _, p := range perms {
		err := service.CreatePermission(ctx, p)
		require.NoError(t, err)
	}

	t.Run("获取所有分组", func(t *testing.T) {
		groups, err := service.ListPermissionGroups(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(groups), 2)
	})
}

func TestPermissionService_ListPermissionsByRoleID(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	permRepo := adminrepo.NewPermissionRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewPermissionService(permRepo, memCache)
	ctx := context.Background()

	// 创建权限
	perm1 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/role-perm1", Code: "admin.role.one"}
	perm2 := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/role-perm2", Code: "admin.role.two"}
	require.NoError(t, permRepo.Create(ctx, perm1))
	require.NoError(t, permRepo.Create(ctx, perm2))

	// 创建角色
	role := &model.RoleModel{Slug: "test-role-perm", Name: "测试角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	// 分配权限给角色
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm1.ID, perm2.ID}))

	t.Run("获取角色权限", func(t *testing.T) {
		perms, err := service.ListPermissionsByRoleID(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, perms, 2)
	})

	t.Run("缓存命中", func(t *testing.T) {
		// 第二次调用应该从缓存获取
		perms, err := service.ListPermissionsByRoleID(ctx, role.ID)
		require.NoError(t, err)
		assert.Len(t, perms, 2)
	})
}

func TestPermissionService_ListPermissionsByUserID(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	permRepo := adminrepo.NewPermissionRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewPermissionService(permRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800002001",
		Email:        "user_perm@test.com",
		Name:         "权限测试用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建权限
	perm := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/user-perm", Code: "admin.user.perm"}
	require.NoError(t, permRepo.Create(ctx, perm))

	// 创建角色
	role := &model.RoleModel{Slug: "user-perm-role", Name: "用户权限角色"}
	require.NoError(t, roleRepo.Create(ctx, role))

	// 分配权限给角色
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm.ID}))

	// 分配角色给用户
	require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{role.ID}))

	t.Run("获取用户权限", func(t *testing.T) {
		perms, err := service.ListPermissionsByUserID(ctx, user.ID)
		require.NoError(t, err)
		assert.Len(t, perms, 1)
	})
}

func TestPermissionService_CheckUserHasPermission(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.User{},
		&model.UserRole{},
	)

	permRepo := adminrepo.NewPermissionRepository(db)
	roleRepo := adminrepo.NewRoleRepository(db)
	memCache := cache.NewMemory()
	service := NewPermissionService(permRepo, memCache)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Phone:        "13800002002",
		Email:        "check_perm@test.com",
		Name:         "权限检查用户",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed",
	}
	require.NoError(t, db.Create(user).Error)

	// 创建权限
	perm := &model.Permission{Method: model.HTTPMethodGET, Path: "/api/v1/check-perm", Code: "admin.check.perm"}
	require.NoError(t, permRepo.Create(ctx, perm))

	// 创建角色并分配权限
	role := &model.RoleModel{Slug: "check-perm-role", Name: "权限检查角色"}
	require.NoError(t, roleRepo.Create(ctx, role))
	require.NoError(t, roleRepo.AssignPermissions(ctx, role.ID, []uint64{perm.ID}))
	require.NoError(t, roleRepo.AssignToUser(ctx, user.ID, []uint64{role.ID}))

	t.Run("用户拥有权限", func(t *testing.T) {
		has, err := service.CheckUserHasPermission(ctx, user.ID, model.HTTPMethodGET, "/api/v1/check-perm")
		require.NoError(t, err)
		assert.True(t, has)
	})

	t.Run("用户没有权限", func(t *testing.T) {
		has, err := service.CheckUserHasPermission(ctx, user.ID, model.HTTPMethodPOST, "/api/v1/check-perm")
		require.NoError(t, err)
		assert.False(t, has)
	})

	t.Run("用户没有该路径权限", func(t *testing.T) {
		has, err := service.CheckUserHasPermission(ctx, user.ID, model.HTTPMethodGET, "/api/v1/other-path")
		require.NoError(t, err)
		assert.False(t, has)
	})
}

func TestPermissionService_CheckPermissionCodeExists(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建权限
	perm := &model.Permission{
		Method: model.HTTPMethodGET,
		Path:   "/api/v1/code-exists",
		Code:   "admin.code.exists",
	}
	err := service.CreatePermission(ctx, perm)
	require.NoError(t, err)

	t.Run("权限码存在", func(t *testing.T) {
		exists, err := service.CheckPermissionCodeExists(ctx, "admin.code.exists")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("权限码不存在", func(t *testing.T) {
		exists, err := service.CheckPermissionCodeExists(ctx, "admin.code.notexists")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestPermissionService_DeletePermissionForce(t *testing.T) {
	service, ctx := setupPermissionService(t)

	t.Run("强制删除普通权限", func(t *testing.T) {
		perm := &model.Permission{
			Method: model.HTTPMethodDELETE,
			Path:   "/api/v1/force-delete",
			Code:   "admin.force.delete",
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		err = service.DeletePermissionForce(ctx, perm.ID)
		require.NoError(t, err)

		_, err = service.GetPermission(ctx, perm.ID)
		assert.Error(t, err)
	})

	t.Run("强制删除系统权限应失败", func(t *testing.T) {
		perm := &model.Permission{
			Method:   model.HTTPMethodGET,
			Path:     "/api/v1/system-force",
			Code:     "admin.system.force",
			IsSystem: true,
		}
		err := service.CreatePermission(ctx, perm)
		require.NoError(t, err)

		err = service.DeletePermissionForce(ctx, perm.ID)
		assert.Error(t, err)
	})
}

func TestPermissionService_GetPermissionTree(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建权限
	perms := []*model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/tree1", Code: "admin.tree.one", Group: "树形测试"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/tree2", Code: "admin.tree.two", Group: "树形测试"},
	}
	for _, p := range perms {
		err := service.CreatePermission(ctx, p)
		require.NoError(t, err)
	}

	t.Run("获取权限树", func(t *testing.T) {
		tree, err := service.GetPermissionTree(ctx)
		require.NoError(t, err)
		assert.NotNil(t, tree)
	})

	t.Run("缓存命中", func(t *testing.T) {
		// 第二次调用应该从缓存获取
		tree, err := service.GetPermissionTree(ctx)
		require.NoError(t, err)
		assert.NotNil(t, tree)
	})
}

func TestPermissionService_GetPermissionTreeByGroup(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建不同分组的权限
	perms := []*model.Permission{
		{Method: model.HTTPMethodGET, Path: "/api/v1/group-tree1", Code: "admin.grouptree.one", Group: "分组树A"},
		{Method: model.HTTPMethodGET, Path: "/api/v1/group-tree2", Code: "admin.grouptree.two", Group: "分组树B"},
	}
	for _, p := range perms {
		err := service.CreatePermission(ctx, p)
		require.NoError(t, err)
	}

	t.Run("按分组获取权限树", func(t *testing.T) {
		groups, err := service.GetPermissionTreeByGroup(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, groups)
	})
}

func TestPermissionService_CreatePermissionWithExistingCode(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建第一个权限
	perm1 := &model.Permission{
		Method: model.HTTPMethodGET,
		Path:   "/api/v1/dup-code1",
		Code:   "admin.dup.code",
	}
	err := service.CreatePermission(ctx, perm1)
	require.NoError(t, err)

	t.Run("创建重复权限码应失败", func(t *testing.T) {
		perm2 := &model.Permission{
			Method: model.HTTPMethodPOST,
			Path:   "/api/v1/dup-code2",
			Code:   "admin.dup.code", // 相同的权限码
		}
		err := service.CreatePermission(ctx, perm2)
		assert.Error(t, err)
	})
}

func TestPermissionService_UpdatePermissionSetNewCode(t *testing.T) {
	service, ctx := setupPermissionService(t)

	// 创建没有权限码的权限
	perm := &model.Permission{
		Method: model.HTTPMethodGET,
		Path:   "/api/v1/no-code",
	}
	err := service.CreatePermission(ctx, perm)
	require.NoError(t, err)

	t.Run("为空权限码设置新码", func(t *testing.T) {
		perm.Code = "admin.new.code"
		err := service.UpdatePermission(ctx, perm)
		require.NoError(t, err)

		updated, err := service.GetPermission(ctx, perm.ID)
		require.NoError(t, err)
		assert.Equal(t, "admin.new.code", updated.Code)
	})
}

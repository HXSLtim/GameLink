package menu

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/pkg/testutil"
)

func setupMenuService(t *testing.T) (*Service, context.Context) {
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db, &model.Menu{})

	menuRepo := adminrepo.NewMenuRepository(db)
	service := NewService(menuRepo)
	return service, context.Background()
}

func TestMenuService_Create(t *testing.T) {
	service, ctx := setupMenuService(t)

	t.Run("create valid menu", func(t *testing.T) {
		menu := &model.Menu{
			Name:      "用户管理",
			Path:      "/admin/users",
			Component: "UserList",
			Icon:      "user",
			Order:     1,
		}

		err := service.Create(ctx, menu)
		require.NoError(t, err)
		assert.NotZero(t, menu.ID)
	})

	t.Run("create child menu", func(t *testing.T) {
		parent := &model.Menu{
			Name:      "系统设置",
			Path:      "/admin/system",
			Component: "SystemLayout",
			Order:     10,
		}
		err := service.Create(ctx, parent)
		require.NoError(t, err)

		child := &model.Menu{
			Name:      "角色管理",
			Path:      "/admin/system/roles",
			Component: "RoleList",
			ParentID:  &parent.ID,
			Order:     1,
		}
		err = service.Create(ctx, child)
		require.NoError(t, err)
		assert.Equal(t, parent.ID, *child.ParentID)
	})
}

func TestMenuService_Get(t *testing.T) {
	service, ctx := setupMenuService(t)

	menu := &model.Menu{
		Name:      "获取测试",
		Path:      "/test",
		Component: "Test",
	}
	err := service.Create(ctx, menu)
	require.NoError(t, err)

	t.Run("get existing menu", func(t *testing.T) {
		found, err := service.Get(ctx, menu.ID)
		require.NoError(t, err)
		assert.Equal(t, menu.Name, found.Name)
	})

	t.Run("get non-existent menu", func(t *testing.T) {
		_, err := service.Get(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestMenuService_Update(t *testing.T) {
	service, ctx := setupMenuService(t)

	menu := &model.Menu{
		Name:      "更新测试",
		Path:      "/update",
		Component: "Update",
	}
	err := service.Create(ctx, menu)
	require.NoError(t, err)

	t.Run("update existing menu", func(t *testing.T) {
		menu.Name = "更新后名称"
		menu.Icon = "new-icon"
		err := service.Update(ctx, menu)
		require.NoError(t, err)

		updated, err := service.Get(ctx, menu.ID)
		require.NoError(t, err)
		assert.Equal(t, "更新后名称", updated.Name)
		assert.Equal(t, "new-icon", updated.Icon)
	})

	t.Run("update with zero ID should fail", func(t *testing.T) {
		invalidMenu := &model.Menu{
			Name: "无ID",
		}
		err := service.Update(ctx, invalidMenu)
		assert.Error(t, err)
	})
}

func TestMenuService_Delete(t *testing.T) {
	service, ctx := setupMenuService(t)

	t.Run("delete existing menu", func(t *testing.T) {
		menu := &model.Menu{
			Name:      "待删除",
			Path:      "/delete",
			Component: "Delete",
		}
		err := service.Create(ctx, menu)
		require.NoError(t, err)

		err = service.Delete(ctx, menu.ID)
		require.NoError(t, err)

		_, err = service.Get(ctx, menu.ID)
		assert.Error(t, err)
	})
}

func TestMenuService_List(t *testing.T) {
	service, ctx := setupMenuService(t)

	// 创建测试菜单
	parent := &model.Menu{Name: "父菜单", Path: "/parent", Component: "Parent"}
	err := service.Create(ctx, parent)
	require.NoError(t, err)

	child1 := &model.Menu{Name: "子菜单1", Path: "/parent/child1", Component: "Child1", ParentID: &parent.ID}
	child2 := &model.Menu{Name: "子菜单2", Path: "/parent/child2", Component: "Child2", ParentID: &parent.ID}
	_ = service.Create(ctx, child1)
	_ = service.Create(ctx, child2)

	t.Run("list root menus", func(t *testing.T) {
		menus, err := service.List(ctx, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(menus), 1)
	})

	t.Run("list child menus", func(t *testing.T) {
		menus, err := service.List(ctx, &parent.ID)
		require.NoError(t, err)
		assert.Len(t, menus, 2)
	})
}

func TestMenuService_ListPaged(t *testing.T) {
	service, ctx := setupMenuService(t)

	// 创建测试菜单
	for i := 0; i < 5; i++ {
		menu := &model.Menu{
			Name:      "分页菜单",
			Path:      "/paged",
			Component: "Paged",
		}
		_ = service.Create(ctx, menu)
	}

	t.Run("list paged", func(t *testing.T) {
		menus, total, err := service.ListPaged(ctx, 1, 3, nil)
		require.NoError(t, err)
		assert.Len(t, menus, 3)
		assert.GreaterOrEqual(t, total, int64(5))
	})
}

func TestMenuService_ListAccessible(t *testing.T) {
	service, ctx := setupMenuService(t)

	t.Run("list with empty codes returns empty", func(t *testing.T) {
		menus, err := service.ListAccessible(ctx, []string{})
		require.NoError(t, err)
		assert.Empty(t, menus)
	})

	t.Run("list with codes", func(t *testing.T) {
		// 创建带权限码的菜单
		menu := &model.Menu{
			Name:       "权限菜单",
			Path:       "/perm",
			Component:  "Perm",
			Permission: "admin.users.read",
		}
		err := service.Create(ctx, menu)
		require.NoError(t, err)

		menus, err := service.ListAccessible(ctx, []string{"admin.users.read"})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(menus), 1)
	})
}

// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/menu"
	menuservice "gamelink/internal/service/menu"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMenuService_Create(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create menu
	newMenu := &model.Menu{
		Name:      "Dashboard",
		Path:      "/dashboard",
		Component: "Dashboard",
		Icon:      "dashboard",
		Order:     1,
	}

	err := svc.Create(ctx, newMenu)
	require.NoError(t, err)
	assert.NotZero(t, newMenu.ID)

	// Verify in database
	saved, err := svc.Get(ctx, newMenu.ID)
	require.NoError(t, err)
	assert.Equal(t, "Dashboard", saved.Name)
	assert.Equal(t, "/dashboard", saved.Path)
}

func TestMenuService_CreateWithParent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create parent menu
	parentMenu := CreateTestMenu(t, db, "System", "/system", nil)

	// Create child menu
	childMenu := &model.Menu{
		Name:      "Users",
		Path:      "/system/users",
		Component: "Users",
		ParentID:  &parentMenu.ID,
		Order:     1,
	}

	err := svc.Create(ctx, childMenu)
	require.NoError(t, err)
	assert.NotZero(t, childMenu.ID)

	// Verify parent relationship
	saved, err := svc.Get(ctx, childMenu.ID)
	require.NoError(t, err)
	assert.Equal(t, parentMenu.ID, *saved.ParentID)
}

func TestMenuService_Update(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create menu
	testMenu := CreateTestMenu(t, db, "Original", "/original", nil)

	// Update menu
	testMenu.Name = "Updated"
	testMenu.Path = "/updated"
	err := svc.Update(ctx, testMenu)
	require.NoError(t, err)

	// Verify update
	saved, err := svc.Get(ctx, testMenu.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", saved.Name)
	assert.Equal(t, "/updated", saved.Path)
}

func TestMenuService_Delete(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create menu
	testMenu := CreateTestMenu(t, db, "ToDelete", "/delete", nil)

	// Delete menu
	err := svc.Delete(ctx, testMenu.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = svc.Get(ctx, testMenu.ID)
	assert.Error(t, err)
}

func TestMenuService_List(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create multiple menus
	CreateTestMenu(t, db, "Menu1", "/menu1", nil)
	CreateTestMenu(t, db, "Menu2", "/menu2", nil)
	CreateTestMenu(t, db, "Menu3", "/menu3", nil)

	// List all menus (no parent filter)
	menus, err := svc.List(ctx, nil)
	require.NoError(t, err)
	assert.Len(t, menus, 3)
}

func TestMenuService_ListByParent(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create parent and child menus
	parentMenu := CreateTestMenu(t, db, "Parent", "/parent", nil)
	CreateTestMenu(t, db, "Child1", "/parent/child1", &parentMenu.ID)
	CreateTestMenu(t, db, "Child2", "/parent/child2", &parentMenu.ID)
	CreateTestMenu(t, db, "Other", "/other", nil)

	// List children of parent
	children, err := svc.List(ctx, &parentMenu.ID)
	require.NoError(t, err)
	assert.Len(t, children, 2)
}

func TestMenuService_ListPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create multiple menus
	for i := 0; i < 5; i++ {
		CreateTestMenu(t, db, "PagedMenu"+string(rune('A'+i)), "/paged"+string(rune('a'+i)), nil)
	}

	// List with pagination
	menus, total, err := svc.ListPaged(ctx, 1, 2, nil)
	require.NoError(t, err)
	assert.Len(t, menus, 2)
	assert.Equal(t, int64(5), total)
}

func TestMenuService_ListAccessible(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Create menus with permissions
	menu1 := &model.Menu{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:       "Admin",
		Path:       "/admin",
		Permission: "admin.dashboard.view",
	}
	db.Create(menu1)

	menu2 := &model.Menu{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:       "Users",
		Path:       "/users",
		Permission: "admin.users.view",
	}
	db.Create(menu2)

	menu3 := &model.Menu{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Name:       "Settings",
		Path:       "/settings",
		Permission: "admin.settings.view",
	}
	db.Create(menu3)

	// List accessible menus with specific permission codes
	codes := []string{"admin.dashboard.view", "admin.users.view"}
	menus, err := svc.ListAccessible(ctx, codes)
	require.NoError(t, err)
	assert.Len(t, menus, 2)
}

func TestMenuService_ListAccessible_EmptyCodes(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// List with empty codes
	menus, err := svc.ListAccessible(ctx, []string{})
	require.NoError(t, err)
	assert.Len(t, menus, 0)
}

func TestMenuService_Get_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Try to get non-existent menu
	_, err := svc.Get(ctx, 99999)
	assert.Error(t, err)
}

func TestMenuService_Update_NotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	menuRepo := menu.NewMenuRepository(db)
	svc := menuservice.NewService(menuRepo)

	// Try to update non-existent menu
	nonExistent := &model.Menu{
		Base: model.Base{
			ID: 0, // ID is 0, should fail validation
		},
		Name: "NonExistent",
		Path: "/nonexistent",
	}

	err := svc.Update(ctx, nonExistent)
	assert.Error(t, err)
}

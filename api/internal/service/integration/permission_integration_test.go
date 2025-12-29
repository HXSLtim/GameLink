// Package integration provides integration tests for services.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	"gamelink/internal/repository/admin"
	"gamelink/internal/repository/permission"
	"gamelink/internal/service/role"
	permissionservice "gamelink/internal/service/permission"
	"gamelink/pkg/cache"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionService_CreatePermission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create permission
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/test",
		Code:        "test.resource.read",
		Group:       "test",
		Description: "Test permission",
	}

	err := svc.CreatePermission(ctx, perm)
	require.NoError(t, err)
	assert.NotZero(t, perm.ID)

	// Verify in database
	saved, err := svc.GetPermission(ctx, perm.ID)
	require.NoError(t, err)
	assert.Equal(t, perm.Code, saved.Code)
	assert.Equal(t, perm.Path, saved.Path)
}

func TestPermissionService_CreatePermission_DuplicateCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create first permission
	perm1 := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/test1",
		Code:        "test.resource.read",
		Group:       "test",
		Description: "Test permission 1",
	}
	err := svc.CreatePermission(ctx, perm1)
	require.NoError(t, err)

	// Try to create with same code
	perm2 := &model.Permission{
		Method:      model.HTTPMethodPOST,
		Path:        "/api/v1/test2",
		Code:        "test.resource.read", // Same code
		Group:       "test",
		Description: "Test permission 2",
	}
	err = svc.CreatePermission(ctx, perm2)
	assert.Error(t, err)
}

func TestPermissionService_CreatePermission_InvalidCode(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create permission with invalid code format
	perm := &model.Permission{
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/test",
		Code:        "invalid-code", // Invalid format
		Group:       "test",
		Description: "Test permission",
	}

	err := svc.CreatePermission(ctx, perm)
	assert.Error(t, err)
}

func TestPermissionService_UpdatePermission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create permission
	perm := CreateTestPermission(t, db, "GET", "/api/v1/update", "test.update.read")

	// Update permission
	perm.Description = "Updated description"
	err := svc.UpdatePermission(ctx, perm)
	require.NoError(t, err)

	// Verify update
	saved, err := svc.GetPermission(ctx, perm.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", saved.Description)
}

func TestPermissionService_DeletePermission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create permission
	perm := CreateTestPermission(t, db, "GET", "/api/v1/delete", "test.delete.read")

	// Delete permission
	err := svc.DeletePermission(ctx, perm.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = svc.GetPermission(ctx, perm.ID)
	assert.Error(t, err)
}

func TestPermissionService_DeletePermission_SystemPermission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create system permission
	perm := &model.Permission{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Method:   model.HTTPMethodGET,
		Path:     "/api/v1/system",
		Code:     "test.system.read",
		Group:    "test",
		IsSystem: true,
	}
	err := db.Create(perm).Error
	require.NoError(t, err)

	// Try to delete system permission
	err = svc.DeletePermission(ctx, perm.ID)
	assert.Error(t, err)
}

func TestPermissionService_ListPermissions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create multiple permissions
	CreateTestPermission(t, db, "GET", "/api/v1/list1", "test.list.read")
	CreateTestPermission(t, db, "POST", "/api/v1/list2", "test.list.create")
	CreateTestPermission(t, db, "PUT", "/api/v1/list3", "test.list.update")

	// List permissions
	perms, err := svc.ListPermissions(ctx)
	require.NoError(t, err)
	assert.Len(t, perms, 3)
}

func TestPermissionService_ListPermissionsPaged(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create multiple permissions
	for i := 0; i < 5; i++ {
		CreateTestPermission(t, db, "GET", "/api/v1/paged"+string(rune('a'+i)), "test.paged"+string(rune('a'+i))+".read")
	}

	// List with pagination
	perms, total, err := svc.ListPermissionsPaged(ctx, 1, 2)
	require.NoError(t, err)
	assert.Len(t, perms, 2)
	assert.Equal(t, int64(5), total)
}

func TestPermissionService_ListPermissionsByRoleID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create role and permissions
	testRole := CreateTestRole(t, db, "test_role", "Test Role")
	perm1 := CreateTestPermission(t, db, "GET", "/api/v1/role1", "test.role.read")
	perm2 := CreateTestPermission(t, db, "POST", "/api/v1/role2", "test.role.create")

	// Assign permissions to role
	AssignPermissionToRole(t, db, testRole.ID, perm1.ID)
	AssignPermissionToRole(t, db, testRole.ID, perm2.ID)

	// List permissions by role
	perms, err := svc.ListPermissionsByRoleID(ctx, testRole.ID)
	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

func TestPermissionService_ListPermissionsByUserID(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create user, role, and permissions
	testUser := CreateUniqueTestUser(t, db, "perm_user")
	testRole := CreateTestRole(t, db, "user_role", "User Role")
	perm1 := CreateTestPermission(t, db, "GET", "/api/v1/user1", "test.user.read")
	perm2 := CreateTestPermission(t, db, "POST", "/api/v1/user2", "test.user.create")

	// Assign permissions to role and role to user
	AssignPermissionToRole(t, db, testRole.ID, perm1.ID)
	AssignPermissionToRole(t, db, testRole.ID, perm2.ID)
	AssignRoleToUser(t, db, testUser.ID, testRole.ID)

	// List permissions by user
	perms, err := svc.ListPermissionsByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

func TestPermissionService_CheckUserHasPermission(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	permRepo := permission.NewPermissionRepository(db)
	memCache := cache.NewMemory()
	svc := permissionservice.NewPermissionService(permRepo, memCache)

	// Create user, role, and permission
	testUser := CreateUniqueTestUser(t, db, "check_user")
	testRole := CreateTestRole(t, db, "check_role", "Check Role")
	perm := CreateTestPermission(t, db, "GET", "/api/v1/check", "test.check.read")

	// Assign permission to role and role to user
	AssignPermissionToRole(t, db, testRole.ID, perm.ID)
	AssignRoleToUser(t, db, testUser.ID, testRole.ID)

	// Check user has permission
	has, err := svc.CheckUserHasPermission(ctx, testUser.ID, model.HTTPMethodGET, "/api/v1/check")
	require.NoError(t, err)
	assert.True(t, has)

	// Check user doesn't have other permission
	has, err = svc.CheckUserHasPermission(ctx, testUser.ID, model.HTTPMethodPOST, "/api/v1/other")
	require.NoError(t, err)
	assert.False(t, has)
}

// Role Service Tests

func TestRoleService_CreateRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create role
	newRole := &model.RoleModel{
		Slug:        "new_role",
		Name:        "New Role",
		Description: "A new test role",
	}

	err := svc.CreateRole(ctx, newRole)
	require.NoError(t, err)
	assert.NotZero(t, newRole.ID)

	// Verify in database
	saved, err := svc.GetRole(ctx, newRole.ID)
	require.NoError(t, err)
	assert.Equal(t, "new_role", saved.Slug)
	assert.Equal(t, "New Role", saved.Name)
}

func TestRoleService_CreateRole_DuplicateSlug(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create first role
	CreateTestRole(t, db, "dup_role", "Duplicate Role")

	// Try to create with same slug
	newRole := &model.RoleModel{
		Slug:        "dup_role",
		Name:        "Another Role",
		Description: "Should fail",
	}

	err := svc.CreateRole(ctx, newRole)
	assert.Error(t, err)
}

func TestRoleService_UpdateRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create role
	testRole := CreateTestRole(t, db, "update_role", "Update Role")

	// Update role
	testRole.Description = "Updated description"
	err := svc.UpdateRole(ctx, testRole)
	require.NoError(t, err)

	// Verify update
	saved, err := svc.GetRole(ctx, testRole.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", saved.Description)
}

func TestRoleService_DeleteRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create role
	testRole := CreateTestRole(t, db, "delete_role", "Delete Role")

	// Delete role
	err := svc.DeleteRole(ctx, testRole.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = svc.GetRole(ctx, testRole.ID)
	assert.Error(t, err)
}

func TestRoleService_AssignPermissionsToRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create role and permissions
	testRole := CreateTestRole(t, db, "assign_role", "Assign Role")
	perm1 := CreateTestPermission(t, db, "GET", "/api/v1/assign1", "test.assign.read")
	perm2 := CreateTestPermission(t, db, "POST", "/api/v1/assign2", "test.assign.create")

	// Assign permissions
	err := svc.AssignPermissionsToRole(ctx, testRole.ID, []uint64{perm1.ID, perm2.ID})
	require.NoError(t, err)

	// Verify assignment
	roleWithPerms, err := svc.GetRoleWithPermissions(ctx, testRole.ID)
	require.NoError(t, err)
	assert.Len(t, roleWithPerms.Permissions, 2)
}

func TestRoleService_AssignRolesToUser(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create user and roles
	testUser := CreateUniqueTestUser(t, db, "role_user")
	role1 := CreateTestRole(t, db, "user_role1", "User Role 1")
	role2 := CreateTestRole(t, db, "user_role2", "User Role 2")

	// Assign roles to user
	err := svc.AssignRolesToUser(ctx, testUser.ID, []uint64{role1.ID, role2.ID})
	require.NoError(t, err)

	// Verify assignment
	roles, err := svc.ListRolesByUserID(ctx, testUser.ID)
	require.NoError(t, err)
	assert.Len(t, roles, 2)
}

func TestRoleService_CheckUserHasRole(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create user and role
	testUser := CreateUniqueTestUser(t, db, "has_role_user")
	testRole := CreateTestRole(t, db, "has_role", "Has Role")

	// Assign role to user
	AssignRoleToUser(t, db, testUser.ID, testRole.ID)

	// Check user has role
	has, err := svc.CheckUserHasRole(ctx, testUser.ID, "has_role")
	require.NoError(t, err)
	assert.True(t, has)

	// Check user doesn't have other role
	has, err = svc.CheckUserHasRole(ctx, testUser.ID, "other_role")
	require.NoError(t, err)
	assert.False(t, has)
}

func TestRoleService_ListRoles(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create multiple roles
	CreateTestRole(t, db, "list_role1", "List Role 1")
	CreateTestRole(t, db, "list_role2", "List Role 2")
	CreateTestRole(t, db, "list_role3", "List Role 3")

	// List roles
	roles, err := svc.ListRoles(ctx)
	require.NoError(t, err)
	assert.Len(t, roles, 3)
}

func TestRoleService_GetEffectivePermissions(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Setup
	roleRepo := admin.NewRoleRepository(db)
	memCache := cache.NewMemory()
	svc := role.NewRoleService(roleRepo, memCache)

	// Create parent role with permissions
	parentRole := CreateTestRole(t, db, "parent_role", "Parent Role")
	parentPerm := CreateTestPermission(t, db, "GET", "/api/v1/parent", "test.parent.read")
	AssignPermissionToRole(t, db, parentRole.ID, parentPerm.ID)

	// Create child role with its own permissions
	childRole := &model.RoleModel{
		Base: model.Base{
			ExtJSON: "{}",
		},
		Slug:     "child_role",
		Name:     "Child Role",
		ParentID: &parentRole.ID,
	}
	err := db.Create(childRole).Error
	require.NoError(t, err)

	childPerm := CreateTestPermission(t, db, "POST", "/api/v1/child", "test.child.create")
	AssignPermissionToRole(t, db, childRole.ID, childPerm.ID)

	// Get effective permissions (should include both parent and child permissions)
	perms, err := svc.GetEffectivePermissions(ctx, childRole.ID)
	require.NoError(t, err)
	assert.Len(t, perms, 2)
}

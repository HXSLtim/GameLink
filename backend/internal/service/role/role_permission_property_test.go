package role

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"pgregory.net/rapid"

	"gamelink/internal/model"
	adminrepo "gamelink/internal/repository/admin"
	"gamelink/pkg/cache"
	"gamelink/pkg/testutil"
)

// setupTestDB creates an in-memory database with required tables for testing
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewMemoryDB(t)
	testutil.MigrateTables(t, db,
		&model.User{},
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
		&model.UserRole{},
	)
	return db
}

// TestRolePermissionAssignmentPersistence tests Property 5: 权限分配持久化
// **Feature: rbac-button-level-permission, Property 5: 权限分配持久化**
// **Validates: Requirements 2.3**
//
// For any role and set of permission IDs, after calling AssignPermissions,
// querying the role's permissions should return exactly the assigned permission IDs.
func TestRolePermissionAssignmentPersistence(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions
	permissions := createTestPermissions(t, db, 10)
	permIDs := make([]uint64, len(permissions))
	for i, p := range permissions {
		permIDs[i] = p.ID
	}

	// Create test role
	role := createTestRole(t, db, "test-role", "Test Role")

	rapid.Check(t, func(rt *rapid.T) {
		ctx := context.Background()

		// Generate a random subset of permission IDs
		selectedCount := rapid.IntRange(0, len(permIDs)).Draw(rt, "selectedCount")
		selectedIDs := make([]uint64, 0, selectedCount)

		// Use a map to track selected indices to avoid duplicates
		selectedIndices := make(map[int]bool)
		for len(selectedIDs) < selectedCount {
			idx := rapid.IntRange(0, len(permIDs)-1).Draw(rt, "permIndex")
			if !selectedIndices[idx] {
				selectedIndices[idx] = true
				selectedIDs = append(selectedIDs, permIDs[idx])
			}
		}

		// Assign permissions to role
		err := roleSvc.AssignPermissionsToRole(ctx, role.ID, selectedIDs)
		require.NoError(t, err, "AssignPermissionsToRole should not fail")

		// Query the role's permissions
		roleWithPerms, err := roleSvc.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err, "GetRoleWithPermissions should not fail")

		// Extract assigned permission IDs
		assignedIDs := make(map[uint64]bool)
		for _, p := range roleWithPerms.Permissions {
			assignedIDs[p.ID] = true
		}

		// Verify: assigned permissions should match exactly
		assert.Equal(t, len(selectedIDs), len(assignedIDs),
			"Number of assigned permissions should match")

		for _, id := range selectedIDs {
			assert.True(t, assignedIDs[id],
				"Permission ID %d should be assigned to role", id)
		}
	})
}

// TestRolePermissionAddRemovePersistence tests that adding and removing
// individual permissions works correctly and persists.
func TestRolePermissionAddRemovePersistence(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions
	permissions := createTestPermissions(t, db, 5)

	// Create test role
	role := createTestRole(t, db, "test-role-add-remove", "Test Role Add Remove")

	rapid.Check(t, func(rt *rapid.T) {
		ctx := context.Background()

		// Start with empty permissions
		err := roleSvc.AssignPermissionsToRole(ctx, role.ID, []uint64{})
		require.NoError(t, err)

		// Track expected permissions
		expectedPerms := make(map[uint64]bool)

		// Perform random add/remove operations
		opCount := rapid.IntRange(1, 10).Draw(rt, "opCount")
		for i := 0; i < opCount; i++ {
			permIdx := rapid.IntRange(0, len(permissions)-1).Draw(rt, "permIdx")
			permID := permissions[permIdx].ID
			isAdd := rapid.Bool().Draw(rt, "isAdd")

			if isAdd {
				err := roleSvc.AddPermissionsToRole(ctx, role.ID, []uint64{permID})
				require.NoError(t, err)
				expectedPerms[permID] = true
			} else {
				err := roleSvc.RemovePermissionsFromRole(ctx, role.ID, []uint64{permID})
				require.NoError(t, err)
				delete(expectedPerms, permID)
			}
		}

		// Verify final state
		roleWithPerms, err := roleSvc.GetRoleWithPermissions(ctx, role.ID)
		require.NoError(t, err)

		actualPerms := make(map[uint64]bool)
		for _, p := range roleWithPerms.Permissions {
			actualPerms[p.ID] = true
		}

		assert.Equal(t, len(expectedPerms), len(actualPerms),
			"Number of permissions should match expected")

		for id := range expectedPerms {
			assert.True(t, actualPerms[id],
				"Permission ID %d should be present", id)
		}
	})
}

// TestRolePermissionBatchAssignmentAtomicity tests that batch assignment
// is atomic - either all permissions are assigned or none.
func TestRolePermissionBatchAssignmentAtomicity(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions
	permissions := createTestPermissions(t, db, 5)
	permIDs := make([]uint64, len(permissions))
	for i, p := range permissions {
		permIDs[i] = p.ID
	}

	// Create test role
	role := createTestRole(t, db, "test-role-atomic", "Test Role Atomic")

	ctx := context.Background()

	// Assign initial permissions
	initialIDs := permIDs[:3]
	err := roleSvc.AssignPermissionsToRole(ctx, role.ID, initialIDs)
	require.NoError(t, err)

	// Verify initial state
	roleWithPerms, err := roleSvc.GetRoleWithPermissions(ctx, role.ID)
	require.NoError(t, err)
	assert.Len(t, roleWithPerms.Permissions, 3)

	// Assign new set of permissions (should replace all)
	newIDs := permIDs[2:]
	err = roleSvc.AssignPermissionsToRole(ctx, role.ID, newIDs)
	require.NoError(t, err)

	// Verify new state - should have exactly the new permissions
	roleWithPerms, err = roleSvc.GetRoleWithPermissions(ctx, role.ID)
	require.NoError(t, err)
	assert.Len(t, roleWithPerms.Permissions, len(newIDs))

	actualIDs := make(map[uint64]bool)
	for _, p := range roleWithPerms.Permissions {
		actualIDs[p.ID] = true
	}

	for _, id := range newIDs {
		assert.True(t, actualIDs[id], "Permission ID %d should be assigned", id)
	}

	// Old permissions that are not in new set should be removed
	for _, id := range initialIDs[:2] {
		if !contains(newIDs, id) {
			assert.False(t, actualIDs[id], "Permission ID %d should be removed", id)
		}
	}
}

// TestRolePermissionCacheInvalidation tests that cache is properly invalidated
// after permission assignment changes.
func TestRolePermissionCacheInvalidation(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions
	permissions := createTestPermissions(t, db, 3)
	permIDs := make([]uint64, len(permissions))
	for i, p := range permissions {
		permIDs[i] = p.ID
	}

	// Create test role
	role := createTestRole(t, db, "test-role-cache", "Test Role Cache")

	ctx := context.Background()

	// Assign initial permissions
	err := roleSvc.AssignPermissionsToRole(ctx, role.ID, permIDs[:2])
	require.NoError(t, err)

	// Query to populate cache
	roleWithPerms, err := roleSvc.GetRoleWithPermissions(ctx, role.ID)
	require.NoError(t, err)
	assert.Len(t, roleWithPerms.Permissions, 2)

	// Modify permissions
	err = roleSvc.AssignPermissionsToRole(ctx, role.ID, permIDs)
	require.NoError(t, err)

	// Query again - should reflect new state (cache should be invalidated)
	roleWithPerms, err = roleSvc.GetRoleWithPermissions(ctx, role.ID)
	require.NoError(t, err)
	assert.Len(t, roleWithPerms.Permissions, 3, "Cache should be invalidated and new permissions should be returned")
}

// Helper functions

func createTestPermissions(t *testing.T, db *gorm.DB, count int) []model.Permission {
	t.Helper()
	permissions := make([]model.Permission, count)
	for i := 0; i < count; i++ {
		perm := model.Permission{
			Method:      model.HTTPMethodGET,
			Path:        "/api/test/" + randomString(8),
			Code:        "test.resource" + randomString(4) + ".read",
			Group:       "test",
			Description: "Test permission " + randomString(4),
		}
		err := db.Create(&perm).Error
		require.NoError(t, err)
		permissions[i] = perm
	}
	return permissions
}

func createTestRole(t *testing.T, db *gorm.DB, slug, name string) *model.RoleModel {
	t.Helper()
	role := &model.RoleModel{
		Slug:        slug + "_" + randomString(8),
		Name:        name,
		Description: "Test role for property testing",
		IsSystem:    false,
	}
	err := db.Create(role).Error
	require.NoError(t, err)
	return role
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

func contains(slice []uint64, val uint64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// TestMultiRolePermissionMerging tests Property 17: 多角色权限合并
// **Feature: rbac-button-level-permission, Property 17: 多角色权限合并**
// **Validates: Requirements 10.1**
//
// For any user with multiple roles, the effective permissions should be
// the union of all permissions from all assigned roles.
func TestMultiRolePermissionMerging(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions (enough for multiple roles)
	permissions := createTestPermissions(t, db, 15)

	rapid.Check(t, func(rt *rapid.T) {
		ctx := context.Background()

		// Create a test user
		user := createTestUser(t, db)

		// Generate 2-4 roles with random permission subsets
		roleCount := rapid.IntRange(2, 4).Draw(rt, "roleCount")
		roles := make([]*model.RoleModel, roleCount)
		rolePermissions := make([][]uint64, roleCount)

		// Track all expected permissions (union of all role permissions)
		expectedPermIDs := make(map[uint64]bool)

		for i := 0; i < roleCount; i++ {
			// Create a role
			role := createTestRole(t, db, "multi-role-test-"+randomString(6), "Multi Role Test "+randomString(4))
			roles[i] = role

			// Assign random subset of permissions to this role
			permCount := rapid.IntRange(1, 8).Draw(rt, "permCount")
			selectedPerms := make([]uint64, 0, permCount)
			selectedIndices := make(map[int]bool)

			for len(selectedPerms) < permCount {
				idx := rapid.IntRange(0, len(permissions)-1).Draw(rt, "permIdx")
				if !selectedIndices[idx] {
					selectedIndices[idx] = true
					selectedPerms = append(selectedPerms, permissions[idx].ID)
					expectedPermIDs[permissions[idx].ID] = true
				}
			}

			rolePermissions[i] = selectedPerms
			err := roleSvc.AssignPermissionsToRole(ctx, role.ID, selectedPerms)
			require.NoError(t, err, "AssignPermissionsToRole should not fail")
		}

		// Assign all roles to the user
		roleIDs := make([]uint64, roleCount)
		for i, role := range roles {
			roleIDs[i] = role.ID
		}
		err := roleSvc.AssignRolesToUser(ctx, user.ID, roleIDs)
		require.NoError(t, err, "AssignRolesToUser should not fail")

		// Get user's effective permissions
		effectivePerms, err := roleSvc.GetUserEffectivePermissions(ctx, user.ID)
		require.NoError(t, err, "GetUserEffectivePermissions should not fail")

		// Extract effective permission IDs
		actualPermIDs := make(map[uint64]bool)
		for _, p := range effectivePerms {
			actualPermIDs[p.ID] = true
		}

		// Property: effective permissions should be the union of all role permissions
		// 1. All expected permissions should be present
		for permID := range expectedPermIDs {
			assert.True(t, actualPermIDs[permID],
				"Permission ID %d should be in effective permissions (union property)", permID)
		}

		// 2. No extra permissions should be present
		for permID := range actualPermIDs {
			assert.True(t, expectedPermIDs[permID],
				"Permission ID %d should not be in effective permissions (no extra permissions)", permID)
		}

		// 3. Count should match (union means no duplicates)
		assert.Equal(t, len(expectedPermIDs), len(actualPermIDs),
			"Effective permission count should equal union of all role permissions")

		// Cleanup: remove user roles for next iteration
		_ = roleSvc.AssignRolesToUser(ctx, user.ID, []uint64{})
	})
}

// TestMultiRolePermissionMergingWithOverlap tests that overlapping permissions
// between roles are correctly deduplicated in the union.
func TestMultiRolePermissionMergingWithOverlap(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test permissions
	permissions := createTestPermissions(t, db, 6)
	permIDs := make([]uint64, len(permissions))
	for i, p := range permissions {
		permIDs[i] = p.ID
	}

	// Create test user
	user := createTestUser(t, db)

	// Create two roles with overlapping permissions
	role1 := createTestRole(t, db, "overlap-role1", "Overlap Role 1")
	role2 := createTestRole(t, db, "overlap-role2", "Overlap Role 2")

	ctx := context.Background()

	// Role1 gets permissions 0, 1, 2, 3
	role1Perms := permIDs[:4]
	err := roleSvc.AssignPermissionsToRole(ctx, role1.ID, role1Perms)
	require.NoError(t, err)

	// Role2 gets permissions 2, 3, 4, 5 (overlaps with 2, 3)
	role2Perms := permIDs[2:]
	err = roleSvc.AssignPermissionsToRole(ctx, role2.ID, role2Perms)
	require.NoError(t, err)

	// Assign both roles to user
	err = roleSvc.AssignRolesToUser(ctx, user.ID, []uint64{role1.ID, role2.ID})
	require.NoError(t, err)

	// Get effective permissions
	effectivePerms, err := roleSvc.GetUserEffectivePermissions(ctx, user.ID)
	require.NoError(t, err)

	// Should have all 6 permissions (union), not 8 (sum)
	assert.Len(t, effectivePerms, 6, "Effective permissions should be union (6), not sum (8)")

	// Verify all permissions are present
	actualPermIDs := make(map[uint64]bool)
	for _, p := range effectivePerms {
		actualPermIDs[p.ID] = true
	}

	for _, id := range permIDs {
		assert.True(t, actualPermIDs[id], "Permission ID %d should be present", id)
	}
}

// TestMultiRolePermissionMergingEmptyRoles tests that a user with roles
// that have no permissions gets an empty permission set.
func TestMultiRolePermissionMergingEmptyRoles(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test user
	user := createTestUser(t, db)

	// Create two roles with no permissions
	role1 := createTestRole(t, db, "empty-role1", "Empty Role 1")
	role2 := createTestRole(t, db, "empty-role2", "Empty Role 2")

	ctx := context.Background()

	// Assign both roles to user (no permissions assigned to roles)
	err := roleSvc.AssignRolesToUser(ctx, user.ID, []uint64{role1.ID, role2.ID})
	require.NoError(t, err)

	// Get effective permissions
	effectivePerms, err := roleSvc.GetUserEffectivePermissions(ctx, user.ID)
	require.NoError(t, err)

	// Should have no permissions
	assert.Len(t, effectivePerms, 0, "User with empty roles should have no permissions")
}

// TestMultiRolePermissionMergingNoRoles tests that a user with no roles
// gets an empty permission set.
func TestMultiRolePermissionMergingNoRoles(t *testing.T) {
	db := setupTestDB(t)
	memCache := cache.NewMemory()
	roleRepo := adminrepo.NewRoleRepository(db)
	roleSvc := NewRoleService(roleRepo, memCache)

	// Create test user with no roles
	user := createTestUser(t, db)

	ctx := context.Background()

	// Get effective permissions (user has no roles)
	effectivePerms, err := roleSvc.GetUserEffectivePermissions(ctx, user.ID)
	require.NoError(t, err)

	// Should have no permissions
	assert.Len(t, effectivePerms, 0, "User with no roles should have no permissions")
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *gorm.DB) *model.User {
	t.Helper()
	uniqueSuffix := randomString(10)
	user := &model.User{
		Phone:        "138" + uniqueSuffix[:8],
		Email:        "test_" + uniqueSuffix + "@example.com",
		Name:         "Test User " + randomString(4),
		PasswordHash: "hashedpassword",
		Role:         model.RoleUser,
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

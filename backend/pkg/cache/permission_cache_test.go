package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPermissionCache(t *testing.T) {
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	assert.NotNil(t, pc)
	assert.NotNil(t, pc.cache)
	assert.Greater(t, pc.version, int64(0))
}

func TestPermissionCache_UserPermissions(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	userID := uint64(123)
	codes := []string{"admin.users.read", "admin.users.create", "admin.orders.read"}

	// Test cache miss
	result, ok, err := pc.GetUserPermissions(ctx, userID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Set permissions
	err = pc.SetUserPermissions(ctx, userID, codes)
	require.NoError(t, err)

	// Test cache hit
	result, ok, err = pc.GetUserPermissions(ctx, userID)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, codes, result)
}

func TestPermissionCache_RolePermissions(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	roleID := uint64(456)
	codes := []string{"admin.roles.read", "admin.roles.update"}

	// Test cache miss
	result, ok, err := pc.GetRolePermissions(ctx, roleID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Set permissions
	err = pc.SetRolePermissions(ctx, roleID, codes)
	require.NoError(t, err)

	// Test cache hit
	result, ok, err = pc.GetRolePermissions(ctx, roleID)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, codes, result)
}

func TestPermissionCache_UserRoles(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	userID := uint64(789)
	roleIDs := []uint64{1, 2, 3}

	// Test cache miss
	result, ok, err := pc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Set roles
	err = pc.SetUserRoles(ctx, userID, roleIDs)
	require.NoError(t, err)

	// Test cache hit
	result, ok, err = pc.GetUserRoles(ctx, userID)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, roleIDs, result)
}

func TestPermissionCache_PermissionTree(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	tree := []PermissionTreeNode{
		{
			ID:          1,
			Code:        "admin",
			Description: "Admin module",
			Group:       "admin",
			Children: []PermissionTreeNode{
				{ID: 2, Code: "admin.users", Description: "User management", Group: "admin"},
			},
		},
	}

	// Test cache miss
	result, ok, err := pc.GetPermissionTree(ctx)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Set tree
	err = pc.SetPermissionTree(ctx, tree)
	require.NoError(t, err)

	// Test cache hit
	result, ok, err = pc.GetPermissionTree(ctx)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, tree, result)
}

func TestPermissionCache_PermissionGroups(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	groups := []string{"admin", "player", "user", "system"}

	// Test cache miss
	result, ok, err := pc.GetPermissionGroups(ctx)
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, result)

	// Set groups
	err = pc.SetPermissionGroups(ctx, groups)
	require.NoError(t, err)

	// Test cache hit
	result, ok, err = pc.GetPermissionGroups(ctx)
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, groups, result)
}

func TestPermissionCache_InvalidateUserCache(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	userID := uint64(123)
	codes := []string{"admin.users.read"}
	roleIDs := []uint64{1, 2}

	// Set both permissions and roles
	err := pc.SetUserPermissions(ctx, userID, codes)
	require.NoError(t, err)
	err = pc.SetUserRoles(ctx, userID, roleIDs)
	require.NoError(t, err)

	// Verify they exist
	_, ok, _ := pc.GetUserPermissions(ctx, userID)
	assert.True(t, ok)
	_, ok, _ = pc.GetUserRoles(ctx, userID)
	assert.True(t, ok)

	// Invalidate
	err = pc.InvalidateUserCache(ctx, userID)
	require.NoError(t, err)

	// Verify they are gone
	_, ok, _ = pc.GetUserPermissions(ctx, userID)
	assert.False(t, ok)
	_, ok, _ = pc.GetUserRoles(ctx, userID)
	assert.False(t, ok)
}

func TestPermissionCache_InvalidateRoleCache(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	roleID := uint64(456)
	codes := []string{"admin.roles.read"}

	// Set permissions
	err := pc.SetRolePermissions(ctx, roleID, codes)
	require.NoError(t, err)

	// Verify it exists
	_, ok, _ := pc.GetRolePermissions(ctx, roleID)
	assert.True(t, ok)

	// Invalidate
	err = pc.InvalidateRoleCache(ctx, roleID)
	require.NoError(t, err)

	// Verify it's gone
	_, ok, _ = pc.GetRolePermissions(ctx, roleID)
	assert.False(t, ok)
}

func TestPermissionCache_InvalidateAllPermissionCaches(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	// Set various caches
	err := pc.SetUserPermissions(ctx, 1, []string{"perm1"})
	require.NoError(t, err)
	err = pc.SetRolePermissions(ctx, 1, []string{"perm2"})
	require.NoError(t, err)
	err = pc.SetPermissionTree(ctx, []PermissionTreeNode{{ID: 1}})
	require.NoError(t, err)

	// Verify they exist
	_, ok, _ := pc.GetUserPermissions(ctx, 1)
	assert.True(t, ok)
	_, ok, _ = pc.GetRolePermissions(ctx, 1)
	assert.True(t, ok)
	_, ok, _ = pc.GetPermissionTree(ctx)
	assert.True(t, ok)

	// Record the old version
	oldVersion := pc.GetVersion()

	// Invalidate all by incrementing version
	pc.InvalidateAllPermissionCaches()

	// Version should have changed
	newVersion := pc.GetVersion()
	assert.NotEqual(t, oldVersion, newVersion, "Version should change after invalidation")

	// All caches should now miss because version changed (keys are different)
	// The old data still exists in the underlying cache but with old version keys
	_, ok, _ = pc.GetUserPermissions(ctx, 1)
	assert.False(t, ok, "User permissions should miss after version change")
	_, ok, _ = pc.GetRolePermissions(ctx, 1)
	assert.False(t, ok, "Role permissions should miss after version change")
	_, ok, _ = pc.GetPermissionTree(ctx)
	assert.False(t, ok, "Permission tree should miss after version change")
}

func TestGetCacheTTLWithJitter(t *testing.T) {
	baseTTL := 30 * time.Minute

	// Run multiple times to verify jitter is applied
	results := make(map[time.Duration]bool)
	for i := 0; i < 100; i++ {
		ttl := GetCacheTTLWithJitter(baseTTL)
		results[ttl] = true

		// Verify TTL is within expected range (±10%)
		minTTL := baseTTL - (baseTTL * TTLJitterPercent / 100)
		maxTTL := baseTTL + (baseTTL * TTLJitterPercent / 100)
		assert.GreaterOrEqual(t, ttl, minTTL, "TTL should be >= min")
		assert.LessOrEqual(t, ttl, maxTTL, "TTL should be <= max")
	}

	// Verify we got different values (jitter is working)
	assert.Greater(t, len(results), 1, "Should have different TTL values due to jitter")
}

func TestPermissionCache_GetVersion(t *testing.T) {
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	v1 := pc.GetVersion()
	assert.Greater(t, v1, int64(0))

	// Version should change after invalidation
	// Add a small delay to ensure time.Now().UnixNano() returns a different value
	time.Sleep(time.Millisecond)
	pc.InvalidateAllPermissionCaches()
	v2 := pc.GetVersion()
	assert.NotEqual(t, v1, v2, "Version should change after InvalidateAllPermissionCaches")
}

func TestPermissionCache_GetUnderlyingCache(t *testing.T) {
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	underlying := pc.GetUnderlyingCache()
	assert.Equal(t, memCache, underlying)
}

// Mock implementation of UserRoleProvider for testing
type mockUserRoleProvider struct {
	userIDs []uint64
	err     error
}

func (m *mockUserRoleProvider) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	return m.userIDs, m.err
}

func TestPermissionCache_InvalidateRolePermissionsAndPropagateToUsers(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	roleID := uint64(1)
	userIDs := []uint64{100, 200, 300}

	// Set up caches for role and users
	err := pc.SetRolePermissions(ctx, roleID, []string{"perm1"})
	require.NoError(t, err)
	for _, userID := range userIDs {
		err = pc.SetUserPermissions(ctx, userID, []string{"perm1"})
		require.NoError(t, err)
	}

	// Verify all caches exist
	_, ok, _ := pc.GetRolePermissions(ctx, roleID)
	assert.True(t, ok)
	for _, userID := range userIDs {
		_, ok, _ = pc.GetUserPermissions(ctx, userID)
		assert.True(t, ok)
	}

	// Create mock provider
	provider := &mockUserRoleProvider{userIDs: userIDs}

	// Invalidate and propagate
	err = pc.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID, provider)
	require.NoError(t, err)

	// Verify all caches are invalidated
	_, ok, _ = pc.GetRolePermissions(ctx, roleID)
	assert.False(t, ok)
	for _, userID := range userIDs {
		_, ok, _ = pc.GetUserPermissions(ctx, userID)
		assert.False(t, ok)
	}
}

func TestPermissionCache_InvalidateMultipleUsers(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	userIDs := []uint64{1, 2, 3}

	// Set up caches
	for _, userID := range userIDs {
		err := pc.SetUserPermissions(ctx, userID, []string{"perm1"})
		require.NoError(t, err)
	}

	// Verify all exist
	for _, userID := range userIDs {
		_, ok, _ := pc.GetUserPermissions(ctx, userID)
		assert.True(t, ok)
	}

	// Invalidate all
	err := pc.InvalidateMultipleUsers(ctx, userIDs)
	require.NoError(t, err)

	// Verify all are gone
	for _, userID := range userIDs {
		_, ok, _ := pc.GetUserPermissions(ctx, userID)
		assert.False(t, ok)
	}
}

func TestPermissionCache_InvalidateMultipleRoles(t *testing.T) {
	ctx := context.Background()
	memCache := NewMemory()
	pc := NewPermissionCache(memCache)

	roleIDs := []uint64{1, 2, 3}

	// Set up caches
	for _, roleID := range roleIDs {
		err := pc.SetRolePermissions(ctx, roleID, []string{"perm1"})
		require.NoError(t, err)
	}

	// Verify all exist
	for _, roleID := range roleIDs {
		_, ok, _ := pc.GetRolePermissions(ctx, roleID)
		assert.True(t, ok)
	}

	// Invalidate all
	err := pc.InvalidateMultipleRoles(ctx, roleIDs)
	require.NoError(t, err)

	// Verify all are gone
	for _, roleID := range roleIDs {
		_, ok, _ := pc.GetRolePermissions(ctx, roleID)
		assert.False(t, ok)
	}
}

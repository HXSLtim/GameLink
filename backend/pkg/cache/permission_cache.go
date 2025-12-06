package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// PermissionCache provides specialized caching for RBAC permission data.
// It implements versioned cache keys and hash-like structures for efficient
// permission storage and invalidation.
type PermissionCache struct {
	cache   Cache
	version int64
	mu      sync.RWMutex
}

// Cache key patterns with version support for quick invalidation
const (
	// User permission cache key pattern: perm:user:{userID}:v{version}
	UserPermissionKeyPattern = "perm:user:%d:v%d"
	// Role permission cache key pattern: perm:role:{roleID}:v{version}
	RolePermissionKeyPattern = "perm:role:%d:v%d"
	// Permission tree cache key: perm:tree:v{version}
	PermissionTreeKeyPattern = "perm:tree:v%d"
	// Permission groups cache key: perm:groups:v{version}
	PermissionGroupsKeyPattern = "perm:groups:v%d"
	// All permissions cache key: perm:all:v{version}
	AllPermissionsKeyPattern = "perm:all:v%d"
	// User roles cache key pattern: perm:user_roles:{userID}:v{version}
	UserRolesKeyPattern = "perm:user_roles:%d:v%d"
	// Role users cache key pattern: perm:role_users:{roleID}:v{version}
	RoleUsersKeyPattern = "perm:role_users:%d:v%d"
)

// Default TTL values
const (
	DefaultPermissionCacheTTL = 30 * time.Minute
	DefaultRoleCacheTTL       = 30 * time.Minute
	DefaultTreeCacheTTL       = 60 * time.Minute
	// TTL jitter percentage (10%)
	TTLJitterPercent = 10
)

// NewPermissionCache creates a new permission cache instance.
func NewPermissionCache(cache Cache) *PermissionCache {
	return &PermissionCache{
		cache:   cache,
		version: time.Now().UnixNano(),
	}
}

// GetCacheTTLWithJitter returns a TTL with random jitter to prevent cache stampede.
// The jitter is ±10% of the base TTL.
func GetCacheTTLWithJitter(baseTTL time.Duration) time.Duration {
	jitterRange := int64(baseTTL) * TTLJitterPercent / 100
	jitter := rand.Int63n(jitterRange*2) - jitterRange // Random value between -jitterRange and +jitterRange
	return baseTTL + time.Duration(jitter)
}

// getVersion returns the current cache version.
func (pc *PermissionCache) getVersion() int64 {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.version
}

// incrementVersion increments the cache version to invalidate all cached data.
func (pc *PermissionCache) incrementVersion() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	// Use a combination of current version + 1 and current time to ensure uniqueness
	// This guarantees the version always changes even if called in quick succession
	pc.version = pc.version + 1 + (time.Now().UnixNano() % 1000000)
}

// userPermissionKey generates the cache key for user permissions.
func (pc *PermissionCache) userPermissionKey(userID uint64) string {
	return fmt.Sprintf(UserPermissionKeyPattern, userID, pc.getVersion())
}

// rolePermissionKey generates the cache key for role permissions.
func (pc *PermissionCache) rolePermissionKey(roleID uint64) string {
	return fmt.Sprintf(RolePermissionKeyPattern, roleID, pc.getVersion())
}

// permissionTreeKey generates the cache key for permission tree.
func (pc *PermissionCache) permissionTreeKey() string {
	return fmt.Sprintf(PermissionTreeKeyPattern, pc.getVersion())
}

// permissionGroupsKey generates the cache key for permission groups.
func (pc *PermissionCache) permissionGroupsKey() string {
	return fmt.Sprintf(PermissionGroupsKeyPattern, pc.getVersion())
}

// allPermissionsKey generates the cache key for all permissions.
func (pc *PermissionCache) allPermissionsKey() string {
	return fmt.Sprintf(AllPermissionsKeyPattern, pc.getVersion())
}

// userRolesKey generates the cache key for user roles.
func (pc *PermissionCache) userRolesKey(userID uint64) string {
	return fmt.Sprintf(UserRolesKeyPattern, userID, pc.getVersion())
}

// roleUsersKey generates the cache key for role users.
func (pc *PermissionCache) roleUsersKey(roleID uint64) string {
	return fmt.Sprintf(RoleUsersKeyPattern, roleID, pc.getVersion())
}

// PermissionCacheData represents cached permission data with metadata.
type PermissionCacheData struct {
	Codes     []string  `json:"codes"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RoleCacheData represents cached role data with metadata.
type RoleCacheData struct {
	RoleIDs   []uint64  `json:"roleIds"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetUserPermissions retrieves user permissions from cache.
func (pc *PermissionCache) GetUserPermissions(ctx context.Context, userID uint64) ([]string, bool, error) {
	key := pc.userPermissionKey(userID)
	value, ok, err := pc.cache.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}

	var data PermissionCacheData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, false, nil // Treat unmarshal error as cache miss
	}

	return data.Codes, true, nil
}

// SetUserPermissions stores user permissions in cache.
func (pc *PermissionCache) SetUserPermissions(ctx context.Context, userID uint64, codes []string) error {
	key := pc.userPermissionKey(userID)
	data := PermissionCacheData{
		Codes:     codes,
		UpdatedAt: time.Now(),
	}

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := GetCacheTTLWithJitter(DefaultPermissionCacheTTL)
	return pc.cache.Set(ctx, key, string(value), ttl)
}

// GetRolePermissions retrieves role permissions from cache.
func (pc *PermissionCache) GetRolePermissions(ctx context.Context, roleID uint64) ([]string, bool, error) {
	key := pc.rolePermissionKey(roleID)
	value, ok, err := pc.cache.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}

	var data PermissionCacheData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, false, nil
	}

	return data.Codes, true, nil
}

// SetRolePermissions stores role permissions in cache.
func (pc *PermissionCache) SetRolePermissions(ctx context.Context, roleID uint64, codes []string) error {
	key := pc.rolePermissionKey(roleID)
	data := PermissionCacheData{
		Codes:     codes,
		UpdatedAt: time.Now(),
	}

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := GetCacheTTLWithJitter(DefaultRoleCacheTTL)
	return pc.cache.Set(ctx, key, string(value), ttl)
}

// GetUserRoles retrieves user roles from cache.
func (pc *PermissionCache) GetUserRoles(ctx context.Context, userID uint64) ([]uint64, bool, error) {
	key := pc.userRolesKey(userID)
	value, ok, err := pc.cache.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}

	var data RoleCacheData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return nil, false, nil
	}

	return data.RoleIDs, true, nil
}

// SetUserRoles stores user roles in cache.
func (pc *PermissionCache) SetUserRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	key := pc.userRolesKey(userID)
	data := RoleCacheData{
		RoleIDs:   roleIDs,
		UpdatedAt: time.Now(),
	}

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := GetCacheTTLWithJitter(DefaultRoleCacheTTL)
	return pc.cache.Set(ctx, key, string(value), ttl)
}

// PermissionTreeNode represents a node in the permission tree.
type PermissionTreeNode struct {
	ID          uint64               `json:"id"`
	Code        string               `json:"code"`
	Description string               `json:"description"`
	Group       string               `json:"group"`
	Children    []PermissionTreeNode `json:"children,omitempty"`
}

// GetPermissionTree retrieves the permission tree from cache.
func (pc *PermissionCache) GetPermissionTree(ctx context.Context) ([]PermissionTreeNode, bool, error) {
	key := pc.permissionTreeKey()
	value, ok, err := pc.cache.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}

	var tree []PermissionTreeNode
	if err := json.Unmarshal([]byte(value), &tree); err != nil {
		return nil, false, nil
	}

	return tree, true, nil
}

// SetPermissionTree stores the permission tree in cache.
func (pc *PermissionCache) SetPermissionTree(ctx context.Context, tree []PermissionTreeNode) error {
	key := pc.permissionTreeKey()
	value, err := json.Marshal(tree)
	if err != nil {
		return err
	}

	ttl := GetCacheTTLWithJitter(DefaultTreeCacheTTL)
	return pc.cache.Set(ctx, key, string(value), ttl)
}

// GetPermissionGroups retrieves permission groups from cache.
func (pc *PermissionCache) GetPermissionGroups(ctx context.Context) ([]string, bool, error) {
	key := pc.permissionGroupsKey()
	value, ok, err := pc.cache.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}

	var groups []string
	if err := json.Unmarshal([]byte(value), &groups); err != nil {
		return nil, false, nil
	}

	return groups, true, nil
}

// SetPermissionGroups stores permission groups in cache.
func (pc *PermissionCache) SetPermissionGroups(ctx context.Context, groups []string) error {
	key := pc.permissionGroupsKey()
	value, err := json.Marshal(groups)
	if err != nil {
		return err
	}

	ttl := GetCacheTTLWithJitter(DefaultTreeCacheTTL)
	return pc.cache.Set(ctx, key, string(value), ttl)
}

// InvalidateUserCache invalidates the cache for a specific user.
// This should be called when user roles change.
func (pc *PermissionCache) InvalidateUserCache(ctx context.Context, userID uint64) error {
	// Delete user permissions cache
	permKey := pc.userPermissionKey(userID)
	if err := pc.cache.Delete(ctx, permKey); err != nil {
		return err
	}

	// Delete user roles cache
	rolesKey := pc.userRolesKey(userID)
	return pc.cache.Delete(ctx, rolesKey)
}

// InvalidateRoleCache invalidates the cache for a specific role.
// This should be called when role permissions change.
func (pc *PermissionCache) InvalidateRoleCache(ctx context.Context, roleID uint64) error {
	key := pc.rolePermissionKey(roleID)
	return pc.cache.Delete(ctx, key)
}

// InvalidateAllPermissionCaches invalidates all permission-related caches.
// This is done by incrementing the version number, which makes all existing
// cache keys obsolete without needing to delete them individually.
func (pc *PermissionCache) InvalidateAllPermissionCaches() {
	pc.incrementVersion()
}

// InvalidatePermissionTree invalidates the permission tree cache.
func (pc *PermissionCache) InvalidatePermissionTree(ctx context.Context) error {
	key := pc.permissionTreeKey()
	return pc.cache.Delete(ctx, key)
}

// InvalidatePermissionGroups invalidates the permission groups cache.
func (pc *PermissionCache) InvalidatePermissionGroups(ctx context.Context) error {
	key := pc.permissionGroupsKey()
	return pc.cache.Delete(ctx, key)
}

// GetVersion returns the current cache version for debugging/monitoring.
func (pc *PermissionCache) GetVersion() int64 {
	return pc.getVersion()
}

// GetUnderlyingCache returns the underlying cache instance.
// This is useful for advanced operations or testing.
func (pc *PermissionCache) GetUnderlyingCache() Cache {
	return pc.cache
}

// WarmupConfig contains configuration for cache warmup.
type WarmupConfig struct {
	// SystemRoleSlugs are the role slugs to preload during warmup
	SystemRoleSlugs []string
	// WarmupPermissionTree indicates whether to warmup the permission tree
	WarmupPermissionTree bool
	// WarmupPermissionGroups indicates whether to warmup permission groups
	WarmupPermissionGroups bool
}

// DefaultWarmupConfig returns the default warmup configuration.
func DefaultWarmupConfig() WarmupConfig {
	return WarmupConfig{
		SystemRoleSlugs:        []string{"superAdmin", "admin"},
		WarmupPermissionTree:   true,
		WarmupPermissionGroups: true,
	}
}

// WarmupResult contains the results of a cache warmup operation.
type WarmupResult struct {
	PermissionTreeWarmed   bool
	PermissionGroupsWarmed bool
	RolesWarmed            []string
	Errors                 []error
}

// PermissionDataProvider is an interface for fetching permission data during warmup.
// This allows the cache to be decoupled from the repository layer.
type PermissionDataProvider interface {
	// GetPermissionTree returns the permission tree for caching
	GetPermissionTree(ctx context.Context) ([]PermissionTreeNode, error)
	// GetPermissionGroups returns all permission groups
	GetPermissionGroups(ctx context.Context) ([]string, error)
	// GetRolePermissionCodes returns permission codes for a role by slug
	GetRolePermissionCodes(ctx context.Context, roleSlug string) ([]string, error)
	// GetRoleIDBySlug returns the role ID for a given slug
	GetRoleIDBySlug(ctx context.Context, roleSlug string) (uint64, error)
}

// Warmup preloads commonly accessed permission data into the cache.
// This should be called during system startup to reduce cold-start latency.
func (pc *PermissionCache) Warmup(ctx context.Context, provider PermissionDataProvider, config WarmupConfig) WarmupResult {
	result := WarmupResult{
		RolesWarmed: make([]string, 0),
		Errors:      make([]error, 0),
	}

	// Warmup permission tree
	if config.WarmupPermissionTree {
		tree, err := provider.GetPermissionTree(ctx)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to warmup permission tree: %w", err))
		} else {
			if err := pc.SetPermissionTree(ctx, tree); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to cache permission tree: %w", err))
			} else {
				result.PermissionTreeWarmed = true
			}
		}
	}

	// Warmup permission groups
	if config.WarmupPermissionGroups {
		groups, err := provider.GetPermissionGroups(ctx)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to warmup permission groups: %w", err))
		} else {
			if err := pc.SetPermissionGroups(ctx, groups); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to cache permission groups: %w", err))
			} else {
				result.PermissionGroupsWarmed = true
			}
		}
	}

	// Warmup system role permissions
	for _, roleSlug := range config.SystemRoleSlugs {
		roleID, err := provider.GetRoleIDBySlug(ctx, roleSlug)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to get role ID for %s: %w", roleSlug, err))
			continue
		}

		codes, err := provider.GetRolePermissionCodes(ctx, roleSlug)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to warmup role %s permissions: %w", roleSlug, err))
			continue
		}

		if err := pc.SetRolePermissions(ctx, roleID, codes); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to cache role %s permissions: %w", roleSlug, err))
			continue
		}

		result.RolesWarmed = append(result.RolesWarmed, roleSlug)
	}

	return result
}

// IsWarmedUp checks if the cache has been warmed up by checking for the permission tree.
func (pc *PermissionCache) IsWarmedUp(ctx context.Context) bool {
	_, ok, _ := pc.GetPermissionTree(ctx)
	return ok
}

// UserRoleProvider is an interface for fetching user-role relationships.
// This is used for cache invalidation propagation.
type UserRoleProvider interface {
	// GetUserIDsByRoleID returns all user IDs that have the specified role
	GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error)
}

// InvalidateRolePermissionsAndPropagateToUsers invalidates the role's permission cache
// and propagates the invalidation to all users who have that role.
// This should be called when role permissions are modified.
func (pc *PermissionCache) InvalidateRolePermissionsAndPropagateToUsers(
	ctx context.Context,
	roleID uint64,
	provider UserRoleProvider,
) error {
	// First, invalidate the role's permission cache
	if err := pc.InvalidateRoleCache(ctx, roleID); err != nil {
		return fmt.Errorf("failed to invalidate role cache: %w", err)
	}

	// Get all users with this role
	userIDs, err := provider.GetUserIDsByRoleID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get users for role %d: %w", roleID, err)
	}

	// Invalidate cache for each user
	var lastErr error
	for _, userID := range userIDs {
		if err := pc.InvalidateUserCache(ctx, userID); err != nil {
			lastErr = err
			// Continue invalidating other users even if one fails
		}
	}

	return lastErr
}

// InvalidateUserRolesAndPermissions invalidates both the user's role cache
// and permission cache. This should be called when user roles are modified.
func (pc *PermissionCache) InvalidateUserRolesAndPermissions(ctx context.Context, userID uint64) error {
	return pc.InvalidateUserCache(ctx, userID)
}

// InvalidateMultipleUsers invalidates caches for multiple users.
// This is useful for batch operations.
func (pc *PermissionCache) InvalidateMultipleUsers(ctx context.Context, userIDs []uint64) error {
	var lastErr error
	for _, userID := range userIDs {
		if err := pc.InvalidateUserCache(ctx, userID); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// InvalidateMultipleRoles invalidates caches for multiple roles.
// This is useful for batch operations.
func (pc *PermissionCache) InvalidateMultipleRoles(ctx context.Context, roleIDs []uint64) error {
	var lastErr error
	for _, roleID := range roleIDs {
		if err := pc.InvalidateRoleCache(ctx, roleID); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// InvalidateAllForRole invalidates the role cache and all related user caches.
// This is a convenience method that combines role and user invalidation.
func (pc *PermissionCache) InvalidateAllForRole(
	ctx context.Context,
	roleID uint64,
	provider UserRoleProvider,
) error {
	return pc.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID, provider)
}

// CacheInvalidationEvent represents a cache invalidation event for logging/monitoring.
type CacheInvalidationEvent struct {
	Type      string // "user", "role", "tree", "groups", "all"
	TargetID  uint64 // User ID or Role ID (0 for tree/groups/all)
	Timestamp time.Time
	Reason    string
}

// CacheInvalidationListener is a callback for cache invalidation events.
type CacheInvalidationListener func(event CacheInvalidationEvent)

// PermissionCacheWithListener wraps PermissionCache with event listeners.
type PermissionCacheWithListener struct {
	*PermissionCache
	listeners []CacheInvalidationListener
	mu        sync.RWMutex
}

// NewPermissionCacheWithListener creates a new permission cache with listener support.
func NewPermissionCacheWithListener(cache Cache) *PermissionCacheWithListener {
	return &PermissionCacheWithListener{
		PermissionCache: NewPermissionCache(cache),
		listeners:       make([]CacheInvalidationListener, 0),
	}
}

// AddListener adds a cache invalidation listener.
func (pcl *PermissionCacheWithListener) AddListener(listener CacheInvalidationListener) {
	pcl.mu.Lock()
	defer pcl.mu.Unlock()
	pcl.listeners = append(pcl.listeners, listener)
}

// notifyListeners notifies all listeners of a cache invalidation event.
func (pcl *PermissionCacheWithListener) notifyListeners(event CacheInvalidationEvent) {
	pcl.mu.RLock()
	listeners := make([]CacheInvalidationListener, len(pcl.listeners))
	copy(listeners, pcl.listeners)
	pcl.mu.RUnlock()

	for _, listener := range listeners {
		listener(event)
	}
}

// InvalidateUserCacheWithEvent invalidates user cache and notifies listeners.
func (pcl *PermissionCacheWithListener) InvalidateUserCacheWithEvent(ctx context.Context, userID uint64, reason string) error {
	err := pcl.InvalidateUserCache(ctx, userID)
	pcl.notifyListeners(CacheInvalidationEvent{
		Type:      "user",
		TargetID:  userID,
		Timestamp: time.Now(),
		Reason:    reason,
	})
	return err
}

// InvalidateRoleCacheWithEvent invalidates role cache and notifies listeners.
func (pcl *PermissionCacheWithListener) InvalidateRoleCacheWithEvent(ctx context.Context, roleID uint64, reason string) error {
	err := pcl.InvalidateRoleCache(ctx, roleID)
	pcl.notifyListeners(CacheInvalidationEvent{
		Type:      "role",
		TargetID:  roleID,
		Timestamp: time.Now(),
		Reason:    reason,
	})
	return err
}

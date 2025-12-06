package cache

import (
	"context"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: rbac-button-level-permission, Property 9: 权限缓存失效传播**
// **Validates: Requirements 5.1, 5.2**
//
// Property 9: 权限缓存失效传播
// *For any* role whose permissions are modified, all users with that role should
// receive updated permissions on their next permission check (cache should be invalidated).

// mockUserRoleProviderForProperty implements UserRoleProvider for property testing
type mockUserRoleProviderForProperty struct {
	roleUsers map[uint64][]uint64
}

func newMockUserRoleProviderForProperty() *mockUserRoleProviderForProperty {
	return &mockUserRoleProviderForProperty{
		roleUsers: make(map[uint64][]uint64),
	}
}

func (m *mockUserRoleProviderForProperty) GetUserIDsByRoleID(ctx context.Context, roleID uint64) ([]uint64, error) {
	return m.roleUsers[roleID], nil
}

func (m *mockUserRoleProviderForProperty) setUsersForRole(roleID uint64, userIDs []uint64) {
	m.roleUsers[roleID] = userIDs
}

// genPermissionCode generates valid permission codes in the format module.resource.action
func genPermissionCode() gopter.Gen {
	return genValidSegment().FlatMap(func(module interface{}) gopter.Gen {
		return genValidSegment().FlatMap(func(resource interface{}) gopter.Gen {
			return genValidSegment().Map(func(action string) string {
				return module.(string) + "." + resource.(string) + "." + action
			})
		}, reflect.TypeOf(""))
	}, reflect.TypeOf(""))
}

// genValidSegment generates a valid segment for permission codes.
// A valid segment starts with a lowercase letter and contains only lowercase letters and digits.
func genValidSegment() gopter.Gen {
	// Generate lowercase letters for the first character (a-z)
	firstChar := gen.IntRange(0, 25).Map(func(i int) rune {
		return rune('a' + i)
	})

	// Generate length for rest of string (0-9 additional characters)
	restLength := gen.IntRange(0, 9)

	// Generate each character as lowercase letter or digit
	charGen := gen.IntRange(0, 35).Map(func(i int) rune {
		if i < 26 {
			return rune('a' + i)
		}
		return rune('0' + (i - 26))
	})

	return gopter.CombineGens(firstChar, restLength).FlatMap(func(vals interface{}) gopter.Gen {
		v := vals.([]interface{})
		first := v[0].(rune)
		length := v[1].(int)

		if length == 0 {
			return gen.Const(string(first))
		}

		return gen.SliceOfN(length, charGen).Map(func(rest []rune) string {
			return string(first) + string(rest)
		})
	}, reflect.TypeOf(""))
}

// TestProperty9_CacheInvalidationPropagation tests that when role permissions are modified,
// all users with that role have their cache invalidated.
func TestProperty9_CacheInvalidationPropagation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: For any role with any number of users, when role permissions are modified,
	// all users with that role should have their cache invalidated
	properties.Property("role permission modification invalidates all user caches", prop.ForAll(
		func(roleID uint64, userIDs []uint64, permCodes []string) bool {
			if len(userIDs) == 0 || len(permCodes) == 0 {
				return true // Skip empty cases
			}

			ctx := context.Background()
			memCache := NewMemory()
			pc := NewPermissionCache(memCache)
			provider := newMockUserRoleProviderForProperty()

			// Set up the provider with the role-user mapping
			provider.setUsersForRole(roleID, userIDs)

			// Set up caches for the role and all users
			if err := pc.SetRolePermissions(ctx, roleID, permCodes); err != nil {
				return false
			}

			for _, userID := range userIDs {
				if err := pc.SetUserPermissions(ctx, userID, permCodes); err != nil {
					return false
				}
			}

			// Verify all caches exist before invalidation
			_, roleOk, _ := pc.GetRolePermissions(ctx, roleID)
			if !roleOk {
				return false
			}

			for _, userID := range userIDs {
				_, userOk, _ := pc.GetUserPermissions(ctx, userID)
				if !userOk {
					return false
				}
			}

			// Simulate role permission modification by invalidating and propagating
			if err := pc.InvalidateRolePermissionsAndPropagateToUsers(ctx, roleID, provider); err != nil {
				return false
			}

			// PROPERTY: After role permission modification, role cache should be invalidated
			_, roleOk, _ = pc.GetRolePermissions(ctx, roleID)
			if roleOk {
				return false // Role cache should be invalidated
			}

			// PROPERTY: All users with that role should have their cache invalidated
			for _, userID := range userIDs {
				_, userOk, _ := pc.GetUserPermissions(ctx, userID)
				if userOk {
					return false // User cache should be invalidated
				}
			}

			return true
		},
		gen.UInt64Range(1, 1000),                    // roleID
		gen.SliceOfN(10, gen.UInt64Range(1, 10000)), // userIDs (up to 10 users)
		gen.SliceOfN(5, genPermissionCode()),        // permCodes
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty9_UserRoleChangeInvalidation tests that when a user's roles are modified,
// their permission cache is invalidated.
func TestProperty9_UserRoleChangeInvalidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property: For any user, when their roles are modified, their permission and role cache
	// should be invalidated
	properties.Property("user role modification invalidates user caches", prop.ForAll(
		func(userID uint64, permCodes []string, roleIDs []uint64) bool {
			if len(permCodes) == 0 || len(roleIDs) == 0 {
				return true // Skip empty cases
			}

			ctx := context.Background()
			memCache := NewMemory()
			pc := NewPermissionCache(memCache)

			// Set up caches for the user
			if err := pc.SetUserPermissions(ctx, userID, permCodes); err != nil {
				return false
			}
			if err := pc.SetUserRoles(ctx, userID, roleIDs); err != nil {
				return false
			}

			// Verify caches exist before invalidation
			_, permOk, _ := pc.GetUserPermissions(ctx, userID)
			if !permOk {
				return false
			}
			_, roleOk, _ := pc.GetUserRoles(ctx, userID)
			if !roleOk {
				return false
			}

			// Simulate user role change by invalidating user cache
			if err := pc.InvalidateUserRolesAndPermissions(ctx, userID); err != nil {
				return false
			}

			// PROPERTY: After user role modification, user's permission cache should be invalidated
			_, permOk, _ = pc.GetUserPermissions(ctx, userID)
			if permOk {
				return false // Permission cache should be invalidated
			}

			// PROPERTY: After user role modification, user's role cache should be invalidated
			_, roleOk, _ = pc.GetUserRoles(ctx, userID)
			if roleOk {
				return false // Role cache should be invalidated
			}

			return true
		},
		gen.UInt64Range(1, 10000),                // userID
		gen.SliceOfN(5, genPermissionCode()),     // permCodes
		gen.SliceOfN(5, gen.UInt64Range(1, 100)), // roleIDs
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty9_MultipleRoleInvalidation tests that invalidating multiple roles
// correctly propagates to all affected users.
func TestProperty9_MultipleRoleInvalidation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 50
	properties := gopter.NewProperties(parameters)

	properties.Property("multiple role invalidation propagates correctly", prop.ForAll(
		func(roleIDs []uint64) bool {
			if len(roleIDs) == 0 {
				return true
			}

			ctx := context.Background()
			memCache := NewMemory()
			pc := NewPermissionCache(memCache)

			// Set up caches for all roles
			for _, roleID := range roleIDs {
				if err := pc.SetRolePermissions(ctx, roleID, []string{"perm1", "perm2"}); err != nil {
					return false
				}
			}

			// Verify all caches exist
			for _, roleID := range roleIDs {
				_, ok, _ := pc.GetRolePermissions(ctx, roleID)
				if !ok {
					return false
				}
			}

			// Invalidate all roles
			if err := pc.InvalidateMultipleRoles(ctx, roleIDs); err != nil {
				return false
			}

			// PROPERTY: All role caches should be invalidated
			for _, roleID := range roleIDs {
				_, ok, _ := pc.GetRolePermissions(ctx, roleID)
				if ok {
					return false // Cache should be invalidated
				}
			}

			return true
		},
		gen.SliceOfN(10, gen.UInt64Range(1, 1000)),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

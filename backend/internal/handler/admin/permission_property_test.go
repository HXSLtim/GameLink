package admin

import (
	"context"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"gamelink/internal/model"
)

// mockPermissionServiceForProperty is a mock implementation of the permission service for property testing.
// **Feature: rbac-button-level-permission, Property 10: 登录权限完整性**
// **Validates: Requirements 5.3**
type mockPermissionServiceForProperty struct {
	// rolePermissions maps roleID to permissions
	rolePermissions map[uint64][]model.Permission
	// userRoles maps userID to roleIDs
	userRoles map[uint64][]uint64
}

func newMockPermissionServiceForProperty() *mockPermissionServiceForProperty {
	return &mockPermissionServiceForProperty{
		rolePermissions: make(map[uint64][]model.Permission),
		userRoles:       make(map[uint64][]uint64),
	}
}

func (m *mockPermissionServiceForProperty) SetRolePermissions(roleID uint64, permissions []model.Permission) {
	m.rolePermissions[roleID] = permissions
}

func (m *mockPermissionServiceForProperty) SetUserRoles(userID uint64, roleIDs []uint64) {
	m.userRoles[userID] = roleIDs
}

// ListPermissionsByUserID returns all permissions for a user by aggregating permissions from all their roles.
// This simulates the actual behavior of the permission service.
func (m *mockPermissionServiceForProperty) ListPermissionsByUserID(ctx context.Context, userID uint64) ([]model.Permission, error) {
	roleIDs := m.userRoles[userID]
	if len(roleIDs) == 0 {
		return []model.Permission{}, nil
	}

	// Aggregate permissions from all roles (union)
	permissionMap := make(map[string]model.Permission)
	for _, roleID := range roleIDs {
		permissions := m.rolePermissions[roleID]
		for _, perm := range permissions {
			// Use code as key for deduplication
			if perm.Code != "" {
				permissionMap[perm.Code] = perm
			}
		}
	}

	// Convert map to slice
	result := make([]model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		result = append(result, perm)
	}
	return result, nil
}

// GetUserPermissionCodes returns all permission codes for a user.
func (m *mockPermissionServiceForProperty) GetUserPermissionCodes(ctx context.Context, userID uint64) ([]string, error) {
	permissions, err := m.ListPermissionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		if perm.Code != "" {
			codes = append(codes, perm.Code)
		}
	}
	return codes, nil
}

// mockRoleServiceForProperty is a mock implementation of the role service for property testing.
type mockRoleServiceForProperty struct {
	superAdminUsers map[uint64]bool
}

func newMockRoleServiceForProperty() *mockRoleServiceForProperty {
	return &mockRoleServiceForProperty{
		superAdminUsers: make(map[uint64]bool),
	}
}

func (m *mockRoleServiceForProperty) SetSuperAdmin(userID uint64, isSuperAdmin bool) {
	m.superAdminUsers[userID] = isSuperAdmin
}

func (m *mockRoleServiceForProperty) CheckUserIsSuperAdmin(ctx context.Context, userID uint64) (bool, error) {
	return m.superAdminUsers[userID], nil
}

// simulateGetCurrentUserPermissions simulates the GetCurrentUserPermissions handler logic.
// Returns the permission codes that would be returned by the API.
func simulateGetCurrentUserPermissions(
	ctx context.Context,
	userID uint64,
	permSvc *mockPermissionServiceForProperty,
	roleSvc *mockRoleServiceForProperty,
) []string {
	// Check if super admin
	if roleSvc != nil {
		isSuperAdmin, err := roleSvc.CheckUserIsSuperAdmin(ctx, userID)
		if err == nil && isSuperAdmin {
			return []string{"*"}
		}
	}

	// Get user permissions
	codes, err := permSvc.GetUserPermissionCodes(ctx, userID)
	if err != nil {
		return []string{}
	}

	if codes == nil {
		return []string{}
	}
	return codes
}

// TestProperty10_LoginPermissionCompleteness tests Property 10: Login Permission Completeness
// For any user who successfully logs in, the response should contain the complete list of
// permission codes derived from all their assigned roles.
// **Feature: rbac-button-level-permission, Property 10: 登录权限完整性**
// **Validates: Requirements 5.3**
func TestProperty10_LoginPermissionCompleteness(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// Property 1: Super admin should receive ['*'] permission
	properties.Property("super admin should receive wildcard permission", prop.ForAll(
		func(userID uint64) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			ctx := context.Background()
			permSvc := newMockPermissionServiceForProperty()
			roleSvc := newMockRoleServiceForProperty()
			roleSvc.SetSuperAdmin(userID, true)

			codes := simulateGetCurrentUserPermissions(ctx, userID, permSvc, roleSvc)

			// Super admin should receive exactly ['*']
			return len(codes) == 1 && codes[0] == "*"
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 2: Non-super admin should receive all permission codes from their roles
	properties.Property("non-super admin should receive all permission codes from their roles", prop.ForAll(
		func(userID uint64, roleIDs []uint64, permCodes []string) bool {
			if userID == 0 || len(roleIDs) == 0 || len(permCodes) == 0 {
				return true // Skip invalid inputs
			}

			ctx := context.Background()
			permSvc := newMockPermissionServiceForProperty()
			roleSvc := newMockRoleServiceForProperty()
			roleSvc.SetSuperAdmin(userID, false)

			// Assign roles to user
			permSvc.SetUserRoles(userID, roleIDs)

			// Create permissions for each role
			// Distribute permissions across roles
			expectedCodes := make(map[string]bool)
			for i, roleID := range roleIDs {
				permissions := make([]model.Permission, 0)
				// Each role gets some permissions
				for j, code := range permCodes {
					if j%len(roleIDs) == i%len(roleIDs) {
						permissions = append(permissions, model.Permission{
							Base: model.Base{ID: uint64(i*100 + j + 1)},
							Code: code,
						})
						expectedCodes[code] = true
					}
				}
				permSvc.SetRolePermissions(roleID, permissions)
			}

			codes := simulateGetCurrentUserPermissions(ctx, userID, permSvc, roleSvc)

			// All expected codes should be present
			actualCodes := make(map[string]bool)
			for _, code := range codes {
				actualCodes[code] = true
			}

			for code := range expectedCodes {
				if !actualCodes[code] {
					return false
				}
			}

			return true
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 })),
		gen.SliceOfN(5, genValidPermissionCodeForProperty()),
	))

	// Property 3: User with no roles should receive empty permission list
	properties.Property("user with no roles should receive empty permission list", prop.ForAll(
		func(userID uint64) bool {
			if userID == 0 {
				return true // Skip invalid user ID
			}

			ctx := context.Background()
			permSvc := newMockPermissionServiceForProperty()
			roleSvc := newMockRoleServiceForProperty()
			roleSvc.SetSuperAdmin(userID, false)
			// Don't assign any roles

			codes := simulateGetCurrentUserPermissions(ctx, userID, permSvc, roleSvc)

			return len(codes) == 0
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
	))

	// Property 4: Permissions from multiple roles should be merged (union)
	properties.Property("permissions from multiple roles should be merged", prop.ForAll(
		func(userID uint64, role1ID, role2ID uint64, codes1, codes2 []string) bool {
			if userID == 0 || role1ID == 0 || role2ID == 0 {
				return true // Skip invalid inputs
			}
			if len(codes1) == 0 || len(codes2) == 0 {
				return true // Skip empty inputs
			}

			ctx := context.Background()
			permSvc := newMockPermissionServiceForProperty()
			roleSvc := newMockRoleServiceForProperty()
			roleSvc.SetSuperAdmin(userID, false)

			// Assign both roles to user
			permSvc.SetUserRoles(userID, []uint64{role1ID, role2ID})

			// Create permissions for role 1
			perms1 := make([]model.Permission, len(codes1))
			for i, code := range codes1 {
				perms1[i] = model.Permission{
					Base: model.Base{ID: uint64(i + 1)},
					Code: code,
				}
			}
			permSvc.SetRolePermissions(role1ID, perms1)

			// Create permissions for role 2
			perms2 := make([]model.Permission, len(codes2))
			for i, code := range codes2 {
				perms2[i] = model.Permission{
					Base: model.Base{ID: uint64(100 + i + 1)},
					Code: code,
				}
			}
			permSvc.SetRolePermissions(role2ID, perms2)

			// Get user permissions
			resultCodes := simulateGetCurrentUserPermissions(ctx, userID, permSvc, roleSvc)

			// Build expected set (union of both roles)
			expectedSet := make(map[string]bool)
			for _, code := range codes1 {
				expectedSet[code] = true
			}
			for _, code := range codes2 {
				expectedSet[code] = true
			}

			// Build actual set
			actualSet := make(map[string]bool)
			for _, code := range resultCodes {
				actualSet[code] = true
			}

			// All expected codes should be present
			for code := range expectedSet {
				if !actualSet[code] {
					return false
				}
			}

			// No extra codes should be present
			for code := range actualSet {
				if !expectedSet[code] {
					return false
				}
			}

			return true
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.SliceOfN(3, genValidPermissionCodeForProperty()),
		gen.SliceOfN(3, genValidPermissionCodeForProperty()),
	))

	// Property 5: Duplicate permissions across roles should be deduplicated
	properties.Property("duplicate permissions across roles should be deduplicated", prop.ForAll(
		func(userID uint64, role1ID, role2ID uint64, sharedCode string) bool {
			if userID == 0 || role1ID == 0 || role2ID == 0 || sharedCode == "" {
				return true // Skip invalid inputs
			}

			ctx := context.Background()
			permSvc := newMockPermissionServiceForProperty()
			roleSvc := newMockRoleServiceForProperty()
			roleSvc.SetSuperAdmin(userID, false)

			// Assign both roles to user
			permSvc.SetUserRoles(userID, []uint64{role1ID, role2ID})

			// Both roles have the same permission
			sharedPerm := model.Permission{
				Base: model.Base{ID: 1},
				Code: sharedCode,
			}
			permSvc.SetRolePermissions(role1ID, []model.Permission{sharedPerm})
			permSvc.SetRolePermissions(role2ID, []model.Permission{sharedPerm})

			// Get user permissions
			resultCodes := simulateGetCurrentUserPermissions(ctx, userID, permSvc, roleSvc)

			// Should have exactly one instance of the shared code
			count := 0
			for _, code := range resultCodes {
				if code == sharedCode {
					count++
				}
			}

			return count == 1
		},
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		gen.UInt64().SuchThat(func(id uint64) bool { return id > 0 }),
		genValidPermissionCodeForProperty(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// genValidPermissionCodeForProperty generates a valid permission code (module.resource.action).
func genValidPermissionCodeForProperty() gopter.Gen {
	// Generate a fixed-length segment to avoid filtering issues
	segmentGen := gen.IntRange(1, 10).FlatMap(func(length interface{}) gopter.Gen {
		return gen.SliceOfN(length.(int), gen.Rune()).Map(func(runes []rune) string {
			result := make([]byte, len(runes))
			for i, r := range runes {
				// Convert to lowercase letter a-z
				result[i] = byte('a' + (r % 26))
			}
			return string(result)
		})
	}, reflect.TypeOf(""))

	return gopter.CombineGens(segmentGen, segmentGen, segmentGen).Map(func(vals []interface{}) string {
		module := vals[0].(string)
		resource := vals[1].(string)
		action := vals[2].(string)
		return module + "." + resource + "." + action
	})
}

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionCodeConstants(t *testing.T) {
	// Test that all NEW permission code constants follow the correct three-part format
	// Note: Legacy codes like "review.list" use two-part format for backward compatibility
	threePartCodes := []string{
		// Admin Permission Module
		PermCodeAdminPermissionsRead,
		PermCodeAdminPermissionsCreate,
		PermCodeAdminPermissionsUpdate,
		PermCodeAdminPermissionsDelete,
		// Admin Role Module
		PermCodeAdminRolesRead,
		PermCodeAdminRolesCreate,
		PermCodeAdminRolesUpdate,
		PermCodeAdminRolesDelete,
		PermCodeAdminRolesAssign,
		// Admin User Module
		PermCodeAdminUsersRead,
		PermCodeAdminUsersCreate,
		PermCodeAdminUsersUpdate,
		PermCodeAdminUsersDelete,
		PermCodeAdminUsersAssign,
		PermCodeAdminUsersStatus,
		PermCodeAdminUsersPoints,
		// Admin Player Module
		PermCodeAdminPlayersRead,
		PermCodeAdminPlayersCreate,
		PermCodeAdminPlayersUpdate,
		PermCodeAdminPlayersDelete,
		PermCodeAdminPlayersVerify,
		// Admin Game Module
		PermCodeAdminGamesRead,
		PermCodeAdminGamesCreate,
		PermCodeAdminGamesUpdate,
		PermCodeAdminGamesDelete,
		// Admin Order Module
		PermCodeAdminOrdersRead,
		PermCodeAdminOrdersCreate,
		PermCodeAdminOrdersUpdate,
		PermCodeAdminOrdersDelete,
		PermCodeAdminOrdersAssign,
		PermCodeAdminOrdersCancel,
		PermCodeAdminOrdersRefund,
		PermCodeAdminOrdersConfirm,
		PermCodeAdminOrdersStart,
		PermCodeAdminOrdersComplete,
		// Admin Payment Module
		PermCodeAdminPaymentsRead,
		PermCodeAdminPaymentsRefund,
		PermCodeAdminPaymentsExport,
		// Admin Withdraw Module
		PermCodeAdminWithdrawsRead,
		PermCodeAdminWithdrawsApprove,
		PermCodeAdminWithdrawsReject,
		PermCodeAdminWithdrawsExport,
		// Admin Audit Module
		PermCodeAdminAuditRead,
		PermCodeAdminAuditExport,
		// Admin Stats Module
		PermCodeAdminStatsRead,
		PermCodeAdminStatsExport,
		// Admin Menu Module
		PermCodeAdminMenusRead,
		PermCodeAdminMenusCreate,
		PermCodeAdminMenusUpdate,
		PermCodeAdminMenusDelete,
		// Content Module (three-part format)
		PermCodeContentFeedList,
		PermCodeContentFeedApprove,
		PermCodeContentChatList,
		PermCodeContentReportList,
	}

	for _, code := range threePartCodes {
		t.Run(code, func(t *testing.T) {
			assert.True(t, IsValidPermissionCode(code), "Permission code %s should be valid", code)
		})
	}

	// Test legacy two-part codes exist (backward compatibility)
	legacyCodes := []string{
		PermCodeReviewList,
		PermCodeReviewGet,
		PermCodeReviewApprove,
		PermCodeReviewReject,
		PermCodeReviewDelete,
		PermCodeContentStats,
	}

	for _, code := range legacyCodes {
		t.Run("legacy_"+code, func(t *testing.T) {
			assert.NotEmpty(t, code, "Legacy permission code should not be empty")
		})
	}
}

func TestIsValidPermissionCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{"valid three-part code", "admin.users.read", true},
		{"valid with numbers", "admin.users2.read", true},
		{"empty code", "", false},
		{"single part", "admin", false},
		{"two parts", "admin.users", false},
		{"four parts", "admin.users.read.extra", false},
		{"uppercase", "Admin.Users.Read", false},
		{"starts with number", "1admin.users.read", false},
		{"contains special chars", "admin.users!.read", false},
		{"contains spaces", "admin users.read", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPermissionCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitPermissionCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected []string
	}{
		{"three parts", "admin.users.read", []string{"admin", "users", "read"}},
		{"single part", "admin", []string{"admin"}},
		{"empty", "", nil},
		{"two parts", "admin.users", []string{"admin", "users"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitPermissionCode(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPermissionCodesByModule(t *testing.T) {
	// Test getting admin module codes
	adminCodes := GetPermissionCodesByModule("admin")
	assert.NotEmpty(t, adminCodes, "Should have admin module codes")

	// Verify all returned codes start with "admin."
	for _, code := range adminCodes {
		parts := splitPermissionCode(code)
		assert.Equal(t, "admin", parts[0], "Code %s should start with admin", code)
	}
}

func TestAllPermissionGroups(t *testing.T) {
	groups := AllPermissionGroups()
	assert.NotEmpty(t, groups, "Should have permission groups")

	// Verify each group has required fields
	for _, group := range groups {
		assert.NotEmpty(t, group.Group, "Group should have Group field")
		assert.NotEmpty(t, group.Name, "Group should have Name field")
		assert.NotEmpty(t, group.Module, "Group should have Module field")
	}
}

func TestGetAllPermissionDefinitions(t *testing.T) {
	defs := GetAllPermissionDefinitions()
	assert.NotEmpty(t, defs, "Should have permission definitions")

	// Verify each definition has required fields
	for _, def := range defs {
		assert.NotEmpty(t, def.Code, "Definition should have Code")
		assert.NotEmpty(t, def.Method, "Definition should have Method")
		assert.NotEmpty(t, def.Path, "Definition should have Path")
		assert.NotEmpty(t, def.Group, "Definition should have Group")
		assert.NotEmpty(t, def.Description, "Definition should have Description")
	}
}

func TestSuperAdminPermission(t *testing.T) {
	assert.Equal(t, "*", SuperAdminPermissionCode)
	assert.True(t, IsSuperAdminPermission("*"))
	assert.False(t, IsSuperAdminPermission("admin.users.read"))
	assert.False(t, IsSuperAdminPermission(""))
}

func TestPermissionGroupInfo(t *testing.T) {
	// Test specific groups
	assert.Equal(t, "/admin/permissions", PermGroupAdminPermissions.Group)
	assert.Equal(t, "权限管理", PermGroupAdminPermissions.Name)
	assert.Equal(t, "admin", PermGroupAdminPermissions.Module)

	assert.Equal(t, "/admin/roles", PermGroupAdminRoles.Group)
	assert.Equal(t, "角色管理", PermGroupAdminRoles.Name)

	assert.Equal(t, "/admin/users", PermGroupAdminUsers.Group)
	assert.Equal(t, "用户管理", PermGroupAdminUsers.Name)
}

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/pkg/testutil"
)

func TestSeedDefaultRoles(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
	)
	require.NoError(t, err)

	// First, seed some permissions that the roles will reference
	testPermissions := []model.Permission{
		{Code: model.PermCodeAdminUsersRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/users", Group: "/admin/users", Description: "查看用户列表"},
		{Code: model.PermCodeAdminUsersCreate, Method: model.HTTPMethodPOST, Path: "/api/v1/admin/users", Group: "/admin/users", Description: "创建用户"},
		{Code: model.PermCodeReviewList, Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews", Group: "/admin/reviews", Description: "查看评价列表"},
		{Code: model.PermCodeReviewGet, Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews/:id", Group: "/admin/reviews", Description: "查看评价详情"},
		{Code: model.PermCodeServiceItemRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/service-items", Group: "/admin/service-items", Description: "查看服务项目"},
		{Code: model.PermCodeWalletRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/wallets", Group: "/admin/wallets", Description: "查看钱包"},
		{Code: model.PermCodeNotificationRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/notifications", Group: "/admin/notifications", Description: "查看通知"},
		{Code: model.PermCodeDisputeRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/disputes", Group: "/admin/disputes", Description: "查看纠纷"},
		{Code: model.PermCodeDisputeCreate, Method: model.HTTPMethodPOST, Path: "/api/v1/admin/disputes", Group: "/admin/disputes", Description: "创建纠纷"},
	}
	for _, perm := range testPermissions {
		require.NoError(t, db.Create(&perm).Error)
	}

	// Execute seed function
	err = seedDefaultRoles(db)
	require.NoError(t, err)

	// Verify roles were created
	var roles []model.RoleModel
	err = db.Find(&roles).Error
	require.NoError(t, err)
	assert.Len(t, roles, 4, "Should have 4 system roles")

	// Verify each role exists with correct properties
	roleMap := make(map[string]*model.RoleModel)
	for i := range roles {
		roleMap[roles[i].Slug] = &roles[i]
	}

	// Check superAdmin role
	superAdmin, ok := roleMap[string(model.RoleSlugSuperAdmin)]
	require.True(t, ok, "superAdmin role should exist")
	assert.Equal(t, "超级管理员", superAdmin.Name)
	assert.True(t, superAdmin.IsSystem)
	assert.Equal(t, 1000, superAdmin.Priority)

	// Check admin role
	admin, ok := roleMap[string(model.RoleSlugAdmin)]
	require.True(t, ok, "admin role should exist")
	assert.Equal(t, "管理员", admin.Name)
	assert.True(t, admin.IsSystem)
	assert.Equal(t, 500, admin.Priority)

	// Check player role
	player, ok := roleMap[string(model.RoleSlugPlayer)]
	require.True(t, ok, "player role should exist")
	assert.Equal(t, "陪玩师", player.Name)
	assert.True(t, player.IsSystem)
	assert.Equal(t, 100, player.Priority)

	// Check user role
	user, ok := roleMap[string(model.RoleSlugUser)]
	require.True(t, ok, "user role should exist")
	assert.Equal(t, "普通用户", user.Name)
	assert.True(t, user.IsSystem)
	assert.Equal(t, 10, user.Priority)

	// Verify admin role has permissions assigned
	var adminPermCount int64
	err = db.Model(&model.RolePermission{}).Where("role_id = ?", admin.ID).Count(&adminPermCount).Error
	require.NoError(t, err)
	assert.Greater(t, adminPermCount, int64(0), "Admin role should have permissions assigned")

	// Verify player role has permissions assigned
	var playerPermCount int64
	err = db.Model(&model.RolePermission{}).Where("role_id = ?", player.ID).Count(&playerPermCount).Error
	require.NoError(t, err)
	assert.Greater(t, playerPermCount, int64(0), "Player role should have permissions assigned")

	// Verify user role has permissions assigned
	var userPermCount int64
	err = db.Model(&model.RolePermission{}).Where("role_id = ?", user.ID).Count(&userPermCount).Error
	require.NoError(t, err)
	assert.Greater(t, userPermCount, int64(0), "User role should have permissions assigned")

	// SuperAdmin should NOT have explicit permissions (uses wildcard)
	var superAdminPermCount int64
	err = db.Model(&model.RolePermission{}).Where("role_id = ?", superAdmin.ID).Count(&superAdminPermCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), superAdminPermCount, "SuperAdmin should not have explicit permissions (uses wildcard)")

	t.Logf("Successfully created %d roles", len(roles))
	t.Logf("Admin role has %d permissions", adminPermCount)
	t.Logf("Player role has %d permissions", playerPermCount)
	t.Logf("User role has %d permissions", userPermCount)
}

func TestSeedDefaultRolesIdempotent(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(
		&model.Permission{},
		&model.RoleModel{},
		&model.RolePermission{},
	)
	require.NoError(t, err)

	// Seed some permissions
	testPermissions := []model.Permission{
		{Code: model.PermCodeAdminUsersRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/users", Group: "/admin/users", Description: "查看用户列表"},
		{Code: model.PermCodeReviewList, Method: model.HTTPMethodGET, Path: "/api/v1/admin/reviews", Group: "/admin/reviews", Description: "查看评价列表"},
	}
	for _, perm := range testPermissions {
		require.NoError(t, db.Create(&perm).Error)
	}

	// Execute seed function twice
	err = seedDefaultRoles(db)
	require.NoError(t, err)

	err = seedDefaultRoles(db)
	require.NoError(t, err, "seedDefaultRoles should be idempotent")

	// Verify still only 4 roles
	var roleCount int64
	err = db.Model(&model.RoleModel{}).Count(&roleCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(4), roleCount, "Should still have exactly 4 roles after running twice")
}

func TestDefaultRolePermissionsDefinition(t *testing.T) {
	// Verify that DefaultRolePermissions has all expected roles
	assert.Contains(t, DefaultRolePermissions, model.RoleSlugSuperAdmin)
	assert.Contains(t, DefaultRolePermissions, model.RoleSlugAdmin)
	assert.Contains(t, DefaultRolePermissions, model.RoleSlugPlayer)
	assert.Contains(t, DefaultRolePermissions, model.RoleSlugUser)

	// Verify superAdmin has wildcard permission
	superAdminPerms := DefaultRolePermissions[model.RoleSlugSuperAdmin]
	assert.Contains(t, superAdminPerms, model.SuperAdminPermissionCode)

	// Verify admin has more permissions than player
	adminPerms := DefaultRolePermissions[model.RoleSlugAdmin]
	playerPerms := DefaultRolePermissions[model.RoleSlugPlayer]
	assert.Greater(t, len(adminPerms), len(playerPerms), "Admin should have more permissions than player")

	// Verify player has more permissions than user
	userPerms := DefaultRolePermissions[model.RoleSlugUser]
	assert.Greater(t, len(playerPerms), len(userPerms), "Player should have more permissions than user")

	t.Logf("SuperAdmin permissions: %d (wildcard)", len(superAdminPerms))
	t.Logf("Admin permissions: %d", len(adminPerms))
	t.Logf("Player permissions: %d", len(playerPerms))
	t.Logf("User permissions: %d", len(userPerms))
}

func TestSystemRoleDefinitions(t *testing.T) {
	// Verify all system roles are defined
	assert.Len(t, SystemRoleDefinitions, 4)

	slugs := make(map[string]bool)
	for _, role := range SystemRoleDefinitions {
		slugs[role.Slug] = true
		assert.True(t, role.IsSystem, "All system roles should have IsSystem=true")
		assert.NotEmpty(t, role.Name, "All roles should have a name")
		assert.NotEmpty(t, role.Description, "All roles should have a description")
	}

	assert.True(t, slugs[string(model.RoleSlugSuperAdmin)])
	assert.True(t, slugs[string(model.RoleSlugAdmin)])
	assert.True(t, slugs[string(model.RoleSlugPlayer)])
	assert.True(t, slugs[string(model.RoleSlugUser)])
}

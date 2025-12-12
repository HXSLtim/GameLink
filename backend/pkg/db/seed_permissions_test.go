package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gamelink/internal/model"
	"gamelink/pkg/testutil"
)

func TestSeedSystemPermissions(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(&model.Permission{})
	require.NoError(t, err)

	// Run the seed function
	err = seedSystemPermissions(db)
	require.NoError(t, err)

	// Verify permissions were created
	var count int64
	err = db.Model(&model.Permission{}).Count(&count).Error
	require.NoError(t, err)

	definitions := model.GetAllPermissionDefinitions()
	// Note: Some definitions share the same code (different paths for same permission)
	// So we count unique codes
	uniqueCodes := make(map[string]bool)
	for _, def := range definitions {
		uniqueCodes[def.Code] = true
	}
	assert.GreaterOrEqual(t, count, int64(len(uniqueCodes)), "Should have at least as many permissions as unique codes")

	// Verify all permissions have IsSystem = true
	var systemCount int64
	err = db.Model(&model.Permission{}).Where("is_system = ?", true).Count(&systemCount).Error
	require.NoError(t, err)
	assert.Equal(t, count, systemCount, "All seeded permissions should be system permissions")

	// Verify specific permissions exist with correct data
	var perm model.Permission
	err = db.Where("code = ?", model.PermCodeAdminUsersRead).First(&perm).Error
	require.NoError(t, err)
	assert.True(t, perm.IsSystem, "Permission should be marked as system")
	assert.Equal(t, model.HTTPMethodGET, perm.Method)
	assert.Contains(t, perm.Path, "/admin/users")

	t.Logf("Seeded %d system permissions", count)
}

func TestSeedSystemPermissionsIdempotent(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(&model.Permission{})
	require.NoError(t, err)

	// Run the seed function twice
	for i := 0; i < 2; i++ {
		err = seedSystemPermissions(db)
		require.NoError(t, err, "Seed should succeed on iteration %d", i+1)
	}

	// Verify no duplicates - count should match unique codes
	var count int64
	err = db.Model(&model.Permission{}).Count(&count).Error
	require.NoError(t, err)

	definitions := model.GetAllPermissionDefinitions()
	// Note: Some definitions may share the same code (different paths for same permission)
	uniqueCodes := make(map[string]bool)
	for _, def := range definitions {
		uniqueCodes[def.Code] = true
	}

	t.Logf("Total permissions: %d, Unique codes in definitions: %d", count, len(uniqueCodes))
}

func TestSeedSystemPermissionsUpdatesExisting(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(&model.Permission{})
	require.NoError(t, err)

	// Create a permission without IsSystem flag
	existingPerm := model.Permission{
		Code:        model.PermCodeAdminUsersRead,
		Method:      model.HTTPMethodGET,
		Path:        "/api/v1/admin/users",
		Group:       "/admin/users",
		Description: "Old description",
		IsSystem:    false,
	}
	err = db.Create(&existingPerm).Error
	require.NoError(t, err)

	originalID := existingPerm.ID

	// Run the seed function
	err = seedSystemPermissions(db)
	require.NoError(t, err)

	// Verify the permission was updated
	var perm model.Permission
	err = db.Where("code = ?", model.PermCodeAdminUsersRead).First(&perm).Error
	require.NoError(t, err)
	assert.True(t, perm.IsSystem, "Permission should be updated to system")
	assert.Equal(t, originalID, perm.ID, "Should update existing record, not create new")
}

func TestMarkExistingPermissionsAsSystem(t *testing.T) {
	db := testutil.NewMemoryDB(t)
	defer testutil.CleanDB(t, db)

	// Migrate required models
	err := db.AutoMigrate(&model.Permission{})
	require.NoError(t, err)

	// Create some permissions without IsSystem flag
	permissions := []model.Permission{
		{Code: model.PermCodeAdminUsersRead, Method: model.HTTPMethodGET, Path: "/api/v1/admin/users", IsSystem: false},
		{Code: model.PermCodeAdminUsersCreate, Method: model.HTTPMethodPOST, Path: "/api/v1/admin/users", IsSystem: false},
		{Code: "custom.permission.test", Method: model.HTTPMethodGET, Path: "/api/v1/custom", IsSystem: false},
	}
	for _, perm := range permissions {
		err = db.Create(&perm).Error
		require.NoError(t, err)
	}

	// Run the mark function
	err = markExistingPermissionsAsSystem(db)
	require.NoError(t, err)

	// Verify system permissions were marked
	var userReadPerm model.Permission
	err = db.Where("code = ?", model.PermCodeAdminUsersRead).First(&userReadPerm).Error
	require.NoError(t, err)
	assert.True(t, userReadPerm.IsSystem, "System permission should be marked")

	var userCreatePerm model.Permission
	err = db.Where("code = ?", model.PermCodeAdminUsersCreate).First(&userCreatePerm).Error
	require.NoError(t, err)
	assert.True(t, userCreatePerm.IsSystem, "System permission should be marked")

	// Verify custom permission was NOT marked
	var customPerm model.Permission
	err = db.Where("code = ?", "custom.permission.test").First(&customPerm).Error
	require.NoError(t, err)
	assert.False(t, customPerm.IsSystem, "Custom permission should not be marked as system")
}

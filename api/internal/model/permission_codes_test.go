package model

import (
	"testing"
)

func TestGameCategoryPermissions(t *testing.T) {
	// Test permission code constants
	tests := []struct {
		name string
		code string
	}{
		{"GameCategoriesRead", PermCodeAdminGameCategoriesRead},
		{"GameCategoriesCreate", PermCodeAdminGameCategoriesCreate},
		{"GameCategoriesUpdate", PermCodeAdminGameCategoriesUpdate},
		{"GameCategoriesDelete", PermCodeAdminGameCategoriesDelete},
		{"GameCategoriesBatch", PermCodeAdminGameCategoriesBatch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code == "" {
				t.Errorf("Permission code %s is empty", tt.name)
			}
		})
	}
}

func TestGameCategoryPermissionGroup(t *testing.T) {
	group := PermGroupAdminGameCategories

	if group.Group != "/admin/game-categories" {
		t.Errorf("Expected group path /admin/game-categories, got %s", group.Group)
	}

	if group.Name != "游戏分类管理" {
		t.Errorf("Expected group name '游戏分类管理', got %s", group.Name)
	}

	if group.Module != "admin" {
		t.Errorf("Expected module 'admin', got %s", group.Module)
	}
}

func TestGameCategoryPermissionsInDefinitions(t *testing.T) {
	defs := GetAllPermissionDefinitions()
	categoryGroup := PermGroupAdminGameCategories.Group

	var found []string
	for _, def := range defs {
		if def.Group == categoryGroup {
			found = append(found, def.Code)
		}
	}

	if len(found) == 0 {
		t.Errorf("No permission definitions found for game categories")
	}

	// Check for expected permission codes
	expectedCodes := []string{
		PermCodeAdminGameCategoriesRead,
		PermCodeAdminGameCategoriesCreate,
		PermCodeAdminGameCategoriesUpdate,
		PermCodeAdminGameCategoriesDelete,
		PermCodeAdminGameCategoriesBatch,
	}

	codeMap := make(map[string]bool)
	for _, code := range found {
		codeMap[code] = true
	}

	for _, expected := range expectedCodes {
		if !codeMap[expected] {
			t.Errorf("Expected permission code %s not found in definitions", expected)
		}
	}
}

func TestGameCategoryPermissionGroupInAllGroups(t *testing.T) {
	groups := AllPermissionGroups()

	categoryGroupFound := false
	for _, group := range groups {
		if group.Group == PermGroupAdminGameCategories.Group {
			categoryGroupFound = true
			break
		}
	}

	if !categoryGroupFound {
		t.Errorf("Game category permission group not found in AllPermissionGroups()")
	}
}

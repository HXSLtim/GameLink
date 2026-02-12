// Package integration provides menu batch operation integration tests.
package integration

import (
	"context"
	"testing"

	"gamelink/internal/model"
	adminrepository "gamelink/internal/repository/admin"
	adminservice "gamelink/internal/service/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// Test Helper Functions
// ============================================================================

// GetMenuCount returns the count of menus in database.
func GetMenuCount(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var count int64
	if err := db.Model(&model.Menu{}).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count menus: %v", err)
	}
	return count
}

// MenuExists checks if a menu exists by ID.
func MenuExists(t *testing.T, db *gorm.DB, menuID uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.Menu{}).Where("id = ?", menuID).Count(&count).Error; err != nil {
		t.Fatalf("Failed to check menu existence: %v", err)
	}
	return count > 0
}

// GetMenuByID retrieves a menu by ID.
func GetMenuByID(t *testing.T, db *gorm.DB, menuID uint64) *model.Menu {
	t.Helper()
	var menu model.Menu
	if err := db.First(&menu, menuID).Error; err != nil {
		return nil
	}
	return &menu
}

// HasChildMenus checks if a menu has child menus.
func HasChildMenus(t *testing.T, db *gorm.DB, parentID uint64) bool {
	t.Helper()
	var count int64
	if err := db.Model(&model.Menu{}).Where("parent_id = ?", parentID).Count(&count).Error; err != nil {
		t.Fatalf("Failed to count child menus: %v", err)
	}
	return count > 0
}

// ============================================================================
// BatchDeleteMenus Tests
// ============================================================================

func TestMenuService_BatchDeleteMenus_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create test menus (leaf menus without children)
	menu1 := CreateTestMenu(t, db, "Menu1", "/menu1", nil)
	menu2 := CreateTestMenu(t, db, "Menu2", "/menu2", nil)
	menu3 := CreateTestMenu(t, db, "Menu3", "/menu3", nil)

	// Verify initial count
	initialCount := GetMenuCount(t, db)
	assert.Equal(t, int64(3), initialCount)

	// Execute batch delete
	input := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{menu1.ID, menu2.ID, menu3.ID},
	}
	result, err := menuService.BatchDeleteMenus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)
	assert.Contains(t, result.SuccessItems, menu1.ID)
	assert.Contains(t, result.SuccessItems, menu2.ID)
	assert.Contains(t, result.SuccessItems, menu3.ID)
	assert.Empty(t, result.FailedItems)

	// Verify database state
	finalCount := GetMenuCount(t, db)
	assert.Equal(t, int64(0), finalCount)
	assert.False(t, MenuExists(t, db, menu1.ID))
	assert.False(t, MenuExists(t, db, menu2.ID))
	assert.False(t, MenuExists(t, db, menu3.ID))
}

func TestMenuService_BatchDeleteMenus_WithChildMenus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create parent and child menus
	parentMenu := CreateTestMenu(t, db, "Parent", "/parent", nil)
	childMenu1 := CreateTestMenu(t, db, "Child1", "/parent/child1", &parentMenu.ID)
	childMenu2 := CreateTestMenu(t, db, "Child2", "/parent/child2", &parentMenu.ID)
	leafMenu := CreateTestMenu(t, db, "Leaf", "/leaf", nil)

	// Verify children exist
	assert.True(t, HasChildMenus(t, db, parentMenu.ID))

	// Execute batch delete - parent has children so should fail
	input := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{parentMenu.ID, leafMenu.ID},
	}
	result, err := menuService.BatchDeleteMenus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.SuccessItems, leafMenu.ID)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item details
	failedItem := result.FailedItems[0]
	assert.Equal(t, parentMenu.ID, failedItem.ID)
	assert.Contains(t, failedItem.Message, "cannot delete menu with child menus")

	// Verify database state - parent should still exist, leaf deleted
	assert.True(t, MenuExists(t, db, parentMenu.ID))
	assert.True(t, MenuExists(t, db, childMenu1.ID))
	assert.True(t, MenuExists(t, db, childMenu2.ID))
	assert.False(t, MenuExists(t, db, leafMenu.ID))
}

func TestMenuService_BatchDeleteMenus_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create one menu
	menu1 := CreateTestMenu(t, db, "Menu1", "/menu1", nil)
	nonExistentID := uint64(999999)

	// Execute batch delete with mixed valid and invalid IDs
	input := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{menu1.ID, nonExistentID, nonExistentID + 1},
	}
	result, err := menuService.BatchDeleteMenus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 2, result.FailedCount)
	assert.Contains(t, result.SuccessItems, menu1.ID)
	assert.Len(t, result.FailedItems, 2)

	// Verify failed items are for non-existent menus
	for _, item := range result.FailedItems {
		assert.Equal(t, "record not found", item.Message)
	}

	// Verify only valid menu was deleted
	assert.False(t, MenuExists(t, db, menu1.ID))
}

func TestMenuService_BatchDeleteMenus_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Execute with empty IDs
	input := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{},
	}
	result, err := menuService.BatchDeleteMenus(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMenuService_BatchDeleteMenus_TooManyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create 101 IDs (exceeds limit of 100)
	ids := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		ids[i] = uint64(i + 1)
	}

	input := adminservice.BatchDeleteMenusInput{
		MenuIDs: ids,
	}
	result, err := menuService.BatchDeleteMenus(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 100")
}

// ============================================================================
// BatchUpdateMenuStatus Tests
// ============================================================================

func TestMenuService_BatchUpdateMenuStatus_EnableSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create menus with various hidden states
	menu1 := CreateTestMenuWithOrder(t, db, "Menu1", "/menu1", nil, 0, true)  // hidden
	menu2 := CreateTestMenuWithOrder(t, db, "Menu2", "/menu2", nil, 1, true)  // hidden
	menu3 := CreateTestMenuWithOrder(t, db, "Menu3", "/menu3", nil, 2, false) // already visible

	// Execute batch enable
	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{menu1.ID, menu2.ID, menu3.ID},
		Status:  "enabled",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)

	// Verify database state - all should be visible (hidden = false)
	updatedMenu1 := GetMenuByID(t, db, menu1.ID)
	updatedMenu2 := GetMenuByID(t, db, menu2.ID)
	updatedMenu3 := GetMenuByID(t, db, menu3.ID)

	assert.NotNil(t, updatedMenu1)
	assert.NotNil(t, updatedMenu2)
	assert.NotNil(t, updatedMenu3)

	assert.False(t, updatedMenu1.Hidden, "Menu1 should be visible")
	assert.False(t, updatedMenu2.Hidden, "Menu2 should be visible")
	assert.False(t, updatedMenu3.Hidden, "Menu3 should be visible")
}

func TestMenuService_BatchUpdateMenuStatus_DisableSuccess(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create menus
	menu1 := CreateTestMenuWithOrder(t, db, "Menu1", "/menu1", nil, 0, false) // visible
	menu2 := CreateTestMenuWithOrder(t, db, "Menu2", "/menu2", nil, 1, false) // visible
	menu3 := CreateTestMenuWithOrder(t, db, "Menu3", "/menu3", nil, 2, true)  // already hidden

	// Execute batch disable
	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{menu1.ID, menu2.ID, menu3.ID},
		Status:  "disabled",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify database state - all should be hidden (hidden = true)
	updatedMenu1 := GetMenuByID(t, db, menu1.ID)
	updatedMenu2 := GetMenuByID(t, db, menu2.ID)
	updatedMenu3 := GetMenuByID(t, db, menu3.ID)

	assert.True(t, updatedMenu1.Hidden, "Menu1 should be hidden")
	assert.True(t, updatedMenu2.Hidden, "Menu2 should be hidden")
	assert.True(t, updatedMenu3.Hidden, "Menu3 should be hidden")
}

func TestMenuService_BatchUpdateMenuStatus_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create one menu
	menu1 := CreateTestMenu(t, db, "Menu1", "/menu1", nil)
	nonExistentID := uint64(999999)

	// Execute batch update with mixed valid and invalid IDs
	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{menu1.ID, nonExistentID},
		Status:  "enabled",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalCount)
	assert.Equal(t, 1, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Contains(t, result.SuccessItems, menu1.ID)

	// Verify failed item
	assert.Len(t, result.FailedItems, 1)
	assert.Equal(t, nonExistentID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "menu not found")

	// Verify valid menu was updated
	updatedMenu := GetMenuByID(t, db, menu1.ID)
	assert.False(t, updatedMenu.Hidden)
}

func TestMenuService_BatchUpdateMenuStatus_InvalidStatus(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create a menu
	menu1 := CreateTestMenu(t, db, "Menu1", "/menu1", nil)

	// Execute with invalid status
	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{menu1.ID},
		Status:  "invalid_status",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "status must be 'enabled' or 'disabled'")
}

func TestMenuService_BatchUpdateMenuStatus_EmptyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Execute with empty IDs
	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{},
		Status:  "enabled",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMenuService_BatchUpdateMenuStatus_TooManyIDs(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create 101 IDs
	ids := make([]uint64, 101)
	for i := 0; i < 101; i++ {
		ids[i] = uint64(i + 1)
	}

	input := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: ids,
		Status:  "enabled",
	}
	result, err := menuService.BatchUpdateMenuStatus(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 100")
}

// ============================================================================
// BatchUpdateMenuOrder Tests
// ============================================================================

func TestMenuService_BatchUpdateMenuOrder_Success(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create menus with initial orders
	menu1 := CreateTestMenuWithOrder(t, db, "Menu1", "/menu1", nil, 0, false)
	menu2 := CreateTestMenuWithOrder(t, db, "Menu2", "/menu2", nil, 1, false)
	menu3 := CreateTestMenuWithOrder(t, db, "Menu3", "/menu3", nil, 2, false)

	// Execute batch order update
	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{
			{MenuID: menu1.ID, Order: 10},
			{MenuID: menu2.ID, Order: 20},
			{MenuID: menu3.ID, Order: 30},
		},
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.SuccessItems, 3)

	// Verify database state - orders should be updated
	updatedMenu1 := GetMenuByID(t, db, menu1.ID)
	updatedMenu2 := GetMenuByID(t, db, menu2.ID)
	updatedMenu3 := GetMenuByID(t, db, menu3.ID)

	assert.Equal(t, 10, updatedMenu1.Order)
	assert.Equal(t, 20, updatedMenu2.Order)
	assert.Equal(t, 30, updatedMenu3.Order)
}

func TestMenuService_BatchUpdateMenuOrder_PartialNotFound(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create two menus
	menu1 := CreateTestMenuWithOrder(t, db, "Menu1", "/menu1", nil, 0, false)
	menu2 := CreateTestMenuWithOrder(t, db, "Menu2", "/menu2", nil, 1, false)
	nonExistentID := uint64(999999)

	// Execute batch order update with mixed valid and invalid IDs
	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{
			{MenuID: menu1.ID, Order: 10},
			{MenuID: menu2.ID, Order: 20},
			{MenuID: nonExistentID, Order: 30},
		},
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 2, result.SuccessCount)
	assert.Equal(t, 1, result.FailedCount)
	assert.Len(t, result.FailedItems, 1)

	// Verify failed item
	assert.Equal(t, nonExistentID, result.FailedItems[0].ID)
	assert.Contains(t, result.FailedItems[0].Message, "menu not found")

	// Verify valid menus were updated
	updatedMenu1 := GetMenuByID(t, db, menu1.ID)
	updatedMenu2 := GetMenuByID(t, db, menu2.ID)
	assert.Equal(t, 10, updatedMenu1.Order)
	assert.Equal(t, 20, updatedMenu2.Order)
}

func TestMenuService_BatchUpdateMenuOrder_DuplicateOrderValues(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create menus
	menu1 := CreateTestMenu(t, db, "Menu1", "/menu1", nil)
	menu2 := CreateTestMenu(t, db, "Menu2", "/menu2", nil)

	// Execute with duplicate order values
	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{
			{MenuID: menu1.ID, Order: 10},
			{MenuID: menu2.ID, Order: 10}, // Duplicate order value
		},
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "duplicate order value")
}

func TestMenuService_BatchUpdateMenuOrder_EmptyOrders(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Execute with empty orders
	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{},
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMenuService_BatchUpdateMenuOrder_TooManyOrders(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create 101 order updates
	orders := make([]adminservice.MenuOrderUpdate, 101)
	for i := 0; i < 101; i++ {
		orders[i] = adminservice.MenuOrderUpdate{
			MenuID: uint64(i + 1),
			Order:  i,
		}
	}

	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: orders,
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "maximum 100")
}

func TestMenuService_BatchUpdateMenuOrder_WithParentChildRelationship(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create parent and child menus
	parentMenu := CreateTestMenuWithOrder(t, db, "Parent", "/parent", nil, 0, false)
	childMenu1 := CreateTestMenuWithOrder(t, db, "Child1", "/parent/child1", &parentMenu.ID, 0, false)
	childMenu2 := CreateTestMenuWithOrder(t, db, "Child2", "/parent/child2", &parentMenu.ID, 1, false)

	// Update orders for both parent and children
	input := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{
			{MenuID: parentMenu.ID, Order: 5},
			{MenuID: childMenu1.ID, Order: 10},
			{MenuID: childMenu2.ID, Order: 20},
		},
	}
	result, err := menuService.BatchUpdateMenuOrder(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, 3, result.TotalCount)
	assert.Equal(t, 3, result.SuccessCount)
	assert.Equal(t, 0, result.FailedCount)

	// Verify orders - children can have same order values as parent
	// since they're in different hierarchy levels
	updatedParent := GetMenuByID(t, db, parentMenu.ID)
	updatedChild1 := GetMenuByID(t, db, childMenu1.ID)
	updatedChild2 := GetMenuByID(t, db, childMenu2.ID)

	assert.Equal(t, 5, updatedParent.Order)
	assert.Equal(t, 10, updatedChild1.Order)
	assert.Equal(t, 20, updatedChild2.Order)
}

// ============================================================================
// Combined Tests
// ============================================================================

func TestMenuService_BatchOperations_ComplexScenario(t *testing.T) {
	SkipIfNoTestDB(t)
	db := SetupTestDB(t)
	ctx := context.Background()

	// Create repository and service
	menuRepo := adminrepository.NewMenuRepository(db)
	menuService := adminservice.NewMenuService(menuRepo)

	// Create complex menu hierarchy
	// Root level
	menu1 := CreateTestMenuWithOrder(t, db, "Dashboard", "/dashboard", nil, 0, false)
	menu2 := CreateTestMenuWithOrder(t, db, "Users", "/users", nil, 1, true)
	menu3 := CreateTestMenuWithOrder(t, db, "Settings", "/settings", nil, 2, false)

	// Children under Users
	userList := CreateTestMenuWithOrder(t, db, "User List", "/users/list", &menu2.ID, 0, false)
	userCreate := CreateTestMenuWithOrder(t, db, "Create User", "/users/create", &menu2.ID, 1, false)

	// Leaf menu (no children, for deletion)
	leafMenu := CreateTestMenuWithOrder(t, db, "Temp Menu", "/temp", nil, 99, true)

	// Test 1: Batch update status - enable hidden menus
	statusInput := adminservice.BatchUpdateMenuStatusInput{
		MenuIDs: []uint64{menu2.ID, leafMenu.ID},
		Status:  "enabled",
	}
	statusResult, err := menuService.BatchUpdateMenuStatus(ctx, statusInput)
	require.NoError(t, err)
	assert.Equal(t, 2, statusResult.SuccessCount)

	// Verify
	assert.False(t, GetMenuByID(t, db, menu2.ID).Hidden)
	assert.False(t, GetMenuByID(t, db, leafMenu.ID).Hidden)

	// Test 2: Batch update order
	orderInput := adminservice.BatchUpdateMenuOrderInput{
		MenuOrders: []adminservice.MenuOrderUpdate{
			{MenuID: menu1.ID, Order: 100},
			{MenuID: menu3.ID, Order: 200},
		},
	}
	orderResult, err := menuService.BatchUpdateMenuOrder(ctx, orderInput)
	require.NoError(t, err)
	assert.Equal(t, 2, orderResult.SuccessCount)

	// Verify
	assert.Equal(t, 100, GetMenuByID(t, db, menu1.ID).Order)
	assert.Equal(t, 200, GetMenuByID(t, db, menu3.ID).Order)

	// Test 3: Batch delete - leaf menu succeeds, parent with children fails
	deleteInput := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{menu2.ID, leafMenu.ID},
	}
	deleteResult, err := menuService.BatchDeleteMenus(ctx, deleteInput)
	require.NoError(t, err)
	assert.Equal(t, 2, deleteResult.TotalCount)
	assert.Equal(t, 1, deleteResult.SuccessCount)
	assert.Equal(t, 1, deleteResult.FailedCount)

	// Verify leaf menu deleted, parent still exists
	assert.False(t, MenuExists(t, db, leafMenu.ID))
	assert.True(t, MenuExists(t, db, menu2.ID))

	// Test 4: Delete children, then parent
	deleteInput2 := adminservice.BatchDeleteMenusInput{
		MenuIDs: []uint64{userList.ID, userCreate.ID, menu2.ID},
	}
	deleteResult2, err := menuService.BatchDeleteMenus(ctx, deleteInput2)
	require.NoError(t, err)
	assert.Equal(t, 3, deleteResult2.TotalCount)
	assert.Equal(t, 3, deleteResult2.SuccessCount)
	assert.Equal(t, 0, deleteResult2.FailedCount)

	// Verify all deleted
	assert.False(t, MenuExists(t, db, userList.ID))
	assert.False(t, MenuExists(t, db, userCreate.ID))
	assert.False(t, MenuExists(t, db, menu2.ID))
}

package permission

import (
	"context"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func testContext() context.Context {
	return context.Background()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Migrate in correct order
	if err := db.AutoMigrate(&model.Permission{}, &model.RoleModel{}, &model.User{}, &model.RolePermission{}, &model.UserRole{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	// Ensure columns used by repository filters exist even if the struct does not define them.
	for _, stmt := range []string{
		"ALTER TABLE permissions ADD COLUMN resource TEXT",
		"ALTER TABLE permissions ADD COLUMN action TEXT",
	} {
		if execErr := db.Exec(stmt).Error; execErr != nil && !strings.Contains(strings.ToLower(execErr.Error()), "duplicate column name") {
			t.Fatalf("failed to alter permissions table: %v", execErr)
		}
	}

	return db
}

// permissionWithResource is a helper struct for tests to seed resource/action columns.
type permissionWithResource struct {
	model.Permission
	Resource string `gorm:"column:resource"`
	Action   string `gorm:"column:action"`
}

func (permissionWithResource) TableName() string {
	return "permissions"
}

func TestPermissionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{
		Method:      "GET",
		Path:        "/api/users",
		Code:        "users.read",
		Group:       "/api/users",
		Description: "读取用户列表",
	}

	err := repo.Create(testContext(), permission)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if permission.ID == 0 {
		t.Error("expected permission ID to be set")
	}

	// Verify created
	retrieved, err := repo.Get(testContext(), permission.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Code != "users.read" {
		t.Errorf("expected code 'users.read', got %s", retrieved.Code)
	}
}

func TestPermissionRepository_CreateMultiple(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permissions := []*model.Permission{
		{Method: "GET", Path: "/api/games", Code: "games.read", Group: "/api/games"},
		{Method: "POST", Path: "/api/games", Code: "games.create", Group: "/api/games"},
		{Method: "PUT", Path: "/api/games/:id", Code: "games.update", Group: "/api/games"},
	}

	for _, p := range permissions {
		err := repo.Create(testContext(), p)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// Verify all created
	list, _ := repo.List(testContext())
	if len(list) != 3 {
		t.Errorf("expected 3 permissions, got %d", len(list))
	}
}

func TestPermissionRepository_CreateBatch(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db).(*permissionRepository)

	perms := []model.Permission{
		{Method: "GET", Path: "/api/batch/orders", Code: "batch.orders.read"},
		{Method: "POST", Path: "/api/batch/orders", Code: "batch.orders.create"},
	}

	if err := repo.CreateBatch(testContext(), perms); err != nil {
		t.Fatalf("CreateBatch failed: %v", err)
	}

	list, err := repo.List(testContext())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(list))
	}

	if err := repo.CreateBatch(testContext(), []model.Permission{}); err != nil {
		t.Fatalf("CreateBatch should ignore empty slice: %v", err)
	}
}

func TestPermissionRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{Method: "DELETE", Path: "/api/orders/:id", Code: "orders.delete"}
	_ = repo.Create(testContext(), permission)

	t.Run("Get existing permission", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), permission.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.Code != "orders.delete" {
			t.Errorf("expected code 'orders.delete', got %s", retrieved.Code)
		}
	})

	t.Run("Get non-existent permission", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPermissionRepository_GetByMethodAndPath(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{Method: "POST", Path: "/api/players", Code: "players.create"}
	_ = repo.Create(testContext(), permission)

	t.Run("Get existing by method and path", func(t *testing.T) {
		retrieved, err := repo.GetByMethodAndPath(testContext(), "POST", "/api/players")
		if err != nil {
			t.Fatalf("GetByMethodAndPath failed: %v", err)
		}

		if retrieved.Code != "players.create" {
			t.Errorf("expected code 'players.create', got %s", retrieved.Code)
		}
	})

	t.Run("Get non-existent by method and path", func(t *testing.T) {
		_, err := repo.GetByMethodAndPath(testContext(), "GET", "/api/nonexistent")
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPermissionRepository_GetByResource(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	record := &permissionWithResource{
		Permission: model.Permission{
			Method: "GET",
			Path:   "/api/resources/orders",
			Code:   "orders.resources.read",
			Group:  "/api/resources",
		},
		Resource: "orders",
		Action:   "read",
	}
	if err := db.Create(record).Error; err != nil {
		t.Fatalf("failed to seed permission with resource: %v", err)
	}

	t.Run("Get existing by resource/action", func(t *testing.T) {
		retrieved, err := repo.GetByResource(testContext(), "orders", "read")
		if err != nil {
			t.Fatalf("GetByResource failed: %v", err)
		}
		if retrieved.Code != "orders.resources.read" {
			t.Errorf("expected code 'orders.resources.read', got %s", retrieved.Code)
		}
	})

	t.Run("GetByResource not found", func(t *testing.T) {
		_, err := repo.GetByResource(testContext(), "orders", "write")
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPermissionRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{
		Method:      "GET",
		Path:        "/api/stats",
		Code:        "stats.read",
		Description: "统计数据读取",
	}
	_ = repo.Create(testContext(), permission)

	t.Run("Update existing permission", func(t *testing.T) {
		permission.Description = "获取统计数据"
		permission.Group = "/api/stats"

		err := repo.Update(testContext(), permission)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Verify updated
		retrieved, _ := repo.Get(testContext(), permission.ID)
		if retrieved.Description != "获取统计数据" {
			t.Errorf("expected updated description, got %s", retrieved.Description)
		}
		if retrieved.Group != "/api/stats" {
			t.Errorf("expected group '/api/stats', got %s", retrieved.Group)
		}
	})

	t.Run("Update non-existent permission", func(t *testing.T) {
		nonExistent := &model.Permission{Base: model.Base{ID: 99999}, Code: "fake"}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPermissionRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{Method: "DELETE", Path: "/api/payments/:id", Code: "payments.delete"}
	_ = repo.Create(testContext(), permission)

	t.Run("Delete existing permission", func(t *testing.T) {
		err := repo.Delete(testContext(), permission.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify deleted
		_, err = repo.Get(testContext(), permission.ID)
		if err != repository.ErrNotFound {
			t.Error("expected permission to be deleted")
		}
	})

	t.Run("Delete non-existent permission", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestPermissionRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	t.Run("List empty", func(t *testing.T) {
		permissions, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(permissions) != 0 {
			t.Errorf("expected 0 permissions, got %d", len(permissions))
		}
	})

	// Create test data
	permissions := []*model.Permission{
		{Method: "GET", Path: "/api/users", Code: "users.read", Group: "/api/users"},
		{Method: "POST", Path: "/api/users", Code: "users.create", Group: "/api/users"},
		{Method: "GET", Path: "/api/games", Code: "games.read", Group: "/api/games"},
	}
	for _, p := range permissions {
		_ = repo.Create(testContext(), p)
	}

	t.Run("List all permissions", func(t *testing.T) {
		result, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(result) != 3 {
			t.Errorf("expected 3 permissions, got %d", len(result))
		}
	})
}

func TestPermissionRepository_ListPaged(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	// Create test data
	for i := 1; i <= 15; i++ {
		permission := &model.Permission{
			Method: "GET",
			Path:   "/api/test/" + string(rune('0'+i)),
			Code:   "test.read." + string(rune('0'+i)),
		}
		_ = repo.Create(testContext(), permission)
	}

	t.Run("First page", func(t *testing.T) {
		permissions, total, err := repo.ListPaged(testContext(), 1, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(permissions) != 10 {
			t.Errorf("expected 10 permissions, got %d", len(permissions))
		}
	})

	t.Run("Second page", func(t *testing.T) {
		permissions, total, err := repo.ListPaged(testContext(), 2, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(permissions) != 5 {
			t.Errorf("expected 5 permissions on second page, got %d", len(permissions))
		}
	})
}

func TestPermissionRepository_ListPagedWithFilter(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	seeds := []*model.Permission{
		{Method: "GET", Path: "/api/orders", Code: "orders.read", Group: "/api/orders", Description: "list orders"},
		{Method: "POST", Path: "/api/orders", Code: "orders.create", Group: "/api/orders"},
		{Method: "DELETE", Path: "/api/users/:id", Code: "users.delete", Group: "/api/users"},
	}
	for _, p := range seeds {
		if err := repo.Create(testContext(), p); err != nil {
			t.Fatalf("failed to seed permission: %v", err)
		}
	}

	t.Run("filter by keyword", func(t *testing.T) {
		result, total, err := repo.ListPagedWithFilter(testContext(), 1, 10, "orders", "", "")
		if err != nil {
			t.Fatalf("ListPagedWithFilter failed: %v", err)
		}
		if total != 2 {
			t.Fatalf("expected total 2, got %d", total)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 records, got %d", len(result))
		}
	})

	t.Run("filter by method and group", func(t *testing.T) {
		result, total, err := repo.ListPagedWithFilter(testContext(), 1, 10, "", "POST", "/api/orders")
		if err != nil {
			t.Fatalf("ListPagedWithFilter failed: %v", err)
		}
		if total != 1 {
			t.Fatalf("expected total 1, got %d", total)
		}
		if len(result) != 1 || result[0].Code != "orders.create" {
			t.Fatalf("expected orders.create record, got %+v", result)
		}
	})
}

func TestPermissionRepository_ListByGroup(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permissions := []*model.Permission{
		{Method: "GET", Path: "/api/users", Code: "users.read", Group: "/api/users"},
		{Method: "POST", Path: "/api/users", Code: "users.create", Group: "/api/users"},
		{Method: "GET", Path: "/api/games", Code: "games.read", Group: "/api/games"},
		{Method: "DELETE", Path: "/api/games/:id", Code: "games.delete", Group: "/api/games"},
	}
	for _, p := range permissions {
		_ = repo.Create(testContext(), p)
	}

	grouped, err := repo.ListByGroup(testContext())
	if err != nil {
		t.Fatalf("ListByGroup failed: %v", err)
	}

	if len(grouped) != 2 {
		t.Errorf("expected 2 groups, got %d", len(grouped))
	}

	if len(grouped["/api/users"]) != 2 {
		t.Errorf("expected 2 permissions in /api/users group, got %d", len(grouped["/api/users"]))
	}

	if len(grouped["/api/games"]) != 2 {
		t.Errorf("expected 2 permissions in /api/games group, got %d", len(grouped["/api/games"]))
	}
}

func TestPermissionRepository_ListGroups(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permissions := []*model.Permission{
		{Method: "GET", Path: "/api/users", Code: "users.read", Group: "/api/users"},
		{Method: "GET", Path: "/api/games", Code: "games.read", Group: "/api/games"},
		{Method: "GET", Path: "/api/orders", Code: "orders.read", Group: "/api/orders"},
		{Method: "POST", Path: "/api/games", Code: "games.create", Group: "/api/games"},
	}
	for _, p := range permissions {
		_ = repo.Create(testContext(), p)
	}

	groups, err := repo.ListGroups(testContext())
	if err != nil {
		t.Fatalf("ListGroups failed: %v", err)
	}

	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}

	// Verify groups are present
	groupMap := make(map[string]bool)
	for _, g := range groups {
		groupMap[g] = true
	}

	if !groupMap["/api/users"] || !groupMap["/api/games"] || !groupMap["/api/orders"] {
		t.Error("expected all three groups to be present")
	}
}

func TestPermissionRepository_UpsertByMethodPath(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	t.Run("Insert new permission", func(t *testing.T) {
		permission := &model.Permission{
			Method:      "GET",
			Path:        "/api/new",
			Code:        "new.read",
			Description: "新权限",
		}

		err := repo.UpsertByMethodPath(testContext(), permission)
		if err != nil {
			t.Fatalf("UpsertByMethodPath failed: %v", err)
		}

		if permission.ID == 0 {
			t.Error("expected permission ID to be set")
		}

		// Verify created
		retrieved, _ := repo.GetByMethodAndPath(testContext(), "GET", "/api/new")
		if retrieved.Code != "new.read" {
			t.Errorf("expected code 'new.read', got %s", retrieved.Code)
		}
	})

	t.Run("Update existing permission", func(t *testing.T) {
		// First create
		permission := &model.Permission{
			Method:      "POST",
			Path:        "/api/update",
			Code:        "update.create",
			Description: "原描述",
		}
		_ = repo.Create(testContext(), permission)

		// Then upsert with same method/path but different code
		upsertPerm := &model.Permission{
			Method:      "POST",
			Path:        "/api/update",
			Code:        "update.create.v2",
			Description: "新描述",
		}

		err := repo.UpsertByMethodPath(testContext(), upsertPerm)
		if err != nil {
			t.Fatalf("UpsertByMethodPath failed: %v", err)
		}

		// Verify updated (ID should be preserved)
		retrieved, _ := repo.GetByMethodAndPath(testContext(), "POST", "/api/update")
		if retrieved.ID != permission.ID {
			t.Errorf("expected ID to be preserved: %d, got %d", permission.ID, retrieved.ID)
		}
		if retrieved.Code != "update.create.v2" {
			t.Errorf("expected code 'update.create.v2', got %s", retrieved.Code)
		}
		if retrieved.Description != "新描述" {
			t.Errorf("expected description '新描述', got %s", retrieved.Description)
		}
	})
}

func TestPermissionRepository_ListByRoleID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	// Create permissions
	perm1 := &model.Permission{Method: "GET", Path: "/api/users", Code: "users.read"}
	perm2 := &model.Permission{Method: "POST", Path: "/api/users", Code: "users.create"}
	perm3 := &model.Permission{Method: "GET", Path: "/api/games", Code: "games.read"}
	_ = repo.Create(testContext(), perm1)
	_ = repo.Create(testContext(), perm2)
	_ = repo.Create(testContext(), perm3)

	// Create role
	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	db.Create(role)

	// Assign permissions to role
	db.Create(&model.RolePermission{RoleID: role.ID, PermissionID: perm1.ID})
	db.Create(&model.RolePermission{RoleID: role.ID, PermissionID: perm2.ID})

	// List permissions by role ID
	permissions, err := repo.ListByRoleID(testContext(), role.ID)
	if err != nil {
		t.Fatalf("ListByRoleID failed: %v", err)
	}

	if len(permissions) != 2 {
		t.Errorf("expected 2 permissions, got %d", len(permissions))
	}

	// Verify permissions
	codes := make(map[string]bool)
	for _, p := range permissions {
		codes[p.Code] = true
	}

	if !codes["users.read"] || !codes["users.create"] {
		t.Error("expected users.read and users.create permissions")
	}
}

func TestPermissionRepository_ListByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	// Create permissions
	perm1 := &model.Permission{Method: "GET", Path: "/api/users", Code: "users.read"}
	perm2 := &model.Permission{Method: "POST", Path: "/api/users", Code: "users.create"}
	perm3 := &model.Permission{Method: "GET", Path: "/api/games", Code: "games.read"}
	perm4 := &model.Permission{Method: "POST", Path: "/api/games", Code: "games.create"}
	_ = repo.Create(testContext(), perm1)
	_ = repo.Create(testContext(), perm2)
	_ = repo.Create(testContext(), perm3)
	_ = repo.Create(testContext(), perm4)

	// Create roles
	role1 := &model.RoleModel{Name: "User Manager", Slug: "user_manager"}
	role2 := &model.RoleModel{Name: "Game Manager", Slug: "game_manager"}
	db.Create(role1)
	db.Create(role2)

	// Assign permissions to roles
	db.Create(&model.RolePermission{RoleID: role1.ID, PermissionID: perm1.ID})
	db.Create(&model.RolePermission{RoleID: role1.ID, PermissionID: perm2.ID})
	db.Create(&model.RolePermission{RoleID: role2.ID, PermissionID: perm3.ID})
	db.Create(&model.RolePermission{RoleID: role2.ID, PermissionID: perm4.ID})

	// Create user and assign roles
	user := &model.User{Email: "admin@example.com", Name: "Admin"}
	db.Create(user)
	db.Create(&model.UserRole{UserID: user.ID, RoleID: role1.ID})
	db.Create(&model.UserRole{UserID: user.ID, RoleID: role2.ID})

	// List permissions by user ID
	permissions, err := repo.ListByUserID(testContext(), user.ID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(permissions) != 4 {
		t.Errorf("expected 4 permissions, got %d", len(permissions))
	}

	// Verify all permissions
	codes := make(map[string]bool)
	for _, p := range permissions {
		codes[p.Code] = true
	}

	if !codes["users.read"] || !codes["users.create"] || !codes["games.read"] || !codes["games.create"] {
		t.Error("expected all four permissions to be present")
	}
}

func TestPermissionRepository_CompleteWorkflow(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	// Create
	permission := &model.Permission{
		Method:      "GET",
		Path:        "/api/complete",
		Code:        "complete.test",
		Group:       "/api/complete",
		Description: "完整测试",
	}
	err := repo.Create(testContext(), permission)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Read by ID
	retrieved, err := repo.Get(testContext(), permission.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.Code != "complete.test" {
		t.Errorf("expected code 'complete.test', got %s", retrieved.Code)
	}

	// Read by method and path
	byMethodPath, err := repo.GetByMethodAndPath(testContext(), "GET", "/api/complete")
	if err != nil {
		t.Fatalf("GetByMethodAndPath failed: %v", err)
	}
	if byMethodPath.ID != permission.ID {
		t.Error("expected same permission ID")
	}

	// Update
	permission.Description = "更新后的描述"
	err = repo.Update(testContext(), permission)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	updated, _ := repo.Get(testContext(), permission.ID)
	if updated.Description != "更新后的描述" {
		t.Errorf("expected updated description, got %s", updated.Description)
	}

	// List
	list, _ := repo.List(testContext())
	if len(list) != 1 {
		t.Errorf("expected 1 permission in list, got %d", len(list))
	}

	// Delete
	err = repo.Delete(testContext(), permission.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.Get(testContext(), permission.ID)
	if err != repository.ErrNotFound {
		t.Error("expected permission to be deleted")
	}
}

// TestPermissionRepository_GetByCode 测试通过code获取权限
func TestPermissionRepository_GetByCode(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	permission := &model.Permission{
		Method: "GET",
		Path:   "/api/users",
		Code:   "users.read",
		Group:  "/api/users",
	}
	err := repo.Create(testContext(), permission)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	t.Run("通过code获取存在的权限", func(t *testing.T) {
		result, err := repo.GetByCode(testContext(), "users.read")
		if err != nil {
			t.Fatalf("GetByCode failed: %v", err)
		}
		if result.Code != "users.read" {
			t.Errorf("expected code 'users.read', got '%s'", result.Code)
		}
		if result.Path != "/api/users" {
			t.Errorf("expected path '/api/users', got '%s'", result.Path)
		}
	})

	t.Run("通过code获取不存在的权限", func(t *testing.T) {
		_, err := repo.GetByCode(testContext(), "nonexistent.code")
		if err == nil {
			t.Error("expected error for nonexistent code")
		}
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestPermissionRepository_ListByGroup_EdgeCases 测试ListByGroup的边界条件
func TestPermissionRepository_ListByGroup_EdgeCases(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	t.Run("空数据库返回空map", func(t *testing.T) {
		grouped, err := repo.ListByGroup(testContext())
		if err != nil {
			t.Fatalf("ListByGroup failed: %v", err)
		}
		if grouped == nil {
			t.Error("expected non-nil map")
		}
		if len(grouped) != 0 {
			t.Errorf("expected empty map, got %d groups", len(grouped))
		}
	})

	t.Run("权限没有group字段时", func(t *testing.T) {
		permission := &model.Permission{
			Method: "GET",
			Path:   "/api/test",
			Code:   "test.read",
			Group:  "", // 空group
		}
		err := repo.Create(testContext(), permission)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		grouped, err := repo.ListByGroup(testContext())
		if err != nil {
			t.Fatalf("ListByGroup failed: %v", err)
		}
		// 空group应该被包含在结果中
		if len(grouped) == 0 {
			t.Error("expected at least one group (empty string group)")
		}
	})
}

// TestPermissionRepository_GetUserPermissions_EdgeCases 测试GetUserPermissions（ListByUserID）的边界条件
func TestPermissionRepository_GetUserPermissions_EdgeCases(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPermissionRepository(db)

	t.Run("查询不存在的用户权限", func(t *testing.T) {
		permissions, err := repo.ListByUserID(testContext(), 99999)
		if err != nil {
			t.Fatalf("ListByUserID should not error for nonexistent user, got: %v", err)
		}
		if len(permissions) != 0 {
			t.Errorf("expected empty permissions for nonexistent user, got %d", len(permissions))
		}
	})

	t.Run("用户没有分配角色时返回空列表", func(t *testing.T) {
		// 创建用户但不分配角色
		user := &model.User{Email: "no-role@example.com", Name: "No Role User"}
		db.Create(user)

		permissions, err := repo.ListByUserID(testContext(), user.ID)
		if err != nil {
			t.Fatalf("ListByUserID failed: %v", err)
		}
		if len(permissions) != 0 {
			t.Errorf("expected empty permissions for user without roles, got %d", len(permissions))
		}
	})

	t.Run("角色没有分配权限时返回空列表", func(t *testing.T) {
		// 创建权限
		perm := &model.Permission{Method: "GET", Path: "/api/test", Code: "test.read"}
		_ = repo.Create(testContext(), perm)

		// 创建角色但不分配权限
		role := &model.RoleModel{Name: "Empty Role", Slug: "empty_role"}
		db.Create(role)

		// 创建用户并分配角色
		user := &model.User{Email: "empty-role@example.com", Name: "Empty Role User"}
		db.Create(user)
		db.Create(&model.UserRole{UserID: user.ID, RoleID: role.ID})

		permissions, err := repo.ListByUserID(testContext(), user.ID)
		if err != nil {
			t.Fatalf("ListByUserID failed: %v", err)
		}
		if len(permissions) != 0 {
			t.Errorf("expected empty permissions for role without permissions, got %d", len(permissions))
		}
	})
}

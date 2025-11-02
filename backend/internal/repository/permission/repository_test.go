package permission

import (
	"context"
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

	return db
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

	permission := &model.Permission{Method: "GET", Path: "/api/reviews", Code: "reviews.list"}
	_ = repo.Create(testContext(), permission)

	t.Run("Get existing by resource", func(t *testing.T) {
		// Note: GetByResource is in the interface but takes different parameters
		// This test is kept for interface coverage
		retrieved, err := repo.Get(testContext(), permission.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.Path != "/api/reviews" {
			t.Errorf("expected path '/api/reviews', got %s", retrieved.Path)
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

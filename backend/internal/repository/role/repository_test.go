package role

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

	// Migrate in correct order to handle foreign keys
	if err := db.AutoMigrate(&model.User{}, &model.RoleModel{}, &model.Permission{}, &model.RolePermission{}, &model.UserRole{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestNewRoleRepository 测试构造函数。
func TestNewRoleRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewRoleRepository(db)
	if repo == nil {
		t.Fatal("NewRoleRepository returned nil")
	}

	if _, ok := repo.(*roleRepository); !ok {
		t.Errorf("expected *roleRepository, got %T", repo)
	}
}

func TestRoleRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{
		Name:        "Admin",
		Slug:        "admin",
		Description: "Administrator role",
		IsSystem:    true,
	}

	err := repo.Create(testContext(), role)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if role.ID == 0 {
		t.Error("expected ID to be set after create")
	}
}

func TestRoleRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin", IsSystem: true}
	_ = repo.Create(testContext(), role)

	t.Run("Get existing role", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), role.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != role.ID {
			t.Errorf("expected ID %d, got %d", role.ID, retrieved.ID)
		}
		if retrieved.Slug != "admin" {
			t.Errorf("expected slug 'admin', got %s", retrieved.Slug)
		}
	})

	t.Run("Get nonexistent role", func(t *testing.T) {
		_, err := repo.Get(testContext(), 999999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestRoleRepository_GetBySlug(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "User", Slug: "user", IsSystem: true}
	_ = repo.Create(testContext(), role)

	t.Run("Get by valid slug", func(t *testing.T) {
		retrieved, err := repo.GetBySlug(testContext(), "user")
		if err != nil {
			t.Fatalf("GetBySlug failed: %v", err)
		}
		if retrieved.Slug != "user" {
			t.Errorf("expected slug 'user', got %s", retrieved.Slug)
		}
	})

	t.Run("Get by invalid slug", func(t *testing.T) {
		_, err := repo.GetBySlug(testContext(), "nonexistent")
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestRoleRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "User", Slug: "user"}
	_ = repo.Create(testContext(), role)

	role.Name = "Updated User"
	role.Description = "Updated description"
	err := repo.Update(testContext(), role)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	retrieved, _ := repo.Get(testContext(), role.ID)
	if retrieved.Name != "Updated User" {
		t.Errorf("expected name 'Updated User', got %s", retrieved.Name)
	}
}

func TestRoleRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	t.Run("Delete non-system role", func(t *testing.T) {
		role := &model.RoleModel{Name: "Custom", Slug: "custom", IsSystem: false}
		_ = repo.Create(testContext(), role)

		err := repo.Delete(testContext(), role.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repo.Get(testContext(), role.ID)
		if err != repository.ErrNotFound {
			t.Error("expected role to be deleted")
		}
	})

	t.Run("Cannot delete system role", func(t *testing.T) {
		role := &model.RoleModel{Name: "Admin", Slug: "admin", IsSystem: true}
		_ = repo.Create(testContext(), role)

		err := repo.Delete(testContext(), role.ID)
		if err == nil {
			t.Error("expected error when deleting system role")
		}
		if err.Error() != "cannot delete system role" {
			t.Errorf("expected 'cannot delete system role' error, got %v", err)
		}
	})

	t.Run("Delete nonexistent role", func(t *testing.T) {
		err := repo.Delete(testContext(), 999999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestRoleRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	// Create test roles
	roles := []*model.RoleModel{
		{Name: "Admin", Slug: "admin", IsSystem: true},
		{Name: "User", Slug: "user", IsSystem: true},
		{Name: "Custom", Slug: "custom", IsSystem: false},
	}
	for _, r := range roles {
		_ = repo.Create(testContext(), r)
	}

	list, err := repo.List(testContext())
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) < 3 {
		t.Errorf("expected at least 3 roles, got %d", len(list))
	}

	// Verify system roles come first
	if len(list) >= 3 && !list[0].IsSystem {
		t.Error("expected system roles to be listed first")
	}
}

func TestRoleRepository_ListPaged(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	// Create test roles
	for i := 1; i <= 5; i++ {
		role := &model.RoleModel{Name: "Role", Slug: "role-" + string(rune(i+'0'))}
		_ = repo.Create(testContext(), role)
	}

	roles, total, err := repo.ListPaged(testContext(), 1, 2)
	if err != nil {
		t.Fatalf("ListPaged failed: %v", err)
	}

	if len(roles) != 2 {
		t.Errorf("expected 2 roles, got %d", len(roles))
	}

	if total < 5 {
		t.Errorf("expected total >= 5, got %d", total)
	}
}

func TestRoleRepository_AssignPermissions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	// Create test permissions
	perm1 := &model.Permission{Code: "read", Method: "GET", Path: "/api/read", Description: "Read permission"}
	perm2 := &model.Permission{Code: "write", Method: "POST", Path: "/api/write", Description: "Write permission"}
	db.Create(perm1)
	db.Create(perm2)

	t.Run("Assign permissions", func(t *testing.T) {
		err := repo.AssignPermissions(testContext(), role.ID, []uint64{perm1.ID, perm2.ID})
		if err != nil {
			t.Fatalf("AssignPermissions failed: %v", err)
		}

		// Verify permissions assigned
		retrieved, _ := repo.GetWithPermissions(testContext(), role.ID)
		if len(retrieved.Permissions) != 2 {
			t.Errorf("expected 2 permissions, got %d", len(retrieved.Permissions))
		}
	})

	t.Run("Replace permissions", func(t *testing.T) {
		// Assign only perm1, should replace perm2
		err := repo.AssignPermissions(testContext(), role.ID, []uint64{perm1.ID})
		if err != nil {
			t.Fatalf("AssignPermissions failed: %v", err)
		}

		retrieved, _ := repo.GetWithPermissions(testContext(), role.ID)
		if len(retrieved.Permissions) != 1 {
			t.Errorf("expected 1 permission after replacement, got %d", len(retrieved.Permissions))
		}
	})
}

func TestRoleRepository_AddPermissions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "User", Slug: "user"}
	_ = repo.Create(testContext(), role)

	perm := &model.Permission{Code: "read", Method: "GET", Path: "/api/read"}
	db.Create(perm)

	err := repo.AddPermissions(testContext(), role.ID, []uint64{perm.ID})
	if err != nil {
		t.Fatalf("AddPermissions failed: %v", err)
	}

	// Add same permission again (should not error)
	err = repo.AddPermissions(testContext(), role.ID, []uint64{perm.ID})
	if err != nil {
		t.Fatalf("AddPermissions duplicate failed: %v", err)
	}
}

func TestRoleRepository_RemovePermissions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	perm1 := &model.Permission{Code: "read", Method: "GET", Path: "/api/read"}
	perm2 := &model.Permission{Code: "write", Method: "POST", Path: "/api/write"}
	db.Create(perm1)
	db.Create(perm2)

	// Assign both permissions
	_ = repo.AssignPermissions(testContext(), role.ID, []uint64{perm1.ID, perm2.ID})

	// Remove one permission
	err := repo.RemovePermissions(testContext(), role.ID, []uint64{perm2.ID})
	if err != nil {
		t.Fatalf("RemovePermissions failed: %v", err)
	}

	retrieved, _ := repo.GetWithPermissions(testContext(), role.ID)
	if len(retrieved.Permissions) != 1 {
		t.Errorf("expected 1 permission after removal, got %d", len(retrieved.Permissions))
	}
}

func TestRoleRepository_AssignToUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role1 := &model.RoleModel{Name: "Admin", Slug: "admin"}
	role2 := &model.RoleModel{Name: "User", Slug: "user"}
	_ = repo.Create(testContext(), role1)
	_ = repo.Create(testContext(), role2)

	userID := uint64(1)

	t.Run("Assign roles to user", func(t *testing.T) {
		err := repo.AssignToUser(testContext(), userID, []uint64{role1.ID, role2.ID})
		if err != nil {
			t.Fatalf("AssignToUser failed: %v", err)
		}

		roles, _ := repo.ListByUserID(testContext(), userID)
		if len(roles) != 2 {
			t.Errorf("expected 2 roles, got %d", len(roles))
		}
	})

	t.Run("Replace user roles", func(t *testing.T) {
		err := repo.AssignToUser(testContext(), userID, []uint64{role1.ID})
		if err != nil {
			t.Fatalf("AssignToUser failed: %v", err)
		}

		roles, _ := repo.ListByUserID(testContext(), userID)
		if len(roles) != 1 {
			t.Errorf("expected 1 role after replacement, got %d", len(roles))
		}
	})
}

func TestRoleRepository_RemoveFromUser(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	userID := uint64(1)
	_ = repo.AssignToUser(testContext(), userID, []uint64{role.ID})

	err := repo.RemoveFromUser(testContext(), userID, []uint64{role.ID})
	if err != nil {
		t.Fatalf("RemoveFromUser failed: %v", err)
	}

	roles, _ := repo.ListByUserID(testContext(), userID)
	if len(roles) != 0 {
		t.Errorf("expected 0 roles after removal, got %d", len(roles))
	}
}

func TestRoleRepository_CheckUserHasRole(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	userID := uint64(1)
	_ = repo.AssignToUser(testContext(), userID, []uint64{role.ID})

	t.Run("User has role", func(t *testing.T) {
		has, err := repo.CheckUserHasRole(testContext(), userID, "admin")
		if err != nil {
			t.Fatalf("CheckUserHasRole failed: %v", err)
		}
		if !has {
			t.Error("expected user to have admin role")
		}
	})

	t.Run("User does not have role", func(t *testing.T) {
		has, err := repo.CheckUserHasRole(testContext(), userID, "nonexistent")
		if err != nil {
			t.Fatalf("CheckUserHasRole failed: %v", err)
		}
		if has {
			t.Error("expected user not to have nonexistent role")
		}
	})
}

func TestRoleRepository_GetWithPermissions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	perm := &model.Permission{Code: "read", Method: "GET", Path: "/api/read"}
	db.Create(perm)
	_ = repo.AssignPermissions(testContext(), role.ID, []uint64{perm.ID})

	retrieved, err := repo.GetWithPermissions(testContext(), role.ID)
	if err != nil {
		t.Fatalf("GetWithPermissions failed: %v", err)
	}

	if len(retrieved.Permissions) != 1 {
		t.Errorf("expected 1 permission, got %d", len(retrieved.Permissions))
	}
}

func TestRoleRepository_ListWithPermissions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role := &model.RoleModel{Name: "Admin", Slug: "admin"}
	_ = repo.Create(testContext(), role)

	perm := &model.Permission{Code: "read", Method: "GET", Path: "/api/read"}
	db.Create(perm)
	_ = repo.AssignPermissions(testContext(), role.ID, []uint64{perm.ID})

	roles, err := repo.ListWithPermissions(testContext())
	if err != nil {
		t.Fatalf("ListWithPermissions failed: %v", err)
	}

	if len(roles) < 1 {
		t.Error("expected at least 1 role")
	}

	// Find our test role and verify permissions are loaded
	for _, r := range roles {
		if r.ID == role.ID {
			if len(r.Permissions) != 1 {
				t.Errorf("expected 1 permission for role, got %d", len(r.Permissions))
			}
			return
		}
	}
	t.Error("test role not found in list")
}

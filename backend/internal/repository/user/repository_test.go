package user

import (
	"context"
	"testing"
	"time"

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

	// 迁移 schema
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestNewUserRepository 测试构造函数。
func TestNewUserRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewUserRepository(db)
	if repo == nil {
		t.Fatal("NewUserRepository returned nil")
	}

	if _, ok := repo.(*gormUserRepository); !ok {
		t.Errorf("expected *gormUserRepository, got %T", repo)
	}
}

// TestUserRepository_Create 测试创建用户
func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:        "13800138000",
		Email:        "test@example.com",
		Name:         "Test User",
		AvatarURL:    "http://example.com/avatar.jpg",
		Role:         model.RoleUser,
		Status:       model.UserStatusActive,
		PasswordHash: "hashed_password",
	}

	err := repo.Create(testContext(), user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected ID to be set after create")
	}

	// 验证创建的用户
	retrieved, err := repo.Get(testContext(), user.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Phone != user.Phone {
		t.Errorf("expected phone %s, got %s", user.Phone, retrieved.Phone)
	}
	if retrieved.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrieved.Email)
	}
	if retrieved.Name != user.Name {
		t.Errorf("expected name %s, got %s", user.Name, retrieved.Name)
	}
}

// TestUserRepository_Get 测试获取用户
func TestUserRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建测试用户
	user := &model.User{
		Phone:  "13800138001",
		Email:  "get@example.com",
		Name:   "Get Test",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	_ = repo.Create(testContext(), user)

	t.Run("Get existing user", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), user.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != user.ID {
			t.Errorf("expected ID %d, got %d", user.ID, retrieved.ID)
		}
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestUserRepository_FindByEmail 测试通过邮箱查找用户
func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:  "13800138002",
		Email:  "find@example.com",
		Name:   "Find Test",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	_ = repo.Create(testContext(), user)

	t.Run("Find existing user by email", func(t *testing.T) {
		found, err := repo.FindByEmail(testContext(), "find@example.com")
		if err != nil {
			t.Fatalf("FindByEmail failed: %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("expected ID %d, got %d", user.ID, found.ID)
		}
	})

	t.Run("Find non-existent user by email", func(t *testing.T) {
		_, err := repo.FindByEmail(testContext(), "nonexistent@example.com")
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestUserRepository_FindByPhone 测试通过手机号查找用户
func TestUserRepository_FindByPhone(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:  "13800138003",
		Email:  "phone@example.com",
		Name:   "Phone Test",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	_ = repo.Create(testContext(), user)

	t.Run("Find existing user by phone", func(t *testing.T) {
		found, err := repo.FindByPhone(testContext(), "13800138003")
		if err != nil {
			t.Fatalf("FindByPhone failed: %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("expected ID %d, got %d", user.ID, found.ID)
		}
	})

	t.Run("GetByPhone also works", func(t *testing.T) {
		found, err := repo.GetByPhone(testContext(), "13800138003")
		if err != nil {
			t.Fatalf("GetByPhone failed: %v", err)
		}
		if found.ID != user.ID {
			t.Errorf("expected ID %d, got %d", user.ID, found.ID)
		}
	})

	t.Run("Find non-existent user by phone", func(t *testing.T) {
		_, err := repo.FindByPhone(testContext(), "99999999999")
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestUserRepository_Update 测试更新用户
func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:  "13800138004",
		Email:  "update@example.com",
		Name:   "Update Test",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	_ = repo.Create(testContext(), user)

	t.Run("Update existing user", func(t *testing.T) {
		user.Name = "Updated Name"
		user.Email = "updated@example.com"
		user.Role = model.RolePlayer

		err := repo.Update(testContext(), user)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// 验证更新
		updated, _ := repo.Get(testContext(), user.ID)
		if updated.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", updated.Name)
		}
		if updated.Email != "updated@example.com" {
			t.Errorf("expected email 'updated@example.com', got %s", updated.Email)
		}
		if updated.Role != model.RolePlayer {
			t.Errorf("expected role player, got %s", updated.Role)
		}
	})

	t.Run("Update non-existent user", func(t *testing.T) {
		nonExistent := &model.User{Base: model.Base{ID: 99999}, Name: "Ghost"}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestUserRepository_Delete 测试删除用户
func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &model.User{
		Phone:  "13800138005",
		Email:  "delete@example.com",
		Name:   "Delete Test",
		Role:   model.RoleUser,
		Status: model.UserStatusActive,
	}
	_ = repo.Create(testContext(), user)

	t.Run("Delete existing user", func(t *testing.T) {
		err := repo.Delete(testContext(), user.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// 验证已删除（软删除）
		_, err = repo.Get(testContext(), user.ID)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestUserRepository_List 测试列表查询
func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建多个用户
	users := []*model.User{
		{Phone: "13800138010", Email: "user1@example.com", Name: "User 1", Role: model.RoleUser, Status: model.UserStatusActive},
		{Phone: "13800138011", Email: "user2@example.com", Name: "User 2", Role: model.RolePlayer, Status: model.UserStatusActive},
		{Phone: "13800138012", Email: "user3@example.com", Name: "User 3", Role: model.RoleUser, Status: model.UserStatusSuspended},
	}
	for _, u := range users {
		_ = repo.Create(testContext(), u)
	}

	t.Run("List all users", func(t *testing.T) {
		list, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(list) < 3 {
			t.Errorf("expected at least 3 users, got %d", len(list))
		}
	})
}

// TestUserRepository_ListPaged 测试分页查询
func TestUserRepository_ListPaged(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建多个用户
	for i := 0; i < 15; i++ {
		user := &model.User{
			Phone:  "1380013" + string(rune('0'+i)),
			Email:  "page" + string(rune('0'+i)) + "@example.com",
			Name:   "Page User " + string(rune('0'+i)),
			Role:   model.RoleUser,
			Status: model.UserStatusActive,
		}
		_ = repo.Create(testContext(), user)
	}

	t.Run("First page", func(t *testing.T) {
		users, total, err := repo.ListPaged(testContext(), 1, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(users) != 10 {
			t.Errorf("expected 10 users, got %d", len(users))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("Second page", func(t *testing.T) {
		users, total, err := repo.ListPaged(testContext(), 2, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(users) < 5 {
			t.Errorf("expected at least 5 users on page 2, got %d", len(users))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("Invalid page defaults to 1", func(t *testing.T) {
		users, _, err := repo.ListPaged(testContext(), 0, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(users) != 10 {
			t.Errorf("expected 10 users, got %d", len(users))
		}
	})
}

// TestUserRepository_ListWithFilters 测试带过滤条件的查询
func TestUserRepository_ListWithFilters(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 创建不同角色和状态的用户
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	users := []*model.User{
		{Phone: "13800138020", Email: "filter1@example.com", Name: "Alice", Role: model.RoleUser, Status: model.UserStatusActive},
		{Phone: "13800138021", Email: "filter2@example.com", Name: "Bob", Role: model.RolePlayer, Status: model.UserStatusActive},
		{Phone: "13800138022", Email: "filter3@example.com", Name: "Charlie", Role: model.RoleAdmin, Status: model.UserStatusSuspended},
		{Phone: "13800138023", Email: "alice@example.com", Name: "Alice2", Role: model.RoleUser, Status: model.UserStatusActive},
	}
	for _, u := range users {
		u.Base.CreatedAt = yesterday
		_ = repo.Create(testContext(), u)
	}

	t.Run("Filter by role", func(t *testing.T) {
		opts := repository.UserListOptions{
			Page:     1,
			PageSize: 20,
			Roles:    []model.Role{model.RolePlayer},
		}
		_, total, err := repo.ListWithFilters(testContext(), opts)
		if err != nil {
			t.Fatalf("ListWithFilters failed: %v", err)
		}
		if total < 1 {
			t.Errorf("expected at least 1 player, got %d", total)
		}
	})

	t.Run("Filter by status", func(t *testing.T) {
		opts := repository.UserListOptions{
			Page:     1,
			PageSize: 20,
			Statuses: []model.UserStatus{model.UserStatusActive},
		}
		_, total, err := repo.ListWithFilters(testContext(), opts)
		if err != nil {
			t.Fatalf("ListWithFilters failed: %v", err)
		}
		if total < 3 {
			t.Errorf("expected at least 3 active users, got %d", total)
		}
	})

	t.Run("Filter by keyword", func(t *testing.T) {
		opts := repository.UserListOptions{
			Page:     1,
			PageSize: 20,
			Keyword:  "Alice",
		}
		_, total, err := repo.ListWithFilters(testContext(), opts)
		if err != nil {
			t.Fatalf("ListWithFilters failed: %v", err)
		}
		if total < 2 {
			t.Errorf("expected at least 2 users with 'Alice', got %d", total)
		}
	})

	t.Run("Filter by date range", func(t *testing.T) {
		twoDaysAgo := now.Add(-48 * time.Hour)
		opts := repository.UserListOptions{
			Page:     1,
			PageSize: 20,
			DateFrom: &twoDaysAgo,
			DateTo:   &now,
		}
		_, total, err := repo.ListWithFilters(testContext(), opts)
		if err != nil {
			t.Fatalf("ListWithFilters failed: %v", err)
		}
		if total < 4 {
			t.Errorf("expected at least 4 users in date range, got %d", total)
		}
	})

	t.Run("Combined filters", func(t *testing.T) {
		opts := repository.UserListOptions{
			Page:     1,
			PageSize: 20,
			Roles:    []model.Role{model.RoleUser},
			Statuses: []model.UserStatus{model.UserStatusActive},
		}
		_, total, err := repo.ListWithFilters(testContext(), opts)
		if err != nil {
			t.Fatalf("ListWithFilters failed: %v", err)
		}
		if total < 2 {
			t.Errorf("expected at least 2 active users with role user, got %d", total)
		}
	})
}

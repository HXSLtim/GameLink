package role

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

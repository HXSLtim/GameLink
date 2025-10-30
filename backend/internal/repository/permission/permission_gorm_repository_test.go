package permission

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewPermissionRepository 测试构造函数。
func TestNewPermissionRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewPermissionRepository(db)
	if repo == nil {
		t.Fatal("NewPermissionRepository returned nil")
	}

	if _, ok := repo.(*permissionRepository); !ok {
		t.Errorf("expected *permissionRepository, got %T", repo)
	}
}

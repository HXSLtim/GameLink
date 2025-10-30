package user

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

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

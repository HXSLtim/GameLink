package order

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewOrderRepository 测试构造函数。
func TestNewOrderRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewOrderRepository(db)
	if repo == nil {
		t.Fatal("NewOrderRepository returned nil")
	}

	if _, ok := repo.(*gormOrderRepository); !ok {
		t.Errorf("expected *gormOrderRepository, got %T", repo)
	}
}

package player

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewPlayerRepository 测试构造函数。
func TestNewPlayerRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewPlayerRepository(db)
	if repo == nil {
		t.Fatal("NewPlayerRepository returned nil")
	}

	if _, ok := repo.(*gormPlayerRepository); !ok {
		t.Errorf("expected *gormPlayerRepository, got %T", repo)
	}
}

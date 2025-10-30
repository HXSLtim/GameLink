package playertag

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewPlayerTagRepository 测试构造函数。
func TestNewPlayerTagRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewPlayerTagRepository(db)
	if repo == nil {
		t.Fatal("NewPlayerTagRepository returned nil")
	}

	if _, ok := repo.(*gormPlayerTagRepository); !ok {
		t.Errorf("expected *gormPlayerTagRepository, got %T", repo)
	}
}

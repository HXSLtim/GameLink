package stats

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewStatsRepository 测试构造函数。
func TestNewStatsRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewStatsRepository(db)
	if repo == nil {
		t.Fatal("NewStatsRepository returned nil")
	}

	if _, ok := repo.(*gormStatsRepository); !ok {
		t.Errorf("expected *gormStatsRepository, got %T", repo)
	}
}

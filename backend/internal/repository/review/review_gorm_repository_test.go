package review

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewReviewRepository 测试构造函数。
func TestNewReviewRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewReviewRepository(db)
	if repo == nil {
		t.Fatal("NewReviewRepository returned nil")
	}

	if _, ok := repo.(*gormReviewRepository); !ok {
		t.Errorf("expected *gormReviewRepository, got %T", repo)
	}
}

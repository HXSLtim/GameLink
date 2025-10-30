package payment

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TestNewPaymentRepository 测试构造函数。
func TestNewPaymentRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewPaymentRepository(db)
	if repo == nil {
		t.Fatal("NewPaymentRepository returned nil")
	}

	if _, ok := repo.(*gormPaymentRepository); !ok {
		t.Errorf("expected *gormPaymentRepository, got %T", repo)
	}
}

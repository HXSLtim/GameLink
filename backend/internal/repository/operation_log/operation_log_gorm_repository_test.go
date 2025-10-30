package operationlog

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
)

func testContext() context.Context {
	return context.Background()
}

// TestNewOperationLogRepository 测试构造函数。
func TestNewOperationLogRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewOperationLogRepository(db)
	if repo == nil {
		t.Fatal("NewOperationLogRepository returned nil")
	}

	if _, ok := repo.(*gormOperationLogRepository); !ok {
		t.Errorf("expected *gormOperationLogRepository, got %T", repo)
	}
}

// TestOperationLogRepositoryWithDB 测试基本操作。
func TestOperationLogRepositoryWithDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.OperationLog{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	repo := NewOperationLogRepository(db)

	t.Run("Append creates log", func(t *testing.T) {
		log := &model.OperationLog{
			EntityType: "game",
			EntityID:   1,
			Action:     "create",
		}

		if err := repo.Append(testContext(), log); err != nil {
			t.Fatalf("Append failed: %v", err)
		}

		if log.ID == 0 {
			t.Error("expected ID to be set after append")
		}
	})
}

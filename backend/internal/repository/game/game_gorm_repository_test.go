package game

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

// TestNewGameRepository 测试构造函数。
func TestNewGameRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewGameRepository(db)
	if repo == nil {
		t.Fatal("NewGameRepository returned nil")
	}

	// 确保返回的是正确的类型
	if _, ok := repo.(*gormGameRepository); !ok {
		t.Errorf("expected *gormGameRepository, got %T", repo)
	}
}

// TestGameRepositoryWithDB 测试基本的CRUD操作（需要数据库迁移）。
func TestGameRepositoryWithDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// 迁移schema
	if err := db.AutoMigrate(&model.Game{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	repo := NewGameRepository(db)

	t.Run("List returns empty slice when no games", func(t *testing.T) {
		games, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(games) != 0 {
			t.Errorf("expected 0 games, got %d", len(games))
		}
	})

	t.Run("Create and Get game", func(t *testing.T) {
		game := &model.Game{
			Key:         "lol",
			Name:        "League of Legends",
			Category:    "MOBA",
			IconURL:     "http://example.com/icon.png",
			Description: "A MOBA game",
		}

		if err := repo.Create(testContext(), game); err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		if game.ID == 0 {
			t.Error("expected ID to be set after create")
		}

		retrieved, err := repo.Get(testContext(), game.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.Key != game.Key {
			t.Errorf("expected key %s, got %s", game.Key, retrieved.Key)
		}
		if retrieved.Name != game.Name {
			t.Errorf("expected name %s, got %s", game.Name, retrieved.Name)
		}
	})
}

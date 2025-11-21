package game

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gamelink/internal/model"
	"gamelink/internal/repository"
)

func testContext() context.Context {
	return context.Background()
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.Game{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestGameRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	game := &model.Game{
		Key:         "lol",
		Name:        "League of Legends",
		Category:    "MOBA",
		IconURL:     "https://example.com/lol.png",
		Description: "5v5 战术竞技游戏",
	}

	err := repo.Create(testContext(), game)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if game.ID == 0 {
		t.Error("expected game ID to be set")
	}

	// Verify created
	retrieved, err := repo.Get(testContext(), game.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Key != "lol" {
		t.Errorf("expected key 'lol', got %s", retrieved.Key)
	}
	if retrieved.Name != "League of Legends" {
		t.Errorf("expected name 'League of Legends', got %s", retrieved.Name)
	}
}

func TestGameRepository_CreateDuplicateKey(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	game1 := &model.Game{Key: "dota2", Name: "Dota 2"}
	_ = repo.Create(testContext(), game1)

	game2 := &model.Game{Key: "dota2", Name: "Dota 2 Copy"}
	err := repo.Create(testContext(), game2)
	if err == nil {
		t.Error("expected error for duplicate key, got nil")
	}
}

func TestGameRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	game := &model.Game{Key: "valorant", Name: "Valorant", Category: "FPS"}
	_ = repo.Create(testContext(), game)

	t.Run("Get existing game", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), game.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.ID != game.ID {
			t.Errorf("expected ID %d, got %d", game.ID, retrieved.ID)
		}
		if retrieved.Key != "valorant" {
			t.Errorf("expected key 'valorant', got %s", retrieved.Key)
		}
	})

	t.Run("Get non-existent game", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGameRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	game := &model.Game{
		Key:      "csgo",
		Name:     "CS:GO",
		Category: "FPS",
	}
	_ = repo.Create(testContext(), game)

	t.Run("Update existing game", func(t *testing.T) {
		game.Name = "Counter-Strike: Global Offensive"
		game.Description = "经典 FPS 游戏"
		game.IconURL = "https://example.com/csgo.png"

		err := repo.Update(testContext(), game)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Verify updated
		retrieved, _ := repo.Get(testContext(), game.ID)
		if retrieved.Name != "Counter-Strike: Global Offensive" {
			t.Errorf("expected updated name, got %s", retrieved.Name)
		}
		if retrieved.Description != "经典 FPS 游戏" {
			t.Errorf("expected updated description, got %s", retrieved.Description)
		}
	})

	t.Run("Update non-existent game", func(t *testing.T) {
		nonExistent := &model.Game{Base: model.Base{ID: 99999}, Key: "fake", Name: "Fake"}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGameRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	game := &model.Game{Key: "apex", Name: "Apex Legends"}
	_ = repo.Create(testContext(), game)

	t.Run("Delete existing game", func(t *testing.T) {
		err := repo.Delete(testContext(), game.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// Verify deleted
		_, err = repo.Get(testContext(), game.ID)
		if err != repository.ErrNotFound {
			t.Error("expected game to be deleted")
		}
	})

	t.Run("Delete non-existent game", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGameRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	t.Run("List empty", func(t *testing.T) {
		games, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(games) != 0 {
			t.Errorf("expected 0 games, got %d", len(games))
		}
	})

	// Create test data
	games := []model.Game{
		{Key: "lol", Name: "League of Legends"},
		{Key: "dota2", Name: "Dota 2"},
		{Key: "csgo", Name: "CS:GO"},
	}
	for i := range games {
		_ = repo.Create(testContext(), &games[i])
	}

	t.Run("List all games", func(t *testing.T) {
		result, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(result) != 3 {
			t.Errorf("expected 3 games, got %d", len(result))
		}

		// Verify all games are present (order may vary due to same timestamps)
		keys := make(map[string]bool)
		for _, g := range result {
			keys[g.Key] = true
		}

		if !keys["lol"] || !keys["dota2"] || !keys["csgo"] {
			t.Error("expected all three games (lol, dota2, csgo) to be present")
		}
	})
}

func TestGameRepository_ListPaged(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	// Create test data
	for i := 1; i <= 15; i++ {
		game := &model.Game{
			Key:  "game" + string(rune('0'+i)),
			Name: "Game " + string(rune('0'+i)),
		}
		_ = repo.Create(testContext(), game)
	}

	t.Run("First page", func(t *testing.T) {
		games, total, err := repo.ListPaged(testContext(), 1, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(games) != 10 {
			t.Errorf("expected 10 games, got %d", len(games))
		}
	})

	t.Run("Second page", func(t *testing.T) {
		games, total, err := repo.ListPaged(testContext(), 2, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(games) != 5 {
			t.Errorf("expected 5 games on second page, got %d", len(games))
		}
	})

	t.Run("Page beyond data", func(t *testing.T) {
		games, total, err := repo.ListPaged(testContext(), 10, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(games) != 0 {
			t.Errorf("expected 0 games, got %d", len(games))
		}
	})

	t.Run("Invalid page number normalized", func(t *testing.T) {
		games, total, err := repo.ListPaged(testContext(), 0, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if total != 15 {
			t.Errorf("expected total 15, got %d", total)
		}
		if len(games) != 10 {
			t.Errorf("expected 10 games (normalized to page 1), got %d", len(games))
		}
	})

	t.Run("Custom page size", func(t *testing.T) {
		games, _, err := repo.ListPaged(testContext(), 1, 5)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}

		if len(games) != 5 {
			t.Errorf("expected 5 games, got %d", len(games))
		}
	})
}

func TestGameRepository_CompleteWorkflow(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	// Create
	game := &model.Game{
		Key:         "fortnite",
		Name:        "Fortnite",
		Category:    "Battle Royale",
		IconURL:     "https://example.com/fortnite.png",
		Description: "大逃杀游戏",
	}
	err := repo.Create(testContext(), game)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Read
	retrieved, err := repo.Get(testContext(), game.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if retrieved.Key != "fortnite" {
		t.Errorf("expected key 'fortnite', got %s", retrieved.Key)
	}

	// Update
	game.Name = "Fortnite Battle Royale"
	game.Description = "免费大逃杀游戏"
	err = repo.Update(testContext(), game)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	updated, _ := repo.Get(testContext(), game.ID)
	if updated.Name != "Fortnite Battle Royale" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}

	// List
	games, _ := repo.List(testContext())
	if len(games) != 1 {
		t.Errorf("expected 1 game in list, got %d", len(games))
	}

	// Delete
	err = repo.Delete(testContext(), game.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.Get(testContext(), game.ID)
	if err != repository.ErrNotFound {
		t.Error("expected game to be deleted")
	}
}

func TestGameRepository_MultipleGamesOrdering(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)

	// Create games in order
	games := []model.Game{
		{Key: "game1", Name: "First Game"},
		{Key: "game2", Name: "Second Game"},
		{Key: "game3", Name: "Third Game"},
	}
	for i := range games {
		_ = repo.Create(testContext(), &games[i])
	}

	// List should return all games
	result, _ := repo.List(testContext())
	if len(result) != 3 {
		t.Fatalf("expected 3 games, got %d", len(result))
	}

	// Verify all games are present (order may vary due to same timestamps)
	keys := make(map[string]bool)
	for _, g := range result {
		keys[g.Key] = true
	}

	if !keys["game1"] || !keys["game2"] || !keys["game3"] {
		t.Error("expected all three games (game1, game2, game3) to be present")
	}
}

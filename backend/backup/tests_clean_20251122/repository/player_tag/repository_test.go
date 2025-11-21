package playertag

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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(&model.PlayerSkillTag{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

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

func TestPlayerTagRepository_GetTags_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerTagRepository(db)

	tags, err := repo.GetTags(testContext(), 1)
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %d tags", len(tags))
	}
}

func TestPlayerTagRepository_ReplaceTags(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerTagRepository(db)

	t.Run("Add new tags", func(t *testing.T) {
		tags := []string{"friendly", "skilled", "patient"}
		err := repo.ReplaceTags(testContext(), 1, tags)
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, err := repo.GetTags(testContext(), 1)
		if err != nil {
			t.Fatalf("GetTags failed: %v", err)
		}

		if len(retrieved) != 3 {
			t.Errorf("expected 3 tags, got %d", len(retrieved))
		}
	})

	t.Run("Replace existing tags", func(t *testing.T) {
		// First set
		_ = repo.ReplaceTags(testContext(), 2, []string{"tag1", "tag2"})

		// Replace
		newTags := []string{"tag3", "tag4", "tag5"}
		err := repo.ReplaceTags(testContext(), 2, newTags)
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, _ := repo.GetTags(testContext(), 2)
		if len(retrieved) != 3 {
			t.Errorf("expected 3 tags after replacement, got %d", len(retrieved))
		}

		// Verify old tags are gone
		for _, tag := range retrieved {
			if tag == "tag1" || tag == "tag2" {
				t.Errorf("old tag %s still exists", tag)
			}
		}
	})

	t.Run("Deduplicate tags", func(t *testing.T) {
		tags := []string{"friendly", "friendly", "skilled", "Friendly"}
		err := repo.ReplaceTags(testContext(), 3, tags)
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, _ := repo.GetTags(testContext(), 3)
		// Should deduplicate to 2 unique tags (friendly, skilled)
		if len(retrieved) != 2 {
			t.Errorf("expected 2 unique tags after deduplication, got %d", len(retrieved))
		}
	})

	t.Run("Normalize tags (trim and lowercase)", func(t *testing.T) {
		tags := []string{"  Friendly  ", "SKILLED", "Patient"}
		err := repo.ReplaceTags(testContext(), 4, tags)
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, _ := repo.GetTags(testContext(), 4)
		if len(retrieved) != 3 {
			t.Errorf("expected 3 tags, got %d", len(retrieved))
		}

		// All should be lowercase and trimmed
		for _, tag := range retrieved {
			if tag == "  Friendly  " || tag == "SKILLED" {
				t.Errorf("tag not normalized: %s", tag)
			}
		}
	})

	t.Run("Handle empty tags", func(t *testing.T) {
		tags := []string{"friendly", "", "  ", "skilled"}
		err := repo.ReplaceTags(testContext(), 5, tags)
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, _ := repo.GetTags(testContext(), 5)
		// Should filter out empty/whitespace tags
		if len(retrieved) != 2 {
			t.Errorf("expected 2 tags after filtering empty, got %d", len(retrieved))
		}
	})

	t.Run("Replace with empty list", func(t *testing.T) {
		// First add some tags
		_ = repo.ReplaceTags(testContext(), 6, []string{"tag1", "tag2"})

		// Replace with empty
		err := repo.ReplaceTags(testContext(), 6, []string{})
		if err != nil {
			t.Fatalf("ReplaceTags failed: %v", err)
		}

		retrieved, _ := repo.GetTags(testContext(), 6)
		if len(retrieved) != 0 {
			t.Errorf("expected 0 tags after clearing, got %d", len(retrieved))
		}
	})
}

func TestPlayerTagRepository_GetTags(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerTagRepository(db)

	// Add tags for player 1
	_ = repo.ReplaceTags(testContext(), 1, []string{"friendly", "skilled"})

	tags, err := repo.GetTags(testContext(), 1)
	if err != nil {
		t.Fatalf("GetTags failed: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	// Verify tags are returned
	for _, tag := range tags {
		if tag == "" {
			t.Error("expected tag to be set")
		}
	}
}

func TestPlayerTagRepository_MultiplePlayer(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerTagRepository(db)

	// Add tags for multiple players
	_ = repo.ReplaceTags(testContext(), 1, []string{"friendly"})
	_ = repo.ReplaceTags(testContext(), 2, []string{"skilled", "patient"})

	// Verify player 1 tags
	tags1, _ := repo.GetTags(testContext(), 1)
	if len(tags1) != 1 {
		t.Errorf("expected 1 tag for player 1, got %d", len(tags1))
	}

	// Verify player 2 tags
	tags2, _ := repo.GetTags(testContext(), 2)
	if len(tags2) != 2 {
		t.Errorf("expected 2 tags for player 2, got %d", len(tags2))
	}

	// Verify tags are sorted
	if len(tags2) == 2 {
		if tags2[0] > tags2[1] {
			t.Error("tags should be sorted alphabetically")
		}
	}
}

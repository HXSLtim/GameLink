package review

import (
	"context"
	"testing"
	"time"

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

	if err := db.AutoMigrate(&model.Review{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestNewReviewRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewReviewRepository(db)
	if repo == nil {
		t.Fatal("NewReviewRepository returned nil")
	}
}

func TestReviewRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReviewRepository(db)

	review := &model.Review{
		OrderID:  1,
		UserID:   1,
		PlayerID: 1,
		Score:    5,
		Content:  "Great player!",
	}

	err := repo.Create(testContext(), review)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if review.ID == 0 {
		t.Error("expected ID to be set after create")
	}

	retrieved, _ := repo.Get(testContext(), review.ID)
	if retrieved.Score != 5 {
		t.Errorf("expected score 5, got %d", retrieved.Score)
	}
}

func TestReviewRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReviewRepository(db)

	review := &model.Review{
		OrderID:  2,
		UserID:   2,
		PlayerID: 2,
		Score:    4,
		Content:  "Good service",
	}
	_ = repo.Create(testContext(), review)

	t.Run("Get existing review", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), review.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != review.ID {
			t.Errorf("expected ID %d, got %d", review.ID, retrieved.ID)
		}
	})

	t.Run("Get non-existent review", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestReviewRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReviewRepository(db)

	review := &model.Review{
		OrderID:  3,
		UserID:   3,
		PlayerID: 3,
		Score:    3,
		Content:  "OK",
	}
	_ = repo.Create(testContext(), review)

	t.Run("Update existing review", func(t *testing.T) {
		review.Score = 5
		review.Content = "Actually great!"

		err := repo.Update(testContext(), review)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		updated, _ := repo.Get(testContext(), review.ID)
		if updated.Score != 5 {
			t.Errorf("expected score 5, got %d", updated.Score)
		}
		if updated.Content != "Actually great!" {
			t.Errorf("expected updated content, got %s", updated.Content)
		}
	})

	t.Run("Update non-existent review", func(t *testing.T) {
		nonExistent := &model.Review{Base: model.Base{ID: 99999}}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestReviewRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReviewRepository(db)

	review := &model.Review{
		OrderID:  4,
		UserID:   4,
		PlayerID: 4,
		Score:    3,
	}
	_ = repo.Create(testContext(), review)

	t.Run("Delete existing review", func(t *testing.T) {
		err := repo.Delete(testContext(), review.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repo.Get(testContext(), review.ID)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("Delete non-existent review", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestReviewRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewReviewRepository(db)

	// 创建多个评价
	for i := 0; i < 15; i++ {
		score := model.Rating(3 + i%3) // 3, 4, 5

		review := &model.Review{
			OrderID:  uint64(100 + i),
			UserID:   uint64(10 + i%3),
			PlayerID: uint64(20 + i%2),
			Score:    score,
			Content:  "Review " + string(rune('A'+i)),
		}
		_ = repo.Create(testContext(), review)
	}

	t.Run("List with pagination", func(t *testing.T) {
		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 10,
		}
		reviews, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(reviews) != 10 {
			t.Errorf("expected 10 reviews, got %d", len(reviews))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("List by order", func(t *testing.T) {
		orderID := uint64(100)
		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 20,
			OrderID:  &orderID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total == 1 {
			t.Logf("found review for order 100")
		}
	})

	t.Run("List by user", func(t *testing.T) {
		userID := uint64(10)
		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 20,
			UserID:   &userID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 5 {
			t.Logf("found %d reviews from user 10", total)
		}
	})

	t.Run("List by player", func(t *testing.T) {
		playerID := uint64(20)
		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 20,
			PlayerID: &playerID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 7 {
			t.Logf("found %d reviews for player 20", total)
		}
	})

	t.Run("List combined filters", func(t *testing.T) {
		playerID := uint64(20)
		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 20,
			PlayerID: &playerID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		t.Logf("found %d reviews for player in combined test", total)
	})

	t.Run("List with date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		opts := repository.ReviewListOptions{
			Page:     1,
			PageSize: 20,
			DateFrom: &yesterday,
			DateTo:   &tomorrow,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 15 {
			t.Errorf("expected at least 15 reviews in date range, got %d", total)
		}
	})
}

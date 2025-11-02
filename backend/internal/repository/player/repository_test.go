package player

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

	// 迁移 schema
	if err := db.AutoMigrate(&model.Player{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

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

// TestPlayerRepository_Create 测试创建陪玩师
func TestPlayerRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	player := &model.Player{
		UserID:             1,
		Nickname:           "TestPlayer",
		Bio:                "Test bio",
		Rank:               "Diamond",
		RatingAverage:      4.5,
		RatingCount:        10,
		HourlyRateCents:    10000,
		MainGameID:         1,
		VerificationStatus: model.VerificationPending,
	}

	err := repo.Create(testContext(), player)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if player.ID == 0 {
		t.Error("expected ID to be set after create")
	}

	// 验证创建的陪玩师
	retrieved, err := repo.Get(testContext(), player.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Nickname != player.Nickname {
		t.Errorf("expected nickname %s, got %s", player.Nickname, retrieved.Nickname)
	}
	if retrieved.Bio != player.Bio {
		t.Errorf("expected bio %s, got %s", player.Bio, retrieved.Bio)
	}
}

// TestPlayerRepository_Get 测试获取陪玩师
func TestPlayerRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	player := &model.Player{
		UserID:             2,
		Nickname:           "GetTest",
		VerificationStatus: model.VerificationVerified,
	}
	_ = repo.Create(testContext(), player)

	t.Run("Get existing player", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), player.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != player.ID {
			t.Errorf("expected ID %d, got %d", player.ID, retrieved.ID)
		}
	})

	t.Run("Get non-existent player", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestPlayerRepository_Update 测试更新陪玩师
func TestPlayerRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	player := &model.Player{
		UserID:             3,
		Nickname:           "UpdateTest",
		Bio:                "Original bio",
		HourlyRateCents:    10000,
		VerificationStatus: model.VerificationPending,
	}
	_ = repo.Create(testContext(), player)

	t.Run("Update existing player", func(t *testing.T) {
		player.Nickname = "Updated Nickname"
		player.Bio = "Updated bio"
		player.HourlyRateCents = 15000
		player.VerificationStatus = model.VerificationVerified

		err := repo.Update(testContext(), player)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// 验证更新
		updated, _ := repo.Get(testContext(), player.ID)
		if updated.Nickname != "Updated Nickname" {
			t.Errorf("expected nickname 'Updated Nickname', got %s", updated.Nickname)
		}
		if updated.Bio != "Updated bio" {
			t.Errorf("expected bio 'Updated bio', got %s", updated.Bio)
		}
		if updated.HourlyRateCents != 15000 {
			t.Errorf("expected hourly rate 15000, got %d", updated.HourlyRateCents)
		}
		if updated.VerificationStatus != model.VerificationVerified {
			t.Errorf("expected verified status, got %s", updated.VerificationStatus)
		}
	})

	t.Run("Update non-existent player", func(t *testing.T) {
		nonExistent := &model.Player{Base: model.Base{ID: 99999}, Nickname: "Ghost"}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestPlayerRepository_Delete 测试删除陪玩师
func TestPlayerRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	player := &model.Player{
		UserID:             4,
		Nickname:           "DeleteTest",
		VerificationStatus: model.VerificationVerified,
	}
	_ = repo.Create(testContext(), player)

	t.Run("Delete existing player", func(t *testing.T) {
		err := repo.Delete(testContext(), player.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// 验证已删除（软删除）
		_, err = repo.Get(testContext(), player.ID)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("Delete non-existent player", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

// TestPlayerRepository_List 测试列表查询
func TestPlayerRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	// 创建多个陪玩师
	players := []*model.Player{
		{UserID: 10, Nickname: "Player1", VerificationStatus: model.VerificationVerified},
		{UserID: 11, Nickname: "Player2", VerificationStatus: model.VerificationVerified},
		{UserID: 12, Nickname: "Player3", VerificationStatus: model.VerificationPending},
	}
	for _, p := range players {
		_ = repo.Create(testContext(), p)
	}

	t.Run("List all players", func(t *testing.T) {
		list, err := repo.List(testContext())
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(list) < 3 {
			t.Errorf("expected at least 3 players, got %d", len(list))
		}
	})
}

// TestPlayerRepository_ListPaged 测试分页查询
func TestPlayerRepository_ListPaged(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	// 创建多个陪玩师
	for i := 0; i < 15; i++ {
		player := &model.Player{
			UserID:             uint64(20 + i),
			Nickname:           "PagedPlayer" + string(rune('0'+i)),
			VerificationStatus: model.VerificationVerified,
		}
		_ = repo.Create(testContext(), player)
	}

	t.Run("First page", func(t *testing.T) {
		players, total, err := repo.ListPaged(testContext(), 1, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(players) != 10 {
			t.Errorf("expected 10 players, got %d", len(players))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("Second page", func(t *testing.T) {
		players, total, err := repo.ListPaged(testContext(), 2, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(players) < 5 {
			t.Errorf("expected at least 5 players on page 2, got %d", len(players))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("Invalid page defaults to 1", func(t *testing.T) {
		players, _, err := repo.ListPaged(testContext(), 0, 10)
		if err != nil {
			t.Fatalf("ListPaged failed: %v", err)
		}
		if len(players) != 10 {
			t.Errorf("expected 10 players, got %d", len(players))
		}
	})
}

// TestPlayerRepository_UpdateRating 测试评分更新
func TestPlayerRepository_UpdateRating(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPlayerRepository(db)

	player := &model.Player{
		UserID:             5,
		Nickname:           "RatingTest",
		RatingAverage:      4.0,
		RatingCount:        5,
		VerificationStatus: model.VerificationVerified,
	}
	_ = repo.Create(testContext(), player)

	t.Run("Update rating", func(t *testing.T) {
		// 更新评分
		player.RatingAverage = 4.5
		player.RatingCount = 10

		err := repo.Update(testContext(), player)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// 验证更新
		updated, _ := repo.Get(testContext(), player.ID)
		if updated.RatingAverage != 4.5 {
			t.Errorf("expected rating average 4.5, got %f", updated.RatingAverage)
		}
		if updated.RatingCount != 10 {
			t.Errorf("expected rating count 10, got %d", updated.RatingCount)
		}
	})
}

package order

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

	if err := db.AutoMigrate(&model.Order{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestNewOrderRepository 测试构造函数。
func TestNewOrderRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	repo := NewOrderRepository(db)
	if repo == nil {
		t.Fatal("NewOrderRepository returned nil")
	}

	if _, ok := repo.(*gormOrderRepository); !ok {
		t.Errorf("expected *gormOrderRepository, got %T", repo)
	}
}

func TestOrderRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	now := time.Now()
	order := &model.Order{
		UserID:         1,
		PlayerID:       1,
		GameID:         1,
		Title:          "Test Order",
		Description:    "Test description",
		Status:         model.OrderStatusPending,
		PriceCents:     10000,
		Currency:       model.CurrencyCNY,
		ScheduledStart: &now,
	}

	err := repo.Create(testContext(), order)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if order.ID == 0 {
		t.Error("expected ID to be set after create")
	}

	retrieved, err := repo.Get(testContext(), order.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.Title != order.Title {
		t.Errorf("expected title %s, got %s", order.Title, retrieved.Title)
	}
	if retrieved.Status != order.Status {
		t.Errorf("expected status %s, got %s", order.Status, retrieved.Status)
	}
}

func TestOrderRepository_Get(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		UserID:     2,
		GameID:     1,
		Title:      "Get Test",
		Status:     model.OrderStatusPending,
		PriceCents: 5000,
	}
	_ = repo.Create(testContext(), order)

	t.Run("Get existing order", func(t *testing.T) {
		retrieved, err := repo.Get(testContext(), order.ID)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.ID != order.ID {
			t.Errorf("expected ID %d, got %d", order.ID, retrieved.ID)
		}
	})

	t.Run("Get non-existent order", func(t *testing.T) {
		_, err := repo.Get(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestOrderRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		UserID:     3,
		GameID:     1,
		Title:      "Update Test",
		Status:     model.OrderStatusPending,
		PriceCents: 8000,
	}
	_ = repo.Create(testContext(), order)

	t.Run("Update existing order", func(t *testing.T) {
		order.Status = model.OrderStatusConfirmed
		order.PriceCents = 12000
		order.CancelReason = "Test reason"

		err := repo.Update(testContext(), order)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		updated, _ := repo.Get(testContext(), order.ID)
		if updated.Status != model.OrderStatusConfirmed {
			t.Errorf("expected status confirmed, got %s", updated.Status)
		}
		if updated.PriceCents != 12000 {
			t.Errorf("expected price 12000, got %d", updated.PriceCents)
		}
	})

	t.Run("Update non-existent order", func(t *testing.T) {
		nonExistent := &model.Order{Base: model.Base{ID: 99999}, Title: "Ghost"}
		err := repo.Update(testContext(), nonExistent)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestOrderRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	order := &model.Order{
		UserID:     4,
		GameID:     1,
		Title:      "Delete Test",
		Status:     model.OrderStatusPending,
		PriceCents: 6000,
	}
	_ = repo.Create(testContext(), order)

	t.Run("Delete existing order", func(t *testing.T) {
		err := repo.Delete(testContext(), order.ID)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err = repo.Get(testContext(), order.ID)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("Delete non-existent order", func(t *testing.T) {
		err := repo.Delete(testContext(), 99999)
		if err != repository.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestOrderRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewOrderRepository(db)

	// 创建多个订单
	for i := 0; i < 15; i++ {
		status := model.OrderStatusPending
		if i%3 == 0 {
			status = model.OrderStatusConfirmed
		} else if i%3 == 1 {
			status = model.OrderStatusCompleted
		}

		order := &model.Order{
			UserID:     uint64(10 + i%3),
			PlayerID:   uint64(20 + i%2),
			GameID:     uint64(1 + i%2),
			Title:      "Order " + string(rune('A'+i)),
			Status:     status,
			PriceCents: int64(1000 * (i + 1)),
		}
		_ = repo.Create(testContext(), order)
	}

	t.Run("List with pagination", func(t *testing.T) {
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 10,
		}
		orders, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(orders) != 10 {
			t.Errorf("expected 10 orders, got %d", len(orders))
		}
		if total < 15 {
			t.Errorf("expected total >= 15, got %d", total)
		}
	})

	t.Run("List with status filter", func(t *testing.T) {
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			Statuses: []model.OrderStatus{model.OrderStatusPending},
		}
		orders, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 5 {
			t.Errorf("expected at least 5 pending orders, got %d", total)
		}
		for _, o := range orders {
			if o.Status != model.OrderStatusPending {
				t.Errorf("expected status pending, got %s", o.Status)
			}
		}
	})

	t.Run("List with user filter", func(t *testing.T) {
		userID := uint64(10)
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			UserID:   &userID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 5 {
			t.Errorf("expected at least 5 orders for user 10, got %d", total)
		}
	})

	t.Run("List with player filter", func(t *testing.T) {
		playerID := uint64(20)
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			PlayerID: &playerID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 1 {
			t.Logf("found %d orders for player 20", total)
		}
	})

	t.Run("List with game filter", func(t *testing.T) {
		gameID := uint64(1)
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			GameID:   &gameID,
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total >= 1 {
			t.Logf("found %d orders for game 1", total)
		}
	})

	t.Run("List with keyword filter", func(t *testing.T) {
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			Keyword:  "Order",
		}
		orders, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if total < 15 {
			t.Errorf("expected at least 15 orders with 'Order' keyword, got %d", total)
		}
		if len(orders) > 0 && orders[0].Title == "" {
			t.Error("expected title to be populated")
		}
	})

	t.Run("List with date range", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		tomorrow := now.Add(24 * time.Hour)

		opts := repository.OrderListOptions{
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
			t.Errorf("expected at least 15 orders in date range, got %d", total)
		}
	})

	t.Run("List with combined filters", func(t *testing.T) {
		userID := uint64(10)
		opts := repository.OrderListOptions{
			Page:     1,
			PageSize: 20,
			UserID:   &userID,
			Statuses: []model.OrderStatus{model.OrderStatusPending, model.OrderStatusConfirmed},
		}
		_, total, err := repo.List(testContext(), opts)
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		t.Logf("found %d orders with combined filters", total)
	})
}
